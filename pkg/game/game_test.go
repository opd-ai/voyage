//go:build headless

package game

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/world"
	"github.com/opd-ai/voyage/pkg/resources"
	"github.com/opd-ai/voyage/pkg/vessel"
)

func TestTimeManager(t *testing.T) {
	tm := NewTimeManager()

	if tm.Turn() != 0 {
		t.Errorf("initial turn = %d, want 0", tm.Turn())
	}
	if tm.Day() != 1 {
		t.Errorf("initial day = %d, want 1", tm.Day())
	}

	// Advance through a day
	for i := 0; i < 4; i++ {
		tm.Advance()
	}

	if tm.Day() != 2 {
		t.Errorf("day after 4 turns = %d, want 2", tm.Day())
	}

	// Check night detection
	tm.SetTurn(3) // Last turn of day 1
	if !tm.IsNight() {
		t.Error("turn 3 should be night")
	}

	tm.SetTurn(0)
	if tm.IsNight() {
		t.Error("turn 0 should not be night")
	}
}

func TestTimePhases(t *testing.T) {
	tm := NewTimeManager()

	tm.SetTurn(0)
	if tm.PhaseOfDay() != "Dawn" {
		t.Errorf("turn 0 phase = %s, want Dawn", tm.PhaseOfDay())
	}

	tm.SetTurn(1)
	if tm.PhaseOfDay() != "Morning" {
		t.Errorf("turn 1 phase = %s, want Morning", tm.PhaseOfDay())
	}

	tm.SetTurn(2)
	if tm.PhaseOfDay() != "Afternoon" {
		t.Errorf("turn 2 phase = %s, want Afternoon", tm.PhaseOfDay())
	}

	tm.SetTurn(3)
	if tm.PhaseOfDay() != "Night" {
		t.Errorf("turn 3 phase = %s, want Night", tm.PhaseOfDay())
	}
}

func TestMovementManager(t *testing.T) {
	mm := NewMovementManager()
	v := vessel.NewVessel(vessel.VesselMedium, engine.GenreFantasy)
	res := resources.NewResources(engine.GenreFantasy)

	// Get plains terrain
	plainsTerrain := world.DefaultTerrainInfo(world.TerrainPlains, engine.GenreFantasy)

	// Check can move
	if !mm.CanMove(plainsTerrain, v, res) {
		t.Error("should be able to move with full resources")
	}

	// Calculate cost
	fuelCost, timeCost := mm.CalculateMoveCost(plainsTerrain, v)
	if fuelCost <= 0 {
		t.Error("fuel cost should be positive")
	}
	if timeCost != 1 {
		t.Errorf("plains time cost = %d, want 1", timeCost)
	}

	// Execute move
	initialFuel := res.Get(resources.ResourceFuel)
	result := mm.Move(plainsTerrain, v, res)
	if !result.Success {
		t.Error("move should succeed")
	}
	if res.Get(resources.ResourceFuel) >= initialFuel {
		t.Error("fuel should be consumed")
	}
}

func TestMovementTerrainCosts(t *testing.T) {
	mm := NewMovementManager()
	v := vessel.NewVessel(vessel.VesselMedium, engine.GenreFantasy)

	terrains := []struct {
		terrain world.TerrainType
		minCost int
	}{
		{world.TerrainPlains, 1},
		{world.TerrainForest, 2},
		{world.TerrainMountain, 3},
	}

	for _, tt := range terrains {
		info := world.DefaultTerrainInfo(tt.terrain, engine.GenreFantasy)
		_, timeCost := mm.CalculateMoveCost(info, v)
		if timeCost < tt.minCost {
			t.Errorf("terrain %d time cost = %d, want >= %d", tt.terrain, timeCost, tt.minCost)
		}
	}
}

func TestRestManager(t *testing.T) {
	rm := NewRestManager()
	res := resources.NewResources(engine.GenreFantasy)

	// Create party with some injured members
	gen := crew.NewGenerator(12345, engine.GenreFantasy)
	party := crew.NewParty(engine.GenreFantasy, 4)
	for i := 0; i < 3; i++ {
		member := gen.Generate()
		member.Health = 50    // Set deterministic health
		member.TakeDamage(20) // Now at 30 health
		party.Add(member)
	}

	// Check can rest
	if !rm.CanRest(res, party.LivingCount()) {
		t.Error("should be able to rest with full resources")
	}

	initialMorale := res.Get(resources.ResourceMorale)
	initialFood := res.Get(resources.ResourceFood)

	// Record health before rest
	healthBefore := make(map[int]float64)
	for _, m := range party.Living() {
		healthBefore[m.ID] = m.Health
	}

	// Rest
	result := rm.Rest(res, party)
	if result.MoraleRecovered <= 0 {
		t.Error("should recover morale from rest")
	}
	if res.Get(resources.ResourceMorale) <= initialMorale {
		t.Error("morale should increase after rest")
	}
	if res.Get(resources.ResourceFood) >= initialFood {
		t.Error("food should be consumed during rest")
	}

	// Check crew healed
	for _, member := range party.Living() {
		if member.Health <= healthBefore[member.ID] {
			t.Error("crew should be healed after rest")
		}
	}
}

func TestRestWithInsufficientResources(t *testing.T) {
	rm := NewRestManager()
	res := resources.NewResources(engine.GenreFantasy)

	gen := crew.NewGenerator(12345, engine.GenreFantasy)
	party := crew.NewParty(engine.GenreFantasy, 4)
	party.Add(gen.Generate())

	// Deplete food
	res.Set(resources.ResourceFood, 0)

	if rm.CanRest(res, party.LivingCount()) {
		t.Error("should not be able to rest without food")
	}

	result := rm.Rest(res, party)
	if result.MoraleRecovered != 0 {
		t.Error("failed rest should not recover morale")
	}
}

func TestCampRest(t *testing.T) {
	rm := NewRestManager()
	res := resources.NewResources(engine.GenreFantasy)

	gen := crew.NewGenerator(12345, engine.GenreFantasy)
	party := crew.NewParty(engine.GenreFantasy, 4)
	party.Add(gen.Generate())

	result := rm.CampRest(res, party)
	if result.TurnsSpent != 2 {
		t.Errorf("camp rest turns = %d, want 2", result.TurnsSpent)
	}
	if result.MoraleRecovered <= 10 {
		t.Error("camp rest should recover more morale than regular rest")
	}
}

func TestTerrainConsumption(t *testing.T) {
	mm := NewMovementManager()
	res := resources.NewResources(engine.GenreFantasy)

	// Desert terrain increases water consumption
	desertTerrain := world.DefaultTerrainInfo(world.TerrainDesert, engine.GenreFantasy)

	initialWater := res.Get(resources.ResourceWater)
	mm.ApplyTerrainConsumption(desertTerrain, res, 3)

	// Desert should consume more water than normal
	waterConsumed := initialWater - res.Get(resources.ResourceWater)
	if waterConsumed <= 0 {
		t.Error("terrain should consume water")
	}
}

func TestSeason(t *testing.T) {
	tm := NewTimeManager()

	// At start, should be spring
	if tm.Season() != SeasonSpring {
		t.Errorf("initial season = %s, want Spring", SeasonName(tm.Season()))
	}

	// Advance to summer (20 days * 4 turns = 80 turns)
	tm.SetTurn(80)
	if tm.Season() != SeasonSummer {
		t.Errorf("season at day 21 = %s, want Summer", SeasonName(tm.Season()))
	}

	// Advance to autumn (40 days * 4 turns = 160 turns)
	tm.SetTurn(160)
	if tm.Season() != SeasonAutumn {
		t.Errorf("season at day 41 = %s, want Autumn", SeasonName(tm.Season()))
	}

	// Advance to winter (60 days * 4 turns = 240 turns)
	tm.SetTurn(240)
	if tm.Season() != SeasonWinter {
		t.Errorf("season at day 61 = %s, want Winter", SeasonName(tm.Season()))
	}

	// Verify year wraps back to spring (80 days * 4 turns = 320 turns)
	tm.SetTurn(320)
	if tm.Season() != SeasonSpring {
		t.Errorf("season at day 81 = %s, want Spring (year 2)", SeasonName(tm.Season()))
	}
}

func TestSeasonDayTracking(t *testing.T) {
	tm := NewTimeManager()

	// Day 1 of spring
	tm.SetTurn(0)
	if tm.DayInSeason() != 1 {
		t.Errorf("day in season at turn 0 = %d, want 1", tm.DayInSeason())
	}

	// Day 10 of spring (36 turns = 9 days, so turn 36 is day 10)
	tm.SetTurn(36)
	if tm.DayInSeason() != 10 {
		t.Errorf("day in season at turn 36 = %d, want 10", tm.DayInSeason())
	}

	// Day 1 of summer
	tm.SetTurn(80) // Day 21 overall, day 1 of summer
	if tm.DayInSeason() != 1 {
		t.Errorf("day in season at turn 80 = %d, want 1", tm.DayInSeason())
	}
}

func TestSeasonModifiers(t *testing.T) {
	// Test winter has highest costs
	winterCost, winterHazard := SeasonModifiers(SeasonWinter)
	summerCost, summerHazard := SeasonModifiers(SeasonSummer)

	if winterCost <= summerCost {
		t.Errorf("winter cost (%f) should be > summer cost (%f)", winterCost, summerCost)
	}
	if winterHazard <= summerHazard {
		t.Errorf("winter hazard (%f) should be > summer hazard (%f)", winterHazard, winterHazard)
	}

	// Test autumn has lowest costs
	autumnCost, _ := SeasonModifiers(SeasonAutumn)
	springCost, _ := SeasonModifiers(SeasonSpring)

	if autumnCost >= springCost {
		t.Errorf("autumn cost (%f) should be < spring cost (%f)", autumnCost, springCost)
	}
}

func TestDaysUntilSeasonChange(t *testing.T) {
	tm := NewTimeManager()

	// At start of spring (day 1), 20 days until summer
	tm.SetTurn(0)
	if tm.DaysUntilSeasonChange() != 20 {
		t.Errorf("days until season change at start = %d, want 20", tm.DaysUntilSeasonChange())
	}

	// At day 15 of spring, 6 days until summer
	tm.SetTurn(56) // Day 15
	if tm.DaysUntilSeasonChange() != 6 {
		t.Errorf("days until season change at day 15 = %d, want 6", tm.DaysUntilSeasonChange())
	}
}

func TestYearTracking(t *testing.T) {
	tm := NewTimeManager()

	// Year 1 at start
	tm.SetTurn(0)
	if tm.Year() != 1 {
		t.Errorf("year at turn 0 = %d, want 1", tm.Year())
	}

	// Year 1 at day 80 (end of first year)
	tm.SetTurn(316) // Day 80
	if tm.Year() != 1 {
		t.Errorf("year at day 80 = %d, want 1", tm.Year())
	}

	// Year 2 at day 81
	tm.SetTurn(320) // Day 81
	if tm.Year() != 2 {
		t.Errorf("year at day 81 = %d, want 2", tm.Year())
	}
}

func TestResourceCostModifier(t *testing.T) {
	tm := NewTimeManager()

	// Spring should be baseline
	tm.SetTurn(0)
	if tm.ResourceCostModifier() != 1.0 {
		t.Errorf("spring resource modifier = %f, want 1.0", tm.ResourceCostModifier())
	}

	// Winter should be higher
	tm.SetTurn(240) // Winter
	if tm.ResourceCostModifier() <= 1.0 {
		t.Errorf("winter resource modifier = %f, want > 1.0", tm.ResourceCostModifier())
	}
}

func TestSeasonProgress(t *testing.T) {
	tm := NewTimeManager()

	// Start of season should be 0%
	tm.SetTurn(0)
	if tm.SeasonProgress() != 0.0 {
		t.Errorf("season progress at start = %f, want 0.0", tm.SeasonProgress())
	}

	// Middle of season should be ~50%
	tm.SetTurn(40) // Day 11, about halfway through spring
	progress := tm.SeasonProgress()
	if progress < 0.4 || progress > 0.6 {
		t.Errorf("season progress at day 11 = %f, want ~0.5", progress)
	}
}

func TestForageManager(t *testing.T) {
	fm := NewForageManager(12345, engine.GenreFantasy)

	// Test action name
	if fm.ActionName() != "Forage" {
		t.Errorf("fantasy action name = %s, want Forage", fm.ActionName())
	}

	fm.SetGenre(engine.GenreScifi)
	if fm.ActionName() != "Salvage" {
		t.Errorf("scifi action name = %s, want Salvage", fm.ActionName())
	}

	fm.SetGenre(engine.GenreHorror)
	if fm.ActionName() != "Scavenge" {
		t.Errorf("horror action name = %s, want Scavenge", fm.ActionName())
	}
}

func TestForageCanForage(t *testing.T) {
	fm := NewForageManager(12345, engine.GenreFantasy)

	// Forest tile should be forageable
	forestTile := &world.Tile{X: 0, Y: 0, Terrain: world.TerrainForest}
	if !fm.CanForage(forestTile) {
		t.Error("should be able to forage in forest")
	}

	// Ruin tile should be forageable
	ruinTile := &world.Tile{X: 0, Y: 0, Terrain: world.TerrainRuin}
	if !fm.CanForage(ruinTile) {
		t.Error("should be able to forage in ruins")
	}

	// River tile should not be forageable
	riverTile := &world.Tile{X: 0, Y: 0, Terrain: world.TerrainRiver}
	if fm.CanForage(riverTile) {
		t.Error("should not be able to forage in river")
	}

	// Nil tile should not be forageable
	if fm.CanForage(nil) {
		t.Error("should not be able to forage nil tile")
	}
}

func TestForageDeterminism(t *testing.T) {
	fm1 := NewForageManager(12345, engine.GenreFantasy)
	fm2 := NewForageManager(12345, engine.GenreFantasy)

	tile := &world.Tile{X: 5, Y: 10, Terrain: world.TerrainForest}

	result1 := fm1.Forage(tile, 0)
	result2 := fm2.Forage(tile, 0)

	// Same seed and position should produce same outcome
	if result1.Outcome != result2.Outcome {
		t.Errorf("outcome mismatch: %d vs %d", result1.Outcome, result2.Outcome)
	}
	if result1.FoodGain != result2.FoodGain {
		t.Errorf("food gain mismatch: %f vs %f", result1.FoodGain, result2.FoodGain)
	}
}

func TestForageDiminishingReturns(t *testing.T) {
	fm := NewForageManager(12345, engine.GenreFantasy)

	tile := &world.Tile{X: 0, Y: 0, Terrain: world.TerrainForest}

	// First forage
	result1 := fm.Forage(tile, 0)
	gain1 := result1.FoodGain + result1.WaterGain + result1.FuelGain + result1.MedsGain

	// Second forage at same tile
	result2 := fm.Forage(tile, 1)
	gain2 := result2.FoodGain + result2.WaterGain + result2.FuelGain + result2.MedsGain

	// Forage count should increase
	if fm.GetForageCount(0, 0) != 2 {
		t.Errorf("forage count = %d, want 2", fm.GetForageCount(0, 0))
	}

	// Note: We can't guarantee second is less due to RNG, but we can test the mechanism
	// The modifier itself can be tested
	mod0 := fm.calculateYieldModifier(0)
	mod1 := fm.calculateYieldModifier(1)
	mod2 := fm.calculateYieldModifier(2)

	if mod0 != 1.0 {
		t.Errorf("first forage modifier = %f, want 1.0", mod0)
	}
	if mod1 >= mod0 {
		t.Error("second forage modifier should be less than first")
	}
	if mod2 >= mod1 {
		t.Error("third forage modifier should be less than second")
	}

	_ = gain1
	_ = gain2
}

func TestForageResetTile(t *testing.T) {
	fm := NewForageManager(12345, engine.GenreFantasy)
	tile := &world.Tile{X: 3, Y: 7, Terrain: world.TerrainForest}

	fm.Forage(tile, 0)
	fm.Forage(tile, 1)

	if fm.GetForageCount(3, 7) != 2 {
		t.Errorf("forage count = %d, want 2", fm.GetForageCount(3, 7))
	}

	fm.ResetTile(3, 7)
	if fm.GetForageCount(3, 7) != 0 {
		t.Errorf("forage count after reset = %d, want 0", fm.GetForageCount(3, 7))
	}
}

func TestForageApplyResult(t *testing.T) {
	fm := NewForageManager(12345, engine.GenreFantasy)
	res := resources.NewResources(engine.GenreFantasy)

	initialFood := res.Get(resources.ResourceFood)

	result := &ForageResult{
		Outcome:  ForageFood,
		FoodGain: 15.0,
	}

	fm.ApplyResult(result, res)

	if res.Get(resources.ResourceFood) != initialFood+15.0 {
		t.Errorf("food after apply = %f, want %f", res.Get(resources.ResourceFood), initialFood+15.0)
	}
}

func TestForageNonForageableTile(t *testing.T) {
	fm := NewForageManager(12345, engine.GenreFantasy)
	tile := &world.Tile{X: 0, Y: 0, Terrain: world.TerrainRiver}

	result := fm.Forage(tile, 0)
	if result.Outcome != ForageNothing {
		t.Errorf("non-forageable tile outcome = %d, want ForageNothing", result.Outcome)
	}
	if result.TurnsSpent != 0 {
		t.Errorf("non-forageable turns = %d, want 0", result.TurnsSpent)
	}
}
