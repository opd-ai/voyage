package game

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/world"
	"github.com/opd-ai/voyage/pkg/resources"
	"github.com/opd-ai/voyage/pkg/vessel"
)

func TestNewConditionChecker(t *testing.T) {
	cc := NewConditionChecker()
	if cc == nil {
		t.Fatal("NewConditionChecker returned nil")
	}
}

func TestConditionCheckerSetPosition(t *testing.T) {
	cc := NewConditionChecker()
	cc.SetVesselPosition(10, 20)
	if cc.vesselX != 10 || cc.vesselY != 20 {
		t.Errorf("expected position (10, 20), got (%d, %d)", cc.vesselX, cc.vesselY)
	}
}

func TestCheckWinReachedDestination(t *testing.T) {
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

func TestCheckWinRequiresLivingCrew(t *testing.T) {
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

func TestCheckLoseAllCrewDead(t *testing.T) {
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

func TestCheckLoseVesselDestroyed(t *testing.T) {
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

func TestCheckLoseMoraleZero(t *testing.T) {
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

func TestCheckLoseStarvation(t *testing.T) {
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

func TestCheckNoLose(t *testing.T) {
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

func TestConditionNames(t *testing.T) {
	winName := WinConditionName(WinReachedDestination)
	if winName == "" {
		t.Error("WinConditionName should return non-empty string")
	}

	loseNames := []LoseCondition{
		LoseAllCrewDead,
		LoseVesselDestroyed,
		LoseMoraleZero,
		LoseStarvation,
	}

	for _, lc := range loseNames {
		name := LoseConditionName(lc)
		if name == "" {
			t.Errorf("LoseConditionName(%v) should return non-empty string", lc)
		}
	}
}

func TestConditionDescriptions(t *testing.T) {
	winDesc := WinConditionDescription(WinReachedDestination)
	if winDesc == "" {
		t.Error("WinConditionDescription should return non-empty string")
	}

	loseDescriptions := []LoseCondition{
		LoseAllCrewDead,
		LoseVesselDestroyed,
		LoseMoraleZero,
		LoseStarvation,
	}

	for _, lc := range loseDescriptions {
		desc := LoseConditionDescription(lc)
		if desc == "" {
			t.Errorf("LoseConditionDescription(%v) should return non-empty string", lc)
		}
	}
}
