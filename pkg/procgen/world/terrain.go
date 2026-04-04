package world

import "github.com/opd-ai/voyage/pkg/engine"

// TerrainType represents a type of terrain on the map.
type TerrainType int

const (
	// TerrainPlains is open, easy-to-traverse terrain.
	TerrainPlains TerrainType = iota
	// TerrainForest is wooded terrain with moderate difficulty.
	TerrainForest
	// TerrainMountain is difficult terrain with high movement cost.
	TerrainMountain
	// TerrainDesert is hot terrain with increased water consumption.
	TerrainDesert
	// TerrainRiver is water terrain requiring crossing.
	TerrainRiver
	// TerrainSwamp is slow, hazardous terrain.
	TerrainSwamp
	// TerrainRuin is explorable terrain with potential discoveries.
	TerrainRuin
)

// TerrainInfo contains metadata about a terrain type.
type TerrainInfo struct {
	Type          TerrainType
	Name          string
	MovementCost  int
	WaterModifier float64
	FoodModifier  float64
	HazardChance  float64
}

// DefaultTerrainInfo returns terrain info for the given type and genre.
func DefaultTerrainInfo(t TerrainType, genre engine.GenreID) TerrainInfo {
	base := terrainBaseInfo[t]
	base.Name = terrainNames[genre][t]
	return base
}

var terrainBaseInfo = map[TerrainType]TerrainInfo{
	TerrainPlains:   {Type: TerrainPlains, MovementCost: 1, WaterModifier: 1.0, FoodModifier: 1.0, HazardChance: 0.05},
	TerrainForest:   {Type: TerrainForest, MovementCost: 2, WaterModifier: 0.8, FoodModifier: 0.8, HazardChance: 0.15},
	TerrainMountain: {Type: TerrainMountain, MovementCost: 3, WaterModifier: 1.5, FoodModifier: 1.5, HazardChance: 0.25},
	TerrainDesert:   {Type: TerrainDesert, MovementCost: 2, WaterModifier: 2.0, FoodModifier: 1.2, HazardChance: 0.20},
	TerrainRiver:    {Type: TerrainRiver, MovementCost: 2, WaterModifier: 0.5, FoodModifier: 0.7, HazardChance: 0.10},
	TerrainSwamp:    {Type: TerrainSwamp, MovementCost: 3, WaterModifier: 0.9, FoodModifier: 1.3, HazardChance: 0.30},
	TerrainRuin:     {Type: TerrainRuin, MovementCost: 2, WaterModifier: 1.0, FoodModifier: 1.0, HazardChance: 0.35},
}

var terrainNames = map[engine.GenreID]map[TerrainType]string{
	engine.GenreFantasy: {
		TerrainPlains:   "Plains",
		TerrainForest:   "Forest",
		TerrainMountain: "Mountain",
		TerrainDesert:   "Desert",
		TerrainRiver:    "River",
		TerrainSwamp:    "Swamp",
		TerrainRuin:     "Ancient Ruins",
	},
	engine.GenreScifi: {
		TerrainPlains:   "Open Space",
		TerrainForest:   "Nebula",
		TerrainMountain: "Asteroid Belt",
		TerrainDesert:   "Radiation Zone",
		TerrainRiver:    "Ion Stream",
		TerrainSwamp:    "Gravity Well",
		TerrainRuin:     "Derelict Station",
	},
	engine.GenreHorror: {
		TerrainPlains:   "Open Road",
		TerrainForest:   "Dead Forest",
		TerrainMountain: "Rubble Pile",
		TerrainDesert:   "Wasteland",
		TerrainRiver:    "Toxic River",
		TerrainSwamp:    "Infected Zone",
		TerrainRuin:     "Abandoned City",
	},
	engine.GenreCyberpunk: {
		TerrainPlains:   "Street Level",
		TerrainForest:   "Neon District",
		TerrainMountain: "Mega Tower",
		TerrainDesert:   "Industrial Zone",
		TerrainRiver:    "Data Stream",
		TerrainSwamp:    "Toxic Slum",
		TerrainRuin:     "Collapsed Sector",
	},
	engine.GenrePostapoc: {
		TerrainPlains:   "Dustbowl",
		TerrainForest:   "Dead Woods",
		TerrainMountain: "Scrap Mountain",
		TerrainDesert:   "Irradiated Zone",
		TerrainRiver:    "Polluted River",
		TerrainSwamp:    "Mutant Bog",
		TerrainRuin:     "Pre-War Ruins",
	},
}

// AllTerrainTypes returns all terrain types.
func AllTerrainTypes() []TerrainType {
	return []TerrainType{
		TerrainPlains,
		TerrainForest,
		TerrainMountain,
		TerrainDesert,
		TerrainRiver,
		TerrainSwamp,
		TerrainRuin,
	}
}
