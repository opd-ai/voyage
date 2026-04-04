//go:build !headless

package ux

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
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
}

// NewHUD creates a new HUD with default settings.
func NewHUD(genre engine.GenreID) *HUD {
	h := &HUD{
		skin:         DefaultSkin(genre),
		genre:        genre,
		panelWidth:   200,
		panelPadding: 8,
		barHeight:    12,
	}
	h.resourcePanel = ebiten.NewImage(h.panelWidth, 160)
	h.crewPanel = ebiten.NewImage(h.panelWidth, 200)
	return h
}

// SetGenre changes the HUD's visual theme.
func (h *HUD) SetGenre(genre engine.GenreID) {
	h.genre = genre
	h.skin = DefaultSkin(genre)
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

// drawPanelBorder draws a border around the panel.
func (h *HUD) drawPanelBorder(panel *ebiten.Image) {
	w, ht := panel.Bounds().Dx(), panel.Bounds().Dy()
	borderColor := h.skin.PanelBorder

	// Top and bottom
	for x := 0; x < w; x++ {
		panel.Set(x, 0, borderColor)
		panel.Set(x, ht-1, borderColor)
	}
	// Left and right
	for y := 0; y < ht; y++ {
		panel.Set(0, y, borderColor)
		panel.Set(w-1, y, borderColor)
	}
}

// drawBar draws a horizontal bar with fill.
func (h *HUD) drawBar(img *ebiten.Image, x, y, w, ht int, ratio float64, status resources.ThresholdStatus) {
	// Background
	for dx := 0; dx < w; dx++ {
		for dy := 0; dy < ht; dy++ {
			img.Set(x+dx, y+dy, h.skin.BarBackground)
		}
	}

	// Fill
	fillWidth := int(float64(w) * ratio)
	fillColor := h.statusColor(status)
	for dx := 0; dx < fillWidth; dx++ {
		for dy := 1; dy < ht-1; dy++ {
			img.Set(x+dx, y+dy, fillColor)
		}
	}
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
