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
func (es *EndStats) CalculateScore() int {
	score := 0

	// Base score for surviving
	if es.IsVictory {
		score += 1000
	}

	// Bonus for crew survival
	if es.CrewStarted > 0 {
		survivalRate := float64(es.CrewSurvived) / float64(es.CrewStarted)
		score += int(survivalRate * 500)
	}

	// Points for distance traveled
	score += es.DistanceTraveled * 5

	// Points for exploration
	score += es.TilesExplored * 2

	// Points for events resolved
	score += es.EventsResolved * 10

	// Efficiency bonus for quick completion
	if es.DaysTraveled > 0 && es.DaysTraveled < 100 {
		efficiencyBonus := (100 - es.DaysTraveled) * 5
		score += efficiencyBonus
	}

	// Economic score
	netCurrency := es.CurrencyEarned - es.CurrencySpent
	if netCurrency > 0 {
		score += int(netCurrency / 10)
	}

	return score
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
