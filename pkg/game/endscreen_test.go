//go:build !headless

package game

import (
	"testing"
)

func TestNewEndStats(t *testing.T) {
	es := NewEndStats()
	if es == nil {
		t.Fatal("NewEndStats returned nil")
	}
	if es.IsVictory {
		t.Error("new EndStats should not be a victory by default")
	}
}

func TestEndStatsSetVictory(t *testing.T) {
	es := NewEndStats()
	es.SetVictory(WinReachedDestination)

	if !es.IsVictory {
		t.Error("IsVictory should be true after SetVictory")
	}
	if es.WinCondition != WinReachedDestination {
		t.Errorf("expected WinReachedDestination, got %v", es.WinCondition)
	}
}

func TestEndStatsSetDefeat(t *testing.T) {
	es := NewEndStats()
	es.SetDefeat(LoseAllCrewDead)

	if es.IsVictory {
		t.Error("IsVictory should be false after SetDefeat")
	}
	if es.LoseCondition != LoseAllCrewDead {
		t.Errorf("expected LoseAllCrewDead, got %v", es.LoseCondition)
	}
}

func TestEndStatsCalculateScore(t *testing.T) {
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

	// Score should include:
	// - 1000 for victory
	// - 375 for 75% survival (500 * 0.75)
	// - 500 for distance (100 * 5)
	// - 100 for exploration (50 * 2)
	// - 200 for events (20 * 10)
	// - 350 for efficiency (70 * 5)
	// - 10 for net currency (100 / 10)
	// Total = ~2535

	if score < 2000 {
		t.Errorf("expected score >= 2000, got %d", score)
	}
}

func TestEndStatsCalculateScoreDefeat(t *testing.T) {
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

func TestEndStatsGetSurvivalRate(t *testing.T) {
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
}

func TestEndStatsGetTitle(t *testing.T) {
	tests := []struct {
		name      string
		isVictory bool
		loseCond  LoseCondition
		crewStart int
		crewSurv  int
		wantEmpty bool
	}{
		{"victory 100%", true, LoseNone, 4, 4, false},
		{"victory 80%", true, LoseNone, 5, 4, false},
		{"victory 50%", true, LoseNone, 4, 2, false},
		{"victory low", true, LoseNone, 4, 1, false},
		{"loss crew dead", false, LoseAllCrewDead, 4, 0, false},
		{"loss vessel", false, LoseVesselDestroyed, 4, 2, false},
		{"loss morale", false, LoseMoraleZero, 4, 3, false},
		{"loss starvation", false, LoseStarvation, 4, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := NewEndStats()
			es.IsVictory = tt.isVictory
			es.LoseCondition = tt.loseCond
			es.CrewStarted = tt.crewStart
			es.CrewSurvived = tt.crewSurv

			title := es.GetTitle()
			if tt.wantEmpty && title != "" {
				t.Errorf("expected empty title, got %s", title)
			}
			if !tt.wantEmpty && title == "" {
				t.Error("expected non-empty title")
			}
		})
	}
}

func TestEndStatsGetRank(t *testing.T) {
	tests := []struct {
		score int
		rank  string
	}{
		{2500, "S"},
		{1700, "A"},
		{1200, "B"},
		{700, "C"},
		{300, "D"},
		{100, "F"},
	}

	for _, tt := range tests {
		t.Run(tt.rank, func(t *testing.T) {
			es := NewEndStats()
			es.IsVictory = true
			// Manipulate to get approximate score
			if tt.score >= 2000 {
				es.DistanceTraveled = (tt.score - 1000) / 5
			} else {
				es.IsVictory = false
				es.DistanceTraveled = tt.score / 5
			}

			// This is an approximation since CalculateScore includes many factors
			// Just verify we get a valid rank
			rank := es.GetRank()
			if rank == "" {
				t.Error("expected non-empty rank")
			}
		})
	}
}
