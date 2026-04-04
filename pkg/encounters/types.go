package encounters

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// EncounterType identifies the category of tactical encounter.
type EncounterType int

const (
	// TypeAmbush is a combat-focused surprise encounter.
	TypeAmbush EncounterType = iota
	// TypeNegotiation is a dialogue-based encounter.
	TypeNegotiation
	// TypeRace is a timed chase or escape encounter.
	TypeRace
	// TypeCrisis is an emergency management encounter.
	TypeCrisis
	// TypePuzzle is a logic/problem-solving encounter.
	TypePuzzle
)

// AllEncounterTypes returns all encounter types.
func AllEncounterTypes() []EncounterType {
	return []EncounterType{
		TypeAmbush,
		TypeNegotiation,
		TypeRace,
		TypeCrisis,
		TypePuzzle,
	}
}

// EncounterOutcome represents the result of an encounter resolution.
type EncounterOutcome int

const (
	// OutcomeVictory is a complete success.
	OutcomeVictory EncounterOutcome = iota
	// OutcomePartialSuccess is a mixed result with some gains and losses.
	OutcomePartialSuccess
	// OutcomeRetreat is a successful escape with minimal losses.
	OutcomeRetreat
	// OutcomeDefeat is a failed encounter with significant losses.
	OutcomeDefeat
)

// OutcomeNames returns genre-appropriate names for outcomes.
func OutcomeNames(genre engine.GenreID) map[EncounterOutcome]string {
	names := outcomeNames[genre]
	if names == nil {
		names = outcomeNames[engine.GenreFantasy]
	}
	return names
}

var outcomeNames = map[engine.GenreID]map[EncounterOutcome]string{
	engine.GenreFantasy: {
		OutcomeVictory:        "Triumph",
		OutcomePartialSuccess: "Pyrrhic Victory",
		OutcomeRetreat:        "Tactical Withdrawal",
		OutcomeDefeat:         "Catastrophe",
	},
	engine.GenreScifi: {
		OutcomeVictory:        "Mission Success",
		OutcomePartialSuccess: "Partial Success",
		OutcomeRetreat:        "Emergency Retreat",
		OutcomeDefeat:         "Mission Failure",
	},
	engine.GenreHorror: {
		OutcomeVictory:        "Survived",
		OutcomePartialSuccess: "Barely Escaped",
		OutcomeRetreat:        "Fled",
		OutcomeDefeat:         "Lost",
	},
	engine.GenreCyberpunk: {
		OutcomeVictory:        "Clean Run",
		OutcomePartialSuccess: "Messy Run",
		OutcomeRetreat:        "Bailed Out",
		OutcomeDefeat:         "Flatlined",
	},
	engine.GenrePostapoc: {
		OutcomeVictory:        "Victory",
		OutcomePartialSuccess: "Survived",
		OutcomeRetreat:        "Got Away",
		OutcomeDefeat:         "Wiped Out",
	},
}

// EncounterRole represents a role a crew member can fill during an encounter.
type EncounterRole int

const (
	// RoleFighter handles combat and defense.
	RoleFighter EncounterRole = iota
	// RoleMedic provides healing and support.
	RoleMedic
	// RoleEngineer handles technical challenges.
	RoleEngineer
	// RoleNegotiator handles dialogue and persuasion.
	RoleNegotiator
)

// AllEncounterRoles returns all encounter roles.
func AllEncounterRoles() []EncounterRole {
	return []EncounterRole{
		RoleFighter,
		RoleMedic,
		RoleEngineer,
		RoleNegotiator,
	}
}

// RoleNames returns genre-appropriate names for roles.
func RoleNames(genre engine.GenreID) map[EncounterRole]string {
	names := roleNames[genre]
	if names == nil {
		names = roleNames[engine.GenreFantasy]
	}
	return names
}

var roleNames = map[engine.GenreID]map[EncounterRole]string{
	engine.GenreFantasy: {
		RoleFighter:    "Defender",
		RoleMedic:      "Healer",
		RoleEngineer:   "Craftsman",
		RoleNegotiator: "Diplomat",
	},
	engine.GenreScifi: {
		RoleFighter:    "Security",
		RoleMedic:      "Medical",
		RoleEngineer:   "Engineering",
		RoleNegotiator: "Communications",
	},
	engine.GenreHorror: {
		RoleFighter:    "Fighter",
		RoleMedic:      "Medic",
		RoleEngineer:   "Mechanic",
		RoleNegotiator: "Talker",
	},
	engine.GenreCyberpunk: {
		RoleFighter:    "Solo",
		RoleMedic:      "Medtech",
		RoleEngineer:   "Tech",
		RoleNegotiator: "Face",
	},
	engine.GenrePostapoc: {
		RoleFighter:    "Muscle",
		RoleMedic:      "Doc",
		RoleEngineer:   "Wrench",
		RoleNegotiator: "Smooth-talker",
	},
}

// EncounterResult holds the complete result of an encounter resolution.
type EncounterResult struct {
	Outcome       EncounterOutcome
	Description   string
	FoodDelta     float64
	WaterDelta    float64
	FuelDelta     float64
	MedicineDelta float64
	MoraleDelta   float64
	CurrencyDelta float64
	CrewDamage    map[int]float64 // Crew member ID -> damage taken
	VesselDamage  float64
	TurnsElapsed  int
	SkillExpGains map[int]float64 // Crew member ID -> exp gained
}

// NewEncounterResult creates a new encounter result.
func NewEncounterResult(outcome EncounterOutcome) *EncounterResult {
	return &EncounterResult{
		Outcome:       outcome,
		CrewDamage:    make(map[int]float64),
		SkillExpGains: make(map[int]float64),
	}
}

// TypeName returns the genre-appropriate name for an encounter type.
func TypeName(t EncounterType, genre engine.GenreID) string {
	names := typeNames[genre]
	if names == nil {
		names = typeNames[engine.GenreFantasy]
	}
	return names[t]
}

var typeNames = map[engine.GenreID]map[EncounterType]string{
	engine.GenreFantasy: {
		TypeAmbush:      "Ambush",
		TypeNegotiation: "Parley",
		TypeRace:        "Chase",
		TypeCrisis:      "Peril",
		TypePuzzle:      "Riddle",
	},
	engine.GenreScifi: {
		TypeAmbush:      "Hostile Contact",
		TypeNegotiation: "Hailing",
		TypeRace:        "Pursuit",
		TypeCrisis:      "Emergency",
		TypePuzzle:      "Anomaly",
	},
	engine.GenreHorror: {
		TypeAmbush:      "Attack",
		TypeNegotiation: "Standoff",
		TypeRace:        "Flight",
		TypeCrisis:      "Crisis",
		TypePuzzle:      "Mystery",
	},
	engine.GenreCyberpunk: {
		TypeAmbush:      "Ambush",
		TypeNegotiation: "Deal",
		TypeRace:        "Chase",
		TypeCrisis:      "Situation",
		TypePuzzle:      "Hack",
	},
	engine.GenrePostapoc: {
		TypeAmbush:      "Raid",
		TypeNegotiation: "Parley",
		TypeRace:        "Escape",
		TypeCrisis:      "Disaster",
		TypePuzzle:      "Salvage",
	},
}
