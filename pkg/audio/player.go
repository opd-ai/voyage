package audio

import "github.com/opd-ai/voyage/pkg/engine"

// Player manages audio playback (stub for Ebitengine integration).
// The actual playback will use Ebitengine's audio package when the
// game loop is running.
type Player struct {
	genre    engine.GenreID
	sfxGen   *SFXGenerator
	musicGen *MusicGenerator
	volume   float64
	muted    bool
}

// NewPlayer creates a new audio player.
func NewPlayer(masterSeed int64, genre engine.GenreID) *Player {
	return &Player{
		genre:    genre,
		sfxGen:   NewSFXGenerator(masterSeed, genre),
		musicGen: NewMusicGenerator(masterSeed, genre),
		volume:   1.0,
		muted:    false,
	}
}

// SetGenre changes the audio theme.
func (p *Player) SetGenre(genre engine.GenreID) {
	p.genre = genre
	p.sfxGen.SetGenre(genre)
	p.musicGen.SetGenre(genre)
}

// Genre returns the current genre.
func (p *Player) Genre() engine.GenreID {
	return p.genre
}

// SetVolume sets the master volume (0.0 to 1.0).
func (p *Player) SetVolume(vol float64) {
	if vol < 0 {
		vol = 0
	}
	if vol > 1 {
		vol = 1
	}
	p.volume = vol
}

// Volume returns the current volume.
func (p *Player) Volume() float64 {
	return p.volume
}

// Mute mutes all audio.
func (p *Player) Mute() {
	p.muted = true
}

// Unmute unmutes audio.
func (p *Player) Unmute() {
	p.muted = false
}

// IsMuted returns true if audio is muted.
func (p *Player) IsMuted() bool {
	return p.muted
}

// PlaySFX generates and plays a sound effect.
// Returns the generated samples (for potential caching).
// Returns empty slice when muted to avoid nil pointer issues.
func (p *Player) PlaySFX(sfxType SFXType) []float64 {
	if p.muted {
		return []float64{}
	}

	samples := p.sfxGen.Generate(sfxType)

	// Apply volume
	for i := range samples {
		samples[i] *= p.volume
	}

	// In actual implementation, this would queue the samples
	// for playback via Ebitengine's audio system.
	// For now, we just return the samples.
	return samples
}

// PreloadSFX generates all SFX for caching.
func (p *Player) PreloadSFX() map[SFXType][]float64 {
	cache := make(map[SFXType][]float64)
	for _, sfx := range AllSFXTypes() {
		cache[sfx] = p.sfxGen.Generate(sfx)
	}
	return cache
}

// SFXGenerator returns the underlying SFX generator.
func (p *Player) SFXGenerator() *SFXGenerator {
	return p.sfxGen
}

// MusicGenerator returns the underlying music generator.
func (p *Player) MusicGenerator() *MusicGenerator {
	return p.musicGen
}

// GenerateAmbientMusic creates a looping ambient music track.
// Returns the generated samples for playback.
// Returns empty AmbientLoop when muted to avoid nil pointer issues.
func (p *Player) GenerateAmbientMusic(bars int) *AmbientLoop {
	if p.muted {
		return &AmbientLoop{Samples: []float64{}}
	}

	loop := p.musicGen.GenerateAmbientLoop(bars)

	// Apply volume to samples
	for i := range loop.Samples {
		loop.Samples[i] *= p.volume
	}

	return loop
}

// SetMusicState changes the music intensity state.
// Use this to transition music based on game events.
func (p *Player) SetMusicState(state MusicState) {
	p.musicGen.SetMusicState(state)
}

// MusicState returns the current music state.
func (p *Player) MusicState() MusicState {
	return p.musicGen.MusicState()
}

// CrossfadeMusicTo generates a crossfade transition to a new music state.
// Returns audio samples for the transition period.
// Returns empty slice when muted to avoid nil pointer issues.
func (p *Player) CrossfadeMusicTo(targetState MusicState, durationMs int) []float64 {
	if p.muted {
		return []float64{}
	}

	samples := p.musicGen.CrossfadeTo(targetState, durationMs)

	// Apply volume
	for i := range samples {
		samples[i] *= p.volume
	}

	return samples
}
