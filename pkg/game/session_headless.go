//go:build headless

package game

import (
	"github.com/opd-ai/voyage/pkg/audio"
	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/events"
	"github.com/opd-ai/voyage/pkg/procgen/world"
	"github.com/opd-ai/voyage/pkg/rendering"
	"github.com/opd-ai/voyage/pkg/resources"
	"github.com/opd-ai/voyage/pkg/vessel"
)

// NewGameSession creates a new game session with all subsystems initialized.
func NewGameSession(cfg SessionConfig) *GameSession {
	return initializeSession(cfg)
}

// Run is a no-op in headless mode (no Ebitengine window).
func (s *GameSession) Run() error {
	return nil
}

// Update advances the game state (headless version).
func (s *GameSession) Update() error {
	return nil
}

// AdvanceTurn processes one turn of gameplay.
func (s *GameSession) AdvanceTurn() {
	s.turn++
	s.consumeResources()
	s.maybeGenerateEvent()
	s.party.AdvanceDay()
	s.checkConditions()
}

// MovePlayer moves the player to a new position if valid.
func (s *GameSession) MovePlayer(newPos world.Point) bool {
	if s.worldMap.IsValidMove(s.playerPos, newPos) {
		s.playerPos = newPos
		s.AdvanceTurn()
		return true
	}
	return false
}

// ResolveEvent processes an event choice.
func (s *GameSession) ResolveEvent(eventID, choiceID int) {
	outcome := s.eventQueue.Resolve(eventID, choiceID)
	if outcome == nil {
		return
	}

	s.applyOutcome(outcome)

	for i := 0; i < outcome.TimeAdvance; i++ {
		s.AdvanceTurn()
	}
}

// Accessors

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

// SetState sets the game state.
func (s *GameSession) SetState(state GameState) {
	s.state = state
}

// setPlayerPosition sets the player's position for testing purposes.
func (s *GameSession) setPlayerPosition(pos world.Point) {
	s.playerPos = pos
}

// SetGenre changes the genre for all subsystems.
func (s *GameSession) SetGenre(genre engine.GenreID) {
	s.propagateGenre(genre)
}
