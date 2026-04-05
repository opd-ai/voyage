// Package factions provides procedural faction generation and management for Voyage.
//
// This package generates 4-6 factions per run with seeded names, ideologies,
// and controlled territory on the overworld. Factions have relationship matrices
// (allied, neutral, hostile) affected by player choices.
//
// # Faction System
//
// Each run procedurally generates factions with:
//   - Unique names based on genre and ideology
//   - Territory blocks on the overworld
//   - Relationship matrix with other factions
//   - Player reputation tracking
//
// # Relationships
//
// Factions can be:
//   - Allied: Actively help each other
//   - Neutral: No strong feelings
//   - Hostile: Attack on sight
//
// Player actions (favors, betrayals, bribes) shift faction standing.
//
// # Genre Support
//
// All faction types implement SetGenre() to swap vocabulary:
//   - Fantasy: guild, duchy, cult, order
//   - Scifi: corporation, colony, pirate fleet
//   - Horror: gang, survivor band, military remnant
//   - Cyberpunk: megacorp, street gang, syndicate
//   - Postapoc: settlement, raider clan, merchant guild
package factions
