package council

import (
	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// Council manages crew voting for critical decisions.
type Council struct {
	genre   engine.GenreID
	seedGen *seed.Generator
}

// Decision represents a pending decision for the crew to vote on.
type Decision struct {
	Type        DecisionType
	Description string
	RiskyDesc   string
	SafeDesc    string
	genre       engine.GenreID
}

// NewCouncil creates a new council for the given genre.
func NewCouncil(genre engine.GenreID) *Council {
	return &Council{
		genre:   genre,
		seedGen: seed.NewGenerator(0, "council"),
	}
}

// SetGenre changes the council's genre.
func (c *Council) SetGenre(genre engine.GenreID) {
	c.genre = genre
}

// CreateDecision generates a decision of the given type.
func (c *Council) CreateDecision(decType DecisionType) *Decision {
	return &Decision{
		Type:        decType,
		Description: c.getDecisionDescription(decType),
		RiskyDesc:   c.getRiskyDescription(decType),
		SafeDesc:    c.getSafeDescription(decType),
		genre:       c.genre,
	}
}

// HoldVote conducts a vote among the party members.
func (c *Council) HoldVote(decision *Decision, party *crew.Party) *VoteResult {
	result := &VoteResult{
		Decision: decision.Type,
		Votes:    make([]Vote, 0),
	}

	members := party.Members()
	for _, member := range members {
		vote := c.generateVote(member, decision)
		result.Votes = append(result.Votes, vote)

		if vote.Option == OptionRisky {
			result.RiskyVotes++
		} else {
			result.SafeVotes++
		}
	}

	result.Unanimous = result.RiskyVotes == 0 || result.SafeVotes == 0
	return result
}

// ApplyChoice applies the player's decision and calculates morale effects.
func (c *Council) ApplyChoice(result *VoteResult, choice VoteOption) {
	result.PlayerChoice = choice
	majorityChoice := c.determineMajorityChoice(result)
	result.Overruled = (choice != majorityChoice) && !result.Unanimous
	result.MoraleDelta = c.calculateMoraleDelta(result, choice, majorityChoice)
}

// determineMajorityChoice returns the option favored by the majority.
func (c *Council) determineMajorityChoice(result *VoteResult) VoteOption {
	if result.RiskyVotes > result.SafeVotes {
		return OptionRisky
	}
	return OptionSafe
}

// calculateMoraleDelta computes morale change based on the voting outcome.
func (c *Council) calculateMoraleDelta(result *VoteResult, choice, majorityChoice VoteOption) float64 {
	if result.Unanimous && choice == majorityChoice {
		return 0.1 // Unanimous agreement bonus
	}
	if result.Overruled {
		return c.calculateDissentPenalty(result, choice)
	}
	return 0 // Followed the vote - no change
}

// calculateDissentPenalty computes the morale penalty for overruling the crew.
func (c *Council) calculateDissentPenalty(result *VoteResult, choice VoteOption) float64 {
	totalVotes := result.RiskyVotes + result.SafeVotes
	if totalVotes == 0 {
		return 0
	}
	var dissent int
	if choice == OptionRisky {
		dissent = result.SafeVotes
	} else {
		dissent = result.RiskyVotes
	}
	dissentRatio := float64(dissent) / float64(totalVotes)
	return -0.15 * dissentRatio
}

// generateVote determines how a crew member votes based on their trait.
func (c *Council) generateVote(member *crew.CrewMember, decision *Decision) Vote {
	vote := Vote{
		CrewID:   member.ID,
		CrewName: member.Name,
	}

	// Determine vote based on trait
	vote.Option, vote.Confidence = c.traitVotePattern(member.Trait, decision.Type)
	vote.Reasoning = c.generateReasoning(member, decision, vote.Option)

	return vote
}

// traitVotePattern returns how a trait typically votes.
func (c *Council) traitVotePattern(trait crew.Trait, decType DecisionType) (VoteOption, float64) {
	switch trait {
	case crew.TraitBrave:
		return OptionRisky, 0.9
	case crew.TraitCautious:
		return OptionSafe, 0.9
	case crew.TraitOptimistic:
		return OptionRisky, 0.7
	case crew.TraitPessimistic:
		return OptionSafe, 0.8
	case crew.TraitGreedy:
		// Greedy votes based on potential profit
		if decType == DecisionTrade {
			return OptionRisky, 0.95
		}
		return OptionRisky, 0.6
	case crew.TraitGenerous:
		return OptionSafe, 0.6
	case crew.TraitStoic:
		// Stoic is logical - often safe
		return OptionSafe, 0.7
	case crew.TraitEmotional:
		// Emotional varies - slight risk preference
		return OptionRisky, 0.5
	case crew.TraitNavigator:
		// Navigator prefers efficient routes (risky shortcuts)
		if decType == DecisionRoute {
			return OptionRisky, 0.85
		}
		return OptionSafe, 0.5
	case crew.TraitScavenger:
		// Scavenger likes opportunities
		return OptionRisky, 0.65
	default:
		return OptionSafe, 0.5
	}
}

// generateReasoning creates trait-appropriate reasoning text.
func (c *Council) generateReasoning(member *crew.CrewMember, decision *Decision, option VoteOption) string {
	reasonings := voteReasonings[c.genre][member.Trait][option]
	if len(reasonings) == 0 {
		// Fallback to generic
		if option == OptionRisky {
			return "I say we take the chance."
		}
		return "Let's play it safe."
	}
	return seed.Choice(c.seedGen, reasonings)
}

// GetScene returns the council scene description.
func (c *Council) GetScene(decision *Decision) *CouncilScene {
	return &CouncilScene{
		SceneName:   sceneNames[c.genre],
		Opening:     c.getSceneOpening(decision),
		RiskyOption: decision.RiskyDesc,
		SafeOption:  decision.SafeDesc,
		Discussion:  c.generateDiscussion(decision),
		Closing:     c.getSceneClosing(),
	}
}

// getDecisionDescription returns a genre-appropriate decision description.
func (c *Council) getDecisionDescription(decType DecisionType) string {
	descriptions := decisionDescriptions[c.genre][decType]
	return seed.Choice(c.seedGen, descriptions)
}

// getRiskyDescription returns the risky option description.
func (c *Council) getRiskyDescription(decType DecisionType) string {
	descriptions := riskyDescriptions[c.genre][decType]
	return seed.Choice(c.seedGen, descriptions)
}

// getSafeDescription returns the safe option description.
func (c *Council) getSafeDescription(decType DecisionType) string {
	descriptions := safeDescriptions[c.genre][decType]
	return seed.Choice(c.seedGen, descriptions)
}

// getSceneOpening returns the opening text for the council scene.
func (c *Council) getSceneOpening(decision *Decision) string {
	openings := sceneOpenings[c.genre]
	return seed.Choice(c.seedGen, openings)
}

// generateDiscussion creates discussion snippets.
func (c *Council) generateDiscussion(decision *Decision) []string {
	snippets := discussionSnippets[c.genre]
	count := 2 + c.seedGen.Intn(2) // 2-3 snippets
	discussion := make([]string, 0, count)
	for i := 0; i < count && i < len(snippets); i++ {
		discussion = append(discussion, snippets[i])
	}
	return discussion
}

// getSceneClosing returns the closing text.
func (c *Council) getSceneClosing() string {
	closings := sceneClosings[c.genre]
	return seed.Choice(c.seedGen, closings)
}
