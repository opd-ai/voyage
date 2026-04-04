package destination

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestDestinationType(t *testing.T) {
	types := AllDestinationTypes()
	if len(types) != 5 {
		t.Errorf("Expected 5 destination types, got %d", len(types))
	}

	// Test String method
	expected := []string{"City", "Sanctuary", "Treasure", "Escape", "Settlement"}
	for i, dt := range types {
		if dt.String() != expected[i] {
			t.Errorf("Type %d: expected %q, got %q", i, expected[i], dt.String())
		}
	}
}

func TestDiscoveryPhase(t *testing.T) {
	phases := []DiscoveryPhase{Distant, Signs, Approach, Arrival}
	expected := []string{"Distant", "Signs", "Approach", "Arrival"}

	for i, p := range phases {
		if p.String() != expected[i] {
			t.Errorf("Phase %d: expected %q, got %q", i, expected[i], p.String())
		}
	}
}

func TestNewGenerator(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		gen := NewGenerator(genre)
		if gen == nil {
			t.Errorf("NewGenerator(%s) returned nil", genre)
		}
		if gen.genre != genre {
			t.Errorf("Generator genre: expected %s, got %s", genre, gen.genre)
		}
	}
}

func TestGeneratorSetGenre(t *testing.T) {
	gen := NewGenerator(engine.GenreFantasy)
	gen.SetGenre(engine.GenreScifi)
	if gen.genre != engine.GenreScifi {
		t.Errorf("SetGenre failed: expected %s, got %s", engine.GenreScifi, gen.genre)
	}
}

func TestGenerate(t *testing.T) {
	gen := NewGenerator(engine.GenreFantasy)
	dest := gen.Generate(City, 10)

	if dest == nil {
		t.Fatal("Generate returned nil")
	}
	if dest.Type != City {
		t.Errorf("Type: expected %v, got %v", City, dest.Type)
	}
	if dest.Distance != 10 {
		t.Errorf("Distance: expected 10, got %d", dest.Distance)
	}
	if dest.Phase != Distant {
		t.Errorf("Phase: expected %v, got %v", Distant, dest.Phase)
	}
	if dest.Name == "" {
		t.Error("Name should not be empty")
	}
	if dest.Description == "" {
		t.Error("Description should not be empty")
	}
}

func TestGenerateRandom(t *testing.T) {
	gen := NewGenerator(engine.GenreScifi)
	dest := gen.GenerateRandom(15)

	if dest == nil {
		t.Fatal("GenerateRandom returned nil")
	}
	if dest.Distance != 15 {
		t.Errorf("Distance: expected 15, got %d", dest.Distance)
	}
}

func TestGenerateSet(t *testing.T) {
	gen := NewGenerator(engine.GenreHorror)
	destinations := gen.GenerateSet(5, 100)

	if len(destinations) != 5 {
		t.Errorf("Expected 5 destinations, got %d", len(destinations))
	}

	// Check distances increase
	for i := 1; i < len(destinations); i++ {
		if destinations[i].Distance <= destinations[i-1].Distance {
			t.Errorf("Destination %d distance (%d) should be > %d",
				i, destinations[i].Distance, destinations[i-1].Distance)
		}
	}
}

func TestDestinationSetGenre(t *testing.T) {
	gen := NewGenerator(engine.GenreFantasy)
	dest := gen.Generate(City, 10)
	originalName := dest.Name

	dest.SetGenre(engine.GenreCyberpunk)

	if dest.Name == originalName {
		t.Error("Name should change after SetGenre")
	}
	if dest.genre != engine.GenreCyberpunk {
		t.Errorf("Genre: expected %s, got %s", engine.GenreCyberpunk, dest.genre)
	}
}

func TestAdvanceTurn(t *testing.T) {
	gen := NewGenerator(engine.GenrePostapoc)
	dest := gen.Generate(Sanctuary, 3)

	// Initial state
	if dest.Phase != Distant {
		t.Errorf("Initial phase: expected %v, got %v", Distant, dest.Phase)
	}

	// Advance to Approach
	event := dest.AdvanceTurn()
	if dest.Distance != 2 {
		t.Errorf("Distance after turn: expected 2, got %d", dest.Distance)
	}
	if dest.Phase != Approach {
		t.Errorf("Phase: expected %v, got %v", Approach, dest.Phase)
	}
	if event == nil {
		t.Error("Expected discovery event on phase change")
	}

	// Advance again
	dest.AdvanceTurn()
	if dest.Distance != 1 {
		t.Errorf("Distance: expected 1, got %d", dest.Distance)
	}

	// Advance to arrival
	dest.AdvanceTurn()
	if dest.Distance != 0 {
		t.Errorf("Distance: expected 0, got %d", dest.Distance)
	}
	if dest.Phase != Arrival {
		t.Errorf("Phase: expected %v, got %v", Arrival, dest.Phase)
	}
	if !dest.Reached {
		t.Error("Destination should be marked as reached")
	}

	// No more events after arrival
	event = dest.AdvanceTurn()
	if event != nil {
		t.Error("Should not generate events after arrival")
	}
}

func TestDiscoveryEvent(t *testing.T) {
	gen := NewGenerator(engine.GenreCyberpunk)
	dest := gen.Generate(Treasure, 3)

	event := dest.AdvanceTurn() // Phase change to Approach
	if event == nil {
		t.Fatal("Expected event on phase change")
	}
	if event.Phase != Approach {
		t.Errorf("Event phase: expected %v, got %v", Approach, event.Phase)
	}
	if event.Title == "" {
		t.Error("Event title should not be empty")
	}
	if event.Description == "" {
		t.Error("Event description should not be empty")
	}
	if event.Destination != dest {
		t.Error("Event should reference the destination")
	}
}

func TestGetArrivalText(t *testing.T) {
	gen := NewGenerator(engine.GenreFantasy)
	dest := gen.Generate(City, 0)
	dest.Phase = Arrival
	dest.Reached = true

	text := dest.GetArrivalText()
	if text == "" {
		t.Error("Arrival text should not be empty")
	}

	// Check that text contains multiple sections
	if len(text) < 100 {
		t.Errorf("Arrival text seems too short: %d chars", len(text))
	}
}

func TestGetArrivalTextNotReached(t *testing.T) {
	gen := NewGenerator(engine.GenreScifi)
	dest := gen.Generate(Escape, 10)

	text := dest.GetArrivalText()
	if text != "" {
		t.Error("Arrival text should be empty when not reached")
	}
}

func TestAllGenresHaveContent(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		// Check destination names
		names, ok := destinationNames[genre]
		if !ok {
			t.Errorf("Missing destination names for genre %s", genre)
			continue
		}
		for _, dt := range AllDestinationTypes() {
			if _, ok := names[dt]; !ok {
				t.Errorf("Missing name for %s/%s", genre, dt)
			}
		}

		// Check destination descriptions
		descs, ok := destinationDescriptions[genre]
		if !ok {
			t.Errorf("Missing destination descriptions for genre %s", genre)
			continue
		}
		for _, dt := range AllDestinationTypes() {
			if _, ok := descs[dt]; !ok {
				t.Errorf("Missing description for %s/%s", genre, dt)
			}
		}

		// Check prefixes
		if _, ok := destinationPrefixes[genre]; !ok {
			t.Errorf("Missing prefixes for genre %s", genre)
		}

		// Check discovery titles
		titles, ok := discoveryTitles[genre]
		if !ok {
			t.Errorf("Missing discovery titles for genre %s", genre)
			continue
		}
		for _, phase := range []DiscoveryPhase{Distant, Signs, Approach, Arrival} {
			if _, ok := titles[phase]; !ok {
				t.Errorf("Missing discovery titles for %s/%v", genre, phase)
			}
		}

		// Check discovery descriptions
		discDescs, ok := discoveryDescriptions[genre]
		if !ok {
			t.Errorf("Missing discovery descriptions for genre %s", genre)
			continue
		}
		for _, phase := range []DiscoveryPhase{Distant, Signs, Approach, Arrival} {
			if _, ok := discDescs[phase]; !ok {
				t.Errorf("Missing discovery descriptions for %s/%v", genre, phase)
			}
		}

		// Check arrival content
		if _, ok := arrivalOpenings[genre]; !ok {
			t.Errorf("Missing arrival openings for genre %s", genre)
		}
		if _, ok := arrivalReflections[genre]; !ok {
			t.Errorf("Missing arrival reflections for genre %s", genre)
		}
		if _, ok := arrivalClosings[genre]; !ok {
			t.Errorf("Missing arrival closings for genre %s", genre)
		}
		narratives, ok := arrivalNarratives[genre]
		if !ok {
			t.Errorf("Missing arrival narratives for genre %s", genre)
			continue
		}
		for _, dt := range AllDestinationTypes() {
			if _, ok := narratives[dt]; !ok {
				t.Errorf("Missing arrival narrative for %s/%s", genre, dt)
			}
		}
	}
}

func TestDeterministicGeneration(t *testing.T) {
	gen1 := NewGenerator(engine.GenreFantasy)
	gen2 := NewGenerator(engine.GenreFantasy)

	// Same seed should produce same results
	dest1 := gen1.Generate(City, 10)
	dest2 := gen2.Generate(City, 10)

	if dest1.Name != dest2.Name {
		t.Errorf("Names should match: %q vs %q", dest1.Name, dest2.Name)
	}
}
