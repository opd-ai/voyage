//go:build !headless

package ux

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/events"
	"github.com/opd-ai/voyage/pkg/resources"
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

func TestGameStats(t *testing.T) {
	stats := GameStats{
		DaysTraveled:     30,
		DistanceTraveled: 150,
		CrewLost:         2,
		EventsResolved:   25,
		Victory:          true,
	}

	if stats.DaysTraveled != 30 {
		t.Errorf("expected DaysTraveled 30, got %d", stats.DaysTraveled)
	}
	if !stats.Victory {
		t.Error("expected Victory true")
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
