package audio

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
	"github.com/opd-ai/voyage/pkg/procgen/world"
)

// AmbientLoopType represents different ambient soundscape categories.
type AmbientLoopType int

const (
	// AmbientWind is a general outdoor wind sound.
	AmbientWind AmbientLoopType = iota
	// AmbientSpaceHum is a low-frequency space station/ship hum.
	AmbientSpaceHum
	// AmbientGroaningMetal is creaking/groaning metallic sounds.
	AmbientGroaningMetal
	// AmbientCityNoise is urban ambient noise (distant traffic, hum).
	AmbientCityNoise
	// AmbientWater is flowing/dripping water sounds.
	AmbientWater
	// AmbientCreature is distant creature/insect sounds.
	AmbientCreature
)

// AmbientLoopTypeName returns a human-readable name for the ambient type.
func AmbientLoopTypeName(t AmbientLoopType) string {
	switch t {
	case AmbientWind:
		return "Wind"
	case AmbientSpaceHum:
		return "Space Hum"
	case AmbientGroaningMetal:
		return "Groaning Metal"
	case AmbientCityNoise:
		return "City Noise"
	case AmbientWater:
		return "Water"
	case AmbientCreature:
		return "Creature"
	default:
		return "Unknown"
	}
}

// AmbientGenerator creates procedural ambient soundscapes.
type AmbientGenerator struct {
	gen        *seed.Generator
	genre      engine.GenreID
	biome      world.BiomeType
	sampleRate float64
	loopType   AmbientLoopType
	volume     float64
}

// NewAmbientGenerator creates a new ambient sound generator.
func NewAmbientGenerator(masterSeed int64, genre engine.GenreID) *AmbientGenerator {
	ag := &AmbientGenerator{
		gen:        seed.NewGenerator(masterSeed, "ambient"),
		genre:      genre,
		biome:      world.BiomeTemperate,
		sampleRate: 44100,
		volume:     0.3,
	}
	ag.updateLoopTypeFromContext()
	return ag
}

// SetGenre changes the ambient sound style.
func (ag *AmbientGenerator) SetGenre(genre engine.GenreID) {
	ag.genre = genre
	ag.updateLoopTypeFromContext()
}

// SetBiome changes the biome for ambient sound selection.
func (ag *AmbientGenerator) SetBiome(biome world.BiomeType) {
	ag.biome = biome
	ag.updateLoopTypeFromContext()
}

// Biome returns the current biome.
func (ag *AmbientGenerator) Biome() world.BiomeType {
	return ag.biome
}

// LoopType returns the current ambient loop type.
func (ag *AmbientGenerator) LoopType() AmbientLoopType {
	return ag.loopType
}

// SetVolume adjusts the ambient volume (0.0 to 1.0).
func (ag *AmbientGenerator) SetVolume(v float64) {
	if v < 0 {
		v = 0
	}
	if v > 1 {
		v = 1
	}
	ag.volume = v
}

// Volume returns the current volume.
func (ag *AmbientGenerator) Volume() float64 {
	return ag.volume
}

// updateLoopTypeFromContext selects the appropriate ambient loop type based on genre and biome.
func (ag *AmbientGenerator) updateLoopTypeFromContext() {
	// Genre takes priority for overall soundscape
	switch ag.genre {
	case engine.GenreScifi:
		ag.loopType = AmbientSpaceHum
	case engine.GenreCyberpunk:
		ag.loopType = AmbientCityNoise
	case engine.GenreHorror:
		ag.loopType = AmbientGroaningMetal
	case engine.GenrePostapoc:
		// Post-apocalyptic uses biome-specific sounds
		ag.loopType = ag.postapocBiomeAmbient()
	default:
		// Fantasy uses biome-specific sounds
		ag.loopType = ag.fantasyBiomeAmbient()
	}
}

// fantasyBiomeAmbient returns the ambient type for fantasy genre based on biome.
func (ag *AmbientGenerator) fantasyBiomeAmbient() AmbientLoopType {
	switch ag.biome {
	case world.BiomeArid:
		return AmbientWind
	case world.BiomeForested:
		return AmbientCreature
	case world.BiomeMountainous:
		return AmbientWind
	case world.BiomeWetland:
		return AmbientWater
	case world.BiomeRuined:
		return AmbientWind
	default:
		return AmbientWind
	}
}

// postapocBiomeAmbient returns the ambient type for post-apocalyptic genre based on biome.
func (ag *AmbientGenerator) postapocBiomeAmbient() AmbientLoopType {
	switch ag.biome {
	case world.BiomeArid:
		return AmbientWind
	case world.BiomeForested:
		return AmbientCreature
	case world.BiomeMountainous:
		return AmbientWind
	case world.BiomeWetland:
		return AmbientWater
	case world.BiomeRuined:
		return AmbientGroaningMetal
	default:
		return AmbientWind
	}
}

// GenerateLoop creates an ambient sound loop of the specified duration in seconds.
func (ag *AmbientGenerator) GenerateLoop(durationSec float64) []float64 {
	if durationSec < 1 {
		durationSec = 1
	}
	if durationSec > 30 {
		durationSec = 30
	}

	totalSamples := int(ag.sampleRate * durationSec)
	result := make([]float64, totalSamples)

	switch ag.loopType {
	case AmbientWind:
		ag.generateWind(result)
	case AmbientSpaceHum:
		ag.generateSpaceHum(result)
	case AmbientGroaningMetal:
		ag.generateGroaningMetal(result)
	case AmbientCityNoise:
		ag.generateCityNoise(result)
	case AmbientWater:
		ag.generateWater(result)
	case AmbientCreature:
		ag.generateCreature(result)
	}

	// Apply volume and normalize
	ag.normalizeAndApplyVolume(result)

	return result
}

// generateWind creates a wind-like ambient sound using filtered noise.
func (ag *AmbientGenerator) generateWind(result []float64) {
	// Low-frequency oscillating filter cutoff for wind gusts
	gustSpeed := 0.15 + ag.gen.Float64()*0.1
	filterPos := 0.0

	for i := range result {
		t := float64(i) / ag.sampleRate

		// Generate noise base
		noise := ag.gen.Float64()*2 - 1

		// Oscillating filter creates wind gusts
		gustMod := 0.3 + 0.7*((1+sinFast(t*gustSpeed*6.28))/2)

		// Simple low-pass filter simulation
		filterPos += (noise*gustMod - filterPos) * 0.02
		result[i] = filterPos
	}
}

// generateSpaceHum creates a deep space hum ambient sound.
func (ag *AmbientGenerator) generateSpaceHum(result []float64) {
	// Multiple low-frequency sine waves with slight detuning
	baseFreq := 30.0 + ag.gen.Float64()*10
	detune1 := baseFreq * 1.01
	detune2 := baseFreq * 0.99
	detune3 := baseFreq * 2.01

	for i := range result {
		t := float64(i) / ag.sampleRate

		// Layered low-frequency hum
		hum := sinFast(t*baseFreq*6.28) * 0.5
		hum += sinFast(t*detune1*6.28) * 0.3
		hum += sinFast(t*detune2*6.28) * 0.3
		hum += sinFast(t*detune3*6.28) * 0.15 // Harmonic

		// Occasional subtle warble
		warble := 1.0 + 0.1*sinFast(t*0.5*6.28)
		result[i] = hum * warble
	}
}

// generateGroaningMetal creates creaking/groaning metal sounds.
func (ag *AmbientGenerator) generateGroaningMetal(result []float64) {
	// Low-frequency modulated noise with resonance
	resonanceFreq := 80.0 + ag.gen.Float64()*40
	filterPos := 0.0
	groanPhase := ag.gen.Float64() * 6.28

	for i := range result {
		t := float64(i) / ag.sampleRate

		// Base noise
		noise := ag.gen.Float64()*2 - 1

		// Resonant groaning
		groan := sinFast(t*resonanceFreq*6.28+groanPhase) * 0.4

		// Modulate intensity for creaking effect
		creak := sinFast(t * 0.3 * 6.28)
		creakMod := 0.3 + 0.7*((creak+1)/2)

		// Filter
		filterPos += (noise*0.2 + groan*creakMod - filterPos) * 0.05
		result[i] = filterPos
	}
}

// generateCityNoise creates urban ambient noise.
func (ag *AmbientGenerator) generateCityNoise(result []float64) {
	// Layered noise with different filter frequencies
	filterLow := 0.0
	filterMid := 0.0

	hum60Hz := 60.0 // Electrical hum
	hum120Hz := 120.0

	for i := range result {
		t := float64(i) / ag.sampleRate

		// Base noise
		noise := ag.gen.Float64()*2 - 1

		// Electrical hum (60Hz and harmonics)
		elecHum := sinFast(t*hum60Hz*6.28)*0.15 + sinFast(t*hum120Hz*6.28)*0.08

		// Low rumble (traffic)
		filterLow += (noise*0.3 - filterLow) * 0.01
		// Mid frequency hum (machinery)
		filterMid += (noise*0.2 - filterMid) * 0.03

		result[i] = filterLow*0.5 + filterMid*0.3 + elecHum
	}
}

// generateWater creates flowing/dripping water ambient sound.
func (ag *AmbientGenerator) generateWater(result []float64) {
	filterPos := 0.0
	bubbleTimer := 0.0
	bubbleActive := false
	bubbleFreq := 800.0
	bubblePhase := 0.0

	for i := range result {
		// Base water flow (filtered noise)
		noise := ag.gen.Float64()*2 - 1
		filterPos += (noise*0.3 - filterPos) * 0.08
		waterFlow := filterPos * 0.6

		// Random bubbles/drips
		bubbleTimer -= 1.0 / ag.sampleRate
		if bubbleTimer <= 0 && ag.gen.Chance(0.001) {
			bubbleActive = true
			bubbleTimer = 0.05 + ag.gen.Float64()*0.1
			bubbleFreq = 600 + ag.gen.Float64()*400
			bubblePhase = 0
		}

		bubbleSound := 0.0
		if bubbleActive {
			// Exponentially decaying bubble
			decay := 1.0 - bubblePhase/(bubbleTimer*ag.sampleRate)
			if decay < 0 {
				bubbleActive = false
				decay = 0
			}
			bubbleSound = sinFast(bubblePhase*bubbleFreq*6.28/ag.sampleRate) * decay * 0.3
			bubblePhase++
		}

		result[i] = waterFlow + bubbleSound
	}
}

// generateCreature creates distant creature/insect ambient sounds.
func (ag *AmbientGenerator) generateCreature(result []float64) {
	// Background chirping with occasional calls
	chirpTimer := 0.0
	chirpActive := false
	chirpFreq := 3000.0
	chirpPhase := 0.0
	chirpDuration := 0.0

	filterPos := 0.0

	for i := range result {
		// Soft background noise (wind in leaves)
		noise := ag.gen.Float64()*2 - 1
		filterPos += (noise*0.1 - filterPos) * 0.01
		background := filterPos * 0.3

		// Random creature sounds
		chirpTimer -= 1.0 / ag.sampleRate
		if chirpTimer <= 0 && ag.gen.Chance(0.0005) {
			chirpActive = true
			chirpDuration = 0.1 + ag.gen.Float64()*0.2
			chirpTimer = chirpDuration
			chirpFreq = 2000 + ag.gen.Float64()*2000
			chirpPhase = 0
		}

		chirpSound := 0.0
		if chirpActive {
			progress := chirpPhase / (chirpDuration * ag.sampleRate)
			if progress >= 1 {
				chirpActive = false
			} else {
				// Amplitude envelope
				env := sinFast(progress * 3.14) // Simple envelope
				// Frequency modulation for natural sound
				freqMod := 1.0 + 0.1*sinFast(progress*20*6.28)
				chirpSound = sinFast(chirpPhase*chirpFreq*freqMod*6.28/ag.sampleRate) * env * 0.2
			}
			chirpPhase++
		}

		result[i] = background + chirpSound
	}
}

// normalizeAndApplyVolume scales audio to prevent clipping and applies volume.
func (ag *AmbientGenerator) normalizeAndApplyVolume(result []float64) {
	maxAmp := 0.0
	for _, s := range result {
		if abs := absFast(s); abs > maxAmp {
			maxAmp = abs
		}
	}
	if maxAmp > 0.9 {
		scale := 0.85 / maxAmp
		for i := range result {
			result[i] *= scale
		}
	}
	// Apply volume
	for i := range result {
		result[i] *= ag.volume
	}
}

// GenerateBytes converts ambient samples to 16-bit PCM bytes.
func (ag *AmbientGenerator) GenerateBytes(durationSec float64) []byte {
	samples := ag.GenerateLoop(durationSec)
	bytes := make([]byte, len(samples)*2)

	for i, sample := range samples {
		sample = clampSample(sample)
		val := int16(sample * 32767)
		bytes[i*2] = byte(val)
		bytes[i*2+1] = byte(val >> 8)
	}

	return bytes
}

// sinFast is a fast sine approximation.
func sinFast(x float64) float64 {
	// Normalize to [-pi, pi]
	const twoPi = 6.283185307179586
	const pi = 3.141592653589793
	x = x - twoPi*float64(int(x/twoPi))
	if x < -pi {
		x += twoPi
	} else if x > pi {
		x -= twoPi
	}
	// Parabolic approximation
	const a = 4 / (pi * pi)
	const p = 0.225
	y := a * x * (pi - absFast(x))
	return p*(y*absFast(y)-y) + y
}

// absFast returns the absolute value of a float64.
func absFast(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// AllAmbientLoopTypes returns all available ambient loop types.
func AllAmbientLoopTypes() []AmbientLoopType {
	return []AmbientLoopType{
		AmbientWind,
		AmbientSpaceHum,
		AmbientGroaningMetal,
		AmbientCityNoise,
		AmbientWater,
		AmbientCreature,
	}
}
