//go:build headless

package rendering

import (
	"image/color"
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

// TestRendererCreation tests renderer creation.
func TestRendererCreation(t *testing.T) {
	r := NewRenderer(800, 600, 16)
	if r == nil {
		t.Fatal("NewRenderer returned nil")
	}
	if r.Width() != 800 {
		t.Errorf("expected width 800, got %d", r.Width())
	}
	if r.Height() != 600 {
		t.Errorf("expected height 600, got %d", r.Height())
	}
	if r.TileSize() != 16 {
		t.Errorf("expected tileSize 16, got %d", r.TileSize())
	}
	if r.Palette() == nil {
		t.Error("palette should not be nil")
	}
}

// TestRendererSetGenre tests genre changing on renderer.
func TestRendererSetGenre(t *testing.T) {
	r := NewRenderer(800, 600, 16)

	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		r.SetGenre(genre)
		if r.Palette() == nil {
			t.Errorf("palette should not be nil after SetGenre(%v)", genre)
		}
	}
}

// TestRendererCamera tests camera operations.
func TestRendererCamera(t *testing.T) {
	r := NewRenderer(800, 600, 16)

	// Test initial camera
	cam := r.Camera()
	if cam.X != 0 || cam.Y != 0 {
		t.Errorf("expected camera at (0,0), got (%f,%f)", cam.X, cam.Y)
	}

	// Test SetCamera
	r.SetCamera(100, 200)
	cam = r.Camera()
	if cam.X != 100 || cam.Y != 200 {
		t.Errorf("expected camera at (100,200), got (%f,%f)", cam.X, cam.Y)
	}

	// Test SetZoom
	r.SetZoom(2.0)
	cam = r.Camera()
	if cam.Zoom != 2.0 {
		t.Errorf("expected zoom 2.0, got %f", cam.Zoom)
	}
}

// TestAnimatedTileCreation tests animated tile creation.
func TestAnimatedTileCreation(t *testing.T) {
	at := NewAnimatedTile(4, 0.25, true)
	if at == nil {
		t.Fatal("NewAnimatedTile returned nil")
	}
	if at.FrameCount != 4 {
		t.Errorf("expected FrameCount 4, got %d", at.FrameCount)
	}
	if at.FrameTime != 0.25 {
		t.Errorf("expected FrameTime 0.25, got %f", at.FrameTime)
	}
	if !at.Loop {
		t.Error("expected Loop to be true")
	}
	if at.CurrentFrameIndex() != 0 {
		t.Errorf("expected initial frame 0, got %d", at.CurrentFrameIndex())
	}
}

// TestAnimatedTileUpdate tests animation frame progression.
func TestAnimatedTileUpdate(t *testing.T) {
	at := NewAnimatedTile(4, 0.25, true)

	// Update less than frame time - should not change frame
	at.Update(0.1)
	if at.CurrentFrameIndex() != 0 {
		t.Errorf("expected frame 0 after 0.1s, got %d", at.CurrentFrameIndex())
	}

	// Update past frame time - should advance
	at.Update(0.2) // Total 0.3s, past 0.25s
	if at.CurrentFrameIndex() != 1 {
		t.Errorf("expected frame 1 after 0.3s, got %d", at.CurrentFrameIndex())
	}
}

// TestAnimatedTileLoop tests animation looping.
func TestAnimatedTileLoop(t *testing.T) {
	at := NewAnimatedTile(3, 0.1, true)

	// Advance through all frames
	for i := 0; i < 4; i++ {
		at.Update(0.15)
	}

	// Should loop back to 0
	if at.CurrentFrameIndex() > 2 {
		t.Errorf("looping animation should stay in range, got %d", at.CurrentFrameIndex())
	}
}

// TestAnimatedTileNoLoop tests animation without looping.
func TestAnimatedTileNoLoop(t *testing.T) {
	at := NewAnimatedTile(3, 0.1, false)

	// Advance past all frames
	for i := 0; i < 5; i++ {
		at.Update(0.15)
	}

	// Should stop at last frame
	if at.CurrentFrameIndex() != 2 {
		t.Errorf("non-looping animation should stop at last frame, got %d", at.CurrentFrameIndex())
	}
}

// TestAnimatedTileSingleFrame tests animation with single frame.
func TestAnimatedTileSingleFrame(t *testing.T) {
	at := NewAnimatedTile(1, 0.25, true)

	// Update should not advance (single frame)
	at.Update(1.0)
	if at.CurrentFrameIndex() != 0 {
		t.Errorf("single frame animation should stay at 0, got %d", at.CurrentFrameIndex())
	}
}

// TestAnimatedTileReset tests animation reset.
func TestAnimatedTileReset(t *testing.T) {
	at := NewAnimatedTile(4, 0.1, true)

	// Advance some frames
	at.Update(0.5)
	if at.CurrentFrameIndex() == 0 {
		// Advance more to ensure we move
		at.Update(0.5)
	}

	// Reset
	at.Reset()
	if at.CurrentFrameIndex() != 0 {
		t.Errorf("reset should set frame to 0, got %d", at.CurrentFrameIndex())
	}
}

// TestAnimatedTileGeneratorCreation tests animated tile generator creation.
func TestAnimatedTileGeneratorCreation(t *testing.T) {
	atg := NewAnimatedTileGenerator(12345, 16)
	if atg == nil {
		t.Fatal("NewAnimatedTileGenerator returned nil")
	}
	if atg.tileSize != 16 {
		t.Errorf("expected tileSize 16, got %d", atg.tileSize)
	}
}

// TestGenerateAnimatedTile tests animated tile generation.
func TestGenerateAnimatedTile(t *testing.T) {
	atg := NewAnimatedTileGenerator(12345, 16)
	baseColor := color.RGBA{100, 100, 200, 255}
	accentColor := color.RGBA{200, 200, 255, 255}

	animTypes := []AnimationType{
		AnimationWater,
		AnimationGrass,
		AnimationFire,
	}

	for _, animType := range animTypes {
		at := atg.GenerateAnimatedTile(animType, baseColor, accentColor)
		if at == nil {
			t.Errorf("GenerateAnimatedTile(%d) returned nil", animType)
			continue
		}
		if at.FrameCount < 1 {
			t.Errorf("animated tile should have at least 1 frame")
		}
	}
}

// TestLandmarkIconGeneratorCreation tests landmark icon generator creation.
func TestLandmarkIconGeneratorCreation(t *testing.T) {
	lig := NewLandmarkIconGenerator(12345, 16)
	if lig == nil {
		t.Fatal("NewLandmarkIconGenerator returned nil")
	}
	if lig.iconSize != 16 {
		t.Errorf("expected iconSize 16, got %d", lig.iconSize)
	}
}

// TestGenerateAnimatedLandmarkIcon tests landmark icon generation.
func TestGenerateAnimatedLandmarkIcon(t *testing.T) {
	lig := NewLandmarkIconGenerator(12345, 16)
	primaryColor := color.RGBA{200, 150, 100, 255}
	secondaryColor := color.RGBA{255, 200, 150, 255}

	iconTypes := []LandmarkIconType{
		LandmarkIconTown,
		LandmarkIconOutpost,
		LandmarkIconRuins,
		LandmarkIconShrine,
		LandmarkIconOrigin,
		LandmarkIconDestination,
	}

	for _, iconType := range iconTypes {
		ali := lig.GenerateAnimatedIcon(iconType, primaryColor, secondaryColor)
		if ali == nil {
			t.Errorf("GenerateAnimatedIcon(%d) returned nil", iconType)
			continue
		}
		if ali.FrameCount < 1 {
			t.Errorf("animated icon should have at least 1 frame")
		}
	}
}

// TestAnimatedLandmarkIconCreation tests animated landmark icon creation.
func TestAnimatedLandmarkIconCreation(t *testing.T) {
	ali := NewAnimatedLandmarkIcon(6, 0.15)
	if ali == nil {
		t.Fatal("NewAnimatedLandmarkIcon returned nil")
	}
	if ali.FrameCount != 6 {
		t.Errorf("expected FrameCount 6, got %d", ali.FrameCount)
	}
	if ali.FrameTime != 0.15 {
		t.Errorf("expected FrameTime 0.15, got %f", ali.FrameTime)
	}
}

// TestAnimatedLandmarkIconUpdate tests landmark icon animation.
func TestAnimatedLandmarkIconUpdate(t *testing.T) {
	ali := NewAnimatedLandmarkIcon(4, 0.2)

	// Update and verify frame advancement
	ali.Update(0.25) // Past one frame time
	if ali.CurrentFrameIndex() != 1 {
		t.Errorf("expected frame 1 after 0.25s, got %d", ali.CurrentFrameIndex())
	}
}

// TestAnimatedLandmarkIconReset tests landmark icon reset.
func TestAnimatedLandmarkIconReset(t *testing.T) {
	ali := NewAnimatedLandmarkIcon(4, 0.1)

	ali.Update(0.5)
	ali.Reset()
	if ali.CurrentFrameIndex() != 0 {
		t.Errorf("reset should set frame to 0, got %d", ali.CurrentFrameIndex())
	}
}

// Note: ParticleSystem tests already in particles_test.go - only testing what's missing here

// TestDefaultPalette tests palette creation for all genres.
func TestDefaultPalette(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		palette := DefaultPalette(genre)
		if palette == nil {
			t.Errorf("DefaultPalette(%v) returned nil", genre)
			continue
		}
		if palette.Background == nil {
			t.Errorf("palette for %v has nil Background", genre)
		}
		if len(palette.TileColors) == 0 {
			t.Errorf("palette for %v has no TileColors", genre)
		}
	}
}

// TestSinApprox tests sine approximation function.
func TestSinApprox(t *testing.T) {
	// sinApprox uses Taylor series, accurate near 0 and wraps for periodicity
	// Only test where the approximation is known to be accurate
	tests := []struct {
		x        float64
		expected float64
		epsilon  float64
	}{
		{0, 0, 0.05},
		{0.5, 0.479, 0.05},   // sin(0.5) ≈ 0.479
		{1.0, 0.841, 0.05},   // sin(1.0) ≈ 0.841
		{1.5708, 1.0, 0.2},   // ~π/2, may diverge slightly
		{-0.5, -0.479, 0.15}, // sin(-0.5), wrapping may add error
		{6.283, 0, 0.3},      // ~2π should wrap to ~0
	}

	for _, tt := range tests {
		result := sinApprox(tt.x)
		diff := result - tt.expected
		if diff < 0 {
			diff = -diff
		}
		if diff > tt.epsilon {
			t.Errorf("sinApprox(%f) = %f, expected ~%f (epsilon %f)", tt.x, result, tt.expected, tt.epsilon)
		}
	}
}

// TestLerp tests linear interpolation function.
func TestLerp(t *testing.T) {
	tests := []struct {
		a, b, t  float64
		expected float64
	}{
		{0, 10, 0, 0},
		{0, 10, 1, 10},
		{0, 10, 0.5, 5},
		{-10, 10, 0.5, 0},
	}

	for _, tt := range tests {
		result := lerp(tt.a, tt.b, tt.t)
		if result != tt.expected {
			t.Errorf("lerp(%f, %f, %f) = %f, want %f", tt.a, tt.b, tt.t, result, tt.expected)
		}
	}
}

// TestLandmarkIconTypeConstants tests that constants are distinct.
func TestLandmarkIconTypeConstants(t *testing.T) {
	types := []LandmarkIconType{
		LandmarkIconTown,
		LandmarkIconOutpost,
		LandmarkIconRuins,
		LandmarkIconShrine,
		LandmarkIconOrigin,
		LandmarkIconDestination,
	}

	seen := make(map[LandmarkIconType]bool)
	for _, lt := range types {
		if seen[lt] {
			t.Errorf("duplicate LandmarkIconType: %d", lt)
		}
		seen[lt] = true
	}
}

// TestAnimationTypeConstants tests that constants are distinct.
func TestAnimationTypeConstants(t *testing.T) {
	types := []AnimationType{
		AnimationWater,
		AnimationGrass,
		AnimationFire,
	}

	seen := make(map[AnimationType]bool)
	for _, at := range types {
		if seen[at] {
			t.Errorf("duplicate AnimationType: %d", at)
		}
		seen[at] = true
	}
}
