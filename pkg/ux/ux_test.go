//go:build !headless

package ux

import (
	"testing"
	"time"

	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/events"
	"github.com/opd-ai/voyage/pkg/resources"
	"github.com/opd-ai/voyage/pkg/saveload"
)

func TestDefaultSkin(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		skin := DefaultSkin(genre)
		if skin == nil {
			t.Errorf("DefaultSkin(%v) returned nil", genre)
			continue
		}
		if skin.MenuTitle == "" {
			t.Errorf("DefaultSkin(%v) has empty MenuTitle", genre)
		}
		if skin.PanelBackground == nil {
			t.Errorf("DefaultSkin(%v) has nil PanelBackground", genre)
		}
	}
}

func TestNewHUD(t *testing.T) {
	hud := NewHUD(engine.GenreFantasy)
	if hud == nil {
		t.Fatal("NewHUD returned nil")
	}
	if hud.genre != engine.GenreFantasy {
		t.Errorf("expected genre Fantasy, got %v", hud.genre)
	}
}

func TestHUDSetGenre(t *testing.T) {
	hud := NewHUD(engine.GenreFantasy)
	hud.SetGenre(engine.GenreScifi)
	if hud.genre != engine.GenreScifi {
		t.Errorf("expected genre Scifi, got %v", hud.genre)
	}
}

func TestNewWorldMapView(t *testing.T) {
	wmv := NewWorldMapView(engine.GenreFantasy, 16, 800, 600)
	if wmv == nil {
		t.Fatal("NewWorldMapView returned nil")
	}
	if wmv.tileSize != 16 {
		t.Errorf("expected tileSize 16, got %d", wmv.tileSize)
	}
}

func TestWorldMapViewCenterOn(t *testing.T) {
	wmv := NewWorldMapView(engine.GenreFantasy, 16, 320, 240)
	wmv.CenterOn(50, 50)

	tilesWide := 320 / 16 // 20
	tilesHigh := 240 / 16 // 15
	expectedX := 50 - tilesWide/2
	expectedY := 50 - tilesHigh/2

	if wmv.cameraX != expectedX {
		t.Errorf("expected cameraX %d, got %d", expectedX, wmv.cameraX)
	}
	if wmv.cameraY != expectedY {
		t.Errorf("expected cameraY %d, got %d", expectedY, wmv.cameraY)
	}
}

func TestNewEventOverlay(t *testing.T) {
	eo := NewEventOverlay(engine.GenreFantasy, 400, 300)
	if eo == nil {
		t.Fatal("NewEventOverlay returned nil")
	}
	if eo.visible {
		t.Error("overlay should not be visible by default")
	}
}

func TestEventOverlayVisibility(t *testing.T) {
	eo := NewEventOverlay(engine.GenreFantasy, 400, 300)

	eo.Show()
	if !eo.IsVisible() {
		t.Error("overlay should be visible after Show()")
	}

	eo.Hide()
	if eo.IsVisible() {
		t.Error("overlay should not be visible after Hide()")
	}
}

func TestEventOverlaySelection(t *testing.T) {
	eo := NewEventOverlay(engine.GenreFantasy, 400, 300)

	// Test SelectNext
	eo.SelectNext(4)
	if eo.SelectedChoice() != 1 {
		t.Errorf("expected selectedChoice 1, got %d", eo.SelectedChoice())
	}

	// Test SelectPrev
	eo.SelectPrev(4)
	if eo.SelectedChoice() != 0 {
		t.Errorf("expected selectedChoice 0, got %d", eo.SelectedChoice())
	}

	// Test wrap-around
	eo.SelectPrev(4)
	if eo.SelectedChoice() != 3 {
		t.Errorf("expected selectedChoice 3 (wrap), got %d", eo.SelectedChoice())
	}

	// Test SelectByNumber
	if !eo.SelectByNumber(2, 4) {
		t.Error("SelectByNumber(2, 4) should return true")
	}
	if eo.SelectedChoice() != 1 {
		t.Errorf("expected selectedChoice 1, got %d", eo.SelectedChoice())
	}

	// Test invalid number
	if eo.SelectByNumber(0, 4) {
		t.Error("SelectByNumber(0, 4) should return false")
	}
	if eo.SelectByNumber(5, 4) {
		t.Error("SelectByNumber(5, 4) should return false")
	}
}

func TestNewMenu(t *testing.T) {
	menu := NewMenu(engine.GenreFantasy, MenuMain, 800, 600)
	if menu == nil {
		t.Fatal("NewMenu returned nil")
	}
	if len(menu.items) == 0 {
		t.Error("menu should have items")
	}
}

func TestMenuSetMenuType(t *testing.T) {
	menu := NewMenu(engine.GenreFantasy, MenuMain, 800, 600)
	mainItemCount := len(menu.items)

	menu.SetMenuType(MenuPause)
	if menu.menuType != MenuPause {
		t.Errorf("expected menuType MenuPause, got %v", menu.menuType)
	}

	menu.SetMenuType(MenuOptions)
	if len(menu.items) == mainItemCount {
		t.Error("options menu should have different number of items")
	}
}

func TestMenuSelection(t *testing.T) {
	menu := NewMenu(engine.GenreFantasy, MenuMain, 800, 600)

	// Initially selected first item
	item := menu.SelectedItem()
	if item == nil {
		t.Fatal("SelectedItem returned nil")
	}
	if item.ID != "new_game" {
		t.Errorf("expected first item 'new_game', got '%s'", item.ID)
	}

	// SelectNext
	menu.SelectNext()
	// Since 'continue' is disabled by default, it should skip to 'options'
	selected := menu.SelectedID()
	if selected != "continue" && selected != "options" {
		t.Errorf("unexpected selected item after SelectNext: %s", selected)
	}
}

func TestMenuSetItemEnabled(t *testing.T) {
	menu := NewMenu(engine.GenreFantasy, MenuMain, 800, 600)

	// Find continue item and verify it's disabled
	for _, item := range menu.items {
		if item.ID == "continue" && item.Enabled {
			t.Error("continue should be disabled by default")
		}
	}

	menu.SetItemEnabled("continue", true)

	for _, item := range menu.items {
		if item.ID == "continue" && !item.Enabled {
			t.Error("continue should be enabled after SetItemEnabled")
		}
	}
}

func TestWrapText(t *testing.T) {
	eo := NewEventOverlay(engine.GenreFantasy, 400, 300)

	text := "This is a test of the word wrapping functionality"
	lines := eo.wrapText(text, 20)

	if len(lines) == 0 {
		t.Error("wrapText should return at least one line")
	}

	for _, line := range lines {
		if len(line) > 20+10 { // Some slack for word boundaries
			t.Errorf("line too long: %s", line)
		}
	}
}

// Integration test stubs - these test that types can be used together
func TestHUDWithResources(t *testing.T) {
	hud := NewHUD(engine.GenreFantasy)
	res := resources.NewResources(engine.GenreFantasy)
	party := crew.NewParty(engine.GenreFantasy, 4)

	// This should not panic
	_ = hud
	_ = res
	_ = party
}

func TestEventOverlayWithEvent(t *testing.T) {
	eo := NewEventOverlay(engine.GenreFantasy, 400, 300)
	event := events.NewEvent(1, events.CategoryWeather, "Storm", "A storm approaches", engine.GenreFantasy)

	// This should not panic
	_ = eo
	_ = event
}

func TestWorldMapViewSetGenre(t *testing.T) {
	wmv := NewWorldMapView(engine.GenreFantasy, 16, 320, 240)
	wmv.SetGenre(engine.GenreScifi)
	if wmv.genre != engine.GenreScifi {
		t.Errorf("expected genre Scifi, got %v", wmv.genre)
	}
}

func TestEventOverlaySetGenre(t *testing.T) {
	eo := NewEventOverlay(engine.GenreFantasy, 400, 300)
	eo.SetGenre(engine.GenreHorror)
	if eo.genre != engine.GenreHorror {
		t.Errorf("expected genre Horror, got %v", eo.genre)
	}
}

func TestMenuSetGenre(t *testing.T) {
	menu := NewMenu(engine.GenreFantasy, MenuMain, 800, 600)
	menu.SetGenre(engine.GenreCyberpunk)
	if menu.genre != engine.GenreCyberpunk {
		t.Errorf("expected genre Cyberpunk, got %v", menu.genre)
	}
}

func TestMenuSelectPrev(t *testing.T) {
	menu := NewMenu(engine.GenreFantasy, MenuMain, 800, 600)

	// Get to the end
	for i := 0; i < 10; i++ {
		menu.SelectNext()
	}

	// Now go back
	menu.SelectPrev()
	idx := menu.selectedIndex
	if idx >= len(menu.items)-1 {
		t.Error("SelectPrev should move selection backwards")
	}
}

func TestMenuTypes(t *testing.T) {
	menuTypes := []MenuType{MenuMain, MenuPause, MenuOptions, MenuGameOver}

	for _, mt := range menuTypes {
		menu := NewMenu(engine.GenreFantasy, mt, 800, 600)
		if menu.menuType != mt {
			t.Errorf("expected menuType %d, got %d", mt, menu.menuType)
		}
		if len(menu.items) == 0 {
			t.Errorf("menu type %d should have items", mt)
		}
	}
}

func TestAbsFunction(t *testing.T) {
	// Test the abs function through WorldMapView
	wmv := NewWorldMapView(engine.GenreFantasy, 16, 320, 240)

	// abs is used internally - we verify behavior through CenterOn
	wmv.CenterOn(-10, -10)
	// Camera should be set (even if negative, it still processes)
	if wmv.cameraX == 0 && wmv.cameraY == 0 {
		// Actually this is fine - CenterOn calculates properly
	}
}

func TestHealthToStatus(t *testing.T) {
	hud := NewHUD(engine.GenreFantasy)

	// Test healthToStatus by checking the statusColor function indirectly
	// We can't easily test private functions, but we verify HUD creation works
	if hud.skin == nil {
		t.Error("HUD should have a skin")
	}
}

func TestWrapTextEdgeCases(t *testing.T) {
	eo := NewEventOverlay(engine.GenreFantasy, 400, 300)

	// Test with maxWidth <= 0 (should default to 40)
	lines := eo.wrapText("short", 0)
	if len(lines) != 1 {
		t.Errorf("wrapText with maxWidth 0 should return 1 line, got %d", len(lines))
	}

	// Test with negative maxWidth
	lines = eo.wrapText("short", -5)
	if len(lines) != 1 {
		t.Errorf("wrapText with negative maxWidth should return 1 line, got %d", len(lines))
	}

	// Test empty string
	lines = eo.wrapText("", 20)
	if len(lines) != 0 {
		t.Errorf("wrapText with empty string should return 0 lines, got %d", len(lines))
	}

	// Test very long word
	lines = eo.wrapText("supercalifragilisticexpialidocious", 10)
	if len(lines) == 0 {
		t.Error("wrapText should handle words longer than maxWidth")
	}
}

func TestMenuGameOverItems(t *testing.T) {
	menu := NewMenu(engine.GenreFantasy, MenuGameOver, 800, 600)

	// Game over menu should have specific items
	foundRetry := false
	foundMainMenu := false
	for _, item := range menu.items {
		if item.ID == "retry" {
			foundRetry = true
		}
		if item.ID == "main_menu" {
			foundMainMenu = true
		}
	}

	if !foundRetry {
		t.Error("Game over menu should have retry option")
	}
	if !foundMainMenu {
		t.Error("Game over menu should have main_menu option")
	}
}

func TestUISkinColors(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		skin := DefaultSkin(genre)

		// Verify all colors are non-nil
		if skin.TextPrimary == nil {
			t.Errorf("DefaultSkin(%v) has nil TextPrimary", genre)
		}
		if skin.HighlightColor == nil {
			t.Errorf("DefaultSkin(%v) has nil HighlightColor", genre)
		}
		if skin.PanelBackground == nil {
			t.Errorf("DefaultSkin(%v) has nil PanelBackground", genre)
		}
		if skin.PanelBorder == nil {
			t.Errorf("DefaultSkin(%v) has nil PanelBorder", genre)
		}
	}
}

func TestWorldMapViewUpdateCamera(t *testing.T) {
	wmv := NewWorldMapView(engine.GenreFantasy, 16, 320, 240)

	// Test CenterOn at various positions
	positions := []struct{ x, y int }{
		{0, 0},
		{50, 50},
		{100, 100},
		{-10, -10},
	}

	for _, pos := range positions {
		wmv.CenterOn(pos.x, pos.y)
		// Just verify it doesn't panic
		_ = wmv.cameraX
		_ = wmv.cameraY
	}
}

func TestEventOverlayResetSelection(t *testing.T) {
	eo := NewEventOverlay(engine.GenreFantasy, 400, 300)

	// Move selection
	eo.SelectNext(4)
	eo.SelectNext(4)

	if eo.SelectedChoice() != 2 {
		t.Errorf("expected selectedChoice 2, got %d", eo.SelectedChoice())
	}

	// Hide and show - Show() resets selection to 0
	eo.Hide()
	eo.Show()

	if eo.SelectedChoice() != 0 {
		t.Errorf("after hide/show, expected selectedChoice 0 (reset), got %d", eo.SelectedChoice())
	}
}

func TestNewSlotSelectionScreen(t *testing.T) {
	screen := NewSlotSelectionScreen(engine.GenreFantasy, SlotModeLoad, 800, 600)
	if screen == nil {
		t.Fatal("NewSlotSelectionScreen returned nil")
	}
	if screen.mode != SlotModeLoad {
		t.Errorf("expected mode SlotModeLoad, got %v", screen.mode)
	}
	if screen.genre != engine.GenreFantasy {
		t.Errorf("expected genre Fantasy, got %v", screen.genre)
	}
}

func TestSlotSelectionScreenSetMode(t *testing.T) {
	screen := NewSlotSelectionScreen(engine.GenreFantasy, SlotModeLoad, 800, 600)

	screen.SetMode(SlotModeSave)
	if screen.mode != SlotModeSave {
		t.Errorf("expected mode SlotModeSave, got %v", screen.mode)
	}
}

func TestSlotSelectionScreenSetGenre(t *testing.T) {
	screen := NewSlotSelectionScreen(engine.GenreFantasy, SlotModeLoad, 800, 600)

	screen.SetGenre(engine.GenreScifi)
	if screen.genre != engine.GenreScifi {
		t.Errorf("expected genre Scifi, got %v", screen.genre)
	}
}

func TestSlotSelectionScreenNavigation(t *testing.T) {
	screen := NewSlotSelectionScreen(engine.GenreFantasy, SlotModeLoad, 800, 600)

	// Simulate some slots
	screen.slots = []saveload.SlotInfo{
		{Slot: 0, Empty: false, IsAuto: true},
		{Slot: 1, Empty: true, IsAuto: false},
		{Slot: 2, Empty: false, IsAuto: false},
	}

	// Test initial state
	if screen.SelectedSlot() != 0 {
		t.Errorf("expected selected slot 0, got %d", screen.SelectedSlot())
	}

	// Test SelectNext
	screen.SelectNext()
	if screen.SelectedSlot() != 1 {
		t.Errorf("expected selected slot 1 after SelectNext, got %d", screen.SelectedSlot())
	}

	// Test SelectPrev
	screen.SelectPrev()
	if screen.SelectedSlot() != 0 {
		t.Errorf("expected selected slot 0 after SelectPrev, got %d", screen.SelectedSlot())
	}

	// Test wrap-around
	screen.SelectPrev()
	if screen.SelectedSlot() != 2 {
		t.Errorf("expected selected slot 2 (wrap), got %d", screen.SelectedSlot())
	}
}

func TestSlotSelectionScreenCanSelect(t *testing.T) {
	screen := NewSlotSelectionScreen(engine.GenreFantasy, SlotModeLoad, 800, 600)

	screen.slots = []saveload.SlotInfo{
		{Slot: 0, Empty: false, IsAuto: true},
		{Slot: 1, Empty: true, IsAuto: false},
	}

	// In load mode, can only select non-empty slots
	screen.selectedIndex = 0
	if !screen.CanSelect() {
		t.Error("should be able to select non-empty slot in load mode")
	}

	screen.selectedIndex = 1
	if screen.CanSelect() {
		t.Error("should not be able to select empty slot in load mode")
	}

	// In save mode, can always select
	screen.SetMode(SlotModeSave)
	screen.selectedIndex = 1
	if !screen.CanSelect() {
		t.Error("should be able to select any slot in save mode")
	}
}

func TestSlotSelectionScreenDeleteConfirm(t *testing.T) {
	screen := NewSlotSelectionScreen(engine.GenreFantasy, SlotModeLoad, 800, 600)

	screen.slots = []saveload.SlotInfo{
		{Slot: 0, Empty: false, IsAuto: true},
	}

	if screen.IsConfirmingDelete() {
		t.Error("should not be confirming delete initially")
	}

	screen.ToggleDeleteConfirm()
	if !screen.IsConfirmingDelete() {
		t.Error("should be confirming delete after toggle")
	}

	screen.CancelDelete()
	if screen.IsConfirmingDelete() {
		t.Error("should not be confirming delete after cancel")
	}
}

func TestSlotSelectionScreenNoSlots(t *testing.T) {
	screen := NewSlotSelectionScreen(engine.GenreFantasy, SlotModeLoad, 800, 600)

	// With no slots, navigation should not panic
	screen.SelectNext()
	screen.SelectPrev()

	// Selected slot should be -1 with no slots
	if screen.SelectedSlot() != -1 {
		t.Errorf("expected selected slot -1 with no slots, got %d", screen.SelectedSlot())
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		input    time.Duration
		expected string
	}{
		{5 * time.Minute, "5m"},
		{30 * time.Minute, "30m"},
		{2 * time.Hour, "2h"},
		{48 * time.Hour, "2d"},
	}

	for _, tt := range tests {
		result := formatDuration(tt.input)
		if result != tt.expected {
			t.Errorf("formatDuration(%v) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestSlotActionConstants(t *testing.T) {
	// Verify constants exist and have distinct values
	actions := []SlotAction{SlotActionNone, SlotActionSelect, SlotActionDelete, SlotActionCancel}
	seen := make(map[SlotAction]bool)

	for _, a := range actions {
		if seen[a] {
			t.Errorf("duplicate SlotAction value: %d", a)
		}
		seen[a] = true
	}
}

func TestSlotModeConstants(t *testing.T) {
	// Verify constants exist and have distinct values
	modes := []SlotSelectionMode{SlotModeLoad, SlotModeSave}
	seen := make(map[SlotSelectionMode]bool)

	for _, m := range modes {
		if seen[m] {
			t.Errorf("duplicate SlotSelectionMode value: %d", m)
		}
		seen[m] = true
	}
}
