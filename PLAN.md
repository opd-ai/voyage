# Implementation Plan: Game Loop Integration & Playability

## Project Context
- **What it does**: 100% procedural travel simulator inspired by Oregon Trail, FTL, and Organ Trail — every map, event, crew, vessel, audio, and narrative generated from a single seed.
- **Current goal**: Complete the game loop integration to make the game playable (GAPS.md #1 — CRITICAL priority)
- **Estimated Scope**: Medium (7 implementation steps, 5–15 items above threshold across metrics)

## Goal-Achievement Status
| Stated Goal | Current Status | This Plan Addresses |
|-------------|---------------|---------------------|
| Playable game loop | ❌ Not integrated | Yes — Primary focus |
| ECS Framework with GenreSwitcher | ✅ Achieved | No |
| All 22 subsystems implemented | ✅ Achieved | Yes — Wire together |
| Packages with tests | ⚠️ 3 packages missing tests | Yes — Step 7 |
| Crew relationship events | ⚠️ Data structure exists, not used | Yes — Step 6 |
| F3 debug toggle bug | ⚠️ Known bug (AUDIT #HIGH) | Yes — Step 5 |
| Animated sprites (v3.0) | ❌ Planned for future | No |
| Adaptive music (v3.0) | ❌ Planned for future | No |

## Metrics Summary
- **Complexity hotspots**: 1 function above threshold (`applyAlignmentVariance` @ 8 cyclomatic — below 9 threshold)
- **Functions above complexity 7**: 1 (0.08% of 1,230 functions)
- **Duplication ratio**: 0.36% (7 clone pairs, 84 duplicated lines)
- **Doc coverage**: 84.2% (above project's 40% minimum)
- **Package coupling hotspots**: `ux` (5.0), `game` (4.5) — these are expected as integration points
- **Test gaps**: `cmd/voyage`, `pkg/procgen/event`, `pkg/world`

## Implementation Steps

### Step 1: Create Game Session Orchestrator
- **Deliverable**: New file `pkg/game/session.go` containing `GameSession` struct that initializes and coordinates all subsystems (world, crew, vessel, resources, events, audio, rendering, UI).
- **Dependencies**: None (foundational step)
- **Goal Impact**: Enables game loop integration — the core blocker preventing playability
- **Acceptance**: 
  - `GameSession` struct compiles and initializes all subsystems
  - Unit test `TestGameSessionInit` passes in `pkg/game/session_test.go`
- **Validation**: 
  ```bash
  go build ./pkg/game/...
  go test -tags headless -v ./pkg/game/... -run TestGameSessionInit
  ```

### Step 2: Wire Subsystems in main.go
- **Deliverable**: Replace TODO block in `cmd/voyage/main.go:80-91` with actual system initialization:
  - Instantiate `GameSession` from Step 1
  - Initialize world map via `pkg/procgen/world.Generator`
  - Create crew party via `pkg/crew.Party`
  - Create vessel via `pkg/vessel.Vessel`
  - Initialize resources via `pkg/resources.Manager`
  - Set up event queue via `pkg/events.Generator`
  - Initialize audio via `pkg/audio.Player`
  - Initialize rendering via `pkg/rendering.Renderer`
  - Initialize UI via `pkg/ux.Manager`
- **Dependencies**: Step 1
- **Goal Impact**: Connects 22 implemented packages to the entry point
- **Acceptance**: 
  - `go build ./cmd/voyage` succeeds
  - Running `./voyage --seed 12345` initializes all systems (verified via log output)
- **Validation**: 
  ```bash
  go build -tags headless ./cmd/voyage && ./voyage --seed 12345 2>&1 | grep -c "initialized"
  ```

### Step 3: Implement Ebitengine Game Loop
- **Deliverable**: Replace `demo()` call in `cmd/voyage/main.go:98` with `ebiten.RunGame(gameSession)`:
  - `GameSession` implements `ebiten.Game` interface (`Update`, `Draw`, `Layout`)
  - `Update` advances game state (turn progression, event resolution, resource consumption)
  - `Draw` renders the current game state via the rendering subsystem
  - `Layout` returns configured window dimensions
- **Dependencies**: Step 2
- **Goal Impact**: Makes the game actually run with a visual window and interactive gameplay
- **Acceptance**: 
  - Running `./voyage` opens a graphical window with the world map visible
  - Pressing Escape opens the pause menu
  - Game loop runs at stable 60 FPS (verified via Ebitengine debug info)
- **Validation**: 
  ```bash
  go build ./cmd/voyage && timeout 5 ./voyage --seed 42 || echo "Window opened"
  go test -tags headless -v ./pkg/game/... -run TestGameLoop
  ```

### Step 4: Connect Turn Progression to Events and Resources
- **Deliverable**: Modify `GameSession.Update()` to:
  - Advance turn counter on player movement
  - Consume resources per turn (food, water, fuel per `pkg/resources.Manager`)
  - Generate and queue events based on position and turn via `pkg/events.Generator`
  - Check win/lose conditions via `pkg/game/conditions.go`
- **Dependencies**: Step 3
- **Goal Impact**: Enables the core gameplay loop of resource management and event resolution
- **Acceptance**: 
  - Moving the vessel consumes fuel
  - Food/water deplete each turn
  - Events trigger based on map position
  - Running out of resources triggers game over
- **Validation**: 
  ```bash
  go test -tags headless -v ./pkg/game/... -run TestTurnProgression
  go test -tags headless -v ./pkg/game/... -run TestResourceDepletion
  go test -tags headless -v ./pkg/game/... -run TestEventGeneration
  ```

### Step 5: Fix F3 Debug Toggle Bug
- **Deliverable**: Modify `pkg/game/game.go:143-145` to use key release detection instead of key press:
  - Add `f3WasPressed bool` field to game struct
  - Toggle debug mode only on key-down transition (was not pressed, now pressed)
  - This is a documented BUG at line 141 in the AUDIT
- **Dependencies**: Step 3 (requires working game loop to test)
- **Goal Impact**: Fixes user-facing bug that prevents proper debug overlay usage
- **Acceptance**: 
  - Pressing F3 once toggles debug mode exactly once
  - Holding F3 does not cause repeated toggles
- **Validation**: 
  ```bash
  go test -tags headless -v ./pkg/game/... -run TestDebugToggle
  ```

### Step 6: Connect Crew Relationship Network to Events
- **Deliverable**: Modify `pkg/events/generator.go` to query `pkg/crew/relationship.go`:
  - Add `GenerateCrewRelationshipEvent(network *crew.RelationshipNetwork) *Event` method
  - Check for strong bonds (generate cooperation events)
  - Check for rivalries (generate conflict events)
  - Add relationship delta as event outcome
- **Dependencies**: Step 4 (requires event system integration)
- **Goal Impact**: Activates the crew relationship system that exists but has no gameplay effect (GAPS.md #7)
- **Acceptance**: 
  - Crew pairs with high bond/rivalry values trigger relationship-specific events
  - Event outcomes modify relationship values
- **Validation**: 
  ```bash
  go test -tags headless -v ./pkg/events/... -run TestRelationshipEvents
  go test -tags headless -v ./pkg/crew/... -run TestRelationshipNetwork
  ```

### Step 7: Add Missing Test Coverage
- **Deliverable**: Create test files for packages flagged as `[no test files]`:
  - `cmd/voyage/main_test.go` — test flag parsing, genre validation, seed handling
  - `pkg/procgen/event/event_test.go` — test deterministic event generation
  - `pkg/world/world_test.go` — test world state management
- **Dependencies**: Steps 1–4 (tests validate the integrated system)
- **Goal Impact**: Fulfills CONTRIBUTING.md requirement of "≥40% coverage per package" for untested packages
- **Acceptance**: 
  - `go test -cover ./...` shows no `[no test files]` output
  - All three packages have ≥40% coverage
- **Validation**: 
  ```bash
  go test -tags headless -v ./cmd/voyage/... ./pkg/procgen/event/... ./pkg/world/...
  go test -tags headless -cover ./cmd/voyage/... ./pkg/procgen/event/... ./pkg/world/... | grep -E "coverage|PASS|FAIL"
  ```

## Dependency Graph

```
Step 1 (GameSession)
    ↓
Step 2 (Wire main.go)
    ↓
Step 3 (Ebitengine loop)
    ↓
Step 4 (Turn/Events/Resources) → Step 5 (F3 bug fix)
    ↓
Step 6 (Relationship events)
    ↓
Step 7 (Test coverage)
```

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Subsystem interface mismatch | Low | Medium | All subsystems implement `SetGenre()` — verify with `go build` |
| Ebitengine headless mode issues | Medium | Low | Use `-tags headless` consistently in CI |
| Event system performance | Low | Low | Current avg complexity 2.7 — well within bounds |
| Package coupling increase | Medium | Medium | Keep `GameSession` thin — delegate to subsystems |

## Out of Scope (v3.0 Features)

The following ROADMAP v3.0 items are **not** included in this plan — they require the playable foundation first:
- Animated overworld tiles
- Crew member portrait animation
- Vessel damage state sprites
- Adaptive multi-layer music
- Positional audio
- Genre post-processing shaders
- Dynamic minimap overlay

## Success Criteria

This plan is complete when:
1. `go build ./cmd/voyage && ./voyage --seed 42` opens a playable window
2. Player can navigate the world map and trigger events
3. Resources deplete over time
4. Win/lose conditions fire correctly
5. `go test -tags headless -race ./...` passes with no `[no test files]` gaps
6. F3 debug toggle works correctly

## Validation Commands Summary

```bash
# Full validation sequence
go build -tags headless ./...
go test -tags headless -race ./...
go vet -tags headless ./...
./scripts/validate-no-assets.sh

# Metrics verification (after implementation)
go-stats-generator analyze . --skip-tests --format json | jq '{
  test_gaps: [.packages[] | select(.documentation.quality_score == 0) | .name],
  complexity_hotspots: [.functions[] | select(.complexity.cyclomatic > 9) | .name],
  duplication: .duplication.duplication_ratio
}'
```

---

*Generated by data-driven analysis on 2026-04-05*
*Metrics source: go-stats-generator v1.0.0*
*Priority derived from: GAPS.md, AUDIT.md, ROADMAP.md*
