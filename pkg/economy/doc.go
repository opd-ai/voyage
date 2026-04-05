// Package economy provides dynamic trade route economics for the voyage game.
//
// This package implements a regional supply and demand model with price propagation,
// speculation mechanics, and genre-appropriate trade goods vocabulary.
//
// # Overview
//
// The economy system simulates realistic market dynamics where player actions
// affect regional prices. Goods sold in a region become cheaper there while
// scarcities drive prices up. Price changes propagate along procedurally
// generated trade routes.
//
// # Core Components
//
//   - TradeGood: A tradeable commodity with base price and current stock
//   - Market: Regional marketplace with supply/demand mechanics
//   - TradeRoute: Connection between markets that propagates price changes
//   - EconomyManager: Manages all markets and routes in the world
//
// # Features
//
//   - Regional supply and demand model
//   - Price propagation along trade routes
//   - Price history tracking with sparkline data
//   - Speculation mechanic for buying cheap and selling dear
//   - Genre-appropriate trade goods (grain→fuel cells→supplies→data→scrap)
//
// # Genre Support
//
// All components implement engine.GenreSwitcher for genre-aware generation:
//
//   - Fantasy: Grain, Spices, Silk, Iron, Gold, Magic Crystals
//   - Sci-Fi: Fuel Cells, Ore, Electronics, Medical Supplies, Data Cores
//   - Horror: Medical Supplies, Ammunition, Preserved Food, Weapons, Safe Water
//   - Cyberpunk: Data Chips, Access Codes, Cyberware, Stims, Black ICE
//   - Post-Apocalyptic: Scrap Metal, Clean Water, Rad-X, Fuel, Seeds
//
// # Usage
//
//	g := economy.NewGenerator(seed, engine.GenreFantasy)
//	manager := g.GenerateEconomy(5) // 5 markets
//
//	// Player sells goods at market
//	market := manager.GetMarket("town_1")
//	market.SellGoods("grain", 10)
//
//	// Price changes propagate
//	manager.Tick()
//
//	// Check price history
//	history := market.GetPriceHistory("grain")
//
// # Price Mechanics
//
//   - Base prices set by good type and genre
//   - Supply affects price: more supply = lower price
//   - Demand affects price: more demand = higher price
//   - Player trades affect local supply/demand
//   - Changes propagate to connected markets over time
package economy
