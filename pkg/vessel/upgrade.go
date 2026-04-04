package vessel

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// UpgradeResult represents the outcome of an upgrade attempt.
type UpgradeResult struct {
	Success      bool
	NewTier      int
	CurrencyUsed float64
	Message      string
}

// UpgradeManager handles module upgrades at supply points.
type UpgradeManager struct {
	genre engine.GenreID
}

// NewUpgradeManager creates a new upgrade manager.
func NewUpgradeManager(genre engine.GenreID) *UpgradeManager {
	return &UpgradeManager{genre: genre}
}

// SetGenre changes the upgrade vocabulary.
func (um *UpgradeManager) SetGenre(genre engine.GenreID) {
	um.genre = genre
}

// UpgradeCost calculates the cost to upgrade a module to the next tier.
func UpgradeCost(moduleType ModuleType, currentTier int) float64 {
	if currentTier >= 5 {
		return 0 // Already maxed
	}
	baseCosts := map[ModuleType]float64{
		ModuleEngine:     50,
		ModuleCargoHold:  40,
		ModuleMedicalBay: 45,
		ModuleNavigation: 35,
		ModuleDefense:    55,
	}
	base := baseCosts[moduleType]
	// Cost increases by 50% per tier
	return base * (1.0 + float64(currentTier)*0.5)
}

// CanAffordUpgrade checks if the party can afford to upgrade a module.
func (um *UpgradeManager) CanAffordUpgrade(ms *ModuleSystem, moduleType ModuleType, currency float64) bool {
	m := ms.GetModule(moduleType)
	if m == nil || m.Tier() >= 5 {
		return false
	}
	return currency >= UpgradeCost(moduleType, m.Tier())
}

// AttemptUpgrade tries to upgrade a module, spending currency.
// Returns the result and the amount of currency used.
func (um *UpgradeManager) AttemptUpgrade(ms *ModuleSystem, moduleType ModuleType, currency float64) UpgradeResult {
	m := ms.GetModule(moduleType)
	if m == nil {
		return UpgradeResult{
			Success: false,
			Message: "Module not found",
		}
	}

	if m.Tier() >= 5 {
		return UpgradeResult{
			Success: false,
			NewTier: m.Tier(),
			Message: "Module already at maximum tier",
		}
	}

	cost := UpgradeCost(moduleType, m.Tier())
	if currency < cost {
		return UpgradeResult{
			Success: false,
			NewTier: m.Tier(),
			Message: "Insufficient currency for upgrade",
		}
	}

	// Perform the upgrade
	ms.UpgradeModule(moduleType)

	return UpgradeResult{
		Success:      true,
		NewTier:      m.Tier(),
		CurrencyUsed: cost,
		Message:      um.upgradeMessage(moduleType, m.Tier()),
	}
}

// upgradeMessage generates a genre-appropriate upgrade message.
func (um *UpgradeManager) upgradeMessage(moduleType ModuleType, newTier int) string {
	moduleName := ModuleTypeName(moduleType, um.genre)
	tierName := TierName(newTier, um.genre)
	return moduleName + " upgraded to " + tierName
}

// TierName returns the genre-appropriate name for a tier level.
func TierName(tier int, genre engine.GenreID) string {
	names, ok := tierNames[genre]
	if !ok {
		names = tierNames[engine.GenreFantasy]
	}
	if tier < 1 || tier > 5 {
		return "Unknown"
	}
	return names[tier-1]
}

var tierNames = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"Basic",
		"Improved",
		"Superior",
		"Exceptional",
		"Legendary",
	},
	engine.GenreScifi: {
		"Standard",
		"Enhanced",
		"Advanced",
		"Military-Grade",
		"Prototype",
	},
	engine.GenreHorror: {
		"Makeshift",
		"Functional",
		"Reliable",
		"Hardened",
		"Fortified",
	},
	engine.GenreCyberpunk: {
		"Stock",
		"Modded",
		"Custom",
		"Black Market",
		"Mil-Spec",
	},
	engine.GenrePostapoc: {
		"Salvaged",
		"Patched",
		"Restored",
		"Reinforced",
		"Pre-War",
	},
}

// GetAvailableUpgrades returns a list of all possible upgrades with costs.
func (um *UpgradeManager) GetAvailableUpgrades(ms *ModuleSystem) []UpgradeOption {
	var options []UpgradeOption
	for _, mt := range AllModuleTypes() {
		m := ms.GetModule(mt)
		if m == nil || m.Tier() >= 5 {
			continue
		}
		options = append(options, UpgradeOption{
			ModuleType:   mt,
			CurrentTier:  m.Tier(),
			NextTier:     m.Tier() + 1,
			Cost:         UpgradeCost(mt, m.Tier()),
			ModuleName:   ModuleTypeName(mt, um.genre),
			NextTierName: TierName(m.Tier()+1, um.genre),
		})
	}
	return options
}

// UpgradeOption represents a possible upgrade at a supply point.
type UpgradeOption struct {
	ModuleType   ModuleType
	CurrentTier  int
	NextTier     int
	Cost         float64
	ModuleName   string
	NextTierName string
}

// UpgradeDescription returns a description of what the upgrade provides.
func (uo *UpgradeOption) UpgradeDescription() string {
	switch uo.ModuleType {
	case ModuleEngine:
		return "Increases vessel speed and fuel efficiency"
	case ModuleCargoHold:
		return "Increases cargo capacity"
	case ModuleMedicalBay:
		return "Improves healing effectiveness"
	case ModuleNavigation:
		return "Improves route accuracy and hazard detection"
	case ModuleDefense:
		return "Increases protection against attacks"
	default:
		return "Improves module performance"
	}
}

// BulkUpgrade attempts to upgrade multiple modules at once.
func (um *UpgradeManager) BulkUpgrade(ms *ModuleSystem, modules []ModuleType, currency float64) ([]UpgradeResult, float64) {
	var results []UpgradeResult
	totalSpent := 0.0

	for _, mt := range modules {
		remaining := currency - totalSpent
		result := um.AttemptUpgrade(ms, mt, remaining)
		results = append(results, result)
		if result.Success {
			totalSpent += result.CurrencyUsed
		}
	}

	return results, totalSpent
}

// AllTierNames returns all tier names for a genre.
func AllTierNames(genre engine.GenreID) []string {
	names, ok := tierNames[genre]
	if !ok {
		return tierNames[engine.GenreFantasy]
	}
	return names
}
