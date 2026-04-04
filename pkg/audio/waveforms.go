package audio

import (
	"math"
)

// Waveform represents different oscillator types.
type Waveform int

const (
	// WaveSine produces a smooth, pure tone.
	WaveSine Waveform = iota
	// WaveSquare produces a harsh, hollow sound.
	WaveSquare
	// WaveSawtooth produces a bright, buzzy sound.
	WaveSawtooth
	// WaveTriangle produces a soft, mellow sound.
	WaveTriangle
	// WaveNoise produces random noise.
	WaveNoise
)

// AllWaveforms returns all available waveforms.
func AllWaveforms() []Waveform {
	return []Waveform{WaveSine, WaveSquare, WaveSawtooth, WaveTriangle, WaveNoise}
}

// WaveformName returns a human-readable name.
func WaveformName(w Waveform) string {
	switch w {
	case WaveSine:
		return "Sine"
	case WaveSquare:
		return "Square"
	case WaveSawtooth:
		return "Sawtooth"
	case WaveTriangle:
		return "Triangle"
	case WaveNoise:
		return "Noise"
	default:
		return "Unknown"
	}
}

// Oscillator generates audio samples for a waveform.
type Oscillator struct {
	waveform   Waveform
	frequency  float64
	amplitude  float64
	phase      float64
	sampleRate float64
	noiseState uint32 // For noise generation
}

// NewOscillator creates an oscillator with the given parameters.
func NewOscillator(waveform Waveform, frequency, amplitude float64) *Oscillator {
	return &Oscillator{
		waveform:   waveform,
		frequency:  frequency,
		amplitude:  amplitude,
		phase:      0,
		sampleRate: 44100,
		noiseState: 12345,
	}
}

// SetFrequency changes the oscillator frequency.
func (o *Oscillator) SetFrequency(freq float64) {
	o.frequency = freq
}

// SetAmplitude changes the oscillator amplitude.
func (o *Oscillator) SetAmplitude(amp float64) {
	o.amplitude = amp
}

// SetWaveform changes the oscillator waveform.
func (o *Oscillator) SetWaveform(w Waveform) {
	o.waveform = w
}

// Sample generates the next audio sample.
func (o *Oscillator) Sample() float64 {
	var sample float64

	switch o.waveform {
	case WaveSine:
		sample = o.sineWave()
	case WaveSquare:
		sample = o.squareWave()
	case WaveSawtooth:
		sample = o.sawtoothWave()
	case WaveTriangle:
		sample = o.triangleWave()
	case WaveNoise:
		sample = o.noise()
	default:
		sample = 0
	}

	// Advance phase
	o.phase += 2 * math.Pi * o.frequency / o.sampleRate
	if o.phase >= 2*math.Pi {
		o.phase -= 2 * math.Pi
	}

	return sample * o.amplitude
}

// Reset resets the oscillator phase.
func (o *Oscillator) Reset() {
	o.phase = 0
}

func (o *Oscillator) sineWave() float64 {
	return math.Sin(o.phase)
}

func (o *Oscillator) squareWave() float64 {
	if math.Sin(o.phase) >= 0 {
		return 1.0
	}
	return -1.0
}

func (o *Oscillator) sawtoothWave() float64 {
	return 2.0*(o.phase/(2*math.Pi)) - 1.0
}

func (o *Oscillator) triangleWave() float64 {
	// Triangle wave from -1 to 1
	t := o.phase / (2 * math.Pi)
	if t < 0.25 {
		return 4 * t
	} else if t < 0.75 {
		return 2 - 4*t
	}
	return 4*t - 4
}

func (o *Oscillator) noise() float64 {
	// Simple LFSR-based noise
	o.noiseState = o.noiseState*1103515245 + 12345
	return (float64(o.noiseState&0x7FFFFFFF)/float64(0x7FFFFFFF))*2 - 1
}

// GenerateSamples generates multiple samples.
func (o *Oscillator) GenerateSamples(count int) []float64 {
	samples := make([]float64, count)
	for i := 0; i < count; i++ {
		samples[i] = o.Sample()
	}
	return samples
}

// NoteToFrequency converts a MIDI note number to frequency.
// Note 69 = A4 = 440 Hz
func NoteToFrequency(note int) float64 {
	return 440.0 * math.Pow(2, float64(note-69)/12.0)
}

// FrequencyToNote converts frequency to the nearest MIDI note.
func FrequencyToNote(freq float64) int {
	return int(math.Round(12*math.Log2(freq/440.0) + 69))
}
