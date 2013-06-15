package main

import (
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

func main() {
	//
	// Set up command-line flags
	//
	flags := flag.NewFlagSet("default", flag.ContinueOnError)
	flags.Usage = printHelp(flags)

	// Input and output files
	in := flags.String("in", "", "The input file")
	out := flags.String("out", "", "The input file")

	// Parse the command line flags
	if err := flags.Parse(os.Args[1:]); err != nil && err != flag.ErrHelp {
		os.Exit(1)
	} else if err == flag.ErrHelp {
		os.Exit(0)
	}

	// Read in the source image
	var src image.Image
	var src_fmt string
	if *in == "" {
		// Read from standard input
		var err error
		if src, src_fmt, err = image.Decode(os.Stdin); err != nil {
			panic(err)
		}
	} else {
		// Try to open the file
		if f, err := os.Open(*in); err == nil {
			var err2 error
			if src, src_fmt, err2 = image.Decode(f); err2 != nil {
				panic(err2)
			}
		} else {
			panic(err)
		}
	}

	// Create the destination image
	dst := image.NewRGBA64(src.Bounds())

	// Approximate the image
	const (
		N      = 2000
		startR = 100
		endR   = 5
		seed   = 0
		alpha  = .75
	)
	rnd := rand.New(rand.NewSource(seed))
	circles := circapprox.DecreasingCircles(src, N, startR, endR, rnd)
	circapprox.Approximate(src, dst, alpha, circles)

	// Get the output write
	var writer io.Writer
	if *out == "" {
		writer = os.Stdout
	} else {
		var err error
		if writer, err = os.Create(*out); err != nil {
			panic(err)
		}
	}

	// Choose the output encoder. If the output file was specified, use
	// its extension to determine the encoding. If we're writing to
	// standard output, use the input encoding; if we can't output in
	// that format, just use PNG.
	var dst_fmt string
	if *out == "" {
		if src_fmt == ".png" {
			dst_fmt = src_fmt
		} else if src_fmt == ".jpeg" {
			dst_fmt = src_fmt
		} else {
			dst_fmt = ".png"
		}
	} else {
		ext := filepath.Ext(*out)
		if ext == ".png" {
			dst_fmt = ".png"
		} else if ext == ".jpeg" {
			dst_fmt = ".jpeg"
		} else {
			msg := fmt.Sprintf("Unknown output format '%s'; please use png or jpeg",
				ext)
			panic(msg)
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
		panic(msg)
	}
}
