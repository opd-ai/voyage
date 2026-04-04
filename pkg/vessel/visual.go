package vessel

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// VisualVariant identifies a procedurally generated hull skin.
type VisualVariant int

const (
	// VisualVariantA is the first hull skin variant.
	VisualVariantA VisualVariant = iota
	// VisualVariantB is the second hull skin variant.
	VisualVariantB
	// VisualVariantC is the third hull skin variant.
	VisualVariantC
)

// AllVisualVariants returns all available visual variants.
func AllVisualVariants() []VisualVariant {
	return []VisualVariant{VisualVariantA, VisualVariantB, VisualVariantC}
}

// VisualVariantName returns the genre-appropriate name for a variant.
func VisualVariantName(v VisualVariant, genre engine.GenreID) string {
	names, ok := variantNames[genre]
	if !ok {
		names = variantNames[engine.GenreFantasy]
	}
	return names[v]
}

var variantNames = map[engine.GenreID]map[VisualVariant]string{
	engine.GenreFantasy: {
		VisualVariantA: "Oak Frame",
		VisualVariantB: "Iron Bound",
		VisualVariantC: "Gilded Trim",
	},
	engine.GenreScifi: {
		VisualVariantA: "Standard Hull",
		VisualVariantB: "Stealth Coating",
		VisualVariantC: "Armored Plating",
	},
	engine.GenreHorror: {
		VisualVariantA: "Rusty Shell",
		VisualVariantB: "Reinforced Frame",
		VisualVariantC: "Spiked Chassis",
	},
	engine.GenreCyberpunk: {
		VisualVariantA: "Matte Black",
		VisualVariantB: "Neon Accents",
		VisualVariantC: "Chrome Finish",
	},
	engine.GenrePostapoc: {
		VisualVariantA: "Salvage Patchwork",
		VisualVariantB: "Dust Runner",
		VisualVariantC: "War Paint",
	},
}

// VisualVariantDescription returns a description of the variant.
func VisualVariantDescription(v VisualVariant) string {
	descriptions := map[VisualVariant]string{
		VisualVariantA: "Classic design with balanced aesthetics",
		VisualVariantB: "Distinctive look with enhanced profile",
		VisualVariantC: "Premium appearance with unique details",
	}
	return descriptions[v]
}

// HullSkinParams defines the parameters for generating a hull skin sprite.
type HullSkinParams struct {
	Variant        VisualVariant
	Genre          engine.GenreID
	VesselType     VesselType
	PrimaryHue     float64 // 0-360 hue value
	SecondaryHue   float64 // 0-360 hue value
	PatternDensity float64 // 0-1 density of secondary color
	Symmetry       bool    // Whether to mirror the sprite
	DetailLevel    int     // Number of detail iterations
}

// HullSkinGenerator creates procedural hull skin parameters from a seed.
type HullSkinGenerator struct {
	gen   *seed.Generator
	genre engine.GenreID
}

// NewHullSkinGenerator creates a new hull skin generator.
func NewHullSkinGenerator(masterSeed int64, genre engine.GenreID) *HullSkinGenerator {
	return &HullSkinGenerator{
		gen:   seed.NewGenerator(masterSeed, "hullskin"),
		genre: genre,
	}
}

// SetGenre changes the generator's genre.
func (g *HullSkinGenerator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// Generate creates hull skin parameters for the given variant and type.
func (g *HullSkinGenerator) Generate(variant VisualVariant, vt VesselType) *HullSkinParams {
	primaryHue, secondaryHue := g.genreBaseHues()
	params := &HullSkinParams{
		Variant:    variant,
		Genre:      g.genre,
		VesselType: vt,
		Symmetry:   true,
	}
	g.applyVariantParams(params, variant, primaryHue, secondaryHue)
	return params
}

// genreBaseHues returns the base hue values for the current genre.
func (g *HullSkinGenerator) genreBaseHues() (primary, secondary float64) {
	switch g.genre {
	case engine.GenreScifi:
		return 210.0, 180.0 // Blue/cyan tones
	case engine.GenreHorror:
		return 30.0, 0.0 // Orange/red tones
	case engine.GenreCyberpunk:
		return 300.0, 180.0 // Magenta/cyan tones
	case engine.GenrePostapoc:
		return 40.0, 30.0 // Brown/tan tones
	default: // Fantasy
		return 120.0, 30.0 // Green/brown tones
	}
}

// applyVariantParams sets the variant-specific parameters.
func (g *HullSkinGenerator) applyVariantParams(params *HullSkinParams, variant VisualVariant, primaryHue, secondaryHue float64) {
	params.PrimaryHue = primaryHue + g.gen.RangeFloat64(-20, 20)
	params.SecondaryHue = secondaryHue + g.gen.RangeFloat64(-20, 20)

	switch variant {
	case VisualVariantA:
		params.PatternDensity = 0.2 + g.gen.RangeFloat64(0, 0.1)
		params.DetailLevel = 2
	case VisualVariantB:
		params.PatternDensity = 0.35 + g.gen.RangeFloat64(0, 0.1)
		params.DetailLevel = 3
	case VisualVariantC:
		params.PatternDensity = 0.5 + g.gen.RangeFloat64(0, 0.1)
		params.DetailLevel = 4
	}
	g.clampHues(params)
}

// clampHues ensures hue values wrap around properly.
func (g *HullSkinGenerator) clampHues(params *HullSkinParams) {
	params.PrimaryHue = wrapHue(params.PrimaryHue)
	params.SecondaryHue = wrapHue(params.SecondaryHue)
}

// wrapHue wraps a hue value to the 0-360 range.
func wrapHue(h float64) float64 {
	for h < 0 {
		h += 360
	}
	for h >= 360 {
		h -= 360
	}
	return h
}

// GenerateAll creates hull skin parameters for all variants of a vessel type.
func (g *HullSkinGenerator) GenerateAll(vt VesselType) []*HullSkinParams {
	skins := make([]*HullSkinParams, 0, 3)
	for _, v := range AllVisualVariants() {
		skins = append(skins, g.Generate(v, vt))
	}
	return skins
}

// VesselVisuals holds the selected visual variant for a vessel.
type VesselVisuals struct {
	Variant VisualVariant
	Params  *HullSkinParams
}

// NewVesselVisuals creates a new visuals holder with the given variant.
func NewVesselVisuals(variant VisualVariant, params *HullSkinParams) *VesselVisuals {
	return &VesselVisuals{
		Variant: variant,
		Params:  params,
	}
}

// VariantName returns the genre-appropriate name.
func (vv *VesselVisuals) VariantName() string {
	if vv.Params == nil {
		return VisualVariantName(vv.Variant, engine.GenreFantasy)
	}
	return VisualVariantName(vv.Variant, vv.Params.Genre)
}

// Description returns the variant description.
func (vv *VesselVisuals) Description() string {
	return VisualVariantDescription(vv.Variant)
}
