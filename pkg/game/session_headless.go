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

// consumeResources depletes resources based on turn progression.
func (s *GameSession) consumeResources() {
	crewCount := float64(s.party.LivingCount())
	s.resources.Consume(resources.ResourceFood, crewCount*0.5)
	s.resources.Consume(resources.ResourceWater, crewCount*0.3)
	s.resources.Consume(resources.ResourceFuel, s.vessel.Speed())

	if s.resources.IsDepleted(resources.ResourceFood) {
		s.resources.Add(resources.ResourceMorale, -5)
	}
	if s.resources.IsDepleted(resources.ResourceWater) {
		s.resources.Add(resources.ResourceMorale, -8)
	}
}

// maybeGenerateEvent potentially generates an event.
func (s *GameSession) maybeGenerateEvent() {
	tile := s.worldMap.GetTile(s.playerPos.X, s.playerPos.Y)
	if tile == nil {
		return
	}

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
	if s.playerPos.X == s.worldMap.Destination.X && s.playerPos.Y == s.worldMap.Destination.Y {
		s.state = StateGameOver
		return
	}

	if s.party.IsEmpty() {
		s.state = StateGameOver
		return
	}

	if s.vessel.IsDestroyed() {
		s.state = StateGameOver
		return
	}

	if s.resources.IsDepleted(resources.ResourceMorale) {
		s.state = StateGameOver
		return
	}
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

	s.resources.Add(resources.ResourceFood, outcome.FoodDelta)
	s.resources.Add(resources.ResourceWater, outcome.WaterDelta)
	s.resources.Add(resources.ResourceFuel, outcome.FuelDelta)
	s.resources.Add(resources.ResourceMedicine, outcome.MedicineDelta)
	s.resources.Add(resources.ResourceMorale, outcome.MoraleDelta)
	s.resources.Add(resources.ResourceCurrency, outcome.CurrencyDelta)

	if outcome.VesselDamage > 0 {
		s.vessel.TakeDamage(outcome.VesselDamage)
	}

	if outcome.CrewDamage > 0 {
		s.party.ApplyDamageToAll(outcome.CrewDamage)
	}

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
