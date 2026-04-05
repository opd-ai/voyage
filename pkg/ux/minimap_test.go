//go:build headless

package ux

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/world"
)

// TestMinimapCreation tests Minimap initialization.
func TestMinimapCreation(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		mm := NewMinimap(genre, 150, 100)
		if mm == nil {
			t.Errorf("NewMinimap(%v) returned nil", genre)
			continue
		}
		if mm.Genre() != genre {
			t.Errorf("expected genre %v, got %v", genre, mm.Genre())
		}
		if mm.Width() != 150 {
			t.Errorf("expected width 150, got %d", mm.Width())
		}
		if mm.Height() != 100 {
			t.Errorf("expected height 100, got %d", mm.Height())
		}
	}
}

// TestMinimapSetGenre tests genre switching.
func TestMinimapSetGenre(t *testing.T) {
	mm := NewMinimap(engine.GenreFantasy, 150, 100)

	mm.SetGenre(engine.GenreScifi)
	if mm.Genre() != engine.GenreScifi {
		t.Errorf("expected genre Scifi, got %v", mm.Genre())
	}

	mm.SetGenre(engine.GenreCyberpunk)
	if mm.Genre() != engine.GenreCyberpunk {
		t.Errorf("expected genre Cyberpunk, got %v", mm.Genre())
	}
}

// TestMinimapCrisisMode tests crisis mode functionality.
func TestMinimapCrisisMode(t *testing.T) {
	mm := NewMinimap(engine.GenreFantasy, 150, 100)

	// Initially not in crisis mode
	if mm.IsCrisisMode() {
		t.Error("Minimap should not be in crisis mode initially")
	}
	if mm.Alpha() != 1.0 {
		t.Errorf("expected alpha 1.0, got %f", mm.Alpha())
	}

	// Enable crisis mode
	mm.SetCrisisMode(true)
	if !mm.IsCrisisMode() {
		t.Error("Minimap should be in crisis mode after SetCrisisMode(true)")
	}
	if mm.Alpha() != 0.3 {
		t.Errorf("expected alpha 0.3 in crisis mode, got %f", mm.Alpha())
	}

	// Disable crisis mode
	mm.SetCrisisMode(false)
	if mm.IsCrisisMode() {
		t.Error("Minimap should not be in crisis mode after SetCrisisMode(false)")
	}
	if mm.Alpha() != 1.0 {
		t.Errorf("expected alpha 1.0 after crisis mode disabled, got %f", mm.Alpha())
	}
}

// TestMinimapCalculateScale tests scale calculation.
func TestMinimapCalculateScale(t *testing.T) {
	mm := NewMinimap(engine.GenreFantasy, 150, 100)

	// Create a test world map
	wm := createTestWorldMap(50, 50)

	scale := mm.CalculateScale(wm)
	if scale <= 0 {
		t.Errorf("expected positive scale, got %f", scale)
	}

	// Scale should fit the world in the minimap
	expectedMaxScale := float64(mm.Width()-2) / float64(wm.Width)
	if scale > expectedMaxScale {
		t.Errorf("scale %f exceeds max expected %f", scale, expectedMaxScale)
	}
}

// TestMinimapCalculateScaleNilMap tests scale calculation with nil map.
func TestMinimapCalculateScaleNilMap(t *testing.T) {
	mm := NewMinimap(engine.GenreFantasy, 150, 100)

	scale := mm.CalculateScale(nil)
	if scale != 1.0 {
		t.Errorf("expected scale 1.0 for nil map, got %f", scale)
	}
}

// TestMinimapWorldToMinimap tests coordinate conversion.
func TestMinimapWorldToMinimap(t *testing.T) {
	mm := NewMinimap(engine.GenreFantasy, 150, 100)
	wm := createTestWorldMap(50, 50)

	// Test origin (0,0)
	px, py := mm.WorldToMinimap(wm, 0, 0)
	if px < 1 || py < 1 {
		t.Errorf("minimap coords should be >= 1, got (%d, %d)", px, py)
	}

	// Test center
	centerX, centerY := wm.Width/2, wm.Height/2
	px, py = mm.WorldToMinimap(wm, centerX, centerY)
	if px < 1 || py < 1 || px >= mm.Width()-1 || py >= mm.Height()-1 {
		t.Errorf("center coords out of bounds: (%d, %d)", px, py)
	}

	// Test that different world coords produce different minimap coords
	px1, py1 := mm.WorldToMinimap(wm, 0, 0)
	px2, py2 := mm.WorldToMinimap(wm, wm.Width-1, wm.Height-1)
	if px1 == px2 && py1 == py2 {
		t.Error("different world coords should produce different minimap coords")
	}
}

// TestMinimapCountExploredTiles tests explored tile counting.
func TestMinimapCountExploredTiles(t *testing.T) {
	mm := NewMinimap(engine.GenreFantasy, 150, 100)
	wm := createTestWorldMap(10, 10)

	// Initially no tiles explored
	count := mm.CountExploredTiles(wm)
	if count != 0 {
		t.Errorf("expected 0 explored tiles, got %d", count)
	}

	// Explore some tiles
	wm.GetTile(0, 0).Explored = true
	wm.GetTile(1, 1).Explored = true
	wm.GetTile(2, 2).Explored = true

	count = mm.CountExploredTiles(wm)
	if count != 3 {
		t.Errorf("expected 3 explored tiles, got %d", count)
	}
}

// TestMinimapCountExploredTilesNilMap tests counting with nil map.
func TestMinimapCountExploredTilesNilMap(t *testing.T) {
	mm := NewMinimap(engine.GenreFantasy, 150, 100)

	count := mm.CountExploredTiles(nil)
	if count != 0 {
		t.Errorf("expected 0 for nil map, got %d", count)
	}
}

// TestMinimapCountLandmarks tests landmark counting.
func TestMinimapCountLandmarks(t *testing.T) {
	mm := NewMinimap(engine.GenreFantasy, 150, 100)
	wm := createTestWorldMap(10, 10)

	// Add some landmarks and explore them
	wm.GetTile(2, 2).Landmark = &world.Landmark{Type: world.LandmarkTown, Name: "Test Town"}
	wm.GetTile(2, 2).Explored = true
	wm.GetTile(5, 5).Landmark = &world.Landmark{Type: world.LandmarkOutpost, Name: "Test Outpost"}
	wm.GetTile(5, 5).Explored = true
	// Add unexplored landmark (should not be counted)
	wm.GetTile(8, 8).Landmark = &world.Landmark{Type: world.LandmarkRuins, Name: "Hidden Ruins"}

	count := mm.CountLandmarks(wm)
	if count != 2 {
		t.Errorf("expected 2 landmarks, got %d", count)
	}
}

// TestMinimapCountLandmarksNilMap tests landmark counting with nil map.
func TestMinimapCountLandmarksNilMap(t *testing.T) {
	mm := NewMinimap(engine.GenreFantasy, 150, 100)

	count := mm.CountLandmarks(nil)
	if count != 0 {
		t.Errorf("expected 0 for nil map, got %d", count)
	}
}

// TestMinimapDimensions tests various minimap sizes.
func TestMinimapDimensions(t *testing.T) {
	tests := []struct {
		width  int
		height int
	}{
		{100, 80},
		{150, 100},
		{200, 150},
		{64, 64},
	}

	for _, tc := range tests {
		mm := NewMinimap(engine.GenreFantasy, tc.width, tc.height)
		if mm.Width() != tc.width {
			t.Errorf("expected width %d, got %d", tc.width, mm.Width())
		}
		if mm.Height() != tc.height {
			t.Errorf("expected height %d, got %d", tc.height, mm.Height())
		}
	}
}

// TestMinimapGenreTheming tests that different genres can be set.
func TestMinimapGenreTheming(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	mm := NewMinimap(engine.GenreFantasy, 150, 100)

	for _, genre := range genres {
		mm.SetGenre(genre)
		if mm.Genre() != genre {
			t.Errorf("failed to set genre to %v", genre)
		}
	}
}

// createTestWorldMap creates a simple world map for testing.
func createTestWorldMap(width, height int) *world.WorldMap {
	tiles := make([][]*world.Tile, height)
	for y := 0; y < height; y++ {
		tiles[y] = make([]*world.Tile, width)
		for x := 0; x < width; x++ {
			tiles[y][x] = &world.Tile{
				X:       x,
				Y:       y,
				Terrain: world.TerrainType(x % 7),
				Explored: false,
			}
		}
	}

	return &world.WorldMap{
		Width:       width,
		Height:      height,
		Tiles:       tiles,
		Origin:      world.Point{X: 0, Y: 0},
		Destination: world.Point{X: width - 1, Y: height - 1},
		Genre:       engine.GenreFantasy,
	}
}
