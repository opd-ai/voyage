// Package world provides procedural world map generation for Voyage.
//
// This package implements:
//   - Voronoi/grid-based overworld with regions and biomes
//   - Origin and destination placement with guaranteed solvable paths
//   - Terrain type assignment based on genre
//   - Waypoint and landmark seeding
//   - Branching path networks with risk/reward tradeoffs
//
// All world generation is deterministic based on the master seed.
package world
