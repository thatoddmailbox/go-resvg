package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/thatoddmailbox/go-resvg"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input.svg> [output.png] [width] [height]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  input.svg  - Path to SVG file to render\n")
		fmt.Fprintf(os.Stderr, "  output.png - Output PNG file (optional, defaults to input name with .png extension)\n")
		fmt.Fprintf(os.Stderr, "  width      - Output width in pixels (optional, uses SVG natural size if not specified)\n")
		fmt.Fprintf(os.Stderr, "  height     - Output height in pixels (optional, uses SVG natural size if not specified)\n")
		os.Exit(1)
	}

	inputFile := os.Args[1]

	// Determine output file name
	var outputFile string
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	} else {
		ext := filepath.Ext(inputFile)
		outputFile = strings.TrimSuffix(inputFile, ext) + ".png"
	}

	// Parse optional width and height
	var width, height uint32
	var useCustomSize bool

	if len(os.Args) > 3 {
		w, err := strconv.ParseUint(os.Args[3], 10, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing width: %v\n", err)
			os.Exit(1)
		}
		width = uint32(w)

		if len(os.Args) > 4 {
			h, err := strconv.ParseUint(os.Args[4], 10, 32)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing height: %v\n", err)
				os.Exit(1)
			}
			height = uint32(h)
			useCustomSize = true
		} else {
			fmt.Fprintf(os.Stderr, "Error: height must be specified when width is provided\n")
			os.Exit(1)
		}
	}

	// Initialize resvg logging to see any warnings
	resvg.InitLog()

	// Read SVG file
	svgData, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading SVG file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Rendering SVG: %s\n", inputFile)

	var img *image.RGBA

	if useCustomSize {
		// Render with custom size
		fmt.Printf("Output size: %dx%d pixels\n", width, height)
		img, err = resvg.RenderWithSize(svgData, width, height)
	} else {
		// Use natural SVG size
		img, err = resvg.Render(svgData)
		if img != nil {
			bounds := img.Bounds()
			fmt.Printf("Output size: %dx%d pixels (natural SVG size)\n",
				bounds.Dx(), bounds.Dy())
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering SVG: %v\n", err)
		os.Exit(1)
	}

	// Create output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	// Encode as PNG
	err = png.Encode(outFile, img)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding PNG: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully rendered to: %s\n", outputFile)
}
