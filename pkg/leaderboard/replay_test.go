package leaderboard

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewReplayInfo(t *testing.T) {
	entry := NewEntry(12345, engine.GenreFantasy, 1000, 30, 4)
	info := NewReplayInfo(entry)

	if info.Seed != 12345 {
		t.Errorf("expected seed 12345, got %d", info.Seed)
	}
	if info.Genre != engine.GenreFantasy {
		t.Errorf("expected genre fantasy, got %s", info.Genre)
	}
	if info.SourceRun == nil {
		t.Error("expected SourceRun to be set")
	}
}

func TestReplayManagerGetReplayableSeeds(t *testing.T) {
	board := NewBoard()
	_ = board.Add(NewEntry(100, engine.GenreFantasy, 1000, 30, 4))
	_ = board.Add(NewEntry(200, engine.GenreScifi, 800, 25, 3))
	_ = board.Add(NewEntry(300, engine.GenreHorror, 1200, 35, 5))

	rm := NewReplayManager(board, nil, nil)

	seeds := rm.GetReplayableSeeds(2)
	if len(seeds) != 2 {
		t.Fatalf("expected 2 seeds, got %d", len(seeds))
	}

	// Should be sorted by score
	if seeds[0].SourceRun.Score != 1200 {
		t.Errorf("expected first seed to have score 1200, got %d", seeds[0].SourceRun.Score)
	}
}

func TestReplayManagerGetReplayableSeedsByGenre(t *testing.T) {
	board := NewBoard()
	_ = board.Add(NewEntry(100, engine.GenreFantasy, 1000, 30, 4))
	_ = board.Add(NewEntry(200, engine.GenreFantasy, 800, 25, 3))
	_ = board.Add(NewEntry(300, engine.GenreScifi, 1200, 35, 5))

	rm := NewReplayManager(board, nil, nil)

	seeds := rm.GetReplayableSeedsByGenre(engine.GenreFantasy, 10)
	if len(seeds) != 2 {
		t.Fatalf("expected 2 fantasy seeds, got %d", len(seeds))
	}

	// Should be sorted by score
	if seeds[0].SourceRun.Score != 1000 {
		t.Errorf("expected top fantasy score 1000, got %d", seeds[0].SourceRun.Score)
	}
}

func TestReplayManagerGetReplayInfoForSeed(t *testing.T) {
	board := NewBoard()
	_ = board.Add(NewEntry(100, engine.GenreFantasy, 500, 30, 2))
	_ = board.Add(NewEntry(100, engine.GenreScifi, 800, 25, 3))

	rm := NewReplayManager(board, nil, nil)

	info := rm.GetReplayInfoForSeed(100)
	if info == nil {
		t.Fatal("expected non-nil replay info")
	}
	if info.Seed != 100 {
		t.Errorf("expected seed 100, got %d", info.Seed)
	}
	// Should return the top-scoring entry
	if info.SourceRun.Score != 800 {
		t.Errorf("expected top score 800, got %d", info.SourceRun.Score)
	}
}

func TestReplayManagerGetReplayInfoForUnknownSeed(t *testing.T) {
	board := NewBoard()
	rm := NewReplayManager(board, nil, nil)

	info := rm.GetReplayInfoForSeed(99999)
	if info == nil {
		t.Fatal("expected non-nil replay info even for unknown seed")
	}
	if info.Seed != 99999 {
		t.Errorf("expected seed 99999, got %d", info.Seed)
	}
	if info.SourceRun != nil {
		t.Error("expected nil SourceRun for unknown seed")
	}
}

func TestReplayManagerGetTopScoreForSeed(t *testing.T) {
	board := NewBoard()
	_ = board.Add(NewEntry(100, engine.GenreFantasy, 500, 30, 2))
	_ = board.Add(NewEntry(100, engine.GenreScifi, 800, 25, 3))
	_ = board.Add(NewEntry(200, engine.GenreHorror, 1000, 35, 4))

	rm := NewReplayManager(board, nil, nil)

	score, exists := rm.GetTopScoreForSeed(100)
	if !exists {
		t.Error("expected score to exist")
	}
	if score != 800 {
		t.Errorf("expected top score 800, got %d", score)
	}

	_, exists = rm.GetTopScoreForSeed(999)
	if exists {
		t.Error("expected non-existent seed to return false")
	}
}

func TestReplayManagerGetTopScoreForSeedAndGenre(t *testing.T) {
	board := NewBoard()
	_ = board.Add(NewEntry(100, engine.GenreFantasy, 500, 30, 2))
	_ = board.Add(NewEntry(100, engine.GenreFantasy, 600, 28, 3))
	_ = board.Add(NewEntry(100, engine.GenreScifi, 800, 25, 4))

	rm := NewReplayManager(board, nil, nil)

	score, exists := rm.GetTopScoreForSeedAndGenre(100, engine.GenreFantasy)
	if !exists {
		t.Error("expected score to exist")
	}
	if score != 600 {
		t.Errorf("expected top fantasy score 600, got %d", score)
	}
}

func TestReplayManagerIsNewHighScore(t *testing.T) {
	board := NewBoard()
	_ = board.Add(NewEntry(100, engine.GenreFantasy, 500, 30, 2))

	rm := NewReplayManager(board, nil, nil)

	// Score that beats existing
	if !rm.IsNewHighScore(100, engine.GenreFantasy, 600) {
		t.Error("expected 600 to be a new high score")
	}

	// Score that doesn't beat existing
	if rm.IsNewHighScore(100, engine.GenreFantasy, 400) {
		t.Error("expected 400 to not be a new high score")
	}

	// Score for new seed should always be a high score
	if !rm.IsNewHighScore(999, engine.GenreFantasy, 100) {
		t.Error("expected any score for new seed to be high score")
	}
}

func TestReplayManagerGetRunCountForSeed(t *testing.T) {
	board := NewBoard()
	_ = board.Add(NewEntry(100, engine.GenreFantasy, 500, 30, 2))
	_ = board.Add(NewEntry(100, engine.GenreScifi, 600, 28, 3))
	_ = board.Add(NewEntry(200, engine.GenreHorror, 800, 25, 4))

	rm := NewReplayManager(board, nil, nil)

	count := rm.GetRunCountForSeed(100)
	if count != 2 {
		t.Errorf("expected 2 runs for seed 100, got %d", count)
	}

	count = rm.GetRunCountForSeed(200)
	if count != 1 {
		t.Errorf("expected 1 run for seed 200, got %d", count)
	}

	count = rm.GetRunCountForSeed(999)
	if count != 0 {
		t.Errorf("expected 0 runs for unknown seed, got %d", count)
	}
}

func TestReplayManagerGetUniqueSeeds(t *testing.T) {
	board := NewBoard()
	_ = board.Add(NewEntry(100, engine.GenreFantasy, 500, 30, 2))
	_ = board.Add(NewEntry(100, engine.GenreScifi, 600, 28, 3))
	_ = board.Add(NewEntry(200, engine.GenreHorror, 800, 25, 4))
	_ = board.Add(NewEntry(300, engine.GenreCyberpunk, 700, 32, 3))

	rm := NewReplayManager(board, nil, nil)

	seeds := rm.GetUniqueSeeds()
	if len(seeds) != 3 {
		t.Errorf("expected 3 unique seeds, got %d", len(seeds))
	}

	// Verify all seeds are present
	seedSet := make(map[int64]bool)
	for _, s := range seeds {
		seedSet[s] = true
	}
	if !seedSet[100] || !seedSet[200] || !seedSet[300] {
		t.Error("missing expected seed in unique seeds")
	}
}

func TestReplayManagerValidateSeed(t *testing.T) {
	rm := NewReplayManager(nil, nil, nil)

	if !rm.ValidateSeed(0) {
		t.Error("expected seed 0 to be valid")
	}
	if !rm.ValidateSeed(-12345) {
		t.Error("expected negative seed to be valid")
	}
	if !rm.ValidateSeed(9999999) {
		t.Error("expected large seed to be valid")
	}
}

func TestReplayManagerGetChallengeSeeds(t *testing.T) {
	board := NewBoard()
	// Add multiple runs for seed 100 (popular seed)
	_ = board.Add(NewEntry(100, engine.GenreFantasy, 500, 30, 2))
	_ = board.Add(NewEntry(100, engine.GenreFantasy, 600, 28, 3))
	_ = board.Add(NewEntry(100, engine.GenreFantasy, 700, 26, 4))
	// Add single run for seed 200
	_ = board.Add(NewEntry(200, engine.GenreScifi, 1000, 25, 5))
	// Add two runs for seed 300
	_ = board.Add(NewEntry(300, engine.GenreHorror, 800, 32, 3))
	_ = board.Add(NewEntry(300, engine.GenreHorror, 900, 30, 4))

	rm := NewReplayManager(board, nil, nil)

	challenges := rm.GetChallengeSeeds(2)
	if len(challenges) != 2 {
		t.Fatalf("expected 2 challenge seeds, got %d", len(challenges))
	}

	// First should be seed 100 (3 attempts)
	if challenges[0].Seed != 100 {
		t.Errorf("expected first challenge seed to be 100 (most attempts), got %d", challenges[0].Seed)
	}

	// Second should be seed 300 (2 attempts)
	if challenges[1].Seed != 300 {
		t.Errorf("expected second challenge seed to be 300, got %d", challenges[1].Seed)
	}
}

func TestReplayManagerNilBoard(t *testing.T) {
	rm := NewReplayManager(nil, nil, nil)

	// All operations should handle nil board gracefully
	seeds := rm.GetReplayableSeeds(10)
	if seeds != nil {
		t.Error("expected nil for nil board")
	}

	info := rm.GetReplayInfoForSeed(100)
	if info != nil {
		t.Error("expected nil info for nil board")
	}

	score, exists := rm.GetTopScoreForSeed(100)
	if exists || score != 0 {
		t.Error("expected no score for nil board")
	}

	count := rm.GetRunCountForSeed(100)
	if count != 0 {
		t.Error("expected 0 count for nil board")
	}

	unique := rm.GetUniqueSeeds()
	if unique != nil {
		t.Error("expected nil unique seeds for nil board")
	}

	challenges := rm.GetChallengeSeeds(10)
	if challenges != nil {
		t.Error("expected nil challenges for nil board")
	}
}

func TestReplayInfoStruct(t *testing.T) {
	entry := NewEntry(42, engine.GenreCyberpunk, 1500, 40, 5)
	entry.WithPlayer("p1", "TestPlayer")

	info := NewReplayInfo(entry)

	if info.Seed != 42 {
		t.Errorf("expected seed 42, got %d", info.Seed)
	}
	if info.Genre != engine.GenreCyberpunk {
		t.Errorf("expected genre cyberpunk, got %s", info.Genre)
	}
	if info.SourceRun.PlayerName != "TestPlayer" {
		t.Errorf("expected player name TestPlayer, got %s", info.SourceRun.PlayerName)
	}
}
