package game

// Season represents a time of year affecting gameplay.
type Season int

const (
	// SeasonSpring has mild conditions with occasional rain.
	SeasonSpring Season = iota
	// SeasonSummer has hot conditions with increased water consumption.
	SeasonSummer
	// SeasonAutumn has moderate conditions with increased foraging yields.
	SeasonAutumn
	// SeasonWinter has harsh conditions with increased hazards and resource costs.
	SeasonWinter
)

// AllSeasons returns all seasons in order.
func AllSeasons() []Season {
	return []Season{SeasonSpring, SeasonSummer, SeasonAutumn, SeasonWinter}
}

// SeasonName returns a human-readable name for the season.
func SeasonName(s Season) string {
	names := map[Season]string{
		SeasonSpring: "Spring",
		SeasonSummer: "Summer",
		SeasonAutumn: "Autumn",
		SeasonWinter: "Winter",
	}
	return names[s]
}

// SeasonModifiers returns modifiers for resource consumption and hazard frequency.
// Returns (resourceCostMod, hazardFreqMod) where 1.0 is baseline.
func SeasonModifiers(s Season) (resourceCost, hazardFreq float64) {
	switch s {
	case SeasonSpring:
		return 1.0, 0.9 // Normal costs, slightly fewer hazards
	case SeasonSummer:
		return 1.2, 0.8 // Higher water/food costs, fewer hazards
	case SeasonAutumn:
		return 0.9, 1.0 // Lower resource costs (good foraging), normal hazards
	case SeasonWinter:
		return 1.4, 1.3 // Much higher costs, more hazards
	default:
		return 1.0, 1.0
	}
}

// TimeManager handles turn-based time progression and day/night cycle.
type TimeManager struct {
	turn         int
	dayLength    int
	seasonLength int // days per season
	isNight      bool
}

// NewTimeManager creates a new time manager.
func NewTimeManager() *TimeManager {
	return &TimeManager{
		turn:         0,
		dayLength:    4,  // 4 turns per day
		seasonLength: 20, // 20 days per season (80 turns)
		isNight:      false,
	}
}

// Turn returns the current turn number.
func (tm *TimeManager) Turn() int {
	return tm.turn
}

// Day returns the current day number (1-indexed).
func (tm *TimeManager) Day() int {
	return (tm.turn / tm.dayLength) + 1
}

// TimeOfDay returns the position within the current day (0 to dayLength-1).
func (tm *TimeManager) TimeOfDay() int {
	return tm.turn % tm.dayLength
}

// IsNight returns true if it's currently night time.
func (tm *TimeManager) IsNight() bool {
	// Night is the last turn of each day
	return tm.TimeOfDay() >= tm.dayLength-1
}

// Advance increments the turn counter and returns the new turn.
func (tm *TimeManager) Advance() int {
	tm.turn++
	return tm.turn
}

// AdvanceMultiple advances time by multiple turns.
func (tm *TimeManager) AdvanceMultiple(turns int) int {
	tm.turn += turns
	return tm.turn
}

// SetTurn sets the turn counter directly.
func (tm *TimeManager) SetTurn(turn int) {
	tm.turn = turn
}

// TurnsUntilNight returns the number of turns until nightfall.
func (tm *TimeManager) TurnsUntilNight() int {
	remaining := tm.dayLength - 1 - tm.TimeOfDay()
	if remaining <= 0 {
		return tm.dayLength - 1
	}
	return remaining
}

// TurnsUntilDawn returns the number of turns until dawn.
func (tm *TimeManager) TurnsUntilDawn() int {
	if !tm.IsNight() {
		return tm.TurnsUntilNight() + 1
	}
	return 1
}

// PhaseOfDay returns a string describing the current time of day.
func (tm *TimeManager) PhaseOfDay() string {
	tod := tm.TimeOfDay()
	switch {
	case tod == 0:
		return "Dawn"
	case tod == tm.dayLength-1:
		return "Night"
	case tod < tm.dayLength/2:
		return "Morning"
	default:
		return "Afternoon"
	}
}

// Season returns the current season.
func (tm *TimeManager) Season() Season {
	day := tm.Day()
	seasonDays := tm.seasonLength * 4 // Total days in a year
	dayInYear := (day - 1) % seasonDays
	seasonIndex := dayInYear / tm.seasonLength
	return Season(seasonIndex % 4)
}

// DayInSeason returns the current day within the season (1-indexed).
func (tm *TimeManager) DayInSeason() int {
	day := tm.Day()
	return ((day - 1) % tm.seasonLength) + 1
}

// SeasonProgress returns the progress through the current season (0.0 to 1.0).
func (tm *TimeManager) SeasonProgress() float64 {
	return float64(tm.DayInSeason()-1) / float64(tm.seasonLength)
}

// DaysUntilSeasonChange returns the number of days until the season changes.
func (tm *TimeManager) DaysUntilSeasonChange() int {
	return tm.seasonLength - tm.DayInSeason() + 1
}

// TurnsUntilSeasonChange returns the number of turns until the season changes.
func (tm *TimeManager) TurnsUntilSeasonChange() int {
	return tm.DaysUntilSeasonChange() * tm.dayLength
}

// Year returns the current year (1-indexed).
func (tm *TimeManager) Year() int {
	totalSeasons := 4
	daysPerYear := tm.seasonLength * totalSeasons
	return ((tm.Day() - 1) / daysPerYear) + 1
}

// ResourceCostModifier returns the current resource cost modifier based on season.
func (tm *TimeManager) ResourceCostModifier() float64 {
	cost, _ := SeasonModifiers(tm.Season())
	return cost
}

// HazardFrequencyModifier returns the current hazard frequency modifier based on season.
func (tm *TimeManager) HazardFrequencyModifier() float64 {
	_, hazard := SeasonModifiers(tm.Season())
	return hazard
}

// SetSeasonLength sets the number of days per season.
func (tm *TimeManager) SetSeasonLength(days int) {
	if days > 0 {
		tm.seasonLength = days
	}
}

// SeasonLength returns the number of days per season.
func (tm *TimeManager) SeasonLength() int {
	return tm.seasonLength
}
