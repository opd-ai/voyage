//go:build headless

package game

import (
	"testing"
)

// TestNewEndStatsHeadless tests EndStats creation.
func TestNewEndStatsHeadless(t *testing.T) {
	es := NewEndStats()
	if es == nil {
		t.Fatal("NewEndStats returned nil")
	}
	if es.IsVictory {
		t.Error("new EndStats should not be a victory by default")
	}
}

// TestEndStatsSetVictoryHeadless tests victory setting.
func TestEndStatsSetVictoryHeadless(t *testing.T) {
	es := NewEndStats()
	es.SetVictory(WinReachedDestination)

	if !es.IsVictory {
		t.Error("IsVictory should be true after SetVictory")
	}
	if es.WinCondition != WinReachedDestination {
		t.Errorf("expected WinReachedDestination, got %v", es.WinCondition)
	}
}

// TestEndStatsSetDefeatHeadless tests defeat setting.
func TestEndStatsSetDefeatHeadless(t *testing.T) {
	es := NewEndStats()
	es.SetDefeat(LoseAllCrewDead)

	if es.IsVictory {
		t.Error("IsVictory should be false after SetDefeat")
	}
	if es.LoseCondition != LoseAllCrewDead {
		t.Errorf("expected LoseAllCrewDead, got %v", es.LoseCondition)
	}
}

// TestEndStatsCalculateScoreHeadless tests score calculation.
func TestEndStatsCalculateScoreHeadless(t *testing.T) {
	es := NewEndStats()
	es.IsVictory = true
	es.CrewStarted = 4
	es.CrewSurvived = 3
	es.DistanceTraveled = 100
	es.TilesExplored = 50
	es.EventsResolved = 20
	es.DaysTraveled = 30
	es.CurrencyEarned = 200
	es.CurrencySpent = 100

	score := es.CalculateScore()

	// Score should be reasonably high with these stats
	if score < 2000 {
		t.Errorf("expected score >= 2000, got %d", score)
	}
}

// TestEndStatsCalculateScoreDefeatHeadless tests defeat score.
func TestEndStatsCalculateScoreDefeatHeadless(t *testing.T) {
	es := NewEndStats()
	es.IsVictory = false
	es.CrewStarted = 4
	es.CrewSurvived = 0
	es.DistanceTraveled = 50

	score := es.CalculateScore()

	// Should not include victory bonus
	if score >= 1000 {
		t.Errorf("defeat score should be < 1000, got %d", score)
	}
}

// TestEndStatsCalculateScoreEdgeCases tests score edge cases.
func TestEndStatsCalculateScoreEdgeCases(t *testing.T) {
	// Test with zero crew started
	es := NewEndStats()
	es.IsVictory = true
	es.CrewStarted = 0
	es.CrewSurvived = 0
	score := es.CalculateScore()
	if score < 1000 {
		t.Errorf("victory with no crew should still have base score")
	}

	// Test with long journey (no efficiency bonus)
	es2 := NewEndStats()
	es2.IsVictory = true
	es2.DaysTraveled = 150
	score2 := es2.CalculateScore()
	// Should have victory bonus but no efficiency bonus
	if score2 < 1000 {
		t.Error("score should include victory bonus")
	}

	// Test with negative net currency
	es3 := NewEndStats()
	es3.IsVictory = true
	es3.CurrencySpent = 500
	es3.CurrencyEarned = 100
	score3 := es3.CalculateScore()
	// Should not subtract for negative currency
	if score3 < 1000 {
		t.Error("negative currency should not reduce below victory bonus")
	}
}

// TestEndStatsGetSurvivalRateHeadless tests survival rate calculation.
func TestEndStatsGetSurvivalRateHeadless(t *testing.T) {
	es := NewEndStats()

	// Test zero crew
	rate := es.GetSurvivalRate()
	if rate != 0 {
		t.Errorf("expected 0%% survival with zero crew, got %.1f%%", rate)
	}

	// Test 75% survival
	es.CrewStarted = 4
	es.CrewSurvived = 3
	rate = es.GetSurvivalRate()
	if rate != 75 {
		t.Errorf("expected 75%% survival, got %.1f%%", rate)
	}

	// Test 100% survival
	es.CrewSurvived = 4
	rate = es.GetSurvivalRate()
	if rate != 100 {
		t.Errorf("expected 100%% survival, got %.1f%%", rate)
	}

	// Test 50% survival
	es.CrewSurvived = 2
	rate = es.GetSurvivalRate()
	if rate != 50 {
		t.Errorf("expected 50%% survival, got %.1f%%", rate)
	}
}

// TestEndStatsGetTitleHeadless tests title generation.
func TestEndStatsGetTitleHeadless(t *testing.T) {
	tests := []struct {
		name      string
		isVictory bool
		loseCond  LoseCondition
		crewStart int
		crewSurv  int
		wantTitle string
	}{
		{"victory 100%", true, LoseNone, 4, 4, "Perfect Journey"},
		{"victory 80%", true, LoseNone, 5, 4, "Triumphant Arrival"},
		{"victory 60%", true, LoseNone, 5, 3, "Bittersweet Victory"},
		{"victory 40%", true, LoseNone, 5, 2, "Pyrrhic Victory"},
		{"loss crew dead", false, LoseAllCrewDead, 4, 0, "The Last Journey"},
		{"loss vessel", false, LoseVesselDestroyed, 4, 2, "Wreckage"},
		{"loss morale", false, LoseMoraleZero, 4, 3, "Abandoned"},
		{"loss starvation", false, LoseStarvation, 4, 1, "The Final Rest"},
		{"loss none", false, LoseNone, 4, 1, "Journey's End"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := NewEndStats()
			es.IsVictory = tt.isVictory
			es.LoseCondition = tt.loseCond
			es.CrewStarted = tt.crewStart
			es.CrewSurvived = tt.crewSurv

			title := es.GetTitle()
			if title != tt.wantTitle {
				t.Errorf("expected title %q, got %q", tt.wantTitle, title)
			}
		})
	}
}

// TestEndStatsGetRankHeadless tests rank calculation.
func TestEndStatsGetRankHeadless(t *testing.T) {
	// Test S rank (score >= 2000)
	es := NewEndStats()
	es.IsVictory = true
	es.DistanceTraveled = 200
	es.TilesExplored = 100
	es.EventsResolved = 50
	es.DaysTraveled = 30
	rank := es.GetRank()
	if rank != "S" {
		t.Errorf("expected S rank for high score, got %s", rank)
	}

	// Test lower ranks
	es2 := NewEndStats()
	es2.IsVictory = false
	es2.DistanceTraveled = 10
	rank2 := es2.GetRank()
	if rank2 == "" {
		t.Error("rank should not be empty")
	}
}

// TestEndStatsAllRanks tests all rank thresholds.
func TestEndStatsAllRanks(t *testing.T) {
	ranks := []string{"S", "A", "B", "C", "D", "F"}
	seenRanks := make(map[string]bool)

	// Generate different score scenarios
	for _, victory := range []bool{true, false} {
		for dist := 0; dist <= 400; dist += 50 {
			es := NewEndStats()
			es.IsVictory = victory
			es.DistanceTraveled = dist
			rank := es.GetRank()
			seenRanks[rank] = true
		}
	}

	// We may not hit all ranks but verify the ones we get are valid
	for rank := range seenRanks {
		found := false
		for _, valid := range ranks {
			if rank == valid {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("invalid rank returned: %s", rank)
		}
	}
}

// TestEndStatsDefeatTypes tests all defeat condition titles.
func TestEndStatsDefeatTypesHeadless(t *testing.T) {
	defeatConds := []LoseCondition{
		LoseAllCrewDead,
		LoseVesselDestroyed,
		LoseMoraleZero,
		LoseStarvation,
	}

	for _, lc := range defeatConds {
		es := NewEndStats()
		es.SetDefeat(lc)
		title := es.GetTitle()
		if title == "" {
			t.Errorf("GetTitle returned empty for %v", lc)
		}
		if title == "Journey's End" {
			t.Errorf("defeat condition %v should have specific title", lc)
		}
	}
}
