//go:build !headless

package rendering

import (
	"image/color"
	"testing"
)

func TestNewAnimatedTileGenerator(t *testing.T) {
	atg := NewAnimatedTileGenerator(12345, 16)
	if atg == nil {
		t.Fatal("NewAnimatedTileGenerator returned nil")
	}
	if atg.tileSize != 16 {
		t.Errorf("expected tileSize 16, got %d", atg.tileSize)
	}
	if atg.gen == nil {
		t.Error("generator should be initialized")
	}
}

func TestGenerateWaterTile(t *testing.T) {
	atg := NewAnimatedTileGenerator(12345, 16)
	baseColor := color.RGBA{30, 60, 120, 255}
	highlightColor := color.RGBA{80, 120, 180, 255}

	anim := atg.GenerateAnimatedTile(AnimationWater, baseColor, highlightColor)
	if anim == nil {
		t.Fatal("GenerateAnimatedTile returned nil for water")
	}
	if len(anim.Frames) != 4 {
		t.Errorf("expected 4 water frames, got %d", len(anim.Frames))
	}
	if !anim.Loop {
		t.Error("water animation should loop")
	}
	for i, frame := range anim.Frames {
		if frame == nil {
			t.Errorf("water frame %d is nil", i)
			continue
		}
		bounds := frame.Bounds()
		if bounds.Dx() != 16 || bounds.Dy() != 16 {
			t.Errorf("water frame %d: expected 16x16, got %dx%d", i, bounds.Dx(), bounds.Dy())
		}
	}
}

func TestGenerateGrassTile(t *testing.T) {
	atg := NewAnimatedTileGenerator(12345, 16)
	baseColor := color.RGBA{40, 80, 40, 255}
	tipColor := color.RGBA{80, 140, 80, 255}

	anim := atg.GenerateAnimatedTile(AnimationGrass, baseColor, tipColor)
	if anim == nil {
		t.Fatal("GenerateAnimatedTile returned nil for grass")
	}
	if len(anim.Frames) != 4 {
		t.Errorf("expected 4 grass frames, got %d", len(anim.Frames))
	}
	if !anim.Loop {
		t.Error("grass animation should loop")
	}
	for i, frame := range anim.Frames {
		if frame == nil {
			t.Errorf("grass frame %d is nil", i)
		}
	}
}

func TestGenerateFireTile(t *testing.T) {
	atg := NewAnimatedTileGenerator(12345, 16)
	baseColor := color.RGBA{80, 30, 10, 255}
	brightColor := color.RGBA{255, 200, 50, 255}

	anim := atg.GenerateAnimatedTile(AnimationFire, baseColor, brightColor)
	if anim == nil {
		t.Fatal("GenerateAnimatedTile returned nil for fire")
	}
	if len(anim.Frames) != 4 {
		t.Errorf("expected 4 fire frames, got %d", len(anim.Frames))
	}
	if !anim.Loop {
		t.Error("fire animation should loop")
	}
}

func TestAnimatedTileUpdate(t *testing.T) {
	atg := NewAnimatedTileGenerator(12345, 16)
	baseColor := color.RGBA{30, 60, 120, 255}
	highlightColor := color.RGBA{80, 120, 180, 255}

	anim := atg.GenerateAnimatedTile(AnimationWater, baseColor, highlightColor)

	if anim.currentFrame != 0 {
		t.Errorf("initial frame should be 0, got %d", anim.currentFrame)
	}

	anim.Update(0.1)
	if anim.currentFrame != 0 {
		t.Errorf("frame should still be 0 after 0.1s (frameTime=0.2), got %d", anim.currentFrame)
	}

	anim.Update(0.15)
	if anim.currentFrame != 1 {
		t.Errorf("frame should be 1 after 0.25s total, got %d", anim.currentFrame)
	}
}

func TestAnimatedTileLooping(t *testing.T) {
	atg := NewAnimatedTileGenerator(12345, 16)
	baseColor := color.RGBA{30, 60, 120, 255}
	highlightColor := color.RGBA{80, 120, 180, 255}

	anim := atg.GenerateAnimatedTile(AnimationWater, baseColor, highlightColor)

	for i := 0; i < 8; i++ {
		anim.Update(0.21)
	}

	if anim.currentFrame >= 4 {
		t.Errorf("frame should loop back, got %d", anim.currentFrame)
	}
}

func TestAnimatedTileReset(t *testing.T) {
	atg := NewAnimatedTileGenerator(12345, 16)
	baseColor := color.RGBA{30, 60, 120, 255}
	highlightColor := color.RGBA{80, 120, 180, 255}

	anim := atg.GenerateAnimatedTile(AnimationWater, baseColor, highlightColor)

	anim.Update(0.5)
	if anim.currentFrame == 0 {
		t.Error("frame should have advanced after update")
	}

	anim.Reset()
	if anim.currentFrame != 0 {
		t.Errorf("frame should be 0 after reset, got %d", anim.currentFrame)
	}
	if anim.elapsed != 0 {
		t.Errorf("elapsed should be 0 after reset, got %f", anim.elapsed)
	}
}

func TestCurrentFrame(t *testing.T) {
	atg := NewAnimatedTileGenerator(12345, 16)
	baseColor := color.RGBA{30, 60, 120, 255}
	highlightColor := color.RGBA{80, 120, 180, 255}

	anim := atg.GenerateAnimatedTile(AnimationWater, baseColor, highlightColor)

	frame := anim.CurrentFrame()
	if frame == nil {
		t.Fatal("CurrentFrame returned nil")
	}

	expectedFrame := anim.Frames[0]
	if frame != expectedFrame {
		t.Error("CurrentFrame should return first frame initially")
	}
}

func TestAnimatedTileDeterminism(t *testing.T) {
	atg1 := NewAnimatedTileGenerator(12345, 16)
	atg2 := NewAnimatedTileGenerator(12345, 16)

	baseColor := color.RGBA{30, 60, 120, 255}
	highlightColor := color.RGBA{80, 120, 180, 255}

	anim1 := atg1.GenerateAnimatedTile(AnimationWater, baseColor, highlightColor)
	anim2 := atg2.GenerateAnimatedTile(AnimationWater, baseColor, highlightColor)

	if len(anim1.Frames) != len(anim2.Frames) {
		t.Errorf("frame counts should match: %d vs %d", len(anim1.Frames), len(anim2.Frames))
	}
	if anim1.FrameTime != anim2.FrameTime {
		t.Errorf("frame times should match: %f vs %f", anim1.FrameTime, anim2.FrameTime)
	}
}

func TestSinApprox(t *testing.T) {
	testCases := []struct {
		input    float64
		expected float64
		epsilon  float64
	}{
		{0, 0, 0.1},
		{1.57, 1.0, 0.15},
		{3.14, 0, 0.1},
	}

	for _, tc := range testCases {
		result := sinApprox(tc.input)
		diff := result - tc.expected
		if diff < 0 {
			diff = -diff
		}
		if diff > tc.epsilon {
			t.Errorf("sinApprox(%f) = %f, expected ~%f (±%f)", tc.input, result, tc.expected, tc.epsilon)
		}
	}
}

func TestLerp(t *testing.T) {
	testCases := []struct {
		a, b, t  float64
		expected float64
	}{
		{0, 10, 0, 0},
		{0, 10, 1, 10},
		{0, 10, 0.5, 5},
		{-5, 5, 0.5, 0},
	}

	for _, tc := range testCases {
		result := lerp(tc.a, tc.b, tc.t)
		if result != tc.expected {
			t.Errorf("lerp(%f, %f, %f) = %f, expected %f", tc.a, tc.b, tc.t, result, tc.expected)
		}
	}
}

func TestClampFloat(t *testing.T) {
	testCases := []struct {
		v, min, max float64
		expected    float64
	}{
		{5, 0, 10, 5},
		{-5, 0, 10, 0},
		{15, 0, 10, 10},
		{0.5, 0, 1, 0.5},
	}

	for _, tc := range testCases {
		result := clampFloat(tc.v, tc.min, tc.max)
		if result != tc.expected {
			t.Errorf("clampFloat(%f, %f, %f) = %f, expected %f", tc.v, tc.min, tc.max, result, tc.expected)
		}
	}
}
