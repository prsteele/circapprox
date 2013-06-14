package circapprox

import (
	"image"
	"image/color"
	"math"
)

// Set a pixel in img by blending c on top of whatever color is
// already there
func blend(pt image.Point, c color.Color, img *image.RGBA64) {
	img.Set(pt.X, pt.Y, over(c, img.At(pt.X, pt.Y)))
}

// The 'over' operator for blending colors. Color a will be painted
// over color b.
func over(a, b color.Color) color.Color {
	const M = 0xFFFF

	aR, aG, aB, aA := a.RGBA()
	bR, bG, bB, bA := b.RGBA()

	ar := float64(aR) / 0xFFFF
	ag := float64(aG) / 0xFFFF
	ab := float64(aB) / 0xFFFF
	aa := float64(aA) / 0xFFFF
	br := float64(bR) / 0xFFFF
	bg := float64(bG) / 0xFFFF
	bb := float64(bB) / 0xFFFF
	ba := float64(bA) / 0xFFFF

	oa := aa + ba*(1-aa)

	or := ar + br*(1-aa)
	og := ag + bg*(1-aa)
	ob := ab + bb*(1-aa)

	oR := uint16(0xFe * or * oa)
	oG := uint16(0xFe * og * oa)
	oB := uint16(0xFe * ob * oa)
	oA := uint16(0xFe * oa)

	return color.RGBA64{R: oR, G: oG, B: oB, A: oA}
}

type Circle struct {
	X int
	Y int
	R float64
}

func (circ Circle) points(img image.Image) chan image.Point {
	c := make(chan image.Point)

	go func() {
		rect := img.Bounds()

		x := circ.X
		y := circ.Y
		r := circ.R

		min_x := x - int(r)
		min_y := y - int(r)
		max_x := x + int(r)
		max_y := y + int(r)

		if min_x < rect.Min.X {
			min_x = rect.Min.X
		}
		if min_y < rect.Min.Y {
			min_y = rect.Min.Y
		}
		if max_x > rect.Max.X {
			max_x = rect.Max.X
		}
		if max_y > rect.Max.Y {
			max_y = rect.Max.Y
		}

		for _x := min_x; _x <= max_x; _x++ {
			for _y := min_y; _y < max_y; _y++ {
				if math.Sqrt(float64((x-_x)*(x-_x)+(y-_y)*(y-_y))) < r {
					c <- image.Point{X: _x, Y: _y}
				}
			}
		}
	}()

	return c
}
