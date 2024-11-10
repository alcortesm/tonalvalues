package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"

	"github.com/alcortesm/tonalvalues/staircase"
)

// run like this:
//
// ./tonalvalues example.jpg
//
// then check out the output.jpg file

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

	minValue, maxValue := valueRange(gray)
	fmt.Printf("min_value=%d, max_value=%d", minValue, maxValue)

	output := mergeImagesHorizontally([]image.Image{original, gray})
	output = outputAppend(output, gray, 2)
	output = outputAppend(output, gray, 3)
	output = outputAppend(output, gray, 4)
	output = outputAppend(output, gray, 5)
	output = outputAppend(output, gray, 6)

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

func toGray(img image.Image) *image.Gray {
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

func valueRange(img *image.Gray) (minV, maxV uint) {
	minV = 255
	maxV = 0

	for x := 0; x < img.Bounds().Max.X; x++ {
		for y := 0; y < img.Bounds().Max.Y; y++ {
			minV = min(minV, uint(img.GrayAt(x, y).Y))
			maxV = max(maxV, uint(img.GrayAt(x, y).Y))
		}
	}

	return
}

// tones generates a new grayscale image from img, by limiting the number of grayscale values to n.
func tones(img *image.Gray, n uint) *image.Gray {
	bounds := img.Bounds()
	result := image.NewGray(bounds)

	minValue, maxValue := valueRange(img)
	transformer, err := staircase.New(minValue, maxValue, n)
	if err != nil {
		panic(err)
	}

	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			g := img.GrayAt(x, y)
			g.Y = uint8(transformer.Transform(int(g.Y)))
			result.SetGray(x, y, g)
		}
	}

	return result
}

func outputAppend(output image.Image, gray *image.Gray, n uint) image.Image {
	grayN := tones(gray, n)
	bottom := mergeImagesHorizontally([]image.Image{gray, grayN})
	return mergeImagesVertically([]image.Image{output, bottom})
}
