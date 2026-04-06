//go:build !headless

package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
	"github.com/opd-ai/voyage/pkg/rendering"
)

// GameState represents the current state of the game.
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

// Game implements ebiten.Game and manages the game lifecycle.
type Game struct {
	world     *engine.World
	renderer  *rendering.Renderer
	state     GameState
	seed      int64
	genre     engine.GenreID
	tileGen   *rendering.TileGenerator
	turn      int
	debugMode bool

	// Screen dimensions
	width  int
	height int

	// Input state tracking for key release detection
	f3WasPressed bool

	// Cached overlay image to avoid per-frame allocations (C-004)
	pauseOverlay *ebiten.Image
}

// Config holds game configuration options.
type Config struct {
	Width    int
	Height   int
	TileSize int
	Seed     int64
	Genre    engine.GenreID
}

// DefaultConfig returns default game configuration.
func DefaultConfig() Config {
	return Config{
		Width:    800,
		Height:   600,
		TileSize: 16,
		Seed:     0,
		Genre:    engine.GenreFantasy,
	}
}

// NewGame creates a new game instance with the given configuration.
func NewGame(cfg Config) *Game {
	registry := engine.NewComponentRegistry()
	world := engine.NewWorld(registry)
	world.SetGenre(cfg.Genre)

	renderer := rendering.NewRenderer(cfg.Width, cfg.Height, cfg.TileSize)
	renderer.SetGenre(cfg.Genre)

	tileGen := rendering.NewTileGenerator(cfg.Seed, cfg.TileSize)

	return &Game{
		world:     world,
		renderer:  renderer,
		state:     StateMenu,
		seed:      cfg.Seed,
		genre:     cfg.Genre,
		tileGen:   tileGen,
		turn:      0,
		debugMode: false,
		width:     cfg.Width,
		height:    cfg.Height,
	}
}

// Update implements ebiten.Game.Update.
func (g *Game) Update() error {
	// Check for window close and perform graceful shutdown (M-010)
	if ebiten.IsWindowBeingClosed() {
		// Trigger autosave before exit if game is in progress
		if g.state == StatePlaying || g.state == StatePaused {
			// Note: actual autosave would be called here if saveload integration exists
		}
		return ebiten.Termination
	}

	g.handleStateInput()
	g.handleDebugToggle()
	return nil
}

// handleStateInput processes input based on current game state.
func (g *Game) handleStateInput() {
	switch g.state {
	case StateMenu:
		g.handleMenuInput()
	case StatePlaying:
		g.handlePlayingInput()
	case StatePaused:
		g.handlePausedInput()
	case StateGameOver:
		g.handleGameOverInput()
	}
}

// handleMenuInput handles input in menu state.
func (g *Game) handleMenuInput() {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.state = StatePlaying
	}
}

// handlePlayingInput handles input during gameplay.
func (g *Game) handlePlayingInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.state = StatePaused
	}
	// Use dynamic delta time based on actual TPS (H-010)
	dt := 1.0 / float64(ebiten.TPS())
	g.world.Update(dt)
}

// handlePausedInput handles input in paused state.
func (g *Game) handlePausedInput() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.state = StatePlaying
	}
}

// handleGameOverInput handles input in game over state.
func (g *Game) handleGameOverInput() {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.state = StateMenu
	}
}

// handleDebugToggle toggles debug mode with F3 key.
// Uses key release detection to prevent multiple toggles per frame.
func (g *Game) handleDebugToggle() {
	if ebiten.IsKeyPressed(ebiten.KeyF3) {
		if !g.f3WasPressed {
			g.debugMode = !g.debugMode
		}
		g.f3WasPressed = true
	} else {
		g.f3WasPressed = false
	}
}

// Draw implements ebiten.Game.Draw.
func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen with background color
	screen.Fill(g.renderer.Palette().Background)

	switch g.state {
	case StateMenu:
		g.drawMenu(screen)
	case StatePlaying:
		g.drawGame(screen)
	case StatePaused:
		g.drawGame(screen)
		g.drawPauseOverlay(screen)
	case StateGameOver:
		g.drawGameOver(screen)
	}

	if g.debugMode {
		g.drawDebug(screen)
	}
}

// Layout implements ebiten.Game.Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}

// World returns the game's ECS world.
func (g *Game) World() *engine.World {
	return g.world
}

// Renderer returns the game's renderer.
func (g *Game) Renderer() *rendering.Renderer {
	return g.renderer
}

// State returns the current game state.
func (g *Game) State() GameState {
	return g.state
}

// SetState changes the game state.
func (g *Game) SetState(state GameState) {
	g.state = state
}

// Turn returns the current turn number.
func (g *Game) Turn() int {
	return g.turn
}

// AdvanceTurn increments the turn counter.
func (g *Game) AdvanceTurn() {
	g.turn++
}

// SetGenre changes the game's genre theme.
func (g *Game) SetGenre(genre engine.GenreID) {
	g.genre = genre
	g.world.SetGenre(genre)
	g.renderer.SetGenre(genre)
}

// Seed returns the master seed.
func (g *Game) Seed() int64 {
	return g.seed
}

// drawMenu renders the main menu.
func (g *Game) drawMenu(screen *ebiten.Image) {
	msg := fmt.Sprintf("VOYAGE\n\nGenre: %s\nSeed: %d\n\nPress ENTER or SPACE to start", g.genre, g.seed)
	ebitenutil.DebugPrintAt(screen, msg, g.width/4, g.height/3)
}

// drawGame renders the main gameplay view.
func (g *Game) drawGame(screen *ebiten.Image) {
	// Draw a simple demo grid
	tileSize := g.renderer.TileSize()
	palette := g.renderer.Palette()

	for y := 0; y < g.height/tileSize; y++ {
		for x := 0; x < g.width/tileSize; x++ {
			// Use seed-based generator for consistent tiles
			gen := seed.NewGenerator(g.seed, fmt.Sprintf("tile_%d_%d", x, y))
			tileType := gen.Intn(len(palette.TileColors))
			g.renderer.DrawTile(screen, x, y, tileType)
		}
	}
}

// drawPauseOverlay renders the pause screen overlay.
func (g *Game) drawPauseOverlay(screen *ebiten.Image) {
	// Use cached overlay to avoid per-frame allocations (C-004)
	if g.pauseOverlay == nil {
		g.pauseOverlay = ebiten.NewImage(g.width, g.height)
		g.pauseOverlay.Fill(color.RGBA{0, 0, 0, 128})
	}
	screen.DrawImage(g.pauseOverlay, nil)

	ebitenutil.DebugPrintAt(screen, "PAUSED\n\nPress ESC to resume", g.width/3, g.height/2)
}

// drawGameOver renders the game over screen.
func (g *Game) drawGameOver(screen *ebiten.Image) {
	msg := fmt.Sprintf("GAME OVER\n\nTurns: %d\n\nPress ENTER to return to menu", g.turn)
	ebitenutil.DebugPrintAt(screen, msg, g.width/3, g.height/3)
}

// drawDebug renders debug information.
func (g *Game) drawDebug(screen *ebiten.Image) {
	debugMsg := fmt.Sprintf("FPS: %.2f\nTPS: %.2f\nEntities: %d\nTurn: %d\nState: %d\nGenre: %s",
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
		g.world.EntityCount(),
		g.turn,
		g.state,
		g.genre)
	ebitenutil.DebugPrint(screen, debugMsg)
}

// Run starts the game with Ebitengine.
func (g *Game) Run() error {
	ebiten.SetWindowSize(g.width, g.height)
	ebiten.SetWindowTitle("Voyage")
	return ebiten.RunGame(g)
}
