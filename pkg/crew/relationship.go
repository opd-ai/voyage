package crew

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// RelationType represents the relationship between two crew members.
type RelationType int

const (
	RelationNeutral RelationType = iota
	RelationFriendly
	RelationRomantic
	RelationRivalry
	RelationMentorship
)

// Relationship tracks the bond between two crew members.
type Relationship struct {
	MemberA      int          // First crew member ID
	MemberB      int          // Second crew member ID
	Type         RelationType // Current relationship type
	Strength     float64      // Strength of relationship [-100, 100]
	Interactions int          // Number of interactions
}

// MoraleModifier returns the morale effect for this relationship.
func (r *Relationship) MoraleModifier() float64 {
	switch r.Type {
	case RelationFriendly:
		return r.Strength * 0.05 // Up to +5 morale at max strength
	case RelationRomantic:
		return r.Strength * 0.08 // Up to +8 morale at max strength
	case RelationRivalry:
		// Rivalry has negative strength, multiply by positive factor for negative result
		return r.Strength * 0.03 // -100 * 0.03 = -3 morale
	case RelationMentorship:
		return r.Strength * 0.04 // Up to +4 morale
	}
	return 0
}

// RelationshipNetwork manages all relationships in a party.
type RelationshipNetwork struct {
	relationships map[string]*Relationship
	genre         engine.GenreID
}

// NewRelationshipNetwork creates a new relationship network.
func NewRelationshipNetwork(genre engine.GenreID) *RelationshipNetwork {
	return &RelationshipNetwork{
		relationships: make(map[string]*Relationship),
		genre:         genre,
	}
}

// SetGenre changes the network's genre.
func (rn *RelationshipNetwork) SetGenre(genre engine.GenreID) {
	rn.genre = genre
}

// pairKey returns a consistent key for a pair of member IDs.
func pairKey(a, b int) string {
	if a > b {
		a, b = b, a
	}
	return string(rune(a)) + "-" + string(rune(b))
}

// GetRelationship returns the relationship between two members.
func (rn *RelationshipNetwork) GetRelationship(a, b int) *Relationship {
	key := pairKey(a, b)
	if rel, ok := rn.relationships[key]; ok {
		return rel
	}
	// Create neutral relationship if none exists
	rel := &Relationship{
		MemberA:  a,
		MemberB:  b,
		Type:     RelationNeutral,
		Strength: 0,
	}
	rn.relationships[key] = rel
	return rel
}

// Interact updates the relationship based on an interaction.
// Positive delta improves the relationship, negative worsens it.
func (rn *RelationshipNetwork) Interact(a, b int, delta float64) {
	rel := rn.GetRelationship(a, b)
	rel.Strength += delta
	rel.Interactions++

	// Clamp strength
	if rel.Strength > 100 {
		rel.Strength = 100
	}
	if rel.Strength < -100 {
		rel.Strength = -100
	}

	// Update relationship type based on strength
	rel.Type = relationTypeFromStrength(rel.Strength, rel.Interactions)
}

// relationTypeFromStrength determines relationship type.
func relationTypeFromStrength(strength float64, interactions int) RelationType {
	switch {
	case strength >= 75 && interactions >= 10:
		return RelationRomantic
	case strength >= 50:
		return RelationFriendly
	case strength >= 25 && interactions >= 5:
		return RelationMentorship
	case strength <= -50:
		return RelationRivalry
	default:
		return RelationNeutral
	}
}

// AllRelationships returns all relationships with non-neutral status or interactions.
func (rn *RelationshipNetwork) AllRelationships() []*Relationship {
	result := make([]*Relationship, 0)
	for _, rel := range rn.relationships {
		if rel.Type != RelationNeutral || rel.Interactions > 0 {
			result = append(result, rel)
		}
	}
	return result
}

// TotalMoraleModifier returns the sum of all morale modifiers.
func (rn *RelationshipNetwork) TotalMoraleModifier() float64 {
	total := 0.0
	for _, rel := range rn.relationships {
		total += rel.MoraleModifier()
	}
	return total
}

// RelationshipsFor returns all relationships involving a specific member.
func (rn *RelationshipNetwork) RelationshipsFor(memberID int) []*Relationship {
	result := make([]*Relationship, 0)
	for _, rel := range rn.relationships {
		if rel.MemberA == memberID || rel.MemberB == memberID {
			result = append(result, rel)
		}
	}
	return result
}

// GenerateInitialRelationships creates random starting relationships for a party.
func (rn *RelationshipNetwork) GenerateInitialRelationships(gen *seed.Generator, members []*CrewMember) {
	for i := 0; i < len(members); i++ {
		for j := i + 1; j < len(members); j++ {
			// 30% chance of pre-existing relationship
			if gen.Float64() < 0.3 {
				// Random initial strength [-30, 50]
				strength := gen.Float64()*80 - 30
				rn.Interact(members[i].ID, members[j].ID, strength)
			}
		}
	}
}

// RelationTypeName returns the display name for a relationship type.
func RelationTypeName(rt RelationType, genre engine.GenreID) string {
	names := relationTypeNames[genre]
	if name, ok := names[rt]; ok {
		return name
	}
	return "Unknown"
}

var relationTypeNames = map[engine.GenreID]map[RelationType]string{
	engine.GenreFantasy: {
		RelationNeutral:    "Acquaintance",
		RelationFriendly:   "Comrade",
		RelationRomantic:   "Beloved",
		RelationRivalry:    "Rival",
		RelationMentorship: "Mentor & Pupil",
	},
	engine.GenreScifi: {
		RelationNeutral:    "Colleague",
		RelationFriendly:   "Crew Mate",
		RelationRomantic:   "Partner",
		RelationRivalry:    "Rival",
		RelationMentorship: "Trainer & Cadet",
	},
	engine.GenreHorror: {
		RelationNeutral:    "Fellow Survivor",
		RelationFriendly:   "Trusted Ally",
		RelationRomantic:   "Bonded",
		RelationRivalry:    "Distrusted",
		RelationMentorship: "Protector & Ward",
	},
	engine.GenreCyberpunk: {
		RelationNeutral:    "Contact",
		RelationFriendly:   "Choom",
		RelationRomantic:   "Output",
		RelationRivalry:    "Gonk",
		RelationMentorship: "Fixer & Rookie",
	},
	engine.GenrePostapoc: {
		RelationNeutral:    "Fellow Traveler",
		RelationFriendly:   "Pack Mate",
		RelationRomantic:   "Bonded",
		RelationRivalry:    "Bad Blood",
		RelationMentorship: "Elder & Youth",
	},
}
