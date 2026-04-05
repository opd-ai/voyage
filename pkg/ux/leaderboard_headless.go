//go:build headless

package ux

import (
	"fmt"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/leaderboard"
)

// LeaderboardScreen displays per-seed leaderboard entries with filtering.
// This is the headless version for testing.
type LeaderboardScreen struct {
	genre        engine.GenreID
	screenWidth  int
	screenHeight int

	// Data
	board   *leaderboard.Board
	entries []*leaderboard.Entry

	// Filters
	filterGenre *engine.GenreID
	filterSeed  *int64

	// Selection state
	selectedIndex int
	scrollOffset  int
	maxVisible    int

	// Tab state
	filterTab int // 0 = all, 1 = by genre, 2 = by seed

	// Available genres for filtering
	availableGenres []engine.GenreID
	genreIndex      int
}

// NewLeaderboardScreen creates a new leaderboard display screen.
func NewLeaderboardScreen(genre engine.GenreID, screenWidth, screenHeight int) *LeaderboardScreen {
	ls := &LeaderboardScreen{
		genre:           genre,
		screenWidth:     screenWidth,
		screenHeight:    screenHeight,
		selectedIndex:   0,
		scrollOffset:    0,
		maxVisible:      10,
		filterTab:       0,
		availableGenres: engine.AllGenres(),
		genreIndex:      0,
	}
	return ls
}

// SetGenre changes the screen's visual theme.
func (ls *LeaderboardScreen) SetGenre(genre engine.GenreID) {
	ls.genre = genre
}

// SetBoard sets the leaderboard data to display.
func (ls *LeaderboardScreen) SetBoard(board *leaderboard.Board) {
	ls.board = board
	ls.refreshEntries()
}

// refreshEntries updates the displayed entries based on current filters.
func (ls *LeaderboardScreen) refreshEntries() {
	if ls.board == nil {
		ls.entries = nil
		return
	}

	switch ls.filterTab {
	case 0: // All entries
		ls.entries = ls.board.GetTopN(100)
	case 1: // By genre
		if ls.filterGenre != nil {
			ls.entries = ls.board.GetByGenre(*ls.filterGenre)
		} else {
			ls.entries = ls.board.GetAll()
		}
	case 2: // By seed
		if ls.filterSeed != nil {
			ls.entries = ls.board.GetBySeed(*ls.filterSeed)
		} else {
			ls.entries = ls.board.GetAll()
		}
	}

	// Reset selection if out of bounds
	if ls.selectedIndex >= len(ls.entries) {
		ls.selectedIndex = 0
	}
	ls.scrollOffset = 0
}

// SetFilterGenre sets the genre filter.
func (ls *LeaderboardScreen) SetFilterGenre(genre *engine.GenreID) {
	ls.filterGenre = genre
	ls.refreshEntries()
}

// SetFilterSeed sets the seed filter.
func (ls *LeaderboardScreen) SetFilterSeed(seed *int64) {
	ls.filterSeed = seed
	ls.refreshEntries()
}

// ClearFilters removes all filters.
func (ls *LeaderboardScreen) ClearFilters() {
	ls.filterGenre = nil
	ls.filterSeed = nil
	ls.filterTab = 0
	ls.refreshEntries()
}

// CycleFilterTab cycles through filter tabs.
func (ls *LeaderboardScreen) CycleFilterTab() {
	ls.filterTab = (ls.filterTab + 1) % 3
	switch ls.filterTab {
	case 0:
		ls.filterGenre = nil
		ls.filterSeed = nil
	case 1:
		g := ls.availableGenres[ls.genreIndex]
		ls.filterGenre = &g
		ls.filterSeed = nil
	case 2:
		ls.filterGenre = nil
		// Seed filter is set externally
	}
	ls.refreshEntries()
}

// CycleGenreFilter cycles through available genres in genre filter mode.
func (ls *LeaderboardScreen) CycleGenreFilter() {
	if ls.filterTab == 1 {
		ls.genreIndex = (ls.genreIndex + 1) % len(ls.availableGenres)
		g := ls.availableGenres[ls.genreIndex]
		ls.filterGenre = &g
		ls.refreshEntries()
	}
}

// SelectNext moves selection down.
func (ls *LeaderboardScreen) SelectNext() {
	if len(ls.entries) == 0 {
		return
	}
	ls.selectedIndex++
	if ls.selectedIndex >= len(ls.entries) {
		ls.selectedIndex = len(ls.entries) - 1
	}
	// Adjust scroll if needed
	if ls.selectedIndex >= ls.scrollOffset+ls.maxVisible {
		ls.scrollOffset = ls.selectedIndex - ls.maxVisible + 1
	}
}

// SelectPrev moves selection up.
func (ls *LeaderboardScreen) SelectPrev() {
	if len(ls.entries) == 0 {
		return
	}
	ls.selectedIndex--
	if ls.selectedIndex < 0 {
		ls.selectedIndex = 0
	}
	// Adjust scroll if needed
	if ls.selectedIndex < ls.scrollOffset {
		ls.scrollOffset = ls.selectedIndex
	}
}

// SelectedEntry returns the currently selected entry.
func (ls *LeaderboardScreen) SelectedEntry() *leaderboard.Entry {
	if ls.selectedIndex >= 0 && ls.selectedIndex < len(ls.entries) {
		return ls.entries[ls.selectedIndex]
	}
	return nil
}

// GetSelectedSeed returns the seed of the selected entry for replay.
func (ls *LeaderboardScreen) GetSelectedSeed() (int64, bool) {
	entry := ls.SelectedEntry()
	if entry != nil {
		return entry.Seed, true
	}
	return 0, false
}

// EntryCount returns the number of displayed entries.
func (ls *LeaderboardScreen) EntryCount() int {
	return len(ls.entries)
}

// TotalEntryCount returns total entries in the board.
func (ls *LeaderboardScreen) TotalEntryCount() int {
	if ls.board == nil {
		return 0
	}
	return ls.board.Count()
}

// GetFilterInfo returns current filter information.
func (ls *LeaderboardScreen) GetFilterInfo() (filterType, filterValue string) {
	switch ls.filterTab {
	case 0:
		return "all", ""
	case 1:
		if ls.filterGenre != nil {
			return "genre", string(*ls.filterGenre)
		}
		return "genre", ""
	case 2:
		if ls.filterSeed != nil {
			return "seed", fmt.Sprintf("%d", *ls.filterSeed)
		}
		return "seed", ""
	}
	return "unknown", ""
}
