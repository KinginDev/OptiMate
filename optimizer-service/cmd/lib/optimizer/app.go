package optimizer

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"optimizer-service/cmd/internal/app/interfaces"
	"optimizer-service/cmd/internal/models"
	"optimizer-service/cmd/internal/storage"
	"optimizer-service/cmd/internal/utils"
	"path/filepath"
)

// Optimizer is a struct that defines the optimizer.
type Optimizer struct {
	Storage storage.Storage
	Utils   utils.IUtils
	Repo    interfaces.IFileRepository
}

const (
	DefaultLevel = "high"
)

// NewOptimizer is a function that creates a new optimizer.
// constructor
func InitOptimizer(storage storage.Storage, repo interfaces.IFileRepository, u utils.IUtils) *Optimizer {
	return &Optimizer{
		Storage: storage,
		Utils:   u,
		Repo:    repo,
	}
}

func (o *Optimizer) Optimize(filePath string, file *models.File, oParam *interfaces.OptimizerParams) (io.ReadCloser, error) {
	log.Printf("Starting optimization for file: %s", filePath)

	// Retrieve the file
	fileReader, err := o.Storage.Retrieve(filePath)
	if err != nil {
		log.Printf("Error retrieving file %s: %v", filePath, err)
		return nil, err
	}
	defer fileReader.Close()

	// Read file contents once
	fileBytes, err := io.ReadAll(fileReader)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		return nil, err
	}
	log.Printf("Successfully read %d bytes from file", len(fileBytes))

	// Create a new reader for file type detection
	fileTypeReader := bytes.NewReader(fileBytes)

	// Detect image format
	config, format, err := image.DecodeConfig(fileTypeReader)
	if err != nil {
		log.Printf("Error detecting image format: %v", err)
		return nil, fmt.Errorf("invalid image format: %v", err)
	}
	log.Printf("Detected image format: %s, dimensions: %dx%d", format, config.Width, config.Height)

	// Set default optimization level
	level := DefaultLevel
	if oParam.Level != nil && *oParam.Level != "" {
		level = *oParam.Level
	}
	log.Printf("Using optimization level: %s", level)

	// Set default crop params
	var cropParams *interfaces.CropParams
	if oParam.CropParams != nil && oParam.CropParams.Width > 0 && oParam.CropParams.Height > 0 {
		cropParams = oParam.CropParams
		log.Printf("Applying crop: x=%d, y=%d, width=%d, height=%d",
			cropParams.X, cropParams.Y, cropParams.Width, cropParams.Height)
	} else {
		log.Printf("No crop parameters provided or invalid dimensions")
	}

	// Create new reader from bytes for image processing
	fileReader = io.NopCloser(bytes.NewReader(fileBytes))

	var optimizedImage image.Image
	var optimizationErr error

	// Process image based on type
	switch format {
	case "jpeg":
		log.Printf("Processing JPEG image")
		optimizedImage, optimizationErr = o.OptimizeJPEG(fileReader, &level, cropParams)
	case "png":
		log.Printf("Processing PNG image")
		optimizedImage, optimizationErr = o.OptimizePNG(fileReader, &level, cropParams)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}

	if optimizationErr != nil {
		log.Printf("Error optimizing image: %v", optimizationErr)
		return nil, optimizationErr
	}

	if optimizedImage == nil {
		log.Printf("Error: optimized image is nil")
		return nil, fmt.Errorf("optimization resulted in nil image")
	}

	bounds := optimizedImage.Bounds()
	log.Printf("Optimized image dimensions: %dx%d", bounds.Dx(), bounds.Dy())

	// Create buffer for optimized image
	var optimizedBuffer bytes.Buffer

	// Encode based on file type
	switch format {
	case "jpeg":
		quality := o.mapJPEGQuality(level)
		log.Printf("Encoding JPEG with quality: %d", quality)
		err = jpeg.Encode(&optimizedBuffer, optimizedImage, &jpeg.Options{Quality: quality})
	case "png":
		log.Printf("Encoding PNG")
		err = png.Encode(&optimizedBuffer, optimizedImage)
	}

	if err != nil {
		log.Printf("Error encoding optimized image: %v", err)
		return nil, err
	}

	log.Printf("Successfully encoded optimized image, size: %d bytes", optimizedBuffer.Len())

	// Update file details
	intOptimizedSize := int64(optimizedBuffer.Len())
	oParam.Path = &filePath
	oParam.Name = &filePath
	oParam.Size = &intOptimizedSize
	oParam.Level = &level

	if err := o.updateOptimizedFileDetails(file, oParam); err != nil {
		log.Printf("Error updating file details: %v", err)
		return nil, err
	}

	// Generate optimized filename
	fileName := filepath.Base(filePath)
	ext := filepath.Ext(fileName)
	optimizedFileName := fileName[:len(fileName)-len(ext)] + "_optimized" + ext

	// Save optimized file
	if err := o.Storage.Save(optimizedFileName, &optimizedBuffer); err != nil {
		log.Printf("Error saving optimized file: %v", err)
		return nil, err
	}
	log.Printf("Successfully saved optimized file: %s", optimizedFileName)

	// Create new buffer for return
	returnBuffer := bytes.NewBuffer(optimizedBuffer.Bytes())
	return io.NopCloser(returnBuffer), nil
}

func (o *Optimizer) updateOptimizedFileDetails(file *models.File, optimizerParams *interfaces.OptimizerParams) error {
	file.OptimizedPath = optimizerParams.Path
	file.OptimizedName = optimizerParams.Name
	file.OptimizedSize = optimizerParams.Size
	file.OptimizationLevel = optimizerParams.Level
	file.Status = models.StatusCompleted

	err := o.Repo.UpdateFile(file)
	if err != nil {
		log.Printf("Error updating optimized file details:- %v", err)
		return err
	}
	return nil
}
