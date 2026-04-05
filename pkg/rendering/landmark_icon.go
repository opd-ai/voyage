//go:build !headless

package rendering

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// LandmarkIconType identifies the type of landmark icon.
type LandmarkIconType int

const (
	LandmarkIconTown LandmarkIconType = iota
	LandmarkIconOutpost
	LandmarkIconRuins
	LandmarkIconShrine
	LandmarkIconOrigin
	LandmarkIconDestination
)

// AnimatedLandmarkIcon represents an animated landmark icon.
type AnimatedLandmarkIcon struct {
	Frames       []*ebiten.Image
	FrameTime    float64
	currentFrame int
	elapsed      float64
}

// NewAnimatedLandmarkIcon creates a new animated landmark icon.
func NewAnimatedLandmarkIcon(frames []*ebiten.Image, frameTime float64) *AnimatedLandmarkIcon {
	return &AnimatedLandmarkIcon{
		Frames:       frames,
		FrameTime:    frameTime,
		currentFrame: 0,
		elapsed:      0,
	}
}

// generateAnimatedIcon creates an animated icon using provided draw functions.
func (lig *LandmarkIconGenerator) generateAnimatedIcon(
	frameCount int,
	frameTime float64,
	drawFrame func(img *ebiten.Image, frame int),
) *AnimatedLandmarkIcon {
	frames := make([]*ebiten.Image, frameCount)
	for f := 0; f < frameCount; f++ {
		img := ebiten.NewImage(lig.iconSize, lig.iconSize)
		drawFrame(img, f)
		frames[f] = img
	}
	return NewAnimatedLandmarkIcon(frames, frameTime)
}

// Update advances the landmark animation by the given delta time.
func (ali *AnimatedLandmarkIcon) Update(dt float64) {
	if len(ali.Frames) <= 1 {
		return
	}
	ali.elapsed += dt
	if ali.elapsed >= ali.FrameTime {
		ali.elapsed -= ali.FrameTime
		ali.currentFrame++
		if ali.currentFrame >= len(ali.Frames) {
			ali.currentFrame = 0
		}
	}
}

// CurrentFrame returns the current animation frame image.
func (ali *AnimatedLandmarkIcon) CurrentFrame() *ebiten.Image {
	if len(ali.Frames) == 0 {
		return nil
	}
	return ali.Frames[ali.currentFrame]
}

// Reset resets the animation to the first frame.
func (ali *AnimatedLandmarkIcon) Reset() {
	ali.currentFrame = 0
	ali.elapsed = 0
}

// LandmarkIconGenerator creates animated landmark icons.
type LandmarkIconGenerator struct {
	gen      *seed.Generator
	iconSize int
}

// drawRect draws a filled rectangle with bounds checking.
func (lig *LandmarkIconGenerator) drawRect(img *ebiten.Image, x1, y1, x2, y2 int, c color.Color) {
	size := lig.iconSize
	for y := y1; y < y2; y++ {
		for x := x1; x <= x2; x++ {
			if x >= 0 && x < size && y >= 0 && y < size {
				img.Set(x, y, c)
			}
		}
	}
}

// drawVerticalLine draws a vertical line with bounds checking.
func (lig *LandmarkIconGenerator) drawVerticalLine(img *ebiten.Image, x, y1, y2 int, c color.Color) {
	size := lig.iconSize
	for y := y1; y < y2; y++ {
		if y >= 0 && y < size {
			img.Set(x, y, c)
		}
	}
}

// NewLandmarkIconGenerator creates a new landmark icon generator.
func NewLandmarkIconGenerator(masterSeed int64, iconSize int) *LandmarkIconGenerator {
	return &LandmarkIconGenerator{
		gen:      seed.NewGenerator(masterSeed, "landmark-icon"),
		iconSize: iconSize,
	}
}

// GenerateAnimatedIcon creates an animated icon for the specified landmark type.
func (lig *LandmarkIconGenerator) GenerateAnimatedIcon(iconType LandmarkIconType, primaryColor, secondaryColor color.Color) *AnimatedLandmarkIcon {
	switch iconType {
	case LandmarkIconRuins:
		return lig.generateRuinsIcon(primaryColor, secondaryColor)
	case LandmarkIconOutpost:
		return lig.generateOutpostIcon(primaryColor, secondaryColor)
	case LandmarkIconTown:
		return lig.generateTownIcon(primaryColor, secondaryColor)
	case LandmarkIconShrine:
		return lig.generateShrineIcon(primaryColor, secondaryColor)
	case LandmarkIconOrigin:
		return lig.generateOriginIcon(primaryColor, secondaryColor)
	case LandmarkIconDestination:
		return lig.generateDestinationIcon(primaryColor, secondaryColor)
	default:
		return lig.generateTownIcon(primaryColor, secondaryColor)
	}
}

// generateRuinsIcon creates a smoking ruins animation.
func (lig *LandmarkIconGenerator) generateRuinsIcon(baseColor, smokeColor color.Color) *AnimatedLandmarkIcon {
	const frameCount = 6
	frames := make([]*ebiten.Image, frameCount)
	size := lig.iconSize

	// Generate base ruins structure
	ruinsPattern := lig.generateRuinsPattern(size)

	for f := 0; f < frameCount; f++ {
		img := ebiten.NewImage(size, size)
		lig.drawRuinsBase(img, ruinsPattern, baseColor)
		lig.drawSmoke(img, smokeColor, f, frameCount)
		frames[f] = img
	}

	return NewAnimatedLandmarkIcon(frames, 0.15)
}

// generateRuinsPattern creates the structural pattern for ruins.
func (lig *LandmarkIconGenerator) generateRuinsPattern(size int) [][]bool {
	pattern := make([][]bool, size)
	for i := range pattern {
		pattern[i] = make([]bool, size)
	}

	centerX := size / 2
	groundY := size * 3 / 4

	// Left wall with decay
	wallHeight := size/3 + lig.gen.Intn(size/6)
	lig.fillWallSection(pattern, centerX-size/4, centerX-size/8, groundY-wallHeight, groundY, 0.85)

	// Right partial wall with more decay
	rightHeight := size/4 + lig.gen.Intn(size/8)
	lig.fillWallSection(pattern, centerX+size/8, centerX+size/4, groundY-rightHeight, groundY, 0.75)

	return pattern
}

// fillWallSection fills a rectangular wall section with optional holes.
func (lig *LandmarkIconGenerator) fillWallSection(pattern [][]bool, x1, x2, y1, y2 int, fillChance float64) {
	size := len(pattern)
	for y := y1; y < y2; y++ {
		for x := x1; x < x2; x++ {
			if lig.isInBounds(x, y, size) && lig.gen.Chance(fillChance) {
				pattern[y][x] = true
			}
		}
	}
}

// isInBounds checks if coordinates are within bounds.
func (lig *LandmarkIconGenerator) isInBounds(x, y, size int) bool {
	return x >= 0 && x < size && y >= 0 && y < size
}

// drawRuinsBase draws the ruins structure.
func (lig *LandmarkIconGenerator) drawRuinsBase(img *ebiten.Image, pattern [][]bool, c color.Color) {
	size := len(pattern)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if pattern[y][x] {
				img.Set(x, y, c)
			}
		}
	}
}

// drawSmoke draws animated smoke rising from ruins.
func (lig *LandmarkIconGenerator) drawSmoke(img *ebiten.Image, smokeColor color.Color, frame, totalFrames int) {
	size := lig.iconSize
	centerX := size / 2
	baseY := size * 2 / 3

	// Multiple smoke particles at different phases
	smokeParticles := 3
	for p := 0; p < smokeParticles; p++ {
		offset := (frame + p*totalFrames/smokeParticles) % totalFrames
		progress := float64(offset) / float64(totalFrames)

		// Smoke rises and spreads
		smokeY := baseY - int(float64(size/2)*progress)
		spreadX := int(float64(size/8) * progress)
		particleX := centerX + (p-1)*size/8

		// Draw smoke particle
		alpha := uint8(200 * (1.0 - progress))
		smokeWithAlpha := color.RGBA{
			R: smokeColor.(color.RGBA).R,
			G: smokeColor.(color.RGBA).G,
			B: smokeColor.(color.RGBA).B,
			A: alpha,
		}

		for dy := -1; dy <= 1; dy++ {
			for dx := -spreadX; dx <= spreadX; dx++ {
				px := particleX + dx
				py := smokeY + dy
				if px >= 0 && px < size && py >= 0 && py < size {
					img.Set(px, py, smokeWithAlpha)
				}
			}
		}
	}
}

// generateOutpostIcon creates an outpost with blinking lights.
func (lig *LandmarkIconGenerator) generateOutpostIcon(buildingColor, lightColor color.Color) *AnimatedLandmarkIcon {
	lightPositions := lig.generateLightPositions(lig.iconSize)
	return lig.generateAnimatedIcon(4, 0.3, func(img *ebiten.Image, frame int) {
		lig.drawOutpostBase(img, buildingColor)
		lig.drawBlinkingLights(img, lightPositions, lightColor, frame)
	})
}

// generateLightPositions returns positions for outpost lights.
func (lig *LandmarkIconGenerator) generateLightPositions(size int) [][2]int {
	positions := make([][2]int, 0, 4)

	// Top of tower
	positions = append(positions, [2]int{size / 2, size / 4})

	// Window lights
	if lig.gen.Chance(0.7) {
		positions = append(positions, [2]int{size / 3, size / 2})
	}
	if lig.gen.Chance(0.7) {
		positions = append(positions, [2]int{size * 2 / 3, size / 2})
	}

	return positions
}

// drawOutpostBase draws the outpost structure.
func (lig *LandmarkIconGenerator) drawOutpostBase(img *ebiten.Image, c color.Color) {
	size := lig.iconSize
	centerX := size / 2
	groundY := size * 3 / 4

	// Main building
	buildingWidth := size / 3
	buildingHeight := size / 2
	lig.drawRect(img, centerX-buildingWidth/2, groundY-buildingHeight, centerX+buildingWidth/2, groundY, c)

	// Tower/antenna
	towerHeight := size / 4
	lig.drawVerticalLine(img, centerX, groundY-buildingHeight-towerHeight, groundY-buildingHeight, c)
}

// drawBlinkingLights draws lights that blink on different frames.
func (lig *LandmarkIconGenerator) drawBlinkingLights(img *ebiten.Image, positions [][2]int, lightColor color.Color, frame int) {
	for i, pos := range positions {
		// Each light has its own blink pattern
		if (frame+i)%2 == 0 {
			// Light is on
			lig.drawLight(img, pos[0], pos[1], lightColor)
		}
	}
}

// drawLight draws a small glowing light.
func (lig *LandmarkIconGenerator) drawLight(img *ebiten.Image, x, y int, c color.Color) {
	size := lig.iconSize

	// Core pixel
	if x >= 0 && x < size && y >= 0 && y < size {
		img.Set(x, y, c)
	}

	// Glow (dimmer neighbors)
	r, g, b, a := c.RGBA()
	dimColor := color.RGBA{
		R: uint8((r >> 8) / 2),
		G: uint8((g >> 8) / 2),
		B: uint8((b >> 8) / 2),
		A: uint8(a >> 8),
	}

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			px, py := x+dx, y+dy
			if px >= 0 && px < size && py >= 0 && py < size {
				img.Set(px, py, dimColor)
			}
		}
	}
}

// generateTownIcon creates a town with flickering window lights.
func (lig *LandmarkIconGenerator) generateTownIcon(buildingColor, lightColor color.Color) *AnimatedLandmarkIcon {
	windowPositions := lig.generateWindowPositions(lig.iconSize)
	return lig.generateAnimatedIcon(4, 0.25, func(img *ebiten.Image, frame int) {
		lig.drawTownBase(img, buildingColor)
		lig.drawFlickeringWindows(img, windowPositions, lightColor, frame)
	})
}

// generateWindowPositions returns positions for town windows.
func (lig *LandmarkIconGenerator) generateWindowPositions(size int) [][2]int {
	positions := make([][2]int, 0, 6)

	// Generate 4-6 window positions
	windowCount := 4 + lig.gen.Intn(3)
	for i := 0; i < windowCount; i++ {
		x := size/4 + lig.gen.Intn(size/2)
		y := size/3 + lig.gen.Intn(size/3)
		positions = append(positions, [2]int{x, y})
	}

	return positions
}

// drawTownBase draws the town buildings.
func (lig *LandmarkIconGenerator) drawTownBase(img *ebiten.Image, c color.Color) {
	size := lig.iconSize
	groundY := size * 3 / 4

	// Multiple buildings of varying heights
	buildingCount := 3
	buildingWidth := size / 4

	for b := 0; b < buildingCount; b++ {
		startX := b*size/buildingCount + size/12
		height := size/3 + lig.gen.Intn(size/6)

		for y := groundY - height; y < groundY; y++ {
			for x := startX; x < startX+buildingWidth && x < size; x++ {
				if x >= 0 && y >= 0 && y < size {
					img.Set(x, y, c)
				}
			}
		}
	}
}

// drawFlickeringWindows draws windows with random flicker.
func (lig *LandmarkIconGenerator) drawFlickeringWindows(img *ebiten.Image, positions [][2]int, lightColor color.Color, frame int) {
	for i, pos := range positions {
		// Each window flickers randomly based on frame and position
		flicker := (frame + i*7) % 4
		if flicker < 3 { // On 75% of the time
			img.Set(pos[0], pos[1], lightColor)
		}
	}
}

// generateShrineIcon creates a shrine with glowing effect.
func (lig *LandmarkIconGenerator) generateShrineIcon(stoneColor, glowColor color.Color) *AnimatedLandmarkIcon {
	const frameCount = 4
	return lig.generateAnimatedIcon(frameCount, 0.2, func(img *ebiten.Image, frame int) {
		lig.drawShrineBase(img, stoneColor)
		lig.drawShrineGlow(img, glowColor, frame, frameCount)
	})
}

// drawShrineBase draws the shrine structure.
func (lig *LandmarkIconGenerator) drawShrineBase(img *ebiten.Image, c color.Color) {
	size := lig.iconSize
	centerX := size / 2
	groundY := size * 3 / 4

	// Altar/pedestal
	pedestalWidth := size / 3
	pedestalHeight := size / 6
	lig.drawRect(img, centerX-pedestalWidth/2, groundY-pedestalHeight, centerX+pedestalWidth/2, groundY, c)

	// Central pillar/artifact (3 pixels wide)
	pillarHeight := size / 3
	pillarTop := groundY - pedestalHeight - pillarHeight
	pillarBottom := groundY - pedestalHeight
	lig.drawVerticalLine(img, centerX, pillarTop, pillarBottom, c)
	lig.drawVerticalLine(img, centerX-1, pillarTop, pillarBottom, c)
	lig.drawVerticalLine(img, centerX+1, pillarTop, pillarBottom, c)
}

// drawShrineGlow draws the pulsing glow effect.
func (lig *LandmarkIconGenerator) drawShrineGlow(img *ebiten.Image, glowColor color.Color, frame, totalFrames int) {
	size := lig.iconSize
	centerX := size / 2
	centerY := size / 2

	// Pulsing intensity
	progress := float64(frame) / float64(totalFrames)
	intensity := 0.5 + 0.5*sinApprox(progress*6.28)

	r, g, b, _ := glowColor.RGBA()
	alpha := uint8(100 * intensity)

	glowRadius := size/6 + int(float64(size/8)*intensity)

	for dy := -glowRadius; dy <= glowRadius; dy++ {
		for dx := -glowRadius; dx <= glowRadius; dx++ {
			dist := dx*dx + dy*dy
			if dist <= glowRadius*glowRadius {
				px, py := centerX+dx, centerY+dy
				if px >= 0 && px < size && py >= 0 && py < size {
					distFactor := 1.0 - float64(dist)/float64(glowRadius*glowRadius)
					pixelAlpha := uint8(float64(alpha) * distFactor)
					glowPixel := color.RGBA{
						R: uint8(r >> 8),
						G: uint8(g >> 8),
						B: uint8(b >> 8),
						A: pixelAlpha,
					}
					img.Set(px, py, glowPixel)
				}
			}
		}
	}
}

// generateOriginIcon creates a starting point marker.
func (lig *LandmarkIconGenerator) generateOriginIcon(markerColor, glowColor color.Color) *AnimatedLandmarkIcon {
	const frameCount = 4
	frames := make([]*ebiten.Image, frameCount)
	size := lig.iconSize

	for f := 0; f < frameCount; f++ {
		img := ebiten.NewImage(size, size)
		lig.drawOriginMarker(img, markerColor, glowColor, f, frameCount)
		frames[f] = img
	}

	return NewAnimatedLandmarkIcon(frames, 0.2)
}

// drawOriginMarker draws a pulsing origin marker.
func (lig *LandmarkIconGenerator) drawOriginMarker(img *ebiten.Image, markerColor, glowColor color.Color, frame, totalFrames int) {
	size := lig.iconSize
	centerX := size / 2
	centerY := size / 2

	// Calculate pulse size based on animation frame
	progress := float64(frame) / float64(totalFrames)
	pulseSize := size/4 + int(float64(size/8)*sinApprox(progress*6.28))

	// Draw glow and center marker
	lig.drawFilledCircle(img, centerX, centerY, pulseSize, glowColor)
	lig.drawFilledCircle(img, centerX, centerY, size/6, markerColor)
}

// drawFilledCircle draws a filled circle at the given center point with the specified radius.
func (lig *LandmarkIconGenerator) drawFilledCircle(img *ebiten.Image, centerX, centerY, radius int, c color.Color) {
	size := lig.iconSize
	for dy := -radius; dy <= radius; dy++ {
		for dx := -radius; dx <= radius; dx++ {
			if dx*dx+dy*dy <= radius*radius {
				px, py := centerX+dx, centerY+dy
				if px >= 0 && px < size && py >= 0 && py < size {
					img.Set(px, py, c)
				}
			}
		}
	}
}

// generateDestinationIcon creates a goal marker with beacon effect.
func (lig *LandmarkIconGenerator) generateDestinationIcon(markerColor, beamColor color.Color) *AnimatedLandmarkIcon {
	const frameCount = 6
	frames := make([]*ebiten.Image, frameCount)
	size := lig.iconSize

	for f := 0; f < frameCount; f++ {
		img := ebiten.NewImage(size, size)
		lig.drawDestinationMarker(img, markerColor, beamColor, f, frameCount)
		frames[f] = img
	}

	return NewAnimatedLandmarkIcon(frames, 0.15)
}

// drawDestinationMarker draws a beacon-style destination marker.
func (lig *LandmarkIconGenerator) drawDestinationMarker(img *ebiten.Image, markerColor, beamColor color.Color, frame, totalFrames int) {
	size := lig.iconSize
	centerX := size / 2
	groundY := size * 3 / 4

	lig.drawDestinationBase(img, centerX, groundY, markerColor)
	lig.drawBeaconBeam(img, centerX, groundY, beamColor, frame, totalFrames)
}

// drawDestinationBase draws the base structure (star/flag shape) for a destination marker.
func (lig *LandmarkIconGenerator) drawDestinationBase(img *ebiten.Image, centerX, groundY int, markerColor color.Color) {
	size := lig.iconSize
	for y := groundY - size/3; y < groundY; y++ {
		width := (groundY - y) / 3
		if width < 1 {
			width = 1
		}
		for x := centerX - width; x <= centerX+width; x++ {
			if x >= 0 && x < size && y >= 0 && y < size {
				img.Set(x, y, markerColor)
			}
		}
	}
}

// drawBeaconBeam draws the rising beam effect for a destination marker.
func (lig *LandmarkIconGenerator) drawBeaconBeam(img *ebiten.Image, centerX, groundY int, beamColor color.Color, frame, totalFrames int) {
	size := lig.iconSize
	beamOffset := (frame * size / totalFrames) % size
	r, g, b, _ := beamColor.RGBA()

	for y := 0; y < size/2; y++ {
		beamY := (groundY - size/3 - y + beamOffset) % size
		if beamY < 0 {
			beamY += size
		}
		if beamY < groundY-size/3 {
			alpha := uint8(150 * (1.0 - float64(y)/float64(size/2)))
			pixelColor := color.RGBA{
				R: uint8(r >> 8),
				G: uint8(g >> 8),
				B: uint8(b >> 8),
				A: alpha,
			}
			img.Set(centerX, beamY, pixelColor)
		}
	}
}
