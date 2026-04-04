//go:build !headless

package ux

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/vessel"
)

// LoadoutScreen displays the module loadout configuration before departure.
type LoadoutScreen struct {
	skin            *UISkin
	genre           engine.GenreID
	screenWidth     int
	screenHeight    int
	panelWidth      int
	panelHeight     int
	selectedSlot    int
	visible         bool
	pointsRemaining int
	moduleSystem    *vessel.ModuleSystem
}

// DefaultStartingPoints is the number of upgrade points available at start.
const DefaultStartingPoints = 5

// NewLoadoutScreen creates a new loadout configuration screen.
func NewLoadoutScreen(genre engine.GenreID, screenWidth, screenHeight int) *LoadoutScreen {
	return &LoadoutScreen{
		skin:            DefaultSkin(genre),
		genre:           genre,
		screenWidth:     screenWidth,
		screenHeight:    screenHeight,
		panelWidth:      450,
		panelHeight:     400,
		selectedSlot:    0,
		visible:         false,
		pointsRemaining: DefaultStartingPoints,
		moduleSystem:    vessel.NewModuleSystem(genre),
	}
}

// SetGenre changes the screen's visual theme.
func (ls *LoadoutScreen) SetGenre(genre engine.GenreID) {
	ls.genre = genre
	ls.skin = DefaultSkin(genre)
	ls.moduleSystem.SetGenre(genre)
}

// Show makes the loadout screen visible.
func (ls *LoadoutScreen) Show() {
	ls.visible = true
	ls.selectedSlot = 0
}

// Hide makes the loadout screen hidden.
func (ls *LoadoutScreen) Hide() {
	ls.visible = false
}

// IsVisible returns whether the screen is currently visible.
func (ls *LoadoutScreen) IsVisible() bool {
	return ls.visible
}

// Reset resets the loadout to defaults.
func (ls *LoadoutScreen) Reset() {
	ls.moduleSystem = vessel.NewModuleSystem(ls.genre)
	ls.pointsRemaining = DefaultStartingPoints
	ls.selectedSlot = 0
}

// GetModuleSystem returns the configured module system.
func (ls *LoadoutScreen) GetModuleSystem() *vessel.ModuleSystem {
	return ls.moduleSystem
}

// PointsRemaining returns the remaining upgrade points.
func (ls *LoadoutScreen) PointsRemaining() int {
	return ls.pointsRemaining
}

// SelectNext moves selection to the next module slot.
func (ls *LoadoutScreen) SelectNext() {
	moduleTypes := vessel.AllModuleTypes()
	ls.selectedSlot = (ls.selectedSlot + 1) % len(moduleTypes)
}

// SelectPrev moves selection to the previous module slot.
func (ls *LoadoutScreen) SelectPrev() {
	moduleTypes := vessel.AllModuleTypes()
	ls.selectedSlot--
	if ls.selectedSlot < 0 {
		ls.selectedSlot = len(moduleTypes) - 1
	}
}

// UpgradeSelected attempts to upgrade the selected module.
func (ls *LoadoutScreen) UpgradeSelected() bool {
	if ls.pointsRemaining <= 0 {
		return false
	}

	moduleTypes := vessel.AllModuleTypes()
	mt := moduleTypes[ls.selectedSlot]
	m := ls.moduleSystem.GetModule(mt)

	if m.Tier() >= 3 { // Max starting tier is 3
		return false
	}

	if ls.moduleSystem.UpgradeModule(mt) {
		ls.pointsRemaining--
		return true
	}
	return false
}

// DowngradeSelected attempts to downgrade the selected module.
func (ls *LoadoutScreen) DowngradeSelected() bool {
	moduleTypes := vessel.AllModuleTypes()
	mt := moduleTypes[ls.selectedSlot]
	m := ls.moduleSystem.GetModule(mt)

	if m.Tier() <= 1 {
		return false
	}

	m.SetTier(m.Tier() - 1)
	ls.pointsRemaining++
	return true
}

// Draw renders the loadout screen.
func (ls *LoadoutScreen) Draw(screen *ebiten.Image) {
	if !ls.visible {
		return
	}

	// Draw semi-transparent background overlay
	overlay := ebiten.NewImage(ls.screenWidth, ls.screenHeight)
	overlay.Fill(ls.skin.PanelBackground)
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(0.7)
	screen.DrawImage(overlay, op)

	// Calculate panel position (centered)
	panelX := (ls.screenWidth - ls.panelWidth) / 2
	panelY := (ls.screenHeight - ls.panelHeight) / 2

	// Create panel
	panel := ebiten.NewImage(ls.panelWidth, ls.panelHeight)
	panel.Fill(ls.skin.PanelBackground)
	ls.drawBorder(panel)

	// Draw title
	title := ls.getTitle()
	titleX := (ls.panelWidth - len(title)*7) / 2
	ebitenutil.DebugPrintAt(panel, title, titleX, 12)

	// Draw points remaining
	pointsStr := fmt.Sprintf("Upgrade Points: %d", ls.pointsRemaining)
	ebitenutil.DebugPrintAt(panel, pointsStr, 16, 40)

	// Draw module slots
	ls.drawModuleSlots(panel)

	// Draw instructions
	instructions := "UP/DOWN: select  LEFT/RIGHT: adjust tier  ENTER: confirm  ESC: cancel"
	instrX := (ls.panelWidth - len(instructions)*7) / 2
	ebitenutil.DebugPrintAt(panel, instructions, instrX, ls.panelHeight-20)

	// Draw panel to screen
	opPanel := &ebiten.DrawImageOptions{}
	opPanel.GeoM.Translate(float64(panelX), float64(panelY))
	screen.DrawImage(panel, opPanel)
}

// getTitle returns the genre-appropriate title.
func (ls *LoadoutScreen) getTitle() string {
	titles := map[engine.GenreID]string{
		engine.GenreFantasy:   "Prepare Your Caravan",
		engine.GenreScifi:     "Configure Ship Systems",
		engine.GenreHorror:    "Outfit Your Vehicle",
		engine.GenreCyberpunk: "Customize Your Rig",
		engine.GenrePostapoc:  "Set Up Your Ride",
	}
	if title, ok := titles[ls.genre]; ok {
		return title
	}
	return "Loadout Configuration"
}

// drawModuleSlots draws the list of module configuration slots.
func (ls *LoadoutScreen) drawModuleSlots(panel *ebiten.Image) {
	padding := 16
	startY := 70
	slotHeight := 50

	moduleTypes := vessel.AllModuleTypes()

	for i, mt := range moduleTypes {
		y := startY + i*slotHeight
		selected := i == ls.selectedSlot

		// Selection indicator
		prefix := "  "
		if selected {
			prefix = "> "
		}

		// Module info
		m := ls.moduleSystem.GetModule(mt)
		moduleName := vessel.ModuleTypeName(mt, ls.genre)
		tierName := vessel.TierName(m.Tier(), ls.genre)

		// Draw module name
		line1 := fmt.Sprintf("%s%s", prefix, moduleName)
		ebitenutil.DebugPrintAt(panel, line1, padding, y)

		// Draw tier indicator with arrows
		tierLine := fmt.Sprintf("   Tier: < %s (%d/3) >", tierName, m.Tier())
		ebitenutil.DebugPrintAt(panel, tierLine, padding, y+16)

		// Draw module description
		desc := ls.moduleDescription(mt)
		ebitenutil.DebugPrintAt(panel, "   "+desc, padding, y+32)
	}
}

// moduleDescription returns a short description of what the module affects.
func (ls *LoadoutScreen) moduleDescription(mt vessel.ModuleType) string {
	descriptions := map[vessel.ModuleType]string{
		vessel.ModuleEngine:     "Affects travel speed",
		vessel.ModuleCargoHold:  "Affects cargo capacity",
		vessel.ModuleMedicalBay: "Affects healing rate",
		vessel.ModuleNavigation: "Affects route accuracy",
		vessel.ModuleDefense:    "Affects damage resistance",
	}
	if desc, ok := descriptions[mt]; ok {
		return desc
	}
	return "General improvement"
}

// drawBorder draws a border around the panel.
func (ls *LoadoutScreen) drawBorder(panel *ebiten.Image) {
	w, h := panel.Bounds().Dx(), panel.Bounds().Dy()
	c := ls.skin.PanelBorder

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

// LoadoutConfiguration represents a saved loadout configuration.
type LoadoutConfiguration struct {
	EngineTier     int
	CargoTier      int
	MedicalTier    int
	NavigationTier int
	DefenseTier    int
}

// GetConfiguration returns the current loadout configuration.
func (ls *LoadoutScreen) GetConfiguration() LoadoutConfiguration {
	return LoadoutConfiguration{
		EngineTier:     ls.moduleSystem.GetModule(vessel.ModuleEngine).Tier(),
		CargoTier:      ls.moduleSystem.GetModule(vessel.ModuleCargoHold).Tier(),
		MedicalTier:    ls.moduleSystem.GetModule(vessel.ModuleMedicalBay).Tier(),
		NavigationTier: ls.moduleSystem.GetModule(vessel.ModuleNavigation).Tier(),
		DefenseTier:    ls.moduleSystem.GetModule(vessel.ModuleDefense).Tier(),
	}
}

// ApplyConfiguration applies a saved configuration to the screen.
func (ls *LoadoutScreen) ApplyConfiguration(cfg LoadoutConfiguration) {
	ls.Reset()

	// Apply each tier, deducting points as needed
	tiers := map[vessel.ModuleType]int{
		vessel.ModuleEngine:     cfg.EngineTier,
		vessel.ModuleCargoHold:  cfg.CargoTier,
		vessel.ModuleMedicalBay: cfg.MedicalTier,
		vessel.ModuleNavigation: cfg.NavigationTier,
		vessel.ModuleDefense:    cfg.DefenseTier,
	}

	for mt, targetTier := range tiers {
		m := ls.moduleSystem.GetModule(mt)
		for m.Tier() < targetTier && ls.pointsRemaining > 0 {
			if ls.moduleSystem.UpgradeModule(mt) {
				ls.pointsRemaining--
			} else {
				break
			}
		}
	}
}
