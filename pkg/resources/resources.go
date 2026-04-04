package resources

import "github.com/opd-ai/voyage/pkg/engine"

// ResourceType identifies a resource category.
type ResourceType int

const (
	// ResourceFood is consumed daily; crew starves if empty.
	ResourceFood ResourceType = iota
	// ResourceWater is consumed daily; faster in hot terrain.
	ResourceWater
	// ResourceFuel is consumed per movement; vessel stops if empty.
	ResourceFuel
	// ResourceMedicine is consumed on injury/disease events.
	ResourceMedicine
	// ResourceMorale affects crew desertion and performance.
	ResourceMorale
	// ResourceCurrency is used at supply points for trading.
	ResourceCurrency
)

// AllResourceTypes returns all resource types.
func AllResourceTypes() []ResourceType {
	return []ResourceType{
		ResourceFood,
		ResourceWater,
		ResourceFuel,
		ResourceMedicine,
		ResourceMorale,
		ResourceCurrency,
	}
}

// Resources holds the current resource levels for a journey.
type Resources struct {
	genre   engine.GenreID
	levels  map[ResourceType]float64
	maxLvls map[ResourceType]float64
}

// NewResources creates a new resource manager with default starting values.
func NewResources(genre engine.GenreID) *Resources {
	r := &Resources{
		genre:   genre,
		levels:  make(map[ResourceType]float64),
		maxLvls: make(map[ResourceType]float64),
	}
	
	// Set default max levels
	r.maxLvls[ResourceFood] = 100
	r.maxLvls[ResourceWater] = 100
	r.maxLvls[ResourceFuel] = 100
	r.maxLvls[ResourceMedicine] = 50
	r.maxLvls[ResourceMorale] = 100
	r.maxLvls[ResourceCurrency] = 1000
	
	// Start with moderate levels
	r.levels[ResourceFood] = 75
	r.levels[ResourceWater] = 75
	r.levels[ResourceFuel] = 80
	r.levels[ResourceMedicine] = 30
	r.levels[ResourceMorale] = 80
	r.levels[ResourceCurrency] = 50
	
	return r
}

// SetGenre changes the resource vocabulary.
func (r *Resources) SetGenre(genre engine.GenreID) {
	r.genre = genre
}

// Genre returns the current genre.
func (r *Resources) Genre() engine.GenreID {
	return r.genre
}

// Get returns the current level of a resource.
func (r *Resources) Get(rt ResourceType) float64 {
	return r.levels[rt]
}

// GetMax returns the maximum level of a resource.
func (r *Resources) GetMax(rt ResourceType) float64 {
	return r.maxLvls[rt]
}

// Set sets the level of a resource, clamping to [0, max].
func (r *Resources) Set(rt ResourceType, value float64) {
	max := r.maxLvls[rt]
	if value < 0 {
		value = 0
	}
	if value > max {
		value = max
	}
	r.levels[rt] = value
}

// Add adds to a resource level, clamping to [0, max].
func (r *Resources) Add(rt ResourceType, delta float64) {
	r.Set(rt, r.levels[rt]+delta)
}

// Consume subtracts from a resource level.
// Returns true if there was enough to consume, false if depleted.
func (r *Resources) Consume(rt ResourceType, amount float64) bool {
	if r.levels[rt] < amount {
		r.levels[rt] = 0
		return false
	}
	r.levels[rt] -= amount
	return true
}

// IsDepleted returns true if the resource is at zero.
func (r *Resources) IsDepleted(rt ResourceType) bool {
	return r.levels[rt] <= 0
}

// GetRatio returns the resource level as a ratio [0, 1].
func (r *Resources) GetRatio(rt ResourceType) float64 {
	max := r.maxLvls[rt]
	if max <= 0 {
		return 0
	}
	return r.levels[rt] / max
}

// GetStatus returns the warning status for a resource.
func (r *Resources) GetStatus(rt ResourceType) ThresholdStatus {
	return GetThresholdStatus(rt, r.GetRatio(rt))
}

// Name returns the genre-appropriate name for a resource.
func (r *Resources) Name(rt ResourceType) string {
	return GetResourceName(rt, r.genre)
}

// All returns a copy of all resource levels.
func (r *Resources) All() map[ResourceType]float64 {
	result := make(map[ResourceType]float64)
	for k, v := range r.levels {
		result[k] = v
	}
	return result
}

// SetMax sets the maximum level for a resource.
func (r *Resources) SetMax(rt ResourceType, max float64) {
	r.maxLvls[rt] = max
	// Clamp current level if needed
	if r.levels[rt] > max {
		r.levels[rt] = max
	}
}

// GetResourceName returns the genre-specific name for a resource type.
func GetResourceName(rt ResourceType, genre engine.GenreID) string {
	names, ok := resourceNames[genre]
	if !ok {
		names = resourceNames[engine.GenreFantasy]
	}
	return names[rt]
}

var resourceNames = map[engine.GenreID]map[ResourceType]string{
	engine.GenreFantasy: {
		ResourceFood:     "Food",
		ResourceWater:    "Water",
		ResourceFuel:     "Stamina",
		ResourceMedicine: "Herbs",
		ResourceMorale:   "Morale",
		ResourceCurrency: "Gold",
	},
	engine.GenreScifi: {
		ResourceFood:     "Rations",
		ResourceWater:    "Hydration",
		ResourceFuel:     "Fuel Cells",
		ResourceMedicine: "Med Packs",
		ResourceMorale:   "Morale",
		ResourceCurrency: "Credits",
	},
	engine.GenreHorror: {
		ResourceFood:     "Supplies",
		ResourceWater:    "Clean Water",
		ResourceFuel:     "Gas",
		ResourceMedicine: "First Aid",
		ResourceMorale:   "Sanity",
		ResourceCurrency: "Barter Goods",
	},
	engine.GenreCyberpunk: {
		ResourceFood:     "Nutrient Paste",
		ResourceWater:    "Filtered Water",
		ResourceFuel:     "Battery",
		ResourceMedicine: "Stims",
		ResourceMorale:   "Edge",
		ResourceCurrency: "Nuyen",
	},
	engine.GenrePostapoc: {
		ResourceFood:     "Canned Food",
		ResourceWater:    "Purified Water",
		ResourceFuel:     "Diesel",
		ResourceMedicine: "Meds",
		ResourceMorale:   "Hope",
		ResourceCurrency: "Caps",
	},
}
