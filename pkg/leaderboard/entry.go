package leaderboard

import (
	"encoding/json"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/opd-ai/voyage/pkg/engine"
)

// Entry represents a single leaderboard entry for a completed run.
type Entry struct {
	// Run identification
	Seed  int64          `json:"seed"`
	Genre engine.GenreID `json:"genre"`

	// Player identification (anonymous by default)
	PlayerID   string `json:"playerId,omitempty"`
	PlayerName string `json:"playerName,omitempty"`

	// Run results
	Score     int `json:"score"`
	Days      int `json:"days"`
	Survivors int `json:"survivors"`

	// Metadata
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
}

// NewEntry creates a new leaderboard entry from run completion data.
func NewEntry(seed int64, genre engine.GenreID, score, days, survivors int) *Entry {
	return &Entry{
		Seed:      seed,
		Genre:     genre,
		Score:     score,
		Days:      days,
		Survivors: survivors,
		Timestamp: time.Now().UTC(),
	}
}

// WithPlayer sets player identification on the entry.
func (e *Entry) WithPlayer(id, name string) *Entry {
	e.PlayerID = id
	e.PlayerName = name
	return e
}

// WithVersion sets the game version on the entry.
func (e *Entry) WithVersion(version string) *Entry {
	e.Version = version
	return e
}

// Marshal serializes the entry to JSON.
func (e *Entry) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

// UnmarshalEntry deserializes a JSON entry.
func UnmarshalEntry(data []byte) (*Entry, error) {
	var entry Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, err
	}
	return &entry, nil
}

// Validate checks if the entry contains valid data.
func (e *Entry) Validate() error {
	if e.Seed == 0 {
		return errors.New("seed cannot be zero")
	}
	if !engine.IsValidGenre(string(e.Genre)) {
		return errors.New("invalid genre")
	}
	if e.Score < 0 {
		return errors.New("score cannot be negative")
	}
	if e.Days < 0 {
		return errors.New("days cannot be negative")
	}
	if e.Survivors < 0 {
		return errors.New("survivors cannot be negative")
	}
	return nil
}

// Board holds leaderboard entries and provides query methods.
type Board struct {
	mu      sync.RWMutex
	entries []*Entry
	// Indices for fast lookups
	bySeed  map[int64][]*Entry
	byGenre map[engine.GenreID][]*Entry
}

// NewBoard creates a new empty leaderboard.
func NewBoard() *Board {
	return &Board{
		entries: make([]*Entry, 0),
		bySeed:  make(map[int64][]*Entry),
		byGenre: make(map[engine.GenreID][]*Entry),
	}
}

// Add adds an entry to the leaderboard.
func (b *Board) Add(entry *Entry) error {
	if err := entry.Validate(); err != nil {
		return err
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.entries = append(b.entries, entry)
	b.bySeed[entry.Seed] = append(b.bySeed[entry.Seed], entry)
	b.byGenre[entry.Genre] = append(b.byGenre[entry.Genre], entry)

	return nil
}

// GetAll returns all entries sorted by score (descending).
func (b *Board) GetAll() []*Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := make([]*Entry, len(b.entries))
	copy(result, b.entries)
	sortByScore(result)
	return result
}

// copyAndSortEntries returns a sorted copy of the given entries slice.
func copyAndSortEntries(entries []*Entry) []*Entry {
	result := make([]*Entry, len(entries))
	copy(result, entries)
	sortByScore(result)
	return result
}

// GetBySeed returns all entries for a specific seed, sorted by score.
func (b *Board) GetBySeed(seed int64) []*Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return copyAndSortEntries(b.bySeed[seed])
}

// GetByGenre returns all entries for a specific genre, sorted by score.
func (b *Board) GetByGenre(genre engine.GenreID) []*Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return copyAndSortEntries(b.byGenre[genre])
}

// GetBySeedAndGenre returns entries matching both seed and genre.
func (b *Board) GetBySeedAndGenre(seed int64, genre engine.GenreID) []*Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var result []*Entry
	for _, e := range b.bySeed[seed] {
		if e.Genre == genre {
			result = append(result, e)
		}
	}
	sortByScore(result)
	return result
}

// GetTopN returns the top N entries sorted by score.
func (b *Board) GetTopN(n int) []*Entry {
	all := b.GetAll()
	if n > len(all) {
		n = len(all)
	}
	return all[:n]
}

// GetTopNBySeed returns the top N entries for a specific seed.
func (b *Board) GetTopNBySeed(seed int64, n int) []*Entry {
	entries := b.GetBySeed(seed)
	if n > len(entries) {
		n = len(entries)
	}
	return entries[:n]
}

// GetTopNByGenre returns the top N entries for a specific genre.
func (b *Board) GetTopNByGenre(genre engine.GenreID, n int) []*Entry {
	entries := b.GetByGenre(genre)
	if n > len(entries) {
		n = len(entries)
	}
	return entries[:n]
}

// Count returns the total number of entries.
func (b *Board) Count() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.entries)
}

// CountBySeed returns the number of entries for a specific seed.
func (b *Board) CountBySeed(seed int64) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.bySeed[seed])
}

// CountByGenre returns the number of entries for a specific genre.
func (b *Board) CountByGenre(genre engine.GenreID) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.byGenre[genre])
}

// GetUniqueSeedCount returns the number of unique seeds in the board.
func (b *Board) GetUniqueSeedCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.bySeed)
}

// GetTopSeed returns the seed with the highest top score.
func (b *Board) GetTopSeed() (int64, int) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var topSeed int64
	topScore := -1
	for seed, entries := range b.bySeed {
		for _, e := range entries {
			if e.Score > topScore {
				topScore = e.Score
				topSeed = seed
			}
		}
	}
	return topSeed, topScore
}

// GetReplayableSeed finds a top-scoring seed that can be replayed.
func (b *Board) GetReplayableSeed(genre engine.GenreID) (int64, *Entry) {
	entries := b.GetTopNByGenre(genre, 1)
	if len(entries) == 0 {
		return 0, nil
	}
	return entries[0].Seed, entries[0]
}

// Marshal serializes the entire board to JSON.
func (b *Board) Marshal() ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return json.Marshal(b.entries)
}

// Unmarshal deserializes JSON data into a board.
func Unmarshal(data []byte) (*Board, error) {
	var entries []*Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	board := NewBoard()
	for _, e := range entries {
		if err := board.Add(e); err != nil {
			continue // Skip invalid entries
		}
	}
	return board, nil
}

// sortByScore sorts entries by score in descending order.
func sortByScore(entries []*Entry) {
	sort.Slice(entries, func(i, j int) bool {
		// Primary: score descending
		if entries[i].Score != entries[j].Score {
			return entries[i].Score > entries[j].Score
		}
		// Secondary: fewer days is better
		if entries[i].Days != entries[j].Days {
			return entries[i].Days < entries[j].Days
		}
		// Tertiary: more survivors is better
		return entries[i].Survivors > entries[j].Survivors
	})
}
