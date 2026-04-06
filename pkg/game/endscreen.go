package game

// EndStats holds statistics for the game end screen.
type EndStats struct {
	// Journey metrics
	DaysTraveled     int
	DistanceTraveled int
	TilesExplored    int

	// Party metrics
	CrewStarted    int
	CrewLost       int
	CrewSurvived   int
	TotalDamageInc int

	// Events
	EventsEncountered int
	EventsResolved    int
	ChoicesMade       int

	// Resources
	FoodConsumed   float64
	WaterConsumed  float64
	FuelConsumed   float64
	MedicineUsed   float64
	CurrencySpent  float64
	CurrencyEarned float64

	// Victory info
	IsVictory     bool
	WinCondition  WinCondition
	LoseCondition LoseCondition
}

// NewEndStats creates a new end stats with default values.
func NewEndStats() *EndStats {
	return &EndStats{}
}

// SetVictory marks the game as won.
func (es *EndStats) SetVictory(wc WinCondition) {
	es.IsVictory = true
	es.WinCondition = wc
}

// SetDefeat marks the game as lost.
func (es *EndStats) SetDefeat(lc LoseCondition) {
	es.IsVictory = false
	es.LoseCondition = lc
}

// CalculateScore computes a final score based on performance.
// Uses int64 for intermediate calculations to prevent overflow (M-016).
func (es *EndStats) CalculateScore() int {
	var score int64 = 0

	score += es.calculateVictoryBonus()
	score += es.calculateSurvivalBonus()
	score += es.calculateExplorationScore()
	score += es.calculateEfficiencyBonus()
	score += es.calculateEconomicScore()

	return clampScore(score)
}

// calculateVictoryBonus returns points for winning the game.
func (es *EndStats) calculateVictoryBonus() int64 {
	if es.IsVictory {
		return 1000
	}
	return 0
}

// calculateSurvivalBonus returns points based on crew survival rate.
func (es *EndStats) calculateSurvivalBonus() int64 {
	if es.CrewStarted == 0 {
		return 0
	}
	survivalRate := float64(es.CrewSurvived) / float64(es.CrewStarted)
	return int64(survivalRate * 500)
}

// calculateExplorationScore returns points for distance and exploration.
func (es *EndStats) calculateExplorationScore() int64 {
	var score int64
	score += int64(es.DistanceTraveled) * 5
	score += int64(es.TilesExplored) * 2
	score += int64(es.EventsResolved) * 10
	return score
}

// calculateEfficiencyBonus returns bonus points for quick completion.
func (es *EndStats) calculateEfficiencyBonus() int64 {
	if es.DaysTraveled > 0 && es.DaysTraveled < 100 {
		return int64(100-es.DaysTraveled) * 5
	}
	return 0
}

// calculateEconomicScore returns points for positive net currency.
func (es *EndStats) calculateEconomicScore() int64 {
	netCurrency := es.CurrencyEarned - es.CurrencySpent
	if netCurrency > 0 {
		return int64(netCurrency) / 10
	}
	return 0
}

// clampScore caps score at maximum int value to prevent overflow (M-016).
func clampScore(score int64) int {
	const maxScore = int64(1<<31 - 1)
	if score > maxScore {
		score = maxScore
	}
	if score < 0 {
		score = 0
	}
	return int(score)
}

// GetSurvivalRate returns the crew survival rate as a percentage.
func (es *EndStats) GetSurvivalRate() float64 {
	if es.CrewStarted == 0 {
		return 0
	}
	return float64(es.CrewSurvived) / float64(es.CrewStarted) * 100
}

// GetTitle returns a title based on performance.
func (es *EndStats) GetTitle() string {
	if !es.IsVictory {
		switch es.LoseCondition {
		case LoseAllCrewDead:
			return "The Last Journey"
		case LoseVesselDestroyed:
			return "Wreckage"
		case LoseMoraleZero:
			return "Abandoned"
		case LoseStarvation:
			return "The Final Rest"
		default:
			return "Journey's End"
		}
	}

	survivalRate := es.GetSurvivalRate()
	switch {
	case survivalRate >= 100:
		return "Perfect Journey"
	case survivalRate >= 80:
		return "Triumphant Arrival"
	case survivalRate >= 50:
		return "Bittersweet Victory"
	default:
		return "Pyrrhic Victory"
	}
}

// GetRank returns a letter rank based on score.
func (es *EndStats) GetRank() string {
	score := es.CalculateScore()
	switch {
	case score >= 2000:
		return "S"
	case score >= 1500:
		return "A"
	case score >= 1000:
		return "B"
	case score >= 500:
		return "C"
	case score >= 200:
		return "D"
	default:
		return "F"
	}
}
