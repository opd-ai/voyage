//go:build !headless

package ux

import (
	"image/color"

	"github.com/opd-ai/voyage/pkg/engine"
)

// UISkin holds genre-themed UI styling.
type UISkin struct {
	PanelBackground   color.Color
	PanelBorder       color.Color
	TextPrimary       color.Color
	TextSecondary     color.Color
	BarFill           color.Color
	BarBackground     color.Color
	ButtonNormal      color.Color
	ButtonHover       color.Color
	ButtonPressed     color.Color
	WarningColor      color.Color
	CriticalColor     color.Color
	HighlightColor    color.Color
	MenuTitle         string
	ResourcePanelName string
	CrewPanelName     string
}

// DefaultSkin returns the UI skin for the given genre.
func DefaultSkin(genre engine.GenreID) *UISkin {
	switch genre {
	case engine.GenreScifi:
		return &UISkin{
			PanelBackground:   color.RGBA{10, 15, 30, 220},
			PanelBorder:       color.RGBA{50, 100, 200, 255},
			TextPrimary:       color.RGBA{200, 220, 255, 255},
			TextSecondary:     color.RGBA{100, 140, 180, 255},
			BarFill:           color.RGBA{0, 200, 255, 255},
			BarBackground:     color.RGBA{20, 30, 50, 255},
			ButtonNormal:      color.RGBA{30, 50, 100, 255},
			ButtonHover:       color.RGBA{50, 80, 150, 255},
			ButtonPressed:     color.RGBA{70, 100, 180, 255},
			WarningColor:      color.RGBA{255, 200, 0, 255},
			CriticalColor:     color.RGBA{255, 50, 50, 255},
			HighlightColor:    color.RGBA{0, 255, 200, 255},
			MenuTitle:         "VOYAGE: STELLAR DRIFT",
			ResourcePanelName: "SHIP SYSTEMS",
			CrewPanelName:     "CREW MANIFEST",
		}
	case engine.GenreHorror:
		return &UISkin{
			PanelBackground:   color.RGBA{20, 15, 15, 220},
			PanelBorder:       color.RGBA{100, 40, 40, 255},
			TextPrimary:       color.RGBA{180, 160, 150, 255},
			TextSecondary:     color.RGBA{120, 100, 90, 255},
			BarFill:           color.RGBA{150, 50, 50, 255},
			BarBackground:     color.RGBA{40, 30, 30, 255},
			ButtonNormal:      color.RGBA{60, 40, 40, 255},
			ButtonHover:       color.RGBA{80, 50, 50, 255},
			ButtonPressed:     color.RGBA{100, 60, 60, 255},
			WarningColor:      color.RGBA{200, 150, 50, 255},
			CriticalColor:     color.RGBA{150, 0, 0, 255},
			HighlightColor:    color.RGBA{200, 50, 50, 255},
			MenuTitle:         "VOYAGE: DEAD ROADS",
			ResourcePanelName: "SUPPLIES",
			CrewPanelName:     "SURVIVORS",
		}
	case engine.GenreCyberpunk:
		return &UISkin{
			PanelBackground:   color.RGBA{15, 15, 25, 220},
			PanelBorder:       color.RGBA{255, 0, 100, 255},
			TextPrimary:       color.RGBA{200, 200, 220, 255},
			TextSecondary:     color.RGBA{100, 150, 200, 255},
			BarFill:           color.RGBA{0, 200, 255, 255},
			BarBackground:     color.RGBA{30, 30, 40, 255},
			ButtonNormal:      color.RGBA{40, 40, 60, 255},
			ButtonHover:       color.RGBA{60, 60, 90, 255},
			ButtonPressed:     color.RGBA{80, 80, 120, 255},
			WarningColor:      color.RGBA{255, 150, 0, 255},
			CriticalColor:     color.RGBA{255, 0, 50, 255},
			HighlightColor:    color.RGBA{255, 255, 0, 255},
			MenuTitle:         "VOYAGE: NEON EXODUS",
			ResourcePanelName: "ASSETS",
			CrewPanelName:     "CREW",
		}
	case engine.GenrePostapoc:
		return &UISkin{
			PanelBackground:   color.RGBA{35, 30, 25, 220},
			PanelBorder:       color.RGBA{150, 100, 50, 255},
			TextPrimary:       color.RGBA{200, 180, 150, 255},
			TextSecondary:     color.RGBA{150, 130, 100, 255},
			BarFill:           color.RGBA{200, 150, 50, 255},
			BarBackground:     color.RGBA{50, 40, 30, 255},
			ButtonNormal:      color.RGBA{80, 60, 40, 255},
			ButtonHover:       color.RGBA{100, 80, 50, 255},
			ButtonPressed:     color.RGBA{120, 100, 60, 255},
			WarningColor:      color.RGBA{200, 100, 50, 255},
			CriticalColor:     color.RGBA{180, 50, 30, 255},
			HighlightColor:    color.RGBA{200, 150, 50, 255},
			MenuTitle:         "VOYAGE: WASTELAND TRAIL",
			ResourcePanelName: "STASH",
			CrewPanelName:     "PARTY",
		}
	default: // Fantasy
		return &UISkin{
			PanelBackground:   color.RGBA{30, 35, 30, 220},
			PanelBorder:       color.RGBA{120, 100, 60, 255},
			TextPrimary:       color.RGBA{230, 220, 200, 255},
			TextSecondary:     color.RGBA{160, 150, 130, 255},
			BarFill:           color.RGBA{80, 150, 80, 255},
			BarBackground:     color.RGBA{40, 45, 40, 255},
			ButtonNormal:      color.RGBA{60, 70, 50, 255},
			ButtonHover:       color.RGBA{80, 100, 60, 255},
			ButtonPressed:     color.RGBA{100, 120, 70, 255},
			WarningColor:      color.RGBA{200, 150, 50, 255},
			CriticalColor:     color.RGBA{180, 50, 50, 255},
			HighlightColor:    color.RGBA{200, 180, 100, 255},
			MenuTitle:         "VOYAGE: SILK ROAD",
			ResourcePanelName: "SUPPLIES",
			CrewPanelName:     "PARTY",
		}
	}
}
