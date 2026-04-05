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

// TestNewConditionCheckerHeadless tests condition checker creation.
func TestNewConditionCheckerHeadless(t *testing.T) {
	cc := NewConditionChecker()
	if cc == nil {
		t.Fatal("NewConditionChecker returned nil")
	}
}

// TestConditionCheckerSetPositionHeadless tests position setting.
func TestConditionCheckerSetPositionHeadless(t *testing.T) {
	cc := NewConditionChecker()
	cc.SetVesselPosition(10, 20)
	if cc.vesselX != 10 || cc.vesselY != 20 {
		t.Errorf("expected position (10, 20), got (%d, %d)", cc.vesselX, cc.vesselY)
	}
}

// TestCheckWinReachedDestinationHeadless tests win condition at destination.
func TestCheckWinReachedDestinationHeadless(t *testing.T) {
	cc := NewConditionChecker()

	// Create a world map
	gen := world.NewGenerator(12345, engine.GenreFantasy)
	worldMap := gen.Generate(50, 50)

	// Create party with living members
	party := crew.NewParty(engine.GenreFantasy, 4)
	party.Add(crew.NewCrewMember(1, "Test", crew.TraitBrave, crew.SkillNone))

	// Position at origin - should not win
	cc.SetVesselPosition(worldMap.Origin.X, worldMap.Origin.Y)
	won, _ := cc.CheckWin(worldMap, party)
	if won {
		t.Error("should not win at origin")
	}

	// Position at destination - should win
	cc.SetVesselPosition(worldMap.Destination.X, worldMap.Destination.Y)
	won, cond := cc.CheckWin(worldMap, party)
	if !won {
		t.Error("should win at destination with living crew")
	}
	if cond != WinReachedDestination {
		t.Errorf("expected WinReachedDestination, got %v", cond)
	}
}

// TestCheckWinRequiresLivingCrewHeadless tests that winning requires living crew.
func TestCheckWinRequiresLivingCrewHeadless(t *testing.T) {
	cc := NewConditionChecker()

	gen := world.NewGenerator(12345, engine.GenreFantasy)
	worldMap := gen.Generate(50, 50)

	// Create party with dead members only
	party := crew.NewParty(engine.GenreFantasy, 4)
	member := crew.NewCrewMember(1, "Test", crew.TraitBrave, crew.SkillNone)
	member.TakeDamage(200) // Kill the member
	party.Add(member)

	// Position at destination - should NOT win without living crew
	cc.SetVesselPosition(worldMap.Destination.X, worldMap.Destination.Y)
	won, _ := cc.CheckWin(worldMap, party)
	if won {
		t.Error("should not win without living crew")
	}
}

// TestCheckWinNilInputs tests CheckWin with nil inputs.
func TestCheckWinNilInputs(t *testing.T) {
	cc := NewConditionChecker()

	// Nil world map should return false
	won, _ := cc.CheckWin(nil, nil)
	if won {
		t.Error("should not win with nil inputs")
	}

	party := crew.NewParty(engine.GenreFantasy, 4)
	won, _ = cc.CheckWin(nil, party)
	if won {
		t.Error("should not win with nil world map")
	}
}

// TestCheckLoseAllCrewDeadHeadless tests loss condition for crew death.
func TestCheckLoseAllCrewDeadHeadless(t *testing.T) {
	cc := NewConditionChecker()

	party := crew.NewParty(engine.GenreFantasy, 4)
	v := vessel.NewVessel(vessel.VesselMedium, engine.GenreFantasy)
	res := resources.NewResources(engine.GenreFantasy)

	// Empty party should trigger loss
	lost, cond := cc.CheckLose(party, v, res)
	if !lost {
		t.Error("should lose with no living crew")
	}
	if cond != LoseAllCrewDead {
		t.Errorf("expected LoseAllCrewDead, got %v", cond)
	}
}

// TestCheckLoseVesselDestroyedHeadless tests loss condition for vessel destruction.
func TestCheckLoseVesselDestroyedHeadless(t *testing.T) {
	cc := NewConditionChecker()

	party := crew.NewParty(engine.GenreFantasy, 4)
	party.Add(crew.NewCrewMember(1, "Test", crew.TraitBrave, crew.SkillNone))

	v := vessel.NewVessel(vessel.VesselMedium, engine.GenreFantasy)
	v.TakeDamage(999) // Destroy vessel

	res := resources.NewResources(engine.GenreFantasy)

	lost, cond := cc.CheckLose(party, v, res)
	if !lost {
		t.Error("should lose with destroyed vessel")
	}
	if cond != LoseVesselDestroyed {
		t.Errorf("expected LoseVesselDestroyed, got %v", cond)
	}
}

// TestCheckLoseMoraleZeroHeadless tests loss condition for zero morale.
func TestCheckLoseMoraleZeroHeadless(t *testing.T) {
	cc := NewConditionChecker()

	party := crew.NewParty(engine.GenreFantasy, 4)
	party.Add(crew.NewCrewMember(1, "Test", crew.TraitBrave, crew.SkillNone))

	v := vessel.NewVessel(vessel.VesselMedium, engine.GenreFantasy)
	res := resources.NewResources(engine.GenreFantasy)
	res.Set(resources.ResourceMorale, 0)

	lost, cond := cc.CheckLose(party, v, res)
	if !lost {
		t.Error("should lose with zero morale")
	}
	if cond != LoseMoraleZero {
		t.Errorf("expected LoseMoraleZero, got %v", cond)
	}
}

// TestCheckLoseStarvationHeadless tests loss condition for starvation.
func TestCheckLoseStarvationHeadless(t *testing.T) {
	cc := NewConditionChecker()

	party := crew.NewParty(engine.GenreFantasy, 4)
	party.Add(crew.NewCrewMember(1, "Test", crew.TraitBrave, crew.SkillNone))

	v := vessel.NewVessel(vessel.VesselMedium, engine.GenreFantasy)
	res := resources.NewResources(engine.GenreFantasy)
	res.Set(resources.ResourceFood, 0)
	res.Set(resources.ResourceWater, 0)
	res.Set(resources.ResourceMorale, 50) // Keep morale to test starvation

	lost, cond := cc.CheckLose(party, v, res)
	if !lost {
		t.Error("should lose with no food and water")
	}
	if cond != LoseStarvation {
		t.Errorf("expected LoseStarvation, got %v", cond)
	}
}

// TestCheckNoLoseHeadless tests no loss condition when healthy.
func TestCheckNoLoseHeadless(t *testing.T) {
	cc := NewConditionChecker()

	party := crew.NewParty(engine.GenreFantasy, 4)
	party.Add(crew.NewCrewMember(1, "Test", crew.TraitBrave, crew.SkillNone))

	v := vessel.NewVessel(vessel.VesselMedium, engine.GenreFantasy)
	res := resources.NewResources(engine.GenreFantasy)

	lost, cond := cc.CheckLose(party, v, res)
	if lost {
		t.Error("should not lose with healthy party, vessel, and resources")
	}
	if cond != LoseNone {
		t.Errorf("expected LoseNone, got %v", cond)
	}
}

// TestCheckLoseNilInputs tests CheckLose with nil inputs.
func TestCheckLoseNilInputs(t *testing.T) {
	cc := NewConditionChecker()

	// Nil party should not panic
	lost, _ := cc.CheckLose(nil, nil, nil)
	if lost {
		t.Error("should not lose with nil inputs")
	}
}

// TestConditionNamesHeadless tests condition name functions.
func TestConditionNamesHeadless(t *testing.T) {
	winName := WinConditionName(WinReachedDestination)
	if winName == "" {
		t.Error("WinConditionName should return non-empty string")
	}
	if winName != "Reached Destination" {
		t.Errorf("expected 'Reached Destination', got '%s'", winName)
	}

	loseNames := []struct {
		cond LoseCondition
		want string
	}{
		{LoseAllCrewDead, "All Crew Lost"},
		{LoseVesselDestroyed, "Vessel Destroyed"},
		{LoseMoraleZero, "Crew Deserted"},
		{LoseStarvation, "Starved"},
		{LoseNone, "Unknown"},
	}

	for _, tt := range loseNames {
		name := LoseConditionName(tt.cond)
		if name != tt.want {
			t.Errorf("LoseConditionName(%v) = %s, want %s", tt.cond, name, tt.want)
		}
	}
}

// TestConditionDescriptionsHeadless tests condition description functions.
func TestConditionDescriptionsHeadless(t *testing.T) {
	winDesc := WinConditionDescription(WinReachedDestination)
	if winDesc == "" {
		t.Error("WinConditionDescription should return non-empty string")
	}

	loseDescriptions := []LoseCondition{
		LoseAllCrewDead,
		LoseVesselDestroyed,
		LoseMoraleZero,
		LoseStarvation,
		LoseNone,
	}

	for _, lc := range loseDescriptions {
		desc := LoseConditionDescription(lc)
		if desc == "" {
			t.Errorf("LoseConditionDescription(%v) should return non-empty string", lc)
		}
	}
}

// TestWinConditionDefault tests default win condition handling.
func TestWinConditionDefault(t *testing.T) {
	// Test non-existent win condition (cast a number)
	name := WinConditionName(WinCondition(99))
	if name != "Victory" {
		t.Errorf("expected 'Victory' for unknown condition, got '%s'", name)
	}

	desc := WinConditionDescription(WinCondition(99))
	if desc == "" {
		t.Error("expected non-empty description for unknown condition")
	}
}

// TestGameStateConstants tests that GameState constants are distinct.
func TestGameStateConstants(t *testing.T) {
	states := []GameState{StateMenu, StatePlaying, StatePaused, StateGameOver}
	seen := make(map[GameState]bool)

	for _, s := range states {
		if seen[s] {
			t.Errorf("duplicate GameState value: %d", s)
		}
		seen[s] = true
	}
}

// TestWinConditionConstants tests that WinCondition constants are distinct.
func TestWinConditionConstants(t *testing.T) {
	conditions := []WinCondition{WinReachedDestination}
	seen := make(map[WinCondition]bool)

	for _, c := range conditions {
		if seen[c] {
			t.Errorf("duplicate WinCondition value: %d", c)
		}
		seen[c] = true
	}
}

// TestLoseConditionConstants tests that LoseCondition constants are distinct.
func TestLoseConditionConstants(t *testing.T) {
	conditions := []LoseCondition{
		LoseNone,
		LoseAllCrewDead,
		LoseVesselDestroyed,
		LoseMoraleZero,
		LoseStarvation,
	}
	seen := make(map[LoseCondition]bool)

	for _, c := range conditions {
		if seen[c] {
			t.Errorf("duplicate LoseCondition value: %d", c)
		}
		seen[c] = true
	}
}
