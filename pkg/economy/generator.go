package economy

import (
	"fmt"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// Generator creates procedural economies
type Generator struct {
	gen      *seed.Generator
	genre    engine.GenreID
	marketID int
	routeID  int
}

// NewGenerator creates an economy generator with the given seed and genre
func NewGenerator(masterSeed int64, genre engine.GenreID) *Generator {
	return &Generator{
		gen:      seed.NewGenerator(masterSeed, "economy"),
		genre:    genre,
		marketID: 0,
		routeID:  0,
	}
}

// SetGenre updates the generator's active genre
func (g *Generator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// GenerateEconomy creates a complete economy with the specified number of markets
func (g *Generator) GenerateEconomy(marketCount int) *EconomyManager {
	marketCount = clampMarketCount(marketCount)
	manager := NewEconomyManager(g.genre)

	g.populateMarkets(manager, marketCount)
	g.connectMarkets(manager)

	return manager
}

// clampMarketCount ensures market count is within valid range [2, 10].
func clampMarketCount(count int) int {
	if count < 2 {
		return 2
	}
	if count > 10 {
		return 10
	}
	return count
}

// populateMarkets creates and adds markets to the manager.
func (g *Generator) populateMarkets(manager *EconomyManager, count int) {
	for i := 0; i < count; i++ {
		x := i*100 + g.gen.Intn(50)
		y := g.gen.Intn(100)
		market := g.GenerateMarket(x, y)
		manager.AddMarket(market)
	}
}

// connectMarkets creates trade routes between markets.
func (g *Generator) connectMarkets(manager *EconomyManager) {
	marketIDs := make([]string, 0, len(manager.Markets))
	for id := range manager.Markets {
		marketIDs = append(marketIDs, id)
	}

	for i := 0; i < len(marketIDs)-1; i++ {
		manager.AddRoute(g.GenerateRoute(marketIDs[i], marketIDs[i+1]))
		g.maybeAddSkipConnection(manager, marketIDs, i)
	}
}

// maybeAddSkipConnection adds a skip route with 30% probability.
func (g *Generator) maybeAddSkipConnection(manager *EconomyManager, marketIDs []string, i int) {
	if i < len(marketIDs)-2 && g.gen.Intn(10) < 3 {
		manager.AddRoute(g.GenerateRoute(marketIDs[i], marketIDs[i+2]))
	}
}

// GenerateMarket creates a single market
func (g *Generator) GenerateMarket(x, y int) *Market {
	g.marketID++

	name := g.generateMarketName()
	id := fmt.Sprintf("market_%d", g.marketID)

	market := NewMarket(id, name, x, y, g.genre)
	market.Prosperity = 0.3 + float64(g.gen.Intn(50))/100.0

	// Assign speciality and scarcity
	categories := AllCategories()
	market.Speciality = categories[g.gen.Intn(len(categories))]
	market.Scarcity = categories[g.gen.Intn(len(categories))]
	// Ensure they're different
	for market.Scarcity == market.Speciality {
		market.Scarcity = categories[g.gen.Intn(len(categories))]
	}

	// Add all trade goods
	goods := g.generateTradeGoods()
	for _, good := range goods {
		market.AddGood(good)
	}

	// Apply market characteristics
	market.ApplySpecialties()

	return market
}

// GenerateRoute creates a trade route between two markets
func (g *Generator) GenerateRoute(fromID, toID string) *TradeRoute {
	g.routeID++

	id := fmt.Sprintf("route_%d", g.routeID)
	distance := 50 + g.gen.Intn(100)

	route := NewTradeRoute(id, fromID, toID, distance, g.genre)
	route.Danger = float64(g.gen.Intn(40)) / 100.0
	route.Name = g.generateRouteName()

	// 20% chance of toll
	if g.gen.Intn(10) < 2 {
		route.Toll = 10 + g.gen.Intn(40)
	}

	return route
}

func (g *Generator) generateMarketName() string {
	names := map[engine.GenreID][]string{
		engine.GenreFantasy: {
			"Riverford Market", "Crossroads Trading Post", "Highcastle Bazaar",
			"Silverhollow Fair", "Dragon's Gate Emporium", "Moonwell Market",
			"Thornwood Exchange", "Crown City Plaza", "Frostpeak Trading Hall",
		},
		engine.GenreScifi: {
			"Station Alpha Market", "Orbital Trade Hub", "Frontier Exchange",
			"Deep Space Depot", "Colony Commerce Center", "Jump Point Bazaar",
			"Mining Station Market", "The Nexus Exchange", "Rim Trade Post",
		},
		engine.GenreHorror: {
			"Survivor's Market", "Safe Zone Trading", "Underground Exchange",
			"Fortress Commerce", "Last Stand Bazaar", "Refugee Market",
			"Bunker Trading Post", "The Haven Market", "Sanctuary Exchange",
		},
		engine.GenreCyberpunk: {
			"Street Market", "Black Market Node", "Corporate Exchange",
			"The Sprawl Bazaar", "Neon District Market", "Fixer's Den",
			"Underground Trade Hub", "Net Exchange", "Chrome Alley Market",
		},
		engine.GenrePostapoc: {
			"Scrap Market", "Wasteland Trading Post", "Survivor's Bazaar",
			"The Settlement Exchange", "Bunker Market", "Caravan Stop",
			"Oasis Trading Hub", "Rust Town Market", "New Hope Exchange",
		},
	}

	nameList := names[g.genre]
	return seed.Choice(g.gen, nameList)
}

func (g *Generator) generateRouteName() string {
	names := map[engine.GenreID][]string{
		engine.GenreFantasy: {
			"The King's Road", "Merchant's Way", "Pilgrim's Path",
			"The Old Trade Route", "Dragon's Pass", "River Road",
			"Forest Trail", "Mountain Pass", "Coastal Way",
		},
		engine.GenreScifi: {
			"Trade Corridor Alpha", "Shipping Lane 7", "Transit Route",
			"Jump Route", "The Freight Line", "Supply Corridor",
			"Express Lane", "Priority Route", "Standard Passage",
		},
		engine.GenreHorror: {
			"The Safe Path", "Patrol Route", "Evacuation Road",
			"Underground Passage", "The Cleared Way", "Refugee Trail",
			"Night Run", "Emergency Route", "Escape Corridor",
		},
		engine.GenreCyberpunk: {
			"The Smuggler's Run", "Data Highway", "Street Route",
			"Underground Link", "Corporate Corridor", "The Shadow Path",
			"Net Cable Run", "Back Alley Route", "Transit Tunnel",
		},
		engine.GenrePostapoc: {
			"The Wasteland Road", "Caravan Trail", "Safe Passage",
			"The Old Highway", "Scavenger's Path", "Radiation-Free Route",
			"Water Route", "Supply Trail", "Survivor's Highway",
		},
	}

	nameList := names[g.genre]
	return seed.Choice(g.gen, nameList)
}

func (g *Generator) generateTradeGoods() []*TradeGood {
	goodData := map[engine.GenreID][]struct {
		id       string
		name     string
		category GoodCategory
		price    int
	}{
		engine.GenreFantasy: {
			{"food", "Grain", CategoryFood, 10},
			{"material", "Iron Ore", CategoryMaterial, 25},
			{"luxury", "Silk", CategoryLuxury, 100},
			{"medical", "Healing Herbs", CategoryMedical, 40},
			{"weapon", "Steel Blades", CategoryWeapon, 75},
			{"special", "Magic Crystals", CategorySpecial, 200},
		},
		engine.GenreScifi: {
			{"food", "Ration Packs", CategoryFood, 15},
			{"material", "Refined Ore", CategoryMaterial, 30},
			{"luxury", "Rare Minerals", CategoryLuxury, 120},
			{"medical", "Med-Kits", CategoryMedical, 50},
			{"weapon", "Plasma Cells", CategoryWeapon, 80},
			{"special", "Data Cores", CategorySpecial, 250},
		},
		engine.GenreHorror: {
			{"food", "Canned Food", CategoryFood, 20},
			{"material", "Scrap Metal", CategoryMaterial, 15},
			{"luxury", "Medicine", CategoryLuxury, 80},
			{"medical", "First Aid Kits", CategoryMedical, 60},
			{"weapon", "Ammunition", CategoryWeapon, 100},
			{"special", "Ritual Components", CategorySpecial, 150},
		},
		engine.GenreCyberpunk: {
			{"food", "Synth-Food", CategoryFood, 8},
			{"material", "Electronics", CategoryMaterial, 35},
			{"luxury", "Designer Stims", CategoryLuxury, 150},
			{"medical", "Trauma Kits", CategoryMedical, 70},
			{"weapon", "Smart Ammo", CategoryWeapon, 90},
			{"special", "Black ICE", CategorySpecial, 300},
		},
		engine.GenrePostapoc: {
			{"food", "Clean Water", CategoryFood, 25},
			{"material", "Scrap Parts", CategoryMaterial, 20},
			{"luxury", "Pre-War Goods", CategoryLuxury, 100},
			{"medical", "Rad-Away", CategoryMedical, 80},
			{"weapon", "Bullets", CategoryWeapon, 50},
			{"special", "Seeds", CategorySpecial, 200},
		},
	}

	goods := make([]*TradeGood, 0, 6)
	for _, data := range goodData[g.genre] {
		good := NewTradeGood(data.id, data.name, data.category, data.price, g.genre)
		// Add some variance to initial stock
		good.Stock = 30 + g.gen.Intn(40)
		goods = append(goods, good)
	}

	return goods
}

// GenerateSpeculationOpportunity finds the best buy/sell pair in the economy
func (g *Generator) GenerateSpeculationOpportunity(manager *EconomyManager) (buyMarket, sellMarket, goodID string, profit int) {
	bestProfit := 0
	bestBuy := ""
	bestSell := ""
	bestGood := ""

	for buyID, buyM := range manager.Markets {
		for sellID := range manager.Markets {
			if buyID == sellID {
				continue
			}

			for gID := range buyM.Goods {
				p, ok := manager.CalculateSpeculation(buyID, sellID, gID, 10)
				if ok && p > bestProfit {
					bestProfit = p
					bestBuy = buyID
					bestSell = sellID
					bestGood = gID
				}
			}
		}
	}

	return bestBuy, bestSell, bestGood, bestProfit
}
