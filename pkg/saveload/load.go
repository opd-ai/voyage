package saveload

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// MaxSlots is the maximum number of save slots.
const MaxSlots = 10

// AutosaveSlot is the dedicated autosave slot.
const AutosaveSlot = 0

// Errors
var (
	ErrInvalidVersion = errors.New("invalid save file version")
	ErrInvalidSeed    = errors.New("invalid master seed")
	ErrInvalidSlot    = errors.New("invalid save slot")
	ErrSlotEmpty      = errors.New("save slot is empty")
	ErrSaveFailed     = errors.New("failed to write save file")
	ErrLoadFailed     = errors.New("failed to read save file")
)

// SaveManager handles save and load operations.
type SaveManager struct {
	savePath string
}

// NewSaveManager creates a new save manager.
func NewSaveManager() *SaveManager {
	return &SaveManager{
		savePath: defaultSavePath(),
	}
}

// defaultSavePath returns the default save directory.
func defaultSavePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".voyage_saves"
	}
	return filepath.Join(homeDir, ".voyage", "saves")
}

// SetSavePath sets a custom save directory.
func (sm *SaveManager) SetSavePath(path string) {
	sm.savePath = path
}

// SavePath returns the current save directory.
func (sm *SaveManager) SavePath() string {
	return sm.savePath
}

// ensureSaveDir creates the save directory if it doesn't exist.
func (sm *SaveManager) ensureSaveDir() error {
	return os.MkdirAll(sm.savePath, 0o755)
}

// slotFilename returns the filename for a save slot.
func (sm *SaveManager) slotFilename(slot int) string {
	if slot == AutosaveSlot {
		return filepath.Join(sm.savePath, "autosave.json")
	}
	return filepath.Join(sm.savePath, fmt.Sprintf("save_%d.json", slot))
}

// Save persists a SaveData to the specified slot.
func (sm *SaveManager) Save(data *SaveData) error {
	if err := data.Validate(); err != nil {
		return err
	}

	if err := sm.ensureSaveDir(); err != nil {
		return fmt.Errorf("%w: %v", ErrSaveFailed, err)
	}

	data.SavedAt = time.Now()
	data.IsAutosave = (data.Slot == AutosaveSlot)

	jsonData, err := data.Marshal()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSaveFailed, err)
	}

	filename := sm.slotFilename(data.Slot)
	if err := os.WriteFile(filename, jsonData, 0o644); err != nil {
		return fmt.Errorf("%w: %v", ErrSaveFailed, err)
	}

	return nil
}

// Load reads a SaveData from the specified slot.
func (sm *SaveManager) Load(slot int) (*SaveData, error) {
	filename := sm.slotFilename(slot)

	jsonData, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrSlotEmpty
		}
		return nil, fmt.Errorf("%w: %v", ErrLoadFailed, err)
	}

	data, err := Unmarshal(jsonData)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLoadFailed, err)
	}

	if err := data.Validate(); err != nil {
		return nil, err
	}

	return data, nil
}

// SlotExists checks if a save slot has a file.
func (sm *SaveManager) SlotExists(slot int) bool {
	filename := sm.slotFilename(slot)
	_, err := os.Stat(filename)
	return err == nil
}

// Delete removes a save file from a slot.
func (sm *SaveManager) Delete(slot int) error {
	filename := sm.slotFilename(slot)
	if err := os.Remove(filename); err != nil {
		if os.IsNotExist(err) {
			return ErrSlotEmpty
		}
		return err
	}
	return nil
}

// ListSlots returns information about all save slots.
func (sm *SaveManager) ListSlots() []SlotInfo {
	slots := make([]SlotInfo, 0, MaxSlots+1)

	for slot := 0; slot <= MaxSlots; slot++ {
		info := SlotInfo{
			Slot:   slot,
			Empty:  !sm.SlotExists(slot),
			IsAuto: slot == AutosaveSlot,
		}

		if !info.Empty {
			data, err := sm.Load(slot)
			if err == nil {
				info.Summary = data.GetSummary()
			}
		}

		slots = append(slots, info)
	}

	return slots
}

// SlotInfo provides metadata about a save slot.
type SlotInfo struct {
	Slot    int
	Empty   bool
	IsAuto  bool
	Summary SaveSummary
}

// GetMostRecent returns the most recently saved slot.
func (sm *SaveManager) GetMostRecent() (int, error) {
	slots := sm.ListSlots()

	var mostRecent *SlotInfo
	for i := range slots {
		if slots[i].Empty {
			continue
		}
		if mostRecent == nil || slots[i].Summary.SavedAt.After(mostRecent.Summary.SavedAt) {
			mostRecent = &slots[i]
		}
	}

	if mostRecent == nil {
		return -1, ErrSlotEmpty
	}

	return mostRecent.Slot, nil
}

// Autosave saves to the autosave slot.
func (sm *SaveManager) Autosave(data *SaveData) error {
	data.Slot = AutosaveSlot
	return sm.Save(data)
}

// LoadAutosave loads from the autosave slot.
func (sm *SaveManager) LoadAutosave() (*SaveData, error) {
	return sm.Load(AutosaveSlot)
}

// HasAutosave checks if an autosave exists.
func (sm *SaveManager) HasAutosave() bool {
	return sm.SlotExists(AutosaveSlot)
}

// CopySlot copies a save from one slot to another.
func (sm *SaveManager) CopySlot(fromSlot, toSlot int) error {
	data, err := sm.Load(fromSlot)
	if err != nil {
		return err
	}

	data.Slot = toSlot
	data.IsAutosave = (toSlot == AutosaveSlot)
	return sm.Save(data)
}

// GetOldestSlot returns the oldest non-autosave slot (for overwriting).
func (sm *SaveManager) GetOldestSlot() int {
	slots := sm.ListSlots()

	// First, look for empty slots (skip autosave)
	for i := 1; i <= MaxSlots; i++ {
		if slots[i].Empty {
			return i
		}
	}

	// All slots full, find oldest
	var oldest *SlotInfo
	for i := 1; i <= MaxSlots; i++ {
		if oldest == nil || slots[i].Summary.SavedAt.Before(oldest.Summary.SavedAt) {
			oldest = &slots[i]
		}
	}

	if oldest != nil {
		return oldest.Slot
	}
	return 1
}

// CleanupOldSaves removes saves older than the specified duration.
func (sm *SaveManager) CleanupOldSaves(maxAge time.Duration) (int, error) {
	slots := sm.ListSlots()
	cutoff := time.Now().Add(-maxAge)
	deleted := 0

	for _, slot := range slots {
		if slot.Empty || slot.IsAuto {
			continue
		}
		if slot.Summary.SavedAt.Before(cutoff) {
			if err := sm.Delete(slot.Slot); err == nil {
				deleted++
			}
		}
	}

	return deleted, nil
}

// ExportSave exports a save to a specified path.
func (sm *SaveManager) ExportSave(slot int, exportPath string) error {
	data, err := sm.Load(slot)
	if err != nil {
		return err
	}

	jsonData, err := data.Marshal()
	if err != nil {
		return err
	}

	return os.WriteFile(exportPath, jsonData, 0o644)
}

// ImportSave imports a save from a specified path.
func (sm *SaveManager) ImportSave(importPath string, slot int) error {
	jsonData, err := os.ReadFile(importPath)
	if err != nil {
		return err
	}

	data, err := Unmarshal(jsonData)
	if err != nil {
		return err
	}

	data.Slot = slot
	return sm.Save(data)
}

// BackupAll creates backup copies of all saves.
func (sm *SaveManager) BackupAll(backupDir string) error {
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return err
	}

	slots := sm.ListSlots()
	for _, slot := range slots {
		if slot.Empty {
			continue
		}
		backupFile := filepath.Join(backupDir, fmt.Sprintf("backup_%d_%d.json", slot.Slot, time.Now().Unix()))
		if err := sm.ExportSave(slot.Slot, backupFile); err != nil {
			return err
		}
	}

	return nil
}

// SortSlotsByDate returns slot indices sorted by save date (newest first).
func (sm *SaveManager) SortSlotsByDate() []int {
	slots := sm.ListSlots()

	// Filter non-empty slots
	nonEmpty := make([]SlotInfo, 0)
	for _, s := range slots {
		if !s.Empty {
			nonEmpty = append(nonEmpty, s)
		}
	}

	// Sort by date descending
	sort.Slice(nonEmpty, func(i, j int) bool {
		return nonEmpty[i].Summary.SavedAt.After(nonEmpty[j].Summary.SavedAt)
	})

	// Extract slot numbers
	result := make([]int, len(nonEmpty))
	for i, s := range nonEmpty {
		result[i] = s.Slot
	}

	return result
}
