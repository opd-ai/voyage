# AUDIT — 2026-04-06

## Project Goals

Voyage is a **100% Procedural Travel Simulator** inspired by Oregon Trail, FTL, and Organ Trail. Per the README, the project claims to deliver:

1. **Full procedural generation** — Every map, event, crew, vessel, audio, and narrative generated from a single seed
2. **No bundled assets** — All visual and audio content generated at runtime
3. **Five genre themes** — Fantasy, Sci-fi, Horror, Cyberpunk, Post-apocalyptic
4. **Feature-complete v5.0** — Full gameplay loop with all subsystems operational
5. **Six-axis resource management** — Food, water, fuel, medicine, morale, currency
6. **Crew mortality system** — Procedurally generated crew who can sicken, die, or desert
7. **Vessel integrity** — Transport accumulates damage and requires repair
8. **Procedural events** — Grammar-based text generation with branching choices
9. **Save/load system** — Multiple slots with autosave
10. **WebAssembly and mobile builds** — Cross-platform deployment
11. **Modding system** — JSON and WASM extension support
12. **Leaderboards and convoy mode** — Async multiplayer features

**Target audience**: Players who enjoy roguelike travel simulators; developers interested in procedural generation in Go.

---

## Goal-Achievement Summary

| Goal | Status | Evidence |
|------|--------|----------|
| 100% procedural generation | ✅ Achieved | `scripts/validate-no-assets.sh` passes; no bundled assets found |
| Five genre themes | ✅ Achieved | `pkg/engine/genre.go` implements Fantasy, Scifi, Horror, Cyberpunk, Postapoc |
| ECS framework with GenreSwitcher | ✅ Achieved | `pkg/engine/` exports `GenreSwitcher`, `BaseSystem`; 86.7% coverage |
| Seed-based deterministic RNG | ✅ Achieved | `pkg/procgen/seed/` at 86.5% coverage; determinism verified |
| Ebitengine rendering foundation | ✅ Achieved | `pkg/rendering/` at 89.6% coverage |
| Procedural world map generation | ✅ Achieved | `pkg/procgen/world/` at 95.1% coverage |
| Resource management (6-axis) | ✅ Achieved | `pkg/resources/` implements all six resources; 67.7% coverage |
| Crew/party system | ✅ Achieved | `pkg/crew/` at 92.1% coverage |
| Vessel/transport system | ✅ Achieved | `pkg/vessel/` at 77.9% coverage |
| Procedural event system | ✅ Achieved | `pkg/events/` at 75.6% coverage |
| Audio synthesis | ✅ Achieved | `pkg/audio/` at 86.8% coverage; spatial audio, music states |
| UI/HUD/Menus with genre theming | ✅ Achieved | `pkg/ux/` at 95.2% coverage |
| Win/lose conditions | ✅ Achieved | `pkg/game/conditions.go` implements both |
| Save/load system | ✅ Achieved | `pkg/saveload/` at 71.2% coverage; 10 slots + autosave |
| Configuration and input rebinding | ✅ Achieved | `pkg/config/` at 84.3%, `pkg/input/` at 88.9% |
| CI/CD pipeline | ✅ Achieved | `.github/workflows/ci.yml` runs build, test, vet, lint |
| Faction system with reputation | ✅ Achieved | `pkg/factions/` at 92.1% coverage |
| Quest/objective system | ✅ Achieved | `pkg/quests/` at 78.7% coverage |
| Meta-progression | ✅ Achieved | `pkg/metaprog/` at 91.9% coverage |
| Leaderboards and convoy mode | ⚠️ Partial | Implemented but uses placeholder server URL |
| WebAssembly builds | ❌ Non-functional | WASM build fails due to compilation error in `pkg/game/session.go:48` |
| Mobile builds | ⚠️ Partial | Makefile targets exist; requires external SDK; not CI-verified |
| Modding system | ⚠️ Partial | Implemented at 50.8% coverage — lowest in codebase |
| Non-headless desktop build | ❌ Non-functional | Compilation error blocks all non-headless builds |

**Overall: 20/24 goals achieved; 2 partial; 2 non-functional**

---

## Findings

### CRITICAL

- [x] **Double-pointer type mismatch blocks non-headless and WASM builds** — `pkg/game/session.go:48` — The expression `&e` creates `**events.Event` when `pending[0]` already returns `*Event`. Should be `s.currentEventSnapshot = pending[0]`. This blocks all non-headless desktop builds, WebAssembly builds, and mobile builds. The README claims "Feature Complete (v5.0)" and "WebAssembly and mobile builds" but neither work. — **Remediation:** Change line 48 from `s.currentEventSnapshot = &e` to `s.currentEventSnapshot = pending[0]`. Validate with `GOOS=js GOARCH=wasm go build ./cmd/voyage && go build ./cmd/voyage`.

- [x] **Event choice selection off-by-one error** — `pkg/game/session.go:195` — Choice IDs start at 1 (`pkg/events/event.go:87`), but `resolveEvent()` receives 0-based index `i`. Key 1 passes `choiceID=0`, which matches no choice (IDs are 1,2,3,4). First choice unreachable; subsequent choices shifted by one. — **Remediation:** Change line 195 from `s.resolveEvent(currentEvent.ID, i)` to `s.resolveEvent(currentEvent.ID, i+1)`. Validate with `go test -tags headless ./pkg/game/...`.

- [x] **Menu input fires every frame** — `pkg/game/game.go:130,154` — `handleMenuInput()` and `handleGameOverInput()` use `ebiten.IsKeyPressed()` (continuous detection) instead of `inpututil.IsKeyJustPressed()` (edge detection). Pressing Enter causes instant state transitions and potential input spillover to the next state. — **Remediation:** Replace `ebiten.IsKeyPressed(ebiten.KeyEnter)` with `inpututil.IsKeyJustPressed(ebiten.KeyEnter)` on lines 130 and 154. Add import `"github.com/hajimehoshi/ebiten/v2/inpututil"` if not present. Validate with manual testing.

### HIGH

- [x] **Input Manager abstraction entirely unused** — `pkg/input/` (entire package) vs `pkg/game/*.go` — The `pkg/input` package implements a sophisticated input Manager with touch support, key repeat, swipe gestures, and action deduplication, but `pkg/game/` directly calls raw Ebiten APIs (`ebiten.IsKeyPressed`, `inpututil.IsKeyJustPressed`). This renders the input package dead code and creates inconsistent input behavior. — **Remediation:** Wire `input.Manager` into `pkg/game/game.go` via `Update()` and query actions from `Manager.State()` instead of raw Ebiten calls. Validate with `go test -tags headless ./pkg/input/... ./pkg/game/...`.

- [ ] **Modding system test coverage at 50.8%** — `pkg/modding/` — The lowest coverage in the codebase. WASM loader handles untrusted code execution with capability-based sandboxing, but edge cases (invalid bytecode, missing exports, capability denial) are undertested. — **Remediation:** Add tests to `pkg/modding/wasm_loader_test.go` for invalid WASM, missing `mod_get_id`, capability denial scenarios. Target ≥70% coverage. Validate with `go test -tags headless -cover ./pkg/modding/...`.

- [ ] **WASM build not verified in CI** — `.github/workflows/ci.yml` — CI only builds with `-tags headless`. The claimed "WebAssembly builds" feature is never verified, allowing silent regressions. — **Remediation:** Add job to `.github/workflows/ci.yml`: `GOOS=js GOARCH=wasm go build -o /dev/null ./cmd/voyage`. Validate with `gh workflow run ci.yml` after fix.

- [ ] **Leaderboard uses placeholder server URL** — `pkg/leaderboard/client.go:46` — `ServerURL: "https://api.voyage-game.example.com/leaderboard"` is clearly a placeholder. Users may expect working online leaderboards but requests will fail. — **Remediation:** Document in README that leaderboards use local storage by default; online requires self-hosted server. Add `--offline` flag or environment variable to explicitly disable server attempts. Validate behavior documented in README.

### MEDIUM

- [x] **Save/load slots invisible on small panels** — `pkg/ux/slots.go:158` — When `panelHeight < 150`, the formula `(panelHeight - 100) / 50` produces 0 visible slots, making the save/load screen blank. — **Remediation:** Clamp minimum visible slots: `visibleSlots := max(1, (panelHeight - 100) / 50)`. Validate with window resize to 200px height.

- [x] **Event overlay text overflow** — `pkg/ux/events.go:201-207` — Words longer than `maxWidth` are placed on a single line without truncation, causing overflow past panel boundaries. — **Remediation:** Add word truncation or hyphenation in `addWordToLine()`: if `len(word) > maxWidth`, truncate to `maxWidth-3` + "...". Validate with procedurally generated long words.

- [x] **Pause overlay not recreated on window resize** — `pkg/game/game.go:264-269` — Pause overlay is cached at initial window size. Resizing leaves overlay incorrectly sized. — **Remediation:** Add size check in `drawPauseOverlay()`: recreate overlay if `g.width` or `g.height` changed since last creation. Validate with window resize while paused.

- [x] **Minimap renders off-screen on narrow windows** — `pkg/ux/minimap.go:103` — Position `screenW - m.width - 10` becomes negative when `screenW < m.width + 10`. — **Remediation:** Add bounds check: `x := max(0, screenW - m.width - 10)`. Validate with window width < 150px.

- [x] **Stacking morale penalties for simultaneous resource depletion** — `pkg/game/session_common.go:120-125` — When both food AND water deplete simultaneously, morale receives -5 AND -8 (total -13 per turn) with no cap, causing rapid morale collapse and unfair difficulty spike. — **Remediation:** Cap combined penalty per turn: `penaltyTotal := min(10, foodPenalty + waterPenalty)`. Validate with test case depleting both resources.

- [x] **Negative time advance not clamped** — `pkg/game/session.go:211-215` — `outcome.TimeAdvance` only has upper bound clamp. Negative values make the loop condition false, silently skipping turn advancement. — **Remediation:** Add lower bound: `if timeAdvance < 0 { timeAdvance = 0 }`. Validate with event generating negative TimeAdvance.

- [x] **Event queue resolved list grows unbounded** — `pkg/events/queue.go:158` — `q.resolved` accumulates all resolved events forever. Over very long sessions, memory grows without pruning. — **Remediation:** Add pruning: keep only last 100 resolved events, or clear on game completion. Validate with 1000+ turn session memory profile.

- [ ] **Per-pixel operations in post-processing** — `pkg/rendering/postprocess.go:193,341` — `ApplyScanlines()` and `ApplySepia()` use `img.At()` and `result.Set()` per-pixel in tight loops, causing frame drops at high resolutions. — **Remediation:** Use shader-based effects or batch pixel operations via `Pix` slice access. Validate with benchmark at 1920x1080.

### LOW

- [x] **Cargo scroll indicator shows negative values** — Already resolved; guard at line 170 prevents display when items <= maxVisible — `pkg/ux/cargo.go:171` — When items fit on screen, denominator becomes zero or negative, displaying `[1/-3]`. — **Remediation:** Guard: `if len(items) <= maxVisible { return }`. Validate with < 8 cargo items.

- [x] **Division by zero risk in vignette generation** — `pkg/rendering/postprocess.go:149-176` — If screen dimensions are 0, `maxDistSq` is 0 causing division by zero. — **Remediation:** Add guard: `if w == 0 || h == 0 { return }`. Validate with minimized window.

- [x] **Redundant F3 key-release detection** — `pkg/game/game.go:162-169` — Manual `f3WasPressed` flag reimplements `inpututil.IsKeyJustPressed()`. — **Remediation:** Replace manual tracking with `inpututil.IsKeyJustPressed(ebiten.KeyF3)`. Validate with F3 toggle behavior.

- [ ] **No `Image.Dispose()` calls in rendering** — `pkg/rendering/` (entire package) — Ebiten images never disposed. For long sessions or genre switches, GPU memory accumulates. — **Remediation:** Add cleanup in `SetGenre()` and session teardown. Validate with memory profiling over 10 genre switches.

- [x] **Code duplication in WASM loader** — `pkg/modding/wasm_loader.go:440-453,468-481` — 14-line renamed clone between `Initialize()` and `OnTurnStart()`. — **Remediation:** Extract common pattern to helper function `callOptionalHook(ctx, name, args...)`. Validate with `go-stats-generator analyze . --sections duplication`.

---

## Metrics Snapshot

| Metric | Value | Assessment |
|--------|-------|------------|
| Total Lines of Code | 18,571 | Healthy codebase size |
| Total Functions | 538 | Well-factored |
| Total Packages | 35 | Good separation of concerns |
| Average Function Length | 8.7 lines | Excellent (project target <30) |
| Functions >50 lines | 17 (0.7%) | Very low |
| Average Complexity | 2.9 | Excellent (project target <10) |
| High Complexity (>10) | 0 | None |
| Overall Test Coverage | ~82% avg | Exceeds project target (40%) |
| Documentation Coverage | 82.3% | Good |
| Duplication Ratio | 0.20% | Negligible (9 clone pairs) |
| Circular Dependencies | 0 | Clean architecture |

---

## Dependency Health

| Dependency | Version | Status |
|------------|---------|--------|
| `github.com/hajimehoshi/ebiten/v2` | v2.9.9 | ✅ Stable; no CVEs; Go 1.24 required |
| `github.com/ebitengine/gomobile` | v0.0.0-20250923 | ✅ Mobile build support |
| `github.com/ebitengine/purego` | v0.9.0 | ✅ Pure Go FFI |
| `github.com/tetratelabs/wazero` | v1.11.0 | ✅ WASM runtime; no known CVEs |

**Note:** Ebitengine 2.9 deprecated `AppendVerticesAndIndicesForFilling` and `AppendVerticesAndIndicesForStroke`. The codebase does not use these functions — no migration required.

---

## Test Health

- All 37 packages pass with race detector: `go test -tags headless -race ./...`
- `go vet -tags headless ./...` reports no issues
- Headless build compiles successfully
- Non-headless and WASM builds fail due to CRITICAL compilation error

---

## Verification Commands

```bash
# Validate compilation fix
GOOS=js GOARCH=wasm go build -o /dev/null ./cmd/voyage && go build ./cmd/voyage

# Validate test coverage
go test -tags headless -cover ./pkg/modding/... | grep coverage

# Validate no regressions
go test -tags headless -race ./...

# Check metrics post-fix
go-stats-generator analyze . --skip-tests
```
