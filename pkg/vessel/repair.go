package vessel

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// RepairResult represents the outcome of a repair attempt.
type RepairResult struct {
	Success       bool
	Repaired      float64
	MaterialsUsed int
	Message       string
}

// RepairManager handles vessel repair mechanics.
type RepairManager struct {
	genre engine.GenreID
}

// NewRepairManager creates a new repair manager.
func NewRepairManager(genre engine.GenreID) *RepairManager {
	return &RepairManager{genre: genre}
}

// SetGenre changes the repair vocabulary theme.
func (rm *RepairManager) SetGenre(genre engine.GenreID) {
	rm.genre = genre
}

// RepairWithMaterials attempts to repair the vessel using repair materials.
// Each material unit repairs a base amount of integrity.
func (rm *RepairManager) RepairWithMaterials(v *Vessel, hold *CargoHold, materials int) RepairResult {
	matName := RepairMaterialName(rm.genre)

	if !hold.Has(matName, materials) {
		return RepairResult{
			Success: false,
			Message: "Insufficient repair materials",
		}
	}

	// Each material repairs 5 integrity points
	repairAmount := float64(materials) * 5.0
	needed := v.MaxIntegrity() - v.Integrity()
	if repairAmount > needed {
		repairAmount = needed
		materials = int(needed / 5.0)
		if materials < 1 && needed > 0 {
			materials = 1
		}
	}

	hold.Remove(matName, materials)
	v.Repair(repairAmount)

	return RepairResult{
		Success:       true,
		Repaired:      repairAmount,
		MaterialsUsed: materials,
		Message:       "Repair successful",
	}
}

// RepairAtStation repairs the vessel at a supply station using currency.
// Returns the cost or -1 if cannot afford.
func (rm *RepairManager) RepairAtStation(v *Vessel, currency float64, costPer int) (float64, bool) {
	needed := v.MaxIntegrity() - v.Integrity()
	if needed <= 0 {
		return 0, true
	}

	cost := (needed / 10.0) * float64(costPer)
	if cost > currency {
		return cost, false
	}

	v.Repair(needed)
	return cost, true
}

// EmergencyRepair performs a minimal field repair at no material cost.
// Only available once per journey and restores minimal integrity.
func (rm *RepairManager) EmergencyRepair(v *Vessel, usedThisJourney bool) RepairResult {
	if usedThisJourney {
		return RepairResult{
			Success: false,
			Message: "Emergency repair already used this journey",
		}
	}

	repairAmount := 10.0 // Minimal emergency repair
	v.Repair(repairAmount)

	return RepairResult{
		Success:       true,
		Repaired:      repairAmount,
		MaterialsUsed: 0,
		Message:       "Makeshift repairs hold together... for now",
	}
}

// RepairMaterialName returns the genre-appropriate name for repair materials.
func RepairMaterialName(genre engine.GenreID) string {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Timber",
		engine.GenreScifi:     "Hull Plates",
		engine.GenreHorror:    "Spare Parts",
		engine.GenreCyberpunk: "Tech Parts",
		engine.GenrePostapoc:  "Scrap Metal",
	}
	if name, ok := names[genre]; ok {
		return name
	}
	return "Repair Materials"
}

// GetRepairCost calculates the repair cost for a given breakdown.
func GetRepairCost(b *Breakdown) int {
	return b.RepairCost
}

// CanAffordRepair checks if the party can afford to repair a breakdown.
func CanAffordRepair(hold *CargoHold, b *Breakdown, genre engine.GenreID) bool {
	matName := RepairMaterialName(genre)
	return hold.Has(matName, b.RepairCost)
}
