package trading

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// SupplyPostType identifies the kind of trading location.
type SupplyPostType int

const (
	// PostTypeMarket is a general trading hub with diverse inventory.
	PostTypeMarket SupplyPostType = iota
	// PostTypeOutpost is a small frontier trading location.
	PostTypeOutpost
	// PostTypeSpecialist is a location with focused inventory.
	PostTypeSpecialist
	// PostTypeBlackMarket offers rare items at premium prices.
	PostTypeBlackMarket
)

// AllPostTypes returns all supply post types.
func AllPostTypes() []SupplyPostType {
	return []SupplyPostType{
		PostTypeMarket,
		PostTypeOutpost,
		PostTypeSpecialist,
		PostTypeBlackMarket,
	}
}

// SupplyPost represents a procedurally generated trading location.
type SupplyPost struct {
	Name          string
	PostType      SupplyPostType
	Genre         engine.GenreID
	RegionID      int    // Used for seeding inventory
	Position      [2]int // Map coordinates
	Inventory     *Inventory
	PriceModifier float64 // Multiplier for all prices (1.0 = standard)
	Reputation    float64 // 0-1, affects prices and available items
	Description   string
}

// SupplyPostGenerator creates procedural supply posts from a seed.
type SupplyPostGenerator struct {
	gen   *seed.Generator
	genre engine.GenreID
}

// NewSupplyPostGenerator creates a new supply post generator.
func NewSupplyPostGenerator(masterSeed int64, genre engine.GenreID) *SupplyPostGenerator {
	return &SupplyPostGenerator{
		gen:   seed.NewGenerator(masterSeed, "supplypost"),
		genre: genre,
	}
}

// SetGenre changes the generator's genre.
func (g *SupplyPostGenerator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// Genre returns the current genre.
func (g *SupplyPostGenerator) Genre() engine.GenreID {
	return g.genre
}

// Generate creates a new supply post at the given position.
func (g *SupplyPostGenerator) Generate(x, y, regionID int) *SupplyPost {
	postType := g.selectPostType()

	post := &SupplyPost{
		PostType:      postType,
		Genre:         g.genre,
		RegionID:      regionID,
		Position:      [2]int{x, y},
		PriceModifier: g.generatePriceModifier(postType),
		Reputation:    0.5, // Start neutral
	}

	post.Name = g.generateName(postType)
	post.Description = g.generateDescription(postType)
	post.Inventory = g.generateInventory(post)

	return post
}

// selectPostType chooses a post type based on weighted probabilities.
func (g *SupplyPostGenerator) selectPostType() SupplyPostType {
	roll := g.gen.Float64()
	switch {
	case roll < 0.4:
		return PostTypeMarket
	case roll < 0.7:
		return PostTypeOutpost
	case roll < 0.9:
		return PostTypeSpecialist
	default:
		return PostTypeBlackMarket
	}
}

// generatePriceModifier creates a price multiplier based on post type.
func (g *SupplyPostGenerator) generatePriceModifier(postType SupplyPostType) float64 {
	baseModifiers := map[SupplyPostType]float64{
		PostTypeMarket:      1.0,
		PostTypeOutpost:     1.2,
		PostTypeSpecialist:  0.9,
		PostTypeBlackMarket: 1.5,
	}
	base := baseModifiers[postType]
	// Add ±10% variation
	return base + g.gen.RangeFloat64(-0.1, 0.1)
}

// generateName creates a procedural name for the supply post.
func (g *SupplyPostGenerator) generateName(postType SupplyPostType) string {
	prefixes := g.getNamePrefixes()
	suffixes := g.getNameSuffixes(postType)

	prefix := seed.Choice(g.gen, prefixes)
	suffix := seed.Choice(g.gen, suffixes)

	return prefix + " " + suffix
}

// getNamePrefixes returns genre-appropriate name prefixes.
func (g *SupplyPostGenerator) getNamePrefixes() []string {
	prefixes := map[engine.GenreID][]string{
		engine.GenreFantasy: {
			"Golden", "Silver", "Iron", "Copper", "Emerald",
			"Northern", "Southern", "Eastern", "Western", "Central",
		},
		engine.GenreScifi: {
			"Alpha", "Beta", "Gamma", "Delta", "Epsilon",
			"Orbital", "Stellar", "Quantum", "Binary", "Void",
		},
		engine.GenreHorror: {
			"Last", "Final", "Broken", "Silent", "Dead",
			"Hidden", "Safe", "Lost", "Fallen", "Forsaken",
		},
		engine.GenreCyberpunk: {
			"Neon", "Chrome", "Digital", "Binary", "Cyber",
			"Shadow", "Ghost", "Null", "Vector", "Grid",
		},
		engine.GenrePostapoc: {
			"Rusty", "Dusty", "Old", "New", "Last",
			"Reclaimed", "Salvaged", "Free", "Lone", "Hardy",
		},
	}
	if p, ok := prefixes[g.genre]; ok {
		return p
	}
	return prefixes[engine.GenreFantasy]
}

// getNameSuffixes returns genre-appropriate name suffixes by post type.
func (g *SupplyPostGenerator) getNameSuffixes(postType SupplyPostType) []string {
	suffixes := map[engine.GenreID]map[SupplyPostType][]string{
		engine.GenreFantasy: {
			PostTypeMarket:      {"Market", "Bazaar", "Trading Post", "Merchant Hall", "Exchange"},
			PostTypeOutpost:     {"Outpost", "Camp", "Rest Stop", "Waystation", "Lodge"},
			PostTypeSpecialist:  {"Smithy", "Apothecary", "Stables", "Armory", "Guild"},
			PostTypeBlackMarket: {"Shadow Market", "Thieves' Den", "Underground", "Hidden Cache", "Smugglers' Hold"},
		},
		engine.GenreScifi: {
			PostTypeMarket:      {"Station", "Hub", "Depot", "Exchange", "Terminal"},
			PostTypeOutpost:     {"Outpost", "Relay", "Beacon", "Array", "Node"},
			PostTypeSpecialist:  {"Tech Bay", "Med Station", "Armory", "Fuel Depot", "Research Lab"},
			PostTypeBlackMarket: {"Black Dock", "Shadow Port", "Gray Market", "Smuggler's Den", "Dead Drop"},
		},
		engine.GenreHorror: {
			PostTypeMarket:      {"Camp", "Settlement", "Compound", "Haven", "Refuge"},
			PostTypeOutpost:     {"Shelter", "Bunker", "Hideout", "Safe House", "Watchtower"},
			PostTypeSpecialist:  {"Med Tent", "Armory", "Garage", "Workshop", "Supply Cache"},
			PostTypeBlackMarket: {"Black Market", "Raider Trade", "Scav Den", "Underground", "Shadow Trade"},
		},
		engine.GenreCyberpunk: {
			PostTypeMarket:      {"Mart", "Depot", "Hub", "Exchange", "Terminal"},
			PostTypeOutpost:     {"Node", "Link", "Access Point", "Relay", "Junction"},
			PostTypeSpecialist:  {"Clinic", "Mod Shop", "Arms Dealer", "Fixer's Den", "Ripper Doc"},
			PostTypeBlackMarket: {"Black ICE", "Shadow Net", "Dark Market", "Smuggler's Deck", "Gray Zone"},
		},
		engine.GenrePostapoc: {
			PostTypeMarket:      {"Bazaar", "Swap Meet", "Trading Post", "Market", "Exchange"},
			PostTypeOutpost:     {"Outpost", "Camp", "Shelter", "Post", "Station"},
			PostTypeSpecialist:  {"Mechanic", "Med Station", "Armory", "Fuel Depot", "Parts Shop"},
			PostTypeBlackMarket: {"Black Market", "Smuggler's Den", "Underground", "Shadow Trade", "Raider Trade"},
		},
	}

	if genreSuffixes, ok := suffixes[g.genre]; ok {
		if typeSuffixes, ok := genreSuffixes[postType]; ok {
			return typeSuffixes
		}
	}
	return []string{"Trading Post"}
}

// generateDescription creates a description for the supply post.
func (g *SupplyPostGenerator) generateDescription(postType SupplyPostType) string {
	descriptions := g.getDescriptions(postType)
	return seed.Choice(g.gen, descriptions)
}

// getDescriptions returns genre-appropriate descriptions by post type.
func (g *SupplyPostGenerator) getDescriptions(postType SupplyPostType) []string {
	descriptions := map[engine.GenreID]map[SupplyPostType][]string{
		engine.GenreFantasy: {
			PostTypeMarket:      {"A bustling market filled with merchants.", "Traders from distant lands gather here.", "A well-established trading hub."},
			PostTypeOutpost:     {"A small frontier settlement.", "A lonely outpost on the road.", "A simple waystation for travelers."},
			PostTypeSpecialist:  {"A skilled craftsperson runs this shop.", "Specialized goods at fair prices.", "Known for quality merchandise."},
			PostTypeBlackMarket: {"Dealings best kept quiet.", "Ask no questions here.", "For those who know what they seek."},
		},
		engine.GenreScifi: {
			PostTypeMarket:      {"A major trading station.", "Ships from many systems dock here.", "Credits flow freely in this hub."},
			PostTypeOutpost:     {"A remote frontier station.", "Basic supplies available.", "An isolated relay post."},
			PostTypeSpecialist:  {"Specialized equipment dealers.", "Technical experts provide services.", "Quality components guaranteed."},
			PostTypeBlackMarket: {"Encrypted transactions only.", "No questions asked.", "Unofficial channels."},
		},
		engine.GenreHorror: {
			PostTypeMarket:      {"Survivors barter what they can.", "A fortified trading compound.", "Safety in numbers here."},
			PostTypeOutpost:     {"A desperate shelter.", "Basic necessities only.", "Barely holding together."},
			PostTypeSpecialist:  {"Precious medical supplies.", "Hard-to-find equipment.", "Worth the detour."},
			PostTypeBlackMarket: {"Raiders sometimes trade here.", "Don't ask where it came from.", "Dangerous deals."},
		},
		engine.GenreCyberpunk: {
			PostTypeMarket:      {"A busy corporate plaza.", "Legal commerce... mostly.", "Standard consumer goods."},
			PostTypeOutpost:     {"A neighborhood node.", "Local access point.", "Small-time traders."},
			PostTypeSpecialist:  {"Expert modifications available.", "Premium tech services.", "Quality guaranteed."},
			PostTypeBlackMarket: {"Off the grid.", "Untraceable transactions.", "The fixer knows a guy."},
		},
		engine.GenrePostapoc: {
			PostTypeMarket:      {"Scavengers gather to trade.", "A hub of bartering activity.", "Goods from the old world."},
			PostTypeOutpost:     {"A hardscrabble trading post.", "Basic supplies only.", "Built from salvage."},
			PostTypeSpecialist:  {"Hard-to-find skills here.", "Pre-war technology repaired.", "Worth the journey."},
			PostTypeBlackMarket: {"No tribe laws here.", "Dangerous but profitable.", "Bring something valuable."},
		},
	}

	if genreDesc, ok := descriptions[g.genre]; ok {
		if typeDesc, ok := genreDesc[postType]; ok {
			return typeDesc
		}
	}
	return []string{"A trading location."}
}

// generateInventory creates the initial inventory for a supply post.
func (g *SupplyPostGenerator) generateInventory(post *SupplyPost) *Inventory {
	inv := NewInventory(g.genre)

	// Base inventory size by post type
	sizes := map[SupplyPostType]int{
		PostTypeMarket:      15,
		PostTypeOutpost:     8,
		PostTypeSpecialist:  10,
		PostTypeBlackMarket: 12,
	}
	size := sizes[post.PostType] + g.gen.Range(-2, 3)

	// Generate items
	itemGen := NewItemGenerator(g.gen.Master()+int64(post.RegionID), g.genre)
	for i := 0; i < size; i++ {
		item := itemGen.Generate(post.PostType)
		inv.AddItem(item)
	}

	return inv
}

// PostTypeName returns the genre-appropriate name for a post type.
func PostTypeName(postType SupplyPostType, genre engine.GenreID) string {
	names := map[engine.GenreID]map[SupplyPostType]string{
		engine.GenreFantasy: {
			PostTypeMarket:      "Market",
			PostTypeOutpost:     "Trading Post",
			PostTypeSpecialist:  "Guild Shop",
			PostTypeBlackMarket: "Shadow Market",
		},
		engine.GenreScifi: {
			PostTypeMarket:      "Space Dock",
			PostTypeOutpost:     "Relay Station",
			PostTypeSpecialist:  "Tech Bay",
			PostTypeBlackMarket: "Gray Market",
		},
		engine.GenreHorror: {
			PostTypeMarket:      "Survivor Camp",
			PostTypeOutpost:     "Safe House",
			PostTypeSpecialist:  "Specialist",
			PostTypeBlackMarket: "Black Market",
		},
		engine.GenreCyberpunk: {
			PostTypeMarket:      "Black Market",
			PostTypeOutpost:     "Data Node",
			PostTypeSpecialist:  "Mod Shop",
			PostTypeBlackMarket: "Shadow Net",
		},
		engine.GenrePostapoc: {
			PostTypeMarket:      "Scrap Bazaar",
			PostTypeOutpost:     "Trading Post",
			PostTypeSpecialist:  "Workshop",
			PostTypeBlackMarket: "Underground",
		},
	}

	if genreNames, ok := names[genre]; ok {
		if name, ok := genreNames[postType]; ok {
			return name
		}
	}
	return "Trading Post"
}

// TypeName returns the genre-appropriate type name for this post.
func (sp *SupplyPost) TypeName() string {
	return PostTypeName(sp.PostType, sp.Genre)
}

// AdjustedPrice calculates the final price for an item at this post.
func (sp *SupplyPost) AdjustedPrice(basePrice float64, isSelling bool) float64 {
	price := basePrice * sp.PriceModifier

	// Reputation affects prices
	// Good rep (>0.5) = better prices
	// Bad rep (<0.5) = worse prices
	repModifier := 1.0 + (0.5-sp.Reputation)*0.2
	price *= repModifier

	// Selling to the post gets lower prices
	if isSelling {
		price *= 0.6
	}

	return price
}

// UpdateReputation changes the post's reputation based on player actions.
func (sp *SupplyPost) UpdateReputation(delta float64) {
	sp.Reputation += delta
	if sp.Reputation < 0 {
		sp.Reputation = 0
	}
	if sp.Reputation > 1 {
		sp.Reputation = 1
	}
}

// ReputationStatus returns a text description of the current reputation.
func (sp *SupplyPost) ReputationStatus() string {
	switch {
	case sp.Reputation >= 0.8:
		return "Honored"
	case sp.Reputation >= 0.6:
		return "Friendly"
	case sp.Reputation >= 0.4:
		return "Neutral"
	case sp.Reputation >= 0.2:
		return "Suspicious"
	default:
		return "Hostile"
	}
}
