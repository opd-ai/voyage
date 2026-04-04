package audio

// Envelope implements an ADSR (Attack, Decay, Sustain, Release) envelope.
type Envelope struct {
	attack  float64 // Attack time in seconds
	decay   float64 // Decay time in seconds
	sustain float64 // Sustain level [0, 1]
	release float64 // Release time in seconds

	state      EnvelopeState
	level      float64
	time       float64
	sampleRate float64
	released   bool
}

// EnvelopeState represents the current phase of the envelope.
type EnvelopeState int

const (
	// EnvelopeIdle is the initial state before note on.
	EnvelopeIdle EnvelopeState = iota
	// EnvelopeAttack is the rising phase.
	EnvelopeAttack
	// EnvelopeDecay is the falling phase to sustain level.
	EnvelopeDecay
	// EnvelopeSustain holds at sustain level until note off.
	EnvelopeSustain
	// EnvelopeRelease falls to zero after note off.
	EnvelopeRelease
)

// NewEnvelope creates an ADSR envelope with the given parameters.
// Times are in seconds, sustain is a level from 0 to 1.
func NewEnvelope(attack, decay, sustain, release float64) *Envelope {
	return &Envelope{
		attack:     attack,
		decay:      decay,
		sustain:    sustain,
		release:    release,
		state:      EnvelopeIdle,
		level:      0,
		time:       0,
		sampleRate: 44100,
		released:   false,
	}
}

// NoteOn triggers the envelope attack phase.
func (e *Envelope) NoteOn() {
	e.state = EnvelopeAttack
	e.time = 0
	e.released = false
}

// NoteOff triggers the envelope release phase.
func (e *Envelope) NoteOff() {
	if e.state != EnvelopeIdle && e.state != EnvelopeRelease {
		e.state = EnvelopeRelease
		e.time = 0
		e.released = true
	}
}

// Sample returns the current envelope level and advances.
func (e *Envelope) Sample() float64 {
	dt := 1.0 / e.sampleRate

	switch e.state {
	case EnvelopeIdle:
		e.level = 0

	case EnvelopeAttack:
		if e.attack > 0 {
			e.level = e.time / e.attack
			if e.level >= 1.0 {
				e.level = 1.0
				e.state = EnvelopeDecay
				e.time = 0
			}
		} else {
			e.level = 1.0
			e.state = EnvelopeDecay
			e.time = 0
		}

	case EnvelopeDecay:
		if e.decay > 0 {
			e.level = 1.0 - (1.0-e.sustain)*(e.time/e.decay)
			if e.time >= e.decay {
				e.level = e.sustain
				e.state = EnvelopeSustain
				e.time = 0
			}
		} else {
			e.level = e.sustain
			e.state = EnvelopeSustain
			e.time = 0
		}

	case EnvelopeSustain:
		e.level = e.sustain

	case EnvelopeRelease:
		if e.release > 0 {
			startLevel := e.sustain
			e.level = startLevel * (1.0 - e.time/e.release)
			if e.level <= 0 {
				e.level = 0
				e.state = EnvelopeIdle
			}
		} else {
			e.level = 0
			e.state = EnvelopeIdle
		}
	}

	e.time += dt
	return e.level
}

// IsActive returns true if the envelope is producing output.
func (e *Envelope) IsActive() bool {
	return e.state != EnvelopeIdle
}

// IsReleased returns true if note off has been triggered.
func (e *Envelope) IsReleased() bool {
	return e.released
}

// Level returns the current envelope level without advancing.
func (e *Envelope) Level() float64 {
	return e.level
}

// State returns the current envelope state.
func (e *Envelope) State() EnvelopeState {
	return e.state
}

// Reset returns the envelope to idle state.
func (e *Envelope) Reset() {
	e.state = EnvelopeIdle
	e.level = 0
	e.time = 0
	e.released = false
}

// SetADSR updates the envelope parameters.
func (e *Envelope) SetADSR(attack, decay, sustain, release float64) {
	e.attack = attack
	e.decay = decay
	e.sustain = sustain
	e.release = release
}

// QuickEnvelope creates an envelope suitable for SFX.
func QuickEnvelope() *Envelope {
	return NewEnvelope(0.01, 0.1, 0.3, 0.2)
}

// SlowEnvelope creates an envelope for ambient sounds.
func SlowEnvelope() *Envelope {
	return NewEnvelope(0.5, 0.3, 0.6, 1.0)
}

// PunchyEnvelope creates an envelope for percussive sounds.
func PunchyEnvelope() *Envelope {
	return NewEnvelope(0.001, 0.05, 0.0, 0.1)
}
