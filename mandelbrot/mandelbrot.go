package main

import (
	"github.com/PerformLine/go-stockutil/colorutil"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
)

const iterations = 255

func main() {
	const (
		xmin, ymin, xmax, ymax = -0.5, 0.0, 0.5, 1.
		//xmin, ymin, xmax, ymax = -1.5, -1.5, 1.5, 1.5
		width, height = 2048, 2048
		supersample   = 2
	)
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		y := float64(py)/height*(ymax-ymin) + ymin
		dy := float64(1) / height * (ymax - ymin)

		for px := 0; px < width; px++ {
			x := float64(px)/width*(xmax-xmin) + xmin
			dx := float64(1) / height * (xmax - xmin)

			var avgN float64
			var avgAbsV float64
			for i := 0; i < supersample; i++ {
				for j := 0; j < supersample; j++ {
					x1 := x + (float64(i)/(supersample-1)-0.5)*dx
					y1 := y + (float64(j)/(supersample-1)-0.5)*dy
					if supersample == 1 {
						x1 = x
						y1 = y
					}
					z := complex(x1, y1)
					n, absV := mandelbrot(z)
					avgN += float64(n)
					avgAbsV += absV
				}
			}
			avgN /= supersample * supersample
			avgAbsV /= supersample * supersample
			hue, saturation, lightness := nAabsVToHSL(avgN, avgAbsV)
			red, green, blue := colorutil.HslToRgb(hue, saturation, lightness)
			var avgColor = color.RGBA{R: red, G: green, B: blue, A: 255}
			img.Set(px, py, avgColor)
		}
	}
	png.Encode(os.Stdout, img)
}

func mandelbrot(z complex128) (uint8, float64) {
	var v complex128
	for n := uint8(0); n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			return n, cmplx.Abs(v)
		}
	}
	return iterations, cmplx.Abs(v)
}

func nAabsVToHSL(n float64, absV float64) (float64, float64, float64) {
	if n >= iterations {
		return 0.0, 0.0, 0.0
	}
	if absV < 1 {
		absV = 1
	}
	hue := 270 * (1 - 10*(float64(n)+1-math.Log(math.Log2(absV)))/float64(iterations))
	if hue > 270 {
		hue = 270
	}
	if hue < 0 {
		hue = 0
	}
	return hue, 1.0, 0.5
}
