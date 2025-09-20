package resvg

/*
#cgo CFLAGS: -I${SRCDIR}/bin
#cgo LDFLAGS: -L${SRCDIR}/bin/linux/amd64 -lresvg -lm -lstdc++
#include "resvg.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"image"
	"math"
	"runtime"
	"unsafe"
)

// Error types
var (
	ErrNotUTF8        = errors.New("not a UTF-8 string")
	ErrFileOpenFailed = errors.New("failed to open file")
	ErrMalformedGzip  = errors.New("malformed gzip")
	ErrElementsLimit  = errors.New("elements limit reached")
	ErrInvalidSize    = errors.New("invalid size")
	ErrParsingFailed  = errors.New("parsing failed")
)

// ImageRenderingMode represents image rendering quality settings
type ImageRenderingMode int

const (
	ImageRenderingOptimizeQuality ImageRenderingMode = C.RESVG_IMAGE_RENDERING_OPTIMIZE_QUALITY
	ImageRenderingOptimizeSpeed   ImageRenderingMode = C.RESVG_IMAGE_RENDERING_OPTIMIZE_SPEED
)

// ShapeRenderingMode represents shape rendering quality settings
type ShapeRenderingMode int

const (
	ShapeRenderingOptimizeSpeed      ShapeRenderingMode = C.RESVG_SHAPE_RENDERING_OPTIMIZE_SPEED
	ShapeRenderingCrispEdges         ShapeRenderingMode = C.RESVG_SHAPE_RENDERING_CRISP_EDGES
	ShapeRenderingGeometricPrecision ShapeRenderingMode = C.RESVG_SHAPE_RENDERING_GEOMETRIC_PRECISION
)

// TextRenderingMode represents text rendering quality settings
type TextRenderingMode int

const (
	TextRenderingOptimizeSpeed      TextRenderingMode = C.RESVG_TEXT_RENDERING_OPTIMIZE_SPEED
	TextRenderingOptimizeLegibility TextRenderingMode = C.RESVG_TEXT_RENDERING_OPTIMIZE_LEGIBILITY
	TextRenderingGeometricPrecision TextRenderingMode = C.RESVG_TEXT_RENDERING_GEOMETRIC_PRECISION
)

// Transform represents a 2D transformation matrix
type Transform struct {
	A, B, C, D, E, F float32
}

// Size represents width and height dimensions
type Size struct {
	Width, Height float32
}

// Rect represents a rectangle with position and dimensions
type Rect struct {
	X, Y, Width, Height float32
}

// Options contains configuration for SVG rendering
type Options struct {
	cOpts *C.resvg_options
}

// NewOptions creates a new Options instance with default settings
func NewOptions() *Options {
	opts := &Options{
		cOpts: C.resvg_options_create(),
	}
	runtime.SetFinalizer(opts, (*Options).destroy)
	return opts
}

// SetResourcesDir sets the directory for resolving relative paths
func (o *Options) SetResourcesDir(path string) {
	if path == "" {
		C.resvg_options_set_resources_dir(o.cOpts, nil)
		return
	}
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))
	C.resvg_options_set_resources_dir(o.cOpts, cPath)
}

// SetDPI sets the target DPI for unit conversion
func (o *Options) SetDPI(dpi float32) {
	C.resvg_options_set_dpi(o.cOpts, C.float(dpi))
}

// SetStylesheet sets a CSS stylesheet to use when resolving attributes
func (o *Options) SetStylesheet(css string) {
	if css == "" {
		C.resvg_options_set_stylesheet(o.cOpts, nil)
		return
	}
	cCSS := C.CString(css)
	defer C.free(unsafe.Pointer(cCSS))
	C.resvg_options_set_stylesheet(o.cOpts, cCSS)
}

// SetFontFamily sets the default font family
func (o *Options) SetFontFamily(family string) {
	cFamily := C.CString(family)
	defer C.free(unsafe.Pointer(cFamily))
	C.resvg_options_set_font_family(o.cOpts, cFamily)
}

// SetFontSize sets the default font size
func (o *Options) SetFontSize(size float32) {
	C.resvg_options_set_font_size(o.cOpts, C.float(size))
}

// SetSerifFamily sets the serif font family
func (o *Options) SetSerifFamily(family string) {
	cFamily := C.CString(family)
	defer C.free(unsafe.Pointer(cFamily))
	C.resvg_options_set_serif_family(o.cOpts, cFamily)
}

// SetSansSerifFamily sets the sans-serif font family
func (o *Options) SetSansSerifFamily(family string) {
	cFamily := C.CString(family)
	defer C.free(unsafe.Pointer(cFamily))
	C.resvg_options_set_sans_serif_family(o.cOpts, cFamily)
}

// SetCursiveFamily sets the cursive font family
func (o *Options) SetCursiveFamily(family string) {
	cFamily := C.CString(family)
	defer C.free(unsafe.Pointer(cFamily))
	C.resvg_options_set_cursive_family(o.cOpts, cFamily)
}

// SetFantasyFamily sets the fantasy font family
func (o *Options) SetFantasyFamily(family string) {
	cFamily := C.CString(family)
	defer C.free(unsafe.Pointer(cFamily))
	C.resvg_options_set_fantasy_family(o.cOpts, cFamily)
}

// SetMonospaceFamily sets the monospace font family
func (o *Options) SetMonospaceFamily(family string) {
	cFamily := C.CString(family)
	defer C.free(unsafe.Pointer(cFamily))
	C.resvg_options_set_monospace_family(o.cOpts, cFamily)
}

// SetShapeRenderingMode sets the shape rendering method
func (o *Options) SetShapeRenderingMode(mode ShapeRenderingMode) {
	C.resvg_options_set_shape_rendering_mode(o.cOpts, C.resvg_shape_rendering(mode))
}

// SetTextRenderingMode sets the text rendering method
func (o *Options) SetTextRenderingMode(mode TextRenderingMode) {
	C.resvg_options_set_text_rendering_mode(o.cOpts, C.resvg_text_rendering(mode))
}

// SetImageRenderingMode sets the image rendering method
func (o *Options) SetImageRenderingMode(mode ImageRenderingMode) {
	C.resvg_options_set_image_rendering_mode(o.cOpts, C.resvg_image_rendering(mode))
}

// LoadFontData loads font data into the internal font database
func (o *Options) LoadFontData(data []byte) {
	if len(data) == 0 {
		return
	}
	C.resvg_options_load_font_data(o.cOpts, (*C.char)(unsafe.Pointer(&data[0])), C.uintptr_t(len(data)))
}

// LoadFontFile loads a font file into the internal font database
func (o *Options) LoadFontFile(path string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	result := C.resvg_options_load_font_file(o.cOpts, cPath)
	return cErrorToGoError(result)
}

// LoadSystemFonts loads system fonts into the internal font database
func (o *Options) LoadSystemFonts() {
	C.resvg_options_load_system_fonts(o.cOpts)
}

func (o *Options) destroy() {
	if o.cOpts != nil {
		C.resvg_options_destroy(o.cOpts)
		o.cOpts = nil
	}
}

// RenderTree represents a parsed SVG render tree
type RenderTree struct {
	cTree *C.resvg_render_tree
}

// ParseFromData parses SVG data into a render tree
func ParseFromData(data []byte, opts *Options) (*RenderTree, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	var cTree *C.resvg_render_tree
	result := C.resvg_parse_tree_from_data(
		(*C.char)(unsafe.Pointer(&data[0])),
		C.uintptr_t(len(data)),
		opts.cOpts,
		&cTree,
	)

	if err := cErrorToGoError(result); err != nil {
		return nil, err
	}

	tree := &RenderTree{cTree: cTree}
	runtime.SetFinalizer(tree, (*RenderTree).destroy)
	return tree, nil
}

// ParseFromFile parses an SVG file into a render tree
func ParseFromFile(path string, opts *Options) (*RenderTree, error) {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	var cTree *C.resvg_render_tree
	result := C.resvg_parse_tree_from_file(cPath, opts.cOpts, &cTree)

	if err := cErrorToGoError(result); err != nil {
		return nil, err
	}

	tree := &RenderTree{cTree: cTree}
	runtime.SetFinalizer(tree, (*RenderTree).destroy)
	return tree, nil
}

// IsEmpty returns true if the tree has no renderable nodes
func (t *RenderTree) IsEmpty() bool {
	return bool(C.resvg_is_image_empty(t.cTree))
}

// GetImageSize returns the natural size of the SVG
func (t *RenderTree) GetImageSize() Size {
	cSize := C.resvg_get_image_size(t.cTree)
	return Size{
		Width:  float32(cSize.width),
		Height: float32(cSize.height),
	}
}

// GetImageBBox returns the bounding box that contains all SVG elements
func (t *RenderTree) GetImageBBox() (Rect, bool) {
	var cRect C.resvg_rect
	exists := bool(C.resvg_get_image_bbox(t.cTree, &cRect))
	return Rect{
		X:      float32(cRect.x),
		Y:      float32(cRect.y),
		Width:  float32(cRect.width),
		Height: float32(cRect.height),
	}, exists
}

// GetObjectBBox returns the object bounding box (without stroke and filters)
func (t *RenderTree) GetObjectBBox() (Rect, bool) {
	var cRect C.resvg_rect
	exists := bool(C.resvg_get_object_bbox(t.cTree, &cRect))
	return Rect{
		X:      float32(cRect.x),
		Y:      float32(cRect.y),
		Width:  float32(cRect.width),
		Height: float32(cRect.height),
	}, exists
}

// Render renders the SVG tree to an RGBA image
func (t *RenderTree) Render(transform Transform, width, height uint32) *image.RGBA {
	// Create RGBA image
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	// Convert transform
	cTransform := C.resvg_transform{
		a: C.float(transform.A),
		b: C.float(transform.B),
		c: C.float(transform.C),
		d: C.float(transform.D),
		e: C.float(transform.E),
		f: C.float(transform.F),
	}

	// Render to the image buffer
	C.resvg_render(
		t.cTree,
		cTransform,
		C.uint32_t(width),
		C.uint32_t(height),
		(*C.char)(unsafe.Pointer(&img.Pix[0])),
	)

	// Convert from premultiplied alpha to straight alpha
	convertFromPremultiplied(img)

	return img
}

// RenderNode renders a specific node by ID to an RGBA image
func (t *RenderTree) RenderNode(id string, transform Transform, width, height uint32) (*image.RGBA, error) {
	cID := C.CString(id)
	defer C.free(unsafe.Pointer(cID))

	// Create RGBA image
	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	// Convert transform
	cTransform := C.resvg_transform{
		a: C.float(transform.A),
		b: C.float(transform.B),
		c: C.float(transform.C),
		d: C.float(transform.D),
		e: C.float(transform.E),
		f: C.float(transform.F),
	}

	// Render to the image buffer
	success := bool(C.resvg_render_node(
		t.cTree,
		cID,
		cTransform,
		C.uint32_t(width),
		C.uint32_t(height),
		(*C.char)(unsafe.Pointer(&img.Pix[0])),
	))

	if !success {
		return nil, fmt.Errorf("failed to render node with id '%s'", id)
	}

	// Convert from premultiplied alpha to straight alpha
	convertFromPremultiplied(img)

	return img, nil
}

func (t *RenderTree) destroy() {
	if t.cTree != nil {
		C.resvg_tree_destroy(t.cTree)
		t.cTree = nil
	}
}

// IdentityTransform returns an identity transformation matrix
func IdentityTransform() Transform {
	cTransform := C.resvg_transform_identity()
	return Transform{
		A: float32(cTransform.a),
		B: float32(cTransform.b),
		C: float32(cTransform.c),
		D: float32(cTransform.d),
		E: float32(cTransform.e),
		F: float32(cTransform.f),
	}
}

// InitLog initializes resvg logging (call once)
func InitLog() {
	C.resvg_init_log()
}

// Render is a convenience function that renders SVG data to an RGBA image
func Render(data []byte) (*image.RGBA, error) {
	opts := NewOptions()
	opts.LoadSystemFonts()
	defer opts.destroy()

	tree, err := ParseFromData(data, opts)
	if err != nil {
		return nil, err
	}
	defer tree.destroy()

	if tree.IsEmpty() {
		return nil, errors.New("SVG contains no renderable elements")
	}

	size := tree.GetImageSize()
	if size.Width <= 0 || size.Height <= 0 {
		return nil, errors.New("SVG has invalid dimensions")
	}

	return tree.Render(IdentityTransform(), uint32(size.Width), uint32(size.Height)), nil
}

// RenderWithSize renders SVG data to an RGBA image with specified dimensions
func RenderWithSize(data []byte, width, height uint32) (*image.RGBA, error) {
	opts := NewOptions()
	opts.LoadSystemFonts()
	defer opts.destroy()

	tree, err := ParseFromData(data, opts)
	if err != nil {
		return nil, err
	}
	defer tree.destroy()

	if tree.IsEmpty() {
		return nil, errors.New("SVG contains no renderable elements")
	}

	return tree.Render(IdentityTransform(), width, height), nil
}

// RenderScaledToSize renders SVG data to an RGBA image, scaling the content to fit the specified dimensions
// while preserving aspect ratio and centering it on the canvas. If the natural aspect ratio doesn't match
// the target, the content will be letterboxed (black bars on sides/top/bottom).
func RenderScaledToSize(data []byte, width, height uint32) (*image.RGBA, error) {
	opts := NewOptions()
	opts.LoadSystemFonts()
	defer opts.destroy()

	tree, err := ParseFromData(data, opts)
	if err != nil {
		return nil, err
	}
	defer tree.destroy()

	if tree.IsEmpty() {
		return nil, errors.New("SVG contains no renderable elements")
	}

	naturalSize := tree.GetImageSize()
	naturalW := float64(naturalSize.Width)
	naturalH := float64(naturalSize.Height)
	if naturalW <= 0 || naturalH <= 0 {
		return nil, errors.New("SVG has invalid natural dimensions")
	}

	targetW := float64(width)
	targetH := float64(height)

	// Compute uniform scale to fit (preserve aspect ratio)
	scaleX := targetW / naturalW
	scaleY := targetH / naturalH
	scale := math.Min(scaleX, scaleY)

	// Scaled dimensions
	scaledW := naturalW * scale
	scaledH := naturalH * scale

	// Center offsets
	tx := (targetW - scaledW) / 2.0
	ty := (targetH - scaledH) / 2.0

	// Build the transform: scale uniformly, then translate to center
	transform := Transform{
		A: float32(scale), // Scale X
		D: float32(scale), // Scale Y
		E: float32(tx),    // Translate X
		F: float32(ty),    // Translate Y
		B: 0,              // Skew X
		C: 0,              // Skew Y
	}

	return tree.Render(transform, width, height), nil
}

// Helper functions

func cErrorToGoError(result C.int32_t) error {
	switch result {
	case C.RESVG_OK:
		return nil
	case C.RESVG_ERROR_NOT_AN_UTF8_STR:
		return ErrNotUTF8
	case C.RESVG_ERROR_FILE_OPEN_FAILED:
		return ErrFileOpenFailed
	case C.RESVG_ERROR_MALFORMED_GZIP:
		return ErrMalformedGzip
	case C.RESVG_ERROR_ELEMENTS_LIMIT_REACHED:
		return ErrElementsLimit
	case C.RESVG_ERROR_INVALID_SIZE:
		return ErrInvalidSize
	case C.RESVG_ERROR_PARSING_FAILED:
		return ErrParsingFailed
	default:
		return fmt.Errorf("unknown resvg error: %d", int(result))
	}
}

// convertFromPremultiplied converts premultiplied RGBA to straight RGBA
func convertFromPremultiplied(img *image.RGBA) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			i := img.PixOffset(x, y)
			r := uint32(img.Pix[i+0])
			g := uint32(img.Pix[i+1])
			b := uint32(img.Pix[i+2])
			a := uint32(img.Pix[i+3])

			if a != 0 && a != 255 {
				// Unpremultiply
				r = (r * 255) / a
				g = (g * 255) / a
				b = (b * 255) / a

				if r > 255 {
					r = 255
				}
				if g > 255 {
					g = 255
				}
				if b > 255 {
					b = 255
				}

				img.Pix[i+0] = uint8(r)
				img.Pix[i+1] = uint8(g)
				img.Pix[i+2] = uint8(b)
			}
		}
	}
}
