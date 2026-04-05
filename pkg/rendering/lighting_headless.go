//go:build headless

package rendering

import (
	"image/color"
	"math"

	"github.com/opd-ai/voyage/pkg/engine"
)

// TimePhase represents a phase of the day/night cycle.
type TimePhase int

const (
	// PhaseDawn is the early morning transition from night to day.
	PhaseDawn TimePhase = iota
	// PhaseDay is full daylight.
	PhaseDay
	// PhaseDusk is the evening transition from day to night.
	PhaseDusk
	// PhaseNight is full darkness.
	PhaseNight
)

// PointLightType categorizes different point light sources.
type PointLightType int

const (
	// LightTypeTown represents lights from towns and settlements.
	LightTypeTown PointLightType = iota
	// LightTypeCampfire represents temporary camp lights.
	LightTypeCampfire
	// LightTypeVessel represents the player's vessel lantern/lights.
	LightTypeVessel
	// LightTypeLantern represents individual lanterns or torches.
	LightTypeLantern
)

// PointLight represents a light source at a specific position.
type PointLight struct {
	X, Y       float64
	Radius     float64
	Intensity  float64
	Color      color.RGBA
	LightType  PointLightType
	Flickering bool
}

// LightingSystem handles day/night cycle lighting on the overworld.
type LightingSystem struct {
	engine.BaseSystem
	phase        TimePhase
	progress     float64
	genrePreset  *LightingPreset
	pointLights  []*PointLight
	flickerPhase float64
}

// LightingPreset defines genre-specific lighting colors.
type LightingPreset struct {
	DawnTint           color.RGBA
	DayTint            color.RGBA
	DuskTint           color.RGBA
	NightTint          color.RGBA
	AmbientDay         float64
	AmbientNight       float64
	TransitionRate     float64
	TownLightColor     color.RGBA
	CampfireLightColor color.RGBA
	VesselLightColor   color.RGBA
	LanternLightColor  color.RGBA
}

// NewLightingSystem creates a new lighting system with default settings.
func NewLightingSystem() *LightingSystem {
	ls := &LightingSystem{
		BaseSystem:  engine.NewBaseSystem(engine.PriorityRender - 1),
		phase:       PhaseDay,
		progress:    0.5,
		genrePreset: defaultLightingPreset(engine.GenreFantasy),
	}
	return ls
}

// SetGenre changes the lighting presets to match the genre.
func (ls *LightingSystem) SetGenre(genreID engine.GenreID) {
	ls.BaseSystem.SetGenre(genreID)
	ls.genrePreset = defaultLightingPreset(genreID)
}

// defaultLightingPreset returns genre-specific lighting presets.
func defaultLightingPreset(genre engine.GenreID) *LightingPreset {
	switch genre {
	case engine.GenreScifi:
		return &LightingPreset{
			DawnTint:           color.RGBA{200, 220, 255, 255},
			DayTint:            color.RGBA{220, 230, 255, 255},
			DuskTint:           color.RGBA{150, 180, 220, 255},
			NightTint:          color.RGBA{40, 60, 100, 255},
			AmbientDay:         0.95,
			AmbientNight:       0.3,
			TransitionRate:     1.0,
			TownLightColor:     color.RGBA{100, 180, 255, 255},
			CampfireLightColor: color.RGBA{200, 220, 255, 255},
			VesselLightColor:   color.RGBA{150, 200, 255, 255},
			LanternLightColor:  color.RGBA{180, 200, 255, 255},
		}
	case engine.GenreHorror:
		return &LightingPreset{
			DawnTint:           color.RGBA{180, 150, 140, 255},
			DayTint:            color.RGBA{200, 180, 170, 255},
			DuskTint:           color.RGBA{150, 100, 90, 255},
			NightTint:          color.RGBA{30, 25, 35, 255},
			AmbientDay:         0.75,
			AmbientNight:       0.15,
			TransitionRate:     0.8,
			TownLightColor:     color.RGBA{200, 100, 80, 255},
			CampfireLightColor: color.RGBA{255, 150, 100, 255},
			VesselLightColor:   color.RGBA{180, 80, 60, 255},
			LanternLightColor:  color.RGBA{200, 130, 80, 255},
		}
	case engine.GenreCyberpunk:
		return &LightingPreset{
			DawnTint:           color.RGBA{255, 200, 255, 255},
			DayTint:            color.RGBA{220, 220, 240, 255},
			DuskTint:           color.RGBA{255, 100, 200, 255},
			NightTint:          color.RGBA{60, 40, 100, 255},
			AmbientDay:         0.85,
			AmbientNight:       0.4,
			TransitionRate:     1.2,
			TownLightColor:     color.RGBA{255, 0, 200, 255},
			CampfireLightColor: color.RGBA{0, 255, 200, 255},
			VesselLightColor:   color.RGBA{200, 100, 255, 255},
			LanternLightColor:  color.RGBA{255, 255, 0, 255},
		}
	case engine.GenrePostapoc:
		return &LightingPreset{
			DawnTint:           color.RGBA{220, 180, 140, 255},
			DayTint:            color.RGBA{240, 220, 180, 255},
			DuskTint:           color.RGBA{200, 120, 80, 255},
			NightTint:          color.RGBA{40, 35, 30, 255},
			AmbientDay:         0.85,
			AmbientNight:       0.2,
			TransitionRate:     0.9,
			TownLightColor:     color.RGBA{200, 150, 80, 255},
			CampfireLightColor: color.RGBA{255, 180, 100, 255},
			VesselLightColor:   color.RGBA{180, 140, 80, 255},
			LanternLightColor:  color.RGBA{200, 160, 100, 255},
		}
	default:
		return &LightingPreset{
			DawnTint:           color.RGBA{255, 220, 180, 255},
			DayTint:            color.RGBA{255, 250, 240, 255},
			DuskTint:           color.RGBA{255, 180, 130, 255},
			NightTint:          color.RGBA{40, 50, 80, 255},
			AmbientDay:         1.0,
			AmbientNight:       0.25,
			TransitionRate:     1.0,
			TownLightColor:     color.RGBA{255, 220, 150, 255},
			CampfireLightColor: color.RGBA{255, 180, 100, 255},
			VesselLightColor:   color.RGBA{255, 200, 120, 255},
			LanternLightColor:  color.RGBA{255, 210, 130, 255},
		}
	}
}

// SetPhase sets the current time phase and progress within that phase.
func (ls *LightingSystem) SetPhase(phase TimePhase, progress float64) {
	ls.phase = phase
	ls.progress = clampFloat(progress, 0.0, 1.0)
}

// Phase returns the current time phase.
func (ls *LightingSystem) Phase() TimePhase {
	return ls.phase
}

// Progress returns the current progress within the phase (0.0-1.0).
func (ls *LightingSystem) Progress() float64 {
	return ls.progress
}

// PhaseName returns a human-readable name for the current phase.
func (ls *LightingSystem) PhaseName() string {
	switch ls.phase {
	case PhaseDawn:
		return "Dawn"
	case PhaseDay:
		return "Day"
	case PhaseDusk:
		return "Dusk"
	case PhaseNight:
		return "Night"
	default:
		return "Unknown"
	}
}

// CurrentTint returns the current lighting tint color.
func (ls *LightingSystem) CurrentTint() color.RGBA {
	var from, to color.RGBA
	switch ls.phase {
	case PhaseDawn:
		from = ls.genrePreset.NightTint
		to = ls.genrePreset.DawnTint
	case PhaseDay:
		from = ls.genrePreset.DawnTint
		to = ls.genrePreset.DayTint
	case PhaseDusk:
		from = ls.genrePreset.DayTint
		to = ls.genrePreset.DuskTint
	case PhaseNight:
		from = ls.genrePreset.DuskTint
		to = ls.genrePreset.NightTint
	}
	return lerpColor(from, to, ls.progress)
}

// AmbientLevel returns the current ambient light level (0.0-1.0).
func (ls *LightingSystem) AmbientLevel() float64 {
	dayLevel := ls.genrePreset.AmbientDay
	nightLevel := ls.genrePreset.AmbientNight
	switch ls.phase {
	case PhaseDawn:
		return lerpFloat(nightLevel, dayLevel, ls.progress)
	case PhaseDay:
		return dayLevel
	case PhaseDusk:
		return lerpFloat(dayLevel, nightLevel, ls.progress)
	case PhaseNight:
		return nightLevel
	}
	return dayLevel
}

// Update implements the System interface.
func (ls *LightingSystem) Update(world *engine.World, dt float64) {
	ls.flickerPhase += dt * 10.0
	if ls.flickerPhase > math.Pi*2 {
		ls.flickerPhase -= math.Pi * 2
	}
}

// AddPointLight adds a point light to the scene.
func (ls *LightingSystem) AddPointLight(light *PointLight) {
	if light == nil {
		return
	}
	if light.Color.A == 0 {
		light.Color = ls.GetLightColorForType(light.LightType)
	}
	ls.pointLights = append(ls.pointLights, light)
}

// RemovePointLight removes a point light from the scene.
func (ls *LightingSystem) RemovePointLight(light *PointLight) {
	for i, l := range ls.pointLights {
		if l == light {
			ls.pointLights = append(ls.pointLights[:i], ls.pointLights[i+1:]...)
			return
		}
	}
}

// ClearPointLights removes all point lights from the scene.
func (ls *LightingSystem) ClearPointLights() {
	ls.pointLights = nil
}

// PointLights returns all current point lights.
func (ls *LightingSystem) PointLights() []*PointLight {
	return ls.pointLights
}

// GetLightColorForType returns the genre-appropriate color for a light type.
func (ls *LightingSystem) GetLightColorForType(lightType PointLightType) color.RGBA {
	switch lightType {
	case LightTypeTown:
		return ls.genrePreset.TownLightColor
	case LightTypeCampfire:
		return ls.genrePreset.CampfireLightColor
	case LightTypeVessel:
		return ls.genrePreset.VesselLightColor
	case LightTypeLantern:
		return ls.genrePreset.LanternLightColor
	default:
		return ls.genrePreset.LanternLightColor
	}
}

// CreateTownLight creates a point light for a town/settlement.
func (ls *LightingSystem) CreateTownLight(x, y, radius float64) *PointLight {
	return &PointLight{
		X:          x,
		Y:          y,
		Radius:     radius,
		Intensity:  0.8,
		Color:      ls.genrePreset.TownLightColor,
		LightType:  LightTypeTown,
		Flickering: false,
	}
}

// CreateCampfireLight creates a point light for a campfire.
func (ls *LightingSystem) CreateCampfireLight(x, y float64) *PointLight {
	return &PointLight{
		X:          x,
		Y:          y,
		Radius:     3.0,
		Intensity:  0.9,
		Color:      ls.genrePreset.CampfireLightColor,
		LightType:  LightTypeCampfire,
		Flickering: true,
	}
}

// CreateVesselLight creates a point light for the player's vessel.
func (ls *LightingSystem) CreateVesselLight(x, y float64) *PointLight {
	return &PointLight{
		X:          x,
		Y:          y,
		Radius:     4.0,
		Intensity:  0.85,
		Color:      ls.genrePreset.VesselLightColor,
		LightType:  LightTypeVessel,
		Flickering: false,
	}
}

// CreateLanternLight creates a point light for a lantern or torch.
func (ls *LightingSystem) CreateLanternLight(x, y float64) *PointLight {
	return &PointLight{
		X:          x,
		Y:          y,
		Radius:     2.5,
		Intensity:  0.7,
		Color:      ls.genrePreset.LanternLightColor,
		LightType:  LightTypeLantern,
		Flickering: true,
	}
}

// GetLightIntensityAt calculates the combined light intensity at a position.
func (ls *LightingSystem) GetLightIntensityAt(x, y float64) float64 {
	totalIntensity := ls.AmbientLevel()
	for _, light := range ls.pointLights {
		intensity := ls.calculateLightContribution(light, x, y)
		totalIntensity += intensity
	}
	return clampFloat(totalIntensity, 0.0, 1.0)
}

// calculateLightContribution calculates a single light's contribution at a position.
func (ls *LightingSystem) calculateLightContribution(light *PointLight, x, y float64) float64 {
	dx := x - light.X
	dy := y - light.Y
	distSq := dx*dx + dy*dy
	radiusSq := light.Radius * light.Radius

	if distSq > radiusSq {
		return 0.0
	}

	dist := math.Sqrt(distSq)
	falloff := 1.0 - (dist / light.Radius)
	falloff = falloff * falloff

	intensity := light.Intensity * falloff

	if light.Flickering {
		flicker := 0.9 + 0.1*math.Sin(ls.flickerPhase*3.0+light.X+light.Y)
		intensity *= flicker
	}

	return intensity
}

// GetLightColorAt returns the combined light color at a position.
func (ls *LightingSystem) GetLightColorAt(x, y float64) color.RGBA {
	ambient := ls.CurrentTint()
	totalR := float64(ambient.R) * ls.AmbientLevel()
	totalG := float64(ambient.G) * ls.AmbientLevel()
	totalB := float64(ambient.B) * ls.AmbientLevel()
	totalWeight := ls.AmbientLevel()

	for _, light := range ls.pointLights {
		intensity := ls.calculateLightContribution(light, x, y)
		if intensity > 0 {
			totalR += float64(light.Color.R) * intensity
			totalG += float64(light.Color.G) * intensity
			totalB += float64(light.Color.B) * intensity
			totalWeight += intensity
		}
	}

	if totalWeight > 0 {
		return color.RGBA{
			R: uint8(clampFloat(totalR/totalWeight, 0, 255)),
			G: uint8(clampFloat(totalG/totalWeight, 0, 255)),
			B: uint8(clampFloat(totalB/totalWeight, 0, 255)),
			A: 255,
		}
	}

	return ambient
}

// ApplyLighting applies the current lighting tint to a color.
func (ls *LightingSystem) ApplyLighting(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	tint := ls.CurrentTint()
	ambient := ls.AmbientLevel()

	outR := uint8(clampFloat(float64(r>>8)*ambient*float64(tint.R)/255.0, 0, 255))
	outG := uint8(clampFloat(float64(g>>8)*ambient*float64(tint.G)/255.0, 0, 255))
	outB := uint8(clampFloat(float64(b>>8)*ambient*float64(tint.B)/255.0, 0, 255))
	outA := uint8(a >> 8)

	return color.RGBA{outR, outG, outB, outA}
}

// lerpColor linearly interpolates between two colors.
func lerpColor(from, to color.RGBA, t float64) color.RGBA {
	t = clampFloat(t, 0.0, 1.0)
	return color.RGBA{
		R: uint8(float64(from.R) + t*(float64(to.R)-float64(from.R))),
		G: uint8(float64(from.G) + t*(float64(to.G)-float64(from.G))),
		B: uint8(float64(from.B) + t*(float64(to.B)-float64(from.B))),
		A: uint8(float64(from.A) + t*(float64(to.A)-float64(from.A))),
	}
}

// lerpFloat linearly interpolates between two float values.
func lerpFloat(from, to, t float64) float64 {
	return from + t*(to-from)
}

// clampFloat restricts a float value to a range.
func clampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
