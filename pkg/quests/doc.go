// Package quests provides procedural quest and objective generation for Voyage.
//
// This package generates side quests for the travel simulator including
// delivery missions, rescue operations, and artifact retrieval. All quest
// text, objectives, and rewards are procedurally generated from seed.
//
// # Quest Types
//
//   - Delivery: Transport goods from one location to another
//   - Rescue: Save stranded crew or civilians
//   - Retrieve: Find and return an artifact or item
//   - Explore: Map a region or investigate a location
//   - Eliminate: Deal with a threat
//
// # Quest Board
//
// Supply points offer a quest board where players can accept or decline
// optional missions. Each quest has requirements, rewards, and time limits.
//
// # Genre Support
//
// All quest types implement SetGenre() to swap vocabulary and theming
// to match the current genre (fantasy scroll → data pad → wanted poster).
package quests
