package vessel

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// InsurancePolicy represents protection for a vessel module against breakdown.
type InsurancePolicy struct {
	ModuleType    ModuleType
	PurchasePrice float64
	Genre         engine.GenreID
	Active        bool
	UsedOnce      bool // Insurance is single-use per breakdown
}

// InsuranceManager handles purchasing and claiming insurance.
type InsuranceManager struct {
	genre    engine.GenreID
	policies map[ModuleType]*InsurancePolicy
}

// NewInsuranceManager creates a new insurance manager.
func NewInsuranceManager(genre engine.GenreID) *InsuranceManager {
	return &InsuranceManager{
		genre:    genre,
		policies: make(map[ModuleType]*InsurancePolicy),
	}
}

// SetGenre changes the insurance vocabulary.
func (im *InsuranceManager) SetGenre(genre engine.GenreID) {
	im.genre = genre
	// Update genre on all existing policies
	for _, p := range im.policies {
		p.Genre = genre
	}
}

// Genre returns the current genre.
func (im *InsuranceManager) Genre() engine.GenreID {
	return im.genre
}

// InsuranceCost calculates the cost to insure a module.
// Higher tier modules cost more to insure.
func InsuranceCost(moduleType ModuleType, tier int) float64 {
	baseCosts := map[ModuleType]float64{
		ModuleEngine:     35,
		ModuleCargoHold:  25,
		ModuleMedicalBay: 30,
		ModuleNavigation: 25,
		ModuleDefense:    40,
	}
	base := baseCosts[moduleType]
	// Cost increases by 25% per tier above 1
	return base * (1.0 + float64(tier-1)*0.25)
}

// PurchaseInsurance buys insurance for a module.
// Returns true if the purchase was successful.
func (im *InsuranceManager) PurchaseInsurance(moduleType ModuleType, tier int, currency float64) (*InsurancePolicy, float64, bool) {
	cost := InsuranceCost(moduleType, tier)
	if currency < cost {
		return nil, 0, false
	}

	// Check if already insured
	if existing, ok := im.policies[moduleType]; ok && existing.Active {
		return nil, 0, false
	}

	policy := &InsurancePolicy{
		ModuleType:    moduleType,
		PurchasePrice: cost,
		Genre:         im.genre,
		Active:        true,
		UsedOnce:      false,
	}
	im.policies[moduleType] = policy

	return policy, cost, true
}

// HasInsurance checks if a module is currently insured.
func (im *InsuranceManager) HasInsurance(moduleType ModuleType) bool {
	if p, ok := im.policies[moduleType]; ok {
		return p.Active && !p.UsedOnce
	}
	return false
}

// GetPolicy returns the insurance policy for a module, if any.
func (im *InsuranceManager) GetPolicy(moduleType ModuleType) *InsurancePolicy {
	return im.policies[moduleType]
}

// ClaimInsurance attempts to claim insurance for a catastrophic breakdown.
// Returns true if the claim was successful and the module was protected.
func (im *InsuranceManager) ClaimInsurance(moduleType ModuleType) bool {
	p, ok := im.policies[moduleType]
	if !ok || !p.Active || p.UsedOnce {
		return false
	}

	// Mark the policy as used
	p.UsedOnce = true
	p.Active = false

	return true
}

// CancelInsurance cancels an active policy.
// Returns partial refund (50% of purchase price).
func (im *InsuranceManager) CancelInsurance(moduleType ModuleType) (float64, bool) {
	p, ok := im.policies[moduleType]
	if !ok || !p.Active || p.UsedOnce {
		return 0, false
	}

	refund := p.PurchasePrice * 0.5
	p.Active = false
	delete(im.policies, moduleType)

	return refund, true
}

// GetAllPolicies returns all active insurance policies.
func (im *InsuranceManager) GetAllPolicies() []*InsurancePolicy {
	var policies []*InsurancePolicy
	for _, p := range im.policies {
		if p.Active && !p.UsedOnce {
			policies = append(policies, p)
		}
	}
	return policies
}

// GetAvailableInsurance returns insurance options for uninsured modules.
func (im *InsuranceManager) GetAvailableInsurance(ms *ModuleSystem) []InsuranceOption {
	var options []InsuranceOption
	for _, mt := range AllModuleTypes() {
		if im.HasInsurance(mt) {
			continue
		}
		m := ms.GetModule(mt)
		if m == nil {
			continue
		}
		options = append(options, InsuranceOption{
			ModuleType: mt,
			Tier:       m.Tier(),
			Cost:       InsuranceCost(mt, m.Tier()),
			ModuleName: ModuleTypeName(mt, im.genre),
		})
	}
	return options
}

// InsuranceOption represents an available insurance purchase.
type InsuranceOption struct {
	ModuleType ModuleType
	Tier       int
	Cost       float64
	ModuleName string
}

// InsuranceTypeName returns the genre-appropriate name for insurance.
func InsuranceTypeName(genre engine.GenreID) string {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Protection Ward",
		engine.GenreScifi:     "Warranty Plan",
		engine.GenreHorror:    "Salvage Insurance",
		engine.GenreCyberpunk: "Repair Contract",
		engine.GenrePostapoc:  "Barter Bond",
	}
	if name, ok := names[genre]; ok {
		return name
	}
	return "Insurance"
}

// InsuranceDescription returns a description of what insurance does.
func InsuranceDescription(genre engine.GenreID) string {
	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Protects one module from a catastrophic failure. The ward absorbs the damage, leaving the module intact.",
		engine.GenreScifi:     "Covers emergency repairs for one catastrophic malfunction. Automated nanites restore module integrity.",
		engine.GenreHorror:    "Guarantees priority salvage rights. If a module is destroyed, spare parts will be found to rebuild it.",
		engine.GenreCyberpunk: "A fixer's guarantee. One free emergency repair when a module goes critical.",
		engine.GenrePostapoc:  "A promise backed by the tribe. If your gear breaks down, the community will help you rebuild.",
	}
	if desc, ok := descriptions[genre]; ok {
		return desc
	}
	return "Protects one module from a catastrophic breakdown."
}

// PolicyName returns the genre-appropriate name for a policy.
func (p *InsurancePolicy) PolicyName() string {
	moduleName := ModuleTypeName(p.ModuleType, p.Genre)
	typeName := InsuranceTypeName(p.Genre)
	return moduleName + " " + typeName
}

// Status returns the current status of the policy.
func (p *InsurancePolicy) Status() string {
	if p.UsedOnce {
		return "Claimed"
	}
	if p.Active {
		return "Active"
	}
	return "Expired"
}

// ProtectionResult represents the outcome of checking module protection.
type ProtectionResult struct {
	WasProtected    bool
	ModuleType      ModuleType
	PreventedDamage float64
	Message         string
}

// CheckProtection checks if a module has insurance and claims it if needed.
// Returns whether the module was protected from a catastrophic failure.
func (im *InsuranceManager) CheckProtection(moduleType ModuleType, damage float64, wouldDestroy bool) ProtectionResult {
	if !wouldDestroy {
		return ProtectionResult{
			WasProtected: false,
			ModuleType:   moduleType,
			Message:      "Damage was not catastrophic",
		}
	}

	if im.ClaimInsurance(moduleType) {
		moduleName := ModuleTypeName(moduleType, im.genre)
		return ProtectionResult{
			WasProtected:    true,
			ModuleType:      moduleType,
			PreventedDamage: damage,
			Message:         moduleName + " was protected by " + InsuranceTypeName(im.genre),
		}
	}

	return ProtectionResult{
		WasProtected: false,
		ModuleType:   moduleType,
		Message:      "No active insurance for this module",
	}
}

// TotalPolicyCost returns the total currency spent on active policies.
func (im *InsuranceManager) TotalPolicyCost() float64 {
	total := 0.0
	for _, p := range im.policies {
		if p.Active && !p.UsedOnce {
			total += p.PurchasePrice
		}
	}
	return total
}

// InsuredModuleCount returns the number of currently insured modules.
func (im *InsuranceManager) InsuredModuleCount() int {
	count := 0
	for _, p := range im.policies {
		if p.Active && !p.UsedOnce {
			count++
		}
	}
	return count
}
