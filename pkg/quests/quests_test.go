package quests

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewQuest(t *testing.T) {
	q := NewQuest(1, TypeDelivery, "Test Quest", "Test description", engine.GenreFantasy)

	if q.ID != 1 {
		t.Errorf("expected ID 1, got %d", q.ID)
	}
	if q.Type != TypeDelivery {
		t.Errorf("expected type Delivery, got %v", q.Type)
	}
	if q.Title != "Test Quest" {
		t.Errorf("expected title 'Test Quest', got %s", q.Title)
	}
	if q.Status != StatusAvailable {
		t.Error("new quest should be available")
	}
}

func TestQuestSetGenre(t *testing.T) {
	q := NewQuest(1, TypeDelivery, "Test", "Desc", engine.GenreFantasy)
	q.SetGenre(engine.GenreScifi)

	if q.Genre != engine.GenreScifi {
		t.Errorf("expected genre Scifi, got %v", q.Genre)
	}
}

func TestQuestObjectives(t *testing.T) {
	q := NewQuest(1, TypeDelivery, "Test", "Desc", engine.GenreFantasy)
	q.AddObjective("Deliver package", 10, 20, "Town")

	if len(q.Objectives) != 1 {
		t.Errorf("expected 1 objective, got %d", len(q.Objectives))
	}
	if q.Objectives[0].TargetX != 10 || q.Objectives[0].TargetY != 20 {
		t.Error("objective position incorrect")
	}
	if q.IsComplete() {
		t.Error("quest should not be complete with uncompleted objective")
	}

	q.CompleteObjective(0)
	if !q.Objectives[0].Completed {
		t.Error("objective should be completed")
	}
	if !q.IsComplete() {
		t.Error("quest should be complete when all objectives done")
	}
}

func TestQuestStatusTransitions(t *testing.T) {
	q := NewQuest(1, TypeDelivery, "Test", "Desc", engine.GenreFantasy)

	if q.Status != StatusAvailable {
		t.Error("initial status should be available")
	}

	q.Accept()
	if q.Status != StatusActive {
		t.Error("status should be active after accepting")
	}

	// Can't accept again
	q.Accept()
	if q.Status != StatusActive {
		t.Error("status should remain active")
	}

	q.Complete()
	if q.Status != StatusCompleted {
		t.Error("status should be completed")
	}
}

func TestQuestDecline(t *testing.T) {
	q := NewQuest(1, TypeDelivery, "Test", "Desc", engine.GenreFantasy)
	q.Decline()

	if q.Status != StatusDeclined {
		t.Error("status should be declined")
	}
}

func TestQuestFail(t *testing.T) {
	q := NewQuest(1, TypeDelivery, "Test", "Desc", engine.GenreFantasy)
	q.Accept()
	q.Fail()

	if q.Status != StatusFailed {
		t.Error("status should be failed")
	}
}

func TestQuestTimeLimit(t *testing.T) {
	q := NewQuest(1, TypeDelivery, "Test", "Desc", engine.GenreFantasy)
	q.TimeLimit = 10
	q.Accept()

	q.AdvanceTime(5)
	if q.TimeLimit != 5 {
		t.Errorf("expected time limit 5, got %d", q.TimeLimit)
	}
	if q.Status != StatusActive {
		t.Error("quest should still be active")
	}

	q.AdvanceTime(6)
	if q.Status != StatusFailed {
		t.Error("quest should fail when time expires")
	}
	if q.TimeLimit != 0 {
		t.Error("time limit should be 0 when expired")
	}
}

func TestQuestTracker(t *testing.T) {
	tracker := NewQuestTracker(engine.GenreFantasy)

	q1 := NewQuest(1, TypeDelivery, "Quest 1", "Desc", engine.GenreFantasy)
	q2 := NewQuest(2, TypeRescue, "Quest 2", "Desc", engine.GenreFantasy)

	tracker.AddQuest(q1)
	tracker.AddQuest(q2)

	if tracker.GetQuest(1) != q1 {
		t.Error("should retrieve quest 1")
	}
	if tracker.GetQuest(2) != q2 {
		t.Error("should retrieve quest 2")
	}

	available := tracker.AvailableQuests()
	if len(available) != 2 {
		t.Errorf("expected 2 available quests, got %d", len(available))
	}
}

func TestQuestTrackerAccept(t *testing.T) {
	tracker := NewQuestTracker(engine.GenreFantasy)
	tracker.ActiveLimit = 2

	q1 := NewQuest(1, TypeDelivery, "Quest 1", "Desc", engine.GenreFantasy)
	q2 := NewQuest(2, TypeRescue, "Quest 2", "Desc", engine.GenreFantasy)
	q3 := NewQuest(3, TypeExplore, "Quest 3", "Desc", engine.GenreFantasy)

	tracker.AddQuest(q1)
	tracker.AddQuest(q2)
	tracker.AddQuest(q3)

	if !tracker.AcceptQuest(1) {
		t.Error("should accept quest 1")
	}
	if !tracker.AcceptQuest(2) {
		t.Error("should accept quest 2")
	}
	if tracker.AcceptQuest(3) {
		t.Error("should not accept quest 3 (over limit)")
	}

	if tracker.ActiveQuestCount() != 2 {
		t.Errorf("expected 2 active quests, got %d", tracker.ActiveQuestCount())
	}
}

func TestQuestTrackerSetGenre(t *testing.T) {
	tracker := NewQuestTracker(engine.GenreFantasy)
	q := NewQuest(1, TypeDelivery, "Test", "Desc", engine.GenreFantasy)
	tracker.AddQuest(q)

	tracker.SetGenre(engine.GenreCyberpunk)

	if q.Genre != engine.GenreCyberpunk {
		t.Error("quest genre should be updated")
	}
}

func TestQuestTrackerCheckObjectives(t *testing.T) {
	tracker := NewQuestTracker(engine.GenreFantasy)
	q := NewQuest(1, TypeDelivery, "Test", "Desc", engine.GenreFantasy)
	q.AddObjective("Deliver", 10, 20, "Town")
	q.Accept()
	tracker.AddQuest(q)

	matches := tracker.CheckObjectivesAt(10, 20)
	if len(matches) != 1 {
		t.Error("should find quest with objective at position")
	}

	matches = tracker.CheckObjectivesAt(50, 50)
	if len(matches) != 0 {
		t.Error("should not find quests at wrong position")
	}
}

func TestGenerator(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)

	quests := g.GenerateQuestBoard(50, 50, 100, 100)

	if len(quests) < 2 || len(quests) > 4 {
		t.Errorf("expected 2-4 quests, got %d", len(quests))
	}

	for _, q := range quests {
		if q.Title == "" {
			t.Error("quest should have title")
		}
		if q.Description == "" {
			t.Error("quest should have description")
		}
		if q.GiverName == "" {
			t.Error("quest should have giver name")
		}
		if len(q.Objectives) == 0 {
			t.Error("quest should have objectives")
		}
		if q.Reward.Currency == 0 && q.Reward.Morale == 0 {
			t.Error("quest should have some reward")
		}
	}
}

func TestGeneratorDeterminism(t *testing.T) {
	g1 := NewGenerator(12345, engine.GenreFantasy)
	g2 := NewGenerator(12345, engine.GenreFantasy)

	q1 := g1.GenerateQuestBoard(50, 50, 100, 100)
	q2 := g2.GenerateQuestBoard(50, 50, 100, 100)

	if len(q1) != len(q2) {
		t.Error("same seed should produce same number of quests")
	}

	for i := range q1 {
		if q1[i].Title != q2[i].Title {
			t.Error("same seed should produce same quest titles")
		}
	}
}

func TestGeneratorSetGenre(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)
	g.SetGenre(engine.GenreScifi)

	quests := g.GenerateQuestBoard(50, 50, 100, 100)

	for _, q := range quests {
		if q.Genre != engine.GenreScifi {
			t.Error("quests should have scifi genre")
		}
	}
}

func TestAllGenreQuestGeneration(t *testing.T) {
	genres := engine.AllGenres()

	for _, genre := range genres {
		g := NewGenerator(12345, genre)
		quests := g.GenerateQuestBoard(50, 50, 100, 100)

		if len(quests) < 2 {
			t.Errorf("genre %s: should generate at least 2 quests", genre)
		}

		for _, q := range quests {
			if q.Title == "" {
				t.Errorf("genre %s: quest should have title", genre)
			}
			typeName := q.TypeDisplayName()
			if typeName == "" {
				t.Errorf("genre %s: quest type should have display name", genre)
			}
		}
	}
}

func TestPrimaryObjective(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)

	primary := g.GeneratePrimaryObjective(100, 100, "Legendary City")

	if primary.ID != 0 {
		t.Error("primary objective should have ID 0")
	}
	if primary.Status != StatusActive {
		t.Error("primary objective should be active")
	}
	if primary.TimeLimit != 0 {
		t.Error("primary objective should have no time limit")
	}
	if len(primary.Objectives) != 1 {
		t.Error("primary objective should have one objective")
	}
}

func TestQuestTypeNames(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		for _, qType := range AllQuestTypes() {
			name := QuestTypeName(qType, genre)
			if name == "" {
				t.Errorf("quest type %v genre %s should have name", qType, genre)
			}
		}
	}
}

func TestStatusNames(t *testing.T) {
	statuses := []QuestStatus{
		StatusAvailable,
		StatusActive,
		StatusCompleted,
		StatusFailed,
		StatusDeclined,
	}

	for _, s := range statuses {
		name := StatusName(s)
		if name == "" || name == "Unknown" {
			t.Errorf("status %v should have a name", s)
		}
	}
}

func TestQuestTrackerAdvanceAllTime(t *testing.T) {
	tracker := NewQuestTracker(engine.GenreFantasy)

	q1 := NewQuest(1, TypeDelivery, "Quest 1", "Desc", engine.GenreFantasy)
	q1.TimeLimit = 10
	q1.Accept()

	q2 := NewQuest(2, TypeRescue, "Quest 2", "Desc", engine.GenreFantasy)
	q2.TimeLimit = 5
	q2.Accept()

	tracker.AddQuest(q1)
	tracker.AddQuest(q2)

	tracker.AdvanceAllTime(3)

	if q1.TimeLimit != 7 {
		t.Errorf("q1 time limit should be 7, got %d", q1.TimeLimit)
	}
	if q2.TimeLimit != 2 {
		t.Errorf("q2 time limit should be 2, got %d", q2.TimeLimit)
	}

	tracker.AdvanceAllTime(5)

	if q2.Status != StatusFailed {
		t.Error("q2 should have failed")
	}
}
