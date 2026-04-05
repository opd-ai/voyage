package convoy

import (
	"sort"
)

// ComparisonResult holds the comparison between convoy members' runs.
type ComparisonResult struct {
	ConvoyID      ConvoyID          `json:"convoyId"`
	Seed          int64             `json:"seed"`
	Rankings      []*PlayerRanking  `json:"rankings"`
	Comparisons   []*StatComparison `json:"comparisons"`
	WinnerID      PlayerID          `json:"winnerId,omitempty"`
	FastestID     PlayerID          `json:"fastestId,omitempty"`
	FirstToFinish PlayerID          `json:"firstToFinish,omitempty"`
}

// PlayerRanking represents a player's overall ranking.
type PlayerRanking struct {
	Rank      int      `json:"rank"`
	PlayerID  PlayerID `json:"playerId"`
	Name      string   `json:"name"`
	Score     int      `json:"score"`
	Victory   bool     `json:"victory"`
	Days      int      `json:"days"`
	Survivors int      `json:"survivors"`
}

// StatComparison compares a specific stat across players.
type StatComparison struct {
	Stat     string           `json:"stat"`
	Values   map[PlayerID]int `json:"values"`
	BestID   PlayerID         `json:"bestId"`
	BestName string           `json:"bestName"`
}

// Compare generates a comparison of all runs in the convoy.
func (c *Convoy) Compare() *ComparisonResult {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := &ComparisonResult{
		ConvoyID:    c.ID,
		Seed:        c.Seed,
		Rankings:    make([]*PlayerRanking, 0),
		Comparisons: make([]*StatComparison, 0),
	}

	if len(c.Runs) == 0 {
		return result
	}

	// Build player name map
	nameMap := make(map[PlayerID]string)
	for _, p := range c.Players {
		nameMap[p.ID] = p.Name
	}

	// Create rankings
	for _, run := range c.Runs {
		result.Rankings = append(result.Rankings, &PlayerRanking{
			PlayerID:  run.PlayerID,
			Name:      nameMap[run.PlayerID],
			Score:     run.Score,
			Victory:   run.IsVictory,
			Days:      run.DaysTraveled,
			Survivors: run.Survivors,
		})
	}

	// Sort rankings by score (descending), then by days (ascending)
	sort.Slice(result.Rankings, func(i, j int) bool {
		if result.Rankings[i].Score != result.Rankings[j].Score {
			return result.Rankings[i].Score > result.Rankings[j].Score
		}
		return result.Rankings[i].Days < result.Rankings[j].Days
	})

	// Assign ranks
	for i, r := range result.Rankings {
		r.Rank = i + 1
	}

	// Determine winner (highest score with victory)
	for _, r := range result.Rankings {
		if r.Victory {
			result.WinnerID = r.PlayerID
			break
		}
	}

	// Find fastest (fewest days with victory)
	fastestDays := int(^uint(0) >> 1) // Max int
	for _, r := range result.Rankings {
		if r.Victory && r.Days < fastestDays {
			fastestDays = r.Days
			result.FastestID = r.PlayerID
		}
	}

	// Find first to finish
	earliestFinish := c.Runs[0].FinishedAt
	result.FirstToFinish = c.Runs[0].PlayerID
	for _, run := range c.Runs {
		if !run.FinishedAt.IsZero() && run.FinishedAt.Before(earliestFinish) {
			earliestFinish = run.FinishedAt
			result.FirstToFinish = run.PlayerID
		}
	}

	// Generate stat comparisons
	result.Comparisons = c.generateStatComparisons(nameMap)

	return result
}

// generateStatComparisons creates detailed comparisons for each stat.
func (c *Convoy) generateStatComparisons(nameMap map[PlayerID]string) []*StatComparison {
	comparisons := make([]*StatComparison, 0)

	// Score comparison
	scoreComp := &StatComparison{
		Stat:   "Score",
		Values: make(map[PlayerID]int),
	}
	var bestScore int
	for _, run := range c.Runs {
		scoreComp.Values[run.PlayerID] = run.Score
		if run.Score > bestScore {
			bestScore = run.Score
			scoreComp.BestID = run.PlayerID
			scoreComp.BestName = nameMap[run.PlayerID]
		}
	}
	comparisons = append(comparisons, scoreComp)

	// Days comparison (lower is better)
	daysComp := &StatComparison{
		Stat:   "Days Traveled",
		Values: make(map[PlayerID]int),
	}
	bestDays := int(^uint(0) >> 1) // Max int
	for _, run := range c.Runs {
		daysComp.Values[run.PlayerID] = run.DaysTraveled
		if run.DaysTraveled < bestDays && run.DaysTraveled > 0 {
			bestDays = run.DaysTraveled
			daysComp.BestID = run.PlayerID
			daysComp.BestName = nameMap[run.PlayerID]
		}
	}
	comparisons = append(comparisons, daysComp)

	// Survivors comparison
	survivorComp := &StatComparison{
		Stat:   "Survivors",
		Values: make(map[PlayerID]int),
	}
	var bestSurvivors int
	for _, run := range c.Runs {
		survivorComp.Values[run.PlayerID] = run.Survivors
		if run.Survivors > bestSurvivors {
			bestSurvivors = run.Survivors
			survivorComp.BestID = run.PlayerID
			survivorComp.BestName = nameMap[run.PlayerID]
		}
	}
	comparisons = append(comparisons, survivorComp)

	// Events seen comparison
	eventsComp := &StatComparison{
		Stat:   "Events Faced",
		Values: make(map[PlayerID]int),
	}
	var mostEvents int
	for _, run := range c.Runs {
		eventsComp.Values[run.PlayerID] = run.EventsSeen
		if run.EventsSeen > mostEvents {
			mostEvents = run.EventsSeen
			eventsComp.BestID = run.PlayerID
			eventsComp.BestName = nameMap[run.PlayerID]
		}
	}
	comparisons = append(comparisons, eventsComp)

	return comparisons
}

// GetWinner returns the winning player and their run data.
func (c *Convoy) GetWinner() (*Player, *RunData) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Find the player with highest score who achieved victory
	var winner *Player
	var winningRun *RunData
	var highestScore int

	for _, run := range c.Runs {
		if run.IsVictory && run.Score > highestScore {
			highestScore = run.Score
			winningRun = run
		}
	}

	if winningRun != nil {
		for _, p := range c.Players {
			if p.ID == winningRun.PlayerID {
				winner = p
				break
			}
		}
	}

	return winner, winningRun
}

// GetFastest returns the player who finished fastest (fewest days).
func (c *Convoy) GetFastest() (*Player, *RunData) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var fastest *Player
	var fastestRun *RunData
	fewestDays := int(^uint(0) >> 1)

	for _, run := range c.Runs {
		if run.IsVictory && run.DaysTraveled < fewestDays {
			fewestDays = run.DaysTraveled
			fastestRun = run
		}
	}

	if fastestRun != nil {
		for _, p := range c.Players {
			if p.ID == fastestRun.PlayerID {
				fastest = p
				break
			}
		}
	}

	return fastest, fastestRun
}

// GetFirstToFinish returns the player who finished their run first.
func (c *Convoy) GetFirstToFinish() (*Player, *RunData) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var first *Player
	var firstRun *RunData

	for _, run := range c.Runs {
		if run.FinishedAt.IsZero() {
			continue
		}
		if firstRun == nil || run.FinishedAt.Before(firstRun.FinishedAt) {
			firstRun = run
		}
	}

	if firstRun != nil {
		for _, p := range c.Players {
			if p.ID == firstRun.PlayerID {
				first = p
				break
			}
		}
	}

	return first, firstRun
}

// GetVictoryCount returns the number of players who achieved victory.
func (c *Convoy) GetVictoryCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	count := 0
	for _, run := range c.Runs {
		if run.IsVictory {
			count++
		}
	}
	return count
}
