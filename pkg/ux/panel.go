//go:build !headless

package ux

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// overlayCache stores pre-allocated overlay images by size and color (H-003).
type overlayCache struct {
	overlays map[overlayCacheKey]*ebiten.Image
}

type overlayCacheKey struct {
	width  int
	height int
	color  color.RGBA
}

var globalOverlayCache = &overlayCache{
	overlays: make(map[overlayCacheKey]*ebiten.Image),
}

// getOverlay returns a cached overlay image, creating it if necessary.
func (c *overlayCache) getOverlay(width, height int, col color.Color) *ebiten.Image {
	rgba := colorToRGBA(col)
	key := overlayCacheKey{width: width, height: height, color: rgba}
	if img, ok := c.overlays[key]; ok {
		return img
	}
	img := ebiten.NewImage(width, height)
	img.Fill(col)
	c.overlays[key] = img
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

// DrawBorder draws a 2-pixel border around the panel using the skin's border color.
func DrawBorder(panel *ebiten.Image, skin *UISkin) {
	w, h := panel.Bounds().Dx(), panel.Bounds().Dy()
	c := skin.PanelBorder

	for x := 0; x < w; x++ {
		panel.Set(x, 0, c)
		panel.Set(x, 1, c)
		panel.Set(x, h-1, c)
		panel.Set(x, h-2, c)
	}
	for y := 0; y < h; y++ {
		panel.Set(0, y, c)
		panel.Set(1, y, c)
		panel.Set(w-1, y, c)
		panel.Set(w-2, y, c)
	}
}

// DrawOverlay creates and draws a semi-transparent background overlay.
// Uses cached overlay image to avoid per-frame allocations (H-003).
func DrawOverlay(screen *ebiten.Image, skin *UISkin, width, height int) {
	overlay := globalOverlayCache.getOverlay(width, height, skin.PanelBackground)
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(0.7)
	screen.DrawImage(overlay, op)
}

// DrawCenteredPanel draws a panel centered on screen and returns the panel image
// and its position for content drawing.
// Uses cached panel image to avoid per-frame allocations (H-003).
func DrawCenteredPanel(screen *ebiten.Image, skin *UISkin, screenWidth, screenHeight, panelWidth, panelHeight int) (*ebiten.Image, int, int) {
	panelX := (screenWidth - panelWidth) / 2
	panelY := (screenHeight - panelHeight) / 2

	panel := globalOverlayCache.getOverlay(panelWidth, panelHeight, skin.PanelBackground)
	// Need to clear and redraw border each time since panel content changes
	panel.Fill(skin.PanelBackground)
	DrawBorder(panel, skin)

	return panel, panelX, panelY
}

// DrawPanelToScreen draws a panel at the given position.
func DrawPanelToScreen(screen, panel *ebiten.Image, panelX, panelY int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(panelX), float64(panelY))
	screen.DrawImage(panel, op)
}

// DrawCenteredText draws text centered horizontally at the given y position.
func DrawCenteredText(panel *ebiten.Image, text string, y int) {
	panelWidth := panel.Bounds().Dx()
	x := (panelWidth - len(text)*7) / 2
	ebitenutil.DebugPrintAt(panel, text, x, y)
}

// DrawInstructions draws instruction text centered at the bottom of a panel.
func DrawInstructions(panel *ebiten.Image, instructions string, panelHeight int) {
	DrawCenteredText(panel, instructions, panelHeight-20)
}
