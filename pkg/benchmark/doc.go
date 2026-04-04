// Package benchmark provides performance benchmarks for Voyage.
//
// This package contains benchmarks to verify the performance targets:
//   - 60 FPS target (16.67ms per frame)
//   - <500MB heap usage
//
// Run benchmarks with:
//
//	go test -bench=. ./pkg/benchmark/...
package benchmark
