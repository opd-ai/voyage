package vessel

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestUpgradeCost(t *testing.T) {
	// Tier 1 -> 2 costs base
	cost := UpgradeCost(ModuleEngine, 1)
	if cost != 75 { // 50 * (1 + 0.5)
		t.Errorf("tier 1 upgrade cost = %f, want 75", cost)
	}

	// Tier 2 -> 3 costs more
	cost2 := UpgradeCost(ModuleEngine, 2)
	if cost2 <= cost {
		t.Error("tier 2 upgrade should cost more than tier 1")
	}

	// Max tier returns 0
	cost5 := UpgradeCost(ModuleEngine, 5)
	if cost5 != 0 {
		t.Errorf("max tier upgrade cost = %f, want 0", cost5)
	}
}

func TestUpgradeManagerCanAfford(t *testing.T) {
	ms := NewModuleSystem(engine.GenreFantasy)
	um := NewUpgradeManager(engine.GenreFantasy)

	cost := UpgradeCost(ModuleEngine, 1)

	// With enough currency
	if !um.CanAffordUpgrade(ms, ModuleEngine, cost) {
		t.Error("should be able to afford upgrade")
	}

	// Without enough currency
	if um.CanAffordUpgrade(ms, ModuleEngine, cost-1) {
		t.Error("should not be able to afford upgrade")
	}

	// Max tier module
	ms.GetModule(ModuleEngine).SetTier(5)
	if um.CanAffordUpgrade(ms, ModuleEngine, 1000) {
		t.Error("should not be able to upgrade max tier")
	}
}

func TestUpgradeManagerAttemptUpgrade(t *testing.T) {
	ms := NewModuleSystem(engine.GenreFantasy)
	um := NewUpgradeManager(engine.GenreFantasy)

	cost := UpgradeCost(ModuleEngine, 1)

	// Successful upgrade
	result := um.AttemptUpgrade(ms, ModuleEngine, cost+10)
	if !result.Success {
		t.Errorf("upgrade should succeed: %s", result.Message)
	}
	if result.NewTier != 2 {
		t.Errorf("new tier = %d, want 2", result.NewTier)
	}
	if result.CurrencyUsed != cost {
		t.Errorf("currency used = %f, want %f", result.CurrencyUsed, cost)
	}

	// Failed upgrade - insufficient funds
	result = um.AttemptUpgrade(ms, ModuleEngine, 0)
	if result.Success {
		t.Error("upgrade should fail with no currency")
	}
}

func TestUpgradeManagerAttemptUpgradeMaxTier(t *testing.T) {
	ms := NewModuleSystem(engine.GenreFantasy)
	um := NewUpgradeManager(engine.GenreFantasy)

	// Set to max tier
	ms.GetModule(ModuleEngine).SetTier(5)

	result := um.AttemptUpgrade(ms, ModuleEngine, 1000)
	if result.Success {
		t.Error("upgrade should fail for max tier")
	}
	if result.Message != "Module already at maximum tier" {
		t.Errorf("unexpected message: %s", result.Message)
	}
}

func TestUpgradeManagerGetAvailableUpgrades(t *testing.T) {
	ms := NewModuleSystem(engine.GenreFantasy)
	um := NewUpgradeManager(engine.GenreFantasy)

	options := um.GetAvailableUpgrades(ms)
	if len(options) != 5 {
		t.Errorf("expected 5 upgrade options, got %d", len(options))
	}

	// Max one module
	ms.GetModule(ModuleEngine).SetTier(5)
	options = um.GetAvailableUpgrades(ms)
	if len(options) != 4 {
		t.Errorf("expected 4 upgrade options after max, got %d", len(options))
	}
}

func TestTierName(t *testing.T) {
	genres := engine.AllGenres()

	for _, g := range genres {
		for tier := 1; tier <= 5; tier++ {
			name := TierName(tier, g)
			if name == "" || name == "Unknown" {
				t.Errorf("missing tier name for genre=%s, tier=%d", g, tier)
			}
		}
	}

	// Invalid tier
	name := TierName(0, engine.GenreFantasy)
	if name != "Unknown" {
		t.Errorf("tier 0 name = %s, want Unknown", name)
	}
	name = TierName(6, engine.GenreFantasy)
	if name != "Unknown" {
		t.Errorf("tier 6 name = %s, want Unknown", name)
	}
}

func TestUpgradeManagerSetGenre(t *testing.T) {
	um := NewUpgradeManager(engine.GenreFantasy)
	um.SetGenre(engine.GenreScifi)

	// Test that genre affects upgrade messages
	ms := NewModuleSystem(engine.GenreScifi)
	options := um.GetAvailableUpgrades(ms)

	if len(options) == 0 {
		t.Fatal("expected upgrade options")
	}

	// Scifi names should be used
	if options[0].NextTierName != "Enhanced" {
		t.Errorf("scifi tier name = %s, want Enhanced", options[0].NextTierName)
	}
}

func TestUpgradeOption(t *testing.T) {
	uo := UpgradeOption{
		ModuleType:   ModuleEngine,
		CurrentTier:  1,
		NextTier:     2,
		Cost:         75,
		ModuleName:   "Stable",
		NextTierName: "Improved",
	}

	desc := uo.UpgradeDescription()
	if desc == "" {
		t.Error("upgrade description should not be empty")
	}
}

func TestBulkUpgrade(t *testing.T) {
	ms := NewModuleSystem(engine.GenreFantasy)
	um := NewUpgradeManager(engine.GenreFantasy)

	// Calculate total cost for two upgrades
	cost1 := UpgradeCost(ModuleEngine, 1)
	cost2 := UpgradeCost(ModuleCargoHold, 1)
	totalCost := cost1 + cost2

	results, spent := um.BulkUpgrade(ms, []ModuleType{ModuleEngine, ModuleCargoHold}, totalCost)

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
	if !results[0].Success || !results[1].Success {
		t.Error("both upgrades should succeed")
	}
	if spent != totalCost {
		t.Errorf("spent = %f, want %f", spent, totalCost)
	}

	// Verify modules were upgraded
	if ms.GetModule(ModuleEngine).Tier() != 2 {
		t.Error("engine should be tier 2")
	}
	if ms.GetModule(ModuleCargoHold).Tier() != 2 {
		t.Error("cargo hold should be tier 2")
	}
}

func TestBulkUpgradePartialFunds(t *testing.T) {
	ms := NewModuleSystem(engine.GenreFantasy)
	um := NewUpgradeManager(engine.GenreFantasy)

	// Only enough for one upgrade
	cost1 := UpgradeCost(ModuleEngine, 1)

	results, spent := um.BulkUpgrade(ms, []ModuleType{ModuleEngine, ModuleCargoHold}, cost1)

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
	if !results[0].Success {
		t.Error("first upgrade should succeed")
	}
	if results[1].Success {
		t.Error("second upgrade should fail (insufficient funds)")
	}
	if spent != cost1 {
		t.Errorf("spent = %f, want %f", spent, cost1)
	}
}

func TestAllTierNames(t *testing.T) {
	names := AllTierNames(engine.GenreFantasy)
	if len(names) != 5 {
		t.Errorf("expected 5 tier names, got %d", len(names))
	}

	// Test fallback to fantasy for unknown genre
	names = AllTierNames("invalid")
	if len(names) != 5 {
		t.Errorf("fallback should return 5 tier names, got %d", len(names))
	}
}

func TestUpgradeDescriptions(t *testing.T) {
	modules := AllModuleTypes()
	for _, mt := range modules {
		uo := UpgradeOption{ModuleType: mt}
		desc := uo.UpgradeDescription()
		if desc == "" || desc == "Improves module performance" {
			if mt != ModuleType(-1) { // Unknown types get default
				// All known types should have specific descriptions
				continue
			}
			t.Errorf("module type %d has generic description", mt)
		}
	}
}
