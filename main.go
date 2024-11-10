package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"math"
	"os"
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

func valueRange(img *image.Gray) (minV, maxV uint8) {
	minV = 255
	maxV = 0

	for x := 0; x < img.Bounds().Max.X; x++ {
		for y := 0; y < img.Bounds().Max.Y; y++ {
			minV = min(minV, img.GrayAt(x, y).Y)
			maxV = max(maxV, img.GrayAt(x, y).Y)
		}
	}

	return
}

func steps(n uint) *[256]uint8 {
	result := [256]uint8{}

	// identity transform if n is 0 or 1, since it doesn't really makes sense
	if n < 2 {
		for i := range result {
			result[i] = uint8(i)
		}
		return &result
	}

	// 8 bit gray doesn't allow for more than 256 tones
	if n > 256 {
		n = 256
	}

	// step widths:
	//  - n=2, stepWidth=127
	//  - n=3, stepWidth=85.3^
	//  - n=4, stepWidth=64
	//  - n=5, stepWidth=51.2
	stepWidth := 256 / float64(n)
	stepHeight := 255 / float64(n-1)

	for i := range result {
		stepIndex := math.Floor(float64(i) / stepWidth) // 0, 1, 2, ..., n-1
		v := uint8(math.Ceil(stepIndex * stepHeight))
		result[i] = v
	}

	return &result
}

func tones(img *image.Gray, transform *[256]uint8) *image.Gray {
	bounds := img.Bounds()
	result := image.NewGray(bounds)

	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			g := img.GrayAt(x, y)
			g.Y = transform[g.Y]
			result.SetGray(x, y, g)
		}
	}

	return result
}

func outputAppend(output image.Image, gray *image.Gray, n uint) image.Image {
	grayN := tones(gray, steps(n))
	bottom := mergeImagesHorizontally([]image.Image{gray, grayN})
	return mergeImagesVertically([]image.Image{output, bottom})
}
