package game

import (
	"github.com/opd-ai/voyage/pkg/audio"
	"github.com/opd-ai/voyage/pkg/config"
	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/events"
	"github.com/opd-ai/voyage/pkg/input"
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

	// Input manager for unified input handling (nil in headless mode)
	inputMgr *input.Manager

	// Game state
	state        GameState
	turn         int
	playerPos    world.Point
	debugMode    bool
	f3WasPressed bool

	// Event snapshot for Draw synchronization (C-004)
	// Updated in Update(), read in Draw() to prevent desync
	currentEventSnapshot *events.Event

	// Cached strings for Draw to reduce allocations (H-003)
	cachedHUDText   string
	cachedEventText string
	hudDirty        bool

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

	// Calculate morale penalties for resource depletion, capped to prevent unfair stacking
	var totalPenalty float64
	if s.resources.IsDepleted(resources.ResourceFood) {
		totalPenalty += 5
	}
	if s.resources.IsDepleted(resources.ResourceWater) {
		totalPenalty += 8
	}
	// Cap combined penalty per turn to prevent rapid morale collapse
	if totalPenalty > 10 {
		totalPenalty = 10
	}
	if totalPenalty > 0 {
		s.resources.Add(resources.ResourceMorale, -totalPenalty)
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
		hudDirty:      true, // Force initial HUD text generation (H-003)
		width:         cfg.Width,
		height:        cfg.Height,
	}
}

// maybeGenerateEvent potentially generates an event at the current position.
// This is shared between headless and non-headless builds.
func (s *GameSession) maybeGenerateEvent() {
	tile := s.worldMap.GetTile(s.playerPos.X, s.playerPos.Y)
	if tile == nil {
		return
	}

	// Higher chance at hazardous terrain
	hazardChance := 0.0
	if tile.Terrain == world.TerrainMountain || tile.Terrain == world.TerrainSwamp {
		hazardChance = 0.2
		// Tense music for hazardous terrain
		s.audioPlayer.SetMusicState(audio.MusicTense)
	} else {
		// Peaceful music for normal travel
		s.audioPlayer.SetMusicState(audio.MusicPeaceful)
	}

	if s.eventQueue.ShouldTrigger(hazardChance) {
		s.eventQueue.Generate(s.playerPos.X, s.playerPos.Y, s.turn)
		// Combat music when an event triggers
		s.audioPlayer.SetMusicState(audio.MusicCombat)
	}
}

// checkConditions checks win/lose conditions.
// This is shared between headless and non-headless builds.
func (s *GameSession) checkConditions() {
	// Win: reached destination
	if s.playerPos.X == s.worldMap.Destination.X && s.playerPos.Y == s.worldMap.Destination.Y {
		s.state = StateGameOver
		s.audioPlayer.SetMusicState(audio.MusicVictory)
		return
	}

	// Lose: all crew dead
	if s.party.IsEmpty() {
		s.state = StateGameOver
		s.audioPlayer.SetMusicState(audio.MusicDeath)
		return
	}

	// Lose: vessel destroyed
	if s.vessel.IsDestroyed() {
		s.state = StateGameOver
		s.audioPlayer.SetMusicState(audio.MusicDeath)
		return
	}

	// Lose: morale collapsed
	if s.resources.IsDepleted(resources.ResourceMorale) {
		s.state = StateGameOver
		s.audioPlayer.SetMusicState(audio.MusicDeath)
		return
	}
}
