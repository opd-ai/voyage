package vessel

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewSpecialization(t *testing.T) {
	s := NewSpecialization()

	if s.Type != SpecNone {
		t.Errorf("type = %d, want SpecNone", s.Type)
	}
	if s.Level != 0 {
		t.Errorf("level = %d, want 0", s.Level)
	}
	if s.SpeedBonus != 1.0 {
		t.Errorf("speed bonus = %f, want 1.0", s.SpeedBonus)
	}
	if s.CargoBonus != 1.0 {
		t.Errorf("cargo bonus = %f, want 1.0", s.CargoBonus)
	}
	if s.DefBonus != 1.0 {
		t.Errorf("defense bonus = %f, want 1.0", s.DefBonus)
	}
}

func TestSpecializationSetType(t *testing.T) {
	s := NewSpecialization()
	s.SetType(SpecSpeed)

	if s.Type != SpecSpeed {
		t.Errorf("type = %d, want SpecSpeed", s.Type)
	}
}

func TestSpecializationUpgrade(t *testing.T) {
	s := NewSpecialization()

	// Can't upgrade without type
	if s.Upgrade() {
		t.Error("should not upgrade without type set")
	}

	s.SetType(SpecSpeed)

	// Upgrade 3 times
	for i := 0; i < 3; i++ {
		if !s.Upgrade() {
			t.Errorf("upgrade %d should succeed", i+1)
		}
	}

	if s.Level != 3 {
		t.Errorf("level = %d, want 3", s.Level)
	}

	// Can't upgrade beyond level 3
	if s.Upgrade() {
		t.Error("should not upgrade beyond level 3")
	}
}

func TestSpecializationBonuses(t *testing.T) {
	tests := []struct {
		specType    SpecializationType
		expectSpeed bool // Higher speed bonus
		expectCargo bool // Higher cargo bonus
		expectDef   bool // Higher defense bonus
	}{
		{SpecSpeed, true, false, false},
		{SpecCargo, false, true, false},
		{SpecDefense, false, false, true},
	}

	for _, tt := range tests {
		s := NewSpecialization()
		s.SetType(tt.specType)
		s.Upgrade() // Level 1

		base := 1.0

		if tt.expectSpeed && s.SpeedBonus <= base {
			t.Errorf("%d speed bonus should be > 1.0, got %f", tt.specType, s.SpeedBonus)
		}
		if tt.expectCargo && s.CargoBonus <= base {
			t.Errorf("%d cargo bonus should be > 1.0, got %f", tt.specType, s.CargoBonus)
		}
		if tt.expectDef && s.DefBonus <= base {
			t.Errorf("%d defense bonus should be > 1.0, got %f", tt.specType, s.DefBonus)
		}
	}
}

func TestSpecializedModule(t *testing.T) {
	sm := NewSpecializedModule(ModuleEngine)

	if sm.Module == nil {
		t.Fatal("module should not be nil")
	}
	if sm.spec == nil {
		t.Fatal("specialization should not be nil")
	}

	// Set specialization
	sm.SetSpecialization(SpecSpeed)
	if sm.Specialization().Type != SpecSpeed {
		t.Error("specialization type should be Speed")
	}

	// Upgrade
	if !sm.UpgradeSpecialization() {
		t.Error("should be able to upgrade specialization")
	}
	if sm.Specialization().Level != 1 {
		t.Errorf("specialization level = %d, want 1", sm.Specialization().Level)
	}
}

func TestSpecializedModuleEffectiveBonuses(t *testing.T) {
	sm := NewSpecializedModule(ModuleEngine)
	sm.SetSpecialization(SpecSpeed)
	sm.UpgradeSpecialization()

	// At full health, bonuses should be spec bonuses
	speedBonus := sm.EffectiveSpeedBonus()
	if speedBonus <= 1.0 {
		t.Errorf("effective speed bonus = %f, should be > 1.0", speedBonus)
	}

	// Damage module
	sm.TakeDamage(sm.MaxIntegrity() * 0.5)

	// Effective bonuses should be reduced
	damagedSpeedBonus := sm.EffectiveSpeedBonus()
	if damagedSpeedBonus >= speedBonus {
		t.Error("damaged module should have lower effective bonus")
	}
}

func TestSpecializationNames(t *testing.T) {
	genres := engine.AllGenres()
	types := AllSpecializationTypes()
	types = append(types, SpecNone)

	for _, g := range genres {
		for _, st := range types {
			name := SpecializationName(st, g)
			if name == "" {
				t.Errorf("missing name for genre=%s, spec=%d", g, st)
			}
		}
	}
}

func TestSpecializationDescription(t *testing.T) {
	types := append(AllSpecializationTypes(), SpecNone)

	for _, st := range types {
		desc := SpecializationDescription(st)
		if desc == "" || desc == "Unknown specialization" {
			t.Errorf("missing description for spec type %d", st)
		}
	}
}

func TestModuleSpecializationCost(t *testing.T) {
	// Test cost increases with level
	cost1 := ModuleSpecializationCost(ModuleEngine, 0)
	cost2 := ModuleSpecializationCost(ModuleEngine, 1)
	cost3 := ModuleSpecializationCost(ModuleEngine, 2)

	if cost2 <= cost1 {
		t.Error("level 1 cost should be higher than level 0")
	}
	if cost3 <= cost2 {
		t.Error("level 2 cost should be higher than level 1")
	}
}

func TestGetRecommendedSpecialization(t *testing.T) {
	// Engine should recommend Speed
	if GetRecommendedSpecialization(ModuleEngine) != SpecSpeed {
		t.Error("engine should recommend Speed")
	}

	// Cargo should recommend Cargo
	if GetRecommendedSpecialization(ModuleCargoHold) != SpecCargo {
		t.Error("cargo hold should recommend Cargo")
	}

	// Defense should recommend Defense
	if GetRecommendedSpecialization(ModuleDefense) != SpecDefense {
		t.Error("defense should recommend Defense")
	}
}

func TestCanSpecialize(t *testing.T) {
	m := NewModule(ModuleEngine)

	// Tier 1 cannot specialize
	if CanSpecialize(m) {
		t.Error("tier 1 module should not be able to specialize")
	}

	// Tier 2 can specialize
	m.SetTier(2)
	if !CanSpecialize(m) {
		t.Error("tier 2 module should be able to specialize")
	}
}

func TestAllSpecializationTypes(t *testing.T) {
	types := AllSpecializationTypes()
	if len(types) != 3 {
		t.Errorf("expected 3 specialization types, got %d", len(types))
	}

	// Should not include SpecNone
	for _, st := range types {
		if st == SpecNone {
			t.Error("AllSpecializationTypes should not include SpecNone")
		}
	}
}
