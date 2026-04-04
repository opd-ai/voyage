package vessel

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// CustomizationVocab holds genre-specific vocabulary for the customization screen.
type CustomizationVocab struct {
	// Screen titles and labels
	ScreenTitle       string
	SubtitleText      string
	UpgradePointsText string
	ConfirmText       string
	CancelText        string

	// Module-related
	ModuleSectionTitle    string
	ModuleTierLabel       string
	ModuleEffectLabel     string
	ModuleUpgradeAction   string
	ModuleDowngradeAction string

	// Loadout-related
	LoadoutSectionTitle string
	LoadoutSelectText   string
	LoadoutPreviewText  string

	// Visual variant-related
	VariantSectionTitle string
	VariantSelectText   string

	// Insignia-related
	InsigniaSectionTitle string
	InsigniaSelectText   string

	// Insurance-related
	InsuranceSectionTitle string
	InsurancePurchaseText string

	// Summary
	SummaryTitle       string
	TotalCostLabel     string
	TotalCapacityLabel string
	TotalSpeedLabel    string
	TotalDefenseLabel  string
}

// GetCustomizationVocab returns genre-specific customization vocabulary.
func GetCustomizationVocab(genre engine.GenreID) *CustomizationVocab {
	vocab, ok := customizationVocabs[genre]
	if !ok {
		return customizationVocabs[engine.GenreFantasy]
	}
	return vocab
}

var customizationVocabs = map[engine.GenreID]*CustomizationVocab{
	engine.GenreFantasy: {
		ScreenTitle:           "Prepare Your Caravan",
		SubtitleText:          "Outfit your wagon for the long road ahead",
		UpgradePointsText:     "Crafting Points",
		ConfirmText:           "Begin Journey",
		CancelText:            "Return to Camp",
		ModuleSectionTitle:    "Caravan Components",
		ModuleTierLabel:       "Quality",
		ModuleEffectLabel:     "Benefit",
		ModuleUpgradeAction:   "Improve",
		ModuleDowngradeAction: "Reduce",
		LoadoutSectionTitle:   "Caravan Type",
		LoadoutSelectText:     "Choose your wagon style",
		LoadoutPreviewText:    "Preview supplies",
		VariantSectionTitle:   "Wagon Appearance",
		VariantSelectText:     "Choose wood and trim",
		InsigniaSectionTitle:  "Caravan Banner",
		InsigniaSelectText:    "Design your heraldry",
		InsuranceSectionTitle: "Protection Wards",
		InsurancePurchaseText: "Purchase magical protection",
		SummaryTitle:          "Journey Preparation",
		TotalCostLabel:        "Gold Spent",
		TotalCapacityLabel:    "Carrying Capacity",
		TotalSpeedLabel:       "Travel Speed",
		TotalDefenseLabel:     "Protection Rating",
	},
	engine.GenreScifi: {
		ScreenTitle:           "Configure Ship Systems",
		SubtitleText:          "Initialize subsystems before departure",
		UpgradePointsText:     "Tech Credits",
		ConfirmText:           "Launch Sequence",
		CancelText:            "Return to Hangar",
		ModuleSectionTitle:    "Ship Modules",
		ModuleTierLabel:       "Rating",
		ModuleEffectLabel:     "Function",
		ModuleUpgradeAction:   "Upgrade",
		ModuleDowngradeAction: "Downgrade",
		LoadoutSectionTitle:   "Ship Configuration",
		LoadoutSelectText:     "Select ship class",
		LoadoutPreviewText:    "Review manifest",
		VariantSectionTitle:   "Hull Design",
		VariantSelectText:     "Choose hull plating",
		InsigniaSectionTitle:  "Ship Insignia",
		InsigniaSelectText:    "Design corps emblem",
		InsuranceSectionTitle: "Warranty Plans",
		InsurancePurchaseText: "Purchase repair warranty",
		SummaryTitle:          "Launch Readiness",
		TotalCostLabel:        "Credits Used",
		TotalCapacityLabel:    "Cargo Tonnage",
		TotalSpeedLabel:       "Warp Factor",
		TotalDefenseLabel:     "Shield Rating",
	},
	engine.GenreHorror: {
		ScreenTitle:           "Outfit Your Vehicle",
		SubtitleText:          "Prepare to survive the wasteland",
		UpgradePointsText:     "Salvage Parts",
		ConfirmText:           "Hit the Road",
		CancelText:            "Stay Here",
		ModuleSectionTitle:    "Vehicle Parts",
		ModuleTierLabel:       "Condition",
		ModuleEffectLabel:     "Effect",
		ModuleUpgradeAction:   "Reinforce",
		ModuleDowngradeAction: "Strip",
		LoadoutSectionTitle:   "Vehicle Type",
		LoadoutSelectText:     "Choose your ride",
		LoadoutPreviewText:    "Check supplies",
		VariantSectionTitle:   "Vehicle Look",
		VariantSelectText:     "Choose body style",
		InsigniaSectionTitle:  "Group Mark",
		InsigniaSelectText:    "Paint your symbol",
		InsuranceSectionTitle: "Salvage Insurance",
		InsurancePurchaseText: "Pre-arrange spare parts",
		SummaryTitle:          "Survival Check",
		TotalCostLabel:        "Parts Used",
		TotalCapacityLabel:    "Storage Space",
		TotalSpeedLabel:       "Max Speed",
		TotalDefenseLabel:     "Armor Rating",
	},
	engine.GenreCyberpunk: {
		ScreenTitle:           "Customize Your Rig",
		SubtitleText:          "Mod your ride before the run",
		UpgradePointsText:     "Mod Points",
		ConfirmText:           "Jack In",
		CancelText:            "Back Out",
		ModuleSectionTitle:    "Rig Systems",
		ModuleTierLabel:       "Spec",
		ModuleEffectLabel:     "Output",
		ModuleUpgradeAction:   "Boost",
		ModuleDowngradeAction: "Nerf",
		LoadoutSectionTitle:   "Rig Profile",
		LoadoutSelectText:     "Select build type",
		LoadoutPreviewText:    "Check loadout",
		VariantSectionTitle:   "Exterior Mods",
		VariantSelectText:     "Choose body kit",
		InsigniaSectionTitle:  "Crew Tag",
		InsigniaSelectText:    "Design gang colors",
		InsuranceSectionTitle: "Fixer Contracts",
		InsurancePurchaseText: "Arrange emergency repairs",
		SummaryTitle:          "Build Summary",
		TotalCostLabel:        "Creds Burned",
		TotalCapacityLabel:    "Data Capacity",
		TotalSpeedLabel:       "Processing Speed",
		TotalDefenseLabel:     "ICE Rating",
	},
	engine.GenrePostapoc: {
		ScreenTitle:           "Set Up Your Ride",
		SubtitleText:          "Scrap together what you can",
		UpgradePointsText:     "Scrap Budget",
		ConfirmText:           "Roll Out",
		CancelText:            "Back to Camp",
		ModuleSectionTitle:    "Rig Components",
		ModuleTierLabel:       "Grade",
		ModuleEffectLabel:     "Function",
		ModuleUpgradeAction:   "Weld On",
		ModuleDowngradeAction: "Strip Off",
		LoadoutSectionTitle:   "Rig Type",
		LoadoutSelectText:     "Pick your hauler",
		LoadoutPreviewText:    "Inventory check",
		VariantSectionTitle:   "War Paint",
		VariantSelectText:     "Choose your look",
		InsigniaSectionTitle:  "Tribe Mark",
		InsigniaSelectText:    "Brand your symbol",
		InsuranceSectionTitle: "Barter Bonds",
		InsurancePurchaseText: "Trade for spare parts guarantee",
		SummaryTitle:          "Ready Check",
		TotalCostLabel:        "Scrap Spent",
		TotalCapacityLabel:    "Hauling Weight",
		TotalSpeedLabel:       "Road Speed",
		TotalDefenseLabel:     "Plating Rating",
	},
}

// ModuleUpgradeVocab holds genre-specific vocabulary for module upgrades.
type ModuleUpgradeVocab struct {
	// Actions
	UpgradeText    string
	DowngradeText  string
	RepairText     string
	SpecializeText string

	// Status
	MaxedOutText     string
	InsufficientText string
	SuccessText      string
	FailureText      string

	// Costs
	CurrencyName  string
	MaterialsName string

	// Module effects
	SpeedBoostText      string
	CargoBoostText      string
	DefenseBoostText    string
	HealingBoostText    string
	NavigationBoostText string
}

// GetModuleUpgradeVocab returns genre-specific module upgrade vocabulary.
func GetModuleUpgradeVocab(genre engine.GenreID) *ModuleUpgradeVocab {
	vocab, ok := moduleUpgradeVocabs[genre]
	if !ok {
		return moduleUpgradeVocabs[engine.GenreFantasy]
	}
	return vocab
}

var moduleUpgradeVocabs = map[engine.GenreID]*ModuleUpgradeVocab{
	engine.GenreFantasy: {
		UpgradeText:         "Enhance",
		DowngradeText:       "Diminish",
		RepairText:          "Mend",
		SpecializeText:      "Enchant",
		MaxedOutText:        "Already legendary quality",
		InsufficientText:    "Insufficient gold",
		SuccessText:         "Enhancement complete",
		FailureText:         "Enhancement failed",
		CurrencyName:        "Gold",
		MaterialsName:       "Materials",
		SpeedBoostText:      "Swifter steeds",
		CargoBoostText:      "Expanded hold",
		DefenseBoostText:    "Stronger walls",
		HealingBoostText:    "Better healers",
		NavigationBoostText: "Keener scouts",
	},
	engine.GenreScifi: {
		UpgradeText:         "Upgrade",
		DowngradeText:       "Downgrade",
		RepairText:          "Repair",
		SpecializeText:      "Optimize",
		MaxedOutText:        "Maximum specification reached",
		InsufficientText:    "Insufficient credits",
		SuccessText:         "Upgrade successful",
		FailureText:         "Upgrade failed",
		CurrencyName:        "Credits",
		MaterialsName:       "Components",
		SpeedBoostText:      "Improved propulsion",
		CargoBoostText:      "Expanded bay",
		DefenseBoostText:    "Enhanced shields",
		HealingBoostText:    "Advanced med-bay",
		NavigationBoostText: "Better sensors",
	},
	engine.GenreHorror: {
		UpgradeText:         "Reinforce",
		DowngradeText:       "Strip",
		RepairText:          "Patch",
		SpecializeText:      "Fortify",
		MaxedOutText:        "Can't improve further",
		InsufficientText:    "Not enough parts",
		SuccessText:         "Reinforcement done",
		FailureText:         "Reinforcement failed",
		CurrencyName:        "Parts",
		MaterialsName:       "Scrap",
		SpeedBoostText:      "Faster escape",
		CargoBoostText:      "More storage",
		DefenseBoostText:    "Better armor",
		HealingBoostText:    "Better first aid",
		NavigationBoostText: "Clearer routes",
	},
	engine.GenreCyberpunk: {
		UpgradeText:         "Mod",
		DowngradeText:       "Unmod",
		RepairText:          "Fix",
		SpecializeText:      "Overclock",
		MaxedOutText:        "Mil-spec maximum",
		InsufficientText:    "Not enough creds",
		SuccessText:         "Mod installed",
		FailureText:         "Mod rejected",
		CurrencyName:        "Creds",
		MaterialsName:       "Tech",
		SpeedBoostText:      "Faster processing",
		CargoBoostText:      "More bandwidth",
		DefenseBoostText:    "Harder ICE",
		HealingBoostText:    "Better trauma care",
		NavigationBoostText: "Net integration",
	},
	engine.GenrePostapoc: {
		UpgradeText:         "Weld",
		DowngradeText:       "Unbolt",
		RepairText:          "Jury-rig",
		SpecializeText:      "Customize",
		MaxedOutText:        "Pre-war quality max",
		InsufficientText:    "Not enough scrap",
		SuccessText:         "Upgrade welded on",
		FailureText:         "Upgrade failed",
		CurrencyName:        "Scrap",
		MaterialsName:       "Salvage",
		SpeedBoostText:      "More horsepower",
		CargoBoostText:      "Bigger bed",
		DefenseBoostText:    "Thicker plating",
		HealingBoostText:    "Better chem station",
		NavigationBoostText: "Clearer maps",
	},
}

// AllCustomizationVocabs returns all available customization vocabularies.
func AllCustomizationVocabs() map[engine.GenreID]*CustomizationVocab {
	return customizationVocabs
}

// AllModuleUpgradeVocabs returns all available module upgrade vocabularies.
func AllModuleUpgradeVocabs() map[engine.GenreID]*ModuleUpgradeVocab {
	return moduleUpgradeVocabs
}
