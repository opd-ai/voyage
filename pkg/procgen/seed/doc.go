// Package seed provides deterministic random number generation for Voyage.
//
// The seed system ensures that every procedural generation subsystem produces
// identical results given the same master seed. This enables reproducible runs
// where the same seed produces the same game world, events, crew, and outcomes.
//
// Usage:
//
//	masterSeed := int64(12345)
//	worldRng := seed.HashSeed(masterSeed, "world")
//	eventRng := seed.HashSeed(masterSeed, "events")
//	crewRng := seed.HashSeed(masterSeed, "crew")
//
// Each subsystem gets its own isolated RNG stream that is deterministically
// derived from the master seed plus a subsystem identifier.
package seed
