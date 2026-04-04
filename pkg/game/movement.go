package game

import (
	"github.com/opd-ai/voyage/pkg/procgen/world"
	"github.com/opd-ai/voyage/pkg/resources"
	"github.com/opd-ai/voyage/pkg/vessel"
)

// MovementResult contains the outcome of a movement action.
type MovementResult struct {
	Success     bool
	FuelUsed    float64
	TimeSpent   int
	TerrainInfo world.TerrainInfo
	Message     string
}

// MovementManager handles vessel movement and resource costs.
type MovementManager struct {
	baseFuelCost float64
	baseTimeCost int
}

// NewMovementManager creates a new movement manager.
func NewMovementManager() *MovementManager {
	return &MovementManager{
		baseFuelCost: 5.0, // Base fuel per move
		baseTimeCost: 1,   // Base turns per move
	}
}

// CalculateMoveCost calculates the cost of moving into terrain.
func (mm *MovementManager) CalculateMoveCost(terrain world.TerrainInfo, v *vessel.Vessel) (float64, int) {
	// Fuel cost = base * terrain modifier / vessel speed
	fuelCost := mm.baseFuelCost * float64(terrain.MovementCost) / v.Speed()

	// Time cost = base * terrain modifier (rounded up)
	timeCost := mm.baseTimeCost * terrain.MovementCost

	return fuelCost, timeCost
}

// CanMove checks if the party can afford to move into terrain.
func (mm *MovementManager) CanMove(terrain world.TerrainInfo, v *vessel.Vessel, res *resources.Resources) bool {
	fuelCost, _ := mm.CalculateMoveCost(terrain, v)
	return res.Get(resources.ResourceFuel) >= fuelCost
}

// Move executes a movement action, consuming resources.
// Returns the result of the movement.
func (mm *MovementManager) Move(terrain world.TerrainInfo, v *vessel.Vessel, res *resources.Resources) MovementResult {
	fuelCost, timeCost := mm.CalculateMoveCost(terrain, v)

	if !res.Consume(resources.ResourceFuel, fuelCost) {
		return MovementResult{
			Success:     false,
			TerrainInfo: terrain,
			Message:     "Not enough fuel to move",
		}
	}

	return MovementResult{
		Success:     true,
		FuelUsed:    fuelCost,
		TimeSpent:   timeCost,
		TerrainInfo: terrain,
		Message:     "Moved successfully",
	}
}

// ApplyTerrainConsumption applies terrain-specific resource consumption.
// Call this after movement to apply additional consumption effects.
func (mm *MovementManager) ApplyTerrainConsumption(terrain world.TerrainInfo, res *resources.Resources, crewCount int) {
	// Base daily consumption per crew
	baseFood := 2.0
	baseWater := 2.0

	// Apply terrain modifiers
	foodCost := baseFood * terrain.FoodModifier * float64(crewCount)
	waterCost := baseWater * terrain.WaterModifier * float64(crewCount)

	res.Consume(resources.ResourceFood, foodCost)
	res.Consume(resources.ResourceWater, waterCost)
}

// GetTerrainMovementDescription returns a description of terrain movement.
func GetTerrainMovementDescription(terrain world.TerrainInfo) string {
	switch terrain.MovementCost {
	case 1:
		return "Easy terrain - quick passage"
	case 2:
		return "Moderate terrain - steady progress"
	case 3:
		return "Difficult terrain - slow going"
	default:
		return "Challenging terrain"
	}
}
