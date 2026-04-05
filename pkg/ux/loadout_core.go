package ux

import "github.com/opd-ai/voyage/pkg/vessel"

// DefaultStartingPoints is the number of upgrade points available at start.
const DefaultStartingPoints = 5

// LoadoutConfiguration represents a saved loadout configuration.
type LoadoutConfiguration struct {
	EngineTier     int
	CargoTier      int
	MedicalTier    int
	NavigationTier int
	DefenseTier    int
}

// TryUpgradeModule attempts to upgrade the selected module.
// Returns true if the upgrade was successful.
func TryUpgradeModule(moduleSystem *vessel.ModuleSystem, selectedSlot int, pointsRemaining *int) bool {
	if *pointsRemaining <= 0 {
		return false
	}

	moduleTypes := vessel.AllModuleTypes()
	mt := moduleTypes[selectedSlot]
	m := moduleSystem.GetModule(mt)

	if m.Tier() >= 3 { // Max starting tier is 3
		return false
	}

	if moduleSystem.UpgradeModule(mt) {
		*pointsRemaining--
		return true
	}
	return false
}

// TryDowngradeModule attempts to downgrade the selected module.
// Returns true if the downgrade was successful.
func TryDowngradeModule(moduleSystem *vessel.ModuleSystem, selectedSlot int, pointsRemaining *int) bool {
	moduleTypes := vessel.AllModuleTypes()
	mt := moduleTypes[selectedSlot]
	m := moduleSystem.GetModule(mt)

	if m.Tier() <= 1 {
		return false
	}

	m.SetTier(m.Tier() - 1)
	*pointsRemaining++
	return true
}
