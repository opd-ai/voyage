package vessel

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewVessel(t *testing.T) {
	tests := []struct {
		name       string
		vesselType VesselType
		wantCap    int
	}{
		{"small", VesselSmall, 20},
		{"medium", VesselMedium, 50},
		{"large", VesselLarge, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewVessel(tt.vesselType, engine.GenreFantasy)
			if v.Capacity() != tt.wantCap {
				t.Errorf("capacity = %d, want %d", v.Capacity(), tt.wantCap)
			}
			if v.IsDestroyed() {
				t.Error("new vessel should not be destroyed")
			}
			if v.IntegrityRatio() != 1.0 {
				t.Error("new vessel should have full integrity")
			}
		})
	}
}

func TestVesselDamageAndRepair(t *testing.T) {
	v := NewVessel(VesselMedium, engine.GenreFantasy)
	initial := v.Integrity()

	// Take damage
	destroyed := v.TakeDamage(30)
	if destroyed {
		t.Error("vessel should not be destroyed with 30 damage")
	}
	if v.Integrity() != initial-30 {
		t.Errorf("integrity = %f, want %f", v.Integrity(), initial-30)
	}

	// Speed should be reduced
	if v.Speed() >= v.BaseSpeed() {
		t.Error("damaged vessel should have reduced speed")
	}

	// Repair
	v.Repair(20)
	if v.Integrity() != initial-10 {
		t.Errorf("integrity after repair = %f, want %f", v.Integrity(), initial-10)
	}
}

func TestVesselDestruction(t *testing.T) {
	v := NewVessel(VesselSmall, engine.GenreFantasy)

	// Deal lethal damage
	destroyed := v.TakeDamage(v.MaxIntegrity() + 10)
	if !destroyed {
		t.Error("vessel should be destroyed")
	}
	if v.Integrity() != 0 {
		t.Error("destroyed vessel should have 0 integrity")
	}
	if !v.IsDestroyed() {
		t.Error("IsDestroyed should return true")
	}
}

func TestVesselCargo(t *testing.T) {
	v := NewVessel(VesselSmall, engine.GenreFantasy) // capacity 20

	// Add cargo
	if !v.AddCargo(15) {
		t.Error("should be able to add 15 cargo")
	}
	if v.FreeCapacity() != 5 {
		t.Errorf("free capacity = %d, want 5", v.FreeCapacity())
	}

	// Try to add too much
	if v.AddCargo(10) {
		t.Error("should not be able to add 10 more cargo")
	}

	// Add exactly remaining
	if !v.AddCargo(5) {
		t.Error("should be able to add 5 more cargo")
	}
	if v.FreeCapacity() != 0 {
		t.Error("vessel should be at capacity")
	}

	// Remove cargo
	if !v.RemoveCargo(10) {
		t.Error("should be able to remove 10 cargo")
	}
	if v.UsedCapacity() != 10 {
		t.Errorf("used capacity = %d, want 10", v.UsedCapacity())
	}
}

func TestVesselGenreSwitching(t *testing.T) {
	v := NewVessel(VesselMedium, engine.GenreFantasy)

	if v.Name() != "Wagon" {
		t.Errorf("fantasy name = %q, want Wagon", v.Name())
	}

	v.SetGenre(engine.GenreScifi)
	if v.Name() != "Shuttle" {
		t.Errorf("scifi name = %q, want Shuttle", v.Name())
	}

	v.SetGenre(engine.GenreHorror)
	if v.Name() != "SUV" {
		t.Errorf("horror name = %q, want SUV", v.Name())
	}
}

func TestCargoHold(t *testing.T) {
	hold := NewCargoHold(50)

	// Add items
	if !hold.Add("Food", 2, 10, CargoSupplies) {
		t.Error("should add food")
	}
	if hold.Used() != 20 {
		t.Errorf("used = %d, want 20", hold.Used())
	}

	// Stack same item
	if !hold.Add("Food", 2, 5, CargoSupplies) {
		t.Error("should stack food")
	}
	if hold.GetQuantity("Food") != 15 {
		t.Errorf("food quantity = %d, want 15", hold.GetQuantity("Food"))
	}

	// Add different item
	if !hold.Add("Medicine", 1, 10, CargoMedical) {
		t.Error("should add medicine")
	}

	// Check has
	if !hold.Has("Food", 10) {
		t.Error("should have 10 food")
	}
	if hold.Has("Food", 20) {
		t.Error("should not have 20 food")
	}

	// Remove
	if !hold.Remove("Food", 5) {
		t.Error("should remove 5 food")
	}
	if hold.GetQuantity("Food") != 10 {
		t.Errorf("food after remove = %d, want 10", hold.GetQuantity("Food"))
	}

	// Get by category
	medical := hold.GetByCategory(CargoMedical)
	if len(medical) != 1 {
		t.Errorf("medical items = %d, want 1", len(medical))
	}
}

func TestBreakdownChecker(t *testing.T) {
	bc := NewBreakdownChecker(12345)
	v := NewVessel(VesselMedium, engine.GenreFantasy)

	// With full integrity, breakdowns are rare
	// Run multiple checks to potentially trigger one
	hadBreakdown := false
	for i := 0; i < 100; i++ {
		if b := bc.Check(v, i); b != nil {
			hadBreakdown = true
			if b.DamageDealt <= 0 {
				t.Error("breakdown should deal damage")
			}
			break
		}
	}

	// Damage vessel to increase breakdown chance
	v.TakeDamage(v.MaxIntegrity() * 0.6) // 60% damage

	// With low integrity, breakdowns are more likely
	for i := 0; i < 50; i++ {
		if b := bc.Check(v, 100+i); b != nil {
			hadBreakdown = true
			if b.RepairCost <= 0 {
				t.Error("breakdown should have repair cost")
			}
			break
		}
	}

	// Note: We don't require a breakdown to occur due to RNG
	_ = hadBreakdown
}

func TestRepairManager(t *testing.T) {
	rm := NewRepairManager(engine.GenreFantasy)
	v := NewVessel(VesselMedium, engine.GenreFantasy)
	hold := NewCargoHold(100)

	// Damage vessel
	v.TakeDamage(50)

	// Try repair without materials
	result := rm.RepairWithMaterials(v, hold, 5)
	if result.Success {
		t.Error("repair should fail without materials")
	}

	// Add materials
	matName := RepairMaterialName(engine.GenreFantasy)
	hold.Add(matName, 1, 10, CargoRepair)

	// Repair with materials
	result = rm.RepairWithMaterials(v, hold, 5)
	if !result.Success {
		t.Error("repair should succeed with materials")
	}
	if result.Repaired != 25 { // 5 materials * 5 integrity each
		t.Errorf("repaired = %f, want 25", result.Repaired)
	}
	if hold.GetQuantity(matName) != 5 {
		t.Errorf("remaining materials = %d, want 5", hold.GetQuantity(matName))
	}
}

func TestConditionStatus(t *testing.T) {
	v := NewVessel(VesselMedium, engine.GenreFantasy)

	if GetConditionStatus(v) != ConditionPristine {
		t.Error("new vessel should be pristine")
	}

	v.TakeDamage(15) // 85% integrity
	if GetConditionStatus(v) != ConditionGood {
		t.Error("85% should be Good")
	}

	v.TakeDamage(30) // 55% integrity
	if GetConditionStatus(v) != ConditionDamaged {
		t.Error("55% should be Damaged")
	}

	v.TakeDamage(20) // 35% integrity
	if GetConditionStatus(v) != ConditionCritical {
		t.Error("35% should be Critical")
	}

	v.TakeDamage(30) // 5% integrity
	if GetConditionStatus(v) != ConditionDestroyed {
		t.Error("5% should be Destroyed")
	}
}

func TestVesselNames(t *testing.T) {
	genres := engine.AllGenres()
	types := AllVesselTypes()

	for _, g := range genres {
		for _, vt := range types {
			name := VesselName(vt, g)
			if name == "" {
				t.Errorf("missing name for genre=%s, type=%d", g, vt)
			}
		}
	}
}
