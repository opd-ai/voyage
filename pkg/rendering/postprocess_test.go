//go:build headless

package rendering

import (
	"image"
	"image/color"
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func createTestImage(w, h int, col color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, col)
		}
	}
	return img
}

func TestNewPostProcessor(t *testing.T) {
	pp := NewPostProcessor(engine.GenreFantasy)
	if pp == nil {
		t.Fatal("NewPostProcessor returned nil")
	}

	if pp.Genre() != engine.GenreFantasy {
		t.Errorf("genre = %s, want fantasy", pp.Genre())
	}
}

func TestPostProcessorSetGenre(t *testing.T) {
	pp := NewPostProcessor(engine.GenreFantasy)

	// Test all genres configure without error
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		pp.SetGenre(genre)
		if pp.Genre() != genre {
			t.Errorf("genre = %s, want %s", pp.Genre(), genre)
		}
	}
}

func TestVignette(t *testing.T) {
	pp := NewPostProcessor(engine.GenreFantasy)

	// Create a white test image
	img := createTestImage(100, 100, color.RGBA{255, 255, 255, 255})

	result := pp.ApplyVignette(img, 0.5)
	if result == nil {
		t.Fatal("ApplyVignette returned nil")
	}

	// Center pixel should be brighter than corner pixel
	centerC := result.RGBAAt(50, 50)
	cornerC := result.RGBAAt(0, 0)

	centerBrightness := int(centerC.R) + int(centerC.G) + int(centerC.B)
	cornerBrightness := int(cornerC.R) + int(cornerC.G) + int(cornerC.B)

	if centerBrightness <= cornerBrightness {
		t.Error("center should be brighter than corner after vignette")
	}
}

func TestVignetteNil(t *testing.T) {
	pp := NewPostProcessor(engine.GenreFantasy)
	result := pp.ApplyVignette(nil, 0.5)
	if result != nil {
		t.Error("ApplyVignette(nil) should return nil")
	}
}

func TestScanlines(t *testing.T) {
	pp := NewPostProcessor(engine.GenreScifi)

	// Create a white test image
	img := createTestImage(100, 100, color.RGBA{255, 255, 255, 255})

	result := pp.ApplyScanlines(img, 2.0, 0.3)
	if result == nil {
		t.Fatal("ApplyScanlines returned nil")
	}

	// Every other row should be darkened
	scanlineC := result.RGBAAt(50, 0)
	regularC := result.RGBAAt(50, 1)

	// Scanline row should be darker
	if scanlineC.R >= regularC.R {
		t.Error("scanline row should be darker than regular row")
	}
}

func TestFilmGrain(t *testing.T) {
	pp := NewPostProcessor(engine.GenreHorror)

	// Create a gray test image
	img := createTestImage(100, 100, color.RGBA{128, 128, 128, 255})

	result := pp.ApplyFilmGrain(img, 12345, 0.2)
	if result == nil {
		t.Fatal("ApplyFilmGrain returned nil")
	}

	// Verify some variation exists (grain was applied)
	hasVariation := false
	firstC := result.RGBAAt(0, 0)
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			c := result.RGBAAt(x, y)
			if c.R != firstC.R || c.G != firstC.G || c.B != firstC.B {
				hasVariation = true
				break
			}
		}
		if hasVariation {
			break
		}
	}

	if !hasVariation {
		t.Error("film grain should add some variation")
	}
}

func TestFilmGrainDeterminism(t *testing.T) {
	pp := NewPostProcessor(engine.GenreHorror)

	img1 := createTestImage(50, 50, color.RGBA{128, 128, 128, 255})
	img2 := createTestImage(50, 50, color.RGBA{128, 128, 128, 255})

	// Same seed should produce same result
	result1 := pp.ApplyFilmGrain(img1, 12345, 0.2)
	result2 := pp.ApplyFilmGrain(img2, 12345, 0.2)

	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			c1 := result1.RGBAAt(x, y)
			c2 := result2.RGBAAt(x, y)
			if c1 != c2 {
				t.Errorf("same seed should produce same result at (%d,%d)", x, y)
				return
			}
		}
	}
}

func TestChromaticAberration(t *testing.T) {
	pp := NewPostProcessor(engine.GenreCyberpunk)

	// Create a colored test image
	img := createTestImage(100, 100, color.RGBA{255, 128, 64, 255})

	result := pp.ApplyChromaticAberration(img, 2.0)
	if result == nil {
		t.Fatal("ApplyChromaticAberration returned nil")
	}

	// Verify the result has the expected dimensions
	bounds := result.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("result dimensions = (%d,%d), want (100,100)", bounds.Dx(), bounds.Dy())
	}
}

func TestSepia(t *testing.T) {
	pp := NewPostProcessor(engine.GenrePostapoc)

	// Create a colored test image
	img := createTestImage(100, 100, color.RGBA{100, 150, 200, 255})

	result := pp.ApplySepia(img, 0.8)
	if result == nil {
		t.Fatal("ApplySepia returned nil")
	}

	// Sepia should warm the colors (increase red relative to blue)
	c := result.RGBAAt(50, 50)
	// In sepia, R > G > B typically for most input colors
	if c.R < c.B {
		t.Error("sepia should make image warmer (more red than blue)")
	}
}

func TestApplyAllEffects(t *testing.T) {
	// Test that Apply() works for each genre
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		pp := NewPostProcessor(genre)
		img := createTestImage(50, 50, color.RGBA{128, 128, 128, 255})

		result := pp.Apply(img, 12345)
		if result == nil {
			t.Errorf("Apply() returned nil for genre %s", genre)
		}
	}
}

func TestApplyNil(t *testing.T) {
	pp := NewPostProcessor(engine.GenreFantasy)
	result := pp.Apply(nil, 12345)
	if result != nil {
		t.Error("Apply(nil) should return nil")
	}
}

func TestGenrePostProcessing(t *testing.T) {
	// Test that different genres produce visually different results
	img := createTestImage(50, 50, color.RGBA{128, 128, 128, 255})

	ppFantasy := NewPostProcessor(engine.GenreFantasy)
	ppCyberpunk := NewPostProcessor(engine.GenreCyberpunk)

	resultFantasy := ppFantasy.Apply(img, 12345)
	resultCyberpunk := ppCyberpunk.Apply(img, 12345)

	// Results should be different due to different effect chains
	same := true
	for y := 0; y < 50; y++ {
		for x := 0; x < 50; x++ {
			c1 := resultFantasy.RGBAAt(x, y)
			c2 := resultCyberpunk.RGBAAt(x, y)
			if c1 != c2 {
				same = false
				break
			}
		}
		if !same {
			break
		}
	}

	if same {
		t.Error("different genres should produce different visual results")
	}
}

func TestSetters(t *testing.T) {
	pp := NewPostProcessor(engine.GenreFantasy)

	pp.SetVignetteIntensity(0.8)
	pp.SetScanlinesEnabled(true)
	pp.SetFilmGrainEnabled(true)
	pp.SetChromaticEnabled(true)
	pp.SetSepiaEnabled(true)

	// Just verify no panics and methods work
	img := createTestImage(50, 50, color.RGBA{128, 128, 128, 255})
	result := pp.Apply(img, 12345)
	if result == nil {
		t.Error("Apply after setters should not return nil")
	}
}
