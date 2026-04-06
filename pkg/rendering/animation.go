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
// Uses WritePixels for efficient bulk pixel upload instead of per-pixel Set().
func (atg *AnimatedTileGenerator) generateWaterTile(baseColor, highlightColor color.Color) *AnimatedTile {
	const frameCount = 4
	frames := make([]*ebiten.Image, frameCount)
	size := atg.tileSize

	for f := 0; f < frameCount; f++ {
		img := ebiten.NewImage(size, size)
		pixels := make([]byte, size*size*4)
		offset := float64(f) * (float64(size) / float64(frameCount))
		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				wave := atg.waterWaveValue(x, y, offset)
				c := atg.blendColors(baseColor, highlightColor, wave)
				r, g, b, a := c.RGBA()
				i := (y*size + x) * 4
				pixels[i] = uint8(r >> 8)
				pixels[i+1] = uint8(g >> 8)
				pixels[i+2] = uint8(b >> 8)
				pixels[i+3] = uint8(a >> 8)
			}
		}
		img.WritePixels(pixels)
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
// Uses bulk pixel operations for efficient initialization.
func (atg *AnimatedTileGenerator) generateGrassTile(baseColor, tipColor color.Color) *AnimatedTile {
	const frameCount = 4
	frames := make([]*ebiten.Image, frameCount)
	size := atg.tileSize

	grassBlades := atg.generateGrassPattern(size)

	// Pre-extract color components
	br, bg, bb, ba := baseColor.RGBA()
	tr, tg, tb, ta := tipColor.RGBA()
	baseR, baseG, baseB, baseA := uint8(br>>8), uint8(bg>>8), uint8(bb>>8), uint8(ba>>8)
	tipR, tipG, tipB, tipA := uint8(tr>>8), uint8(tg>>8), uint8(tb>>8), uint8(ta>>8)

	for f := 0; f < frameCount; f++ {
		img := ebiten.NewImage(size, size)
		pixels := make([]byte, size*size*4)

		// Fill with base color
		for i := 0; i < size*size; i++ {
			pixels[i*4] = baseR
			pixels[i*4+1] = baseG
			pixels[i*4+2] = baseB
			pixels[i*4+3] = baseA
		}

		// Draw grass blades
		windOffset := float64(f) / float64(frameCount)
		atg.drawGrassBladesToPixels(pixels, grassBlades, tipR, tipG, tipB, tipA, windOffset)

		img.WritePixels(pixels)
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
// Uses per-pixel Set for sparse drawing (few blade pixels).
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

// drawGrassBladesToPixels draws grass blades directly to a pixel buffer.
func (atg *AnimatedTileGenerator) drawGrassBladesToPixels(pixels []byte, blades []grassBlade, r, g, b, a uint8, windOffset float64) {
	size := atg.tileSize
	for _, blade := range blades {
		sway := int(sinApprox((float64(blade.x)/float64(size)+windOffset)*6.28) * 2)
		for h := 0; h < blade.height; h++ {
			py := blade.y - h
			px := blade.x + (sway * h / blade.height)
			if py >= 0 && py < size && px >= 0 && px < size {
				i := (py*size + px) * 4
				pixels[i] = r
				pixels[i+1] = g
				pixels[i+2] = b
				pixels[i+3] = a
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
// Uses WritePixels for efficient bulk pixel upload instead of per-pixel Set().
func (atg *AnimatedTileGenerator) drawFireFrame(img *ebiten.Image, baseColor, brightColor color.Color, frame int) {
	size := atg.tileSize
	pixels := make([]byte, size*size*4)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			intensity := atg.fireIntensity(x, y, frame)
			c := atg.blendColors(baseColor, brightColor, intensity)
			r, g, b, a := c.RGBA()
			i := (y*size + x) * 4
			pixels[i] = uint8(r >> 8)
			pixels[i+1] = uint8(g >> 8)
			pixels[i+2] = uint8(b >> 8)
			pixels[i+3] = uint8(a >> 8)
		}
	}
	img.WritePixels(pixels)
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

// Note: sinApprox and lerp are defined in animation_core.go which is available in all builds.
// Note: clampFloat is defined in lighting_core.go which is available in all builds.
