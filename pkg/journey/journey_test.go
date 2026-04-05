package journey

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewLeg(t *testing.T) {
	leg := NewLeg(1, "Test Leg", "Origin", "Destination", 100, DifficultyNormal, engine.GenreFantasy)

	if leg.ID != 1 {
		t.Errorf("expected ID 1, got %d", leg.ID)
	}
	if leg.Name != "Test Leg" {
		t.Error("name mismatch")
	}
	if leg.Distance != 100 {
		t.Error("distance mismatch")
	}
	if leg.EnemyStrength != 1.0 {
		t.Error("normal difficulty should have 1.0 enemy strength")
	}
	if leg.Started || leg.Completed {
		t.Error("should not be started or completed initially")
	}
}

func TestLegDifficultyMultipliers(t *testing.T) {
	testCases := []struct {
		difficulty  DifficultyLevel
		minStrength float64
		maxStrength float64
	}{
		{DifficultyEasy, 0.7, 0.9},
		{DifficultyNormal, 0.9, 1.1},
		{DifficultyHard, 1.2, 1.4},
		{DifficultyExtreme, 1.5, 1.7},
	}

	for _, tc := range testCases {
		leg := NewLeg(1, "Test", "A", "B", 100, tc.difficulty, engine.GenreFantasy)
		if leg.EnemyStrength < tc.minStrength || leg.EnemyStrength > tc.maxStrength {
			t.Errorf("difficulty %s: enemy strength %f out of expected range [%f, %f]",
				DifficultyName(tc.difficulty), leg.EnemyStrength, tc.minStrength, tc.maxStrength)
		}
	}
}

func TestLegStartComplete(t *testing.T) {
	leg := NewLeg(1, "Test", "A", "B", 100, DifficultyNormal, engine.GenreFantasy)

	leg.Start()
	if !leg.Started {
		t.Error("should be started")
	}

	leg.Complete(10, 5)
	if !leg.Completed {
		t.Error("should be completed")
	}
	if leg.DaysTaken != 10 {
		t.Error("days taken mismatch")
	}
	if leg.Survivors != 5 {
		t.Error("survivors mismatch")
	}
}

func TestLegSetGenre(t *testing.T) {
	leg := NewLeg(1, "Test", "A", "B", 100, DifficultyNormal, engine.GenreFantasy)
	leg.SetGenre(engine.GenreScifi)

	if leg.Genre != engine.GenreScifi {
		t.Error("genre should be updated")
	}
}

func TestDifficultyName(t *testing.T) {
	for _, d := range AllDifficultyLevels() {
		name := DifficultyName(d)
		if name == "Unknown" || name == "" {
			t.Errorf("difficulty %v should have a name", d)
		}
	}
}

func TestNewStopover(t *testing.T) {
	stopover := NewStopover(1, "Test Stop", "Description", 1, engine.GenreFantasy)

	if stopover.ID != 1 {
		t.Error("ID mismatch")
	}
	if stopover.Name != "Test Stop" {
		t.Error("name mismatch")
	}
	if stopover.AfterLeg != 1 {
		t.Error("after leg mismatch")
	}
	if stopover.Visited {
		t.Error("should not be visited initially")
	}
}

func TestStopoverServices(t *testing.T) {
	stopover := NewStopover(1, "Test", "Desc", 1, engine.GenreFantasy)

	stopover.AddService(ServiceTrading)
	stopover.AddService(ServiceRepairs)

	if !stopover.HasService(ServiceTrading) {
		t.Error("should have trading")
	}
	if !stopover.HasService(ServiceRepairs) {
		t.Error("should have repairs")
	}
	if stopover.HasService(ServiceHealing) {
		t.Error("should not have healing")
	}
}

func TestStopoverVisit(t *testing.T) {
	stopover := NewStopover(1, "Test", "Desc", 1, engine.GenreFantasy)
	stopover.Visit()

	if !stopover.Visited {
		t.Error("should be visited")
	}
}

func TestServiceNameAllGenres(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		for _, service := range AllStopoverServices() {
			name := ServiceName(service, genre)
			if name == "" {
				t.Errorf("service %v genre %s should have name", service, genre)
			}
		}
	}
}

func TestCampaignState(t *testing.T) {
	state := NewCampaignState()

	if state.Gold != 100 {
		t.Error("starting gold should be 100")
	}
	if state.CrewCount != 4 {
		t.Error("starting crew should be 4")
	}
	if state.CrewHealth != 1.0 {
		t.Error("starting health should be 1.0")
	}
}

func TestCampaignStateApplyResults(t *testing.T) {
	state := NewCampaignState()
	leg := NewLeg(1, "Test", "A", "B", 100, DifficultyNormal, engine.GenreFantasy)

	state.ApplyLegResults(leg, 3, 10, 50, 20)

	if state.TotalDistance != 100 {
		t.Error("distance should be applied")
	}
	if state.TotalDays != 10 {
		t.Error("days should be applied")
	}
	if state.TotalDeaths != 1 {
		t.Error("should have 1 death (4 crew - 3 survivors)")
	}
	if state.Gold != 150 {
		t.Error("gold should be 150 (100 + 50)")
	}
	if state.Food != 30 {
		t.Error("food should be 30 (50 - 20)")
	}
}

func TestCampaignStateAchievements(t *testing.T) {
	state := NewCampaignState()
	state.AddAchievement("First Journey")
	state.AddAchievement("Survivor")

	if len(state.Achievements) != 2 {
		t.Error("should have 2 achievements")
	}
}

func TestNewCampaign(t *testing.T) {
	campaign := NewCampaign("test_1", "Test Campaign", "Description", engine.GenreFantasy)

	if campaign.ID != "test_1" {
		t.Error("ID mismatch")
	}
	if campaign.Name != "Test Campaign" {
		t.Error("name mismatch")
	}
	if campaign.State == nil {
		t.Error("state should be initialized")
	}
	if campaign.CurrentLegIndex != 0 {
		t.Error("should start at leg 0")
	}
}

func TestCampaignAddLeg(t *testing.T) {
	campaign := NewCampaign("test", "Test", "Desc", engine.GenreFantasy)
	leg1 := NewLeg(1, "Leg 1", "A", "B", 100, DifficultyEasy, engine.GenreFantasy)
	leg2 := NewLeg(2, "Leg 2", "B", "C", 150, DifficultyNormal, engine.GenreFantasy)

	campaign.AddLeg(leg1)
	campaign.AddLeg(leg2)

	if campaign.LegCount() != 2 {
		t.Error("should have 2 legs")
	}
	if campaign.TotalDistance() != 250 {
		t.Error("total distance should be 250")
	}
}

func TestCampaignProgress(t *testing.T) {
	campaign := NewCampaign("test", "Test", "Desc", engine.GenreFantasy)
	campaign.AddLeg(NewLeg(1, "L1", "A", "B", 100, DifficultyEasy, engine.GenreFantasy))
	campaign.AddLeg(NewLeg(2, "L2", "B", "C", 100, DifficultyNormal, engine.GenreFantasy))

	if campaign.Progress() != 0 {
		t.Error("progress should be 0% at start")
	}

	campaign.CompleteLeg(10, 4)

	if campaign.Progress() != 50 {
		t.Error("progress should be 50% after completing 1 of 2 legs")
	}
}

func TestCampaignCompleteLeg(t *testing.T) {
	campaign := NewCampaign("test", "Test", "Desc", engine.GenreFantasy)
	campaign.AddLeg(NewLeg(1, "L1", "A", "B", 100, DifficultyEasy, engine.GenreFantasy))
	campaign.AddLeg(NewLeg(2, "L2", "B", "C", 100, DifficultyNormal, engine.GenreFantasy))

	// Complete first leg
	hasMore := campaign.CompleteLeg(10, 4)
	if !hasMore {
		t.Error("should have more legs")
	}
	if campaign.CurrentLegIndex != 1 {
		t.Error("should advance to leg index 1")
	}

	// Complete second leg
	hasMore = campaign.CompleteLeg(15, 3)
	if hasMore {
		t.Error("should not have more legs")
	}
	if !campaign.IsComplete() {
		t.Error("campaign should be complete")
	}
}

func TestCampaignSetGenre(t *testing.T) {
	campaign := NewCampaign("test", "Test", "Desc", engine.GenreFantasy)
	leg := NewLeg(1, "L1", "A", "B", 100, DifficultyEasy, engine.GenreFantasy)
	stopover := NewStopover(1, "S1", "Desc", 1, engine.GenreFantasy)

	campaign.AddLeg(leg)
	campaign.AddStopover(stopover)

	campaign.SetGenre(engine.GenreScifi)

	if campaign.Genre != engine.GenreScifi {
		t.Error("campaign genre should be updated")
	}
	if leg.Genre != engine.GenreScifi {
		t.Error("leg genre should be updated")
	}
	if stopover.Genre != engine.GenreScifi {
		t.Error("stopover genre should be updated")
	}
}

func TestCampaignGenreShifts(t *testing.T) {
	campaign := NewCampaign("test", "Test", "Desc", engine.GenreFantasy)
	leg1 := NewLeg(1, "L1", "A", "B", 100, DifficultyEasy, engine.GenreFantasy)
	leg2 := NewLeg(2, "L2", "B", "C", 100, DifficultyNormal, engine.GenreScifi)

	campaign.AddLeg(leg1)
	campaign.AddLeg(leg2)
	campaign.EnableGenreShifts()

	// SetGenre should NOT update individual legs when genre shifts enabled
	campaign.SetGenre(engine.GenreHorror)

	if leg1.Genre != engine.GenreFantasy {
		t.Error("leg1 genre should remain fantasy with genre shifts enabled")
	}
	if leg2.Genre != engine.GenreScifi {
		t.Error("leg2 genre should remain scifi with genre shifts enabled")
	}
}

func TestGenerator(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)

	campaign := g.GenerateCampaign(3)

	if campaign.LegCount() != 3 {
		t.Errorf("expected 3 legs, got %d", campaign.LegCount())
	}
	if len(campaign.Stopovers) != 2 {
		t.Errorf("expected 2 stopovers (between legs), got %d", len(campaign.Stopovers))
	}
	if campaign.Name == "" {
		t.Error("campaign should have a name")
	}
}

func TestGeneratorDeterminism(t *testing.T) {
	g1 := NewGenerator(12345, engine.GenreFantasy)
	g2 := NewGenerator(12345, engine.GenreFantasy)

	c1 := g1.GenerateCampaign(2)
	c2 := g2.GenerateCampaign(2)

	if c1.Name != c2.Name {
		t.Error("same seed should produce same campaign name")
	}
	if c1.Legs[0].Name != c2.Legs[0].Name {
		t.Error("same seed should produce same leg names")
	}
}

func TestGeneratorLegCount(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)

	// Test minimum clamping
	campaign := g.GenerateCampaign(1)
	if campaign.LegCount() != 2 {
		t.Error("minimum leg count should be 2")
	}

	// Test maximum clamping
	g2 := NewGenerator(12346, engine.GenreFantasy)
	campaign2 := g2.GenerateCampaign(10)
	if campaign2.LegCount() != 4 {
		t.Error("maximum leg count should be 4")
	}
}

func TestGeneratorAllGenres(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		g := NewGenerator(12345, genre)
		campaign := g.GenerateCampaign(2)

		if campaign.Name == "" {
			t.Errorf("genre %s: campaign should have name", genre)
		}
		if len(campaign.Legs) != 2 {
			t.Errorf("genre %s: should have 2 legs", genre)
		}

		for _, leg := range campaign.Legs {
			if leg.Description == "" {
				t.Errorf("genre %s: leg should have description", genre)
			}
			if leg.TerrainType == "" {
				t.Errorf("genre %s: leg should have terrain", genre)
			}
		}

		for _, stopover := range campaign.Stopovers {
			if stopover.Description == "" {
				t.Errorf("genre %s: stopover should have description", genre)
			}
			if len(stopover.Services) == 0 {
				t.Errorf("genre %s: stopover should have services", genre)
			}
		}
	}
}

func TestGeneratorEscalatingDifficulty(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)
	campaign := g.GenerateCampaign(4)

	// First leg should be easier than last
	if campaign.Legs[0].Difficulty >= campaign.Legs[3].Difficulty {
		t.Error("difficulty should escalate over legs")
	}

	// Distance should generally increase
	if campaign.Legs[0].Distance >= campaign.Legs[3].Distance {
		t.Error("distance should increase over legs")
	}
}

func TestGeneratorSetGenre(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)
	g.SetGenre(engine.GenreCyberpunk)

	campaign := g.GenerateCampaign(2)

	if campaign.Genre != engine.GenreCyberpunk {
		t.Error("campaign should have cyberpunk genre")
	}
	for _, leg := range campaign.Legs {
		if leg.Genre != engine.GenreCyberpunk {
			t.Error("legs should have cyberpunk genre")
		}
	}
}

func TestGenerateCampaignWithGenreShifts(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)
	genres := []engine.GenreID{engine.GenreFantasy, engine.GenreScifi, engine.GenreHorror}

	campaign := g.GenerateCampaignWithGenreShifts(3, genres)

	if !campaign.AllowGenreShifts {
		t.Error("genre shifts should be enabled")
	}

	for i, leg := range campaign.Legs {
		if leg.Genre != genres[i] {
			t.Errorf("leg %d should have genre %s, got %s", i, genres[i], leg.Genre)
		}
	}
}

func TestGeneratorLeg(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)

	leg := g.GenerateLeg("Start", "End", 200, DifficultyHard)

	if leg.OriginName != "Start" {
		t.Error("origin should match")
	}
	if leg.DestinationName != "End" {
		t.Error("destination should match")
	}
	if leg.Distance != 200 {
		t.Error("distance should match")
	}
	if leg.Difficulty != DifficultyHard {
		t.Error("difficulty should match")
	}
	if len(leg.Hazards) == 0 {
		t.Error("leg should have hazards")
	}
}

func TestGeneratorStopover(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)

	stopover := g.GenerateStopover(1, "Midpoint")

	if stopover.Name == "" {
		t.Error("stopover should have name")
	}
	if stopover.Description == "" {
		t.Error("stopover should have description")
	}
	if !stopover.HasService(ServiceTrading) {
		t.Error("stopover should always have trading")
	}
	if stopover.Inhabitants == "" {
		t.Error("stopover should have inhabitants description")
	}
}
