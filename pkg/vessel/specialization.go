package vessel

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// SpecializationType identifies upgrade specialization tracks.
type SpecializationType int

const (
	// SpecNone indicates no specialization selected.
	SpecNone SpecializationType = iota
	// SpecSpeed focuses on improving travel speed and efficiency.
	SpecSpeed
	// SpecCargo focuses on improving storage and capacity.
	SpecCargo
	// SpecDefense focuses on improving protection and durability.
	SpecDefense
)

// AllSpecializationTypes returns all specialization types except None.
func AllSpecializationTypes() []SpecializationType {
	return []SpecializationType{SpecSpeed, SpecCargo, SpecDefense}
}

// Specialization represents a module's upgrade specialization.
type Specialization struct {
	Type       SpecializationType
	Level      int     // 0-3 specialization levels
	SpeedBonus float64 // Multiplier for speed-related effects
	CargoBonus float64 // Multiplier for cargo-related effects
	DefBonus   float64 // Multiplier for defense-related effects
}

// NewSpecialization creates a new empty specialization.
func NewSpecialization() *Specialization {
	return &Specialization{
		Type:       SpecNone,
		Level:      0,
		SpeedBonus: 1.0,
		CargoBonus: 1.0,
		DefBonus:   1.0,
	}
}

// SetType sets the specialization type and recalculates bonuses.
func (s *Specialization) SetType(specType SpecializationType) {
	s.Type = specType
	s.recalculateBonuses()
}

// Upgrade increases the specialization level and recalculates bonuses.
func (s *Specialization) Upgrade() bool {
	if s.Type == SpecNone {
		return false
	}
	if s.Level >= 3 {
		return false
	}
	s.Level++
	s.recalculateBonuses()
	return true
}

// recalculateBonuses updates bonuses based on type and level.
func (s *Specialization) recalculateBonuses() {
	// Reset to base
	s.SpeedBonus = 1.0
	s.CargoBonus = 1.0
	s.DefBonus = 1.0

	if s.Type == SpecNone || s.Level == 0 {
		return
	}

	// Each level adds 10% to primary stat, 5% to secondary
	levelBonus := float64(s.Level) * 0.10
	secondaryBonus := float64(s.Level) * 0.05

	switch s.Type {
	case SpecSpeed:
		s.SpeedBonus = 1.0 + levelBonus
		s.CargoBonus = 1.0 - secondaryBonus*0.5 // Small cargo penalty
		s.DefBonus = 1.0 + secondaryBonus*0.5   // Slight defense bonus
	case SpecCargo:
		s.CargoBonus = 1.0 + levelBonus
		s.SpeedBonus = 1.0 - secondaryBonus*0.5 // Small speed penalty
		s.DefBonus = 1.0 + secondaryBonus*0.5   // Slight defense bonus
	case SpecDefense:
		s.DefBonus = 1.0 + levelBonus
		s.SpeedBonus = 1.0 - secondaryBonus*0.5 // Small speed penalty
		s.CargoBonus = 1.0 + secondaryBonus*0.5 // Slight cargo bonus
	}
}

// SpecializedModule extends Module with specialization support.
type SpecializedModule struct {
	*Module
	spec *Specialization
}

// NewSpecializedModule creates a new module with specialization support.
func NewSpecializedModule(moduleType ModuleType) *SpecializedModule {
	return &SpecializedModule{
		Module: NewModule(moduleType),
		spec:   NewSpecialization(),
	}
}

// Specialization returns the module's specialization.
func (sm *SpecializedModule) Specialization() *Specialization {
	return sm.spec
}

// SetSpecialization sets the module's specialization type.
func (sm *SpecializedModule) SetSpecialization(specType SpecializationType) {
	sm.spec.SetType(specType)
}

// UpgradeSpecialization upgrades the specialization level.
func (sm *SpecializedModule) UpgradeSpecialization() bool {
	return sm.spec.Upgrade()
}

// EffectiveSpeedBonus returns the total speed bonus for this module.
func (sm *SpecializedModule) EffectiveSpeedBonus() float64 {
	return sm.spec.SpeedBonus * sm.Efficiency()
}

// EffectiveCargoBonus returns the total cargo bonus for this module.
func (sm *SpecializedModule) EffectiveCargoBonus() float64 {
	return sm.spec.CargoBonus * sm.Efficiency()
}

// EffectiveDefenseBonus returns the total defense bonus for this module.
func (sm *SpecializedModule) EffectiveDefenseBonus() float64 {
	return sm.spec.DefBonus * sm.Efficiency()
}

// SpecializationName returns the genre-appropriate name for a specialization.
func SpecializationName(specType SpecializationType, genre engine.GenreID) string {
	names, ok := specNames[genre]
	if !ok {
		names = specNames[engine.GenreFantasy]
	}
	return names[specType]
}

var specNames = map[engine.GenreID]map[SpecializationType]string{
	engine.GenreFantasy: {
		SpecNone:    "Unspecialized",
		SpecSpeed:   "Swift",
		SpecCargo:   "Laden",
		SpecDefense: "Fortified",
	},
	engine.GenreScifi: {
		SpecNone:    "Standard",
		SpecSpeed:   "Streamlined",
		SpecCargo:   "Expanded",
		SpecDefense: "Reinforced",
	},
	engine.GenreHorror: {
		SpecNone:    "Stock",
		SpecSpeed:   "Tuned",
		SpecCargo:   "Converted",
		SpecDefense: "Armored",
	},
	engine.GenreCyberpunk: {
		SpecNone:    "Factory",
		SpecSpeed:   "Boosted",
		SpecCargo:   "Modded",
		SpecDefense: "Hardened",
	},
	engine.GenrePostapoc: {
		SpecNone:    "Salvaged",
		SpecSpeed:   "Stripped",
		SpecCargo:   "Expanded",
		SpecDefense: "Plated",
	},
}

// SpecializationDescription returns a description of what the specialization does.
func SpecializationDescription(specType SpecializationType) string {
	descriptions := map[SpecializationType]string{
		SpecNone:    "No specialization selected",
		SpecSpeed:   "Increases speed, slightly reduces cargo capacity",
		SpecCargo:   "Increases cargo capacity, slightly reduces speed",
		SpecDefense: "Increases protection, slightly reduces speed",
	}
	if desc, ok := descriptions[specType]; ok {
		return desc
	}
	return "Unknown specialization"
}

// ModuleSpecializationCost returns the currency cost to specialize a module.
func ModuleSpecializationCost(moduleType ModuleType, specLevel int) float64 {
	baseCosts := map[ModuleType]float64{
		ModuleEngine:     30,
		ModuleCargoHold:  25,
		ModuleMedicalBay: 35,
		ModuleNavigation: 30,
		ModuleDefense:    40,
	}
	base := baseCosts[moduleType]
	// Cost increases with level
	return base * float64(specLevel+1)
}

// GetRecommendedSpecialization returns the recommended specialization for a module type.
func GetRecommendedSpecialization(moduleType ModuleType) SpecializationType {
	recommendations := map[ModuleType]SpecializationType{
		ModuleEngine:     SpecSpeed,
		ModuleCargoHold:  SpecCargo,
		ModuleMedicalBay: SpecDefense,
		ModuleNavigation: SpecSpeed,
		ModuleDefense:    SpecDefense,
	}
	if spec, ok := recommendations[moduleType]; ok {
		return spec
	}
	return SpecNone
}

// CanSpecialize checks if a module can be specialized (must be tier 2+).
func CanSpecialize(m *Module) bool {
	return m.Tier() >= 2
}
