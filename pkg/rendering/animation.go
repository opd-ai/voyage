//go:build !headless

package rendering

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// AnimatedTile represents a tile with multiple animation frames.
type AnimatedTile struct {
	Frames       []*ebiten.Image
	FrameTime    float64 // Duration per frame in seconds
	Loop         bool
	currentFrame int
	elapsed      float64
}

// NewAnimatedTile creates a new animated tile with the given frames.
func NewAnimatedTile(frames []*ebiten.Image, frameTime float64, loop bool) *AnimatedTile {
	return &AnimatedTile{
		Frames:       frames,
		FrameTime:    frameTime,
		Loop:         loop,
		currentFrame: 0,
		elapsed:      0,
	}
}

// Update advances the animation by the given delta time.
func (at *AnimatedTile) Update(dt float64) {
	if len(at.Frames) <= 1 {
		return
	}
	at.elapsed += dt
	if at.elapsed >= at.FrameTime {
		at.elapsed -= at.FrameTime
		at.currentFrame++
		if at.currentFrame >= len(at.Frames) {
			if at.Loop {
				at.currentFrame = 0
			} else {
				at.currentFrame = len(at.Frames) - 1
			}
		}
	}
}

// CurrentFrame returns the current animation frame image.
func (at *AnimatedTile) CurrentFrame() *ebiten.Image {
	if len(at.Frames) == 0 {
		return nil
	}
	return at.Frames[at.currentFrame]
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

// AnimatedTileGenerator creates animated overworld tiles.
type AnimatedTileGenerator struct {
	gen      *seed.Generator
	tileSize int
}

// NewAnimatedTileGenerator creates a new animated tile generator.
func NewAnimatedTileGenerator(masterSeed int64, tileSize int) *AnimatedTileGenerator {
	return &AnimatedTileGenerator{
		gen:      seed.NewGenerator(masterSeed, "animation"),
		tileSize: tileSize,
	}
}

// GenerateAnimatedTile creates an animated tile of the specified type.
func (atg *AnimatedTileGenerator) GenerateAnimatedTile(animType AnimationType, baseColor, accentColor color.Color) *AnimatedTile {
	switch animType {
	case AnimationWater:
		return atg.generateWaterTile(baseColor, accentColor)
	case AnimationGrass:
		return atg.generateGrassTile(baseColor, accentColor)
	case AnimationFire:
		return atg.generateFireTile(baseColor, accentColor)
	default:
		return atg.generateWaterTile(baseColor, accentColor)
	}
}

// generateWaterTile creates a flowing water animation with 4 frames.
func (atg *AnimatedTileGenerator) generateWaterTile(baseColor, highlightColor color.Color) *AnimatedTile {
	const frameCount = 4
	frames := make([]*ebiten.Image, frameCount)
	size := atg.tileSize

	for f := 0; f < frameCount; f++ {
		img := ebiten.NewImage(size, size)
		offset := float64(f) * (float64(size) / float64(frameCount))
		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				wave := atg.waterWaveValue(x, y, offset)
				c := atg.blendColors(baseColor, highlightColor, wave)
				img.Set(x, y, c)
			}
		}
		frames[f] = img
	}
	return NewAnimatedTile(frames, 0.2, true)
}

// waterWaveValue computes a wave pattern for water animation.
func (atg *AnimatedTileGenerator) waterWaveValue(x, y int, offset float64) float64 {
	size := float64(atg.tileSize)
	fx := float64(x) / size
	fy := float64(y) / size
	wave1 := sinApprox((fx*4 + fy*2 + offset/size) * 6.28)
	wave2 := sinApprox((fx*2 - fy*3 + offset/size*1.5) * 6.28)
	return (wave1*0.5 + wave2*0.5 + 1.0) * 0.5
}

// generateGrassTile creates a wind-swept grass animation with 4 frames.
func (atg *AnimatedTileGenerator) generateGrassTile(baseColor, tipColor color.Color) *AnimatedTile {
	const frameCount = 4
	frames := make([]*ebiten.Image, frameCount)
	size := atg.tileSize

	grassBlades := atg.generateGrassPattern(size)

	for f := 0; f < frameCount; f++ {
		img := ebiten.NewImage(size, size)
		img.Fill(baseColor)
		windOffset := float64(f) / float64(frameCount)
		atg.drawGrassBlades(img, grassBlades, tipColor, windOffset)
		frames[f] = img
	}
	return NewAnimatedTile(frames, 0.15, true)
}

// grassBlade represents a single blade of grass.
type grassBlade struct {
	x, y   int
	height int
}

// generateGrassPattern creates random grass blade positions.
func (atg *AnimatedTileGenerator) generateGrassPattern(size int) []grassBlade {
	bladeCount := size * size / 8
	blades := make([]grassBlade, bladeCount)
	for i := range blades {
		blades[i] = grassBlade{
			x:      atg.gen.Intn(size),
			y:      atg.gen.Intn(size),
			height: 2 + atg.gen.Intn(4),
		}
	}
	return blades
}

// drawGrassBlades draws grass blades with wind sway effect.
func (atg *AnimatedTileGenerator) drawGrassBlades(img *ebiten.Image, blades []grassBlade, tipColor color.Color, windOffset float64) {
	for _, blade := range blades {
		sway := int(sinApprox((float64(blade.x)/float64(atg.tileSize)+windOffset)*6.28) * 2)
		for h := 0; h < blade.height; h++ {
			py := blade.y - h
			px := blade.x + (sway * h / blade.height)
			if py >= 0 && py < atg.tileSize && px >= 0 && px < atg.tileSize {
				img.Set(px, py, tipColor)
			}
		}
	}
}

// generateFireTile creates a flickering fire animation with 4 frames.
func (atg *AnimatedTileGenerator) generateFireTile(baseColor, brightColor color.Color) *AnimatedTile {
	const frameCount = 4
	frames := make([]*ebiten.Image, frameCount)
	size := atg.tileSize

	for f := 0; f < frameCount; f++ {
		img := ebiten.NewImage(size, size)
		atg.drawFireFrame(img, baseColor, brightColor, f)
		frames[f] = img
	}
	return NewAnimatedTile(frames, 0.1, true)
}

// drawFireFrame draws a single frame of fire animation.
func (atg *AnimatedTileGenerator) drawFireFrame(img *ebiten.Image, baseColor, brightColor color.Color, frame int) {
	size := atg.tileSize
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			intensity := atg.fireIntensity(x, y, frame)
			c := atg.blendColors(baseColor, brightColor, intensity)
			img.Set(x, y, c)
		}
	}
}

// fireIntensity calculates fire brightness with random flickering.
func (atg *AnimatedTileGenerator) fireIntensity(x, y, frame int) float64 {
	size := float64(atg.tileSize)
	fx := float64(x) / size
	fy := float64(y) / size
	centerDist := (0.5-fx)*(0.5-fx) + (0.5-fy)*(0.5-fy)
	baseIntensity := 1.0 - centerDist*4
	flicker := (float64(frame%4) + 1.0) / 4.0
	noise := atg.gen.Float64()*0.3 - 0.15
	return clampFloat(baseIntensity*flicker+noise, 0, 1)
}

// blendColors blends two colors by a factor t (0=a, 1=b).
func (atg *AnimatedTileGenerator) blendColors(a, b color.Color, t float64) color.Color {
	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	return color.RGBA{
		R: uint8(lerp(float64(ar>>8), float64(br>>8), t)),
		G: uint8(lerp(float64(ag>>8), float64(bg>>8), t)),
		B: uint8(lerp(float64(ab>>8), float64(bb>>8), t)),
		A: uint8(lerp(float64(aa>>8), float64(ba>>8), t)),
	}
}

// sinApprox provides a fast sine approximation.
func sinApprox(x float64) float64 {
	const twoPi = 6.283185307179586
	x = x - float64(int(x/twoPi))*twoPi
	if x < 0 {
		x += twoPi
	}
	if x > 3.14159 {
		x -= twoPi
	}
	return x - (x*x*x)/6.0 + (x*x*x*x*x)/120.0
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
