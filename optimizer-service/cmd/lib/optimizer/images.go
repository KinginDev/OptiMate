package optimizer

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"log"
	"path/filepath"

	"github.com/disintegration/imaging"
)

type CropParams struct {
	X, Y, Width, Height int
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

	fileType, err := o.Utils.CheckFileType(fileBytes)
	if err != nil {
		log.Printf("Error checking file type:- %v", err)
		return nil, err
	}
	// Optimize the image based on the file type
	switch fileType {
	case "jpeg":

		// Create a new io.Reader from fileBytes to be used for optimization
		fileReader = io.NopCloser(bytes.NewReader(fileBytes))

		// Optimize the jpeg file
		optimizedImage, err = o.OptimizeJpeg(fileReader, oParams.Level, cropParams)
		if err != nil {
			log.Printf("Error optimizing jpeg file:- %v", err)
			return nil, err
		}
		//Encode the optimized image and wrtite it to the buffer
		err = jpeg.Encode(&optimizedBuffer, optimizedImage,
			&jpeg.Options{Quality: o.mapJPEGQuality(*oParams.Level)})
		if err != nil {
			log.Printf("Error encoding optimized image:- %v", err)
			return nil, err
		}

	case "png":

	case "webp":

	default:
		log.Printf("Unsupported file type:- %s", fileType)
	}

	fileName := filepath.Base(filePath)
	ext := filepath.Ext(fileName)
	optimizedFileName := fileName[:len(fileName)-len(ext)] + "_optimized" + ext
	err = o.Storage.Save(optimizedFileName, &optimizedBuffer)
	if err != nil {
		log.Printf("Error saving optimized file:- %v", err)
		return nil, err
	}

	// Create a pipe to return the optimized image as a io.ReadCloser
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		err := imaging.Encode(pw, optimizedImage, imaging.JPEG)
		if err != nil {
			log.Printf("Error encoding optimized image %v", err)
			return
		}
	}()

	return pr, nil
}

func (o *Optimizer) SupportedFormats() []string {
	return []string{"jpeg", "png", "webp"}
}

func (o *Optimizer) OptimizeJpeg(fileReader io.ReadCloser, level *string, cropParams *CropParams) (image.Image, error) {
	//Decode the jpeg file
	img, err := imaging.Decode(fileReader)
	if err != nil {
		log.Printf("Error decoding jpeg file %v", err)
		return nil, err
	}

	// Check if crop parameters are provided
	// Then apply the crop
	if cropParams != nil {
		img = imaging.Crop(img, image.Rect(cropParams.X, cropParams.Y, cropParams.Width, cropParams.Height))
	}

	//resize the image
	// Apply different techniques based on the optimization level
	switch *level {
	case "low":
		// Low compression, resize to the same size
		imaging.Resize(img, img.Bounds().Dx(), img.Bounds().Dy(), imaging.Lanczos)
	case "medium":
		// Moderrate  compression, resize to half the size
		imaging.Resize(img, img.Bounds().Dx()/2, 0, imaging.CatmullRom)
	case "high":
		// High compression, agressive resize and lossy
		imaging.Resize(img, img.Bounds().Dx()/4, 0, imaging.Lanczos)
	default:
		// default to medium compression if no level is provided
		imaging.Resize(img, img.Bounds().Dx()/2, 0, imaging.Lanczos)
	}
	return img, nil
}

func (o *Optimizer) mapJPEGQuality(level string) int {
	switch level {
	case "low":
		return 85
	case "medium":
		return 70
	case "high":
		return 50
	default:
		return 75
	}
}
