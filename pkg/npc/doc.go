// Package npc provides procedural NPC generation for Voyage.
//
// This package generates wandering NPCs for random encounters including
// traders, refugees, bandits, and lost travelers. All NPCs have
// faction affiliations with alignment tags, and all names, dialogue,
// and descriptions are procedurally generated.
//
// # NPC Types
//
//   - Trader: Offers goods for sale/trade
//   - Refugee: Needs help or offers information
//   - Bandit: Hostile encounter
//   - Traveler: Neutral wanderer, may trade or share info
//   - Scout: May offer route information
//   - Guard: Patrols faction territory
//
// # Faction Support
//
// NPCs can be affiliated with factions which affect their behavior
// and disposition toward the player.
//
// # Genre Support
//
// All NPC types implement SetGenre() to swap vocabulary and
// archetypes to match the current genre theme.
package npc
