package economy

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// GoodCategory represents the type of tradeable good
type GoodCategory int

const (
	CategoryFood GoodCategory = iota
	CategoryMaterial
	CategoryLuxury
	CategoryMedical
	CategoryWeapon
	CategorySpecial
)

// PriceHistoryLength is the number of ticks to track for sparkline
const PriceHistoryLength = 20

// TradeGood represents a tradeable commodity
type TradeGood struct {
	ID       string
	Name     string
	Category GoodCategory
	Genre    engine.GenreID

	// Base economics
	BasePrice int
	MinPrice  int
	MaxPrice  int

	// Current state
	CurrentPrice int
	Stock        int
	Demand       int // Demand pressure (positive = high demand)

	// History for sparkline
	PriceHistory []int
}

// NewTradeGood creates a new trade good
func NewTradeGood(id, name string, category GoodCategory, basePrice int, genre engine.GenreID) *TradeGood {
	return &TradeGood{
		ID:           id,
		Name:         name,
		Category:     category,
		Genre:        genre,
		BasePrice:    basePrice,
		MinPrice:     basePrice / 4,
		MaxPrice:     basePrice * 4,
		CurrentPrice: basePrice,
		Stock:        50,
		Demand:       0,
		PriceHistory: make([]int, 0, PriceHistoryLength),
	}
}

// SetGenre updates the good's genre
func (g *TradeGood) SetGenre(genre engine.GenreID) {
	g.Genre = genre
}

// RecalculatePrice updates price based on supply and demand
func (g *TradeGood) RecalculatePrice() {
	// Base calculation: price inversely proportional to stock
	supplyFactor := 1.0
	if g.Stock > 100 {
		supplyFactor = 0.5 + 50.0/float64(g.Stock)
	} else if g.Stock < 20 {
		supplyFactor = 2.0 - float64(g.Stock)/20.0
	}

	// Demand adjustment
	demandFactor := 1.0 + float64(g.Demand)*0.05

	newPrice := int(float64(g.BasePrice) * supplyFactor * demandFactor)

	// Clamp to bounds
	if newPrice < g.MinPrice {
		newPrice = g.MinPrice
	}
	if newPrice > g.MaxPrice {
		newPrice = g.MaxPrice
	}

	g.CurrentPrice = newPrice
}

// RecordPrice adds current price to history
func (g *TradeGood) RecordPrice() {
	g.PriceHistory = append(g.PriceHistory, g.CurrentPrice)
	if len(g.PriceHistory) > PriceHistoryLength {
		g.PriceHistory = g.PriceHistory[1:]
	}
}

// GetSparklineData returns price history as percentages of max for display
func (g *TradeGood) GetSparklineData() []float64 {
	data := make([]float64, len(g.PriceHistory))
	priceRange := float64(g.MaxPrice - g.MinPrice)
	if priceRange == 0 {
		priceRange = 1
	}
	for i, price := range g.PriceHistory {
		data[i] = float64(price-g.MinPrice) / priceRange
	}
	return data
}

// Market represents a regional marketplace
type Market struct {
	ID    string
	Name  string
	Genre engine.GenreID
	X, Y  int // Position in world

	// Goods available
	Goods map[string]*TradeGood

	// Market characteristics
	Prosperity float64      // 0.0-1.0 affects prices and stock
	Speciality GoodCategory // What this market produces cheaply
	Scarcity   GoodCategory // What this market needs

	// Connected markets
	ConnectedTo []string
}

// NewMarket creates a new market
func NewMarket(id, name string, x, y int, genre engine.GenreID) *Market {
	return &Market{
		ID:          id,
		Name:        name,
		Genre:       genre,
		X:           x,
		Y:           y,
		Goods:       make(map[string]*TradeGood),
		Prosperity:  0.5,
		ConnectedTo: make([]string, 0),
	}
}

// SetGenre updates the market and all goods' genre
func (m *Market) SetGenre(genre engine.GenreID) {
	m.Genre = genre
	for _, good := range m.Goods {
		good.SetGenre(genre)
	}
}

// AddGood registers a trade good in this market
func (m *Market) AddGood(good *TradeGood) {
	m.Goods[good.ID] = good
}

// GetGood retrieves a good by ID
func (m *Market) GetGood(id string) *TradeGood {
	return m.Goods[id]
}

// BuyGoods player buys from market (increases price)
func (m *Market) BuyGoods(goodID string, quantity int) (int, bool) {
	good := m.Goods[goodID]
	if good == nil || good.Stock < quantity {
		return 0, false
	}

	cost := good.CurrentPrice * quantity
	good.Stock -= quantity
	good.Demand += quantity / 5 // Buying increases demand
	good.RecalculatePrice()

	return cost, true
}

// SellGoods player sells to market (decreases price)
func (m *Market) SellGoods(goodID string, quantity int) (int, bool) {
	good := m.Goods[goodID]
	if good == nil {
		return 0, false
	}

	revenue := good.CurrentPrice * quantity
	good.Stock += quantity
	good.Demand -= quantity / 5 // Selling decreases demand
	good.RecalculatePrice()

	return revenue, true
}

// GetPriceHistory returns price history for a good
func (m *Market) GetPriceHistory(goodID string) []int {
	good := m.Goods[goodID]
	if good == nil {
		return nil
	}
	return good.PriceHistory
}

// Tick processes one time unit
func (m *Market) Tick() {
	for _, good := range m.Goods {
		// Natural demand decay toward 0
		if good.Demand > 0 {
			good.Demand--
		} else if good.Demand < 0 {
			good.Demand++
		}

		// Natural stock recovery toward baseline
		if good.Stock < 50 {
			good.Stock++
		}

		good.RecalculatePrice()
		good.RecordPrice()
	}
}

// ApplySpecialties adjusts prices based on market speciality/scarcity
func (m *Market) ApplySpecialties() {
	for _, good := range m.Goods {
		if good.Category == m.Speciality {
			// Produce this cheaply - more stock, lower base
			good.Stock += 20
			good.BasePrice = good.BasePrice * 80 / 100
		}
		if good.Category == m.Scarcity {
			// Need this badly - less stock, higher base
			good.Stock -= 10
			if good.Stock < 5 {
				good.Stock = 5
			}
			good.BasePrice = good.BasePrice * 120 / 100
		}
		good.RecalculatePrice()
	}
}

// Connect links this market to another
func (m *Market) Connect(otherID string) {
	for _, id := range m.ConnectedTo {
		if id == otherID {
			return // Already connected
		}
	}
	m.ConnectedTo = append(m.ConnectedTo, otherID)
}

// TradeRoute represents a connection between markets
type TradeRoute struct {
	ID       string
	FromID   string
	ToID     string
	Distance int
	Genre    engine.GenreID

	// Route characteristics
	Danger float64 // 0.0-1.0 affects trade volume
	Toll   int     // Fixed cost to use route
	Name   string
}

// NewTradeRoute creates a new route between markets
func NewTradeRoute(id, fromID, toID string, distance int, genre engine.GenreID) *TradeRoute {
	return &TradeRoute{
		ID:       id,
		FromID:   fromID,
		ToID:     toID,
		Distance: distance,
		Genre:    genre,
		Danger:   0.1,
		Toll:     0,
	}
}

// SetGenre updates the route's genre
func (r *TradeRoute) SetGenre(genre engine.GenreID) {
	r.Genre = genre
}

// EconomyManager manages all markets and routes
type EconomyManager struct {
	Markets map[string]*Market
	Routes  []*TradeRoute
	Genre   engine.GenreID

	// Configuration
	PropagationRate float64 // How fast prices spread (0.0-1.0)
	TickCount       int
}

// NewEconomyManager creates a new economy manager
func NewEconomyManager(genre engine.GenreID) *EconomyManager {
	return &EconomyManager{
		Markets:         make(map[string]*Market),
		Routes:          make([]*TradeRoute, 0),
		Genre:           genre,
		PropagationRate: 0.2,
		TickCount:       0,
	}
}

// SetGenre updates all markets and routes
func (e *EconomyManager) SetGenre(genre engine.GenreID) {
	e.Genre = genre
	for _, market := range e.Markets {
		market.SetGenre(genre)
	}
	for _, route := range e.Routes {
		route.SetGenre(genre)
	}
}

// AddMarket registers a market
func (e *EconomyManager) AddMarket(market *Market) {
	e.Markets[market.ID] = market
}

// GetMarket retrieves a market by ID
func (e *EconomyManager) GetMarket(id string) *Market {
	return e.Markets[id]
}

// AddRoute registers a trade route
func (e *EconomyManager) AddRoute(route *TradeRoute) {
	e.Routes = append(e.Routes, route)

	// Connect markets bidirectionally
	if from := e.Markets[route.FromID]; from != nil {
		from.Connect(route.ToID)
	}
	if to := e.Markets[route.ToID]; to != nil {
		to.Connect(route.FromID)
	}
}

// GetRoutesBetween finds routes connecting two markets
func (e *EconomyManager) GetRoutesBetween(fromID, toID string) []*TradeRoute {
	routes := make([]*TradeRoute, 0)
	for _, route := range e.Routes {
		if (route.FromID == fromID && route.ToID == toID) ||
			(route.FromID == toID && route.ToID == fromID) {
			routes = append(routes, route)
		}
	}
	return routes
}

// Tick processes one time unit for all markets and propagates prices
func (e *EconomyManager) Tick() {
	e.TickCount++

	// Update each market
	for _, market := range e.Markets {
		market.Tick()
	}

	// Propagate price changes along routes
	e.propagatePrices()
}

// propagatePrices spreads price changes along trade routes
func (e *EconomyManager) propagatePrices() {
	for _, route := range e.Routes {
		from := e.Markets[route.FromID]
		to := e.Markets[route.ToID]
		if from == nil || to == nil {
			continue
		}
		e.propagatePricesBetweenMarkets(from, to, route.Distance)
	}
}

// propagatePricesBetweenMarkets equalizes prices for shared goods between two markets.
func (e *EconomyManager) propagatePricesBetweenMarkets(from, to *Market, distance int) {
	for goodID, fromGood := range from.Goods {
		toGood := to.Goods[goodID]
		if toGood == nil {
			continue
		}
		e.adjustPricePair(fromGood, toGood, distance)
	}
}

// adjustPricePair equalizes prices between two instances of the same good.
func (e *EconomyManager) adjustPricePair(fromGood, toGood *TradeGood, distance int) {
	priceDiff := fromGood.CurrentPrice - toGood.CurrentPrice
	adjustment := int(float64(priceDiff) * e.PropagationRate / float64(distance+1))
	if adjustment == 0 {
		return
	}
	fromGood.CurrentPrice -= adjustment / 2
	toGood.CurrentPrice += adjustment / 2
	clampPrice(fromGood)
	clampPrice(toGood)
}

// clampPrice ensures a good's price stays within its min/max bounds.
func clampPrice(g *TradeGood) {
	if g.CurrentPrice < g.MinPrice {
		g.CurrentPrice = g.MinPrice
	}
	if g.CurrentPrice > g.MaxPrice {
		g.CurrentPrice = g.MaxPrice
	}
}

// GetMarketsByProximity returns markets sorted by distance from a point
func (e *EconomyManager) GetMarketsByProximity(x, y int) []*Market {
	markets := make([]*Market, 0, len(e.Markets))
	for _, m := range e.Markets {
		markets = append(markets, m)
	}

	// Simple bubble sort by distance (small N expected)
	for i := 0; i < len(markets)-1; i++ {
		for j := i + 1; j < len(markets); j++ {
			distI := (markets[i].X-x)*(markets[i].X-x) + (markets[i].Y-y)*(markets[i].Y-y)
			distJ := (markets[j].X-x)*(markets[j].X-x) + (markets[j].Y-y)*(markets[j].Y-y)
			if distJ < distI {
				markets[i], markets[j] = markets[j], markets[i]
			}
		}
	}

	return markets
}

// CalculateSpeculation calculates potential profit for buying at one market and selling at another
func (e *EconomyManager) CalculateSpeculation(buyMarketID, sellMarketID, goodID string, quantity int) (int, bool) {
	buyMarket := e.Markets[buyMarketID]
	sellMarket := e.Markets[sellMarketID]

	if buyMarket == nil || sellMarket == nil {
		return 0, false
	}

	buyGood := buyMarket.Goods[goodID]
	sellGood := sellMarket.Goods[goodID]

	if buyGood == nil || sellGood == nil {
		return 0, false
	}

	if buyGood.Stock < quantity {
		return 0, false
	}

	buyCost := buyGood.CurrentPrice * quantity
	sellRevenue := sellGood.CurrentPrice * quantity
	profit := sellRevenue - buyCost

	return profit, true
}

// GoodName returns the genre-appropriate name for a good category
func GoodName(id string, genre engine.GenreID) string {
	goodNames := map[engine.GenreID]map[string]string{
		engine.GenreFantasy: {
			"food":     "Grain",
			"material": "Iron Ore",
			"luxury":   "Silk",
			"medical":  "Healing Herbs",
			"weapon":   "Steel Blades",
			"special":  "Magic Crystals",
		},
		engine.GenreScifi: {
			"food":     "Ration Packs",
			"material": "Refined Ore",
			"luxury":   "Rare Minerals",
			"medical":  "Med-Kits",
			"weapon":   "Plasma Cells",
			"special":  "Data Cores",
		},
		engine.GenreHorror: {
			"food":     "Canned Food",
			"material": "Scrap Metal",
			"luxury":   "Medicine",
			"medical":  "First Aid Kits",
			"weapon":   "Ammunition",
			"special":  "Ritual Components",
		},
		engine.GenreCyberpunk: {
			"food":     "Synth-Food",
			"material": "Electronics",
			"luxury":   "Designer Stims",
			"medical":  "Trauma Kits",
			"weapon":   "Smart Ammo",
			"special":  "Black ICE",
		},
		engine.GenrePostapoc: {
			"food":     "Clean Water",
			"material": "Scrap Parts",
			"luxury":   "Pre-War Goods",
			"medical":  "Rad-Away",
			"weapon":   "Bullets",
			"special":  "Seeds",
		},
	}

	if genreGoods, ok := goodNames[genre]; ok {
		if name, ok := genreGoods[id]; ok {
			return name
		}
	}
	return goodNames[engine.GenreFantasy][id]
}

// CategoryName returns the human-readable category name
func CategoryName(cat GoodCategory) string {
	names := map[GoodCategory]string{
		CategoryFood:     "Food",
		CategoryMaterial: "Materials",
		CategoryLuxury:   "Luxury",
		CategoryMedical:  "Medical",
		CategoryWeapon:   "Weapons",
		CategorySpecial:  "Special",
	}
	if name, ok := names[cat]; ok {
		return name
	}
	return "Unknown"
}

// AllCategories returns all good categories
func AllCategories() []GoodCategory {
	return []GoodCategory{
		CategoryFood, CategoryMaterial, CategoryLuxury,
		CategoryMedical, CategoryWeapon, CategorySpecial,
	}
}
