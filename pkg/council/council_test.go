package council

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
)

func TestDecisionType(t *testing.T) {
	types := AllDecisionTypes()
	if len(types) != 5 {
		t.Errorf("Expected 5 decision types, got %d", len(types))
	}

	expected := []string{"Route", "Camp", "Trade", "Fight", "Split"}
	for i, dt := range types {
		if dt.String() != expected[i] {
			t.Errorf("Type %d: expected %q, got %q", i, expected[i], dt.String())
		}
	}
}

func TestVoteOption(t *testing.T) {
	if OptionRisky.String() != "Risky" {
		t.Errorf("OptionRisky: expected %q, got %q", "Risky", OptionRisky.String())
	}
	if OptionSafe.String() != "Safe" {
		t.Errorf("OptionSafe: expected %q, got %q", "Safe", OptionSafe.String())
	}
}

func TestNewCouncil(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		c := NewCouncil(genre)
		if c == nil {
			t.Errorf("NewCouncil(%s) returned nil", genre)
		}
		if c.genre != genre {
			t.Errorf("Genre: expected %s, got %s", genre, c.genre)
		}
	}
}

func TestCouncilSetGenre(t *testing.T) {
	c := NewCouncil(engine.GenreFantasy)
	c.SetGenre(engine.GenreScifi)
	if c.genre != engine.GenreScifi {
		t.Errorf("SetGenre failed: expected %s, got %s", engine.GenreScifi, c.genre)
	}
}

func TestCreateDecision(t *testing.T) {
	c := NewCouncil(engine.GenreFantasy)
	decision := c.CreateDecision(DecisionRoute)

	if decision == nil {
		t.Fatal("CreateDecision returned nil")
	}
	if decision.Type != DecisionRoute {
		t.Errorf("Type: expected %v, got %v", DecisionRoute, decision.Type)
	}
	if decision.Description == "" {
		t.Error("Description should not be empty")
	}
	if decision.RiskyDesc == "" {
		t.Error("RiskyDesc should not be empty")
	}
	if decision.SafeDesc == "" {
		t.Error("SafeDesc should not be empty")
	}
}

func TestHoldVote(t *testing.T) {
	c := NewCouncil(engine.GenreHorror)
	decision := c.CreateDecision(DecisionFight)

	// Create a test party
	party := crew.NewParty(engine.GenreHorror, 5)
	gen := crew.NewGenerator(0, engine.GenreHorror)
	for i := 0; i < 3; i++ {
		member := gen.Generate()
		party.Add(member)
	}

	result := c.HoldVote(decision, party)

	if result == nil {
		t.Fatal("HoldVote returned nil")
	}
	if result.Decision != DecisionFight {
		t.Errorf("Decision: expected %v, got %v", DecisionFight, result.Decision)
	}
	if len(result.Votes) != 3 {
		t.Errorf("Votes: expected 3, got %d", len(result.Votes))
	}
	if result.RiskyVotes+result.SafeVotes != 3 {
		t.Errorf("Vote count mismatch: %d + %d != 3", result.RiskyVotes, result.SafeVotes)
	}
}

func TestApplyChoiceUnanimous(t *testing.T) {
	c := NewCouncil(engine.GenreCyberpunk)

	result := &VoteResult{
		Decision:   DecisionTrade,
		RiskyVotes: 3,
		SafeVotes:  0,
		Unanimous:  true,
	}

	c.ApplyChoice(result, OptionRisky)

	if result.MoraleDelta <= 0 {
		t.Errorf("Unanimous agreement should give positive morale: got %f", result.MoraleDelta)
	}
	if result.Overruled {
		t.Error("Should not be overruled when following unanimous vote")
	}
}

func TestApplyChoiceOverrule(t *testing.T) {
	c := NewCouncil(engine.GenrePostapoc)

	result := &VoteResult{
		Decision:   DecisionCamp,
		RiskyVotes: 1,
		SafeVotes:  4,
		Unanimous:  false,
	}

	c.ApplyChoice(result, OptionRisky)

	if !result.Overruled {
		t.Error("Should be overruled when going against majority")
	}
	if result.MoraleDelta >= 0 {
		t.Errorf("Overruling should give negative morale: got %f", result.MoraleDelta)
	}
}

func TestApplyChoiceFollowVote(t *testing.T) {
	c := NewCouncil(engine.GenreFantasy)

	result := &VoteResult{
		Decision:   DecisionSplit,
		RiskyVotes: 4,
		SafeVotes:  2,
		Unanimous:  false,
	}

	c.ApplyChoice(result, OptionRisky)

	if result.Overruled {
		t.Error("Should not be overruled when following majority")
	}
	if result.MoraleDelta != 0 {
		t.Errorf("Following vote should give 0 morale delta: got %f", result.MoraleDelta)
	}
}

func TestVoteReasoning(t *testing.T) {
	c := NewCouncil(engine.GenreScifi)
	decision := c.CreateDecision(DecisionRoute)

	member := crew.NewCrewMember(1, "Test", crew.TraitBrave, crew.SkillWarrior)
	vote := c.generateVote(member, decision)

	if vote.Reasoning == "" {
		t.Error("Vote should have reasoning text")
	}
	if vote.CrewID != 1 {
		t.Errorf("CrewID: expected 1, got %d", vote.CrewID)
	}
}

func TestTraitVotePatterns(t *testing.T) {
	c := NewCouncil(engine.GenreFantasy)

	// Brave should vote risky
	opt, conf := c.traitVotePattern(crew.TraitBrave, DecisionFight)
	if opt != OptionRisky {
		t.Errorf("TraitBrave should vote Risky, got %v", opt)
	}
	if conf <= 0 {
		t.Error("Confidence should be positive")
	}

	// Cautious should vote safe
	opt, _ = c.traitVotePattern(crew.TraitCautious, DecisionRoute)
	if opt != OptionSafe {
		t.Errorf("TraitCautious should vote Safe, got %v", opt)
	}
}

func TestGetScene(t *testing.T) {
	c := NewCouncil(engine.GenreHorror)
	decision := c.CreateDecision(DecisionFight)
	scene := c.GetScene(decision)

	if scene == nil {
		t.Fatal("GetScene returned nil")
	}
	if scene.SceneName == "" {
		t.Error("Scene name should not be empty")
	}
	if scene.Opening == "" {
		t.Error("Opening should not be empty")
	}
	if scene.RiskyOption == "" {
		t.Error("RiskyOption should not be empty")
	}
	if scene.SafeOption == "" {
		t.Error("SafeOption should not be empty")
	}
}

func TestAllGenresHaveContent(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		// Check scene names
		if _, ok := sceneNames[genre]; !ok {
			t.Errorf("Missing scene name for genre %s", genre)
		}

		// Check decision descriptions
		descs, ok := decisionDescriptions[genre]
		if !ok {
			t.Errorf("Missing decision descriptions for genre %s", genre)
			continue
		}
		for _, dt := range AllDecisionTypes() {
			if _, ok := descs[dt]; !ok {
				t.Errorf("Missing description for %s/%s", genre, dt)
			}
		}

		// Check risky descriptions
		risky, ok := riskyDescriptions[genre]
		if !ok {
			t.Errorf("Missing risky descriptions for genre %s", genre)
			continue
		}
		for _, dt := range AllDecisionTypes() {
			if _, ok := risky[dt]; !ok {
				t.Errorf("Missing risky description for %s/%s", genre, dt)
			}
		}

		// Check safe descriptions
		safe, ok := safeDescriptions[genre]
		if !ok {
			t.Errorf("Missing safe descriptions for genre %s", genre)
			continue
		}
		for _, dt := range AllDecisionTypes() {
			if _, ok := safe[dt]; !ok {
				t.Errorf("Missing safe description for %s/%s", genre, dt)
			}
		}

		// Check vote reasonings
		reasonings, ok := voteReasonings[genre]
		if !ok {
			t.Errorf("Missing vote reasonings for genre %s", genre)
			continue
		}
		for _, trait := range crew.AllTraits() {
			if _, ok := reasonings[trait]; !ok {
				t.Errorf("Missing reasonings for trait %v in genre %s", trait, genre)
			}
		}

		// Check scene content
		if _, ok := sceneOpenings[genre]; !ok {
			t.Errorf("Missing scene openings for genre %s", genre)
		}
		if _, ok := discussionSnippets[genre]; !ok {
			t.Errorf("Missing discussion snippets for genre %s", genre)
		}
		if _, ok := sceneClosings[genre]; !ok {
			t.Errorf("Missing scene closings for genre %s", genre)
		}
	}
}

func TestDeterministicVoting(t *testing.T) {
	c1 := NewCouncil(engine.GenreFantasy)
	c2 := NewCouncil(engine.GenreFantasy)

	member := crew.NewCrewMember(1, "Test", crew.TraitBrave, crew.SkillWarrior)
	decision := c1.CreateDecision(DecisionRoute)

	vote1 := c1.generateVote(member, decision)
	vote2 := c2.generateVote(member, decision)

	if vote1.Option != vote2.Option {
		t.Errorf("Vote options should match: %v vs %v", vote1.Option, vote2.Option)
	}
}
