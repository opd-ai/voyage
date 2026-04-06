//go:build !headless

package ux

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/world"
)

// Minimap displays a small overview of the world map in the corner of the screen.
// Uses caching to avoid full redraws every frame (H-010).
type Minimap struct {
	skin       *UISkin
	genre      engine.GenreID
	width      int
	height     int
	tileSize   int
	image      *ebiten.Image
	alpha      float64
	crisisMode bool

	// Caching state (H-010)
	cachedBase     *ebiten.Image // Cached static elements (border, terrain, landmarks)
	lastPlayerX    int           // Last player X position
	lastPlayerY    int           // Last player Y position
	lastWorldHash  int64         // Hash of explored tiles for change detection
	needsBaseRegen bool          // Flag to regenerate base cache
}

// NewMinimap creates a new minimap with the given dimensions.
func NewMinimap(genre engine.GenreID, width, height int) *Minimap {
	return &Minimap{
		skin:           DefaultSkin(genre),
		genre:          genre,
		width:          width,
		height:         height,
		tileSize:       4,
		image:          ebiten.NewImage(width, height),
		alpha:          1.0,
		cachedBase:     ebiten.NewImage(width, height),
		needsBaseRegen: true,
		lastPlayerX:    -1,
		lastPlayerY:    -1,
	}
}

// SetGenre changes the minimap's visual theme.
func (m *Minimap) SetGenre(genre engine.GenreID) {
	m.genre = genre
	m.skin = DefaultSkin(genre)
	m.needsBaseRegen = true // Regenerate on genre change (H-010)
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

// Draw renders the minimap to the screen using cached base when possible (H-010).
func (m *Minimap) Draw(screen *ebiten.Image, wm *world.WorldMap, playerX, playerY int) {
	if wm == nil {
		return
	}

	// Check if we need to regenerate the base cache
	worldHash := m.computeWorldHash(wm)
	if m.needsBaseRegen || worldHash != m.lastWorldHash {
		m.regenerateBaseCache(wm)
		m.lastWorldHash = worldHash
		m.needsBaseRegen = false
	}

	// Only redraw if player position changed (H-010)
	if playerX != m.lastPlayerX || playerY != m.lastPlayerY {
		// Clear the main image and draw cached base
		m.image.Clear()
		m.image.DrawImage(m.cachedBase, nil)

		// Draw dynamic elements (player, origin, destination)
		m.drawPlayer(wm, playerX, playerY)
		m.drawOriginAndDestination(wm)

		m.lastPlayerX = playerX
		m.lastPlayerY = playerY
	}

	// Position in top-right corner
	screenW, _ := screen.Bounds().Dx(), screen.Bounds().Dy()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(screenW-m.width-10), 10)
	op.ColorScale.ScaleAlpha(float32(m.alpha))
	screen.DrawImage(m.image, op)
}

// computeWorldHash creates a simple hash to detect world changes (H-010).
func (m *Minimap) computeWorldHash(wm *world.WorldMap) int64 {
	var hash int64
	for y := 0; y < wm.Height; y++ {
		for x := 0; x < wm.Width; x++ {
			tile := wm.GetTile(x, y)
			if tile != nil && tile.Explored {
				hash += int64(x*1000 + y)
			}
		}
	}
	return hash
}

// regenerateBaseCache rebuilds the static minimap elements (H-010).
func (m *Minimap) regenerateBaseCache(wm *world.WorldMap) {
	m.cachedBase.Fill(m.skin.PanelBackground)
	m.drawBorderTo(m.cachedBase)
	m.drawTilesTo(m.cachedBase, wm)
	m.drawLandmarksTo(m.cachedBase, wm)
}

// drawBorderTo draws the minimap frame to the specified image.
func (m *Minimap) drawBorderTo(img *ebiten.Image) {
	borderColor := m.skin.PanelBorder
	w, h := m.width, m.height

	for x := 0; x < w; x++ {
		img.Set(x, 0, borderColor)
		img.Set(x, h-1, borderColor)
	}
	for y := 0; y < h; y++ {
		img.Set(0, y, borderColor)
		img.Set(w-1, y, borderColor)
	}
}

// drawTilesTo renders explored and unexplored tiles to the specified image.
func (m *Minimap) drawTilesTo(img *ebiten.Image, wm *world.WorldMap) {
	scale := m.calculateScale(wm)
	for y := 0; y < wm.Height; y++ {
		for x := 0; x < wm.Width; x++ {
			m.drawTileAtTo(img, wm, x, y, scale)
		}
	}
}

// calculateScale computes the scale factor for mapping world to minimap.
func (m *Minimap) calculateScale(wm *world.WorldMap) float64 {
	scaleX := float64(m.width-2) / float64(wm.Width)
	scaleY := float64(m.height-2) / float64(wm.Height)
	if scaleY < scaleX {
		return scaleY
	}
	return scaleX
}

// drawTileAtTo draws a single tile at the given world coordinates to the specified image.
func (m *Minimap) drawTileAtTo(img *ebiten.Image, wm *world.WorldMap, x, y int, scale float64) {
	tile := wm.GetTile(x, y)
	if tile == nil {
		return
	}

	px := 1 + int(float64(x)*scale)
	py := 1 + int(float64(y)*scale)

	if tile.Explored {
		m.drawExploredTileTo(img, px, py, tile)
	} else {
		m.drawFogTo(img, px, py)
	}
}

// drawExploredTileTo draws a single explored tile with terrain color to the specified image.
func (m *Minimap) drawExploredTileTo(img *ebiten.Image, px, py int, tile *world.Tile) {
	c := m.terrainColor(tile.Terrain)
	for dx := 0; dx < m.tileSize && px+dx < m.width-1; dx++ {
		for dy := 0; dy < m.tileSize && py+dy < m.height-1; dy++ {
			img.Set(px+dx, py+dy, c)
		}
	}
}

// drawFogTo draws fog overlay for unexplored areas to the specified image.
func (m *Minimap) drawFogTo(img *ebiten.Image, px, py int) {
	fogColor := m.skin.TextSecondary
	for dx := 0; dx < m.tileSize && px+dx < m.width-1; dx++ {
		for dy := 0; dy < m.tileSize && py+dy < m.height-1; dy++ {
			img.Set(px+dx, py+dy, fogColor)
		}
	}
}

// drawLandmarksTo draws icons for notable locations to the specified image.
func (m *Minimap) drawLandmarksTo(img *ebiten.Image, wm *world.WorldMap) {
	scale := m.calculateScale(wm)
	for y := 0; y < wm.Height; y++ {
		for x := 0; x < wm.Width; x++ {
			tile := wm.GetTile(x, y)
			if tile == nil || tile.Landmark == nil || !tile.Explored {
				continue
			}

			px := 1 + int(float64(x)*scale) + m.tileSize/2
			py := 1 + int(float64(y)*scale) + m.tileSize/2

			m.drawLandmarkIconTo(img, px, py, tile.Landmark.Type)
		}
	}
}

// drawLandmarkIconTo draws a specific landmark icon to the specified image.
func (m *Minimap) drawLandmarkIconTo(img *ebiten.Image, px, py int, lt world.LandmarkType) {
	c := m.landmarkColor(lt)
	// Draw a small marker (cross or dot)
	img.Set(px, py, c)
	img.Set(px-1, py, c)
	img.Set(px+1, py, c)
	img.Set(px, py-1, c)
	img.Set(px, py+1, c)
}

// drawPlayer draws the player position indicator.
func (m *Minimap) drawPlayer(wm *world.WorldMap, playerX, playerY int) {
	scale := m.calculateScale(wm)
	px := 1 + int(float64(playerX)*scale) + m.tileSize/2
	py := 1 + int(float64(playerY)*scale) + m.tileSize/2

	c := m.skin.HighlightColor
	m.drawFilledSquare(px, py, c)
}

// drawFilledSquare draws a small filled square marker.
func (m *Minimap) drawFilledSquare(px, py int, c color.Color) {
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
	scale := m.calculateScale(wm)

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
