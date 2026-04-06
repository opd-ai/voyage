package achievements

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// AchievementCategory groups achievements by type
type AchievementCategory int

const (
	CategorySurvival AchievementCategory = iota
	CategoryTrade
	CategoryExploration
	CategoryCombat
	CategorySocial
	CategorySpecial
)

// AchievementTier represents the difficulty/prestige level
type AchievementTier int

const (
	TierBronze AchievementTier = iota
	TierSilver
	TierGold
	TierLegendary
)

// Achievement represents a single milestone that can be earned
type Achievement struct {
	ID          string
	Name        string
	Description string
	Genre       engine.GenreID
	Category    AchievementCategory
	Tier        AchievementTier

	// Requirements
	Required int  // Threshold value to earn
	Hidden   bool // Don't show until earned

	// State
	Progress int
	Earned   bool
	EarnedAt int // Day earned (0 if not earned)

	// Meta rewards
	UnlockReward string // What this achievement unlocks
	Points       int    // Achievement points value
}

// NewAchievement creates a new achievement
func NewAchievement(id, name, description string, category AchievementCategory, tier AchievementTier, required int, genre engine.GenreID) *Achievement {
	points := 10
	switch tier {
	case TierSilver:
		points = 25
	case TierGold:
		points = 50
	case TierLegendary:
		points = 100
	}

	return &Achievement{
		ID:          id,
		Name:        name,
		Description: description,
		Genre:       genre,
		Category:    category,
		Tier:        tier,
		Required:    required,
		Points:      points,
	}
}

// SetGenre updates the achievement's genre
func (a *Achievement) SetGenre(g engine.GenreID) {
	a.Genre = g
}

// UpdateProgress sets progress and checks if earned
func (a *Achievement) UpdateProgress(value int) bool {
	if a.Earned {
		return false
	}

	a.Progress = value
	if a.Progress >= a.Required {
		a.Earned = true
		return true
	}
	return false
}

// Earn marks the achievement as earned on a specific day.
// Validates day is non-negative to prevent corrupt timestamps (L-010).
func (a *Achievement) Earn(day int) {
	if !a.Earned {
		a.Earned = true
		if day < 0 {
			day = 0
		}
		a.EarnedAt = day
		a.Progress = a.Required
	}
}

// ProgressPercent returns completion percentage (0-100).
// Uses int64 for intermediate calculation to prevent overflow (L-005).
func (a *Achievement) ProgressPercent() int {
	if a.Required == 0 {
		return 100
	}
	percent := int(int64(a.Progress) * 100 / int64(a.Required))
	if percent > 100 {
		return 100
	}
	return percent
}

// RunStatistics holds all stats collected during a run
type RunStatistics struct {
	// Survival
	DaysSurvived    int
	CrewStarted     int
	CrewSurvived    int
	DaysWithoutLoss int
	LowestHealth    float64
	HighestMorale   float64

	// Trade
	TotalGoldEarned int
	TotalGoldSpent  int
	TradesCompleted int
	BestSingleTrade int
	RegionsTraded   int

	// Exploration
	DistanceTraveled int
	DiscoveriesMade  int
	RegionsVisited   int
	SecretsFound     int
	LoreCollected    int

	// Combat
	EnemiesDefeated   int
	BattlesWon        int
	FlawlessVictories int
	DamageTaken       int
	DamageDealt       int

	// Social
	FactionAllies       int
	NPCsHelped          int
	QuestsCompleted     int
	ReputationGained    int
	CompanionsRecruited int

	// Special
	PerfectDays        int // Days with no negative events
	CloseCallsSurvived int
	CriticalSuccesses  int
	GenreSpecificStat1 int
	GenreSpecificStat2 int
}

// NewRunStatistics creates a fresh statistics tracker
func NewRunStatistics() *RunStatistics {
	return &RunStatistics{
		LowestHealth: 1.0,
	}
}

// AchievementTracker manages all achievements for a run
type AchievementTracker struct {
	Achievements []*Achievement
	Stats        *RunStatistics
	Genre        engine.GenreID
	CurrentDay   int

	// Callbacks
	OnEarned func(a *Achievement)
}

// NewAchievementTracker creates a new tracker with the given genre
func NewAchievementTracker(genre engine.GenreID) *AchievementTracker {
	return &AchievementTracker{
		Achievements: make([]*Achievement, 0),
		Stats:        NewRunStatistics(),
		Genre:        genre,
		CurrentDay:   0,
	}
}

// SetGenre updates all achievements to new genre
func (t *AchievementTracker) SetGenre(g engine.GenreID) {
	t.Genre = g
	for _, a := range t.Achievements {
		a.SetGenre(g)
	}
}

// AddAchievement registers an achievement
func (t *AchievementTracker) AddAchievement(a *Achievement) {
	t.Achievements = append(t.Achievements, a)
}

// GetAchievement retrieves an achievement by ID
func (t *AchievementTracker) GetAchievement(id string) *Achievement {
	for _, a := range t.Achievements {
		if a.ID == id {
			return a
		}
	}
	return nil
}

// GetByCategory returns all achievements in a category
func (t *AchievementTracker) GetByCategory(category AchievementCategory) []*Achievement {
	result := make([]*Achievement, 0)
	for _, a := range t.Achievements {
		if a.Category == category {
			result = append(result, a)
		}
	}
	return result
}

// GetEarned returns all earned achievements
func (t *AchievementTracker) GetEarned() []*Achievement {
	earned := make([]*Achievement, 0)
	for _, a := range t.Achievements {
		if a.Earned {
			earned = append(earned, a)
		}
	}
	return earned
}

// GetUnearned returns all unearned achievements (excluding hidden)
func (t *AchievementTracker) GetUnearned() []*Achievement {
	unearned := make([]*Achievement, 0)
	for _, a := range t.Achievements {
		if !a.Earned && !a.Hidden {
			unearned = append(unearned, a)
		}
	}
	return unearned
}

// TotalPoints returns total achievement points earned
func (t *AchievementTracker) TotalPoints() int {
	total := 0
	for _, a := range t.Achievements {
		if a.Earned {
			total += a.Points
		}
	}
	return total
}

// CheckAchievements evaluates all achievements against current stats
func (t *AchievementTracker) CheckAchievements() []*Achievement {
	newlyEarned := make([]*Achievement, 0)

	for _, a := range t.Achievements {
		if a.Earned {
			continue
		}
		if t.checkAchievement(a) {
			a.EarnedAt = t.CurrentDay
			newlyEarned = append(newlyEarned, a)
			if t.OnEarned != nil {
				t.OnEarned(a)
			}
		}
	}

	return newlyEarned
}

// checkAchievement evaluates a single achievement against current stats.
func (t *AchievementTracker) checkAchievement(a *Achievement) bool {
	switch a.ID {
	// Survival achievements
	case "survive_10", "survive_30", "survive_100":
		return a.UpdateProgress(t.Stats.DaysSurvived)
	case "no_losses":
		return a.UpdateProgress(t.Stats.DaysWithoutLoss)
	case "full_crew":
		return t.checkFullCrew(a)

	// Trade achievements
	case "trader", "merchant":
		return a.UpdateProgress(t.Stats.TradesCompleted)
	case "tycoon":
		return a.UpdateProgress(t.Stats.TotalGoldEarned)
	case "regional_trader":
		return a.UpdateProgress(t.Stats.RegionsTraded)

	// Exploration achievements
	case "explorer":
		return a.UpdateProgress(t.Stats.DistanceTraveled)
	case "cartographer":
		return a.UpdateProgress(t.Stats.RegionsVisited)
	case "discoverer":
		return a.UpdateProgress(t.Stats.DiscoveriesMade)
	case "lore_keeper":
		return a.UpdateProgress(t.Stats.LoreCollected)
	case "secret_finder":
		return a.UpdateProgress(t.Stats.SecretsFound)

	// Combat achievements
	case "warrior":
		return a.UpdateProgress(t.Stats.EnemiesDefeated)
	case "champion":
		return a.UpdateProgress(t.Stats.BattlesWon)
	case "flawless":
		return a.UpdateProgress(t.Stats.FlawlessVictories)

	// Social achievements
	case "diplomat":
		return a.UpdateProgress(t.Stats.FactionAllies)
	case "helper":
		return a.UpdateProgress(t.Stats.NPCsHelped)
	case "questor":
		return a.UpdateProgress(t.Stats.QuestsCompleted)
	case "recruiter":
		return a.UpdateProgress(t.Stats.CompanionsRecruited)

	// Special achievements
	case "perfect_run":
		return t.checkPerfectRun(a)
	case "close_calls":
		return a.UpdateProgress(t.Stats.CloseCallsSurvived)
	case "critical_master":
		return a.UpdateProgress(t.Stats.CriticalSuccesses)
	}
	return false
}

// checkFullCrew checks if full crew survived at least one day.
func (t *AchievementTracker) checkFullCrew(a *Achievement) bool {
	if t.Stats.CrewSurvived >= t.Stats.CrewStarted && t.Stats.DaysSurvived > 0 {
		return a.UpdateProgress(1)
	}
	return false
}

// checkPerfectRun checks if a perfect run was achieved.
func (t *AchievementTracker) checkPerfectRun(a *Achievement) bool {
	if t.Stats.CrewSurvived >= t.Stats.CrewStarted && t.Stats.DaysSurvived >= 30 {
		return a.UpdateProgress(1)
	}
	return false
}

// Tick advances the day counter
func (t *AchievementTracker) Tick() {
	t.CurrentDay++
}

// CompletionPercent returns overall achievement completion
func (t *AchievementTracker) CompletionPercent() int {
	if len(t.Achievements) == 0 {
		return 0
	}
	earned := len(t.GetEarned())
	return (earned * 100) / len(t.Achievements)
}

// CategoryName returns the human-readable category name
func CategoryName(c AchievementCategory) string {
	names := map[AchievementCategory]string{
		CategorySurvival:    "Survival",
		CategoryTrade:       "Trade",
		CategoryExploration: "Exploration",
		CategoryCombat:      "Combat",
		CategorySocial:      "Social",
		CategorySpecial:     "Special",
	}
	if name, ok := names[c]; ok {
		return name
	}
	return "Unknown"
}

// TierName returns the human-readable tier name
func TierName(tier AchievementTier) string {
	names := map[AchievementTier]string{
		TierBronze:    "Bronze",
		TierSilver:    "Silver",
		TierGold:      "Gold",
		TierLegendary: "Legendary",
	}
	if name, ok := names[tier]; ok {
		return name
	}
	return "Unknown"
}

// TierNameByGenre returns genre-appropriate tier name
func TierNameByGenre(tier AchievementTier, genre engine.GenreID) string {
	tierNames := map[engine.GenreID]map[AchievementTier]string{
		engine.GenreFantasy: {
			TierBronze:    "Apprentice",
			TierSilver:    "Journeyman",
			TierGold:      "Master",
			TierLegendary: "Legendary",
		},
		engine.GenreScifi: {
			TierBronze:    "Cadet",
			TierSilver:    "Officer",
			TierGold:      "Commander",
			TierLegendary: "Admiral",
		},
		engine.GenreHorror: {
			TierBronze:    "Survivor",
			TierSilver:    "Veteran",
			TierGold:      "Expert",
			TierLegendary: "Legend",
		},
		engine.GenreCyberpunk: {
			TierBronze:    "Street",
			TierSilver:    "Professional",
			TierGold:      "Elite",
			TierLegendary: "Legend",
		},
		engine.GenrePostapoc: {
			TierBronze:    "Wanderer",
			TierSilver:    "Scavenger",
			TierGold:      "Wasteland Veteran",
			TierLegendary: "Legend of the Wastes",
		},
	}

	if genreTiers, ok := tierNames[genre]; ok {
		if name, ok := genreTiers[tier]; ok {
			return name
		}
	}
	return TierName(tier)
}

// AllCategories returns all achievement categories
func AllCategories() []AchievementCategory {
	return []AchievementCategory{
		CategorySurvival, CategoryTrade, CategoryExploration,
		CategoryCombat, CategorySocial, CategorySpecial,
	}
}

// AllTiers returns all achievement tiers
func AllTiers() []AchievementTier {
	return []AchievementTier{TierBronze, TierSilver, TierGold, TierLegendary}
}

// AchievementSummary provides end-screen summary data
type AchievementSummary struct {
	TotalAchievements int
	EarnedCount       int
	TotalPoints       int
	EarnedPoints      int
	NewlyEarned       []*Achievement
	ByCategory        map[AchievementCategory]int
	ByTier            map[AchievementTier]int
}

// GetSummary generates an end-screen summary
func (t *AchievementTracker) GetSummary() *AchievementSummary {
	summary := &AchievementSummary{
		TotalAchievements: len(t.Achievements),
		ByCategory:        make(map[AchievementCategory]int),
		ByTier:            make(map[AchievementTier]int),
		NewlyEarned:       make([]*Achievement, 0),
	}

	for _, a := range t.Achievements {
		summary.TotalPoints += a.Points
		if a.Earned {
			summary.EarnedCount++
			summary.EarnedPoints += a.Points
			summary.ByCategory[a.Category]++
			summary.ByTier[a.Tier]++
			if a.EarnedAt == t.CurrentDay {
				summary.NewlyEarned = append(summary.NewlyEarned, a)
			}
		}
	}

	return summary
}
