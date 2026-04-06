package npc

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// NPCType identifies different types of NPCs.
type NPCType int

const (
	// TypeTrader buys and sells goods.
	TypeTrader NPCType = iota
	// TypeRefugee seeks help or shares information.
	TypeRefugee
	// TypeBandit is hostile and may attack.
	TypeBandit
	// TypeTraveler is a neutral wanderer.
	TypeTraveler
	// TypeScout offers route information.
	TypeScout
	// TypeGuard patrols faction territory.
	TypeGuard
)

// AllNPCTypes returns all NPC types.
func AllNPCTypes() []NPCType {
	return []NPCType{
		TypeTrader,
		TypeRefugee,
		TypeBandit,
		TypeTraveler,
		TypeScout,
		TypeGuard,
	}
}

// NPCTypeName returns the genre-appropriate name for an NPC type.
// Uses comma-ok pattern to properly detect missing genres (L-012).
func NPCTypeName(t NPCType, genre engine.GenreID) string {
	names, ok := typeNames[genre]
	if !ok {
		names = typeNames[engine.GenreFantasy]
	}
	return names[t]
}

var typeNames = map[engine.GenreID]map[NPCType]string{
	engine.GenreFantasy: {
		TypeTrader:   "Merchant",
		TypeRefugee:  "Refugee",
		TypeBandit:   "Bandit",
		TypeTraveler: "Traveler",
		TypeScout:    "Ranger",
		TypeGuard:    "Guard",
	},
	engine.GenreScifi: {
		TypeTrader:   "Trader",
		TypeRefugee:  "Refugee",
		TypeBandit:   "Pirate",
		TypeTraveler: "Spacer",
		TypeScout:    "Scout",
		TypeGuard:    "Patrol",
	},
	engine.GenreHorror: {
		TypeTrader:   "Survivor",
		TypeRefugee:  "Refugee",
		TypeBandit:   "Raider",
		TypeTraveler: "Wanderer",
		TypeScout:    "Scout",
		TypeGuard:    "Watch",
	},
	engine.GenreCyberpunk: {
		TypeTrader:   "Fixer",
		TypeRefugee:  "Runaway",
		TypeBandit:   "Ganger",
		TypeTraveler: "Drifter",
		TypeScout:    "Netrunner",
		TypeGuard:    "Corpo Sec",
	},
	engine.GenrePostapoc: {
		TypeTrader:   "Scavenger",
		TypeRefugee:  "Wanderer",
		TypeBandit:   "Raider",
		TypeTraveler: "Drifter",
		TypeScout:    "Tracker",
		TypeGuard:    "Patrol",
	},
}

// Alignment represents an NPC's disposition.
type Alignment int

const (
	// AlignmentHostile attacks on sight.
	AlignmentHostile Alignment = iota
	// AlignmentSuspicious is wary but might trade.
	AlignmentSuspicious
	// AlignmentNeutral has no strong feelings.
	AlignmentNeutral
	// AlignmentFriendly offers good deals.
	AlignmentFriendly
	// AlignmentAllied will help actively.
	AlignmentAllied
)

// AlignmentName returns the display name for an alignment.
func AlignmentName(a Alignment) string {
	names := map[Alignment]string{
		AlignmentHostile:    "Hostile",
		AlignmentSuspicious: "Suspicious",
		AlignmentNeutral:    "Neutral",
		AlignmentFriendly:   "Friendly",
		AlignmentAllied:     "Allied",
	}
	if name, ok := names[a]; ok {
		return name
	}
	return "Unknown"
}

// NPC represents a procedurally generated non-player character.
type NPC struct {
	ID          int
	Name        string
	NPCType     NPCType
	Genre       engine.GenreID
	Alignment   Alignment
	FactionID   int // 0 = no faction
	Description string
	Dialogue    []string
	TradeGoods  []TradeGood // For traders only
}

// TradeGood represents an item an NPC trader has.
type TradeGood struct {
	Name     string
	Quantity int
	Price    float64
}

// NewNPC creates a new NPC.
func NewNPC(id int, name string, npcType NPCType, genre engine.GenreID) *NPC {
	return &NPC{
		ID:        id,
		Name:      name,
		NPCType:   npcType,
		Genre:     genre,
		Alignment: AlignmentNeutral,
	}
}

// SetGenre updates the NPC's genre.
func (n *NPC) SetGenre(genre engine.GenreID) {
	n.Genre = genre
}

// TypeDisplayName returns the genre-appropriate type name.
func (n *NPC) TypeDisplayName() string {
	return NPCTypeName(n.NPCType, n.Genre)
}

// IsHostile returns true if the NPC is hostile.
func (n *NPC) IsHostile() bool {
	return n.Alignment == AlignmentHostile
}

// CanTrade returns true if the NPC can trade.
func (n *NPC) CanTrade() bool {
	return n.NPCType == TypeTrader && n.Alignment != AlignmentHostile
}

// GetDefaultAlignment returns the typical alignment for an NPC type.
func GetDefaultAlignment(npcType NPCType) Alignment {
	defaults := map[NPCType]Alignment{
		TypeTrader:   AlignmentNeutral,
		TypeRefugee:  AlignmentFriendly,
		TypeBandit:   AlignmentHostile,
		TypeTraveler: AlignmentNeutral,
		TypeScout:    AlignmentSuspicious,
		TypeGuard:    AlignmentSuspicious,
	}
	if align, ok := defaults[npcType]; ok {
		return align
	}
	return AlignmentNeutral
}
