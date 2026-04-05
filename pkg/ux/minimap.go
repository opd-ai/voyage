//go:build !headless

package ux

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/world"
)

// Minimap displays a small overview of the world map in the corner of the screen.
type Minimap struct {
	skin       *UISkin
	genre      engine.GenreID
	width      int
	height     int
	tileSize   int
	image      *ebiten.Image
	alpha      float64
	crisisMode bool
}

// NewMinimap creates a new minimap with the given dimensions.
func NewMinimap(genre engine.GenreID, width, height int) *Minimap {
	return &Minimap{
		skin:     DefaultSkin(genre),
		genre:    genre,
		width:    width,
		height:   height,
		tileSize: 4,
		image:    ebiten.NewImage(width, height),
		alpha:    1.0,
	}
}

// SetGenre changes the minimap's visual theme.
func (m *Minimap) SetGenre(genre engine.GenreID) {
	m.genre = genre
	m.skin = DefaultSkin(genre)
}

// SetCrisisMode enables/disables crisis mode (fades the minimap during encounters).
func (m *Minimap) SetCrisisMode(enabled bool) {
	m.crisisMode = enabled
	if enabled {
		m.alpha = 0.3
	} else {
		m.alpha = 1.0
	}
}

// IsCrisisMode returns whether crisis mode is enabled.
func (m *Minimap) IsCrisisMode() bool {
	return m.crisisMode
}

// Draw renders the minimap to the screen.
func (m *Minimap) Draw(screen *ebiten.Image, wm *world.WorldMap, playerX, playerY int) {
	if wm == nil {
		return
	}

	m.image.Fill(m.skin.PanelBackground)
	m.drawBorder()
	m.drawTiles(wm)
	m.drawLandmarks(wm)
	m.drawPlayer(wm, playerX, playerY)
	m.drawOriginAndDestination(wm)

	// Position in top-right corner
	screenW, _ := screen.Bounds().Dx(), screen.Bounds().Dy()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(screenW-m.width-10), 10)
	op.ColorScale.ScaleAlpha(float32(m.alpha))
	screen.DrawImage(m.image, op)
}

// drawBorder draws the minimap frame.
func (m *Minimap) drawBorder() {
	borderColor := m.skin.PanelBorder
	w, h := m.width, m.height

	for x := 0; x < w; x++ {
		m.image.Set(x, 0, borderColor)
		m.image.Set(x, h-1, borderColor)
	}
	for y := 0; y < h; y++ {
		m.image.Set(0, y, borderColor)
		m.image.Set(w-1, y, borderColor)
	}
}

// drawTiles renders explored and unexplored tiles.
func (m *Minimap) drawTiles(wm *world.WorldMap) {
	scaleX := float64(m.width-2) / float64(wm.Width)
	scaleY := float64(m.height-2) / float64(wm.Height)
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}

	for y := 0; y < wm.Height; y++ {
		for x := 0; x < wm.Width; x++ {
			tile := wm.GetTile(x, y)
			if tile == nil {
				continue
			}

			px := 1 + int(float64(x)*scale)
			py := 1 + int(float64(y)*scale)

			if tile.Explored {
				m.drawExploredTile(px, py, tile)
			} else {
				m.drawFog(px, py)
			}
		}
	}
}

// drawExploredTile draws a single explored tile with terrain color.
func (m *Minimap) drawExploredTile(px, py int, tile *world.Tile) {
	c := m.terrainColor(tile.Terrain)
	for dx := 0; dx < m.tileSize && px+dx < m.width-1; dx++ {
		for dy := 0; dy < m.tileSize && py+dy < m.height-1; dy++ {
			m.image.Set(px+dx, py+dy, c)
		}
	}
}

// drawFog draws fog overlay for unexplored areas.
func (m *Minimap) drawFog(px, py int) {
	fogColor := m.skin.TextSecondary
	for dx := 0; dx < m.tileSize && px+dx < m.width-1; dx++ {
		for dy := 0; dy < m.tileSize && py+dy < m.height-1; dy++ {
			m.image.Set(px+dx, py+dy, fogColor)
		}
	}
}

// drawLandmarks draws icons for notable locations.
func (m *Minimap) drawLandmarks(wm *world.WorldMap) {
	scaleX := float64(m.width-2) / float64(wm.Width)
	scaleY := float64(m.height-2) / float64(wm.Height)
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}

	for y := 0; y < wm.Height; y++ {
		for x := 0; x < wm.Width; x++ {
			tile := wm.GetTile(x, y)
			if tile == nil || tile.Landmark == nil || !tile.Explored {
				continue
			}

			px := 1 + int(float64(x)*scale) + m.tileSize/2
			py := 1 + int(float64(y)*scale) + m.tileSize/2

			m.drawLandmarkIcon(px, py, tile.Landmark.Type)
		}
	}
}

// drawLandmarkIcon draws a specific landmark icon.
func (m *Minimap) drawLandmarkIcon(px, py int, lt world.LandmarkType) {
	c := m.landmarkColor(lt)
	// Draw a small marker (cross or dot)
	m.image.Set(px, py, c)
	m.image.Set(px-1, py, c)
	m.image.Set(px+1, py, c)
	m.image.Set(px, py-1, c)
	m.image.Set(px, py+1, c)
}

// drawPlayer draws the player position indicator.
func (m *Minimap) drawPlayer(wm *world.WorldMap, playerX, playerY int) {
	scaleX := float64(m.width-2) / float64(wm.Width)
	scaleY := float64(m.height-2) / float64(wm.Height)
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}

	px := 1 + int(float64(playerX)*scale) + m.tileSize/2
	py := 1 + int(float64(playerY)*scale) + m.tileSize/2

	c := m.skin.HighlightColor
	// Draw player as a small filled square
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if px+dx > 0 && px+dx < m.width-1 && py+dy > 0 && py+dy < m.height-1 {
				m.image.Set(px+dx, py+dy, c)
			}
		}
	}
}

// drawOriginAndDestination draws special markers for start and end points.
func (m *Minimap) drawOriginAndDestination(wm *world.WorldMap) {
	scaleX := float64(m.width-2) / float64(wm.Width)
	scaleY := float64(m.height-2) / float64(wm.Height)
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}

	// Draw origin marker
	ox := 1 + int(float64(wm.Origin.X)*scale) + m.tileSize/2
	oy := 1 + int(float64(wm.Origin.Y)*scale) + m.tileSize/2
	m.drawOriginMarker(ox, oy)

	// Draw destination marker
	dx := 1 + int(float64(wm.Destination.X)*scale) + m.tileSize/2
	dy := 1 + int(float64(wm.Destination.Y)*scale) + m.tileSize/2
	m.drawDestinationMarker(dx, dy)
}

// drawOriginMarker draws the starting point indicator.
func (m *Minimap) drawOriginMarker(px, py int) {
	c := m.skin.BarFill
	// Draw a small circle-like pattern
	offsets := []struct{ dx, dy int }{{0, -2}, {0, 2}, {-2, 0}, {2, 0}, {-1, -1}, {1, -1}, {-1, 1}, {1, 1}}
	for _, off := range offsets {
		x, y := px+off.dx, py+off.dy
		if x > 0 && x < m.width-1 && y > 0 && y < m.height-1 {
			m.image.Set(x, y, c)
		}
	}
}

// drawDestinationMarker draws the goal indicator.
func (m *Minimap) drawDestinationMarker(px, py int) {
	c := m.skin.CriticalColor
	// Draw a star-like pattern
	for d := -2; d <= 2; d++ {
		if px+d > 0 && px+d < m.width-1 {
			m.image.Set(px+d, py, c)
		}
		if py+d > 0 && py+d < m.height-1 {
			m.image.Set(px, py+d, c)
		}
	}
}

// terrainColor returns the color for a terrain type.
func (m *Minimap) terrainColor(t world.TerrainType) color.Color {
	switch t {
	case world.TerrainPlains:
		return m.skin.BarFill
	case world.TerrainForest:
		return m.skin.ButtonNormal
	case world.TerrainMountain:
		return m.skin.TextSecondary
	case world.TerrainDesert:
		return m.skin.WarningColor
	case world.TerrainRiver:
		return m.skin.HighlightColor
	case world.TerrainSwamp:
		return m.skin.ButtonPressed
	case world.TerrainRuin:
		return m.skin.CriticalColor
	default:
		return m.skin.PanelBackground
	}
}

// landmarkColor returns the color for a landmark type.
func (m *Minimap) landmarkColor(lt world.LandmarkType) color.Color {
	switch lt {
	case world.LandmarkTown:
		return m.skin.TextPrimary
	case world.LandmarkOutpost:
		return m.skin.BarFill
	case world.LandmarkRuins:
		return m.skin.WarningColor
	case world.LandmarkShrine:
		return m.skin.HighlightColor
	case world.LandmarkOrigin:
		return m.skin.BarFill
	case world.LandmarkDestination:
		return m.skin.CriticalColor
	default:
		return m.skin.TextSecondary
	}
}

// Width returns the minimap width.
func (m *Minimap) Width() int {
	return m.width
}

// Height returns the minimap height.
func (m *Minimap) Height() int {
	return m.height
}

// Genre returns the current genre.
func (m *Minimap) Genre() engine.GenreID {
	return m.genre
}
