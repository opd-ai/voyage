//go:build !headless

package ux

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/resources"
	"github.com/opd-ai/voyage/pkg/vessel"
)

// CargoScreen displays cargo hold contents with weight/volume limits.
type CargoScreen struct {
	skin         *UISkin
	genre        engine.GenreID
	screenWidth  int
	screenHeight int
	panelWidth   int
	panelHeight  int
	scrollOffset int
	visible      bool
}

// NewCargoScreen creates a new cargo management screen.
func NewCargoScreen(genre engine.GenreID, screenWidth, screenHeight int) *CargoScreen {
	return &CargoScreen{
		skin:         DefaultSkin(genre),
		genre:        genre,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		panelWidth:   400,
		panelHeight:  350,
		scrollOffset: 0,
		visible:      false,
	}
}

// SetGenre changes the screen's visual theme.
func (cs *CargoScreen) SetGenre(genre engine.GenreID) {
	cs.genre = genre
	cs.skin = DefaultSkin(genre)
}

// Show makes the cargo screen visible.
func (cs *CargoScreen) Show() {
	cs.visible = true
	cs.scrollOffset = 0
}

// Hide makes the cargo screen hidden.
func (cs *CargoScreen) Hide() {
	cs.visible = false
}

// IsVisible returns whether the screen is currently visible.
func (cs *CargoScreen) IsVisible() bool {
	return cs.visible
}

// ScrollUp scrolls the cargo list up.
func (cs *CargoScreen) ScrollUp() {
	if cs.scrollOffset > 0 {
		cs.scrollOffset--
	}
}

// ScrollDown scrolls the cargo list down.
func (cs *CargoScreen) ScrollDown() {
	cs.scrollOffset++
}

// Draw renders the cargo screen.
func (cs *CargoScreen) Draw(screen *ebiten.Image, hold *vessel.CargoHold, modules *vessel.ModuleSystem) {
	if !cs.visible {
		return
	}

	DrawOverlay(screen, cs.skin, cs.screenWidth, cs.screenHeight)
	panel, panelX, panelY := DrawCenteredPanel(screen, cs.skin, cs.screenWidth, cs.screenHeight, cs.panelWidth, cs.panelHeight)

	DrawCenteredText(panel, cs.getTitle(), 12)
	cs.drawCapacityBars(panel, hold, modules)
	cs.drawCargoList(panel, hold)
	DrawInstructions(panel, "UP/DOWN to scroll, ESC to close", cs.panelHeight)

	DrawPanelToScreen(screen, panel, panelX, panelY)
}

// getTitle returns the genre-appropriate title.
func (cs *CargoScreen) getTitle() string {
	titles := map[engine.GenreID]string{
		engine.GenreFantasy:   "Wagon Inventory",
		engine.GenreScifi:     "Cargo Bay",
		engine.GenreHorror:    "Vehicle Storage",
		engine.GenreCyberpunk: "Smuggler's Hold",
		engine.GenrePostapoc:  "Scrap Hold",
	}
	if title, ok := titles[cs.genre]; ok {
		return title
	}
	return "Cargo Hold"
}

// drawCapacityBars draws the weight and volume capacity bars.
func (cs *CargoScreen) drawCapacityBars(panel *ebiten.Image, hold *vessel.CargoHold, modules *vessel.ModuleSystem) {
	padding := 16
	barWidth := cs.panelWidth - padding*2
	barHeight := 14

	// Weight bar
	weightY := 40
	weightLabel := fmt.Sprintf("Weight: %d / %d", hold.UsedWeight(), hold.WeightLimit())
	ebitenutil.DebugPrintAt(panel, weightLabel, padding, weightY)
	weightRatio := hold.WeightRatio()
	cs.drawBar(panel, padding, weightY+16, barWidth, barHeight, weightRatio, cs.ratioToStatus(weightRatio))

	// Volume bar
	volumeY := 80
	volumeLabel := fmt.Sprintf("Volume: %d / %d", hold.UsedVolume(), hold.VolumeLimit())
	ebitenutil.DebugPrintAt(panel, volumeLabel, padding, volumeY)
	volumeRatio := hold.VolumeRatio()
	cs.drawBar(panel, padding, volumeY+16, barWidth, barHeight, volumeRatio, cs.ratioToStatus(volumeRatio))

	// Tier info
	tierY := 120
	tierLabel := fmt.Sprintf("Hold Tier: %d / 5", hold.Tier())
	ebitenutil.DebugPrintAt(panel, tierLabel, padding, tierY)
}

// drawCargoList draws the list of cargo items.
func (cs *CargoScreen) drawCargoList(panel *ebiten.Image, hold *vessel.CargoHold) {
	padding := 16
	startY := 145
	lineHeight := 20
	maxVisible := 8

	items := hold.Items()
	if len(items) == 0 {
		ebitenutil.DebugPrintAt(panel, "No cargo", padding, startY)
		return
	}

	// Header
	ebitenutil.DebugPrintAt(panel, "Item                Qty    Wt    Vol", padding, startY)
	startY += lineHeight

	// Ensure scroll doesn't exceed items
	if cs.scrollOffset > len(items)-maxVisible {
		cs.scrollOffset = len(items) - maxVisible
	}
	if cs.scrollOffset < 0 {
		cs.scrollOffset = 0
	}

	// Draw visible items
	for i := cs.scrollOffset; i < len(items) && i < cs.scrollOffset+maxVisible; i++ {
		item := items[i]
		catName := cs.categoryName(item.Category)
		line := fmt.Sprintf("%-15s %5d %5d %5d", truncate(item.Name, 15), item.Quantity, item.Weight*item.Quantity, item.Volume*item.Quantity)
		ebitenutil.DebugPrintAt(panel, line, padding, startY)
		// Category on same line
		ebitenutil.DebugPrintAt(panel, catName, cs.panelWidth-padding-len(catName)*7, startY)
		startY += lineHeight
	}

	// Scroll indicator
	if len(items) > maxVisible {
		indicator := fmt.Sprintf("[%d/%d]", cs.scrollOffset+1, len(items)-maxVisible+1)
		ebitenutil.DebugPrintAt(panel, indicator, cs.panelWidth-padding-len(indicator)*7, 145)
	}
}

// categoryName returns a short category name.
func (cs *CargoScreen) categoryName(cat vessel.CargoCategory) string {
	names := map[vessel.CargoCategory]string{
		vessel.CargoSupplies: "[SUP]",
		vessel.CargoMedical:  "[MED]",
		vessel.CargoRepair:   "[REP]",
		vessel.CargoTrade:    "[TRD]",
		vessel.CargoSpecial:  "[SPC]",
	}
	if name, ok := names[cat]; ok {
		return name
	}
	return "[???]"
}

// drawBar draws a horizontal capacity bar.
func (cs *CargoScreen) drawBar(img *ebiten.Image, x, y, w, h int, ratio float64, status resources.ThresholdStatus) {
	// Background
	for dx := 0; dx < w; dx++ {
		for dy := 0; dy < h; dy++ {
			img.Set(x+dx, y+dy, cs.skin.BarBackground)
		}
	}

	// Fill
	fillWidth := int(float64(w) * ratio)
	fillColor := cs.skin.BarFill
	switch status {
	case resources.StatusCritical, resources.StatusDepleted:
		fillColor = cs.skin.CriticalColor
	case resources.StatusLow:
		fillColor = cs.skin.WarningColor
	}

	for dx := 0; dx < fillWidth; dx++ {
		for dy := 1; dy < h-1; dy++ {
			img.Set(x+dx, y+dy, fillColor)
		}
	}
}

// ratioToStatus converts a usage ratio to a threshold status.
func (cs *CargoScreen) ratioToStatus(ratio float64) resources.ThresholdStatus {
	if ratio >= 0.95 {
		return resources.StatusCritical
	}
	if ratio >= 0.80 {
		return resources.StatusLow
	}
	return resources.StatusNormal
}

// drawBorder draws a border around the panel.
func (cs *CargoScreen) drawBorder(panel *ebiten.Image) {
	DrawBorder(panel, cs.skin)
}

// truncate shortens a string to max length with ellipsis.
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 3 {
		return s[:max]
	}
	return s[:max-3] + "..."
}

// CargoSummary provides cargo statistics for display.
type CargoSummary struct {
	TotalItems     int
	TotalWeight    int
	TotalVolume    int
	WeightLimit    int
	VolumeLimit    int
	WeightRatio    float64
	VolumeRatio    float64
	Tier           int
	CategoryCounts map[vessel.CargoCategory]int
}

// GetCargoSummary returns a summary of cargo hold contents.
func GetCargoSummary(hold *vessel.CargoHold) CargoSummary {
	summary := CargoSummary{
		TotalItems:     len(hold.Items()),
		TotalWeight:    hold.UsedWeight(),
		TotalVolume:    hold.UsedVolume(),
		WeightLimit:    hold.WeightLimit(),
		VolumeLimit:    hold.VolumeLimit(),
		WeightRatio:    hold.WeightRatio(),
		VolumeRatio:    hold.VolumeRatio(),
		Tier:           hold.Tier(),
		CategoryCounts: make(map[vessel.CargoCategory]int),
	}

	for _, item := range hold.Items() {
		summary.CategoryCounts[item.Category] += item.Quantity
	}

	return summary
}
