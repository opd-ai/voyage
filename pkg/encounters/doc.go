// Package encounters provides tactical encounter resolution for Voyage.
//
// This package implements FTL-style pausable encounters with crew role assignment
// and branching outcomes. All encounters are procedurally generated and support
// genre-specific reskinning via the GenreSwitcher interface.
//
// # Encounter Types
//
//   - Ambush: Combat-focused encounter requiring fighters and quick decisions
//   - Negotiation: Dialogue-based encounter benefiting from negotiators
//   - Race: Timed challenges requiring navigation and engineering
//   - Crisis: Emergency situations requiring varied skills
//   - Puzzle: Logic challenges that benefit from scouts and engineers
//
// # Resolution Flow
//
// 1. Encounter triggers based on map position and seed
// 2. Player assigns crew members to encounter roles
// 3. Resolution phase runs (turn-based or pausable)
// 4. Outcome determined: victory, partial success, retreat, defeat
// 5. Resource deltas applied to game state
//
// # Genre Support
//
// All encounter types implement SetGenre() to swap vocabulary, visuals,
// and sound design to match the current genre theme.
package encounters
