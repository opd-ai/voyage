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

	nameMap := c.buildPlayerNameMap()
	result.Rankings = c.buildRankings(nameMap)
	c.sortAndAssignRanks(result.Rankings)
	result.WinnerID = c.findWinner(result.Rankings)
	result.FastestID = c.findFastest(result.Rankings)
	result.FirstToFinish = c.findFirstToFinish()
	result.Comparisons = c.generateStatComparisons(nameMap)

	return result
}

// buildPlayerNameMap creates a map of player IDs to names.
func (c *Convoy) buildPlayerNameMap() map[PlayerID]string {
	nameMap := make(map[PlayerID]string)
	for _, p := range c.Players {
		nameMap[p.ID] = p.Name
	}
	return nameMap
}

// buildRankings creates player rankings from run data.
func (c *Convoy) buildRankings(nameMap map[PlayerID]string) []*PlayerRanking {
	rankings := make([]*PlayerRanking, 0, len(c.Runs))
	for _, run := range c.Runs {
		rankings = append(rankings, &PlayerRanking{
			PlayerID:  run.PlayerID,
			Name:      nameMap[run.PlayerID],
			Score:     run.Score,
			Victory:   run.IsVictory,
			Days:      run.DaysTraveled,
			Survivors: run.Survivors,
		})
	}
	return rankings
}

// sortAndAssignRanks sorts rankings by score/days and assigns rank numbers.
func (c *Convoy) sortAndAssignRanks(rankings []*PlayerRanking) {
	sort.Slice(rankings, func(i, j int) bool {
		if rankings[i].Score != rankings[j].Score {
			return rankings[i].Score > rankings[j].Score
		}
		return rankings[i].Days < rankings[j].Days
	})
	for i, r := range rankings {
		r.Rank = i + 1
	}
}

// findWinner returns the player ID of the highest-scoring victor.
func (c *Convoy) findWinner(rankings []*PlayerRanking) PlayerID {
	for _, r := range rankings {
		if r.Victory {
			return r.PlayerID
		}
	}
	return ""
}

// findFastest returns the player ID with fewest days among victors.
func (c *Convoy) findFastest(rankings []*PlayerRanking) PlayerID {
	fastestDays := int(^uint(0) >> 1)
	var fastestID PlayerID
	for _, r := range rankings {
		if r.Victory && r.Days < fastestDays {
			fastestDays = r.Days
			fastestID = r.PlayerID
		}
	}
	return fastestID
}

// findFirstToFinish returns the player ID who finished their run first.
func (c *Convoy) findFirstToFinish() PlayerID {
	if len(c.Runs) == 0 {
		return ""
	}
	earliestFinish := c.Runs[0].FinishedAt
	firstID := c.Runs[0].PlayerID
	for _, run := range c.Runs {
		if !run.FinishedAt.IsZero() && run.FinishedAt.Before(earliestFinish) {
			earliestFinish = run.FinishedAt
			firstID = run.PlayerID
		}
	}
	return firstID
}

// generateStatComparisons creates detailed comparisons for each stat.
func (c *Convoy) generateStatComparisons(nameMap map[PlayerID]string) []*StatComparison {
	comparisons := make([]*StatComparison, 0, 4)
	comparisons = append(comparisons, c.buildScoreComparison(nameMap))
	comparisons = append(comparisons, c.buildDaysComparison(nameMap))
	comparisons = append(comparisons, c.buildSurvivorComparison(nameMap))
	comparisons = append(comparisons, c.buildEventsComparison(nameMap))
	return comparisons
}

// buildScoreComparison creates a score comparison across all runs.
func (c *Convoy) buildScoreComparison(nameMap map[PlayerID]string) *StatComparison {
	comp := &StatComparison{Stat: "Score", Values: make(map[PlayerID]int)}
	var bestScore int
	for _, run := range c.Runs {
		comp.Values[run.PlayerID] = run.Score
		if run.Score > bestScore {
			bestScore = run.Score
			comp.BestID = run.PlayerID
			comp.BestName = nameMap[run.PlayerID]
		}
	}
	return comp
}

// buildDaysComparison creates a days-traveled comparison (lower is better).
func (c *Convoy) buildDaysComparison(nameMap map[PlayerID]string) *StatComparison {
	comp := &StatComparison{Stat: "Days Traveled", Values: make(map[PlayerID]int)}
	bestDays := int(^uint(0) >> 1)
	for _, run := range c.Runs {
		comp.Values[run.PlayerID] = run.DaysTraveled
		if run.DaysTraveled < bestDays && run.DaysTraveled > 0 {
			bestDays = run.DaysTraveled
			comp.BestID = run.PlayerID
			comp.BestName = nameMap[run.PlayerID]
		}
	}
	return comp
}

// buildSurvivorComparison creates a survivors comparison across all runs.
func (c *Convoy) buildSurvivorComparison(nameMap map[PlayerID]string) *StatComparison {
	comp := &StatComparison{Stat: "Survivors", Values: make(map[PlayerID]int)}
	var bestSurvivors int
	for _, run := range c.Runs {
		comp.Values[run.PlayerID] = run.Survivors
		if run.Survivors > bestSurvivors {
			bestSurvivors = run.Survivors
			comp.BestID = run.PlayerID
			comp.BestName = nameMap[run.PlayerID]
		}
	}
	return comp
}

// buildEventsComparison creates an events-seen comparison across all runs.
func (c *Convoy) buildEventsComparison(nameMap map[PlayerID]string) *StatComparison {
	comp := &StatComparison{Stat: "Events Faced", Values: make(map[PlayerID]int)}
	var mostEvents int
	for _, run := range c.Runs {
		comp.Values[run.PlayerID] = run.EventsSeen
		if run.EventsSeen > mostEvents {
			mostEvents = run.EventsSeen
			comp.BestID = run.PlayerID
			comp.BestName = nameMap[run.PlayerID]
		}
	}
	return comp
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

	firstRun := c.findEarliestFinishedRun()
	if firstRun == nil {
		return nil, nil
	}
	return c.findPlayerByID(firstRun.PlayerID), firstRun
}

// findEarliestFinishedRun returns the run with the earliest finish time.
func (c *Convoy) findEarliestFinishedRun() *RunData {
	var firstRun *RunData
	for _, run := range c.Runs {
		if run.FinishedAt.IsZero() {
			continue
		}
		if firstRun == nil || run.FinishedAt.Before(firstRun.FinishedAt) {
			firstRun = run
		}
	}
	return firstRun
}

// findPlayerByID looks up a player by ID from the convoy's player list.
func (c *Convoy) findPlayerByID(id PlayerID) *Player {
	for _, p := range c.Players {
		if p.ID == id {
			return p
		}
	}
	return nil
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
