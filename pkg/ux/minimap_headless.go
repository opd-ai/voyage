//go:build headless

package ux

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/world"
)

// Minimap displays a small overview of the world map in the corner of the screen.
// This is the headless implementation for testing without graphics.
type Minimap struct {
	genre      engine.GenreID
	width      int
	height     int
	tileSize   int
	alpha      float64
	crisisMode bool
}

// NewMinimap creates a new minimap with the given dimensions.
func NewMinimap(genre engine.GenreID, width, height int) *Minimap {
	return &Minimap{
		genre:    genre,
		width:    width,
		height:   height,
		tileSize: 4,
		alpha:    1.0,
	}
}

// SetGenre changes the minimap's visual theme.
func (m *Minimap) SetGenre(genre engine.GenreID) {
	m.genre = genre
}

// SetCrisisMode enables/disables crisis mode (fades the minimap during encounters).
func (m *Minimap) SetCrisisMode(enabled bool) {
	m.crisisMode = enabled
	if enabled {
		m.alpha = 0.3
	} else {
		m.alpha = 1.0
	}
}

// IsCrisisMode returns whether crisis mode is enabled.
func (m *Minimap) IsCrisisMode() bool {
	return m.crisisMode
}

// Width returns the minimap width.
func (m *Minimap) Width() int {
	return m.width
}

// Height returns the minimap height.
func (m *Minimap) Height() int {
	return m.height
}

// Genre returns the current genre.
func (m *Minimap) Genre() engine.GenreID {
	return m.genre
}

// Alpha returns the current alpha transparency.
func (m *Minimap) Alpha() float64 {
	return m.alpha
}

// CalculateScale computes the scale factor for mapping world to minimap coordinates.
func (m *Minimap) CalculateScale(wm *world.WorldMap) float64 {
	if wm == nil {
		return 1.0
	}
	scaleX := float64(m.width-2) / float64(wm.Width)
	scaleY := float64(m.height-2) / float64(wm.Height)
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}
	return scale
}

// WorldToMinimap converts world coordinates to minimap coordinates.
func (m *Minimap) WorldToMinimap(wm *world.WorldMap, worldX, worldY int) (int, int) {
	scale := m.CalculateScale(wm)
	px := 1 + int(float64(worldX)*scale) + m.tileSize/2
	py := 1 + int(float64(worldY)*scale) + m.tileSize/2
	return px, py
}

// CountExploredTiles counts the number of explored tiles in the world map.
func (m *Minimap) CountExploredTiles(wm *world.WorldMap) int {
	if wm == nil {
		return 0
	}
	count := 0
	for y := 0; y < wm.Height; y++ {
		for x := 0; x < wm.Width; x++ {
			tile := wm.GetTile(x, y)
			if tile != nil && tile.Explored {
				count++
			}
		}
	}
	return count
}

// CountLandmarks counts the number of explored landmarks in the world map.
func (m *Minimap) CountLandmarks(wm *world.WorldMap) int {
	if wm == nil {
		return 0
	}
	count := 0
	for y := 0; y < wm.Height; y++ {
		for x := 0; x < wm.Width; x++ {
			tile := wm.GetTile(x, y)
			if tile != nil && tile.Landmark != nil && tile.Explored {
				count++
			}
		}
	}
	return count
}
