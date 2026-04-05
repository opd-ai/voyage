//go:build headless

package rendering

import (
	"image/color"
)

// AnimatedTile represents a tile with multiple animation frames (headless stub).
type AnimatedTile struct {
	FrameCount   int
	FrameTime    float64
	Loop         bool
	currentFrame int
	elapsed      float64
}

// NewAnimatedTile creates a new animated tile stub.
func NewAnimatedTile(frameCount int, frameTime float64, loop bool) *AnimatedTile {
	return &AnimatedTile{
		FrameCount:   frameCount,
		FrameTime:    frameTime,
		Loop:         loop,
		currentFrame: 0,
		elapsed:      0,
	}
}

// Update advances the animation by the given delta time.
func (at *AnimatedTile) Update(dt float64) {
	if at.FrameCount <= 1 {
		return
	}
	at.elapsed += dt
	if at.elapsed >= at.FrameTime {
		at.elapsed -= at.FrameTime
		at.currentFrame++
		if at.currentFrame >= at.FrameCount {
			if at.Loop {
				at.currentFrame = 0
			} else {
				at.currentFrame = at.FrameCount - 1
			}
		}
	}
}

// CurrentFrameIndex returns the current animation frame index.
func (at *AnimatedTile) CurrentFrameIndex() int {
	return at.currentFrame
}

// Reset resets the animation to the first frame.
func (at *AnimatedTile) Reset() {
	at.currentFrame = 0
	at.elapsed = 0
}

// AnimationType defines the kind of animated tile.
type AnimationType int

const (
	AnimationWater AnimationType = iota
	AnimationGrass
	AnimationFire
)

// AnimatedTileGenerator creates animated overworld tiles (headless stub).
type AnimatedTileGenerator struct {
	tileSize int
}

// NewAnimatedTileGenerator creates a new animated tile generator stub.
func NewAnimatedTileGenerator(masterSeed int64, tileSize int) *AnimatedTileGenerator {
	return &AnimatedTileGenerator{
		tileSize: tileSize,
	}
}

// GenerateAnimatedTile creates an animated tile stub of the specified type.
func (atg *AnimatedTileGenerator) GenerateAnimatedTile(animType AnimationType, baseColor, accentColor color.Color) *AnimatedTile {
	var frameCount int
	var frameTime float64

	switch animType {
	case AnimationWater:
		frameCount = 4
		frameTime = 0.2
	case AnimationGrass:
		frameCount = 4
		frameTime = 0.15
	case AnimationFire:
		frameCount = 4
		frameTime = 0.1
	default:
		frameCount = 4
		frameTime = 0.2
	}

	return NewAnimatedTile(frameCount, frameTime, true)
}

// sinApprox provides a fast sine approximation.
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

// Note: clampFloat is defined in lighting_core.go which is available in all builds.
