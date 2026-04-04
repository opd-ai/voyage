//go:build !headless

package ux

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/events"
)

// EventOverlay displays event text and choices.
type EventOverlay struct {
	skin           *UISkin
	genre          engine.GenreID
	overlayWidth   int
	overlayHeight  int
	selectedChoice int
	visible        bool
}

// NewEventOverlay creates a new event overlay.
func NewEventOverlay(genre engine.GenreID, width, height int) *EventOverlay {
	return &EventOverlay{
		skin:           DefaultSkin(genre),
		genre:          genre,
		overlayWidth:   width,
		overlayHeight:  height,
		selectedChoice: 0,
		visible:        false,
	}
}

// SetGenre changes the overlay's visual theme.
func (eo *EventOverlay) SetGenre(genre engine.GenreID) {
	eo.genre = genre
	eo.skin = DefaultSkin(genre)
}

// Show makes the overlay visible.
func (eo *EventOverlay) Show() {
	eo.visible = true
	eo.selectedChoice = 0
}

// Hide hides the overlay.
func (eo *EventOverlay) Hide() {
	eo.visible = false
}

// IsVisible returns true if the overlay is visible.
func (eo *EventOverlay) IsVisible() bool {
	return eo.visible
}

// SelectNext moves selection to the next choice.
func (eo *EventOverlay) SelectNext(maxChoices int) {
	eo.selectedChoice = (eo.selectedChoice + 1) % maxChoices
}

// SelectPrev moves selection to the previous choice.
func (eo *EventOverlay) SelectPrev(maxChoices int) {
	eo.selectedChoice--
	if eo.selectedChoice < 0 {
		eo.selectedChoice = maxChoices - 1
	}
}

// SelectByNumber selects a choice by number (1-indexed).
func (eo *EventOverlay) SelectByNumber(num, maxChoices int) bool {
	if num >= 1 && num <= maxChoices {
		eo.selectedChoice = num - 1
		return true
	}
	return false
}

// SelectedChoice returns the currently selected choice index.
func (eo *EventOverlay) SelectedChoice() int {
	return eo.selectedChoice
}

// Draw renders the event overlay to the screen.
func (eo *EventOverlay) Draw(screen *ebiten.Image, event *events.Event, screenWidth, screenHeight int) {
	if !eo.visible || event == nil {
		return
	}

	// Calculate overlay position (centered)
	overlayX := (screenWidth - eo.overlayWidth) / 2
	overlayY := (screenHeight - eo.overlayHeight) / 2

	// Create overlay image
	overlay := ebiten.NewImage(eo.overlayWidth, eo.overlayHeight)
	overlay.Fill(eo.skin.PanelBackground)
	eo.drawBorder(overlay)

	// Draw event category and title
	padding := 12
	y := padding
	categoryText := fmt.Sprintf("[%s]", events.CategoryName(event.Category, eo.genre))
	ebitenutil.DebugPrintAt(overlay, categoryText, padding, y)
	y += 20

	// Title
	ebitenutil.DebugPrintAt(overlay, event.Title, padding, y)
	y += 24

	// Description (word-wrapped)
	desc := eo.wrapText(event.Description, (eo.overlayWidth-padding*2)/7)
	for _, line := range desc {
		ebitenutil.DebugPrintAt(overlay, line, padding, y)
		y += 16
	}
	y += 16

	// Separator
	ebitenutil.DebugPrintAt(overlay, "---", padding, y)
	y += 20

	// Choices
	for i, choice := range event.Choices {
		prefix := "  "
		if i == eo.selectedChoice {
			prefix = "> "
		}
		choiceText := fmt.Sprintf("%s%d. %s", prefix, i+1, choice.Text)
		ebitenutil.DebugPrintAt(overlay, choiceText, padding, y)
		y += 20
	}

	// Instructions
	y = eo.overlayHeight - padding - 16
	ebitenutil.DebugPrintAt(overlay, "Press 1-4 or ENTER to select", padding, y)

	// Draw overlay to screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(overlayX), float64(overlayY))
	screen.DrawImage(overlay, op)
}

// drawBorder draws a border around the overlay.
func (eo *EventOverlay) drawBorder(img *ebiten.Image) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	c := eo.skin.PanelBorder

	for x := 0; x < w; x++ {
		img.Set(x, 0, c)
		img.Set(x, 1, c)
		img.Set(x, h-1, c)
		img.Set(x, h-2, c)
	}
	for y := 0; y < h; y++ {
		img.Set(0, y, c)
		img.Set(1, y, c)
		img.Set(w-1, y, c)
		img.Set(w-2, y, c)
	}
}

// wrapText wraps text to fit within a maximum character width.
func (eo *EventOverlay) wrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		maxWidth = 40
	}

	var lines []string
	var currentLine string

	for _, word := range splitWords(text) {
		if len(currentLine)+len(word)+1 > maxWidth {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		} else {
			if currentLine != "" {
				currentLine += " "
			}
			currentLine += word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}

// splitWords splits text into words.
func splitWords(text string) []string {
	var words []string
	var current string

	for _, c := range text {
		if c == ' ' || c == '\n' {
			if current != "" {
				words = append(words, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		words = append(words, current)
	}

	return words
}
