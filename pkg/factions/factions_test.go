package factions

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewFaction(t *testing.T) {
	f := NewFaction(1, "Test Faction", IdeologyMerchant, engine.GenreFantasy)

	if f.ID != 1 {
		t.Errorf("expected ID 1, got %d", f.ID)
	}
	if f.Name != "Test Faction" {
		t.Errorf("expected name 'Test Faction', got %s", f.Name)
	}
	if f.Ideology != IdeologyMerchant {
		t.Errorf("expected ideology Merchant, got %v", f.Ideology)
	}
	if f.Genre != engine.GenreFantasy {
		t.Errorf("expected genre Fantasy, got %v", f.Genre)
	}
}

func TestFactionSetGenre(t *testing.T) {
	f := NewFaction(1, "Test", IdeologyMerchant, engine.GenreFantasy)
	f.SetGenre(engine.GenreScifi)

	if f.Genre != engine.GenreScifi {
		t.Errorf("expected genre Scifi, got %v", f.Genre)
	}
}

func TestFactionTerritory(t *testing.T) {
	f := NewFaction(1, "Test", IdeologyMerchant, engine.GenreFantasy)
	f.AddTerritory(10, 10, 5)

	if !f.ControlsPosition(10, 10) {
		t.Error("faction should control center of territory")
	}
	if !f.ControlsPosition(12, 10) {
		t.Error("faction should control position within radius")
	}
	if f.ControlsPosition(20, 20) {
		t.Error("faction should not control distant position")
	}
}

func TestFactionRelations(t *testing.T) {
	f := NewFaction(1, "Test", IdeologyMerchant, engine.GenreFantasy)

	// Default should be neutral
	if f.GetRelation(2) != RelationshipNeutral {
		t.Error("default relation should be neutral")
	}

	f.SetRelation(2, RelationshipHostile)
	if f.GetRelation(2) != RelationshipHostile {
		t.Error("relation should be hostile after setting")
	}
}

func TestFactionManager(t *testing.T) {
	m := NewFactionManager(engine.GenreFantasy)

	f1 := NewFaction(1, "Faction One", IdeologyMerchant, engine.GenreFantasy)
	f2 := NewFaction(2, "Faction Two", IdeologyMilitary, engine.GenreFantasy)

	m.AddFaction(f1)
	m.AddFaction(f2)

	if m.GetFaction(1) != f1 {
		t.Error("should retrieve faction 1")
	}
	if m.GetFaction(2) != f2 {
		t.Error("should retrieve faction 2")
	}
}

func TestFactionManagerReputation(t *testing.T) {
	m := NewFactionManager(engine.GenreFantasy)
	f := NewFaction(1, "Test", IdeologyMerchant, engine.GenreFantasy)
	m.AddFaction(f)

	if m.GetReputation(1) != 0 {
		t.Error("initial reputation should be 0")
	}

	m.AdjustReputation(1, 30)
	if m.GetReputation(1) != 30 {
		t.Errorf("reputation should be 30, got %d", m.GetReputation(1))
	}

	// Test clamping at max
	m.AdjustReputation(1, 100)
	if m.GetReputation(1) != 100 {
		t.Errorf("reputation should be clamped at 100, got %d", m.GetReputation(1))
	}

	// Test clamping at min
	m.AdjustReputation(1, -250)
	if m.GetReputation(1) != -100 {
		t.Errorf("reputation should be clamped at -100, got %d", m.GetReputation(1))
	}
}

func TestReputationToRelationship(t *testing.T) {
	tests := []struct {
		rep      int
		expected Relationship
	}{
		{100, RelationshipAllied},
		{75, RelationshipAllied},
		{50, RelationshipFriendly},
		{25, RelationshipFriendly},
		{0, RelationshipNeutral},
		{-25, RelationshipNeutral},
		{-50, RelationshipSuspicious},
		{-75, RelationshipSuspicious},
		{-100, RelationshipHostile},
	}

	for _, tc := range tests {
		result := ReputationToRelationship(tc.rep)
		if result != tc.expected {
			t.Errorf("rep %d: expected %v, got %v", tc.rep, tc.expected, result)
		}
	}
}

func TestFactionManagerSetGenre(t *testing.T) {
	m := NewFactionManager(engine.GenreFantasy)
	f := NewFaction(1, "Test", IdeologyMerchant, engine.GenreFantasy)
	m.AddFaction(f)

	m.SetGenre(engine.GenreScifi)

	if f.Genre != engine.GenreScifi {
		t.Error("faction genre should be updated when manager genre changes")
	}
}

func TestGenerator(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)

	manager := g.GenerateFactions(100, 100)

	factions := manager.AllFactions()
	if len(factions) < 4 || len(factions) > 6 {
		t.Errorf("expected 4-6 factions, got %d", len(factions))
	}

	// Check all factions have names and territories
	for _, f := range factions {
		if f.Name == "" {
			t.Error("faction should have a name")
		}
		if f.Description == "" {
			t.Error("faction should have a description")
		}
	}
}

func TestGeneratorDeterminism(t *testing.T) {
	g1 := NewGenerator(12345, engine.GenreFantasy)
	g2 := NewGenerator(12345, engine.GenreFantasy)

	m1 := g1.GenerateFactions(100, 100)
	m2 := g2.GenerateFactions(100, 100)

	factions1 := m1.AllFactions()
	factions2 := m2.AllFactions()

	if len(factions1) != len(factions2) {
		t.Error("same seed should produce same number of factions")
	}
}

func TestGeneratorSetGenre(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)
	g.SetGenre(engine.GenreCyberpunk)

	manager := g.GenerateFactions(100, 100)
	factions := manager.AllFactions()

	for _, f := range factions {
		if f.Genre != engine.GenreCyberpunk {
			t.Error("factions should have cyberpunk genre")
		}
	}
}

func TestAllGenreFactionGeneration(t *testing.T) {
	genres := engine.AllGenres()

	for _, genre := range genres {
		g := NewGenerator(12345, genre)
		manager := g.GenerateFactions(100, 100)

		factions := manager.AllFactions()
		if len(factions) < 4 {
			t.Errorf("genre %s: should generate at least 4 factions", genre)
		}

		for _, f := range factions {
			if f.Name == "" {
				t.Errorf("genre %s: faction should have name", genre)
			}
			if f.Description == "" {
				t.Errorf("genre %s: faction should have description", genre)
			}
			name := f.IdeologyDisplayName()
			if name == "" {
				t.Errorf("genre %s: ideology should have display name", genre)
			}
		}
	}
}

func TestRelationshipNames(t *testing.T) {
	for _, r := range AllRelationships() {
		name := RelationshipName(r)
		if name == "" || name == "Unknown" {
			t.Errorf("relationship %v should have a name", r)
		}
	}
}

func TestIdeologyNames(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		for _, ideology := range AllIdeologies() {
			name := IdeologyName(ideology, genre)
			if name == "" {
				t.Errorf("ideology %v genre %s should have name", ideology, genre)
			}
		}
	}
}

func TestFactionManagerGetFactionsAt(t *testing.T) {
	m := NewFactionManager(engine.GenreFantasy)

	f1 := NewFaction(1, "Faction One", IdeologyMerchant, engine.GenreFantasy)
	f1.AddTerritory(10, 10, 5)

	f2 := NewFaction(2, "Faction Two", IdeologyMilitary, engine.GenreFantasy)
	f2.AddTerritory(50, 50, 5)

	m.AddFaction(f1)
	m.AddFaction(f2)

	factions := m.GetFactionsAt(10, 10)
	if len(factions) != 1 || factions[0].ID != 1 {
		t.Error("should find faction 1 at position 10,10")
	}

	factions = m.GetFactionsAt(50, 50)
	if len(factions) != 1 || factions[0].ID != 2 {
		t.Error("should find faction 2 at position 50,50")
	}

	factions = m.GetFactionsAt(100, 100)
	if len(factions) != 0 {
		t.Error("should find no factions at distant position")
	}
}

func TestAreOpposingIdeologies(t *testing.T) {
	if !areOpposingIdeologies(IdeologyMerchant, IdeologyCriminal) {
		t.Error("merchant and criminal should be opposing")
	}
	if !areOpposingIdeologies(IdeologyCriminal, IdeologyMerchant) {
		t.Error("criminal and merchant should be opposing (reverse)")
	}
	if !areOpposingIdeologies(IdeologyMilitary, IdeologyCriminal) {
		t.Error("military and criminal should be opposing")
	}
	if !areOpposingIdeologies(IdeologyReligious, IdeologyScientific) {
		t.Error("religious and scientific should be opposing")
	}
	if areOpposingIdeologies(IdeologyMerchant, IdeologyMilitary) {
		t.Error("merchant and military should not be opposing")
	}
}
