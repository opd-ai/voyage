//go:build headless

package rendering

import (
	"image/color"
)

// PortraitAnimState represents the animation state of a crew portrait.
type PortraitAnimState int

const (
	PortraitIdle PortraitAnimState = iota
	PortraitHurt
	PortraitDeath
)

// AnimatedPortrait represents a crew member portrait with animation states (headless stub).
type AnimatedPortrait struct {
	IdleFrameCount  int
	HurtFrameCount  int
	DeathFrameCount int
	FrameTime       float64

	state        PortraitAnimState
	currentFrame int
	elapsed      float64
}

// NewAnimatedPortrait creates a new animated portrait stub.
func NewAnimatedPortrait(idleCount, hurtCount, deathCount int, frameTime float64) *AnimatedPortrait {
	return &AnimatedPortrait{
		IdleFrameCount:  idleCount,
		HurtFrameCount:  hurtCount,
		DeathFrameCount: deathCount,
		FrameTime:       frameTime,
		state:           PortraitIdle,
		currentFrame:    0,
		elapsed:         0,
	}
}

// SetState changes the portrait animation state.
func (ap *AnimatedPortrait) SetState(state PortraitAnimState) {
	if ap.state != state {
		ap.state = state
		ap.currentFrame = 0
		ap.elapsed = 0
	}
}

// State returns the current animation state.
func (ap *AnimatedPortrait) State() PortraitAnimState {
	return ap.state
}

// Update advances the portrait animation by the given delta time.
func (ap *AnimatedPortrait) Update(dt float64) {
	frameCount := ap.currentFrameCount()
	if frameCount <= 1 {
		return
	}

	ap.elapsed += dt
	if ap.elapsed >= ap.FrameTime {
		ap.elapsed -= ap.FrameTime
		ap.currentFrame++

		if ap.state == PortraitDeath {
			if ap.currentFrame >= frameCount {
				ap.currentFrame = frameCount - 1
			}
		} else {
			if ap.currentFrame >= frameCount {
				ap.currentFrame = 0
			}
		}
	}
}

// CurrentFrameIndex returns the current animation frame index.
func (ap *AnimatedPortrait) CurrentFrameIndex() int {
	return ap.currentFrame
}

// currentFrameCount returns the frame count for the current state.
func (ap *AnimatedPortrait) currentFrameCount() int {
	switch ap.state {
	case PortraitHurt:
		return ap.HurtFrameCount
	case PortraitDeath:
		return ap.DeathFrameCount
	default:
		return ap.IdleFrameCount
	}
}

// Reset resets the animation to the first frame of the current state.
func (ap *AnimatedPortrait) Reset() {
	ap.currentFrame = 0
	ap.elapsed = 0
}

// PortraitGenerator creates animated crew member portraits (headless stub).
type PortraitGenerator struct {
	portraitSize int
}

// NewPortraitGenerator creates a new portrait generator stub.
func NewPortraitGenerator(masterSeed int64, portraitSize int) *PortraitGenerator {
	return &PortraitGenerator{
		portraitSize: portraitSize,
	}
}

// GenerateAnimatedPortrait creates an animated portrait stub with all animation states.
func (pg *PortraitGenerator) GenerateAnimatedPortrait(primaryColor, secondaryColor, skinColor color.Color) *AnimatedPortrait {
	const idleFrameCount = 4
	const hurtFrameCount = 4
	const deathFrameCount = 8

	return NewAnimatedPortrait(idleFrameCount, hurtFrameCount, deathFrameCount, 0.25)
}
