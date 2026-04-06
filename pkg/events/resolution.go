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
// Returns an empty ResolutionResult if outcome is nil.
func (r *Resolver) Apply(outcome *EventOutcome, res *resources.Resources, party *crew.Party, v *vessel.Vessel) ResolutionResult {
	if outcome == nil {
		return ResolutionResult{}
	}

	result := ResolutionResult{
		Outcome: *outcome,
		Message: outcome.Description,
	}

	r.applyResourceDeltas(outcome, res)
	result.Deaths = r.applyCrewDamage(outcome, party)
	r.applyVesselDamage(outcome, v)

	return result
}

// applyResourceDeltas applies all resource changes from an outcome.
func (r *Resolver) applyResourceDeltas(outcome *EventOutcome, res *resources.Resources) {
	resourceDeltas := []struct {
		resource resources.ResourceType
		delta    float64
	}{
		{resources.ResourceFood, outcome.FoodDelta},
		{resources.ResourceWater, outcome.WaterDelta},
		{resources.ResourceFuel, outcome.FuelDelta},
		{resources.ResourceMedicine, outcome.MedicineDelta},
		{resources.ResourceMorale, outcome.MoraleDelta},
		{resources.ResourceCurrency, outcome.CurrencyDelta},
	}

	for _, rd := range resourceDeltas {
		if rd.delta != 0 {
			res.Add(rd.resource, rd.delta)
		}
	}
}

// applyCrewDamage applies damage to all crew members if needed.
func (r *Resolver) applyCrewDamage(outcome *EventOutcome, party *crew.Party) []string {
	if outcome.CrewDamage > 0 {
		return party.ApplyDamageToAll(outcome.CrewDamage)
	}
	return nil
}

// applyVesselDamage applies damage to the vessel if needed.
func (r *Resolver) applyVesselDamage(outcome *EventOutcome, v *vessel.Vessel) {
	if outcome.VesselDamage > 0 {
		v.TakeDamage(outcome.VesselDamage)
	}
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

	// Reduce negative effects
	modified.CrewDamage = r.reduceNegativeEffect(modified.CrewDamage, bonus)
	modified.VesselDamage = r.reduceNegativeEffect(modified.VesselDamage, bonus)
	modified.MoraleDelta = r.reduceNegativeDelta(modified.MoraleDelta, bonus)

	// Increase positive effects
	modified.FoodDelta = r.increasePositiveDelta(modified.FoodDelta, bonus)
	modified.WaterDelta = r.increasePositiveDelta(modified.WaterDelta, bonus)
	modified.CurrencyDelta = r.increasePositiveDelta(modified.CurrencyDelta, bonus)

	return modified
}

// reduceNegativeEffect reduces a positive damage value by bonus percentage.
func (r *Resolver) reduceNegativeEffect(value, bonus float64) float64 {
	if value > 0 {
		return value * (1 - bonus)
	}
	return value
}

// reduceNegativeDelta reduces the magnitude of a negative delta.
func (r *Resolver) reduceNegativeDelta(value, bonus float64) float64 {
	if value < 0 {
		return value * (1 - bonus)
	}
	return value
}

// increasePositiveDelta increases a positive delta by bonus percentage.
func (r *Resolver) increasePositiveDelta(value, bonus float64) float64 {
	if value > 0 {
		return value * (1 + bonus)
	}
	return value
}
