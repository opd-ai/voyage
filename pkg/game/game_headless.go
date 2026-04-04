//go:build headless

package game

// GameState represents the current state of the game.
// This is a stub for headless builds.
type GameState int

const (
	// StateMenu is the main menu state.
	StateMenu GameState = iota
	// StatePlaying is the active gameplay state.
	StatePlaying
	// StatePaused is the paused state.
	StatePaused
	// StateGameOver is the game over state.
	StateGameOver
)

// Config holds game configuration options.
// This is a stub for headless builds.
type Config struct {
	Width    int
	Height   int
	TileSize int
	Seed     int64
}
