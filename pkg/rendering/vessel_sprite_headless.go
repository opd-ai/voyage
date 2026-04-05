//go:build headless

package rendering

import (
	"image/color"
)

// VesselDamageState represents the visual damage state of a vessel.
type VesselDamageState int

const (
	VesselPristine VesselDamageState = iota
	VesselWorn
	VesselDamaged
	VesselCritical
)

// VesselDamageStateFromRatio determines the damage state from an integrity ratio.
func VesselDamageStateFromRatio(ratio float64) VesselDamageState {
	switch {
	case ratio >= 0.90:
		return VesselPristine
	case ratio >= 0.50:
		return VesselWorn
	case ratio >= 0.25:
		return VesselDamaged
	default:
		return VesselCritical
	}
}

// VesselSprite holds metadata for all damage states of a vessel (headless stub).
type VesselSprite struct {
	SpriteSize int
}

// GetDamageState returns the current damage state for a given ratio.
func (vs *VesselSprite) GetDamageState(ratio float64) VesselDamageState {
	return VesselDamageStateFromRatio(ratio)
}

// VesselSpriteGenerator creates vessel sprites with damage state variations (headless stub).
type VesselSpriteGenerator struct {
	spriteSize int
}

// NewVesselSpriteGenerator creates a new vessel sprite generator stub.
func NewVesselSpriteGenerator(masterSeed int64, spriteSize int) *VesselSpriteGenerator {
	return &VesselSpriteGenerator{
		spriteSize: spriteSize,
	}
}

// GenerateVesselSprite creates a vessel sprite stub with all damage state variants.
func (vsg *VesselSpriteGenerator) GenerateVesselSprite(hullColor, accentColor color.Color) *VesselSprite {
	return &VesselSprite{
		SpriteSize: vsg.spriteSize,
	}
}

// absInt returns the absolute value of an integer.
func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
