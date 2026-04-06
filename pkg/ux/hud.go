//go:build !headless

package ux

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/world"
	"github.com/opd-ai/voyage/pkg/resources"
)

// HUD manages the heads-up display showing resources and crew.
type HUD struct {
	skin          *UISkin
	genre         engine.GenreID
	panelWidth    int
	panelPadding  int
	barHeight     int
	resourcePanel *ebiten.Image
	crewPanel     *ebiten.Image
	minimap       *Minimap
	pixelSprite   *ebiten.Image // Cached 1x1 white pixel for efficient drawing (M-002, M-003)
}

// NewHUD creates a new HUD with default settings.
func NewHUD(genre engine.GenreID) *HUD {
	// Create cached pixel sprite for efficient drawing (M-002, M-003)
	pixelSprite := ebiten.NewImage(1, 1)
	pixelSprite.Fill(color.White)

	h := &HUD{
		skin:         DefaultSkin(genre),
		genre:        genre,
		panelWidth:   200,
		panelPadding: 8,
		barHeight:    12,
		minimap:      NewMinimap(genre, 150, 100),
		pixelSprite:  pixelSprite,
	}
	h.resourcePanel = ebiten.NewImage(h.panelWidth, 160)
	h.crewPanel = ebiten.NewImage(h.panelWidth, 200)
	return h
}

// SetGenre changes the HUD's visual theme.
func (h *HUD) SetGenre(genre engine.GenreID) {
	h.genre = genre
	h.skin = DefaultSkin(genre)
	if h.minimap != nil {
		h.minimap.SetGenre(genre)
	}
}

// GetMinimap returns the HUD's minimap component.
func (h *HUD) GetMinimap() *Minimap {
	return h.minimap
}

// SetCrisisMode enables/disables crisis mode on the minimap.
func (h *HUD) SetCrisisMode(enabled bool) {
	if h.minimap != nil {
		h.minimap.SetCrisisMode(enabled)
	}
}

// Draw renders the HUD to the screen.
func (h *HUD) Draw(screen *ebiten.Image, res *resources.Resources, party *crew.Party) {
	h.drawResourcePanel(res)
	h.drawCrewPanel(party)

	// Position resource panel at top-left
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(10, 10)
	screen.DrawImage(h.resourcePanel, op)

	// Position crew panel below resource panel
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(10, 180)
	screen.DrawImage(h.crewPanel, op2)
}

// DrawWithMinimap renders the HUD including the minimap to the screen.
func (h *HUD) DrawWithMinimap(screen *ebiten.Image, res *resources.Resources, party *crew.Party, wm *world.WorldMap, playerX, playerY int) {
	h.Draw(screen, res, party)
	if h.minimap != nil && wm != nil {
		h.minimap.Draw(screen, wm, playerX, playerY)
	}
}

// drawResourcePanel renders the resource bars.
func (h *HUD) drawResourcePanel(res *resources.Resources) {
	h.resourcePanel.Fill(h.skin.PanelBackground)
	h.drawPanelBorder(h.resourcePanel)

	// Panel title
	ebitenutil.DebugPrintAt(h.resourcePanel, h.skin.ResourcePanelName, h.panelPadding, 4)

	y := 24
	for _, rt := range resources.AllResourceTypes() {
		name := res.Name(rt)
		ratio := res.GetRatio(rt)
		status := res.GetStatus(rt)

		// Resource name and value
		value := res.Get(rt)
		maxVal := res.GetMax(rt)
		label := fmt.Sprintf("%s: %.0f/%.0f", name, value, maxVal)
		ebitenutil.DebugPrintAt(h.resourcePanel, label, h.panelPadding, y)

		// Resource bar
		barY := y + 14
		h.drawBar(h.resourcePanel, h.panelPadding, barY, h.panelWidth-h.panelPadding*2, h.barHeight, ratio, status)

		y += 22
	}
}

// drawCrewPanel renders the crew roster.
func (h *HUD) drawCrewPanel(party *crew.Party) {
	h.crewPanel.Fill(h.skin.PanelBackground)
	h.drawPanelBorder(h.crewPanel)

	// Panel title
	ebitenutil.DebugPrintAt(h.crewPanel, h.skin.CrewPanelName, h.panelPadding, 4)

	y := 24
	if party == nil {
		ebitenutil.DebugPrintAt(h.crewPanel, "No party", h.panelPadding, y)
		return
	}

	for _, member := range party.Living() {
		// Name and health
		healthPct := int(member.HealthRatio() * 100)
		label := fmt.Sprintf("%s (%d%%)", member.Name, healthPct)
		ebitenutil.DebugPrintAt(h.crewPanel, label, h.panelPadding, y)

		// Health bar
		barY := y + 14
		status := h.healthToStatus(member.HealthRatio())
		h.drawBar(h.crewPanel, h.panelPadding, barY, h.panelWidth-h.panelPadding*2, 8, member.HealthRatio(), status)

		// Skill
		skillName := crew.SkillName(member.Skill, h.genre)
		ebitenutil.DebugPrintAt(h.crewPanel, skillName, h.panelPadding, barY+10)

		y += 36
	}

	// Dead count
	deadCount := party.DeadCount()
	if deadCount > 0 {
		msg := fmt.Sprintf("Lost: %d", deadCount)
		ebitenutil.DebugPrintAt(h.crewPanel, msg, h.panelPadding, y)
	}
}

// drawPanelBorder draws a border around the panel using DrawImage for efficiency (M-002).
func (h *HUD) drawPanelBorder(panel *ebiten.Image) {
	w, ht := panel.Bounds().Dx(), panel.Bounds().Dy()
	h.drawRect(panel, 0, 0, w, 1, h.skin.PanelBorder)    // Top
	h.drawRect(panel, 0, ht-1, w, 1, h.skin.PanelBorder) // Bottom
	h.drawRect(panel, 0, 0, 1, ht, h.skin.PanelBorder)   // Left
	h.drawRect(panel, w-1, 0, 1, ht, h.skin.PanelBorder) // Right
}

// drawBar draws a horizontal bar with fill using DrawImage for efficiency (M-003).
func (h *HUD) drawBar(img *ebiten.Image, x, y, w, ht int, ratio float64, status resources.ThresholdStatus) {
	// Background
	h.drawRect(img, x, y, w, ht, h.skin.BarBackground)

	// Fill (with 1px inset on top/bottom for border effect)
	fillWidth := int(float64(w) * ratio)
	if fillWidth > 0 {
		fillColor := h.statusColor(status)
		h.drawRect(img, x, y+1, fillWidth, ht-2, fillColor)
	}
}

// drawRect draws a filled rectangle using a scaled pixel sprite (M-002, M-003).
func (h *HUD) drawRect(img *ebiten.Image, x, y, w, ht int, c color.Color) {
	if w <= 0 || ht <= 0 {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(w), float64(ht))
	op.GeoM.Translate(float64(x), float64(y))

	// Apply color using ColorScale
	r, g, b, a := c.RGBA()
	op.ColorScale.Scale(
		float32(r)/65535.0,
		float32(g)/65535.0,
		float32(b)/65535.0,
		float32(a)/65535.0,
	)

	img.DrawImage(h.pixelSprite, op)
}

// statusColor returns the appropriate color for a threshold status.
func (h *HUD) statusColor(status resources.ThresholdStatus) color.Color {
	switch status {
	case resources.StatusCritical, resources.StatusDepleted:
		return h.skin.CriticalColor
	case resources.StatusLow:
		return h.skin.WarningColor
	default:
		return h.skin.BarFill
	}
}

// healthToStatus converts health ratio to threshold status.
func (h *HUD) healthToStatus(ratio float64) resources.ThresholdStatus {
	if ratio <= 0.2 {
		return resources.StatusCritical
	}
	if ratio <= 0.4 {
		return resources.StatusLow
	}
	return resources.StatusNormal
}
