package economy

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewTradeGood(t *testing.T) {
	good := NewTradeGood("food", "Grain", CategoryFood, 100, engine.GenreFantasy)

	if good.ID != "food" {
		t.Error("ID mismatch")
	}
	if good.BasePrice != 100 {
		t.Error("base price mismatch")
	}
	if good.CurrentPrice != 100 {
		t.Error("current price should equal base")
	}
	if good.MinPrice != 25 {
		t.Error("min price should be 1/4 of base")
	}
	if good.MaxPrice != 400 {
		t.Error("max price should be 4x base")
	}
}

func TestTradeGoodRecalculatePrice(t *testing.T) {
	good := NewTradeGood("test", "Test", CategoryFood, 100, engine.GenreFantasy)

	// High stock should lower price
	good.Stock = 200
	good.RecalculatePrice()
	if good.CurrentPrice >= 100 {
		t.Error("high stock should lower price")
	}

	// Low stock should raise price
	good.Stock = 10
	good.Demand = 0
	good.RecalculatePrice()
	if good.CurrentPrice <= 100 {
		t.Error("low stock should raise price")
	}

	// High demand should raise price
	good.Stock = 50
	good.Demand = 10
	good.RecalculatePrice()
	if good.CurrentPrice <= 100 {
		t.Error("high demand should raise price")
	}
}

func TestTradeGoodPriceHistory(t *testing.T) {
	good := NewTradeGood("test", "Test", CategoryFood, 100, engine.GenreFantasy)

	for i := 0; i < 25; i++ {
		good.CurrentPrice = 100 + i
		good.RecordPrice()
	}

	if len(good.PriceHistory) != PriceHistoryLength {
		t.Errorf("history should be limited to %d entries, got %d", PriceHistoryLength, len(good.PriceHistory))
	}
}

func TestTradeGoodSparklineData(t *testing.T) {
	good := NewTradeGood("test", "Test", CategoryFood, 100, engine.GenreFantasy)
	good.MinPrice = 0
	good.MaxPrice = 100

	good.PriceHistory = []int{0, 50, 100}
	data := good.GetSparklineData()

	if len(data) != 3 {
		t.Error("sparkline data should have 3 entries")
	}
	if data[0] != 0.0 {
		t.Errorf("expected 0.0, got %f", data[0])
	}
	if data[1] != 0.5 {
		t.Errorf("expected 0.5, got %f", data[1])
	}
	if data[2] != 1.0 {
		t.Errorf("expected 1.0, got %f", data[2])
	}
}

func TestNewMarket(t *testing.T) {
	market := NewMarket("test_1", "Test Market", 10, 20, engine.GenreFantasy)

	if market.ID != "test_1" {
		t.Error("ID mismatch")
	}
	if market.X != 10 || market.Y != 20 {
		t.Error("position mismatch")
	}
	if len(market.Goods) != 0 {
		t.Error("should start with no goods")
	}
}

func TestMarketBuySell(t *testing.T) {
	market := NewMarket("test", "Test", 0, 0, engine.GenreFantasy)
	good := NewTradeGood("grain", "Grain", CategoryFood, 10, engine.GenreFantasy)
	good.Stock = 50
	market.AddGood(good)

	// Buy goods
	initialPrice := good.CurrentPrice
	cost, ok := market.BuyGoods("grain", 10)
	if !ok {
		t.Error("buy should succeed")
	}
	if cost != initialPrice*10 {
		t.Error("cost calculation wrong")
	}
	if good.Stock != 40 {
		t.Error("stock should decrease")
	}

	// Sell goods
	revenue, ok := market.SellGoods("grain", 5)
	if !ok {
		t.Error("sell should succeed")
	}
	if revenue == 0 {
		t.Error("should get revenue")
	}
	if good.Stock != 45 {
		t.Error("stock should increase")
	}
}

func TestMarketBuyInsufficientStock(t *testing.T) {
	market := NewMarket("test", "Test", 0, 0, engine.GenreFantasy)
	good := NewTradeGood("grain", "Grain", CategoryFood, 10, engine.GenreFantasy)
	good.Stock = 5
	market.AddGood(good)

	_, ok := market.BuyGoods("grain", 10)
	if ok {
		t.Error("buy should fail with insufficient stock")
	}
}

func TestMarketTick(t *testing.T) {
	market := NewMarket("test", "Test", 0, 0, engine.GenreFantasy)
	good := NewTradeGood("grain", "Grain", CategoryFood, 10, engine.GenreFantasy)
	good.Stock = 30
	good.Demand = 5
	market.AddGood(good)

	market.Tick()

	// Stock should recover toward 50
	if good.Stock != 31 {
		t.Errorf("stock should recover, got %d", good.Stock)
	}
	// Demand should decay toward 0
	if good.Demand != 4 {
		t.Errorf("demand should decay, got %d", good.Demand)
	}
	// Price history should have entry
	if len(good.PriceHistory) != 1 {
		t.Error("should record price")
	}
}

func TestMarketConnect(t *testing.T) {
	market := NewMarket("test", "Test", 0, 0, engine.GenreFantasy)

	market.Connect("other_1")
	market.Connect("other_2")
	market.Connect("other_1") // Duplicate

	if len(market.ConnectedTo) != 2 {
		t.Errorf("should have 2 connections, got %d", len(market.ConnectedTo))
	}
}

func TestMarketApplySpecialties(t *testing.T) {
	market := NewMarket("test", "Test", 0, 0, engine.GenreFantasy)
	good := NewTradeGood("grain", "Grain", CategoryFood, 100, engine.GenreFantasy)
	market.AddGood(good)

	market.Speciality = CategoryFood
	market.Scarcity = CategoryMaterial

	initialStock := good.Stock
	market.ApplySpecialties()

	if good.Stock <= initialStock {
		t.Error("specialty goods should have increased stock")
	}
}

func TestNewTradeRoute(t *testing.T) {
	route := NewTradeRoute("route_1", "market_1", "market_2", 100, engine.GenreFantasy)

	if route.ID != "route_1" {
		t.Error("ID mismatch")
	}
	if route.Distance != 100 {
		t.Error("distance mismatch")
	}
}

func TestEconomyManager(t *testing.T) {
	manager := NewEconomyManager(engine.GenreFantasy)

	m1 := NewMarket("m1", "Market 1", 0, 0, engine.GenreFantasy)
	m2 := NewMarket("m2", "Market 2", 100, 0, engine.GenreFantasy)

	manager.AddMarket(m1)
	manager.AddMarket(m2)

	if len(manager.Markets) != 2 {
		t.Error("should have 2 markets")
	}

	if manager.GetMarket("m1") != m1 {
		t.Error("should retrieve market by ID")
	}
}

func TestEconomyManagerRoutes(t *testing.T) {
	manager := NewEconomyManager(engine.GenreFantasy)

	m1 := NewMarket("m1", "Market 1", 0, 0, engine.GenreFantasy)
	m2 := NewMarket("m2", "Market 2", 100, 0, engine.GenreFantasy)
	manager.AddMarket(m1)
	manager.AddMarket(m2)

	route := NewTradeRoute("r1", "m1", "m2", 50, engine.GenreFantasy)
	manager.AddRoute(route)

	routes := manager.GetRoutesBetween("m1", "m2")
	if len(routes) != 1 {
		t.Error("should find 1 route")
	}

	// Should work in reverse direction too
	routes = manager.GetRoutesBetween("m2", "m1")
	if len(routes) != 1 {
		t.Error("should find route in reverse direction")
	}
}

func TestEconomyManagerTick(t *testing.T) {
	manager := NewEconomyManager(engine.GenreFantasy)

	m1 := NewMarket("m1", "Market 1", 0, 0, engine.GenreFantasy)
	good := NewTradeGood("grain", "Grain", CategoryFood, 100, engine.GenreFantasy)
	good.Stock = 30
	m1.AddGood(good)
	manager.AddMarket(m1)

	manager.Tick()

	if manager.TickCount != 1 {
		t.Error("tick count should increment")
	}
}

func TestEconomyManagerPricePropagation(t *testing.T) {
	manager := NewEconomyManager(engine.GenreFantasy)
	manager.PropagationRate = 0.5

	m1 := NewMarket("m1", "Market 1", 0, 0, engine.GenreFantasy)
	m2 := NewMarket("m2", "Market 2", 100, 0, engine.GenreFantasy)

	g1 := NewTradeGood("grain", "Grain", CategoryFood, 100, engine.GenreFantasy)
	g1.CurrentPrice = 50
	g2 := NewTradeGood("grain", "Grain", CategoryFood, 100, engine.GenreFantasy)
	g2.CurrentPrice = 150

	m1.AddGood(g1)
	m2.AddGood(g2)

	manager.AddMarket(m1)
	manager.AddMarket(m2)

	route := NewTradeRoute("r1", "m1", "m2", 1, engine.GenreFantasy)
	manager.AddRoute(route)

	// Prices should move toward each other
	manager.Tick()

	if g1.CurrentPrice >= 150 {
		t.Error("lower price should have increased")
	}
	if g2.CurrentPrice <= 50 {
		t.Error("higher price should have decreased")
	}
}

func TestEconomyManagerSetGenre(t *testing.T) {
	manager := NewEconomyManager(engine.GenreFantasy)
	m1 := NewMarket("m1", "Market 1", 0, 0, engine.GenreFantasy)
	good := NewTradeGood("grain", "Grain", CategoryFood, 100, engine.GenreFantasy)
	m1.AddGood(good)
	manager.AddMarket(m1)

	route := NewTradeRoute("r1", "m1", "m2", 50, engine.GenreFantasy)
	manager.AddRoute(route)

	manager.SetGenre(engine.GenreScifi)

	if manager.Genre != engine.GenreScifi {
		t.Error("manager genre should update")
	}
	if m1.Genre != engine.GenreScifi {
		t.Error("market genre should update")
	}
	if good.Genre != engine.GenreScifi {
		t.Error("good genre should update")
	}
	if route.Genre != engine.GenreScifi {
		t.Error("route genre should update")
	}
}

func TestEconomyManagerSpeculation(t *testing.T) {
	manager := NewEconomyManager(engine.GenreFantasy)

	m1 := NewMarket("m1", "Market 1", 0, 0, engine.GenreFantasy)
	m2 := NewMarket("m2", "Market 2", 100, 0, engine.GenreFantasy)

	g1 := NewTradeGood("grain", "Grain", CategoryFood, 50, engine.GenreFantasy)
	g1.Stock = 100
	g2 := NewTradeGood("grain", "Grain", CategoryFood, 100, engine.GenreFantasy)

	m1.AddGood(g1)
	m2.AddGood(g2)

	manager.AddMarket(m1)
	manager.AddMarket(m2)

	profit, ok := manager.CalculateSpeculation("m1", "m2", "grain", 10)
	if !ok {
		t.Error("speculation calculation should succeed")
	}
	if profit <= 0 {
		t.Error("should show profit buying cheap and selling dear")
	}
}

func TestEconomyManagerMarketsByProximity(t *testing.T) {
	manager := NewEconomyManager(engine.GenreFantasy)

	m1 := NewMarket("m1", "Far", 100, 100, engine.GenreFantasy)
	m2 := NewMarket("m2", "Near", 10, 10, engine.GenreFantasy)
	m3 := NewMarket("m3", "Medium", 50, 50, engine.GenreFantasy)

	manager.AddMarket(m1)
	manager.AddMarket(m2)
	manager.AddMarket(m3)

	markets := manager.GetMarketsByProximity(0, 0)

	if markets[0].ID != "m2" {
		t.Error("nearest market should be first")
	}
	if markets[2].ID != "m1" {
		t.Error("farthest market should be last")
	}
}

func TestGoodName(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		goodIDs := []string{"food", "material", "luxury", "medical", "weapon", "special"}
		for _, id := range goodIDs {
			name := GoodName(id, genre)
			if name == "" {
				t.Errorf("good %s genre %s should have name", id, genre)
			}
		}
	}
}

func TestCategoryNameFunc(t *testing.T) {
	for _, cat := range AllCategories() {
		name := CategoryName(cat)
		if name == "Unknown" || name == "" {
			t.Errorf("category %v should have name", cat)
		}
	}
}

func TestGenerator(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)

	economy := g.GenerateEconomy(5)

	if len(economy.Markets) != 5 {
		t.Errorf("expected 5 markets, got %d", len(economy.Markets))
	}
	if len(economy.Routes) < 4 {
		t.Error("should have at least 4 routes (linear connections)")
	}

	for _, market := range economy.Markets {
		if len(market.Goods) == 0 {
			t.Error("markets should have goods")
		}
	}
}

func TestGeneratorDeterminism(t *testing.T) {
	g1 := NewGenerator(12345, engine.GenreFantasy)
	g2 := NewGenerator(12345, engine.GenreFantasy)

	e1 := g1.GenerateEconomy(3)
	e2 := g2.GenerateEconomy(3)

	if len(e1.Markets) != len(e2.Markets) {
		t.Error("same seed should produce same number of markets")
	}
}

func TestGeneratorAllGenres(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		g := NewGenerator(12345, genre)
		economy := g.GenerateEconomy(3)

		for _, market := range economy.Markets {
			if market.Name == "" {
				t.Errorf("genre %s: market should have name", genre)
			}
			for _, good := range market.Goods {
				if good.Name == "" {
					t.Errorf("genre %s: good should have name", genre)
				}
			}
		}

		for _, route := range economy.Routes {
			if route.Name == "" {
				t.Errorf("genre %s: route should have name", genre)
			}
		}
	}
}

func TestGeneratorSetGenre(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)
	g.SetGenre(engine.GenreCyberpunk)

	economy := g.GenerateEconomy(2)

	if economy.Genre != engine.GenreCyberpunk {
		t.Error("economy should have cyberpunk genre")
	}
}

func TestGeneratorSpeculationOpportunity(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)
	economy := g.GenerateEconomy(3)

	// Artificially create price difference
	var firstMarket, secondMarket *Market
	for _, m := range economy.Markets {
		if firstMarket == nil {
			firstMarket = m
		} else if secondMarket == nil {
			secondMarket = m
			break
		}
	}

	if firstMarket != nil && secondMarket != nil {
		// Set up price differential
		if g1 := firstMarket.GetGood("food"); g1 != nil {
			g1.CurrentPrice = 10
			g1.Stock = 100
		}
		if g2 := secondMarket.GetGood("food"); g2 != nil {
			g2.CurrentPrice = 50
		}
	}

	buyMarket, sellMarket, goodID, profit := g.GenerateSpeculationOpportunity(economy)

	if buyMarket == "" || sellMarket == "" || goodID == "" {
		t.Error("should find speculation opportunity")
	}
	if profit <= 0 {
		t.Error("should find profitable opportunity")
	}
}
