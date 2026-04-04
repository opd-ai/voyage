package saveload

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewSaveData(t *testing.T) {
	sd := NewSaveData(1, 12345, engine.GenreFantasy)
	if sd == nil {
		t.Fatal("NewSaveData returned nil")
	}
	if sd.Version != CurrentVersion {
		t.Errorf("expected version %d, got %d", CurrentVersion, sd.Version)
	}
	if sd.Slot != 1 {
		t.Errorf("expected slot 1, got %d", sd.Slot)
	}
	if sd.MasterSeed != 12345 {
		t.Errorf("expected seed 12345, got %d", sd.MasterSeed)
	}
	if sd.Genre != engine.GenreFantasy {
		t.Errorf("expected genre Fantasy, got %v", sd.Genre)
	}
}

func TestSaveDataMarshal(t *testing.T) {
	sd := NewSaveData(1, 12345, engine.GenreFantasy)
	sd.Turn = 10
	sd.Day = 3

	data, err := sd.Marshal()
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	if len(data) == 0 {
		t.Error("Marshal returned empty data")
	}
}

func TestSaveDataUnmarshal(t *testing.T) {
	original := NewSaveData(1, 12345, engine.GenreFantasy)
	original.Turn = 10
	original.Day = 3

	data, err := original.Marshal()
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	loaded, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if loaded.Slot != original.Slot {
		t.Errorf("expected slot %d, got %d", original.Slot, loaded.Slot)
	}
	if loaded.MasterSeed != original.MasterSeed {
		t.Errorf("expected seed %d, got %d", original.MasterSeed, loaded.MasterSeed)
	}
	if loaded.Turn != original.Turn {
		t.Errorf("expected turn %d, got %d", original.Turn, loaded.Turn)
	}
}

func TestSaveDataValidate(t *testing.T) {
	tests := []struct {
		name    string
		modify  func(*SaveData)
		wantErr error
	}{
		{"valid", func(sd *SaveData) {}, nil},
		{"invalid version", func(sd *SaveData) { sd.Version = 0 }, ErrInvalidVersion},
		{"invalid seed", func(sd *SaveData) { sd.MasterSeed = 0 }, ErrInvalidSeed},
		{"invalid slot", func(sd *SaveData) { sd.Slot = -1 }, ErrInvalidSlot},
		{"slot too high", func(sd *SaveData) { sd.Slot = MaxSlots + 1 }, ErrInvalidSlot},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := NewSaveData(1, 12345, engine.GenreFantasy)
			tt.modify(sd)
			err := sd.Validate()
			if err != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestSaveDataGetSummary(t *testing.T) {
	sd := NewSaveData(1, 12345, engine.GenreFantasy)
	sd.Turn = 50
	sd.Day = 12
	sd.Party.Members = []CrewState{
		{ID: 1, Name: "Test", IsAlive: true},
		{ID: 2, Name: "Dead", IsAlive: false},
	}

	summary := sd.GetSummary()
	if summary.Slot != 1 {
		t.Errorf("expected slot 1, got %d", summary.Slot)
	}
	if summary.Turn != 50 {
		t.Errorf("expected turn 50, got %d", summary.Turn)
	}
	if summary.CrewCount != 1 {
		t.Errorf("expected crew count 1, got %d", summary.CrewCount)
	}
}

func TestNewSaveManager(t *testing.T) {
	sm := NewSaveManager()
	if sm == nil {
		t.Fatal("NewSaveManager returned nil")
	}
	if sm.savePath == "" {
		t.Error("SaveManager should have default path")
	}
}

func TestSaveManagerSetPath(t *testing.T) {
	sm := NewSaveManager()
	customPath := "/tmp/test_saves"
	sm.SetSavePath(customPath)
	if sm.SavePath() != customPath {
		t.Errorf("expected path %s, got %s", customPath, sm.SavePath())
	}
}

func TestSaveManagerSaveAndLoad(t *testing.T) {
	// Use temp directory
	tempDir, err := os.MkdirTemp("", "voyage_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sm := NewSaveManager()
	sm.SetSavePath(tempDir)

	// Create save data
	sd := NewSaveData(1, 12345, engine.GenreFantasy)
	sd.Turn = 25
	sd.Day = 5

	// Save
	if err := sm.Save(sd); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	if !sm.SlotExists(1) {
		t.Error("slot should exist after save")
	}

	// Load
	loaded, err := sm.Load(1)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Turn != 25 {
		t.Errorf("expected turn 25, got %d", loaded.Turn)
	}
	if loaded.MasterSeed != 12345 {
		t.Errorf("expected seed 12345, got %d", loaded.MasterSeed)
	}
}

func TestSaveManagerSlotEmpty(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "voyage_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sm := NewSaveManager()
	sm.SetSavePath(tempDir)

	if sm.SlotExists(5) {
		t.Error("slot 5 should not exist")
	}

	_, err = sm.Load(5)
	if err != ErrSlotEmpty {
		t.Errorf("expected ErrSlotEmpty, got %v", err)
	}
}

func TestSaveManagerDelete(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "voyage_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sm := NewSaveManager()
	sm.SetSavePath(tempDir)

	// Save
	sd := NewSaveData(1, 12345, engine.GenreFantasy)
	if err := sm.Save(sd); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Delete
	if err := sm.Delete(1); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if sm.SlotExists(1) {
		t.Error("slot should not exist after delete")
	}
}

func TestSaveManagerAutosave(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "voyage_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sm := NewSaveManager()
	sm.SetSavePath(tempDir)

	sd := NewSaveData(5, 12345, engine.GenreFantasy) // Any slot
	sd.Turn = 100

	// Autosave should use slot 0
	if err := sm.Autosave(sd); err != nil {
		t.Fatalf("Autosave failed: %v", err)
	}

	if !sm.HasAutosave() {
		t.Error("should have autosave")
	}

	loaded, err := sm.LoadAutosave()
	if err != nil {
		t.Fatalf("LoadAutosave failed: %v", err)
	}

	if loaded.Slot != AutosaveSlot {
		t.Errorf("expected slot %d, got %d", AutosaveSlot, loaded.Slot)
	}
	if !loaded.IsAutosave {
		t.Error("should be marked as autosave")
	}
}

func TestSaveManagerListSlots(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "voyage_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sm := NewSaveManager()
	sm.SetSavePath(tempDir)

	// Save to a few slots
	for i := 1; i <= 3; i++ {
		sd := NewSaveData(i, int64(12345+i), engine.GenreFantasy)
		if err := sm.Save(sd); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	slots := sm.ListSlots()
	if len(slots) != MaxSlots+1 {
		t.Errorf("expected %d slots, got %d", MaxSlots+1, len(slots))
	}

	// Verify first 3 non-autosave slots are not empty
	nonEmptyCount := 0
	for _, s := range slots {
		if !s.Empty {
			nonEmptyCount++
		}
	}
	if nonEmptyCount != 3 {
		t.Errorf("expected 3 non-empty slots, got %d", nonEmptyCount)
	}
}

func TestSaveManagerGetMostRecent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "voyage_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sm := NewSaveManager()
	sm.SetSavePath(tempDir)

	// Save to slots with delays
	for i := 1; i <= 3; i++ {
		sd := NewSaveData(i, int64(12345+i), engine.GenreFantasy)
		if err := sm.Save(sd); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	mostRecent, err := sm.GetMostRecent()
	if err != nil {
		t.Fatalf("GetMostRecent failed: %v", err)
	}

	if mostRecent != 3 {
		t.Errorf("expected most recent slot 3, got %d", mostRecent)
	}
}

func TestSaveManagerCopySlot(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "voyage_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sm := NewSaveManager()
	sm.SetSavePath(tempDir)

	// Save to slot 1
	sd := NewSaveData(1, 12345, engine.GenreFantasy)
	sd.Turn = 50
	if err := sm.Save(sd); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Copy to slot 2
	if err := sm.CopySlot(1, 2); err != nil {
		t.Fatalf("CopySlot failed: %v", err)
	}

	// Verify copy
	copied, err := sm.Load(2)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if copied.Slot != 2 {
		t.Errorf("copied save should have slot 2, got %d", copied.Slot)
	}
	if copied.Turn != 50 {
		t.Errorf("expected turn 50, got %d", copied.Turn)
	}
}

func TestSaveManagerExportImport(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "voyage_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	sm := NewSaveManager()
	sm.SetSavePath(tempDir)

	// Save
	sd := NewSaveData(1, 12345, engine.GenreFantasy)
	sd.Turn = 75
	if err := sm.Save(sd); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Export
	exportPath := filepath.Join(tempDir, "export.json")
	if err := sm.ExportSave(1, exportPath); err != nil {
		t.Fatalf("ExportSave failed: %v", err)
	}

	// Import to different slot
	if err := sm.ImportSave(exportPath, 5); err != nil {
		t.Fatalf("ImportSave failed: %v", err)
	}

	// Verify import
	imported, err := sm.Load(5)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if imported.Turn != 75 {
		t.Errorf("expected turn 75, got %d", imported.Turn)
	}
}
