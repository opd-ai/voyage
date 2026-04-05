package destination

import "github.com/opd-ai/voyage/pkg/engine"

// DestinationType represents a category of destination.
type DestinationType int

const (
	City DestinationType = iota
	Sanctuary
	Treasure
	Escape
	Settlement
)

// String returns the base type name.
func (d DestinationType) String() string {
	return [...]string{"City", "Sanctuary", "Treasure", "Escape", "Settlement"}[d]
}

// AllDestinationTypes returns all destination types.
func AllDestinationTypes() []DestinationType {
	return []DestinationType{City, Sanctuary, Treasure, Escape, Settlement}
}

// DiscoveryPhase represents how close the party is to the destination.
type DiscoveryPhase int

const (
	Distant  DiscoveryPhase = iota // Far away, vague hints
	Signs                          // Medium range, signs of activity
	Approach                       // Close range, clear view
	Arrival                        // At the destination
)

// String returns the phase name.
func (p DiscoveryPhase) String() string {
	return [...]string{"Distant", "Signs", "Approach", "Arrival"}[p]
}

// destinationNames maps genre to destination type names.
var destinationNames = map[engine.GenreID]map[DestinationType]string{
	engine.GenreFantasy: {
		City:       "Walled City",
		Sanctuary:  "Sacred Grove",
		Treasure:   "Dragon's Hoard",
		Escape:     "Portal Gate",
		Settlement: "Village",
	},
	engine.GenreScifi: {
		City:       "Space Station",
		Sanctuary:  "Orbital Haven",
		Treasure:   "Derelict Cargo Hold",
		Escape:     "Jump Gate",
		Settlement: "Colony Dome",
	},
	engine.GenreHorror: {
		City:       "Abandoned Town",
		Sanctuary:  "Fortified Church",
		Treasure:   "Sealed Crypt",
		Escape:     "Last Bridge",
		Settlement: "Survivor Camp",
	},
	engine.GenreCyberpunk: {
		City:       "Megacity Sector",
		Sanctuary:  "Off-Grid Safehouse",
		Treasure:   "Data Vault",
		Escape:     "Black Market Port",
		Settlement: "Free Zone",
	},
	engine.GenrePostapoc: {
		City:       "Survivor Camp",
		Sanctuary:  "Underground Bunker",
		Treasure:   "Supply Cache",
		Escape:     "Evacuation Point",
		Settlement: "Reclaimed Zone",
	},
}

// destinationDescriptions maps genre to destination type descriptions.
var destinationDescriptions = map[engine.GenreID]map[DestinationType]string{
	engine.GenreFantasy: {
		City:       "A great walled city with towers and spires",
		Sanctuary:  "A sacred grove blessed by the old gods",
		Treasure:   "A dragon's hoard of legendary wealth",
		Escape:     "An ancient portal gate to distant lands",
		Settlement: "A humble village of farmers and craftsmen",
	},
	engine.GenreScifi: {
		City:       "A massive orbital station teeming with commerce",
		Sanctuary:  "A protected habitat with life support for all",
		Treasure:   "A derelict vessel's cargo hold full of salvage",
		Escape:     "A FTL jump gate to safety",
		Settlement: "A domed colony on a terraformed world",
	},
	engine.GenreHorror: {
		City:       "A fog-shrouded town with boarded windows",
		Sanctuary:  "A fortified church warded against evil",
		Treasure:   "A sealed crypt rumored to hold relics",
		Escape:     "The last bridge out of the cursed land",
		Settlement: "A desperate camp of survivors",
	},
	engine.GenreCyberpunk: {
		City:       "A neon-lit megacity sector controlled by corps",
		Sanctuary:  "An off-grid safehouse hidden from the net",
		Treasure:   "A secure data vault with priceless intel",
		Escape:     "A black market port with forged papers ready",
		Settlement: "A free zone outside corporate jurisdiction",
	},
	engine.GenrePostapoc: {
		City:       "A fortified camp of survivors",
		Sanctuary:  "A sealed bunker untouched by the apocalypse",
		Treasure:   "A pre-war supply cache with vital resources",
		Escape:     "The last evacuation point before the zone closes",
		Settlement: "A reclaimed area where survivors have rebuilt",
	},
}
