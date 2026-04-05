// Package leaderboard provides per-seed leaderboard functionality.
//
// This package implements leaderboard submission and retrieval for
// completed game runs. Leaderboards are organized by seed and genre,
// allowing players to compare their performance on the same
// procedurally generated world.
//
// The leaderboard system supports:
//   - Submitting run completions with seed, genre, score, days, and survivors
//   - Querying global leaderboards filtered by genre and/or seed
//   - Replaying seeds from top-score runs
//   - Local caching of leaderboard data when offline
//
// Note: The actual server communication requires network connectivity.
// When offline, submissions are queued locally and synced when
// connectivity is restored.
package leaderboard
