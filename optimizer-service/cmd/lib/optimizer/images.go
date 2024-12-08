package optimizer

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"io"
	"log"
	"optimizer-service/cmd/internal/app/interfaces"

	"github.com/disintegration/imaging"
)

type LevelConfig struct {
	ScaleFactor    float64
	ApplySharpen   bool
	ApplyBlur      bool
	JPEGQuality    int
	PNGCompression png.CompressionLevel
}

var LevelConfigs = map[string]LevelConfig{
	"low": {
		ScaleFactor:    1.0,
		ApplySharpen:   false,
		ApplyBlur:      false,
		JPEGQuality:    85,
		PNGCompression: png.DefaultCompression,
	},
	"medium": {
		ScaleFactor:    0.5,
		ApplySharpen:   true,
		ApplyBlur:      false,
		JPEGQuality:    70,
		PNGCompression: png.BestCompression,
	},
	"high": {
		ScaleFactor:    0.25,
		ApplySharpen:   true,
		ApplyBlur:      true,
		JPEGQuality:    50,
		PNGCompression: png.BestCompression,
	},
}

func (o *Optimizer) SupportedFormats() []string {
	return []string{"jpeg", "png", "webp"}
}

func (o *Optimizer) OptimizeJPEG(fileReader io.ReadCloser, level *string, cropParams *interfaces.CropParams) (image.Image, error) {
	//Decode the jpeg file
	img, err := imaging.Decode(fileReader)
	defer fileReader.Close()
	if err != nil {
		log.Printf("Error decoding jpeg file %v", err)
		return nil, err
	}

	img, err = o.optimizeImage(img, level, cropParams)
	if err != nil {
		log.Printf("Error optimizing jpeg file %v", err)
		return nil, err
	}

	// Further optimize the jpeg file
	var buf bytes.Buffer
	config := LevelConfigs[getLevelOrDefault(level)]
	err = imaging.Encode(&buf, img, imaging.JPEG, imaging.JPEGQuality(config.JPEGQuality))
	if err != nil {
		log.Printf("Error encoding jpeg file %v", err)
		return nil, err
	}

	//convert the bytes to an image
	finalImage, err := imaging.Decode(&buf)
	if err != nil {
		log.Printf("Error decoding jpeg file %v", err)
		return nil, err
	}

	return finalImage, nil
}

func (o *Optimizer) OptimizePNG(fileReader io.ReadCloser, level *string, cropParams *interfaces.CropParams) (image.Image, error) {
	// Decode the png file
	img, err := imaging.Decode(fileReader)
	defer fileReader.Close()
	if err != nil {
		log.Printf("Error decoding png file %v", err)
		return nil, err
	}

	img, err = o.optimizeImage(img, level, cropParams)
	if err != nil {
		log.Printf("Error optimizing png file %v", err)
		return nil, err
	}

	config := LevelConfigs[getLevelOrDefault(level)]
	var buf bytes.Buffer
	err = imaging.Encode(&buf, img, imaging.PNG, imaging.PNGCompressionLevel(config.PNGCompression))
	if err != nil {
		log.Printf("Error encoding png file %v", err)
		return nil, err
	}

	//convert the bytes to an image
	finalImage, err := imaging.Decode(&buf)
	if err != nil {
		log.Printf("Error decoding png file %v", err)
		return nil, err
	}

	return finalImage, nil
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
	if img == nil {
		return nil, errors.New("image must be provided")
	}
	bounds := img.Bounds()
	config := LevelConfigs[getLevelOrDefault(level)]

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

	}
	// Apply the crop
	img = imaging.Crop(img, image.Rect(cropParams.X, cropParams.Y, cropParams.Width, cropParams.Height))

	newWidth := int(float64(bounds.Dx()) * config.ScaleFactor)
	if newWidth < 1 {
		newWidth = 1
	}

	img = imaging.Resize(img, newWidth, 0, imaging.Lanczos)

	if config.ApplySharpen {
		img = imaging.Sharpen(img, 0.5)
	}

	if config.ApplyBlur {
		img = imaging.Blur(img, 0.5)
	}

	return img, nil
}

func getLevelOrDefault(level *string) string {
	if level == nil {
		return "medium"
	}

	if *level == "low" || *level == "medium" || *level == "high" {
		return *level
	}

	return "medium"
}
