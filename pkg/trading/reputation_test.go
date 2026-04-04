package trading

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewTownReputationManager(t *testing.T) {
	trm := NewTownReputationManager(12345, engine.GenreFantasy)
	if trm == nil {
		t.Fatal("NewTownReputationManager returned nil")
	}
	if trm.Genre() != engine.GenreFantasy {
		t.Errorf("Expected fantasy genre, got %v", trm.Genre())
	}
	if trm.GlobalReputation() != 0.5 {
		t.Errorf("Expected initial global rep 0.5, got %f", trm.GlobalReputation())
	}
}

func TestTownReputationManager_SetGenre(t *testing.T) {
	trm := NewTownReputationManager(12345, engine.GenreFantasy)
	trm.GetReputation("town1", "Test Town")

	trm.SetGenre(engine.GenreScifi)

	if trm.Genre() != engine.GenreScifi {
		t.Errorf("Expected scifi genre, got %v", trm.Genre())
	}

	rep := trm.GetReputation("town1", "Test Town")
	if rep.Genre != engine.GenreScifi {
		t.Errorf("Town genre should have updated to scifi")
	}
}

func TestTownReputationManager_GetReputation(t *testing.T) {
	trm := NewTownReputationManager(12345, engine.GenreFantasy)

	rep1 := trm.GetReputation("town1", "First Town")
	rep2 := trm.GetReputation("town2", "Second Town")
	rep1Again := trm.GetReputation("town1", "First Town")

	if rep1 == nil || rep2 == nil {
		t.Fatal("GetReputation returned nil")
	}

	// Same town should return same object
	if rep1 != rep1Again {
		t.Error("Expected same reputation object for same town ID")
	}

	// Different towns should have different objects
	if rep1 == rep2 {
		t.Error("Expected different reputation objects for different towns")
	}
}

func TestTownReputation_ModifyReputation(t *testing.T) {
	trm := NewTownReputationManager(12345, engine.GenreFantasy)
	rep := trm.GetReputation("town1", "Test Town")

	initial := rep.Reputation
	rep.ModifyReputation(0.1)

	if rep.Reputation != initial+0.1 {
		t.Errorf("Expected reputation %f, got %f", initial+0.1, rep.Reputation)
	}

	// Test clamping
	rep.ModifyReputation(2.0) // Should clamp at 1.0
	if rep.Reputation > 1.0 {
		t.Errorf("Reputation should be clamped at 1.0, got %f", rep.Reputation)
	}

	rep.ModifyReputation(-3.0) // Should clamp at 0.0
	if rep.Reputation < 0 {
		t.Errorf("Reputation should be clamped at 0.0, got %f", rep.Reputation)
	}
}

func TestTownReputation_Behavior(t *testing.T) {
	trm := NewTownReputationManager(12345, engine.GenreFantasy)
	rep := trm.GetReputation("town1", "Test Town")

	testCases := []struct {
		reputation float64
		expected   TownBehavior
	}{
		{0.9, BehaviorWelcoming},
		{0.7, BehaviorFriendly},
		{0.5, BehaviorNeutral},
		{0.3, BehaviorSuspicious},
		{0.1, BehaviorHostile},
	}

	for _, tc := range testCases {
		rep.Reputation = tc.reputation
		rep.updateBehavior()
		if rep.Behavior != tc.expected {
			t.Errorf("Reputation %f should have behavior %v, got %v",
				tc.reputation, tc.expected, rep.Behavior)
		}
	}
}

func TestTownReputation_PriceModifier(t *testing.T) {
	trm := NewTownReputationManager(12345, engine.GenreFantasy)
	rep := trm.GetReputation("town1", "Test Town")

	// Test welcoming (discount)
	rep.Reputation = 0.9
	rep.updateBehavior()
	if rep.PriceModifier() >= 1.0 {
		t.Errorf("Welcoming should offer discount, got %f", rep.PriceModifier())
	}

	// Test hostile (markup)
	rep.Reputation = 0.1
	rep.updateBehavior()
	if rep.PriceModifier() <= 1.0 {
		t.Errorf("Hostile should have markup, got %f", rep.PriceModifier())
	}
}

func TestTownReputation_WillTrade(t *testing.T) {
	trm := NewTownReputationManager(12345, engine.GenreFantasy)
	rep := trm.GetReputation("town1", "Test Town")

	// Friendly should trade
	rep.Reputation = 0.7
	rep.updateBehavior()
	if !rep.WillTrade() {
		t.Error("Friendly town should trade")
	}

	// Very hostile should refuse
	rep.Reputation = 0.05
	rep.updateBehavior()
	if rep.WillTrade() {
		t.Error("Very hostile town should refuse trade")
	}

	// Moderately hostile might still trade
	rep.Reputation = 0.15
	rep.updateBehavior()
	if !rep.WillTrade() {
		t.Error("Moderately hostile town should still trade")
	}
}

func TestTownReputation_WillAllowEntry(t *testing.T) {
	trm := NewTownReputationManager(12345, engine.GenreFantasy)
	rep := trm.GetReputation("town1", "Test Town")

	// Normal reputation allows entry
	rep.Reputation = 0.5
	if !rep.WillAllowEntry() {
		t.Error("Neutral town should allow entry")
	}

	// Extremely hostile denies entry
	rep.Reputation = 0.02
	if rep.WillAllowEntry() {
		t.Error("Extremely hostile town should deny entry")
	}
}

func TestTownReputation_MayAttack(t *testing.T) {
	trm := NewTownReputationManager(12345, engine.GenreFantasy)
	rep := trm.GetReputation("town1", "Test Town")

	// Friendly should not attack
	rep.Reputation = 0.7
	rep.updateBehavior()
	if rep.MayAttack() {
		t.Error("Friendly town should not attack")
	}

	// Very hostile may attack
	rep.Reputation = 0.1
	rep.updateBehavior()
	if !rep.MayAttack() {
		t.Error("Very hostile town may attack")
	}

	// Hostile but above threshold should not attack
	rep.Reputation = 0.18
	rep.updateBehavior()
	if rep.MayAttack() {
		t.Error("Hostile town above 15% should not attack")
	}
}

func TestTownReputation_AttackChance(t *testing.T) {
	trm := NewTownReputationManager(12345, engine.GenreFantasy)
	rep := trm.GetReputation("town1", "Test Town")

	// Non-hostile has no attack chance
	rep.Reputation = 0.5
	rep.updateBehavior()
	if rep.AttackChance() != 0 {
		t.Errorf("Non-hostile should have 0 attack chance, got %f", rep.AttackChance())
	}

	// Lower reputation = higher attack chance
	rep.Reputation = 0.0
	rep.updateBehavior()
	chance0 := rep.AttackChance()

	rep.Reputation = 0.1
	rep.updateBehavior()
	chance10 := rep.AttackChance()

	if chance10 >= chance0 {
		t.Errorf("Attack chance should decrease with higher rep: 0%%=%f, 10%%=%f",
			chance0, chance10)
	}
}

func TestTownReputation_RecordEvents(t *testing.T) {
	trm := NewTownReputationManager(12345, engine.GenreFantasy)
	rep := trm.GetReputation("town1", "Test Town")
	rep.Reputation = 0.5 // Start at neutral

	initial := rep.Reputation
	rep.RecordEvent(EventSuccessfulTrade)

	if rep.Reputation <= initial {
		t.Error("Successful trade should increase reputation")
	}
	if rep.TradeCount != 1 {
		t.Errorf("Expected trade count 1, got %d", rep.TradeCount)
	}

	rep.Reputation = 0.5
	rep.RecordEvent(EventHelped)
	if rep.Reputation <= 0.5 {
		t.Error("Helping should increase reputation")
	}

	rep.Reputation = 0.5
	rep.RecordEvent(EventHarmed)
	if rep.Reputation >= 0.5 {
		t.Error("Harming should decrease reputation")
	}
}

func TestTownReputation_BehaviorChanged(t *testing.T) {
	trm := NewTownReputationManager(12345, engine.GenreFantasy)
	rep := trm.GetReputation("town1", "Test Town")

	rep.Reputation = 0.5
	rep.updateBehavior()

	// No change
	rep.Reputation = 0.45
	rep.updateBehavior()
	if rep.BehaviorChanged() {
		t.Error("Behavior should not have changed within same tier")
	}

	// Big change
	rep.LastBehavior = BehaviorNeutral
	rep.Reputation = 0.1
	rep.updateBehavior()
	if !rep.BehaviorChanged() {
		t.Error("Behavior should have changed between tiers")
	}
}

func TestTownReputationManager_GetFriendlyTowns(t *testing.T) {
	trm := NewTownReputationManager(12345, engine.GenreFantasy)

	rep1 := trm.GetReputation("town1", "Friendly Town")
	rep1.Reputation = 0.8
	rep1.updateBehavior()

	rep2 := trm.GetReputation("town2", "Hostile Town")
	rep2.Reputation = 0.1
	rep2.updateBehavior()

	rep3 := trm.GetReputation("town3", "Another Friendly")
	rep3.Reputation = 0.7
	rep3.updateBehavior()

	friendly := trm.GetFriendlyTowns()
	if len(friendly) != 2 {
		t.Errorf("Expected 2 friendly towns, got %d", len(friendly))
	}
}

func TestTownReputationManager_GetHostileTowns(t *testing.T) {
	trm := NewTownReputationManager(12345, engine.GenreFantasy)

	rep1 := trm.GetReputation("town1", "Friendly Town")
	rep1.Reputation = 0.8

	rep2 := trm.GetReputation("town2", "Hostile Town")
	rep2.Reputation = 0.1

	hostile := trm.GetHostileTowns()
	if len(hostile) != 1 {
		t.Errorf("Expected 1 hostile town, got %d", len(hostile))
	}
}

func TestTownReputationManager_UpdateGlobalReputation(t *testing.T) {
	trm := NewTownReputationManager(12345, engine.GenreFantasy)

	rep1 := trm.GetReputation("town1", "Town 1")
	rep1.Reputation = 0.8

	rep2 := trm.GetReputation("town2", "Town 2")
	rep2.Reputation = 0.4

	trm.UpdateGlobalReputation()

	// Expected is the average of the two explicit values
	expected := (rep1.Reputation + rep2.Reputation) / 2
	if trm.GlobalReputation() != expected {
		t.Errorf("Expected global rep %f, got %f", expected, trm.GlobalReputation())
	}
}

func TestBehaviorName_AllGenres(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	behaviors := []TownBehavior{
		BehaviorWelcoming,
		BehaviorFriendly,
		BehaviorNeutral,
		BehaviorSuspicious,
		BehaviorHostile,
	}

	for _, genre := range genres {
		trm := NewTownReputationManager(12345, genre)
		rep := trm.GetReputation("test", "Test")

		for _, behavior := range behaviors {
			rep.Behavior = behavior
			name := rep.BehaviorName()
			if name == "" || name == "Unknown" {
				t.Errorf("Genre %v behavior %v has empty/unknown name", genre, behavior)
			}
		}
	}
}

func TestReputationDescription_AllGenres(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	behaviors := []TownBehavior{
		BehaviorWelcoming,
		BehaviorFriendly,
		BehaviorNeutral,
		BehaviorSuspicious,
		BehaviorHostile,
	}

	for _, genre := range genres {
		trm := NewTownReputationManager(12345, genre)
		rep := trm.GetReputation("test", "Test")

		for _, behavior := range behaviors {
			rep.Behavior = behavior
			desc := rep.ReputationDescription()
			if desc == "" || desc == "Unknown standing" {
				t.Errorf("Genre %v behavior %v has empty/unknown description", genre, behavior)
			}
		}
	}
}

func TestHostileWarning_AllGenres(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		trm := NewTownReputationManager(12345, genre)
		rep := trm.GetReputation("test", "Test")
		rep.Reputation = 0.1
		rep.updateBehavior()

		warning := rep.HostileWarning()
		if warning == "" {
			t.Errorf("Genre %v should have hostile warning", genre)
		}
	}
}

func TestGetReputationVocab_AllGenres(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		vocab := GetReputationVocab(genre)
		if vocab == nil {
			t.Errorf("GetReputationVocab(%v) returned nil", genre)
			continue
		}
		if vocab.ReputationLabel == "" {
			t.Errorf("Genre %v has empty ReputationLabel", genre)
		}
		if vocab.BehaviorLabel == "" {
			t.Errorf("Genre %v has empty BehaviorLabel", genre)
		}
		if vocab.FriendlyLabel == "" {
			t.Errorf("Genre %v has empty FriendlyLabel", genre)
		}
		if vocab.HostileLabel == "" {
			t.Errorf("Genre %v has empty HostileLabel", genre)
		}
	}
}
