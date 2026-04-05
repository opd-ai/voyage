package weather

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// WeatherType identifies different weather conditions.
type WeatherType int

const (
	// WeatherClear is normal conditions.
	WeatherClear WeatherType = iota
	// WeatherStorm is high wind and precipitation.
	WeatherStorm
	// WeatherBlizzard is extreme cold and low visibility.
	WeatherBlizzard
	// WeatherHeatwave is extreme heat.
	WeatherHeatwave
	// WeatherFlood is water hazard.
	WeatherFlood
	// WeatherFog is low visibility.
	WeatherFog
	// WeatherMeteorShower is space debris.
	WeatherMeteorShower
	// WeatherDustStorm is choking particulates.
	WeatherDustStorm
	// WeatherAcidRain is corrosive precipitation.
	WeatherAcidRain
)

// AllWeatherTypes returns all weather types.
func AllWeatherTypes() []WeatherType {
	return []WeatherType{
		WeatherClear,
		WeatherStorm,
		WeatherBlizzard,
		WeatherHeatwave,
		WeatherFlood,
		WeatherFog,
		WeatherMeteorShower,
		WeatherDustStorm,
		WeatherAcidRain,
	}
}

// HazardousWeatherTypes returns only hazardous weather types (excludes Clear).
func HazardousWeatherTypes() []WeatherType {
	return []WeatherType{
		WeatherStorm,
		WeatherBlizzard,
		WeatherHeatwave,
		WeatherFlood,
		WeatherFog,
		WeatherMeteorShower,
		WeatherDustStorm,
		WeatherAcidRain,
	}
}

// WeatherName returns the genre-appropriate name for a weather type.
func WeatherName(w WeatherType, genre engine.GenreID) string {
	names := weatherNames[genre]
	if names == nil {
		names = weatherNames[engine.GenreFantasy]
	}
	return names[w]
}

var weatherNames = map[engine.GenreID]map[WeatherType]string{
	engine.GenreFantasy: {
		WeatherClear:        "Clear Skies",
		WeatherStorm:        "Tempest",
		WeatherBlizzard:     "Blizzard",
		WeatherHeatwave:     "Scorching Heat",
		WeatherFlood:        "Flash Flood",
		WeatherFog:          "Mystic Fog",
		WeatherMeteorShower: "Falling Stars",
		WeatherDustStorm:    "Sand Tempest",
		WeatherAcidRain:     "Cursed Rain",
	},
	engine.GenreScifi: {
		WeatherClear:        "Nominal Conditions",
		WeatherStorm:        "Ion Storm",
		WeatherBlizzard:     "Cryo Event",
		WeatherHeatwave:     "Solar Flare",
		WeatherFlood:        "Coolant Breach",
		WeatherFog:          "Sensor Interference",
		WeatherMeteorShower: "Meteor Shower",
		WeatherDustStorm:    "Particle Storm",
		WeatherAcidRain:     "Corrosive Atmosphere",
	},
	engine.GenreHorror: {
		WeatherClear:        "Ominous Calm",
		WeatherStorm:        "Storm",
		WeatherBlizzard:     "Blizzard",
		WeatherHeatwave:     "Oppressive Heat",
		WeatherFlood:        "Flood Waters",
		WeatherFog:          "Thick Fog",
		WeatherMeteorShower: "Blood Rain",
		WeatherDustStorm:    "Ash Storm",
		WeatherAcidRain:     "Toxic Rain",
	},
	engine.GenreCyberpunk: {
		WeatherClear:        "Clear",
		WeatherStorm:        "Electric Storm",
		WeatherBlizzard:     "Ice Storm",
		WeatherHeatwave:     "Heat Dome",
		WeatherFlood:        "Flash Flood",
		WeatherFog:          "Smog",
		WeatherMeteorShower: "Debris Fall",
		WeatherDustStorm:    "Toxic Cloud",
		WeatherAcidRain:     "Acid Rain",
	},
	engine.GenrePostapoc: {
		WeatherClear:        "Clear",
		WeatherStorm:        "Rad Storm",
		WeatherBlizzard:     "Nuclear Winter",
		WeatherHeatwave:     "Heat Blast",
		WeatherFlood:        "Toxic Flood",
		WeatherFog:          "Rad Fog",
		WeatherMeteorShower: "Fallout",
		WeatherDustStorm:    "Dust Storm",
		WeatherAcidRain:     "Acid Rain",
	},
}

// WeatherDescription returns a genre-appropriate description.
func WeatherDescription(w WeatherType, genre engine.GenreID) string {
	descs := weatherDescriptions[genre]
	if descs == nil {
		descs = weatherDescriptions[engine.GenreFantasy]
	}
	return descs[w]
}

var weatherDescriptions = map[engine.GenreID]map[WeatherType]string{
	engine.GenreFantasy: {
		WeatherClear:        "The sky is clear and the wind is fair.",
		WeatherStorm:        "Dark clouds gather and lightning splits the sky.",
		WeatherBlizzard:     "A freezing blizzard descends with blinding snow.",
		WeatherHeatwave:     "The sun beats down mercilessly on the land.",
		WeatherFlood:        "Waters rise quickly, threatening to sweep all away.",
		WeatherFog:          "An unnatural mist shrouds the land in mystery.",
		WeatherMeteorShower: "Stars fall from the heavens like tears of the gods.",
		WeatherDustStorm:    "Sands rise in a choking, blinding wall.",
		WeatherAcidRain:     "The rain burns where it touches, cursed by dark magic.",
	},
	engine.GenreScifi: {
		WeatherClear:        "All environmental readings are within normal parameters.",
		WeatherStorm:        "Electromagnetic interference detected across all spectrums.",
		WeatherBlizzard:     "Cryogenic atmospheric conditions are rapidly developing.",
		WeatherHeatwave:     "Solar activity is exceeding safe exposure limits.",
		WeatherFlood:        "Catastrophic coolant system failure detected.",
		WeatherFog:          "Sensor systems compromised by unknown interference.",
		WeatherMeteorShower: "Multiple impact warnings from orbital debris.",
		WeatherDustStorm:    "Particulate density exceeding filtration capacity.",
		WeatherAcidRain:     "Atmospheric acidity reaching critical levels.",
	},
	engine.GenreHorror: {
		WeatherClear:        "It's quiet. Too quiet.",
		WeatherStorm:        "Thunder rolls across the land like angry drums.",
		WeatherBlizzard:     "The cold cuts through everything.",
		WeatherHeatwave:     "The air is thick and hard to breathe.",
		WeatherFlood:        "Dark waters rise, hiding what lurks beneath.",
		WeatherFog:          "You can't see more than a few feet ahead.",
		WeatherMeteorShower: "The sky bleeds red as fire falls from above.",
		WeatherDustStorm:    "Ash fills the air, choking the living.",
		WeatherAcidRain:     "The rain burns. Don't let it touch your skin.",
	},
	engine.GenreCyberpunk: {
		WeatherClear:        "A rare day without smog. Enjoy it.",
		WeatherStorm:        "EMPs and lightning. Stay off the grid.",
		WeatherBlizzard:     "Temperature crash. Ice coats everything.",
		WeatherHeatwave:     "Heat dome over the city. Stay inside.",
		WeatherFlood:        "Streets are rivers. Seek high ground.",
		WeatherFog:          "The smog is thick today. Filters recommended.",
		WeatherMeteorShower: "Orbital debris incoming. Take cover.",
		WeatherDustStorm:    "Industrial waste cloud. Seal your vehicle.",
		WeatherAcidRain:     "Acid rain advisory. Protective gear required.",
	},
	engine.GenrePostapoc: {
		WeatherClear:        "A rare clear day in the wastes.",
		WeatherStorm:        "Radiation spikes with the storm clouds.",
		WeatherBlizzard:     "Nuclear winter grips the land in ice.",
		WeatherHeatwave:     "The sun scorches the irradiated earth.",
		WeatherFlood:        "Toxic waters rise from the poisoned ground.",
		WeatherFog:          "A radioactive mist rolls across the wastes.",
		WeatherMeteorShower: "Something's falling from the sky again.",
		WeatherDustStorm:    "Choking dust obscures everything.",
		WeatherAcidRain:     "The rain is poison. Find shelter.",
	},
}

// TerrainHazard identifies types of terrain-based hazards.
type TerrainHazard int

const (
	// HazardNone means no special hazard.
	HazardNone TerrainHazard = iota
	// HazardMountainPass poses injury risk.
	HazardMountainPass
	// HazardRiverCrossing costs extra fuel.
	HazardRiverCrossing
	// HazardDesert causes water crisis.
	HazardDesert
	// HazardRuin offers loot but danger.
	HazardRuin
	// HazardSwamp slows movement and health risk.
	HazardSwamp
	// HazardRadiation causes health damage.
	HazardRadiation
	// HazardMineField poses explosion risk.
	HazardMineField
)

// AllTerrainHazards returns all terrain hazard types.
func AllTerrainHazards() []TerrainHazard {
	return []TerrainHazard{
		HazardNone,
		HazardMountainPass,
		HazardRiverCrossing,
		HazardDesert,
		HazardRuin,
		HazardSwamp,
		HazardRadiation,
		HazardMineField,
	}
}

// HazardName returns the genre-appropriate name for a terrain hazard.
func HazardName(h TerrainHazard, genre engine.GenreID) string {
	names := hazardNames[genre]
	if names == nil {
		names = hazardNames[engine.GenreFantasy]
	}
	return names[h]
}

var hazardNames = map[engine.GenreID]map[TerrainHazard]string{
	engine.GenreFantasy: {
		HazardNone:          "Safe Passage",
		HazardMountainPass:  "Treacherous Pass",
		HazardRiverCrossing: "River Ford",
		HazardDesert:        "Parched Wastes",
		HazardRuin:          "Ancient Ruins",
		HazardSwamp:         "Fetid Marsh",
		HazardRadiation:     "Cursed Ground",
		HazardMineField:     "Trapped Path",
	},
	engine.GenreScifi: {
		HazardNone:          "Clear Route",
		HazardMountainPass:  "Asteroid Field",
		HazardRiverCrossing: "Gravity Well",
		HazardDesert:        "Barren World",
		HazardRuin:          "Derelict Station",
		HazardSwamp:         "Bio-Hazard Zone",
		HazardRadiation:     "Radiation Belt",
		HazardMineField:     "Mine Field",
	},
	engine.GenreHorror: {
		HazardNone:          "Safe Path",
		HazardMountainPass:  "Cliff Road",
		HazardRiverCrossing: "Bridge Crossing",
		HazardDesert:        "Wasteland",
		HazardRuin:          "Abandoned Town",
		HazardSwamp:         "Swamp",
		HazardRadiation:     "Hot Zone",
		HazardMineField:     "Booby Traps",
	},
	engine.GenreCyberpunk: {
		HazardNone:          "Clear Lane",
		HazardMountainPass:  "Elevated Highway",
		HazardRiverCrossing: "Flooded Tunnel",
		HazardDesert:        "Dead Zone",
		HazardRuin:          "Abandoned Sector",
		HazardSwamp:         "Toxic Pit",
		HazardRadiation:     "Rad Zone",
		HazardMineField:     "Corporate Defenses",
	},
	engine.GenrePostapoc: {
		HazardNone:          "Safe Route",
		HazardMountainPass:  "Mountain Pass",
		HazardRiverCrossing: "River Crossing",
		HazardDesert:        "Scorched Earth",
		HazardRuin:          "Pre-War Ruins",
		HazardSwamp:         "Toxic Marsh",
		HazardRadiation:     "Radiation Zone",
		HazardMineField:     "Minefield",
	},
}
