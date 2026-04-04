package events

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/resources"
	"github.com/opd-ai/voyage/pkg/vessel"
)

func TestNewEvent(t *testing.T) {
	event := NewEvent(1, CategoryWeather, "Storm", "A storm approaches", engine.GenreFantasy)

	if event.ID != 1 {
		t.Errorf("ID = %d, want 1", event.ID)
	}
	if event.Category != CategoryWeather {
		t.Errorf("Category = %d, want %d", event.Category, CategoryWeather)
	}
	if event.Title != "Storm" {
		t.Errorf("Title = %s, want Storm", event.Title)
	}
	if len(event.Choices) != 0 {
		t.Errorf("Choices = %d, want 0", len(event.Choices))
	}
}

func TestEventChoices(t *testing.T) {
	event := NewEvent(1, CategoryWeather, "Storm", "A storm approaches", engine.GenreFantasy)

	event.AddChoice("Hide", EventOutcome{MoraleDelta: -5})
	event.AddChoice("Press on", EventOutcome{VesselDamage: 10})

	if len(event.Choices) != 2 {
		t.Errorf("Choices = %d, want 2", len(event.Choices))
	}

	choice := event.GetChoice(1)
	if choice == nil {
		t.Fatal("choice 1 should exist")
	}
	if choice.Text != "Hide" {
		t.Errorf("choice 1 text = %s, want Hide", choice.Text)
	}

	choice2 := event.GetChoice(2)
	if choice2.Outcome.VesselDamage != 10 {
		t.Errorf("choice 2 vessel damage = %f, want 10", choice2.Outcome.VesselDamage)
	}

	nilChoice := event.GetChoice(99)
	if nilChoice != nil {
		t.Error("choice 99 should not exist")
	}
}

func TestEventQueue(t *testing.T) {
	queue := NewQueue(12345, engine.GenreFantasy)

	// Generate events at different positions
	event1 := queue.Generate(0, 0, 1)
	event2 := queue.Generate(1, 0, 1)
	event3 := queue.Generate(0, 0, 1)

	if event1.ID == event2.ID {
		t.Error("different positions should produce different event IDs")
	}

	// Same position should produce same event content (determinism)
	// but queue assigns incrementing IDs
	if event3.ID == event1.ID {
		t.Error("should have different IDs even for same position")
	}

	// Check pending
	if len(queue.Pending()) != 3 {
		t.Errorf("pending = %d, want 3", len(queue.Pending()))
	}
}

func TestEventQueueDeterminism(t *testing.T) {
	// Two queues with same seed should generate same events
	queue1 := NewQueue(12345, engine.GenreFantasy)
	queue2 := NewQueue(12345, engine.GenreFantasy)

	event1 := queue1.Generate(5, 5, 10)
	event2 := queue2.Generate(5, 5, 10)

	if event1.Category != event2.Category {
		t.Error("same seed should produce same category")
	}
	if event1.Title != event2.Title {
		t.Error("same seed should produce same title")
	}
}

func TestEventQueueResolution(t *testing.T) {
	queue := NewQueue(12345, engine.GenreFantasy)

	event := queue.Generate(0, 0, 1)
	eventID := event.ID

	if len(event.Choices) == 0 {
		t.Fatal("event should have choices")
	}

	// Resolve the event
	outcome := queue.Resolve(eventID, 1)
	if outcome == nil {
		t.Fatal("resolution should return outcome")
	}

	// Event should be removed from pending
	if queue.HasPending() {
		// There might still be pending events if we generated multiple
		for _, e := range queue.Pending() {
			if e.ID == eventID {
				t.Error("resolved event should not be in pending")
			}
		}
	}

	// Resolved count should increase
	if queue.ResolvedCount() != 1 {
		t.Errorf("resolved count = %d, want 1", queue.ResolvedCount())
	}
}

func TestEventResolver(t *testing.T) {
	resolver := NewResolver()
	res := resources.NewResources(engine.GenreFantasy)
	gen := crew.NewGenerator(12345, engine.GenreFantasy)
	party := crew.NewParty(engine.GenreFantasy, 4)
	party.Add(gen.Generate())
	v := vessel.NewVessel(vessel.VesselMedium, engine.GenreFantasy)

	initialFood := res.Get(resources.ResourceFood)
	initialMorale := res.Get(resources.ResourceMorale)
	initialIntegrity := v.Integrity()

	outcome := &EventOutcome{
		FoodDelta:    -10,
		MoraleDelta:  5,
		VesselDamage: 15,
	}

	result := resolver.Apply(outcome, res, party, v)

	// Check resources changed
	if res.Get(resources.ResourceFood) != initialFood-10 {
		t.Errorf("food = %f, want %f", res.Get(resources.ResourceFood), initialFood-10)
	}
	if res.Get(resources.ResourceMorale) != initialMorale+5 {
		t.Errorf("morale = %f, want %f", res.Get(resources.ResourceMorale), initialMorale+5)
	}
	if v.Integrity() != initialIntegrity-15 {
		t.Errorf("integrity = %f, want %f", v.Integrity(), initialIntegrity-15)
	}

	_ = result
}

func TestEventResolverCrewDamage(t *testing.T) {
	resolver := NewResolver()
	res := resources.NewResources(engine.GenreFantasy)
	gen := crew.NewGenerator(12345, engine.GenreFantasy)
	party := crew.NewParty(engine.GenreFantasy, 4)

	member := gen.Generate()
	member.Health = 10 // Low health
	party.Add(member)

	v := vessel.NewVessel(vessel.VesselMedium, engine.GenreFantasy)

	outcome := &EventOutcome{
		CrewDamage: 15, // Lethal to low-health crew
	}

	result := resolver.Apply(outcome, res, party, v)

	// Should have deaths
	if len(result.Deaths) == 0 {
		t.Error("should have crew deaths")
	}
}

func TestCategoryNames(t *testing.T) {
	genres := engine.AllGenres()
	categories := AllEventCategories()

	for _, g := range genres {
		for _, c := range categories {
			name := CategoryName(c, g)
			if name == "" {
				t.Errorf("missing category name for genre=%s, cat=%d", g, c)
			}
		}
	}
}

func TestOutcomeSeverity(t *testing.T) {
	// Positive outcome
	positive := &EventOutcome{
		FoodDelta:  20,
		WaterDelta: 10,
	}
	if OutcomeSeverity(positive) <= 0 {
		t.Error("positive outcome should have positive severity")
	}

	// Negative outcome
	negative := &EventOutcome{
		CrewDamage:   20,
		VesselDamage: 20,
	}
	if OutcomeSeverity(negative) >= 0 {
		t.Error("negative outcome should have negative severity")
	}

	// Neutral outcome
	neutral := &EventOutcome{}
	if OutcomeSeverity(neutral) != 0 {
		t.Error("empty outcome should have zero severity")
	}
}

func TestShouldTrigger(t *testing.T) {
	queue := NewQueue(12345, engine.GenreFantasy)

	// Run multiple times to check trigger mechanism
	triggerCount := 0
	for i := 0; i < 100; i++ {
		if queue.ShouldTrigger(0.5) { // 50% hazard terrain
			triggerCount++
		}
	}

	// Should trigger some but not all
	if triggerCount == 0 {
		t.Error("should trigger at least some events")
	}
	if triggerCount == 100 {
		t.Error("should not trigger every time")
	}
}

func TestGenreTemplates(t *testing.T) {
	genres := engine.AllGenres()

	for _, g := range genres {
		queue := NewQueue(12345, g)

		// Generate multiple events to test templates
		for i := 0; i < 10; i++ {
			event := queue.Generate(i, 0, 1)
			if event.Title == "" {
				t.Errorf("genre %s: event should have title", g)
			}
			if len(event.Choices) == 0 {
				t.Errorf("genre %s: event should have choices", g)
			}
		}
	}
}
