package events

import "github.com/opd-ai/voyage/pkg/engine"

// SetGenre is a convenience method to satisfy GenreSwitcher interface.
func (e *Event) SetGenre(genre engine.GenreID) {
	e.Genre = genre
}

// GetGenreDescription returns a genre-specific category description.
func GetGenreDescription(cat EventCategory, genre engine.GenreID) string {
	desc := categoryDescriptions[genre]
	if desc == nil {
		desc = categoryDescriptions[engine.GenreFantasy]
	}
	return desc[cat]
}

var categoryDescriptions = map[engine.GenreID]map[EventCategory]string{
	engine.GenreFantasy: {
		CategoryWeather:   "The elements test your resolve",
		CategoryEncounter: "You meet others on the road",
		CategoryDiscovery: "Something catches your eye",
		CategoryHardship:  "Misfortune strikes the caravan",
		CategoryWindfall:  "Fortune smiles upon you",
	},
	engine.GenreScifi: {
		CategoryWeather:   "Space conditions challenge navigation",
		CategoryEncounter: "Contact established with another vessel",
		CategoryDiscovery: "Scanners detect something interesting",
		CategoryHardship:  "Ship systems report anomalies",
		CategoryWindfall:  "A favorable opportunity presents itself",
	},
	engine.GenreHorror: {
		CategoryWeather:   "The environment turns hostile",
		CategoryEncounter: "You're not alone out here",
		CategoryDiscovery: "You find evidence of those who came before",
		CategoryHardship:  "Things go from bad to worse",
		CategoryWindfall:  "A brief respite from the horror",
	},
	engine.GenreCyberpunk: {
		CategoryWeather:   "Environmental hazards detected",
		CategoryEncounter: "Someone wants your attention",
		CategoryDiscovery: "Intel opportunity identified",
		CategoryHardship:  "Complications arise",
		CategoryWindfall:  "A lucrative opportunity emerges",
	},
	engine.GenrePostapoc: {
		CategoryWeather:   "The wasteland is unforgiving",
		CategoryEncounter: "Others traverse the wastes",
		CategoryDiscovery: "Remnants of the old world",
		CategoryHardship:  "Survival becomes more difficult",
		CategoryWindfall:  "A rare piece of good luck",
	},
}

// OutcomeSeverity calculates how good or bad an outcome is.
// Positive = good, negative = bad, 0 = neutral.
func OutcomeSeverity(o *EventOutcome) float64 {
	score := 0.0

	// Resources
	score += o.FoodDelta * 0.5
	score += o.WaterDelta * 0.5
	score += o.FuelDelta * 0.3
	score += o.MedicineDelta * 1.0
	score += o.MoraleDelta * 0.8
	score += o.CurrencyDelta * 0.2

	// Damage is always bad
	score -= o.CrewDamage * 2.0
	score -= o.VesselDamage * 1.5

	// Time lost is moderately bad
	score -= float64(o.TimeAdvance) * 3.0

	return score
}

// OutcomeRisk returns a risk assessment of a choice.
func OutcomeRisk(o *EventOutcome) string {
	severity := OutcomeSeverity(o)

	switch {
	case severity >= 20:
		return "Highly Favorable"
	case severity >= 5:
		return "Favorable"
	case severity >= -5:
		return "Neutral"
	case severity >= -20:
		return "Risky"
	default:
		return "Dangerous"
	}
}
