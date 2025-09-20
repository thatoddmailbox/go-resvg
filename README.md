# go-resvg

A Go wrapper for the [resvg](https://github.com/linebender/resvg) library. Currently only supports linux/amd64.

## Installation

```bash
go get github.com/thatoddmailbox/go-resvg
```

## Quick start

### Simple SVG rendering

```go
package main

import (
    "fmt"
    "image/png"
    "os"

    "github.com/thatoddmailbox/go-resvg"
)

func main() {
    // Read SVG data
    svgData := []byte(`<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg">
        <circle cx="50" cy="50" r="40" fill="red"/>
    </svg>`)

    // Render to RGBA image using natural SVG size
    img, err := resvg.Render(svgData)
    if err != nil {
        panic(err)
    }

    // Save as PNG
    file, _ := os.Create("output.png")
    defer file.Close()
    png.Encode(file, img)
}
```

> **Note:** The `Render()` and `RenderWithSize()` convenience functions use your system fonts and do not set default fonts for the various font styles (serif, sans-serif, ...). If your SVG contains text elements, it's recommended to create your own `Options` struct to specify the appropriate font settings.

### Custom size rendering

```go
// Render with specific dimensions
img, err := resvg.RenderWithSize(svgData, 400, 300)
if err != nil {
    panic(err)
}
```

## Advanced usage

### Custom rendering options

```go
package main

import (
    "github.com/thatoddmailbox/go-resvg"
)

func main() {
    // Create custom options
    opts := resvg.NewOptions()

    // Set DPI for better quality
    opts.SetDPI(192.0)

    // Set rendering modes for speed vs. quality
    opts.SetImageRenderingMode(resvg.ImageRenderingOptimizeQuality)
    opts.SetShapeRenderingMode(resvg.ShapeRenderingGeometricPrecision)
    opts.SetTextRenderingMode(resvg.TextRenderingOptimizeLegibility)

    // Load system fonts
    opts.LoadSystemFonts()

    // Set font preferences
    opts.SetFontFamily("Arial")
    opts.SetFontSize(14.0)

    // Set specific font families for different CSS font types
    // You should make sure these fonts exist on your system!
    opts.SetSerifFamily("Times New Roman")
    opts.SetSansSerifFamily("Arial")
    opts.SetMonospaceFamily("Courier New")
    opts.SetCursiveFamily("Comic Sans MS")
    opts.SetFantasyFamily("Impact")

    // Parse SVG with custom options
    tree, err := resvg.ParseFromData(svgData, opts)
    if err != nil {
        panic(err)
    }

    // Get SVG information
    size := tree.GetImageSize()
    fmt.Printf("SVG size: %.1fx%.1f\n", size.Width, size.Height)

    // Render with identity transform
    img := tree.Render(resvg.IdentityTransform(), uint32(size.Width), uint32(size.Height))
}
```

### Transform and scaling

```go
// Create a scaling transform
transform := resvg.Transform{
    A: 2.0, B: 0.0,  // 2x scale
    C: 0.0, D: 2.0,
    E: 0.0, F: 0.0,  // no translation
}

// Render with transform
img := tree.Render(transform, 400, 300)
```

### Working with files

```go
// Parse from file
tree, err := resvg.ParseFromFile("input.svg", opts)
if err != nil {
    panic(err)
}

// Load custom font
err = opts.LoadFontFile("/path/to/font.ttf")
if err != nil {
    fmt.Printf("Failed to load font: %v\n", err)
}
```

## API reference

### Types

- **`Options`** - Configuration for SVG parsing and rendering
- **`RenderTree`** - Parsed SVG representation
- **`Transform`** - 2D transformation matrix
- **`Size`** - Width and height dimensions
- **`Rect`** - Rectangle with position and size

### Rendering modes

#### Image rendering
- `ImageRenderingOptimizeQuality` - Better quality (default)
- `ImageRenderingOptimizeSpeed` - Faster rendering

#### Shape rendering
- `ShapeRenderingOptimizeSpeed` - Faster rendering
- `ShapeRenderingCrispEdges` - Sharp edges
- `ShapeRenderingGeometricPrecision` - Best quality (default)

#### Text rendering
- `TextRenderingOptimizeSpeed` - Faster rendering
- `TextRenderingOptimizeLegibility` - Better readability (default)
- `TextRenderingGeometricPrecision` - Best quality

### Functions

#### Simple API
- `Render(data []byte) (*image.RGBA, error)` - Render SVG at natural size
- `RenderWithSize(data []byte, width, height uint32) (*image.RGBA, error)` - Render at custom size

#### Advanced API
- `NewOptions() *Options` - Create new options
- `ParseFromData(data []byte, opts *Options) (*RenderTree, error)` - Parse SVG from data
- `ParseFromFile(path string, opts *Options) (*RenderTree, error)` - Parse SVG from file
- `IdentityTransform() Transform` - Create identity transformation
- `InitLog()` - Initialize resvg logging

#### Options methods
- `SetDPI(dpi float32)` - Set target DPI
- `SetResourcesDir(path string)` - Set directory for relative paths
- `SetStylesheet(css string)` - Set CSS stylesheet for attribute resolution
- `SetFontFamily(family string)` - Set default font family
- `SetFontSize(size float32)` - Set default font size
- `SetSerifFamily(family string)` - Set serif font family (default: "Times New Roman")
- `SetSansSerifFamily(family string)` - Set sans-serif font family (default: "Arial")
- `SetCursiveFamily(family string)` - Set cursive font family (default: "Comic Sans MS")
- `SetFantasyFamily(family string)` - Set fantasy font family (default: "Papyrus" on macOS, "Impact" elsewhere)
- `SetMonospaceFamily(family string)` - Set monospace font family (default: "Courier New")
- `SetShapeRenderingMode(mode ShapeRenderingMode)` - Set shape rendering mode (see above list)
- `SetTextRenderingMode(mode TextRenderingMode)` - Set text rendering mode (see above list)
- `SetImageRenderingMode(mode ImageRenderingMode)` - Set image rendering mode (see above list)
- `LoadSystemFonts()` - Load system fonts
- `LoadFontFile(path string) error` - Load font from file
- `LoadFontData(data []byte)` - Load font from memory

#### RenderTree methods
- `Render(transform Transform, width, height uint32) *image.RGBA` - Render full SVG
- `RenderNode(id string, transform Transform, width, height uint32) (*image.RGBA, error)` - Render specific node
- `GetImageSize() Size` - Get natural SVG size
- `GetImageBBox() (Rect, bool)` - Get bounding box including all elements
- `GetObjectBBox() (Rect, bool)` - Get object bounding box (excludes stroke/filters)
- `IsEmpty() bool` - Check if SVG has renderable content

## Examples

The `examples/` directory contains several demonstration programs:

- **`svg2png.go`** - Simple SVG to PNG converter with optional custom sizing
- **`advanced.go`** - Demonstrates various rendering options and features
- **`test.svg`** - Sample SVG file for testing

### Running Examples

```bash
# Simple conversion
cd examples
go run svg2png.go test.svg output.png

# Custom size
go run svg2png.go test.svg output.png 400 300

# Advanced examples with multiple outputs
go run advanced.go test.svg demo
```

## Performance tips

1. **Reuse Options objects** when rendering multiple SVGs with the same settings
2. **Load system fonts once** and reuse the Options object
3. **Use speed-optimized rendering modes** for real-time applications
4. **Cache parsed RenderTree objects** when rendering the same SVG multiple times

## Error handling

The library provides specific error types for different failure modes:

```go
var (
    ErrNotUTF8        = errors.New("not a UTF-8 string")
    ErrFileOpenFailed = errors.New("failed to open file")
    ErrMalformedGzip  = errors.New("malformed gzip")
    ErrElementsLimit  = errors.New("elements limit reached")
    ErrInvalidSize    = errors.New("invalid size")
    ErrParsingFailed  = errors.New("parsing failed")
)
```

## Platform support

Currently, this package includes pre-compiled resvg binaries for:

- Linux/amd64

For other platforms, you'll need to:

1. Compile resvg as a static C library
2. Place the library in `bin/{platform}/`
3. Update the cgo directives accordingly

## Dependencies

- **resvg** - High-quality SVG rendering library (included as binary)
- **Go 1.19+** - Required for the Go implementation

## Troubleshooting

### Common issues

**"library not found" errors**: Ensure the resvg library is in the correct location for your platform.

**Font rendering issues**: Make sure to load appropriate fonts using `LoadSystemFonts()` or `LoadFontFile()`.

**Build errors**: Ensure you have a C compiler installed (gcc, clang, etc.).