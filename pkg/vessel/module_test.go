package vessel

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewModule(t *testing.T) {
	tests := []struct {
		name       string
		moduleType ModuleType
		wantMax    float64
	}{
		{"engine", ModuleEngine, 50},
		{"cargo", ModuleCargoHold, 40},
		{"medical", ModuleMedicalBay, 30},
		{"navigation", ModuleNavigation, 30},
		{"defense", ModuleDefense, 60},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModule(tt.moduleType)
			if m.MaxIntegrity() != tt.wantMax {
				t.Errorf("max integrity = %f, want %f", m.MaxIntegrity(), tt.wantMax)
			}
			if m.IsDisabled() {
				t.Error("new module should not be disabled")
			}
			if m.IntegrityRatio() != 1.0 {
				t.Error("new module should have full integrity")
			}
			if m.Efficiency() != 1.0 {
				t.Error("new module should have full efficiency")
			}
		})
	}
}

func TestModuleDamageAndRepair(t *testing.T) {
	m := NewModule(ModuleEngine)
	initial := m.Integrity()

	// Take damage
	disabled := m.TakeDamage(20)
	if disabled {
		t.Error("module should not be disabled with 20 damage")
	}
	if m.Integrity() != initial-20 {
		t.Errorf("integrity = %f, want %f", m.Integrity(), initial-20)
	}

	// Efficiency should be reduced
	if m.Efficiency() >= 1.0 {
		t.Error("damaged module should have reduced efficiency")
	}

	// Repair
	m.Repair(10)
	if m.Integrity() != initial-10 {
		t.Errorf("integrity after repair = %f, want %f", m.Integrity(), initial-10)
	}

	// Full repair
	m.RepairFull()
	if m.Integrity() != m.MaxIntegrity() {
		t.Error("module should be at full integrity after RepairFull")
	}
	if m.Efficiency() != 1.0 {
		t.Error("module should have full efficiency after RepairFull")
	}
}

func TestModuleDisabling(t *testing.T) {
	m := NewModule(ModuleMedicalBay)

	// Deal lethal damage
	disabled := m.TakeDamage(m.MaxIntegrity() + 10)
	if !disabled {
		t.Error("module should be disabled")
	}
	if m.Integrity() != 0 {
		t.Error("disabled module should have 0 integrity")
	}
	if !m.IsDisabled() {
		t.Error("IsDisabled should return true")
	}
	if m.Efficiency() != 0 {
		t.Error("disabled module should have 0 efficiency")
	}
}

func TestModuleTier(t *testing.T) {
	m := NewModule(ModuleEngine)
	baseMI := m.MaxIntegrity()

	// Upgrade tier
	m.SetTier(2)
	if m.Tier() != 2 {
		t.Errorf("tier = %d, want 2", m.Tier())
	}
	expectedMax := baseMI * 1.25 // 25% increase per tier
	if m.MaxIntegrity() != expectedMax {
		t.Errorf("max integrity at tier 2 = %f, want %f", m.MaxIntegrity(), expectedMax)
	}

	// Max tier is 5
	m.SetTier(10)
	if m.Tier() != 5 {
		t.Errorf("tier should cap at 5, got %d", m.Tier())
	}

	// Min tier is 1
	m.SetTier(0)
	if m.Tier() != 1 {
		t.Errorf("tier should min at 1, got %d", m.Tier())
	}
}

func TestModuleCondition(t *testing.T) {
	m := NewModule(ModuleEngine)

	if GetModuleCondition(m) != ModuleConditionPristine {
		t.Error("new module should be pristine")
	}

	m.TakeDamage(10) // 80% integrity
	if GetModuleCondition(m) != ModuleConditionOperational {
		t.Error("80% should be Operational")
	}

	m.TakeDamage(15) // 50% integrity
	if GetModuleCondition(m) != ModuleConditionDamaged {
		t.Error("50% should be Damaged")
	}

	m.TakeDamage(15) // 20% integrity
	if GetModuleCondition(m) != ModuleConditionDisabled {
		t.Error("20% should be Disabled")
	}
}

func TestNewModuleSystem(t *testing.T) {
	ms := NewModuleSystem(engine.GenreFantasy)

	// All modules should be present
	for _, mt := range AllModuleTypes() {
		if m := ms.GetModule(mt); m == nil {
			t.Errorf("missing module type: %d", mt)
		}
	}

	// Check genre
	if ms.Genre() != engine.GenreFantasy {
		t.Errorf("genre = %s, want fantasy", ms.Genre())
	}
}

func TestModuleSystemGenreSwitching(t *testing.T) {
	ms := NewModuleSystem(engine.GenreFantasy)

	ms.SetGenre(engine.GenreScifi)
	if ms.Genre() != engine.GenreScifi {
		t.Errorf("genre = %s, want scifi", ms.Genre())
	}
}

func TestModuleSystemEfficiencies(t *testing.T) {
	ms := NewModuleSystem(engine.GenreFantasy)

	// All new modules should be at full efficiency
	if ms.EngineEfficiency() != 1.0 {
		t.Errorf("engine efficiency = %f, want 1.0", ms.EngineEfficiency())
	}
	if ms.MedicalEfficiency() != 1.0 {
		t.Errorf("medical efficiency = %f, want 1.0", ms.MedicalEfficiency())
	}
	if ms.NavigationAccuracy() != 1.0 {
		t.Errorf("navigation accuracy = %f, want 1.0", ms.NavigationAccuracy())
	}

	// Damage engine
	ms.GetModule(ModuleEngine).TakeDamage(25) // 50% integrity
	if ms.EngineEfficiency() >= 1.0 {
		t.Error("damaged engine should have reduced efficiency")
	}
}

func TestModuleSystemCargoCapacity(t *testing.T) {
	ms := NewModuleSystem(engine.GenreFantasy)

	// Full health tier 1 should give 1.0 multiplier
	mult := ms.CargoCapacityMultiplier()
	if mult != 1.0 {
		t.Errorf("cargo multiplier = %f, want 1.0", mult)
	}

	// Upgrade cargo hold
	ms.UpgradeModule(ModuleCargoHold)
	mult = ms.CargoCapacityMultiplier()
	if mult <= 1.0 {
		t.Error("upgraded cargo should have higher multiplier")
	}
}

func TestModuleSystemDefense(t *testing.T) {
	ms := NewModuleSystem(engine.GenreFantasy)

	rating := ms.DefenseRating()
	if rating <= 0 {
		t.Error("defense rating should be positive")
	}

	// Upgrade defense
	ms.UpgradeModule(ModuleDefense)
	newRating := ms.DefenseRating()
	if newRating <= rating {
		t.Error("upgraded defense should have higher rating")
	}
}

func TestModuleSystemIntegrity(t *testing.T) {
	ms := NewModuleSystem(engine.GenreFantasy)

	totalMax := ms.TotalMaxIntegrity()
	total := ms.TotalIntegrity()

	if total != totalMax {
		t.Error("new system should have full integrity")
	}
	if ms.OverallIntegrityRatio() != 1.0 {
		t.Error("new system should have 100% integrity ratio")
	}

	// Damage a module
	ms.GetModule(ModuleEngine).TakeDamage(25)
	if ms.TotalIntegrity() >= totalMax {
		t.Error("damaged system should have reduced total integrity")
	}
	if ms.OverallIntegrityRatio() >= 1.0 {
		t.Error("damaged system should have reduced integrity ratio")
	}
}

func TestModuleSystemDistributeDamage(t *testing.T) {
	ms := NewModuleSystem(engine.GenreFantasy)
	initialTotal := ms.TotalIntegrity()

	// Distribute 10 damage
	rng := func(n int) int { return 0 } // Always pick first module
	ms.DistributeDamage(10, rng)

	if ms.TotalIntegrity() >= initialTotal {
		t.Error("total integrity should decrease after damage")
	}
}

func TestModuleSystemRepairModule(t *testing.T) {
	ms := NewModuleSystem(engine.GenreFantasy)
	m := ms.GetModule(ModuleEngine)
	m.TakeDamage(20)

	integrityBefore := m.Integrity()
	ms.RepairModule(ModuleEngine, 10)

	if m.Integrity() != integrityBefore+10 {
		t.Errorf("integrity = %f, want %f", m.Integrity(), integrityBefore+10)
	}
}

func TestModuleSystemUpgradeModule(t *testing.T) {
	ms := NewModuleSystem(engine.GenreFantasy)
	m := ms.GetModule(ModuleEngine)

	if m.Tier() != 1 {
		t.Errorf("initial tier = %d, want 1", m.Tier())
	}

	if !ms.UpgradeModule(ModuleEngine) {
		t.Error("should be able to upgrade from tier 1")
	}
	if m.Tier() != 2 {
		t.Errorf("upgraded tier = %d, want 2", m.Tier())
	}

	// Upgrade to max
	m.SetTier(5)
	if ms.UpgradeModule(ModuleEngine) {
		t.Error("should not be able to upgrade beyond tier 5")
	}
}

func TestModuleTypeNames(t *testing.T) {
	genres := engine.AllGenres()
	types := AllModuleTypes()

	for _, g := range genres {
		for _, mt := range types {
			name := ModuleTypeName(mt, g)
			if name == "" {
				t.Errorf("missing name for genre=%s, module=%d", g, mt)
			}
		}
	}
}

func TestModuleConditionNames(t *testing.T) {
	conditions := []ModuleCondition{
		ModuleConditionPristine,
		ModuleConditionOperational,
		ModuleConditionDamaged,
		ModuleConditionCritical,
		ModuleConditionDisabled,
	}

	for _, mc := range conditions {
		name := ModuleConditionName(mc)
		if name == "" || name == "Unknown" {
			t.Errorf("missing or invalid name for condition %d", mc)
		}
	}
}

func TestAllModuleTypes(t *testing.T) {
	types := AllModuleTypes()
	if len(types) != 5 {
		t.Errorf("expected 5 module types, got %d", len(types))
	}

	expected := map[ModuleType]bool{
		ModuleEngine:     false,
		ModuleCargoHold:  false,
		ModuleMedicalBay: false,
		ModuleNavigation: false,
		ModuleDefense:    false,
	}
	for _, mt := range types {
		expected[mt] = true
	}
	for mt, found := range expected {
		if !found {
			t.Errorf("module type %d not in AllModuleTypes", mt)
		}
	}
}

func TestModuleName(t *testing.T) {
	m := NewModule(ModuleEngine)

	name := m.Name(engine.GenreFantasy)
	if name != "Stable" {
		t.Errorf("fantasy engine name = %s, want Stable", name)
	}

	name = m.Name(engine.GenreScifi)
	if name != "Engine Room" {
		t.Errorf("scifi engine name = %s, want Engine Room", name)
	}
}

func TestGetAllModules(t *testing.T) {
	ms := NewModuleSystem(engine.GenreFantasy)
	modules := ms.GetAllModules()

	if len(modules) != 5 {
		t.Errorf("expected 5 modules, got %d", len(modules))
	}
}

func TestModuleEfficiencyCalculation(t *testing.T) {
	m := NewModule(ModuleEngine)

	// At 100% integrity, efficiency should be 1.0
	if m.Efficiency() != 1.0 {
		t.Errorf("efficiency at 100%% = %f, want 1.0", m.Efficiency())
	}

	// At 50% integrity
	m.TakeDamage(m.MaxIntegrity() * 0.5)
	// efficiency = 0.5 * (0.4 + 0.6*0.5) = 0.5 * 0.7 = 0.35
	expectedEff := 0.5 * (0.4 + 0.6*0.5)
	if m.Efficiency() != expectedEff {
		t.Errorf("efficiency at 50%% = %f, want %f", m.Efficiency(), expectedEff)
	}

	// At 0% integrity
	m.TakeDamage(m.MaxIntegrity())
	if m.Efficiency() != 0 {
		t.Errorf("efficiency at 0%% = %f, want 0", m.Efficiency())
	}
}

func TestNewModuleWithTier(t *testing.T) {
	m := NewModuleWithTier(ModuleEngine, 3)

	if m.Tier() != 3 {
		t.Errorf("tier = %d, want 3", m.Tier())
	}

	// Tier 3 should have 50% more integrity than base
	expectedMax := DefaultModuleStats[ModuleEngine].MaxIntegrity * 1.5
	if m.MaxIntegrity() != expectedMax {
		t.Errorf("max integrity = %f, want %f", m.MaxIntegrity(), expectedMax)
	}
}
