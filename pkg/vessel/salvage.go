package vessel

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// SalvageType identifies types of salvageable wrecks.
type SalvageType int

const (
	// SalvageWreck is a destroyed vessel wreck.
	SalvageWreck SalvageType = iota
	// SalvageAbandoned is an abandoned vessel.
	SalvageAbandoned
	// SalvageCrash is a crash site.
	SalvageCrash
	// SalvageRuin is a ruined structure with salvage.
	SalvageRuin
)

// AllSalvageTypes returns all salvage types.
func AllSalvageTypes() []SalvageType {
	return []SalvageType{SalvageWreck, SalvageAbandoned, SalvageCrash, SalvageRuin}
}

// SalvageItem represents an item obtainable from salvage.
type SalvageItem struct {
	Name     string
	Weight   int
	Volume   int
	Quantity int
	Category CargoCategory
	Value    float64 // Trade value
}

// SalvageResult represents the outcome of a salvage attempt.
type SalvageResult struct {
	Success bool
	Items   []SalvageItem
	Message string
	Danger  bool // Encountered danger during salvage
}

// SalvageSite represents a location that can be salvaged.
type SalvageSite struct {
	ID          int
	Type        SalvageType
	Salvaged    bool
	Richness    float64 // 0.0-1.0, affects loot quality
	DangerLevel float64 // 0.0-1.0, chance of encounter
}

// SalvageManager handles salvage operations.
type SalvageManager struct {
	genre engine.GenreID
	gen   *seed.Generator
}

// NewSalvageManager creates a new salvage manager.
func NewSalvageManager(masterSeed int64, genre engine.GenreID) *SalvageManager {
	return &SalvageManager{
		genre: genre,
		gen:   seed.NewGenerator(masterSeed, "salvage"),
	}
}

// SetGenre changes the salvage vocabulary.
func (sm *SalvageManager) SetGenre(genre engine.GenreID) {
	sm.genre = genre
}

// GenerateSite creates a new salvage site.
func (sm *SalvageManager) GenerateSite(id int) *SalvageSite {
	types := AllSalvageTypes()
	sType := types[sm.gen.Intn(len(types))]

	return &SalvageSite{
		ID:          id,
		Type:        sType,
		Salvaged:    false,
		Richness:    0.3 + sm.gen.Float64()*0.7, // 0.3-1.0
		DangerLevel: sm.gen.Float64() * 0.6,     // 0.0-0.6
	}
}

// AttemptSalvage tries to salvage a site.
func (sm *SalvageManager) AttemptSalvage(site *SalvageSite) SalvageResult {
	if site.Salvaged {
		return SalvageResult{
			Success: false,
			Message: "This site has already been salvaged",
		}
	}

	// Check for danger
	danger := sm.gen.Float64() < site.DangerLevel
	if danger {
		// Still get some items, but flag danger
		items := sm.generateLoot(site, 0.5) // Reduced loot
		site.Salvaged = true
		return SalvageResult{
			Success: true,
			Items:   items,
			Message: sm.dangerMessage(),
			Danger:  true,
		}
	}

	// Successful salvage
	items := sm.generateLoot(site, 1.0)
	site.Salvaged = true

	return SalvageResult{
		Success: true,
		Items:   items,
		Message: sm.successMessage(site.Type),
		Danger:  false,
	}
}

// generateLoot creates salvage items based on site richness.
func (sm *SalvageManager) generateLoot(site *SalvageSite, multiplier float64) []SalvageItem {
	var items []SalvageItem

	// Number of items based on richness
	itemCount := 1 + sm.gen.Intn(int(site.Richness*4)+1)
	effectiveRichness := site.Richness * multiplier

	for i := 0; i < itemCount; i++ {
		item := sm.generateItem(effectiveRichness)
		items = append(items, item)
	}

	return items
}

// generateItem creates a single salvage item.
func (sm *SalvageManager) generateItem(richness float64) SalvageItem {
	// Choose category based on probability
	roll := sm.gen.Float64()
	var cat CargoCategory
	switch {
	case roll < 0.4:
		cat = CargoRepair // Most common: repair materials
	case roll < 0.7:
		cat = CargoSupplies
	case roll < 0.85:
		cat = CargoTrade
	case roll < 0.95:
		cat = CargoMedical
	default:
		cat = CargoSpecial
	}

	name := sm.itemName(cat)
	quantity := 1 + sm.gen.Intn(int(richness*5)+1)
	weight := 1 + sm.gen.Intn(3)
	volume := weight
	value := float64(weight*quantity) * (0.5 + richness*2)

	return SalvageItem{
		Name:     name,
		Weight:   weight,
		Volume:   volume,
		Quantity: quantity,
		Category: cat,
		Value:    value,
	}
}

// itemName returns a genre-appropriate item name.
func (sm *SalvageManager) itemName(cat CargoCategory) string {
	names, ok := salvageItemNames[sm.genre]
	if !ok {
		names = salvageItemNames[engine.GenreFantasy]
	}
	catNames := names[cat]
	if len(catNames) == 0 {
		return "Salvage"
	}
	return seed.Choice(sm.gen, catNames)
}

var salvageItemNames = map[engine.GenreID]map[CargoCategory][]string{
	engine.GenreFantasy: {
		CargoRepair:   {"Timber", "Iron Nails", "Rope", "Canvas", "Wheel Spokes"},
		CargoSupplies: {"Dried Meat", "Grain Sack", "Water Skin", "Lamp Oil"},
		CargoTrade:    {"Cloth Bolt", "Spices", "Silver Coin", "Gemstone"},
		CargoMedical:  {"Healing Herbs", "Bandages", "Salve"},
		CargoSpecial:  {"Ancient Map", "Mysterious Amulet", "Sealed Letter"},
	},
	engine.GenreScifi: {
		CargoRepair:   {"Hull Plates", "Power Cells", "Conduits", "Circuitry", "Sealant"},
		CargoSupplies: {"Rations", "Water Packs", "Oxygen Tanks", "Battery Packs"},
		CargoTrade:    {"Rare Minerals", "Data Chips", "Tech Components", "Alloys"},
		CargoMedical:  {"Med-Kits", "Stims", "Bio-Gel"},
		CargoSpecial:  {"Star Map", "AI Core", "Black Box"},
	},
	engine.GenreHorror: {
		CargoRepair:   {"Spare Parts", "Duct Tape", "Wire", "Motor Oil", "Tires"},
		CargoSupplies: {"Canned Food", "Water Bottles", "Gasoline", "Batteries"},
		CargoTrade:    {"Ammo Box", "Tools", "Generator Parts", "Solar Panels"},
		CargoMedical:  {"First Aid Kit", "Antibiotics", "Painkillers"},
		CargoSpecial:  {"Survivor's Journal", "Security Keycard", "Radio Parts"},
	},
	engine.GenreCyberpunk: {
		CargoRepair:   {"Tech Parts", "Fiber Optics", "Neural Chips", "Actuators"},
		CargoSupplies: {"Nutrient Packs", "Filtered Water", "Power Cells", "Coolant"},
		CargoTrade:    {"Data Shards", "Crypto Keys", "Black ICE", "Cyberware"},
		CargoMedical:  {"Trauma Kits", "Boosters", "Neural Stabilizers"},
		CargoSpecial:  {"Corporate Intel", "Prototype Tech", "Netrunner Deck"},
	},
	engine.GenrePostapoc: {
		CargoRepair:   {"Scrap Metal", "Rubber", "Chain", "Bolts", "Sheet Metal"},
		CargoSupplies: {"Canned Goods", "Purified Water", "Diesel", "Propane"},
		CargoTrade:    {"Pre-War Tech", "Ammunition", "Tools", "Seeds"},
		CargoMedical:  {"Rad-Away", "Stimpaks", "Antibiotics"},
		CargoSpecial:  {"Pre-War Map", "Vault Key", "Working Radio"},
	},
}

// successMessage returns a genre-appropriate success message.
func (sm *SalvageManager) successMessage(sType SalvageType) string {
	messages := map[engine.GenreID]map[SalvageType]string{
		engine.GenreFantasy: {
			SalvageWreck:     "You strip the wrecked wagon of useful materials.",
			SalvageAbandoned: "The abandoned camp yields supplies.",
			SalvageCrash:     "You salvage what you can from the wreckage.",
			SalvageRuin:      "Among the ruins, you find useful items.",
		},
		engine.GenreScifi: {
			SalvageWreck:     "You extract components from the derelict hull.",
			SalvageAbandoned: "The abandoned station has salvageable tech.",
			SalvageCrash:     "The crash site yields valuable parts.",
			SalvageRuin:      "You recover equipment from the ruins.",
		},
		engine.GenreHorror: {
			SalvageWreck:     "You quickly grab supplies from the wrecked vehicle.",
			SalvageAbandoned: "The abandoned building has supplies.",
			SalvageCrash:     "The crash site has useful items.",
			SalvageRuin:      "You scavenge what remains in the ruins.",
		},
		engine.GenreCyberpunk: {
			SalvageWreck:     "You jack tech from the totaled vehicle.",
			SalvageAbandoned: "The abandoned complex has gear to grab.",
			SalvageCrash:     "You strip the crash site for parts.",
			SalvageRuin:      "The ruined building yields salvage.",
		},
		engine.GenrePostapoc: {
			SalvageWreck:     "You strip the wreck for scrap and parts.",
			SalvageAbandoned: "The abandoned shelter has supplies.",
			SalvageCrash:     "You salvage from the crash site.",
			SalvageRuin:      "The ruins hold valuable finds.",
		},
	}

	if genreMsgs, ok := messages[sm.genre]; ok {
		if msg, ok := genreMsgs[sType]; ok {
			return msg
		}
	}
	return "You salvage useful items."
}

// dangerMessage returns a genre-appropriate danger message.
func (sm *SalvageManager) dangerMessage() string {
	messages := map[engine.GenreID]string{
		engine.GenreFantasy:   "Bandits ambush you during salvage!",
		engine.GenreScifi:     "Hostiles detected during salvage!",
		engine.GenreHorror:    "Infected attack during salvage!",
		engine.GenreCyberpunk: "Gangers jump you during salvage!",
		engine.GenrePostapoc:  "Raiders attack during salvage!",
	}
	if msg, ok := messages[sm.genre]; ok {
		return msg
	}
	return "You're attacked during salvage!"
}

// SalvageTypeName returns the genre-appropriate name for a salvage type.
func SalvageTypeName(st SalvageType, genre engine.GenreID) string {
	names, ok := salvageTypeNames[genre]
	if !ok {
		names = salvageTypeNames[engine.GenreFantasy]
	}
	return names[st]
}

var salvageTypeNames = map[engine.GenreID]map[SalvageType]string{
	engine.GenreFantasy: {
		SalvageWreck:     "Wrecked Wagon",
		SalvageAbandoned: "Abandoned Camp",
		SalvageCrash:     "Collapsed Cart",
		SalvageRuin:      "Crumbled Ruin",
	},
	engine.GenreScifi: {
		SalvageWreck:     "Derelict Ship",
		SalvageAbandoned: "Abandoned Station",
		SalvageCrash:     "Crash Site",
		SalvageRuin:      "Ruined Outpost",
	},
	engine.GenreHorror: {
		SalvageWreck:     "Wrecked Vehicle",
		SalvageAbandoned: "Abandoned Building",
		SalvageCrash:     "Crash Site",
		SalvageRuin:      "Collapsed Structure",
	},
	engine.GenreCyberpunk: {
		SalvageWreck:     "Totaled Vehicle",
		SalvageAbandoned: "Abandoned Complex",
		SalvageCrash:     "Crash Site",
		SalvageRuin:      "Ruined Building",
	},
	engine.GenrePostapoc: {
		SalvageWreck:     "Vehicle Wreck",
		SalvageAbandoned: "Abandoned Shelter",
		SalvageCrash:     "Crash Site",
		SalvageRuin:      "Collapsed Ruin",
	},
}

// AddSalvageToHold adds salvage items to a cargo hold.
func AddSalvageToHold(items []SalvageItem, hold *CargoHold) (added, failed int) {
	for _, item := range items {
		if hold.AddWithVolume(item.Name, item.Weight, item.Volume, item.Quantity, item.Category) {
			added++
		} else {
			failed++
		}
	}
	return added, failed
}
