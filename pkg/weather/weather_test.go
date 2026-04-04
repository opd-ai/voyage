package weather

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/resources"
	"github.com/opd-ai/voyage/pkg/vessel"
)

func TestAllWeatherTypes(t *testing.T) {
	types := AllWeatherTypes()
	if len(types) != 9 {
		t.Errorf("Expected 9 weather types, got %d", len(types))
	}
}

func TestHazardousWeatherTypes(t *testing.T) {
	types := HazardousWeatherTypes()
	if len(types) != 8 {
		t.Errorf("Expected 8 hazardous weather types, got %d", len(types))
	}

	// Verify clear is not included
	for _, w := range types {
		if w == WeatherClear {
			t.Error("Clear weather should not be in hazardous list")
		}
	}
}

func TestAllTerrainHazards(t *testing.T) {
	hazards := AllTerrainHazards()
	if len(hazards) != 8 {
		t.Errorf("Expected 8 terrain hazards, got %d", len(hazards))
	}
}

func TestWeatherName(t *testing.T) {
	tests := []struct {
		weather WeatherType
		genre   engine.GenreID
		want    string
	}{
		{WeatherClear, engine.GenreFantasy, "Clear Skies"},
		{WeatherStorm, engine.GenreScifi, "Ion Storm"},
		{WeatherBlizzard, engine.GenreHorror, "Blizzard"},
		{WeatherAcidRain, engine.GenreCyberpunk, "Acid Rain"},
		{WeatherDustStorm, engine.GenrePostapoc, "Dust Storm"},
	}

	for _, tt := range tests {
		got := WeatherName(tt.weather, tt.genre)
		if got != tt.want {
			t.Errorf("WeatherName(%d, %s) = %q, want %q", tt.weather, tt.genre, got, tt.want)
		}
	}
}

func TestWeatherDescription(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		for _, w := range AllWeatherTypes() {
			desc := WeatherDescription(w, genre)
			if desc == "" {
				t.Errorf("Missing description for weather %d in genre %s", w, genre)
			}
		}
	}
}

func TestHazardName(t *testing.T) {
	tests := []struct {
		hazard TerrainHazard
		genre  engine.GenreID
		want   string
	}{
		{HazardNone, engine.GenreFantasy, "Safe Passage"},
		{HazardMountainPass, engine.GenreScifi, "Asteroid Field"},
		{HazardRuin, engine.GenreHorror, "Abandoned Town"},
		{HazardRadiation, engine.GenreCyberpunk, "Rad Zone"},
		{HazardDesert, engine.GenrePostapoc, "Scorched Earth"},
	}

	for _, tt := range tests {
		got := HazardName(tt.hazard, tt.genre)
		if got != tt.want {
			t.Errorf("HazardName(%d, %s) = %q, want %q", tt.hazard, tt.genre, got, tt.want)
		}
	}
}

func TestGetWeatherEffect(t *testing.T) {
	effect := GetWeatherEffect(WeatherStorm)
	if effect.MovementCostMultiplier <= 1.0 {
		t.Error("Storm should increase movement cost")
	}
	if effect.CrewHealthDamage <= 0 {
		t.Error("Storm should cause crew damage")
	}

	clearEffect := GetWeatherEffect(WeatherClear)
	if clearEffect.MovementCostMultiplier != 1.0 {
		t.Error("Clear weather should have normal movement cost")
	}
}

func TestGetHazardEffect(t *testing.T) {
	effect := GetHazardEffect(HazardMountainPass)
	if effect.InjuryChance <= 0 {
		t.Error("Mountain pass should have injury chance")
	}
	if effect.FuelCostMultiplier <= 1.0 {
		t.Error("Mountain pass should increase fuel cost")
	}

	noneEffect := GetHazardEffect(HazardNone)
	if noneEffect.InjuryChance != 0 {
		t.Error("No hazard should have zero injury chance")
	}
}

func TestGenreWeatherSubset(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		subset := GenreWeatherSubset(genre)
		if len(subset) < 5 {
			t.Errorf("Genre %s should have at least 5 weather types, got %d", genre, len(subset))
		}

		// Verify clear is included
		hasClean := false
		for _, w := range subset {
			if w == WeatherClear {
				hasClean = true
				break
			}
		}
		if !hasClean {
			t.Errorf("Genre %s subset should include clear weather", genre)
		}
	}
}

func TestGenreHazardSubset(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		subset := GenreHazardSubset(genre)
		if len(subset) < 4 {
			t.Errorf("Genre %s should have at least 4 hazard types, got %d", genre, len(subset))
		}

		// Verify none is included
		hasNone := false
		for _, h := range subset {
			if h == HazardNone {
				hasNone = true
				break
			}
		}
		if !hasNone {
			t.Errorf("Genre %s subset should include no hazard", genre)
		}
	}
}

func TestNewSystem(t *testing.T) {
	sys := NewSystem(12345, engine.GenreFantasy)
	if sys.currentWeather != WeatherClear {
		t.Error("System should start with clear weather")
	}
	if sys.weatherTurns <= 0 {
		t.Error("System should have positive turns remaining")
	}
}

func TestSystemSetGenre(t *testing.T) {
	sys := NewSystem(12345, engine.GenreFantasy)
	sys.SetGenre(engine.GenreScifi)

	if sys.genre != engine.GenreScifi {
		t.Errorf("Expected genre scifi, got %s", sys.genre)
	}
}

func TestSystemAdvanceTurn(t *testing.T) {
	sys := NewSystem(12345, engine.GenreFantasy)

	// Advance until weather changes
	changed := false
	for i := 0; i < 20; i++ {
		if sys.AdvanceTurn() {
			changed = true
			break
		}
	}

	if !changed {
		t.Error("Weather should change eventually")
	}
}

func TestSystemForceWeather(t *testing.T) {
	sys := NewSystem(12345, engine.GenreFantasy)
	sys.ForceWeather(WeatherStorm, 5)

	if sys.CurrentWeather() != WeatherStorm {
		t.Errorf("Expected storm, got %d", sys.CurrentWeather())
	}
	if sys.TurnsRemaining() != 5 {
		t.Errorf("Expected 5 turns, got %d", sys.TurnsRemaining())
	}
}

func TestSystemApplyWeatherEffects(t *testing.T) {
	sys := NewSystem(12345, engine.GenreFantasy)
	sys.ForceWeather(WeatherStorm, 5)

	res := resources.NewResources(engine.GenreFantasy)
	party := crew.NewParty(engine.GenreFantasy, 4)
	crewGen := crew.NewGenerator(12345, engine.GenreFantasy)
	for i := 0; i < 3; i++ {
		party.Add(crewGen.Generate())
	}
	v := vessel.NewVessel(vessel.VesselMedium, engine.GenreFantasy)

	initialHealth := party.AverageHealth()
	initialIntegrity := v.Integrity()

	result := sys.ApplyWeatherEffects(res, party, v)

	if result.Weather != WeatherStorm {
		t.Error("Result should report storm")
	}
	if result.CrewDamage <= 0 {
		t.Error("Storm should cause crew damage")
	}
	if party.AverageHealth() >= initialHealth {
		t.Error("Crew health should decrease in storm")
	}
	if v.Integrity() >= initialIntegrity {
		t.Error("Vessel integrity should decrease in storm")
	}
}

func TestSystemApplyHazardEffects(t *testing.T) {
	sys := NewSystem(12345, engine.GenreFantasy)

	res := resources.NewResources(engine.GenreFantasy)
	party := crew.NewParty(engine.GenreFantasy, 4)
	crewGen := crew.NewGenerator(12345, engine.GenreFantasy)
	for i := 0; i < 3; i++ {
		party.Add(crewGen.Generate())
	}
	v := vessel.NewVessel(vessel.VesselMedium, engine.GenreFantasy)

	result := sys.ApplyHazardEffects(HazardMountainPass, res, party, v)

	if result.Hazard != HazardMountainPass {
		t.Error("Result should report mountain pass")
	}
	if result.FuelMultiplier <= 1.0 {
		t.Error("Mountain pass should increase fuel cost")
	}
}

func TestSystemGetMovementCost(t *testing.T) {
	sys := NewSystem(12345, engine.GenreFantasy)

	// Clear weather, no hazard
	cost1 := sys.GetMovementCost(10, HazardNone)
	if cost1 != 10 {
		t.Errorf("Expected base cost 10, got %f", cost1)
	}

	// Storm weather
	sys.ForceWeather(WeatherStorm, 5)
	cost2 := sys.GetMovementCost(10, HazardNone)
	if cost2 <= cost1 {
		t.Error("Storm should increase movement cost")
	}

	// Storm + hazard
	cost3 := sys.GetMovementCost(10, HazardMountainPass)
	if cost3 <= cost2 {
		t.Error("Hazard should further increase movement cost")
	}
}

func TestSystemGetVisibility(t *testing.T) {
	sys := NewSystem(12345, engine.GenreFantasy)

	// Clear weather
	vis1 := sys.GetVisibility()
	if vis1 != 0 {
		t.Errorf("Clear weather should have 0 visibility penalty, got %f", vis1)
	}

	// Fog
	sys.ForceWeather(WeatherFog, 5)
	vis2 := sys.GetVisibility()
	if vis2 <= vis1 {
		t.Error("Fog should increase visibility penalty")
	}
}

func TestSystemIsDangerous(t *testing.T) {
	sys := NewSystem(12345, engine.GenreFantasy)

	if sys.IsDangerous() {
		t.Error("Clear weather should not be dangerous")
	}

	sys.ForceWeather(WeatherStorm, 5)
	if !sys.IsDangerous() {
		t.Error("Storm should be dangerous")
	}
}

func TestSystemGetWeatherSeverity(t *testing.T) {
	sys := NewSystem(12345, engine.GenreFantasy)

	sev1 := sys.GetWeatherSeverity()
	if sev1 != 0 {
		t.Errorf("Clear weather severity should be 0, got %f", sev1)
	}

	sys.ForceWeather(WeatherMeteorShower, 5)
	sev2 := sys.GetWeatherSeverity()
	if sev2 <= sev1 {
		t.Error("Meteor shower should have higher severity")
	}
}

func TestLootApply(t *testing.T) {
	loot := Loot{
		Food:     10,
		Water:    5,
		Currency: 20,
	}

	res := resources.NewResources(engine.GenreFantasy)
	initialFood := res.Get(resources.ResourceFood)
	initialWater := res.Get(resources.ResourceWater)
	initialCurrency := res.Get(resources.ResourceCurrency)

	loot.ApplyLoot(res)

	if res.Get(resources.ResourceFood) != initialFood+10 {
		t.Error("Food should increase by loot amount")
	}
	if res.Get(resources.ResourceWater) != initialWater+5 {
		t.Error("Water should increase by loot amount")
	}
	if res.Get(resources.ResourceCurrency) != initialCurrency+20 {
		t.Error("Currency should increase by loot amount")
	}
}

func TestSystemDeterminism(t *testing.T) {
	sys1 := NewSystem(42, engine.GenreFantasy)
	sys2 := NewSystem(42, engine.GenreFantasy)

	// Advance both systems
	for i := 0; i < 10; i++ {
		sys1.AdvanceTurn()
		sys2.AdvanceTurn()
	}

	if sys1.CurrentWeather() != sys2.CurrentWeather() {
		t.Error("Same seed should produce same weather")
	}
}
