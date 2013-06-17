package circapprox

import (
	"image"
	"math"
	"math/rand"
)

// A Circle represents a circle in a plane; Circle objects can be used
// as a bitwise image mask, blocking pixels outside Circle.
type Circle struct {
	X int
	Y int
	R float64
}

func (circ Circle) Points(img image.Image) []image.Point {
	c := make([]image.Point, 0)

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

	for _x := min_x; _x < max_x; _x++ {
		for _y := min_y; _y < max_y; _y++ {
			if math.Sqrt(float64((x-_x)*(x-_x)+(y-_y)*(y-_y))) < r {
				c = append(c, image.Point{X: _x, Y: _y})
			}
		}
	}

	return c
}

// Returns a point chosen such that each pixel on the image has an
// equally likely chance of being painted; we do this by sampling
// points up to one radius away from the edges of the image.
func unifPoint(img image.Image, r int, rnd *rand.Rand) (int, int) {
	bounds := img.Bounds()

	dx := (bounds.Max.X - bounds.Min.X)
	dy := (bounds.Max.Y - bounds.Min.Y)
	x := bounds.Min.X + rnd.Int()%dx
	y := bounds.Min.Y + rnd.Int()%dy

	return x, y
}

// Produces `n' randomly places Circle objects with radius `r'.
func UniformCircles(img image.Image, n int, r float64, rnd *rand.Rand) []Circle {
	c := make([]Circle, 0)
	for i := 0; i < n; i++ {
		x, y := unifPoint(img, int(r), rnd)
		c = append(c, Circle{X: x, Y: y, R: r})
	}

	return c
}

// Returns a channel of `n' Circle objects at random locations within
// the image `img'. The radius of the circles begins as startR and
// decreases linearly until the radius is endR. Uses `rnd' as a random
// number source.
func DecreasingCircles(img image.Image, n int, startR, endR float64, rnd *rand.Rand) []Circle {
	c := make([]Circle, 0)

	r := startR
	var rstep float64
	if n > 1 {
		rstep = (startR - endR) / float64(n-1)
	} else {
		rstep = 0.0
	}

	for i := 0; i < n; i++ {
		x, y := unifPoint(img, int(r), rnd)
		c = append(c, Circle{X: x, Y: y, R: r})
		r -= rstep
	}

	return c
}
