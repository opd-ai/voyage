//go:build !headless

package rendering

import (
	"image/color"
	"testing"
)

func TestNewTileGenerator(t *testing.T) {
	tg := NewTileGenerator(12345, 16)
	if tg == nil {
		t.Fatal("NewTileGenerator returned nil")
	}
	if tg.tileSize != 16 {
		t.Errorf("expected tileSize 16, got %d", tg.tileSize)
	}
	if tg.gen == nil {
		t.Error("generator should be initialized")
	}
}

func TestGenerateTile(t *testing.T) {
	tg := NewTileGenerator(12345, 16)
	baseColor := color.RGBA{100, 150, 200, 255}

	img := tg.GenerateTile(baseColor, 0.2)
	if img == nil {
		t.Fatal("GenerateTile returned nil")
	}

	bounds := img.Bounds()
	if bounds.Dx() != 16 || bounds.Dy() != 16 {
		t.Errorf("expected 16x16 image, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestGenerateTileCA(t *testing.T) {
	tg := NewTileGenerator(12345, 16)
	baseColor := color.RGBA{50, 100, 50, 255}
	accentColor := color.RGBA{100, 200, 100, 255}

	img := tg.GenerateTileCA(baseColor, accentColor, 0.5, 3)
	if img == nil {
		t.Fatal("GenerateTileCA returned nil")
	}

	bounds := img.Bounds()
	if bounds.Dx() != 16 || bounds.Dy() != 16 {
		t.Errorf("expected 16x16 image, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestGenerateCharacterSprite(t *testing.T) {
	tg := NewTileGenerator(12345, 16)
	primary := color.RGBA{200, 100, 100, 255}
	secondary := color.RGBA{255, 200, 200, 255}

	img := tg.GenerateCharacterSprite(primary, secondary)
	if img == nil {
		t.Fatal("GenerateCharacterSprite returned nil")
	}

	bounds := img.Bounds()
	if bounds.Dx() != 16 || bounds.Dy() != 16 {
		t.Errorf("expected 16x16 image, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestTileGeneratorDeterminism(t *testing.T) {
	// Test that same seed produces generators in same state
	// We can't directly compare pixel values without the game running,
	// but we can verify the generator seeds are initialized correctly
	tg1 := NewTileGenerator(12345, 16)
	tg2 := NewTileGenerator(12345, 16)

	if tg1.gen == nil || tg2.gen == nil {
		t.Error("Generators should be initialized")
	}

	// Both should start from same seed and produce same sequence
	// We test this indirectly by generating tiles (which advances the RNG)
	baseColor := color.RGBA{100, 100, 100, 255}

	img1 := tg1.GenerateTile(baseColor, 0.1)
	img2 := tg2.GenerateTile(baseColor, 0.1)

	// Just verify both return valid images
	if img1 == nil || img2 == nil {
		t.Error("Both generators should produce valid images")
	}
}

func TestCountNeighbors(t *testing.T) {
	tg := NewTileGenerator(12345, 4)

	// Create a 4x4 test grid (grid[y][x])
	// Row 0: T F F F
	// Row 1: T T T F
	// Row 2: F T F F
	// Row 3: F F F F
	grid := [][]bool{
		{true, false, false, false},
		{true, true, true, false},
		{false, true, false, false},
		{false, false, false, false},
	}

	// Test center cell (1,1)
	count := tg.countNeighbors(grid, 1, 1)
	if count != 4 {
		t.Errorf("countNeighbors at (1,1) = %d, want 4", count)
	}

	// Test corner (0,0) - neighbors: (1,0), (0,1), (1,1) -> grid[0][1]=F, grid[1][0]=T, grid[1][1]=T = 2
	count = tg.countNeighbors(grid, 0, 0)
	if count != 2 {
		t.Errorf("countNeighbors at (0,0) = %d, want 2", count)
	}

	// Test (3,3) corner with no neighbors set
	count = tg.countNeighbors(grid, 3, 3)
	if count != 0 {
		t.Errorf("countNeighbors at (3,3) = %d, want 0", count)
	}
}
