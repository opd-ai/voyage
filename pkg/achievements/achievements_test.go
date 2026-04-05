package achievements

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewAchievement(t *testing.T) {
	a := NewAchievement("test", "Test Achievement", "Description", CategorySurvival, TierBronze, 10, engine.GenreFantasy)

	if a.ID != "test" {
		t.Error("ID mismatch")
	}
	if a.Required != 10 {
		t.Error("required mismatch")
	}
	if a.Points != 10 {
		t.Error("bronze should be worth 10 points")
	}
	if a.Earned {
		t.Error("should not be earned initially")
	}
}

func TestAchievementTierPoints(t *testing.T) {
	testCases := []struct {
		tier     AchievementTier
		expected int
	}{
		{TierBronze, 10},
		{TierSilver, 25},
		{TierGold, 50},
		{TierLegendary, 100},
	}

	for _, tc := range testCases {
		a := NewAchievement("test", "Test", "Desc", CategorySurvival, tc.tier, 1, engine.GenreFantasy)
		if a.Points != tc.expected {
			t.Errorf("tier %s: expected %d points, got %d", TierName(tc.tier), tc.expected, a.Points)
		}
	}
}

func TestAchievementUpdateProgress(t *testing.T) {
	a := NewAchievement("test", "Test", "Desc", CategorySurvival, TierBronze, 10, engine.GenreFantasy)

	if a.UpdateProgress(5) {
		t.Error("should not be earned at 5/10")
	}
	if a.Earned {
		t.Error("should not be earned yet")
	}

	if !a.UpdateProgress(10) {
		t.Error("should be earned at 10/10")
	}
	if !a.Earned {
		t.Error("should be earned")
	}
}

func TestAchievementProgressPercent(t *testing.T) {
	a := NewAchievement("test", "Test", "Desc", CategorySurvival, TierBronze, 100, engine.GenreFantasy)

	a.Progress = 50
	if a.ProgressPercent() != 50 {
		t.Errorf("expected 50%%, got %d%%", a.ProgressPercent())
	}

	a.Progress = 150
	if a.ProgressPercent() != 100 {
		t.Error("should cap at 100%")
	}
}

func TestAchievementEarn(t *testing.T) {
	a := NewAchievement("test", "Test", "Desc", CategorySurvival, TierBronze, 10, engine.GenreFantasy)

	a.Earn(5)

	if !a.Earned {
		t.Error("should be earned")
	}
	if a.EarnedAt != 5 {
		t.Error("earned at day should be set")
	}

	// Can't earn twice
	a.Earn(10)
	if a.EarnedAt != 5 {
		t.Error("should not change earned day")
	}
}

func TestRunStatistics(t *testing.T) {
	stats := NewRunStatistics()

	if stats.LowestHealth != 1.0 {
		t.Error("lowest health should start at 1.0")
	}

	stats.DaysSurvived = 10
	stats.TradesCompleted = 5

	if stats.DaysSurvived != 10 {
		t.Error("days survived should be set")
	}
}

func TestAchievementTracker(t *testing.T) {
	tracker := NewAchievementTracker(engine.GenreFantasy)

	a1 := NewAchievement("test1", "Test 1", "Desc", CategorySurvival, TierBronze, 10, engine.GenreFantasy)
	a2 := NewAchievement("test2", "Test 2", "Desc", CategoryTrade, TierSilver, 20, engine.GenreFantasy)

	tracker.AddAchievement(a1)
	tracker.AddAchievement(a2)

	if len(tracker.Achievements) != 2 {
		t.Error("should have 2 achievements")
	}

	if tracker.GetAchievement("test1") != a1 {
		t.Error("should retrieve by ID")
	}
}

func TestAchievementTrackerGetByCategory(t *testing.T) {
	tracker := NewAchievementTracker(engine.GenreFantasy)

	tracker.AddAchievement(NewAchievement("s1", "Survival 1", "D", CategorySurvival, TierBronze, 1, engine.GenreFantasy))
	tracker.AddAchievement(NewAchievement("s2", "Survival 2", "D", CategorySurvival, TierSilver, 1, engine.GenreFantasy))
	tracker.AddAchievement(NewAchievement("t1", "Trade 1", "D", CategoryTrade, TierBronze, 1, engine.GenreFantasy))

	survival := tracker.GetByCategory(CategorySurvival)
	if len(survival) != 2 {
		t.Errorf("expected 2 survival achievements, got %d", len(survival))
	}

	trade := tracker.GetByCategory(CategoryTrade)
	if len(trade) != 1 {
		t.Errorf("expected 1 trade achievement, got %d", len(trade))
	}
}

func TestAchievementTrackerGetEarnedUnearned(t *testing.T) {
	tracker := NewAchievementTracker(engine.GenreFantasy)

	a1 := NewAchievement("a1", "A1", "D", CategorySurvival, TierBronze, 1, engine.GenreFantasy)
	a2 := NewAchievement("a2", "A2", "D", CategorySurvival, TierBronze, 1, engine.GenreFantasy)
	a3 := NewAchievement("a3", "A3", "D", CategorySurvival, TierBronze, 1, engine.GenreFantasy)
	a3.Hidden = true

	tracker.AddAchievement(a1)
	tracker.AddAchievement(a2)
	tracker.AddAchievement(a3)

	a1.Earn(1)

	earned := tracker.GetEarned()
	if len(earned) != 1 {
		t.Error("should have 1 earned")
	}

	unearned := tracker.GetUnearned()
	if len(unearned) != 1 {
		t.Error("should have 1 unearned (hidden not included)")
	}
}

func TestAchievementTrackerTotalPoints(t *testing.T) {
	tracker := NewAchievementTracker(engine.GenreFantasy)

	a1 := NewAchievement("a1", "A1", "D", CategorySurvival, TierBronze, 1, engine.GenreFantasy) // 10 pts
	a2 := NewAchievement("a2", "A2", "D", CategorySurvival, TierSilver, 1, engine.GenreFantasy) // 25 pts

	tracker.AddAchievement(a1)
	tracker.AddAchievement(a2)

	a1.Earn(1)

	if tracker.TotalPoints() != 10 {
		t.Errorf("expected 10 points, got %d", tracker.TotalPoints())
	}

	a2.Earn(2)

	if tracker.TotalPoints() != 35 {
		t.Errorf("expected 35 points, got %d", tracker.TotalPoints())
	}
}

func TestAchievementTrackerCheckAchievements(t *testing.T) {
	tracker := NewAchievementTracker(engine.GenreFantasy)

	a := NewAchievement("survive_10", "Survive 10", "D", CategorySurvival, TierBronze, 10, engine.GenreFantasy)
	tracker.AddAchievement(a)

	tracker.Stats.DaysSurvived = 5
	earned := tracker.CheckAchievements()
	if len(earned) != 0 {
		t.Error("should not earn at 5 days")
	}

	tracker.Stats.DaysSurvived = 10
	earned = tracker.CheckAchievements()
	if len(earned) != 1 {
		t.Error("should earn at 10 days")
	}

	// Should not earn again
	earned = tracker.CheckAchievements()
	if len(earned) != 0 {
		t.Error("should not earn twice")
	}
}

func TestAchievementTrackerCallback(t *testing.T) {
	tracker := NewAchievementTracker(engine.GenreFantasy)

	a := NewAchievement("survive_10", "Survive 10", "D", CategorySurvival, TierBronze, 10, engine.GenreFantasy)
	tracker.AddAchievement(a)

	callbackCalled := false
	tracker.OnEarned = func(earned *Achievement) {
		callbackCalled = true
	}

	tracker.Stats.DaysSurvived = 10
	tracker.CheckAchievements()

	if !callbackCalled {
		t.Error("callback should be called")
	}
}

func TestAchievementTrackerSetGenre(t *testing.T) {
	tracker := NewAchievementTracker(engine.GenreFantasy)
	a := NewAchievement("test", "Test", "D", CategorySurvival, TierBronze, 1, engine.GenreFantasy)
	tracker.AddAchievement(a)

	tracker.SetGenre(engine.GenreScifi)

	if tracker.Genre != engine.GenreScifi {
		t.Error("tracker genre should be updated")
	}
	if a.Genre != engine.GenreScifi {
		t.Error("achievement genre should be updated")
	}
}

func TestAchievementTrackerCompletionPercent(t *testing.T) {
	tracker := NewAchievementTracker(engine.GenreFantasy)

	tracker.AddAchievement(NewAchievement("a1", "A1", "D", CategorySurvival, TierBronze, 1, engine.GenreFantasy))
	tracker.AddAchievement(NewAchievement("a2", "A2", "D", CategorySurvival, TierBronze, 1, engine.GenreFantasy))

	if tracker.CompletionPercent() != 0 {
		t.Error("should be 0% at start")
	}

	tracker.Achievements[0].Earn(1)

	if tracker.CompletionPercent() != 50 {
		t.Error("should be 50% after earning 1 of 2")
	}
}

func TestAchievementSummary(t *testing.T) {
	tracker := NewAchievementTracker(engine.GenreFantasy)
	tracker.CurrentDay = 5

	a1 := NewAchievement("a1", "A1", "D", CategorySurvival, TierBronze, 1, engine.GenreFantasy)
	a2 := NewAchievement("a2", "A2", "D", CategoryTrade, TierSilver, 1, engine.GenreFantasy)
	tracker.AddAchievement(a1)
	tracker.AddAchievement(a2)

	a1.Earn(5)

	summary := tracker.GetSummary()

	if summary.TotalAchievements != 2 {
		t.Error("total should be 2")
	}
	if summary.EarnedCount != 1 {
		t.Error("earned should be 1")
	}
	if summary.EarnedPoints != 10 {
		t.Error("earned points should be 10")
	}
	if len(summary.NewlyEarned) != 1 {
		t.Error("should have 1 newly earned")
	}
	if summary.ByCategory[CategorySurvival] != 1 {
		t.Error("should have 1 survival earned")
	}
}

func TestCategoryName(t *testing.T) {
	for _, cat := range AllCategories() {
		name := CategoryName(cat)
		if name == "Unknown" || name == "" {
			t.Errorf("category %v should have name", cat)
		}
	}
}

func TestTierName(t *testing.T) {
	for _, tier := range AllTiers() {
		name := TierName(tier)
		if name == "Unknown" || name == "" {
			t.Errorf("tier %v should have name", tier)
		}
	}
}

func TestTierNameByGenre(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		for _, tier := range AllTiers() {
			name := TierNameByGenre(tier, genre)
			if name == "" {
				t.Errorf("tier %v genre %s should have name", tier, genre)
			}
		}
	}
}

func TestGenerator(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)

	tracker := g.GenerateAchievementTracker()

	if len(tracker.Achievements) < 20 {
		t.Errorf("expected at least 20 achievements, got %d", len(tracker.Achievements))
	}

	// Check categories have achievements
	for _, cat := range AllCategories() {
		achievements := tracker.GetByCategory(cat)
		if len(achievements) == 0 {
			t.Errorf("category %s should have achievements", CategoryName(cat))
		}
	}
}

func TestGeneratorDeterminism(t *testing.T) {
	g1 := NewGenerator(12345, engine.GenreFantasy)
	g2 := NewGenerator(12345, engine.GenreFantasy)

	t1 := g1.GenerateAchievementTracker()
	t2 := g2.GenerateAchievementTracker()

	if len(t1.Achievements) != len(t2.Achievements) {
		t.Error("same seed should produce same number of achievements")
	}

	for i := range t1.Achievements {
		if t1.Achievements[i].Name != t2.Achievements[i].Name {
			t.Error("same seed should produce same achievement names")
		}
	}
}

func TestGeneratorAllGenres(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		g := NewGenerator(12345, genre)
		tracker := g.GenerateAchievementTracker()

		if len(tracker.Achievements) < 20 {
			t.Errorf("genre %s: expected at least 20 achievements", genre)
		}

		for _, a := range tracker.Achievements {
			if a.Name == "" {
				t.Errorf("genre %s: achievement should have name", genre)
			}
			if a.Description == "" {
				t.Errorf("genre %s: achievement should have description", genre)
			}
		}
	}
}

func TestGeneratorSetGenre(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)
	g.SetGenre(engine.GenreCyberpunk)

	tracker := g.GenerateAchievementTracker()

	if tracker.Genre != engine.GenreCyberpunk {
		t.Error("tracker should have cyberpunk genre")
	}
}

func TestFullCrewAchievement(t *testing.T) {
	tracker := NewAchievementTracker(engine.GenreFantasy)

	a := NewAchievement("full_crew", "Full Crew", "D", CategorySurvival, TierGold, 1, engine.GenreFantasy)
	tracker.AddAchievement(a)

	tracker.Stats.CrewStarted = 4
	tracker.Stats.CrewSurvived = 4
	tracker.Stats.DaysSurvived = 10

	earned := tracker.CheckAchievements()
	if len(earned) != 1 {
		t.Error("should earn full crew achievement")
	}
}

func TestPerfectRunAchievement(t *testing.T) {
	tracker := NewAchievementTracker(engine.GenreFantasy)

	a := NewAchievement("perfect_run", "Perfect Run", "D", CategorySpecial, TierLegendary, 1, engine.GenreFantasy)
	tracker.AddAchievement(a)

	// Not enough days
	tracker.Stats.CrewStarted = 4
	tracker.Stats.CrewSurvived = 4
	tracker.Stats.DaysSurvived = 20

	earned := tracker.CheckAchievements()
	if len(earned) != 0 {
		t.Error("should not earn with only 20 days")
	}

	// Enough days but crew loss
	tracker.Stats.DaysSurvived = 30
	tracker.Stats.CrewSurvived = 3

	earned = tracker.CheckAchievements()
	if len(earned) != 0 {
		t.Error("should not earn with crew loss")
	}

	// Perfect conditions
	tracker.Stats.CrewSurvived = 4

	earned = tracker.CheckAchievements()
	if len(earned) != 1 {
		t.Error("should earn perfect run")
	}
}
