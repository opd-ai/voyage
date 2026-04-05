//go:build headless

package ux

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/vessel"
)

// TestLoadoutScreenCreation tests LoadoutScreen initialization.
func TestLoadoutScreenCreation(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		ls := NewLoadoutScreen(genre, 800, 600)
		if ls == nil {
			t.Errorf("NewLoadoutScreen(%v) returned nil", genre)
			continue
		}
		if ls.genre != genre {
			t.Errorf("expected genre %v, got %v", genre, ls.genre)
		}
		if ls.screenWidth != 800 {
			t.Errorf("expected width 800, got %d", ls.screenWidth)
		}
		if ls.screenHeight != 600 {
			t.Errorf("expected height 600, got %d", ls.screenHeight)
		}
		if ls.PointsRemaining() != DefaultStartingPoints {
			t.Errorf("expected %d points, got %d", DefaultStartingPoints, ls.PointsRemaining())
		}
	}
}

// TestLoadoutScreenVisibility tests Show/Hide functionality.
func TestLoadoutScreenVisibility(t *testing.T) {
	ls := NewLoadoutScreen(engine.GenreFantasy, 800, 600)

	if ls.IsVisible() {
		t.Error("LoadoutScreen should not be visible initially")
	}

	ls.Show()
	if !ls.IsVisible() {
		t.Error("LoadoutScreen should be visible after Show()")
	}

	ls.Hide()
	if ls.IsVisible() {
		t.Error("LoadoutScreen should not be visible after Hide()")
	}
}

// TestLoadoutScreenSetGenre tests genre switching.
func TestLoadoutScreenSetGenre(t *testing.T) {
	ls := NewLoadoutScreen(engine.GenreFantasy, 800, 600)

	ls.SetGenre(engine.GenreScifi)
	if ls.genre != engine.GenreScifi {
		t.Errorf("expected genre Scifi, got %v", ls.genre)
	}
}

// TestLoadoutScreenNavigation tests slot navigation.
func TestLoadoutScreenNavigation(t *testing.T) {
	ls := NewLoadoutScreen(engine.GenreFantasy, 800, 600)
	numModules := len(vessel.AllModuleTypes())

	// Initial selection should be 0
	if ls.selectedSlot != 0 {
		t.Errorf("expected selectedSlot 0, got %d", ls.selectedSlot)
	}

	// SelectNext should move forward
	ls.SelectNext()
	if ls.selectedSlot != 1 {
		t.Errorf("expected selectedSlot 1 after SelectNext, got %d", ls.selectedSlot)
	}

	// SelectPrev should move backward
	ls.SelectPrev()
	if ls.selectedSlot != 0 {
		t.Errorf("expected selectedSlot 0 after SelectPrev, got %d", ls.selectedSlot)
	}

	// SelectPrev from 0 should wrap to last
	ls.SelectPrev()
	if ls.selectedSlot != numModules-1 {
		t.Errorf("expected selectedSlot %d (wrap), got %d", numModules-1, ls.selectedSlot)
	}

	// SelectNext from last should wrap to 0
	ls.SelectNext()
	if ls.selectedSlot != 0 {
		t.Errorf("expected selectedSlot 0 (wrap), got %d", ls.selectedSlot)
	}
}

// TestLoadoutScreenUpgradeDowngrade tests module upgrade/downgrade.
func TestLoadoutScreenUpgradeDowngrade(t *testing.T) {
	ls := NewLoadoutScreen(engine.GenreFantasy, 800, 600)
	initialPoints := ls.PointsRemaining()

	// Upgrade should succeed and cost a point
	if !ls.UpgradeSelected() {
		t.Error("UpgradeSelected should succeed initially")
	}
	if ls.PointsRemaining() != initialPoints-1 {
		t.Errorf("expected %d points after upgrade, got %d", initialPoints-1, ls.PointsRemaining())
	}

	// Downgrade should succeed and refund a point
	if !ls.DowngradeSelected() {
		t.Error("DowngradeSelected should succeed after upgrade")
	}
	if ls.PointsRemaining() != initialPoints {
		t.Errorf("expected %d points after downgrade, got %d", initialPoints, ls.PointsRemaining())
	}

	// Downgrade at tier 1 should fail
	if ls.DowngradeSelected() {
		t.Error("DowngradeSelected should fail at tier 1")
	}
}

// TestLoadoutScreenUpgradeLimit tests upgrade limit enforcement.
func TestLoadoutScreenUpgradeLimit(t *testing.T) {
	ls := NewLoadoutScreen(engine.GenreFantasy, 800, 600)

	// Upgrade to max tier (3 for starting modules)
	for i := 0; i < 3; i++ {
		ls.UpgradeSelected()
	}

	// Further upgrade should fail at max tier
	prevPoints := ls.PointsRemaining()
	if ls.UpgradeSelected() {
		t.Error("UpgradeSelected should fail at max tier")
	}
	if ls.PointsRemaining() != prevPoints {
		t.Error("Points should not change when upgrade fails")
	}
}

// TestLoadoutScreenReset tests reset functionality.
func TestLoadoutScreenReset(t *testing.T) {
	ls := NewLoadoutScreen(engine.GenreFantasy, 800, 600)

	// Modify state
	ls.UpgradeSelected()
	ls.SelectNext()
	ls.SelectNext()

	// Reset should restore defaults
	ls.Reset()

	if ls.PointsRemaining() != DefaultStartingPoints {
		t.Errorf("expected %d points after reset, got %d", DefaultStartingPoints, ls.PointsRemaining())
	}
	if ls.selectedSlot != 0 {
		t.Errorf("expected selectedSlot 0 after reset, got %d", ls.selectedSlot)
	}
}

// TestLoadoutScreenGetModuleSystem tests module system retrieval.
func TestLoadoutScreenGetModuleSystem(t *testing.T) {
	ls := NewLoadoutScreen(engine.GenreFantasy, 800, 600)

	ms := ls.GetModuleSystem()
	if ms == nil {
		t.Fatal("GetModuleSystem returned nil")
	}

	// Verify module system has all module types
	for _, mt := range vessel.AllModuleTypes() {
		m := ms.GetModule(mt)
		if m == nil {
			t.Errorf("module type %v not found in system", mt)
		}
	}
}

// TestLoadoutConfiguration tests configuration save/restore.
func TestLoadoutConfiguration(t *testing.T) {
	ls := NewLoadoutScreen(engine.GenreFantasy, 800, 600)

	// Upgrade some modules
	ls.UpgradeSelected() // Engine tier 2
	ls.SelectNext()
	ls.UpgradeSelected() // Cargo tier 2

	// Get configuration
	cfg := ls.GetConfiguration()
	if cfg.EngineTier != 2 {
		t.Errorf("expected EngineTier 2, got %d", cfg.EngineTier)
	}
	if cfg.CargoTier != 2 {
		t.Errorf("expected CargoTier 2, got %d", cfg.CargoTier)
	}

	// Reset and apply configuration
	ls.Reset()
	ls.ApplyConfiguration(cfg)

	newCfg := ls.GetConfiguration()
	if newCfg.EngineTier != cfg.EngineTier {
		t.Errorf("EngineTier mismatch after apply: expected %d, got %d", cfg.EngineTier, newCfg.EngineTier)
	}
	if newCfg.CargoTier != cfg.CargoTier {
		t.Errorf("CargoTier mismatch after apply: expected %d, got %d", cfg.CargoTier, newCfg.CargoTier)
	}
}

// TestCargoScreenCreation tests CargoScreen initialization.
func TestCargoScreenCreation(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		cs := NewCargoScreen(genre, 800, 600)
		if cs == nil {
			t.Errorf("NewCargoScreen(%v) returned nil", genre)
			continue
		}
		if cs.genre != genre {
			t.Errorf("expected genre %v, got %v", genre, cs.genre)
		}
		if cs.IsVisible() {
			t.Errorf("CargoScreen should not be visible initially for %v", genre)
		}
	}
}

// TestCargoScreenVisibility tests Show/Hide functionality.
func TestCargoScreenVisibility(t *testing.T) {
	cs := NewCargoScreen(engine.GenreFantasy, 800, 600)

	if cs.IsVisible() {
		t.Error("CargoScreen should not be visible initially")
	}

	cs.Show()
	if !cs.IsVisible() {
		t.Error("CargoScreen should be visible after Show()")
	}
	if cs.scrollOffset != 0 {
		t.Error("scrollOffset should reset to 0 on Show()")
	}

	cs.Hide()
	if cs.IsVisible() {
		t.Error("CargoScreen should not be visible after Hide()")
	}
}

// TestCargoScreenSetGenre tests genre switching.
func TestCargoScreenSetGenre(t *testing.T) {
	cs := NewCargoScreen(engine.GenreFantasy, 800, 600)

	cs.SetGenre(engine.GenreCyberpunk)
	if cs.genre != engine.GenreCyberpunk {
		t.Errorf("expected genre Cyberpunk, got %v", cs.genre)
	}
}

// TestCargoScreenScrolling tests scroll functionality.
func TestCargoScreenScrolling(t *testing.T) {
	cs := NewCargoScreen(engine.GenreFantasy, 800, 600)
	cs.Show()

	// Initial scroll should be 0
	if cs.scrollOffset != 0 {
		t.Errorf("expected scrollOffset 0, got %d", cs.scrollOffset)
	}

	// ScrollDown should increase offset
	cs.ScrollDown()
	if cs.scrollOffset != 1 {
		t.Errorf("expected scrollOffset 1 after ScrollDown, got %d", cs.scrollOffset)
	}

	// ScrollUp should decrease offset
	cs.ScrollUp()
	if cs.scrollOffset != 0 {
		t.Errorf("expected scrollOffset 0 after ScrollUp, got %d", cs.scrollOffset)
	}

	// ScrollUp at 0 should stay at 0
	cs.ScrollUp()
	if cs.scrollOffset != 0 {
		t.Errorf("scrollOffset should not go negative, got %d", cs.scrollOffset)
	}
}

// TestGetCargoSummary tests cargo summary generation.
func TestGetCargoSummary(t *testing.T) {
	hold := vessel.NewCargoHoldWithTier(1)

	// Add some cargo
	hold.AddWithVolume("Food", 10, 5, 2, vessel.CargoSupplies)

	summary := GetCargoSummary(hold)

	if summary.TotalItems != 1 {
		t.Errorf("expected 1 item, got %d", summary.TotalItems)
	}
	if summary.TotalWeight != 20 { // 10 weight * 2 quantity
		t.Errorf("expected weight 20, got %d", summary.TotalWeight)
	}
	if summary.TotalVolume != 10 { // 5 volume * 2 quantity
		t.Errorf("expected volume 10, got %d", summary.TotalVolume)
	}
	if summary.Tier != 1 {
		t.Errorf("expected tier 1, got %d", summary.Tier)
	}
	if summary.CategoryCounts[vessel.CargoSupplies] != 2 {
		t.Errorf("expected 2 supplies, got %d", summary.CategoryCounts[vessel.CargoSupplies])
	}
}

// TestTryUpgradeModule tests the shared upgrade logic.
func TestTryUpgradeModule(t *testing.T) {
	ms := vessel.NewModuleSystem(engine.GenreFantasy)
	points := 5

	// Upgrade should succeed
	if !TryUpgradeModule(ms, 0, &points) {
		t.Error("TryUpgradeModule should succeed")
	}
	if points != 4 {
		t.Errorf("expected 4 points, got %d", points)
	}

	// Upgrade with 0 points should fail
	points = 0
	if TryUpgradeModule(ms, 1, &points) {
		t.Error("TryUpgradeModule should fail with 0 points")
	}
}

// TestTryDowngradeModule tests the shared downgrade logic.
func TestTryDowngradeModule(t *testing.T) {
	ms := vessel.NewModuleSystem(engine.GenreFantasy)
	points := 5

	// Upgrade first
	TryUpgradeModule(ms, 0, &points)
	// Now points = 4

	// Downgrade should succeed
	if !TryDowngradeModule(ms, 0, &points) {
		t.Error("TryDowngradeModule should succeed")
	}
	if points != 5 {
		t.Errorf("expected 5 points after downgrade, got %d", points)
	}

	// Downgrade at tier 1 should fail
	if TryDowngradeModule(ms, 0, &points) {
		t.Error("TryDowngradeModule should fail at tier 1")
	}
}

// TestDefaultStartingPoints verifies the constant value.
func TestDefaultStartingPoints(t *testing.T) {
	if DefaultStartingPoints != 5 {
		t.Errorf("expected DefaultStartingPoints 5, got %d", DefaultStartingPoints)
	}
}

// TestLoadoutConfigurationDefaults tests default configuration values.
func TestLoadoutConfigurationDefaults(t *testing.T) {
	cfg := LoadoutConfiguration{}

	// Zero values are fine as defaults
	if cfg.EngineTier != 0 {
		t.Errorf("expected EngineTier 0 default, got %d", cfg.EngineTier)
	}
	if cfg.CargoTier != 0 {
		t.Errorf("expected CargoTier 0 default, got %d", cfg.CargoTier)
	}
}

// TestCargoSummaryEmpty tests summary with empty cargo hold.
func TestCargoSummaryEmpty(t *testing.T) {
	hold := vessel.NewCargoHoldWithTier(1)
	summary := GetCargoSummary(hold)

	if summary.TotalItems != 0 {
		t.Errorf("expected 0 items, got %d", summary.TotalItems)
	}
	if summary.TotalWeight != 0 {
		t.Errorf("expected weight 0, got %d", summary.TotalWeight)
	}
	if summary.TotalVolume != 0 {
		t.Errorf("expected volume 0, got %d", summary.TotalVolume)
	}
}
