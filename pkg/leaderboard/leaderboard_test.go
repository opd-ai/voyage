package leaderboard

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewEntry(t *testing.T) {
	entry := NewEntry(12345, engine.GenreFantasy, 1000, 30, 4)

	if entry.Seed != 12345 {
		t.Errorf("expected seed 12345, got %d", entry.Seed)
	}
	if entry.Genre != engine.GenreFantasy {
		t.Errorf("expected genre fantasy, got %s", entry.Genre)
	}
	if entry.Score != 1000 {
		t.Errorf("expected score 1000, got %d", entry.Score)
	}
	if entry.Days != 30 {
		t.Errorf("expected days 30, got %d", entry.Days)
	}
	if entry.Survivors != 4 {
		t.Errorf("expected survivors 4, got %d", entry.Survivors)
	}
	if entry.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestEntryWithPlayer(t *testing.T) {
	entry := NewEntry(12345, engine.GenreScifi, 500, 20, 2)
	entry.WithPlayer("player-123", "TestPlayer")

	if entry.PlayerID != "player-123" {
		t.Errorf("expected player ID player-123, got %s", entry.PlayerID)
	}
	if entry.PlayerName != "TestPlayer" {
		t.Errorf("expected player name TestPlayer, got %s", entry.PlayerName)
	}
}

func TestEntryWithVersion(t *testing.T) {
	entry := NewEntry(12345, engine.GenreHorror, 750, 25, 3)
	entry.WithVersion("1.0.0")

	if entry.Version != "1.0.0" {
		t.Errorf("expected version 1.0.0, got %s", entry.Version)
	}
}

func TestEntryValidate(t *testing.T) {
	tests := []struct {
		name    string
		entry   *Entry
		wantErr bool
	}{
		{
			name:    "valid entry",
			entry:   NewEntry(12345, engine.GenreFantasy, 1000, 30, 4),
			wantErr: false,
		},
		{
			name:    "zero seed",
			entry:   NewEntry(0, engine.GenreFantasy, 1000, 30, 4),
			wantErr: true,
		},
		{
			name:    "invalid genre",
			entry:   NewEntry(12345, "invalid", 1000, 30, 4),
			wantErr: true,
		},
		{
			name:    "negative score",
			entry:   NewEntry(12345, engine.GenreScifi, -100, 30, 4),
			wantErr: true,
		},
		{
			name:    "negative days",
			entry:   NewEntry(12345, engine.GenreHorror, 1000, -5, 4),
			wantErr: true,
		},
		{
			name:    "negative survivors",
			entry:   NewEntry(12345, engine.GenreCyberpunk, 1000, 30, -1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entry.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEntryMarshalUnmarshal(t *testing.T) {
	entry := NewEntry(12345, engine.GenrePostapoc, 1500, 45, 5)
	entry.WithPlayer("p1", "Player One")
	entry.WithVersion("2.0.0")

	data, err := entry.Marshal()
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	restored, err := UnmarshalEntry(data)
	if err != nil {
		t.Fatalf("UnmarshalEntry() error = %v", err)
	}

	if restored.Seed != entry.Seed {
		t.Errorf("seed mismatch: got %d, want %d", restored.Seed, entry.Seed)
	}
	if restored.Genre != entry.Genre {
		t.Errorf("genre mismatch: got %s, want %s", restored.Genre, entry.Genre)
	}
	if restored.Score != entry.Score {
		t.Errorf("score mismatch: got %d, want %d", restored.Score, entry.Score)
	}
	if restored.Days != entry.Days {
		t.Errorf("days mismatch: got %d, want %d", restored.Days, entry.Days)
	}
	if restored.Survivors != entry.Survivors {
		t.Errorf("survivors mismatch: got %d, want %d", restored.Survivors, entry.Survivors)
	}
	if restored.PlayerID != entry.PlayerID {
		t.Errorf("playerID mismatch: got %s, want %s", restored.PlayerID, entry.PlayerID)
	}
}

func TestBoardAdd(t *testing.T) {
	board := NewBoard()

	entry := NewEntry(12345, engine.GenreFantasy, 1000, 30, 4)
	err := board.Add(entry)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	if board.Count() != 1 {
		t.Errorf("expected count 1, got %d", board.Count())
	}
}

func TestBoardAddInvalid(t *testing.T) {
	board := NewBoard()

	entry := NewEntry(0, engine.GenreFantasy, 1000, 30, 4)
	err := board.Add(entry)
	if err == nil {
		t.Error("expected error for invalid entry")
	}
}

func TestBoardGetAll(t *testing.T) {
	board := NewBoard()

	// Add entries with different scores
	_ = board.Add(NewEntry(1, engine.GenreFantasy, 500, 30, 2))
	_ = board.Add(NewEntry(2, engine.GenreScifi, 1000, 25, 3))
	_ = board.Add(NewEntry(3, engine.GenreHorror, 750, 35, 4))

	all := board.GetAll()

	if len(all) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(all))
	}

	// Should be sorted by score descending
	if all[0].Score != 1000 {
		t.Errorf("expected first entry score 1000, got %d", all[0].Score)
	}
	if all[1].Score != 750 {
		t.Errorf("expected second entry score 750, got %d", all[1].Score)
	}
	if all[2].Score != 500 {
		t.Errorf("expected third entry score 500, got %d", all[2].Score)
	}
}

func TestBoardGetBySeed(t *testing.T) {
	board := NewBoard()

	_ = board.Add(NewEntry(100, engine.GenreFantasy, 500, 30, 2))
	_ = board.Add(NewEntry(100, engine.GenreScifi, 800, 25, 3))
	_ = board.Add(NewEntry(200, engine.GenreHorror, 1000, 35, 4))

	entries := board.GetBySeed(100)

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries for seed 100, got %d", len(entries))
	}

	// Should be sorted by score
	if entries[0].Score != 800 {
		t.Errorf("expected first entry score 800, got %d", entries[0].Score)
	}
}

func TestBoardGetByGenre(t *testing.T) {
	board := NewBoard()

	_ = board.Add(NewEntry(1, engine.GenreFantasy, 500, 30, 2))
	_ = board.Add(NewEntry(2, engine.GenreFantasy, 800, 25, 3))
	_ = board.Add(NewEntry(3, engine.GenreScifi, 1000, 35, 4))

	entries := board.GetByGenre(engine.GenreFantasy)

	if len(entries) != 2 {
		t.Fatalf("expected 2 fantasy entries, got %d", len(entries))
	}

	// Should be sorted by score
	if entries[0].Score != 800 {
		t.Errorf("expected first entry score 800, got %d", entries[0].Score)
	}
}

func TestBoardGetBySeedAndGenre(t *testing.T) {
	board := NewBoard()

	_ = board.Add(NewEntry(100, engine.GenreFantasy, 500, 30, 2))
	_ = board.Add(NewEntry(100, engine.GenreScifi, 800, 25, 3))
	_ = board.Add(NewEntry(100, engine.GenreFantasy, 600, 28, 4))
	_ = board.Add(NewEntry(200, engine.GenreFantasy, 1000, 35, 5))

	entries := board.GetBySeedAndGenre(100, engine.GenreFantasy)

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	// Should be sorted by score
	if entries[0].Score != 600 {
		t.Errorf("expected first entry score 600, got %d", entries[0].Score)
	}
}

func TestBoardGetTopN(t *testing.T) {
	board := NewBoard()

	for i := 1; i <= 10; i++ {
		_ = board.Add(NewEntry(int64(i), engine.GenreFantasy, i*100, 30, 2))
	}

	top3 := board.GetTopN(3)

	if len(top3) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(top3))
	}

	if top3[0].Score != 1000 {
		t.Errorf("expected top score 1000, got %d", top3[0].Score)
	}
	if top3[2].Score != 800 {
		t.Errorf("expected third score 800, got %d", top3[2].Score)
	}
}

func TestBoardCounts(t *testing.T) {
	board := NewBoard()

	_ = board.Add(NewEntry(100, engine.GenreFantasy, 500, 30, 2))
	_ = board.Add(NewEntry(100, engine.GenreScifi, 800, 25, 3))
	_ = board.Add(NewEntry(200, engine.GenreFantasy, 1000, 35, 4))

	if board.Count() != 3 {
		t.Errorf("expected total count 3, got %d", board.Count())
	}
	if board.CountBySeed(100) != 2 {
		t.Errorf("expected count by seed 100 = 2, got %d", board.CountBySeed(100))
	}
	if board.CountByGenre(engine.GenreFantasy) != 2 {
		t.Errorf("expected count by fantasy = 2, got %d", board.CountByGenre(engine.GenreFantasy))
	}
	if board.GetUniqueSeedCount() != 2 {
		t.Errorf("expected unique seed count 2, got %d", board.GetUniqueSeedCount())
	}
}

func TestBoardGetTopSeed(t *testing.T) {
	board := NewBoard()

	_ = board.Add(NewEntry(100, engine.GenreFantasy, 500, 30, 2))
	_ = board.Add(NewEntry(200, engine.GenreScifi, 1500, 25, 3))
	_ = board.Add(NewEntry(300, engine.GenreHorror, 1000, 35, 4))

	topSeed, topScore := board.GetTopSeed()

	if topSeed != 200 {
		t.Errorf("expected top seed 200, got %d", topSeed)
	}
	if topScore != 1500 {
		t.Errorf("expected top score 1500, got %d", topScore)
	}
}

func TestBoardGetReplayableSeed(t *testing.T) {
	board := NewBoard()

	_ = board.Add(NewEntry(100, engine.GenreFantasy, 500, 30, 2))
	_ = board.Add(NewEntry(200, engine.GenreFantasy, 800, 25, 3))
	_ = board.Add(NewEntry(300, engine.GenreScifi, 1000, 35, 4))

	seed, entry := board.GetReplayableSeed(engine.GenreFantasy)

	if seed != 200 {
		t.Errorf("expected replayable seed 200, got %d", seed)
	}
	if entry.Score != 800 {
		t.Errorf("expected entry score 800, got %d", entry.Score)
	}
}

func TestBoardMarshalUnmarshal(t *testing.T) {
	board := NewBoard()

	_ = board.Add(NewEntry(100, engine.GenreFantasy, 500, 30, 2))
	_ = board.Add(NewEntry(200, engine.GenreScifi, 800, 25, 3))

	data, err := board.Marshal()
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	restored, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if restored.Count() != board.Count() {
		t.Errorf("count mismatch: got %d, want %d", restored.Count(), board.Count())
	}
}

func TestSortByScore(t *testing.T) {
	entries := []*Entry{
		{Score: 100, Days: 30, Survivors: 2},
		{Score: 200, Days: 25, Survivors: 3},
		{Score: 100, Days: 20, Survivors: 4}, // Same score, fewer days
		{Score: 200, Days: 25, Survivors: 5}, // Same score & days, more survivors
	}

	sortByScore(entries)

	// First should be score 200 with 5 survivors (more survivors for tiebreak)
	if entries[0].Score != 200 || entries[0].Survivors != 5 {
		t.Errorf("expected first entry: score=200, survivors=5, got score=%d, survivors=%d",
			entries[0].Score, entries[0].Survivors)
	}

	// Third should be score 100 with 20 days (fewer days for tiebreak)
	if entries[2].Score != 100 || entries[2].Days != 20 {
		t.Errorf("expected third entry: score=100, days=20, got score=%d, days=%d",
			entries[2].Score, entries[2].Days)
	}
}

func TestLocalStorageCreateAndLoad(t *testing.T) {
	tempDir := t.TempDir()

	storage, err := NewLocalStorage(tempDir)
	if err != nil {
		t.Fatalf("NewLocalStorage() error = %v", err)
	}

	entry := NewEntry(12345, engine.GenreFantasy, 1000, 30, 4)
	err = storage.AddEntry(entry)
	if err != nil {
		t.Fatalf("AddEntry() error = %v", err)
	}

	// Save and reload
	err = storage.Save()
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	storage2, err := NewLocalStorage(tempDir)
	if err != nil {
		t.Fatalf("NewLocalStorage() reload error = %v", err)
	}

	board := storage2.GetBoard()
	if board.Count() != 1 {
		t.Errorf("expected 1 entry after reload, got %d", board.Count())
	}
}

func TestLocalStoragePending(t *testing.T) {
	tempDir := t.TempDir()

	storage, err := NewLocalStorage(tempDir)
	if err != nil {
		t.Fatalf("NewLocalStorage() error = %v", err)
	}

	entry1 := NewEntry(100, engine.GenreFantasy, 500, 30, 2)
	entry2 := NewEntry(200, engine.GenreScifi, 800, 25, 3)

	_ = storage.AddEntry(entry1)
	_ = storage.AddEntry(entry2)

	if storage.PendingCount() != 2 {
		t.Errorf("expected 2 pending, got %d", storage.PendingCount())
	}

	pending := storage.GetPendingEntries()
	if len(pending) != 2 {
		t.Errorf("expected 2 pending entries, got %d", len(pending))
	}

	storage.RemovePendingEntry(entry1)
	if storage.PendingCount() != 1 {
		t.Errorf("expected 1 pending after removal, got %d", storage.PendingCount())
	}

	storage.ClearPending()
	if storage.PendingCount() != 0 {
		t.Errorf("expected 0 pending after clear, got %d", storage.PendingCount())
	}
}

func TestLocalStorageMerge(t *testing.T) {
	tempDir := t.TempDir()

	storage, err := NewLocalStorage(tempDir)
	if err != nil {
		t.Fatalf("NewLocalStorage() error = %v", err)
	}

	_ = storage.AddEntry(NewEntry(100, engine.GenreFantasy, 500, 30, 2))

	other := NewBoard()
	_ = other.Add(NewEntry(200, engine.GenreScifi, 800, 25, 3))
	_ = other.Add(NewEntry(300, engine.GenreHorror, 1000, 35, 4))

	storage.MergeBoard(other)

	board := storage.GetBoard()
	if board.Count() != 3 {
		t.Errorf("expected 3 entries after merge, got %d", board.Count())
	}
}

func TestLocalStorageExportImport(t *testing.T) {
	tempDir := t.TempDir()

	storage, err := NewLocalStorage(tempDir)
	if err != nil {
		t.Fatalf("NewLocalStorage() error = %v", err)
	}

	_ = storage.AddEntry(NewEntry(100, engine.GenreFantasy, 500, 30, 2))
	_ = storage.AddEntry(NewEntry(200, engine.GenreScifi, 800, 25, 3))

	// Export
	var buf bytes.Buffer
	err = storage.Export(&buf)
	if err != nil {
		t.Fatalf("Export() error = %v", err)
	}

	// Create new storage and import
	tempDir2 := t.TempDir()
	storage2, err := NewLocalStorage(tempDir2)
	if err != nil {
		t.Fatalf("NewLocalStorage() error = %v", err)
	}

	err = storage2.Import(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("Import() error = %v", err)
	}

	board := storage2.GetBoard()
	if board.Count() != 2 {
		t.Errorf("expected 2 entries after import, got %d", board.Count())
	}
}

func TestLocalStorageInvalidPath(t *testing.T) {
	_, err := NewLocalStorage("")
	if err != ErrStoragePathInvalid {
		t.Errorf("expected ErrStoragePathInvalid, got %v", err)
	}
}

func TestLocalStoragePersistence(t *testing.T) {
	tempDir := t.TempDir()

	// Create storage and add entries
	storage, _ := NewLocalStorage(tempDir)
	_ = storage.AddEntry(NewEntry(100, engine.GenreFantasy, 500, 30, 2))
	_ = storage.Save()

	// Verify files exist
	if _, err := os.Stat(filepath.Join(tempDir, "leaderboard.json")); os.IsNotExist(err) {
		t.Error("leaderboard.json should exist")
	}
	if _, err := os.Stat(filepath.Join(tempDir, "pending_sync.json")); os.IsNotExist(err) {
		t.Error("pending_sync.json should exist")
	}
}

func TestClientConfig(t *testing.T) {
	config := DefaultConfig()

	if config.ServerURL == "" {
		t.Error("ServerURL should not be empty")
	}
	if config.Timeout == 0 {
		t.Error("Timeout should not be zero")
	}
	if config.RetryCount == 0 {
		t.Error("RetryCount should not be zero")
	}
	if config.RetryDelay == 0 {
		t.Error("RetryDelay should not be zero")
	}
}

func TestClientOfflineSubmit(t *testing.T) {
	tempDir := t.TempDir()
	storage, _ := NewLocalStorage(tempDir)

	config := &ClientConfig{
		ServerURL:    "http://invalid.example.com",
		Timeout:      100 * time.Millisecond,
		RetryCount:   0,
		LocalStorage: storage,
	}

	client := NewClient(config)

	// Submission should succeed locally even if server is unavailable
	entry := NewEntry(12345, engine.GenreFantasy, 1000, 30, 4)
	err := client.Submit(entry)
	if err != nil {
		t.Errorf("Submit() should succeed locally, got error = %v", err)
	}

	// Entry should be in local storage
	board := storage.GetBoard()
	if board.Count() != 1 {
		t.Errorf("expected 1 entry in local storage, got %d", board.Count())
	}

	// Entry should be pending sync
	if storage.PendingCount() != 1 {
		t.Errorf("expected 1 pending entry, got %d", storage.PendingCount())
	}
}

func TestClientQueryOptions(t *testing.T) {
	config := DefaultConfig()
	client := NewClient(config)

	seed := int64(12345)
	genre := engine.GenreFantasy
	opts := QueryOptions{
		Seed:   &seed,
		Genre:  &genre,
		Limit:  10,
		Offset: 5,
	}

	url, err := client.buildQueryURL(opts)
	if err != nil {
		t.Fatalf("buildQueryURL() error = %v", err)
	}

	if url == "" {
		t.Error("URL should not be empty")
	}
}

func TestClientLocalQuery(t *testing.T) {
	tempDir := t.TempDir()
	storage, _ := NewLocalStorage(tempDir)

	_ = storage.AddEntry(NewEntry(100, engine.GenreFantasy, 500, 30, 2))
	_ = storage.AddEntry(NewEntry(100, engine.GenreScifi, 800, 25, 3))
	_ = storage.AddEntry(NewEntry(200, engine.GenreFantasy, 1000, 35, 4))

	config := &ClientConfig{
		ServerURL:    "http://invalid.example.com",
		Timeout:      100 * time.Millisecond,
		RetryCount:   0,
		LocalStorage: storage,
	}

	client := NewClient(config)

	// Query should fall back to local storage
	seed := int64(100)
	opts := QueryOptions{Seed: &seed}
	board, err := client.Query(opts)
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}

	if board.Count() != 2 {
		t.Errorf("expected 2 entries for seed 100, got %d", board.Count())
	}
}

func TestClientOnlineStatus(t *testing.T) {
	client := NewClient(nil)

	// Initially should be offline (no server contact)
	if client.IsOnline() {
		t.Error("client should start offline")
	}
}

func TestSubmitResponseParse(t *testing.T) {
	data := []byte(`{"success": true, "rank": 42, "message": "Entry submitted"}`)

	resp, err := ParseSubmitResponse(data)
	if err != nil {
		t.Fatalf("ParseSubmitResponse() error = %v", err)
	}

	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Rank != 42 {
		t.Errorf("expected rank=42, got %d", resp.Rank)
	}
	if resp.Message != "Entry submitted" {
		t.Errorf("expected message='Entry submitted', got %s", resp.Message)
	}
}

func TestEntryJSON(t *testing.T) {
	entry := NewEntry(12345, engine.GenreFantasy, 1000, 30, 4)
	entry.WithPlayer("p1", "Player One")

	data, err := json.Marshal(entry)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var decoded Entry
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if decoded.Seed != entry.Seed {
		t.Errorf("seed mismatch")
	}
	if decoded.PlayerID != entry.PlayerID {
		t.Errorf("playerID mismatch")
	}
}

func TestBoardEmptyOperations(t *testing.T) {
	board := NewBoard()

	// Operations on empty board should not panic
	all := board.GetAll()
	if len(all) != 0 {
		t.Errorf("expected empty list, got %d entries", len(all))
	}

	entries := board.GetBySeed(12345)
	if len(entries) != 0 {
		t.Errorf("expected empty list, got %d entries", len(entries))
	}

	seed, entry := board.GetReplayableSeed(engine.GenreFantasy)
	if seed != 0 || entry != nil {
		t.Error("expected zero seed and nil entry for empty board")
	}

	topSeed, topScore := board.GetTopSeed()
	if topSeed != 0 || topScore != -1 {
		t.Errorf("expected topSeed=0, topScore=-1, got topSeed=%d, topScore=%d",
			topSeed, topScore)
	}
}

func TestBoardTopNBeyondSize(t *testing.T) {
	board := NewBoard()
	_ = board.Add(NewEntry(1, engine.GenreFantasy, 100, 10, 1))
	_ = board.Add(NewEntry(2, engine.GenreFantasy, 200, 10, 1))

	// Request more than available
	top10 := board.GetTopN(10)
	if len(top10) != 2 {
		t.Errorf("expected 2 entries, got %d", len(top10))
	}
}
