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
		wmv.createDestinationMarker()
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(wmv.destinationMarker, op)
}

// createDestinationMarker initializes the destination marker image.
func (wmv *WorldMapView) createDestinationMarker() {
	wmv.destinationMarker = ebiten.NewImage(wmv.tileSize, wmv.tileSize)
	cR, cG, cB, cA := wmv.extractSkinColor(wmv.skin.HighlightColor)
	pixels := wmv.buildStarPixels(cR, cG, cB, cA)
	wmv.destinationMarker.WritePixels(pixels)
}

// extractSkinColor converts a color.Color to RGBA uint8 components.
func (wmv *WorldMapView) extractSkinColor(c interface {
	RGBA() (uint32, uint32, uint32, uint32)
}) (uint8, uint8, uint8, uint8) {
	r32, g32, b32, a32 := c.RGBA()
	return uint8(r32 >> 8), uint8(g32 >> 8), uint8(b32 >> 8), uint8(a32 >> 8)
}

// buildStarPixels creates a star pattern for the destination marker.
func (wmv *WorldMapView) buildStarPixels(cR, cG, cB, cA uint8) []byte {
	size := wmv.tileSize
	half := size / 2
	pixels := make([]byte, size*size*4)
	wmv.drawCrossLines(pixels, size, half, cR, cG, cB, cA)
	wmv.drawDiagonalLines(pixels, size, cR, cG, cB, cA)
	return pixels
}

// drawCrossLines draws vertical and horizontal lines on the pixel buffer.
func (wmv *WorldMapView) drawCrossLines(pixels []byte, size, half int, cR, cG, cB, cA uint8) {
	for i := 0; i < size; i++ {
		wmv.setPixel(pixels, size, i, half, cR, cG, cB, cA) // Vertical
		wmv.setPixel(pixels, size, half, i, cR, cG, cB, cA) // Horizontal
	}
}

// drawDiagonalLines draws diagonal lines on the pixel buffer.
func (wmv *WorldMapView) drawDiagonalLines(pixels []byte, size int, cR, cG, cB, cA uint8) {
	for i := 0; i < size; i++ {
		wmv.setPixel(pixels, size, i, i, cR, cG, cB, cA)
		if size-1-i >= 0 {
			wmv.setPixel(pixels, size, i, size-1-i, cR, cG, cB, cA)
		}
	}
}

// setPixel sets an RGBA pixel at (row, col) in the pixel buffer.
func (wmv *WorldMapView) setPixel(pixels []byte, size, row, col int, cR, cG, cB, cA uint8) {
	idx := (row*size + col) * 4
	pixels[idx] = cR
	pixels[idx+1] = cG
	pixels[idx+2] = cB
	pixels[idx+3] = cA
}

// drawVesselMarker draws the vessel's position marker.
// Uses cached marker image to avoid per-frame allocations (H-002).
// Uses WritePixels for efficient initialization instead of per-pixel Set().
func (wmv *WorldMapView) drawVesselMarker(screen *ebiten.Image, x, y int) {
	if !wmv.isVesselOnScreen(x, y) {
		return
	}
	if wmv.vesselMarker == nil {
		wmv.createVesselMarker()
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(wmv.vesselMarker, op)
}

// isVesselOnScreen checks if the vessel position is within screen bounds.
func (wmv *WorldMapView) isVesselOnScreen(x, y int) bool {
	return x >= 0 && y >= 0 && x < wmv.viewWidth && y < wmv.viewHeight
}

// createVesselMarker initializes the vessel marker image.
func (wmv *WorldMapView) createVesselMarker() {
	wmv.vesselMarker = ebiten.NewImage(wmv.tileSize, wmv.tileSize)
	cR, cG, cB, cA := wmv.extractSkinColor(wmv.skin.TextPrimary)
	pixels := wmv.buildTrianglePixels(cR, cG, cB, cA)
	wmv.vesselMarker.WritePixels(pixels)
}

// buildTrianglePixels creates a filled triangle pattern for the vessel marker.
func (wmv *WorldMapView) buildTrianglePixels(cR, cG, cB, cA uint8) []byte {
	size := wmv.tileSize
	half := size / 2
	quarter := size / 4
	pixels := make([]byte, size*size*4)

	for dy := -quarter; dy <= quarter; dy++ {
		width := quarter - abs(dy)
		for dx := 0; dx <= width; dx++ {
			py := half + dy
			px := half - quarter + dx
			if py >= 0 && py < size && px >= 0 && px < size {
				wmv.setPixel(pixels, size, py, px, cR, cG, cB, cA)
			}
		}
	}
	return pixels
}

// abs returns the absolute value of an integer.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
