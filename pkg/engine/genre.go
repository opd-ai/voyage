package engine

// GenreID identifies one of the five supported genre themes.
type GenreID string

const (
	// GenreFantasy represents the enchanted realm / Silk Road setting.
	GenreFantasy GenreID = "fantasy"
	// GenreScifi represents the deep space / star lanes setting.
	GenreScifi GenreID = "scifi"
	// GenreHorror represents the zombie apocalypse wasteland setting.
	GenreHorror GenreID = "horror"
	// GenreCyberpunk represents the megacity sprawl setting.
	GenreCyberpunk GenreID = "cyberpunk"
	// GenrePostapoc represents the irradiated dust-bowl wastes setting.
	GenrePostapoc GenreID = "postapoc"
)

// GenreSwitcher is implemented by all Systems to enable genre switching at runtime.
// When SetGenre is called, the system updates its thematic presentation (vocabulary,
// palettes, sound presets, etc.) to match the specified genre.
type GenreSwitcher interface {
	SetGenre(genreID GenreID)
}

// AllGenres returns a slice of all supported genre IDs.
func AllGenres() []GenreID {
	return []GenreID{
		GenreFantasy,
		GenreScifi,
		GenreHorror,
		GenreCyberpunk,
		GenrePostapoc,
	}
}

// IsValidGenre checks if the given string is a valid GenreID.
func IsValidGenre(id string) bool {
	switch GenreID(id) {
	case GenreFantasy, GenreScifi, GenreHorror, GenreCyberpunk, GenrePostapoc:
		return true
	default:
		return false
	}
}
