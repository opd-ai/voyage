package crew

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// StatusType identifies a status effect.
type StatusType int

const (
	// StatusHealthy means no status effects.
	StatusHealthy StatusType = iota
	// StatusDisease spreads between crew and slows recovery.
	StatusDisease
	// StatusInjury reduces action effectiveness.
	StatusInjury
	// StatusExhaustion from overtravel reduces skill effectiveness.
	StatusExhaustion
	// StatusDespair from low morale increases desertion chance.
	StatusDespair
	// StatusGenreAffliction is genre-specific (cursed, irradiated, infected, etc).
	StatusGenreAffliction
)

// StatusEffect represents an active status on a crew member.
type StatusEffect struct {
	Type       StatusType
	Severity   float64 // 0-100, higher is worse
	Duration   int     // Days remaining (-1 for permanent until cured)
	Contagious bool    // Can spread to other crew
}

// AllStatusTypes returns all status types except Healthy.
func AllStatusTypes() []StatusType {
	return []StatusType{
		StatusDisease,
		StatusInjury,
		StatusExhaustion,
		StatusDespair,
		StatusGenreAffliction,
	}
}

// StatusTypeName returns the display name for a status type.
func StatusTypeName(st StatusType, genre engine.GenreID) string {
	if st == StatusHealthy {
		return "Healthy"
	}
	if st == StatusGenreAffliction {
		return genreAfflictionNames[genre]
	}
	return statusTypeNames[st]
}

var statusTypeNames = map[StatusType]string{
	StatusDisease:    "Disease",
	StatusInjury:     "Injured",
	StatusExhaustion: "Exhausted",
	StatusDespair:    "Despairing",
}

var genreAfflictionNames = map[engine.GenreID]string{
	engine.GenreFantasy:   "Cursed",
	engine.GenreScifi:     "Irradiated",
	engine.GenreHorror:    "Infected",
	engine.GenreCyberpunk: "Glitched",
	engine.GenrePostapoc:  "Mutated",
}

// StatusModifiers returns the gameplay modifiers for a status.
type StatusModifiers struct {
	SkillMod       float64 // Modifier to skill effectiveness
	HealMod        float64 // Modifier to healing rate
	DesertionMod   float64 // Modifier to desertion chance
	ContagionRate  float64 // Chance to spread per day
	DamagePerDay   float64 // Health damage per day
	MedicineNeeded float64 // Medicine to cure
}

// GetStatusModifiers returns modifiers for a status effect.
func GetStatusModifiers(st StatusType, severity float64) StatusModifiers {
	base := statusModifiers[st]
	// Scale modifiers by severity (0-100)
	scale := severity / 100.0
	return StatusModifiers{
		SkillMod:       base.SkillMod * scale,
		HealMod:        base.HealMod * scale,
		DesertionMod:   base.DesertionMod * scale,
		ContagionRate:  base.ContagionRate * scale,
		DamagePerDay:   base.DamagePerDay * scale,
		MedicineNeeded: base.MedicineNeeded,
	}
}

var statusModifiers = map[StatusType]StatusModifiers{
	StatusDisease: {
		SkillMod:       -0.3, // 30% skill reduction at max severity
		HealMod:        -0.5, // 50% slower healing
		DesertionMod:   0,
		ContagionRate:  0.2, // 20% spread chance per day
		DamagePerDay:   5,   // 5 damage per day
		MedicineNeeded: 15,
	},
	StatusInjury: {
		SkillMod:       -0.4, // 40% skill reduction
		HealMod:        -0.3, // 30% slower healing
		DesertionMod:   0,
		ContagionRate:  0, // Not contagious
		DamagePerDay:   2, // 2 damage per day without treatment
		MedicineNeeded: 10,
	},
	StatusExhaustion: {
		SkillMod:       -0.25, // 25% skill reduction
		HealMod:        -0.2,
		DesertionMod:   0.1, // 10% more likely to desert
		ContagionRate:  0,
		DamagePerDay:   0, // No direct damage
		MedicineNeeded: 0, // Rest cures, not medicine
	},
	StatusDespair: {
		SkillMod:       -0.15,
		HealMod:        -0.1,
		DesertionMod:   0.3, // 30% more likely to desert
		ContagionRate:  0.1, // Can spread mood
		DamagePerDay:   0,
		MedicineNeeded: 0, // Morale cures, not medicine
	},
	StatusGenreAffliction: {
		SkillMod:       -0.35,
		HealMod:        -0.4,
		DesertionMod:   0.2,
		ContagionRate:  0.15,
		DamagePerDay:   3,
		MedicineNeeded: 20,
	},
}

// StatusTracker manages status effects for a crew member.
type StatusTracker struct {
	effects []StatusEffect
}

// NewStatusTracker creates an empty status tracker.
func NewStatusTracker() *StatusTracker {
	return &StatusTracker{
		effects: make([]StatusEffect, 0),
	}
}

// AddEffect adds a status effect.
func (st *StatusTracker) AddEffect(effect StatusEffect) {
	// Check if already has this effect type
	for i, e := range st.effects {
		if e.Type == effect.Type {
			// Stack severity up to 100
			st.effects[i].Severity += effect.Severity
			if st.effects[i].Severity > 100 {
				st.effects[i].Severity = 100
			}
			// Take longer duration
			if effect.Duration > st.effects[i].Duration {
				st.effects[i].Duration = effect.Duration
			}
			return
		}
	}
	st.effects = append(st.effects, effect)
}

// RemoveEffect removes a status effect by type.
func (st *StatusTracker) RemoveEffect(t StatusType) {
	for i, e := range st.effects {
		if e.Type == t {
			st.effects = append(st.effects[:i], st.effects[i+1:]...)
			return
		}
	}
}

// HasEffect returns true if the member has this status.
func (st *StatusTracker) HasEffect(t StatusType) bool {
	for _, e := range st.effects {
		if e.Type == t {
			return true
		}
	}
	return false
}

// GetEffect returns a copy of the effect if present, nil otherwise.
// Returns a copy rather than a pointer to avoid slice reallocation issues (M-009).
func (st *StatusTracker) GetEffect(t StatusType) *StatusEffect {
	for i := range st.effects {
		if st.effects[i].Type == t {
			// Return a copy to avoid pointer invalidation (M-009)
			eff := st.effects[i]
			return &eff
		}
	}
	return nil
}

// AllEffects returns a copy of all active effects to prevent external mutation (L-009).
func (st *StatusTracker) AllEffects() []StatusEffect {
	return append([]StatusEffect(nil), st.effects...)
}

// IsHealthy returns true if no negative status effects.
func (st *StatusTracker) IsHealthy() bool {
	return len(st.effects) == 0
}

// TotalSkillModifier returns combined skill modifier from all effects.
func (st *StatusTracker) TotalSkillModifier() float64 {
	total := 0.0
	for _, e := range st.effects {
		mods := GetStatusModifiers(e.Type, e.Severity)
		total += mods.SkillMod
	}
	return total
}

// TotalHealModifier returns combined heal modifier from all effects.
func (st *StatusTracker) TotalHealModifier() float64 {
	total := 0.0
	for _, e := range st.effects {
		mods := GetStatusModifiers(e.Type, e.Severity)
		total += mods.HealMod
	}
	return total
}

// TotalDesertionModifier returns combined desertion modifier.
func (st *StatusTracker) TotalDesertionModifier() float64 {
	total := 0.0
	for _, e := range st.effects {
		mods := GetStatusModifiers(e.Type, e.Severity)
		total += mods.DesertionMod
	}
	return total
}

// DailyDamage returns total damage per day from status effects.
func (st *StatusTracker) DailyDamage() float64 {
	total := 0.0
	for _, e := range st.effects {
		mods := GetStatusModifiers(e.Type, e.Severity)
		total += mods.DamagePerDay
	}
	return total
}

// AdvanceDay processes status effects for a new day.
// Returns list of effects that expired.
func (st *StatusTracker) AdvanceDay() []StatusType {
	expired := make([]StatusType, 0)
	remaining := make([]StatusEffect, 0)

	for _, e := range st.effects {
		if e.Duration > 0 {
			e.Duration--
			if e.Duration == 0 {
				expired = append(expired, e.Type)
				continue
			}
		}
		remaining = append(remaining, e)
	}

	st.effects = remaining
	return expired
}

// ContagiousEffects returns effects that can spread to others.
func (st *StatusTracker) ContagiousEffects() []StatusEffect {
	result := make([]StatusEffect, 0)
	for _, e := range st.effects {
		mods := GetStatusModifiers(e.Type, e.Severity)
		if mods.ContagionRate > 0 {
			result = append(result, e)
		}
	}
	return result
}
