package audio

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// SFXType identifies different sound effect types.
type SFXType int

const (
	// SFXTravel is the ambient travel sound.
	SFXTravel SFXType = iota
	// SFXEvent is the event notification sound.
	SFXEvent
	// SFXCrisis is the crisis/danger alert.
	SFXCrisis
	// SFXSuccess is the positive outcome sound.
	SFXSuccess
	// SFXDeath is the crew death sound.
	SFXDeath
	// SFXClick is the UI click sound.
	SFXClick
)

// AllSFXTypes returns all SFX types.
func AllSFXTypes() []SFXType {
	return []SFXType{SFXTravel, SFXEvent, SFXCrisis, SFXSuccess, SFXDeath, SFXClick}
}

// SFXTypeName returns a human-readable name.
func SFXTypeName(s SFXType) string {
	switch s {
	case SFXTravel:
		return "Travel"
	case SFXEvent:
		return "Event"
	case SFXCrisis:
		return "Crisis"
	case SFXSuccess:
		return "Success"
	case SFXDeath:
		return "Death"
	case SFXClick:
		return "Click"
	default:
		return "Unknown"
	}
}

// SFXGenerator creates procedural sound effects.
type SFXGenerator struct {
	gen        *seed.Generator
	genre      engine.GenreID
	sampleRate float64
}

// NewSFXGenerator creates a new SFX generator.
func NewSFXGenerator(masterSeed int64, genre engine.GenreID) *SFXGenerator {
	return &SFXGenerator{
		gen:        seed.NewGenerator(masterSeed, "sfx"),
		genre:      genre,
		sampleRate: 44100,
	}
}

// SetGenre changes the SFX timbre presets.
func (g *SFXGenerator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// Generate creates audio samples for the given SFX type.
// Returns samples in the range [-1, 1].
func (g *SFXGenerator) Generate(sfxType SFXType) []float64 {
	switch sfxType {
	case SFXTravel:
		return g.generateTravel()
	case SFXEvent:
		return g.generateEvent()
	case SFXCrisis:
		return g.generateCrisis()
	case SFXSuccess:
		return g.generateSuccess()
	case SFXDeath:
		return g.generateDeath()
	case SFXClick:
		return g.generateClick()
	default:
		return []float64{}
	}
}

// GenerateBytes creates 16-bit PCM audio bytes.
func (g *SFXGenerator) GenerateBytes(sfxType SFXType) []byte {
	samples := g.Generate(sfxType)
	bytes := make([]byte, len(samples)*2)

	for i, sample := range samples {
		// Clamp to [-1, 1]
		if sample > 1 {
			sample = 1
		} else if sample < -1 {
			sample = -1
		}
		// Convert to 16-bit signed integer
		val := int16(sample * 32767)
		bytes[i*2] = byte(val)
		bytes[i*2+1] = byte(val >> 8)
	}

	return bytes
}

// sfxContext holds common setup for SFX generation.
type sfxContext struct {
	samples []float64
	osc     *Oscillator
	env     *Envelope
	count   int
}

// prepareSFX creates common audio generation context.
func (g *SFXGenerator) prepareSFX(duration float64, wave Waveform, freq, amp float64, env *Envelope) *sfxContext {
	count := int(g.sampleRate * duration)
	return &sfxContext{
		samples: make([]float64, count),
		osc:     NewOscillator(wave, freq, amp),
		env:     env,
		count:   count,
	}
}

// generateTravel creates ambient travel sound.
func (g *SFXGenerator) generateTravel() []float64 {
	preset := g.getGenrePreset()
	env := SlowEnvelope()
	env.NoteOn()
	ctx := g.prepareSFX(0.5, preset.TravelWave, preset.TravelFreq, 0.3, env)

	for i := 0; i < ctx.count; i++ {
		ctx.samples[i] = ctx.osc.Sample() * ctx.env.Sample()
	}

	return ctx.samples
}

// generateEvent creates event notification sound.
func (g *SFXGenerator) generateEvent() []float64 {
	preset := g.getGenrePreset()
	env := QuickEnvelope()
	env.NoteOn()
	ctx := g.prepareSFX(0.3, preset.EventWave, preset.EventFreq, 0.5, env)

	// Two-tone alert
	freqMod := 1.0
	for i := 0; i < ctx.count; i++ {
		if i > ctx.count/2 {
			freqMod = 1.25 // Pitch up second half
		}
		ctx.osc.SetFrequency(preset.EventFreq * freqMod)
		ctx.samples[i] = ctx.osc.Sample() * ctx.env.Sample()
	}

	return ctx.samples
}

// generateCrisis creates crisis alert sound.
func (g *SFXGenerator) generateCrisis() []float64 {
	preset := g.getGenrePreset()
	env := QuickEnvelope()
	env.NoteOn()
	ctx := g.prepareSFX(0.5, preset.CrisisWave, preset.CrisisFreq, 0.7, env)

	// Rapid pulse effect
	pulseRate := 15.0 // Hz
	for i := 0; i < ctx.count; i++ {
		t := float64(i) / g.sampleRate
		pulse := 0.5 + 0.5*float64(int(t*pulseRate)%2)
		ctx.samples[i] = ctx.osc.Sample() * ctx.env.Sample() * pulse
	}

	return ctx.samples
}

// generateSuccess creates success sound.
func (g *SFXGenerator) generateSuccess() []float64 {
	duration := 0.4
	samples := int(g.sampleRate * duration)
	result := make([]float64, samples)

	preset := g.getGenrePreset()
	osc := NewOscillator(WaveSine, preset.SuccessFreq, 0.4)
	env := QuickEnvelope()
	env.NoteOn()

	// Rising arpeggio
	notes := []float64{1.0, 1.25, 1.5, 2.0}
	noteLength := samples / len(notes)
	// Guard against zero noteLength (H-013)
	if noteLength == 0 {
		noteLength = 1
	}

	for i := 0; i < samples; i++ {
		noteIndex := i / noteLength
		if noteIndex >= len(notes) {
			noteIndex = len(notes) - 1
		}
		osc.SetFrequency(preset.SuccessFreq * notes[noteIndex])
		result[i] = osc.Sample() * env.Sample()
	}

	return result
}

// generateDeath creates death sound.
func (g *SFXGenerator) generateDeath() []float64 {
	duration := 0.6
	samples := int(g.sampleRate * duration)
	result := make([]float64, samples)

	preset := g.getGenrePreset()
	osc := NewOscillator(preset.DeathWave, preset.DeathFreq, 0.5)
	env := NewEnvelope(0.01, 0.1, 0.4, 0.5)
	env.NoteOn()

	// Descending pitch
	for i := 0; i < samples; i++ {
		t := float64(i) / float64(samples)
		freqMod := 1.0 - 0.5*t // Pitch drops to 50%
		osc.SetFrequency(preset.DeathFreq * freqMod)
		result[i] = osc.Sample() * env.Sample()
	}

	return result
}

// generateClick creates UI click sound.
func (g *SFXGenerator) generateClick() []float64 {
	duration := 0.05
	samples := int(g.sampleRate * duration)
	result := make([]float64, samples)

	osc := NewOscillator(WaveSquare, 1000, 0.3)
	env := PunchyEnvelope()
	env.NoteOn()

	for i := 0; i < samples; i++ {
		result[i] = osc.Sample() * env.Sample()
	}

	return result
}

// GenrePreset contains timbre settings for a genre.
type GenrePreset struct {
	TravelWave  Waveform
	TravelFreq  float64
	EventWave   Waveform
	EventFreq   float64
	CrisisWave  Waveform
	CrisisFreq  float64
	SuccessFreq float64
	DeathWave   Waveform
	DeathFreq   float64
}

func (g *SFXGenerator) getGenrePreset() GenrePreset {
	presets := map[engine.GenreID]GenrePreset{
		engine.GenreFantasy: {
			TravelWave: WaveTriangle, TravelFreq: 220,
			EventWave: WaveSine, EventFreq: 440,
			CrisisWave: WaveSquare, CrisisFreq: 330,
			SuccessFreq: 523, DeathWave: WaveTriangle, DeathFreq: 220,
		},
		engine.GenreScifi: {
			TravelWave: WaveSine, TravelFreq: 110,
			EventWave: WaveSawtooth, EventFreq: 880,
			CrisisWave: WaveSquare, CrisisFreq: 440,
			SuccessFreq: 660, DeathWave: WaveSawtooth, DeathFreq: 165,
		},
		engine.GenreHorror: {
			TravelWave: WaveSawtooth, TravelFreq: 80,
			EventWave: WaveSquare, EventFreq: 330,
			CrisisWave: WaveNoise, CrisisFreq: 200,
			SuccessFreq: 392, DeathWave: WaveSawtooth, DeathFreq: 110,
		},
		engine.GenreCyberpunk: {
			TravelWave: WaveSawtooth, TravelFreq: 165,
			EventWave: WaveSquare, EventFreq: 660,
			CrisisWave: WaveSawtooth, CrisisFreq: 440,
			SuccessFreq: 880, DeathWave: WaveSquare, DeathFreq: 110,
		},
		engine.GenrePostapoc: {
			TravelWave: WaveTriangle, TravelFreq: 110,
			EventWave: WaveSquare, EventFreq: 440,
			CrisisWave: WaveNoise, CrisisFreq: 330,
			SuccessFreq: 440, DeathWave: WaveTriangle, DeathFreq: 165,
		},
	}

	if preset, ok := presets[g.genre]; ok {
		return preset
	}
	return presets[engine.GenreFantasy]
}
