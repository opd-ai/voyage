// Package lore provides environmental storytelling and codex systems for Voyage.
//
// This package handles procedural generation of:
//   - World-map lore inscriptions (ruins, grave markers, burnt signs)
//   - Abandoned vessel/camp discoveries with item inventories
//   - Lore codex with world history, faction bios, and route legends
//
// # Environmental Storytelling
//
// Discoveries on the map include procedurally generated vignette text
// that reveals snippets of the world's history and current state.
// All text is generated algorithmically from seed.
//
// # Lore Codex
//
// The codex collects lore entries unlocked through:
//   - Exploring ruins and landmarks
//   - Completing events
//   - NPC conversations
//
// # Genre Support
//
// All lore types implement SetGenre() to adjust vocabulary
// and tone to match the current genre theme.
package lore
