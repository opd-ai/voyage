package narrative

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewStoryArc(t *testing.T) {
	arc := NewStoryArc(engine.GenreFantasy)

	if arc.Genre != engine.GenreFantasy {
		t.Errorf("expected genre Fantasy, got %v", arc.Genre)
	}
	if arc.CurrentAct != ActDeparture {
		t.Error("initial act should be Departure")
	}
	if len(arc.Beats) != 0 {
		t.Error("initial beats should be empty")
	}
}

func TestStoryArcSetGenre(t *testing.T) {
	arc := NewStoryArc(engine.GenreFantasy)
	arc.RecurringNPC = NewRecurringNPC("Test", RoleFriend, "Desc", engine.GenreFantasy)

	arc.SetGenre(engine.GenreScifi)

	if arc.Genre != engine.GenreScifi {
		t.Error("arc genre should be updated")
	}
	if arc.RecurringNPC.Genre != engine.GenreScifi {
		t.Error("NPC genre should be updated")
	}
}

func TestStoryBeatTrigger(t *testing.T) {
	arc := NewStoryArc(engine.GenreFantasy)
	arc.AddBeat(&StoryBeat{Act: ActDeparture, Title: "Test"})

	if !arc.TriggerBeat(ActDeparture, 1) {
		t.Error("should successfully trigger beat")
	}

	beat := arc.GetBeatForAct(ActDeparture)
	if !beat.Triggered {
		t.Error("beat should be triggered")
	}
	if beat.TriggerTurn != 1 {
		t.Errorf("trigger turn should be 1, got %d", beat.TriggerTurn)
	}

	// Can't trigger again
	if arc.TriggerBeat(ActDeparture, 2) {
		t.Error("should not trigger already triggered beat")
	}
}

func TestStoryArcAdvance(t *testing.T) {
	arc := NewStoryArc(engine.GenreFantasy)

	if arc.CurrentAct != ActDeparture {
		t.Error("should start at departure")
	}

	if !arc.AdvanceAct() {
		t.Error("should advance from departure")
	}
	if arc.CurrentAct != ActMidJourney {
		t.Error("should be at mid-journey")
	}

	if !arc.AdvanceAct() {
		t.Error("should advance from mid-journey")
	}
	if arc.CurrentAct != ActArrival {
		t.Error("should be at arrival")
	}

	if arc.AdvanceAct() {
		t.Error("should not advance past arrival")
	}
}

func TestStoryArcComplete(t *testing.T) {
	arc := NewStoryArc(engine.GenreFantasy)
	arc.AddBeat(&StoryBeat{Act: ActDeparture})
	arc.AddBeat(&StoryBeat{Act: ActMidJourney})
	arc.AddBeat(&StoryBeat{Act: ActArrival})

	if arc.IsComplete() {
		t.Error("should not be complete initially")
	}

	arc.TriggerBeat(ActDeparture, 1)
	arc.TriggerBeat(ActMidJourney, 10)

	if arc.IsComplete() {
		t.Error("should not be complete without all beats")
	}

	arc.TriggerBeat(ActArrival, 20)

	if !arc.IsComplete() {
		t.Error("should be complete with all beats triggered")
	}
}

func TestRecurringNPC(t *testing.T) {
	npc := NewRecurringNPC("Test NPC", RoleFriend, "A friend", engine.GenreFantasy)

	if npc.Name != "Test NPC" {
		t.Error("name should match")
	}
	if npc.Role != RoleFriend {
		t.Error("role should be friend")
	}
	if npc.Appearances != 0 {
		t.Error("initial appearances should be 0")
	}

	npc.RecordAppearance()
	if npc.Appearances != 1 {
		t.Error("appearances should increment")
	}

	npc.AddDialogue("Hello!")
	if len(npc.Dialogues) != 1 {
		t.Error("dialogue should be added")
	}
}

func TestCrewBackstory(t *testing.T) {
	bs := NewCrewBackstory(1, "Alice", "Hook", "Full", "Link")

	if bs.CrewID != 1 {
		t.Error("crew ID should match")
	}
	if bs.CrewName != "Alice" {
		t.Error("crew name should match")
	}
	if bs.Revealed {
		t.Error("should not be revealed initially")
	}

	bs.Reveal()
	if !bs.Revealed {
		t.Error("should be revealed after Reveal()")
	}
}

func TestGenerator(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)

	crewNames := []string{"Alice", "Bob", "Charlie"}
	arc := g.GenerateStoryArc(crewNames, "Legendary City")

	// Check three-act structure
	if len(arc.Beats) != 3 {
		t.Errorf("expected 3 beats, got %d", len(arc.Beats))
	}

	departure := arc.GetBeatForAct(ActDeparture)
	if departure == nil || departure.Title == "" {
		t.Error("departure beat should exist and have title")
	}

	midJourney := arc.GetBeatForAct(ActMidJourney)
	if midJourney == nil || midJourney.Title == "" {
		t.Error("mid-journey beat should exist and have title")
	}

	arrival := arc.GetBeatForAct(ActArrival)
	if arrival == nil || arrival.Title == "" {
		t.Error("arrival beat should exist and have title")
	}

	// Check recurring NPC
	if arc.RecurringNPC == nil {
		t.Error("recurring NPC should be generated")
	}
	if arc.RecurringNPC.Name == "" {
		t.Error("NPC should have name")
	}
	if len(arc.RecurringNPC.Dialogues) == 0 {
		t.Error("NPC should have dialogues")
	}
}

func TestGeneratorDeterminism(t *testing.T) {
	g1 := NewGenerator(12345, engine.GenreFantasy)
	g2 := NewGenerator(12345, engine.GenreFantasy)

	crewNames := []string{"Alice", "Bob"}
	arc1 := g1.GenerateStoryArc(crewNames, "Destination")
	arc2 := g2.GenerateStoryArc(crewNames, "Destination")

	// Same seed should produce same story beats
	for i := range arc1.Beats {
		if arc1.Beats[i].Title != arc2.Beats[i].Title {
			t.Error("same seed should produce same beat titles")
		}
	}

	// Same seed should produce same NPC
	if arc1.RecurringNPC.Name != arc2.RecurringNPC.Name {
		t.Error("same seed should produce same NPC name")
	}
}

func TestGeneratorSetGenre(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)
	g.SetGenre(engine.GenreCyberpunk)

	arc := g.GenerateStoryArc([]string{"V"}, "Free Zone")

	if arc.Genre != engine.GenreCyberpunk {
		t.Error("arc should have cyberpunk genre")
	}
	if arc.RecurringNPC.Genre != engine.GenreCyberpunk {
		t.Error("NPC should have cyberpunk genre")
	}
}

func TestAllGenreNarrativeGeneration(t *testing.T) {
	genres := engine.AllGenres()

	for _, genre := range genres {
		g := NewGenerator(12345, genre)
		arc := g.GenerateStoryArc([]string{"Crew1", "Crew2"}, "Destination")

		if len(arc.Beats) != 3 {
			t.Errorf("genre %s: should have 3 beats", genre)
		}

		for _, beat := range arc.Beats {
			if beat.Title == "" {
				t.Errorf("genre %s: beat should have title", genre)
			}
			if beat.Description == "" {
				t.Errorf("genre %s: beat should have description", genre)
			}
		}

		if arc.RecurringNPC == nil {
			t.Errorf("genre %s: should have recurring NPC", genre)
		}
		if arc.RecurringNPC.Description == "" {
			t.Errorf("genre %s: NPC should have description", genre)
		}
	}
}

func TestActNames(t *testing.T) {
	for _, act := range AllStoryActs() {
		name := ActName(act)
		if name == "" || name == "Unknown" {
			t.Errorf("act %v should have name", act)
		}
	}
}

func TestRoleNames(t *testing.T) {
	for _, role := range AllNPCRoles() {
		name := RoleName(role)
		if name == "" || name == "Unknown" {
			t.Errorf("role %v should have name", role)
		}
	}
}

func TestStoryArcGetCrewBackstory(t *testing.T) {
	arc := NewStoryArc(engine.GenreFantasy)
	bs1 := NewCrewBackstory(1, "Alice", "Hook1", "Full1", "Link1")
	bs2 := NewCrewBackstory(2, "Bob", "Hook2", "Full2", "Link2")

	arc.AddCrewBackstory(bs1)
	arc.AddCrewBackstory(bs2)

	found := arc.GetCrewBackstory(1)
	if found != bs1 {
		t.Error("should find backstory for crew 1")
	}

	found = arc.GetCrewBackstory(2)
	if found != bs2 {
		t.Error("should find backstory for crew 2")
	}

	found = arc.GetCrewBackstory(99)
	if found != nil {
		t.Error("should return nil for unknown crew")
	}
}

func TestGetActiveBeats(t *testing.T) {
	arc := NewStoryArc(engine.GenreFantasy)
	arc.AddBeat(&StoryBeat{Act: ActDeparture})
	arc.AddBeat(&StoryBeat{Act: ActMidJourney})
	arc.AddBeat(&StoryBeat{Act: ActArrival})

	active := arc.GetActiveBeats()
	if len(active) != 0 {
		t.Error("no beats should be active initially")
	}

	arc.TriggerBeat(ActDeparture, 1)
	active = arc.GetActiveBeats()
	if len(active) != 1 {
		t.Errorf("expected 1 active beat, got %d", len(active))
	}

	arc.TriggerBeat(ActMidJourney, 10)
	active = arc.GetActiveBeats()
	if len(active) != 2 {
		t.Errorf("expected 2 active beats, got %d", len(active))
	}
}
