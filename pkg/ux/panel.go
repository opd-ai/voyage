//go:build !headless

package ux

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

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
func DrawOverlay(screen *ebiten.Image, skin *UISkin, width, height int) {
	overlay := ebiten.NewImage(width, height)
	overlay.Fill(skin.PanelBackground)
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(0.7)
	screen.DrawImage(overlay, op)
}

// DrawCenteredPanel draws a panel centered on screen and returns the panel image
// and its position for content drawing.
func DrawCenteredPanel(screen *ebiten.Image, skin *UISkin, screenWidth, screenHeight, panelWidth, panelHeight int) (*ebiten.Image, int, int) {
	panelX := (screenWidth - panelWidth) / 2
	panelY := (screenHeight - panelHeight) / 2

	panel := ebiten.NewImage(panelWidth, panelHeight)
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
