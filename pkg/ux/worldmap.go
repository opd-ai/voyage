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
	// Cached images to avoid per-frame allocations (H-001, H-002)
	fogOverlay        *ebiten.Image
	destinationMarker *ebiten.Image
	vesselMarker      *ebiten.Image
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
	// Clear cached images when skin changes (H-001, H-002)
	wmv.fogOverlay = nil
	wmv.destinationMarker = nil
	wmv.vesselMarker = nil
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

	wmv.drawVisibleTiles(screen, wm, tilesWide, tilesHigh, palette, renderer)
	wmv.drawVesselAtPosition(screen, vesselX, vesselY)
}

// drawVisibleTiles renders all tiles visible within the viewport.
func (wmv *WorldMapView) drawVisibleTiles(screen *ebiten.Image, wm *world.WorldMap, tilesWide, tilesHigh int, palette *rendering.Palette, renderer *rendering.Renderer) {
	for dy := 0; dy < tilesHigh; dy++ {
		for dx := 0; dx < tilesWide; dx++ {
			wmv.drawTileAtOffset(screen, wm, dx, dy, palette, renderer)
		}
	}
}

// drawTileAtOffset renders a single tile at the given viewport offset.
func (wmv *WorldMapView) drawTileAtOffset(screen *ebiten.Image, wm *world.WorldMap, dx, dy int, palette *rendering.Palette, renderer *rendering.Renderer) {
	worldX := wmv.cameraX + dx
	worldY := wmv.cameraY + dy

	if !wmv.isWorldCoordInBounds(worldX, worldY, wm) {
		return
	}

	tile := wm.GetTile(worldX, worldY)
	if tile == nil {
		return
	}

	screenX := dx * wmv.tileSize
	screenY := dy * wmv.tileSize

	tileType := int(tile.Terrain) % len(palette.TileColors)
	renderer.DrawTile(screen, dx, dy, tileType)

	wmv.drawTileOverlays(screen, tile, screenX, screenY)
}

// isWorldCoordInBounds checks if world coordinates are within map bounds.
func (wmv *WorldMapView) isWorldCoordInBounds(x, y int, wm *world.WorldMap) bool {
	return x >= 0 && y >= 0 && x < wm.Width && y < wm.Height
}

// drawTileOverlays renders fog and destination markers.
func (wmv *WorldMapView) drawTileOverlays(screen *ebiten.Image, tile *world.Tile, screenX, screenY int) {
	if !tile.Explored {
		wmv.drawFogOverlay(screen, screenX, screenY)
	}
	if tile.Landmark != nil && tile.Landmark.Type == world.LandmarkDestination {
		wmv.drawDestinationMarker(screen, screenX, screenY)
	}
}

// drawVesselAtPosition renders the vessel marker at the appropriate screen position.
func (wmv *WorldMapView) drawVesselAtPosition(screen *ebiten.Image, vesselX, vesselY int) {
	vesselScreenX := (vesselX - wmv.cameraX) * wmv.tileSize
	vesselScreenY := (vesselY - wmv.cameraY) * wmv.tileSize
	wmv.drawVesselMarker(screen, vesselScreenX, vesselScreenY)
}

// drawFogOverlay draws a semi-transparent fog over unexplored tiles.
// Uses cached fog image to avoid per-frame allocations (H-001).
func (wmv *WorldMapView) drawFogOverlay(screen *ebiten.Image, x, y int) {
	if wmv.fogOverlay == nil {
		wmv.fogOverlay = ebiten.NewImage(wmv.tileSize, wmv.tileSize)
		wmv.fogOverlay.Fill(wmv.skin.PanelBackground)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleAlpha(0.7)
	screen.DrawImage(wmv.fogOverlay, op)
}

// drawDestinationMarker draws a marker for the destination tile.
// Uses cached marker image to avoid per-frame allocations (H-002).
// Uses WritePixels for efficient initialization instead of per-pixel Set().
func (wmv *WorldMapView) drawDestinationMarker(screen *ebiten.Image, x, y int) {
	if wmv.destinationMarker == nil {
		wmv.destinationMarker = ebiten.NewImage(wmv.tileSize, wmv.tileSize)

		// Draw a star-like pattern using bulk pixel data
		c := wmv.skin.HighlightColor
		half := wmv.tileSize / 2
		size := wmv.tileSize
		pixels := make([]byte, size*size*4)

		// Vertical line
		for i := 0; i < size; i++ {
			idx := (i*size + half) * 4
			pixels[idx] = c.R
			pixels[idx+1] = c.G
			pixels[idx+2] = c.B
			pixels[idx+3] = c.A
		}
		// Horizontal line
		for i := 0; i < size; i++ {
			idx := (half*size + i) * 4
			pixels[idx] = c.R
			pixels[idx+1] = c.G
			pixels[idx+2] = c.B
			pixels[idx+3] = c.A
		}
		// Diagonals
		for i := 0; i < size; i++ {
			// Top-left to bottom-right
			idx := (i*size + i) * 4
			pixels[idx] = c.R
			pixels[idx+1] = c.G
			pixels[idx+2] = c.B
			pixels[idx+3] = c.A

			// Top-right to bottom-left
			if size-1-i >= 0 {
				idx2 := (i*size + (size - 1 - i)) * 4
				pixels[idx2] = c.R
				pixels[idx2+1] = c.G
				pixels[idx2+2] = c.B
				pixels[idx2+3] = c.A
			}
		}
		wmv.destinationMarker.WritePixels(pixels)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(wmv.destinationMarker, op)
}

// drawVesselMarker draws the vessel's position marker.
// Uses cached marker image to avoid per-frame allocations (H-002).
// Uses WritePixels for efficient initialization instead of per-pixel Set().
func (wmv *WorldMapView) drawVesselMarker(screen *ebiten.Image, x, y int) {
	// Only draw if on screen
	if x < 0 || y < 0 || x >= wmv.viewWidth || y >= wmv.viewHeight {
		return
	}

	if wmv.vesselMarker == nil {
		wmv.vesselMarker = ebiten.NewImage(wmv.tileSize, wmv.tileSize)

		// Draw a simple arrow or dot for the vessel using bulk pixel data
		c := wmv.skin.TextPrimary
		half := wmv.tileSize / 2
		quarter := wmv.tileSize / 4
		size := wmv.tileSize
		pixels := make([]byte, size*size*4)

		// Draw a filled triangle pointing right
		for dy := -quarter; dy <= quarter; dy++ {
			width := quarter - abs(dy)
			for dx := 0; dx <= width; dx++ {
				py := half + dy
				px := half - quarter + dx
				if py >= 0 && py < size && px >= 0 && px < size {
					idx := (py*size + px) * 4
					pixels[idx] = c.R
					pixels[idx+1] = c.G
					pixels[idx+2] = c.B
					pixels[idx+3] = c.A
				}
			}
		}
		wmv.vesselMarker.WritePixels(pixels)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(wmv.vesselMarker, op)
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
