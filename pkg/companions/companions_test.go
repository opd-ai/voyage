package companions

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewAbility(t *testing.T) {
	ability := NewAbility("test", "Test Ability", "A test", AbilityActive, 5, engine.GenreFantasy)

	if ability.ID != "test" {
		t.Error("ID mismatch")
	}
	if ability.MinSkillLevel != 5 {
		t.Error("min skill level mismatch")
	}
	if ability.Unlocked {
		t.Error("should not be unlocked initially")
	}
}

func TestAbilityUseCooldown(t *testing.T) {
	ability := NewAbility("test", "Test", "Desc", AbilityActive, 1, engine.GenreFantasy)
	ability.Cooldown = 3
	ability.Unlocked = true

	if !ability.IsReady() {
		t.Error("should be ready when unlocked and no cooldown")
	}

	ability.Use()

	if ability.IsReady() {
		t.Error("should not be ready after use")
	}
	if ability.CurrentCooldown != 3 {
		t.Error("cooldown should be set")
	}

	ability.Tick()
	ability.Tick()
	ability.Tick()

	if !ability.IsReady() {
		t.Error("should be ready after cooldown expires")
	}
}

func TestNewCompanion(t *testing.T) {
	companion := NewCompanion(1, "Test", "Wizard", RoleGuide, engine.GenreFantasy)

	if companion.ID != 1 {
		t.Error("ID mismatch")
	}
	if companion.Name != "Test" {
		t.Error("name mismatch")
	}
	if companion.Role != RoleGuide {
		t.Error("role mismatch")
	}
	if companion.SkillLevel != 1 {
		t.Error("should start at skill level 1")
	}
	if !companion.IsActive {
		t.Error("should be active initially")
	}
}

func TestCompanionTraits(t *testing.T) {
	companion := NewCompanion(1, "Test", "Title", RoleGuide, engine.GenreFantasy)

	companion.AddTrait(TraitBrave)
	companion.AddTrait(TraitLoyal)

	if !companion.HasTrait(TraitBrave) {
		t.Error("should have brave trait")
	}
	if !companion.HasTrait(TraitLoyal) {
		t.Error("should have loyal trait")
	}
	if companion.HasTrait(TraitCautious) {
		t.Error("should not have cautious trait")
	}
}

func TestCompanionAbilities(t *testing.T) {
	companion := NewCompanion(1, "Test", "Title", RoleGuide, engine.GenreFantasy)
	companion.SkillLevel = 5

	ability1 := NewAbility("a1", "Ability 1", "Desc", AbilityActive, 3, engine.GenreFantasy)
	ability2 := NewAbility("a2", "Ability 2", "Desc", AbilityActive, 8, engine.GenreFantasy)

	companion.AddAbility(ability1)
	companion.AddAbility(ability2)

	if !ability1.Unlocked {
		t.Error("ability 1 should be unlocked (skill 5 >= min 3)")
	}
	if ability2.Unlocked {
		t.Error("ability 2 should not be unlocked (skill 5 < min 8)")
	}

	unlocked := companion.GetUnlockedAbilities()
	if len(unlocked) != 1 {
		t.Errorf("expected 1 unlocked ability, got %d", len(unlocked))
	}
}

func TestCompanionGainExperience(t *testing.T) {
	companion := NewCompanion(1, "Test", "Title", RoleGuide, engine.GenreFantasy)
	ability := NewAbility("a1", "Ability", "Desc", AbilityActive, 3, engine.GenreFantasy)
	companion.AddAbility(ability)

	// Gain enough XP to level up to 3
	companion.GainExperience(200)

	if companion.SkillLevel != 3 {
		t.Errorf("expected skill level 3, got %d", companion.SkillLevel)
	}
	if !ability.Unlocked {
		t.Error("ability should be unlocked at skill level 3")
	}
}

func TestCompanionMoraleAndLoyalty(t *testing.T) {
	companion := NewCompanion(1, "Test", "Title", RoleGuide, engine.GenreFantasy)

	companion.AdjustMorale(0.5)
	if companion.Morale > 1.0 {
		t.Error("morale should be clamped to 1.0")
	}

	companion.AdjustMorale(-2.0)
	if companion.Morale < 0.0 {
		t.Error("morale should be clamped to 0.0")
	}

	companion.AdjustLoyalty(0.5)
	companion.AdjustLoyalty(-2.0)
	if companion.Loyalty < 0.0 {
		t.Error("loyalty should be clamped to 0.0")
	}
}

func TestCompanionRelationship(t *testing.T) {
	companion := NewCompanion(1, "Test", "Title", RoleGuide, engine.GenreFantasy)

	companion.AdjustRelationshipWithPlayer(2.0)
	if companion.RelationshipWithPlayer > 1.0 {
		t.Error("relationship should be clamped to 1.0")
	}

	companion.AdjustRelationshipWithPlayer(-3.0)
	if companion.RelationshipWithPlayer < -1.0 {
		t.Error("relationship should be clamped to -1.0")
	}
}

func TestCompanionSetGenre(t *testing.T) {
	companion := NewCompanion(1, "Test", "Title", RoleGuide, engine.GenreFantasy)
	ability := NewAbility("a1", "Ability", "Desc", AbilityActive, 1, engine.GenreFantasy)
	companion.AddAbility(ability)

	companion.SetGenre(engine.GenreScifi)

	if companion.Genre != engine.GenreScifi {
		t.Error("companion genre should be updated")
	}
	if ability.Genre != engine.GenreScifi {
		t.Error("ability genre should be updated")
	}
}

func TestCompanionEvent(t *testing.T) {
	event := NewCompanionEvent(1, 1, "Test Event", "Description", "Dialogue", engine.GenreFantasy)

	if event.Triggered {
		t.Error("should not be triggered initially")
	}

	event.RequiredTrait = TraitBrave
	event.MoraleChange = 0.1
	event.LoyaltyChange = 0.05

	companion := NewCompanion(1, "Test", "Title", RoleGuide, engine.GenreFantasy)
	companion.AddTrait(TraitBrave)

	if !event.CanTrigger(companion) {
		t.Error("event should be triggerable")
	}

	initialMorale := companion.Morale
	event.Trigger(companion)

	if !event.Triggered {
		t.Error("event should be triggered")
	}
	if companion.Morale <= initialMorale {
		t.Error("morale should increase")
	}
}

func TestCompanionEventCantTrigger(t *testing.T) {
	event := NewCompanionEvent(1, 1, "Test", "Desc", "Dialogue", engine.GenreFantasy)
	event.RequiredTrait = TraitBrave

	companion := NewCompanion(2, "Test", "Title", RoleGuide, engine.GenreFantasy)
	companion.AddTrait(TraitCautious)

	if event.CanTrigger(companion) {
		t.Error("event should not trigger for different companion ID")
	}

	companion.ID = 1
	if event.CanTrigger(companion) {
		t.Error("event should not trigger without required trait")
	}
}

func TestCompanionManager(t *testing.T) {
	manager := NewCompanionManager(engine.GenreFantasy, 5)

	c1 := NewCompanion(1, "Test1", "Title", RoleGuide, engine.GenreFantasy)
	c2 := NewCompanion(2, "Test2", "Title", RoleScout, engine.GenreFantasy)

	if !manager.AddCompanion(c1) {
		t.Error("should add companion")
	}
	if !manager.AddCompanion(c2) {
		t.Error("should add second companion")
	}

	if manager.ActiveCount() != 2 {
		t.Error("should have 2 active companions")
	}

	if manager.GetCompanion(1) != c1 {
		t.Error("should retrieve companion by ID")
	}

	if manager.GetCompanionByRole(RoleScout) != c2 {
		t.Error("should retrieve companion by role")
	}
}

func TestCompanionManagerMaxSize(t *testing.T) {
	manager := NewCompanionManager(engine.GenreFantasy, 2)

	manager.AddCompanion(NewCompanion(1, "Test1", "T", RoleGuide, engine.GenreFantasy))
	manager.AddCompanion(NewCompanion(2, "Test2", "T", RoleScout, engine.GenreFantasy))

	if manager.AddCompanion(NewCompanion(3, "Test3", "T", RoleMedic, engine.GenreFantasy)) {
		t.Error("should not exceed max size")
	}
}

func TestCompanionManagerRemove(t *testing.T) {
	manager := NewCompanionManager(engine.GenreFantasy, 5)
	c1 := NewCompanion(1, "Test1", "Title", RoleGuide, engine.GenreFantasy)
	manager.AddCompanion(c1)

	if !manager.RemoveCompanion(1) {
		t.Error("should remove companion")
	}
	if manager.RemoveCompanion(1) {
		t.Error("should not remove already removed companion")
	}
}

func TestCompanionManagerStats(t *testing.T) {
	manager := NewCompanionManager(engine.GenreFantasy, 5)

	c1 := NewCompanion(1, "Test1", "Title", RoleGuide, engine.GenreFantasy)
	c1.SkillLevel = 5
	c1.Morale = 0.8

	c2 := NewCompanion(2, "Test2", "Title", RoleScout, engine.GenreFantasy)
	c2.SkillLevel = 3
	c2.Morale = 0.6

	manager.AddCompanion(c1)
	manager.AddCompanion(c2)

	if manager.TotalSkillLevel() != 8 {
		t.Errorf("expected total skill 8, got %d", manager.TotalSkillLevel())
	}

	avgMorale := manager.AverageMorale()
	if avgMorale < 0.69 || avgMorale > 0.71 {
		t.Errorf("expected avg morale ~0.7, got %f", avgMorale)
	}
}

func TestCompanionManagerSetGenre(t *testing.T) {
	manager := NewCompanionManager(engine.GenreFantasy, 5)
	c1 := NewCompanion(1, "Test", "Title", RoleGuide, engine.GenreFantasy)
	manager.AddCompanion(c1)

	manager.SetGenre(engine.GenreScifi)

	if manager.Genre != engine.GenreScifi {
		t.Error("manager genre should be updated")
	}
	if c1.Genre != engine.GenreScifi {
		t.Error("companion genre should be updated")
	}
}

func TestRoleNameAllGenres(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		for _, role := range AllCompanionRoles() {
			name := RoleName(role, genre)
			if name == "" {
				t.Errorf("role %v genre %s should have name", role, genre)
			}
		}
	}
}

func TestTraitName(t *testing.T) {
	for _, trait := range AllPersonalityTraits() {
		name := TraitName(trait)
		if name == "Unknown" || name == "" {
			t.Errorf("trait %v should have name", trait)
		}
	}
}

func TestGenerator(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)

	companion := g.GenerateCompanion(RoleGuide)

	if companion.Name == "" {
		t.Error("companion should have name")
	}
	if companion.Backstory == "" {
		t.Error("companion should have backstory")
	}
	if len(companion.Traits) < 2 {
		t.Error("companion should have at least 2 traits")
	}
	if len(companion.Abilities) == 0 {
		t.Error("companion should have abilities")
	}
}

func TestGeneratorDeterminism(t *testing.T) {
	g1 := NewGenerator(12345, engine.GenreFantasy)
	g2 := NewGenerator(12345, engine.GenreFantasy)

	c1 := g1.GenerateCompanion(RoleGuide)
	c2 := g2.GenerateCompanion(RoleGuide)

	if c1.Name != c2.Name {
		t.Error("same seed should produce same name")
	}
	if c1.Backstory != c2.Backstory {
		t.Error("same seed should produce same backstory")
	}
}

func TestGeneratorAllGenres(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		g := NewGenerator(12345, genre)

		for _, role := range AllCompanionRoles() {
			companion := g.GenerateCompanion(role)

			if companion.Name == "" {
				t.Errorf("genre %s role %v: should have name", genre, role)
			}
			if companion.Backstory == "" {
				t.Errorf("genre %s role %v: should have backstory", genre, role)
			}
			if len(companion.Abilities) == 0 {
				t.Errorf("genre %s role %v: should have abilities", genre, role)
			}
		}
	}
}

func TestGeneratorRandomCompanion(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)

	companion := g.GenerateRandomCompanion()

	if companion.Name == "" {
		t.Error("random companion should have name")
	}
}

func TestGeneratorCompanionEvent(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)

	companion := g.GenerateCompanion(RoleGuide)
	event := g.GenerateCompanionEvent(companion)

	if event == nil {
		t.Fatal("should generate event")
	}
	if event.Title == "" {
		t.Error("event should have title")
	}
	if event.Description == "" {
		t.Error("event should have description")
	}
	if event.Dialogue == "" {
		t.Error("event should have dialogue")
	}
}

func TestGeneratorEventDialogueAllGenres(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		g := NewGenerator(12345, genre)

		companion := g.GenerateCompanion(RoleGuide)
		event := g.GenerateCompanionEvent(companion)

		if event == nil {
			t.Errorf("genre %s: should generate event", genre)
			continue
		}
		if event.Dialogue == "" {
			t.Errorf("genre %s: event should have dialogue", genre)
		}
	}
}

func TestGeneratorParty(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)

	party := g.GenerateParty(4)

	if party.ActiveCount() != 4 {
		t.Errorf("expected 4 companions, got %d", party.ActiveCount())
	}
	if len(party.Events) < 4 {
		t.Error("should generate events for companions")
	}
}

func TestGeneratorSetGenre(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)
	g.SetGenre(engine.GenreCyberpunk)

	companion := g.GenerateCompanion(RoleGuide)

	if companion.Genre != engine.GenreCyberpunk {
		t.Error("companion should have cyberpunk genre")
	}
}

func TestCompanionManagerCheckEvents(t *testing.T) {
	manager := NewCompanionManager(engine.GenreFantasy, 5)

	c1 := NewCompanion(1, "Test", "Title", RoleGuide, engine.GenreFantasy)
	c1.AddTrait(TraitBrave)
	manager.AddCompanion(c1)

	event := NewCompanionEvent(1, 1, "Test Event", "Desc", "Dialogue", engine.GenreFantasy)
	event.RequiredTrait = TraitBrave
	manager.AddEvent(event)

	triggerable := manager.CheckEvents()
	if len(triggerable) != 1 {
		t.Errorf("expected 1 triggerable event, got %d", len(triggerable))
	}
}
