package events

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// EventCategory identifies the type of event.
type EventCategory int

const (
	// CategoryWeather represents weather-related events.
	CategoryWeather EventCategory = iota
	// CategoryEncounter represents meeting NPCs or creatures.
	CategoryEncounter
	// CategoryDiscovery represents finding locations or items.
	CategoryDiscovery
	// CategoryHardship represents negative events like illness or accidents.
	CategoryHardship
	// CategoryWindfall represents positive random events.
	CategoryWindfall
)

// AllEventCategories returns all event categories.
func AllEventCategories() []EventCategory {
	return []EventCategory{
		CategoryWeather,
		CategoryEncounter,
		CategoryDiscovery,
		CategoryHardship,
		CategoryWindfall,
	}
}

// EventOutcome represents the result of a choice.
type EventOutcome struct {
	Description   string
	FoodDelta     float64
	WaterDelta    float64
	FuelDelta     float64
	MedicineDelta float64
	MoraleDelta   float64
	CurrencyDelta float64
	CrewDamage    float64
	VesselDamage  float64
	TimeAdvance   int
}

// Choice represents an option the player can select.
type Choice struct {
	ID           int
	Text         string
	Outcome      EventOutcome
	RequireSkill string // Optional skill requirement
}

// Event represents a procedural event with choices.
type Event struct {
	ID          int
	Category    EventCategory
	Title       string
	Description string
	Choices     []Choice
	Genre       engine.GenreID
}

// NewEvent creates a new event.
func NewEvent(id int, category EventCategory, title, description string, genre engine.GenreID) *Event {
	return &Event{
		ID:          id,
		Category:    category,
		Title:       title,
		Description: description,
		Choices:     make([]Choice, 0),
		Genre:       genre,
	}
}

// AddChoice adds a choice to the event.
func (e *Event) AddChoice(text string, outcome EventOutcome) {
	e.Choices = append(e.Choices, Choice{
		ID:      len(e.Choices) + 1,
		Text:    text,
		Outcome: outcome,
	})
}

// AddSkillChoice adds a choice requiring a skill.
func (e *Event) AddSkillChoice(text string, outcome EventOutcome, skill string) {
	e.Choices = append(e.Choices, Choice{
		ID:           len(e.Choices) + 1,
		Text:         text,
		Outcome:      outcome,
		RequireSkill: skill,
	})
}

// GetChoice returns a choice by ID.
func (e *Event) GetChoice(id int) *Choice {
	for i := range e.Choices {
		if e.Choices[i].ID == id {
			return &e.Choices[i]
		}
	}
	return nil
}

// CategoryName returns a human-readable name for the category.
func CategoryName(cat EventCategory, genre engine.GenreID) string {
	names, ok := categoryNames[genre]
	if !ok {
		names = categoryNames[engine.GenreFantasy]
	}
	return names[cat]
}

var categoryNames = map[engine.GenreID]map[EventCategory]string{
	engine.GenreFantasy: {
		CategoryWeather:   "Weather",
		CategoryEncounter: "Encounter",
		CategoryDiscovery: "Discovery",
		CategoryHardship:  "Hardship",
		CategoryWindfall:  "Fortune",
	},
	engine.GenreScifi: {
		CategoryWeather:   "Space Weather",
		CategoryEncounter: "Contact",
		CategoryDiscovery: "Scan Result",
		CategoryHardship:  "Malfunction",
		CategoryWindfall:  "Lucky Find",
	},
	engine.GenreHorror: {
		CategoryWeather:   "Weather",
		CategoryEncounter: "Encounter",
		CategoryDiscovery: "Discovery",
		CategoryHardship:  "Crisis",
		CategoryWindfall:  "Relief",
	},
	engine.GenreCyberpunk: {
		CategoryWeather:   "Environment",
		CategoryEncounter: "Contact",
		CategoryDiscovery: "Intel",
		CategoryHardship:  "Complication",
		CategoryWindfall:  "Score",
	},
	engine.GenrePostapoc: {
		CategoryWeather:   "Weather",
		CategoryEncounter: "Sighting",
		CategoryDiscovery: "Salvage",
		CategoryHardship:  "Crisis",
		CategoryWindfall:  "Lucky Break",
	},
}
