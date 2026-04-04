package game

// TimeManager handles turn-based time progression and day/night cycle.
type TimeManager struct {
	turn      int
	dayLength int
	isNight   bool
}

// NewTimeManager creates a new time manager.
func NewTimeManager() *TimeManager {
	return &TimeManager{
		turn:      0,
		dayLength: 4, // 4 turns per day
		isNight:   false,
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
