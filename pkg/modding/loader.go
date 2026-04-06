package modding

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Loader handles loading and managing mods.
type Loader struct {
	mu       sync.RWMutex
	mods     map[string]*Mod
	modOrder []string
	basePath string
}

// NewLoader creates a new mod loader.
func NewLoader() *Loader {
	return &Loader{
		mods:     make(map[string]*Mod),
		modOrder: make([]string, 0),
	}
}

// SetBasePath sets the directory to search for mod files.
func (l *Loader) SetBasePath(path string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.basePath = path
}

// LoadFromFile loads a mod from a JSON file.
func (l *Loader) LoadFromFile(path string) (*Mod, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}
	defer file.Close()

	return l.LoadFromReader(file)
}

// LoadFromReader loads a mod from an io.Reader.
func (l *Loader) LoadFromReader(r io.Reader) (*Mod, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return l.LoadFromBytes(data)
}

// LoadFromBytes loads a mod from JSON bytes.
func (l *Loader) LoadFromBytes(data []byte) (*Mod, error) {
	var mod Mod
	if err := json.Unmarshal(data, &mod); err != nil {
		return nil, ErrInvalidJSON
	}

	if err := mod.Validate(); err != nil {
		return nil, err
	}

	return l.Register(&mod)
}

// Register adds a mod to the loader.
func (l *Loader) Register(mod *Mod) (*Mod, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.mods[mod.ID]; exists {
		return nil, ErrModAlreadyLoaded
	}

	mod.LoadedAt = time.Now()
	mod.Enabled = true
	l.mods[mod.ID] = mod
	l.modOrder = append(l.modOrder, mod.ID)

	return mod, nil
}

// Unload removes a mod from the loader.
func (l *Loader) Unload(modID string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.mods[modID]; !exists {
		return ErrModNotFound
	}

	delete(l.mods, modID)

	// Remove from order
	for i, id := range l.modOrder {
		if id == modID {
			l.modOrder = append(l.modOrder[:i], l.modOrder[i+1:]...)
			break
		}
	}

	return nil
}

// Get returns a loaded mod by ID.
func (l *Loader) Get(modID string) (*Mod, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	mod, exists := l.mods[modID]
	if !exists {
		return nil, ErrModNotFound
	}

	return mod, nil
}

// List returns all loaded mods in load order.
func (l *Loader) List() []*Mod {
	l.mu.RLock()
	defer l.mu.RUnlock()

	result := make([]*Mod, 0, len(l.modOrder))
	for _, id := range l.modOrder {
		if mod, exists := l.mods[id]; exists {
			result = append(result, mod)
		}
	}

	return result
}

// ListEnabled returns only enabled mods.
func (l *Loader) ListEnabled() []*Mod {
	l.mu.RLock()
	defer l.mu.RUnlock()

	result := make([]*Mod, 0)
	for _, id := range l.modOrder {
		if mod, exists := l.mods[id]; exists && mod.Enabled {
			result = append(result, mod)
		}
	}

	return result
}

// setModEnabled sets the enabled state of a mod.
func (l *Loader) setModEnabled(modID string, enabled bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	mod, exists := l.mods[modID]
	if !exists {
		return ErrModNotFound
	}

	mod.Enabled = enabled
	return nil
}

// Enable enables a mod.
func (l *Loader) Enable(modID string) error {
	return l.setModEnabled(modID, true)
}

// Disable disables a mod.
func (l *Loader) Disable(modID string) error {
	return l.setModEnabled(modID, false)
}

// LoadDirectory loads all .json files from a directory.
// Validates paths to prevent directory traversal attacks (C-003).
func (l *Loader) LoadDirectory(dirPath string) ([]*Mod, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	// Resolve the base directory path for validation
	absBase, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, err
	}

	var mods []*Mod
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		// Reject entries containing path traversal sequences (C-003)
		if containsPathTraversal(entry.Name()) {
			continue
		}

		modPath := filepath.Join(dirPath, entry.Name())

		// Validate the resolved path stays within the base directory (C-003)
		absPath, err := filepath.Abs(modPath)
		if err != nil {
			continue
		}
		if !isSubpath(absBase, absPath) {
			continue
		}

		mod, err := l.LoadFromFile(modPath)
		if err != nil {
			// Skip invalid mods but continue loading
			continue
		}
		mods = append(mods, mod)
	}

	return mods, nil
}

// containsPathTraversal checks if a filename contains path traversal sequences.
func containsPathTraversal(name string) bool {
	// Check for ".." anywhere in the name
	if len(name) >= 2 {
		for i := 0; i < len(name)-1; i++ {
			if name[i] == '.' && name[i+1] == '.' {
				return true
			}
		}
	}
	// Check for absolute path indicators
	if filepath.IsAbs(name) {
		return true
	}
	return false
}

// isSubpath checks if child is within parent directory.
func isSubpath(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	// If the relative path starts with "..", child is outside parent
	return len(rel) < 2 || (rel[0] != '.' || rel[1] != '.')
}

// GetAllEvents returns all events from all enabled mods.
func (l *Loader) GetAllEvents() []EventDef {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var events []EventDef
	for _, id := range l.modOrder {
		if mod, exists := l.mods[id]; exists && mod.Enabled {
			events = append(events, mod.Events...)
		}
	}

	return events
}

// GetEventsForGenre returns events for a specific genre from all enabled mods.
func (l *Loader) GetEventsForGenre(genre string) []EventDef {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var events []EventDef
	for _, id := range l.modOrder {
		if mod, exists := l.mods[id]; exists && mod.Enabled {
			for _, e := range mod.Events {
				if e.Genre == genre || e.Genre == "" {
					events = append(events, e)
				}
			}
		}
	}

	return events
}

// GetCustomGenres returns all custom genre definitions from enabled mods.
func (l *Loader) GetCustomGenres() []GenreDef {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var genres []GenreDef
	for _, id := range l.modOrder {
		if mod, exists := l.mods[id]; exists && mod.Enabled {
			genres = append(genres, mod.Genres...)
		}
	}

	return genres
}

// Count returns the number of loaded mods.
func (l *Loader) Count() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.mods)
}

// Clear removes all loaded mods.
func (l *Loader) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.mods = make(map[string]*Mod)
	l.modOrder = make([]string, 0)
}
