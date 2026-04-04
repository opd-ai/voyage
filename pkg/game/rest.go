package game

import (
	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/resources"
)

// RestResult contains the outcome of a rest action.
type RestResult struct {
	MoraleRecovered float64
	HealthRecovered map[int]float64 // crew ID -> health recovered
	FoodConsumed    float64
	WaterConsumed   float64
	TurnsSpent      int
	Message         string
}

// RestManager handles rest mechanics.
type RestManager struct {
	baseMoraleRecovery float64
	baseHealthRecovery float64
	baseFoodCost       float64
	baseWaterCost      float64
}

// NewRestManager creates a new rest manager.
func NewRestManager() *RestManager {
	return &RestManager{
		baseMoraleRecovery: 10.0,
		baseHealthRecovery: 15.0,
		baseFoodCost:       3.0, // Extra food consumed during rest
		baseWaterCost:      2.0, // Extra water consumed during rest
	}
}

// CanRest checks if the party has enough resources to rest.
func (rm *RestManager) CanRest(res *resources.Resources, crewCount int) bool {
	foodNeeded := rm.baseFoodCost * float64(crewCount)
	waterNeeded := rm.baseWaterCost * float64(crewCount)

	return res.Get(resources.ResourceFood) >= foodNeeded &&
		res.Get(resources.ResourceWater) >= waterNeeded
}

// Rest performs a rest action, recovering morale and health.
func (rm *RestManager) Rest(res *resources.Resources, party *crew.Party) RestResult {
	crewCount := party.LivingCount()
	foodCost := rm.baseFoodCost * float64(crewCount)
	waterCost := rm.baseWaterCost * float64(crewCount)

	// Check resources
	if res.Get(resources.ResourceFood) < foodCost {
		return RestResult{
			Message: "Not enough food to rest",
		}
	}
	if res.Get(resources.ResourceWater) < waterCost {
		return RestResult{
			Message: "Not enough water to rest",
		}
	}

	// Consume resources
	res.Consume(resources.ResourceFood, foodCost)
	res.Consume(resources.ResourceWater, waterCost)

	// Recover morale
	moraleRecovered := rm.baseMoraleRecovery
	res.Add(resources.ResourceMorale, moraleRecovered)

	// Heal crew
	healthRecovered := make(map[int]float64)
	for _, member := range party.Living() {
		if member.Health < member.MaxHealth {
			healAmount := rm.baseHealthRecovery
			// Medics heal more
			if member.Skill == crew.SkillMedic {
				healAmount *= 1.5
			}
			member.Heal(healAmount)
			healthRecovered[member.ID] = healAmount
		}
	}

	return RestResult{
		MoraleRecovered: moraleRecovered,
		HealthRecovered: healthRecovered,
		FoodConsumed:    foodCost,
		WaterConsumed:   waterCost,
		TurnsSpent:      1,
		Message:         "The party rested and recovered",
	}
}

// CampRest performs an extended camp rest with better recovery.
func (rm *RestManager) CampRest(res *resources.Resources, party *crew.Party) RestResult {
	result := rm.Rest(res, party)
	if result.Message != "The party rested and recovered" {
		return result
	}

	// Additional recovery from camping
	extraMorale := 5.0
	res.Add(resources.ResourceMorale, extraMorale)
	result.MoraleRecovered += extraMorale
	result.TurnsSpent = 2 // Camping takes longer
	result.Message = "The party set up camp and rested well"

	// Additional healing
	for _, member := range party.Living() {
		if member.Health < member.MaxHealth {
			extraHeal := 5.0
			member.Heal(extraHeal)
			result.HealthRecovered[member.ID] += extraHeal
		}
	}

	return result
}

// GetRestBenefits returns a description of rest benefits.
func GetRestBenefits() string {
	return "Resting recovers morale and heals injured crew members at the cost of extra supplies."
}
