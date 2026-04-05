package audio

import "math"

// SpatialAudioConfig holds configuration for spatial audio processing.
type SpatialAudioConfig struct {
	// ListenerX is the X position of the listener.
	ListenerX float64
	// ListenerY is the Y position of the listener.
	ListenerY float64
	// SourceX is the X position of the audio source.
	SourceX float64
	// SourceY is the Y position of the audio source.
	SourceY float64
	// MaxDistance is the distance at which audio is fully attenuated.
	MaxDistance float64
	// RolloffFactor controls how quickly sound attenuates with distance.
	// 1.0 is realistic, <1.0 is slower, >1.0 is faster.
	RolloffFactor float64
	// RefDistance is the distance at which audio is at full volume.
	RefDistance float64
}

// DefaultSpatialConfig returns sensible defaults for spatial audio.
func DefaultSpatialConfig() SpatialAudioConfig {
	return SpatialAudioConfig{
		ListenerX:     0,
		ListenerY:     0,
		SourceX:       0,
		SourceY:       0,
		MaxDistance:   100.0,
		RolloffFactor: 1.0,
		RefDistance:   1.0,
	}
}

// ApplySpatialAudio applies spatial audio effects to mono samples.
// Returns stereo (left, right) sample arrays with distance attenuation and stereo panning.
func ApplySpatialAudio(samples []float64, config SpatialAudioConfig) (left, right []float64) {
	// Calculate distance
	dx := config.SourceX - config.ListenerX
	dy := config.SourceY - config.ListenerY
	distance := math.Sqrt(dx*dx + dy*dy)

	// Calculate attenuation using inverse distance model
	attenuation := calculateAttenuation(distance, config)

	// Calculate stereo panning
	panLeft, panRight := calculatePanning(dx, distance)

	// Apply to samples
	left = make([]float64, len(samples))
	right = make([]float64, len(samples))

	for i, s := range samples {
		attenuatedSample := s * attenuation
		left[i] = attenuatedSample * panLeft
		right[i] = attenuatedSample * panRight
	}

	return left, right
}

// calculateAttenuation computes the volume attenuation based on distance.
// Uses inverse distance clamped model.
func calculateAttenuation(distance float64, config SpatialAudioConfig) float64 {
	if distance <= config.RefDistance {
		return 1.0
	}
	if distance >= config.MaxDistance {
		return 0.0
	}

	// Inverse distance model: gain = refDist / (refDist + rolloff * (dist - refDist))
	gain := config.RefDistance / (config.RefDistance + config.RolloffFactor*(distance-config.RefDistance))

	// Clamp to [0, 1]
	if gain < 0 {
		return 0
	}
	if gain > 1 {
		return 1
	}
	return gain
}

// calculatePanning computes the stereo panning based on horizontal offset.
// Returns left and right channel gains using constant-power panning.
func calculatePanning(dx, distance float64) (left, right float64) {
	if distance <= 0 {
		// Source at listener position - center pan
		return 0.707, 0.707
	}

	// Calculate pan position in range [-1, 1]
	// Negative dx = source to the left, positive = source to the right
	pan := dx / distance
	if pan < -1 {
		pan = -1
	}
	if pan > 1 {
		pan = 1
	}

	// Constant-power panning for smooth transitions
	// pan = -1 -> left only, pan = 0 -> center, pan = 1 -> right only
	angle := (pan + 1) * math.Pi / 4 // Map to [0, π/2]
	left = math.Cos(angle)
	right = math.Sin(angle)

	return left, right
}

// ApplySpatialAudioStereo applies spatial effects to already-stereo samples.
func ApplySpatialAudioStereo(leftIn, rightIn []float64, config SpatialAudioConfig) (leftOut, rightOut []float64) {
	// First mix to mono
	mono := make([]float64, len(leftIn))
	for i := range mono {
		mono[i] = (leftIn[i] + rightIn[i]) * 0.5
	}

	// Then apply spatial
	return ApplySpatialAudio(mono, config)
}

// Distance calculates the distance between listener and source.
func (c SpatialAudioConfig) Distance() float64 {
	dx := c.SourceX - c.ListenerX
	dy := c.SourceY - c.ListenerY
	return math.Sqrt(dx*dx + dy*dy)
}

// Attenuation returns the current attenuation factor for this config.
func (c SpatialAudioConfig) Attenuation() float64 {
	return calculateAttenuation(c.Distance(), c)
}

// Pan returns the current pan values (left, right) for this config.
func (c SpatialAudioConfig) Pan() (left, right float64) {
	dx := c.SourceX - c.ListenerX
	return calculatePanning(dx, c.Distance())
}

// SpatialAudioProcessor manages spatial audio for multiple sources.
type SpatialAudioProcessor struct {
	listenerX float64
	listenerY float64
}

// NewSpatialAudioProcessor creates a new spatial audio processor.
func NewSpatialAudioProcessor() *SpatialAudioProcessor {
	return &SpatialAudioProcessor{
		listenerX: 0,
		listenerY: 0,
	}
}

// SetListenerPosition updates the listener's position.
func (p *SpatialAudioProcessor) SetListenerPosition(x, y float64) {
	p.listenerX = x
	p.listenerY = y
}

// ListenerPosition returns the current listener position.
func (p *SpatialAudioProcessor) ListenerPosition() (x, y float64) {
	return p.listenerX, p.listenerY
}

// ProcessSource applies spatial audio to a source at the given position.
func (p *SpatialAudioProcessor) ProcessSource(samples []float64, sourceX, sourceY, maxDistance float64) (left, right []float64) {
	config := SpatialAudioConfig{
		ListenerX:     p.listenerX,
		ListenerY:     p.listenerY,
		SourceX:       sourceX,
		SourceY:       sourceY,
		MaxDistance:   maxDistance,
		RolloffFactor: 1.0,
		RefDistance:   1.0,
	}
	return ApplySpatialAudio(samples, config)
}
