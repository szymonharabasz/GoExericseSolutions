package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

var pallette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0
	blackIndex = 1
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cycles := getParameterInt(r, "cycles", 5)
		size := getParameterInt(r, "size", 100)
		nframes := getParameterInt(r, "nframes", 64)
		delay := getParameterInt(r, "delay", 8)
		res := getParameterFloat(r, "res", 0.001)
		lissajous(w, cycles, size, nframes, delay, res)
	})
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func getParameterInt(r *http.Request, name string, defaultValue int) int {
	value, err := strconv.Atoi(r.URL.Query().Get(name))
	if err != nil {
		fmt.Println(err)
		return defaultValue
	}
	return value
}

func getParameterFloat(r *http.Request, name string, defaultValue float64) float64 {
	value, err := strconv.ParseFloat(r.URL.Query().Get(name), 64)
	if err != nil {
		fmt.Println(err)
		return defaultValue
	}
	return value
}

func lissajous(out io.Writer, cycles int, size int, nframes int, delay int, res float64) {
	freq := rand.Float64() * 3.0
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, pallette)
		for t := 0.0; t < float64(cycles)*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5), blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim)
}
