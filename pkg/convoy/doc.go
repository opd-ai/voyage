// Package convoy provides async convoy mode for shared-seed multiplayer.
//
// Convoy mode allows multiple players to run the same procedurally generated
// world (same seed) simultaneously, with their progress tracked independently.
// At the end of each run, players can compare their results to see who
// survived best or reached the destination first.
//
// Key features:
//   - Shared-seed co-op: all convoy members play the same world
//   - Async event resolution: each player's run can diverge
//   - End-of-run comparison across convoy members
//   - High-latency tolerant design (200-5000ms latency supported)
//
// The convoy system is designed to work well over Tor and other
// high-latency networks.
package convoy
