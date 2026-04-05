//go:build !headless

package game

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/opd-ai/voyage/pkg/audio"
)

// NewGameSession creates a new game session with all subsystems initialized.
func NewGameSession(cfg SessionConfig) *GameSession {
	return initializeSession(cfg)
}

// Update implements ebiten.Game.Update.
func (s *GameSession) Update() error {
	s.handleDebugToggle()
	s.handleStateInput()
	return nil
}

// handleDebugToggle toggles debug mode with F3 using proper key release detection.
func (s *GameSession) handleDebugToggle() {
	if ebiten.IsKeyPressed(ebiten.KeyF3) {
		if !s.f3WasPressed {
			s.debugMode = !s.debugMode
		}
		s.f3WasPressed = true
	} else {
		s.f3WasPressed = false
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
func (s *GameSession) handleMenuInput() {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeySpace) {
		s.state = StatePlaying
	}
}

// handlePlayingInput handles input during gameplay.
func (s *GameSession) handlePlayingInput() {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		s.state = StatePaused
		return
	}

	// Handle movement
	moved := s.handleMovement()
	if moved {
		s.advanceTurn()
	}

	// Handle event choices (1-4 keys)
	s.handleEventChoices()
}

// handleMovement processes arrow key input for vessel movement.
// Returns true if the player moved.
func (s *GameSession) handleMovement() bool {
	var newPos world.Point
	moved := false

	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		newPos = world.Point{X: s.playerPos.X, Y: s.playerPos.Y - 1}
		moved = true
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		newPos = world.Point{X: s.playerPos.X, Y: s.playerPos.Y + 1}
		moved = true
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		newPos = world.Point{X: s.playerPos.X - 1, Y: s.playerPos.Y}
		moved = true
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		newPos = world.Point{X: s.playerPos.X + 1, Y: s.playerPos.Y}
		moved = true
	}

	if moved && s.worldMap.IsValidMove(s.playerPos, newPos) {
		s.playerPos = newPos
		return true
	}
	return false
}

// handleEventChoices processes number key input for event choice selection.
func (s *GameSession) handleEventChoices() {
	if !s.eventQueue.HasPending() {
		return
	}

	pending := s.eventQueue.Pending()
	if len(pending) == 0 {
		return
	}

	currentEvent := pending[0]
	choiceKeys := []ebiten.Key{ebiten.Key1, ebiten.Key2, ebiten.Key3, ebiten.Key4}

	for i, key := range choiceKeys {
		if ebiten.IsKeyPressed(key) && i < len(currentEvent.Choices) {
			s.resolveEvent(currentEvent.ID, i)
			break
		}
	}
}

// resolveEvent processes an event choice and applies the outcome.
func (s *GameSession) resolveEvent(eventID, choiceID int) {
	outcome := s.eventQueue.Resolve(eventID, choiceID)
	if outcome == nil {
		return
	}

	// Apply resource changes
	s.resources.Add(resources.ResourceFood, outcome.FoodDelta)
	s.resources.Add(resources.ResourceWater, outcome.WaterDelta)
	s.resources.Add(resources.ResourceFuel, outcome.FuelDelta)
	s.resources.Add(resources.ResourceMedicine, outcome.MedicineDelta)
	s.resources.Add(resources.ResourceMorale, outcome.MoraleDelta)
	s.resources.Add(resources.ResourceCurrency, outcome.CurrencyDelta)

	// Apply vessel damage
	if outcome.VesselDamage > 0 {
		s.vessel.TakeDamage(outcome.VesselDamage)
	}

	// Apply crew damage
	if outcome.CrewDamage > 0 {
		s.party.ApplyDamageToAll(outcome.CrewDamage)
	}

	// Advance time if needed
	for i := 0; i < outcome.TimeAdvance; i++ {
		s.advanceTurn()
	}
}

// advanceTurn processes one turn of gameplay.
func (s *GameSession) advanceTurn() {
	s.turn++

	// Consume resources
	s.consumeResources()

	// Maybe generate event
	s.maybeGenerateEvent()

	// Advance party day
	s.party.AdvanceDay()

	// Check win/lose conditions
	s.checkConditions()
}

// consumeResources depletes resources based on turn progression.
func (s *GameSession) consumeResources() {
	// Consume food and water based on crew size
	crewCount := float64(s.party.LivingCount())
	s.resources.Consume(resources.ResourceFood, crewCount*0.5)
	s.resources.Consume(resources.ResourceWater, crewCount*0.3)

	// Consume fuel based on vessel speed
	s.resources.Consume(resources.ResourceFuel, s.vessel.Speed())

	// Morale changes based on resource status
	if s.resources.IsDepleted(resources.ResourceFood) {
		s.resources.Add(resources.ResourceMorale, -5)
	}
	if s.resources.IsDepleted(resources.ResourceWater) {
		s.resources.Add(resources.ResourceMorale, -8)
	}
}

// maybeGenerateEvent potentially generates an event at the current position.
func (s *GameSession) maybeGenerateEvent() {
	tile := s.worldMap.GetTile(s.playerPos.X, s.playerPos.Y)
	if tile == nil {
		return
	}

	// Higher chance at hazardous terrain
	hazardChance := 0.0
	if tile.Terrain == world.TerrainMountain || tile.Terrain == world.TerrainSwamp {
		hazardChance = 0.2
	}

	if s.eventQueue.ShouldTrigger(hazardChance) {
		s.eventQueue.Generate(s.playerPos.X, s.playerPos.Y, s.turn)
	}
}

// checkConditions checks win/lose conditions.
func (s *GameSession) checkConditions() {
	// Win: reached destination
	if s.playerPos.X == s.worldMap.Destination.X && s.playerPos.Y == s.worldMap.Destination.Y {
		s.state = StateGameOver
		return
	}

	// Lose: all crew dead
	if s.party.IsEmpty() {
		s.state = StateGameOver
		return
	}

	// Lose: vessel destroyed
	if s.vessel.IsDestroyed() {
		s.state = StateGameOver
		return
	}

	// Lose: morale collapsed
	if s.resources.IsDepleted(resources.ResourceMorale) {
		s.state = StateGameOver
		return
	}
}

// handlePausedInput handles input in paused state.
func (s *GameSession) handlePausedInput() {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		s.state = StatePlaying
	}
}

// handleGameOverInput handles input in game over state.
func (s *GameSession) handleGameOverInput() {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
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
	viewHeight := (s.height - 100) / tileSize // Leave space for HUD

	// Calculate view offset to center on player
	offsetX := s.playerPos.X - viewWidth/2
	offsetY := s.playerPos.Y - viewHeight/2

	for screenY := 0; screenY < viewHeight; screenY++ {
		for screenX := 0; screenX < viewWidth; screenX++ {
			mapX := screenX + offsetX
			mapY := screenY + offsetY

			tile := s.worldMap.GetTile(mapX, mapY)
			if tile == nil {
				continue
			}

			// Draw terrain
			tileType := int(tile.Terrain)
			s.renderer.DrawTile(screen, screenX, screenY, tileType)

			// Draw player marker
			if mapX == s.playerPos.X && mapY == s.playerPos.Y {
				s.renderer.DrawTile(screen, screenX, screenY, 10) // Player tile type
			}

			// Draw destination marker
			if mapX == s.worldMap.Destination.X && mapY == s.worldMap.Destination.Y {
				s.renderer.DrawTile(screen, screenX, screenY, 11) // Destination tile type
			}
		}
	}
}

// drawHUD renders the heads-up display.
func (s *GameSession) drawHUD(screen *ebiten.Image) {
	hudY := s.height - 80
	msg := fmt.Sprintf("Turn: %d | Pos: (%d,%d) | Crew: %d/%d | Vessel: %.0f%%\nFood: %.0f | Water: %.0f | Fuel: %.0f | Morale: %.0f | Gold: %.0f",
		s.turn,
		s.playerPos.X, s.playerPos.Y,
		s.party.LivingCount(), s.party.Count(),
		s.vessel.IntegrityRatio()*100,
		s.resources.Get(resources.ResourceFood),
		s.resources.Get(resources.ResourceWater),
		s.resources.Get(resources.ResourceFuel),
		s.resources.Get(resources.ResourceMorale),
		s.resources.Get(resources.ResourceCurrency))
	drawCenteredText(screen, msg, 10, hudY)
}

// drawEventOverlay renders the current event dialog.
func (s *GameSession) drawEventOverlay(screen *ebiten.Image) {
	pending := s.eventQueue.Pending()
	if len(pending) == 0 {
		return
	}

	event := pending[0]
	msg := fmt.Sprintf("=== %s ===\n%s\n\n", event.Title, event.Description)
	for i, choice := range event.Choices {
		msg += fmt.Sprintf("[%d] %s\n", i+1, choice.Text)
	}

	// Draw semi-transparent background
	drawCenteredText(screen, msg, s.width/4, s.height/4)
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
	s.config.Genre = genre
	s.ecsWorld.SetGenre(genre)
	s.party.SetGenre(genre)
	s.relationships.SetGenre(genre)
	s.vessel.SetGenre(genre)
	s.resources.SetGenre(genre)
	s.eventQueue.SetGenre(genre)
	s.audioPlayer.SetGenre(genre)
	s.renderer.SetGenre(genre)
}

// drawCenteredText is a helper to draw debug text.
func drawCenteredText(screen *ebiten.Image, msg string, x, y int) {
	// Use Ebitengine's debug print for simplicity
	// In production, this would use proper font rendering
	ebitenutil.DebugPrintAt(screen, msg, x, y)
}
