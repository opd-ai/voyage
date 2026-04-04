# Implementation Gaps — 2026-04-04

This document enumerates the gaps between the stated goals in README.md and ROADMAP.md and the current implementation state of the Voyage project.

---

## Total Implementation Gap

- **Stated Goal**: A fully playable, 100% procedural travel simulator with five genre themes, ECS architecture, procedural audio/visual/narrative generation, and cross-platform support.
- **Current State**: **Zero implementation.** The repository contains only documentation files (`README.md`, `ROADMAP.md`, `LICENSE`, `.gitignore`). No Go source code, no `go.mod`, no packages, no executable.
- **Impact**: The project cannot be built, run, played, tested, or evaluated. Users expecting a game find only a design document.
- **Closing the Gap**: Begin implementation following ROADMAP.md v1.0 milestones:
  1. Initialize module: `go mod init github.com/opd-ai/voyage`
  2. Add Ebitengine: `go get github.com/hajimehoshi/ebiten/v2`
  3. Create entry point: `cmd/voyage/main.go`
  4. Implement ECS framework: `pkg/engine/`
  5. Implement procedural world generation: `pkg/procgen/world/`
  6. Continue through remaining v1.0 tasks

---

## Go Module Definition

- **Stated Goal**: "Single binary: `go build ./cmd/voyage` succeeds" (ROADMAP.md:460)
- **Current State**: No `go.mod` file exists. The project is not a valid Go module.
- **Impact**: Cannot compile, cannot add dependencies, cannot import packages, cannot run tests.
- **Closing the Gap**: 
  ```bash
  cd /home/user/go/src/github.com/opd-ai/voyage
  go mod init github.com/opd-ai/voyage
  go get github.com/hajimehoshi/ebiten/v2
  go mod tidy
  ```

---

## ECS Framework

- **Stated Goal**: "Component / Entity / System interfaces (`SetGenre(genreID GenreID)` required on every **System**)" (ROADMAP.md:61)
- **Current State**: No `pkg/engine/` directory. The `GenreSwitcher` interface exists only as a code block in documentation.
- **Impact**: No foundation exists for the entire game architecture. All other systems depend on ECS.
- **Closing the Gap**: Create the following files in `pkg/engine/`:
  - `genre.go` — `GenreID` type and `GenreSwitcher` interface
  - `component.go` — Component interface and registry
  - `entity.go` — Entity struct with component attachment
  - `system.go` — System interface with `Update()` and `SetGenre()`
  - `world.go` — World struct managing entities and systems
  
  Verification: `go build ./pkg/engine/`

---

## Seed-Based Deterministic RNG

- **Stated Goal**: "Master seed → subsystem seed derivation (`HashSeed` via SHA-256), Per-subsystem isolated `math/rand` sources, Determinism test suite (same seed → same game)" (ROADMAP.md:66-69)
- **Current State**: No seed system, no `HashSeed` function, no RNG isolation, no determinism tests.
- **Impact**: Runs are not reproducible. The core promise of "every run generated from a single seed" cannot be fulfilled.
- **Closing the Gap**: Create `pkg/procgen/seed/seed.go`:
  ```go
  package seed

  import (
      "crypto/sha256"
      "encoding/binary"
      "math/rand"
  )

  func HashSeed(master int64, subsystem string) *rand.Rand {
      h := sha256.New()
      binary.Write(h, binary.LittleEndian, master)
      h.Write([]byte(subsystem))
      sum := h.Sum(nil)
      derived := int64(binary.LittleEndian.Uint64(sum[:8]))
      return rand.New(rand.NewSource(derived))
  }
  ```
  
  Create `pkg/procgen/seed/seed_test.go` with determinism verification tests.

---

## Procedural World Map Generation

- **Stated Goal**: "Voronoi / grid-based overworld with regions and biomes, Origin → destination placement with guaranteed solvable path" (ROADMAP.md:76-81)
- **Current State**: No `pkg/procgen/world/` directory exists.
- **Impact**: No game world can be generated. No terrain, no waypoints, no routes.
- **Closing the Gap**: Implement `pkg/procgen/world/`:
  - `generator.go` — Voronoi/grid world generation
  - `terrain.go` — Terrain type definitions
  - `biome.go` — Biome assignment logic
  - `pathfinding.go` — Solvable path guarantee
  
  Verification: `go test ./pkg/procgen/world/...`

---

## Overworld Rendering

- **Stated Goal**: "Ebiten tile renderer for the world map, Procedural tile sprite generation (cellular automata + palette)" (ROADMAP.md:83-89)
- **Current State**: No `pkg/rendering/` directory. No Ebitengine dependency.
- **Impact**: Even if a world existed, it could not be displayed.
- **Closing the Gap**: Create `pkg/rendering/`:
  - `renderer.go` — Ebitengine `Game` interface implementation
  - `tilemap.go` — Tile rendering with camera viewport
  - `sprites.go` — Procedural sprite generation using cellular automata
  
  Verification: `go build ./cmd/voyage && ./voyage`

---

## Resource Management System

- **Stated Goal**: Six-axis resource model: "Food, Water, Fuel/Stamina, Medicine, Morale, Currency" (ROADMAP.md:97-105)
- **Current State**: No `pkg/resources/` implementation.
- **Impact**: No survival mechanics. The core "resource attrition" design pillar is absent.
- **Closing the Gap**: Create `pkg/resources/`:
  - `resources.go` — Resource struct with six fields
  - `consumption.go` — Daily depletion logic
  - `thresholds.go` — Warning and critical thresholds
  
  Verification: `go test ./pkg/resources/...`

---

## Party / Crew System

- **Stated Goal**: "Party entity with 2–6 crew member slots, Procedurally generated crew member (name, portrait-sprite, trait, skill)" (ROADMAP.md:107-112)
- **Current State**: No `pkg/crew/` implementation.
- **Impact**: No crew members, no mortality, no survival drama.
- **Closing the Gap**: Create `pkg/crew/`:
  - `crew.go` — Crew member struct (name, health, trait, skill)
  - `generator.go` — Procedural crew generation
  - `party.go` — Party management (2-6 slots)
  
  Verification: `go test ./pkg/crew/...`

---

## Vessel / Transport System

- **Stated Goal**: "Vessel entity with hull integrity, speed, and cargo capacity stats" (ROADMAP.md:114-119)
- **Current State**: No `pkg/vessel/` implementation.
- **Impact**: No transport, no journey progression, no cargo system.
- **Closing the Gap**: Create `pkg/vessel/`:
  - `vessel.go` — Vessel struct with integrity, speed, capacity
  - `cargo.go` — Cargo inventory management
  - `breakdown.go` — Random breakdown event logic
  
  Verification: `go test ./pkg/vessel/...`

---

## Procedural Event System

- **Stated Goal**: "Event queue seeded from master seed + current map position, All event text procedurally generated at runtime from grammar templates" (ROADMAP.md:127-132)
- **Current State**: No `pkg/events/` or `pkg/procgen/event/` implementation.
- **Impact**: No events, no choices, no branching narrative.
- **Closing the Gap**: Create `pkg/events/`:
  - `queue.go` — Seeded event queue
  - `event.go` — Event struct with choices and outcomes
  - `generator.go` — Grammar-based text generation
  
  Create `pkg/procgen/event/grammar.go` for text templates.
  
  Verification: `go test ./pkg/events/...`

---

## Audio Synthesis

- **Stated Goal**: "Sine / square / sawtooth / triangle / noise waveforms, ADSR envelope system, SFX generation, Ambient travel music" (ROADMAP.md:134-139)
- **Current State**: No `pkg/audio/` implementation.
- **Impact**: Silent game. No audio feedback for events, travel, or atmosphere.
- **Closing the Gap**: Create `pkg/audio/`:
  - `waveforms.go` — Oscillator implementations
  - `envelope.go` — ADSR envelope
  - `sfx.go` — SFX generation functions
  - `music.go` — Procedural ambient music
  
  Verification: `go test ./pkg/audio/...`

---

## UI / HUD / Menus

- **Stated Goal**: "World map screen, Resource panel, Crew roster panel, Event overlay, Main menu, pause menu, options screen" (ROADMAP.md:141-146)
- **Current State**: No `pkg/ux/` implementation.
- **Impact**: No user interface. Users cannot interact with any system.
- **Closing the Gap**: Create `pkg/ux/`:
  - `hud.go` — Resource and crew panels
  - `menus.go` — Main menu, pause, options
  - `events.go` — Event overlay with choices
  - `worldmap.go` — Map view and navigation
  
  Verification: `go build ./cmd/voyage && ./voyage`

---

## Save / Load System

- **Stated Goal**: "Multiple save slots with autosave on turn advance, Seed embedded in save for reproducibility" (ROADMAP.md:148-150)
- **Current State**: No `pkg/saveload/` implementation.
- **Impact**: No persistence. Game state lost on exit.
- **Closing the Gap**: Create `pkg/saveload/`:
  - `save.go` — Serialization logic
  - `load.go` — Deserialization logic
  - `slots.go` — Multi-slot management with autosave
  
  Verification: `go test ./pkg/saveload/...`

---

## Win / Lose Conditions

- **Stated Goal**: "Win: vessel reaches destination tile with ≥1 living crew member. Lose: vessel destroyed, or all crew dead, or morale hits zero" (ROADMAP.md:158-161)
- **Current State**: No game loop, no win/lose detection.
- **Impact**: No game can be completed or failed.
- **Closing the Gap**: Implement win/lose detection in `pkg/game/`:
  - `conditions.go` — Win/lose condition checks
  - `endscreen.go` — Run summary display
  
  Verification: Integration test that simulates a complete journey.

---

## Input System

- **Stated Goal**: "Keyboard / gamepad mapping, Rebindable controls, Modal input handling" (ROADMAP.md:71-73)
- **Current State**: No `pkg/config/` with input handling.
- **Impact**: Cannot control the game.
- **Closing the Gap**: Create `pkg/input/`:
  - `keyboard.go` — Key mappings
  - `gamepad.go` — Controller support
  - `rebinding.go` — User-configurable bindings
  
  Verification: `go build ./cmd/voyage && ./voyage` with keyboard/gamepad test.

---

## Validation Infrastructure

- **Stated Goal**: "Confirmed by `scripts/validate-no-assets.sh`" (ROADMAP.md:447-448)
- **Current State**: No `scripts/` directory.
- **Impact**: Cannot verify the "no bundled assets" constraint automatically.
- **Closing the Gap**: Create `scripts/validate-no-assets.sh`:
  ```bash
  #!/bin/bash
  set -e
  PROHIBITED_EXTENSIONS="png jpg jpeg gif svg bmp mp3 wav ogg flac aac"
  for ext in $PROHIBITED_EXTENSIONS; do
      if find . -name "*.$ext" -type f | grep -q .; then
          echo "FAIL: Found prohibited .$ext files"
          exit 1
      fi
  done
  echo "PASS: No bundled assets found"
  ```
  
  Verification: `chmod +x scripts/validate-no-assets.sh && ./scripts/validate-no-assets.sh`

---

## CI/CD Pipeline

- **Stated Goal**: "CI build job" (ROADMAP.md:460)
- **Current State**: No `.github/workflows/` directory.
- **Impact**: No automated build verification. PRs are not automatically tested.
- **Closing the Gap**: Create `.github/workflows/ci.yml`:
  ```yaml
  name: CI
  on: [push, pull_request]
  jobs:
    build:
      runs-on: ubuntu-latest
      steps:
        - uses: actions/checkout@v4
        - uses: actions/setup-go@v5
          with:
            go-version: '1.22'
        - run: go build ./cmd/voyage
        - run: go test ./...
        - run: go vet ./...
        - run: ./scripts/validate-no-assets.sh
  ```
  
  Verification: Push to GitHub and check Actions tab.

---

## Test Coverage

- **Stated Goal**: "Test coverage: ≥40% per package" (ROADMAP.md:463)
- **Current State**: 0% coverage. No tests exist because no code exists.
- **Impact**: Cannot verify correctness of any system.
- **Closing the Gap**: Write `*_test.go` files for every package as implementation proceeds. Use table-driven tests and property-based testing for procedural generators.
  
  Verification: `go test -cover ./pkg/... | grep -v "no test files"`

---

## Genre System Implementation

- **Stated Goal**: "Every system must implement the `GenreSwitcher` interface... This interface is not yet implemented — creating it is the first task in v1.0 ECS Framework" (ROADMAP.md:13)
- **Current State**: Interface exists only as markdown code block.
- **Impact**: No genre switching possible. All five theme variants are inaccessible.
- **Closing the Gap**: After creating ECS framework, ensure every System struct includes:
  ```go
  func (s *SomeSystem) SetGenre(genreID GenreID) {
      s.genre = genreID
      // Update genre-specific parameters
  }
  ```
  
  Create interface conformance test in `pkg/engine/genre_test.go`.

---

## Summary Table

| Gap Category | Files Needed | Priority |
|--------------|--------------|----------|
| Go Module | `go.mod` | P0 (Blocker) |
| ECS Framework | `pkg/engine/*.go` | P0 (Foundation) |
| Seed System | `pkg/procgen/seed/*.go` | P0 (Core) |
| World Generation | `pkg/procgen/world/*.go` | P1 |
| Rendering | `pkg/rendering/*.go` | P1 |
| Resources | `pkg/resources/*.go` | P1 |
| Crew | `pkg/crew/*.go` | P1 |
| Vessel | `pkg/vessel/*.go` | P1 |
| Events | `pkg/events/*.go` | P1 |
| Audio | `pkg/audio/*.go` | P2 |
| UI/HUD | `pkg/ux/*.go` | P2 |
| Save/Load | `pkg/saveload/*.go` | P2 |
| Input | `pkg/input/*.go` | P2 |
| Config/CLI | `pkg/config/*.go` | P2 |
| Entry Point | `cmd/voyage/main.go` | P1 |
| Scripts | `scripts/*.sh` | P3 |
| CI/CD | `.github/workflows/ci.yml` | P3 |

---

## Recommended Implementation Order

1. **Week 1**: `go.mod`, `pkg/engine/` (ECS + GenreSwitcher), `pkg/procgen/seed/`
2. **Week 2**: `cmd/voyage/main.go`, `pkg/rendering/` (basic Ebitengine loop)
3. **Week 3**: `pkg/procgen/world/`, `pkg/world/` (map generation + state)
4. **Week 4**: `pkg/resources/`, `pkg/crew/`, `pkg/vessel/`
5. **Week 5**: `pkg/events/`, `pkg/procgen/event/`
6. **Week 6**: `pkg/ux/` (HUD, menus, event overlay)
7. **Week 7**: `pkg/audio/` (synthesis, SFX)
8. **Week 8**: `pkg/saveload/`, `pkg/config/`, CLI flags
9. **Week 9**: Integration, playtesting, determinism verification
10. **Week 10**: CI/CD, documentation, scripts, release

This timeline assumes one developer working part-time. Adjust based on team size and availability.
