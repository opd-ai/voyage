//go:build !headless

package rendering

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/opd-ai/voyage/pkg/engine"
)

// Renderer handles all Ebitengine rendering for the game.
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

// NewRenderer creates a new renderer with the given dimensions.
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

// Camera returns the current camera state.
func (r *Renderer) Camera() *Camera {
	return &r.camera
}

// SetGenre changes the renderer's palette to match the genre.
func (r *Renderer) SetGenre(genreID engine.GenreID) {
	r.BaseSystem.SetGenre(genreID)
	r.palette = DefaultPalette(genreID)
}

// Palette returns the current color palette.
func (r *Renderer) Palette() *Palette {
	return r.palette
}

// Update implements the System interface.
func (r *Renderer) Update(world *engine.World, dt float64) {
	// Rendering updates are handled in Draw, not Update
}

// Draw renders all visible entities to the screen.
func (r *Renderer) Draw(screen *ebiten.Image) {
	// Clear to background color
	screen.Fill(r.palette.Background)
}

// DrawTile draws a single tile at the given screen position.
func (r *Renderer) DrawTile(screen *ebiten.Image, x, y, tileType int) {
	img := ebiten.NewImage(r.tileSize, r.tileSize)
	fillColor := r.palette.GetTileColor(tileType)
	img.Fill(fillColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x*r.tileSize), float64(y*r.tileSize))
	screen.DrawImage(img, op)
}

// WorldToScreen converts world coordinates to screen coordinates.
func (r *Renderer) WorldToScreen(worldX, worldY float64) (screenX, screenY float64) {
	screenX = (worldX - r.camera.X) * r.camera.Zoom
	screenY = (worldY - r.camera.Y) * r.camera.Zoom
	return screenX, screenY
}

// ScreenToWorld converts screen coordinates to world coordinates.
func (r *Renderer) ScreenToWorld(screenX, screenY float64) (worldX, worldY float64) {
	worldX = screenX/r.camera.Zoom + r.camera.X
	worldY = screenY/r.camera.Zoom + r.camera.Y
	return worldX, worldY
}
