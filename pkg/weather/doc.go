// Package weather provides weather and environmental hazard systems for Voyage.
//
// This package implements a procedural weather system with 8+ weather types and
// terrain hazards. Weather affects movement cost, resource consumption, visibility,
// and crew health. All weather types support genre-specific reskinning.
//
// # Weather Types
//
//   - Storm: High wind and rain/debris
//   - Blizzard: Extreme cold and low visibility
//   - Heatwave: High temperatures and water consumption
//   - Flood: Water damage and movement impediment
//   - Fog: Low visibility
//   - MeteorShower: Space debris (scifi genre)
//   - DustStorm: Choking particulates
//   - AcidRain: Corrosive precipitation
//
// # Terrain Hazards
//
//   - Mountain passes: Injury risk
//   - River crossings: Fuel cost
//   - Desert: Water crisis
//   - Ruin: Random loot and danger
//
// # Genre Support
//
// All weather types implement SetGenre() to swap vocabulary and visual effects
// to match the current genre theme.
package weather
