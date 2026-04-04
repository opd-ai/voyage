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

func TestTraitEffects(t *testing.T) {
	// Test that all traits have effects defined
	for _, trait := range AllTraits() {
		effect := GetTraitEffect(trait)
		// Effects can be zero values, just verify we can get them
		_ = effect
	}

	// Test specific trait effects
	braveEffect := GetTraitEffect(TraitBrave)
	if braveEffect.CombatModifier <= 0 {
		t.Error("TraitBrave should have positive combat modifier")
	}

	cautiousEffect := GetTraitEffect(TraitCautious)
	if cautiousEffect.FuelModifier >= 0 {
		t.Error("TraitCautious should have negative fuel modifier (saves fuel)")
	}

	navigatorEffect := GetTraitEffect(TraitNavigator)
	if navigatorEffect.TravelModifier <= 0 {
		t.Error("TraitNavigator should have positive travel modifier")
	}
	if navigatorEffect.FuelModifier >= 0 {
		t.Error("TraitNavigator should have negative fuel modifier (saves fuel)")
	}

	scavengerEffect := GetTraitEffect(TraitScavenger)
	if scavengerEffect.ScavengeModifier <= 0 {
		t.Error("TraitScavenger should have positive scavenge modifier")
	}
}

func TestAllTraitsCount(t *testing.T) {
	traits := AllTraits()
	if len(traits) != 10 {
		t.Errorf("Expected 10 traits, got %d", len(traits))
	}
}

func TestNavigatorAndScavengerTraits(t *testing.T) {
	// Test that navigator and scavenger exist
	found := make(map[Trait]bool)
	for _, trait := range AllTraits() {
		found[trait] = true
	}
	if !found[TraitNavigator] {
		t.Error("TraitNavigator should be in AllTraits()")
	}
	if !found[TraitScavenger] {
		t.Error("TraitScavenger should be in AllTraits()")
	}

	// Test trait names for all genres
	for _, genre := range engine.AllGenres() {
		navName := TraitName(TraitNavigator, genre)
		if navName == "" {
			t.Errorf("TraitNavigator should have name for genre %s", genre)
		}
		scavName := TraitName(TraitScavenger, genre)
		if scavName == "" {
			t.Errorf("TraitScavenger should have name for genre %s", genre)
		}
	}
}

func TestSkillEffectiveness(t *testing.T) {
	// Unskilled worker
	unskilled := NewCrewMember(1, "Test", TraitBrave, SkillNone)
	if unskilled.SkillEffectiveness() != 0.5 {
		t.Errorf("Unskilled effectiveness = %f, want 0.5", unskilled.SkillEffectiveness())
	}

	// New skilled worker (level 0)
	skilled := NewCrewMember(2, "Test", TraitBrave, SkillMedic)
	if skilled.SkillEffectiveness() != 1.0 {
		t.Errorf("Level 0 effectiveness = %f, want 1.0", skilled.SkillEffectiveness())
	}

	// Level up and test
	skilled.SkillLevel = 3
	if skilled.SkillEffectiveness() != 1.3 {
		t.Errorf("Level 3 effectiveness = %f, want 1.3", skilled.SkillEffectiveness())
	}

	// Max level
	skilled.SkillLevel = MaxSkillLevel
	if skilled.SkillEffectiveness() != 1.5 {
		t.Errorf("Level 5 effectiveness = %f, want 1.5", skilled.SkillEffectiveness())
	}
}

func TestSkillExpGain(t *testing.T) {
	member := NewCrewMember(1, "Test", TraitBrave, SkillMedic)
	if member.SkillLevel != 0 {
		t.Error("Should start at level 0")
	}

	// Partial experience
	leveledUp := member.GainSkillExp(50)
	if leveledUp {
		t.Error("Should not level up from 50 exp")
	}
	if member.SkillLevel != 0 {
		t.Error("Should still be level 0")
	}

	// Level up
	leveledUp = member.GainSkillExp(60)
	if !leveledUp {
		t.Error("Should level up after 110 total exp (threshold 100)")
	}
	if member.SkillLevel != 1 {
		t.Errorf("Should be level 1, got %d", member.SkillLevel)
	}
	// Should have overflow exp
	if member.SkillExp != 10 {
		t.Errorf("Overflow exp = %f, want 10", member.SkillExp)
	}
}

func TestSkillExpThresholds(t *testing.T) {
	expected := []float64{100, 150, 225, 337.5, 506.25}
	for level, exp := range expected {
		got := SkillExpThreshold(level)
		if got != exp {
			t.Errorf("Level %d threshold = %f, want %f", level, got, exp)
		}
	}
}

func TestSkillExpNoGainForUnskilled(t *testing.T) {
	member := NewCrewMember(1, "Test", TraitBrave, SkillNone)
	leveledUp := member.GainSkillExp(1000)
	if leveledUp {
		t.Error("Unskilled members should not gain exp")
	}
	if member.SkillLevel != 0 {
		t.Error("Unskilled members should stay at level 0")
	}
}

func TestSkillExpProgress(t *testing.T) {
	member := NewCrewMember(1, "Test", TraitBrave, SkillMedic)
	member.SkillExp = 50
	progress := member.SkillExpProgress()
	if progress != 0.5 {
		t.Errorf("Progress = %f, want 0.5", progress)
	}
}

func TestSkillLevelName(t *testing.T) {
	names := []string{"Novice", "Apprentice", "Journeyman", "Expert", "Master", "Grandmaster"}
	for level, expectedName := range names {
		name := SkillLevelName(level)
		if name != expectedName {
			t.Errorf("Level %d name = %s, want %s", level, name, expectedName)
		}
	}
}

func TestRelationshipNetwork(t *testing.T) {
	network := NewRelationshipNetwork(engine.GenreFantasy)

	// Test getting non-existent relationship
	rel := network.GetRelationship(1, 2)
	if rel.Type != RelationNeutral {
		t.Error("New relationship should be neutral")
	}
	if rel.Strength != 0 {
		t.Error("New relationship should have 0 strength")
	}
}

func TestRelationshipInteraction(t *testing.T) {
	network := NewRelationshipNetwork(engine.GenreFantasy)

	// Positive interactions
	for i := 0; i < 10; i++ {
		network.Interact(1, 2, 8)
	}
	rel := network.GetRelationship(1, 2)
	if rel.Strength != 80 {
		t.Errorf("Strength = %f, want 80", rel.Strength)
	}
	if rel.Interactions != 10 {
		t.Errorf("Interactions = %d, want 10", rel.Interactions)
	}
	if rel.Type != RelationRomantic {
		t.Errorf("Type = %d, want %d (Romantic)", rel.Type, RelationRomantic)
	}
}

func TestRelationshipRivalry(t *testing.T) {
	network := NewRelationshipNetwork(engine.GenreFantasy)

	// Negative interactions
	for i := 0; i < 10; i++ {
		network.Interact(1, 2, -10)
	}
	rel := network.GetRelationship(1, 2)
	if rel.Strength != -100 {
		t.Errorf("Strength = %f, want -100", rel.Strength)
	}
	if rel.Type != RelationRivalry {
		t.Errorf("Type = %d, want %d (Rivalry)", rel.Type, RelationRivalry)
	}
}

func TestRelationshipMoraleModifier(t *testing.T) {
	network := NewRelationshipNetwork(engine.GenreFantasy)

	// Create friendly relationship
	for i := 0; i < 5; i++ {
		network.Interact(1, 2, 15)
	}
	rel := network.GetRelationship(1, 2)
	modifier := rel.MoraleModifier()
	if modifier <= 0 {
		t.Errorf("Friendly relationship should have positive morale modifier, got %f", modifier)
	}

	// Create rivalry
	network2 := NewRelationshipNetwork(engine.GenreFantasy)
	for i := 0; i < 10; i++ {
		network2.Interact(3, 4, -10)
	}
	rel2 := network2.GetRelationship(3, 4)
	modifier2 := rel2.MoraleModifier()
	if modifier2 >= 0 {
		t.Errorf("Rivalry should have negative morale modifier, got %f", modifier2)
	}
}

func TestRelationTypeName(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		types := []RelationType{RelationNeutral, RelationFriendly, RelationRomantic, RelationRivalry, RelationMentorship}
		for _, rt := range types {
			name := RelationTypeName(rt, genre)
			if name == "" || name == "Unknown" {
				t.Errorf("RelationType %d should have name for genre %s", rt, genre)
			}
		}
	}
}

func TestRelationshipsFor(t *testing.T) {
	network := NewRelationshipNetwork(engine.GenreFantasy)

	network.Interact(1, 2, 10)
	network.Interact(1, 3, 10)
	network.Interact(2, 3, 10)

	rels := network.RelationshipsFor(1)
	if len(rels) != 2 {
		t.Errorf("Member 1 should have 2 relationships, got %d", len(rels))
	}
}

func TestAllRelationships(t *testing.T) {
	network := NewRelationshipNetwork(engine.GenreFantasy)

	network.Interact(1, 2, 10)
	network.Interact(3, 4, -10)

	all := network.AllRelationships()
	if len(all) != 2 {
		t.Errorf("Should have 2 relationships, got %d", len(all))
	}
}
