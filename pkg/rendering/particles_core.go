package rendering

import (
	"image/color"
	"math"
	"math/rand"
)

// ParticleType categorizes different particle effects.
type ParticleType int

const (
	// ParticleTypeDust represents dust clouds from movement on land.
	ParticleTypeDust ParticleType = iota
	// ParticleTypeThruster represents exhaust from space vessels.
	ParticleTypeThruster
	// ParticleTypeTireTrack represents tracks from wheeled vehicles.
	ParticleTypeTireTrack
	// ParticleTypeRain represents rain drops.
	ParticleTypeRain
	// ParticleTypeSnow represents snowflakes.
	ParticleTypeSnow
	// ParticleTypeSand represents sandstorm particles.
	ParticleTypeSand
	// ParticleTypeEmbers represents floating fire embers.
	ParticleTypeEmbers
	// ParticleTypeAsh represents falling ash.
	ParticleTypeAsh
	// ParticleTypeSparks represents combat/impact sparks.
	ParticleTypeSparks
	// ParticleTypeHeal represents healing glow particles.
	ParticleTypeHeal
	// ParticleTypeExplosion represents explosion debris.
	ParticleTypeExplosion
)

// Particle represents a single particle in the system.
type Particle struct {
	X, Y     float64      // Position
	VX, VY   float64      // Velocity
	Life     float64      // Remaining lifetime (0-1)
	MaxLife  float64      // Initial lifetime
	Size     float64      // Particle size
	Color    color.RGBA   // Particle color
	Type     ParticleType // Type of particle
	Rotation float64      // Rotation angle
	RotSpeed float64      // Rotation speed
	Alpha    float64      // Current alpha (modified by life)
	FadeIn   bool         // Whether particle fades in at start
}

// ParticleEmitter generates particles at a location.
type ParticleEmitter struct {
	X, Y         float64      // Emitter position
	Type         ParticleType // Type of particles to emit
	Rate         float64      // Particles per second
	Burst        int          // Particles per burst (0 for continuous)
	SpreadAngle  float64      // Spread angle in radians
	BaseVelocity float64      // Base particle velocity
	VelocityVar  float64      // Velocity variation
	BaseLife     float64      // Base particle lifetime
	LifeVar      float64      // Lifetime variation
	BaseSize     float64      // Base particle size
	SizeVar      float64      // Size variation
	Active       bool         // Whether emitter is active
	Accumulator  float64      // Time accumulator for emission
	BaseColor    color.RGBA   // Base particle color
}

// ParticlePreset contains genre-specific particle settings.
type ParticlePreset struct {
	MovementTrailType ParticleType // Type for movement trails
	DustColor         color.RGBA
	ThrusterColor     color.RGBA
	TireTrackColor    color.RGBA
	RainColor         color.RGBA
	SnowColor         color.RGBA
	SandColor         color.RGBA
	EmberColor        color.RGBA
	AshColor          color.RGBA
	SparksColor       color.RGBA
	HealColor         color.RGBA
	ExplosionColor    color.RGBA
}

// BurstParams holds computed parameters for a particle burst.
type BurstParams struct {
	Life     float64
	Velocity float64
	Size     float64
}

// ComputeBurstParams calculates the life, velocity, and size for a burst particle.
func ComputeBurstParams(pType ParticleType, rng *rand.Rand) BurstParams {
	velocity := 20.0 + rng.Float64()*30.0
	size := 2.0 + rng.Float64()*3.0
	var life float64

	switch pType {
	case ParticleTypeSparks:
		life = 0.3 + rng.Float64()*0.2
		velocity = 50.0 + rng.Float64()*50.0
	case ParticleTypeExplosion:
		life = 0.5 + rng.Float64()*0.5
		velocity = 30.0 + rng.Float64()*40.0
		size = 4.0 + rng.Float64()*4.0
	case ParticleTypeHeal:
		life = 1.0 + rng.Float64()*0.5
		velocity = 10.0 + rng.Float64()*10.0
	default:
		life = 0.5 + rng.Float64()*0.5
	}

	return BurstParams{
		Life:     life,
		Velocity: velocity,
		Size:     size,
	}
}

// CreateBurstParticle creates a particle for a burst effect at the given position.
func CreateBurstParticle(x, y float64, pType ParticleType, particleColor color.RGBA, rng *rand.Rand) *Particle {
	params := ComputeBurstParams(pType, rng)
	angle := rng.Float64() * math.Pi * 2

	return &Particle{
		X:        x + (rng.Float64()-0.5)*8,
		Y:        y + (rng.Float64()-0.5)*8,
		VX:       params.Velocity * math.Cos(angle),
		VY:       params.Velocity * math.Sin(angle),
		Life:     1.0,
		MaxLife:  params.Life,
		Size:     params.Size,
		Color:    particleColor,
		Type:     pType,
		Rotation: rng.Float64() * math.Pi * 2,
		RotSpeed: (rng.Float64() - 0.5) * 4,
		Alpha:    1.0,
		FadeIn:   pType == ParticleTypeHeal,
	}
}

// CreateEmitterParticle creates a particle from an emitter.
func CreateEmitterParticle(emitter *ParticleEmitter, particleColor color.RGBA, rng *rand.Rand) *Particle {
	angle := (rng.Float64() - 0.5) * emitter.SpreadAngle
	velocity := emitter.BaseVelocity + (rng.Float64()-0.5)*2*emitter.VelocityVar
	life := emitter.BaseLife + (rng.Float64()-0.5)*2*emitter.LifeVar
	size := emitter.BaseSize + (rng.Float64()-0.5)*2*emitter.SizeVar

	return &Particle{
		X:        emitter.X + (rng.Float64()-0.5)*4,
		Y:        emitter.Y + (rng.Float64()-0.5)*4,
		VX:       velocity * math.Cos(angle),
		VY:       velocity * math.Sin(angle),
		Life:     1.0,
		MaxLife:  life,
		Size:     size,
		Color:    particleColor,
		Type:     emitter.Type,
		Rotation: rng.Float64() * math.Pi * 2,
		RotSpeed: (rng.Float64() - 0.5) * 2,
		Alpha:    1.0,
		FadeIn:   emitter.Type == ParticleTypeHeal,
	}
}
