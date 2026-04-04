//go:build !headless

package ux

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/opd-ai/voyage/pkg/engine"
)

// MenuType identifies different menu screens.
type MenuType int

const (
	MenuMain MenuType = iota
	MenuPause
	MenuOptions
	MenuGameOver
)

// MenuItem represents a selectable menu option.
type MenuItem struct {
	ID      string
	Label   string
	Action  func()
	Enabled bool
}

// Menu manages menu screens and selection.
type Menu struct {
	skin          *UISkin
	genre         engine.GenreID
	menuType      MenuType
	items         []MenuItem
	selectedIndex int
	screenWidth   int
	screenHeight  int
}

// NewMenu creates a new menu instance.
func NewMenu(genre engine.GenreID, menuType MenuType, screenWidth, screenHeight int) *Menu {
	m := &Menu{
		skin:          DefaultSkin(genre),
		genre:         genre,
		menuType:      menuType,
		selectedIndex: 0,
		screenWidth:   screenWidth,
		screenHeight:  screenHeight,
	}
	m.initItems()
	return m
}

// SetGenre changes the menu's visual theme.
func (m *Menu) SetGenre(genre engine.GenreID) {
	m.genre = genre
	m.skin = DefaultSkin(genre)
}

// SetMenuType changes the menu type and reinitializes items.
func (m *Menu) SetMenuType(menuType MenuType) {
	m.menuType = menuType
	m.selectedIndex = 0
	m.initItems()
}

// initItems initializes menu items based on menu type.
func (m *Menu) initItems() {
	switch m.menuType {
	case MenuMain:
		m.items = []MenuItem{
			{ID: "new_game", Label: "New Journey", Enabled: true},
			{ID: "continue", Label: "Continue", Enabled: false},
			{ID: "options", Label: "Options", Enabled: true},
			{ID: "quit", Label: "Quit", Enabled: true},
		}
	case MenuPause:
		m.items = []MenuItem{
			{ID: "resume", Label: "Resume", Enabled: true},
			{ID: "save", Label: "Save Game", Enabled: true},
			{ID: "options", Label: "Options", Enabled: true},
			{ID: "main_menu", Label: "Main Menu", Enabled: true},
		}
	case MenuOptions:
		m.items = []MenuItem{
			{ID: "volume", Label: "Volume", Enabled: true},
			{ID: "controls", Label: "Controls", Enabled: true},
			{ID: "back", Label: "Back", Enabled: true},
		}
	case MenuGameOver:
		m.items = []MenuItem{
			{ID: "retry", Label: "Try Again", Enabled: true},
			{ID: "main_menu", Label: "Main Menu", Enabled: true},
		}
	}
}

// SelectNext moves selection to the next enabled item.
func (m *Menu) SelectNext() {
	for i := 0; i < len(m.items); i++ {
		m.selectedIndex = (m.selectedIndex + 1) % len(m.items)
		if m.items[m.selectedIndex].Enabled {
			return
		}
	}
}

// SelectPrev moves selection to the previous enabled item.
func (m *Menu) SelectPrev() {
	for i := 0; i < len(m.items); i++ {
		m.selectedIndex--
		if m.selectedIndex < 0 {
			m.selectedIndex = len(m.items) - 1
		}
		if m.items[m.selectedIndex].Enabled {
			return
		}
	}
}

// SelectedItem returns the currently selected item.
func (m *Menu) SelectedItem() *MenuItem {
	if m.selectedIndex >= 0 && m.selectedIndex < len(m.items) {
		return &m.items[m.selectedIndex]
	}
	return nil
}

// SelectedID returns the ID of the currently selected item.
func (m *Menu) SelectedID() string {
	item := m.SelectedItem()
	if item != nil {
		return item.ID
	}
	return ""
}

// SetItemEnabled enables or disables a menu item by ID.
func (m *Menu) SetItemEnabled(id string, enabled bool) {
	for i := range m.items {
		if m.items[i].ID == id {
			m.items[i].Enabled = enabled
			return
		}
	}
}

// Draw renders the menu to the screen.
func (m *Menu) Draw(screen *ebiten.Image) {
	// Draw background overlay for pause/options
	if m.menuType != MenuMain {
		overlay := ebiten.NewImage(m.screenWidth, m.screenHeight)
		overlay.Fill(m.skin.PanelBackground)
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleAlpha(0.5)
		screen.DrawImage(overlay, op)
	}

	// Calculate menu dimensions
	menuWidth := 300
	menuHeight := 60 + len(m.items)*40
	menuX := (m.screenWidth - menuWidth) / 2
	menuY := (m.screenHeight - menuHeight) / 2

	// Draw menu panel
	panel := ebiten.NewImage(menuWidth, menuHeight)
	panel.Fill(m.skin.PanelBackground)
	m.drawBorder(panel)

	// Draw title
	title := m.getTitle()
	titleX := (menuWidth - len(title)*7) / 2
	ebitenutil.DebugPrintAt(panel, title, titleX, 16)

	// Draw items
	y := 50
	for i, item := range m.items {
		prefix := "  "
		if i == m.selectedIndex {
			prefix = "> "
		}

		label := prefix + item.Label
		if !item.Enabled {
			label = prefix + "[" + item.Label + "]"
		}

		ebitenutil.DebugPrintAt(panel, label, 40, y)
		y += 30
	}

	// Draw panel to screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(menuX), float64(menuY))
	screen.DrawImage(panel, op)

	// Draw instructions at bottom
	instructions := m.getInstructions()
	instrX := (m.screenWidth - len(instructions)*7) / 2
	ebitenutil.DebugPrintAt(screen, instructions, instrX, m.screenHeight-30)
}

// getTitle returns the title for the current menu type.
func (m *Menu) getTitle() string {
	switch m.menuType {
	case MenuMain:
		return m.skin.MenuTitle
	case MenuPause:
		return "PAUSED"
	case MenuOptions:
		return "OPTIONS"
	case MenuGameOver:
		return "JOURNEY'S END"
	default:
		return "MENU"
	}
}

// getInstructions returns the instructions for the current menu.
func (m *Menu) getInstructions() string {
	return "UP/DOWN to select, ENTER to confirm"
}

// drawBorder draws a border around the panel.
func (m *Menu) drawBorder(panel *ebiten.Image) {
	w, h := panel.Bounds().Dx(), panel.Bounds().Dy()
	c := m.skin.PanelBorder

	for x := 0; x < w; x++ {
		panel.Set(x, 0, c)
		panel.Set(x, 1, c)
		panel.Set(x, h-1, c)
		panel.Set(x, h-2, c)
	}
	for y := 0; y < h; y++ {
		panel.Set(0, y, c)
		panel.Set(1, y, c)
		panel.Set(w-1, y, c)
		panel.Set(w-2, y, c)
	}
}

// DrawGameOverScreen draws the game over screen with stats.
func (m *Menu) DrawGameOverScreen(screen *ebiten.Image, stats GameStats) {
	m.Draw(screen)

	// Draw stats above the menu
	statsY := m.screenHeight/2 - 150
	statsX := m.screenWidth/2 - 100

	lines := []string{
		fmt.Sprintf("Days Traveled: %d", stats.DaysTraveled),
		fmt.Sprintf("Distance: %d tiles", stats.DistanceTraveled),
		fmt.Sprintf("Crew Lost: %d", stats.CrewLost),
		fmt.Sprintf("Events Faced: %d", stats.EventsResolved),
	}

	for i, line := range lines {
		ebitenutil.DebugPrintAt(screen, line, statsX, statsY+i*20)
	}
}

// GameStats holds end-of-game statistics.
type GameStats struct {
	DaysTraveled     int
	DistanceTraveled int
	CrewLost         int
	EventsResolved   int
	Victory          bool
}
