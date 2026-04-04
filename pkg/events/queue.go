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

// categoryTemplates maps event categories to their template sources.
var categoryTemplates = map[EventCategory]map[engine.GenreID][]EventTemplate{
	CategoryWeather:   weatherTemplates,
	CategoryEncounter: encounterTemplates,
	CategoryDiscovery: discoveryTemplates,
	CategoryHardship:  hardshipTemplates,
	CategoryWindfall:  windfallTemplates,
}

// generateForCategory creates an event of the given category using the
// appropriate template set. Falls back to fantasy genre if no templates
// exist for the current genre.
func (q *Queue) generateForCategory(gen *seed.Generator, cat EventCategory) *Event {
	templateMap, ok := categoryTemplates[cat]
	if !ok {
		templateMap = categoryTemplates[CategoryHardship]
	}
	return q.generateEventFromTemplates(gen, cat, templateMap)
}

// generateEventFromTemplates creates an event by selecting from the given
// template map based on the queue's current genre.
func (q *Queue) generateEventFromTemplates(gen *seed.Generator, cat EventCategory, templateMap map[engine.GenreID][]EventTemplate) *Event {
	templates := templateMap[q.genre]
	if len(templates) == 0 {
		templates = templateMap[engine.GenreFantasy]
	}
	tmpl := seed.Choice(gen, templates)

	event := NewEvent(0, cat, tmpl.Title, tmpl.Description, q.genre)
	for _, c := range tmpl.Choices {
		event.AddChoice(c.Text, c.Outcome)
	}
	return event
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
