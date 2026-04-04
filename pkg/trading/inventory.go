package trading

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// ItemCategory represents the category of a tradeable item.
type ItemCategory int

const (
	// CategoryFood represents food supplies.
	CategoryFood ItemCategory = iota
	// CategoryWater represents water supplies.
	CategoryWater
	// CategoryFuel represents fuel supplies.
	CategoryFuel
	// CategoryMedicine represents medical supplies.
	CategoryMedicine
	// CategoryParts represents repair parts.
	CategoryParts
	// CategoryTrade represents trade goods.
	CategoryTrade
	// CategoryRare represents rare/special items.
	CategoryRare
)

// AllItemCategories returns all item categories.
func AllItemCategories() []ItemCategory {
	return []ItemCategory{
		CategoryFood,
		CategoryWater,
		CategoryFuel,
		CategoryMedicine,
		CategoryParts,
		CategoryTrade,
		CategoryRare,
	}
}

// Item represents a tradeable item at a supply post.
type Item struct {
	Name        string
	Description string
	Category    ItemCategory
	BasePrice   float64
	Quantity    int
	Quality     float64 // 0-1, affects effectiveness
	Genre       engine.GenreID
}

// Inventory holds the items available at a supply post.
type Inventory struct {
	Items []*Item
	Genre engine.GenreID
}

// NewInventory creates a new empty inventory.
func NewInventory(genre engine.GenreID) *Inventory {
	return &Inventory{
		Items: make([]*Item, 0),
		Genre: genre,
	}
}

// AddItem adds an item to the inventory.
func (inv *Inventory) AddItem(item *Item) {
	// Check if item already exists (merge quantities)
	for _, existing := range inv.Items {
		if existing.Name == item.Name && existing.Category == item.Category {
			existing.Quantity += item.Quantity
			return
		}
	}
	inv.Items = append(inv.Items, item)
}

// RemoveItem reduces or removes an item from inventory.
func (inv *Inventory) RemoveItem(name string, quantity int) bool {
	for i, item := range inv.Items {
		if item.Name == name {
			if item.Quantity < quantity {
				return false
			}
			item.Quantity -= quantity
			if item.Quantity <= 0 {
				inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
			}
			return true
		}
	}
	return false
}

// GetItem returns an item by name.
func (inv *Inventory) GetItem(name string) *Item {
	for _, item := range inv.Items {
		if item.Name == name {
			return item
		}
	}
	return nil
}

// GetByCategory returns all items in a category.
func (inv *Inventory) GetByCategory(category ItemCategory) []*Item {
	var result []*Item
	for _, item := range inv.Items {
		if item.Category == category {
			result = append(result, item)
		}
	}
	return result
}

// ItemCount returns the total number of unique items.
func (inv *Inventory) ItemCount() int {
	return len(inv.Items)
}

// TotalValue returns the total value of all items.
func (inv *Inventory) TotalValue() float64 {
	total := 0.0
	for _, item := range inv.Items {
		total += item.BasePrice * float64(item.Quantity)
	}
	return total
}

// ItemGenerator creates procedural items from a seed.
type ItemGenerator struct {
	gen   *seed.Generator
	genre engine.GenreID
}

// NewItemGenerator creates a new item generator.
func NewItemGenerator(masterSeed int64, genre engine.GenreID) *ItemGenerator {
	return &ItemGenerator{
		gen:   seed.NewGenerator(masterSeed, "items"),
		genre: genre,
	}
}

// SetGenre changes the generator's genre.
func (g *ItemGenerator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// Generate creates a new procedural item appropriate for the post type.
func (g *ItemGenerator) Generate(postType SupplyPostType) *Item {
	category := g.selectCategory(postType)

	item := &Item{
		Category: category,
		Genre:    g.genre,
		Quality:  g.generateQuality(postType),
		Quantity: g.generateQuantity(category),
	}

	item.Name = g.generateName(category)
	item.Description = g.generateDescription(category)
	item.BasePrice = g.generatePrice(category, item.Quality)

	return item
}

// selectCategory chooses an item category based on post type.
func (g *ItemGenerator) selectCategory(postType SupplyPostType) ItemCategory {
	weights := map[SupplyPostType][]float64{
		PostTypeMarket:      {0.2, 0.2, 0.15, 0.15, 0.15, 0.1, 0.05},   // Balanced
		PostTypeOutpost:     {0.3, 0.3, 0.2, 0.1, 0.05, 0.05, 0.0},     // Basics focus
		PostTypeSpecialist:  {0.05, 0.05, 0.1, 0.3, 0.3, 0.1, 0.1},     // Parts/meds
		PostTypeBlackMarket: {0.05, 0.05, 0.1, 0.15, 0.15, 0.25, 0.25}, // Trade/rare
	}

	w := weights[postType]
	categories := AllItemCategories()
	return seed.WeightedChoice(g.gen, categories, w)
}

// generateQuality creates a quality value based on post type.
func (g *ItemGenerator) generateQuality(postType SupplyPostType) float64 {
	baseQuality := map[SupplyPostType]float64{
		PostTypeMarket:      0.6,
		PostTypeOutpost:     0.4,
		PostTypeSpecialist:  0.8,
		PostTypeBlackMarket: 0.7,
	}
	base := baseQuality[postType]
	return clampFloat(base+g.gen.RangeFloat64(-0.2, 0.2), 0.1, 1.0)
}

// generateQuantity creates a quantity based on category.
func (g *ItemGenerator) generateQuantity(category ItemCategory) int {
	baseQuantities := map[ItemCategory]int{
		CategoryFood:     10,
		CategoryWater:    10,
		CategoryFuel:     8,
		CategoryMedicine: 5,
		CategoryParts:    3,
		CategoryTrade:    2,
		CategoryRare:     1,
	}
	base := baseQuantities[category]
	return base + g.gen.Range(-base/2, base/2)
}

// generateName creates a procedural item name.
func (g *ItemGenerator) generateName(category ItemCategory) string {
	names := g.getCategoryNames(category)
	baseName := seed.Choice(g.gen, names)

	// Sometimes add a qualifier
	if g.gen.Chance(0.3) {
		qualifiers := g.getQualifiers()
		qualifier := seed.Choice(g.gen, qualifiers)
		return qualifier + " " + baseName
	}
	return baseName
}

// getCategoryNames returns genre-appropriate item names by category.
func (g *ItemGenerator) getCategoryNames(category ItemCategory) []string {
	names := map[engine.GenreID]map[ItemCategory][]string{
		engine.GenreFantasy: {
			CategoryFood:     {"Bread", "Dried Meat", "Cheese", "Fruit", "Grain", "Salted Fish"},
			CategoryWater:    {"Water Flask", "Wine", "Ale", "Spring Water", "Tea"},
			CategoryFuel:     {"Torch", "Lamp Oil", "Firewood", "Coal", "Charcoal"},
			CategoryMedicine: {"Healing Herbs", "Potion", "Bandages", "Salve", "Tonic"},
			CategoryParts:    {"Rope", "Nails", "Leather", "Wheel", "Axle"},
			CategoryTrade:    {"Silk", "Spices", "Gems", "Gold Dust", "Artifacts"},
			CategoryRare:     {"Enchanted Ring", "Ancient Scroll", "Magic Crystal", "Blessed Amulet"},
		},
		engine.GenreScifi: {
			CategoryFood:     {"Ration Pack", "Nutrient Bar", "Protein Gel", "Synth-Meat", "Hydro-Veg"},
			CategoryWater:    {"Purified Water", "Electrolyte Mix", "Condensate", "Ice Core"},
			CategoryFuel:     {"Fuel Cell", "Plasma Core", "Deuterium", "Power Pack", "Solar Cell"},
			CategoryMedicine: {"Med-Kit", "Stim-Pack", "Nano-Healers", "Bio-Gel", "Rad-Away"},
			CategoryParts:    {"Circuit Board", "Power Conduit", "Hull Patch", "Servo Motor", "Sensor Array"},
			CategoryTrade:    {"Data Chip", "Rare Alloy", "Alien Artifact", "Star Maps", "Tech Blueprint"},
			CategoryRare:     {"AI Core", "Quantum Processor", "Alien Tech", "Prototype Module"},
		},
		engine.GenreHorror: {
			CategoryFood:     {"Canned Food", "MRE", "Jerky", "Crackers", "Preserved Fruit"},
			CategoryWater:    {"Bottled Water", "Water Filter", "Purification Tabs", "Canteen"},
			CategoryFuel:     {"Gasoline", "Diesel", "Propane", "Battery", "Generator Fuel"},
			CategoryMedicine: {"First Aid Kit", "Antibiotics", "Painkillers", "Bandages", "Disinfectant"},
			CategoryParts:    {"Scrap Metal", "Wire", "Tools", "Duct Tape", "Spare Tire"},
			CategoryTrade:    {"Cigarettes", "Liquor", "Ammo", "Batteries", "Medicine"},
			CategoryRare:     {"Military Gear", "Vaccine", "Radio Equipment", "Working Vehicle"},
		},
		engine.GenreCyberpunk: {
			CategoryFood:     {"Synth-Food", "Kibble", "Protein Bar", "Nutri-Paste", "Street Food"},
			CategoryWater:    {"Filtered Water", "Energy Drink", "Stim-Juice", "Pure Water"},
			CategoryFuel:     {"Power Cell", "E-Charge", "Hydrogen Cell", "Battery Pack"},
			CategoryMedicine: {"Med-Stim", "Pain Block", "Trauma Kit", "Bio-Gel", "Cyberware Cleaner"},
			CategoryParts:    {"Chrome Parts", "Neural Link", "Deck Upgrade", "ICE Chip", "Optic Mod"},
			CategoryTrade:    {"Data Shard", "Black ICE", "Corporate Intel", "Stolen Goods", "Fake IDs"},
			CategoryRare:     {"Military Cyberware", "Black Ops Gear", "Experimental Tech", "AI Fragment"},
		},
		engine.GenrePostapoc: {
			CategoryFood:     {"Canned Goods", "Dried Food", "Mutfruit", "Rad-Meat", "Preserved Rations"},
			CategoryWater:    {"Clean Water", "Rad-Free Water", "Purifier", "Dew Collector"},
			CategoryFuel:     {"Guzzoline", "Bio-Diesel", "Ethanol", "Fusion Core", "Solar Cell"},
			CategoryMedicine: {"Rad-X", "Stimpak", "Med-X", "Antibiotics", "Blood Pack"},
			CategoryParts:    {"Scrap", "Salvage", "Pre-War Parts", "Circuitry", "Gear"},
			CategoryTrade:    {"Caps", "Pre-War Money", "Ammo", "Chems", "Tech"},
			CategoryRare:     {"Pre-War Tech", "Vault Gear", "Military Hardware", "Clean Seed"},
		},
	}

	if genreNames, ok := names[g.genre]; ok {
		if catNames, ok := genreNames[category]; ok {
			return catNames
		}
	}
	return []string{"Item"}
}

// getQualifiers returns genre-appropriate item qualifiers.
func (g *ItemGenerator) getQualifiers() []string {
	qualifiers := map[engine.GenreID][]string{
		engine.GenreFantasy:   {"Fine", "Common", "Crude", "Exotic", "Ancient", "Blessed"},
		engine.GenreScifi:     {"Standard", "Military", "Civilian", "Prototype", "Surplus"},
		engine.GenreHorror:    {"Salvaged", "Damaged", "Pristine", "Expired", "Military"},
		engine.GenreCyberpunk: {"Street", "Corporate", "Military", "Bootleg", "Custom"},
		engine.GenrePostapoc:  {"Salvaged", "Pre-War", "Irradiated", "Clean", "Homemade"},
	}
	if q, ok := qualifiers[g.genre]; ok {
		return q
	}
	return []string{"Standard"}
}

// generateDescription creates a procedural item description.
func (g *ItemGenerator) generateDescription(category ItemCategory) string {
	descriptions := g.getCategoryDescriptions(category)
	return seed.Choice(g.gen, descriptions)
}

// getCategoryDescriptions returns descriptions by category.
func (g *ItemGenerator) getCategoryDescriptions(category ItemCategory) []string {
	descriptions := map[ItemCategory][]string{
		CategoryFood:     {"Provides sustenance for the journey.", "Nutritious and filling.", "Essential for survival."},
		CategoryWater:    {"Clean drinking water.", "Prevents dehydration.", "Life-giving liquid."},
		CategoryFuel:     {"Keeps things running.", "Powers your transport.", "Essential for travel."},
		CategoryMedicine: {"Heals wounds and illness.", "Medical supplies.", "Could save a life."},
		CategoryParts:    {"Useful for repairs.", "Keeps your transport running.", "Valuable components."},
		CategoryTrade:    {"Valuable for trading.", "Sought after goods.", "Worth its weight."},
		CategoryRare:     {"Exceptionally rare find.", "Highly valuable.", "One of a kind."},
	}
	if d, ok := descriptions[category]; ok {
		return d
	}
	return []string{"A tradeable item."}
}

// generatePrice creates a base price based on category and quality.
func (g *ItemGenerator) generatePrice(category ItemCategory, quality float64) float64 {
	basePrices := map[ItemCategory]float64{
		CategoryFood:     5.0,
		CategoryWater:    5.0,
		CategoryFuel:     8.0,
		CategoryMedicine: 15.0,
		CategoryParts:    12.0,
		CategoryTrade:    25.0,
		CategoryRare:     50.0,
	}
	base := basePrices[category]
	// Quality affects price significantly
	return base * (0.5 + quality)
}

// CategoryName returns the genre-appropriate name for a category.
func CategoryName(category ItemCategory, genre engine.GenreID) string {
	names := map[engine.GenreID]map[ItemCategory]string{
		engine.GenreFantasy: {
			CategoryFood:     "Provisions",
			CategoryWater:    "Drink",
			CategoryFuel:     "Supplies",
			CategoryMedicine: "Remedies",
			CategoryParts:    "Materials",
			CategoryTrade:    "Valuables",
			CategoryRare:     "Treasures",
		},
		engine.GenreScifi: {
			CategoryFood:     "Rations",
			CategoryWater:    "Hydration",
			CategoryFuel:     "Power",
			CategoryMedicine: "Medical",
			CategoryParts:    "Components",
			CategoryTrade:    "Cargo",
			CategoryRare:     "Artifacts",
		},
		engine.GenreHorror: {
			CategoryFood:     "Food",
			CategoryWater:    "Water",
			CategoryFuel:     "Fuel",
			CategoryMedicine: "Medical",
			CategoryParts:    "Parts",
			CategoryTrade:    "Trade Goods",
			CategoryRare:     "Rare Finds",
		},
		engine.GenreCyberpunk: {
			CategoryFood:     "Nutrients",
			CategoryWater:    "Hydration",
			CategoryFuel:     "Power",
			CategoryMedicine: "Meds",
			CategoryParts:    "Hardware",
			CategoryTrade:    "Data",
			CategoryRare:     "Black Market",
		},
		engine.GenrePostapoc: {
			CategoryFood:     "Grub",
			CategoryWater:    "Water",
			CategoryFuel:     "Juice",
			CategoryMedicine: "Meds",
			CategoryParts:    "Scrap",
			CategoryTrade:    "Barter",
			CategoryRare:     "Pre-War",
		},
	}

	if genreNames, ok := names[genre]; ok {
		if name, ok := genreNames[category]; ok {
			return name
		}
	}
	return "Items"
}

// clampFloat clamps a float64 value to the specified range.
func clampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
