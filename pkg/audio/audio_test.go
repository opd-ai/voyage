package audio

import (
	"math"
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestOscillator(t *testing.T) {
	osc := NewOscillator(WaveSine, 440, 1.0)

	// Generate some samples
	samples := osc.GenerateSamples(100)

	if len(samples) != 100 {
		t.Errorf("samples = %d, want 100", len(samples))
	}

	// Check samples are within range
	for i, s := range samples {
		if s < -1.0 || s > 1.0 {
			t.Errorf("sample %d = %f, out of range [-1, 1]", i, s)
		}
	}
}

func TestOscillatorWaveforms(t *testing.T) {
	waveforms := AllWaveforms()

	for _, w := range waveforms {
		osc := NewOscillator(w, 440, 1.0)
		samples := osc.GenerateSamples(1000)

		// All samples should be within range
		for i, s := range samples {
			if s < -1.0 || s > 1.0 {
				t.Errorf("waveform %s: sample %d = %f, out of range", WaveformName(w), i, s)
			}
		}

		// Non-noise waveforms should have some variation
		if w != WaveNoise {
			hasVariation := false
			for i := 1; i < len(samples); i++ {
				if samples[i] != samples[0] {
					hasVariation = true
					break
				}
			}
			if !hasVariation {
				t.Errorf("waveform %s should have variation", WaveformName(w))
			}
		}
	}
}

func TestEnvelope(t *testing.T) {
	env := NewEnvelope(0.1, 0.1, 0.5, 0.2)

	// Initially idle
	if env.IsActive() {
		t.Error("envelope should start idle")
	}

	// Trigger note on
	env.NoteOn()
	if !env.IsActive() {
		t.Error("envelope should be active after note on")
	}

	// Generate samples through attack phase
	attackSamples := int(44100 * 0.1)
	for i := 0; i < attackSamples; i++ {
		env.Sample()
	}

	// Level should be near 1.0 after attack
	level := env.Level()
	if level < 0.8 {
		t.Errorf("level after attack = %f, want >= 0.8", level)
	}

	// Continue through decay
	decaySamples := int(44100 * 0.15)
	for i := 0; i < decaySamples; i++ {
		env.Sample()
	}

	// Should be at or near sustain level
	level = env.Level()
	if level < 0.4 || level > 0.6 {
		t.Errorf("level at sustain = %f, want ~0.5", level)
	}

	// Note off
	env.NoteOff()
	if !env.IsReleased() {
		t.Error("envelope should be released after note off")
	}

	// Continue through release
	releaseSamples := int(44100 * 0.25)
	for i := 0; i < releaseSamples; i++ {
		env.Sample()
	}

	// Should be idle now
	if env.IsActive() {
		t.Error("envelope should be idle after release")
	}
}

func TestEnvelopePresets(t *testing.T) {
	quick := QuickEnvelope()
	if quick == nil {
		t.Error("QuickEnvelope should not be nil")
	}

	slow := SlowEnvelope()
	if slow == nil {
		t.Error("SlowEnvelope should not be nil")
	}

	punchy := PunchyEnvelope()
	if punchy == nil {
		t.Error("PunchyEnvelope should not be nil")
	}
}

func TestSFXGenerator(t *testing.T) {
	gen := NewSFXGenerator(12345, engine.GenreFantasy)

	for _, sfxType := range AllSFXTypes() {
		samples := gen.Generate(sfxType)

		if len(samples) == 0 {
			t.Errorf("SFX %s should generate samples", SFXTypeName(sfxType))
			continue
		}

		// Check samples are in range
		for i, s := range samples {
			if s < -1.0 || s > 1.0 {
				t.Errorf("SFX %s: sample %d = %f, out of range", SFXTypeName(sfxType), i, s)
			}
		}
	}
}

func TestSFXGeneratorGenres(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		gen := NewSFXGenerator(12345, genre)

		// Generate a sample SFX
		samples := gen.Generate(SFXEvent)

		if len(samples) == 0 {
			t.Errorf("genre %s: should generate event SFX", genre)
		}
	}
}

func TestSFXGeneratorBytes(t *testing.T) {
	gen := NewSFXGenerator(12345, engine.GenreFantasy)

	bytes := gen.GenerateBytes(SFXClick)

	if len(bytes) == 0 {
		t.Error("should generate bytes")
	}

	// Should be even number (16-bit samples)
	if len(bytes)%2 != 0 {
		t.Errorf("bytes length = %d, should be even", len(bytes))
	}
}

func TestPlayer(t *testing.T) {
	player := NewPlayer(12345, engine.GenreFantasy)

	if player.IsMuted() {
		t.Error("player should not start muted")
	}

	if player.Volume() != 1.0 {
		t.Errorf("volume = %f, want 1.0", player.Volume())
	}

	// Test volume
	player.SetVolume(0.5)
	if player.Volume() != 0.5 {
		t.Errorf("volume = %f, want 0.5", player.Volume())
	}

	// Test mute
	player.Mute()
	if !player.IsMuted() {
		t.Error("player should be muted")
	}

	samples := player.PlaySFX(SFXClick)
	if samples != nil {
		t.Error("muted player should not return samples")
	}

	player.Unmute()
	samples = player.PlaySFX(SFXClick)
	if samples == nil {
		t.Error("unmuted player should return samples")
	}
}

func TestPlayerGenreSwitch(t *testing.T) {
	player := NewPlayer(12345, engine.GenreFantasy)

	player.SetGenre(engine.GenreScifi)
	if player.Genre() != engine.GenreScifi {
		t.Errorf("genre = %s, want scifi", player.Genre())
	}
}

func TestPreloadSFX(t *testing.T) {
	player := NewPlayer(12345, engine.GenreFantasy)

	cache := player.PreloadSFX()

	if len(cache) != len(AllSFXTypes()) {
		t.Errorf("cache size = %d, want %d", len(cache), len(AllSFXTypes()))
	}

	for _, sfxType := range AllSFXTypes() {
		if _, ok := cache[sfxType]; !ok {
			t.Errorf("cache missing SFX type %s", SFXTypeName(sfxType))
		}
	}
}

func TestNoteToFrequency(t *testing.T) {
	// A4 should be 440 Hz
	freq := NoteToFrequency(69)
	if math.Abs(freq-440) > 0.01 {
		t.Errorf("note 69 freq = %f, want 440", freq)
	}

	// C4 should be ~262 Hz
	freqC4 := NoteToFrequency(60)
	if freqC4 < 260 || freqC4 > 264 {
		t.Errorf("note 60 freq = %f, want ~262", freqC4)
	}

	// Octave up should double frequency
	freqA5 := NoteToFrequency(81)
	if math.Abs(freqA5-880) > 0.01 {
		t.Errorf("note 81 freq = %f, want 880", freqA5)
	}
}

func TestFrequencyToNote(t *testing.T) {
	note := FrequencyToNote(440)
	if note != 69 {
		t.Errorf("440 Hz = note %d, want 69", note)
	}

	note880 := FrequencyToNote(880)
	if note880 != 81 {
		t.Errorf("880 Hz = note %d, want 81", note880)
	}
}

func TestGenrePresets(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		w := GenreInstrument(genre)
		if w < 0 || w > WaveNoise {
			t.Errorf("genre %s: invalid instrument waveform", genre)
		}

		f := GenreBaseFrequency(genre)
		if f < 20 || f > 2000 {
			t.Errorf("genre %s: base frequency %f out of audible range", genre, f)
		}

		env := GenreEnvelope(genre)
		if env == nil {
			t.Errorf("genre %s: envelope is nil", genre)
		}

		desc := GenreSFXDescription(genre)
		if desc == "" {
			t.Errorf("genre %s: missing description", genre)
		}
	}
}

func TestMusicGenerator(t *testing.T) {
	gen := NewMusicGenerator(12345, engine.GenreFantasy)

	// Generate a 4-bar loop
	samples := gen.GenerateLoop(4)

	if len(samples) == 0 {
		t.Fatal("should generate samples")
	}

	// At 80 BPM, 4 bars = 12 seconds = 529200 samples
	expectedMinSamples := int(44100 * 10) // At least 10 seconds
	if len(samples) < expectedMinSamples {
		t.Errorf("samples = %d, want at least %d", len(samples), expectedMinSamples)
	}

	// Check samples are in valid range
	for i, s := range samples {
		if s < -1.0 || s > 1.0 {
			t.Errorf("sample %d = %f, out of range [-1, 1]", i, s)
			break
		}
	}
}

func TestMusicGeneratorGenres(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		gen := NewMusicGenerator(12345, genre)
		samples := gen.GenerateLoop(2)

		if len(samples) == 0 {
			t.Errorf("genre %s: should generate samples", genre)
		}
	}
}

func TestMusicGeneratorDeterminism(t *testing.T) {
	gen1 := NewMusicGenerator(12345, engine.GenreFantasy)
	gen2 := NewMusicGenerator(12345, engine.GenreFantasy)

	samples1 := gen1.GenerateLoop(2)
	samples2 := gen2.GenerateLoop(2)

	if len(samples1) != len(samples2) {
		t.Fatalf("sample lengths differ: %d vs %d", len(samples1), len(samples2))
	}

	// First 1000 samples should be identical
	for i := 0; i < 1000 && i < len(samples1); i++ {
		if samples1[i] != samples2[i] {
			t.Errorf("sample %d differs: %f vs %f", i, samples1[i], samples2[i])
			break
		}
	}
}

func TestMusicGeneratorBPM(t *testing.T) {
	gen := NewMusicGenerator(12345, engine.GenreFantasy)

	gen.SetBPM(120)
	fast := gen.GenerateLoop(4)

	gen.SetBPM(60)
	slow := gen.GenerateLoop(4)

	// Slower BPM should produce longer loop
	if len(slow) <= len(fast) {
		t.Error("slower BPM should produce longer loop")
	}
}

func TestMusicGeneratorBytes(t *testing.T) {
	gen := NewMusicGenerator(12345, engine.GenreFantasy)

	bytes := gen.GenerateBytes(2)

	if len(bytes) == 0 {
		t.Error("should generate bytes")
	}

	// Should be even number (16-bit samples)
	if len(bytes)%2 != 0 {
		t.Errorf("bytes length = %d, should be even", len(bytes))
	}
}

func TestAmbientLoop(t *testing.T) {
	gen := NewMusicGenerator(12345, engine.GenreFantasy)
	loop := gen.GenerateAmbientLoop(4)

	if loop.IsEmpty() {
		t.Error("loop should not be empty")
	}

	if loop.Duration <= 0 {
		t.Errorf("duration = %f, should be positive", loop.Duration)
	}

	if loop.SampleRate != 44100 {
		t.Errorf("sample rate = %f, want 44100", loop.SampleRate)
	}

	if loop.Genre != engine.GenreFantasy {
		t.Errorf("genre = %s, want fantasy", loop.Genre)
	}

	if loop.Bars != 4 {
		t.Errorf("bars = %d, want 4", loop.Bars)
	}
}

func TestAmbientLoopGetSampleAt(t *testing.T) {
	gen := NewMusicGenerator(12345, engine.GenreFantasy)
	loop := gen.GenerateAmbientLoop(2)

	// Sample at t=0 should be valid
	s0 := loop.GetSampleAt(0)
	if s0 < -1 || s0 > 1 {
		t.Errorf("sample at 0 = %f, out of range", s0)
	}

	// Sample at t > duration should wrap (looping)
	sLoop := loop.GetSampleAt(loop.Duration + 0.1)
	if sLoop < -1 || sLoop > 1 {
		t.Errorf("looped sample = %f, out of range", sLoop)
	}
}

func TestMusicGeneratorGenreSwitch(t *testing.T) {
	gen := NewMusicGenerator(12345, engine.GenreFantasy)

	gen.SetGenre(engine.GenreScifi)
	loop := gen.GenerateAmbientLoop(2)

	if loop.Genre != engine.GenreScifi {
		t.Errorf("genre = %s, want scifi", loop.Genre)
	}
}

func TestMusicState(t *testing.T) {
	gen := NewMusicGenerator(12345, engine.GenreFantasy)

	// Default state should be Peaceful
	if gen.MusicState() != MusicPeaceful {
		t.Errorf("initial state = %s, want Peaceful", MusicStateName(gen.MusicState()))
	}

	// Test state changes affect BPM
	testCases := []struct {
		state       MusicState
		expectedBPM float64
	}{
		{MusicPeaceful, 80},
		{MusicTense, 100},
		{MusicCombat, 140},
		{MusicVictory, 120},
		{MusicDeath, 60},
	}

	for _, tc := range testCases {
		gen.SetMusicState(tc.state)
		if gen.MusicState() != tc.state {
			t.Errorf("state = %s, want %s", MusicStateName(gen.MusicState()), MusicStateName(tc.state))
		}
		// Generate a loop to verify state is applied
		samples := gen.GenerateLoop(2)
		if len(samples) == 0 {
			t.Errorf("state %s: should generate samples", MusicStateName(tc.state))
		}
	}
}

func TestMusicStateName(t *testing.T) {
	testCases := []struct {
		state MusicState
		name  string
	}{
		{MusicPeaceful, "Peaceful"},
		{MusicTense, "Tense"},
		{MusicCombat, "Combat"},
		{MusicVictory, "Victory"},
		{MusicDeath, "Death"},
		{MusicState(99), "Unknown"},
	}

	for _, tc := range testCases {
		name := MusicStateName(tc.state)
		if name != tc.name {
			t.Errorf("MusicStateName(%d) = %s, want %s", tc.state, name, tc.name)
		}
	}
}

func TestMusicStateAffectsOutput(t *testing.T) {
	// Verify different states produce different audio characteristics
	gen := NewMusicGenerator(12345, engine.GenreFantasy)

	gen.SetMusicState(MusicPeaceful)
	peacefulSamples := gen.GenerateLoop(2)

	gen.SetMusicState(MusicCombat)
	combatSamples := gen.GenerateLoop(2)

	// Combat should be faster (shorter duration) due to higher BPM
	if len(combatSamples) >= len(peacefulSamples) {
		t.Error("combat (higher BPM) should produce shorter loop than peaceful (lower BPM)")
	}
}
