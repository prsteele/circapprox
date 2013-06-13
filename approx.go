package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
)

func over(a, b color.Color) color.RGBA {
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

	oR := uint8(0xFe * or * oa)
	oG := uint8(0xFe * og * oa)
	oB := uint8(0xFe * ob * oa)
	oA := uint8(0xFe * oa)

	return color.RGBA{R: oR, G: oG, B: oB, A: oA}
}

func blendCircle(x, y int, r float64, c color.Color, img *image.RGBA) {
	bounds := img.Bounds()

	min_x := x - int(r)
	min_y := y - int(r)
	max_x := x + int(r)
	max_y := y + int(r)

	if min_x < bounds.Min.X {
		min_x = bounds.Min.X
	}
	if min_y < bounds.Min.Y {
		min_y = bounds.Min.Y
	}
	if max_x > bounds.Max.X {
		max_x = bounds.Max.X
	}
	if max_y > bounds.Max.Y {
		max_y = bounds.Max.Y
	}

	for xx := min_x; xx <= max_x; xx++ {
		for yy := min_y; yy <= max_y; yy++ {
			if math.Sqrt(float64((x-xx)*(x-xx)+(y-yy)*(y-yy))) < r {
				img.SetRGBA(xx, yy, over(c, img.At(xx, yy)))
			}
		}
	}
}

func approx(src image.Image, n int, r *rand.Rand) image.Image {
	rect := src.Bounds()
	dst := image.NewRGBA(rect)

	white := color.NRGBA{0xFF, 0xFF, 0xFF, 0xFF}
	white_rgba := over(white, white)
	for x := rect.Min.X; x <= rect.Max.X; x++ {
		for y := rect.Min.Y; y <= rect.Max.Y; y++ {
			dst.SetRGBA(x, y, white_rgba)
		}
	}

	for i := 0; i < n; i++ {
		x := rect.Min.X + r.Int()%(rect.Max.X-rect.Min.X+1)
		y := rect.Min.Y + r.Int()%(rect.Max.Y-rect.Min.Y+1)

		c := src.At(x, y)
		cr, cg, cb, ca := c.RGBA()

		if ca > 0 {
			factor := float64(ca) / 0xFFFF
			cr = uint32(float64(cr) / factor)
			cg = uint32(float64(cg) / factor)
			cb = uint32(float64(cb) / factor)
		}

		blendCircle(x, y, 10, color.NRGBA64{R: uint16(cr), G: uint16(cg), B: uint16(cb), A: 0xFFFF / 2}, dst)
	}

	return dst
}

func main() {
	const (
		N    = 100000
		seed = 0
	)

	if len(os.Args) != 3 {
		fmt.Println("Usage: approx.go input.png output.png")
		return
	}

	inname := os.Args[1]
	outname := os.Args[2]

	infile, err := os.Open(inname)
	if err != nil {
		panic(err)
	}

	outfile, err := os.Create(outname)
	if err != nil {
		panic(err)
	}

	src, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}

	r := rand.New(rand.NewSource(seed))
	dst := approx(src, N, r)

	err = png.Encode(outfile, dst)
}
