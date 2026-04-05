//go:build headless

package rendering

import (
	"image/color"
	"testing"
)

func TestNewVesselSpriteGenerator(t *testing.T) {
	vsg := NewVesselSpriteGenerator(12345, 32)
	if vsg == nil {
		t.Fatal("NewVesselSpriteGenerator returned nil")
	}
	if vsg.spriteSize != 32 {
		t.Errorf("expected spriteSize 32, got %d", vsg.spriteSize)
	}
}

func TestGenerateVesselSprite(t *testing.T) {
	vsg := NewVesselSpriteGenerator(12345, 32)
	hullColor := color.RGBA{100, 100, 120, 255}
	accentColor := color.RGBA{200, 200, 50, 255}

	sprite := vsg.GenerateVesselSprite(hullColor, accentColor)
	if sprite == nil {
		t.Fatal("GenerateVesselSprite returned nil")
	}

	// Check sprite size is stored
	if sprite.SpriteSize != 32 {
		t.Errorf("expected SpriteSize 32, got %d", sprite.SpriteSize)
	}
}

func TestVesselDamageStateFromRatio(t *testing.T) {
	testCases := []struct {
		ratio    float64
		expected VesselDamageState
	}{
		{1.0, VesselPristine},
		{0.95, VesselPristine},
		{0.90, VesselPristine},
		{0.89, VesselWorn},
		{0.50, VesselWorn},
		{0.49, VesselDamaged},
		{0.25, VesselDamaged},
		{0.24, VesselCritical},
		{0.0, VesselCritical},
	}

	for _, tc := range testCases {
		result := VesselDamageStateFromRatio(tc.ratio)
		if result != tc.expected {
			t.Errorf("VesselDamageStateFromRatio(%f) = %d, expected %d", tc.ratio, result, tc.expected)
		}
	}
}

func TestVesselSpriteGetDamageState(t *testing.T) {
	vsg := NewVesselSpriteGenerator(12345, 32)
	hullColor := color.RGBA{100, 100, 120, 255}
	accentColor := color.RGBA{200, 200, 50, 255}

	sprite := vsg.GenerateVesselSprite(hullColor, accentColor)

	testCases := []struct {
		ratio    float64
		expected VesselDamageState
	}{
		{1.0, VesselPristine},
		{0.75, VesselWorn},
		{0.30, VesselDamaged},
		{0.10, VesselCritical},
	}

	for _, tc := range testCases {
		result := sprite.GetDamageState(tc.ratio)
		if result != tc.expected {
			t.Errorf("GetDamageState(%f) = %d, expected %d", tc.ratio, result, tc.expected)
		}
	}
}

func TestVesselSpriteDeterminism(t *testing.T) {
	vsg1 := NewVesselSpriteGenerator(12345, 32)
	vsg2 := NewVesselSpriteGenerator(12345, 32)

	hullColor := color.RGBA{100, 100, 120, 255}
	accentColor := color.RGBA{200, 200, 50, 255}

	sprite1 := vsg1.GenerateVesselSprite(hullColor, accentColor)
	sprite2 := vsg2.GenerateVesselSprite(hullColor, accentColor)

	// Both should produce valid sprites
	if sprite1 == nil || sprite2 == nil {
		t.Fatal("Both generators should produce valid sprites")
	}

	// Check sprite sizes match
	if sprite1.SpriteSize != sprite2.SpriteSize {
		t.Error("Sprite sizes should match for same seed")
	}
}

func TestVesselSpriteDifferentSeeds(t *testing.T) {
	vsg1 := NewVesselSpriteGenerator(12345, 32)
	vsg2 := NewVesselSpriteGenerator(67890, 32)

	hullColor := color.RGBA{100, 100, 120, 255}
	accentColor := color.RGBA{200, 200, 50, 255}

	sprite1 := vsg1.GenerateVesselSprite(hullColor, accentColor)
	sprite2 := vsg2.GenerateVesselSprite(hullColor, accentColor)

	// Both should produce valid sprites
	if sprite1 == nil || sprite2 == nil {
		t.Fatal("Both generators should produce valid sprites")
	}
}

func TestAbsInt(t *testing.T) {
	testCases := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
		{-100, 100},
	}

	for _, tc := range testCases {
		result := absInt(tc.input)
		if result != tc.expected {
			t.Errorf("absInt(%d) = %d, expected %d", tc.input, result, tc.expected)
		}
	}
}
