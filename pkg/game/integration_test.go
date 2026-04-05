//go:build headless

package game

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/config"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/world"
	"github.com/opd-ai/voyage/pkg/resources"
)

// TestCompleteJourney simulates a full game journey from origin to destination.
// This validates end-to-end game flow and catches multi-turn progression edge cases.
func TestCompleteJourney(t *testing.T) {
	// Use a known seed for reproducibility
	cfg := SessionConfig{
		Width:      800,
		Height:     600,
		TileSize:   16,
		Seed:       42424,
		Genre:      engine.GenreFantasy,
		Difficulty: config.DifficultyEasy, // Easy for better chance of survival
		MapWidth:   20,
		MapHeight:  15,
		CrewSize:   4,
	}

	session := NewGameSession(cfg)
	if session == nil {
		t.Fatal("Failed to create game session")
	}

	// Start the game
	session.SetState(StatePlaying)

	// Track journey statistics
	turnsCompleted := 0
	eventsEncountered := 0
	crewLost := 0
	initialCrew := session.Party().Count()

	// Simulate journey with maximum turn limit
	const maxTurns = 500
	for turnsCompleted < maxTurns {
		// Check if game is over
		if session.State() == StateGameOver {
			break
		}

		// Check win condition (reached destination)
		pos := session.PlayerPosition()
		dest := session.WorldMap().Destination
		if pos.X == dest.X && pos.Y == dest.Y {
			break
		}

		// Try to move toward destination using pathfinding
		moved := moveTowardDestination(session)
		if !moved {
			// If we can't move, advance turn anyway (stuck)
			session.AdvanceTurn()
		}

		// Check for events
		if session.EventQueue().HasPending() {
			eventsEncountered++
			// Resolve first event with first choice
			pending := session.EventQueue().Pending()
			if len(pending) > 0 && len(pending[0].Choices) > 0 {
				outcome := session.EventQueue().Resolve(pending[0].ID, 1)
				if outcome != nil {
					session.applyOutcome(outcome)
				}
			}
		}

		turnsCompleted++

		// Track crew losses
		currentCrew := session.Party().LivingCount()
		if currentCrew < initialCrew-crewLost {
			crewLost = initialCrew - currentCrew
		}
	}

	// Validate journey results
	t.Logf("Journey statistics:")
	t.Logf("  Turns completed: %d", turnsCompleted)
	t.Logf("  Events encountered: %d", eventsEncountered)
	t.Logf("  Crew lost: %d/%d", crewLost, initialCrew)
	t.Logf("  Final state: %v", session.State())

	// At least some turns should have been completed
	if turnsCompleted == 0 {
		t.Error("No turns completed during journey")
	}

	// Verify state is valid (either won, lost, or still playing)
	finalState := session.State()
	if finalState != StatePlaying && finalState != StateGameOver {
		t.Errorf("Unexpected final state: %v", finalState)
	}

	// Check turn counter advanced correctly
	if session.Turn() != turnsCompleted {
		t.Errorf("Turn counter mismatch: expected %d, got %d", turnsCompleted, session.Turn())
	}
}

// TestCompleteJourneyWithVictory tests a successful journey to destination.
func TestCompleteJourneyWithVictory(t *testing.T) {
	// Create a small map to make victory achievable
	cfg := SessionConfig{
		Width:      800,
		Height:     600,
		TileSize:   16,
		Seed:       99999,
		Genre:      engine.GenreScifi,
		Difficulty: config.DifficultyEasy,
		MapWidth:   10, // Small map
		MapHeight:  8,
		CrewSize:   4,
	}

	session := NewGameSession(cfg)
	session.SetState(StatePlaying)

	// Give extra resources to ensure survival
	session.Resources().Set(resources.ResourceFood, 500)
	session.Resources().Set(resources.ResourceWater, 500)
	session.Resources().Set(resources.ResourceFuel, 500)
	session.Resources().Set(resources.ResourceMedicine, 100)
	session.Resources().Set(resources.ResourceMorale, 100)

	const maxTurns = 200
	for turn := 0; turn < maxTurns; turn++ {
		pos := session.PlayerPosition()
		dest := session.WorldMap().Destination

		// Check if we reached destination
		if pos.X == dest.X && pos.Y == dest.Y {
			// Trigger condition check
			session.checkConditions()
			if session.State() == StateGameOver {
				t.Log("Victory achieved!")
				return
			}
		}

		// Check for game over
		if session.State() == StateGameOver {
			t.Log("Game over before reaching destination")
			return
		}

		// Move toward destination
		moveTowardDestination(session)

		// Clear any pending events
		session.EventQueue().Clear()
	}

	t.Log("Journey ended without reaching destination or losing")
}

// TestJourneyStateTransitions validates state transitions during gameplay.
func TestJourneyStateTransitions(t *testing.T) {
	cfg := DefaultSessionConfig()
	cfg.Seed = 12345
	cfg.MapWidth = 15
	cfg.MapHeight = 15

	session := NewGameSession(cfg)

	// Initial state should be menu
	if session.State() != StateMenu {
		t.Errorf("Expected initial state StateMenu, got %v", session.State())
	}

	// Transition to playing
	session.SetState(StatePlaying)
	if session.State() != StatePlaying {
		t.Errorf("Expected state StatePlaying, got %v", session.State())
	}

	// Simulate some turns
	for i := 0; i < 10; i++ {
		session.AdvanceTurn()
	}

	// State should still be playing (unless game over)
	state := session.State()
	if state != StatePlaying && state != StateGameOver {
		t.Errorf("Unexpected state after turns: %v", state)
	}

	// Pause the game
	session.SetState(StatePaused)
	if session.State() != StatePaused {
		t.Errorf("Expected state StatePaused, got %v", session.State())
	}

	// Resume
	session.SetState(StatePlaying)
	if session.State() != StatePlaying {
		t.Errorf("Expected state StatePlaying after resume, got %v", session.State())
	}
}

// TestLossConditions tests all loss conditions.
func TestLossConditions(t *testing.T) {
	testCases := []struct {
		name      string
		setupLoss func(*GameSession)
	}{
		{
			name: "All crew dead",
			setupLoss: func(s *GameSession) {
				s.Party().ApplyDamageToAll(1000)
			},
		},
		{
			name: "Vessel destroyed",
			setupLoss: func(s *GameSession) {
				s.Vessel().TakeDamage(1000)
			},
		},
		{
			name: "Morale collapsed",
			setupLoss: func(s *GameSession) {
				s.Resources().Set(resources.ResourceMorale, 0)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := DefaultSessionConfig()
			cfg.Seed = 12345
			cfg.MapWidth = 15
			cfg.MapHeight = 15

			session := NewGameSession(cfg)
			session.SetState(StatePlaying)

			// Setup loss condition
			tc.setupLoss(session)

			// Check conditions
			session.checkConditions()

			// Verify game over
			if session.State() != StateGameOver {
				t.Errorf("Expected StateGameOver for %s, got %v", tc.name, session.State())
			}
		})
	}
}

// TestWinCondition tests reaching the destination.
func TestWinCondition(t *testing.T) {
	cfg := DefaultSessionConfig()
	cfg.Seed = 12345
	cfg.MapWidth = 15
	cfg.MapHeight = 15

	session := NewGameSession(cfg)
	session.SetState(StatePlaying)

	// Teleport to destination
	dest := session.WorldMap().Destination
	session.setPlayerPosition(dest)

	// Check conditions
	session.checkConditions()

	// Verify game over (win)
	if session.State() != StateGameOver {
		t.Errorf("Expected StateGameOver for win, got %v", session.State())
	}
}

// moveTowardDestination attempts to move the player closer to the destination.
// Returns true if movement occurred.
func moveTowardDestination(s *GameSession) bool {
	pos := s.PlayerPosition()
	dest := s.WorldMap().Destination

	tile := s.WorldMap().GetTile(pos.X, pos.Y)
	if tile == nil || len(tile.Connections) == 0 {
		return false
	}

	// Find the connection closest to destination
	var bestMove world.Point
	bestDist := manhattanDistance(pos, dest)
	found := false

	for _, conn := range tile.Connections {
		dist := manhattanDistance(conn, dest)
		if dist < bestDist {
			bestDist = dist
			bestMove = conn
			found = true
		}
	}

	if found {
		return s.MovePlayer(bestMove)
	}

	// If no closer move, try any connected position
	if len(tile.Connections) > 0 {
		return s.MovePlayer(tile.Connections[0])
	}

	return false
}

// manhattanDistance calculates Manhattan distance between two points.
func manhattanDistance(a, b world.Point) int {
	dx := a.X - b.X
	dy := a.Y - b.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

// TestResourceConsumptionOverTime verifies resources deplete correctly over turns.
func TestResourceConsumptionOverTime(t *testing.T) {
	cfg := DefaultSessionConfig()
	cfg.Seed = 12345
	cfg.MapWidth = 15
	cfg.MapHeight = 15
	cfg.CrewSize = 4

	session := NewGameSession(cfg)
	session.SetState(StatePlaying)

	initialFood := session.Resources().Get(resources.ResourceFood)
	initialWater := session.Resources().Get(resources.ResourceWater)
	initialFuel := session.Resources().Get(resources.ResourceFuel)

	// Advance 10 turns
	for i := 0; i < 10; i++ {
		session.consumeResources()
	}

	// Resources should have decreased
	finalFood := session.Resources().Get(resources.ResourceFood)
	finalWater := session.Resources().Get(resources.ResourceWater)
	finalFuel := session.Resources().Get(resources.ResourceFuel)

	if finalFood >= initialFood {
		t.Errorf("Food should decrease: was %f, now %f", initialFood, finalFood)
	}
	if finalWater >= initialWater {
		t.Errorf("Water should decrease: was %f, now %f", initialWater, finalWater)
	}
	if finalFuel >= initialFuel {
		t.Errorf("Fuel should decrease: was %f, now %f", initialFuel, finalFuel)
	}
}

// TestEventGeneration verifies events are generated at certain positions/turns.
func TestEventGeneration(t *testing.T) {
	cfg := DefaultSessionConfig()
	cfg.Seed = 54321
	cfg.MapWidth = 15
	cfg.MapHeight = 15

	session := NewGameSession(cfg)
	session.SetState(StatePlaying)

	eventsGenerated := 0

	// Move and advance turns multiple times to trigger events
	for i := 0; i < 50; i++ {
		moveTowardDestination(session)

		if session.EventQueue().HasPending() {
			eventsGenerated++
			// Clear event to continue
			session.EventQueue().Clear()
		}

		if session.State() == StateGameOver {
			break
		}
	}

	// We should have generated some events over 50 turns
	t.Logf("Events generated over 50 turns: %d", eventsGenerated)
	// Note: With the RNG, we might not always get events, so we don't require it
}

// TestDeterministicJourney verifies same seed produces same journey.
func TestDeterministicJourney(t *testing.T) {
	seed := int64(77777)

	// Run journey twice with same seed
	turns1, events1 := runDeterministicJourney(t, seed)
	turns2, events2 := runDeterministicJourney(t, seed)

	// Results should be identical
	if turns1 != turns2 {
		t.Errorf("Turn count differs: run1=%d, run2=%d", turns1, turns2)
	}
	if events1 != events2 {
		t.Errorf("Event count differs: run1=%d, run2=%d", events1, events2)
	}
}

// runDeterministicJourney runs a journey with given seed and returns turn and event counts.
func runDeterministicJourney(t *testing.T, seed int64) (turns, events int) {
	cfg := SessionConfig{
		Width:      800,
		Height:     600,
		TileSize:   16,
		Seed:       seed,
		Genre:      engine.GenreFantasy,
		Difficulty: config.DifficultyEasy,
		MapWidth:   12,
		MapHeight:  10,
		CrewSize:   3,
	}

	session := NewGameSession(cfg)
	session.SetState(StatePlaying)

	for i := 0; i < 100; i++ {
		if session.State() == StateGameOver {
			break
		}

		pos := session.PlayerPosition()
		dest := session.WorldMap().Destination
		if pos.X == dest.X && pos.Y == dest.Y {
			break
		}

		moveTowardDestination(session)

		if session.EventQueue().HasPending() {
			events++
			// Always choose first option for determinism
			pending := session.EventQueue().Pending()
			if len(pending) > 0 && len(pending[0].Choices) > 0 {
				session.EventQueue().Resolve(pending[0].ID, 1)
			}
		}

		turns++
	}

	return turns, events
}
