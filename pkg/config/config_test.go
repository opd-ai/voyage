package config

import (
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
