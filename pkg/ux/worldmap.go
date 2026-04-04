//go:build !headless

package ux

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/world"
	"github.com/opd-ai/voyage/pkg/rendering"
)

// WorldMapView manages the world map display.
type WorldMapView struct {
	skin       *UISkin
	genre      engine.GenreID
	tileSize   int
	viewWidth  int
	viewHeight int
	cameraX    int
	cameraY    int
}

// NewWorldMapView creates a new world map view.
func NewWorldMapView(genre engine.GenreID, tileSize, viewWidth, viewHeight int) *WorldMapView {
	return &WorldMapView{
		skin:       DefaultSkin(genre),
		genre:      genre,
		tileSize:   tileSize,
		viewWidth:  viewWidth,
		viewHeight: viewHeight,
		cameraX:    0,
		cameraY:    0,
	}
}

// SetGenre changes the map view's visual theme.
func (wmv *WorldMapView) SetGenre(genre engine.GenreID) {
	wmv.genre = genre
	wmv.skin = DefaultSkin(genre)
}

// CenterOn centers the camera on the given tile coordinates.
func (wmv *WorldMapView) CenterOn(tileX, tileY int) {
	tilesWide := wmv.viewWidth / wmv.tileSize
	tilesHigh := wmv.viewHeight / wmv.tileSize
	wmv.cameraX = tileX - tilesWide/2
	wmv.cameraY = tileY - tilesHigh/2
}

// Draw renders the world map to the screen.
func (wmv *WorldMapView) Draw(screen *ebiten.Image, wm *world.WorldMap, vesselX, vesselY int, renderer *rendering.Renderer) {
	if wm == nil {
		return
	}

	tilesWide := wmv.viewWidth / wmv.tileSize
	tilesHigh := wmv.viewHeight / wmv.tileSize
	palette := renderer.Palette()

	// Draw visible tiles
	for dy := 0; dy < tilesHigh; dy++ {
		for dx := 0; dx < tilesWide; dx++ {
			worldX := wmv.cameraX + dx
			worldY := wmv.cameraY + dy

			if worldX < 0 || worldY < 0 || worldX >= wm.Width || worldY >= wm.Height {
				continue
			}

			tile := wm.GetTile(worldX, worldY)
			if tile == nil {
				continue
			}
			screenX := dx * wmv.tileSize
			screenY := dy * wmv.tileSize

			// Draw terrain
			tileType := int(tile.Terrain) % len(palette.TileColors)
			renderer.DrawTile(screen, dx, dy, tileType)

			// Mark explored/unexplored (fog of war)
			if !tile.Explored {
				wmv.drawFogOverlay(screen, screenX, screenY)
			}

			// Mark destination
			if tile.Landmark != nil && tile.Landmark.Type == world.LandmarkDestination {
				wmv.drawDestinationMarker(screen, screenX, screenY)
			}
		}
	}

	// Draw vessel position
	vesselScreenX := (vesselX - wmv.cameraX) * wmv.tileSize
	vesselScreenY := (vesselY - wmv.cameraY) * wmv.tileSize
	wmv.drawVesselMarker(screen, vesselScreenX, vesselScreenY)
}

// drawFogOverlay draws a semi-transparent fog over unexplored tiles.
func (wmv *WorldMapView) drawFogOverlay(screen *ebiten.Image, x, y int) {
	fog := ebiten.NewImage(wmv.tileSize, wmv.tileSize)
	fog.Fill(wmv.skin.PanelBackground)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleAlpha(0.7)
	screen.DrawImage(fog, op)
}

// drawDestinationMarker draws a marker for the destination tile.
func (wmv *WorldMapView) drawDestinationMarker(screen *ebiten.Image, x, y int) {
	marker := ebiten.NewImage(wmv.tileSize, wmv.tileSize)

	// Draw a star-like pattern
	c := wmv.skin.HighlightColor
	half := wmv.tileSize / 2
	for i := 0; i < wmv.tileSize; i++ {
		marker.Set(half, i, c)
		marker.Set(i, half, c)
	}
	// Diagonals
	for i := 0; i < wmv.tileSize; i++ {
		marker.Set(i, i, c)
		if wmv.tileSize-1-i >= 0 {
			marker.Set(i, wmv.tileSize-1-i, c)
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(marker, op)
}

// drawVesselMarker draws the vessel's position marker.
func (wmv *WorldMapView) drawVesselMarker(screen *ebiten.Image, x, y int) {
	// Only draw if on screen
	if x < 0 || y < 0 || x >= wmv.viewWidth || y >= wmv.viewHeight {
		return
	}

	marker := ebiten.NewImage(wmv.tileSize, wmv.tileSize)

	// Draw a simple arrow or dot for the vessel
	c := wmv.skin.TextPrimary
	half := wmv.tileSize / 2
	quarter := wmv.tileSize / 4

	// Draw a filled triangle pointing right
	for dy := -quarter; dy <= quarter; dy++ {
		width := quarter - abs(dy)
		for dx := 0; dx <= width; dx++ {
			marker.Set(half-quarter+dx, half+dy, c)
		}
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(marker, op)
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
