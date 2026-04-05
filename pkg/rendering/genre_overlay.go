//go:build !headless

package rendering

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/opd-ai/voyage/pkg/engine"
)

// GenreOverlay applies genre-specific visual overlays to sprites.
type GenreOverlay struct {
	genre engine.GenreID
}

// NewGenreOverlay creates a new genre overlay applicator.
func NewGenreOverlay(genre engine.GenreID) *GenreOverlay {
	return &GenreOverlay{genre: genre}
}

// SetGenre changes the overlay's genre.
func (go_ *GenreOverlay) SetGenre(genre engine.GenreID) {
	go_.genre = genre
}

// Genre returns the current genre.
func (go_ *GenreOverlay) Genre() engine.GenreID {
	return go_.genre
}

// ApplyOverlay applies the genre-specific color overlay to an image.
func (go_ *GenreOverlay) ApplyOverlay(img *ebiten.Image) *ebiten.Image {
	bounds := img.Bounds()
	result := ebiten.NewImage(bounds.Dx(), bounds.Dy())

	op := &ebiten.DrawImageOptions{}
	overlay := go_.getOverlayParams()

	// Apply color tinting based on genre
	op.ColorScale.Scale(float32(overlay.tintR), float32(overlay.tintG), float32(overlay.tintB), 1.0)

	result.DrawImage(img, op)

	// Apply additional genre effects
	go_.applyGenreEffects(result, overlay)

	return result
}

// genreOverlayParams holds the overlay parameters for a genre.
type genreOverlayParams struct {
	tintR, tintG, tintB float64
	saturation          float64
	brightness          float64
	noise               float64
}

// getOverlayParams returns the overlay parameters for the current genre.
func (go_ *GenreOverlay) getOverlayParams() genreOverlayParams {
	switch go_.genre {
	case engine.GenreScifi:
		return genreOverlayParams{
			tintR:      0.9,
			tintG:      1.0,
			tintB:      1.1,
			saturation: 0.9,
			brightness: 1.05,
			noise:      0.02,
		}
	case engine.GenreHorror:
		return genreOverlayParams{
			tintR:      1.1,
			tintG:      0.9,
			tintB:      0.9,
			saturation: 0.7,
			brightness: 0.85,
			noise:      0.05,
		}
	case engine.GenreCyberpunk:
		return genreOverlayParams{
			tintR:      1.1,
			tintG:      0.95,
			tintB:      1.15,
			saturation: 1.2,
			brightness: 1.1,
			noise:      0.03,
		}
	case engine.GenrePostapoc:
		return genreOverlayParams{
			tintR:      1.15,
			tintG:      1.1,
			tintB:      0.9,
			saturation: 0.6,
			brightness: 0.9,
			noise:      0.08,
		}
	default: // GenreFantasy
		return genreOverlayParams{
			tintR:      1.05,
			tintG:      1.0,
			tintB:      0.95,
			saturation: 0.95,
			brightness: 1.0,
			noise:      0.01,
		}
	}
}

// applyGenreEffects applies additional genre-specific effects.
func (go_ *GenreOverlay) applyGenreEffects(img *ebiten.Image, params genreOverlayParams) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// Apply saturation and brightness adjustments
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()
			if a == 0 {
				continue
			}

			// Convert to float
			fr := float64(r>>8) / 255.0
			fg := float64(g>>8) / 255.0
			fb := float64(b>>8) / 255.0

			// Apply saturation
			gray := 0.299*fr + 0.587*fg + 0.114*fb
			fr = gray + params.saturation*(fr-gray)
			fg = gray + params.saturation*(fg-gray)
			fb = gray + params.saturation*(fb-gray)

			// Apply brightness
			fr *= params.brightness
			fg *= params.brightness
			fb *= params.brightness

			// Clamp
			fr = clampFloat(fr, 0, 1)
			fg = clampFloat(fg, 0, 1)
			fb = clampFloat(fb, 0, 1)

			img.Set(x, y, color.RGBA{
				R: uint8(fr * 255),
				G: uint8(fg * 255),
				B: uint8(fb * 255),
				A: uint8(a >> 8),
			})
		}
	}
}

// GenrePalettes provides pre-defined color palettes for each genre.
var GenrePalettes = map[engine.GenreID]GenreColorPalette{
	engine.GenreFantasy: {
		Primary:    color.RGBA{80, 120, 80, 255},   // Forest green
		Secondary:  color.RGBA{200, 180, 100, 255}, // Golden
		Accent:     color.RGBA{160, 120, 80, 255},  // Brown
		Highlight:  color.RGBA{255, 240, 200, 255}, // Warm white
		Shadow:     color.RGBA{40, 50, 40, 255},    // Dark green
		Skin:       color.RGBA{210, 170, 140, 255}, // Warm skin
		Fire:       color.RGBA{255, 150, 50, 255},  // Orange flame
		Water:      color.RGBA{50, 100, 150, 255},  // Blue water
		Vegetation: color.RGBA{60, 100, 50, 255},   // Green
		Stone:      color.RGBA{100, 95, 90, 255},   // Gray stone
		Metal:      color.RGBA{140, 130, 110, 255}, // Bronze
		Smoke:      color.RGBA{120, 115, 110, 255}, // Gray smoke
	},
	engine.GenreScifi: {
		Primary:    color.RGBA{50, 100, 180, 255},  // Tech blue
		Secondary:  color.RGBA{0, 200, 200, 255},   // Cyan
		Accent:     color.RGBA{150, 50, 200, 255},  // Purple
		Highlight:  color.RGBA{220, 240, 255, 255}, // Cool white
		Shadow:     color.RGBA{20, 30, 50, 255},    // Dark blue
		Skin:       color.RGBA{180, 200, 200, 255}, // Cool skin
		Fire:       color.RGBA{100, 200, 255, 255}, // Plasma blue
		Water:      color.RGBA{30, 80, 120, 255},   // Dark water
		Vegetation: color.RGBA{50, 80, 100, 255},   // Bioluminescent
		Stone:      color.RGBA{80, 85, 95, 255},    // Metal hull
		Metal:      color.RGBA{150, 160, 180, 255}, // Chrome
		Smoke:      color.RGBA{100, 120, 140, 255}, // Exhaust
	},
	engine.GenreHorror: {
		Primary:    color.RGBA{100, 50, 50, 255},   // Blood red
		Secondary:  color.RGBA{60, 50, 40, 255},    // Rot brown
		Accent:     color.RGBA{150, 40, 40, 255},   // Bright red
		Highlight:  color.RGBA{200, 180, 170, 255}, // Sickly white
		Shadow:     color.RGBA{30, 25, 25, 255},    // Deep black
		Skin:       color.RGBA{170, 160, 150, 255}, // Pale skin
		Fire:       color.RGBA{200, 80, 30, 255},   // Dying embers
		Water:      color.RGBA{40, 50, 50, 255},    // Murky
		Vegetation: color.RGBA{40, 50, 40, 255},    // Dead plants
		Stone:      color.RGBA{70, 65, 60, 255},    // Crumbling
		Metal:      color.RGBA{100, 90, 80, 255},   // Rusted
		Smoke:      color.RGBA{80, 70, 70, 255},    // Ash
	},
	engine.GenreCyberpunk: {
		Primary:    color.RGBA{255, 0, 100, 255},   // Neon pink
		Secondary:  color.RGBA{0, 200, 255, 255},   // Neon cyan
		Accent:     color.RGBA{255, 255, 0, 255},   // Neon yellow
		Highlight:  color.RGBA{255, 255, 255, 255}, // Pure white
		Shadow:     color.RGBA{20, 15, 30, 255},    // Deep purple
		Skin:       color.RGBA{200, 180, 170, 255}, // Varied skin
		Fire:       color.RGBA{255, 100, 0, 255},   // Industrial
		Water:      color.RGBA{20, 40, 50, 255},    // Polluted
		Vegetation: color.RGBA{30, 60, 30, 255},    // Rare plants
		Stone:      color.RGBA{50, 50, 60, 255},    // Concrete
		Metal:      color.RGBA{80, 80, 100, 255},   // Chrome
		Smoke:      color.RGBA{60, 60, 80, 255},    // Smog
	},
	engine.GenrePostapoc: {
		Primary:    color.RGBA{150, 100, 60, 255},  // Rust
		Secondary:  color.RGBA{100, 80, 60, 255},   // Dirt
		Accent:     color.RGBA{180, 140, 50, 255},  // Rad warning
		Highlight:  color.RGBA{220, 200, 170, 255}, // Dusty white
		Shadow:     color.RGBA{40, 35, 30, 255},    // Dark earth
		Skin:       color.RGBA{190, 160, 130, 255}, // Weathered
		Fire:       color.RGBA{255, 120, 30, 255},  // Barrel fire
		Water:      color.RGBA{70, 80, 60, 255},    // Irradiated
		Vegetation: color.RGBA{80, 70, 50, 255},    // Dead grass
		Stone:      color.RGBA{90, 85, 75, 255},    // Rubble
		Metal:      color.RGBA{120, 100, 80, 255},  // Scrap
		Smoke:      color.RGBA{100, 90, 80, 255},   // Dust
	},
}

// GenreColorPalette defines the color scheme for a genre.
type GenreColorPalette struct {
	Primary    color.RGBA
	Secondary  color.RGBA
	Accent     color.RGBA
	Highlight  color.RGBA
	Shadow     color.RGBA
	Skin       color.RGBA
	Fire       color.RGBA
	Water      color.RGBA
	Vegetation color.RGBA
	Stone      color.RGBA
	Metal      color.RGBA
	Smoke      color.RGBA
}

// GetGenrePalette returns the color palette for a genre.
func GetGenrePalette(genre engine.GenreID) GenreColorPalette {
	if palette, ok := GenrePalettes[genre]; ok {
		return palette
	}
	return GenrePalettes[engine.GenreFantasy]
}

// AnimationColorScheme returns appropriate colors for animations based on genre.
type AnimationColorScheme struct {
	TileWaterBase      color.RGBA
	TileWaterHighlight color.RGBA
	TileGrassBase      color.RGBA
	TileGrassTip       color.RGBA
	TileFireBase       color.RGBA
	TileFireBright     color.RGBA
	PortraitPrimary    color.RGBA
	PortraitSecondary  color.RGBA
	PortraitSkin       color.RGBA
	VesselHull         color.RGBA
	VesselAccent       color.RGBA
	LandmarkPrimary    color.RGBA
	LandmarkSecondary  color.RGBA
}

// GetAnimationColorScheme returns colors for all animation types based on genre.
func GetAnimationColorScheme(genre engine.GenreID) AnimationColorScheme {
	p := GetGenrePalette(genre)
	return AnimationColorScheme{
		TileWaterBase:      p.Water,
		TileWaterHighlight: blendRGBA(p.Water, p.Highlight, 0.3),
		TileGrassBase:      p.Vegetation,
		TileGrassTip:       blendRGBA(p.Vegetation, p.Highlight, 0.4),
		TileFireBase:       blendRGBA(p.Fire, p.Shadow, 0.3),
		TileFireBright:     p.Fire,
		PortraitPrimary:    p.Primary,
		PortraitSecondary:  p.Secondary,
		PortraitSkin:       p.Skin,
		VesselHull:         p.Metal,
		VesselAccent:       p.Accent,
		LandmarkPrimary:    p.Stone,
		LandmarkSecondary:  p.Highlight,
	}
}

// blendRGBA blends two RGBA colors.
func blendRGBA(a, b color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		R: uint8(float64(a.R)*(1-t) + float64(b.R)*t),
		G: uint8(float64(a.G)*(1-t) + float64(b.G)*t),
		B: uint8(float64(a.B)*(1-t) + float64(b.B)*t),
		A: uint8(float64(a.A)*(1-t) + float64(b.A)*t),
	}
}
