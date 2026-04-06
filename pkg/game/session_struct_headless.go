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

	// HUD dirty flag for cached string updates (H-003)
	hudDirty bool

	// Tutorial manager for onboarding hints
	tutorial *TutorialManager

	// Screen dimensions
	width  int
	height int
}
