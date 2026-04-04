package vessel

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// ModuleType identifies the type of vessel module.
type ModuleType int

const (
	// ModuleEngine provides propulsion and affects vessel speed.
	ModuleEngine ModuleType = iota
	// ModuleCargoHold provides storage capacity for supplies and goods.
	ModuleCargoHold
	// ModuleMedicalBay provides healing capabilities for crew.
	ModuleMedicalBay
	// ModuleNavigation provides routing assistance and hazard detection.
	ModuleNavigation
	// ModuleDefense provides combat and hazard protection.
	ModuleDefense
)

// AllModuleTypes returns all module types.
func AllModuleTypes() []ModuleType {
	return []ModuleType{
		ModuleEngine,
		ModuleCargoHold,
		ModuleMedicalBay,
		ModuleNavigation,
		ModuleDefense,
	}
}

// Module represents a vessel subsystem with its own integrity.
type Module struct {
	moduleType   ModuleType
	tier         int
	integrity    float64
	maxIntegrity float64
	efficiency   float64 // Effectiveness based on integrity [0, 1]
}

// DefaultModuleStats returns base stats for each module type.
var DefaultModuleStats = map[ModuleType]struct {
	MaxIntegrity float64
	BaseTier     int
}{
	ModuleEngine:     {MaxIntegrity: 50, BaseTier: 1},
	ModuleCargoHold:  {MaxIntegrity: 40, BaseTier: 1},
	ModuleMedicalBay: {MaxIntegrity: 30, BaseTier: 1},
	ModuleNavigation: {MaxIntegrity: 30, BaseTier: 1},
	ModuleDefense:    {MaxIntegrity: 60, BaseTier: 1},
}

// NewModule creates a new module of the given type.
func NewModule(moduleType ModuleType) *Module {
	stats := DefaultModuleStats[moduleType]
	return &Module{
		moduleType:   moduleType,
		tier:         stats.BaseTier,
		integrity:    stats.MaxIntegrity,
		maxIntegrity: stats.MaxIntegrity,
		efficiency:   1.0,
	}
}

// NewModuleWithTier creates a new module with a specific tier.
func NewModuleWithTier(moduleType ModuleType, tier int) *Module {
	m := NewModule(moduleType)
	m.SetTier(tier)
	return m
}

// Type returns the module type.
func (m *Module) Type() ModuleType {
	return m.moduleType
}

// Tier returns the current module tier.
func (m *Module) Tier() int {
	return m.tier
}

// SetTier changes the module tier, adjusting max integrity.
func (m *Module) SetTier(tier int) {
	if tier < 1 {
		tier = 1
	}
	if tier > 5 {
		tier = 5
	}
	m.tier = tier
	baseStats := DefaultModuleStats[m.moduleType]
	m.maxIntegrity = baseStats.MaxIntegrity * (1.0 + float64(tier-1)*0.25)
	if m.integrity > m.maxIntegrity {
		m.integrity = m.maxIntegrity
	}
	m.updateEfficiency()
}

// Integrity returns the current module integrity.
func (m *Module) Integrity() float64 {
	return m.integrity
}

// MaxIntegrity returns the maximum module integrity.
func (m *Module) MaxIntegrity() float64 {
	return m.maxIntegrity
}

// IntegrityRatio returns integrity as a ratio [0, 1].
func (m *Module) IntegrityRatio() float64 {
	if m.maxIntegrity <= 0 {
		return 0
	}
	return m.integrity / m.maxIntegrity
}

// Efficiency returns the module effectiveness based on integrity.
func (m *Module) Efficiency() float64 {
	return m.efficiency
}

// TakeDamage reduces module integrity by the given amount.
// Returns true if the module is disabled (0 integrity).
func (m *Module) TakeDamage(amount float64) bool {
	m.integrity -= amount
	if m.integrity < 0 {
		m.integrity = 0
	}
	m.updateEfficiency()
	return m.integrity <= 0
}

// Repair increases module integrity by the given amount.
func (m *Module) Repair(amount float64) {
	m.integrity += amount
	if m.integrity > m.maxIntegrity {
		m.integrity = m.maxIntegrity
	}
	m.updateEfficiency()
}

// RepairFull restores module to full integrity.
func (m *Module) RepairFull() {
	m.integrity = m.maxIntegrity
	m.efficiency = 1.0
}

// IsDisabled returns true if the module has zero integrity.
func (m *Module) IsDisabled() bool {
	return m.integrity <= 0
}

// IsCritical returns true if module integrity is below 25%.
func (m *Module) IsCritical() bool {
	return m.IntegrityRatio() < 0.25
}

// IsDamaged returns true if module integrity is below 75%.
func (m *Module) IsDamaged() bool {
	return m.IntegrityRatio() < 0.75
}

// updateEfficiency recalculates efficiency based on integrity ratio.
func (m *Module) updateEfficiency() {
	ratio := m.IntegrityRatio()
	// Efficiency drops more steeply as integrity falls
	// At 50% integrity, efficiency is 60%
	// At 25% integrity, efficiency is 30%
	// At 0% integrity, efficiency is 0%
	m.efficiency = ratio * (0.4 + 0.6*ratio)
}

// ModuleSystem manages the modular vessel system.
type ModuleSystem struct {
	genre   engine.GenreID
	modules map[ModuleType]*Module
}

// NewModuleSystem creates a new module system with default modules.
func NewModuleSystem(genre engine.GenreID) *ModuleSystem {
	ms := &ModuleSystem{
		genre:   genre,
		modules: make(map[ModuleType]*Module),
	}
	// Initialize all modules with default configuration
	for _, mt := range AllModuleTypes() {
		ms.modules[mt] = NewModule(mt)
	}
	return ms
}

// SetGenre changes the module vocabulary theme.
func (ms *ModuleSystem) SetGenre(genre engine.GenreID) {
	ms.genre = genre
}

// Genre returns the current genre.
func (ms *ModuleSystem) Genre() engine.GenreID {
	return ms.genre
}

// GetModule returns the module of the given type.
func (ms *ModuleSystem) GetModule(moduleType ModuleType) *Module {
	return ms.modules[moduleType]
}

// GetAllModules returns all modules.
func (ms *ModuleSystem) GetAllModules() []*Module {
	result := make([]*Module, 0, len(ms.modules))
	for _, mt := range AllModuleTypes() {
		if m := ms.modules[mt]; m != nil {
			result = append(result, m)
		}
	}
	return result
}

// EngineEfficiency returns the engine module's efficiency.
func (ms *ModuleSystem) EngineEfficiency() float64 {
	if m := ms.modules[ModuleEngine]; m != nil {
		return m.Efficiency()
	}
	return 0
}

// CargoCapacityMultiplier returns a multiplier based on cargo hold.
func (ms *ModuleSystem) CargoCapacityMultiplier() float64 {
	if m := ms.modules[ModuleCargoHold]; m != nil {
		return 0.5 + 0.5*m.Efficiency() + float64(m.Tier()-1)*0.15
	}
	return 0.5
}

// MedicalEfficiency returns the medical bay's efficiency.
func (ms *ModuleSystem) MedicalEfficiency() float64 {
	if m := ms.modules[ModuleMedicalBay]; m != nil {
		return m.Efficiency()
	}
	return 0
}

// NavigationAccuracy returns the navigation module's accuracy.
func (ms *ModuleSystem) NavigationAccuracy() float64 {
	if m := ms.modules[ModuleNavigation]; m != nil {
		return m.Efficiency()
	}
	return 0
}

// DefenseRating returns the defense module's effectiveness.
func (ms *ModuleSystem) DefenseRating() float64 {
	if m := ms.modules[ModuleDefense]; m != nil {
		return m.Efficiency() * float64(m.Tier()) * 0.1
	}
	return 0
}

// TotalIntegrity returns the sum of all module integrities.
func (ms *ModuleSystem) TotalIntegrity() float64 {
	total := 0.0
	for _, m := range ms.modules {
		total += m.Integrity()
	}
	return total
}

// TotalMaxIntegrity returns the sum of all module max integrities.
func (ms *ModuleSystem) TotalMaxIntegrity() float64 {
	total := 0.0
	for _, m := range ms.modules {
		total += m.MaxIntegrity()
	}
	return total
}

// OverallIntegrityRatio returns the overall system health.
func (ms *ModuleSystem) OverallIntegrityRatio() float64 {
	maxTotal := ms.TotalMaxIntegrity()
	if maxTotal <= 0 {
		return 0
	}
	return ms.TotalIntegrity() / maxTotal
}

// DistributeDamage applies damage across modules randomly.
func (ms *ModuleSystem) DistributeDamage(amount float64, rng func(n int) int) {
	if amount <= 0 {
		return
	}
	moduleList := ms.GetAllModules()
	if len(moduleList) == 0 {
		return
	}
	// Pick a random module to take the damage
	idx := rng(len(moduleList))
	moduleList[idx].TakeDamage(amount)
}

// RepairModule repairs a specific module type by the given amount.
func (ms *ModuleSystem) RepairModule(moduleType ModuleType, amount float64) bool {
	if m := ms.modules[moduleType]; m != nil {
		m.Repair(amount)
		return true
	}
	return false
}

// UpgradeModule increases the tier of a specific module.
func (ms *ModuleSystem) UpgradeModule(moduleType ModuleType) bool {
	if m := ms.modules[moduleType]; m != nil {
		if m.Tier() >= 5 {
			return false
		}
		m.SetTier(m.Tier() + 1)
		return true
	}
	return false
}

// ModuleTypeName returns the genre-appropriate name for a module type.
func ModuleTypeName(mt ModuleType, genre engine.GenreID) string {
	names, ok := moduleNames[genre]
	if !ok {
		names = moduleNames[engine.GenreFantasy]
	}
	return names[mt]
}

var moduleNames = map[engine.GenreID]map[ModuleType]string{
	engine.GenreFantasy: {
		ModuleEngine:     "Stable",
		ModuleCargoHold:  "Wagon Bed",
		ModuleMedicalBay: "Healer's Cart",
		ModuleNavigation: "Scout's Perch",
		ModuleDefense:    "Guard Platform",
	},
	engine.GenreScifi: {
		ModuleEngine:     "Engine Room",
		ModuleCargoHold:  "Cargo Bay",
		ModuleMedicalBay: "Sickbay",
		ModuleNavigation: "Bridge",
		ModuleDefense:    "Weapons Array",
	},
	engine.GenreHorror: {
		ModuleEngine:     "Engine Bay",
		ModuleCargoHold:  "Trunk Space",
		ModuleMedicalBay: "First Aid Station",
		ModuleNavigation: "GPS System",
		ModuleDefense:    "Armored Plates",
	},
	engine.GenreCyberpunk: {
		ModuleEngine:     "Core Systems",
		ModuleCargoHold:  "Smuggler's Hold",
		ModuleMedicalBay: "Trauma Bay",
		ModuleNavigation: "Netlink Array",
		ModuleDefense:    "ICE Suite",
	},
	engine.GenrePostapoc: {
		ModuleEngine:     "Reactor",
		ModuleCargoHold:  "Scrap Hold",
		ModuleMedicalBay: "Chem Station",
		ModuleNavigation: "Watchtower",
		ModuleDefense:    "Gun Mount",
	},
}

// ModuleName returns the name for a module using its system's genre.
func (m *Module) Name(genre engine.GenreID) string {
	return ModuleTypeName(m.moduleType, genre)
}

// ModuleCondition describes the state of a module.
type ModuleCondition int

const (
	// ModuleConditionPristine means the module is at full integrity.
	ModuleConditionPristine ModuleCondition = iota
	// ModuleConditionOperational means the module is above 75% integrity.
	ModuleConditionOperational
	// ModuleConditionDamaged means the module is between 50% and 75%.
	ModuleConditionDamaged
	// ModuleConditionCritical means the module is between 25% and 50%.
	ModuleConditionCritical
	// ModuleConditionDisabled means the module is below 25% or at 0.
	ModuleConditionDisabled
)

// GetModuleCondition returns the condition status of a module.
func GetModuleCondition(m *Module) ModuleCondition {
	ratio := m.IntegrityRatio()
	switch {
	case ratio >= 1.0:
		return ModuleConditionPristine
	case ratio >= 0.75:
		return ModuleConditionOperational
	case ratio >= 0.5:
		return ModuleConditionDamaged
	case ratio >= 0.25:
		return ModuleConditionCritical
	default:
		return ModuleConditionDisabled
	}
}

// ModuleConditionName returns a human-readable condition name.
func ModuleConditionName(mc ModuleCondition) string {
	switch mc {
	case ModuleConditionPristine:
		return "Pristine"
	case ModuleConditionOperational:
		return "Operational"
	case ModuleConditionDamaged:
		return "Damaged"
	case ModuleConditionCritical:
		return "Critical"
	case ModuleConditionDisabled:
		return "Disabled"
	default:
		return "Unknown"
	}
}
