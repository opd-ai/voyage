package vessel

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// VesselType identifies the type of transport.
type VesselType int

const (
	// VesselSmall is a small transport with low capacity but fast speed.
	VesselSmall VesselType = iota
	// VesselMedium is a balanced transport.
	VesselMedium
	// VesselLarge is a large transport with high capacity but slow speed.
	VesselLarge
)

// AllVesselTypes returns all vessel types.
func AllVesselTypes() []VesselType {
	return []VesselType{VesselSmall, VesselMedium, VesselLarge}
}

// Vessel represents the party's transport with integrity, speed, and capacity.
type Vessel struct {
	genre        engine.GenreID
	vesselType   VesselType
	integrity    float64
	maxIntegrity float64
	speed        float64
	baseSpeed    float64
	capacity     int
	usedCapacity int
}

// DefaultVesselStats returns the base stats for each vessel type.
var DefaultVesselStats = map[VesselType]struct {
	MaxIntegrity float64
	BaseSpeed    float64
	Capacity     int
}{
	VesselSmall:  {MaxIntegrity: 80, BaseSpeed: 1.5, Capacity: 20},
	VesselMedium: {MaxIntegrity: 100, BaseSpeed: 1.0, Capacity: 50},
	VesselLarge:  {MaxIntegrity: 150, BaseSpeed: 0.7, Capacity: 100},
}

// NewVessel creates a new vessel of the given type.
func NewVessel(vesselType VesselType, genre engine.GenreID) *Vessel {
	stats := DefaultVesselStats[vesselType]
	return &Vessel{
		genre:        genre,
		vesselType:   vesselType,
		integrity:    stats.MaxIntegrity,
		maxIntegrity: stats.MaxIntegrity,
		speed:        stats.BaseSpeed,
		baseSpeed:    stats.BaseSpeed,
		capacity:     stats.Capacity,
		usedCapacity: 0,
	}
}

// SetGenre changes the vessel vocabulary theme.
func (v *Vessel) SetGenre(genre engine.GenreID) {
	v.genre = genre
}

// Genre returns the current genre.
func (v *Vessel) Genre() engine.GenreID {
	return v.genre
}

// Type returns the vessel type.
func (v *Vessel) Type() VesselType {
	return v.vesselType
}

// Integrity returns the current hull integrity.
func (v *Vessel) Integrity() float64 {
	return v.integrity
}

// MaxIntegrity returns the maximum hull integrity.
func (v *Vessel) MaxIntegrity() float64 {
	return v.maxIntegrity
}

// IntegrityRatio returns integrity as a ratio [0, 1].
func (v *Vessel) IntegrityRatio() float64 {
	if v.maxIntegrity <= 0 {
		return 0
	}
	return v.integrity / v.maxIntegrity
}

// Speed returns the current movement speed multiplier.
func (v *Vessel) Speed() float64 {
	return v.speed
}

// BaseSpeed returns the base movement speed.
func (v *Vessel) BaseSpeed() float64 {
	return v.baseSpeed
}

// Capacity returns the total cargo capacity.
func (v *Vessel) Capacity() int {
	return v.capacity
}

// UsedCapacity returns the currently used cargo capacity.
func (v *Vessel) UsedCapacity() int {
	return v.usedCapacity
}

// FreeCapacity returns the remaining cargo capacity.
func (v *Vessel) FreeCapacity() int {
	return v.capacity - v.usedCapacity
}

// TakeDamage reduces integrity by the given amount.
// Returns true if the vessel is destroyed.
func (v *Vessel) TakeDamage(amount float64) bool {
	v.integrity -= amount
	if v.integrity <= 0 {
		v.integrity = 0
		return true
	}
	v.updateSpeed()
	return false
}

// Repair increases integrity by the given amount.
func (v *Vessel) Repair(amount float64) {
	v.integrity += amount
	if v.integrity > v.maxIntegrity {
		v.integrity = v.maxIntegrity
	}
	v.updateSpeed()
}

// updateSpeed adjusts speed based on damage level.
func (v *Vessel) updateSpeed() {
	ratio := v.IntegrityRatio()
	// Speed scales with integrity: at 50% integrity, speed is 75%
	// At 25% integrity, speed is 50%
	v.speed = v.baseSpeed * (0.5 + 0.5*ratio)
}

// IsDestroyed returns true if the vessel has zero integrity.
func (v *Vessel) IsDestroyed() bool {
	return v.integrity <= 0
}

// IsCritical returns true if integrity is below 25%.
func (v *Vessel) IsCritical() bool {
	return v.IntegrityRatio() < 0.25
}

// IsDamaged returns true if integrity is below 75%.
func (v *Vessel) IsDamaged() bool {
	return v.IntegrityRatio() < 0.75
}

// AddCargo adds cargo weight. Returns true if successful.
func (v *Vessel) AddCargo(weight int) bool {
	if v.usedCapacity+weight > v.capacity {
		return false
	}
	v.usedCapacity += weight
	return true
}

// RemoveCargo removes cargo weight. Returns true if successful.
func (v *Vessel) RemoveCargo(weight int) bool {
	if weight > v.usedCapacity {
		return false
	}
	v.usedCapacity -= weight
	return true
}

// Name returns the genre-appropriate name for this vessel type.
func (v *Vessel) Name() string {
	return VesselName(v.vesselType, v.genre)
}

// VesselName returns the genre-specific name for a vessel type.
func VesselName(vt VesselType, genre engine.GenreID) string {
	names, ok := vesselNames[genre]
	if !ok {
		names = vesselNames[engine.GenreFantasy]
	}
	return names[vt]
}

var vesselNames = map[engine.GenreID]map[VesselType]string{
	engine.GenreFantasy: {
		VesselSmall:  "Pony Cart",
		VesselMedium: "Wagon",
		VesselLarge:  "Caravan",
	},
	engine.GenreScifi: {
		VesselSmall:  "Scout Pod",
		VesselMedium: "Shuttle",
		VesselLarge:  "Freighter",
	},
	engine.GenreHorror: {
		VesselSmall:  "Motorcycle",
		VesselMedium: "SUV",
		VesselLarge:  "Bus",
	},
	engine.GenreCyberpunk: {
		VesselSmall:  "Speedbike",
		VesselMedium: "Aerodyne",
		VesselLarge:  "Road Train",
	},
	engine.GenrePostapoc: {
		VesselSmall:  "Dirt Bike",
		VesselMedium: "War Rig",
		VesselLarge:  "Convoy",
	},
}
