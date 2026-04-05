package ux

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/leaderboard"
)

func TestNewLeaderboardScreen(t *testing.T) {
	ls := NewLeaderboardScreen(engine.GenreFantasy, 800, 600)

	if ls == nil {
		t.Fatal("NewLeaderboardScreen returned nil")
	}

	if ls.EntryCount() != 0 {
		t.Errorf("expected 0 entries initially, got %d", ls.EntryCount())
	}
}

func TestLeaderboardScreenSetBoard(t *testing.T) {
	ls := NewLeaderboardScreen(engine.GenreFantasy, 800, 600)

	board := leaderboard.NewBoard()
	_ = board.Add(leaderboard.NewEntry(100, engine.GenreFantasy, 1000, 30, 4))
	_ = board.Add(leaderboard.NewEntry(200, engine.GenreScifi, 800, 25, 3))
	_ = board.Add(leaderboard.NewEntry(300, engine.GenreHorror, 1200, 35, 5))

	ls.SetBoard(board)

	if ls.EntryCount() != 3 {
		t.Errorf("expected 3 entries, got %d", ls.EntryCount())
	}

	if ls.TotalEntryCount() != 3 {
		t.Errorf("expected total 3, got %d", ls.TotalEntryCount())
	}
}

func TestLeaderboardScreenSelection(t *testing.T) {
	ls := NewLeaderboardScreen(engine.GenreFantasy, 800, 600)

	board := leaderboard.NewBoard()
	_ = board.Add(leaderboard.NewEntry(100, engine.GenreFantasy, 1000, 30, 4))
	_ = board.Add(leaderboard.NewEntry(200, engine.GenreScifi, 800, 25, 3))
	_ = board.Add(leaderboard.NewEntry(300, engine.GenreHorror, 1200, 35, 5))
	ls.SetBoard(board)

	// Initial selection should be first entry (highest score)
	entry := ls.SelectedEntry()
	if entry == nil {
		t.Fatal("SelectedEntry returned nil")
	}
	if entry.Score != 1200 {
		t.Errorf("expected top score 1200, got %d", entry.Score)
	}

	// Move selection down
	ls.SelectNext()
	entry = ls.SelectedEntry()
	if entry.Score != 1000 {
		t.Errorf("expected second score 1000, got %d", entry.Score)
	}

	// Move back up
	ls.SelectPrev()
	entry = ls.SelectedEntry()
	if entry.Score != 1200 {
		t.Errorf("expected back to 1200, got %d", entry.Score)
	}
}

func TestLeaderboardScreenSelectionBounds(t *testing.T) {
	ls := NewLeaderboardScreen(engine.GenreFantasy, 800, 600)

	board := leaderboard.NewBoard()
	_ = board.Add(leaderboard.NewEntry(100, engine.GenreFantasy, 1000, 30, 4))
	_ = board.Add(leaderboard.NewEntry(200, engine.GenreScifi, 800, 25, 3))
	ls.SetBoard(board)

	// Try to go beyond bounds
	ls.SelectPrev() // Should stay at 0
	entry := ls.SelectedEntry()
	if entry.Score != 1000 {
		t.Errorf("expected to stay at first entry, got score %d", entry.Score)
	}

	ls.SelectNext()
	ls.SelectNext() // Should stay at last
	entry = ls.SelectedEntry()
	if entry.Score != 800 {
		t.Errorf("expected to stay at last entry, got score %d", entry.Score)
	}
}

func TestLeaderboardScreenFilterGenre(t *testing.T) {
	ls := NewLeaderboardScreen(engine.GenreFantasy, 800, 600)

	board := leaderboard.NewBoard()
	_ = board.Add(leaderboard.NewEntry(100, engine.GenreFantasy, 1000, 30, 4))
	_ = board.Add(leaderboard.NewEntry(200, engine.GenreScifi, 800, 25, 3))
	_ = board.Add(leaderboard.NewEntry(300, engine.GenreFantasy, 600, 35, 5))
	ls.SetBoard(board)

	// Set genre filter
	genre := engine.GenreFantasy
	ls.SetFilterGenre(&genre)

	// Should still show all because filter tab isn't set
	// Use CycleFilterTab to activate genre filter
	ls.CycleFilterTab() // Now in genre filter mode

	if ls.EntryCount() != 2 {
		t.Errorf("expected 2 fantasy entries, got %d", ls.EntryCount())
	}
}

func TestLeaderboardScreenFilterSeed(t *testing.T) {
	ls := NewLeaderboardScreen(engine.GenreFantasy, 800, 600)

	board := leaderboard.NewBoard()
	_ = board.Add(leaderboard.NewEntry(100, engine.GenreFantasy, 1000, 30, 4))
	_ = board.Add(leaderboard.NewEntry(100, engine.GenreScifi, 800, 25, 3))
	_ = board.Add(leaderboard.NewEntry(200, engine.GenreFantasy, 600, 35, 5))
	ls.SetBoard(board)

	// Cycle to seed filter mode
	ls.CycleFilterTab() // genre mode
	ls.CycleFilterTab() // seed mode

	seed := int64(100)
	ls.SetFilterSeed(&seed)

	if ls.EntryCount() != 2 {
		t.Errorf("expected 2 entries for seed 100, got %d", ls.EntryCount())
	}
}

func TestLeaderboardScreenClearFilters(t *testing.T) {
	ls := NewLeaderboardScreen(engine.GenreFantasy, 800, 600)

	board := leaderboard.NewBoard()
	_ = board.Add(leaderboard.NewEntry(100, engine.GenreFantasy, 1000, 30, 4))
	_ = board.Add(leaderboard.NewEntry(200, engine.GenreScifi, 800, 25, 3))
	ls.SetBoard(board)

	// Set filter
	ls.CycleFilterTab()

	// Clear filters
	ls.ClearFilters()

	filterType, _ := ls.GetFilterInfo()
	if filterType != "all" {
		t.Errorf("expected filter type 'all', got %s", filterType)
	}
}

func TestLeaderboardScreenGetSelectedSeed(t *testing.T) {
	ls := NewLeaderboardScreen(engine.GenreFantasy, 800, 600)

	board := leaderboard.NewBoard()
	_ = board.Add(leaderboard.NewEntry(12345, engine.GenreFantasy, 1000, 30, 4))
	ls.SetBoard(board)

	seed, ok := ls.GetSelectedSeed()
	if !ok {
		t.Fatal("expected GetSelectedSeed to return true")
	}
	if seed != 12345 {
		t.Errorf("expected seed 12345, got %d", seed)
	}
}

func TestLeaderboardScreenGetSelectedSeedEmpty(t *testing.T) {
	ls := NewLeaderboardScreen(engine.GenreFantasy, 800, 600)

	seed, ok := ls.GetSelectedSeed()
	if ok {
		t.Error("expected GetSelectedSeed to return false for empty board")
	}
	if seed != 0 {
		t.Errorf("expected seed 0, got %d", seed)
	}
}

func TestLeaderboardScreenSetGenre(t *testing.T) {
	ls := NewLeaderboardScreen(engine.GenreFantasy, 800, 600)

	ls.SetGenre(engine.GenreScifi)

	// No panic = success (headless version doesn't have skin)
}

func TestLeaderboardScreenCycleGenreFilter(t *testing.T) {
	ls := NewLeaderboardScreen(engine.GenreFantasy, 800, 600)

	board := leaderboard.NewBoard()
	_ = board.Add(leaderboard.NewEntry(100, engine.GenreFantasy, 1000, 30, 4))
	_ = board.Add(leaderboard.NewEntry(200, engine.GenreScifi, 800, 25, 3))
	ls.SetBoard(board)

	// Enter genre filter mode
	ls.CycleFilterTab()

	// Cycle through genres
	ls.CycleGenreFilter()

	filterType, _ := ls.GetFilterInfo()
	if filterType != "genre" {
		t.Errorf("expected filter type 'genre', got %s", filterType)
	}
}

func TestLeaderboardScreenFilterInfo(t *testing.T) {
	ls := NewLeaderboardScreen(engine.GenreFantasy, 800, 600)

	// Test all filter type
	filterType, filterValue := ls.GetFilterInfo()
	if filterType != "all" {
		t.Errorf("expected 'all', got %s", filterType)
	}
	if filterValue != "" {
		t.Errorf("expected empty value, got %s", filterValue)
	}

	// Test genre filter type
	ls.CycleFilterTab()
	filterType, filterValue = ls.GetFilterInfo()
	if filterType != "genre" {
		t.Errorf("expected 'genre', got %s", filterType)
	}

	// Test seed filter type
	ls.CycleFilterTab()
	filterType, _ = ls.GetFilterInfo()
	if filterType != "seed" {
		t.Errorf("expected 'seed', got %s", filterType)
	}
}

func TestLeaderboardScreenEmptySelection(t *testing.T) {
	ls := NewLeaderboardScreen(engine.GenreFantasy, 800, 600)

	// Test operations on empty screen
	ls.SelectNext()
	ls.SelectPrev()

	entry := ls.SelectedEntry()
	if entry != nil {
		t.Error("expected nil entry for empty board")
	}
}

func TestLeaderboardScreenNilBoard(t *testing.T) {
	ls := NewLeaderboardScreen(engine.GenreFantasy, 800, 600)

	if ls.TotalEntryCount() != 0 {
		t.Errorf("expected 0 for nil board, got %d", ls.TotalEntryCount())
	}
}

func TestLeaderboardScreenScrolling(t *testing.T) {
	ls := NewLeaderboardScreen(engine.GenreFantasy, 800, 600)

	board := leaderboard.NewBoard()
	// Add more entries than maxVisible (which is 10)
	for i := 0; i < 15; i++ {
		_ = board.Add(leaderboard.NewEntry(int64(i+1), engine.GenreFantasy, (15-i)*100, 30, 4))
	}
	ls.SetBoard(board)

	// Navigate down to trigger scrolling
	for i := 0; i < 12; i++ {
		ls.SelectNext()
	}

	// Selection should still work
	entry := ls.SelectedEntry()
	if entry == nil {
		t.Fatal("SelectedEntry should not be nil after scrolling")
	}

	// Navigate back up
	for i := 0; i < 12; i++ {
		ls.SelectPrev()
	}

	entry = ls.SelectedEntry()
	if entry == nil {
		t.Fatal("SelectedEntry should not be nil after scrolling back")
	}
}
