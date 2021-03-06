package util

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"io"
)

// Grayscale returns black and white image
func Grayscale(file io.Reader) (*image.Gray, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	xMax, yMax := bounds.Max.X, bounds.Max.Y

	for x := 0; x < xMax; x++ {
		for y := 0; y < yMax; y++ {
			rgba := img.At(x, y)
			gray.Set(x, y, rgba)
		}
	}

	return gray, nil
}

// OtsuThreshold computes the Otsu's threshold for a b&w image
func OtsuThreshold(m *image.Gray) uint8 {
	hist := histogram(m)
	sum := 0
	for i, v := range hist {
		sum += i * v
	}
	wB, wF := 0, len(m.Pix)
	var sumB, sumF int
	maxVariance := 0.0
	threshold := uint8(0)
	for t := 0; t < 256; t++ {
		wB += hist[t]
		wF -= hist[t]
		if wB == 0 {
			continue
		}
		if wF == 0 {
			return threshold
		}
		sumB += t * hist[t]
		sumF = sum - sumB
		mB := float64(sumB) / float64(wB)
		mF := float64(sumF) / float64(wF)
		betweenVariance := float64(wB*wF) * (mB - mF) * (mB - mF)
		if betweenVariance > maxVariance {
			maxVariance = betweenVariance
			threshold = uint8(t)
		}
	}
	return threshold
}

func histogram(m *image.Gray) []int {
	hist := make([]int, 256)
	count := len(m.Pix)
	for i := 0; i < count; i++ {
		hist[m.Pix[i]]++
	}
	return hist
}
