package circapprox

import (
	"image"
	"image/color"
)

// Approximates the image `src' based on the provided circles. Each
// circle will result in a colors from the source image being blended
// into the resulting image `dst'.
func Approximate(src image.Image, dst *image.RGBA64, alpha float64, circles chan *Circle) {
	for circle := range circles {
		// Get the color of the image at this point, and adjust the
		// alpha channel
		c := src.At(circle.X, circle.Y)

		r, g, b, a := c.RGBA()

		if a > 0 {
			factor := float64(a) / 0xFFFF
			r = uint32(float64(r) / factor)
			g = uint32(float64(g) / factor)
			b = uint32(float64(b) / factor)
		}
		a = uint32(alpha * 0xFFFF)

		cc := color.NRGBA64{
			R: uint16(r), G: uint16(g), B: uint16(b), A: uint16(a)}

		for pt := range circle.Points(src) {
			blend(pt, cc, dst)
		}
	}
}

// Set a pixel in img by blending c on top of whatever color is
// already there
func blend(pt image.Point, c color.Color, img *image.RGBA64) {
	cc := over(c, img.At(pt.X, pt.Y))
	img.Set(pt.X, pt.Y, cc)
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

	oR := uint16(0xFFFF * or * oa)
	oG := uint16(0xFFFF * og * oa)
	oB := uint16(0xFFFF * ob * oa)
	oA := uint16(0xFFFF * oa)

	return color.RGBA64{R: oR, G: oG, B: oB, A: oA}
}
