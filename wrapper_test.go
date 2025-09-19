package resvg

import (
	"image"
	"testing"
)

func TestRender(t *testing.T) {
	// Simple SVG data
	svgData := []byte(`<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg">
		<circle cx="50" cy="50" r="40" fill="red"/>
	</svg>`)

	// Test basic render
	img, err := Render(svgData)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if img == nil {
		t.Fatal("Render returned nil image")
	}

	bounds := img.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Fatalf("Expected 100x100 image, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestRenderWithSize(t *testing.T) {
	svgData := []byte(`<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg">
		<rect x="10" y="10" width="80" height="80" fill="blue"/>
	</svg>`)

	img, err := RenderWithSize(svgData, 200, 150)
	if err != nil {
		t.Fatalf("RenderWithSize failed: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 200 || bounds.Dy() != 150 {
		t.Fatalf("Expected 200x150 image, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestOptions(t *testing.T) {
	opts := NewOptions()
	if opts == nil {
		t.Fatal("NewOptions returned nil")
	}

	// Test setting various options
	opts.SetDPI(96.0)
	opts.SetFontSize(12.0)
	opts.SetFontFamily("Arial")
	opts.SetShapeRenderingMode(ShapeRenderingGeometricPrecision)
	opts.SetTextRenderingMode(TextRenderingOptimizeLegibility)
	opts.SetImageRenderingMode(ImageRenderingOptimizeQuality)
}

func TestParseFromData(t *testing.T) {
	svgData := []byte(`<svg width="50" height="50" xmlns="http://www.w3.org/2000/svg">
		<rect width="50" height="50" fill="green"/>
	</svg>`)

	opts := NewOptions()
	tree, err := ParseFromData(svgData, opts)
	if err != nil {
		t.Fatalf("ParseFromData failed: %v", err)
	}

	if tree == nil {
		t.Fatal("ParseFromData returned nil tree")
	}

	if tree.IsEmpty() {
		t.Fatal("Tree should not be empty")
	}

	size := tree.GetImageSize()
	if size.Width != 50.0 || size.Height != 50.0 {
		t.Fatalf("Expected size 50x50, got %.1fx%.1f", size.Width, size.Height)
	}

	// Test rendering
	img := tree.Render(IdentityTransform(), 50, 50)
	if img == nil {
		t.Fatal("Render returned nil image")
	}
}

func TestInvalidSVG(t *testing.T) {
	// Test with invalid SVG
	_, err := Render([]byte("not svg"))
	if err == nil {
		t.Fatal("Expected error for invalid SVG")
	}

	// Test with empty data
	_, err = Render([]byte(""))
	if err == nil {
		t.Fatal("Expected error for empty data")
	}
}

func TestTransform(t *testing.T) {
	identity := IdentityTransform()
	if identity.A != 1.0 || identity.D != 1.0 {
		t.Fatalf("Identity transform incorrect: A=%f, D=%f", identity.A, identity.D)
	}
	if identity.B != 0.0 || identity.C != 0.0 || identity.E != 0.0 || identity.F != 0.0 {
		t.Fatal("Identity transform should have zero values for B, C, E, F")
	}
}

func TestBoundingBox(t *testing.T) {
	svgData := []byte(`<svg width="100" height="100" xmlns="http://www.w3.org/2000/svg">
		<rect x="20" y="20" width="60" height="60" fill="blue"/>
	</svg>`)

	opts := NewOptions()
	tree, err := ParseFromData(svgData, opts)
	if err != nil {
		t.Fatalf("ParseFromData failed: %v", err)
	}

	// Test object bounding box
	objBBox, exists := tree.GetObjectBBox()
	if !exists {
		t.Fatal("Object bounding box should exist")
	}

	if objBBox.X != 20.0 || objBBox.Y != 20.0 || objBBox.Width != 60.0 || objBBox.Height != 60.0 {
		t.Fatalf("Object bbox incorrect: (%.1f,%.1f) %.1fx%.1f",
			objBBox.X, objBBox.Y, objBBox.Width, objBBox.Height)
	}

	// Test image bounding box
	_, exists = tree.GetImageBBox()
	if !exists {
		t.Fatal("Image bounding box should exist")
	}
}

func TestColorChannels(t *testing.T) {
	// Test rendering with known colors
	svgData := []byte(`<svg width="2" height="2" xmlns="http://www.w3.org/2000/svg">
		<rect x="0" y="0" width="1" height="1" fill="red"/>
		<rect x="1" y="0" width="1" height="1" fill="green"/>
		<rect x="0" y="1" width="1" height="1" fill="blue"/>
		<rect x="1" y="1" width="1" height="1" fill="white"/>
	</svg>`)

	img, err := RenderWithSize(svgData, 2, 2)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Check that we have an RGBA image
	if _, ok := img.(*image.RGBA); !ok {
		t.Fatal("Expected RGBA image")
	}

	// Basic sanity check - image should have some non-zero pixels
	hasColor := false
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			if r > 0 || g > 0 || b > 0 || a > 0 {
				hasColor = true
				break
			}
		}
		if hasColor {
			break
		}
	}

	if !hasColor {
		t.Fatal("Rendered image appears to be completely transparent/black")
	}
}
