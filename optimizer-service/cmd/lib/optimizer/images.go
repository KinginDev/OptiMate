package optimizer

import (
	"image"
	"io"
	"log"

	"github.com/disintegration/imaging"
)

type CropParams struct {
	X, Y, Width, Height int
}

func (o *Optimizer) SupportedFormats() []string {
	return []string{"jpeg", "png", "webp"}
}

func (o *Optimizer) OptimizeJPEG(fileReader io.ReadCloser, level *string, cropParams *CropParams) (image.Image, error) {
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
		img = imaging.Resize(img, img.Bounds().Dx(), img.Bounds().Dy(), imaging.Lanczos)
	case "medium":
		// Moderrate  compression, resize to half the size
		img = imaging.Resize(img, img.Bounds().Dx()/2, 0, imaging.Lanczos)
	case "high":
		// High compression, agressive resize and lossy
		img = imaging.Resize(img, img.Bounds().Dx()/4, 0, imaging.Lanczos)
	default:
		// default to medium compression if no level is provided
		img = imaging.Resize(img, img.Bounds().Dx()/2, 0, imaging.Lanczos)
	}
	return img, nil
}

func (o *Optimizer) OptimizePNG(fileReader io.ReadCloser, level *string, cropParams *CropParams) (image.Image, error) {
	// Decode the png file
	img, err := imaging.Decode((fileReader))
	if err != nil {
		log.Printf("Error decoding png file %v", err)
		return nil, err
	}

	// Check if crop parameters are provided
	if cropParams != nil {
		img = imaging.Crop(img, image.Rect(cropParams.X, cropParams.Y, cropParams.Width, cropParams.Height))
	}

	// Apply different techniques based on the optimization level
	switch *level {
	case "low":
		img = imaging.Resize(img, img.Bounds().Dx(), img.Bounds().Dy(), imaging.Lanczos)
	case "medium":
		img = imaging.Resize(img, img.Bounds().Dx()/2, 0, imaging.Lanczos)
	case "high":
		img = imaging.Resize(img, img.Bounds().Dx()/4, 0, imaging.Lanczos)
	default:
		// Default should be medium compression
		img = imaging.Resize(img, img.Bounds().Dx()/2, 0, imaging.Lanczos)
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
