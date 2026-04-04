package world

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// BiomeType represents a regional climate/theme.
type BiomeType int

const (
	// BiomeTemperate has balanced terrain distribution.
	BiomeTemperate BiomeType = iota
	// BiomeArid has desert-dominant terrain.
	BiomeArid
	// BiomeForested has forest-dominant terrain.
	BiomeForested
	// BiomeMountainous has mountain-dominant terrain.
	BiomeMountainous
	// BiomeWetland has swamp and river terrain.
	BiomeWetland
	// BiomeRuined has ruin-dominant terrain.
	BiomeRuined
)

// BiomeInfo contains metadata about a biome.
type BiomeInfo struct {
	Type             BiomeType
	Name             string
	TerrainWeights   map[TerrainType]float64
	HazardModifier   float64
	ResourceModifier float64
}

// DefaultBiomeInfo returns biome info for the given type and genre.
func DefaultBiomeInfo(b BiomeType, genre engine.GenreID) BiomeInfo {
	base := biomeBaseInfo[b]
	base.Name = biomeNames[genre][b]
	return base
}

var biomeBaseInfo = map[BiomeType]BiomeInfo{
	BiomeTemperate: {
		Type: BiomeTemperate,
		TerrainWeights: map[TerrainType]float64{
			TerrainPlains:   0.40,
			TerrainForest:   0.30,
			TerrainMountain: 0.10,
			TerrainRiver:    0.10,
			TerrainRuin:     0.10,
		},
		HazardModifier:   1.0,
		ResourceModifier: 1.0,
	},
	BiomeArid: {
		Type: BiomeArid,
		TerrainWeights: map[TerrainType]float64{
			TerrainPlains:   0.20,
			TerrainDesert:   0.50,
			TerrainMountain: 0.15,
			TerrainRuin:     0.15,
		},
		HazardModifier:   1.2,
		ResourceModifier: 0.7,
	},
	BiomeForested: {
		Type: BiomeForested,
		TerrainWeights: map[TerrainType]float64{
			TerrainPlains: 0.15,
			TerrainForest: 0.55,
			TerrainRiver:  0.15,
			TerrainSwamp:  0.10,
			TerrainRuin:   0.05,
		},
		HazardModifier:   1.1,
		ResourceModifier: 1.2,
	},
	BiomeMountainous: {
		Type: BiomeMountainous,
		TerrainWeights: map[TerrainType]float64{
			TerrainPlains:   0.10,
			TerrainMountain: 0.50,
			TerrainForest:   0.20,
			TerrainRiver:    0.10,
			TerrainRuin:     0.10,
		},
		HazardModifier:   1.3,
		ResourceModifier: 0.8,
	},
	BiomeWetland: {
		Type: BiomeWetland,
		TerrainWeights: map[TerrainType]float64{
			TerrainPlains: 0.15,
			TerrainSwamp:  0.40,
			TerrainRiver:  0.30,
			TerrainForest: 0.10,
			TerrainRuin:   0.05,
		},
		HazardModifier:   1.4,
		ResourceModifier: 0.9,
	},
	BiomeRuined: {
		Type: BiomeRuined,
		TerrainWeights: map[TerrainType]float64{
			TerrainPlains:   0.20,
			TerrainRuin:     0.50,
			TerrainDesert:   0.15,
			TerrainMountain: 0.15,
		},
		HazardModifier:   1.5,
		ResourceModifier: 1.1,
	},
}

var biomeNames = map[engine.GenreID]map[BiomeType]string{
	engine.GenreFantasy: {
		BiomeTemperate:   "Verdant Lands",
		BiomeArid:        "Sunscorched Wastes",
		BiomeForested:    "Enchanted Woods",
		BiomeMountainous: "Dragon Peaks",
		BiomeWetland:     "Misty Fens",
		BiomeRuined:      "Fallen Kingdom",
	},
	engine.GenreScifi: {
		BiomeTemperate:   "Habitable Zone",
		BiomeArid:        "Stellar Desert",
		BiomeForested:    "Nebula Cloud",
		BiomeMountainous: "Asteroid Field",
		BiomeWetland:     "Ion Storm",
		BiomeRuined:      "Debris Field",
	},
	engine.GenreHorror: {
		BiomeTemperate:   "Survivor Territory",
		BiomeArid:        "Scorched Earth",
		BiomeForested:    "Haunted Forest",
		BiomeMountainous: "Corpse Mountain",
		BiomeWetland:     "Plague Marsh",
		BiomeRuined:      "Dead City",
	},
	engine.GenreCyberpunk: {
		BiomeTemperate:   "Corporate Zone",
		BiomeArid:        "Industrial Sprawl",
		BiomeForested:    "Vertical Gardens",
		BiomeMountainous: "Mega Towers",
		BiomeWetland:     "Flooded District",
		BiomeRuined:      "War Zone",
	},
	engine.GenrePostapoc: {
		BiomeTemperate:   "Reclaimed Lands",
		BiomeArid:        "Rad Desert",
		BiomeForested:    "Mutant Woods",
		BiomeMountainous: "Scrap Peaks",
		BiomeWetland:     "Toxic Bayou",
		BiomeRuined:      "Ground Zero",
	},
}

// AllBiomeTypes returns all biome types.
func AllBiomeTypes() []BiomeType {
	return []BiomeType{
		BiomeTemperate,
		BiomeArid,
		BiomeForested,
		BiomeMountainous,
		BiomeWetland,
		BiomeRuined,
	}
}

// SelectTerrain picks a terrain type based on biome weights.
func SelectTerrain(gen *seed.Generator, biome BiomeType) TerrainType {
	info := biomeBaseInfo[biome]

	// Use fixed order for determinism (maps iterate non-deterministically)
	allTerrains := AllTerrainTypes()
	terrains := make([]TerrainType, 0)
	weights := make([]float64, 0)

	for _, terrain := range allTerrains {
		if weight, ok := info.TerrainWeights[terrain]; ok {
			terrains = append(terrains, terrain)
			weights = append(weights, weight)
		}
	}

	return seed.WeightedChoice(gen, terrains, weights)
}
