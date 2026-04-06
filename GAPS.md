# Implementation Gaps — 2026-04-06

This document identifies gaps between the project's stated goals (per README) and the current implementation state.

---

## Non-Headless and WASM Builds Do Not Compile

- **Stated Goal**: README claims "Feature Complete (v5.0)" and "WebAssembly and mobile builds ✅"
- **Current State**: A compilation error in `pkg/game/session.go:48` (`cannot use &e (value of type **events.Event) as *events.Event`) blocks all non-headless desktop builds, WebAssembly builds, and mobile builds. Only the headless build compiles.
- **Impact**: Users cannot run the game on desktop with rendering, cannot deploy to web via WASM, and cannot build for mobile. The core product is non-functional outside of headless testing.
- **Closing the Gap**:
  1. Change line 48 from `s.currentEventSnapshot = &e` to `s.currentEventSnapshot = pending[0]`
  2. Validate with `go build ./cmd/voyage && GOOS=js GOARCH=wasm go build ./cmd/voyage`
  3. Add WASM build verification to CI to prevent regression

---

## Event Choice Selection Is Off-by-One

- **Stated Goal**: "Procedural Event Stream — Grammar-based text generation with branching choices" should allow players to select any choice presented
- **Current State**: In `pkg/game/session.go:195`, the loop index `i` (0-based) is passed to `resolveEvent()`, but choice IDs start at 1 (`pkg/events/event.go:87`). Pressing key 1 sends `choiceID=0`, which matches no choice. Key 2 selects choice 1, etc.
- **Impact**: The first choice in every event is unreachable. All other choices are shifted by one. Core gameplay is broken.
- **Closing the Gap**:
  1. Change `s.resolveEvent(currentEvent.ID, i)` to `s.resolveEvent(currentEvent.ID, i+1)` on line 195
  2. Add unit test asserting key 1 selects choice ID 1
  3. Validate with `go test -tags headless ./pkg/game/...`

---

## Menu and Game Over Input Fires Continuously

- **Stated Goal**: Clean state transitions controlled by player input
- **Current State**: `pkg/game/game.go:130,154` use `ebiten.IsKeyPressed()` (continuous detection) for menu transitions. Pressing Enter fires every frame, causing instant state changes and potential input spillover.
- **Impact**: Menu navigation is unreliable. State may skip immediately or register multiple times.
- **Closing the Gap**:
  1. Replace `ebiten.IsKeyPressed()` with `inpututil.IsKeyJustPressed()` on lines 130 and 154
  2. Add import for `inpututil` if missing
  3. Manually test menu navigation flow

---

## Input Manager Package Is Dead Code

- **Stated Goal**: Robust input handling with touch, swipe, and key repeat support
- **Current State**: `pkg/input/` implements a comprehensive Manager with 88.9% test coverage, but `pkg/game/` directly calls raw Ebiten APIs, completely bypassing the Manager.
- **Impact**: Touch and swipe gestures don't work. Key repeat timing is inconsistent. The 400+ lines of input abstraction provide no value to users.
- **Closing the Gap**:
  1. Instantiate `input.Manager` in `Game` struct
  2. Call `manager.Update()` at start of `Game.Update()`
  3. Query `manager.State()` instead of raw Ebiten calls throughout `pkg/game/`
  4. Remove redundant direct Ebiten input calls
  5. Validate touch input on a touchscreen device or emulator

---

## Leaderboard Server URL Is Placeholder

- **Stated Goal**: README claims "Leaderboards and async convoy mode ✅"
- **Current State**: `pkg/leaderboard/client.go:46` uses `https://api.voyage-game.example.com/leaderboard` — clearly a placeholder domain that will fail on any real network request.
- **Impact**: Users expecting online leaderboards will get silent failures. Feature appears broken.
- **Closing the Gap**:
  1. Document in README that leaderboards default to local storage
  2. Explain that online leaderboards require self-hosted server
  3. Add `--offline` flag or environment variable for explicit mode control
  4. Optionally: provide reference server implementation or link to hosting guide

---

## WASM Build Not Verified in CI

- **Stated Goal**: README lists "WebAssembly and mobile builds ✅" as implemented
- **Current State**: `.github/workflows/ci.yml` only builds with `-tags headless`. WASM target is never compiled in CI.
- **Impact**: The current compilation error went undetected. Future regressions could silently break the web version.
- **Closing the Gap**:
  1. Add CI step: `GOOS=js GOARCH=wasm go build -o /dev/null ./cmd/voyage`
  2. Optionally add smoke test loading WASM in headless browser
  3. Document mobile build prerequisites (Android/iOS SDK) in CONTRIBUTING.md

---

## Modding System Has Lowest Test Coverage

- **Stated Goal**: "Modding system ✅" with JSON and WASM extension support
- **Current State**: `pkg/modding/` at 50.8% coverage — the only package below the project's stated 40% minimum that handles meaningful functionality. WASM loader executes untrusted code with capability sandboxing, but edge cases are untested.
- **Impact**: Bugs in mod loading, capability enforcement, or error handling may exist undetected. Security-critical code paths need higher confidence.
- **Closing the Gap**:
  1. Add tests for invalid WASM bytecode handling
  2. Add tests for missing required exports (`mod_get_id`)
  3. Add tests for capability denial scenarios
  4. Add tests for multi-mod loading with conflicting IDs
  5. Target ≥70% coverage; validate with `go test -cover ./pkg/modding/...`

---

## Resource Package Below Average Coverage

- **Stated Goal**: "Resource management system (six-axis model) ✅" is a core design pillar
- **Current State**: `pkg/resources/` at 67.7% coverage is below the project average (82%) and handles critical gameplay logic.
- **Impact**: Resource exhaustion edge cases, overflow/underflow, and genre-specific naming may have untested bugs.
- **Closing the Gap**:
  1. Add tests for 0-resource scenarios (depletion handling)
  2. Add tests for resource overflow past max values
  3. Add tests for genre-specific resource names
  4. Target ≥80% coverage

---

## `--mods-dir` Flag Undocumented in README

- **Stated Goal**: Modding system should be discoverable and usable
- **Current State**: `docs/MODDING.md` references `--mods-dir` flag, but README's "Available Options" table doesn't include it.
- **Impact**: Users may not discover the ability to customize mod directory location.
- **Closing the Gap**:
  1. Add `--mods-dir` to README Available Options table
  2. Add usage example: `./voyage --mods-dir ~/.local/share/voyage/mods/`

---

## Summary

| Gap | Severity | Blocks Goal |
|-----|----------|-------------|
| Compilation error in session.go | CRITICAL | WASM builds, non-headless builds |
| Event choice off-by-one | CRITICAL | Event system functionality |
| Menu input continuous fire | CRITICAL | State machine integrity |
| Input Manager unused | HIGH | Touch/mobile support |
| Leaderboard placeholder URL | HIGH | Online leaderboards |
| WASM CI verification missing | HIGH | Build reliability |
| Modding test coverage | HIGH | Modding system confidence |
| Resource test coverage | MEDIUM | Core gameplay confidence |
| `--mods-dir` documentation | LOW | Feature discoverability |

**Priority order for closing gaps:**
1. Fix compilation error (unblocks all other testing)
2. Fix event choice off-by-one (restores core gameplay)
3. Fix menu input detection (restores state machine)
4. Add WASM CI verification (prevents regression)
5. Wire input Manager (enables touch support)
6. Improve modding test coverage
7. Document leaderboard behavior
8. Improve resource test coverage
9. Document `--mods-dir` flag
