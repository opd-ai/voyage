//go:build !headless

package ux

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/saveload"
)

// SlotSelectionMode determines the slot screen's purpose.
type SlotSelectionMode int

const (
	// SlotModeLoad displays slots for loading a saved game.
	SlotModeLoad SlotSelectionMode = iota
	// SlotModeSave displays slots for saving the current game.
	SlotModeSave
)

// SlotSelectionScreen manages the save/load slot selection UI.
type SlotSelectionScreen struct {
	skin          *UISkin
	genre         engine.GenreID
	mode          SlotSelectionMode
	slots         []saveload.SlotInfo
	selectedIndex int
	screenWidth   int
	screenHeight  int
	confirmDelete bool
}

// NewSlotSelectionScreen creates a new slot selection screen.
func NewSlotSelectionScreen(genre engine.GenreID, mode SlotSelectionMode, screenWidth, screenHeight int) *SlotSelectionScreen {
	return &SlotSelectionScreen{
		skin:          DefaultSkin(genre),
		genre:         genre,
		mode:          mode,
		slots:         nil,
		selectedIndex: 0,
		screenWidth:   screenWidth,
		screenHeight:  screenHeight,
		confirmDelete: false,
	}
}

// SetGenre changes the screen's visual theme.
func (s *SlotSelectionScreen) SetGenre(genre engine.GenreID) {
	s.genre = genre
	s.skin = DefaultSkin(genre)
}

// SetMode changes between load and save modes.
func (s *SlotSelectionScreen) SetMode(mode SlotSelectionMode) {
	s.mode = mode
	s.selectedIndex = 0
	s.confirmDelete = false
}

// UpdateSlots refreshes the slot information from the save manager.
func (s *SlotSelectionScreen) UpdateSlots(manager *saveload.SaveManager) {
	s.slots = manager.ListSlots()
}

// SelectNext moves selection to the next slot.
func (s *SlotSelectionScreen) SelectNext() {
	if len(s.slots) == 0 {
		return
	}
	s.selectedIndex = (s.selectedIndex + 1) % len(s.slots)
	s.confirmDelete = false
}

// SelectPrev moves selection to the previous slot.
func (s *SlotSelectionScreen) SelectPrev() {
	if len(s.slots) == 0 {
		return
	}
	s.selectedIndex--
	if s.selectedIndex < 0 {
		s.selectedIndex = len(s.slots) - 1
	}
	s.confirmDelete = false
}

// SelectedSlot returns the currently selected slot number.
func (s *SlotSelectionScreen) SelectedSlot() int {
	if s.selectedIndex >= 0 && s.selectedIndex < len(s.slots) {
		return s.slots[s.selectedIndex].Slot
	}
	return -1
}

// SelectedSlotInfo returns the info for the selected slot.
func (s *SlotSelectionScreen) SelectedSlotInfo() *saveload.SlotInfo {
	if s.selectedIndex >= 0 && s.selectedIndex < len(s.slots) {
		return &s.slots[s.selectedIndex]
	}
	return nil
}

// CanSelect returns true if the current selection can be used.
func (s *SlotSelectionScreen) CanSelect() bool {
	info := s.SelectedSlotInfo()
	if info == nil {
		return false
	}

	if s.mode == SlotModeLoad {
		return !info.Empty
	}
	// Save mode: can always select (will overwrite if not empty)
	return true
}

// ToggleDeleteConfirm toggles delete confirmation for the selected slot.
func (s *SlotSelectionScreen) ToggleDeleteConfirm() {
	info := s.SelectedSlotInfo()
	if info != nil && !info.Empty {
		s.confirmDelete = !s.confirmDelete
	}
}

// IsConfirmingDelete returns true if awaiting delete confirmation.
func (s *SlotSelectionScreen) IsConfirmingDelete() bool {
	return s.confirmDelete
}

// CancelDelete cancels the delete confirmation.
func (s *SlotSelectionScreen) CancelDelete() {
	s.confirmDelete = false
}

// Draw renders the slot selection screen.
func (s *SlotSelectionScreen) Draw(screen *ebiten.Image) {
	// Draw semi-transparent background
	overlay := ebiten.NewImage(s.screenWidth, s.screenHeight)
	overlay.Fill(s.skin.PanelBackground)
	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(0.8)
	screen.DrawImage(overlay, op)

	// Calculate panel dimensions
	panelWidth := 450
	panelHeight := 60 + len(s.slots)*50 + 40
	if panelHeight > s.screenHeight-60 {
		panelHeight = s.screenHeight - 60
	}
	panelX := (s.screenWidth - panelWidth) / 2
	panelY := (s.screenHeight - panelHeight) / 2

	// Draw main panel
	panel := ebiten.NewImage(panelWidth, panelHeight)
	panel.Fill(s.skin.PanelBackground)
	s.drawBorder(panel)

	// Draw title
	title := s.getTitle()
	titleX := (panelWidth - len(title)*7) / 2
	ebitenutil.DebugPrintAt(panel, title, titleX, 16)

	// Draw slots
	y := 50
	visibleSlots := (panelHeight - 100) / 50
	startIdx := 0
	if s.selectedIndex >= visibleSlots {
		startIdx = s.selectedIndex - visibleSlots + 1
	}

	for i := startIdx; i < len(s.slots) && y < panelHeight-50; i++ {
		slot := s.slots[i]
		s.drawSlotEntry(panel, slot, i == s.selectedIndex, 20, y)
		y += 45
	}

	// Draw panel to screen
	opPanel := &ebiten.DrawImageOptions{}
	opPanel.GeoM.Translate(float64(panelX), float64(panelY))
	screen.DrawImage(panel, opPanel)

	// Draw instructions at bottom
	instructions := s.getInstructions()
	instrX := (s.screenWidth - len(instructions)*7) / 2
	ebitenutil.DebugPrintAt(screen, instructions, instrX, s.screenHeight-30)

	// Draw delete confirmation if active
	if s.confirmDelete {
		s.drawDeleteConfirm(screen)
	}
}

// drawSlotEntry renders a single slot entry.
func (s *SlotSelectionScreen) drawSlotEntry(panel *ebiten.Image, slot saveload.SlotInfo, selected bool, x, y int) {
	prefix := "  "
	if selected {
		prefix = "> "
	}

	// Slot label
	slotLabel := s.getSlotLabel(slot)
	ebitenutil.DebugPrintAt(panel, prefix+slotLabel, x, y)

	// Slot details on second line
	if !slot.Empty {
		details := s.getSlotDetails(slot)
		ebitenutil.DebugPrintAt(panel, "    "+details, x, y+15)
	}
}

// getSlotLabel returns the label for a slot.
func (s *SlotSelectionScreen) getSlotLabel(slot saveload.SlotInfo) string {
	if slot.IsAuto {
		if slot.Empty {
			return "Autosave: [Empty]"
		}
		return "Autosave"
	}

	if slot.Empty {
		return fmt.Sprintf("Slot %d: [Empty]", slot.Slot)
	}
	return fmt.Sprintf("Slot %d: %s", slot.Slot, slot.Summary.Genre)
}

// getSlotDetails returns details for a non-empty slot.
func (s *SlotSelectionScreen) getSlotDetails(slot saveload.SlotInfo) string {
	sum := slot.Summary
	savedAgo := time.Since(sum.SavedAt).Round(time.Minute)

	agoStr := formatDuration(savedAgo)
	return fmt.Sprintf("Day %d, Crew: %d, Saved: %s ago", sum.Day, sum.CrewCount, agoStr)
}

// formatDuration formats a duration in human-readable form.
func formatDuration(d time.Duration) string {
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}

// getTitle returns the title for the current mode.
func (s *SlotSelectionScreen) getTitle() string {
	if s.mode == SlotModeLoad {
		return "LOAD GAME"
	}
	return "SAVE GAME"
}

// getInstructions returns the control instructions.
func (s *SlotSelectionScreen) getInstructions() string {
	if s.confirmDelete {
		return "ENTER to confirm delete, ESC to cancel"
	}
	if s.mode == SlotModeLoad {
		return "UP/DOWN select, ENTER load, DEL delete, ESC back"
	}
	return "UP/DOWN select, ENTER save, DEL delete, ESC back"
}

// drawDeleteConfirm draws the delete confirmation dialog.
func (s *SlotSelectionScreen) drawDeleteConfirm(screen *ebiten.Image) {
	dialogWidth := 300
	dialogHeight := 80
	dialogX := (s.screenWidth - dialogWidth) / 2
	dialogY := (s.screenHeight - dialogHeight) / 2

	// Draw dialog background
	dialog := ebiten.NewImage(dialogWidth, dialogHeight)
	dialog.Fill(s.skin.PanelBackground)
	s.drawBorder(dialog)

	// Draw warning text
	slot := s.SelectedSlot()
	warning := fmt.Sprintf("Delete Slot %d?", slot)
	if slot == saveload.AutosaveSlot {
		warning = "Delete Autosave?"
	}
	warningX := (dialogWidth - len(warning)*7) / 2
	ebitenutil.DebugPrintAt(dialog, warning, warningX, 20)

	confirm := "Press ENTER to confirm"
	confirmX := (dialogWidth - len(confirm)*7) / 2
	ebitenutil.DebugPrintAt(dialog, confirm, confirmX, 45)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(dialogX), float64(dialogY))
	screen.DrawImage(dialog, op)
}

// drawBorder draws a border around the panel.
func (s *SlotSelectionScreen) drawBorder(panel *ebiten.Image) {
	DrawBorder(panel, s.skin)
}

// SlotAction represents the result of a slot screen interaction.
type SlotAction int

const (
	// SlotActionNone means no action taken.
	SlotActionNone SlotAction = iota
	// SlotActionSelect means a slot was selected for load/save.
	SlotActionSelect
	// SlotActionDelete means a slot should be deleted.
	SlotActionDelete
	// SlotActionCancel means the screen was cancelled.
	SlotActionCancel
)
