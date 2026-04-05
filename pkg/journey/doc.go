// Package journey provides multi-leg journey support for the voyage game.
//
// This package implements campaign mode with chained journey legs, intermediate
// stopover cities, and escalating difficulty progression.
//
// # Overview
//
// The journey system manages multi-leg campaigns where players travel through
// 2-4 connected journey legs with persistent state between legs. Each leg
// features increasing difficulty with longer distances, harsher terrain,
// and stronger faction opposition.
//
// # Core Components
//
//   - Campaign: Manages a multi-leg journey with persistent state
//   - Leg: Individual journey segment with origin, destination, and difficulty
//   - Stopover: Intermediate hub city between legs for resupply and recruitment
//   - CampaignState: Persistent state that carries between journey legs
//
// # Features
//
//   - Chain 2-4 journey legs with state persisting between legs
//   - Intermediate stopover cities with trading, upgrading, and recruiting
//   - Escalating difficulty per leg (longer distances, harsher terrain)
//   - Optional genre shifts between legs for narrative variety
//   - All stopover names, descriptions, and features procedurally generated
//
// # Genre Support
//
// All components implement engine.GenreSwitcher for genre-aware generation:
//
//   - Fantasy: Fortified towns, trade posts, sacred shrines
//   - Sci-Fi: Space stations, orbital platforms, fuel depots
//   - Horror: Refuge camps, barricaded settlements, safe houses
//   - Cyberpunk: Megahubs, freeports, data havens
//   - Post-Apocalyptic: Survivor camps, bunker outposts, trade markets
//
// # Usage
//
//	g := journey.NewGenerator(seed, engine.GenreFantasy)
//	campaign := g.GenerateCampaign(3) // Generate 3-leg campaign
//
//	for _, leg := range campaign.Legs {
//	    // Execute each leg
//	    campaign.CompleteLeg(leg.ID)
//	}
//
// # Stopover Features
//
// Stopovers provide essential services between legs:
//
//   - Trading: Buy/sell resources at regional prices
//   - Repairs: Fix vessel damage, restore equipment
//   - Recruitment: Hire new crew members
//   - Upgrades: Improve vessel and equipment
//   - Information: Gather intel about upcoming leg
package journey
