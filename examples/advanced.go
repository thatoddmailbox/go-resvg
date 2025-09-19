package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/thatoddmailbox/go-resvg"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input.svg> [output_prefix]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  input.svg     - Path to SVG file to render\n")
		fmt.Fprintf(os.Stderr, "  output_prefix - Prefix for output files (optional)\n")
		fmt.Fprintf(os.Stderr, "\nThis example demonstrates various rendering options and generates multiple output files.\n")
		os.Exit(1)
	}

	inputFile := os.Args[1]

	// Determine output prefix
	var outputPrefix string
	if len(os.Args) > 2 {
		outputPrefix = os.Args[2]
	} else {
		ext := filepath.Ext(inputFile)
		outputPrefix = inputFile[:len(inputFile)-len(ext)]
	}

	// Initialize resvg logging
	resvg.InitLog()

	// Read SVG file
	svgData, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading SVG file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Processing SVG: %s\n", inputFile)
	fmt.Println("Generating multiple renderings with different options...")

	// Example 1: Default options
	fmt.Println("\n1. Rendering with default options...")
	renderWithDefaults(svgData, outputPrefix+"_default.png")

	// Example 2: Custom DPI
	fmt.Println("\n2. Rendering with high DPI (192)...")
	renderWithHighDPI(svgData, outputPrefix+"_highdpi.png")

	// Example 3: Custom rendering modes
	fmt.Println("\n3. Rendering optimized for speed...")
	renderOptimizedForSpeed(svgData, outputPrefix+"_speed.png")

	// Example 4: Custom size with scaling
	fmt.Println("\n4. Rendering at 2x scale...")
	renderAtScale(svgData, outputPrefix+"_2x.png", 2.0)

	// Example 5: Get SVG information
	fmt.Println("\n5. SVG Information:")
	displaySVGInfo(svgData)

	// Example 6: Load system fonts
	fmt.Println("\n6. Rendering with system fonts loaded...")
	renderWithSystemFonts(svgData, outputPrefix+"_with_fonts.png")

	fmt.Println("\nAll renderings complete!")
}

func renderWithDefaults(svgData []byte, outputFile string) {
	img, err := resvg.Render(svgData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	saveImage(img, outputFile)
	fmt.Printf("Saved: %s (%dx%d)\n", outputFile, img.Bounds().Dx(), img.Bounds().Dy())
}

func renderWithHighDPI(svgData []byte, outputFile string) {
	opts := resvg.NewOptions()
	opts.SetDPI(192.0) // Double the default DPI

	tree, err := resvg.ParseFromData(svgData, opts)
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}

	if tree.IsEmpty() {
		fmt.Println("SVG is empty")
		return
	}

	size := tree.GetImageSize()
	img := tree.Render(resvg.IdentityTransform(), uint32(size.Width), uint32(size.Height))

	saveImage(img, outputFile)
	fmt.Printf("Saved: %s (%dx%d)\n", outputFile, img.Bounds().Dx(), img.Bounds().Dy())
}

func renderOptimizedForSpeed(svgData []byte, outputFile string) {
	opts := resvg.NewOptions()
	opts.SetImageRenderingMode(resvg.ImageRenderingOptimizeSpeed)
	opts.SetShapeRenderingMode(resvg.ShapeRenderingOptimizeSpeed)
	opts.SetTextRenderingMode(resvg.TextRenderingOptimizeSpeed)

	tree, err := resvg.ParseFromData(svgData, opts)
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}

	if tree.IsEmpty() {
		fmt.Println("SVG is empty")
		return
	}

	size := tree.GetImageSize()
	img := tree.Render(resvg.IdentityTransform(), uint32(size.Width), uint32(size.Height))

	saveImage(img, outputFile)
	fmt.Printf("Saved: %s (%dx%d)\n", outputFile, img.Bounds().Dx(), img.Bounds().Dy())
}

func renderAtScale(svgData []byte, outputFile string, scale float32) {
	opts := resvg.NewOptions()

	tree, err := resvg.ParseFromData(svgData, opts)
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}

	if tree.IsEmpty() {
		fmt.Println("SVG is empty")
		return
	}

	size := tree.GetImageSize()
	newWidth := uint32(float32(size.Width) * scale)
	newHeight := uint32(float32(size.Height) * scale)

	// Create a scaling transform
	transform := resvg.Transform{
		A: scale, B: 0,
		C: 0, D: scale,
		E: 0, F: 0,
	}

	img := tree.Render(transform, newWidth, newHeight)

	saveImage(img, outputFile)
	fmt.Printf("Saved: %s (%dx%d at %.1fx scale)\n", outputFile, img.Bounds().Dx(), img.Bounds().Dy(), scale)
}

func displaySVGInfo(svgData []byte) {
	opts := resvg.NewOptions()
	tree, err := resvg.ParseFromData(svgData, opts)
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}

	if tree.IsEmpty() {
		fmt.Println("  SVG is empty")
		return
	}

	size := tree.GetImageSize()
	fmt.Printf("  Natural size: %.1f x %.1f\n", size.Width, size.Height)

	if objBBox, exists := tree.GetObjectBBox(); exists {
		fmt.Printf("  Object bbox: (%.1f, %.1f) %.1f x %.1f\n",
			objBBox.X, objBBox.Y, objBBox.Width, objBBox.Height)
	}

	if imgBBox, exists := tree.GetImageBBox(); exists {
		fmt.Printf("  Image bbox: (%.1f, %.1f) %.1f x %.1f\n",
			imgBBox.X, imgBBox.Y, imgBBox.Width, imgBBox.Height)
	}
}

func renderWithSystemFonts(svgData []byte, outputFile string) {
	opts := resvg.NewOptions()

	// Load system fonts (this can take a moment)
	fmt.Println("  Loading system fonts...")
	opts.LoadSystemFonts()

	// Set some font preferences
	opts.SetFontFamily("Arial")
	opts.SetFontSize(12.0)

	tree, err := resvg.ParseFromData(svgData, opts)
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}

	if tree.IsEmpty() {
		fmt.Println("SVG is empty")
		return
	}

	size := tree.GetImageSize()
	img := tree.Render(resvg.IdentityTransform(), uint32(size.Width), uint32(size.Height))

	saveImage(img, outputFile)
	fmt.Printf("Saved: %s (%dx%d with system fonts)\n", outputFile, img.Bounds().Dx(), img.Bounds().Dy())
}

func saveImage(img *image.RGBA, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		fmt.Printf("Error encoding PNG %s: %v\n", filename, err)
		return
	}
}
