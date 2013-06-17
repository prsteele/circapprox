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
)

// Returns a function that will print a help string
func printHelp(flags *flag.FlagSet) func() {
	return func() {
		fmt.Printf("usage: %s [--in=INPUT_IMAGE] [--out=OUTPUT_IMAGE]\n", os.Args[0])

		fmt.Printf("\nOptions:\n")
		flags.PrintDefaults()

		fmt.Printf("\nProduces a new image by applying a filter to the input\n")
		fmt.Printf("If --in is not specified, we read from standard input.\n")
		fmt.Printf("If --out is not specified, we write from standard output.\n")
	}
}

// Attempt to decode the input as an Image. Returns a pointer to the
// decoded image, a string representing the encoding used, and an
// error, if any.
func readInput(in string) (image.Image, string, error) {
	var src image.Image
	var src_fmt string
	if in == "" {
		// Read from standard input
		var err error
		if src, src_fmt, err = image.Decode(os.Stdin); err != nil {
			return src, "", err
		}
	} else {
		// Try to open the file
		if f, err := os.Open(in); err != nil {
			return src, "", err
		} else {
			var dec_err error
			if src, src_fmt, dec_err = image.Decode(f); dec_err != nil {
				return src, "", dec_err
			}
		}
	}

	return src, src_fmt, nil
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
		var err error
		if writer, err = os.Create(out); err != nil {
			return err
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

func main() {
	// Set up command-line flags
	flags := flag.NewFlagSet("default", flag.ContinueOnError)
	flags.Usage = printHelp(flags)

	// Input and output file flags
	in := flags.String("in", "", "The input file")
	out := flags.String("out", "", "The input file")

	// Approximation flags
	var n int
	var r float64
	var a float64
	var s int64
	flags.IntVar(&n, "n", 100, "The number of points to approximate")
	flags.Float64Var(&r, "r", 10, "The radius of approximation points")
	flags.Float64Var(&a, "a", .75, "The alpha-value of approximated points")
	flags.Int64Var(&s, "s", 0, "The random number seed")

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
	src, src_fmt, src_err := readInput(*in)
	if src_err != nil {
		panic(src_err)
	}

	// The random number stream to use
	rnd := rand.New(rand.NewSource(s))

	// Create the destination image and approximate the source image
	dst := image.NewRGBA64(src.Bounds())
	circles := circapprox.UniformCircles(src, n, r, rnd)
	circapprox.Approximate(src, dst, a, circles)

	// Write the output
	if dst_err := writeOutput(dst, *out, src_fmt); dst_err != nil {
		panic(dst_err)
	}
}
