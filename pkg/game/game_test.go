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
