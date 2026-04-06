//go:build !headless

package rendering

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/opd-ai/voyage/pkg/engine"
)

// Renderer handles all Ebitengine rendering for the game.
type Renderer struct {
	engine.BaseSystem
	width     int
	height    int
	palette   *Palette
	tileSize  int
	camera    Camera
	tileCache map[tileCacheKey]*ebiten.Image
}

// tileCacheKey uniquely identifies a cached tile by type and color.
type tileCacheKey struct {
	tileType int
	color    color.RGBA
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
		tileCache: make(map[tileCacheKey]*ebiten.Image),
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
	// Dispose old tile cache images before clearing
	for _, img := range r.tileCache {
		if img != nil {
			img.Dispose()
		}
	}
	// Clear tile cache when palette changes
	r.tileCache = make(map[tileCacheKey]*ebiten.Image)
}

// Dispose releases all GPU resources held by the renderer.
// Call this when the renderer is no longer needed.
func (r *Renderer) Dispose() {
	for _, img := range r.tileCache {
		if img != nil {
			img.Dispose()
		}
	}
	r.tileCache = nil
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
	fillColor := r.palette.GetTileColor(tileType)
	img := r.getTileCached(tileType, fillColor)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x*r.tileSize), float64(y*r.tileSize))
	screen.DrawImage(img, op)
}

// getTileCached returns a cached tile image, creating it if necessary.
func (r *Renderer) getTileCached(tileType int, c color.Color) *ebiten.Image {
	rgba := colorToRGBA(c)
	key := tileCacheKey{tileType: tileType, color: rgba}
	if img, ok := r.tileCache[key]; ok {
		return img
	}
	img := ebiten.NewImage(r.tileSize, r.tileSize)
	img.Fill(c)
	r.tileCache[key] = img
	return img
}

// colorToRGBA converts a color.Color to color.RGBA for use as a cache key.
func colorToRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
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
