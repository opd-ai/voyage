package seed

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
)

// HashSeed derives a deterministic RNG from a master seed and subsystem identifier.
// The same master seed and subsystem always produce the same RNG stream.
func HashSeed(master int64, subsystem string) *rand.Rand {
	h := sha256.New()

	// Write master seed as bytes
	seedBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(seedBytes, uint64(master))
	h.Write(seedBytes)

	// Write subsystem identifier
	h.Write([]byte(subsystem))

	// Extract derived seed from hash
	sum := h.Sum(nil)
	derived := int64(binary.LittleEndian.Uint64(sum[:8]))

	return rand.New(rand.NewSource(derived))
}

// Generator provides a seeded RNG for a specific subsystem.
type Generator struct {
	rng       *rand.Rand
	master    int64
	subsystem string
}

// NewGenerator creates a new generator for the given subsystem.
func NewGenerator(master int64, subsystem string) *Generator {
	return &Generator{
		rng:       HashSeed(master, subsystem),
		master:    master,
		subsystem: subsystem,
	}
}

// Rng returns the underlying RNG for direct access.
func (g *Generator) Rng() *rand.Rand {
	return g.rng
}

// Master returns the master seed used to create this generator.
func (g *Generator) Master() int64 {
	return g.master
}

// Subsystem returns the subsystem identifier.
func (g *Generator) Subsystem() string {
	return g.subsystem
}

// Reset recreates the RNG with the same seeds, returning to the initial state.
func (g *Generator) Reset() {
	g.rng = HashSeed(g.master, g.subsystem)
}

// Int returns a non-negative pseudo-random int.
func (g *Generator) Int() int {
	return g.rng.Int()
}

// Intn returns a non-negative pseudo-random int in [0, n).
func (g *Generator) Intn(n int) int {
	return g.rng.Intn(n)
}

// Int63 returns a non-negative pseudo-random 63-bit integer as int64.
func (g *Generator) Int63() int64 {
	return g.rng.Int63()
}

// Int63n returns a non-negative pseudo-random int64 in [0, n).
func (g *Generator) Int63n(n int64) int64 {
	return g.rng.Int63n(n)
}

// Float64 returns a pseudo-random float64 in [0.0, 1.0).
func (g *Generator) Float64() float64 {
	return g.rng.Float64()
}

// Float32 returns a pseudo-random float32 in [0.0, 1.0).
func (g *Generator) Float32() float32 {
	return g.rng.Float32()
}

// Shuffle pseudo-randomizes the order of elements using the Fisher-Yates shuffle.
func (g *Generator) Shuffle(n int, swap func(i, j int)) {
	g.rng.Shuffle(n, swap)
}

// Perm returns a pseudo-random permutation of integers [0, n).
func (g *Generator) Perm(n int) []int {
	return g.rng.Perm(n)
}

// Choice returns a random element from the slice.
func Choice[T any](g *Generator, slice []T) T {
	return slice[g.Intn(len(slice))]
}

// WeightedChoice returns a random element based on weights.
// The weights slice must have the same length as the choices slice.
func WeightedChoice[T any](g *Generator, choices []T, weights []float64) T {
	total := 0.0
	for _, w := range weights {
		total += w
	}

	r := g.Float64() * total
	cumulative := 0.0
	for i, w := range weights {
		cumulative += w
		if r < cumulative {
			return choices[i]
		}
	}
	return choices[len(choices)-1]
}

// Range returns a random integer in [min, max].
func (g *Generator) Range(min, max int) int {
	if min >= max {
		return min
	}
	return min + g.Intn(max-min+1)
}

// RangeFloat64 returns a random float64 in [min, max).
func (g *Generator) RangeFloat64(min, max float64) float64 {
	return min + g.Float64()*(max-min)
}

// Chance returns true with the given probability (0.0 to 1.0).
func (g *Generator) Chance(probability float64) bool {
	return g.Float64() < probability
}
