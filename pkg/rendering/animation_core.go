package rendering

// sinApprox provides a fast sine approximation using Taylor series.
// This avoids importing math for simple animation calculations.
func sinApprox(x float64) float64 {
	const pi = 3.14159265358979323846
	const twoPi = 6.283185307179586

	// Normalize x to [0, 2π)
	x = x - float64(int(x/twoPi))*twoPi
	if x < 0 {
		x += twoPi
	}

	// Determine sign and map to [0, π]
	sign := 1.0
	if x > pi {
		x -= pi
		sign = -1.0
	}

	// Map to [-π/2, π/2] for better Taylor accuracy
	if x > pi/2 {
		x = pi - x
	}

	// Taylor series: sin(x) ≈ x - x³/6 + x⁵/120 - x⁷/5040
	x2 := x * x
	result := x * (1 - x2/6*(1-x2/20*(1-x2/42)))
	return sign * result
}

// lerp performs linear interpolation between a and b.
func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

// clampFloat restricts a float64 value to a range.
func clampFloat(v, minVal, maxVal float64) float64 {
	if v < minVal {
		return minVal
	}
	if v > maxVal {
		return maxVal
	}
	return v
}
