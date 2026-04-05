//go:build !headless

package rendering

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewGenreOverlay(t *testing.T) {
	overlay := NewGenreOverlay(engine.GenreFantasy)
	if overlay == nil {
		t.Fatal("NewGenreOverlay returned nil")
	}
	if overlay.Genre() != engine.GenreFantasy {
		t.Errorf("expected genre fantasy, got %s", overlay.Genre())
	}
}

func TestGenreOverlaySetGenre(t *testing.T) {
	overlay := NewGenreOverlay(engine.GenreFantasy)

	overlay.SetGenre(engine.GenreScifi)
	if overlay.Genre() != engine.GenreScifi {
		t.Errorf("expected genre scifi, got %s", overlay.Genre())
	}

	overlay.SetGenre(engine.GenreHorror)
	if overlay.Genre() != engine.GenreHorror {
		t.Errorf("expected genre horror, got %s", overlay.Genre())
	}
}

func TestApplyOverlay(t *testing.T) {
	overlay := NewGenreOverlay(engine.GenreFantasy)

	// Create a simple test image
	img := ebiten.NewImage(16, 16)
	img.Fill(color.RGBA{100, 100, 100, 255})

	result := overlay.ApplyOverlay(img)
	if result == nil {
		t.Fatal("ApplyOverlay returned nil")
	}

	bounds := result.Bounds()
	if bounds.Dx() != 16 || bounds.Dy() != 16 {
		t.Errorf("expected 16x16, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestGetOverlayParams(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		overlay := NewGenreOverlay(genre)
		params := overlay.getOverlayParams()

		// Check that all params are within reasonable ranges
		if params.tintR <= 0 || params.tintR > 2 {
			t.Errorf("genre %s: tintR out of range: %f", genre, params.tintR)
		}
		if params.tintG <= 0 || params.tintG > 2 {
			t.Errorf("genre %s: tintG out of range: %f", genre, params.tintG)
		}
		if params.tintB <= 0 || params.tintB > 2 {
			t.Errorf("genre %s: tintB out of range: %f", genre, params.tintB)
		}
		if params.saturation <= 0 || params.saturation > 2 {
			t.Errorf("genre %s: saturation out of range: %f", genre, params.saturation)
		}
		if params.brightness <= 0 || params.brightness > 2 {
			t.Errorf("genre %s: brightness out of range: %f", genre, params.brightness)
		}
		if params.noise < 0 || params.noise > 1 {
			t.Errorf("genre %s: noise out of range: %f", genre, params.noise)
		}
	}
}

func TestGetGenrePalette(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		palette := GetGenrePalette(genre)

		// Check that all palette colors are valid (non-zero alpha)
		if palette.Primary.A == 0 {
			t.Errorf("genre %s: Primary has zero alpha", genre)
		}
		if palette.Secondary.A == 0 {
			t.Errorf("genre %s: Secondary has zero alpha", genre)
		}
		if palette.Accent.A == 0 {
			t.Errorf("genre %s: Accent has zero alpha", genre)
		}
		if palette.Highlight.A == 0 {
			t.Errorf("genre %s: Highlight has zero alpha", genre)
		}
		if palette.Shadow.A == 0 {
			t.Errorf("genre %s: Shadow has zero alpha", genre)
		}
	}
}

func TestGetGenrePaletteFallback(t *testing.T) {
	// Test with invalid genre
	palette := GetGenrePalette("invalid")

	// Should return fantasy palette as default
	fantasyPalette := GetGenrePalette(engine.GenreFantasy)

	if palette.Primary != fantasyPalette.Primary {
		t.Error("invalid genre should fall back to fantasy palette")
	}
}

func TestGetAnimationColorScheme(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		scheme := GetAnimationColorScheme(genre)

		// Check all colors are populated
		if scheme.TileWaterBase.A == 0 {
			t.Errorf("genre %s: TileWaterBase has zero alpha", genre)
		}
		if scheme.TileWaterHighlight.A == 0 {
			t.Errorf("genre %s: TileWaterHighlight has zero alpha", genre)
		}
		if scheme.TileGrassBase.A == 0 {
			t.Errorf("genre %s: TileGrassBase has zero alpha", genre)
		}
		if scheme.PortraitPrimary.A == 0 {
			t.Errorf("genre %s: PortraitPrimary has zero alpha", genre)
		}
		if scheme.VesselHull.A == 0 {
			t.Errorf("genre %s: VesselHull has zero alpha", genre)
		}
		if scheme.LandmarkPrimary.A == 0 {
			t.Errorf("genre %s: LandmarkPrimary has zero alpha", genre)
		}
	}
}

func TestBlendRGBA(t *testing.T) {
	testCases := []struct {
		a, b     color.RGBA
		t        float64
		expected color.RGBA
	}{
		{
			color.RGBA{0, 0, 0, 255},
			color.RGBA{100, 100, 100, 255},
			0.5,
			color.RGBA{50, 50, 50, 255},
		},
		{
			color.RGBA{100, 100, 100, 255},
			color.RGBA{100, 100, 100, 255},
			0.5,
			color.RGBA{100, 100, 100, 255},
		},
		{
			color.RGBA{0, 0, 0, 255},
			color.RGBA{200, 200, 200, 255},
			0.0,
			color.RGBA{0, 0, 0, 255},
		},
		{
			color.RGBA{0, 0, 0, 255},
			color.RGBA{200, 200, 200, 255},
			1.0,
			color.RGBA{200, 200, 200, 255},
		},
	}

	for _, tc := range testCases {
		result := blendRGBA(tc.a, tc.b, tc.t)
		if result != tc.expected {
			t.Errorf("blendRGBA(%v, %v, %f) = %v, expected %v",
				tc.a, tc.b, tc.t, result, tc.expected)
		}
	}
}

func TestGenrePalettesExist(t *testing.T) {
	expectedGenres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range expectedGenres {
		_, ok := GenrePalettes[genre]
		if !ok {
			t.Errorf("missing palette for genre %s", genre)
		}
	}
}

func TestGenrePalettesAreDistinct(t *testing.T) {
	palettes := make(map[engine.GenreID]GenreColorPalette)
	for _, genre := range engine.AllGenres() {
		palettes[genre] = GetGenrePalette(genre)
	}

	// Check that primary colors are distinct across genres
	for g1, p1 := range palettes {
		for g2, p2 := range palettes {
			if g1 >= g2 {
				continue
			}
			if p1.Primary == p2.Primary {
				t.Errorf("genres %s and %s have identical primary colors", g1, g2)
			}
		}
	}
}
