package audio

import (
	"math"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// ADSR envelope constants for bass line.
const (
	bassAttack  = 0.5 // seconds to reach peak volume
	bassDecay   = 0.3 // seconds to decay to sustain level
	bassSustain = 0.6 // sustain level (0-1)
	bassRelease = 0.5 // seconds to fade to silence
	bassVolume  = 0.3 // mix volume for bass layer
)

// ADSR envelope constants for pad layer.
const (
	padAttack  = 1.0 // slow attack for ambient feel
	padDecay   = 0.5
	padSustain = 0.7
	padRelease = 1.0 // long release for atmosphere
	padVolume  = 0.2
)

// ADSR envelope constants for melody hints.
const (
	melodyAttack  = 0.1 // quick attack for plucky melody
	melodyDecay   = 0.2
	melodySustain = 0.3
	melodyRelease = 0.4
	melodyVolume  = 0.15
)

// MusicGenerator creates procedural ambient music.
type MusicGenerator struct {
	gen        *seed.Generator
	genre      engine.GenreID
	sampleRate float64
	bpm        float64
}

// NewMusicGenerator creates a new music generator.
func NewMusicGenerator(masterSeed int64, genre engine.GenreID) *MusicGenerator {
	return &MusicGenerator{
		gen:        seed.NewGenerator(masterSeed, "music"),
		genre:      genre,
		sampleRate: 44100,
		bpm:        80,
	}
}

// SetGenre changes the music style.
func (m *MusicGenerator) SetGenre(genre engine.GenreID) {
	m.genre = genre
}

// SetBPM changes the tempo.
func (m *MusicGenerator) SetBPM(bpm float64) {
	if bpm < 40 {
		bpm = 40
	}
	if bpm > 200 {
		bpm = 200
	}
	m.bpm = bpm
}

// GenerateLoop creates a looping ambient music segment.
// Returns samples in the range [-1, 1].
func (m *MusicGenerator) GenerateLoop(bars int) []float64 {
	params := m.getGenreMusicParams()
	beatsPerBar := 4
	beatDuration := 60.0 / m.bpm
	barDuration := beatDuration * float64(beatsPerBar)
	totalDuration := barDuration * float64(bars)
	totalSamples := int(m.sampleRate * totalDuration)
	result := make([]float64, totalSamples)

	m.generateBassLine(result, params, bars, barDuration)
	m.generatePadLayer(result, params, bars, barDuration)
	m.generateMelodyHints(result, params, bars, barDuration)
	m.normalizeAudio(result)

	return result
}

// generateBassLine adds a low drone/bass layer to the mix.
func (m *MusicGenerator) generateBassLine(result []float64, params *MusicParams, bars int, barDuration float64) {
	osc := NewOscillator(params.BassWave, params.RootNote, 0.25)
	env := NewEnvelope(bassAttack, bassDecay, bassSustain, bassRelease)
	env.NoteOn()

	for i := range result {
		t := float64(i) / m.sampleRate
		bar := int(t / barDuration)
		if bar >= bars {
			bar = bars - 1
		}
		bassNote := m.getBassNote(params, bar)
		osc.SetFrequency(bassNote)
		result[i] += osc.Sample() * env.Sample() * bassVolume
	}
}

// generatePadLayer adds ambient pad sounds.
func (m *MusicGenerator) generatePadLayer(result []float64, params *MusicParams, bars int, barDuration float64) {
	chord := m.generateChordNotes(params)
	oscs := make([]*Oscillator, len(chord))
	for i, note := range chord {
		oscs[i] = NewOscillator(params.PadWave, note, 0.15)
	}
	env := NewEnvelope(padAttack, padDecay, padSustain, padRelease)
	env.NoteOn()

	for i := range result {
		padSample := 0.0
		for _, osc := range oscs {
			padSample += osc.Sample()
		}
		result[i] += padSample * env.Sample() * padVolume
	}
}

// generateMelodyHints adds subtle melodic elements.
func (m *MusicGenerator) generateMelodyHints(result []float64, params *MusicParams, bars int, barDuration float64) {
	osc := NewOscillator(params.MelodyWave, params.RootNote*2, 0.2)
	env := NewEnvelope(melodyAttack, melodyDecay, melodySustain, melodyRelease)
	noteActive := false
	nextNoteTime := 0.0
	noteDuration := 60.0 / m.bpm

	for i := range result {
		t := float64(i) / m.sampleRate
		if t >= nextNoteTime && m.gen.Chance(params.MelodyDensity) {
			noteActive = true
			env.Reset()
			env.NoteOn()
			melodyNote := m.getMelodyNote(params)
			osc.SetFrequency(melodyNote)
			nextNoteTime = t + noteDuration*(0.5+m.gen.Float64())
		}
		if noteActive {
			result[i] += osc.Sample() * env.Sample() * melodyVolume
			if !env.IsActive() {
				noteActive = false
			}
		}
	}
}

// normalizeAudio prevents clipping by scaling the result.
func (m *MusicGenerator) normalizeAudio(result []float64) {
	maxAmp := 0.0
	for _, s := range result {
		if abs := math.Abs(s); abs > maxAmp {
			maxAmp = abs
		}
	}
	if maxAmp > 0.9 {
		scale := 0.85 / maxAmp
		for i := range result {
			result[i] *= scale
		}
	}
}

// getBassNote returns the bass note for a given bar.
func (m *MusicGenerator) getBassNote(params *MusicParams, bar int) float64 {
	progression := params.ChordProgression
	chordIndex := bar % len(progression)
	return params.RootNote * progression[chordIndex]
}

// generateChordNotes creates chord notes based on the root.
func (m *MusicGenerator) generateChordNotes(params *MusicParams) []float64 {
	root := params.RootNote * 2
	return []float64{root, root * params.ChordIntervals[0], root * params.ChordIntervals[1]}
}

// getMelodyNote returns a random scale note for melody.
func (m *MusicGenerator) getMelodyNote(params *MusicParams) float64 {
	scaleIndex := m.gen.Range(0, len(params.ScaleNotes)-1)
	octaveMod := 1.0
	if m.gen.Chance(0.3) {
		octaveMod = 2.0
	}
	return params.RootNote * params.ScaleNotes[scaleIndex] * octaveMod
}

// GenerateBytes converts samples to 16-bit PCM bytes.
func (m *MusicGenerator) GenerateBytes(bars int) []byte {
	samples := m.GenerateLoop(bars)
	bytes := make([]byte, len(samples)*2)

	for i, sample := range samples {
		sample = clampSample(sample)
		val := int16(sample * 32767)
		bytes[i*2] = byte(val)
		bytes[i*2+1] = byte(val >> 8)
	}

	return bytes
}

func clampSample(s float64) float64 {
	if s > 1 {
		return 1
	}
	if s < -1 {
		return -1
	}
	return s
}

// MusicParams holds genre-specific music generation parameters.
type MusicParams struct {
	BassWave         Waveform
	PadWave          Waveform
	MelodyWave       Waveform
	RootNote         float64
	ChordProgression []float64
	ChordIntervals   []float64
	ScaleNotes       []float64
	MelodyDensity    float64
}

func (m *MusicGenerator) getGenreMusicParams() *MusicParams {
	params := genreMusicParams[m.genre]
	if params == nil {
		params = genreMusicParams[engine.GenreFantasy]
	}
	return params
}

var genreMusicParams = map[engine.GenreID]*MusicParams{
	engine.GenreFantasy: {
		BassWave:         WaveTriangle,
		PadWave:          WaveSine,
		MelodyWave:       WaveTriangle,
		RootNote:         NoteToFrequency(48),                              // C3
		ChordProgression: []float64{1.0, 1.0, 1.5, 1.333},                  // i-i-v-iv
		ChordIntervals:   []float64{1.2, 1.5},                              // minor third, fifth
		ScaleNotes:       []float64{1.0, 1.125, 1.2, 1.333, 1.5, 1.6, 1.8}, // minor scale
		MelodyDensity:    0.15,
	},
	engine.GenreScifi: {
		BassWave:         WaveSine,
		PadWave:          WaveSawtooth,
		MelodyWave:       WaveSine,
		RootNote:         NoteToFrequency(52),                                      // E3
		ChordProgression: []float64{1.0, 1.0, 1.189, 1.0},                          // sustained root
		ChordIntervals:   []float64{1.189, 1.498},                                  // fourth, fifth
		ScaleNotes:       []float64{1.0, 1.122, 1.189, 1.335, 1.498, 1.587, 1.782}, // lydian
		MelodyDensity:    0.1,
	},
	engine.GenreHorror: {
		BassWave:         WaveSawtooth,
		PadWave:          WaveTriangle,
		MelodyWave:       WaveSawtooth,
		RootNote:         NoteToFrequency(43),                                     // G2
		ChordProgression: []float64{1.0, 1.059, 1.0, 1.122},                       // chromatic tension
		ChordIntervals:   []float64{1.189, 1.414},                                 // tritone tension
		ScaleNotes:       []float64{1.0, 1.059, 1.189, 1.26, 1.414, 1.498, 1.682}, // locrian
		MelodyDensity:    0.08,
	},
	engine.GenreCyberpunk: {
		BassWave:         WaveSawtooth,
		PadWave:          WaveSquare,
		MelodyWave:       WaveSawtooth,
		RootNote:         NoteToFrequency(55),                                      // G3
		ChordProgression: []float64{1.0, 1.0, 0.891, 0.944},                        // minor key
		ChordIntervals:   []float64{1.189, 1.498},                                  // power chord
		ScaleNotes:       []float64{1.0, 1.122, 1.189, 1.335, 1.498, 1.682, 1.888}, // phrygian
		MelodyDensity:    0.2,
	},
	engine.GenrePostapoc: {
		BassWave:         WaveTriangle,
		PadWave:          WaveTriangle,
		MelodyWave:       WaveTriangle,
		RootNote:         NoteToFrequency(45),                                // A2
		ChordProgression: []float64{1.0, 0.891, 0.944, 1.0},                  // sparse
		ChordIntervals:   []float64{1.2, 1.5},                                // minor
		ScaleNotes:       []float64{1.0, 1.122, 1.2, 1.335, 1.5, 1.6, 1.782}, // dorian
		MelodyDensity:    0.12,
	},
}

// AmbientLoop represents a generated ambient music loop.
type AmbientLoop struct {
	Samples    []float64
	SampleRate float64
	Duration   float64
	Genre      engine.GenreID
	BPM        float64
	Bars       int
}

// GenerateAmbientLoop creates an AmbientLoop structure.
func (m *MusicGenerator) GenerateAmbientLoop(bars int) *AmbientLoop {
	samples := m.GenerateLoop(bars)
	duration := float64(len(samples)) / m.sampleRate

	return &AmbientLoop{
		Samples:    samples,
		SampleRate: m.sampleRate,
		Duration:   duration,
		Genre:      m.genre,
		BPM:        m.bpm,
		Bars:       bars,
	}
}

// GetSampleAt returns the sample at a given time, wrapping for looping.
func (al *AmbientLoop) GetSampleAt(t float64) float64 {
	sampleIndex := int(t*al.SampleRate) % len(al.Samples)
	if sampleIndex < 0 {
		sampleIndex += len(al.Samples)
	}
	return al.Samples[sampleIndex]
}

// IsEmpty returns true if the loop has no samples.
func (al *AmbientLoop) IsEmpty() bool {
	return len(al.Samples) == 0
}
