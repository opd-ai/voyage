//go:build headless

package rendering

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/opd-ai/voyage/pkg/engine"
)

// ParticleType categorizes different particle effects.
type ParticleType int

const (
	ParticleTypeDust ParticleType = iota
	ParticleTypeThruster
	ParticleTypeTireTrack
	ParticleTypeRain
	ParticleTypeSnow
	ParticleTypeSand
	ParticleTypeEmbers
	ParticleTypeAsh
	ParticleTypeSparks
	ParticleTypeHeal
	ParticleTypeExplosion
)

// Particle represents a single particle in the system.
type Particle struct {
	X, Y     float64
	VX, VY   float64
	Life     float64
	MaxLife  float64
	Size     float64
	Color    color.RGBA
	Type     ParticleType
	Rotation float64
	RotSpeed float64
	Alpha    float64
	FadeIn   bool
}

// ParticleEmitter generates particles at a location.
type ParticleEmitter struct {
	X, Y         float64
	Type         ParticleType
	Rate         float64
	Burst        int
	SpreadAngle  float64
	BaseVelocity float64
	VelocityVar  float64
	BaseLife     float64
	LifeVar      float64
	BaseSize     float64
	SizeVar      float64
	Active       bool
	accumulator  float64
	baseColor    color.RGBA
}

// ParticleSystem manages all particles and emitters.
type ParticleSystem struct {
	engine.BaseSystem
	particles    []*Particle
	emitters     []*ParticleEmitter
	rng          *rand.Rand
	genrePreset  *ParticlePreset
	maxParticles int
}

// ParticlePreset contains genre-specific particle settings.
type ParticlePreset struct {
	MovementTrailType ParticleType
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

// NewParticleSystem creates a new particle system.
func NewParticleSystem(seed int64) *ParticleSystem {
	return &ParticleSystem{
		BaseSystem:   engine.NewBaseSystem(engine.PriorityRender - 2),
		particles:    make([]*Particle, 0, 1000),
		emitters:     make([]*ParticleEmitter, 0, 50),
		rng:          rand.New(rand.NewSource(seed)),
		genrePreset:  defaultParticlePreset(engine.GenreFantasy),
		maxParticles: 2000,
	}
}

// SetGenre changes the particle presets to match the genre.
func (ps *ParticleSystem) SetGenre(genreID engine.GenreID) {
	ps.BaseSystem.SetGenre(genreID)
	ps.genrePreset = defaultParticlePreset(genreID)
}

// defaultParticlePreset returns genre-specific particle settings.
func defaultParticlePreset(genre engine.GenreID) *ParticlePreset {
	switch genre {
	case engine.GenreScifi:
		return &ParticlePreset{
			MovementTrailType: ParticleTypeThruster,
			DustColor:         color.RGBA{100, 100, 120, 200},
			ThrusterColor:     color.RGBA{100, 180, 255, 255},
			TireTrackColor:    color.RGBA{80, 80, 100, 150},
			RainColor:         color.RGBA{150, 200, 255, 180},
			SnowColor:         color.RGBA{200, 220, 255, 200},
			SandColor:         color.RGBA{180, 160, 140, 180},
			EmberColor:        color.RGBA{255, 200, 100, 255},
			AshColor:          color.RGBA{100, 100, 100, 180},
			SparksColor:       color.RGBA{200, 230, 255, 255},
			HealColor:         color.RGBA{100, 255, 200, 200},
			ExplosionColor:    color.RGBA{255, 200, 150, 255},
		}
	case engine.GenreHorror:
		return &ParticlePreset{
			MovementTrailType: ParticleTypeDust,
			DustColor:         color.RGBA{80, 60, 50, 180},
			ThrusterColor:     color.RGBA{200, 100, 80, 200},
			TireTrackColor:    color.RGBA{60, 50, 40, 150},
			RainColor:         color.RGBA{100, 80, 80, 180},
			SnowColor:         color.RGBA{180, 170, 160, 200},
			SandColor:         color.RGBA{150, 120, 100, 180},
			EmberColor:        color.RGBA{255, 100, 50, 255},
			AshColor:          color.RGBA{60, 50, 50, 180},
			SparksColor:       color.RGBA{255, 150, 100, 255},
			HealColor:         color.RGBA{150, 255, 150, 180},
			ExplosionColor:    color.RGBA{200, 80, 50, 255},
		}
	case engine.GenreCyberpunk:
		return &ParticlePreset{
			MovementTrailType: ParticleTypeThruster,
			DustColor:         color.RGBA{100, 80, 120, 180},
			ThrusterColor:     color.RGBA{255, 100, 200, 255},
			TireTrackColor:    color.RGBA{80, 60, 100, 150},
			RainColor:         color.RGBA{150, 180, 200, 180},
			SnowColor:         color.RGBA{200, 200, 220, 200},
			SandColor:         color.RGBA{160, 140, 120, 180},
			EmberColor:        color.RGBA{255, 200, 50, 255},
			AshColor:          color.RGBA{80, 70, 90, 180},
			SparksColor:       color.RGBA{0, 255, 200, 255},
			HealColor:         color.RGBA{200, 100, 255, 200},
			ExplosionColor:    color.RGBA{255, 150, 200, 255},
		}
	case engine.GenrePostapoc:
		return &ParticlePreset{
			MovementTrailType: ParticleTypeDust,
			DustColor:         color.RGBA{180, 150, 100, 200},
			ThrusterColor:     color.RGBA{200, 150, 80, 200},
			TireTrackColor:    color.RGBA{140, 120, 80, 150},
			RainColor:         color.RGBA{180, 170, 150, 180},
			SnowColor:         color.RGBA{200, 190, 170, 200},
			SandColor:         color.RGBA{200, 160, 100, 200},
			EmberColor:        color.RGBA{255, 180, 80, 255},
			AshColor:          color.RGBA{120, 110, 90, 200},
			SparksColor:       color.RGBA{255, 200, 100, 255},
			HealColor:         color.RGBA{100, 200, 100, 180},
			ExplosionColor:    color.RGBA{255, 180, 100, 255},
		}
	default:
		return &ParticlePreset{
			MovementTrailType: ParticleTypeDust,
			DustColor:         color.RGBA{180, 160, 140, 180},
			ThrusterColor:     color.RGBA{255, 200, 100, 200},
			TireTrackColor:    color.RGBA{120, 100, 80, 150},
			RainColor:         color.RGBA{150, 180, 220, 180},
			SnowColor:         color.RGBA{240, 245, 255, 220},
			SandColor:         color.RGBA{200, 180, 140, 180},
			EmberColor:        color.RGBA{255, 200, 100, 255},
			AshColor:          color.RGBA{100, 100, 100, 180},
			SparksColor:       color.RGBA{255, 220, 150, 255},
			HealColor:         color.RGBA{150, 255, 200, 200},
			ExplosionColor:    color.RGBA{255, 180, 100, 255},
		}
	}
}

// Update updates all particles and emitters.
func (ps *ParticleSystem) Update(world *engine.World, dt float64) {
	for _, emitter := range ps.emitters {
		if emitter.Active {
			ps.updateEmitter(emitter, dt)
		}
	}
	ps.updateParticles(dt)
}

// updateEmitter processes particle emission.
func (ps *ParticleSystem) updateEmitter(emitter *ParticleEmitter, dt float64) {
	emitter.accumulator += dt
	particlesToEmit := int(emitter.accumulator * emitter.Rate)
	if particlesToEmit > 0 {
		emitter.accumulator -= float64(particlesToEmit) / emitter.Rate
		for i := 0; i < particlesToEmit && len(ps.particles) < ps.maxParticles; i++ {
			ps.emitParticle(emitter)
		}
	}
}

// emitParticle creates a new particle from an emitter.
func (ps *ParticleSystem) emitParticle(emitter *ParticleEmitter) {
	angle := (ps.rng.Float64() - 0.5) * emitter.SpreadAngle
	velocity := emitter.BaseVelocity + (ps.rng.Float64()-0.5)*2*emitter.VelocityVar
	life := emitter.BaseLife + (ps.rng.Float64()-0.5)*2*emitter.LifeVar
	size := emitter.BaseSize + (ps.rng.Float64()-0.5)*2*emitter.SizeVar

	particle := &Particle{
		X:        emitter.X + (ps.rng.Float64()-0.5)*4,
		Y:        emitter.Y + (ps.rng.Float64()-0.5)*4,
		VX:       velocity * math.Cos(angle),
		VY:       velocity * math.Sin(angle),
		Life:     1.0,
		MaxLife:  life,
		Size:     size,
		Color:    ps.getColorForType(emitter.Type),
		Type:     emitter.Type,
		Rotation: ps.rng.Float64() * math.Pi * 2,
		RotSpeed: (ps.rng.Float64() - 0.5) * 2,
		Alpha:    1.0,
		FadeIn:   emitter.Type == ParticleTypeHeal,
	}
	ps.particles = append(ps.particles, particle)
}

// updateParticles updates all particles and removes dead ones.
func (ps *ParticleSystem) updateParticles(dt float64) {
	alive := ps.particles[:0]
	for _, p := range ps.particles {
		p.Life -= dt / p.MaxLife
		if p.Life <= 0 {
			continue
		}
		ps.updateParticlePhysics(p, dt)
		ps.updateParticleAlpha(p)
		alive = append(alive, p)
	}
	ps.particles = alive
}

// updateParticlePhysics applies movement, gravity, and drag to a particle.
func (ps *ParticleSystem) updateParticlePhysics(p *Particle, dt float64) {
	p.X += p.VX * dt
	p.Y += p.VY * dt

	if p.Type == ParticleTypeRain || p.Type == ParticleTypeSnow || p.Type == ParticleTypeAsh {
		p.VY += 50 * dt
	}
	if p.Type == ParticleTypeDust || p.Type == ParticleTypeSand {
		p.VX *= 0.98
		p.VY *= 0.98
	}
	p.Rotation += p.RotSpeed * dt
}

// updateParticleAlpha updates particle transparency based on lifecycle.
func (ps *ParticleSystem) updateParticleAlpha(p *Particle) {
	if p.FadeIn && p.Life > 0.8 {
		p.Alpha = (1.0 - p.Life) / 0.2
	} else if p.Life < 0.3 {
		p.Alpha = p.Life / 0.3
	} else {
		p.Alpha = 1.0
	}
}

// getColorForType returns the preset color for a particle type.
func (ps *ParticleSystem) getColorForType(pType ParticleType) color.RGBA {
	switch pType {
	case ParticleTypeDust:
		return ps.genrePreset.DustColor
	case ParticleTypeThruster:
		return ps.genrePreset.ThrusterColor
	case ParticleTypeTireTrack:
		return ps.genrePreset.TireTrackColor
	case ParticleTypeRain:
		return ps.genrePreset.RainColor
	case ParticleTypeSnow:
		return ps.genrePreset.SnowColor
	case ParticleTypeSand:
		return ps.genrePreset.SandColor
	case ParticleTypeEmbers:
		return ps.genrePreset.EmberColor
	case ParticleTypeAsh:
		return ps.genrePreset.AshColor
	case ParticleTypeSparks:
		return ps.genrePreset.SparksColor
	case ParticleTypeHeal:
		return ps.genrePreset.HealColor
	case ParticleTypeExplosion:
		return ps.genrePreset.ExplosionColor
	default:
		return color.RGBA{255, 255, 255, 255}
	}
}

// AddEmitter adds a particle emitter to the system.
func (ps *ParticleSystem) AddEmitter(emitter *ParticleEmitter) {
	if emitter != nil {
		ps.emitters = append(ps.emitters, emitter)
	}
}

// RemoveEmitter removes a particle emitter from the system.
func (ps *ParticleSystem) RemoveEmitter(emitter *ParticleEmitter) {
	for i, e := range ps.emitters {
		if e == emitter {
			ps.emitters = append(ps.emitters[:i], ps.emitters[i+1:]...)
			return
		}
	}
}

// ClearEmitters removes all emitters.
func (ps *ParticleSystem) ClearEmitters() {
	ps.emitters = ps.emitters[:0]
}

// ClearParticles removes all particles.
func (ps *ParticleSystem) ClearParticles() {
	ps.particles = ps.particles[:0]
}

// ParticleCount returns the current number of active particles.
func (ps *ParticleSystem) ParticleCount() int {
	return len(ps.particles)
}

// EmitterCount returns the current number of emitters.
func (ps *ParticleSystem) EmitterCount() int {
	return len(ps.emitters)
}

// CreateMovementTrailEmitter creates an emitter for movement trails.
func (ps *ParticleSystem) CreateMovementTrailEmitter(x, y float64) *ParticleEmitter {
	return &ParticleEmitter{
		X:            x,
		Y:            y,
		Type:         ps.genrePreset.MovementTrailType,
		Rate:         15.0,
		SpreadAngle:  math.Pi / 3,
		BaseVelocity: -20.0,
		VelocityVar:  5.0,
		BaseLife:     0.5,
		LifeVar:      0.2,
		BaseSize:     3.0,
		SizeVar:      1.0,
		Active:       true,
	}
}

// CreateWeatherEmitter creates an emitter for weather particles.
func (ps *ParticleSystem) CreateWeatherEmitter(pType ParticleType, x, y, width float64) *ParticleEmitter {
	var velocity, life, size, rate float64
	switch pType {
	case ParticleTypeRain:
		velocity = 200.0
		life = 1.0
		size = 2.0
		rate = 50.0
	case ParticleTypeSnow:
		velocity = 30.0
		life = 3.0
		size = 3.0
		rate = 20.0
	case ParticleTypeSand:
		velocity = 100.0
		life = 1.5
		size = 2.0
		rate = 40.0
	case ParticleTypeAsh:
		velocity = 20.0
		life = 4.0
		size = 2.0
		rate = 15.0
	default:
		velocity = 50.0
		life = 1.0
		size = 2.0
		rate = 20.0
	}
	return &ParticleEmitter{
		X:            x + width/2,
		Y:            y,
		Type:         pType,
		Rate:         rate,
		SpreadAngle:  math.Pi / 6,
		BaseVelocity: velocity,
		VelocityVar:  velocity * 0.2,
		BaseLife:     life,
		LifeVar:      life * 0.3,
		BaseSize:     size,
		SizeVar:      size * 0.3,
		Active:       true,
	}
}

// EmitBurst emits a burst of particles at a location.
func (ps *ParticleSystem) EmitBurst(x, y float64, pType ParticleType, count int) {
	for i := 0; i < count && len(ps.particles) < ps.maxParticles; i++ {
		angle := ps.rng.Float64() * math.Pi * 2
		velocity := 20.0 + ps.rng.Float64()*30.0
		size := 2.0 + ps.rng.Float64()*3.0
		var life float64
		switch pType {
		case ParticleTypeSparks:
			life = 0.3 + ps.rng.Float64()*0.2
			velocity = 50.0 + ps.rng.Float64()*50.0
		case ParticleTypeExplosion:
			life = 0.5 + ps.rng.Float64()*0.5
			velocity = 30.0 + ps.rng.Float64()*40.0
			size = 4.0 + ps.rng.Float64()*4.0
		case ParticleTypeHeal:
			life = 1.0 + ps.rng.Float64()*0.5
			velocity = 10.0 + ps.rng.Float64()*10.0
		default:
			life = 0.5 + ps.rng.Float64()*0.5
		}
		particle := &Particle{
			X:        x + (ps.rng.Float64()-0.5)*8,
			Y:        y + (ps.rng.Float64()-0.5)*8,
			VX:       velocity * math.Cos(angle),
			VY:       velocity * math.Sin(angle),
			Life:     1.0,
			MaxLife:  life,
			Size:     size,
			Color:    ps.getColorForType(pType),
			Type:     pType,
			Rotation: ps.rng.Float64() * math.Pi * 2,
			RotSpeed: (ps.rng.Float64() - 0.5) * 4,
			Alpha:    1.0,
			FadeIn:   pType == ParticleTypeHeal,
		}
		ps.particles = append(ps.particles, particle)
	}
}

// Particles returns all current particles.
func (ps *ParticleSystem) Particles() []*Particle {
	return ps.particles
}

// Emitters returns all current emitters.
func (ps *ParticleSystem) Emitters() []*ParticleEmitter {
	return ps.emitters
}
