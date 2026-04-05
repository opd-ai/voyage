//go:build headless

package ux

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/vessel"
)

// LoadoutScreen displays the module loadout configuration before departure.
// This is a headless stub for testing.
type LoadoutScreen struct {
	genre           engine.GenreID
	screenWidth     int
	screenHeight    int
	selectedSlot    int
	visible         bool
	pointsRemaining int
	moduleSystem    *vessel.ModuleSystem
}

// NewLoadoutScreen creates a new loadout configuration screen (headless stub).
func NewLoadoutScreen(genre engine.GenreID, screenWidth, screenHeight int) *LoadoutScreen {
	return &LoadoutScreen{
		genre:           genre,
		screenWidth:     screenWidth,
		screenHeight:    screenHeight,
		selectedSlot:    0,
		visible:         false,
		pointsRemaining: DefaultStartingPoints,
		moduleSystem:    vessel.NewModuleSystem(genre),
	}
}

// SetGenre changes the screen's visual theme.
func (ls *LoadoutScreen) SetGenre(genre engine.GenreID) {
	ls.genre = genre
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
	return TryUpgradeModule(ls.moduleSystem, ls.selectedSlot, &ls.pointsRemaining)
}

// DowngradeSelected attempts to downgrade the selected module.
func (ls *LoadoutScreen) DowngradeSelected() bool {
	return TryDowngradeModule(ls.moduleSystem, ls.selectedSlot, &ls.pointsRemaining)
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
