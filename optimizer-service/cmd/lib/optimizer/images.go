package optimizer

import (
	"image"
	"io"
	"log"
	"optimizer-service/cmd/internal/app/interfaces"

	"github.com/disintegration/imaging"
)

func (o *Optimizer) SupportedFormats() []string {
	return []string{"jpeg", "png", "webp"}
}

func (o *Optimizer) OptimizeJPEG(fileReader io.ReadCloser, level *string, cropParams *interfaces.CropParams) (image.Image, error) {
	//Decode the jpeg file
	img, err := imaging.Decode(fileReader)
	if err != nil {
		log.Printf("Error decoding jpeg file %v", err)
		return nil, err
	}

	img, err = o.optimizeImage(img, level, cropParams)
	if err != nil {
		log.Printf("Error optimizing jpeg file %v", err)
		return nil, err
	}

	return img, nil
}

func (o *Optimizer) OptimizePNG(fileReader io.ReadCloser, level *string, cropParams *interfaces.CropParams) (image.Image, error) {
	// Decode the png file
	img, err := imaging.Decode(fileReader)
	if err != nil {
		log.Printf("Error decoding png file %v", err)
		return nil, err
	}

	img, err = o.optimizeImage(img, level, cropParams)
	if err != nil {
		log.Printf("Error optimizing png file %v", err)
		return nil, err
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

// optimizeImage is a helper function that optimizes an image based on the provided parameters
// It handles the cropping and resizing of the image based on the optimization level
func (p *Optimizer) optimizeImage(img image.Image, level *string, cropParams *interfaces.CropParams) (image.Image, error) {
	bounds := img.Bounds()

	// If crop parameters are not provided or invalid, use default or full image dimensions
	if cropParams == nil || (cropParams.Width <= 0 || cropParams.Height <= 0) {
		// Try to use default dimensions while maintaining aspect ratio
		cropParams = &interfaces.CropParams{
			X:      0,
			Y:      0,
			Width:  bounds.Dx(),
			Height: bounds.Dy(),
		}

		log.Printf("Using default crop dimensions: %dx%d", cropParams.Width, cropParams.Height)

		// Apply the crop
	}

	img = imaging.Crop(img, image.Rect(cropParams.X, cropParams.Y, cropParams.X+cropParams.Width, cropParams.Y+cropParams.Height))

	//resize the image
	// Apply different techniques based on the optimization level
	switch *level { // derefrence the level pointer, so we can use the value directly
	case "low":
		// Low compression, resize to the same size
		img = imaging.Resize(img, img.Bounds().Dx(), img.Bounds().Dy(), imaging.Lanczos)
	case "medium":
		// Moderate compression, resize to half the size
		img = imaging.Resize(img, img.Bounds().Dx()/2, 0, imaging.Lanczos)
	case "high":
		// High compression, aggressive resize and lossy
		img = imaging.Resize(img, img.Bounds().Dx()/4, 0, imaging.Lanczos)
	default:
		// default to medium compression if no level is provided
		img = imaging.Resize(img, img.Bounds().Dx()/2, 0, imaging.Lanczos)
	}

	return img, nil
}
