//go:build !headless

package rendering

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/opd-ai/voyage/pkg/engine"
)

// PostProcessor applies post-processing effects to rendered images.
// Uses cached overlay images and DrawImage for GPU-accelerated compositing (H-004, H-005, H-006, M-001).
type PostProcessor struct {
	genre  engine.GenreID
	config PostProcessorConfig

	// Cached overlays for efficient rendering (H-004, M-001)
	vignetteCache     *ebiten.Image
	vignetteCacheSize [2]int // [width, height]
	vignetteIntensity float64

	grainTexture     *ebiten.Image
	grainTextureSize int // Size of the grain texture (typically 256x256)
	lastGrainSeed    int64
}

// NewPostProcessor creates a new post processor.
func NewPostProcessor(genre engine.GenreID) *PostProcessor {
	pp := &PostProcessor{
		config:           DefaultPostProcessorConfig(),
		grainTextureSize: 256, // Small repeating texture
	}
	pp.SetGenre(genre)
	return pp
}

// SetGenre configures the post processor for a specific genre.
func (pp *PostProcessor) SetGenre(genre engine.GenreID) {
	pp.genre = genre
	pp.config = ConfigureForGenre(genre)
}

// Genre returns the current genre.
func (pp *PostProcessor) Genre() engine.GenreID {
	return pp.genre
}

// copyImage creates a new image copy with the same dimensions and content.
// Returns nil if the input is nil.
func copyImage(img *ebiten.Image) (*ebiten.Image, int, int) {
	if img == nil {
		return nil, 0, 0
	}
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	result := ebiten.NewImage(w, h)
	result.DrawImage(img, nil)
	return result, w, h
}

// Apply processes an image with all enabled post-processing effects.
func (pp *PostProcessor) Apply(img *ebiten.Image, seed int64) *ebiten.Image {
	if img == nil {
		return nil
	}
	return pp.applyEffectsInOrder(img, seed)
}

// applyEffectsInOrder applies each enabled effect in sequence.
func (pp *PostProcessor) applyEffectsInOrder(img *ebiten.Image, seed int64) *ebiten.Image {
	result := img
	result = pp.maybeApplyVignette(result)
	result = pp.maybeApplyScanlines(result)
	result = pp.maybeApplyFilmGrain(result, seed)
	result = pp.maybeApplyChromatic(result)
	result = pp.maybeApplySepia(result)
	return result
}

func (pp *PostProcessor) maybeApplyVignette(img *ebiten.Image) *ebiten.Image {
	if pp.config.VignetteOn {
		return pp.ApplyVignette(img, pp.config.VignetteInt)
	}
	return img
}

func (pp *PostProcessor) maybeApplyScanlines(img *ebiten.Image) *ebiten.Image {
	if pp.config.ScanlinesOn {
		return pp.ApplyScanlines(img, pp.config.ScanlinesDen, 0.15)
	}
	return img
}

func (pp *PostProcessor) maybeApplyFilmGrain(img *ebiten.Image, seed int64) *ebiten.Image {
	if pp.config.FilmGrainOn {
		return pp.ApplyFilmGrain(img, seed, pp.config.FilmGrainInt)
	}
	return img
}

func (pp *PostProcessor) maybeApplyChromatic(img *ebiten.Image) *ebiten.Image {
	if pp.config.ChromaticOn {
		return pp.ApplyChromaticAberration(img, pp.config.ChromaticOff)
	}
	return img
}

func (pp *PostProcessor) maybeApplySepia(img *ebiten.Image) *ebiten.Image {
	if pp.config.SepiaOn {
		return pp.ApplySepia(img, pp.config.SepiaInt)
	}
	return img
}

// ApplyVignette darkens the edges of an image using a pre-rendered overlay (H-004).
// This replaces per-pixel Set() operations with a single DrawImage composite.
func (pp *PostProcessor) ApplyVignette(img *ebiten.Image, intensity float64) *ebiten.Image {
	if img == nil {
		return nil
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Check if we need to regenerate the vignette cache
	if pp.vignetteCache == nil ||
		pp.vignetteCacheSize[0] != w ||
		pp.vignetteCacheSize[1] != h ||
		pp.vignetteIntensity != intensity {
		pp.generateVignetteCache(w, h, intensity)
	}

	// Create result by drawing original then compositing vignette overlay
	result := ebiten.NewImage(w, h)
	result.DrawImage(img, nil)

	// Composite the vignette overlay using multiply blend via ColorScale
	op := &ebiten.DrawImageOptions{}
	op.Blend = ebiten.BlendSourceOver
	result.DrawImage(pp.vignetteCache, op)

	return result
}

// generateVignetteCache creates a pre-rendered vignette overlay image (H-004).
// The overlay uses alpha channel to darken edges when composited.
func (pp *PostProcessor) generateVignetteCache(w, h int, intensity float64) {
	pp.vignetteCache = ebiten.NewImage(w, h)
	pp.vignetteCacheSize = [2]int{w, h}
	pp.vignetteIntensity = intensity

	centerX := float64(w) / 2
	centerY := float64(h) / 2
	maxDist := math.Max(centerX, centerY)
	maxDistSq := maxDist * maxDist

	// Generate the vignette as a darkening overlay
	// We use RGBA where RGB is black and A controls the darkening amount
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := float64(x) - centerX
			dy := float64(y) - centerY
			dist := (dx*dx + dy*dy) / maxDistSq
			// Calculate darkening: 0 at center, increases toward edges
			darkness := dist * intensity
			if darkness > 1 {
				darkness = 1
			}
			// Alpha controls how much black is applied
			alpha := uint8(darkness * 255)
			pp.vignetteCache.Set(x, y, color.RGBA{0, 0, 0, alpha})
		}
	}
}

// ApplyScanlines adds horizontal scanline effect for retro/sci-fi feel.
func (pp *PostProcessor) ApplyScanlines(img *ebiten.Image, density, alpha float64) *ebiten.Image {
	result, w, h := copyImage(img)
	if result == nil {
		return nil
	}

	// Add scanlines every N pixels
	spacing := int(density)
	if spacing < 1 {
		spacing = 1
	}

	for y := 0; y < h; y += spacing {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			factor := 1.0 - alpha
			newR := uint8(float64(r>>8) * factor)
			newG := uint8(float64(g>>8) * factor)
			newB := uint8(float64(b>>8) * factor)
			result.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}

	return result
}

// ApplyFilmGrain adds random noise using a pre-rendered noise texture (M-001).
// This replaces per-pixel noise generation with a tiled texture overlay.
func (pp *PostProcessor) ApplyFilmGrain(img *ebiten.Image, seed int64, intensity float64) *ebiten.Image {
	if img == nil {
		return nil
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Regenerate grain texture if seed changed
	if pp.grainTexture == nil || pp.lastGrainSeed != seed {
		pp.generateGrainTexture(seed)
	}

	// Copy the original image
	result := ebiten.NewImage(w, h)
	result.DrawImage(img, nil)

	// Tile the grain texture across the result
	// Use overlay blend mode for proper grain effect
	grainSize := pp.grainTextureSize
	for ty := 0; ty < h; ty += grainSize {
		for tx := 0; tx < w; tx += grainSize {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(tx), float64(ty))
			// Scale alpha by intensity
			op.ColorScale.Scale(1, 1, 1, float32(intensity))
			op.Blend = ebiten.BlendSourceOver
			result.DrawImage(pp.grainTexture, op)
		}
	}

	return result
}

// generateGrainTexture creates a pre-rendered noise texture for film grain (M-001).
func (pp *PostProcessor) generateGrainTexture(seed int64) {
	size := pp.grainTextureSize
	pp.grainTexture = ebiten.NewImage(size, size)
	pp.lastGrainSeed = seed

	// Simple LCG for deterministic noise
	rng := uint64(seed)
	nextRand := func() float64 {
		rng = rng*6364136223846793005 + 1442695040888963407
		return float64(rng>>33) / float64(1<<31)
	}

	// Generate grayscale noise
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// Random value centered at 128 (neutral gray)
			noise := (nextRand() - 0.5) * 128
			gray := uint8(clampFloat64(128+noise, 0, 255))
			// Use low alpha so it blends subtly with the underlying image
			pp.grainTexture.Set(x, y, color.RGBA{gray, gray, gray, 64})
		}
	}
}

// clampFloat64 clamps a float64 value to [min, max].
func clampFloat64(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// ApplyChromaticAberration offsets color channels using DrawImage with ColorM (H-005).
// This replaces per-pixel triple sampling with three GPU-accelerated draw calls.
func (pp *PostProcessor) ApplyChromaticAberration(img *ebiten.Image, offset float64) *ebiten.Image {
	if img == nil {
		return nil
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	result := ebiten.NewImage(w, h)
	off := offset

	// Draw red channel shifted left
	opR := &ebiten.DrawImageOptions{}
	opR.GeoM.Translate(-off, 0)
	var cmR ebiten.ColorM
	cmR.Scale(1, 0, 0, 1) // Keep only red channel
	opR.ColorM = cmR
	opR.Blend = ebiten.BlendSourceOver
	result.DrawImage(img, opR)

	// Draw green channel centered (additive blend)
	opG := &ebiten.DrawImageOptions{}
	var cmG ebiten.ColorM
	cmG.Scale(0, 1, 0, 1) // Keep only green channel
	opG.ColorM = cmG
	opG.Blend = ebiten.BlendLighter // Additive blending
	result.DrawImage(img, opG)

	// Draw blue channel shifted right (additive blend)
	opB := &ebiten.DrawImageOptions{}
	opB.GeoM.Translate(off, 0)
	var cmB ebiten.ColorM
	cmB.Scale(0, 0, 1, 1) // Keep only blue channel
	opB.ColorM = cmB
	opB.Blend = ebiten.BlendLighter // Additive blending
	result.DrawImage(img, opB)

	return result
}

// clampInt clamps v to the range [min, max].
func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// ApplySepia applies a warm sepia tone for vintage/dusty look.
func (pp *PostProcessor) ApplySepia(img *ebiten.Image, intensity float64) *ebiten.Image {
	if img == nil {
		return nil
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	result := ebiten.NewImage(w, h)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			rf := float64(r >> 8)
			gf := float64(g >> 8)
			bf := float64(b >> 8)

			// Standard sepia transformation
			sepiaR := rf*0.393 + gf*0.769 + bf*0.189
			sepiaG := rf*0.349 + gf*0.686 + bf*0.168
			sepiaB := rf*0.272 + gf*0.534 + bf*0.131

			// Blend original with sepia
			newR := clampUint8(rf*(1-intensity) + sepiaR*intensity)
			newG := clampUint8(gf*(1-intensity) + sepiaG*intensity)
			newB := clampUint8(bf*(1-intensity) + sepiaB*intensity)

			result.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}

	return result
}

// SetVignetteIntensity sets the vignette darkening strength.
func (pp *PostProcessor) SetVignetteIntensity(intensity float64) {
	pp.config.VignetteInt = intensity
}

// SetScanlinesEnabled toggles scanline effect.
func (pp *PostProcessor) SetScanlinesEnabled(enabled bool) {
	pp.config.ScanlinesOn = enabled
}

// SetFilmGrainEnabled toggles film grain effect.
func (pp *PostProcessor) SetFilmGrainEnabled(enabled bool) {
	pp.config.FilmGrainOn = enabled
}

// SetChromaticEnabled toggles chromatic aberration effect.
func (pp *PostProcessor) SetChromaticEnabled(enabled bool) {
	pp.config.ChromaticOn = enabled
}

// SetSepiaEnabled toggles sepia effect.
func (pp *PostProcessor) SetSepiaEnabled(enabled bool) {
	pp.config.SepiaOn = enabled
}
