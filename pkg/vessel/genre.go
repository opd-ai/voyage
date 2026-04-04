package vessel

import "github.com/opd-ai/voyage/pkg/engine"

// VesselTypeName returns the genre-specific name for a vessel type.
// This is an alias for VesselName for consistency with other packages.
func VesselTypeName(vt VesselType, genre engine.GenreID) string {
	return VesselName(vt, genre)
}

// ConditionStatus represents the overall vessel condition.
type ConditionStatus int

const (
	// ConditionPristine is above 90% integrity.
	ConditionPristine ConditionStatus = iota
	// ConditionGood is 75-90% integrity.
	ConditionGood
	// ConditionDamaged is 50-74% integrity.
	ConditionDamaged
	// ConditionCritical is 25-49% integrity.
	ConditionCritical
	// ConditionDestroyed is below 25% or at 0% integrity.
	ConditionDestroyed
)

// GetConditionStatus returns the condition status for a vessel.
func GetConditionStatus(v *Vessel) ConditionStatus {
	ratio := v.IntegrityRatio()
	switch {
	case ratio >= 0.9:
		return ConditionPristine
	case ratio >= 0.75:
		return ConditionGood
	case ratio >= 0.5:
		return ConditionDamaged
	case ratio >= 0.25:
		return ConditionCritical
	default:
		return ConditionDestroyed
	}
}

// ConditionName returns a human-readable condition name.
func ConditionName(cs ConditionStatus) string {
	switch cs {
	case ConditionPristine:
		return "Pristine"
	case ConditionGood:
		return "Good"
	case ConditionDamaged:
		return "Damaged"
	case ConditionCritical:
		return "Critical"
	case ConditionDestroyed:
		return "Destroyed"
	default:
		return "Unknown"
	}
}

// GetGenreSpecificConditionDesc returns a genre-flavored condition description.
func GetGenreSpecificConditionDesc(cs ConditionStatus, genre engine.GenreID) string {
	descs := conditionDescriptions[genre]
	if descs == nil {
		descs = conditionDescriptions[engine.GenreFantasy]
	}
	return descs[cs]
}

var conditionDescriptions = map[engine.GenreID]map[ConditionStatus]string{
	engine.GenreFantasy: {
		ConditionPristine:  "Your vessel gleams like new",
		ConditionGood:      "Minor wear and tear visible",
		ConditionDamaged:   "The wagon creaks ominously",
		ConditionCritical:  "Held together by prayer",
		ConditionDestroyed: "Reduced to splinters",
	},
	engine.GenreScifi: {
		ConditionPristine:  "All systems nominal",
		ConditionGood:      "Minor hull scoring detected",
		ConditionDamaged:   "Warning lights flashing",
		ConditionCritical:  "Hull breach imminent",
		ConditionDestroyed: "Catastrophic decompression",
	},
	engine.GenreHorror: {
		ConditionPristine:  "Running smoothly... for now",
		ConditionGood:      "A few dents and scratches",
		ConditionDamaged:   "Smoke from under the hood",
		ConditionCritical:  "Engine coughing blood",
		ConditionDestroyed: "A burning wreck",
	},
	engine.GenreCyberpunk: {
		ConditionPristine:  "Chrome shining bright",
		ConditionGood:      "Minor cosmetic damage",
		ConditionDamaged:   "Systems glitching",
		ConditionCritical:  "Red-lining everything",
		ConditionDestroyed: "Flatlined",
	},
	engine.GenrePostapoc: {
		ConditionPristine:  "Good as wasteland gets",
		ConditionGood:      "Some rust, still runs",
		ConditionDamaged:   "Rattling like bones",
		ConditionCritical:  "Running on fumes and hope",
		ConditionDestroyed: "Just scrap now",
	},
}
