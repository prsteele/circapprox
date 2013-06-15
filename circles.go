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

func (circ Circle) Points(img image.Image) chan image.Point {
	c := make(chan image.Point)

	go func() {
		defer close(c)
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

// Returns a channel of `n' Circle objects at random locations within
// the image `img'. The radius of the circles begins as startR and
// decreases linearly until the radius is endR. Uses `rnd' as a random
// number source.
func DecreasingCircles(img image.Image, n int, startR, endR float64, rnd *rand.Rand) chan *Circle {
	c := make(chan *Circle)

	go func() {
		defer close(c)
		r := startR
		var rstep float64
		if n > 1 {
			rstep = (startR - endR) / float64(n-1)
		} else {
			rstep = 0.0
		}

		bounds := img.Bounds()
		dx := (bounds.Max.X - bounds.Min.X)
		dy := (bounds.Max.Y - bounds.Min.Y)
		for i := 0; i < n; i++ {
			x := bounds.Min.X + rnd.Int()%dx
			y := bounds.Min.Y + rnd.Int()%dy
			c <- &Circle{X: x, Y: y, R: r}
			r -= rstep
		}
	}()

	return c
}
