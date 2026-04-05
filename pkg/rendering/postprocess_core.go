package rendering

import "github.com/opd-ai/voyage/pkg/engine"

// PostProcessorConfig holds genre-specific effect settings.
type PostProcessorConfig struct {
	VignetteOn   bool
	VignetteInt  float64
	ScanlinesOn  bool
	ScanlinesDen float64
	FilmGrainOn  bool
	FilmGrainInt float64
	ChromaticOn  bool
	ChromaticOff float64
	SepiaOn      bool
	SepiaInt     float64
}

// DefaultPostProcessorConfig returns default values for post processing.
func DefaultPostProcessorConfig() PostProcessorConfig {
	return PostProcessorConfig{
		VignetteInt:  0.3,
		ScanlinesDen: 2.0,
		FilmGrainInt: 0.15,
		ChromaticOff: 2.0,
		SepiaInt:     0.5,
	}
}

// ConfigureForGenre returns effect settings for the given genre.
func ConfigureForGenre(genre engine.GenreID) PostProcessorConfig {
	cfg := DefaultPostProcessorConfig()

	switch genre {
	case engine.GenreFantasy:
		cfg.VignetteOn = true
		cfg.VignetteInt = 0.2
	case engine.GenreScifi:
		cfg.VignetteOn = true
		cfg.VignetteInt = 0.3
		cfg.ScanlinesOn = true
	case engine.GenreHorror:
		cfg.VignetteOn = true
		cfg.VignetteInt = 0.5
		cfg.FilmGrainOn = true
	case engine.GenreCyberpunk:
		cfg.VignetteOn = true
		cfg.VignetteInt = 0.4
		cfg.ScanlinesOn = true
		cfg.ChromaticOn = true
	case engine.GenrePostapoc:
		cfg.VignetteOn = true
		cfg.VignetteInt = 0.35
		cfg.SepiaOn = true
	}

	return cfg
}

// clampUint8 clamps a float64 to uint8 range.
func clampUint8(v float64) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}
