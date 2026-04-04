package vessel

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewInsigniaGenerator(t *testing.T) {
	gen := NewInsigniaGenerator(12345, engine.GenreFantasy)
	if gen == nil {
		t.Fatal("expected non-nil generator")
	}
	if gen.Genre() != engine.GenreFantasy {
		t.Errorf("expected fantasy genre, got %s", gen.Genre())
	}
}

func TestInsigniaGenerator_SetGenre(t *testing.T) {
	gen := NewInsigniaGenerator(12345, engine.GenreFantasy)
	gen.SetGenre(engine.GenreScifi)
	if gen.Genre() != engine.GenreScifi {
		t.Errorf("expected scifi genre, got %s", gen.Genre())
	}
}

func TestInsigniaGenerator_Generate(t *testing.T) {
	gen := NewInsigniaGenerator(12345, engine.GenreFantasy)
	insignia := gen.Generate()

	if insignia == nil {
		t.Fatal("expected non-nil insignia")
	}
	if insignia.Name == "" {
		t.Error("expected non-empty insignia name")
	}
	if insignia.PrimaryHue < 0 || insignia.PrimaryHue >= 360 {
		t.Errorf("primary hue out of range: %f", insignia.PrimaryHue)
	}
	if insignia.SecondaryHue < 0 || insignia.SecondaryHue >= 360 {
		t.Errorf("secondary hue out of range: %f", insignia.SecondaryHue)
	}
	if insignia.AccentHue < 0 || insignia.AccentHue >= 360 {
		t.Errorf("accent hue out of range: %f", insignia.AccentHue)
	}
	if insignia.Saturation < 0 || insignia.Saturation > 1 {
		t.Errorf("saturation out of range: %f", insignia.Saturation)
	}
	if insignia.Genre != engine.GenreFantasy {
		t.Errorf("expected fantasy genre, got %s", insignia.Genre)
	}
}

func TestInsigniaGenerator_Determinism(t *testing.T) {
	gen1 := NewInsigniaGenerator(42, engine.GenreHorror)
	gen2 := NewInsigniaGenerator(42, engine.GenreHorror)

	i1 := gen1.Generate()
	i2 := gen2.Generate()

	if i1.Shape != i2.Shape {
		t.Error("same seed should produce same shape")
	}
	if i1.Pattern != i2.Pattern {
		t.Error("same seed should produce same pattern")
	}
	if i1.Symbol != i2.Symbol {
		t.Error("same seed should produce same symbol")
	}
	if i1.Name != i2.Name {
		t.Errorf("same seed should produce same name: %s vs %s", i1.Name, i2.Name)
	}
}

func TestInsigniaGenerator_GenerateVariants(t *testing.T) {
	gen := NewInsigniaGenerator(12345, engine.GenreCyberpunk)
	variants := gen.GenerateVariants(5)

	if len(variants) != 5 {
		t.Errorf("expected 5 variants, got %d", len(variants))
	}

	for i, v := range variants {
		if v == nil {
			t.Errorf("variant %d is nil", i)
		}
	}
}

func TestInsigniaGenerator_AllGenres(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		gen := NewInsigniaGenerator(99999, genre)
		insignia := gen.Generate()

		if insignia.Genre != genre {
			t.Errorf("genre mismatch: expected %s, got %s", genre, insignia.Genre)
		}
		if insignia.Name == "" {
			t.Errorf("genre %s produced empty insignia name", genre)
		}
	}
}

func TestInsigniaShapeName(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		for _, shape := range AllInsigniaShapes() {
			name := InsigniaShapeName(shape, genre)
			if name == "" {
				t.Errorf("empty shape name for %v in genre %s", shape, genre)
			}
		}
	}
}

func TestInsigniaPatternName(t *testing.T) {
	for _, pattern := range AllInsigniaPatterns() {
		name := InsigniaPatternName(pattern)
		if name == "" {
			t.Errorf("empty pattern name for %v", pattern)
		}
	}
}

func TestInsigniaSymbolName(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		for _, symbol := range AllInsigniaSymbols() {
			name := InsigniaSymbolName(symbol, genre)
			if name == "" {
				t.Errorf("empty symbol name for %v in genre %s", symbol, genre)
			}
		}
	}
}

func TestInsignia_Description(t *testing.T) {
	gen := NewInsigniaGenerator(12345, engine.GenreFantasy)

	// Generate multiple to test both with and without symbols
	for i := 0; i < 10; i++ {
		insignia := gen.Generate()
		desc := insignia.Description()
		if desc == "" {
			t.Error("expected non-empty description")
		}
	}
}

func TestAllInsigniaShapes(t *testing.T) {
	shapes := AllInsigniaShapes()
	if len(shapes) != 5 {
		t.Errorf("expected 5 shapes, got %d", len(shapes))
	}
}

func TestAllInsigniaPatterns(t *testing.T) {
	patterns := AllInsigniaPatterns()
	if len(patterns) != 6 {
		t.Errorf("expected 6 patterns, got %d", len(patterns))
	}
}

func TestAllInsigniaSymbols(t *testing.T) {
	symbols := AllInsigniaSymbols()
	if len(symbols) != 8 {
		t.Errorf("expected 8 symbols, got %d", len(symbols))
	}
}
