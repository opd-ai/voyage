package world

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestGeneratorDeterminism(t *testing.T) {
	seed := int64(12345)
	genre := engine.GenreFantasy

	// Each generator uses the same seed and should produce identical output
	gen1 := NewGenerator(seed, genre)
	w1 := gen1.Generate(20, 15)
	
	gen2 := NewGenerator(seed, genre)
	w2 := gen2.Generate(20, 15)

	// Verify same dimensions
	if w1.Width != w2.Width || w1.Height != w2.Height {
		t.Error("Dimensions should match")
	}

	// Verify same origin and destination
	if w1.Origin != w2.Origin {
		t.Errorf("Origin mismatch: %v vs %v", w1.Origin, w2.Origin)
	}
	if w1.Destination != w2.Destination {
		t.Errorf("Destination mismatch: %v vs %v", w1.Destination, w2.Destination)
	}

	// Verify same terrain
	for y := 0; y < w1.Height; y++ {
		for x := 0; x < w1.Width; x++ {
			if w1.Tiles[y][x].Terrain != w2.Tiles[y][x].Terrain {
				t.Errorf("Terrain mismatch at (%d,%d)", x, y)
			}
			if w1.Tiles[y][x].Biome != w2.Tiles[y][x].Biome {
				t.Errorf("Biome mismatch at (%d,%d)", x, y)
			}
		}
	}
}

func TestGuaranteedPath(t *testing.T) {
	seeds := []int64{1, 42, 12345, 99999, 123456789}
	genres := engine.AllGenres()

	for _, seed := range seeds {
		for _, genre := range genres {
			gen := NewGenerator(seed, genre)
			w := gen.Generate(30, 20)

			result := w.FindPath(w.Origin, w.Destination)
			if !result.Found {
				t.Errorf("No path found for seed %d, genre %s", seed, genre)
			}
			if len(result.Path) < 2 {
				t.Errorf("Path too short for seed %d, genre %s", seed, genre)
			}
			if result.Path[0] != w.Origin {
				t.Errorf("Path doesn't start at origin for seed %d, genre %s", seed, genre)
			}
			if result.Path[len(result.Path)-1] != w.Destination {
				t.Errorf("Path doesn't end at destination for seed %d, genre %s", seed, genre)
			}
		}
	}
}

func TestGenreSwitching(t *testing.T) {
	seed := int64(42)
	gen := NewGenerator(seed, engine.GenreFantasy)

	// Generate with fantasy
	w1 := gen.Generate(10, 10)
	if w1.Genre != engine.GenreFantasy {
		t.Errorf("Expected fantasy genre, got %s", w1.Genre)
	}

	// Switch to scifi
	gen.SetGenre(engine.GenreScifi)
	w2 := gen.Generate(10, 10)
	if w2.Genre != engine.GenreScifi {
		t.Errorf("Expected scifi genre, got %s", w2.Genre)
	}
}

func TestTerrainInfo(t *testing.T) {
	for _, terrain := range AllTerrainTypes() {
		for _, genre := range engine.AllGenres() {
			info := DefaultTerrainInfo(terrain, genre)
			if info.Name == "" {
				t.Errorf("Empty name for terrain %d, genre %s", terrain, genre)
			}
			if info.MovementCost <= 0 {
				t.Errorf("Invalid movement cost for terrain %d", terrain)
			}
		}
	}
}

func TestBiomeInfo(t *testing.T) {
	for _, biome := range AllBiomeTypes() {
		for _, genre := range engine.AllGenres() {
			info := DefaultBiomeInfo(biome, genre)
			if info.Name == "" {
				t.Errorf("Empty name for biome %d, genre %s", biome, genre)
			}
			if len(info.TerrainWeights) == 0 {
				t.Errorf("No terrain weights for biome %d", biome)
			}
		}
	}
}

func TestMapAccess(t *testing.T) {
	gen := NewGenerator(42, engine.GenreFantasy)
	w := gen.Generate(10, 10)

	// Valid access
	tile := w.GetTile(5, 5)
	if tile == nil {
		t.Error("Expected valid tile at (5,5)")
	}

	// Out of bounds
	if w.GetTile(-1, 0) != nil {
		t.Error("Expected nil for negative x")
	}
	if w.GetTile(0, -1) != nil {
		t.Error("Expected nil for negative y")
	}
	if w.GetTile(w.Width, 0) != nil {
		t.Error("Expected nil for x >= width")
	}
	if w.GetTile(0, w.Height) != nil {
		t.Error("Expected nil for y >= height")
	}
}

func TestLandmarks(t *testing.T) {
	gen := NewGenerator(42, engine.GenreFantasy)
	w := gen.Generate(30, 20)

	// Check origin and destination have landmarks
	originTile := w.GetTile(w.Origin.X, w.Origin.Y)
	if originTile.Landmark == nil {
		t.Error("Origin should have a landmark")
	}
	if originTile.Landmark.Type != LandmarkOrigin {
		t.Error("Origin landmark should be LandmarkOrigin")
	}

	destTile := w.GetTile(w.Destination.X, w.Destination.Y)
	if destTile.Landmark == nil {
		t.Error("Destination should have a landmark")
	}
	if destTile.Landmark.Type != LandmarkDestination {
		t.Error("Destination landmark should be LandmarkDestination")
	}

	// Count landmarks
	landmarkCount := 0
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.Tiles[y][x].Landmark != nil {
				landmarkCount++
			}
		}
	}

	// Should have at least origin + destination + some others
	if landmarkCount < 5 {
		t.Errorf("Expected at least 5 landmarks, got %d", landmarkCount)
	}
}

func BenchmarkWorldGeneration(b *testing.B) {
	gen := NewGenerator(42, engine.GenreFantasy)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gen.Generate(50, 40)
	}
}

func BenchmarkPathfinding(b *testing.B) {
	gen := NewGenerator(42, engine.GenreFantasy)
	w := gen.Generate(50, 40)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.FindPath(w.Origin, w.Destination)
	}
}
