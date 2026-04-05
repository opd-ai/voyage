//go:build headless

package rendering

import (
	"image"
	"image/color"

	"github.com/opd-ai/voyage/pkg/engine"
)

// PostProcessor applies post-processing effects to rendered images.
type PostProcessor struct {
	genre        engine.GenreID
	vignetteOn   bool
	vignetteInt  float64
	scanlinesOn  bool
	scanlinesDen float64
	filmGrainOn  bool
	filmGrainInt float64
	chromaticOn  bool
	chromaticOff float64
	sepiaOn      bool
	sepiaInt     float64
}

// NewPostProcessor creates a new post processor.
func NewPostProcessor(genre engine.GenreID) *PostProcessor {
	pp := &PostProcessor{
		vignetteInt:  0.3,
		scanlinesDen: 2.0,
		filmGrainInt: 0.15,
		chromaticOff: 2.0,
		sepiaInt:     0.5,
	}
	pp.SetGenre(genre)
	return pp
}

// SetGenre configures the post processor for a specific genre.
func (pp *PostProcessor) SetGenre(genre engine.GenreID) {
	pp.genre = genre
	// Reset all effects
	pp.vignetteOn = false
	pp.scanlinesOn = false
	pp.filmGrainOn = false
	pp.chromaticOn = false
	pp.sepiaOn = false

	// Enable genre-specific effects
	switch genre {
	case engine.GenreFantasy:
		pp.vignetteOn = true
		pp.vignetteInt = 0.2
	case engine.GenreScifi:
		pp.vignetteOn = true
		pp.vignetteInt = 0.3
		pp.scanlinesOn = true
	case engine.GenreHorror:
		pp.vignetteOn = true
		pp.vignetteInt = 0.5
		pp.filmGrainOn = true
	case engine.GenreCyberpunk:
		pp.vignetteOn = true
		pp.vignetteInt = 0.4
		pp.scanlinesOn = true
		pp.chromaticOn = true
	case engine.GenrePostapoc:
		pp.vignetteOn = true
		pp.vignetteInt = 0.35
		pp.sepiaOn = true
	}
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

	result := img

	// Apply effects in order
	if pp.vignetteOn {
		result = pp.ApplyVignette(result, pp.vignetteInt)
	}
	if pp.scanlinesOn {
		result = pp.ApplyScanlines(result, pp.scanlinesDen, 0.15)
	}
	if pp.filmGrainOn {
		result = pp.ApplyFilmGrain(result, seed, pp.filmGrainInt)
	}
	if pp.chromaticOn {
		result = pp.ApplyChromaticAberration(result, pp.chromaticOff)
	}
	if pp.sepiaOn {
		result = pp.ApplySepia(result, pp.sepiaInt)
	}

	return result
}

// ApplyVignette darkens the edges of an image to focus attention on center.
func (pp *PostProcessor) ApplyVignette(img *image.RGBA, intensity float64) *image.RGBA {
	if img == nil {
		return nil
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	result := image.NewRGBA(bounds)

	// Apply vignette by darkening pixels based on distance from center
	centerX := float64(w) / 2
	centerY := float64(h) / 2
	maxDist := centerX
	if centerY > maxDist {
		maxDist = centerY
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			dx := float64(x) - centerX
			dy := float64(y) - centerY
			dist := (dx*dx + dy*dy) / (maxDist * maxDist)

			// Calculate darkening factor
			factor := 1.0 - (dist * intensity)
			if factor < 0 {
				factor = 0
			}

			c := img.RGBAAt(x, y)
			result.SetRGBA(x, y, color.RGBA{
				R: uint8(float64(c.R) * factor),
				G: uint8(float64(c.G) * factor),
				B: uint8(float64(c.B) * factor),
				A: c.A,
			})
		}
	}

	return result
}

// ApplyScanlines adds horizontal scanline effect for retro/sci-fi feel.
func (pp *PostProcessor) ApplyScanlines(img *image.RGBA, density, alpha float64) *image.RGBA {
	if img == nil {
		return nil
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	result := image.NewRGBA(bounds)

	// Copy original
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			result.Set(x, y, img.At(x, y))
		}
	}

	// Add scanlines every N pixels
	spacing := int(density)
	if spacing < 1 {
		spacing = 1
	}

	for y := 0; y < h; y += spacing {
		for x := 0; x < w; x++ {
			c := img.RGBAAt(x, y)
			factor := 1.0 - alpha
			result.SetRGBA(x, y, color.RGBA{
				R: uint8(float64(c.R) * factor),
				G: uint8(float64(c.G) * factor),
				B: uint8(float64(c.B) * factor),
				A: c.A,
			})
		}
	}

	return result
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
			// Red channel shifted left
			rX := x - off
			if rX < 0 {
				rX = 0
			}
			cR := img.RGBAAt(rX, y)

			// Green channel stays centered
			cG := img.RGBAAt(x, y)

			// Blue channel shifted right
			bX := x + off
			if bX >= w {
				bX = w - 1
			}
			cB := img.RGBAAt(bX, y)

			result.SetRGBA(x, y, color.RGBA{
				R: cR.R,
				G: cG.G,
				B: cB.B,
				A: cG.A,
			})
		}
	}

	return result
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

// clampUint8 clamps a float64 to uint8 range.
func clampUint8(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}

// SetVignetteIntensity sets the vignette darkening strength.
func (pp *PostProcessor) SetVignetteIntensity(intensity float64) {
	pp.vignetteInt = intensity
}

// SetScanlinesEnabled toggles scanline effect.
func (pp *PostProcessor) SetScanlinesEnabled(enabled bool) {
	pp.scanlinesOn = enabled
}

// SetFilmGrainEnabled toggles film grain effect.
func (pp *PostProcessor) SetFilmGrainEnabled(enabled bool) {
	pp.filmGrainOn = enabled
}

// SetChromaticEnabled toggles chromatic aberration effect.
func (pp *PostProcessor) SetChromaticEnabled(enabled bool) {
	pp.chromaticOn = enabled
}

// SetSepiaEnabled toggles sepia effect.
func (pp *PostProcessor) SetSepiaEnabled(enabled bool) {
	pp.sepiaOn = enabled
}
