package weather

import (
	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
	"github.com/opd-ai/voyage/pkg/resources"
	"github.com/opd-ai/voyage/pkg/vessel"
)

// System manages weather and environmental hazards.
type System struct {
	gen            *seed.Generator
	genre          engine.GenreID
	currentWeather WeatherType
	weatherTurns   int // Turns until weather changes
	minDuration    int
	maxDuration    int
}

// NewSystem creates a new weather system.
func NewSystem(masterSeed int64, genre engine.GenreID) *System {
	return &System{
		gen:            seed.NewGenerator(masterSeed, "weather"),
		genre:          genre,
		currentWeather: WeatherClear,
		weatherTurns:   3,
		minDuration:    2,
		maxDuration:    5,
	}
}

// SetGenre updates the weather system's genre.
func (s *System) SetGenre(genre engine.GenreID) {
	s.genre = genre
}

// CurrentWeather returns the current weather type.
func (s *System) CurrentWeather() WeatherType {
	return s.currentWeather
}

// TurnsRemaining returns turns until weather might change.
func (s *System) TurnsRemaining() int {
	return s.weatherTurns
}

// AdvanceTurn processes weather for a new turn.
// Returns true if weather changed.
func (s *System) AdvanceTurn() bool {
	s.weatherTurns--
	if s.weatherTurns <= 0 {
		return s.changeWeather()
	}
	return false
}

// changeWeather randomly selects new weather.
func (s *System) changeWeather() bool {
	oldWeather := s.currentWeather

	// Higher chance of clear weather
	if s.gen.Float64() < 0.4 {
		s.currentWeather = WeatherClear
	} else {
		subset := GenreWeatherSubset(s.genre)
		s.currentWeather = seed.Choice(s.gen, subset)
	}

	// Set new duration
	s.weatherTurns = s.minDuration + s.gen.Intn(s.maxDuration-s.minDuration+1)

	return s.currentWeather != oldWeather
}

// ForceWeather sets the weather to a specific type.
func (s *System) ForceWeather(w WeatherType, duration int) {
	s.currentWeather = w
	s.weatherTurns = duration
}

// ApplyWeatherEffects applies weather effects to game state.
func (s *System) ApplyWeatherEffects(res *resources.Resources, party *crew.Party, v *vessel.Vessel) WeatherResult {
	effect := GetWeatherEffect(s.currentWeather)
	result := WeatherResult{
		Weather:     s.currentWeather,
		WeatherName: WeatherName(s.currentWeather, s.genre),
	}

	// Apply morale change
	if effect.MoraleModifier != 0 {
		res.Add(resources.ResourceMorale, effect.MoraleModifier)
		result.MoraleChange = effect.MoraleModifier
	}

	// Apply crew health damage
	if effect.CrewHealthDamage > 0 {
		deaths := party.ApplyDamageToAll(effect.CrewHealthDamage)
		result.CrewDamage = effect.CrewHealthDamage
		result.Deaths = deaths
	}

	// Apply vessel damage
	if effect.VesselDamagePerTurn > 0 {
		v.TakeDamage(effect.VesselDamagePerTurn)
		result.VesselDamage = effect.VesselDamagePerTurn
	}

	// Resource consumption modifiers are returned but applied elsewhere
	result.WaterModifier = effect.WaterConsumptionMod
	result.FuelModifier = effect.FuelConsumptionMod
	result.MovementMultiplier = effect.MovementCostMultiplier
	result.VisibilityPenalty = effect.VisibilityPenalty

	return result
}

// WeatherResult contains the results of applying weather effects.
type WeatherResult struct {
	Weather            WeatherType
	WeatherName        string
	MoraleChange       float64
	CrewDamage         float64
	VesselDamage       float64
	Deaths             []string
	WaterModifier      float64
	FuelModifier       float64
	MovementMultiplier float64
	VisibilityPenalty  float64
}

// ApplyHazardEffects applies terrain hazard effects.
func (s *System) ApplyHazardEffects(hazard TerrainHazard, res *resources.Resources, party *crew.Party, v *vessel.Vessel) HazardResult {
	effect := GetHazardEffect(hazard)
	result := HazardResult{
		Hazard:     hazard,
		HazardName: HazardName(hazard, s.genre),
	}

	// Check for injury
	if effect.InjuryChance > 0 && s.gen.Float64() < effect.InjuryChance {
		result.InjuryOccurred = true
		result.InjuryDamage = effect.InjuryDamage

		// Injure a random living crew member
		living := party.Living()
		if len(living) > 0 {
			victim := seed.Choice(s.gen, living)
			if victim.TakeDamage(effect.InjuryDamage) {
				result.Deaths = append(result.Deaths, victim.Name)
			}
			result.InjuredCrew = victim.Name
		}
	}

	// Apply water cost
	if effect.WaterCostIncrease != 0 {
		res.Add(resources.ResourceWater, -effect.WaterCostIncrease)
		result.WaterCost = effect.WaterCostIncrease
	}

	// Check for loot
	if effect.LootChance > 0 && s.gen.Float64() < effect.LootChance {
		result.LootFound = true
		result.LootAmount = s.generateLoot()
	}

	// Check for danger encounter
	if effect.DangerChance > 0 && s.gen.Float64() < effect.DangerChance {
		result.DangerTriggered = true
	}

	// Return fuel and movement costs for application elsewhere
	result.FuelMultiplier = effect.FuelCostMultiplier
	result.MovementPenalty = effect.MovementPenalty

	return result
}

// generateLoot creates random loot rewards.
func (s *System) generateLoot() Loot {
	loot := Loot{}

	// Random resource rewards
	switch s.gen.Intn(5) {
	case 0:
		loot.Food = 5 + float64(s.gen.Intn(10))
	case 1:
		loot.Water = 5 + float64(s.gen.Intn(10))
	case 2:
		loot.Fuel = 5 + float64(s.gen.Intn(10))
	case 3:
		loot.Medicine = 2 + float64(s.gen.Intn(5))
	case 4:
		loot.Currency = 10 + float64(s.gen.Intn(20))
	}

	return loot
}

// Loot represents rewards found at hazard locations.
type Loot struct {
	Food     float64
	Water    float64
	Fuel     float64
	Medicine float64
	Currency float64
}

// ApplyLoot applies loot rewards to resources.
func (l *Loot) ApplyLoot(res *resources.Resources) {
	if l.Food > 0 {
		res.Add(resources.ResourceFood, l.Food)
	}
	if l.Water > 0 {
		res.Add(resources.ResourceWater, l.Water)
	}
	if l.Fuel > 0 {
		res.Add(resources.ResourceFuel, l.Fuel)
	}
	if l.Medicine > 0 {
		res.Add(resources.ResourceMedicine, l.Medicine)
	}
	if l.Currency > 0 {
		res.Add(resources.ResourceCurrency, l.Currency)
	}
}

// HazardResult contains the results of applying hazard effects.
type HazardResult struct {
	Hazard          TerrainHazard
	HazardName      string
	InjuryOccurred  bool
	InjuryDamage    float64
	InjuredCrew     string
	Deaths          []string
	WaterCost       float64
	FuelMultiplier  float64
	MovementPenalty float64
	LootFound       bool
	LootAmount      Loot
	DangerTriggered bool
}

// GetMovementCost calculates total movement cost considering weather and hazard.
func (s *System) GetMovementCost(baseCost float64, hazard TerrainHazard) float64 {
	weatherEffect := GetWeatherEffect(s.currentWeather)
	hazardEffect := GetHazardEffect(hazard)

	cost := baseCost
	cost *= weatherEffect.MovementCostMultiplier
	cost *= hazardEffect.FuelCostMultiplier
	cost += hazardEffect.MovementPenalty

	return cost
}

// GetVisibility returns current visibility (0=full, 1=blind).
func (s *System) GetVisibility() float64 {
	effect := GetWeatherEffect(s.currentWeather)
	return effect.VisibilityPenalty
}

// GetEncounterChanceModifier returns the encounter chance modifier.
func (s *System) GetEncounterChanceModifier() float64 {
	effect := GetWeatherEffect(s.currentWeather)
	return effect.EncounterChanceMod
}

// IsDangerous returns true if current weather is hazardous.
func (s *System) IsDangerous() bool {
	return s.currentWeather != WeatherClear
}

// GetWeatherSeverity returns 0-1 severity rating.
func (s *System) GetWeatherSeverity() float64 {
	effect := GetWeatherEffect(s.currentWeather)

	// Calculate composite severity
	severity := 0.0
	severity += effect.CrewHealthDamage / 10.0
	severity += effect.VesselDamagePerTurn / 10.0
	severity += effect.VisibilityPenalty
	severity += (effect.MovementCostMultiplier - 1.0) / 2.0

	if severity > 1.0 {
		severity = 1.0
	}
	return severity
}
