//go:build headless

package rendering

import (
	"image/color"

	"github.com/opd-ai/voyage/pkg/engine"
)

// GenreOverlay applies genre-specific visual overlays to sprites (headless stub).
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

// GenrePalettes provides pre-defined color palettes for each genre.
var GenrePalettes = map[engine.GenreID]GenreColorPalette{
	engine.GenreFantasy: {
		Primary:    color.RGBA{80, 120, 80, 255},
		Secondary:  color.RGBA{200, 180, 100, 255},
		Accent:     color.RGBA{160, 120, 80, 255},
		Highlight:  color.RGBA{255, 240, 200, 255},
		Shadow:     color.RGBA{40, 50, 40, 255},
		Skin:       color.RGBA{210, 170, 140, 255},
		Fire:       color.RGBA{255, 150, 50, 255},
		Water:      color.RGBA{50, 100, 150, 255},
		Vegetation: color.RGBA{60, 100, 50, 255},
		Stone:      color.RGBA{100, 95, 90, 255},
		Metal:      color.RGBA{140, 130, 110, 255},
		Smoke:      color.RGBA{120, 115, 110, 255},
	},
	engine.GenreScifi: {
		Primary:    color.RGBA{50, 100, 180, 255},
		Secondary:  color.RGBA{0, 200, 200, 255},
		Accent:     color.RGBA{150, 50, 200, 255},
		Highlight:  color.RGBA{220, 240, 255, 255},
		Shadow:     color.RGBA{20, 30, 50, 255},
		Skin:       color.RGBA{180, 200, 200, 255},
		Fire:       color.RGBA{100, 200, 255, 255},
		Water:      color.RGBA{30, 80, 120, 255},
		Vegetation: color.RGBA{50, 80, 100, 255},
		Stone:      color.RGBA{80, 85, 95, 255},
		Metal:      color.RGBA{150, 160, 180, 255},
		Smoke:      color.RGBA{100, 120, 140, 255},
	},
	engine.GenreHorror: {
		Primary:    color.RGBA{100, 50, 50, 255},
		Secondary:  color.RGBA{60, 50, 40, 255},
		Accent:     color.RGBA{150, 40, 40, 255},
		Highlight:  color.RGBA{200, 180, 170, 255},
		Shadow:     color.RGBA{30, 25, 25, 255},
		Skin:       color.RGBA{170, 160, 150, 255},
		Fire:       color.RGBA{200, 80, 30, 255},
		Water:      color.RGBA{40, 50, 50, 255},
		Vegetation: color.RGBA{40, 50, 40, 255},
		Stone:      color.RGBA{70, 65, 60, 255},
		Metal:      color.RGBA{100, 90, 80, 255},
		Smoke:      color.RGBA{80, 70, 70, 255},
	},
	engine.GenreCyberpunk: {
		Primary:    color.RGBA{255, 0, 100, 255},
		Secondary:  color.RGBA{0, 200, 255, 255},
		Accent:     color.RGBA{255, 255, 0, 255},
		Highlight:  color.RGBA{255, 255, 255, 255},
		Shadow:     color.RGBA{20, 15, 30, 255},
		Skin:       color.RGBA{200, 180, 170, 255},
		Fire:       color.RGBA{255, 100, 0, 255},
		Water:      color.RGBA{20, 40, 50, 255},
		Vegetation: color.RGBA{30, 60, 30, 255},
		Stone:      color.RGBA{50, 50, 60, 255},
		Metal:      color.RGBA{80, 80, 100, 255},
		Smoke:      color.RGBA{60, 60, 80, 255},
	},
	engine.GenrePostapoc: {
		Primary:    color.RGBA{150, 100, 60, 255},
		Secondary:  color.RGBA{100, 80, 60, 255},
		Accent:     color.RGBA{180, 140, 50, 255},
		Highlight:  color.RGBA{220, 200, 170, 255},
		Shadow:     color.RGBA{40, 35, 30, 255},
		Skin:       color.RGBA{190, 160, 130, 255},
		Fire:       color.RGBA{255, 120, 30, 255},
		Water:      color.RGBA{70, 80, 60, 255},
		Vegetation: color.RGBA{80, 70, 50, 255},
		Stone:      color.RGBA{90, 85, 75, 255},
		Metal:      color.RGBA{120, 100, 80, 255},
		Smoke:      color.RGBA{100, 90, 80, 255},
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
