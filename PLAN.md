# Implementation Plan: Voyage v1.0 — Core Engine + Playable Journey

## Project Context
- **What it does**: A 100% procedural travel simulator inspired by Oregon Trail, FTL, and Organ Trail — every map, event, crew, vessel, audio, and narrative generated from a single seed.
- **Current goal**: Bootstrap the project from documentation-only to a playable single-genre journey (fantasy baseline).
- **Estimated Scope**: **Large** — Zero implementation exists; the entire codebase must be created from scratch.

## Goal-Achievement Status

| Stated Goal | Current Status | This Plan Addresses |
|-------------|---------------|---------------------|
| Go Module Definition | ❌ No `go.mod` exists | Yes |
| Single Binary (`go build ./cmd/voyage`) | ❌ No entry point | Yes |
| ECS Framework with GenreSwitcher | ❌ Interface only in docs | Yes |
| Seed-Based Deterministic RNG | ❌ Missing | Yes |
| Procedural World Map Generation | ❌ Missing | Yes |
| Overworld Rendering (Ebitengine) | ❌ Missing | Yes |
| Six-Axis Resource Model | ❌ Missing | Yes |
| Party/Crew System | ❌ Missing | Yes |
| Vessel/Transport System | ❌ Missing | Yes |
| Procedural Event System | ❌ Missing | Yes |
| Audio Synthesis | ❌ Missing | Partial (foundation) |
| UI/HUD/Menus | ❌ Missing | Yes |
| Save/Load System | ❌ Missing | Yes |
| Win/Lose Conditions | ❌ Missing | Yes |
| CI/CD Pipeline | ❌ Missing | Yes |
| Test Coverage ≥40% | 0% (no code) | Yes |
| Five Genre Support | ❌ Only fantasy in v1.0 | Foundation only |
| No Bundled Assets | ⚠️ Trivially true | Validation script |

## Metrics Summary

**Note**: `go-stats-generator` returns no metrics because no Go code exists.

```
Error: analysis failed: no Go files found in /home/user/go/src/github.com/opd-ai/voyage
```

- **Complexity hotspots**: N/A (0 functions)
- **Duplication ratio**: N/A (0 code)
- **Doc coverage**: N/A (0 packages)
- **Package coupling**: N/A (0 packages)
- **Lines of Code**: 0
- **Test Coverage**: 0%

**Baseline for Post-Implementation Validation**:
- Target complexity: All functions < 10 cyclomatic complexity
- Target duplication: < 5%
- Target doc coverage: ≥ 70% for exported symbols
- Target test coverage: ≥ 40% per package

---

## Implementation Steps

### Step 1: Initialize Go Module and Project Structure

- **Deliverable**: 
  - `go.mod` with module path `github.com/opd-ai/voyage`
  - `go.sum` with Ebitengine and core dependencies
  - Directory skeleton per ROADMAP.md architecture
- **Dependencies**: None (first step)
- **Goal Impact**: Enables "Single Binary" goal; unblocks all subsequent steps
- **Acceptance**: 
  - `go mod tidy` succeeds
  - `go list ./...` returns expected package paths
- **Validation**: 
  ```bash
  cd /home/user/go/src/github.com/opd-ai/voyage
  go mod tidy && go list ./... | head -20
  ```

**Files to create**:
```
go.mod
cmd/voyage/main.go          (stub)
pkg/engine/doc.go
pkg/procgen/seed/doc.go
pkg/procgen/world/doc.go
pkg/rendering/doc.go
pkg/resources/doc.go
pkg/crew/doc.go
pkg/vessel/doc.go
pkg/events/doc.go
pkg/audio/doc.go
pkg/ux/doc.go
pkg/saveload/doc.go
pkg/game/doc.go
pkg/config/doc.go
scripts/validate-no-assets.sh
```

---

### Step 2: Implement ECS Framework with GenreSwitcher

- **Deliverable**: 
  - `pkg/engine/genre.go` — `GenreID` type, constants, `GenreSwitcher` interface
  - `pkg/engine/component.go` — Component interface and registry
  - `pkg/engine/entity.go` — Entity struct with component attachment
  - `pkg/engine/system.go` — System interface with `Update()` and `SetGenre()`
  - `pkg/engine/world.go` — World struct managing entities and systems
  - `pkg/engine/*_test.go` — Unit tests with ≥40% coverage
- **Dependencies**: Step 1 (module initialization)
- **Goal Impact**: Foundation for all game systems; required for GenreSwitcher compliance
- **Acceptance**: 
  - `go build ./pkg/engine/` succeeds
  - `go test ./pkg/engine/...` passes with ≥40% coverage
  - All System implementations accept `SetGenre(GenreID)`
- **Validation**: 
  ```bash
  go test -cover ./pkg/engine/... | grep -E 'coverage|ok|FAIL'
  go-stats-generator analyze ./pkg/engine --format json --sections functions | jq '.functions | map(select(.complexity > 9)) | length'
  ```

---

### Step 3: Implement Seed-Based Deterministic RNG

- **Deliverable**:
  - `pkg/procgen/seed/seed.go` — `HashSeed(master int64, subsystem string) *rand.Rand`
  - `pkg/procgen/seed/seed_test.go` — Determinism verification tests
- **Dependencies**: Step 1 (module initialization)
- **Goal Impact**: Core promise "same seed → same game"; enables all procedural generators
- **Acceptance**:
  - `HashSeed(12345, "world")` returns identical RNG stream across runs
  - Determinism test generates 1000 values and verifies reproducibility
  - `go test ./pkg/procgen/seed/...` passes
- **Validation**: 
  ```bash
  go test -v ./pkg/procgen/seed/... -run TestDeterminism
  ```

---

### Step 4: Implement Procedural World Map Generation

- **Deliverable**:
  - `pkg/procgen/world/generator.go` — Voronoi/grid world generator
  - `pkg/procgen/world/terrain.go` — Terrain type definitions
  - `pkg/procgen/world/biome.go` — Biome assignment logic
  - `pkg/procgen/world/pathfinding.go` — Guaranteed solvable path algorithm
  - `pkg/procgen/world/genre.go` — `SetGenre()` to swap biome vocabulary
  - `pkg/world/state.go` — Runtime world state management
  - `pkg/world/fog.go` — Fog-of-war tracking
  - Tests with ≥40% coverage
- **Dependencies**: Step 2 (ECS for System integration), Step 3 (seed for determinism)
- **Goal Impact**: Enables "Fully Procedural World" pillar
- **Acceptance**:
  - Generator produces valid map with origin and destination
  - Pathfinding guarantees at least one solvable route
  - Same seed produces identical map layout
  - `SetGenre()` changes biome vocabulary
- **Validation**: 
  ```bash
  go test -cover ./pkg/procgen/world/... ./pkg/world/...
  go-stats-generator analyze ./pkg/procgen/world --format json --sections functions | jq '.functions | map(select(.complexity > 9))'
  ```

---

### Step 5: Implement Ebitengine Game Loop and Tile Renderer

- **Deliverable**:
  - `cmd/voyage/main.go` — Ebitengine game entry point
  - `pkg/game/game.go` — Game struct implementing `ebiten.Game` interface
  - `pkg/game/state.go` — State machine (menu, playing, paused, gameover)
  - `pkg/rendering/renderer.go` — Tile renderer with camera viewport
  - `pkg/rendering/tilemap.go` — Procedural tile sprite generation
  - `pkg/rendering/sprites.go` — Cellular automata + palette sprite generator
  - `pkg/rendering/genre.go` — `SetGenre()` for palette/theme switching
- **Dependencies**: Step 2 (ECS), Step 4 (world for something to render)
- **Goal Impact**: Enables visual feedback; critical path to "Playable Game"
- **Acceptance**:
  - `go build ./cmd/voyage` produces executable
  - Running `./voyage` displays generated world map
  - Fog-of-war masks unexplored regions
- **Validation**: 
  ```bash
  go build ./cmd/voyage && ls -la voyage
  go-stats-generator analyze ./pkg/rendering --format json --sections documentation | jq '.documentation.coverage'
  ```

---

### Step 6: Implement Six-Axis Resource Management

- **Deliverable**:
  - `pkg/resources/resources.go` — Resource struct (Food, Water, Fuel, Medicine, Morale, Currency)
  - `pkg/resources/consumption.go` — Daily depletion logic per terrain
  - `pkg/resources/thresholds.go` — Warning and critical threshold definitions
  - `pkg/resources/genre.go` — `SetGenre()` to rename resources per theme
  - Tests with ≥40% coverage
- **Dependencies**: Step 2 (ECS for System integration)
- **Goal Impact**: Enables "Resource Attrition" design pillar
- **Acceptance**:
  - Resources deplete correctly over turns
  - Thresholds trigger warning states
  - `SetGenre("fantasy")` uses food/water/stamina vocabulary
- **Validation**: 
  ```bash
  go test -cover ./pkg/resources/...
  ```

---

### Step 7: Implement Party/Crew System

- **Deliverable**:
  - `pkg/crew/crew.go` — Crew member struct (name, health, trait, skill)
  - `pkg/crew/generator.go` — Procedural crew generation from seed
  - `pkg/crew/party.go` — Party management (2-6 slots)
  - `pkg/crew/mortality.go` — Starvation, disease, injury death logic
  - `pkg/crew/genre.go` — `SetGenre()` for name/portrait palette
  - Tests with ≥40% coverage
- **Dependencies**: Step 3 (seed), Step 6 (resources affect crew health)
- **Goal Impact**: Enables "Party/Crew Mortality" design pillar
- **Acceptance**:
  - Crew members generated with consistent names from same seed
  - Crew can die from zero health/food/medicine
  - Party tracks living/dead members
- **Validation**: 
  ```bash
  go test -cover ./pkg/crew/...
  ```

---

### Step 8: Implement Vessel/Transport System

- **Deliverable**:
  - `pkg/vessel/vessel.go` — Vessel struct (integrity, speed, capacity)
  - `pkg/vessel/cargo.go` — Cargo inventory management
  - `pkg/vessel/breakdown.go` — Random breakdown events
  - `pkg/vessel/repair.go` — Repair mechanic (spend materials)
  - `pkg/vessel/genre.go` — `SetGenre()` for vessel vocabulary
  - Tests with ≥40% coverage
- **Dependencies**: Step 3 (seed), Step 6 (resources for repairs)
- **Goal Impact**: Enables "Vessel Integrity" design pillar
- **Acceptance**:
  - Vessel degrades over travel turns
  - Breakdowns trigger based on integrity + RNG
  - Repair consumes cargo materials
  - Destroyed vessel triggers game-over
- **Validation**: 
  ```bash
  go test -cover ./pkg/vessel/...
  ```

---

### Step 9: Implement Time Progression and Movement

- **Deliverable**:
  - `pkg/game/time.go` — Turn-based day/night cycle
  - `pkg/game/movement.go` — Movement costs per terrain type
  - `pkg/game/rest.go` — Rest mechanic for recovery
  - Integration with resources (fuel depletes on move)
- **Dependencies**: Step 4 (world), Step 5 (game loop), Step 6 (resources), Step 8 (vessel)
- **Goal Impact**: Links world traversal to resource consumption
- **Acceptance**:
  - Moving vessel advances turn counter
  - Different terrain costs different fuel amounts
  - Rest action restores morale/health at time cost
- **Validation**: 
  ```bash
  go test ./pkg/game/... -run 'TestMovement|TestTime|TestRest'
  ```

---

### Step 10: Implement Procedural Event System

- **Deliverable**:
  - `pkg/events/queue.go` — Seeded event queue from position + turn
  - `pkg/events/event.go` — Event struct with choices and outcomes
  - `pkg/events/generator.go` — Grammar-based text generation
  - `pkg/procgen/event/grammar.go` — Event text templates
  - `pkg/events/resolution.go` — Choice resolution and outcome application
  - `pkg/events/genre.go` — `SetGenre()` for event vocabulary
  - Tests with ≥40% coverage
- **Dependencies**: Step 3 (seed), Step 6-8 (resources, crew, vessel for outcomes)
- **Goal Impact**: Enables "Procedural Event Stream" design pillar
- **Acceptance**:
  - Events trigger based on position and turn
  - Same seed + position produces same event
  - Choices apply resource/crew/vessel deltas
  - All text generated from grammar, not hardcoded
- **Validation**: 
  ```bash
  go test -cover ./pkg/events/... ./pkg/procgen/event/...
  go-stats-generator analyze ./pkg/events --format json --sections duplication | jq '.duplication.ratio'
  ```

---

### Step 11: Implement Basic Audio Synthesis

- **Deliverable**:
  - `pkg/audio/waveforms.go` — Sine, square, sawtooth, triangle, noise oscillators
  - `pkg/audio/envelope.go` — ADSR envelope system
  - `pkg/audio/sfx.go` — SFX generation (travel, event, crisis, success, death)
  - `pkg/audio/player.go` — Ebitengine audio integration
  - `pkg/audio/genre.go` — `SetGenre()` for timbre presets
- **Dependencies**: Step 5 (game loop for audio playback)
- **Goal Impact**: Enables audio feedback; supports "No Bundled Assets" for audio
- **Acceptance**:
  - SFX play on events (movement, crisis, death)
  - All audio synthesized at runtime (no `.mp3`/`.wav`)
  - Genre affects instrument timbre
- **Validation**: 
  ```bash
  go build ./pkg/audio/...
  ./scripts/validate-no-assets.sh
  ```

---

### Step 12: Implement UI/HUD/Menus

- **Deliverable**:
  - `pkg/ux/hud.go` — Resource panel with bars, crew roster
  - `pkg/ux/worldmap.go` — World map view with vessel position
  - `pkg/ux/events.go` — Event overlay with choices
  - `pkg/ux/menus.go` — Main menu, pause menu, options
  - `pkg/ux/genre.go` — `SetGenre()` for UI skin
- **Dependencies**: Step 5 (rendering), Step 6-10 (game state to display)
- **Goal Impact**: Enables player interaction with all systems
- **Acceptance**:
  - HUD displays all six resources with warning colors
  - Event overlay shows text and clickable choices
  - Menus functional (start game, pause, quit)
- **Validation**: 
  ```bash
  go build ./cmd/voyage && ./voyage  # manual verification
  go-stats-generator analyze ./pkg/ux --format json --sections functions | jq '.functions | length'
  ```

---

### Step 13: Implement Input System

- **Deliverable**:
  - `pkg/config/input.go` — Keyboard/gamepad mapping
  - `pkg/config/rebinding.go` — Rebindable controls
  - `pkg/game/input.go` — Modal input handling (overworld vs event vs menu)
- **Dependencies**: Step 5 (game loop), Step 12 (UI to receive input)
- **Goal Impact**: Enables keyboard/gamepad control
- **Acceptance**:
  - Arrow keys move vessel on overworld
  - Number keys select event choices
  - Escape opens pause menu
  - Controls rebindable via options screen
- **Validation**: 
  ```bash
  go build ./cmd/voyage && ./voyage --help  # check CLI flags
  ```

---

### Step 14: Implement Win/Lose Conditions

- **Deliverable**:
  - `pkg/game/conditions.go` — Win/lose detection logic
  - `pkg/game/endscreen.go` — Run summary display
- **Dependencies**: Step 4 (destination tile), Step 7 (crew survival), Step 8 (vessel destruction)
- **Goal Impact**: Completes the gameplay loop
- **Acceptance**:
  - Win: vessel reaches destination with ≥1 living crew
  - Lose: vessel destroyed OR all crew dead OR morale zero
  - End screen shows days traveled, crew lost, score
- **Validation**: 
  ```bash
  go test ./pkg/game/... -run 'TestWin|TestLose|TestEndConditions'
  ```

---

### Step 15: Implement Save/Load System

- **Deliverable**:
  - `pkg/saveload/save.go` — Serialization (JSON or gob)
  - `pkg/saveload/load.go` — Deserialization
  - `pkg/saveload/slots.go` — Multiple slots with autosave
  - Seed embedded in save for reproducibility
- **Dependencies**: Step 2-14 (all game state to serialize)
- **Goal Impact**: Enables persistence across sessions
- **Acceptance**:
  - Autosave triggers on turn advance
  - Manual save/load via menu
  - Loading save restores identical game state
- **Validation**: 
  ```bash
  go test -cover ./pkg/saveload/...
  ```

---

### Step 16: Implement CLI Flags and Config Persistence

- **Deliverable**:
  - `pkg/config/config.go` — Configuration struct
  - `pkg/config/persistence.go` — Save/load config to disk
  - CLI flags: `--seed`, `--genre`, `--difficulty`
- **Dependencies**: Step 1 (module), Step 13 (input config)
- **Goal Impact**: Enables seed sharing and customization
- **Acceptance**:
  - `./voyage --seed 12345 --genre fantasy` starts with that seed
  - Config persists volume, resolution, key bindings
- **Validation**: 
  ```bash
  ./voyage --help | grep -E 'seed|genre|difficulty'
  ```

---

### Step 17: Create Validation Scripts and CI Pipeline

- **Deliverable**:
  - `scripts/validate-no-assets.sh` — Scan for prohibited file extensions
  - `.github/workflows/ci.yml` — Build, test, vet, validate
- **Dependencies**: Step 1-16 (code to validate)
- **Goal Impact**: Automated quality gates; ensures "No Bundled Assets"
- **Acceptance**:
  - `./scripts/validate-no-assets.sh` passes
  - CI runs on push and PR
  - CI fails if tests fail or vet errors
- **Validation**: 
  ```bash
  ./scripts/validate-no-assets.sh
  # After push: check GitHub Actions
  ```

---

### Step 18: Documentation and Test Coverage Audit

- **Deliverable**:
  - Package-level `doc.go` with examples
  - Exported symbol documentation ≥70%
  - Test coverage ≥40% per package
  - Updated README.md with build/run instructions
- **Dependencies**: Step 1-17 (all code complete)
- **Goal Impact**: Meets stated success criteria
- **Acceptance**:
  - `go doc ./pkg/...` shows documented APIs
  - `go test -cover ./...` reports ≥40% per package
  - README includes quick-start guide
- **Validation**: 
  ```bash
  go-stats-generator analyze . --format json --sections documentation,packages | jq '{doc_coverage: .documentation.coverage, packages: [.packages[].name]}'
  go test -cover ./pkg/... | grep -E 'coverage'
  ```

---

## Dependency Graph

```
Step 1 (Module) ──────────────┬─────────────────────────────────────────────────────┐
        │                     │                                                     │
        ▼                     ▼                                                     ▼
    Step 2 (ECS)          Step 3 (Seed)                                     Step 17 (CI)
        │                     │
        ├─────────────────────┤
        │                     │
        ▼                     ▼
    Step 4 (World) ◄──────────┘
        │
        ▼
    Step 5 (Renderer + Game Loop)
        │
        ├───────────────────────────────────┬─────────────────────┐
        │                                   │                     │
        ▼                                   ▼                     ▼
    Step 6 (Resources)              Step 11 (Audio)          Step 12 (UI)
        │                                                         │
        ├─────────────────────────────────────────────────────────┤
        │                                                         │
        ▼                                                         ▼
    Step 7 (Crew) ◄───────────────────────────────────────► Step 13 (Input)
        │
        ▼
    Step 8 (Vessel)
        │
        ▼
    Step 9 (Time/Movement)
        │
        ▼
    Step 10 (Events)
        │
        ▼
    Step 14 (Win/Lose)
        │
        ▼
    Step 15 (Save/Load)
        │
        ▼
    Step 16 (Config/CLI)
        │
        ▼
    Step 18 (Docs/Tests)
```

---

## Risk Mitigation

| Risk | Mitigation |
|------|------------|
| Ebitengine API changes | Pin to specific version in `go.mod` (v2.9.x) |
| Complexity creep in procedural generators | Apply cyclomatic complexity threshold (≤10) early |
| Grammar-based text becomes repetitive | Design extensible grammar system; parameterize by seed |
| Cross-platform audio issues | Test on Linux, macOS, Windows early in Step 11 |
| Performance with large maps | Implement chunked rendering; profile early |

---

## Sibling Project Reference

The opd-ai/venture project provides a reference implementation for:
- ECS architecture (`pkg/engine/`)
- Procedural generation patterns (`pkg/procgen/`)
- Ebitengine integration (`cmd/venture/`, `pkg/rendering/`)
- Save/load serialization (`pkg/saveload/`)

Voyage should follow similar package structure while adapting for travel-simulator mechanics (resources, crew, vessel) vs action-RPG mechanics.

---

## Post-Implementation Metrics Targets

After completing all steps, validate with:

```bash
go-stats-generator analyze . --format json --sections functions,duplication,documentation,packages | jq '{
  "functions_above_complexity_9": [.functions[] | select(.complexity > 9) | .name],
  "duplication_ratio": .duplication.ratio,
  "doc_coverage": .documentation.coverage,
  "package_count": (.packages | length)
}'
```

**Expected Results**:
- Functions above complexity 9: 0 (or explicitly acknowledged)
- Duplication ratio: < 5%
- Doc coverage: ≥ 70%
- Package count: ~15+ (per ROADMAP architecture)

---

## Estimated Timeline

| Phase | Steps | Duration |
|-------|-------|----------|
| Foundation | 1-3 | 1 week |
| Core Systems | 4-8 | 2 weeks |
| Integration | 9-14 | 2 weeks |
| Polish | 15-18 | 1 week |

**Total**: ~6 weeks for single developer (part-time).

---

*Generated: 2026-04-04 | Source: ROADMAP.md v1.0 milestones + GAPS.md analysis*
