package events

import (
	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// RelationshipEventType identifies the type of relationship event.
type RelationshipEventType int

const (
	// RelationEventBond represents a positive relationship event.
	RelationEventBond RelationshipEventType = iota
	// RelationEventRivalry represents a negative relationship event.
	RelationEventRivalry
	// RelationEventMentorship represents a teaching/learning event.
	RelationEventMentorship
	// RelationEventRomance represents a romantic relationship event.
	RelationEventRomance
)

// RelationshipEventTemplate defines a template for relationship events.
type RelationshipEventTemplate struct {
	Type        RelationshipEventType
	Title       string
	Description string
	Choices     []ChoiceTemplate
}

// GenerateCrewRelationshipEvent creates an event based on crew relationships.
// Returns nil if no strong relationships warrant an event.
func GenerateCrewRelationshipEvent(gen *seed.Generator, network *crew.RelationshipNetwork, party *crew.Party, genre engine.GenreID) *Event {
	if network == nil || party == nil {
		return nil
	}

	rel := findSignificantRelationship(gen, network)
	if rel == nil {
		return nil
	}

	memberA, memberB := getAliveMembers(party, rel)
	if memberA == nil || memberB == nil {
		return nil
	}

	tmpl := selectRelationshipTemplate(gen, genre, rel.Type)
	if tmpl == nil {
		return nil
	}

	return buildRelationshipEvent(tmpl, memberA.Name, memberB.Name, rel.Strength, genre)
}

// findSignificantRelationship selects a relationship strong enough to trigger an event.
func findSignificantRelationship(gen *seed.Generator, network *crew.RelationshipNetwork) *crew.Relationship {
	relationships := network.AllRelationships()
	if len(relationships) == 0 {
		return nil
	}

	strongRelationships := filterStrongRelationships(relationships)
	if len(strongRelationships) == 0 {
		return nil
	}

	return seed.Choice(gen, strongRelationships)
}

// filterStrongRelationships returns relationships with significant positive or negative strength.
func filterStrongRelationships(relationships []*crew.Relationship) []*crew.Relationship {
	var result []*crew.Relationship
	for _, rel := range relationships {
		if rel.Strength > 50 || rel.Strength < -30 {
			result = append(result, rel)
		}
	}
	return result
}

// getAliveMembers retrieves the alive crew members involved in a relationship.
func getAliveMembers(party *crew.Party, rel *crew.Relationship) (*crew.CrewMember, *crew.CrewMember) {
	memberA := party.Get(rel.MemberA)
	memberB := party.Get(rel.MemberB)
	if memberA == nil || memberB == nil {
		return nil, nil
	}
	if !memberA.IsAlive || !memberB.IsAlive {
		return nil, nil
	}
	return memberA, memberB
}

// selectRelationshipTemplate chooses an appropriate template for the genre and relation type.
func selectRelationshipTemplate(gen *seed.Generator, genre engine.GenreID, relType crew.RelationType) *RelationshipEventTemplate {
	templates := relationshipTemplates[genre][relType]
	if len(templates) == 0 {
		templates = relationshipTemplates[engine.GenreFantasy][relType]
	}
	if len(templates) == 0 {
		return nil
	}
	tmpl := seed.Choice(gen, templates)
	return &tmpl
}

// buildRelationshipEvent constructs the final event from a template.
func buildRelationshipEvent(tmpl *RelationshipEventTemplate, nameA, nameB string, strength float64, genre engine.GenreID) *Event {
	event := NewEvent(0, CategoryCrew, tmpl.Title, tmpl.Description, genre)
	event.Title = formatRelationshipText(event.Title, nameA, nameB)
	event.Description = formatRelationshipText(event.Description, nameA, nameB)

	for _, c := range tmpl.Choices {
		outcome := c.Outcome
		outcome.MoraleDelta *= (1 + strength/200)
		event.AddChoice(c.Text, outcome)
	}

	return event
}

// formatRelationshipText replaces %A and %B with crew member names.
func formatRelationshipText(text, nameA, nameB string) string {
	result := text
	// Simple replacement - in production use proper template
	for i := 0; i < len(result)-1; i++ {
		if result[i] == '%' && result[i+1] == 'A' {
			result = result[:i] + nameA + result[i+2:]
		} else if result[i] == '%' && result[i+1] == 'B' {
			result = result[:i] + nameB + result[i+2:]
		}
	}
	return result
}

// relationshipTemplates maps genres and relationship types to event templates.
var relationshipTemplates = map[engine.GenreID]map[crew.RelationType][]RelationshipEventTemplate{
	engine.GenreFantasy: {
		crew.RelationFriendly: {
			{
				Type:        RelationEventBond,
				Title:       "%A and %B Share Stories",
				Description: "Around the campfire, %A and %B share tales of their pasts. The bond between them grows stronger.",
				Choices: []ChoiceTemplate{
					{Text: "Join in the storytelling", Outcome: EventOutcome{MoraleDelta: 8}},
					{Text: "Let them have their moment", Outcome: EventOutcome{MoraleDelta: 5}},
				},
			},
			{
				Type:        RelationEventBond,
				Title:       "A Shared Meal",
				Description: "%A prepared something special for %B. The gesture speaks volumes.",
				Choices: []ChoiceTemplate{
					{Text: "Share in the feast", Outcome: EventOutcome{MoraleDelta: 6, FoodDelta: -2}},
					{Text: "Let them eat privately", Outcome: EventOutcome{MoraleDelta: 4}},
				},
			},
		},
		crew.RelationRivalry: {
			{
				Type:        RelationEventRivalry,
				Title:       "Tension Erupts",
				Description: "%A and %B are at each other's throats again. The argument threatens to split the party.",
				Choices: []ChoiceTemplate{
					{Text: "Mediate the dispute", Outcome: EventOutcome{MoraleDelta: 5, TimeAdvance: 1}},
					{Text: "Let them fight it out", Outcome: EventOutcome{MoraleDelta: -8, CrewDamage: 5}},
					{Text: "Take sides", Outcome: EventOutcome{MoraleDelta: -15}},
				},
			},
		},
		crew.RelationMentorship: {
			{
				Type:        RelationEventMentorship,
				Title:       "A Lesson Learned",
				Description: "%A teaches %B an important skill. The student makes progress.",
				Choices: []ChoiceTemplate{
					{Text: "Observe the lesson", Outcome: EventOutcome{MoraleDelta: 6}},
					{Text: "Ask to learn too", Outcome: EventOutcome{MoraleDelta: 4, TimeAdvance: 1}},
				},
			},
		},
		crew.RelationRomantic: {
			{
				Type:        RelationEventRomance,
				Title:       "A Quiet Moment",
				Description: "%A and %B steal a moment away from the others. Love blooms even in hardship.",
				Choices: []ChoiceTemplate{
					{Text: "Give them privacy", Outcome: EventOutcome{MoraleDelta: 10}},
					{Text: "We need to keep moving", Outcome: EventOutcome{MoraleDelta: -5}},
				},
			},
		},
	},
	engine.GenreScifi: {
		crew.RelationFriendly: {
			{
				Type:        RelationEventBond,
				Title:       "Crew Bonding",
				Description: "%A and %B spend their off-shift together in the rec room.",
				Choices: []ChoiceTemplate{
					{Text: "Join the crew", Outcome: EventOutcome{MoraleDelta: 7}},
					{Text: "Let them rest", Outcome: EventOutcome{MoraleDelta: 4}},
				},
			},
		},
		crew.RelationRivalry: {
			{
				Type:        RelationEventRivalry,
				Title:       "Professional Disagreement",
				Description: "%A and %B have incompatible ideas about ship operations.",
				Choices: []ChoiceTemplate{
					{Text: "Hold a crew meeting", Outcome: EventOutcome{MoraleDelta: 3, TimeAdvance: 1}},
					{Text: "Override them both", Outcome: EventOutcome{MoraleDelta: -10}},
				},
			},
		},
	},
	engine.GenreHorror: {
		crew.RelationFriendly: {
			{
				Type:        RelationEventBond,
				Title:       "Shared Survival",
				Description: "%A and %B watch each other's backs. In this nightmare, trust is everything.",
				Choices: []ChoiceTemplate{
					{Text: "Encourage the bond", Outcome: EventOutcome{MoraleDelta: 8}},
					{Text: "Stay focused on survival", Outcome: EventOutcome{MoraleDelta: 3}},
				},
			},
		},
		crew.RelationRivalry: {
			{
				Type:        RelationEventRivalry,
				Title:       "Blame Game",
				Description: "%A blames %B for attracting those things. Accusations fly.",
				Choices: []ChoiceTemplate{
					{Text: "Shut it down hard", Outcome: EventOutcome{MoraleDelta: -5}},
					{Text: "Let them vent", Outcome: EventOutcome{MoraleDelta: -10, TimeAdvance: 1}},
				},
			},
		},
	},
	engine.GenreCyberpunk: {
		crew.RelationFriendly: {
			{
				Type:        RelationEventBond,
				Title:       "Chooms Forever",
				Description: "%A and %B share a drink and swap war stories from the streets.",
				Choices: []ChoiceTemplate{
					{Text: "Buy a round", Outcome: EventOutcome{MoraleDelta: 8, CurrencyDelta: -5}},
					{Text: "Nod along", Outcome: EventOutcome{MoraleDelta: 4}},
				},
			},
		},
		crew.RelationRivalry: {
			{
				Type:        RelationEventRivalry,
				Title:       "Bad Blood",
				Description: "%A and %B have history. The kind that ends in flatlines.",
				Choices: []ChoiceTemplate{
					{Text: "Lay down the law", Outcome: EventOutcome{MoraleDelta: -5}},
					{Text: "Offer a cut to settle it", Outcome: EventOutcome{MoraleDelta: 2, CurrencyDelta: -10}},
				},
			},
		},
	},
	engine.GenrePostapoc: {
		crew.RelationFriendly: {
			{
				Type:        RelationEventBond,
				Title:       "Pack Loyalty",
				Description: "%A and %B share rations they didn't have to. That's trust.",
				Choices: []ChoiceTemplate{
					{Text: "Honor the gesture", Outcome: EventOutcome{MoraleDelta: 8}},
					{Text: "Remind them to ration carefully", Outcome: EventOutcome{MoraleDelta: 2}},
				},
			},
		},
		crew.RelationRivalry: {
			{
				Type:        RelationEventRivalry,
				Title:       "Resource Dispute",
				Description: "%A caught %B skimming extra water. Words are exchanged.",
				Choices: []ChoiceTemplate{
					{Text: "Public accounting", Outcome: EventOutcome{MoraleDelta: -5, TimeAdvance: 1}},
					{Text: "Private warning", Outcome: EventOutcome{MoraleDelta: -3}},
				},
			},
		},
	},
}
