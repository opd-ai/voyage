//go:build headless

package rendering

import (
	"image/color"
)

// LandmarkIconType identifies the type of landmark icon.
type LandmarkIconType int

const (
	LandmarkIconTown LandmarkIconType = iota
	LandmarkIconOutpost
	LandmarkIconRuins
	LandmarkIconShrine
	LandmarkIconOrigin
	LandmarkIconDestination
)

// AnimatedLandmarkIcon represents an animated landmark icon (headless stub).
type AnimatedLandmarkIcon struct {
	FrameCount   int
	FrameTime    float64
	currentFrame int
	elapsed      float64
}

// NewAnimatedLandmarkIcon creates a new animated landmark icon stub.
func NewAnimatedLandmarkIcon(frameCount int, frameTime float64) *AnimatedLandmarkIcon {
	return &AnimatedLandmarkIcon{
		FrameCount:   frameCount,
		FrameTime:    frameTime,
		currentFrame: 0,
		elapsed:      0,
	}
}

// Update advances the landmark animation by the given delta time.
func (ali *AnimatedLandmarkIcon) Update(dt float64) {
	if ali.FrameCount <= 1 {
		return
	}
	ali.elapsed += dt
	if ali.elapsed >= ali.FrameTime {
		ali.elapsed -= ali.FrameTime
		ali.currentFrame++
		if ali.currentFrame >= ali.FrameCount {
			ali.currentFrame = 0
		}
	}
}

// CurrentFrameIndex returns the current animation frame index.
func (ali *AnimatedLandmarkIcon) CurrentFrameIndex() int {
	return ali.currentFrame
}

// Reset resets the animation to the first frame.
func (ali *AnimatedLandmarkIcon) Reset() {
	ali.currentFrame = 0
	ali.elapsed = 0
}

// LandmarkIconGenerator creates animated landmark icons (headless stub).
type LandmarkIconGenerator struct {
	iconSize int
}

// NewLandmarkIconGenerator creates a new landmark icon generator stub.
func NewLandmarkIconGenerator(masterSeed int64, iconSize int) *LandmarkIconGenerator {
	return &LandmarkIconGenerator{
		iconSize: iconSize,
	}
}

// GenerateAnimatedIcon creates an animated icon stub for the specified landmark type.
func (lig *LandmarkIconGenerator) GenerateAnimatedIcon(iconType LandmarkIconType, primaryColor, secondaryColor color.Color) *AnimatedLandmarkIcon {
	var frameCount int
	var frameTime float64

	switch iconType {
	case LandmarkIconRuins:
		frameCount = 6
		frameTime = 0.15
	case LandmarkIconOutpost:
		frameCount = 4
		frameTime = 0.3
	case LandmarkIconTown:
		frameCount = 4
		frameTime = 0.25
	case LandmarkIconShrine:
		frameCount = 4
		frameTime = 0.2
	case LandmarkIconOrigin:
		frameCount = 4
		frameTime = 0.2
	case LandmarkIconDestination:
		frameCount = 6
		frameTime = 0.15
	default:
		frameCount = 4
		frameTime = 0.25
	}

	return NewAnimatedLandmarkIcon(frameCount, frameTime)
}
