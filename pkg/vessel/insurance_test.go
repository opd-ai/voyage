package vessel

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewInsuranceManager(t *testing.T) {
	im := NewInsuranceManager(engine.GenreFantasy)
	if im == nil {
		t.Fatal("expected non-nil manager")
	}
	if im.Genre() != engine.GenreFantasy {
		t.Errorf("expected fantasy genre, got %s", im.Genre())
	}
}

func TestInsuranceManager_SetGenre(t *testing.T) {
	im := NewInsuranceManager(engine.GenreFantasy)
	im.SetGenre(engine.GenreScifi)
	if im.Genre() != engine.GenreScifi {
		t.Errorf("expected scifi genre, got %s", im.Genre())
	}
}

func TestInsuranceCost(t *testing.T) {
	// Test that costs are positive and increase with tier
	for _, mt := range AllModuleTypes() {
		prevCost := 0.0
		for tier := 1; tier <= 5; tier++ {
			cost := InsuranceCost(mt, tier)
			if cost <= 0 {
				t.Errorf("expected positive cost for %v tier %d", mt, tier)
			}
			if cost <= prevCost && tier > 1 {
				t.Errorf("expected cost to increase with tier for %v", mt)
			}
			prevCost = cost
		}
	}
}

func TestInsuranceManager_PurchaseInsurance(t *testing.T) {
	im := NewInsuranceManager(engine.GenreFantasy)

	// Test successful purchase
	policy, cost, ok := im.PurchaseInsurance(ModuleEngine, 1, 100)
	if !ok {
		t.Error("expected successful purchase")
	}
	if policy == nil {
		t.Fatal("expected non-nil policy")
	}
	if cost <= 0 {
		t.Error("expected positive cost")
	}
	if !policy.Active {
		t.Error("expected active policy")
	}

	// Test duplicate purchase fails
	_, _, ok = im.PurchaseInsurance(ModuleEngine, 1, 100)
	if ok {
		t.Error("expected duplicate purchase to fail")
	}

	// Test insufficient funds
	_, _, ok = im.PurchaseInsurance(ModuleDefense, 1, 1)
	if ok {
		t.Error("expected purchase to fail with insufficient funds")
	}
}

func TestInsuranceManager_HasInsurance(t *testing.T) {
	im := NewInsuranceManager(engine.GenreFantasy)

	if im.HasInsurance(ModuleEngine) {
		t.Error("expected no insurance initially")
	}

	im.PurchaseInsurance(ModuleEngine, 1, 100)
	if !im.HasInsurance(ModuleEngine) {
		t.Error("expected insurance after purchase")
	}
}

func TestInsuranceManager_ClaimInsurance(t *testing.T) {
	im := NewInsuranceManager(engine.GenreFantasy)

	// Cannot claim without policy
	if im.ClaimInsurance(ModuleEngine) {
		t.Error("expected claim to fail without policy")
	}

	im.PurchaseInsurance(ModuleEngine, 1, 100)

	// First claim succeeds
	if !im.ClaimInsurance(ModuleEngine) {
		t.Error("expected first claim to succeed")
	}

	// Second claim fails (single-use)
	if im.ClaimInsurance(ModuleEngine) {
		t.Error("expected second claim to fail")
	}

	// Policy should no longer be active
	if im.HasInsurance(ModuleEngine) {
		t.Error("expected no insurance after claim")
	}
}

func TestInsuranceManager_CancelInsurance(t *testing.T) {
	im := NewInsuranceManager(engine.GenreFantasy)

	// Cannot cancel without policy
	_, ok := im.CancelInsurance(ModuleEngine)
	if ok {
		t.Error("expected cancel to fail without policy")
	}

	im.PurchaseInsurance(ModuleEngine, 1, 100)

	// Cancel succeeds and returns refund
	refund, ok := im.CancelInsurance(ModuleEngine)
	if !ok {
		t.Error("expected cancel to succeed")
	}
	if refund <= 0 {
		t.Error("expected positive refund")
	}

	// Cannot cancel again
	_, ok = im.CancelInsurance(ModuleEngine)
	if ok {
		t.Error("expected second cancel to fail")
	}
}

func TestInsuranceManager_GetAllPolicies(t *testing.T) {
	im := NewInsuranceManager(engine.GenreFantasy)

	policies := im.GetAllPolicies()
	if len(policies) != 0 {
		t.Error("expected no policies initially")
	}

	im.PurchaseInsurance(ModuleEngine, 1, 100)
	im.PurchaseInsurance(ModuleDefense, 1, 100)

	policies = im.GetAllPolicies()
	if len(policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(policies))
	}
}

func TestInsuranceManager_GetAvailableInsurance(t *testing.T) {
	im := NewInsuranceManager(engine.GenreFantasy)
	ms := NewModuleSystem(engine.GenreFantasy)

	options := im.GetAvailableInsurance(ms)
	if len(options) != 5 {
		t.Errorf("expected 5 options, got %d", len(options))
	}

	// Purchase one
	im.PurchaseInsurance(ModuleEngine, 1, 100)

	options = im.GetAvailableInsurance(ms)
	if len(options) != 4 {
		t.Errorf("expected 4 options after purchase, got %d", len(options))
	}
}

func TestInsuranceManager_CheckProtection(t *testing.T) {
	im := NewInsuranceManager(engine.GenreFantasy)

	// Non-catastrophic damage is not protected
	result := im.CheckProtection(ModuleEngine, 10, false)
	if result.WasProtected {
		t.Error("expected non-catastrophic damage to not be protected")
	}

	// Catastrophic without insurance
	result = im.CheckProtection(ModuleEngine, 100, true)
	if result.WasProtected {
		t.Error("expected no protection without insurance")
	}

	// Catastrophic with insurance
	im.PurchaseInsurance(ModuleEngine, 1, 100)
	result = im.CheckProtection(ModuleEngine, 100, true)
	if !result.WasProtected {
		t.Error("expected protection with insurance")
	}
	if result.PreventedDamage != 100 {
		t.Errorf("expected prevented damage of 100, got %f", result.PreventedDamage)
	}

	// Insurance is used up
	result = im.CheckProtection(ModuleEngine, 100, true)
	if result.WasProtected {
		t.Error("expected no protection after claim")
	}
}

func TestInsuranceTypeName_AllGenres(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		name := InsuranceTypeName(genre)
		if name == "" {
			t.Errorf("empty insurance type name for genre %s", genre)
		}
	}
}

func TestInsuranceDescription_AllGenres(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		desc := InsuranceDescription(genre)
		if desc == "" {
			t.Errorf("empty insurance description for genre %s", genre)
		}
	}
}

func TestInsurancePolicy_PolicyName(t *testing.T) {
	im := NewInsuranceManager(engine.GenreScifi)
	policy, _, _ := im.PurchaseInsurance(ModuleEngine, 1, 100)

	name := policy.PolicyName()
	if name == "" {
		t.Error("expected non-empty policy name")
	}
}

func TestInsurancePolicy_Status(t *testing.T) {
	im := NewInsuranceManager(engine.GenreFantasy)
	policy, _, _ := im.PurchaseInsurance(ModuleEngine, 1, 100)

	if policy.Status() != "Active" {
		t.Errorf("expected Active status, got %s", policy.Status())
	}

	im.ClaimInsurance(ModuleEngine)
	if policy.Status() != "Claimed" {
		t.Errorf("expected Claimed status, got %s", policy.Status())
	}
}

func TestInsuranceManager_TotalPolicyCost(t *testing.T) {
	im := NewInsuranceManager(engine.GenreFantasy)

	if im.TotalPolicyCost() != 0 {
		t.Error("expected zero initial cost")
	}

	im.PurchaseInsurance(ModuleEngine, 1, 100)
	im.PurchaseInsurance(ModuleDefense, 1, 100)

	total := im.TotalPolicyCost()
	if total <= 0 {
		t.Error("expected positive total cost")
	}
}

func TestInsuranceManager_InsuredModuleCount(t *testing.T) {
	im := NewInsuranceManager(engine.GenreFantasy)

	if im.InsuredModuleCount() != 0 {
		t.Error("expected zero initial count")
	}

	im.PurchaseInsurance(ModuleEngine, 1, 100)
	im.PurchaseInsurance(ModuleDefense, 1, 100)

	if im.InsuredModuleCount() != 2 {
		t.Errorf("expected 2 insured modules, got %d", im.InsuredModuleCount())
	}

	im.ClaimInsurance(ModuleEngine)
	if im.InsuredModuleCount() != 1 {
		t.Errorf("expected 1 insured module after claim, got %d", im.InsuredModuleCount())
	}
}
