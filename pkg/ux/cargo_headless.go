//go:build headless

package ux

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/vessel"
)

// CargoScreen displays cargo hold contents with weight/volume limits.
// This is a headless stub for testing.
type CargoScreen struct {
	genre        engine.GenreID
	screenWidth  int
	screenHeight int
	scrollOffset int
	visible      bool
}

// NewCargoScreen creates a new cargo management screen (headless stub).
func NewCargoScreen(genre engine.GenreID, screenWidth, screenHeight int) *CargoScreen {
	return &CargoScreen{
		genre:        genre,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
		scrollOffset: 0,
		visible:      false,
	}
}

// SetGenre changes the screen's visual theme.
func (cs *CargoScreen) SetGenre(genre engine.GenreID) {
	cs.genre = genre
}

// Show makes the cargo screen visible.
func (cs *CargoScreen) Show() {
	cs.visible = true
	cs.scrollOffset = 0
}

// Hide makes the cargo screen hidden.
func (cs *CargoScreen) Hide() {
	cs.visible = false
}

// IsVisible returns whether the screen is currently visible.
func (cs *CargoScreen) IsVisible() bool {
	return cs.visible
}

// ScrollUp scrolls the cargo list up.
func (cs *CargoScreen) ScrollUp() {
	if cs.scrollOffset > 0 {
		cs.scrollOffset--
	}
}

// ScrollDown scrolls the cargo list down.
func (cs *CargoScreen) ScrollDown() {
	cs.scrollOffset++
}

// CargoSummary provides cargo statistics for display.
type CargoSummary struct {
	TotalItems     int
	TotalWeight    int
	TotalVolume    int
	WeightLimit    int
	VolumeLimit    int
	WeightRatio    float64
	VolumeRatio    float64
	Tier           int
	CategoryCounts map[vessel.CargoCategory]int
}

// GetCargoSummary returns a summary of cargo hold contents.
func GetCargoSummary(hold *vessel.CargoHold) CargoSummary {
	summary := CargoSummary{
		TotalItems:     len(hold.Items()),
		TotalWeight:    hold.UsedWeight(),
		TotalVolume:    hold.UsedVolume(),
		WeightLimit:    hold.WeightLimit(),
		VolumeLimit:    hold.VolumeLimit(),
		WeightRatio:    hold.WeightRatio(),
		VolumeRatio:    hold.VolumeRatio(),
		Tier:           hold.Tier(),
		CategoryCounts: make(map[vessel.CargoCategory]int),
	}

	for _, item := range hold.Items() {
		summary.CategoryCounts[item.Category] += item.Quantity
	}

	return summary
}
