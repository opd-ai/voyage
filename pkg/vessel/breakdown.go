package vessel

import (
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// BreakdownType identifies the type of vessel malfunction.
type BreakdownType int

const (
	// BreakdownMinor is easily repaired with minimal resources.
	BreakdownMinor BreakdownType = iota
	// BreakdownModerate requires significant repair effort.
	BreakdownModerate
	// BreakdownSevere requires major repairs and may strand the party.
	BreakdownSevere
)

// AllBreakdownTypes returns all breakdown types.
func AllBreakdownTypes() []BreakdownType {
	return []BreakdownType{BreakdownMinor, BreakdownModerate, BreakdownSevere}
}

// Breakdown represents a vessel malfunction event.
type Breakdown struct {
	Type        BreakdownType
	Description string
	DamageDealt float64
	RepairCost  int
	DaysLost    int
}

// BreakdownChecker manages breakdown event generation.
type BreakdownChecker struct {
	gen *seed.Generator
}

// NewBreakdownChecker creates a breakdown checker with the given seed.
func NewBreakdownChecker(masterSeed int64) *BreakdownChecker {
	return &BreakdownChecker{
		gen: seed.NewGenerator(masterSeed, "breakdown"),
	}
}

// Check determines if a breakdown occurs based on vessel condition.
// Returns nil if no breakdown occurs.
func (bc *BreakdownChecker) Check(v *Vessel, turn int) *Breakdown {
	// Base breakdown chance increases with damage
	baseChance := 0.02 // 2% base chance per turn
	damageModifier := 1.0 - v.IntegrityRatio()
	// At 50% integrity, chance doubles. At 25%, triples.
	chance := baseChance * (1.0 + 2.0*damageModifier)

	if bc.gen.Float64() > chance {
		return nil
	}

	// Determine breakdown severity
	// More damaged vessels are more likely to have severe breakdowns
	severityRoll := bc.gen.Float64() + damageModifier*0.3
	var bType BreakdownType
	switch {
	case severityRoll < 0.6:
		bType = BreakdownMinor
	case severityRoll < 0.9:
		bType = BreakdownModerate
	default:
		bType = BreakdownSevere
	}

	return bc.generateBreakdown(bType, v)
}

// generateBreakdown creates a breakdown event of the given type.
func (bc *BreakdownChecker) generateBreakdown(bType BreakdownType, v *Vessel) *Breakdown {
	var damage float64
	var cost int
	var days int
	var desc string

	switch bType {
	case BreakdownMinor:
		damage = 5 + float64(bc.gen.Intn(10))
		cost = 1 + bc.gen.Intn(3)
		days = 0
		desc = bc.minorBreakdownDesc(v)
	case BreakdownModerate:
		damage = 15 + float64(bc.gen.Intn(15))
		cost = 3 + bc.gen.Intn(5)
		days = 1
		desc = bc.moderateBreakdownDesc(v)
	case BreakdownSevere:
		damage = 30 + float64(bc.gen.Intn(20))
		cost = 8 + bc.gen.Intn(7)
		days = 2 + bc.gen.Intn(2)
		desc = bc.severeBreakdownDesc(v)
	}

	return &Breakdown{
		Type:        bType,
		Description: desc,
		DamageDealt: damage,
		RepairCost:  cost,
		DaysLost:    days,
	}
}

func (bc *BreakdownChecker) minorBreakdownDesc(v *Vessel) string {
	descs := []string{
		"A wheel came loose",
		"The axle is squeaking",
		"A minor leak was found",
		"The hitches need adjustment",
		"Some cargo shifted dangerously",
	}
	return seed.Choice(bc.gen, descs)
}

func (bc *BreakdownChecker) moderateBreakdownDesc(v *Vessel) string {
	descs := []string{
		"The main axle cracked",
		"A wheel shattered",
		"The frame is warped",
		"Critical components are damaged",
		"The hull is breached",
	}
	return seed.Choice(bc.gen, descs)
}

func (bc *BreakdownChecker) severeBreakdownDesc(v *Vessel) string {
	descs := []string{
		"Catastrophic structural failure",
		"The engine exploded",
		"Total drive system failure",
		"Multiple critical systems failed",
		"The vessel nearly collapsed",
	}
	return seed.Choice(bc.gen, descs)
}

// BreakdownTypeName returns a human-readable name for the breakdown type.
func BreakdownTypeName(bt BreakdownType) string {
	switch bt {
	case BreakdownMinor:
		return "Minor"
	case BreakdownModerate:
		return "Moderate"
	case BreakdownSevere:
		return "Severe"
	default:
		return "Unknown"
	}
}
