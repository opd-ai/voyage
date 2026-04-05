package rendering

import (
	"image/color"
	"math"
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
	X, Y       float64        // World position
	Radius     float64        // Light radius in world units
	Intensity  float64        // Light intensity (0.0-1.0)
	Color      color.RGBA     // Light color
	LightType  PointLightType // Type of light source
	Flickering bool           // Whether the light flickers (e.g., campfires)
}

// LightingPreset defines genre-specific lighting colors.
type LightingPreset struct {
	DawnTint       color.RGBA // Warm sunrise color
	DayTint        color.RGBA // Full daylight (neutral or slight tint)
	DuskTint       color.RGBA // Sunset/evening color
	NightTint      color.RGBA // Night darkness color
	AmbientDay     float64    // Ambient light level during day (0.0-1.0)
	AmbientNight   float64    // Ambient light level at night (0.0-1.0)
	TransitionRate float64    // Speed of light transitions
	// Genre-specific point light colors
	TownLightColor     color.RGBA
	CampfireLightColor color.RGBA
	VesselLightColor   color.RGBA
	LanternLightColor  color.RGBA
}

// VisibilityRange represents the visible range at a position.
type VisibilityRange struct {
	BaseRange       float64 // Visibility range without lights (affected by darkness)
	EffectiveRange  float64 // Actual visibility range including light sources
	DarknessPenalty float64 // Penalty factor (0.0-1.0, where 1.0 = no penalty)
	LightBonus      float64 // Additional range from light sources
}

// BaseVisibilityRange is the maximum visibility in full daylight.
const BaseVisibilityRange = 10.0

// MinVisibilityRange is the minimum visibility in total darkness without lights.
const MinVisibilityRange = 2.0

// clampFloat restricts a value to a range.
func clampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// CalculateLightContribution calculates a single light's contribution at a position.
// flickerPhase is used for flickering effect animation.
func CalculateLightContribution(light *PointLight, x, y, flickerPhase float64) float64 {
	dx := x - light.X
	dy := y - light.Y
	distSq := dx*dx + dy*dy
	radiusSq := light.Radius * light.Radius

	if distSq > radiusSq {
		return 0.0
	}

	dist := math.Sqrt(distSq)
	falloff := 1.0 - (dist / light.Radius)
	falloff = falloff * falloff // Quadratic falloff for smoother light

	intensity := light.Intensity * falloff

	if light.Flickering {
		flicker := 0.9 + 0.1*math.Sin(flickerPhase*3.0+light.X+light.Y)
		intensity *= flicker
	}

	return intensity
}

// CalculateLightColorAt returns the combined light color at a position.
func CalculateLightColorAt(pointLights []*PointLight, ambient color.RGBA, ambientLevel, flickerPhase, x, y float64) color.RGBA {
	totalR := float64(ambient.R) * ambientLevel
	totalG := float64(ambient.G) * ambientLevel
	totalB := float64(ambient.B) * ambientLevel
	totalWeight := ambientLevel

	for _, light := range pointLights {
		intensity := CalculateLightContribution(light, x, y, flickerPhase)
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

// CalculateVisibilityAt calculates the visibility range at a given position.
func CalculateVisibilityAt(pointLights []*PointLight, ambientLevel, flickerPhase, x, y float64) VisibilityRange {
	darknessPenalty := ambientLevel

	lightBonus := 0.0
	for _, light := range pointLights {
		contribution := CalculateLightContribution(light, x, y, flickerPhase)
		if contribution > 0 {
			lightBonus += contribution * light.Radius * 0.5
		}
	}

	baseRange := MinVisibilityRange + (BaseVisibilityRange-MinVisibilityRange)*darknessPenalty
	effectiveRange := clampFloat(baseRange+lightBonus, MinVisibilityRange, BaseVisibilityRange*2)

	return VisibilityRange{
		BaseRange:       baseRange,
		EffectiveRange:  effectiveRange,
		DarknessPenalty: darknessPenalty,
		LightBonus:      lightBonus,
	}
}

// ApplyLightingToColor applies lighting tint to a color.
func ApplyLightingToColor(c color.Color, tint color.RGBA, ambient float64) color.RGBA {
	r, g, b, a := c.RGBA()

	outR := uint8(clampFloat(float64(r>>8)*ambient*float64(tint.R)/255.0, 0, 255))
	outG := uint8(clampFloat(float64(g>>8)*ambient*float64(tint.G)/255.0, 0, 255))
	outB := uint8(clampFloat(float64(b>>8)*ambient*float64(tint.B)/255.0, 0, 255))
	outA := uint8(a >> 8)

	return color.RGBA{outR, outG, outB, outA}
}
