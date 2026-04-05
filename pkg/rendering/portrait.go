//go:build !headless

package rendering

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// PortraitAnimState represents the animation state of a crew portrait.
type PortraitAnimState int

const (
	PortraitIdle PortraitAnimState = iota
	PortraitHurt
	PortraitDeath
)

// AnimatedPortrait represents a crew member portrait with animation frames.
type AnimatedPortrait struct {
	IdleFrames  []*ebiten.Image
	HurtFrames  []*ebiten.Image
	DeathFrames []*ebiten.Image
	FrameTime   float64

	state        PortraitAnimState
	currentFrame int
	elapsed      float64
}

// NewAnimatedPortrait creates a new animated portrait from pre-generated frames.
func NewAnimatedPortrait(idle, hurt, death []*ebiten.Image, frameTime float64) *AnimatedPortrait {
	return &AnimatedPortrait{
		IdleFrames:   idle,
		HurtFrames:   hurt,
		DeathFrames:  death,
		FrameTime:    frameTime,
		state:        PortraitIdle,
		currentFrame: 0,
		elapsed:      0,
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
	frames := ap.currentFrameSet()
	if len(frames) <= 1 {
		return
	}
	ap.advanceFrame(dt, len(frames))
}

// advanceFrame increments the frame counter based on elapsed time.
func (ap *AnimatedPortrait) advanceFrame(dt float64, frameCount int) {
	ap.elapsed += dt
	if ap.elapsed < ap.FrameTime {
		return
	}
	ap.elapsed -= ap.FrameTime
	ap.currentFrame++
	ap.wrapFrame(frameCount)
}

// wrapFrame handles frame wrapping based on animation state.
func (ap *AnimatedPortrait) wrapFrame(frameCount int) {
	if ap.currentFrame < frameCount {
		return
	}
	if ap.state == PortraitDeath {
		ap.currentFrame = frameCount - 1
	} else {
		ap.currentFrame = 0
	}
}

// CurrentFrame returns the current animation frame image.
func (ap *AnimatedPortrait) CurrentFrame() *ebiten.Image {
	frames := ap.currentFrameSet()
	if len(frames) == 0 {
		return nil
	}
	return frames[ap.currentFrame]
}

// currentFrameSet returns the frame set for the current state.
func (ap *AnimatedPortrait) currentFrameSet() []*ebiten.Image {
	switch ap.state {
	case PortraitHurt:
		return ap.HurtFrames
	case PortraitDeath:
		return ap.DeathFrames
	default:
		return ap.IdleFrames
	}
}

// Reset resets the animation to the first frame of the current state.
func (ap *AnimatedPortrait) Reset() {
	ap.currentFrame = 0
	ap.elapsed = 0
}

// PortraitGenerator creates animated crew member portraits.
type PortraitGenerator struct {
	gen          *seed.Generator
	portraitSize int
}

// NewPortraitGenerator creates a new portrait generator.
func NewPortraitGenerator(masterSeed int64, portraitSize int) *PortraitGenerator {
	return &PortraitGenerator{
		gen:          seed.NewGenerator(masterSeed, "portrait"),
		portraitSize: portraitSize,
	}
}

// GenerateAnimatedPortrait creates an animated portrait with all animation states.
func (pg *PortraitGenerator) GenerateAnimatedPortrait(primaryColor, secondaryColor, skinColor color.Color) *AnimatedPortrait {
	basePortrait := pg.generateBasePortrait(primaryColor, secondaryColor, skinColor)

	idleFrames := pg.generateIdleAnimation(basePortrait)
	hurtFrames := pg.generateHurtAnimation(basePortrait)
	deathFrames := pg.generateDeathAnimation(basePortrait)

	return NewAnimatedPortrait(idleFrames, hurtFrames, deathFrames, 0.25)
}

// generateBasePortrait creates the base portrait image.
func (pg *PortraitGenerator) generateBasePortrait(primaryColor, secondaryColor, skinColor color.Color) *ebiten.Image {
	size := pg.portraitSize
	img := ebiten.NewImage(size, size)

	faceTop := size / 6
	faceBottom := size * 5 / 6
	faceLeft := size / 4
	faceRight := size * 3 / 4

	pg.drawFaceRegion(img, faceTop, faceBottom, faceLeft, faceRight, skinColor)
	pg.drawBodyRegion(img, faceBottom, primaryColor)
	pg.drawHairRegion(img, faceTop, secondaryColor)
	pg.drawEyes(img, faceTop, faceBottom, faceLeft, faceRight)

	return img
}

// drawFaceRegion draws the face area of a portrait.
func (pg *PortraitGenerator) drawFaceRegion(img *ebiten.Image, top, bottom, left, right int, skinColor color.Color) {
	for y := top; y < bottom; y++ {
		for x := left; x < right; x++ {
			img.Set(x, y, skinColor)
		}
	}
}

// drawBodyRegion draws the body area of a portrait.
func (pg *PortraitGenerator) drawBodyRegion(img *ebiten.Image, bodyTop int, primaryColor color.Color) {
	size := pg.portraitSize
	for y := bodyTop; y < size; y++ {
		for x := size / 6; x < size*5/6; x++ {
			img.Set(x, y, primaryColor)
		}
	}
}

// drawHairRegion draws the hair area with variation.
func (pg *PortraitGenerator) drawHairRegion(img *ebiten.Image, faceTop int, secondaryColor color.Color) {
	size := pg.portraitSize
	hairBottom := faceTop + size/8
	for y := 0; y < hairBottom; y++ {
		hairWidth := size/4 + (size/4)*y/hairBottom
		for x := size/2 - hairWidth; x < size/2+hairWidth; x++ {
			if x >= 0 && x < size && pg.gen.Chance(0.8) {
				img.Set(x, y, secondaryColor)
			}
		}
	}
}

// drawEyes draws both eyes on a portrait.
func (pg *PortraitGenerator) drawEyes(img *ebiten.Image, faceTop, faceBottom, faceLeft, faceRight int) {
	eyeY := faceTop + (faceBottom-faceTop)/3
	leftEyeX := faceLeft + (faceRight-faceLeft)/4
	rightEyeX := faceRight - (faceRight-faceLeft)/4
	eyeColor := color.RGBA{30, 30, 30, 255}
	pg.drawEye(img, leftEyeX, eyeY, eyeColor)
	pg.drawEye(img, rightEyeX, eyeY, eyeColor)
}

// drawEye draws a simple eye at the given position.
func (pg *PortraitGenerator) drawEye(img *ebiten.Image, x, y int, c color.Color) {
	eyeSize := pg.portraitSize / 10
	if eyeSize < 2 {
		eyeSize = 2
	}
	pg.drawFilledCircleAt(img, x, y, eyeSize/2, c)
}

// drawFilledCircleAt draws a filled circle at the given center with specified radius.
func (pg *PortraitGenerator) drawFilledCircleAt(img *ebiten.Image, centerX, centerY, radius int, c color.Color) {
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			if dx*dx+dy*dy <= radius*radius {
				px, py := centerX+dx, centerY+dy
				if px >= 0 && px < pg.portraitSize && py >= 0 && py < pg.portraitSize {
					img.Set(px, py, c)
				}
			}
		}
	}
}

// generateIdleAnimation creates idle breathing animation frames.
func (pg *PortraitGenerator) generateIdleAnimation(base *ebiten.Image) []*ebiten.Image {
	const frameCount = 4
	frames := make([]*ebiten.Image, frameCount)
	size := pg.portraitSize

	for f := 0; f < frameCount; f++ {
		frame := ebiten.NewImage(size, size)
		frame.DrawImage(base, nil)

		breathOffset := 0
		if f == 1 || f == 2 {
			breathOffset = 1
		}

		if breathOffset > 0 {
			pg.applyBreathingEffect(frame, breathOffset)
		}

		frames[f] = frame
	}

	return frames
}

// applyBreathingEffect applies a subtle vertical shift for breathing.
func (pg *PortraitGenerator) applyBreathingEffect(img *ebiten.Image, offset int) {
	size := pg.portraitSize
	bodyStart := size * 5 / 6

	if offset <= 0 || bodyStart+offset >= size {
		return
	}

	tempImg := pg.copyBodyRegion(img, size, bodyStart, offset)
	pg.clearBodyRegion(img, size, bodyStart)
	pg.drawShiftedBody(img, tempImg, bodyStart, offset)
}

// copyBodyRegion copies the body portion of the image to a temporary buffer.
func (pg *PortraitGenerator) copyBodyRegion(img *ebiten.Image, size, bodyStart, offset int) *ebiten.Image {
	tempImg := ebiten.NewImage(size, size-bodyStart)
	for y := bodyStart; y < size-offset; y++ {
		for x := 0; x < size; x++ {
			c := img.At(x, y)
			tempImg.Set(x, y-bodyStart, c)
		}
	}
	return tempImg
}

// clearBodyRegion clears the body area to transparent.
func (pg *PortraitGenerator) clearBodyRegion(img *ebiten.Image, size, bodyStart int) {
	for y := bodyStart; y < size; y++ {
		for x := 0; x < size; x++ {
			img.Set(x, y, color.Transparent)
		}
	}
}

// drawShiftedBody draws the body with the breathing offset applied.
func (pg *PortraitGenerator) drawShiftedBody(img, tempImg *ebiten.Image, bodyStart, offset int) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, float64(bodyStart+offset))
	img.DrawImage(tempImg, op)
}

// generateHurtAnimation creates hurt flinch animation frames.
func (pg *PortraitGenerator) generateHurtAnimation(base *ebiten.Image) []*ebiten.Image {
	const frameCount = 4
	frames := make([]*ebiten.Image, frameCount)
	size := pg.portraitSize

	hurtTint := color.RGBA{200, 100, 100, 255}

	for f := 0; f < frameCount; f++ {
		frame := ebiten.NewImage(size, size)

		// Horizontal shake offsets
		shakeOffsets := []int{0, -2, 2, 0}
		shakeX := shakeOffsets[f]

		// Draw with shake and tint
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(shakeX), 0)

		// Apply red tint intensity based on frame
		tintIntensity := []float64{0.0, 0.3, 0.2, 0.0}
		if tintIntensity[f] > 0 {
			op.ColorScale.Scale(
				float32(1.0-tintIntensity[f]+tintIntensity[f]*float64(hurtTint.R)/255),
				float32(1.0-tintIntensity[f]+tintIntensity[f]*float64(hurtTint.G)/255),
				float32(1.0-tintIntensity[f]+tintIntensity[f]*float64(hurtTint.B)/255),
				1.0,
			)
		}

		frame.DrawImage(base, op)
		frames[f] = frame
	}

	return frames
}

// generateDeathAnimation creates death fade animation frames.
func (pg *PortraitGenerator) generateDeathAnimation(base *ebiten.Image) []*ebiten.Image {
	const frameCount = 8
	frames := make([]*ebiten.Image, frameCount)
	size := pg.portraitSize

	for f := 0; f < frameCount; f++ {
		frame := ebiten.NewImage(size, size)

		// Calculate fade progress
		fadeProgress := float64(f) / float64(frameCount-1)

		op := &ebiten.DrawImageOptions{}

		// Desaturate and darken as death progresses
		grayScale := 1.0 - fadeProgress*0.7
		alphaScale := 1.0 - fadeProgress*0.8

		op.ColorScale.Scale(float32(grayScale), float32(grayScale), float32(grayScale), float32(alphaScale))

		frame.DrawImage(base, op)
		frames[f] = frame
	}

	return frames
}
