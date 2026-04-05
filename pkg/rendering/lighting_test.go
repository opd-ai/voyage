package rendering

import (
	"image/color"
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewLightingSystem(t *testing.T) {
	ls := NewLightingSystem()
	if ls == nil {
		t.Fatal("NewLightingSystem returned nil")
	}
	if ls.phase != PhaseDay {
		t.Errorf("expected default phase PhaseDay, got %v", ls.phase)
	}
	if ls.genrePreset == nil {
		t.Error("genrePreset should be initialized")
	}
}

func TestLightingSystemSetGenre(t *testing.T) {
	ls := NewLightingSystem()

	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		ls.SetGenre(genre)
		if ls.Genre() != genre {
			t.Errorf("expected genre %s, got %s", genre, ls.Genre())
		}
		if ls.genrePreset == nil {
			t.Errorf("genrePreset should not be nil for genre %s", genre)
		}
	}
}

func TestLightingSystemSetPhase(t *testing.T) {
	ls := NewLightingSystem()

	tests := []struct {
		phase    TimePhase
		progress float64
		wantName string
	}{
		{PhaseDawn, 0.0, "Dawn"},
		{PhaseDay, 0.5, "Day"},
		{PhaseDusk, 1.0, "Dusk"},
		{PhaseNight, 0.25, "Night"},
	}

	for _, tc := range tests {
		ls.SetPhase(tc.phase, tc.progress)
		if ls.Phase() != tc.phase {
			t.Errorf("expected phase %v, got %v", tc.phase, ls.Phase())
		}
		if ls.PhaseName() != tc.wantName {
			t.Errorf("expected phase name %s, got %s", tc.wantName, ls.PhaseName())
		}
	}
}

func TestLightingSystemProgressClamping(t *testing.T) {
	ls := NewLightingSystem()

	// Test clamping below 0
	ls.SetPhase(PhaseDay, -0.5)
	if ls.Progress() != 0.0 {
		t.Errorf("progress should be clamped to 0.0, got %f", ls.Progress())
	}

	// Test clamping above 1
	ls.SetPhase(PhaseDay, 1.5)
	if ls.Progress() != 1.0 {
		t.Errorf("progress should be clamped to 1.0, got %f", ls.Progress())
	}

	// Test normal value
	ls.SetPhase(PhaseDay, 0.5)
	if ls.Progress() != 0.5 {
		t.Errorf("progress should be 0.5, got %f", ls.Progress())
	}
}

func TestLightingSystemCurrentTint(t *testing.T) {
	ls := NewLightingSystem()
	ls.SetGenre(engine.GenreFantasy)

	// Test dawn tint at start of dawn
	ls.SetPhase(PhaseDawn, 0.0)
	tint := ls.CurrentTint()
	if tint.A != 255 {
		t.Errorf("tint alpha should be 255, got %d", tint.A)
	}

	// Test that tint changes with progress
	ls.SetPhase(PhaseDawn, 0.0)
	tint1 := ls.CurrentTint()
	ls.SetPhase(PhaseDawn, 1.0)
	tint2 := ls.CurrentTint()

	if tint1 == tint2 {
		t.Error("tint should change between progress 0.0 and 1.0")
	}
}

func TestLightingSystemAmbientLevel(t *testing.T) {
	ls := NewLightingSystem()
	ls.SetGenre(engine.GenreFantasy)

	// Day should have high ambient
	ls.SetPhase(PhaseDay, 0.5)
	dayAmbient := ls.AmbientLevel()
	if dayAmbient < 0.5 {
		t.Errorf("day ambient should be high, got %f", dayAmbient)
	}

	// Night should have low ambient
	ls.SetPhase(PhaseNight, 0.5)
	nightAmbient := ls.AmbientLevel()
	if nightAmbient > 0.5 {
		t.Errorf("night ambient should be low, got %f", nightAmbient)
	}

	// Day should be brighter than night
	if dayAmbient <= nightAmbient {
		t.Error("day should be brighter than night")
	}
}

func TestLightingSystemApplyLighting(t *testing.T) {
	ls := NewLightingSystem()
	ls.SetGenre(engine.GenreFantasy)

	// Test that lighting is applied to a color
	originalColor := color.RGBA{200, 200, 200, 255}

	// During day, colors should be mostly preserved
	ls.SetPhase(PhaseDay, 0.5)
	dayColor := ls.ApplyLighting(originalColor)
	if dayColor.A != 255 {
		t.Errorf("alpha should be preserved, got %d", dayColor.A)
	}

	// During night, colors should be darker
	ls.SetPhase(PhaseNight, 0.5)
	nightColor := ls.ApplyLighting(originalColor)

	// Night should be darker than day
	dayBrightness := int(dayColor.R) + int(dayColor.G) + int(dayColor.B)
	nightBrightness := int(nightColor.R) + int(nightColor.G) + int(nightColor.B)
	if nightBrightness >= dayBrightness {
		t.Errorf("night should be darker than day: night=%d, day=%d",
			nightBrightness, dayBrightness)
	}
}

func TestLightingSystemGenrePresets(t *testing.T) {
	// Verify each genre has distinct lighting characteristics
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	presets := make(map[engine.GenreID]*LightingPreset)
	for _, genre := range genres {
		preset := defaultLightingPreset(genre)
		if preset == nil {
			t.Errorf("preset for %s should not be nil", genre)
			continue
		}
		presets[genre] = preset

		// Verify basic constraints
		if preset.AmbientDay <= preset.AmbientNight {
			t.Errorf("%s: day ambient (%f) should be greater than night (%f)",
				genre, preset.AmbientDay, preset.AmbientNight)
		}
		if preset.AmbientDay < 0.0 || preset.AmbientDay > 1.0 {
			t.Errorf("%s: day ambient should be 0-1, got %f", genre, preset.AmbientDay)
		}
		if preset.AmbientNight < 0.0 || preset.AmbientNight > 1.0 {
			t.Errorf("%s: night ambient should be 0-1, got %f", genre, preset.AmbientNight)
		}
	}

	// Verify horror is darker than fantasy
	if presets[engine.GenreHorror].AmbientNight >= presets[engine.GenreFantasy].AmbientNight {
		t.Error("horror should have darker nights than fantasy")
	}
}

func TestLerpColor(t *testing.T) {
	from := color.RGBA{0, 0, 0, 255}
	to := color.RGBA{100, 200, 50, 255}

	// Test t=0 returns from
	result := lerpColor(from, to, 0.0)
	if result != from {
		t.Errorf("lerpColor(0.0) should return from, got %v", result)
	}

	// Test t=1 returns to
	result = lerpColor(from, to, 1.0)
	if result != to {
		t.Errorf("lerpColor(1.0) should return to, got %v", result)
	}

	// Test t=0.5 is midpoint
	result = lerpColor(from, to, 0.5)
	if result.R != 50 || result.G != 100 || result.B != 25 {
		t.Errorf("lerpColor(0.5) should be midpoint, got R=%d G=%d B=%d",
			result.R, result.G, result.B)
	}
}

func TestLerpFloat(t *testing.T) {
	if lerpFloat(0.0, 10.0, 0.0) != 0.0 {
		t.Error("lerpFloat(0, 10, 0) should be 0")
	}
	if lerpFloat(0.0, 10.0, 1.0) != 10.0 {
		t.Error("lerpFloat(0, 10, 1) should be 10")
	}
	if lerpFloat(0.0, 10.0, 0.5) != 5.0 {
		t.Error("lerpFloat(0, 10, 0.5) should be 5")
	}
}

func TestClampFloat(t *testing.T) {
	tests := []struct {
		v, min, max, want float64
	}{
		{0.5, 0.0, 1.0, 0.5},
		{-0.5, 0.0, 1.0, 0.0},
		{1.5, 0.0, 1.0, 1.0},
		{0.0, 0.0, 1.0, 0.0},
		{1.0, 0.0, 1.0, 1.0},
	}

	for _, tc := range tests {
		got := clampFloat(tc.v, tc.min, tc.max)
		if got != tc.want {
			t.Errorf("clampFloat(%f, %f, %f) = %f, want %f",
				tc.v, tc.min, tc.max, got, tc.want)
		}
	}
}

func TestLightingSystemDawnTransition(t *testing.T) {
	ls := NewLightingSystem()
	ls.SetGenre(engine.GenreFantasy)

	// Dawn should transition from night ambient to day ambient
	ls.SetPhase(PhaseDawn, 0.0)
	startAmbient := ls.AmbientLevel()

	ls.SetPhase(PhaseDawn, 1.0)
	endAmbient := ls.AmbientLevel()

	if endAmbient <= startAmbient {
		t.Errorf("dawn should brighten: start=%f, end=%f", startAmbient, endAmbient)
	}
}

func TestLightingSystemDuskTransition(t *testing.T) {
	ls := NewLightingSystem()
	ls.SetGenre(engine.GenreFantasy)

	// Dusk should transition from day ambient to night ambient
	ls.SetPhase(PhaseDusk, 0.0)
	startAmbient := ls.AmbientLevel()

	ls.SetPhase(PhaseDusk, 1.0)
	endAmbient := ls.AmbientLevel()

	if endAmbient >= startAmbient {
		t.Errorf("dusk should darken: start=%f, end=%f", startAmbient, endAmbient)
	}
}
