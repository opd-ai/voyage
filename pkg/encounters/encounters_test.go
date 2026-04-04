package encounters

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/resources"
	"github.com/opd-ai/voyage/pkg/vessel"
)

func TestAllEncounterTypes(t *testing.T) {
	types := AllEncounterTypes()
	if len(types) != 5 {
		t.Errorf("Expected 5 encounter types, got %d", len(types))
	}

	expected := []EncounterType{TypeAmbush, TypeNegotiation, TypeRace, TypeCrisis, TypePuzzle}
	for i, typ := range expected {
		if types[i] != typ {
			t.Errorf("Expected type %d at index %d, got %d", typ, i, types[i])
		}
	}
}

func TestAllEncounterRoles(t *testing.T) {
	roles := AllEncounterRoles()
	if len(roles) != 4 {
		t.Errorf("Expected 4 encounter roles, got %d", len(roles))
	}
}

func TestNewEncounter(t *testing.T) {
	enc := NewEncounter(1, TypeAmbush, engine.GenreFantasy)

	if enc.ID != 1 {
		t.Errorf("Expected ID 1, got %d", enc.ID)
	}
	if enc.Type != TypeAmbush {
		t.Errorf("Expected type TypeAmbush, got %d", enc.Type)
	}
	if enc.Genre != engine.GenreFantasy {
		t.Errorf("Expected genre fantasy, got %s", enc.Genre)
	}
	if enc.State != StatePending {
		t.Errorf("Expected state Pending, got %d", enc.State)
	}
}

func TestEncounterSetGenre(t *testing.T) {
	enc := NewEncounter(1, TypeAmbush, engine.GenreFantasy)
	enc.SetGenre(engine.GenreScifi)

	if enc.Genre != engine.GenreScifi {
		t.Errorf("Expected genre scifi, got %s", enc.Genre)
	}
}

func TestEncounterAssignment(t *testing.T) {
	enc := NewEncounter(1, TypeAmbush, engine.GenreFantasy)

	enc.AssignCrew(RoleFighter, 100)

	memberID, ok := enc.GetAssignment(RoleFighter)
	if !ok {
		t.Error("Expected fighter role to be assigned")
	}
	if memberID != 100 {
		t.Errorf("Expected member ID 100, got %d", memberID)
	}

	if !enc.IsRoleFilled(RoleFighter) {
		t.Error("Expected fighter role to be filled")
	}

	if enc.IsRoleFilled(RoleMedic) {
		t.Error("Expected medic role to be empty")
	}

	enc.UnassignCrew(RoleFighter)
	if enc.IsRoleFilled(RoleFighter) {
		t.Error("Expected fighter role to be empty after unassignment")
	}
}

func TestEncounterRequiredRolesFilled(t *testing.T) {
	enc := NewEncounter(1, TypeAmbush, engine.GenreFantasy)
	enc.RequiredRoles = []EncounterRole{RoleFighter, RoleMedic}

	if enc.RequiredRolesFilled() {
		t.Error("Expected required roles not filled")
	}

	enc.AssignCrew(RoleFighter, 1)
	if enc.RequiredRolesFilled() {
		t.Error("Expected required roles still not filled")
	}

	enc.AssignCrew(RoleMedic, 2)
	if !enc.RequiredRolesFilled() {
		t.Error("Expected required roles to be filled")
	}
}

func TestEncounterStart(t *testing.T) {
	enc := NewEncounter(1, TypeAmbush, engine.GenreFantasy)
	enc.RequiredRoles = []EncounterRole{RoleFighter}

	// Should fail without required roles
	if enc.Start() {
		t.Error("Expected start to fail without required roles")
	}

	enc.AssignCrew(RoleFighter, 1)
	if !enc.Start() {
		t.Error("Expected start to succeed with required roles")
	}

	if enc.State != StateResolution {
		t.Error("Expected state to be Resolution")
	}
}

func TestEncounterPauseResume(t *testing.T) {
	enc := NewEncounter(1, TypeAmbush, engine.GenreFantasy)
	enc.State = StateResolution

	enc.Pause()
	if !enc.IsPaused {
		t.Error("Expected encounter to be paused")
	}

	enc.Resume()
	if enc.IsPaused {
		t.Error("Expected encounter to be resumed")
	}
}

func TestGenerator(t *testing.T) {
	gen := NewGenerator(12345, engine.GenreFantasy)

	enc := gen.Generate(TypeAmbush)
	if enc.Type != TypeAmbush {
		t.Errorf("Expected TypeAmbush, got %d", enc.Type)
	}
	if enc.Title == "" {
		t.Error("Expected encounter to have a title")
	}

	enc2 := gen.GenerateRandom()
	if enc2.ID <= enc.ID {
		t.Error("Expected incrementing IDs")
	}
}

func TestGeneratorSetGenre(t *testing.T) {
	gen := NewGenerator(12345, engine.GenreFantasy)
	gen.SetGenre(engine.GenreScifi)

	enc := gen.Generate(TypeAmbush)
	if enc.Genre != engine.GenreScifi {
		t.Errorf("Expected genre scifi, got %s", enc.Genre)
	}
}

func TestGeneratorDeterminism(t *testing.T) {
	gen1 := NewGenerator(42, engine.GenreFantasy)
	gen2 := NewGenerator(42, engine.GenreFantasy)

	enc1 := gen1.Generate(TypeAmbush)
	enc2 := gen2.Generate(TypeAmbush)

	if enc1.Title != enc2.Title {
		t.Errorf("Expected same title, got %q vs %q", enc1.Title, enc2.Title)
	}
	if enc1.Difficulty != enc2.Difficulty {
		t.Errorf("Expected same difficulty, got %f vs %f", enc1.Difficulty, enc2.Difficulty)
	}
}

func TestTypeName(t *testing.T) {
	tests := []struct {
		encType EncounterType
		genre   engine.GenreID
		want    string
	}{
		{TypeAmbush, engine.GenreFantasy, "Ambush"},
		{TypeAmbush, engine.GenreScifi, "Hostile Contact"},
		{TypeNegotiation, engine.GenreCyberpunk, "Deal"},
		{TypeRace, engine.GenreHorror, "Flight"},
		{TypeCrisis, engine.GenrePostapoc, "Disaster"},
	}

	for _, tt := range tests {
		got := TypeName(tt.encType, tt.genre)
		if got != tt.want {
			t.Errorf("TypeName(%d, %s) = %q, want %q", tt.encType, tt.genre, got, tt.want)
		}
	}
}

func TestRoleNames(t *testing.T) {
	names := RoleNames(engine.GenreCyberpunk)

	if names[RoleFighter] != "Solo" {
		t.Errorf("Expected 'Solo', got %q", names[RoleFighter])
	}
	if names[RoleMedic] != "Medtech" {
		t.Errorf("Expected 'Medtech', got %q", names[RoleMedic])
	}
}

func TestOutcomeNames(t *testing.T) {
	names := OutcomeNames(engine.GenreHorror)

	if names[OutcomeVictory] != "Survived" {
		t.Errorf("Expected 'Survived', got %q", names[OutcomeVictory])
	}
	if names[OutcomeDefeat] != "Lost" {
		t.Errorf("Expected 'Lost', got %q", names[OutcomeDefeat])
	}
}

func TestCalculateRoleEffectiveness(t *testing.T) {
	member := crew.NewCrewMember(1, "Test", crew.TraitBrave, crew.SkillWarrior)
	member.Health = 100

	eff := CalculateRoleEffectiveness(member, RoleFighter)
	if eff < 0.8 {
		t.Errorf("Expected effectiveness >= 0.8 for warrior in fighter role, got %f", eff)
	}

	// Lower health should reduce effectiveness
	member.Health = 50
	effLow := CalculateRoleEffectiveness(member, RoleFighter)
	if effLow >= eff {
		t.Error("Expected lower effectiveness with lower health")
	}

	// Dead member should have 0 effectiveness
	member.IsAlive = false
	effDead := CalculateRoleEffectiveness(member, RoleFighter)
	if effDead != 0 {
		t.Errorf("Expected 0 effectiveness for dead member, got %f", effDead)
	}
}

func TestResolver(t *testing.T) {
	resolver := NewResolver(12345, engine.GenreFantasy)
	gen := NewGenerator(12345, engine.GenreFantasy)
	crewGen := crew.NewGenerator(12345, engine.GenreFantasy)

	enc := gen.Generate(TypeAmbush)

	party := crew.NewParty(engine.GenreFantasy, 4)
	for i := 0; i < 3; i++ {
		party.Add(crewGen.Generate())
	}

	// Assign crew to required roles
	living := party.Living()
	if len(living) > 0 {
		enc.AssignCrew(RoleFighter, living[0].ID)
	}

	enc.Start()

	// Resolve all phases
	for enc.State == StateResolution {
		result := resolver.ResolvePhase(enc, party)
		if result.PhaseNumber == 0 {
			t.Error("Expected valid phase number")
		}
	}

	if enc.State != StateComplete {
		t.Errorf("Expected state Complete, got %d", enc.State)
	}

	// Get final result
	finalResult := resolver.ResolveComplete(enc, party)
	if finalResult == nil {
		t.Error("Expected non-nil final result")
	}
}

func TestResolverApplyResult(t *testing.T) {
	resolver := NewResolver(12345, engine.GenreFantasy)

	result := NewEncounterResult(OutcomeVictory)
	result.MoraleDelta = 10
	result.FoodDelta = 5

	res := resources.NewResources(engine.GenreFantasy)
	initialMorale := res.Get(resources.ResourceMorale)
	initialFood := res.Get(resources.ResourceFood)

	party := crew.NewParty(engine.GenreFantasy, 4)
	v := vessel.New(engine.GenreFantasy, "Test Vessel")

	resolver.ApplyResult(result, res, party, v)

	if res.Get(resources.ResourceMorale) != initialMorale+10 {
		t.Error("Expected morale to increase by 10")
	}
	if res.Get(resources.ResourceFood) != initialFood+5 {
		t.Error("Expected food to increase by 5")
	}
}

func TestResolverSetGenre(t *testing.T) {
	resolver := NewResolver(12345, engine.GenreFantasy)
	resolver.SetGenre(engine.GenreScifi)

	if resolver.genre != engine.GenreScifi {
		t.Errorf("Expected genre scifi, got %s", resolver.genre)
	}
}

func TestAllGenres(t *testing.T) {
	genres := engine.AllGenres()

	for _, genre := range genres {
		// Test that all genres have encounter templates
		templates := encounterTemplates[genre]
		if templates == nil {
			t.Errorf("Missing encounter templates for genre %s", genre)
			continue
		}

		// Test each encounter type has templates
		for _, encType := range AllEncounterTypes() {
			typeTemplates := templates[encType]
			if len(typeTemplates) == 0 {
				t.Errorf("Missing templates for genre %s, type %d", genre, encType)
			}
		}

		// Test role names exist
		roleNamesForGenre := RoleNames(genre)
		for _, role := range AllEncounterRoles() {
			if roleNamesForGenre[role] == "" {
				t.Errorf("Missing role name for genre %s, role %d", genre, role)
			}
		}

		// Test outcome names exist
		outcomeNamesForGenre := OutcomeNames(genre)
		outcomes := []EncounterOutcome{OutcomeVictory, OutcomePartialSuccess, OutcomeRetreat, OutcomeDefeat}
		for _, outcome := range outcomes {
			if outcomeNamesForGenre[outcome] == "" {
				t.Errorf("Missing outcome name for genre %s, outcome %d", genre, outcome)
			}
		}

		// Test type names exist
		for _, encType := range AllEncounterTypes() {
			name := TypeName(encType, genre)
			if name == "" {
				t.Errorf("Missing type name for genre %s, type %d", genre, encType)
			}
		}
	}
}
