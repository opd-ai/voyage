//go:build headless

package ux

// splitWords splits text into words.
// This is a headless stub providing the same pure-logic function.
func splitWords(text string) []string {
	var words []string
	var current string

	for _, c := range text {
		if c == ' ' || c == '\n' {
			if current != "" {
				words = append(words, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		words = append(words, current)
	}

	return words
}

// GameStats holds end-of-game statistics.
// This is a headless stub providing the same type.
type GameStats struct {
	DaysTraveled     int
	DistanceTraveled int
	CrewLost         int
	EventsResolved   int
	Victory          bool
}
