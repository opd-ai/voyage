package game

import (
	"github.com/opd-ai/voyage/pkg/audio"
	"github.com/opd-ai/voyage/pkg/config"
	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/events"
	"github.com/opd-ai/voyage/pkg/procgen/world"
	"github.com/opd-ai/voyage/pkg/rendering"
	"github.com/opd-ai/voyage/pkg/resources"
	"github.com/opd-ai/voyage/pkg/vessel"
)

// SessionConfig holds configuration for a game session.
type SessionConfig struct {
	Width      int
	Height     int
	TileSize   int
	Seed       int64
	Genre      engine.GenreID
	Difficulty config.Difficulty
	MapWidth   int
	MapHeight  int
	CrewSize   int
}

// DefaultSessionConfig returns sensible defaults for a game session.
func DefaultSessionConfig() SessionConfig {
	return SessionConfig{
		Width:      800,
		Height:     600,
		TileSize:   16,
		Seed:       0,
		Genre:      engine.GenreFantasy,
		Difficulty: config.DifficultyNormal,
		MapWidth:   50,
		MapHeight:  50,
		CrewSize:   4,
	}
}

// GameSession orchestrates all game subsystems.
type GameSession struct {
	// Core configuration
	config SessionConfig

	// ECS world
	ecsWorld *engine.World

	// Subsystems
	worldMap      *world.WorldMap
	party         *crew.Party
	relationships *crew.RelationshipNetwork
	vessel        *vessel.Vessel
	resources     *resources.Resources
	eventQueue    *events.Queue
	audioPlayer   *audio.Player
	renderer      *rendering.Renderer

	// Game state
	state        GameState
	turn         int
	playerPos    world.Point
	debugMode    bool
	f3WasPressed bool

	// Screen dimensions
	width  int
	height int
}

// applyOutcome applies an event outcome to game state.
func (s *GameSession) applyOutcome(outcome *events.EventOutcome) {
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
}

// propagateGenre sets genre on all subsystems.
func (s *GameSession) propagateGenre(genre engine.GenreID) {
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

// initializeSession creates and initializes all subsystems for a game session.
// This is the common initialization logic shared by both headless and non-headless builds.
func initializeSession(cfg SessionConfig) *GameSession {
	// Initialize ECS world
	registry := engine.NewComponentRegistry()
	ecsWorld := engine.NewWorld(registry)
	ecsWorld.SetGenre(cfg.Genre)

	// Generate world map
	worldGen := world.NewGenerator(cfg.Seed, cfg.Genre)
	worldMap := worldGen.Generate(cfg.MapWidth, cfg.MapHeight)

	// Create crew party
	party := crew.NewParty(cfg.Genre, cfg.CrewSize)
	crewGen := crew.NewGenerator(cfg.Seed, cfg.Genre)
	for i := 0; i < cfg.CrewSize; i++ {
		member := crewGen.Generate()
		party.Add(member)
	}

	// Initialize relationship network
	relationships := crew.NewRelationshipNetwork(cfg.Genre)

	// Create vessel (medium by default)
	vesselInstance := vessel.NewVessel(vessel.VesselMedium, cfg.Genre)

	// Initialize resources
	resourceMgr := resources.NewResources(cfg.Genre)

	// Create event queue
	eventQueue := events.NewQueue(cfg.Seed, cfg.Genre)

	// Initialize audio
	audioPlayer := audio.NewPlayer(cfg.Seed, cfg.Genre)

	// Initialize renderer
	renderer := rendering.NewRenderer(cfg.Width, cfg.Height, cfg.TileSize)
	renderer.SetGenre(cfg.Genre)

	return &GameSession{
		config:        cfg,
		ecsWorld:      ecsWorld,
		worldMap:      worldMap,
		party:         party,
		relationships: relationships,
		vessel:        vesselInstance,
		resources:     resourceMgr,
		eventQueue:    eventQueue,
		audioPlayer:   audioPlayer,
		renderer:      renderer,
		state:         StateMenu,
		turn:          0,
		playerPos:     worldMap.Origin,
		debugMode:     false,
		f3WasPressed:  false,
		width:         cfg.Width,
		height:        cfg.Height,
	}
}
