package main

import (
	"flag"
	"fmt"
	"github.com/prsteele/circapprox"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// Returns a help string
func printHelp() {

	fmt.Printf("usage: %s [--in=INPUT_IMAGE] [--out=OUTPUT_IMAGE]\n", os.Args[0])

	fmt.Printf("\nOptions:\n")
	flag.PrintDefaults()

	fmt.Printf("\nProduces a new image by applying a filter to the input\n")
	fmt.Printf("If --in is not specified, we read from standard input.\n")
	fmt.Printf("If --out is not specified, we write from standard output.\n")
}

func main() {
	//
	// Set up command-line flags
	//
	flags := flag.NewFlagSet("default", flag.ContinueOnError)
	flags.Usage = printHelp

	// Input and output files
	in := flags.String("in", "", "The input file")
	out := flags.String("out", "", "The input file")

	// Parse the command line flags
	if err := flags.Parse(os.Args[1:]); err != nil && err != flag.ErrHelp {
		os.Exit(1)
	}

	var src *image.RGBA64
	var dst *image.RGBA64

	if *in == "" {
		// Read from standard input
		if src, _, err := image.Decode(os.Stdin); err != nil {
			panic(err)
		}
	} else {
		// Try to open the file
		if f, err := os.Open(*in); err == nil {
			if src, _, err2 := image.Decode(f); err2 != nil {
				panic(err2)
			}
		} else {
			panic(err)
		}
	}

}
