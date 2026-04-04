package audio

import "github.com/opd-ai/voyage/pkg/engine"

// GenreInstrument returns the primary instrument sound for a genre.
func GenreInstrument(genre engine.GenreID) Waveform {
	instruments := map[engine.GenreID]Waveform{
		engine.GenreFantasy:   WaveTriangle,
		engine.GenreScifi:     WaveSawtooth,
		engine.GenreHorror:    WaveSawtooth,
		engine.GenreCyberpunk: WaveSquare,
		engine.GenrePostapoc:  WaveTriangle,
	}
	if w, ok := instruments[genre]; ok {
		return w
	}
	return WaveSine
}

// GenreBaseFrequency returns the base frequency for a genre.
func GenreBaseFrequency(genre engine.GenreID) float64 {
	frequencies := map[engine.GenreID]float64{
		engine.GenreFantasy:   262, // C4 - Medieval, warm
		engine.GenreScifi:     330, // E4 - Bright, futuristic
		engine.GenreHorror:    196, // G3 - Low, ominous
		engine.GenreCyberpunk: 392, // G4 - High-tech edge
		engine.GenrePostapoc:  220, // A3 - Gritty, raw
	}
	if f, ok := frequencies[genre]; ok {
		return f
	}
	return 262 // Default C4
}

// GenreEnvelope returns an envelope suited for the genre.
func GenreEnvelope(genre engine.GenreID) *Envelope {
	switch genre {
	case engine.GenreFantasy:
		return NewEnvelope(0.1, 0.2, 0.5, 0.5) // Soft, flowing
	case engine.GenreScifi:
		return NewEnvelope(0.01, 0.1, 0.7, 0.3) // Sharp, synthetic
	case engine.GenreHorror:
		return NewEnvelope(0.3, 0.5, 0.3, 1.0) // Slow, creeping
	case engine.GenreCyberpunk:
		return NewEnvelope(0.005, 0.05, 0.6, 0.2) // Punchy, electronic
	case engine.GenrePostapoc:
		return NewEnvelope(0.05, 0.3, 0.4, 0.6) // Rough, raw
	default:
		return QuickEnvelope()
	}
}

// GenreSFXDescription returns a description of the genre's sound.
func GenreSFXDescription(genre engine.GenreID) string {
	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Warm, melodic tones with flowing envelopes",
		engine.GenreScifi:     "Synthetic, bright sounds with sharp attacks",
		engine.GenreHorror:    "Dark, ominous tones with slow builds",
		engine.GenreCyberpunk: "Electronic, punchy sounds with digital edge",
		engine.GenrePostapoc:  "Raw, gritty tones with organic decay",
	}
	if desc, ok := descriptions[genre]; ok {
		return desc
	}
	return "Standard audio profile"
}
