package vessel

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
	"github.com/opd-ai/voyage/pkg/resources"
)

// LoadoutType identifies the preset loadout configuration.
type LoadoutType int

const (
	// LoadoutBalanced provides moderate stats across all categories.
	LoadoutBalanced LoadoutType = iota
	// LoadoutFastLight prioritizes speed over capacity.
	LoadoutFastLight
	// LoadoutSlowHeavy prioritizes capacity over speed.
	LoadoutSlowHeavy
)

// AllLoadoutTypes returns all available loadout types.
func AllLoadoutTypes() []LoadoutType {
	return []LoadoutType{LoadoutBalanced, LoadoutFastLight, LoadoutSlowHeavy}
}

// LoadoutName returns the genre-appropriate name for a loadout type.
func LoadoutName(lt LoadoutType, genre engine.GenreID) string {
	names, ok := loadoutNames[genre]
	if !ok {
		names = loadoutNames[engine.GenreFantasy]
	}
	return names[lt]
}

var loadoutNames = map[engine.GenreID]map[LoadoutType]string{
	engine.GenreFantasy: {
		LoadoutBalanced:  "Merchant's Caravan",
		LoadoutFastLight: "Swift Courier",
		LoadoutSlowHeavy: "Armored Wagon Train",
	},
	engine.GenreScifi: {
		LoadoutBalanced:  "Survey Vessel",
		LoadoutFastLight: "Scout Ship",
		LoadoutSlowHeavy: "Heavy Freighter",
	},
	engine.GenreHorror: {
		LoadoutBalanced:  "Survivor's Ride",
		LoadoutFastLight: "Escape Vehicle",
		LoadoutSlowHeavy: "Fortified Bus",
	},
	engine.GenreCyberpunk: {
		LoadoutBalanced:  "Standard Runner",
		LoadoutFastLight: "Speed Demon",
		LoadoutSlowHeavy: "Cargo Hauler",
	},
	engine.GenrePostapoc: {
		LoadoutBalanced:  "Road Warrior",
		LoadoutFastLight: "Scout Bike",
		LoadoutSlowHeavy: "War Rig",
	},
}

// LoadoutDescription returns a brief description of the loadout.
func LoadoutDescription(lt LoadoutType) string {
	descriptions := map[LoadoutType]string{
		LoadoutBalanced:  "Balanced speed, capacity, and durability. Good for beginners.",
		LoadoutFastLight: "Fast travel but limited cargo. Best for experienced players.",
		LoadoutSlowHeavy: "Maximum cargo and durability but slow. Best for cautious players.",
	}
	return descriptions[lt]
}

// Loadout represents a starting configuration for vessel and resources.
type Loadout struct {
	Type          LoadoutType
	Genre         engine.GenreID
	VesselType    VesselType
	SpeedMod      float64 // Multiplier applied to base speed
	CapacityMod   float64 // Multiplier applied to base capacity
	IntegrityMod  float64 // Multiplier applied to base integrity
	StartFood     float64
	StartWater    float64
	StartFuel     float64
	StartMeds     float64
	StartMorale   float64
	StartCurrency float64
	CrewCount     int
}

// LoadoutGenerator creates procedural loadouts from a seed.
type LoadoutGenerator struct {
	gen   *seed.Generator
	genre engine.GenreID
}

// NewLoadoutGenerator creates a new loadout generator.
func NewLoadoutGenerator(masterSeed int64, genre engine.GenreID) *LoadoutGenerator {
	return &LoadoutGenerator{
		gen:   seed.NewGenerator(masterSeed, "loadout"),
		genre: genre,
	}
}

// SetGenre changes the generator's genre.
func (g *LoadoutGenerator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// Generate creates a loadout for the given type with procedural variation.
func (g *LoadoutGenerator) Generate(lt LoadoutType) *Loadout {
	base := baseLoadouts[lt]

	// Add procedural variation (+/- 10%)
	variation := func(base float64) float64 {
		variance := 0.1 * base
		return base + g.gen.RangeFloat64(-variance, variance)
	}

	return &Loadout{
		Type:          lt,
		Genre:         g.genre,
		VesselType:    base.VesselType,
		SpeedMod:      variation(base.SpeedMod),
		CapacityMod:   variation(base.CapacityMod),
		IntegrityMod:  variation(base.IntegrityMod),
		StartFood:     variation(base.StartFood),
		StartWater:    variation(base.StartWater),
		StartFuel:     variation(base.StartFuel),
		StartMeds:     variation(base.StartMeds),
		StartMorale:   variation(base.StartMorale),
		StartCurrency: variation(base.StartCurrency),
		CrewCount:     base.CrewCount + g.gen.Range(-1, 1),
	}
}

// GenerateAll generates all three loadout types.
func (g *LoadoutGenerator) GenerateAll() []*Loadout {
	loadouts := make([]*Loadout, 0, 3)
	for _, lt := range AllLoadoutTypes() {
		loadouts = append(loadouts, g.Generate(lt))
	}
	return loadouts
}

// baseLoadouts defines the base stats for each loadout type.
var baseLoadouts = map[LoadoutType]struct {
	VesselType    VesselType
	SpeedMod      float64
	CapacityMod   float64
	IntegrityMod  float64
	StartFood     float64
	StartWater    float64
	StartFuel     float64
	StartMeds     float64
	StartMorale   float64
	StartCurrency float64
	CrewCount     int
}{
	LoadoutBalanced: {
		VesselType:    VesselMedium,
		SpeedMod:      1.0,
		CapacityMod:   1.0,
		IntegrityMod:  1.0,
		StartFood:     75,
		StartWater:    75,
		StartFuel:     80,
		StartMeds:     30,
		StartMorale:   80,
		StartCurrency: 50,
		CrewCount:     4,
	},
	LoadoutFastLight: {
		VesselType:    VesselSmall,
		SpeedMod:      1.3,
		CapacityMod:   0.7,
		IntegrityMod:  0.8,
		StartFood:     60,
		StartWater:    60,
		StartFuel:     100,
		StartMeds:     20,
		StartMorale:   90,
		StartCurrency: 30,
		CrewCount:     3,
	},
	LoadoutSlowHeavy: {
		VesselType:    VesselLarge,
		SpeedMod:      0.8,
		CapacityMod:   1.4,
		IntegrityMod:  1.3,
		StartFood:     100,
		StartWater:    100,
		StartFuel:     70,
		StartMeds:     50,
		StartMorale:   70,
		StartCurrency: 80,
		CrewCount:     5,
	},
}

// ApplyToVessel configures a vessel based on this loadout.
func (l *Loadout) ApplyToVessel(v *Vessel) {
	baseStats := DefaultVesselStats[l.VesselType]

	v.vesselType = l.VesselType
	v.maxIntegrity = baseStats.MaxIntegrity * l.IntegrityMod
	v.integrity = v.maxIntegrity
	v.baseSpeed = baseStats.BaseSpeed * l.SpeedMod
	v.speed = v.baseSpeed
	v.capacity = int(float64(baseStats.Capacity) * l.CapacityMod)
	v.usedCapacity = 0
	v.genre = l.Genre
}

// ApplyToResources configures resources based on this loadout.
func (l *Loadout) ApplyToResources(r *resources.Resources) {
	r.Set(resources.ResourceFood, l.StartFood)
	r.Set(resources.ResourceWater, l.StartWater)
	r.Set(resources.ResourceFuel, l.StartFuel)
	r.Set(resources.ResourceMedicine, l.StartMeds)
	r.Set(resources.ResourceMorale, l.StartMorale)
	r.Set(resources.ResourceCurrency, l.StartCurrency)
}

// Name returns the genre-appropriate name for this loadout.
func (l *Loadout) Name() string {
	return LoadoutName(l.Type, l.Genre)
}

// Description returns the description for this loadout.
func (l *Loadout) Description() string {
	return LoadoutDescription(l.Type)
}
