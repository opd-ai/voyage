package events

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

func TestGenerateCrewRelationshipEvent(t *testing.T) {
	gen := seed.NewGenerator(12345, "test")
	party := crew.NewParty(engine.GenreFantasy, 4)
	network := crew.NewRelationshipNetwork(engine.GenreFantasy)

	// Add crew members
	crewGen := crew.NewGenerator(12345, engine.GenreFantasy)
	for i := 0; i < 4; i++ {
		party.Add(crewGen.Generate())
	}

	// Initially no strong relationships
	event := GenerateCrewRelationshipEvent(gen, network, party, engine.GenreFantasy)
	if event != nil {
		t.Error("expected nil event with no strong relationships")
	}

	// Create a strong friendly relationship
	members := party.Members()
	network.Interact(members[0].ID, members[1].ID, 60) // Strong positive

	// Now should generate an event
	event = GenerateCrewRelationshipEvent(gen, network, party, engine.GenreFantasy)
	if event == nil {
		t.Error("expected event for strong relationship")
		return
	}

	if event.Category != CategoryCrew {
		t.Errorf("expected CategoryCrew, got %v", event.Category)
	}

	if len(event.Choices) == 0 {
		t.Error("expected event to have choices")
	}
}

func TestGenerateCrewRelationshipEventRivalry(t *testing.T) {
	gen := seed.NewGenerator(54321, "rivalry")
	party := crew.NewParty(engine.GenreFantasy, 4)
	network := crew.NewRelationshipNetwork(engine.GenreFantasy)

	// Add crew members
	crewGen := crew.NewGenerator(54321, engine.GenreFantasy)
	for i := 0; i < 4; i++ {
		party.Add(crewGen.Generate())
	}

	// Create a rivalry
	members := party.Members()
	network.Interact(members[0].ID, members[1].ID, -50) // Strong negative

	event := GenerateCrewRelationshipEvent(gen, network, party, engine.GenreFantasy)
	if event == nil {
		t.Error("expected event for rivalry relationship")
		return
	}

	// Event should have choices
	if len(event.Choices) == 0 {
		t.Error("expected event to have choices")
	}
}

func TestGenerateCrewRelationshipEventNilInputs(t *testing.T) {
	gen := seed.NewGenerator(12345, "test")

	// Test nil network
	event := GenerateCrewRelationshipEvent(gen, nil, nil, engine.GenreFantasy)
	if event != nil {
		t.Error("expected nil for nil network")
	}

	// Test nil party
	network := crew.NewRelationshipNetwork(engine.GenreFantasy)
	event = GenerateCrewRelationshipEvent(gen, network, nil, engine.GenreFantasy)
	if event != nil {
		t.Error("expected nil for nil party")
	}
}

func TestFormatRelationshipText(t *testing.T) {
	tests := []struct {
		input    string
		nameA    string
		nameB    string
		expected string
	}{
		{"%A and %B are friends", "Alice", "Bob", "Alice and Bob are friends"},
		{"%A helped %B", "Charlie", "Diana", "Charlie helped Diana"},
		{"No names here", "X", "Y", "No names here"},
		{"%A%B", "A", "B", "AB"},
	}

	for _, tt := range tests {
		result := formatRelationshipText(tt.input, tt.nameA, tt.nameB)
		if result != tt.expected {
			t.Errorf("formatRelationshipText(%q, %q, %q) = %q, want %q",
				tt.input, tt.nameA, tt.nameB, result, tt.expected)
		}
	}
}

func TestRelationshipEventAllGenres(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		gen := seed.NewGenerator(12345, "genre_test")
		party := crew.NewParty(genre, 4)
		network := crew.NewRelationshipNetwork(genre)

		crewGen := crew.NewGenerator(12345, genre)
		for i := 0; i < 4; i++ {
			party.Add(crewGen.Generate())
		}

		// Create strong relationship
		members := party.Members()
		network.Interact(members[0].ID, members[1].ID, 70)

		event := GenerateCrewRelationshipEvent(gen, network, party, genre)
		if event == nil {
			t.Errorf("expected event for genre %v", genre)
			continue
		}

		if event.Genre != genre {
			t.Errorf("expected genre %v, got %v", genre, event.Genre)
		}
	}
}
