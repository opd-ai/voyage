package leaderboard

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// Storage errors.
var (
	ErrStorageNotInitialized = errors.New("storage not initialized")
	ErrStoragePathInvalid    = errors.New("storage path is invalid")
)

// LocalStorage provides persistent local storage for leaderboard data.
type LocalStorage struct {
	mu       sync.RWMutex
	basePath string
	board    *Board
	pending  []*Entry // Entries waiting to be synced to server
}

// NewLocalStorage creates a new local storage instance.
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if basePath == "" {
		return nil, ErrStoragePathInvalid
	}

	// Ensure directory exists
	if err := os.MkdirAll(basePath, 0o755); err != nil {
		return nil, err
	}

	ls := &LocalStorage{
		basePath: basePath,
		board:    NewBoard(),
		pending:  make([]*Entry, 0),
	}

	// Try to load existing data
	if err := ls.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return ls, nil
}

// leaderboardFile returns the path to the leaderboard file.
func (ls *LocalStorage) leaderboardFile() string {
	return filepath.Join(ls.basePath, "leaderboard.json")
}

// pendingFile returns the path to the pending sync file.
func (ls *LocalStorage) pendingFile() string {
	return filepath.Join(ls.basePath, "pending_sync.json")
}

// load reads stored leaderboard data from disk.
func (ls *LocalStorage) load() error {
	// Load main leaderboard
	data, err := os.ReadFile(ls.leaderboardFile())
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	board, err := Unmarshal(data)
	if err != nil {
		return err
	}
	ls.board = board

	// Load pending entries
	pendingData, err := os.ReadFile(ls.pendingFile())
	if err == nil {
		var pending []*Entry
		if err := json.Unmarshal(pendingData, &pending); err == nil {
			ls.pending = pending
		}
	}

	return nil
}

// Save persists leaderboard data to disk.
func (ls *LocalStorage) Save() error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	// Save main leaderboard
	data, err := ls.board.Marshal()
	if err != nil {
		return err
	}
	if err := ls.atomicWrite(ls.leaderboardFile(), data); err != nil {
		return err
	}

	// Save pending entries
	pendingData, err := json.Marshal(ls.pending)
	if err != nil {
		return err
	}
	return ls.atomicWrite(ls.pendingFile(), pendingData)
}

// atomicWrite writes data atomically using a temp file.
func (ls *LocalStorage) atomicWrite(path string, data []byte) error {
	tempFile := path + ".tmp"
	if err := os.WriteFile(tempFile, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tempFile, path)
}

// AddEntry adds an entry to local storage and marks it pending.
func (ls *LocalStorage) AddEntry(entry *Entry) error {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	if err := ls.board.Add(entry); err != nil {
		return err
	}

	ls.pending = append(ls.pending, entry)
	return nil
}

// GetBoard returns the local leaderboard.
func (ls *LocalStorage) GetBoard() *Board {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	return ls.board
}

// GetPendingEntries returns entries waiting to be synced.
func (ls *LocalStorage) GetPendingEntries() []*Entry {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	result := make([]*Entry, len(ls.pending))
	copy(result, ls.pending)
	return result
}

// ClearPending clears the pending sync queue.
func (ls *LocalStorage) ClearPending() {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.pending = make([]*Entry, 0)
}

// RemovePendingEntry removes a specific entry from the pending queue.
func (ls *LocalStorage) RemovePendingEntry(entry *Entry) {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	for i, e := range ls.pending {
		if e.Seed == entry.Seed && e.Timestamp.Equal(entry.Timestamp) {
			ls.pending = append(ls.pending[:i], ls.pending[i+1:]...)
			return
		}
	}
}

// PendingCount returns the number of entries waiting to be synced.
func (ls *LocalStorage) PendingCount() int {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	return len(ls.pending)
}

// MergeBoard merges entries from another board into local storage.
func (ls *LocalStorage) MergeBoard(other *Board) {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	for _, entry := range other.GetAll() {
		// Ignore errors for duplicates
		_ = ls.board.Add(entry)
	}
}

// Export writes the leaderboard to a writer.
func (ls *LocalStorage) Export(w io.Writer) error {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	data, err := ls.board.Marshal()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// Import reads leaderboard entries from a reader and merges them.
func (ls *LocalStorage) Import(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	board, err := Unmarshal(data)
	if err != nil {
		return err
	}

	ls.MergeBoard(board)
	return nil
}
