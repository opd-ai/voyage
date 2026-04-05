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
	"github.com/opd-ai/voyage/pkg/game"
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

	// Create session configuration
	cfg := game.SessionConfig{
		Width:      800,
		Height:     600,
		TileSize:   16,
		Seed:       masterSeed,
		Genre:      genre,
		Difficulty: difficulty,
		MapWidth:   50,
		MapHeight:  50,
		CrewSize:   4,
	}

	// Create and run game session
	session := game.NewGameSession(cfg)

	fmt.Printf("World map generated: %dx%d\n", cfg.MapWidth, cfg.MapHeight)
	fmt.Printf("Crew size: %d\n", session.Party().Count())
	fmt.Printf("Vessel: %s\n", session.Vessel().Name())
	fmt.Println("All systems initialized. Starting game...")

	if err := session.Run(); err != nil {
		log.Fatal(err)
	}
}
