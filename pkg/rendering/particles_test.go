package rendering

import (
	"image/color"
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewParticleSystem(t *testing.T) {
	ps := NewParticleSystem(12345)
	if ps == nil {
		t.Fatal("NewParticleSystem returned nil")
	}
	if ps.ParticleCount() != 0 {
		t.Error("should start with no particles")
	}
	if ps.EmitterCount() != 0 {
		t.Error("should start with no emitters")
	}
}

func TestParticleSystemSetGenre(t *testing.T) {
	ps := NewParticleSystem(12345)
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		ps.SetGenre(genre)
		if ps.Genre() != genre {
			t.Errorf("expected genre %s, got %s", genre, ps.Genre())
		}
	}
}

func TestParticleEmitterManagement(t *testing.T) {
	ps := NewParticleSystem(12345)

	emitter1 := ps.CreateMovementTrailEmitter(0, 0)
	emitter2 := ps.CreateMovementTrailEmitter(10, 10)

	ps.AddEmitter(emitter1)
	if ps.EmitterCount() != 1 {
		t.Errorf("expected 1 emitter, got %d", ps.EmitterCount())
	}

	ps.AddEmitter(emitter2)
	if ps.EmitterCount() != 2 {
		t.Errorf("expected 2 emitters, got %d", ps.EmitterCount())
	}

	// Adding nil should not add anything
	ps.AddEmitter(nil)
	if ps.EmitterCount() != 2 {
		t.Error("adding nil should not add an emitter")
	}

	ps.RemoveEmitter(emitter1)
	if ps.EmitterCount() != 1 {
		t.Errorf("expected 1 emitter after removal, got %d", ps.EmitterCount())
	}

	ps.ClearEmitters()
	if ps.EmitterCount() != 0 {
		t.Error("expected 0 emitters after clear")
	}
}

func TestCreateMovementTrailEmitter(t *testing.T) {
	ps := NewParticleSystem(12345)

	emitter := ps.CreateMovementTrailEmitter(5.0, 10.0)
	if emitter.X != 5.0 || emitter.Y != 10.0 {
		t.Error("emitter position incorrect")
	}
	if !emitter.Active {
		t.Error("emitter should be active by default")
	}
	if emitter.Rate <= 0 {
		t.Error("emitter should have positive rate")
	}

	// Genre should affect movement trail type
	ps.SetGenre(engine.GenreFantasy)
	fantasyEmitter := ps.CreateMovementTrailEmitter(0, 0)
	if fantasyEmitter.Type != ParticleTypeDust {
		t.Error("fantasy should use dust for movement trails")
	}

	ps.SetGenre(engine.GenreScifi)
	scifiEmitter := ps.CreateMovementTrailEmitter(0, 0)
	if scifiEmitter.Type != ParticleTypeThruster {
		t.Error("scifi should use thruster for movement trails")
	}
}

func TestCreateWeatherEmitter(t *testing.T) {
	ps := NewParticleSystem(12345)

	weatherTypes := []ParticleType{
		ParticleTypeRain,
		ParticleTypeSnow,
		ParticleTypeSand,
		ParticleTypeAsh,
	}

	for _, wType := range weatherTypes {
		emitter := ps.CreateWeatherEmitter(wType, 0, 0, 100)
		if emitter.Type != wType {
			t.Errorf("expected type %d, got %d", wType, emitter.Type)
		}
		if !emitter.Active {
			t.Error("weather emitter should be active")
		}
		if emitter.Rate <= 0 {
			t.Error("weather emitter should have positive rate")
		}
	}
}

func TestParticleEmission(t *testing.T) {
	ps := NewParticleSystem(12345)

	emitter := &ParticleEmitter{
		X:            0,
		Y:            0,
		Type:         ParticleTypeDust,
		Rate:         100.0, // High rate for testing
		SpreadAngle:  1.0,
		BaseVelocity: 10.0,
		VelocityVar:  1.0,
		BaseLife:     1.0,
		LifeVar:      0.1,
		BaseSize:     2.0,
		SizeVar:      0.5,
		Active:       true,
	}
	ps.AddEmitter(emitter)

	// Update should emit particles
	ps.Update(nil, 0.1)
	if ps.ParticleCount() == 0 {
		t.Error("emitter should have produced particles")
	}
}

func TestParticleLifecycle(t *testing.T) {
	ps := NewParticleSystem(12345)

	// Emit a burst of short-lived particles
	ps.EmitBurst(0, 0, ParticleTypeSparks, 10)
	initialCount := ps.ParticleCount()
	if initialCount < 10 {
		t.Errorf("expected at least 10 particles, got %d", initialCount)
	}

	// Update with large dt to kill particles
	ps.Update(nil, 2.0)
	if ps.ParticleCount() >= initialCount {
		t.Error("particles should have died")
	}
}

func TestEmitBurst(t *testing.T) {
	ps := NewParticleSystem(12345)

	burstTypes := []ParticleType{
		ParticleTypeSparks,
		ParticleTypeHeal,
		ParticleTypeExplosion,
	}

	for _, bType := range burstTypes {
		ps.ClearParticles()
		ps.EmitBurst(50, 50, bType, 20)
		if ps.ParticleCount() < 20 {
			t.Errorf("burst type %d should have produced 20 particles, got %d",
				bType, ps.ParticleCount())
		}

		// Verify particles have correct type
		for _, p := range ps.Particles() {
			if p.Type != bType {
				t.Errorf("particle should have type %d, got %d", bType, p.Type)
			}
		}
	}
}

func TestParticleUpdate(t *testing.T) {
	ps := NewParticleSystem(12345)

	ps.EmitBurst(0, 0, ParticleTypeDust, 5)
	particles := ps.Particles()

	// Record initial positions
	initialPositions := make([]struct{ x, y float64 }, len(particles))
	for i, p := range particles {
		initialPositions[i] = struct{ x, y float64 }{p.X, p.Y}
	}

	// Update
	ps.Update(nil, 0.1)

	// Particles should have moved
	moved := false
	for i, p := range ps.Particles() {
		if i < len(initialPositions) {
			if p.X != initialPositions[i].x || p.Y != initialPositions[i].y {
				moved = true
				break
			}
		}
	}
	if !moved && ps.ParticleCount() > 0 {
		t.Error("particles should have moved during update")
	}
}

func TestParticleGravity(t *testing.T) {
	ps := NewParticleSystem(12345)

	// Rain should fall down (gravity)
	ps.EmitBurst(0, 0, ParticleTypeRain, 5)
	particles := ps.Particles()
	initialVY := make([]float64, len(particles))
	for i, p := range particles {
		initialVY[i] = p.VY
	}

	ps.Update(nil, 0.1)

	// VY should have increased (downward)
	for i, p := range ps.Particles() {
		if i < len(initialVY) && p.VY <= initialVY[i] {
			t.Error("rain should accelerate downward due to gravity")
		}
	}
}

func TestParticleAlphaFade(t *testing.T) {
	ps := NewParticleSystem(12345)

	ps.EmitBurst(0, 0, ParticleTypeDust, 1)
	particles := ps.Particles()
	if len(particles) == 0 {
		t.Fatal("no particles emitted")
	}

	// Particle should have full alpha initially
	p := particles[0]
	if p.Alpha != 1.0 {
		t.Errorf("initial alpha should be 1.0, got %f", p.Alpha)
	}

	// Simulate particle near end of life
	p.Life = 0.1
	ps.Update(nil, 0.01)

	// Alpha should be reduced when life is low
	if p.Alpha >= 1.0 {
		t.Error("alpha should fade when particle life is low")
	}
}

func TestParticleTypeColors(t *testing.T) {
	ps := NewParticleSystem(12345)
	ps.SetGenre(engine.GenreFantasy)

	types := []ParticleType{
		ParticleTypeDust,
		ParticleTypeThruster,
		ParticleTypeRain,
		ParticleTypeSnow,
		ParticleTypeSparks,
		ParticleTypeHeal,
		ParticleTypeExplosion,
	}

	for _, pType := range types {
		color := ps.getColorForType(pType)
		if color.A == 0 {
			t.Errorf("particle type %d should have non-zero alpha", pType)
		}
	}
}

func TestParticleGenreColors(t *testing.T) {
	ps := NewParticleSystem(12345)

	// Get fantasy sparks color
	ps.SetGenre(engine.GenreFantasy)
	fantasySparks := ps.getColorForType(ParticleTypeSparks)

	// Get cyberpunk sparks color
	ps.SetGenre(engine.GenreCyberpunk)
	cyberpunkSparks := ps.getColorForType(ParticleTypeSparks)

	// Colors should be different
	if fantasySparks == cyberpunkSparks {
		t.Error("different genres should have different spark colors")
	}
}

func TestClearParticles(t *testing.T) {
	ps := NewParticleSystem(12345)

	ps.EmitBurst(0, 0, ParticleTypeDust, 50)
	if ps.ParticleCount() == 0 {
		t.Error("should have particles after burst")
	}

	ps.ClearParticles()
	if ps.ParticleCount() != 0 {
		t.Error("should have no particles after clear")
	}
}

func TestMaxParticleLimit(t *testing.T) {
	ps := NewParticleSystem(12345)

	// Try to emit more than max
	for i := 0; i < 300; i++ {
		ps.EmitBurst(0, 0, ParticleTypeDust, 10)
	}

	// Should be capped at maxParticles
	if ps.ParticleCount() > 2000 {
		t.Errorf("particle count should be capped at 2000, got %d", ps.ParticleCount())
	}
}

func TestParticlePresets(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		preset := defaultParticlePreset(genre)
		if preset == nil {
			t.Errorf("preset for %s should not be nil", genre)
			continue
		}

		// Verify all colors have non-zero alpha
		colors := []struct {
			name  string
			color color.RGBA
		}{
			{"DustColor", preset.DustColor},
			{"ThrusterColor", preset.ThrusterColor},
			{"RainColor", preset.RainColor},
			{"SnowColor", preset.SnowColor},
			{"SparksColor", preset.SparksColor},
			{"HealColor", preset.HealColor},
		}

		for _, c := range colors {
			if c.color.A == 0 {
				t.Errorf("%s %s should have non-zero alpha", genre, c.name)
			}
		}
	}
}
