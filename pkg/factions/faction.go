package factions

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// FactionID uniquely identifies a faction within a run.
type FactionID int

// Relationship represents the state between two factions.
type Relationship int

const (
	// RelationshipAllied indicates factions actively help each other.
	RelationshipAllied Relationship = iota
	// RelationshipFriendly indicates positive relations.
	RelationshipFriendly
	// RelationshipNeutral indicates no strong feelings.
	RelationshipNeutral
	// RelationshipSuspicious indicates wary but not hostile.
	RelationshipSuspicious
	// RelationshipHostile indicates factions attack on sight.
	RelationshipHostile
)

// AllRelationships returns all relationship types.
func AllRelationships() []Relationship {
	return []Relationship{
		RelationshipAllied,
		RelationshipFriendly,
		RelationshipNeutral,
		RelationshipSuspicious,
		RelationshipHostile,
	}
}

// RelationshipName returns the display name for a relationship.
func RelationshipName(r Relationship) string {
	names := map[Relationship]string{
		RelationshipAllied:     "Allied",
		RelationshipFriendly:   "Friendly",
		RelationshipNeutral:    "Neutral",
		RelationshipSuspicious: "Suspicious",
		RelationshipHostile:    "Hostile",
	}
	if name, ok := names[r]; ok {
		return name
	}
	return "Unknown"
}

// Ideology represents a faction's core belief system.
type Ideology int

const (
	// IdeologyMerchant focuses on trade and profit.
	IdeologyMerchant Ideology = iota
	// IdeologyMilitary focuses on conquest and defense.
	IdeologyMilitary
	// IdeologyReligious focuses on faith and tradition.
	IdeologyReligious
	// IdeologyCriminal focuses on lawlessness and gain.
	IdeologyCriminal
	// IdeologyScientific focuses on knowledge and progress.
	IdeologyScientific
	// IdeologySurvivalist focuses on survival at any cost.
	IdeologySurvivalist
)

// AllIdeologies returns all ideology types.
func AllIdeologies() []Ideology {
	return []Ideology{
		IdeologyMerchant,
		IdeologyMilitary,
		IdeologyReligious,
		IdeologyCriminal,
		IdeologyScientific,
		IdeologySurvivalist,
	}
}

// IdeologyName returns the display name for an ideology by genre.
func IdeologyName(i Ideology, genre engine.GenreID) string {
	names := ideologyNames[genre]
	if names == nil {
		names = ideologyNames[engine.GenreFantasy]
	}
	return names[i]
}

var ideologyNames = map[engine.GenreID]map[Ideology]string{
	engine.GenreFantasy: {
		IdeologyMerchant:    "Merchant Guild",
		IdeologyMilitary:    "Knight Order",
		IdeologyReligious:   "Sacred Cult",
		IdeologyCriminal:    "Thieves' Guild",
		IdeologyScientific:  "Mage Circle",
		IdeologySurvivalist: "Ranger Band",
	},
	engine.GenreScifi: {
		IdeologyMerchant:    "Trade Corporation",
		IdeologyMilitary:    "Military Fleet",
		IdeologyReligious:   "Tech Cult",
		IdeologyCriminal:    "Pirate Syndicate",
		IdeologyScientific:  "Research Collective",
		IdeologySurvivalist: "Colonist Union",
	},
	engine.GenreHorror: {
		IdeologyMerchant:    "Trader Network",
		IdeologyMilitary:    "Militia",
		IdeologyReligious:   "Doomsday Cult",
		IdeologyCriminal:    "Raider Gang",
		IdeologyScientific:  "Research Bunker",
		IdeologySurvivalist: "Survivor Band",
	},
	engine.GenreCyberpunk: {
		IdeologyMerchant:    "Megacorporation",
		IdeologyMilitary:    "Private Military",
		IdeologyReligious:   "Neo-Faith",
		IdeologyCriminal:    "Crime Syndicate",
		IdeologyScientific:  "Tech Collective",
		IdeologySurvivalist: "Street Gang",
	},
	engine.GenrePostapoc: {
		IdeologyMerchant:    "Caravan Guild",
		IdeologyMilitary:    "Warlord Clan",
		IdeologyReligious:   "Atom Cult",
		IdeologyCriminal:    "Raider Horde",
		IdeologyScientific:  "Vault Dwellers",
		IdeologySurvivalist: "Settler Coalition",
	},
}

// Faction represents a procedurally generated faction.
type Faction struct {
	ID          FactionID
	Name        string
	Genre       engine.GenreID
	Ideology    Ideology
	Description string
	Territory   []TerritoryBlock
	Relations   map[FactionID]Relationship
}

// TerritoryBlock represents a region controlled by a faction.
type TerritoryBlock struct {
	X      int
	Y      int
	Radius int
}

// NewFaction creates a new faction with default values.
func NewFaction(id FactionID, name string, ideology Ideology, genre engine.GenreID) *Faction {
	return &Faction{
		ID:        id,
		Name:      name,
		Genre:     genre,
		Ideology:  ideology,
		Relations: make(map[FactionID]Relationship),
		Territory: make([]TerritoryBlock, 0),
	}
}

// SetGenre updates the faction's genre theme.
func (f *Faction) SetGenre(genre engine.GenreID) {
	f.Genre = genre
}

// IdeologyDisplayName returns the genre-appropriate ideology name.
func (f *Faction) IdeologyDisplayName() string {
	return IdeologyName(f.Ideology, f.Genre)
}

// GetRelation returns the relationship with another faction.
func (f *Faction) GetRelation(other FactionID) Relationship {
	if r, ok := f.Relations[other]; ok {
		return r
	}
	return RelationshipNeutral
}

// SetRelation sets the relationship with another faction.
func (f *Faction) SetRelation(other FactionID, rel Relationship) {
	f.Relations[other] = rel
}

// AddTerritory adds a territory block to this faction.
func (f *Faction) AddTerritory(x, y, radius int) {
	f.Territory = append(f.Territory, TerritoryBlock{X: x, Y: y, Radius: radius})
}

// ControlsPosition returns true if the faction controls the given position.
func (f *Faction) ControlsPosition(x, y int) bool {
	for _, t := range f.Territory {
		dx := x - t.X
		dy := y - t.Y
		if dx*dx+dy*dy <= t.Radius*t.Radius {
			return true
		}
	}
	return false
}

// FactionManager tracks all factions and player reputation.
type FactionManager struct {
	Factions         map[FactionID]*Faction
	PlayerReputation map[FactionID]int // -100 to +100
	genre            engine.GenreID
}

// NewFactionManager creates a new faction manager.
func NewFactionManager(genre engine.GenreID) *FactionManager {
	return &FactionManager{
		Factions:         make(map[FactionID]*Faction),
		PlayerReputation: make(map[FactionID]int),
		genre:            genre,
	}
}

// SetGenre updates the manager's genre and all managed factions.
func (m *FactionManager) SetGenre(genre engine.GenreID) {
	m.genre = genre
	for _, f := range m.Factions {
		f.SetGenre(genre)
	}
}

// AddFaction adds a faction to the manager.
func (m *FactionManager) AddFaction(f *Faction) {
	m.Factions[f.ID] = f
	m.PlayerReputation[f.ID] = 0
}

// GetFaction returns a faction by ID.
func (m *FactionManager) GetFaction(id FactionID) *Faction {
	return m.Factions[id]
}

// GetReputation returns the player's reputation with a faction.
func (m *FactionManager) GetReputation(id FactionID) int {
	return m.PlayerReputation[id]
}

// AdjustReputation changes the player's reputation with a faction.
func (m *FactionManager) AdjustReputation(id FactionID, delta int) {
	current := m.PlayerReputation[id]
	newRep := current + delta
	if newRep > 100 {
		newRep = 100
	}
	if newRep < -100 {
		newRep = -100
	}
	m.PlayerReputation[id] = newRep
}

// ReputationToRelationship converts player reputation to relationship level.
func ReputationToRelationship(rep int) Relationship {
	switch {
	case rep >= 75:
		return RelationshipAllied
	case rep >= 25:
		return RelationshipFriendly
	case rep >= -25:
		return RelationshipNeutral
	case rep >= -75:
		return RelationshipSuspicious
	default:
		return RelationshipHostile
	}
}

// GetPlayerRelation returns the player's relationship level with a faction.
func (m *FactionManager) GetPlayerRelation(id FactionID) Relationship {
	return ReputationToRelationship(m.PlayerReputation[id])
}

// GetFactionsAt returns all factions that control the given position.
func (m *FactionManager) GetFactionsAt(x, y int) []*Faction {
	var result []*Faction
	for _, f := range m.Factions {
		if f.ControlsPosition(x, y) {
			result = append(result, f)
		}
	}
	return result
}

// AllFactions returns all factions as a slice.
func (m *FactionManager) AllFactions() []*Faction {
	result := make([]*Faction, 0, len(m.Factions))
	for _, f := range m.Factions {
		result = append(result, f)
	}
	return result
}
