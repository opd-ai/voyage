// Package resources provides the six-axis resource management system for Voyage.
//
// The core resources are:
//   - Food: Depletes daily; crew starves if empty
//   - Water: Depletes daily; faster in desert/hot biomes
//   - Fuel/Stamina: Depletes per movement; vessel stops if empty
//   - Medicine: Consumed on injury/disease events
//   - Morale: Falls on hardship, rises on rest/success
//   - Currency: Used at supply points for trading
//
// Resource names change based on genre via the GenreSwitcher interface.
package resources
