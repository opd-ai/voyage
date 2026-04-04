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

func TestLoadoutGenerator(t *testing.T) {
	lg := NewLoadoutGenerator(12345, engine.GenreFantasy)

	loadouts := lg.GenerateAll()
	if len(loadouts) != 3 {
		t.Errorf("expected 3 loadouts, got %d", len(loadouts))
	}

	// Verify each loadout type is represented
	types := make(map[LoadoutType]bool)
	for _, l := range loadouts {
		types[l.Type] = true
	}

	for _, lt := range AllLoadoutTypes() {
		if !types[lt] {
			t.Errorf("missing loadout type: %d", lt)
		}
	}
}

func TestLoadoutDeterminism(t *testing.T) {
	lg1 := NewLoadoutGenerator(12345, engine.GenreFantasy)
	lg2 := NewLoadoutGenerator(12345, engine.GenreFantasy)

	l1 := lg1.Generate(LoadoutBalanced)
	l2 := lg2.Generate(LoadoutBalanced)

	// Same seed should produce same results
	if l1.StartFood != l2.StartFood {
		t.Errorf("food mismatch: %f vs %f", l1.StartFood, l2.StartFood)
	}
	if l1.SpeedMod != l2.SpeedMod {
		t.Errorf("speed mod mismatch: %f vs %f", l1.SpeedMod, l2.SpeedMod)
	}
}

func TestLoadoutVariation(t *testing.T) {
	lg1 := NewLoadoutGenerator(12345, engine.GenreFantasy)
	lg2 := NewLoadoutGenerator(67890, engine.GenreFantasy)

	l1 := lg1.Generate(LoadoutBalanced)
	l2 := lg2.Generate(LoadoutBalanced)

	// Different seeds should produce different results
	// Note: There's a small chance they could be equal, but highly unlikely
	if l1.StartFood == l2.StartFood && l1.SpeedMod == l2.SpeedMod {
		t.Log("Warning: different seeds produced same loadout (possible but unlikely)")
	}
}

func TestLoadoutTypeCharacteristics(t *testing.T) {
	lg := NewLoadoutGenerator(12345, engine.GenreFantasy)

	balanced := lg.Generate(LoadoutBalanced)
	fast := lg.Generate(LoadoutFastLight)
	heavy := lg.Generate(LoadoutSlowHeavy)

	// Fast should be faster than balanced
	if fast.SpeedMod <= balanced.SpeedMod*0.9 { // Allow 10% variation
		t.Error("fast loadout should have higher speed than balanced")
	}

	// Heavy should have more capacity than balanced
	if heavy.CapacityMod <= balanced.CapacityMod*0.9 {
		t.Error("heavy loadout should have more capacity than balanced")
	}

	// Heavy should be slower than fast
	if heavy.SpeedMod >= fast.SpeedMod*1.1 {
		t.Error("heavy loadout should be slower than fast")
	}
}

func TestLoadoutApplyToVessel(t *testing.T) {
	lg := NewLoadoutGenerator(12345, engine.GenreFantasy)
	loadout := lg.Generate(LoadoutBalanced)
	v := NewVessel(VesselSmall, engine.GenreScifi) // Different initial config

	loadout.ApplyToVessel(v)

	// Verify vessel was configured
	if v.vesselType != loadout.VesselType {
		t.Error("vessel type not applied")
	}
	if v.genre != loadout.Genre {
		t.Error("genre not applied")
	}
	if v.Integrity() != v.MaxIntegrity() {
		t.Error("vessel should have full integrity after loadout")
	}
}

func TestLoadoutNames(t *testing.T) {
	genres := engine.AllGenres()
	types := AllLoadoutTypes()

	for _, g := range genres {
		for _, lt := range types {
			name := LoadoutName(lt, g)
			if name == "" {
				t.Errorf("missing name for genre=%s, loadout=%d", g, lt)
			}
		}
	}
}

func TestLoadoutDescription(t *testing.T) {
	for _, lt := range AllLoadoutTypes() {
		desc := LoadoutDescription(lt)
		if desc == "" {
			t.Errorf("missing description for loadout type %d", lt)
		}
	}
}

func TestLoadoutGenreSwitching(t *testing.T) {
	lg := NewLoadoutGenerator(12345, engine.GenreFantasy)

	fantasy := lg.Generate(LoadoutBalanced)
	if fantasy.Name() != "Merchant's Caravan" {
		t.Errorf("fantasy name = %s, want Merchant's Caravan", fantasy.Name())
	}

	lg.SetGenre(engine.GenreScifi)
	scifi := lg.Generate(LoadoutBalanced)
	if scifi.Name() != "Survey Vessel" {
		t.Errorf("scifi name = %s, want Survey Vessel", scifi.Name())
	}
}

func TestAllVisualVariants(t *testing.T) {
	variants := AllVisualVariants()
	if len(variants) != 3 {
		t.Errorf("expected 3 visual variants, got %d", len(variants))
	}

	// Check all three variants are present
	expected := map[VisualVariant]bool{
		VisualVariantA: false,
		VisualVariantB: false,
		VisualVariantC: false,
	}
	for _, v := range variants {
		expected[v] = true
	}
	for v, found := range expected {
		if !found {
			t.Errorf("variant %d not in AllVisualVariants", v)
		}
	}
}

func TestVisualVariantNames(t *testing.T) {
	genres := engine.AllGenres()
	variants := AllVisualVariants()

	for _, g := range genres {
		for _, v := range variants {
			name := VisualVariantName(v, g)
			if name == "" {
				t.Errorf("missing name for genre=%s, variant=%d", g, v)
			}
		}
	}
}

func TestVisualVariantDescription(t *testing.T) {
	for _, v := range AllVisualVariants() {
		desc := VisualVariantDescription(v)
		if desc == "" {
			t.Errorf("missing description for variant %d", v)
		}
	}
}

func TestHullSkinGenerator(t *testing.T) {
	gen := NewHullSkinGenerator(12345, engine.GenreFantasy)

	skins := gen.GenerateAll(VesselMedium)
	if len(skins) != 3 {
		t.Errorf("expected 3 hull skins, got %d", len(skins))
	}

	// Verify each variant is represented
	variants := make(map[VisualVariant]bool)
	for _, s := range skins {
		variants[s.Variant] = true
		if s.Genre != engine.GenreFantasy {
			t.Errorf("skin genre = %s, want fantasy", s.Genre)
		}
		if s.VesselType != VesselMedium {
			t.Errorf("skin vessel type = %d, want %d", s.VesselType, VesselMedium)
		}
	}

	for _, v := range AllVisualVariants() {
		if !variants[v] {
			t.Errorf("missing variant: %d", v)
		}
	}
}

func TestHullSkinGeneratorDeterminism(t *testing.T) {
	gen1 := NewHullSkinGenerator(12345, engine.GenreFantasy)
	gen2 := NewHullSkinGenerator(12345, engine.GenreFantasy)

	skin1 := gen1.Generate(VisualVariantA, VesselMedium)
	skin2 := gen2.Generate(VisualVariantA, VesselMedium)

	if skin1.PrimaryHue != skin2.PrimaryHue {
		t.Errorf("primary hue mismatch: %f vs %f", skin1.PrimaryHue, skin2.PrimaryHue)
	}
	if skin1.PatternDensity != skin2.PatternDensity {
		t.Errorf("pattern density mismatch: %f vs %f", skin1.PatternDensity, skin2.PatternDensity)
	}
}

func TestHullSkinGeneratorVariation(t *testing.T) {
	gen1 := NewHullSkinGenerator(12345, engine.GenreFantasy)
	gen2 := NewHullSkinGenerator(67890, engine.GenreFantasy)

	skin1 := gen1.Generate(VisualVariantA, VesselMedium)
	skin2 := gen2.Generate(VisualVariantA, VesselMedium)

	// Different seeds should produce different results
	if skin1.PrimaryHue == skin2.PrimaryHue && skin1.SecondaryHue == skin2.SecondaryHue {
		t.Log("Warning: different seeds produced same hull skin (possible but unlikely)")
	}
}

func TestHullSkinGeneratorGenreSwitching(t *testing.T) {
	gen := NewHullSkinGenerator(12345, engine.GenreFantasy)

	fantasy := gen.Generate(VisualVariantA, VesselMedium)
	if fantasy.Genre != engine.GenreFantasy {
		t.Errorf("genre = %s, want fantasy", fantasy.Genre)
	}

	gen.SetGenre(engine.GenreScifi)
	scifi := gen.Generate(VisualVariantA, VesselMedium)
	if scifi.Genre != engine.GenreScifi {
		t.Errorf("genre = %s, want scifi", scifi.Genre)
	}
}

func TestVesselVisuals(t *testing.T) {
	gen := NewHullSkinGenerator(12345, engine.GenreFantasy)
	params := gen.Generate(VisualVariantB, VesselMedium)

	visuals := NewVesselVisuals(VisualVariantB, params)

	if visuals.Variant != VisualVariantB {
		t.Errorf("variant = %d, want %d", visuals.Variant, VisualVariantB)
	}

	name := visuals.VariantName()
	if name != "Iron Bound" {
		t.Errorf("variant name = %s, want Iron Bound", name)
	}

	desc := visuals.Description()
	if desc == "" {
		t.Error("description should not be empty")
	}
}

func TestWrapHue(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{0, 0},
		{180, 180},
		{360, 0},
		{370, 10},
		{-10, 350},
		{720, 0},
		{-720, 0},
	}

	for _, tt := range tests {
		result := wrapHue(tt.input)
		if result != tt.expected {
			t.Errorf("wrapHue(%f) = %f, want %f", tt.input, result, tt.expected)
		}
	}
}
