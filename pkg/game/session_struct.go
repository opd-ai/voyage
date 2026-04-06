//go:build !headless

package game

import (
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

	// Input manager for unified input handling
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
	cachedHUDText           string
	cachedEventText         string
	cachedDestinationText   string
	cachedTutorialHintText  string
	cachedTutorialHintPhase TutorialPhase
	cachedTutorialHintValid bool
	hudDirty                bool

	// Tutorial manager for onboarding hints
	tutorial *TutorialManager

	// Screen dimensions
	width  int
	height int
}
