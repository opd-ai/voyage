// Package main provides the entry point for the Voyage travel simulator.
//
// Voyage is a 100% procedural travel simulator where every map, event, crew,
// vessel, audio, and narrative is generated from a single seed. The game
// supports five genre themes: fantasy, sci-fi, horror, cyberpunk, and post-apocalyptic.
//
// Usage:
//
//	voyage                           # Start with random seed
//	voyage --seed 12345              # Start with specific seed
//	voyage --genre scifi             # Start with sci-fi theme
//	voyage --difficulty hard         # Start with hard difficulty
//	voyage --help                    # Show all options
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/opd-ai/voyage/pkg/config"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

var (
	// Command line flags
	seedFlag       = flag.Int64("seed", 0, "Master seed for procedural generation (0 = random)")
	genreFlag      = flag.String("genre", "fantasy", "Genre theme: fantasy, scifi, horror, cyberpunk, postapoc")
	difficultyFlag = flag.String("difficulty", "normal", "Difficulty: easy, normal, hard, nightmare")
	versionFlag    = flag.Bool("version", false, "Print version information")
)

// Version information (set via ldflags in release builds)
var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Voyage %s (built %s)\n", Version, BuildTime)
		os.Exit(0)
	}

	// Validate genre
	if !engine.IsValidGenre(*genreFlag) {
		log.Fatalf("Invalid genre: %s. Valid genres: fantasy, scifi, horror, cyberpunk, postapoc", *genreFlag)
	}
	genre := engine.GenreID(*genreFlag)

	// Validate difficulty
	if !config.IsValidDifficulty(*difficultyFlag) {
		log.Fatalf("Invalid difficulty: %s. Valid: easy, normal, hard, nightmare", *difficultyFlag)
	}
	difficulty, _ := config.ParseDifficulty(*difficultyFlag)

	// Initialize seed
	masterSeed := *seedFlag
	if masterSeed == 0 {
		masterSeed = time.Now().UnixNano()
	}
	fmt.Printf("Voyage starting with seed: %d, genre: %s, difficulty: %s\n",
		masterSeed, genre, config.DifficultyName(difficulty))

	// Initialize world
	registry := engine.NewComponentRegistry()
	world := engine.NewWorld(registry)
	world.SetGenre(genre)

	// Create subsystem generators
	_ = seed.NewGenerator(masterSeed, "world")
	_ = seed.NewGenerator(masterSeed, "events")
	_ = seed.NewGenerator(masterSeed, "crew")

	// TODO: Initialize game systems
	// - World map generation
	// - Rendering system
	// - Resource management
	// - Crew system
	// - Vessel system
	// - Event system
	// - Audio synthesis
	// - UI/HUD

	fmt.Println("Voyage v1.0 - Core engine initialized")
	fmt.Println("Note: This is a foundation build. Full gameplay coming soon.")
	fmt.Printf("World has %d entities\n", world.EntityCount())

	// Placeholder: In full implementation, this would start the Ebitengine game loop
	fmt.Println("\nPress Ctrl+C to exit")

	// For now, just run a simple demonstration
	demo(world, masterSeed)
}

// demo demonstrates the basic ECS functionality
func demo(world *engine.World, masterSeed int64) {
	fmt.Println("\n=== ECS Demo ===")

	// Create some entities
	player := world.SpawnImmediate()
	player.AddTag("player")
	fmt.Printf("Created player entity with ID: %d\n", player.ID())

	// Demonstrate deterministic generation
	g := seed.NewGenerator(masterSeed, "demo")
	fmt.Printf("\nFirst 5 random values from seed %d:\n", masterSeed)
	for i := 0; i < 5; i++ {
		fmt.Printf("  %d: %d\n", i+1, g.Intn(100))
	}

	// Show genre
	fmt.Printf("\nCurrent genre: %s\n", world.Genre())
	fmt.Printf("All genres: %v\n", engine.AllGenres())

	fmt.Println("\n=== Demo Complete ===")
}
