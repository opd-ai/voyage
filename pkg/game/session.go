//go:build !headless

package game

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/opd-ai/voyage/pkg/audio"
	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/events"
	"github.com/opd-ai/voyage/pkg/input"
	"github.com/opd-ai/voyage/pkg/procgen/world"
	"github.com/opd-ai/voyage/pkg/rendering"
	"github.com/opd-ai/voyage/pkg/resources"
	"github.com/opd-ai/voyage/pkg/vessel"
)

// NewGameSession creates a new game session with all subsystems initialized.
func NewGameSession(cfg SessionConfig) *GameSession {
	session := initializeSession(cfg)
	// Initialize input manager for non-headless builds
	session.inputMgr = input.NewManager()
	return session
}

// Update implements ebiten.Game.Update.
func (s *GameSession) Update() error {
	// Update input manager to process this frame's input
	s.inputMgr.Update()

	s.handleDebugToggle()
	s.handleStateInput()

	// Snapshot current event for Draw() synchronization (C-004)
	s.snapshotCurrentEvent()

	// Cache display strings for Draw() to reduce allocations (H-003)
	s.updateCachedStrings()

	return nil
}

// snapshotCurrentEvent captures the current event state for Draw() to use.
// This prevents desynchronization between Update() and Draw() (C-004).
func (s *GameSession) snapshotCurrentEvent() {
	pending := s.eventQueue.Pending()
	if len(pending) > 0 {
		// Copy the event to prevent race conditions
		s.currentEventSnapshot = pending[0]
	} else {
		s.currentEventSnapshot = nil
	}
}

// updateCachedStrings pre-builds display strings to reduce Draw() allocations (H-003).
func (s *GameSession) updateCachedStrings() {
	// Only rebuild HUD text when state changes (marked dirty by advanceTurn, movement, etc.)
	if s.hudDirty || s.cachedHUDText == "" {
		s.cachedHUDText = fmt.Sprintf("Turn: %d | Pos: (%d,%d) | Crew: %d/%d | Vessel: %.0f%%\nFood: %.0f | Water: %.0f | Fuel: %.0f | Morale: %.0f | Gold: %.0f",
			s.turn,
			s.playerPos.X, s.playerPos.Y,
			s.party.LivingCount(), s.party.Count(),
			s.vessel.IntegrityRatio()*100,
			s.resources.Get(resources.ResourceFood),
			s.resources.Get(resources.ResourceWater),
			s.resources.Get(resources.ResourceFuel),
			s.resources.Get(resources.ResourceMorale),
			s.resources.Get(resources.ResourceCurrency))
		s.hudDirty = false
	}

	// Cache event text using strings.Builder to reduce allocations (H-003)
	if event := s.currentEventSnapshot; event != nil {
		var builder strings.Builder
		builder.WriteString("=== ")
		builder.WriteString(event.Title)
		builder.WriteString(" ===\n")
		builder.WriteString(event.Description)
		builder.WriteString("\n\n")
		for i, choice := range event.Choices {
			builder.WriteString("[")
			builder.WriteString(fmt.Sprintf("%d", i+1))
			builder.WriteString("] ")
			builder.WriteString(choice.Text)
			builder.WriteString("\n")
		}
		s.cachedEventText = builder.String()
	} else {
		s.cachedEventText = ""
	}
}

// handleDebugToggle toggles debug mode with F3 using proper key release detection.
func (s *GameSession) handleDebugToggle() {
	if s.inputMgr.JustDebugToggled() {
		s.debugMode = !s.debugMode
	}
}

// handleStateInput processes input based on current game state.
func (s *GameSession) handleStateInput() {
	switch s.state {
	case StateMenu:
		s.handleMenuInput()
	case StatePlaying:
		s.handlePlayingInput()
	case StatePaused:
		s.handlePausedInput()
	case StateGameOver:
		s.handleGameOverInput()
	}
}

// handleMenuInput handles input in menu state.
// Uses Input Manager for clean single-press detection.
func (s *GameSession) handleMenuInput() {
	if s.inputMgr.JustConfirmed() {
		s.state = StatePlaying
	}
}

// handlePlayingInput handles input during gameplay.
func (s *GameSession) handlePlayingInput() {
	// Use Input Manager for ESC to prevent rapid toggling
	if s.inputMgr.JustCancelled() {
		s.state = StatePaused
		return
	}

	// Per-frame action budget: only one action per Update() call (C-003)
	// Handle movement first
	moved := s.handleMovement()
	if moved {
		s.advanceTurn()
		return // Prevent event handling in the same frame
	}

	// Handle event choices (1-4 keys) only if no movement occurred
	s.handleEventChoices()
}

// handleMovement processes arrow key input for vessel movement.
// Returns true if the player moved.
func (s *GameSession) handleMovement() bool {
	newPos, moved := s.getMovementInput()
	if moved && s.worldMap.IsValidMove(s.playerPos, newPos) {
		s.playerPos = newPos
		s.hudDirty = true // Mark HUD for refresh (H-003)
		return true
	}
	return false
}

// getMovementInput checks for directional key presses and returns the target position.
// getMovementInput reads directional input from the Input Manager.
func (s *GameSession) getMovementInput() (world.Point, bool) {
	dir := s.inputMgr.GetDirection()
	switch dir {
	case input.DirectionUp:
		return world.Point{X: s.playerPos.X, Y: s.playerPos.Y - 1}, true
	case input.DirectionDown:
		return world.Point{X: s.playerPos.X, Y: s.playerPos.Y + 1}, true
	case input.DirectionLeft:
		return world.Point{X: s.playerPos.X - 1, Y: s.playerPos.Y}, true
	case input.DirectionRight:
		return world.Point{X: s.playerPos.X + 1, Y: s.playerPos.Y}, true
	default:
		return world.Point{}, false
	}
}

// handleEventChoices processes number key input for event choice selection.
// Uses Input Manager to prevent duplicate resource consumption.
func (s *GameSession) handleEventChoices() {
	if !s.eventQueue.HasPending() {
		return
	}

	pending := s.eventQueue.Pending()
	if len(pending) == 0 {
		return
	}

	currentEvent := pending[0]

	// Get option pressed from input manager (returns 1-9 or 0 if none)
	option := s.inputMgr.GetOptionPressed()
	if option > 0 && option <= len(currentEvent.Choices) {
		// Option numbers are 1-indexed, choice IDs are also 1-indexed
		s.resolveEvent(currentEvent.ID, option)
	}
}

// resolveEvent processes an event choice and applies the outcome.
func (s *GameSession) resolveEvent(eventID, choiceID int) {
	outcome := s.eventQueue.Resolve(eventID, choiceID)
	if outcome == nil {
		return
	}

	s.applyOutcome(outcome)

	// Advance time if needed, clamped to prevent malformed data from freezing game (M-009)
	const maxTimeAdvance = 100
	timeAdvance := outcome.TimeAdvance
	if timeAdvance < 0 {
		timeAdvance = 0
	}
	if timeAdvance > maxTimeAdvance {
		timeAdvance = maxTimeAdvance
	}
	for i := 0; i < timeAdvance; i++ {
		s.advanceTurn()
	}
}

// advanceTurn processes one turn of gameplay.
func (s *GameSession) advanceTurn() {
	s.turn++
	s.hudDirty = true // Mark HUD for refresh (H-003)

	// Consume resources
	s.consumeResources()

	// Advance party day
	s.party.AdvanceDay()

	// Check win/lose conditions before generating new events
	// This prevents events from being generated/queued after game over
	s.checkConditions()

	// Only generate events if game is still active
	if s.state == StatePlaying {
		s.maybeGenerateEvent()
	}
}

// handlePausedInput handles input in paused state.
func (s *GameSession) handlePausedInput() {
	// Use Input Manager for ESC to prevent rapid toggling
	if s.inputMgr.JustCancelled() {
		s.state = StatePlaying
	}
}

// handleGameOverInput handles input in game over state.
func (s *GameSession) handleGameOverInput() {
	if s.inputMgr.JustConfirmed() {
		s.state = StateMenu
		// Reset would go here for new game
	}
}

// Draw implements ebiten.Game.Draw.
func (s *GameSession) Draw(screen *ebiten.Image) {
	// Clear screen
	screen.Fill(s.renderer.Palette().Background)

	switch s.state {
	case StateMenu:
		s.drawMenu(screen)
	case StatePlaying:
		s.drawGame(screen)
	case StatePaused:
		s.drawGame(screen)
		s.drawPauseOverlay(screen)
	case StateGameOver:
		s.drawGameOver(screen)
	}

	if s.debugMode {
		s.drawDebug(screen)
	}
}

// drawMenu renders the main menu.
func (s *GameSession) drawMenu(screen *ebiten.Image) {
	msg := fmt.Sprintf("VOYAGE\n\nGenre: %s\nSeed: %d\nCrew: %d\nVessel: %s\n\nPress ENTER or SPACE to start",
		s.config.Genre, s.config.Seed, s.party.Count(), s.vessel.Name())
	drawCenteredText(screen, msg, s.width/4, s.height/3)
}

// drawGame renders the main gameplay view.
func (s *GameSession) drawGame(screen *ebiten.Image) {
	// Draw world map centered on player
	s.drawWorldMap(screen)

	// Draw HUD
	s.drawHUD(screen)

	// Draw pending events
	if s.eventQueue.HasPending() {
		s.drawEventOverlay(screen)
	}
}

// drawWorldMap renders the world map.
func (s *GameSession) drawWorldMap(screen *ebiten.Image) {
	tileSize := s.renderer.TileSize()
	viewWidth := s.width / tileSize
	viewHeight := (s.height - 100) / tileSize

	offsetX := s.playerPos.X - viewWidth/2
	offsetY := s.playerPos.Y - viewHeight/2

	for screenY := 0; screenY < viewHeight; screenY++ {
		for screenX := 0; screenX < viewWidth; screenX++ {
			s.drawTileAt(screen, screenX, screenY, offsetX, offsetY)
		}
	}
}

// drawTileAt renders a single tile at the given screen position.
func (s *GameSession) drawTileAt(screen *ebiten.Image, screenX, screenY, offsetX, offsetY int) {
	mapX := screenX + offsetX
	mapY := screenY + offsetY

	tile := s.worldMap.GetTile(mapX, mapY)
	if tile == nil {
		return
	}

	s.renderer.DrawTile(screen, screenX, screenY, int(tile.Terrain))

	if mapX == s.playerPos.X && mapY == s.playerPos.Y {
		s.renderer.DrawTile(screen, screenX, screenY, 10)
	}
	if mapX == s.worldMap.Destination.X && mapY == s.worldMap.Destination.Y {
		s.renderer.DrawTile(screen, screenX, screenY, 11)
	}
}

// drawHUD renders the heads-up display using cached text (H-003).
func (s *GameSession) drawHUD(screen *ebiten.Image) {
	hudY := s.height - 80
	drawCenteredText(screen, s.cachedHUDText, 10, hudY)
}

// drawEventOverlay renders the current event dialog.
// Uses snapshotted event and cached text to prevent allocations (C-004, H-003).
func (s *GameSession) drawEventOverlay(screen *ebiten.Image) {
	if s.cachedEventText == "" {
		return
	}
	drawCenteredText(screen, s.cachedEventText, s.width/4, s.height/4)
}

// drawPauseOverlay renders the pause screen.
func (s *GameSession) drawPauseOverlay(screen *ebiten.Image) {
	msg := "=== PAUSED ===\n\nPress ESC to resume"
	drawCenteredText(screen, msg, s.width/3, s.height/2)
}

// drawGameOver renders the game over screen.
func (s *GameSession) drawGameOver(screen *ebiten.Image) {
	result := "JOURNEY ENDED"
	if s.playerPos.X == s.worldMap.Destination.X && s.playerPos.Y == s.worldMap.Destination.Y {
		result = "VICTORY!"
	}

	msg := fmt.Sprintf("=== %s ===\n\nTurns: %d\nCrew Survived: %d/%d\nVessel: %.0f%%\n\nPress ENTER to return to menu",
		result, s.turn, s.party.LivingCount(), s.party.Count(), s.vessel.IntegrityRatio()*100)
	drawCenteredText(screen, msg, s.width/4, s.height/3)
}

// drawDebug renders debug information.
func (s *GameSession) drawDebug(screen *ebiten.Image) {
	msg := fmt.Sprintf("FPS: %.2f | TPS: %.2f | Entities: %d | Events: %d",
		ebiten.ActualFPS(), ebiten.ActualTPS(), s.ecsWorld.EntityCount(), s.eventQueue.ResolvedCount())
	drawCenteredText(screen, msg, 10, 10)
}

// Layout implements ebiten.Game.Layout.
func (s *GameSession) Layout(outsideWidth, outsideHeight int) (int, int) {
	return s.width, s.height
}

// Run starts the game session with Ebitengine.
func (s *GameSession) Run() error {
	ebiten.SetWindowSize(s.width, s.height)
	ebiten.SetWindowTitle("Voyage")
	return ebiten.RunGame(s)
}

// Accessors for subsystems

// World returns the ECS world.
func (s *GameSession) World() *engine.World {
	return s.ecsWorld
}

// WorldMap returns the generated world map.
func (s *GameSession) WorldMap() *world.WorldMap {
	return s.worldMap
}

// Party returns the crew party.
func (s *GameSession) Party() *crew.Party {
	return s.party
}

// Vessel returns the vessel.
func (s *GameSession) Vessel() *vessel.Vessel {
	return s.vessel
}

// Resources returns the resource manager.
func (s *GameSession) Resources() *resources.Resources {
	return s.resources
}

// EventQueue returns the event queue.
func (s *GameSession) EventQueue() *events.Queue {
	return s.eventQueue
}

// AudioPlayer returns the audio player.
func (s *GameSession) AudioPlayer() *audio.Player {
	return s.audioPlayer
}

// Renderer returns the renderer.
func (s *GameSession) Renderer() *rendering.Renderer {
	return s.renderer
}

// State returns the current game state.
func (s *GameSession) State() GameState {
	return s.state
}

// Turn returns the current turn number.
func (s *GameSession) Turn() int {
	return s.turn
}

// PlayerPosition returns the player's current position.
func (s *GameSession) PlayerPosition() world.Point {
	return s.playerPos
}

// SetGenre changes the genre for all subsystems.
func (s *GameSession) SetGenre(genre engine.GenreID) {
	s.propagateGenre(genre)
}

// drawCenteredText is a helper to draw debug text.
func drawCenteredText(screen *ebiten.Image, msg string, x, y int) {
	// Use Ebitengine's debug print for simplicity
	// In production, this would use proper font rendering
	ebitenutil.DebugPrintAt(screen, msg, x, y)
}
