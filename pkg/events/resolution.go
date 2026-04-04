package events

import (
	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/resources"
	"github.com/opd-ai/voyage/pkg/vessel"
)

// ResolutionResult contains the complete outcome of event resolution.
type ResolutionResult struct {
	EventID  int
	ChoiceID int
	Outcome  EventOutcome
	Deaths   []string
	Message  string
}

// Resolver handles applying event outcomes to game state.
type Resolver struct{}

// NewResolver creates a new event resolver.
func NewResolver() *Resolver {
	return &Resolver{}
}

// Apply applies an event outcome to the game state.
func (r *Resolver) Apply(outcome *EventOutcome, res *resources.Resources, party *crew.Party, v *vessel.Vessel) ResolutionResult {
	result := ResolutionResult{
		Outcome: *outcome,
		Message: outcome.Description,
	}

	// Apply resource changes
	if outcome.FoodDelta != 0 {
		res.Add(resources.ResourceFood, outcome.FoodDelta)
	}
	if outcome.WaterDelta != 0 {
		res.Add(resources.ResourceWater, outcome.WaterDelta)
	}
	if outcome.FuelDelta != 0 {
		res.Add(resources.ResourceFuel, outcome.FuelDelta)
	}
	if outcome.MedicineDelta != 0 {
		res.Add(resources.ResourceMedicine, outcome.MedicineDelta)
	}
	if outcome.MoraleDelta != 0 {
		res.Add(resources.ResourceMorale, outcome.MoraleDelta)
	}
	if outcome.CurrencyDelta != 0 {
		res.Add(resources.ResourceCurrency, outcome.CurrencyDelta)
	}

	// Apply crew damage
	if outcome.CrewDamage > 0 {
		deaths := party.ApplyDamageToAll(outcome.CrewDamage)
		result.Deaths = deaths
	}

	// Apply vessel damage
	if outcome.VesselDamage > 0 {
		v.TakeDamage(outcome.VesselDamage)
	}

	return result
}

// CanChoose checks if the party can select a choice requiring a skill.
func (r *Resolver) CanChoose(choice *Choice, party *crew.Party) bool {
	if choice.RequireSkill == "" {
		return true
	}
	// Check if party has the required skill
	for _, member := range party.Living() {
		skillName := crew.SkillName(member.Skill, party.Genre())
		if skillName == choice.RequireSkill {
			return true
		}
	}
	return false
}

// GetSkillBonus returns a bonus modifier if the party has a relevant skill.
func (r *Resolver) GetSkillBonus(event *Event, party *crew.Party) float64 {
	// Some events benefit from having certain skills
	skillBonuses := map[EventCategory]crew.Skill{
		CategoryHardship:  crew.SkillMedic,
		CategoryEncounter: crew.SkillTrader,
		CategoryDiscovery: crew.SkillScout,
	}

	if skill, ok := skillBonuses[event.Category]; ok {
		if party.HasSkill(skill) {
			return 0.2 // 20% bonus/reduction
		}
	}
	return 0
}

// ModifyOutcome applies skill bonuses to an outcome.
func (r *Resolver) ModifyOutcome(outcome *EventOutcome, bonus float64) EventOutcome {
	modified := *outcome

	// Reduce negative effects, increase positive effects
	if modified.CrewDamage > 0 {
		modified.CrewDamage *= (1 - bonus)
	}
	if modified.VesselDamage > 0 {
		modified.VesselDamage *= (1 - bonus)
	}
	if modified.MoraleDelta < 0 {
		modified.MoraleDelta *= (1 - bonus)
	}
	if modified.FoodDelta > 0 {
		modified.FoodDelta *= (1 + bonus)
	}
	if modified.WaterDelta > 0 {
		modified.WaterDelta *= (1 + bonus)
	}
	if modified.CurrencyDelta > 0 {
		modified.CurrencyDelta *= (1 + bonus)
	}

	return modified
}
