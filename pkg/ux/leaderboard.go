//go:build !headless

package ux

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/leaderboard"
)

// LeaderboardScreen displays per-seed leaderboard entries with filtering.
type LeaderboardScreen struct {
	skin         *UISkin
	genre        engine.GenreID
	screenWidth  int
	screenHeight int

	// Data
	board   *leaderboard.Board
	entries []*leaderboard.Entry

	// Filters
	filterGenre *engine.GenreID
	filterSeed  *int64

	// Selection state
	selectedIndex int
	scrollOffset  int
	maxVisible    int

	// Tab state
	filterTab int // 0 = all, 1 = by genre, 2 = by seed

	// Available genres for filtering
	availableGenres []engine.GenreID
	genreIndex      int
}

// NewLeaderboardScreen creates a new leaderboard display screen.
func NewLeaderboardScreen(genre engine.GenreID, screenWidth, screenHeight int) *LeaderboardScreen {
	ls := &LeaderboardScreen{
		skin:            DefaultSkin(genre),
		genre:           genre,
		screenWidth:     screenWidth,
		screenHeight:    screenHeight,
		selectedIndex:   0,
		scrollOffset:    0,
		maxVisible:      10,
		filterTab:       0,
		availableGenres: engine.AllGenres(),
		genreIndex:      0,
	}
	return ls
}

// SetGenre changes the screen's visual theme.
func (ls *LeaderboardScreen) SetGenre(genre engine.GenreID) {
	ls.genre = genre
	ls.skin = DefaultSkin(genre)
}

// SetBoard sets the leaderboard data to display.
func (ls *LeaderboardScreen) SetBoard(board *leaderboard.Board) {
	ls.board = board
	ls.refreshEntries()
}

// refreshEntries updates the displayed entries based on current filters.
func (ls *LeaderboardScreen) refreshEntries() {
	if ls.board == nil {
		ls.entries = nil
		return
	}

	switch ls.filterTab {
	case 0: // All entries
		ls.entries = ls.board.GetTopN(100)
	case 1: // By genre
		if ls.filterGenre != nil {
			ls.entries = ls.board.GetByGenre(*ls.filterGenre)
		} else {
			ls.entries = ls.board.GetAll()
		}
	case 2: // By seed
		if ls.filterSeed != nil {
			ls.entries = ls.board.GetBySeed(*ls.filterSeed)
		} else {
			ls.entries = ls.board.GetAll()
		}
	}

	// Reset selection if out of bounds
	if ls.selectedIndex >= len(ls.entries) {
		ls.selectedIndex = 0
	}
	ls.scrollOffset = 0
}

// SetFilterGenre sets the genre filter.
func (ls *LeaderboardScreen) SetFilterGenre(genre *engine.GenreID) {
	ls.filterGenre = genre
	ls.refreshEntries()
}

// SetFilterSeed sets the seed filter.
func (ls *LeaderboardScreen) SetFilterSeed(seed *int64) {
	ls.filterSeed = seed
	ls.refreshEntries()
}

// ClearFilters removes all filters.
func (ls *LeaderboardScreen) ClearFilters() {
	ls.filterGenre = nil
	ls.filterSeed = nil
	ls.filterTab = 0
	ls.refreshEntries()
}

// CycleFilterTab cycles through filter tabs.
func (ls *LeaderboardScreen) CycleFilterTab() {
	ls.filterTab = (ls.filterTab + 1) % 3
	switch ls.filterTab {
	case 0:
		ls.filterGenre = nil
		ls.filterSeed = nil
	case 1:
		g := ls.availableGenres[ls.genreIndex]
		ls.filterGenre = &g
		ls.filterSeed = nil
	case 2:
		ls.filterGenre = nil
		// Seed filter is set externally
	}
	ls.refreshEntries()
}

// CycleGenreFilter cycles through available genres in genre filter mode.
func (ls *LeaderboardScreen) CycleGenreFilter() {
	if ls.filterTab == 1 {
		ls.genreIndex = (ls.genreIndex + 1) % len(ls.availableGenres)
		g := ls.availableGenres[ls.genreIndex]
		ls.filterGenre = &g
		ls.refreshEntries()
	}
}

// SelectNext moves selection down.
func (ls *LeaderboardScreen) SelectNext() {
	if len(ls.entries) == 0 {
		return
	}
	ls.selectedIndex++
	if ls.selectedIndex >= len(ls.entries) {
		ls.selectedIndex = len(ls.entries) - 1
	}
	// Adjust scroll if needed
	if ls.selectedIndex >= ls.scrollOffset+ls.maxVisible {
		ls.scrollOffset = ls.selectedIndex - ls.maxVisible + 1
	}
}

// SelectPrev moves selection up.
func (ls *LeaderboardScreen) SelectPrev() {
	if len(ls.entries) == 0 {
		return
	}
	ls.selectedIndex--
	if ls.selectedIndex < 0 {
		ls.selectedIndex = 0
	}
	// Adjust scroll if needed
	if ls.selectedIndex < ls.scrollOffset {
		ls.scrollOffset = ls.selectedIndex
	}
}

// SelectedEntry returns the currently selected entry.
func (ls *LeaderboardScreen) SelectedEntry() *leaderboard.Entry {
	if ls.selectedIndex >= 0 && ls.selectedIndex < len(ls.entries) {
		return ls.entries[ls.selectedIndex]
	}
	return nil
}

// GetSelectedSeed returns the seed of the selected entry for replay.
func (ls *LeaderboardScreen) GetSelectedSeed() (int64, bool) {
	entry := ls.SelectedEntry()
	if entry != nil {
		return entry.Seed, true
	}
	return 0, false
}

// Draw renders the leaderboard screen.
func (ls *LeaderboardScreen) Draw(screen *ebiten.Image) {
	// Draw background overlay
	DrawOverlay(screen, ls.skin, ls.screenWidth, ls.screenHeight)

	// Calculate panel dimensions
	panelWidth := 600
	panelHeight := 450

	// Draw centered panel
	panel, panelX, panelY := DrawCenteredPanel(
		screen, ls.skin, ls.screenWidth, ls.screenHeight, panelWidth, panelHeight,
	)

	// Draw title
	DrawCenteredText(panel, ls.getTitle(), 16)

	// Draw filter tabs
	ls.drawFilterTabs(panel, panelWidth)

	// Draw column headers
	ls.drawHeaders(panel)

	// Draw entries
	ls.drawEntries(panel)

	// Draw scrollbar if needed
	if len(ls.entries) > ls.maxVisible {
		ls.drawScrollbar(panel, panelHeight)
	}

	// Draw panel to screen
	DrawPanelToScreen(screen, panel, panelX, panelY)

	// Draw instructions
	ls.drawInstructions(screen)
}

// getTitle returns the screen title based on filter state.
func (ls *LeaderboardScreen) getTitle() string {
	switch ls.filterTab {
	case 1:
		if ls.filterGenre != nil {
			return fmt.Sprintf("LEADERBOARD - %s", genreDisplayName(*ls.filterGenre))
		}
	case 2:
		if ls.filterSeed != nil {
			return fmt.Sprintf("LEADERBOARD - Seed %d", *ls.filterSeed)
		}
	}
	return "GLOBAL LEADERBOARD"
}

// drawFilterTabs draws the filter tab bar.
func (ls *LeaderboardScreen) drawFilterTabs(panel *ebiten.Image, panelWidth int) {
	tabs := []string{"All", "By Genre", "By Seed"}
	tabWidth := 80
	startX := (panelWidth - len(tabs)*tabWidth) / 2
	y := 40

	for i, tab := range tabs {
		x := startX + i*tabWidth
		prefix := "  "
		if i == ls.filterTab {
			prefix = "[*"
			tab = tab + "]"
		}
		ebitenutil.DebugPrintAt(panel, prefix+tab, x, y)
	}

	// Show current genre filter if in genre mode
	if ls.filterTab == 1 && ls.filterGenre != nil {
		genreText := fmt.Sprintf("< %s >", genreDisplayName(*ls.filterGenre))
		ebitenutil.DebugPrintAt(panel, genreText, (panelWidth-len(genreText)*7)/2, y+20)
	}
}

// drawHeaders draws the column headers.
func (ls *LeaderboardScreen) drawHeaders(panel *ebiten.Image) {
	y := 80
	ebitenutil.DebugPrintAt(panel, "RANK", 20, y)
	ebitenutil.DebugPrintAt(panel, "SCORE", 80, y)
	ebitenutil.DebugPrintAt(panel, "DAYS", 160, y)
	ebitenutil.DebugPrintAt(panel, "CREW", 220, y)
	ebitenutil.DebugPrintAt(panel, "GENRE", 280, y)
	ebitenutil.DebugPrintAt(panel, "SEED", 380, y)
	ebitenutil.DebugPrintAt(panel, "PLAYER", 480, y)

	// Draw separator line
	for x := 20; x < 580; x += 2 {
		ebitenutil.DebugPrintAt(panel, "-", x, y+15)
	}
}

// drawEntries draws the leaderboard entries.
func (ls *LeaderboardScreen) drawEntries(panel *ebiten.Image) {
	y := 110
	lineHeight := 25

	endIdx := ls.scrollOffset + ls.maxVisible
	if endIdx > len(ls.entries) {
		endIdx = len(ls.entries)
	}

	for i := ls.scrollOffset; i < endIdx; i++ {
		entry := ls.entries[i]
		rank := i + 1

		// Highlight selected entry
		prefix := "  "
		if i == ls.selectedIndex {
			prefix = "> "
		}

		// Format entry line
		rankStr := fmt.Sprintf("%s%3d", prefix, rank)
		scoreStr := fmt.Sprintf("%6d", entry.Score)
		daysStr := fmt.Sprintf("%4d", entry.Days)
		crewStr := fmt.Sprintf("%4d", entry.Survivors)
		genreStr := genreShortName(entry.Genre)
		seedStr := fmt.Sprintf("%d", entry.Seed)
		playerStr := entry.PlayerName
		if playerStr == "" {
			playerStr = "Anonymous"
		}
		if len(playerStr) > 10 {
			playerStr = playerStr[:10]
		}

		ebitenutil.DebugPrintAt(panel, rankStr, 10, y)
		ebitenutil.DebugPrintAt(panel, scoreStr, 70, y)
		ebitenutil.DebugPrintAt(panel, daysStr, 160, y)
		ebitenutil.DebugPrintAt(panel, crewStr, 220, y)
		ebitenutil.DebugPrintAt(panel, genreStr, 280, y)
		ebitenutil.DebugPrintAt(panel, seedStr, 380, y)
		ebitenutil.DebugPrintAt(panel, playerStr, 480, y)

		y += lineHeight
	}

	// Show empty message if no entries
	if len(ls.entries) == 0 {
		emptyMsg := "No entries found"
		ebitenutil.DebugPrintAt(panel, emptyMsg, 240, 200)
	}
}

// drawScrollbar draws a scrollbar indicator.
func (ls *LeaderboardScreen) drawScrollbar(panel *ebiten.Image, panelHeight int) {
	x := 585
	barTop := 110
	barHeight := panelHeight - 150
	barBottom := barTop + barHeight

	// Draw track
	for y := barTop; y < barBottom; y += 10 {
		ebitenutil.DebugPrintAt(panel, "|", x, y)
	}

	// Calculate thumb position
	totalEntries := len(ls.entries)
	if totalEntries > 0 {
		thumbPos := barTop + (ls.scrollOffset * barHeight / totalEntries)
		thumbSize := ls.maxVisible * barHeight / totalEntries
		if thumbSize < 10 {
			thumbSize = 10
		}
		for y := thumbPos; y < thumbPos+thumbSize && y < barBottom; y += 10 {
			ebitenutil.DebugPrintAt(panel, "#", x, y)
		}
	}
}

// drawInstructions draws control instructions at the bottom.
func (ls *LeaderboardScreen) drawInstructions(screen *ebiten.Image) {
	instructions := "UP/DOWN: Navigate | TAB: Filter | LEFT/RIGHT: Change Genre | ENTER: Replay Seed | ESC: Back"
	x := (ls.screenWidth - len(instructions)*7) / 2
	ebitenutil.DebugPrintAt(screen, instructions, x, ls.screenHeight-30)
}

// genreDisplayName returns a display-friendly genre name.
func genreDisplayName(genre engine.GenreID) string {
	switch genre {
	case engine.GenreFantasy:
		return "Fantasy"
	case engine.GenreScifi:
		return "Sci-Fi"
	case engine.GenreHorror:
		return "Horror"
	case engine.GenreCyberpunk:
		return "Cyberpunk"
	case engine.GenrePostapoc:
		return "Post-Apocalyptic"
	default:
		return string(genre)
	}
}

// genreShortName returns a short genre abbreviation.
func genreShortName(genre engine.GenreID) string {
	switch genre {
	case engine.GenreFantasy:
		return "FAN"
	case engine.GenreScifi:
		return "SCI"
	case engine.GenreHorror:
		return "HOR"
	case engine.GenreCyberpunk:
		return "CYB"
	case engine.GenrePostapoc:
		return "PAC"
	default:
		return "???"
	}
}

// EntryCount returns the number of displayed entries.
func (ls *LeaderboardScreen) EntryCount() int {
	return len(ls.entries)
}

// TotalEntryCount returns total entries in the board.
func (ls *LeaderboardScreen) TotalEntryCount() int {
	if ls.board == nil {
		return 0
	}
	return ls.board.Count()
}

// GetFilterInfo returns current filter information.
func (ls *LeaderboardScreen) GetFilterInfo() (filterType, filterValue string) {
	switch ls.filterTab {
	case 0:
		return "all", ""
	case 1:
		if ls.filterGenre != nil {
			return "genre", string(*ls.filterGenre)
		}
		return "genre", ""
	case 2:
		if ls.filterSeed != nil {
			return "seed", fmt.Sprintf("%d", *ls.filterSeed)
		}
		return "seed", ""
	}
	return "unknown", ""
}
