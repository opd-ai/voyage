package encounters

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// SFXType identifies encounter-specific sound effects.
type SFXType int

const (
	// SFXEncounterStart plays when an encounter begins.
	SFXEncounterStart SFXType = iota
	// SFXPhaseSuccess plays on successful phase resolution.
	SFXPhaseSuccess
	// SFXPhaseFail plays on failed phase resolution.
	SFXPhaseFail
	// SFXVictory plays on encounter victory.
	SFXVictory
	// SFXDefeat plays on encounter defeat.
	SFXDefeat
	// SFXRoleAssign plays when crew is assigned to a role.
	SFXRoleAssign
	// SFXPause plays when encounter is paused.
	SFXPause
	// SFXResume plays when encounter resumes.
	SFXResume
)

// SFXGenerator creates procedural sound effects for encounters.
type SFXGenerator struct {
	gen        *seed.Generator
	genre      engine.GenreID
	sampleRate float64
}

// NewSFXGenerator creates a new encounter SFX generator.
func NewSFXGenerator(masterSeed int64, genre engine.GenreID) *SFXGenerator {
	return &SFXGenerator{
		gen:        seed.NewGenerator(masterSeed, "encounter_sfx"),
		genre:      genre,
		sampleRate: 44100,
	}
}

// SetGenre updates the SFX generator's genre.
func (g *SFXGenerator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// Generate creates audio samples for the given SFX type.
// Returns samples in the range [-1, 1].
func (g *SFXGenerator) Generate(sfxType SFXType) []float64 {
	switch sfxType {
	case SFXEncounterStart:
		return g.generateEncounterStart()
	case SFXPhaseSuccess:
		return g.generatePhaseSuccess()
	case SFXPhaseFail:
		return g.generatePhaseFail()
	case SFXVictory:
		return g.generateVictory()
	case SFXDefeat:
		return g.generateDefeat()
	case SFXRoleAssign:
		return g.generateRoleAssign()
	case SFXPause:
		return g.generatePause()
	case SFXResume:
		return g.generateResume()
	default:
		return []float64{}
	}
}

// GenerateBytes creates 16-bit PCM audio bytes.
func (g *SFXGenerator) GenerateBytes(sfxType SFXType) []byte {
	samples := g.Generate(sfxType)
	bytes := make([]byte, len(samples)*2)

	for i, sample := range samples {
		if sample > 1 {
			sample = 1
		} else if sample < -1 {
			sample = -1
		}
		val := int16(sample * 32767)
		bytes[i*2] = byte(val)
		bytes[i*2+1] = byte(val >> 8)
	}

	return bytes
}

func (g *SFXGenerator) generateEncounterStart() []float64 {
	duration := 0.4
	samples := int(g.sampleRate * duration)
	result := make([]float64, samples)

	preset := g.getGenrePreset()

	// Three-tone alert rising
	freqs := []float64{preset.BaseFreq, preset.BaseFreq * 1.25, preset.BaseFreq * 1.5}
	noteLen := samples / 3

	for i := 0; i < samples; i++ {
		noteIdx := i / noteLen
		if noteIdx >= 3 {
			noteIdx = 2
		}
		t := float64(i) / g.sampleRate
		freq := freqs[noteIdx]

		// Basic waveform based on genre
		var wave float64
		switch preset.WaveType {
		case 0: // Sine
			wave = sin(2 * pi * freq * t)
		case 1: // Square
			if sin(2*pi*freq*t) > 0 {
				wave = 0.5
			} else {
				wave = -0.5
			}
		case 2: // Sawtooth
			wave = 2*(freq*t-float64(int(freq*t))) - 1
		default:
			wave = sin(2 * pi * freq * t)
		}

		// Envelope
		env := 1.0 - float64(i%noteLen)/float64(noteLen)
		result[i] = wave * env * 0.4
	}

	return result
}

func (g *SFXGenerator) generatePhaseSuccess() []float64 {
	duration := 0.2
	samples := int(g.sampleRate * duration)
	result := make([]float64, samples)

	preset := g.getGenrePreset()
	freq := preset.SuccessFreq

	for i := 0; i < samples; i++ {
		t := float64(i) / g.sampleRate
		env := 1.0 - float64(i)/float64(samples)
		wave := sin(2 * pi * freq * t)
		result[i] = wave * env * 0.3
	}

	return result
}

func (g *SFXGenerator) generatePhaseFail() []float64 {
	duration := 0.25
	samples := int(g.sampleRate * duration)
	result := make([]float64, samples)

	preset := g.getGenrePreset()
	freq := preset.FailFreq

	for i := 0; i < samples; i++ {
		t := float64(i) / g.sampleRate
		env := 1.0 - float64(i)/float64(samples)
		// Descending tone
		freqMod := 1.0 - 0.3*float64(i)/float64(samples)
		wave := sin(2 * pi * freq * freqMod * t)
		result[i] = wave * env * 0.35
	}

	return result
}

func (g *SFXGenerator) generateVictory() []float64 {
	duration := 0.6
	samples := int(g.sampleRate * duration)
	result := make([]float64, samples)

	preset := g.getGenrePreset()

	// Triumphant arpeggio
	notes := []float64{preset.SuccessFreq, preset.SuccessFreq * 1.25, preset.SuccessFreq * 1.5, preset.SuccessFreq * 2}
	noteLen := samples / 4

	for i := 0; i < samples; i++ {
		noteIdx := i / noteLen
		if noteIdx >= 4 {
			noteIdx = 3
		}
		t := float64(i) / g.sampleRate
		freq := notes[noteIdx]
		env := 0.8 - 0.3*float64(i)/float64(samples)
		wave := sin(2 * pi * freq * t)
		result[i] = wave * env * 0.35
	}

	return result
}

func (g *SFXGenerator) generateDefeat() []float64 {
	duration := 0.5
	samples := int(g.sampleRate * duration)
	result := make([]float64, samples)

	preset := g.getGenrePreset()

	// Descending tones
	for i := 0; i < samples; i++ {
		t := float64(i) / g.sampleRate
		progress := float64(i) / float64(samples)
		freq := preset.FailFreq * (1.0 - 0.5*progress)
		env := 1.0 - progress
		wave := sin(2 * pi * freq * t)
		result[i] = wave * env * 0.4
	}

	return result
}

func (g *SFXGenerator) generateRoleAssign() []float64 {
	duration := 0.1
	samples := int(g.sampleRate * duration)
	result := make([]float64, samples)

	for i := 0; i < samples; i++ {
		t := float64(i) / g.sampleRate
		env := 1.0 - float64(i)/float64(samples)
		wave := sin(2 * pi * 800 * t)
		result[i] = wave * env * 0.2
	}

	return result
}

func (g *SFXGenerator) generatePause() []float64 {
	duration := 0.15
	samples := int(g.sampleRate * duration)
	result := make([]float64, samples)

	// Descending two-tone
	for i := 0; i < samples; i++ {
		t := float64(i) / g.sampleRate
		env := 1.0 - float64(i)/float64(samples)
		freq := 600.0
		if i > samples/2 {
			freq = 400.0
		}
		wave := sin(2 * pi * freq * t)
		result[i] = wave * env * 0.25
	}

	return result
}

func (g *SFXGenerator) generateResume() []float64 {
	duration := 0.15
	samples := int(g.sampleRate * duration)
	result := make([]float64, samples)

	// Rising two-tone
	for i := 0; i < samples; i++ {
		t := float64(i) / g.sampleRate
		env := 1.0 - float64(i)/float64(samples)
		freq := 400.0
		if i > samples/2 {
			freq = 600.0
		}
		wave := sin(2 * pi * freq * t)
		result[i] = wave * env * 0.25
	}

	return result
}

type genrePreset struct {
	BaseFreq    float64
	SuccessFreq float64
	FailFreq    float64
	WaveType    int // 0=sine, 1=square, 2=sawtooth
}

func (g *SFXGenerator) getGenrePreset() genrePreset {
	presets := map[engine.GenreID]genrePreset{
		engine.GenreFantasy: {
			BaseFreq: 330, SuccessFreq: 440, FailFreq: 220, WaveType: 0,
		},
		engine.GenreScifi: {
			BaseFreq: 440, SuccessFreq: 660, FailFreq: 220, WaveType: 2,
		},
		engine.GenreHorror: {
			BaseFreq: 220, SuccessFreq: 330, FailFreq: 165, WaveType: 1,
		},
		engine.GenreCyberpunk: {
			BaseFreq: 440, SuccessFreq: 880, FailFreq: 330, WaveType: 2,
		},
		engine.GenrePostapoc: {
			BaseFreq: 330, SuccessFreq: 440, FailFreq: 220, WaveType: 1,
		},
	}

	if preset, ok := presets[g.genre]; ok {
		return preset
	}
	return presets[engine.GenreFantasy]
}

// pi constant for waveform calculations.
const pi = 3.14159265358979323846

// sin calculates sine using Taylor series approximation.
func sin(x float64) float64 {
	// Normalize x to [-pi, pi]
	for x > pi {
		x -= 2 * pi
	}
	for x < -pi {
		x += 2 * pi
	}

	// Taylor series approximation
	result := x
	term := x
	for i := 1; i < 10; i++ {
		term *= -x * x / float64((2*i)*(2*i+1))
		result += term
	}
	return result
}
