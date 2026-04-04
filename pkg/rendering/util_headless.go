//go:build headless

package rendering

// clamp restricts a value to a range.
// This is a headless stub that provides the same pure-logic function.
func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
