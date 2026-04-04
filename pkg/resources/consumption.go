package resources

// ConsumptionRate holds the daily consumption rate for a resource.
type ConsumptionRate struct {
	BaseRate         float64
	TerrainModifier  map[int]float64 // terrain type -> modifier
	PartyModifier    float64         // per crew member
	SeasonalModifier float64         // seasonal adjustment
}

// DefaultConsumptionRates returns the default consumption rates.
func DefaultConsumptionRates() map[ResourceType]ConsumptionRate {
	return map[ResourceType]ConsumptionRate{
		ResourceFood: {
			BaseRate:      2.0, // per day
			PartyModifier: 0.5, // per crew member
		},
		ResourceWater: {
			BaseRate:      2.5, // per day
			PartyModifier: 0.6, // per crew member
			TerrainModifier: map[int]float64{
				3: 2.0, // desert doubles water consumption
			},
		},
		ResourceFuel: {
			BaseRate:      0, // only consumed on movement
			PartyModifier: 0,
		},
		ResourceMedicine: {
			BaseRate:      0, // only consumed on events
			PartyModifier: 0,
		},
		ResourceMorale: {
			BaseRate:      0.5, // slow decay
			PartyModifier: 0,
		},
		ResourceCurrency: {
			BaseRate:      0, // only spent at trading posts
			PartyModifier: 0,
		},
	}
}

// CalculateDailyConsumption calculates how much of a resource is consumed per day.
func CalculateDailyConsumption(rt ResourceType, crewCount int, terrainType int) float64 {
	rates := DefaultConsumptionRates()
	rate, ok := rates[rt]
	if !ok {
		return 0
	}

	base := rate.BaseRate + (rate.PartyModifier * float64(crewCount))

	// Apply terrain modifier
	if modifier, exists := rate.TerrainModifier[terrainType]; exists {
		base *= modifier
	}

	return base
}

// CalculateMovementCost calculates fuel consumed for movement.
func CalculateMovementCost(terrainCost int, vesselEfficiency float64) float64 {
	baseCost := float64(terrainCost) * 2.0
	return baseCost / vesselEfficiency
}

// RestRecovery calculates morale/health recovery from resting.
func RestRecovery(rt ResourceType, turnsRested int) float64 {
	switch rt {
	case ResourceMorale:
		return float64(turnsRested) * 5.0 // +5 morale per rest turn
	default:
		return 0
	}
}
