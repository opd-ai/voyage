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

	// Cached scanlines overlay for GPU-accelerated rendering (M-003)
	scanlinesCache        *ebiten.Image
	scanlinesCacheSize    [2]int
	scanlinesCacheDensity float64
	scanlinesCacheAlpha   float64

	// Double buffer for ping-pong rendering to avoid per-frame allocations (M-004)
	bufferA     *ebiten.Image
	bufferB     *ebiten.Image
	bufferSize  [2]int
	useBufferA  bool
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
	// Dispose cached images when changing genre since visual effects change
	pp.disposeCache()
}

// disposeCache releases cached images to free GPU memory.
func (pp *PostProcessor) disposeCache() {
	if pp.vignetteCache != nil {
		pp.vignetteCache.Dispose()
		pp.vignetteCache = nil
	}
	if pp.grainTexture != nil {
		pp.grainTexture.Dispose()
		pp.grainTexture = nil
	}
	if pp.scanlinesCache != nil {
		pp.scanlinesCache.Dispose()
		pp.scanlinesCache = nil
	}
	if pp.bufferA != nil {
		pp.bufferA.Dispose()
		pp.bufferA = nil
	}
	if pp.bufferB != nil {
		pp.bufferB.Dispose()
		pp.bufferB = nil
	}
}

// Dispose releases all GPU resources held by the post processor.
// Call this when the post processor is no longer needed.
func (pp *PostProcessor) Dispose() {
	pp.disposeCache()
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

// getResultBuffer returns a cached result buffer of the specified size.
// Reuses existing buffer if dimensions match, otherwise creates a new one.
// This reduces per-frame allocations for post-processing effects (M-004).
func (pp *PostProcessor) ensureBuffers(w, h int) {
	if pp.bufferA == nil || pp.bufferSize[0] != w || pp.bufferSize[1] != h {
		if pp.bufferA != nil {
			pp.bufferA.Dispose()
		}
		if pp.bufferB != nil {
			pp.bufferB.Dispose()
		}
		pp.bufferA = ebiten.NewImage(w, h)
		pp.bufferB = ebiten.NewImage(w, h)
		pp.bufferSize = [2]int{w, h}
	}
}

// getNextBuffer returns the next buffer in the ping-pong sequence and clears it.
func (pp *PostProcessor) getNextBuffer() *ebiten.Image {
	pp.useBufferA = !pp.useBufferA
	var buf *ebiten.Image
	if pp.useBufferA {
		buf = pp.bufferA
	} else {
		buf = pp.bufferB
	}
	buf.Clear()
	return buf
}

// Apply processes an image with all enabled post-processing effects.
// Uses cached double buffers to avoid per-frame GPU memory allocations (M-004).
func (pp *PostProcessor) Apply(img *ebiten.Image, seed int64) *ebiten.Image {
	if img == nil {
		return nil
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Ensure buffers are sized correctly
	pp.ensureBuffers(w, h)

	return pp.applyEffectsInOrder(img, seed)
}

// applyEffectsInOrder applies each enabled effect in sequence using ping-pong buffers.
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
// Uses cached ping-pong buffers to avoid per-frame allocations (M-004).
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

	// Use cached buffer instead of creating new image
	result := pp.getNextBuffer()
	result.DrawImage(img, nil)

	// Composite the vignette overlay using multiply blend via ColorScale
	op := &ebiten.DrawImageOptions{}
	op.Blend = ebiten.BlendSourceOver
	result.DrawImage(pp.vignetteCache, op)

	return result
}

// generateVignetteCache creates a pre-rendered vignette overlay image (H-004).
// The overlay uses alpha channel to darken edges when composited.
// Uses WritePixels for efficient bulk pixel upload instead of per-pixel Set().
func (pp *PostProcessor) generateVignetteCache(w, h int, intensity float64) {
	// Guard against zero dimensions to prevent division by zero
	if w <= 0 || h <= 0 {
		return
	}

	pp.vignetteCache = ebiten.NewImage(w, h)
	pp.vignetteCacheSize = [2]int{w, h}
	pp.vignetteIntensity = intensity

	centerX := float64(w) / 2
	centerY := float64(h) / 2
	maxDist := math.Max(centerX, centerY)
	maxDistSq := maxDist * maxDist

	// Guard against zero maxDistSq (should not happen with positive w,h but be safe)
	if maxDistSq <= 0 {
		return
	}

	// Generate the vignette as a darkening overlay using bulk pixel data
	// We use RGBA where RGB is black and A controls the darkening amount
	pixels := make([]byte, w*h*4)
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
			i := (y*w + x) * 4
			pixels[i] = 0     // R
			pixels[i+1] = 0   // G
			pixels[i+2] = 0   // B
			pixels[i+3] = alpha
		}
	}
	pp.vignetteCache.WritePixels(pixels)
}

// ApplyScanlines adds horizontal scanline effect for retro/sci-fi feel.
// Uses a cached scanline overlay for GPU-accelerated compositing (M-003).
// Uses cached ping-pong buffers to avoid per-frame allocations (M-004).
func (pp *PostProcessor) ApplyScanlines(img *ebiten.Image, density, alpha float64) *ebiten.Image {
	if img == nil {
		return nil
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Check if we need to regenerate the scanlines cache
	if pp.scanlinesCache == nil ||
		pp.scanlinesCacheSize[0] != w ||
		pp.scanlinesCacheSize[1] != h ||
		pp.scanlinesCacheDensity != density ||
		pp.scanlinesCacheAlpha != alpha {
		pp.generateScanlinesCache(w, h, density, alpha)
	}

	// Use cached buffer instead of creating new image
	result := pp.getNextBuffer()
	result.DrawImage(img, nil)

	// Composite the scanlines overlay using multiply blend
	op := &ebiten.DrawImageOptions{}
	op.Blend = ebiten.BlendSourceOver
	result.DrawImage(pp.scanlinesCache, op)

	return result
}

// generateScanlinesCache creates a pre-rendered scanlines overlay image (M-003).
// The overlay uses alpha channel to darken scanline rows when composited.
// Uses WritePixels for efficient bulk pixel upload instead of per-pixel Set().
func (pp *PostProcessor) generateScanlinesCache(w, h int, density, alpha float64) {
	if w <= 0 || h <= 0 {
		return
	}

	pp.scanlinesCache = ebiten.NewImage(w, h)
	pp.scanlinesCacheSize = [2]int{w, h}
	pp.scanlinesCacheDensity = density
	pp.scanlinesCacheAlpha = alpha

	spacing := int(density)
	if spacing < 1 {
		spacing = 1
	}

	// Create darkening overlay for scanline rows using bulk pixel data
	// Alpha controls how much the underlying image is darkened
	darkenAlpha := uint8(alpha * 255)
	pixels := make([]byte, w*h*4)

	for y := 0; y < h; y += spacing {
		for x := 0; x < w; x++ {
			i := (y*w + x) * 4
			pixels[i] = 0            // R
			pixels[i+1] = 0          // G
			pixels[i+2] = 0          // B
			pixels[i+3] = darkenAlpha // A
		}
	}
	pp.scanlinesCache.WritePixels(pixels)
}

// ApplyFilmGrain adds random noise using a pre-rendered noise texture (M-001).
// This replaces per-pixel noise generation with a tiled texture overlay.
// Uses cached ping-pong buffers to avoid per-frame allocations (M-004).
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

	// Use cached buffer instead of creating new image
	result := pp.getNextBuffer()
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
// Uses WritePixels for efficient bulk pixel upload instead of per-pixel Set().
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

	// Generate grayscale noise using bulk pixel data
	pixels := make([]byte, size*size*4)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// Random value centered at 128 (neutral gray)
			noise := (nextRand() - 0.5) * 128
			gray := uint8(clampFloat64(128+noise, 0, 255))
			// Use low alpha so it blends subtly with the underlying image
			i := (y*size + x) * 4
			pixels[i] = gray   // R
			pixels[i+1] = gray // G
			pixels[i+2] = gray // B
			pixels[i+3] = 64   // A
		}
	}
	pp.grainTexture.WritePixels(pixels)
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
// Uses cached ping-pong buffers to avoid per-frame allocations (M-004).
func (pp *PostProcessor) ApplyChromaticAberration(img *ebiten.Image, offset float64) *ebiten.Image {
	if img == nil {
		return nil
	}

	// Use cached buffer instead of creating new image
	result := pp.getNextBuffer()
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
// Uses ColorM matrix for GPU-accelerated color transformation (M-003).
// Uses cached ping-pong buffers to avoid per-frame allocations (M-004).
func (pp *PostProcessor) ApplySepia(img *ebiten.Image, intensity float64) *ebiten.Image {
	if img == nil {
		return nil
	}

	// Use cached buffer instead of creating new image
	result := pp.getNextBuffer()

	// Standard sepia transformation matrix coefficients
	// R' = R*0.393 + G*0.769 + B*0.189
	// G' = R*0.349 + G*0.686 + B*0.168
	// B' = R*0.272 + G*0.534 + B*0.131
	// Blended with identity matrix based on intensity

	// Identity matrix coefficients
	i := 1.0 - intensity

	// Build the blended color transformation matrix
	var cm ebiten.ColorM
	cm.Reset()

	// Row for Red output: lerp(identity, sepia)
	cm.SetElement(0, 0, i+0.393*intensity) // R contribution to R'
	cm.SetElement(0, 1, 0.769*intensity)   // G contribution to R'
	cm.SetElement(0, 2, 0.189*intensity)   // B contribution to R'

	// Row for Green output
	cm.SetElement(1, 0, 0.349*intensity)   // R contribution to G'
	cm.SetElement(1, 1, i+0.686*intensity) // G contribution to G'
	cm.SetElement(1, 2, 0.168*intensity)   // B contribution to G'

	// Row for Blue output
	cm.SetElement(2, 0, 0.272*intensity)   // R contribution to B'
	cm.SetElement(2, 1, 0.534*intensity)   // G contribution to B'
	cm.SetElement(2, 2, i+0.131*intensity) // B contribution to B'

	op := &ebiten.DrawImageOptions{}
	op.ColorM = cm
	result.DrawImage(img, op)

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
