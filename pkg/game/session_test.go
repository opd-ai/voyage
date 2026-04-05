//go:build headless

package game

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/audio"
	"github.com/opd-ai/voyage/pkg/config"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/world"
	"github.com/opd-ai/voyage/pkg/resources"
)

func TestNewGameSession(t *testing.T) {
	cfg := SessionConfig{
		Width:      800,
		Height:     600,
		TileSize:   16,
		Seed:       12345,
		Genre:      engine.GenreFantasy,
		Difficulty: config.DifficultyNormal,
		MapWidth:   20,
		MapHeight:  20,
		CrewSize:   4,
	}

	session := NewGameSession(cfg)
	if session == nil {
		t.Fatal("NewGameSession returned nil")
	}

	// Verify subsystems are initialized
	if session.World() == nil {
		t.Error("ECS world not initialized")
	}
	if session.WorldMap() == nil {
		t.Error("World map not initialized")
	}
	if session.Party() == nil {
		t.Error("Party not initialized")
	}
	if session.Vessel() == nil {
		t.Error("Vessel not initialized")
	}
	if session.Resources() == nil {
		t.Error("Resources not initialized")
	}
	if session.EventQueue() == nil {
		t.Error("Event queue not initialized")
	}
	if session.AudioPlayer() == nil {
		t.Error("Audio player not initialized")
	}
	if session.Renderer() == nil {
		t.Error("Renderer not initialized")
	}

	// Verify initial state
	if session.State() != StateMenu {
		t.Errorf("expected initial state StateMenu, got %v", session.State())
	}
	if session.Turn() != 0 {
		t.Errorf("expected turn 0, got %d", session.Turn())
	}

	// Verify crew was generated
	if session.Party().Count() != cfg.CrewSize {
		t.Errorf("expected %d crew members, got %d", cfg.CrewSize, session.Party().Count())
	}

	// Verify player is at origin
	origin := session.WorldMap().Origin
	pos := session.PlayerPosition()
	if pos.X != origin.X || pos.Y != origin.Y {
		t.Errorf("expected player at origin (%d,%d), got (%d,%d)", origin.X, origin.Y, pos.X, pos.Y)
	}
}

func TestGameSessionAdvanceTurn(t *testing.T) {
	cfg := DefaultSessionConfig()
	cfg.Seed = 12345
	cfg.MapWidth = 20
	cfg.MapHeight = 20

	session := NewGameSession(cfg)
	initialTurn := session.Turn()
	initialFood := session.Resources().Get(resources.ResourceFood)

	// Advance turn
	session.AdvanceTurn()

	// Verify turn advanced
	if session.Turn() != initialTurn+1 {
		t.Errorf("expected turn %d, got %d", initialTurn+1, session.Turn())
	}

	// Verify resources were consumed
	currentFood := session.Resources().Get(resources.ResourceFood)
	if currentFood >= initialFood {
		t.Error("expected food to be consumed after turn advance")
	}
}

func TestGameSessionMovePlayer(t *testing.T) {
	cfg := DefaultSessionConfig()
	cfg.Seed = 12345
	cfg.MapWidth = 20
	cfg.MapHeight = 20

	session := NewGameSession(cfg)
	startPos := session.PlayerPosition()

	// Find a valid adjacent position
	worldMap := session.WorldMap()
	tile := worldMap.GetTile(startPos.X, startPos.Y)
	if len(tile.Connections) == 0 {
		t.Skip("Origin has no connections for testing movement")
	}

	// Move to first connected position
	newPos := tile.Connections[0]
	moved := session.MovePlayer(world.Point{X: newPos.X, Y: newPos.Y})

	if !moved {
		t.Error("expected move to succeed")
	}

	currentPos := session.PlayerPosition()
	if currentPos.X != newPos.X || currentPos.Y != newPos.Y {
		t.Errorf("expected player at (%d,%d), got (%d,%d)", newPos.X, newPos.Y, currentPos.X, currentPos.Y)
	}
}

func TestGameSessionSetGenre(t *testing.T) {
	cfg := DefaultSessionConfig()
	cfg.Genre = engine.GenreFantasy

	session := NewGameSession(cfg)

	// Change genre
	session.SetGenre(engine.GenreScifi)

	// Verify all subsystems updated
	if session.Party().Genre() != engine.GenreScifi {
		t.Error("Party genre not updated")
	}
	if session.Vessel().Genre() != engine.GenreScifi {
		t.Error("Vessel genre not updated")
	}
	if session.Resources().Genre() != engine.GenreScifi {
		t.Error("Resources genre not updated")
	}
}

func TestDefaultSessionConfig(t *testing.T) {
	cfg := DefaultSessionConfig()

	if cfg.Width != 800 {
		t.Errorf("expected Width 800, got %d", cfg.Width)
	}
	if cfg.Height != 600 {
		t.Errorf("expected Height 600, got %d", cfg.Height)
	}
	if cfg.TileSize != 16 {
		t.Errorf("expected TileSize 16, got %d", cfg.TileSize)
	}
	if cfg.Genre != engine.GenreFantasy {
		t.Errorf("expected Genre Fantasy, got %v", cfg.Genre)
	}
	if cfg.Difficulty != config.DifficultyNormal {
		t.Errorf("expected Difficulty Normal, got %v", cfg.Difficulty)
	}
	if cfg.MapWidth != 50 {
		t.Errorf("expected MapWidth 50, got %d", cfg.MapWidth)
	}
	if cfg.MapHeight != 50 {
		t.Errorf("expected MapHeight 50, got %d", cfg.MapHeight)
	}
	if cfg.CrewSize != 4 {
		t.Errorf("expected CrewSize 4, got %d", cfg.CrewSize)
	}
}

func TestMusicStateTransitions(t *testing.T) {
	cfg := DefaultSessionConfig()
	cfg.Seed = 12345
	cfg.MapWidth = 20
	cfg.MapHeight = 20

	session := NewGameSession(cfg)

	// Initial music state should be peaceful
	initialState := session.AudioPlayer().MusicState()
	if initialState != audio.MusicPeaceful {
		t.Errorf("expected initial music state Peaceful, got %s", audio.MusicStateName(initialState))
	}

	// Simulate advancing turns - music may change based on terrain/events
	for i := 0; i < 5; i++ {
		session.AdvanceTurn()
	}

	// Verify audio player still works
	if session.AudioPlayer() == nil {
		t.Error("audio player should not be nil after turns")
	}
}

func TestMusicStateOnGameOver(t *testing.T) {
	cfg := DefaultSessionConfig()
	cfg.Seed = 12345
	cfg.MapWidth = 20
	cfg.MapHeight = 20

	session := NewGameSession(cfg)

	// Force a loss condition by killing all crew
	party := session.Party()
	for party.LivingCount() > 0 {
		// Apply damage until all crew are dead
		party.ApplyDamageToAll(100)
	}

	// Trigger condition check
	session.AdvanceTurn()

	// State should be game over with death music
	if session.State() != StateGameOver {
		t.Errorf("expected StateGameOver, got %v", session.State())
	}

	// Music should be death state
	musicState := session.AudioPlayer().MusicState()
	if musicState != audio.MusicDeath {
		t.Errorf("expected music state Death on game over, got %s", audio.MusicStateName(musicState))
	}
}
