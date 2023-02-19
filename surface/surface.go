package main

import (
	"fmt"
	"github.com/PerformLine/go-stockutil/colorutil"
	"log"
	"math"
	"net/http"
)

const (
	width, height = 600, 320
	cells         = 400
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6
)

type polygon struct {
	ax, ay, bx, by, cx, cy, dx, dy float64
	z                              float64
}

var sinAngle, cosAngle = math.Sin(angle), math.Cos(angle)
var minZ, maxZ = 1e8, -1e8

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		polygons := make([]polygon, 0)

		writer.Header().Set("Content-Type", "image/svg+xml")
		fmt.Fprintf(writer, "<svg xmlns='http://www.w3.org/2000/svg' "+
			"style='stroke: grey; fill: white; stroke-width: 0.0' "+
			"width='%d' height='%d'>", width, height)
		for i := 0; i < cells; i++ {
			for j := 0; j < cells; j++ {
				poly := polygon{}
				var az, bz, cz, dz float64
				poly.ax, poly.ay, az = corner(i+1, j)
				poly.bx, poly.by, bz = corner(i, j)
				poly.cx, poly.cy, cz = corner(i, j+1)
				poly.dx, poly.dy, dz = corner(i+1, j+1)
				poly.z = 0.25 * (az + bz + cz + dz)
				if poly.isCorrect() {
					polygons = append(polygons, poly)
				}
			}
		}
		for _, poly := range polygons {
			z := poly.z
			hue := 270 * (z - maxZ) / (minZ - maxZ)
			red, green, blue := colorutil.HslToRgb(hue, 1.0, 0.5)

			fmt.Fprintf(writer, "<polygon points='%g,%g %g,%g %g,%g %g,%g' "+
				"style='fill:#%02x%02x%02x'/>\n",
				poly.ax, poly.ay, poly.bx, poly.by, poly.cx, poly.cy, poly.dx, poly.dy, red, green, blue)
		}
		fmt.Fprintln(writer, "</svg>")
	})
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func correctValue(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 1) && !math.IsInf(v, -1)
}

func (p polygon) isCorrect() bool {
	return correctValue(p.ax) && correctValue(p.ay) &&
		correctValue(p.bx) && correctValue(p.by) &&
		correctValue(p.cx) && correctValue(p.cy) &&
		correctValue(p.dx) && correctValue(p.dy)
}

func corner(i, j int) (float64, float64, float64) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	z := f(x, y)
	if z > maxZ {
		maxZ = z
	}
	if z < minZ {
		minZ = z
	}

	sx := width/2 + (x-y)*cosAngle*xyscale
	sy := height/2 + (x+y)*sinAngle*xyscale - z*zscale

	return sx, sy, z

}

func f(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}
