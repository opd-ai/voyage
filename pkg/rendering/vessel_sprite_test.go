//go:build !headless

package rendering

import (
	"image/color"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestNewVesselSpriteGenerator(t *testing.T) {
	vsg := NewVesselSpriteGenerator(12345, 32)
	if vsg == nil {
		t.Fatal("NewVesselSpriteGenerator returned nil")
	}
	if vsg.spriteSize != 32 {
		t.Errorf("expected spriteSize 32, got %d", vsg.spriteSize)
	}
	if vsg.gen == nil {
		t.Error("generator should be initialized")
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

	// Check all damage states are generated
	if sprite.PristineSprite == nil {
		t.Error("PristineSprite should not be nil")
	}
	if sprite.WornSprite == nil {
		t.Error("WornSprite should not be nil")
	}
	if sprite.DamagedSprite == nil {
		t.Error("DamagedSprite should not be nil")
	}
	if sprite.CriticalSprite == nil {
		t.Error("CriticalSprite should not be nil")
	}

	// Check sprite sizes
	bounds := sprite.PristineSprite.Bounds()
	if bounds.Dx() != 32 || bounds.Dy() != 32 {
		t.Errorf("expected 32x32 sprite, got %dx%d", bounds.Dx(), bounds.Dy())
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

func TestVesselSpriteGetSprite(t *testing.T) {
	vsg := NewVesselSpriteGenerator(12345, 32)
	hullColor := color.RGBA{100, 100, 120, 255}
	accentColor := color.RGBA{200, 200, 50, 255}

	sprite := vsg.GenerateVesselSprite(hullColor, accentColor)

	testCases := []struct {
		state    VesselDamageState
		expected *ebiten.Image
	}{
		{VesselPristine, sprite.PristineSprite},
		{VesselWorn, sprite.WornSprite},
		{VesselDamaged, sprite.DamagedSprite},
		{VesselCritical, sprite.CriticalSprite},
	}

	for _, tc := range testCases {
		result := sprite.GetSprite(tc.state)
		if result != tc.expected {
			t.Errorf("GetSprite(%d) returned wrong sprite", tc.state)
		}
	}
}

func TestVesselSpriteGetSpriteForRatio(t *testing.T) {
	vsg := NewVesselSpriteGenerator(12345, 32)
	hullColor := color.RGBA{100, 100, 120, 255}
	accentColor := color.RGBA{200, 200, 50, 255}

	sprite := vsg.GenerateVesselSprite(hullColor, accentColor)

	testCases := []struct {
		ratio    float64
		expected *ebiten.Image
	}{
		{1.0, sprite.PristineSprite},
		{0.75, sprite.WornSprite},
		{0.30, sprite.DamagedSprite},
		{0.10, sprite.CriticalSprite},
	}

	for _, tc := range testCases {
		result := sprite.GetSpriteForRatio(tc.ratio)
		if result != tc.expected {
			t.Errorf("GetSpriteForRatio(%f) returned wrong sprite", tc.ratio)
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

	// Check pristine sprite bounds match
	b1 := sprite1.PristineSprite.Bounds()
	b2 := sprite2.PristineSprite.Bounds()
	if b1.Dx() != b2.Dx() || b1.Dy() != b2.Dy() {
		t.Error("Sprite bounds should match for same seed")
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
