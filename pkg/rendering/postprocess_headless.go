//go:build headless

package rendering

import (
	"image"
	"image/color"

	"github.com/opd-ai/voyage/pkg/engine"
)

// PostProcessor applies post-processing effects to rendered images.
type PostProcessor struct {
	genre  engine.GenreID
	config PostProcessorConfig
}

// NewPostProcessor creates a new post processor.
func NewPostProcessor(genre engine.GenreID) *PostProcessor {
	pp := &PostProcessor{
		config: DefaultPostProcessorConfig(),
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

// Apply processes an image with all enabled post-processing effects.
func (pp *PostProcessor) Apply(img *image.RGBA, seed int64) *image.RGBA {
	if img == nil {
		return nil
	}
	return pp.applyEffectsInOrder(img, seed)
}

// applyEffectsInOrder applies each enabled effect in sequence.
func (pp *PostProcessor) applyEffectsInOrder(img *image.RGBA, seed int64) *image.RGBA {
	result := img
	result = pp.maybeApplyVignette(result)
	result = pp.maybeApplyScanlines(result)
	result = pp.maybeApplyFilmGrain(result, seed)
	result = pp.maybeApplyChromatic(result)
	result = pp.maybeApplySepia(result)
	return result
}

func (pp *PostProcessor) maybeApplyVignette(img *image.RGBA) *image.RGBA {
	if pp.config.VignetteOn {
		return pp.ApplyVignette(img, pp.config.VignetteInt)
	}
	return img
}

func (pp *PostProcessor) maybeApplyScanlines(img *image.RGBA) *image.RGBA {
	if pp.config.ScanlinesOn {
		return pp.ApplyScanlines(img, pp.config.ScanlinesDen, 0.15)
	}
	return img
}

func (pp *PostProcessor) maybeApplyFilmGrain(img *image.RGBA, seed int64) *image.RGBA {
	if pp.config.FilmGrainOn {
		return pp.ApplyFilmGrain(img, seed, pp.config.FilmGrainInt)
	}
	return img
}

func (pp *PostProcessor) maybeApplyChromatic(img *image.RGBA) *image.RGBA {
	if pp.config.ChromaticOn {
		return pp.ApplyChromaticAberration(img, pp.config.ChromaticOff)
	}
	return img
}

func (pp *PostProcessor) maybeApplySepia(img *image.RGBA) *image.RGBA {
	if pp.config.SepiaOn {
		return pp.ApplySepia(img, pp.config.SepiaInt)
	}
	return img
}

// ApplyVignette darkens the edges of an image to focus attention on center.
func (pp *PostProcessor) ApplyVignette(img *image.RGBA, intensity float64) *image.RGBA {
	if img == nil {
		return nil
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	result := image.NewRGBA(bounds)

	vc := newVignetteCalcHeadless(w, h, intensity)
	vc.applyToImage(result, img)
	return result
}

// vignetteCalcHeadless handles vignette darkening calculations for headless mode.
type vignetteCalcHeadless struct {
	centerX, centerY float64
	maxDistSq        float64
	intensity        float64
	w, h             int
}

func newVignetteCalcHeadless(w, h int, intensity float64) *vignetteCalcHeadless {
	centerX := float64(w) / 2
	centerY := float64(h) / 2
	maxDist := centerX
	if centerY > maxDist {
		maxDist = centerY
	}
	return &vignetteCalcHeadless{
		centerX:   centerX,
		centerY:   centerY,
		maxDistSq: maxDist * maxDist,
		intensity: intensity,
		w:         w,
		h:         h,
	}
}

func (vc *vignetteCalcHeadless) applyToImage(result, img *image.RGBA) {
	for y := 0; y < vc.h; y++ {
		for x := 0; x < vc.w; x++ {
			factor := vc.darkeningFactor(x, y)
			c := img.RGBAAt(x, y)
			result.SetRGBA(x, y, color.RGBA{
				R: uint8(float64(c.R) * factor),
				G: uint8(float64(c.G) * factor),
				B: uint8(float64(c.B) * factor),
				A: c.A,
			})
		}
	}
}

func (vc *vignetteCalcHeadless) darkeningFactor(x, y int) float64 {
	dx := float64(x) - vc.centerX
	dy := float64(y) - vc.centerY
	dist := (dx*dx + dy*dy) / vc.maxDistSq
	factor := 1.0 - (dist * vc.intensity)
	if factor < 0 {
		return 0
	}
	return factor
}

// ApplyScanlines adds horizontal scanline effect for retro/sci-fi feel.
func (pp *PostProcessor) ApplyScanlines(img *image.RGBA, density, alpha float64) *image.RGBA {
	if img == nil {
		return nil
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	result := image.NewRGBA(bounds)

	copyImageRGBA(result, img, w, h)
	applyScanlineEffect(result, img, w, h, density, alpha)

	return result
}

// copyImageRGBA copies all pixels from src to dst.
func copyImageRGBA(dst, src *image.RGBA, w, h int) {
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dst.Set(x, y, src.At(x, y))
		}
	}
}

// applyScanlineEffect darkens pixels at scanline intervals.
func applyScanlineEffect(result, src *image.RGBA, w, h int, density, alpha float64) {
	spacing := int(density)
	if spacing < 1 {
		spacing = 1
	}
	factor := 1.0 - alpha

	for y := 0; y < h; y += spacing {
		for x := 0; x < w; x++ {
			c := src.RGBAAt(x, y)
			result.SetRGBA(x, y, color.RGBA{
				R: uint8(float64(c.R) * factor),
				G: uint8(float64(c.G) * factor),
				B: uint8(float64(c.B) * factor),
				A: c.A,
			})
		}
	}
}

// ApplyFilmGrain adds random noise for a gritty film effect.
func (pp *PostProcessor) ApplyFilmGrain(img *image.RGBA, seed int64, intensity float64) *image.RGBA {
	if img == nil {
		return nil
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	result := image.NewRGBA(bounds)

	// Simple LCG for deterministic noise
	rng := uint64(seed)
	nextRand := func() float64 {
		rng = rng*6364136223846793005 + 1442695040888963407
		return float64(rng>>33) / float64(1<<31)
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := img.RGBAAt(x, y)
			noise := (nextRand() - 0.5) * 2.0 * intensity * 255
			result.SetRGBA(x, y, color.RGBA{
				R: clampUint8(float64(c.R) + noise),
				G: clampUint8(float64(c.G) + noise),
				B: clampUint8(float64(c.B) + noise),
				A: c.A,
			})
		}
	}

	return result
}

// ApplyChromaticAberration offsets color channels for a digital glitch effect.
func (pp *PostProcessor) ApplyChromaticAberration(img *image.RGBA, offset float64) *image.RGBA {
	if img == nil {
		return nil
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	result := image.NewRGBA(bounds)
	off := int(offset)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			result.SetRGBA(x, y, chromaticPixelHeadless(img, x, y, w, off))
		}
	}
	return result
}

// chromaticPixelHeadless computes the chromatic-shifted color at a pixel position.
func chromaticPixelHeadless(img *image.RGBA, x, y, w, off int) color.RGBA {
	rX := clampIntHeadless(x-off, 0, w-1)
	bX := clampIntHeadless(x+off, 0, w-1)

	cR := img.RGBAAt(rX, y)
	cG := img.RGBAAt(x, y)
	cB := img.RGBAAt(bX, y)

	return color.RGBA{R: cR.R, G: cG.G, B: cB.B, A: cG.A}
}

// clampIntHeadless clamps v to the range [min, max].
func clampIntHeadless(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// ApplySepia applies a warm sepia tone for vintage/dusty look.
func (pp *PostProcessor) ApplySepia(img *image.RGBA, intensity float64) *image.RGBA {
	if img == nil {
		return nil
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	result := image.NewRGBA(bounds)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := img.RGBAAt(x, y)
			rf := float64(c.R)
			gf := float64(c.G)
			bf := float64(c.B)

			// Standard sepia transformation
			sepiaR := rf*0.393 + gf*0.769 + bf*0.189
			sepiaG := rf*0.349 + gf*0.686 + bf*0.168
			sepiaB := rf*0.272 + gf*0.534 + bf*0.131

			// Blend original with sepia
			result.SetRGBA(x, y, color.RGBA{
				R: clampUint8(rf*(1-intensity) + sepiaR*intensity),
				G: clampUint8(gf*(1-intensity) + sepiaG*intensity),
				B: clampUint8(bf*(1-intensity) + sepiaB*intensity),
				A: c.A,
			})
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
