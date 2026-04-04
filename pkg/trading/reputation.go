package trading

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// TownBehavior represents how a town acts toward the player.
type TownBehavior int

const (
	// BehaviorWelcoming - town offers discounts and bonuses
	BehaviorWelcoming TownBehavior = iota
	// BehaviorFriendly - town is helpful, normal or better prices
	BehaviorFriendly
	// BehaviorNeutral - town acts normally
	BehaviorNeutral
	// BehaviorSuspicious - town has higher prices, may deny services
	BehaviorSuspicious
	// BehaviorHostile - town may attack or deny entry
	BehaviorHostile
)

// TownReputationManager tracks player reputation across multiple towns.
type TownReputationManager struct {
	reputations map[string]*TownReputation
	globalRep   float64 // Average reputation affects new towns
	genre       engine.GenreID
	gen         *seed.Generator
}

// TownReputation tracks reputation and interactions with a specific town.
type TownReputation struct {
	TownID       string
	TownName     string
	Reputation   float64 // 0-1
	Behavior     TownBehavior
	VisitCount   int
	TradeCount   int
	LastBehavior TownBehavior // Previous behavior (for detecting changes)
	Genre        engine.GenreID
}

// NewTownReputationManager creates a new reputation manager.
func NewTownReputationManager(masterSeed int64, genre engine.GenreID) *TownReputationManager {
	return &TownReputationManager{
		reputations: make(map[string]*TownReputation),
		globalRep:   0.5,
		genre:       genre,
		gen:         seed.NewGenerator(masterSeed, "reputation"),
	}
}

// SetGenre changes the genre for vocabulary.
func (trm *TownReputationManager) SetGenre(genre engine.GenreID) {
	trm.genre = genre
	for _, rep := range trm.reputations {
		rep.Genre = genre
	}
}

// Genre returns the current genre.
func (trm *TownReputationManager) Genre() engine.GenreID {
	return trm.genre
}

// GetReputation retrieves or creates reputation for a town.
func (trm *TownReputationManager) GetReputation(townID, townName string) *TownReputation {
	rep, exists := trm.reputations[townID]
	if !exists {
		// New town starts near global reputation with some variation
		startRep := trm.globalRep + (trm.gen.Float64()-0.5)*0.2
		if startRep < 0.1 {
			startRep = 0.1
		}
		if startRep > 0.9 {
			startRep = 0.9
		}

		rep = &TownReputation{
			TownID:     townID,
			TownName:   townName,
			Reputation: startRep,
			Genre:      trm.genre,
		}
		rep.updateBehavior()
		trm.reputations[townID] = rep
	}
	return rep
}

// UpdateGlobalReputation recalculates global reputation from all towns.
func (trm *TownReputationManager) UpdateGlobalReputation() {
	if len(trm.reputations) == 0 {
		return
	}

	total := 0.0
	for _, rep := range trm.reputations {
		total += rep.Reputation
	}
	trm.globalRep = total / float64(len(trm.reputations))
}

// GlobalReputation returns the average reputation across all towns.
func (trm *TownReputationManager) GlobalReputation() float64 {
	return trm.globalRep
}

// GetAllTowns returns all tracked towns.
func (trm *TownReputationManager) GetAllTowns() []*TownReputation {
	towns := make([]*TownReputation, 0, len(trm.reputations))
	for _, rep := range trm.reputations {
		towns = append(towns, rep)
	}
	return towns
}

// GetFriendlyTowns returns towns with friendly or better reputation.
func (trm *TownReputationManager) GetFriendlyTowns() []*TownReputation {
	friendly := make([]*TownReputation, 0)
	for _, rep := range trm.reputations {
		if rep.Reputation >= 0.6 {
			friendly = append(friendly, rep)
		}
	}
	return friendly
}

// GetHostileTowns returns towns with hostile reputation.
func (trm *TownReputationManager) GetHostileTowns() []*TownReputation {
	hostile := make([]*TownReputation, 0)
	for _, rep := range trm.reputations {
		if rep.Reputation < 0.2 {
			hostile = append(hostile, rep)
		}
	}
	return hostile
}

// updateBehavior updates the town behavior based on reputation.
func (tr *TownReputation) updateBehavior() {
	tr.LastBehavior = tr.Behavior

	switch {
	case tr.Reputation >= 0.8:
		tr.Behavior = BehaviorWelcoming
	case tr.Reputation >= 0.6:
		tr.Behavior = BehaviorFriendly
	case tr.Reputation >= 0.4:
		tr.Behavior = BehaviorNeutral
	case tr.Reputation >= 0.2:
		tr.Behavior = BehaviorSuspicious
	default:
		tr.Behavior = BehaviorHostile
	}
}

// ModifyReputation changes reputation and updates behavior.
func (tr *TownReputation) ModifyReputation(delta float64) {
	tr.Reputation += delta
	if tr.Reputation < 0 {
		tr.Reputation = 0
	}
	if tr.Reputation > 1 {
		tr.Reputation = 1
	}
	tr.updateBehavior()
}

// RecordVisit records a visit to the town.
func (tr *TownReputation) RecordVisit() {
	tr.VisitCount++
	// Small reputation boost for visiting
	tr.ModifyReputation(0.01)
}

// RecordTrade records a trade at the town.
func (tr *TownReputation) RecordTrade(success bool) {
	tr.TradeCount++
	if success {
		tr.ModifyReputation(0.02)
	} else {
		tr.ModifyReputation(-0.01)
	}
}

// BehaviorChanged returns true if behavior changed since last check.
func (tr *TownReputation) BehaviorChanged() bool {
	return tr.Behavior != tr.LastBehavior
}

// PriceModifier returns the price multiplier based on reputation.
func (tr *TownReputation) PriceModifier() float64 {
	// Welcoming: 10% discount
	// Friendly: 5% discount
	// Neutral: Normal
	// Suspicious: 15% markup
	// Hostile: 30% markup (if they trade at all)
	switch tr.Behavior {
	case BehaviorWelcoming:
		return 0.9
	case BehaviorFriendly:
		return 0.95
	case BehaviorNeutral:
		return 1.0
	case BehaviorSuspicious:
		return 1.15
	case BehaviorHostile:
		return 1.3
	default:
		return 1.0
	}
}

// WillTrade returns whether the town will trade with the player.
func (tr *TownReputation) WillTrade() bool {
	// Hostile towns may refuse to trade
	if tr.Behavior == BehaviorHostile {
		return tr.Reputation > 0.1 // Only refuse below 10%
	}
	return true
}

// WillAllowEntry returns whether the town allows entry.
func (tr *TownReputation) WillAllowEntry() bool {
	// Only truly hostile towns (< 5%) deny entry
	return tr.Reputation >= 0.05
}

// MayAttack returns whether the town might attack the player.
func (tr *TownReputation) MayAttack() bool {
	// Hostile towns below 15% reputation may attack
	return tr.Behavior == BehaviorHostile && tr.Reputation < 0.15
}

// AttackChance returns the probability of attack (0-1).
func (tr *TownReputation) AttackChance() float64 {
	if !tr.MayAttack() {
		return 0
	}
	// Lower reputation = higher attack chance
	// At 0% rep = 50% attack chance
	// At 15% rep = 0% attack chance
	return (0.15 - tr.Reputation) / 0.15 * 0.5
}

// TownEvent represents something that happened at a town.
type TownEvent int

const (
	EventVisit TownEvent = iota
	EventSuccessfulTrade
	EventFailedTrade
	EventHelped   // Player helped the town
	EventHarmed   // Player harmed the town
	EventDefended // Player defended against attack
	EventAttacked // Town attacked player
)

// RecordEvent records a town event and updates reputation.
func (tr *TownReputation) RecordEvent(event TownEvent) {
	switch event {
	case EventVisit:
		tr.RecordVisit()
	case EventSuccessfulTrade:
		tr.RecordTrade(true)
	case EventFailedTrade:
		tr.RecordTrade(false)
	case EventHelped:
		tr.ModifyReputation(0.1)
	case EventHarmed:
		tr.ModifyReputation(-0.15)
	case EventDefended:
		tr.ModifyReputation(0.05)
	case EventAttacked:
		tr.ModifyReputation(-0.05) // Attack further damages relationship
	}
}

// BehaviorName returns a genre-appropriate name for the behavior.
func (tr *TownReputation) BehaviorName() string {
	names, ok := behaviorNames[tr.Genre]
	if !ok {
		names = behaviorNames[engine.GenreFantasy]
	}
	name, ok := names[tr.Behavior]
	if !ok {
		return "Unknown"
	}
	return name
}

// ReputationDescription returns a genre-appropriate reputation description.
func (tr *TownReputation) ReputationDescription() string {
	descriptions, ok := reputationDescriptions[tr.Genre]
	if !ok {
		descriptions = reputationDescriptions[engine.GenreFantasy]
	}
	desc, ok := descriptions[tr.Behavior]
	if !ok {
		return "Unknown standing"
	}
	return desc
}

// HostileWarning returns a warning message if town is hostile.
func (tr *TownReputation) HostileWarning() string {
	if tr.Behavior != BehaviorHostile {
		return ""
	}
	warnings, ok := hostileWarnings[tr.Genre]
	if !ok {
		warnings = hostileWarnings[engine.GenreFantasy]
	}
	return warnings
}

var behaviorNames = map[engine.GenreID]map[TownBehavior]string{
	engine.GenreFantasy: {
		BehaviorWelcoming:  "Revered",
		BehaviorFriendly:   "Friendly",
		BehaviorNeutral:    "Neutral",
		BehaviorSuspicious: "Wary",
		BehaviorHostile:    "Hostile",
	},
	engine.GenreScifi: {
		BehaviorWelcoming:  "Allied",
		BehaviorFriendly:   "Cooperative",
		BehaviorNeutral:    "Neutral",
		BehaviorSuspicious: "Monitored",
		BehaviorHostile:    "Enemy",
	},
	engine.GenreHorror: {
		BehaviorWelcoming:  "Trusted",
		BehaviorFriendly:   "Accepted",
		BehaviorNeutral:    "Tolerated",
		BehaviorSuspicious: "Suspected",
		BehaviorHostile:    "Hunted",
	},
	engine.GenreCyberpunk: {
		BehaviorWelcoming:  "VIP",
		BehaviorFriendly:   "Connected",
		BehaviorNeutral:    "Anonymous",
		BehaviorSuspicious: "Flagged",
		BehaviorHostile:    "Blacklisted",
	},
	engine.GenrePostapoc: {
		BehaviorWelcoming:  "Family",
		BehaviorFriendly:   "Ally",
		BehaviorNeutral:    "Outsider",
		BehaviorSuspicious: "Threat",
		BehaviorHostile:    "Enemy",
	},
}

var reputationDescriptions = map[engine.GenreID]map[TownBehavior]string{
	engine.GenreFantasy: {
		BehaviorWelcoming:  "The people here celebrate your arrival",
		BehaviorFriendly:   "You are welcomed warmly",
		BehaviorNeutral:    "The townsfolk regard you with indifference",
		BehaviorSuspicious: "Suspicious eyes follow your every move",
		BehaviorHostile:    "This place is dangerous for you",
	},
	engine.GenreScifi: {
		BehaviorWelcoming:  "You have full station privileges",
		BehaviorFriendly:   "Your docking request is approved",
		BehaviorNeutral:    "Standard protocols apply",
		BehaviorSuspicious: "Your activities are being monitored",
		BehaviorHostile:    "Warning: Security alert active",
	},
	engine.GenreHorror: {
		BehaviorWelcoming:  "They trust you with their secrets",
		BehaviorFriendly:   "They seem to accept you",
		BehaviorNeutral:    "They watch you cautiously",
		BehaviorSuspicious: "Something is wrong here",
		BehaviorHostile:    "They know what you are",
	},
	engine.GenreCyberpunk: {
		BehaviorWelcoming:  "Access: Unrestricted",
		BehaviorFriendly:   "Access: Standard Plus",
		BehaviorNeutral:    "Access: Standard",
		BehaviorSuspicious: "Access: Limited",
		BehaviorHostile:    "Access: DENIED",
	},
	engine.GenrePostapoc: {
		BehaviorWelcoming:  "You're one of us now",
		BehaviorFriendly:   "Come on in, friend",
		BehaviorNeutral:    "State your business",
		BehaviorSuspicious: "We're watching you",
		BehaviorHostile:    "Turn back now",
	},
}

var hostileWarnings = map[engine.GenreID]string{
	engine.GenreFantasy:   "Beware - the guards may attack on sight",
	engine.GenreScifi:     "Warning: Defense systems may engage",
	engine.GenreHorror:    "They're hunting for you here",
	engine.GenreCyberpunk: "Alert: Enforcement units dispatched",
	engine.GenrePostapoc:  "They'll shoot first, ask questions never",
}

// ReputationVocab holds genre-specific reputation vocabulary.
type ReputationVocab struct {
	ReputationLabel string
	BehaviorLabel   string
	FriendlyLabel   string
	HostileLabel    string
	VisitsLabel     string
	TradesLabel     string
}

// GetReputationVocab returns genre-specific reputation vocabulary.
func GetReputationVocab(genre engine.GenreID) *ReputationVocab {
	vocab, ok := reputationVocabs[genre]
	if !ok {
		return reputationVocabs[engine.GenreFantasy]
	}
	return vocab
}

var reputationVocabs = map[engine.GenreID]*ReputationVocab{
	engine.GenreFantasy: {
		ReputationLabel: "Standing",
		BehaviorLabel:   "Attitude",
		FriendlyLabel:   "Allies",
		HostileLabel:    "Enemies",
		VisitsLabel:     "Visits",
		TradesLabel:     "Trades",
	},
	engine.GenreScifi: {
		ReputationLabel: "Reputation Index",
		BehaviorLabel:   "Diplomatic Status",
		FriendlyLabel:   "Allies",
		HostileLabel:    "Hostiles",
		VisitsLabel:     "Dockings",
		TradesLabel:     "Transactions",
	},
	engine.GenreHorror: {
		ReputationLabel: "Trust",
		BehaviorLabel:   "Perception",
		FriendlyLabel:   "Safe Havens",
		HostileLabel:    "Danger Zones",
		VisitsLabel:     "Encounters",
		TradesLabel:     "Exchanges",
	},
	engine.GenreCyberpunk: {
		ReputationLabel: "Rep Score",
		BehaviorLabel:   "Access Level",
		FriendlyLabel:   "Contacts",
		HostileLabel:    "Blacklist",
		VisitsLabel:     "Check-ins",
		TradesLabel:     "Deals",
	},
	engine.GenrePostapoc: {
		ReputationLabel: "Standing",
		BehaviorLabel:   "Relations",
		FriendlyLabel:   "Friends",
		HostileLabel:    "Foes",
		VisitsLabel:     "Stops",
		TradesLabel:     "Barters",
	},
}
