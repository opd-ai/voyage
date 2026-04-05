// Package metaprog provides meta-progression between runs for Voyage.
//
// This package tracks persistent progress across multiple runs including:
//   - Unlock log of event types and destinations seen
//   - Unlockable starting crew archetypes
//   - Unlockable vessel configurations
//   - Hall of Records with best run summaries per genre
//
// # Unlock System
//
// Progress is tracked via a cumulative game-state hash that deterministically
// unlocks new content as players complete runs. Unlocks are not random but
// based on achievements and discoveries.
//
// # Hall of Records
//
// The Hall of Records preserves the best run summary per genre:
//   - Days traveled
//   - Crew survivors
//   - Final score
//   - Run date
//
// # Persistence
//
// Meta-progression data is saved to disk and persists between game sessions.
package metaprog
