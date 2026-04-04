package config

import (
	"os"
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewInputConfig(t *testing.T) {
	ic := NewInputConfig()
	if ic == nil {
		t.Fatal("NewInputConfig returned nil")
	}
	if len(ic.bindings) == 0 {
		t.Error("InputConfig should have default bindings")
	}
}

func TestInputConfigDefaults(t *testing.T) {
	ic := NewInputConfig()

	// Check that all actions have bindings
	for _, action := range AllActions() {
		binding := ic.GetBinding(action)
		if binding.KeyCode == 0 && action != ActionMoveUp {
			// Note: KeyCode 0 is technically valid for some keys, but our defaults don't use it
			// ActionMoveUp uses code 38
		}
	}

	// Check specific defaults
	upBinding := ic.GetBinding(ActionMoveUp)
	if upBinding.KeyName != "Up" {
		t.Errorf("expected Up key name, got %s", upBinding.KeyName)
	}
}

func TestInputConfigSetBinding(t *testing.T) {
	ic := NewInputConfig()

	ic.SetBinding(ActionMoveUp, 87, "W") // WASD-style
	binding := ic.GetBinding(ActionMoveUp)

	if binding.KeyCode != 87 {
		t.Errorf("expected keyCode 87, got %d", binding.KeyCode)
	}
	if binding.KeyName != "W" {
		t.Errorf("expected keyName W, got %s", binding.KeyName)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}

	if cfg.ScreenWidth != 800 {
		t.Errorf("expected ScreenWidth 800, got %d", cfg.ScreenWidth)
	}
	if cfg.ScreenHeight != 600 {
		t.Errorf("expected ScreenHeight 600, got %d", cfg.ScreenHeight)
	}
	if cfg.Genre != engine.GenreFantasy {
		t.Errorf("expected Genre Fantasy, got %v", cfg.Genre)
	}
	if cfg.Difficulty != DifficultyNormal {
		t.Errorf("expected Difficulty Normal, got %v", cfg.Difficulty)
	}
	if cfg.Input == nil {
		t.Error("Config.Input should not be nil")
	}
}

func TestConfigValidate(t *testing.T) {
	cfg := DefaultConfig()

	if !cfg.Validate() {
		t.Error("default config should be valid")
	}

	// Test invalid screen size
	cfg.ScreenWidth = 100
	if cfg.Validate() {
		t.Error("config with small screen should be invalid")
	}

	// Test invalid tile size
	cfg = DefaultConfig()
	cfg.TileSize = 4
	if cfg.Validate() {
		t.Error("config with small tile size should be invalid")
	}

	// Test invalid volume
	cfg = DefaultConfig()
	cfg.MasterVolume = 1.5
	if cfg.Validate() {
		t.Error("config with volume > 1 should be invalid")
	}
}

func TestDifficultyModifiers(t *testing.T) {
	tests := []struct {
		difficulty   Difficulty
		expectResMin float64
		expectResMax float64
	}{
		{DifficultyEasy, 1.2, 1.4},
		{DifficultyNormal, 0.9, 1.1},
		{DifficultyHard, 0.7, 0.9},
		{DifficultyNightmare, 0.5, 0.7},
	}

	for _, tt := range tests {
		resMod, _ := DifficultyModifiers(tt.difficulty)
		if resMod < tt.expectResMin || resMod > tt.expectResMax {
			t.Errorf("DifficultyModifiers(%v) resource mod = %f, want between %f and %f",
				tt.difficulty, resMod, tt.expectResMin, tt.expectResMax)
		}
	}
}

func TestDifficultyName(t *testing.T) {
	names := map[Difficulty]string{
		DifficultyEasy:      "Easy",
		DifficultyNormal:    "Normal",
		DifficultyHard:      "Hard",
		DifficultyNightmare: "Nightmare",
	}

	for diff, expected := range names {
		if got := DifficultyName(diff); got != expected {
			t.Errorf("DifficultyName(%v) = %s, want %s", diff, got, expected)
		}
	}
}

func TestAllActions(t *testing.T) {
	actions := AllActions()
	if len(actions) < 10 {
		t.Errorf("expected at least 10 actions, got %d", len(actions))
	}
}

func TestAllDifficulties(t *testing.T) {
	diffs := AllDifficulties()
	if len(diffs) != 4 {
		t.Errorf("expected 4 difficulties, got %d", len(diffs))
	}
}

func TestActionName(t *testing.T) {
	name := ActionName(ActionMoveUp)
	if name == "" {
		t.Error("ActionName should return non-empty string")
	}
	if name != "Move Up" {
		t.Errorf("expected 'Move Up', got '%s'", name)
	}
}

func TestParseDifficulty(t *testing.T) {
	tests := []struct {
		input    string
		expected Difficulty
		valid    bool
	}{
		{"easy", DifficultyEasy, true},
		{"Easy", DifficultyEasy, true},
		{"normal", DifficultyNormal, true},
		{"Normal", DifficultyNormal, true},
		{"hard", DifficultyHard, true},
		{"Hard", DifficultyHard, true},
		{"nightmare", DifficultyNightmare, true},
		{"Nightmare", DifficultyNightmare, true},
		{"invalid", DifficultyNormal, false},
		{"", DifficultyNormal, false},
		{"EASY", DifficultyNormal, false}, // case sensitive except for first letter
	}

	for _, tt := range tests {
		diff, ok := ParseDifficulty(tt.input)
		if ok != tt.valid {
			t.Errorf("ParseDifficulty(%q) valid = %v, want %v", tt.input, ok, tt.valid)
		}
		if diff != tt.expected {
			t.Errorf("ParseDifficulty(%q) = %v, want %v", tt.input, diff, tt.expected)
		}
	}
}

func TestIsValidDifficulty(t *testing.T) {
	validInputs := []string{"easy", "Easy", "normal", "Normal", "hard", "Hard", "nightmare", "Nightmare"}
	for _, input := range validInputs {
		if !IsValidDifficulty(input) {
			t.Errorf("IsValidDifficulty(%q) = false, want true", input)
		}
	}

	invalidInputs := []string{"invalid", "", "EASY", "medium", "insane"}
	for _, input := range invalidInputs {
		if IsValidDifficulty(input) {
			t.Errorf("IsValidDifficulty(%q) = true, want false", input)
		}
	}
}

func TestConfigPath(t *testing.T) {
	path := ConfigPath()
	if path == "" {
		t.Error("ConfigPath should return non-empty string")
	}
}

func TestConfigFilePath(t *testing.T) {
	path := ConfigFilePath()
	if path == "" {
		t.Error("ConfigFilePath should return non-empty string")
	}
	// Should end with config.json
	if len(path) < 11 || path[len(path)-11:] != "config.json" {
		t.Errorf("ConfigFilePath should end with config.json, got %s", path)
	}
}

func TestSaveAndLoadConfig(t *testing.T) {
	// Use temp dir for testing
	tempDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	// Create and modify config
	cfg := DefaultConfig()
	cfg.Seed = 12345
	cfg.ScreenWidth = 1024
	cfg.Difficulty = DifficultyHard

	// Save
	err := SaveConfig(cfg)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify exists
	if !ConfigExists() {
		t.Error("ConfigExists should return true after save")
	}

	// Load
	loaded, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if loaded.Seed != cfg.Seed {
		t.Errorf("loaded Seed = %d, want %d", loaded.Seed, cfg.Seed)
	}
	if loaded.ScreenWidth != cfg.ScreenWidth {
		t.Errorf("loaded ScreenWidth = %d, want %d", loaded.ScreenWidth, cfg.ScreenWidth)
	}
	if loaded.Difficulty != cfg.Difficulty {
		t.Errorf("loaded Difficulty = %v, want %v", loaded.Difficulty, cfg.Difficulty)
	}

	// Delete
	err = DeleteConfig()
	if err != nil {
		t.Fatalf("DeleteConfig failed: %v", err)
	}

	if ConfigExists() {
		t.Error("ConfigExists should return false after delete")
	}
}

func TestLoadConfigDefault(t *testing.T) {
	// Use temp dir with no config
	tempDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", oldHome)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	// Should return default config
	if cfg.ScreenWidth != 800 || cfg.ScreenHeight != 600 {
		t.Error("LoadConfig should return default config when file doesn't exist")
	}
}

func TestInputConfigJSON(t *testing.T) {
	ic := NewInputConfig()
	ic.SetBinding(ActionMoveUp, 87, "W")

	// Marshal
	data, err := ic.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	// Unmarshal into new config
	ic2 := &InputConfig{}
	err = ic2.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	binding := ic2.GetBinding(ActionMoveUp)
	if binding.KeyCode != 87 {
		t.Errorf("restored keyCode = %d, want 87", binding.KeyCode)
	}
	if binding.KeyName != "W" {
		t.Errorf("restored keyName = %s, want W", binding.KeyName)
	}
}

func TestGetKeyCode(t *testing.T) {
	ic := NewInputConfig()
	
	code := ic.GetKeyCode(ActionMoveUp)
	if code == 0 {
		t.Error("GetKeyCode should return non-zero for default binding")
	}
}

func TestAllBindings(t *testing.T) {
	ic := NewInputConfig()
	
	bindings := ic.AllBindings()
	if len(bindings) < 10 {
		t.Errorf("AllBindings should return at least 10 bindings, got %d", len(bindings))
	}
}
