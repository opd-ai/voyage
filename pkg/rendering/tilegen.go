//go:build !headless

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
	grid := tg.initializeCAGrid(size, density)
	grid = tg.runCAIterations(grid, iterations)
	return tg.renderCAGrid(grid, baseColor, accentColor)
}

// initializeCAGrid creates a new grid with random cell distribution.
func (tg *TileGenerator) initializeCAGrid(size int, density float64) [][]bool {
	grid := make([][]bool, size)
	for i := range grid {
		grid[i] = make([]bool, size)
	}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			grid[y][x] = tg.gen.Chance(density)
		}
	}
	return grid
}

// runCAIterations applies cellular automata rules for the specified iterations.
func (tg *TileGenerator) runCAIterations(grid [][]bool, iterations int) [][]bool {
	size := len(grid)
	for iter := 0; iter < iterations; iter++ {
		grid = tg.stepCA(grid, size)
	}
	return grid
}

// stepCA performs a single cellular automata iteration.
func (tg *TileGenerator) stepCA(grid [][]bool, size int) [][]bool {
	newGrid := make([][]bool, size)
	for i := range newGrid {
		newGrid[i] = make([]bool, size)
	}
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			neighbors := tg.countNeighbors(grid, x, y)
			newGrid[y][x] = tg.applyCARules(grid[y][x], neighbors)
		}
	}
	return newGrid
}

// applyCARules determines cell state based on neighbor count.
func (tg *TileGenerator) applyCARules(alive bool, neighbors int) bool {
	if alive {
		return neighbors >= 4
	}
	return neighbors >= 5
}

// renderCAGrid converts the boolean grid to an image.
func (tg *TileGenerator) renderCAGrid(grid [][]bool, baseColor, accentColor color.Color) *ebiten.Image {
	size := len(grid)
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
			if tg.isAliveNeighbor(grid, x+dx, y+dy, size) {
				count++
			}
		}
	}
	return count
}

// isAliveNeighbor checks if a neighbor cell is alive and in bounds.
func (tg *TileGenerator) isAliveNeighbor(grid [][]bool, nx, ny, size int) bool {
	return nx >= 0 && nx < size && ny >= 0 && ny < size && grid[ny][nx]
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
