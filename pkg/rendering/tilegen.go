package rendering

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// TileGenerator creates procedural tile sprites using cellular automata.
type TileGenerator struct {
	gen      *seed.Generator
	tileSize int
}

// NewTileGenerator creates a new procedural tile generator.
func NewTileGenerator(masterSeed int64, tileSize int) *TileGenerator {
	return &TileGenerator{
		gen:      seed.NewGenerator(masterSeed, "tiles"),
		tileSize: tileSize,
	}
}

// GenerateTile creates a procedural tile image using cellular automata.
func (tg *TileGenerator) GenerateTile(baseColor color.Color, variation float64) *ebiten.Image {
	img := ebiten.NewImage(tg.tileSize, tg.tileSize)

	r, g, b, a := baseColor.RGBA()
	baseR := uint8(r >> 8)
	baseG := uint8(g >> 8)
	baseB := uint8(b >> 8)
	baseA := uint8(a >> 8)

	// Generate noise pattern
	for y := 0; y < tg.tileSize; y++ {
		for x := 0; x < tg.tileSize; x++ {
			// Apply variation
			varR := int(float64(baseR) * (1.0 + (tg.gen.Float64()-0.5)*variation))
			varG := int(float64(baseG) * (1.0 + (tg.gen.Float64()-0.5)*variation))
			varB := int(float64(baseB) * (1.0 + (tg.gen.Float64()-0.5)*variation))

			// Clamp values
			varR = clamp(varR, 0, 255)
			varG = clamp(varG, 0, 255)
			varB = clamp(varB, 0, 255)

			img.Set(x, y, color.RGBA{uint8(varR), uint8(varG), uint8(varB), baseA})
		}
	}

	return img
}

// GenerateTileCA creates a tile using cellular automata patterns.
func (tg *TileGenerator) GenerateTileCA(baseColor, accentColor color.Color, density float64, iterations int) *ebiten.Image {
	size := tg.tileSize
	grid := make([][]bool, size)
	for i := range grid {
		grid[i] = make([]bool, size)
	}

	// Initialize with random cells
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			grid[y][x] = tg.gen.Chance(density)
		}
	}

	// Run cellular automata
	for iter := 0; iter < iterations; iter++ {
		newGrid := make([][]bool, size)
		for i := range newGrid {
			newGrid[i] = make([]bool, size)
		}

		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				neighbors := tg.countNeighbors(grid, x, y)
				if grid[y][x] {
					newGrid[y][x] = neighbors >= 4
				} else {
					newGrid[y][x] = neighbors >= 5
				}
			}
		}
		grid = newGrid
	}

	// Render to image
	img := ebiten.NewImage(size, size)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if grid[y][x] {
				img.Set(x, y, accentColor)
			} else {
				img.Set(x, y, baseColor)
			}
		}
	}

	return img
}

// countNeighbors counts alive neighbors for cellular automata.
func (tg *TileGenerator) countNeighbors(grid [][]bool, x, y int) int {
	count := 0
	size := len(grid)
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < size && ny >= 0 && ny < size {
				if grid[ny][nx] {
					count++
				}
			}
		}
	}
	return count
}

// GenerateCharacterSprite creates a procedural character sprite.
func (tg *TileGenerator) GenerateCharacterSprite(primaryColor, secondaryColor color.Color) *ebiten.Image {
	size := tg.tileSize
	img := ebiten.NewImage(size, size)

	// Create a simple symmetric sprite
	halfWidth := size / 2
	for y := 0; y < size; y++ {
		for x := 0; x < halfWidth; x++ {
			if tg.gen.Chance(0.4) {
				col := primaryColor
				if tg.gen.Chance(0.3) {
					col = secondaryColor
				}
				img.Set(x, y, col)
				img.Set(size-1-x, y, col) // Mirror
			}
		}
	}

	return img
}

// clamp restricts a value to a range.
func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
