package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/prsteele/circapprox"
	"image"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Returns a function that will print a help string
func printHelp(flags *flag.FlagSet) func() {
	return func() {
		fmt.Printf("usage: %s [--in=INPUT_IMAGE] [--out=OUTPUT_IMAGE] [options]\n", os.Args[0])

		fmt.Printf("\nOptions:\n")
		flags.PrintDefaults()
	}
}

// Attempt to decode the input as an Image, where `fname' refers to a file if nonempty and standard input otherwise. Returns a pointer to the
// decoded image, a string representing the encoding used, and an
// error, if any.
func readImage(fname string) (image.Image, string, error) {
	if fname == "" {
		// Read from standard input
		if img, img_fmt, err := image.Decode(os.Stdin); err != nil {
			return nil, "", err
		} else {
			return img, img_fmt, nil
		}
	} else {
		// Try to open the file
		f, err := os.Open(fname)
		if err != nil {
			return nil, "", err
		}
		defer f.Close()

		// Try to decode the file
		if img, img_fmt, dec_err := image.Decode(f); dec_err != nil {
			return nil, "", dec_err
		} else {
			return img, img_fmt, nil
		}
	}
}

// Returns the Image to be used for output operations. If patch is
// true and `out' describes an existing image, we will load and use
// that image; note that if the image is a different size than the
// `src' image we return an error. If patch is true and `out' does
// not describe a valid image, we return an error. Otherwise, we
// return a blank image.
func outputImage(src image.Image, out string, patch bool) (*image.RGBA64, error) {
	if patch {
		if out == "" {
			return nil, errors.New("Cannot patch without an explicit output file")
		}

		if img, _, err := readImage(out); err != nil {
			return nil, err
		} else {
			src_bounds := src.Bounds()
			img_bounds := img.Bounds()
			// Check the size of the image
			if img_bounds.Size() != src_bounds.Size() {
				img_w := img_bounds.Size().X
				img_h := img_bounds.Size().Y
				src_w := src_bounds.Size().X
				src_h := src_bounds.Size().Y
				unf_msg := "Output image's size (%d x %d) does not match " +
					"the input image's size (%d x %d)"
				msg := fmt.Sprintf(unf_msg, img_w, img_h, src_w, src_h)
				return nil, errors.New(msg)
			}

			// We now need to coerce the image into the image.RGBA64
			// format. We use src.Bounds() so the coordinate systems are
			// consistent.
			ret := image.NewRGBA64(src.Bounds())
			for x := src_bounds.Min.X; x < src_bounds.Max.X; x++ {
				for y := src_bounds.Min.Y; y < src_bounds.Max.Y; y++ {
					ret.Set(x, y, img.At(x, y))
				}
			}

			return ret, nil
		}
	} else {
		// We use src.Bounds() so the coordinate systems are consistent.
		return image.NewRGBA64(src.Bounds()), nil
	}
}

// Attempt to write the Image `dst' to the file described by `out'. If
// `out' is the empty string, we write to standard output, guessing at
// the encoding using the input image format described by `src_fmt'.
func writeOutput(dst image.Image, out, src_fmt string) error {
	// Get the output writer
	var writer io.Writer
	if out == "" {
		writer = os.Stdout
	} else {
		if f, err := os.Create(out); err != nil {
			return err
		} else {
			defer f.Close()
			writer = f
		}
	}

	// Choose the output encoder. If the output file was specified, use
	// its extension to determine the encoding. If we're writing to
	// standard output, use the input encoding; if we can't output in
	// that format, just use PNG.
	var dst_fmt string
	if out == "" {
		if src_fmt == ".png" {
			dst_fmt = src_fmt
		} else if src_fmt == ".jpeg" || src_fmt == ".jpg" {
			dst_fmt = src_fmt
		} else {
			dst_fmt = ".png"
		}
	} else {
		ext := filepath.Ext(out)
		if ext == ".png" {
			dst_fmt = ".png"
		} else if ext == ".jpeg" || ext == ".jpg" {
			dst_fmt = ".jpeg"
		} else {
			msg := fmt.Sprintf("Unknown output format '%s'; please use png or jpeg",
				ext)
			return errors.New(msg)
		}
	}

	// Write the output
	if dst_fmt == ".png" {
		png.Encode(writer, dst)
	} else if dst_fmt == ".jpeg" {
		jpeg.Encode(writer, dst, &jpeg.Options{Quality: 100})
	} else {
		msg := fmt.Sprintf("Unknown output format '%s'; please use png or jpeg",
			dst_fmt)
		return errors.New(msg)
	}

	return nil
}

func validate(n int, r, a float64) error {
	if n < 0 {
		msg := fmt.Sprintf("Please specify a nonnegative number of "+
			"approximation points, not %i", n)
		return errors.New(msg)
	}

	if r < 0 {
		msg := fmt.Sprintf("Please specify a nonnegative radius of "+
			"approximation points, not %f", r)
		return errors.New(msg)
	}

	if a < 0 || a > 1 {
		msg := fmt.Sprintf("Please specify an alpha value between zero and one "+
			", not %i", a)
		return errors.New(msg)
	}

	return nil
}

// A command line flag that stores an int64 value, as well as
// recording that a value was set
type CheckedInt64 struct {
	IsSet bool
	I     int64
}

func (checked *CheckedInt64) Set(s string) error {
	// Try to parse the value as an int64
	if v, err := strconv.ParseInt(s, 10, 64); err != nil {
		return err
	} else {
		checked.IsSet = true
		checked.I = v
		return nil
	}
}

func (checked *CheckedInt64) String() string {
	return "0"
}

func main() {
	// Set up command-line flags
	flags := flag.NewFlagSet("default", flag.ContinueOnError)
	flags.Usage = printHelp(flags)

	// Input and output file flags
	in := flags.String("in", "", "The input file, or standard in if blank")
	out := flags.String("out", "", "The input file, or standard out if blank")

	// Approximation flags
	var n int
	var r float64
	var a float64
	var s CheckedInt64
	flags.IntVar(&n, "n", 100, "The number of points to approximate")
	flags.Float64Var(&r, "r", 10, "The radius of approximation points")
	flags.Float64Var(&a, "a", .75, "The alpha-value of approximated points")
	flags.Var(&s, "s", "The random number seed")

	// Should we apply our operations to an existing image?
	patch := flags.Bool("p", false, "Modify the output image, rather "+
		"than overwriting it")
	flags.BoolVar(patch, "patch", false, "Modify the output image, rather "+
		"than overwriting it")

	// Parse the command line flags
	if err := flags.Parse(os.Args[1:]); err != nil && err != flag.ErrHelp {
		os.Exit(1)
	} else if err == flag.ErrHelp {
		os.Exit(0)
	}

	// Validate the inputs
	if err := validate(n, r, a); err != nil {
		panic(err)
	}

	// Read in the source image
	src, src_fmt, src_err := readImage(*in)
	if src_err != nil {
		panic(src_err)
	}

	// The random number stream to use
	var rnd *rand.Rand
	if s.IsSet {
		rnd = rand.New(rand.NewSource(s.I))
	} else {
		rnd = rand.New(rand.NewSource(time.Now().Unix()))
	}

	// Create the destination image and approximate the source image
	var dst *image.RGBA64
	if tmp_dst, err := outputImage(src, *out, *patch); err != nil {
		panic(err)
	} else {
		dst = tmp_dst
	}

	// Perform the approximation
	circles := circapprox.UniformCircles(src, n, r, rnd)
	circapprox.Approximate(src, dst, a, circles)

	// Write the output
	if dst_err := writeOutput(dst, *out, src_fmt); dst_err != nil {
		panic(dst_err)
	}
}
