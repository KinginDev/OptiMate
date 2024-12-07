package optimizer

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"log"
	"optimizer-service/cmd/internal/app/interfaces"
	"optimizer-service/cmd/internal/models"
	"optimizer-service/cmd/internal/storage"
	"optimizer-service/cmd/internal/utils"
	"path/filepath"
	"sync"

	"github.com/disintegration/imaging"
)

type IOptimizer interface {
	Optimize(filePath string) (io.ReadCloser, error)
	SupportedFormats() []string
}

// Optimizer is a struct that defines the optimizer.
type Optimizer struct {
	Storage storage.Storage
	Utils   utils.IUtils
	Repo    interfaces.IFileRepository
}

// OptimizerParams is a struct that defines the parameters for the optimizer.
type OptimizerParams struct {
	Level *string
	Path  *string
	Name  *string
	Size  *int64
}

// NewOptimizer is a function that creates a new optimizer.
// constructor
func InitOptimizer(storage storage.Storage, repo interfaces.IFileRepository, u utils.IUtils) *Optimizer {
	return &Optimizer{
		Storage: storage,
		Utils:   u,
		Repo:    repo,
	}
}

func (o *Optimizer) Optimize(filePath string) (io.ReadCloser, error) {
	// Retrieve the file
	fileReader, err := o.Storage.Retrieve(filePath)
	if err != nil {
		log.Printf("Error retrieving file %s: %v", filePath, err)
		return nil, err
	}
	// close the opened file
	defer fileReader.Close()

	// prepare buffer to store the optimized image as a reader
	var optimizedBuffer bytes.Buffer
	var optimizedImage image.Image

	// default optimization level
	level := "high"
	oParams := OptimizerParams{
		Level: &level,
	}

	cropParams := &CropParams{
		X:      10,
		Y:      12,
		Width:  300,
		Height: 500,
	}

	// Create a buffer to store the unoptimized image as it will be used multiple times
	fileBytes, err := io.ReadAll(fileReader)
	if err != nil {
		log.Printf("Error reading file:- %v", err)
		return nil, err
	}

	// Create a new io.Reader from fileBytes to be used for optimization
	fileReader = io.NopCloser(bytes.NewReader(fileBytes))

	fileType, err := o.Utils.CheckFileType(fileBytes)
	if err != nil {
		log.Printf("Error checking file type:- %v", err)
		return nil, err
	}

	// Create a wait group to wait for the optimization to complete
	wg := sync.WaitGroup{}

	// Add 1 wait group to wait for the go routine to complete
	wg.Add(1)

	switch fileType {
	case "jpeg":
		go func() {
			// Mark the wait group as done when the go routine completes
			defer wg.Done()

			optimizedImage, err = o.OptimizeJPEG(fileReader, oParams.Level, cropParams)
			if err != nil {
				log.Printf("Error optimizing jpeg file:- %v", err)
				return
			}

			// Encode the optimized image and wrtite it to the buffer
			err = jpeg.Encode(&optimizedBuffer, optimizedImage,
				&jpeg.Options{Quality: o.mapJPEGQuality(*oParams.Level)})

			// get the optimizer size
			intOptimizedSize := int64(optimizedBuffer.Len())

			// update optimizer params
			oParams.Path = &filePath
			oParams.Name = &filePath
			oParams.Size = &intOptimizedSize
			oParams.Level = &level

			if err != nil {
				log.Printf("Error encoding optimized image:- %v", err)
				return
			}
		}()

	case "png":
		go func() {
			// Mark the wait group as done when the go routine completes
			defer wg.Done()
			optimizedImage, err = o.OptimizePNG(fileReader, oParams.Level, cropParams)
			if err != nil {
				log.Printf("Error optimizing png file:- %v", err)
				return
			}
			// encode the optimized image and write it to the buffer
			err = imaging.Encode(&optimizedBuffer, optimizedImage, imaging.PNG)

			// get the optimizer size
			intOptimizedSize := int64(optimizedBuffer.Len())

			// update optimizer params
			oParams.Path = &filePath
			oParams.Name = &filePath
			oParams.Size = &intOptimizedSize
			oParams.Level = &level

			if err != nil {
				log.Printf("Error encoding optimized image:- %v", err)
				return
			}
		}()
	case "webp":

	default:
		log.Printf("Unsupported file type:- %s", fileType)
	}

	// Wait for the optimization to complete
	wg.Wait()

	// generate the optimized file name
	fileName := filepath.Base(filePath)
	ext := filepath.Ext(fileName)
	optimizedFileName := fileName[:len(fileName)-len(ext)] + "_optimized" + ext

	// Save the optimized image to the storage
	err = o.Storage.Save(optimizedFileName, &optimizedBuffer)
	if err != nil {
		log.Printf("Error saving optimized file:- %v", err)
		return nil, err
	}

	// Create a pipe to return the optimized image as a io.ReadCloser
	rc, pw := io.Pipe()
	go func() {
		defer pw.Close()
		err := imaging.Encode(pw, optimizedImage, imaging.JPEG)
		if err != nil {
			log.Printf("Error encoding optimized image %v", err)
			return
		}
	}()

	return rc, nil
}

func (o *Optimizer) updateOptimizedFileDetails(file *models.File, optimizerParams OptimizerParams) error {
	file.OptimizedPath = optimizerParams.Path
	file.OptimizedName = optimizerParams.Name
	file.OptimizedSize = optimizerParams.Size
	file.OptimizationLevel = optimizerParams.Level
	file.Status = models.StatusCompleleted

	err := o.Repo.UpdateFile(file)
	if err != nil {
		log.Printf("Error updating optimized file details:- %v", err)
		return err
	}
	return nil
}
