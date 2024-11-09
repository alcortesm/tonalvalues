package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := do(*cfg); err != nil {
		log.Fatal(err)
	}
}

func do(cfg config) error {
	if err := validateImageConfig(cfg.imagePath); err != nil {
		return fmt.Errorf("validating image config: %v", err)
	}

	original, err := loadImage(cfg.imagePath)
	if err != nil {
		return fmt.Errorf("loading image: %v", err)
	}

	gray := toGray(original)

	rowColor := mergeImagesHorizontally([]image.Image{original, gray})
	row2Tones := mergeImagesHorizontally([]image.Image{gray, gray})
	output := mergeImagesVertically([]image.Image{rowColor, row2Tones})

	if err := saveImage(output); err != nil {
		return fmt.Errorf("saving output image: %v", err)
	}

	return nil
}

func validateImageConfig(path string) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening image: %v", err)
	}

	defer func() {
		errClose := file.Close()
		if err == nil {
			err = errClose
		}
	}()

	cfg, _, err := image.DecodeConfig(file)
	if err != nil {
		return fmt.Errorf("decoding image config: %v", err)
	}

	const (
		maxWidth  = 5000
		maxHeight = 5000
		minWidth  = 10
		minHeight = 10
	)

	if cfg.Width > maxWidth {
		return fmt.Errorf("image too wide, max %d, got %d", maxWidth, cfg.Width)
	}

	if cfg.Height > maxHeight {
		return fmt.Errorf("image too tall, max %d, got %d", maxHeight, cfg.Height)
	}

	if cfg.Width < minWidth {
		return fmt.Errorf("image too thin, min %d, got %d", minWidth, cfg.Width)
	}

	if cfg.Height < minHeight {
		return fmt.Errorf("image too short, min %d, got %d", minHeight, cfg.Height)
	}

	return nil
}

func loadImage(path string) (_ image.Image, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening image: %v", err)
	}

	defer func() {
		errClose := file.Close()
		if err == nil {
			err = errClose
		}
	}()

	image, err := jpeg.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("decoding image: %v", err)
	}

	return image, nil
}

func debugImage(img image.Image) {
	fmt.Printf("image bounds: %v\n", img.Bounds())
	for x := 0; x < 2; x++ {
		for y := 0; y < 2; y++ {
			fmt.Printf("pixel at %d, %d: %v\n", x, y, img.At(x, y))
		}
	}
}

func toGray(img image.Image) image.Image {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)

	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			rgba := img.At(x, y)
			gray.Set(x, y, rgba)
		}
	}

	return gray
}

func saveImage(img image.Image) (err error) {
	file, err := os.Create("./output.jpg")
	if err != nil {
		return fmt.Errorf("creating file: %v", err)
	}

	defer func() {
		errClose := file.Close()
		if err == nil {
			err = errClose
		}
	}()

	opts := (*jpeg.Options)(nil)
	if err := jpeg.Encode(file, img, opts); err != nil {
		return fmt.Errorf("encoding jpeg: %v", err)
	}

	return nil
}

func mergeImagesHorizontally(imgs []image.Image) image.Image {
	var width, height int

	for _, i := range imgs {
		width += i.Bounds().Max.X
		height = max(height, i.Bounds().Max.Y)
	}

	result := image.NewRGBA(image.Rect(0, 0, width, height))

	offsetX := 0
	for _, img := range imgs {
		for x := 0; x < img.Bounds().Max.X; x++ {
			for y := 0; y < img.Bounds().Max.Y; y++ {
				result.Set(offsetX+x, y, img.At(x, y))
			}
		}
		offsetX += img.Bounds().Max.X
	}

	return result
}

func mergeImagesVertically(imgs []image.Image) image.Image {
	var width, height int

	for _, i := range imgs {
		width = max(width, i.Bounds().Max.X)
		height += i.Bounds().Max.Y
	}

	result := image.NewRGBA(image.Rect(0, 0, width, height))

	offsetY := 0
	for _, img := range imgs {
		for x := 0; x < img.Bounds().Max.X; x++ {
			for y := 0; y < img.Bounds().Max.Y; y++ {
				result.Set(x, offsetY+y, img.At(x, y))
			}
		}
		offsetY += img.Bounds().Max.Y
	}

	return result
}
