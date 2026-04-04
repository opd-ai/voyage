package vessel

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// InsigniaShape identifies the base shape of a vessel insignia.
type InsigniaShape int

const (
	// InsigniaShapeShield is a traditional shield shape.
	InsigniaShapeShield InsigniaShape = iota
	// InsigniaShapeCircle is a circular emblem.
	InsigniaShapeCircle
	// InsigniaShapeDiamond is a diamond/rhombus shape.
	InsigniaShapeDiamond
	// InsigniaShapeHexagon is a hexagonal shape.
	InsigniaShapeHexagon
	// InsigniaShapeBanner is a vertical banner shape.
	InsigniaShapeBanner
)

// AllInsigniaShapes returns all available insignia shapes.
func AllInsigniaShapes() []InsigniaShape {
	return []InsigniaShape{
		InsigniaShapeShield,
		InsigniaShapeCircle,
		InsigniaShapeDiamond,
		InsigniaShapeHexagon,
		InsigniaShapeBanner,
	}
}

// InsigniaPattern identifies the internal pattern of the insignia.
type InsigniaPattern int

const (
	// InsigniaPatternSolid is a single solid color.
	InsigniaPatternSolid InsigniaPattern = iota
	// InsigniaPatternHorizontal has horizontal stripes.
	InsigniaPatternHorizontal
	// InsigniaPatternVertical has vertical stripes.
	InsigniaPatternVertical
	// InsigniaPatternDiagonal has diagonal stripes.
	InsigniaPatternDiagonal
	// InsigniaPatternQuarters is divided into four sections.
	InsigniaPatternQuarters
	// InsigniaPatternChevron has a chevron/arrow pattern.
	InsigniaPatternChevron
)

// AllInsigniaPatterns returns all available insignia patterns.
func AllInsigniaPatterns() []InsigniaPattern {
	return []InsigniaPattern{
		InsigniaPatternSolid,
		InsigniaPatternHorizontal,
		InsigniaPatternVertical,
		InsigniaPatternDiagonal,
		InsigniaPatternQuarters,
		InsigniaPatternChevron,
	}
}

// InsigniaSymbol identifies the central symbol on the insignia.
type InsigniaSymbol int

const (
	// InsigniaSymbolNone has no central symbol.
	InsigniaSymbolNone InsigniaSymbol = iota
	// InsigniaSymbolStar is a star shape.
	InsigniaSymbolStar
	// InsigniaSymbolCross is a cross shape.
	InsigniaSymbolCross
	// InsigniaSymbolCircle is a small circle.
	InsigniaSymbolCircle
	// InsigniaSymbolTriangle is a triangle.
	InsigniaSymbolTriangle
	// InsigniaSymbolArrow is an arrow or direction indicator.
	InsigniaSymbolArrow
	// InsigniaSymbolWave is a wave or flowing pattern.
	InsigniaSymbolWave
	// InsigniaSymbolGear is a gear/cog shape.
	InsigniaSymbolGear
)

// AllInsigniaSymbols returns all available insignia symbols.
func AllInsigniaSymbols() []InsigniaSymbol {
	return []InsigniaSymbol{
		InsigniaSymbolNone,
		InsigniaSymbolStar,
		InsigniaSymbolCross,
		InsigniaSymbolCircle,
		InsigniaSymbolTriangle,
		InsigniaSymbolArrow,
		InsigniaSymbolWave,
		InsigniaSymbolGear,
	}
}

// Insignia represents a procedurally generated vessel emblem.
type Insignia struct {
	Shape          InsigniaShape
	Pattern        InsigniaPattern
	Symbol         InsigniaSymbol
	PrimaryHue     float64 // 0-360 hue value
	SecondaryHue   float64 // 0-360 hue value
	AccentHue      float64 // 0-360 hue value (for symbol)
	Saturation     float64 // 0-1 saturation level
	BorderWidth    float64 // 0-1 relative border thickness
	SymbolScale    float64 // 0.3-0.8 scale of central symbol
	PatternDensity float64 // 0-1 density of pattern elements
	HasBorder      bool
	HasShadow      bool
	Genre          engine.GenreID
	Name           string // Procedurally generated insignia name
}

// InsigniaGenerator creates procedural insignias from a seed.
type InsigniaGenerator struct {
	gen   *seed.Generator
	genre engine.GenreID
}

// NewInsigniaGenerator creates a new insignia generator.
func NewInsigniaGenerator(masterSeed int64, genre engine.GenreID) *InsigniaGenerator {
	return &InsigniaGenerator{
		gen:   seed.NewGenerator(masterSeed, "insignia"),
		genre: genre,
	}
}

// SetGenre changes the generator's genre.
func (g *InsigniaGenerator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// Genre returns the current genre.
func (g *InsigniaGenerator) Genre() engine.GenreID {
	return g.genre
}

// Generate creates a new procedural insignia.
func (g *InsigniaGenerator) Generate() *Insignia {
	shape := seed.Choice(g.gen, AllInsigniaShapes())
	pattern := seed.Choice(g.gen, AllInsigniaPatterns())
	symbol := seed.Choice(g.gen, AllInsigniaSymbols())

	primaryHue, secondaryHue, accentHue := g.genreBaseHues()

	insignia := &Insignia{
		Shape:          shape,
		Pattern:        pattern,
		Symbol:         symbol,
		PrimaryHue:     wrapHue(primaryHue + g.gen.RangeFloat64(-30, 30)),
		SecondaryHue:   wrapHue(secondaryHue + g.gen.RangeFloat64(-30, 30)),
		AccentHue:      wrapHue(accentHue + g.gen.RangeFloat64(-20, 20)),
		Saturation:     g.gen.RangeFloat64(0.5, 1.0),
		BorderWidth:    g.gen.RangeFloat64(0.05, 0.15),
		SymbolScale:    g.gen.RangeFloat64(0.3, 0.7),
		PatternDensity: g.gen.RangeFloat64(0.3, 0.8),
		HasBorder:      g.gen.Chance(0.7),
		HasShadow:      g.gen.Chance(0.5),
		Genre:          g.genre,
	}

	insignia.Name = g.generateName(insignia)
	return insignia
}

// genreBaseHues returns base hue values appropriate for the current genre.
func (g *InsigniaGenerator) genreBaseHues() (primary, secondary, accent float64) {
	switch g.genre {
	case engine.GenreScifi:
		return 210.0, 180.0, 60.0 // Blue/cyan with gold accent
	case engine.GenreHorror:
		return 0.0, 30.0, 120.0 // Red/orange with sickly green
	case engine.GenreCyberpunk:
		return 300.0, 180.0, 60.0 // Magenta/cyan with yellow
	case engine.GenrePostapoc:
		return 30.0, 45.0, 0.0 // Brown/tan with rust red
	default: // Fantasy
		return 45.0, 120.0, 210.0 // Gold/green with blue
	}
}

// generateName creates a procedural name for the insignia.
func (g *InsigniaGenerator) generateName(i *Insignia) string {
	prefixes := g.genrePrefixes()
	suffixes := g.genreSuffixes()
	descriptors := g.shapeDescriptors(i.Shape)

	prefix := seed.Choice(g.gen, prefixes)
	suffix := seed.Choice(g.gen, suffixes)
	descriptor := seed.Choice(g.gen, descriptors)

	return prefix + " " + descriptor + " " + suffix
}

// genrePrefixes returns name prefixes appropriate for the genre.
func (g *InsigniaGenerator) genrePrefixes() []string {
	switch g.genre {
	case engine.GenreScifi:
		return []string{"Alpha", "Nova", "Stellar", "Quantum", "Void", "Astral"}
	case engine.GenreHorror:
		return []string{"Blood", "Shadow", "Dread", "Grim", "Hollow", "Cursed"}
	case engine.GenreCyberpunk:
		return []string{"Neon", "Chrome", "Cyber", "Digital", "Ghost", "Razor"}
	case engine.GenrePostapoc:
		return []string{"Rust", "Ash", "Iron", "Dust", "Scrap", "Storm"}
	default: // Fantasy
		return []string{"Golden", "Silver", "Royal", "Ancient", "Noble", "Mystic"}
	}
}

// genreSuffixes returns name suffixes appropriate for the genre.
func (g *InsigniaGenerator) genreSuffixes() []string {
	switch g.genre {
	case engine.GenreScifi:
		return []string{"Corps", "Fleet", "Division", "Command", "Syndicate"}
	case engine.GenreHorror:
		return []string{"Horde", "Legion", "Cult", "Order", "Covenant"}
	case engine.GenreCyberpunk:
		return []string{"Collective", "Network", "Syndicate", "Crew", "Clan"}
	case engine.GenrePostapoc:
		return []string{"Tribe", "Clan", "Pack", "Caravan", "Convoy"}
	default: // Fantasy
		return []string{"Company", "Guild", "Order", "House", "Banner"}
	}
}

// shapeDescriptors returns descriptive words based on insignia shape.
func (g *InsigniaGenerator) shapeDescriptors(shape InsigniaShape) []string {
	switch shape {
	case InsigniaShapeShield:
		return []string{"Shield", "Aegis", "Guard", "Bulwark", "Ward"}
	case InsigniaShapeCircle:
		return []string{"Ring", "Circle", "Wheel", "Orb", "Sun"}
	case InsigniaShapeDiamond:
		return []string{"Diamond", "Crystal", "Gem", "Prism", "Shard"}
	case InsigniaShapeHexagon:
		return []string{"Hex", "Cell", "Hive", "Grid", "Matrix"}
	case InsigniaShapeBanner:
		return []string{"Banner", "Flag", "Standard", "Pennant", "Crest"}
	default:
		return []string{"Mark", "Sign", "Seal", "Emblem", "Symbol"}
	}
}

// GenerateVariants creates multiple insignia variants for selection.
func (g *InsigniaGenerator) GenerateVariants(count int) []*Insignia {
	variants := make([]*Insignia, 0, count)
	for i := 0; i < count; i++ {
		variants = append(variants, g.Generate())
	}
	return variants
}

// InsigniaShapeName returns the display name for an insignia shape.
func InsigniaShapeName(shape InsigniaShape, genre engine.GenreID) string {
	names, ok := insigniaShapeNames[genre]
	if !ok {
		names = insigniaShapeNames[engine.GenreFantasy]
	}
	return names[shape]
}

var insigniaShapeNames = map[engine.GenreID]map[InsigniaShape]string{
	engine.GenreFantasy: {
		InsigniaShapeShield:  "Heraldic Shield",
		InsigniaShapeCircle:  "Seal",
		InsigniaShapeDiamond: "Gemstone",
		InsigniaShapeHexagon: "Arcane Hex",
		InsigniaShapeBanner:  "Battle Banner",
	},
	engine.GenreScifi: {
		InsigniaShapeShield:  "Tactical Badge",
		InsigniaShapeCircle:  "Corps Insignia",
		InsigniaShapeDiamond: "Command Mark",
		InsigniaShapeHexagon: "Fleet Emblem",
		InsigniaShapeBanner:  "Division Flag",
	},
	engine.GenreHorror: {
		InsigniaShapeShield:  "Survivor's Mark",
		InsigniaShapeCircle:  "Blood Seal",
		InsigniaShapeDiamond: "Warning Sign",
		InsigniaShapeHexagon: "Quarantine Badge",
		InsigniaShapeBanner:  "Faction Flag",
	},
	engine.GenreCyberpunk: {
		InsigniaShapeShield:  "Corp Badge",
		InsigniaShapeCircle:  "Access Token",
		InsigniaShapeDiamond: "Runner Mark",
		InsigniaShapeHexagon: "Grid Icon",
		InsigniaShapeBanner:  "Gang Colors",
	},
	engine.GenrePostapoc: {
		InsigniaShapeShield:  "Warlord's Mark",
		InsigniaShapeCircle:  "Tribe Sign",
		InsigniaShapeDiamond: "Trade Seal",
		InsigniaShapeHexagon: "Salvage Mark",
		InsigniaShapeBanner:  "War Flag",
	},
}

// InsigniaPatternName returns the display name for an insignia pattern.
func InsigniaPatternName(pattern InsigniaPattern) string {
	switch pattern {
	case InsigniaPatternSolid:
		return "Solid"
	case InsigniaPatternHorizontal:
		return "Horizontal Bands"
	case InsigniaPatternVertical:
		return "Vertical Stripes"
	case InsigniaPatternDiagonal:
		return "Diagonal Lines"
	case InsigniaPatternQuarters:
		return "Quartered"
	case InsigniaPatternChevron:
		return "Chevron"
	default:
		return "Unknown"
	}
}

// InsigniaSymbolName returns the display name for an insignia symbol.
func InsigniaSymbolName(symbol InsigniaSymbol, genre engine.GenreID) string {
	names, ok := insigniaSymbolNames[genre]
	if !ok {
		names = insigniaSymbolNames[engine.GenreFantasy]
	}
	return names[symbol]
}

var insigniaSymbolNames = map[engine.GenreID]map[InsigniaSymbol]string{
	engine.GenreFantasy: {
		InsigniaSymbolNone:     "None",
		InsigniaSymbolStar:     "Guiding Star",
		InsigniaSymbolCross:    "Holy Cross",
		InsigniaSymbolCircle:   "Moon Disc",
		InsigniaSymbolTriangle: "Mountain Peak",
		InsigniaSymbolArrow:    "Hunter's Arrow",
		InsigniaSymbolWave:     "River Flow",
		InsigniaSymbolGear:     "Craftsman's Wheel",
	},
	engine.GenreScifi: {
		InsigniaSymbolNone:     "None",
		InsigniaSymbolStar:     "Nova Burst",
		InsigniaSymbolCross:    "Navigation Mark",
		InsigniaSymbolCircle:   "Planet Icon",
		InsigniaSymbolTriangle: "Delta Formation",
		InsigniaSymbolArrow:    "Vector Sign",
		InsigniaSymbolWave:     "Warp Signature",
		InsigniaSymbolGear:     "Engineering Cog",
	},
	engine.GenreHorror: {
		InsigniaSymbolNone:     "None",
		InsigniaSymbolStar:     "Warning Star",
		InsigniaSymbolCross:    "Medical Cross",
		InsigniaSymbolCircle:   "Biohazard Ring",
		InsigniaSymbolTriangle: "Danger Sign",
		InsigniaSymbolArrow:    "Escape Route",
		InsigniaSymbolWave:     "Infection Wave",
		InsigniaSymbolGear:     "Salvage Gear",
	},
	engine.GenreCyberpunk: {
		InsigniaSymbolNone:     "None",
		InsigniaSymbolStar:     "Access Point",
		InsigniaSymbolCross:    "Targeting Reticle",
		InsigniaSymbolCircle:   "Data Node",
		InsigniaSymbolTriangle: "Upload Icon",
		InsigniaSymbolArrow:    "Direction Hack",
		InsigniaSymbolWave:     "Signal Pulse",
		InsigniaSymbolGear:     "Tech Mod",
	},
	engine.GenrePostapoc: {
		InsigniaSymbolNone:     "None",
		InsigniaSymbolStar:     "Survivor's Star",
		InsigniaSymbolCross:    "Trader's Cross",
		InsigniaSymbolCircle:   "Settlement Mark",
		InsigniaSymbolTriangle: "Hazard Warning",
		InsigniaSymbolArrow:    "Road Sign",
		InsigniaSymbolWave:     "Rad Wave",
		InsigniaSymbolGear:     "Mechanic's Badge",
	},
}

// ShapeName returns the genre-appropriate shape name for this insignia.
func (i *Insignia) ShapeName() string {
	return InsigniaShapeName(i.Shape, i.Genre)
}

// PatternName returns the pattern name for this insignia.
func (i *Insignia) PatternName() string {
	return InsigniaPatternName(i.Pattern)
}

// SymbolName returns the genre-appropriate symbol name for this insignia.
func (i *Insignia) SymbolName() string {
	return InsigniaSymbolName(i.Symbol, i.Genre)
}

// Description returns a full description of the insignia.
func (i *Insignia) Description() string {
	if i.Symbol == InsigniaSymbolNone {
		return i.PatternName() + " " + i.ShapeName()
	}
	return i.PatternName() + " " + i.ShapeName() + " with " + i.SymbolName()
}
