package modding

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadExampleSteampunkMod(t *testing.T) {
	// Find the example mod file
	examplePath := filepath.Join("..", "..", "examples", "mods", "steampunk-genre.json")

	// Check if file exists
	if _, err := os.Stat(examplePath); os.IsNotExist(err) {
		t.Skip("Example mod file not found, skipping")
	}

	loader := NewLoader()

	mod, err := loader.LoadFromFile(examplePath)
	if err != nil {
		t.Fatalf("Failed to load example mod: %v", err)
	}

	// Verify mod metadata
	if mod.ID != "steampunk-genre" {
		t.Errorf("expected ID 'steampunk-genre', got %q", mod.ID)
	}
	if mod.Name != "Steampunk Genre Preset" {
		t.Errorf("expected name 'Steampunk Genre Preset', got %q", mod.Name)
	}
	if mod.Version != "1.0.0" {
		t.Errorf("expected version '1.0.0', got %q", mod.Version)
	}
	if mod.Author != "Voyage Team" {
		t.Errorf("expected author 'Voyage Team', got %q", mod.Author)
	}

	// Verify genres
	if len(mod.Genres) != 1 {
		t.Fatalf("expected 1 genre, got %d", len(mod.Genres))
	}

	genre := mod.Genres[0]
	if genre.ID != "steampunk" {
		t.Errorf("expected genre ID 'steampunk', got %q", genre.ID)
	}
	if len(genre.Biomes) < 5 {
		t.Errorf("expected at least 5 biomes, got %d", len(genre.Biomes))
	}
	if len(genre.Resources) < 4 {
		t.Errorf("expected at least 4 resources, got %d", len(genre.Resources))
	}
	if len(genre.Factions) < 4 {
		t.Errorf("expected at least 4 factions, got %d", len(genre.Factions))
	}

	// Verify events
	if len(mod.Events) < 5 {
		t.Fatalf("expected at least 5 events, got %d", len(mod.Events))
	}

	// Check that all events have valid categories
	for i, event := range mod.Events {
		if !isValidCategory(event.Category) {
			t.Errorf("event %d has invalid category: %q", i, event.Category)
		}
		if event.Genre != "steampunk" {
			t.Errorf("event %d has wrong genre: %q", i, event.Genre)
		}
		if len(event.Choices) == 0 {
			t.Errorf("event %d has no choices", i)
		}
	}

	// Verify biome additions
	if len(mod.Biomes) < 1 {
		t.Errorf("expected at least 1 biome addition, got %d", len(mod.Biomes))
	}

	// Verify custom resources
	if len(mod.Resources) < 2 {
		t.Errorf("expected at least 2 resources, got %d", len(mod.Resources))
	}

	// Verify factions
	if len(mod.Factions) < 3 {
		t.Errorf("expected at least 3 factions, got %d", len(mod.Factions))
	}
}

func TestExampleModEventsForGenre(t *testing.T) {
	examplePath := filepath.Join("..", "..", "examples", "mods", "steampunk-genre.json")

	if _, err := os.Stat(examplePath); os.IsNotExist(err) {
		t.Skip("Example mod file not found, skipping")
	}

	loader := NewLoader()

	_, err := loader.LoadFromFile(examplePath)
	if err != nil {
		t.Fatalf("Failed to load example mod: %v", err)
	}

	// Get events for steampunk genre
	events := loader.GetEventsForGenre("steampunk")
	if len(events) < 5 {
		t.Errorf("expected at least 5 steampunk events, got %d", len(events))
	}

	// Verify event categories are diverse
	categories := make(map[string]int)
	for _, e := range events {
		categories[e.Category]++
	}

	if len(categories) < 3 {
		t.Errorf("expected at least 3 different categories, got %d", len(categories))
	}
}

func TestExampleModCustomGenres(t *testing.T) {
	examplePath := filepath.Join("..", "..", "examples", "mods", "steampunk-genre.json")

	if _, err := os.Stat(examplePath); os.IsNotExist(err) {
		t.Skip("Example mod file not found, skipping")
	}

	loader := NewLoader()

	_, err := loader.LoadFromFile(examplePath)
	if err != nil {
		t.Fatalf("Failed to load example mod: %v", err)
	}

	// Get custom genres
	genres := loader.GetCustomGenres()
	if len(genres) != 1 {
		t.Fatalf("expected 1 custom genre, got %d", len(genres))
	}

	steampunk := genres[0]
	if steampunk.ID != "steampunk" {
		t.Errorf("expected steampunk genre, got %q", steampunk.ID)
	}

	// Verify genre has all required vocabulary
	if len(steampunk.VesselTypes) == 0 {
		t.Error("steampunk genre missing vessel types")
	}
	if len(steampunk.CrewRoles) == 0 {
		t.Error("steampunk genre missing crew roles")
	}
	if len(steampunk.CategoryNames) == 0 {
		t.Error("steampunk genre missing category names")
	}
}
