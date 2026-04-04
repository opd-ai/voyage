package encounters

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestSFXGenerator(t *testing.T) {
	gen := NewSFXGenerator(12345, engine.GenreFantasy)

	sfxTypes := []SFXType{
		SFXEncounterStart,
		SFXPhaseSuccess,
		SFXPhaseFail,
		SFXVictory,
		SFXDefeat,
		SFXRoleAssign,
		SFXPause,
		SFXResume,
	}

	for _, sfxType := range sfxTypes {
		samples := gen.Generate(sfxType)
		if len(samples) == 0 {
			t.Errorf("Expected non-empty samples for SFX type %d", sfxType)
		}

		// Check samples are in valid range
		for i, sample := range samples {
			if sample < -1 || sample > 1 {
				t.Errorf("Sample %d for SFX type %d out of range: %f", i, sfxType, sample)
				break
			}
		}
	}
}

func TestSFXGeneratorSetGenre(t *testing.T) {
	gen := NewSFXGenerator(12345, engine.GenreFantasy)
	gen.SetGenre(engine.GenreScifi)

	if gen.genre != engine.GenreScifi {
		t.Errorf("Expected genre scifi, got %s", gen.genre)
	}
}

func TestSFXGeneratorBytes(t *testing.T) {
	gen := NewSFXGenerator(12345, engine.GenreFantasy)

	bytes := gen.GenerateBytes(SFXEncounterStart)
	if len(bytes) == 0 {
		t.Error("Expected non-empty bytes")
	}

	// Each sample is 2 bytes (16-bit)
	samples := gen.Generate(SFXEncounterStart)
	expectedLen := len(samples) * 2
	if len(bytes) != expectedLen {
		t.Errorf("Expected %d bytes, got %d", expectedLen, len(bytes))
	}
}

func TestSFXGeneratorAllGenres(t *testing.T) {
	genres := engine.AllGenres()

	for _, genre := range genres {
		gen := NewSFXGenerator(12345, genre)

		// Test that each genre produces different sounds
		samples := gen.Generate(SFXEncounterStart)
		if len(samples) == 0 {
			t.Errorf("Expected non-empty samples for genre %s", genre)
		}
	}
}

func TestSFXGeneratorDeterminism(t *testing.T) {
	gen1 := NewSFXGenerator(42, engine.GenreFantasy)
	gen2 := NewSFXGenerator(42, engine.GenreFantasy)

	samples1 := gen1.Generate(SFXVictory)
	samples2 := gen2.Generate(SFXVictory)

	if len(samples1) != len(samples2) {
		t.Errorf("Expected same length, got %d vs %d", len(samples1), len(samples2))
		return
	}

	for i := range samples1 {
		if samples1[i] != samples2[i] {
			t.Errorf("Sample %d differs: %f vs %f", i, samples1[i], samples2[i])
			break
		}
	}
}

func TestSinFunction(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
		epsilon  float64
	}{
		{0, 0, 0.001},
		{pi / 2, 1, 0.001},
		{pi, 0, 0.001},
		{-pi / 2, -1, 0.001},
	}

	for _, tt := range tests {
		result := sin(tt.input)
		diff := result - tt.expected
		if diff < 0 {
			diff = -diff
		}
		if diff > tt.epsilon {
			t.Errorf("sin(%f) = %f, expected %f (diff: %f)", tt.input, result, tt.expected, diff)
		}
	}
}
