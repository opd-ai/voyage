//go:build !headless

package rendering

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
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

// VesselSprite holds the sprite images for all damage states of a vessel.
type VesselSprite struct {
	PristineSprite *ebiten.Image
	WornSprite     *ebiten.Image
	DamagedSprite  *ebiten.Image
	CriticalSprite *ebiten.Image
}

// GetSprite returns the sprite for the given damage state.
func (vs *VesselSprite) GetSprite(state VesselDamageState) *ebiten.Image {
	switch state {
	case VesselWorn:
		return vs.WornSprite
	case VesselDamaged:
		return vs.DamagedSprite
	case VesselCritical:
		return vs.CriticalSprite
	default:
		return vs.PristineSprite
	}
}

// GetSpriteForRatio returns the sprite based on an integrity ratio.
func (vs *VesselSprite) GetSpriteForRatio(ratio float64) *ebiten.Image {
	return vs.GetSprite(VesselDamageStateFromRatio(ratio))
}

// VesselSpriteGenerator creates vessel sprites with damage state variations.
type VesselSpriteGenerator struct {
	gen        *seed.Generator
	spriteSize int
}

// NewVesselSpriteGenerator creates a new vessel sprite generator.
func NewVesselSpriteGenerator(masterSeed int64, spriteSize int) *VesselSpriteGenerator {
	return &VesselSpriteGenerator{
		gen:        seed.NewGenerator(masterSeed, "vessel-sprite"),
		spriteSize: spriteSize,
	}
}

// GenerateVesselSprite creates a vessel sprite with all damage state variants.
func (vsg *VesselSpriteGenerator) GenerateVesselSprite(hullColor, accentColor color.Color) *VesselSprite {
	pristine := vsg.generatePristineSprite(hullColor, accentColor)

	return &VesselSprite{
		PristineSprite: pristine,
		WornSprite:     vsg.generateWornSprite(pristine, hullColor),
		DamagedSprite:  vsg.generateDamagedSprite(pristine, hullColor),
		CriticalSprite: vsg.generateCriticalSprite(pristine, hullColor),
	}
}

// generatePristineSprite creates the base vessel sprite.
func (vsg *VesselSpriteGenerator) generatePristineSprite(hullColor, accentColor color.Color) *ebiten.Image {
	size := vsg.spriteSize
	img := ebiten.NewImage(size, size)

	// Generate a symmetric vessel shape
	halfWidth := size / 2

	// Hull body (diamond/boat shape)
	centerY := size / 2
	for y := 0; y < size; y++ {
		// Calculate width at this y position (diamond shape)
		distFromCenter := absInt(y - centerY)
		widthAtY := halfWidth - (distFromCenter * halfWidth / centerY)
		if widthAtY < 0 {
			widthAtY = 0
		}

		for x := halfWidth - widthAtY; x <= halfWidth+widthAtY; x++ {
			if x >= 0 && x < size {
				img.Set(x, y, hullColor)
			}
		}
	}

	// Add accent details (windows, markings)
	vsg.addAccentDetails(img, accentColor)

	return img
}

// absInt returns the absolute value of an integer.
func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// addAccentDetails adds decorative elements to the vessel sprite.
func (vsg *VesselSpriteGenerator) addAccentDetails(img *ebiten.Image, accentColor color.Color) {
	size := vsg.spriteSize
	centerX := size / 2

	// Add a stripe down the center
	stripeWidth := size / 8
	if stripeWidth < 1 {
		stripeWidth = 1
	}
	for y := size / 4; y < size*3/4; y++ {
		for x := centerX - stripeWidth; x <= centerX+stripeWidth; x++ {
			if x >= 0 && x < size {
				img.Set(x, y, accentColor)
			}
		}
	}

	// Add window lights
	windowY := size / 3
	windowSpacing := size / 4
	for i := -1; i <= 1; i++ {
		wx := centerX + i*windowSpacing
		if wx >= 0 && wx < size && vsg.gen.Chance(0.7) {
			vsg.drawPixel(img, wx, windowY, accentColor)
			vsg.drawPixel(img, wx, windowY+1, accentColor)
		}
	}
}

// drawPixel safely draws a single pixel.
func (vsg *VesselSpriteGenerator) drawPixel(img *ebiten.Image, x, y int, c color.Color) {
	if x >= 0 && x < vsg.spriteSize && y >= 0 && y < vsg.spriteSize {
		img.Set(x, y, c)
	}
}

// generateWornSprite creates a slightly weathered version of the vessel.
func (vsg *VesselSpriteGenerator) generateWornSprite(pristine *ebiten.Image, hullColor color.Color) *ebiten.Image {
	size := vsg.spriteSize
	img := ebiten.NewImage(size, size)
	img.DrawImage(pristine, nil)

	// Add minor scuff marks
	wornColor := vsg.darkenColor(hullColor, 0.8)
	scuffCount := size / 4

	for i := 0; i < scuffCount; i++ {
		x := vsg.gen.Intn(size)
		y := vsg.gen.Intn(size)
		if vsg.isOnSprite(pristine, x, y) {
			img.Set(x, y, wornColor)
		}
	}

	return img
}

// generateDamagedSprite creates a visibly damaged version of the vessel.
func (vsg *VesselSpriteGenerator) generateDamagedSprite(pristine *ebiten.Image, hullColor color.Color) *ebiten.Image {
	size := vsg.spriteSize
	img := ebiten.NewImage(size, size)
	img.DrawImage(pristine, nil)

	// Add damage marks and dents
	damageColor := vsg.darkenColor(hullColor, 0.6)
	charColor := color.RGBA{30, 30, 30, 255}
	damageCount := size / 2

	for i := 0; i < damageCount; i++ {
		x := vsg.gen.Intn(size)
		y := vsg.gen.Intn(size)
		if vsg.isOnSprite(pristine, x, y) {
			// Some pixels are darkened, some are charred
			if vsg.gen.Chance(0.3) {
				img.Set(x, y, charColor)
			} else {
				img.Set(x, y, damageColor)
			}
		}
	}

	// Add breach holes (small transparent areas)
	breachCount := 2 + vsg.gen.Intn(3)
	for i := 0; i < breachCount; i++ {
		bx := vsg.gen.Intn(size)
		by := vsg.gen.Intn(size)
		if vsg.isOnSprite(pristine, bx, by) {
			vsg.drawBreach(img, bx, by, 1+vsg.gen.Intn(2))
		}
	}

	return img
}

// generateCriticalSprite creates a heavily damaged version of the vessel.
func (vsg *VesselSpriteGenerator) generateCriticalSprite(pristine *ebiten.Image, hullColor color.Color) *ebiten.Image {
	size := vsg.spriteSize
	img := ebiten.NewImage(size, size)
	img.DrawImage(pristine, nil)

	// Heavy damage coloring
	damageColor := vsg.darkenColor(hullColor, 0.4)
	charColor := color.RGBA{20, 20, 20, 255}
	fireColor := color.RGBA{200, 80, 30, 255}
	damageCount := size

	for i := 0; i < damageCount; i++ {
		x := vsg.gen.Intn(size)
		y := vsg.gen.Intn(size)
		if vsg.isOnSprite(pristine, x, y) {
			roll := vsg.gen.Float64()
			if roll < 0.2 {
				img.Set(x, y, fireColor)
			} else if roll < 0.5 {
				img.Set(x, y, charColor)
			} else {
				img.Set(x, y, damageColor)
			}
		}
	}

	// Multiple breach holes
	breachCount := 4 + vsg.gen.Intn(4)
	for i := 0; i < breachCount; i++ {
		bx := vsg.gen.Intn(size)
		by := vsg.gen.Intn(size)
		if vsg.isOnSprite(pristine, bx, by) {
			vsg.drawBreach(img, bx, by, 2+vsg.gen.Intn(3))
		}
	}

	return img
}

// drawBreach draws a breach hole at the given position.
func (vsg *VesselSpriteGenerator) drawBreach(img *ebiten.Image, cx, cy, radius int) {
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			if dx*dx+dy*dy <= radius*radius {
				x, y := cx+dx, cy+dy
				if x >= 0 && x < vsg.spriteSize && y >= 0 && y < vsg.spriteSize {
					img.Set(x, y, color.Transparent)
				}
			}
		}
	}
}

// isOnSprite checks if a pixel is part of the visible sprite.
func (vsg *VesselSpriteGenerator) isOnSprite(img *ebiten.Image, x, y int) bool {
	if x < 0 || x >= vsg.spriteSize || y < 0 || y >= vsg.spriteSize {
		return false
	}
	_, _, _, a := img.At(x, y).RGBA()
	return a > 0
}

// darkenColor returns a darkened version of the color.
func (vsg *VesselSpriteGenerator) darkenColor(c color.Color, factor float64) color.Color {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(float64(r>>8) * factor),
		G: uint8(float64(g>>8) * factor),
		B: uint8(float64(b>>8) * factor),
		A: uint8(a >> 8),
	}
}
