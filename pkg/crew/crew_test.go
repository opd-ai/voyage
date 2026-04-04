package crew

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestCrewMember(t *testing.T) {
	member := NewCrewMember(1, "Test Person", TraitBrave, SkillMedic)

	if member.ID != 1 {
		t.Errorf("Expected ID 1, got %d", member.ID)
	}
	if member.Name != "Test Person" {
		t.Errorf("Expected name 'Test Person', got %s", member.Name)
	}
	if !member.IsAlive {
		t.Error("New member should be alive")
	}
	if member.Health != 100 {
		t.Errorf("Expected health 100, got %f", member.Health)
	}
}

func TestCrewMemberDamage(t *testing.T) {
	member := NewCrewMember(1, "Test", TraitBrave, SkillNone)

	// Take non-lethal damage
	died := member.TakeDamage(30)
	if died {
		t.Error("Should not die from 30 damage")
	}
	if member.Health != 70 {
		t.Errorf("Expected health 70, got %f", member.Health)
	}

	// Take lethal damage
	died = member.TakeDamage(100)
	if !died {
		t.Error("Should die from 100 damage")
	}
	if member.IsAlive {
		t.Error("Should be dead")
	}
	if member.Health != 0 {
		t.Error("Health should be 0 when dead")
	}
}

func TestCrewMemberHeal(t *testing.T) {
	member := NewCrewMember(1, "Test", TraitBrave, SkillNone)
	member.TakeDamage(50)

	member.Heal(30)
	if member.Health != 80 {
		t.Errorf("Expected health 80, got %f", member.Health)
	}

	// Heal past max
	member.Heal(50)
	if member.Health != member.MaxHealth {
		t.Error("Should not heal past max health")
	}
}

func TestParty(t *testing.T) {
	party := NewParty(engine.GenreFantasy, 4)

	if party.Capacity() != 4 {
		t.Errorf("Expected capacity 4, got %d", party.Capacity())
	}
	if !party.IsEmpty() {
		t.Error("New party should be empty")
	}

	// Add members
	m1 := NewCrewMember(1, "Member 1", TraitBrave, SkillMedic)
	m2 := NewCrewMember(2, "Member 2", TraitCautious, SkillScout)

	if !party.Add(m1) {
		t.Error("Should be able to add first member")
	}
	if !party.Add(m2) {
		t.Error("Should be able to add second member")
	}

	if party.Count() != 2 {
		t.Errorf("Expected count 2, got %d", party.Count())
	}
	if party.LivingCount() != 2 {
		t.Errorf("Expected living count 2, got %d", party.LivingCount())
	}

	// Test skill lookup
	if !party.HasSkill(SkillMedic) {
		t.Error("Party should have medic")
	}
	if party.HasSkill(SkillWarrior) {
		t.Error("Party should not have warrior")
	}
}

func TestPartyCapacity(t *testing.T) {
	party := NewParty(engine.GenreFantasy, 2)

	party.Add(NewCrewMember(1, "A", TraitBrave, SkillNone))
	party.Add(NewCrewMember(2, "B", TraitBrave, SkillNone))

	if !party.IsFull() {
		t.Error("Party should be full")
	}

	ok := party.Add(NewCrewMember(3, "C", TraitBrave, SkillNone))
	if ok {
		t.Error("Should not be able to add to full party")
	}
}

func TestPartyDeaths(t *testing.T) {
	party := NewParty(engine.GenreFantasy, 4)

	m1 := NewCrewMember(1, "Alive", TraitBrave, SkillNone)
	m2 := NewCrewMember(2, "Dead", TraitBrave, SkillNone)

	party.Add(m1)
	party.Add(m2)

	m2.TakeDamage(150)

	if party.Count() != 2 {
		t.Error("Count should include dead")
	}
	if party.LivingCount() != 1 {
		t.Error("Living count should be 1")
	}
	if len(party.Living()) != 1 {
		t.Error("Should have 1 living member")
	}
	if len(party.Dead()) != 1 {
		t.Error("Should have 1 dead member")
	}
}

func TestGenerator(t *testing.T) {
	gen := NewGenerator(12345, engine.GenreFantasy)

	m1 := gen.Generate()
	m2 := gen.Generate()

	if m1.ID == m2.ID {
		t.Error("Members should have unique IDs")
	}
	if m1.Name == "" || m2.Name == "" {
		t.Error("Members should have names")
	}
}

func TestGeneratorDeterminism(t *testing.T) {
	gen1 := NewGenerator(12345, engine.GenreFantasy)
	gen2 := NewGenerator(12345, engine.GenreFantasy)

	for i := 0; i < 10; i++ {
		m1 := gen1.Generate()
		m2 := gen2.Generate()

		if m1.Name != m2.Name {
			t.Errorf("Names should match: %s vs %s", m1.Name, m2.Name)
		}
		if m1.Trait != m2.Trait {
			t.Errorf("Traits should match")
		}
		if m1.Skill != m2.Skill {
			t.Errorf("Skills should match")
		}
	}
}

func TestNamesAndSkills(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		for _, trait := range AllTraits() {
			name := TraitName(trait, genre)
			if name == "" {
				t.Errorf("Empty trait name for trait %d, genre %s", trait, genre)
			}
		}
		for _, skill := range AllSkills() {
			name := SkillName(skill, genre)
			if name == "" {
				t.Errorf("Empty skill name for skill %d, genre %s", skill, genre)
			}
		}
	}
}
