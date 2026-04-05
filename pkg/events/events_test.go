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

func TestHazardVocabulary(t *testing.T) {
	expected := map[engine.GenreID][]string{
		engine.GenreFantasy:   {"Magic Storm", "Cursed Grounds"},
		engine.GenreScifi:     {"Asteroid Field", "Ion Storm"},
		engine.GenreHorror:    {"Zombie Horde", "Infected Zone"},
		engine.GenreCyberpunk: {"Netrunner Ambush", "Corporate Drone Swarm"},
		engine.GenrePostapoc:  {"Radiation Storm", "Mutant Swarm"},
	}

	for genre, expectedNames := range expected {
		names := HazardVocabulary(genre)
		if len(names) != len(expectedNames) {
			t.Errorf("genre %s: got %d hazards, want %d", genre, len(names), len(expectedNames))
			continue
		}
		for i, name := range expectedNames {
			if names[i] != name {
				t.Errorf("genre %s: hazard %d = %s, want %s", genre, i, names[i], name)
			}
		}
	}
}

func TestCategoryHazardExists(t *testing.T) {
	categories := AllEventCategories()
	found := false
	for _, c := range categories {
		if c == CategoryHazard {
			found = true
			break
		}
	}
	if !found {
		t.Error("CategoryHazard should be in AllEventCategories")
	}
}

func TestCategoryHazardNames(t *testing.T) {
	genres := engine.AllGenres()
	for _, g := range genres {
		name := CategoryName(CategoryHazard, g)
		if name == "" {
			t.Errorf("genre %s: CategoryHazard should have a name", g)
		}
	}
}

func TestAllEventTemplates(t *testing.T) {
	genres := engine.AllGenres()
	for _, g := range genres {
		templates := AllEventTemplates(g)
		if len(templates) != 6 {
			t.Errorf("genre %s: got %d categories, want 6", g, len(templates))
		}
		// Check hazard templates exist
		hazards, ok := templates[CategoryHazard]
		if !ok {
			t.Errorf("genre %s: missing hazard templates", g)
		}
		if len(hazards) < 2 {
			t.Errorf("genre %s: got %d hazard templates, want at least 2", g, len(hazards))
		}
	}
}

func TestCrewEventTemplates(t *testing.T) {
	genres := engine.AllGenres()
	for _, g := range genres {
		templates := GetCrewEventTemplates(g)
		if len(templates) < 3 {
			t.Errorf("genre %s: got %d crew templates, want at least 3", g, len(templates))
		}

		// Check for each event type
		hasCrisis := false
		hasMilestone := false
		hasSacrifice := false
		for _, tmpl := range templates {
			switch tmpl.Type {
			case CrewEventCrisis:
				hasCrisis = true
			case CrewEventMilestone:
				hasMilestone = true
			case CrewEventSacrifice:
				hasSacrifice = true
			}
			// Check templates have placeholders
			if tmpl.Title == "" {
				t.Errorf("genre %s: crew event should have title", g)
			}
			if len(tmpl.Choices) == 0 {
				t.Errorf("genre %s: crew event %s should have choices", g, tmpl.Title)
			}
		}

		if !hasCrisis {
			t.Errorf("genre %s: should have crisis event", g)
		}
		if !hasMilestone {
			t.Errorf("genre %s: should have milestone event", g)
		}
		if !hasSacrifice {
			t.Errorf("genre %s: should have sacrifice event", g)
		}
	}
}

func TestCategoryCrewNames(t *testing.T) {
	genres := engine.AllGenres()
	for _, g := range genres {
		name := CategoryName(CategoryCrew, g)
		if name == "" {
			t.Errorf("genre %s: CategoryCrew should have a name", g)
		}
	}
}

func TestCategoryCrewInAllCategories(t *testing.T) {
	categories := AllEventCategories()
	found := false
	for _, c := range categories {
		if c == CategoryCrew {
			found = true
			break
		}
	}
	if !found {
		t.Error("CategoryCrew should be in AllEventCategories")
	}
}

func TestQueueGenre(t *testing.T) {
	queue := NewQueue(12345, engine.GenreFantasy)

	if queue.Genre() != engine.GenreFantasy {
		t.Errorf("Genre() = %v, want %v", queue.Genre(), engine.GenreFantasy)
	}

	queue.SetGenre(engine.GenreScifi)
	if queue.Genre() != engine.GenreScifi {
		t.Errorf("Genre() after SetGenre = %v, want %v", queue.Genre(), engine.GenreScifi)
	}
}

func TestQueueClear(t *testing.T) {
	queue := NewQueue(12345, engine.GenreFantasy)

	// Generate some events
	queue.Generate(0, 0, 1)
	queue.Generate(1, 0, 1)
	queue.Generate(0, 1, 1)

	if !queue.HasPending() {
		t.Fatal("queue should have pending events")
	}

	queue.Clear()

	if queue.HasPending() {
		t.Error("queue should be empty after Clear")
	}
	if len(queue.Pending()) != 0 {
		t.Errorf("Pending() = %d, want 0", len(queue.Pending()))
	}
}

func TestResolverCanChoose(t *testing.T) {
	resolver := NewResolver()
	gen := crew.NewGenerator(12345, engine.GenreFantasy)
	party := crew.NewParty(engine.GenreFantasy, 4)
	member := gen.Generate()
	party.Add(member)

	// Choice without skill requirement
	choiceNoSkill := &Choice{
		Text:         "Regular choice",
		RequireSkill: "",
	}
	if !resolver.CanChoose(choiceNoSkill, party) {
		t.Error("should be able to choose when no skill required")
	}

	// Choice requiring a skill the party has
	skillName := crew.SkillName(member.Skill, party.Genre())
	choiceWithSkill := &Choice{
		Text:         "Skilled choice",
		RequireSkill: skillName,
	}
	if !resolver.CanChoose(choiceWithSkill, party) {
		t.Error("should be able to choose when party has the skill")
	}

	// Choice requiring a skill the party doesn't have
	choiceImpossible := &Choice{
		Text:         "Impossible choice",
		RequireSkill: "NonExistentSkill",
	}
	if resolver.CanChoose(choiceImpossible, party) {
		t.Error("should not be able to choose when party lacks the skill")
	}
}

func TestResolverGetSkillBonus(t *testing.T) {
	resolver := NewResolver()

	// Create party with a medic
	gen := crew.NewGenerator(12345, engine.GenreFantasy)
	party := crew.NewParty(engine.GenreFantasy, 4)
	member := gen.Generate()
	member.Skill = crew.SkillMedic
	party.Add(member)

	// Hardship events benefit from medic
	hardshipEvent := NewEvent(1, CategoryHardship, "Test", "Test", engine.GenreFantasy)
	bonus := resolver.GetSkillBonus(hardshipEvent, party)
	if bonus != 0.2 {
		t.Errorf("medic bonus for hardship = %f, want 0.2", bonus)
	}

	// Weather events don't benefit from medic
	weatherEvent := NewEvent(2, CategoryWeather, "Test", "Test", engine.GenreFantasy)
	noBonus := resolver.GetSkillBonus(weatherEvent, party)
	if noBonus != 0 {
		t.Errorf("medic bonus for weather = %f, want 0", noBonus)
	}
}

func TestResolverModifyOutcome(t *testing.T) {
	resolver := NewResolver()

	outcome := &EventOutcome{
		CrewDamage:    10,
		VesselDamage:  20,
		MoraleDelta:   -15,
		FoodDelta:     50,
		WaterDelta:    25,
		CurrencyDelta: 100,
	}

	bonus := 0.2 // 20% bonus
	modified := resolver.ModifyOutcome(outcome, bonus)

	// Negative effects should be reduced
	expectedCrewDamage := 10 * (1 - 0.2) // 8
	if modified.CrewDamage != expectedCrewDamage {
		t.Errorf("CrewDamage = %f, want %f", modified.CrewDamage, expectedCrewDamage)
	}

	expectedVesselDamage := 20 * (1 - 0.2) // 16
	if modified.VesselDamage != expectedVesselDamage {
		t.Errorf("VesselDamage = %f, want %f", modified.VesselDamage, expectedVesselDamage)
	}

	expectedMorale := -15 * (1 - 0.2) // -12
	if modified.MoraleDelta != expectedMorale {
		t.Errorf("MoraleDelta = %f, want %f", modified.MoraleDelta, expectedMorale)
	}

	// Positive effects should be increased
	expectedFood := 50 * (1 + 0.2) // 60
	if modified.FoodDelta != expectedFood {
		t.Errorf("FoodDelta = %f, want %f", modified.FoodDelta, expectedFood)
	}

	expectedWater := 25 * (1 + 0.2) // 30
	if modified.WaterDelta != expectedWater {
		t.Errorf("WaterDelta = %f, want %f", modified.WaterDelta, expectedWater)
	}

	expectedCurrency := 100 * (1 + 0.2) // 120
	if modified.CurrencyDelta != expectedCurrency {
		t.Errorf("CurrencyDelta = %f, want %f", modified.CurrencyDelta, expectedCurrency)
	}
}

func TestResolverReduceFunctions(t *testing.T) {
	resolver := NewResolver()

	// reduceNegativeEffect only reduces positive damage values
	damage := 10.0
	reducedDamage := resolver.reduceNegativeEffect(damage, 0.2)
	if reducedDamage != 8.0 {
		t.Errorf("reduced damage = %f, want 8.0", reducedDamage)
	}

	// Zero damage unchanged
	zeroDamage := resolver.reduceNegativeEffect(0, 0.2)
	if zeroDamage != 0 {
		t.Errorf("zero damage should stay 0, got %f", zeroDamage)
	}

	// reduceNegativeDelta only reduces negative deltas
	negativeDelta := -10.0
	reducedDelta := resolver.reduceNegativeDelta(negativeDelta, 0.2)
	if reducedDelta != -8.0 {
		t.Errorf("reduced delta = %f, want -8.0", reducedDelta)
	}

	// Positive delta unchanged
	positiveDelta := resolver.reduceNegativeDelta(10.0, 0.2)
	if positiveDelta != 10.0 {
		t.Errorf("positive delta should stay 10, got %f", positiveDelta)
	}

	// increasePositiveDelta only increases positive deltas
	positiveIncrease := resolver.increasePositiveDelta(10.0, 0.2)
	if positiveIncrease != 12.0 {
		t.Errorf("increased positive = %f, want 12.0", positiveIncrease)
	}

	// Negative delta unchanged
	negativeUnchanged := resolver.increasePositiveDelta(-10.0, 0.2)
	if negativeUnchanged != -10.0 {
		t.Errorf("negative delta should stay -10, got %f", negativeUnchanged)
	}
}
