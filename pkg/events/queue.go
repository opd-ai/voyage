package events

import (
	"fmt"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// Queue manages the procedural event queue.
type Queue struct {
	gen      *seed.Generator
	genre    engine.GenreID
	nextID   int
	pending  []*Event
	resolved []*Event
}

// NewQueue creates a new event queue with the given seed.
func NewQueue(masterSeed int64, genre engine.GenreID) *Queue {
	return &Queue{
		gen:      seed.NewGenerator(masterSeed, "events"),
		genre:    genre,
		nextID:   1,
		pending:  make([]*Event, 0),
		resolved: make([]*Event, 0),
	}
}

// SetGenre changes the queue's genre.
func (q *Queue) SetGenre(genre engine.GenreID) {
	q.genre = genre
}

// Genre returns the current genre.
func (q *Queue) Genre() engine.GenreID {
	return q.genre
}

// Generate creates a new event based on position and turn.
// The event is deterministic given the same seed, position, and turn.
func (q *Queue) Generate(x, y, turn int) *Event {
	// Create position-specific generator for determinism
	posGen := seed.NewGenerator(q.gen.Int63(), q.positionKey(x, y, turn))

	// Choose category based on weighted distribution
	weights := []float64{0.15, 0.25, 0.20, 0.25, 0.15}
	category := seed.WeightedChoice(posGen, AllEventCategories(), weights)

	event := q.generateForCategory(posGen, category)
	event.ID = q.nextID
	q.nextID++

	q.pending = append(q.pending, event)
	return event
}

// positionKey creates a unique key for position-based generation.
func (q *Queue) positionKey(x, y, turn int) string {
	return fmt.Sprintf("pos_%d_%d_%d", x, y, turn)
}

// generateForCategory creates an event of the given category.
func (q *Queue) generateForCategory(gen *seed.Generator, cat EventCategory) *Event {
	switch cat {
	case CategoryWeather:
		return q.generateWeatherEvent(gen)
	case CategoryEncounter:
		return q.generateEncounterEvent(gen)
	case CategoryDiscovery:
		return q.generateDiscoveryEvent(gen)
	case CategoryHardship:
		return q.generateHardshipEvent(gen)
	case CategoryWindfall:
		return q.generateWindfallEvent(gen)
	default:
		return q.generateHardshipEvent(gen)
	}
}

// Resolve marks an event as resolved and returns the outcome.
func (q *Queue) Resolve(eventID, choiceID int) *EventOutcome {
	for i, event := range q.pending {
		if event.ID == eventID {
			choice := event.GetChoice(choiceID)
			if choice == nil {
				return nil
			}
			// Move to resolved
			q.pending = append(q.pending[:i], q.pending[i+1:]...)
			q.resolved = append(q.resolved, event)
			return &choice.Outcome
		}
	}
	return nil
}

// Pending returns all pending events.
func (q *Queue) Pending() []*Event {
	return q.pending
}

// HasPending returns true if there are pending events.
func (q *Queue) HasPending() bool {
	return len(q.pending) > 0
}

// Clear removes all pending events.
func (q *Queue) Clear() {
	q.pending = make([]*Event, 0)
}

// ResolvedCount returns the number of resolved events.
func (q *Queue) ResolvedCount() int {
	return len(q.resolved)
}

// ShouldTrigger determines if an event should trigger based on position.
// Returns probability from 0 to 1.
func (q *Queue) ShouldTrigger(hazardChance float64) bool {
	baseChance := 0.15 // Base 15% event chance per move
	chance := baseChance + hazardChance*0.5
	return q.gen.Float64() < chance
}

func (q *Queue) generateWeatherEvent(gen *seed.Generator) *Event {
	templates := weatherTemplates[q.genre]
	if len(templates) == 0 {
		templates = weatherTemplates[engine.GenreFantasy]
	}
	tmpl := seed.Choice(gen, templates)

	event := NewEvent(0, CategoryWeather, tmpl.Title, tmpl.Description, q.genre)
	for _, c := range tmpl.Choices {
		event.AddChoice(c.Text, c.Outcome)
	}
	return event
}

func (q *Queue) generateEncounterEvent(gen *seed.Generator) *Event {
	templates := encounterTemplates[q.genre]
	if len(templates) == 0 {
		templates = encounterTemplates[engine.GenreFantasy]
	}
	tmpl := seed.Choice(gen, templates)

	event := NewEvent(0, CategoryEncounter, tmpl.Title, tmpl.Description, q.genre)
	for _, c := range tmpl.Choices {
		event.AddChoice(c.Text, c.Outcome)
	}
	return event
}

func (q *Queue) generateDiscoveryEvent(gen *seed.Generator) *Event {
	templates := discoveryTemplates[q.genre]
	if len(templates) == 0 {
		templates = discoveryTemplates[engine.GenreFantasy]
	}
	tmpl := seed.Choice(gen, templates)

	event := NewEvent(0, CategoryDiscovery, tmpl.Title, tmpl.Description, q.genre)
	for _, c := range tmpl.Choices {
		event.AddChoice(c.Text, c.Outcome)
	}
	return event
}

func (q *Queue) generateHardshipEvent(gen *seed.Generator) *Event {
	templates := hardshipTemplates[q.genre]
	if len(templates) == 0 {
		templates = hardshipTemplates[engine.GenreFantasy]
	}
	tmpl := seed.Choice(gen, templates)

	event := NewEvent(0, CategoryHardship, tmpl.Title, tmpl.Description, q.genre)
	for _, c := range tmpl.Choices {
		event.AddChoice(c.Text, c.Outcome)
	}
	return event
}

func (q *Queue) generateWindfallEvent(gen *seed.Generator) *Event {
	templates := windfallTemplates[q.genre]
	if len(templates) == 0 {
		templates = windfallTemplates[engine.GenreFantasy]
	}
	tmpl := seed.Choice(gen, templates)

	event := NewEvent(0, CategoryWindfall, tmpl.Title, tmpl.Description, q.genre)
	for _, c := range tmpl.Choices {
		event.AddChoice(c.Text, c.Outcome)
	}
	return event
}
