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

func TestClampFloatLighting(t *testing.T) {
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
		got := clampFloatLighting(tc.v, tc.min, tc.max)
		if got != tc.want {
			t.Errorf("clampFloatLighting(%f, %f, %f) = %f, want %f",
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

func TestPointLightCreation(t *testing.T) {
	ls := NewLightingSystem()
	ls.SetGenre(engine.GenreFantasy)

	// Test creating different light types
	townLight := ls.CreateTownLight(10.0, 20.0, 5.0)
	if townLight.X != 10.0 || townLight.Y != 20.0 {
		t.Error("town light position incorrect")
	}
	if townLight.Radius != 5.0 {
		t.Errorf("town light radius should be 5.0, got %f", townLight.Radius)
	}
	if townLight.LightType != LightTypeTown {
		t.Error("town light type incorrect")
	}
	if townLight.Flickering {
		t.Error("town light should not flicker")
	}

	campfire := ls.CreateCampfireLight(5.0, 5.0)
	if campfire.LightType != LightTypeCampfire {
		t.Error("campfire light type incorrect")
	}
	if !campfire.Flickering {
		t.Error("campfire should flicker")
	}

	vessel := ls.CreateVesselLight(0.0, 0.0)
	if vessel.LightType != LightTypeVessel {
		t.Error("vessel light type incorrect")
	}

	lantern := ls.CreateLanternLight(1.0, 1.0)
	if lantern.LightType != LightTypeLantern {
		t.Error("lantern light type incorrect")
	}
	if !lantern.Flickering {
		t.Error("lantern should flicker")
	}
}

func TestPointLightManagement(t *testing.T) {
	ls := NewLightingSystem()

	if len(ls.PointLights()) != 0 {
		t.Error("should start with no point lights")
	}

	// Add lights
	light1 := ls.CreateTownLight(0, 0, 5)
	light2 := ls.CreateCampfireLight(10, 10)
	ls.AddPointLight(light1)
	ls.AddPointLight(light2)

	if len(ls.PointLights()) != 2 {
		t.Errorf("expected 2 lights, got %d", len(ls.PointLights()))
	}

	// Remove one light
	ls.RemovePointLight(light1)
	if len(ls.PointLights()) != 1 {
		t.Errorf("expected 1 light after removal, got %d", len(ls.PointLights()))
	}

	// Clear all lights
	ls.ClearPointLights()
	if len(ls.PointLights()) != 0 {
		t.Errorf("expected 0 lights after clear, got %d", len(ls.PointLights()))
	}

	// Adding nil should not panic
	ls.AddPointLight(nil)
	if len(ls.PointLights()) != 0 {
		t.Error("adding nil should not add a light")
	}
}

func TestGetLightColorForType(t *testing.T) {
	ls := NewLightingSystem()
	ls.SetGenre(engine.GenreFantasy)

	// Each light type should return a valid color
	types := []PointLightType{LightTypeTown, LightTypeCampfire, LightTypeVessel, LightTypeLantern}
	for _, lt := range types {
		col := ls.GetLightColorForType(lt)
		if col.A != 255 {
			t.Errorf("light type %d should have full alpha", lt)
		}
	}

	// Change genre and verify colors change
	ls.SetGenre(engine.GenreScifi)
	scifiColor := ls.GetLightColorForType(LightTypeTown)

	ls.SetGenre(engine.GenreFantasy)
	fantasyColor := ls.GetLightColorForType(LightTypeTown)

	if scifiColor == fantasyColor {
		t.Error("different genres should have different light colors")
	}
}

func TestGetLightIntensityAt(t *testing.T) {
	ls := NewLightingSystem()
	ls.SetPhase(PhaseNight, 0.5) // Low ambient

	// Without lights, should return ambient level
	intensity := ls.GetLightIntensityAt(0, 0)
	if intensity != ls.AmbientLevel() {
		t.Errorf("without lights, intensity should equal ambient: got %f, want %f",
			intensity, ls.AmbientLevel())
	}

	// Add a light at origin
	light := &PointLight{
		X:         0,
		Y:         0,
		Radius:    5.0,
		Intensity: 1.0,
		Color:     color.RGBA{255, 200, 100, 255},
	}
	ls.AddPointLight(light)

	// At the light position, intensity should be higher
	intensityAtLight := ls.GetLightIntensityAt(0, 0)
	if intensityAtLight <= ls.AmbientLevel() {
		t.Error("intensity at light should be higher than ambient")
	}

	// Far from light, intensity should be ambient
	intensityFar := ls.GetLightIntensityAt(100, 100)
	if intensityFar != ls.AmbientLevel() {
		t.Errorf("intensity far from light should be ambient: got %f, want %f",
			intensityFar, ls.AmbientLevel())
	}
}

func TestGetLightColorAt(t *testing.T) {
	ls := NewLightingSystem()
	ls.SetPhase(PhaseNight, 0.5)

	// Without lights, should return ambient tint
	col := ls.GetLightColorAt(0, 0)
	if col.A != 255 {
		t.Error("color should have full alpha")
	}

	// Add a colored light
	light := &PointLight{
		X:         0,
		Y:         0,
		Radius:    5.0,
		Intensity: 1.0,
		Color:     color.RGBA{255, 0, 0, 255}, // Bright red
	}
	ls.AddPointLight(light)

	// At the light, color should be influenced by the red light
	colAtLight := ls.GetLightColorAt(0, 0)
	ambient := ls.CurrentTint()
	// The red channel should be higher than the ambient
	if colAtLight.R <= ambient.R {
		t.Error("red channel at light should be influenced by red light")
	}
}

func TestPointLightFalloff(t *testing.T) {
	ls := NewLightingSystem()
	ls.SetPhase(PhaseNight, 0.5) // Low ambient
	ls.ClearPointLights()

	light := &PointLight{
		X:         0,
		Y:         0,
		Radius:    10.0,
		Intensity: 1.0,
		Color:     color.RGBA{255, 255, 255, 255},
	}
	ls.AddPointLight(light)

	// Intensity should decrease with distance
	intensityClose := ls.GetLightIntensityAt(1, 0)
	intensityMid := ls.GetLightIntensityAt(5, 0)
	intensityEdge := ls.GetLightIntensityAt(9, 0)

	if intensityClose <= intensityMid {
		t.Error("intensity should be higher closer to light")
	}
	if intensityMid <= intensityEdge {
		t.Error("intensity should decrease toward edge")
	}
}

func TestPointLightGenreColors(t *testing.T) {
	// Test that each genre has distinct light colors
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	townColors := make(map[engine.GenreID]color.RGBA)
	for _, genre := range genres {
		ls := NewLightingSystem()
		ls.SetGenre(genre)
		townColors[genre] = ls.GetLightColorForType(LightTypeTown)
	}

	// Horror should have darker/redder town lights than fantasy
	if townColors[engine.GenreHorror].R < townColors[engine.GenreHorror].G &&
		townColors[engine.GenreHorror].R < townColors[engine.GenreHorror].B {
		t.Error("horror town lights should be reddish")
	}

	// Cyberpunk should have vibrant neon colors
	cyberpunk := townColors[engine.GenreCyberpunk]
	if cyberpunk.R < 200 && cyberpunk.G < 200 && cyberpunk.B < 200 {
		t.Error("cyberpunk should have at least one high color channel (neon)")
	}
}

func TestVisibilityAtDayNight(t *testing.T) {
	ls := NewLightingSystem()

	// Day should have good visibility
	ls.SetPhase(PhaseDay, 0.5)
	dayVisibility := ls.GetVisibilityAt(0, 0)
	if dayVisibility.DarknessPenalty < 0.8 {
		t.Errorf("day should have high darkness penalty (low penalty): got %f",
			dayVisibility.DarknessPenalty)
	}
	if dayVisibility.BaseRange < 8.0 {
		t.Errorf("day base range should be high: got %f", dayVisibility.BaseRange)
	}

	// Night should have reduced visibility
	ls.SetPhase(PhaseNight, 0.5)
	nightVisibility := ls.GetVisibilityAt(0, 0)
	if nightVisibility.DarknessPenalty > 0.5 {
		t.Errorf("night should have significant darkness penalty: got %f",
			nightVisibility.DarknessPenalty)
	}
	if nightVisibility.BaseRange > 5.0 {
		t.Errorf("night base range should be reduced: got %f", nightVisibility.BaseRange)
	}

	// Day should have better visibility than night
	if dayVisibility.EffectiveRange <= nightVisibility.EffectiveRange {
		t.Error("day should have better visibility than night")
	}
}

func TestVisibilityWithLightSource(t *testing.T) {
	ls := NewLightingSystem()
	ls.SetPhase(PhaseNight, 0.5)

	// Without light
	visibilityNoLight := ls.GetVisibilityAt(0, 0)

	// Add a light
	light := ls.CreateLanternLight(0, 0)
	ls.AddPointLight(light)

	// With light, visibility should be better
	visibilityWithLight := ls.GetVisibilityAt(0, 0)
	if visibilityWithLight.LightBonus <= 0 {
		t.Error("light should provide a visibility bonus")
	}
	if visibilityWithLight.EffectiveRange <= visibilityNoLight.EffectiveRange {
		t.Error("light should improve effective range")
	}
}

func TestIsVisibleAt(t *testing.T) {
	ls := NewLightingSystem()
	ls.SetPhase(PhaseNight, 0.5)

	// Get visibility range
	visibility := ls.GetVisibilityAt(0, 0)

	// Target within range should be visible
	withinRange := visibility.EffectiveRange * 0.5
	if !ls.IsVisibleAt(0, 0, withinRange, 0) {
		t.Error("target within range should be visible")
	}

	// Target beyond range should not be visible
	beyondRange := visibility.EffectiveRange * 2
	if ls.IsVisibleAt(0, 0, beyondRange, 0) {
		t.Error("target beyond range should not be visible")
	}
}

func TestVisibilityPenaltyDescription(t *testing.T) {
	ls := NewLightingSystem()

	// Day should have clear visibility
	ls.SetPhase(PhaseDay, 0.5)
	dayDesc := ls.GetVisibilityPenaltyDescription(0, 0)
	if dayDesc != "Clear visibility" {
		t.Errorf("day description should be 'Clear visibility', got %s", dayDesc)
	}

	// Night should show poor visibility
	ls.SetPhase(PhaseNight, 0.5)
	ls.ClearPointLights()
	nightDesc := ls.GetVisibilityPenaltyDescription(0, 0)
	if nightDesc == "Clear visibility" {
		t.Error("night should not have clear visibility")
	}

	// Dusk should show reduced visibility
	ls.SetPhase(PhaseDusk, 0.5)
	duskDesc := ls.GetVisibilityPenaltyDescription(0, 0)
	if duskDesc == "Clear visibility" {
		t.Error("dusk should have reduced visibility")
	}
}

func TestHasAdequateLighting(t *testing.T) {
	ls := NewLightingSystem()

	// Day should have adequate lighting
	ls.SetPhase(PhaseDay, 0.5)
	if !ls.HasAdequateLighting(0, 0) {
		t.Error("day should have adequate lighting")
	}

	// Deep night without lights should not have adequate lighting
	ls.SetPhase(PhaseNight, 0.5)
	ls.SetGenre(engine.GenreHorror) // Horror has dark nights
	ls.ClearPointLights()
	if ls.HasAdequateLighting(0, 0) {
		t.Error("dark night without lights should not have adequate lighting")
	}

	// Adding a bright light should provide adequate lighting
	light := &PointLight{
		X:         0,
		Y:         0,
		Radius:    5.0,
		Intensity: 1.0,
		Color:     color.RGBA{255, 255, 255, 255},
	}
	ls.AddPointLight(light)
	if !ls.HasAdequateLighting(0, 0) {
		t.Error("light source should provide adequate lighting")
	}
}

func TestNeedsLightSource(t *testing.T) {
	ls := NewLightingSystem()

	// Day should not need light source
	ls.SetPhase(PhaseDay, 0.5)
	if ls.NeedsLightSource(0, 0) {
		t.Error("day should not need a light source")
	}

	// Night should need light source
	ls.SetPhase(PhaseNight, 0.5)
	ls.SetGenre(engine.GenreHorror) // Dark nights
	ls.ClearPointLights()
	if !ls.NeedsLightSource(0, 0) {
		t.Error("dark night should need a light source")
	}

	// With a nearby light, should not need additional light
	light := ls.CreateVesselLight(0, 0)
	ls.AddPointLight(light)
	if ls.NeedsLightSource(0, 0) {
		t.Error("with vessel light, should not need additional light source")
	}
}

func TestVisibilityRangeConstants(t *testing.T) {
	// Verify constant relationships
	if MinVisibilityRange >= BaseVisibilityRange {
		t.Error("min visibility should be less than base visibility")
	}
	if MinVisibilityRange < 1.0 {
		t.Error("min visibility should be at least 1.0")
	}
	if BaseVisibilityRange < 5.0 {
		t.Error("base visibility should be at least 5.0 for playability")
	}
}
