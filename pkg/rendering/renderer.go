package rendering

import (
	"image/color"

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
func (r *Renderer) DrawTile(screen *ebiten.Image, x, y int, tileType int) {
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
	return
}

// ScreenToWorld converts screen coordinates to world coordinates.
func (r *Renderer) ScreenToWorld(screenX, screenY float64) (worldX, worldY float64) {
	worldX = screenX/r.camera.Zoom + r.camera.X
	worldY = screenY/r.camera.Zoom + r.camera.Y
	return
}

// Palette defines the color scheme for rendering.
type Palette struct {
	Background color.Color
	Foreground color.Color
	Primary    color.Color
	Secondary  color.Color
	Accent     color.Color
	Warning    color.Color
	Danger     color.Color
	TileColors []color.Color
}

// DefaultPalette returns the default palette for the given genre.
func DefaultPalette(genre engine.GenreID) *Palette {
	switch genre {
	case engine.GenreScifi:
		return &Palette{
			Background: color.RGBA{10, 10, 30, 255},
			Foreground: color.RGBA{200, 220, 255, 255},
			Primary:    color.RGBA{50, 100, 200, 255},
			Secondary:  color.RGBA{100, 50, 150, 255},
			Accent:     color.RGBA{0, 255, 200, 255},
			Warning:    color.RGBA{255, 200, 0, 255},
			Danger:     color.RGBA{255, 50, 50, 255},
			TileColors: []color.Color{
				color.RGBA{5, 5, 20, 255},   // void
				color.RGBA{30, 50, 80, 255}, // nebula
				color.RGBA{50, 50, 60, 255}, // asteroid
				color.RGBA{80, 100, 120, 255}, // station
			},
		}
	case engine.GenreHorror:
		return &Palette{
			Background: color.RGBA{20, 15, 15, 255},
			Foreground: color.RGBA{180, 160, 150, 255},
			Primary:    color.RGBA{100, 40, 40, 255},
			Secondary:  color.RGBA{60, 50, 40, 255},
			Accent:     color.RGBA{200, 50, 50, 255},
			Warning:    color.RGBA{200, 150, 50, 255},
			Danger:     color.RGBA{150, 0, 0, 255},
			TileColors: []color.Color{
				color.RGBA{30, 25, 25, 255},  // wasteland
				color.RGBA{50, 40, 35, 255},  // ruins
				color.RGBA{40, 50, 40, 255},  // toxic
				color.RGBA{60, 55, 50, 255},  // shelter
			},
		}
	case engine.GenreCyberpunk:
		return &Palette{
			Background: color.RGBA{15, 15, 25, 255},
			Foreground: color.RGBA{200, 200, 220, 255},
			Primary:    color.RGBA{255, 0, 100, 255},
			Secondary:  color.RGBA{0, 200, 255, 255},
			Accent:     color.RGBA{255, 255, 0, 255},
			Warning:    color.RGBA{255, 150, 0, 255},
			Danger:     color.RGBA{255, 0, 50, 255},
			TileColors: []color.Color{
				color.RGBA{20, 20, 30, 255},  // slum
				color.RGBA{40, 40, 60, 255},  // street
				color.RGBA{60, 50, 80, 255},  // market
				color.RGBA{80, 80, 100, 255}, // tower
			},
		}
	case engine.GenrePostapoc:
		return &Palette{
			Background: color.RGBA{35, 30, 25, 255},
			Foreground: color.RGBA{200, 180, 150, 255},
			Primary:    color.RGBA{150, 100, 50, 255},
			Secondary:  color.RGBA{100, 80, 60, 255},
			Accent:     color.RGBA{200, 150, 50, 255},
			Warning:    color.RGBA{200, 100, 50, 255},
			Danger:     color.RGBA{180, 50, 30, 255},
			TileColors: []color.Color{
				color.RGBA{45, 40, 35, 255},   // dust
				color.RGBA{60, 55, 45, 255},   // sand
				color.RGBA{80, 70, 55, 255},   // scrapyard
				color.RGBA{100, 90, 70, 255},  // settlement
			},
		}
	default: // Fantasy
		return &Palette{
			Background: color.RGBA{30, 40, 30, 255},
			Foreground: color.RGBA{230, 220, 200, 255},
			Primary:    color.RGBA{80, 120, 80, 255},
			Secondary:  color.RGBA{120, 100, 60, 255},
			Accent:     color.RGBA{200, 180, 100, 255},
			Warning:    color.RGBA{200, 150, 50, 255},
			Danger:     color.RGBA{180, 50, 50, 255},
			TileColors: []color.Color{
				color.RGBA{60, 90, 60, 255},   // plains
				color.RGBA{40, 80, 40, 255},   // forest
				color.RGBA{100, 90, 80, 255},  // mountain
				color.RGBA{80, 100, 80, 255},  // town
			},
		}
	}
}

// GetTileColor returns the color for a tile type.
func (p *Palette) GetTileColor(tileType int) color.Color {
	if tileType >= 0 && tileType < len(p.TileColors) {
		return p.TileColors[tileType]
	}
	return p.Background
}
