package leaderboard

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// ReplayInfo contains the information needed to replay a run.
type ReplayInfo struct {
	Seed      int64          `json:"seed"`
	Genre     engine.GenreID `json:"genre"`
	SourceRun *Entry         `json:"sourceRun,omitempty"`
}

// NewReplayInfo creates replay info from a leaderboard entry.
func NewReplayInfo(entry *Entry) *ReplayInfo {
	return &ReplayInfo{
		Seed:      entry.Seed,
		Genre:     entry.Genre,
		SourceRun: entry,
	}
}

// ReplayManager handles seed replay functionality.
type ReplayManager struct {
	board   *Board
	client  *Client
	storage *LocalStorage
}

// NewReplayManager creates a new replay manager.
func NewReplayManager(board *Board, client *Client, storage *LocalStorage) *ReplayManager {
	return &ReplayManager{
		board:   board,
		client:  client,
		storage: storage,
	}
}

// GetReplayableSeeds returns seeds that can be replayed, sorted by score.
func (rm *ReplayManager) GetReplayableSeeds(limit int) []*ReplayInfo {
	if rm.board == nil {
		return nil
	}

	entries := rm.board.GetTopN(limit)
	result := make([]*ReplayInfo, len(entries))
	for i, e := range entries {
		result[i] = NewReplayInfo(e)
	}
	return result
}

// GetReplayableSeedsByGenre returns seeds for a specific genre.
func (rm *ReplayManager) GetReplayableSeedsByGenre(genre engine.GenreID, limit int) []*ReplayInfo {
	if rm.board == nil {
		return nil
	}

	entries := rm.board.GetTopNByGenre(genre, limit)
	result := make([]*ReplayInfo, len(entries))
	for i, e := range entries {
		result[i] = NewReplayInfo(e)
	}
	return result
}

// GetReplayInfoForSeed returns replay info for a specific seed.
func (rm *ReplayManager) GetReplayInfoForSeed(seed int64) *ReplayInfo {
	if rm.board == nil {
		return nil
	}

	entries := rm.board.GetBySeed(seed)
	if len(entries) == 0 {
		// Return basic info even if no run exists for this seed
		return &ReplayInfo{
			Seed:  seed,
			Genre: engine.GenreFantasy, // Default genre
		}
	}

	// Return info from the top-scoring run for this seed
	return NewReplayInfo(entries[0])
}

// GetTopScoreForSeed returns the top score achieved on a specific seed.
func (rm *ReplayManager) GetTopScoreForSeed(seed int64) (int, bool) {
	if rm.board == nil {
		return 0, false
	}

	entries := rm.board.GetBySeed(seed)
	if len(entries) == 0 {
		return 0, false
	}
	return entries[0].Score, true
}

// GetTopScoreForSeedAndGenre returns the top score for a specific seed/genre combo.
func (rm *ReplayManager) GetTopScoreForSeedAndGenre(seed int64, genre engine.GenreID) (int, bool) {
	if rm.board == nil {
		return 0, false
	}

	entries := rm.board.GetBySeedAndGenre(seed, genre)
	if len(entries) == 0 {
		return 0, false
	}
	return entries[0].Score, true
}

// FetchTopSeedFromServer fetches the top-scoring seed from the server.
func (rm *ReplayManager) FetchTopSeedFromServer(genre engine.GenreID) (*ReplayInfo, error) {
	if rm.client == nil {
		return nil, ErrServerUnavailable
	}

	seed, entry, err := rm.client.GetReplayableSeed(genre)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	return &ReplayInfo{
		Seed:      seed,
		Genre:     genre,
		SourceRun: entry,
	}, nil
}

// IsNewHighScore checks if a score beats the current top for a seed/genre.
func (rm *ReplayManager) IsNewHighScore(seed int64, genre engine.GenreID, score int) bool {
	topScore, exists := rm.GetTopScoreForSeedAndGenre(seed, genre)
	if !exists {
		return true
	}
	return score > topScore
}

// GetRunCountForSeed returns how many runs exist for a specific seed.
func (rm *ReplayManager) GetRunCountForSeed(seed int64) int {
	if rm.board == nil {
		return 0
	}
	return rm.board.CountBySeed(seed)
}

// GetUniqueSeeds returns a list of all unique seeds in the leaderboard.
func (rm *ReplayManager) GetUniqueSeeds() []int64 {
	if rm.board == nil {
		return nil
	}

	// Get all entries and extract unique seeds
	entries := rm.board.GetAll()
	seenSeeds := make(map[int64]bool)
	var seeds []int64

	for _, e := range entries {
		if !seenSeeds[e.Seed] {
			seenSeeds[e.Seed] = true
			seeds = append(seeds, e.Seed)
		}
	}

	return seeds
}

// ValidateSeed checks if a seed produces a valid game world.
// Always returns true since any int64 is a valid seed.
func (rm *ReplayManager) ValidateSeed(seed int64) bool {
	return true
}

// GetChallengeSeeds returns a curated list of "challenge" seeds.
// These are seeds with interesting characteristics (high scores, many attempts).
func (rm *ReplayManager) GetChallengeSeeds(limit int) []*ReplayInfo {
	if rm.board == nil {
		return nil
	}

	entries := rm.board.GetAll()
	stats := rm.buildSeedStats(entries)
	rm.sortSeedStatsByCount(stats)
	return rm.buildReplayInfos(stats, limit)
}

// seedStats holds aggregated statistics for a single seed.
type seedStats struct {
	seed  int64
	count int
	entry *Entry
}

// buildSeedStats aggregates attempt counts and top scores per seed.
func (rm *ReplayManager) buildSeedStats(entries []*Entry) []seedStats {
	seedCounts := make(map[int64]int)
	seedTopScore := make(map[int64]*Entry)

	for _, e := range entries {
		seedCounts[e.Seed]++
		if current, exists := seedTopScore[e.Seed]; !exists || e.Score > current.Score {
			seedTopScore[e.Seed] = e
		}
	}

	stats := make([]seedStats, 0, len(seedCounts))
	for seed, count := range seedCounts {
		stats = append(stats, seedStats{seed, count, seedTopScore[seed]})
	}
	return stats
}

// sortSeedStatsByCount sorts seed stats by attempt count descending.
func (rm *ReplayManager) sortSeedStatsByCount(stats []seedStats) {
	for i := 0; i < len(stats); i++ {
		for j := i + 1; j < len(stats); j++ {
			if stats[j].count > stats[i].count {
				stats[i], stats[j] = stats[j], stats[i]
			}
		}
	}
}

// buildReplayInfos converts seed stats to ReplayInfo slice.
func (rm *ReplayManager) buildReplayInfos(stats []seedStats, limit int) []*ReplayInfo {
	if limit > len(stats) {
		limit = len(stats)
	}
	result := make([]*ReplayInfo, limit)
	for i := 0; i < limit; i++ {
		result[i] = NewReplayInfo(stats[i].entry)
	}
	return result
}
