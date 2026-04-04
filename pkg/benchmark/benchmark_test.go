package benchmark

import (
	"runtime"
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// BenchmarkWorldUpdate measures the ECS world update performance.
func BenchmarkWorldUpdate(b *testing.B) {
	registry := engine.NewComponentRegistry()
	world := engine.NewWorld(registry)

	// Spawn entities
	for i := 0; i < 1000; i++ {
		world.SpawnImmediate()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.Update(1.0 / 60.0)
	}
}

// BenchmarkEntitySpawnDespawn measures entity lifecycle performance.
func BenchmarkEntitySpawnDespawn(b *testing.B) {
	registry := engine.NewComponentRegistry()
	world := engine.NewWorld(registry)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e := world.SpawnImmediate()
		world.DespawnImmediate(e.ID())
	}
}

// BenchmarkSeedGeneration measures seed derivation performance.
func BenchmarkSeedGeneration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		seed.HashSeed(int64(i), "world")
	}
}

// BenchmarkRandomGeneration measures random value generation.
func BenchmarkRandomGeneration(b *testing.B) {
	gen := seed.NewGenerator(12345, "bench")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gen.Int()
	}
}

// BenchmarkEntityQuery measures entity querying performance.
func BenchmarkEntityQuery(b *testing.B) {
	registry := engine.NewComponentRegistry()
	world := engine.NewWorld(registry)

	// Register a test component
	posID := engine.ComponentID("position")

	// Spawn entities, half with the component
	for i := 0; i < 1000; i++ {
		e := world.SpawnImmediate()
		if i%2 == 0 {
			e.Add(&testComponent{id: posID})
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = world.EntitiesWith(posID)
	}
}

type testComponent struct {
	id engine.ComponentID
}

func (c *testComponent) ID() engine.ComponentID { return c.id }

// TestMemoryUsage verifies memory usage stays under the 500MB target.
func TestMemoryUsage(t *testing.T) {
	registry := engine.NewComponentRegistry()
	world := engine.NewWorld(registry)

	// Create a moderately large world
	for i := 0; i < 10000; i++ {
		world.SpawnImmediate()
	}

	// Force GC and measure
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	heapMB := float64(m.HeapAlloc) / 1024 / 1024
	t.Logf("Heap usage with 10k entities: %.2f MB", heapMB)

	// Target is <500MB, we should be way under with just entities
	if heapMB > 100 {
		t.Errorf("Heap usage %.2f MB exceeds expected baseline", heapMB)
	}
}

// BenchmarkGenreSwitch measures genre switching performance.
func BenchmarkGenreSwitch(b *testing.B) {
	registry := engine.NewComponentRegistry()
	world := engine.NewWorld(registry)

	genres := engine.AllGenres()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		world.SetGenre(genres[i%len(genres)])
	}
}
