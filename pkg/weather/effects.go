package weather

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// WeatherEffect contains the gameplay effects of a weather condition.
type WeatherEffect struct {
	MovementCostMultiplier float64 // Multiplier on movement fuel cost (1.0 = normal)
	WaterConsumptionMod    float64 // Modifier to daily water consumption
	FuelConsumptionMod     float64 // Modifier to daily fuel consumption
	VisibilityPenalty      float64 // 0-1 where 1 is totally blind
	CrewHealthDamage       float64 // Per-turn damage to crew health
	VesselDamagePerTurn    float64 // Per-turn damage to vessel
	MoraleModifier         float64 // Per-turn morale change
	EncounterChanceMod     float64 // Modifier to random encounter chance
}

// WeatherEffects maps weather types to their gameplay effects.
var WeatherEffects = map[WeatherType]WeatherEffect{
	WeatherClear: {
		MovementCostMultiplier: 1.0,
		WaterConsumptionMod:    0,
		FuelConsumptionMod:     0,
		VisibilityPenalty:      0,
		CrewHealthDamage:       0,
		VesselDamagePerTurn:    0,
		MoraleModifier:         0.5,
		EncounterChanceMod:     0,
	},
	WeatherStorm: {
		MovementCostMultiplier: 1.5,
		WaterConsumptionMod:    0,
		FuelConsumptionMod:     0.2,
		VisibilityPenalty:      0.3,
		CrewHealthDamage:       2,
		VesselDamagePerTurn:    3,
		MoraleModifier:         -2,
		EncounterChanceMod:     -0.3,
	},
	WeatherBlizzard: {
		MovementCostMultiplier: 2.0,
		WaterConsumptionMod:    -0.2,
		FuelConsumptionMod:     0.5,
		VisibilityPenalty:      0.6,
		CrewHealthDamage:       5,
		VesselDamagePerTurn:    2,
		MoraleModifier:         -3,
		EncounterChanceMod:     -0.5,
	},
	WeatherHeatwave: {
		MovementCostMultiplier: 1.2,
		WaterConsumptionMod:    0.5,
		FuelConsumptionMod:     0.1,
		VisibilityPenalty:      0.1,
		CrewHealthDamage:       3,
		VesselDamagePerTurn:    1,
		MoraleModifier:         -2,
		EncounterChanceMod:     0,
	},
	WeatherFlood: {
		MovementCostMultiplier: 1.8,
		WaterConsumptionMod:    -0.3,
		FuelConsumptionMod:     0.4,
		VisibilityPenalty:      0.2,
		CrewHealthDamage:       2,
		VesselDamagePerTurn:    5,
		MoraleModifier:         -2,
		EncounterChanceMod:     -0.4,
	},
	WeatherFog: {
		MovementCostMultiplier: 1.3,
		WaterConsumptionMod:    0,
		FuelConsumptionMod:     0.1,
		VisibilityPenalty:      0.7,
		CrewHealthDamage:       0,
		VesselDamagePerTurn:    0,
		MoraleModifier:         -1,
		EncounterChanceMod:     0.3,
	},
	WeatherMeteorShower: {
		MovementCostMultiplier: 1.4,
		WaterConsumptionMod:    0,
		FuelConsumptionMod:     0.2,
		VisibilityPenalty:      0.2,
		CrewHealthDamage:       4,
		VesselDamagePerTurn:    8,
		MoraleModifier:         -4,
		EncounterChanceMod:     -0.2,
	},
	WeatherDustStorm: {
		MovementCostMultiplier: 1.6,
		WaterConsumptionMod:    0.2,
		FuelConsumptionMod:     0.3,
		VisibilityPenalty:      0.8,
		CrewHealthDamage:       3,
		VesselDamagePerTurn:    4,
		MoraleModifier:         -3,
		EncounterChanceMod:     -0.4,
	},
	WeatherAcidRain: {
		MovementCostMultiplier: 1.3,
		WaterConsumptionMod:    0.1,
		FuelConsumptionMod:     0.1,
		VisibilityPenalty:      0.3,
		CrewHealthDamage:       5,
		VesselDamagePerTurn:    6,
		MoraleModifier:         -3,
		EncounterChanceMod:     -0.2,
	},
}

// GetWeatherEffect returns the effect for a weather type.
func GetWeatherEffect(w WeatherType) WeatherEffect {
	if effect, ok := WeatherEffects[w]; ok {
		return effect
	}
	return WeatherEffects[WeatherClear]
}

// HazardEffect contains the gameplay effects of a terrain hazard.
type HazardEffect struct {
	InjuryChance       float64 // Chance of crew injury when traversing
	InjuryDamage       float64 // Damage dealt if injury occurs
	FuelCostMultiplier float64 // Multiplier on fuel cost (1.0 = normal)
	WaterCostIncrease  float64 // Additional water cost
	LootChance         float64 // Chance of finding loot
	DangerChance       float64 // Chance of triggering a hostile encounter
	MovementPenalty    float64 // Additional movement time cost
}

// HazardEffects maps terrain hazards to their effects.
var HazardEffects = map[TerrainHazard]HazardEffect{
	HazardNone: {
		InjuryChance:       0,
		InjuryDamage:       0,
		FuelCostMultiplier: 1.0,
		WaterCostIncrease:  0,
		LootChance:         0,
		DangerChance:       0,
		MovementPenalty:    0,
	},
	HazardMountainPass: {
		InjuryChance:       0.3,
		InjuryDamage:       15,
		FuelCostMultiplier: 1.5,
		WaterCostIncrease:  0,
		LootChance:         0.1,
		DangerChance:       0.2,
		MovementPenalty:    2,
	},
	HazardRiverCrossing: {
		InjuryChance:       0.1,
		InjuryDamage:       10,
		FuelCostMultiplier: 2.0,
		WaterCostIncrease:  -5,
		LootChance:         0.05,
		DangerChance:       0.1,
		MovementPenalty:    1,
	},
	HazardDesert: {
		InjuryChance:       0.15,
		InjuryDamage:       10,
		FuelCostMultiplier: 1.3,
		WaterCostIncrease:  10,
		LootChance:         0.05,
		DangerChance:       0.1,
		MovementPenalty:    1,
	},
	HazardRuin: {
		InjuryChance:       0.2,
		InjuryDamage:       12,
		FuelCostMultiplier: 1.2,
		WaterCostIncrease:  0,
		LootChance:         0.4,
		DangerChance:       0.35,
		MovementPenalty:    1,
	},
	HazardSwamp: {
		InjuryChance:       0.15,
		InjuryDamage:       8,
		FuelCostMultiplier: 1.8,
		WaterCostIncrease:  -3,
		LootChance:         0.1,
		DangerChance:       0.25,
		MovementPenalty:    2,
	},
	HazardRadiation: {
		InjuryChance:       0.5,
		InjuryDamage:       20,
		FuelCostMultiplier: 1.2,
		WaterCostIncrease:  5,
		LootChance:         0.2,
		DangerChance:       0.15,
		MovementPenalty:    1,
	},
	HazardMineField: {
		InjuryChance:       0.4,
		InjuryDamage:       30,
		FuelCostMultiplier: 2.0,
		WaterCostIncrease:  0,
		LootChance:         0.1,
		DangerChance:       0.1,
		MovementPenalty:    3,
	},
}

// GetHazardEffect returns the effect for a terrain hazard.
func GetHazardEffect(h TerrainHazard) HazardEffect {
	if effect, ok := HazardEffects[h]; ok {
		return effect
	}
	return HazardEffects[HazardNone]
}

// GenreWeatherSubset returns appropriate weather types for a genre.
func GenreWeatherSubset(genre engine.GenreID) []WeatherType {
	subsets := map[engine.GenreID][]WeatherType{
		engine.GenreFantasy: {
			WeatherClear, WeatherStorm, WeatherBlizzard, WeatherHeatwave,
			WeatherFlood, WeatherFog, WeatherDustStorm,
		},
		engine.GenreScifi: {
			WeatherClear, WeatherStorm, WeatherBlizzard, WeatherHeatwave,
			WeatherFog, WeatherMeteorShower, WeatherDustStorm, WeatherAcidRain,
		},
		engine.GenreHorror: {
			WeatherClear, WeatherStorm, WeatherBlizzard, WeatherFog,
			WeatherDustStorm, WeatherAcidRain, WeatherFlood,
		},
		engine.GenreCyberpunk: {
			WeatherClear, WeatherStorm, WeatherHeatwave, WeatherFog,
			WeatherFlood, WeatherDustStorm, WeatherAcidRain,
		},
		engine.GenrePostapoc: {
			WeatherClear, WeatherStorm, WeatherBlizzard, WeatherHeatwave,
			WeatherDustStorm, WeatherAcidRain, WeatherMeteorShower,
		},
	}

	if subset, ok := subsets[genre]; ok {
		return subset
	}
	return subsets[engine.GenreFantasy]
}

// GenreHazardSubset returns appropriate hazard types for a genre.
func GenreHazardSubset(genre engine.GenreID) []TerrainHazard {
	subsets := map[engine.GenreID][]TerrainHazard{
		engine.GenreFantasy: {
			HazardNone, HazardMountainPass, HazardRiverCrossing,
			HazardDesert, HazardRuin, HazardSwamp,
		},
		engine.GenreScifi: {
			HazardNone, HazardMountainPass, HazardRiverCrossing,
			HazardDesert, HazardRuin, HazardRadiation, HazardMineField,
		},
		engine.GenreHorror: {
			HazardNone, HazardMountainPass, HazardRiverCrossing,
			HazardRuin, HazardSwamp, HazardMineField,
		},
		engine.GenreCyberpunk: {
			HazardNone, HazardRiverCrossing, HazardDesert,
			HazardRuin, HazardRadiation, HazardMineField,
		},
		engine.GenrePostapoc: {
			HazardNone, HazardMountainPass, HazardRiverCrossing,
			HazardDesert, HazardRuin, HazardRadiation, HazardMineField,
		},
	}

	if subset, ok := subsets[genre]; ok {
		return subset
	}
	return subsets[engine.GenreFantasy]
}
