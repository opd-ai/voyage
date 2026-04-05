package config

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// Action represents a bindable game action.
type Action int

const (
	// ActionMoveUp moves the vessel up.
	ActionMoveUp Action = iota
	// ActionMoveDown moves the vessel down.
	ActionMoveDown
	// ActionMoveLeft moves the vessel left.
	ActionMoveLeft
	// ActionMoveRight moves the vessel right.
	ActionMoveRight
	// ActionConfirm confirms a selection.
	ActionConfirm
	// ActionCancel cancels or opens menu.
	ActionCancel
	// ActionChoice1 selects choice 1 in events.
	ActionChoice1
	// ActionChoice2 selects choice 2 in events.
	ActionChoice2
	// ActionChoice3 selects choice 3 in events.
	ActionChoice3
	// ActionChoice4 selects choice 4 in events.
	ActionChoice4
	// ActionDebug toggles debug mode.
	ActionDebug
	// ActionRest triggers rest action.
	ActionRest
)

// AllActions returns all available actions.
func AllActions() []Action {
	return []Action{
		ActionMoveUp,
		ActionMoveDown,
		ActionMoveLeft,
		ActionMoveRight,
		ActionConfirm,
		ActionCancel,
		ActionChoice1,
		ActionChoice2,
		ActionChoice3,
		ActionChoice4,
		ActionDebug,
		ActionRest,
	}
}

// ActionName returns the human-readable name for an action.
func ActionName(a Action) string {
	names := map[Action]string{
		ActionMoveUp:    "Move Up",
		ActionMoveDown:  "Move Down",
		ActionMoveLeft:  "Move Left",
		ActionMoveRight: "Move Right",
		ActionConfirm:   "Confirm",
		ActionCancel:    "Cancel/Menu",
		ActionChoice1:   "Choice 1",
		ActionChoice2:   "Choice 2",
		ActionChoice3:   "Choice 3",
		ActionChoice4:   "Choice 4",
		ActionDebug:     "Toggle Debug",
		ActionRest:      "Rest",
	}
	return names[a]
}

// KeyBinding maps an action to a key code.
type KeyBinding struct {
	Action  Action
	KeyCode int
	KeyName string
}

// InputConfig holds the input mapping configuration.
type InputConfig struct {
	bindings map[Action]KeyBinding
}

// NewInputConfig creates a new input config with default bindings.
func NewInputConfig() *InputConfig {
	ic := &InputConfig{
		bindings: make(map[Action]KeyBinding),
	}
	ic.SetDefaults()
	return ic
}

// SetDefaults sets the default key bindings.
func (ic *InputConfig) SetDefaults() {
	// Default bindings use Ebitengine key codes
	// Up = 38, Down = 40, Left = 37, Right = 39
	// Enter = 13, Escape = 27
	// 1 = 49, 2 = 50, 3 = 51, 4 = 52
	// F3 = 114, R = 82
	defaults := []KeyBinding{
		{ActionMoveUp, 38, "Up"},
		{ActionMoveDown, 40, "Down"},
		{ActionMoveLeft, 37, "Left"},
		{ActionMoveRight, 39, "Right"},
		{ActionConfirm, 13, "Enter"},
		{ActionCancel, 27, "Escape"},
		{ActionChoice1, 49, "1"},
		{ActionChoice2, 50, "2"},
		{ActionChoice3, 51, "3"},
		{ActionChoice4, 52, "4"},
		{ActionDebug, 114, "F3"},
		{ActionRest, 82, "R"},
	}

	for _, b := range defaults {
		ic.bindings[b.Action] = b
	}
}

// GetBinding returns the key binding for an action.
func (ic *InputConfig) GetBinding(action Action) KeyBinding {
	return ic.bindings[action]
}

// SetBinding sets the key binding for an action.
func (ic *InputConfig) SetBinding(action Action, keyCode int, keyName string) {
	ic.bindings[action] = KeyBinding{
		Action:  action,
		KeyCode: keyCode,
		KeyName: keyName,
	}
}

// GetKeyCode returns the key code for an action.
func (ic *InputConfig) GetKeyCode(action Action) int {
	return ic.bindings[action].KeyCode
}

// AllBindings returns all current key bindings.
func (ic *InputConfig) AllBindings() []KeyBinding {
	result := make([]KeyBinding, 0, len(ic.bindings))
	for _, b := range ic.bindings {
		result = append(result, b)
	}
	return result
}

// Config holds all game configuration.
type Config struct {
	// Game settings
	Seed       int64
	Genre      engine.GenreID
	Difficulty Difficulty

	// Display settings
	ScreenWidth  int
	ScreenHeight int
	TileSize     int
	Fullscreen   bool

	// Audio settings
	MasterVolume float64
	MusicVolume  float64
	SFXVolume    float64

	// Input settings
	Input *InputConfig
}

// Difficulty levels for the game.
type Difficulty int

const (
	// DifficultyEasy provides more resources and easier events.
	DifficultyEasy Difficulty = iota
	// DifficultyNormal is the standard experience.
	DifficultyNormal
	// DifficultyHard reduces resources and increases hazards.
	DifficultyHard
	// DifficultyNightmare is extreme difficulty.
	DifficultyNightmare
)

// DifficultyName returns the name of a difficulty level.
func DifficultyName(d Difficulty) string {
	names := map[Difficulty]string{
		DifficultyEasy:      "Easy",
		DifficultyNormal:    "Normal",
		DifficultyHard:      "Hard",
		DifficultyNightmare: "Nightmare",
	}
	return names[d]
}

// AllDifficulties returns all difficulty levels.
func AllDifficulties() []Difficulty {
	return []Difficulty{
		DifficultyEasy,
		DifficultyNormal,
		DifficultyHard,
		DifficultyNightmare,
	}
}

// ParseDifficulty converts a string to a Difficulty level.
// Returns DifficultyNormal and false if the string is invalid.
func ParseDifficulty(s string) (Difficulty, bool) {
	switch s {
	case "easy", "Easy":
		return DifficultyEasy, true
	case "normal", "Normal":
		return DifficultyNormal, true
	case "hard", "Hard":
		return DifficultyHard, true
	case "nightmare", "Nightmare":
		return DifficultyNightmare, true
	default:
		return DifficultyNormal, false
	}
}

// IsValidDifficulty checks if a string is a valid difficulty level.
func IsValidDifficulty(s string) bool {
	_, ok := ParseDifficulty(s)
	return ok
}

// DefaultConfig returns the default game configuration.
func DefaultConfig() *Config {
	return &Config{
		Seed:         0,
		Genre:        engine.GenreFantasy,
		Difficulty:   DifficultyNormal,
		ScreenWidth:  800,
		ScreenHeight: 600,
		TileSize:     16,
		Fullscreen:   false,
		MasterVolume: 1.0,
		MusicVolume:  0.7,
		SFXVolume:    0.8,
		Input:        NewInputConfig(),
	}
}

// Validate checks the configuration for valid values.
func (c *Config) Validate() bool {
	if c.ScreenWidth < 320 || c.ScreenHeight < 240 {
		return false
	}
	if c.TileSize < 8 || c.TileSize > 64 {
		return false
	}
	if c.MasterVolume < 0 || c.MasterVolume > 1 {
		return false
	}
	return true
}

// DifficultyModifiers returns resource and event modifiers for a difficulty.
func DifficultyModifiers(d Difficulty) (resourceMod, eventMod float64) {
	switch d {
	case DifficultyEasy:
		return 1.3, 0.7
	case DifficultyNormal:
		return 1.0, 1.0
	case DifficultyHard:
		return 0.8, 1.2
	case DifficultyNightmare:
		return 0.6, 1.5
	default:
		return 1.0, 1.0
	}
}
