//go:build !headless

package rendering

import (
	"image/color"
	"testing"
)

func TestNewLandmarkIconGenerator(t *testing.T) {
	lig := NewLandmarkIconGenerator(12345, 24)
	if lig == nil {
		t.Fatal("NewLandmarkIconGenerator returned nil")
	}
	if lig.iconSize != 24 {
		t.Errorf("expected iconSize 24, got %d", lig.iconSize)
	}
	if lig.gen == nil {
		t.Error("generator should be initialized")
	}
}

func TestGenerateRuinsIcon(t *testing.T) {
	lig := NewLandmarkIconGenerator(12345, 24)
	baseColor := color.RGBA{80, 70, 60, 255}
	smokeColor := color.RGBA{150, 150, 150, 200}

	icon := lig.GenerateAnimatedIcon(LandmarkIconRuins, baseColor, smokeColor)
	if icon == nil {
		t.Fatal("GenerateAnimatedIcon returned nil for ruins")
	}
	if len(icon.Frames) != 6 {
		t.Errorf("expected 6 ruins frames, got %d", len(icon.Frames))
	}
	for i, frame := range icon.Frames {
		if frame == nil {
			t.Errorf("ruins frame %d is nil", i)
			continue
		}
		bounds := frame.Bounds()
		if bounds.Dx() != 24 || bounds.Dy() != 24 {
			t.Errorf("ruins frame %d: expected 24x24, got %dx%d", i, bounds.Dx(), bounds.Dy())
		}
	}
}

func TestGenerateOutpostIcon(t *testing.T) {
	lig := NewLandmarkIconGenerator(12345, 24)
	buildingColor := color.RGBA{100, 80, 60, 255}
	lightColor := color.RGBA{255, 200, 50, 255}

	icon := lig.GenerateAnimatedIcon(LandmarkIconOutpost, buildingColor, lightColor)
	if icon == nil {
		t.Fatal("GenerateAnimatedIcon returned nil for outpost")
	}
	if len(icon.Frames) != 4 {
		t.Errorf("expected 4 outpost frames, got %d", len(icon.Frames))
	}
}

func TestGenerateTownIcon(t *testing.T) {
	lig := NewLandmarkIconGenerator(12345, 24)
	buildingColor := color.RGBA{90, 80, 70, 255}
	lightColor := color.RGBA{255, 220, 100, 255}

	icon := lig.GenerateAnimatedIcon(LandmarkIconTown, buildingColor, lightColor)
	if icon == nil {
		t.Fatal("GenerateAnimatedIcon returned nil for town")
	}
	if len(icon.Frames) != 4 {
		t.Errorf("expected 4 town frames, got %d", len(icon.Frames))
	}
}

func TestGenerateShrineIcon(t *testing.T) {
	lig := NewLandmarkIconGenerator(12345, 24)
	stoneColor := color.RGBA{120, 110, 100, 255}
	glowColor := color.RGBA{200, 200, 100, 255}

	icon := lig.GenerateAnimatedIcon(LandmarkIconShrine, stoneColor, glowColor)
	if icon == nil {
		t.Fatal("GenerateAnimatedIcon returned nil for shrine")
	}
	if len(icon.Frames) != 4 {
		t.Errorf("expected 4 shrine frames, got %d", len(icon.Frames))
	}
}

func TestGenerateOriginIcon(t *testing.T) {
	lig := NewLandmarkIconGenerator(12345, 24)
	markerColor := color.RGBA{50, 150, 50, 255}
	glowColor := color.RGBA{100, 200, 100, 255}

	icon := lig.GenerateAnimatedIcon(LandmarkIconOrigin, markerColor, glowColor)
	if icon == nil {
		t.Fatal("GenerateAnimatedIcon returned nil for origin")
	}
	if len(icon.Frames) != 4 {
		t.Errorf("expected 4 origin frames, got %d", len(icon.Frames))
	}
}

func TestGenerateDestinationIcon(t *testing.T) {
	lig := NewLandmarkIconGenerator(12345, 24)
	markerColor := color.RGBA{200, 50, 50, 255}
	beamColor := color.RGBA{255, 200, 50, 255}

	icon := lig.GenerateAnimatedIcon(LandmarkIconDestination, markerColor, beamColor)
	if icon == nil {
		t.Fatal("GenerateAnimatedIcon returned nil for destination")
	}
	if len(icon.Frames) != 6 {
		t.Errorf("expected 6 destination frames, got %d", len(icon.Frames))
	}
}

func TestLandmarkIconUpdate(t *testing.T) {
	lig := NewLandmarkIconGenerator(12345, 24)
	baseColor := color.RGBA{80, 70, 60, 255}
	smokeColor := color.RGBA{150, 150, 150, 200}

	icon := lig.GenerateAnimatedIcon(LandmarkIconRuins, baseColor, smokeColor)

	if icon.currentFrame != 0 {
		t.Errorf("initial frame should be 0, got %d", icon.currentFrame)
	}

	// Frame time for ruins is 0.15
	icon.Update(0.1)
	if icon.currentFrame != 0 {
		t.Errorf("frame should still be 0 after 0.1s, got %d", icon.currentFrame)
	}

	icon.Update(0.1)
	if icon.currentFrame != 1 {
		t.Errorf("frame should be 1 after 0.2s total, got %d", icon.currentFrame)
	}
}

func TestLandmarkIconLooping(t *testing.T) {
	lig := NewLandmarkIconGenerator(12345, 24)
	buildingColor := color.RGBA{90, 80, 70, 255}
	lightColor := color.RGBA{255, 220, 100, 255}

	icon := lig.GenerateAnimatedIcon(LandmarkIconTown, buildingColor, lightColor)

	// Advance through multiple cycles
	for i := 0; i < 10; i++ {
		icon.Update(0.26) // More than frame time (0.25)
	}

	// Should have looped back
	if icon.currentFrame >= 4 {
		t.Errorf("animation should loop, frame %d >= 4", icon.currentFrame)
	}
}

func TestLandmarkIconReset(t *testing.T) {
	lig := NewLandmarkIconGenerator(12345, 24)
	baseColor := color.RGBA{80, 70, 60, 255}
	smokeColor := color.RGBA{150, 150, 150, 200}

	icon := lig.GenerateAnimatedIcon(LandmarkIconRuins, baseColor, smokeColor)

	icon.Update(0.5)
	if icon.currentFrame == 0 {
		t.Error("frame should have advanced after update")
	}

	icon.Reset()
	if icon.currentFrame != 0 {
		t.Errorf("frame should be 0 after reset, got %d", icon.currentFrame)
	}
	if icon.elapsed != 0 {
		t.Errorf("elapsed should be 0 after reset, got %f", icon.elapsed)
	}
}

func TestLandmarkIconCurrentFrame(t *testing.T) {
	lig := NewLandmarkIconGenerator(12345, 24)
	baseColor := color.RGBA{80, 70, 60, 255}
	smokeColor := color.RGBA{150, 150, 150, 200}

	icon := lig.GenerateAnimatedIcon(LandmarkIconRuins, baseColor, smokeColor)

	frame := icon.CurrentFrame()
	if frame == nil {
		t.Fatal("CurrentFrame returned nil")
	}

	expectedFrame := icon.Frames[0]
	if frame != expectedFrame {
		t.Error("CurrentFrame should return first frame initially")
	}
}

func TestLandmarkIconDeterminism(t *testing.T) {
	lig1 := NewLandmarkIconGenerator(12345, 24)
	lig2 := NewLandmarkIconGenerator(12345, 24)

	baseColor := color.RGBA{80, 70, 60, 255}
	smokeColor := color.RGBA{150, 150, 150, 200}

	icon1 := lig1.GenerateAnimatedIcon(LandmarkIconRuins, baseColor, smokeColor)
	icon2 := lig2.GenerateAnimatedIcon(LandmarkIconRuins, baseColor, smokeColor)

	if len(icon1.Frames) != len(icon2.Frames) {
		t.Errorf("frame counts should match: %d vs %d", len(icon1.Frames), len(icon2.Frames))
	}
	if icon1.FrameTime != icon2.FrameTime {
		t.Errorf("frame times should match: %f vs %f", icon1.FrameTime, icon2.FrameTime)
	}
}

func TestAllLandmarkTypes(t *testing.T) {
	lig := NewLandmarkIconGenerator(12345, 24)

	types := []LandmarkIconType{
		LandmarkIconTown,
		LandmarkIconOutpost,
		LandmarkIconRuins,
		LandmarkIconShrine,
		LandmarkIconOrigin,
		LandmarkIconDestination,
	}

	primaryColor := color.RGBA{100, 100, 100, 255}
	secondaryColor := color.RGBA{200, 200, 50, 255}

	for _, iconType := range types {
		icon := lig.GenerateAnimatedIcon(iconType, primaryColor, secondaryColor)
		if icon == nil {
			t.Errorf("GenerateAnimatedIcon returned nil for type %d", iconType)
			continue
		}
		if len(icon.Frames) == 0 {
			t.Errorf("type %d should have at least one frame", iconType)
		}
		if icon.FrameTime <= 0 {
			t.Errorf("type %d should have positive frame time", iconType)
		}
	}
}
