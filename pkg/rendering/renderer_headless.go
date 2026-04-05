//go:build headless

package rendering

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// Renderer is a stub for headless builds.
type Renderer struct {
	engine.BaseSystem
	width    int
	height   int
	palette  *Palette
	tileSize int
	camera   Camera
}

// Camera represents the viewport into the game world.
type Camera struct {
	X, Y   float64
	Width  int
	Height int
	Zoom   float64
}

// NewRenderer creates a new renderer stub for headless builds.
func NewRenderer(width, height, tileSize int) *Renderer {
	return &Renderer{
		BaseSystem: engine.NewBaseSystem(engine.PriorityRender),
		width:      width,
		height:     height,
		palette:    DefaultPalette(engine.GenreFantasy),
		tileSize:   tileSize,
		camera: Camera{
			X:      0,
			Y:      0,
			Width:  width,
			Height: height,
			Zoom:   1.0,
		},
	}
}

// Width returns the screen width.
func (r *Renderer) Width() int {
	return r.width
}

// Height returns the screen height.
func (r *Renderer) Height() int {
	return r.height
}

// TileSize returns the tile size in pixels.
func (r *Renderer) TileSize() int {
	return r.tileSize
}

// Palette returns the current color palette.
func (r *Renderer) Palette() *Palette {
	return r.palette
}

// SetGenre changes the color palette to match the genre.
func (r *Renderer) SetGenre(genre engine.GenreID) {
	r.palette = DefaultPalette(genre)
}

// SetCamera updates the camera position.
func (r *Renderer) SetCamera(x, y float64) {
	r.camera.X = x
	r.camera.Y = y
}

// Camera returns the current camera state.
func (r *Renderer) Camera() Camera {
	return r.camera
}

// SetZoom adjusts the camera zoom level.
func (r *Renderer) SetZoom(zoom float64) {
	r.camera.Zoom = zoom
}
