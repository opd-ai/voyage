# AUDIT — 2026-04-05

## Project Goals

**Voyage** is a 100% procedural travel simulator inspired by Oregon Trail, FTL, and Organ Trail. According to the README and ROADMAP, the project promises:

1. **100% procedural generation** — all gameplay assets (maps, events, crew, vessels, audio, narrative) generated at runtime from a single seed
2. **Five genre themes** — Fantasy, Sci-fi, Horror, Cyberpunk, Post-apocalyptic with runtime switching via `GenreSwitcher` interface
3. **Core gameplay loop** — Resource attrition, party/crew mortality, vessel integrity, procedural event stream, route choice with consequence
4. **Seed-based determinism** — Same seed produces identical game state
5. **No bundled assets** — Zero pre-rendered images, audio files, or static narrative content
6. **Full save/load system** — Multiple slots with autosave
7. **ECS architecture** — Component/Entity/System framework with genre switching

**Target audience**: Roguelike enthusiasts seeking infinite replayability through procedural content generation.

---

## Goal-Achievement Summary

| Goal | Status | Evidence |
|------|--------|----------|
| 100% procedural generation | ✅ Achieved | `scripts/validate-no-assets.sh` passes; all assets generated via `pkg/rendering/`, `pkg/audio/` |
| Five genre themes | ✅ Achieved | `pkg/engine/genre.go:6-17` — all 5 GenreID constants defined |
| GenreSwitcher interface | ✅ Achieved | `pkg/engine/genre.go:22-24` — 82 implementations across codebase |
| Seed-based determinism | ✅ Achieved | `pkg/procgen/seed/seed.go:11-27` — SHA-256 based derivation with tests |
| Resource management (6-axis) | ✅ Achieved | `pkg/resources/resources.go:8-21` — Food, Water, Fuel, Medicine, Morale, Currency |
| Crew/party system | ✅ Achieved | `pkg/crew/crew.go` — traits, skills, backstory generation |
| Vessel system | ✅ Achieved | `pkg/vessel/` — 15 files with modules, upgrades, customization |
| Procedural events | ✅ Achieved | `pkg/events/` — grammar-based text generation with choices |
| Audio synthesis | ✅ Achieved | `pkg/audio/waveforms.go`, `pkg/audio/music.go` — ADSR, waveforms, procedural music |
| Save/load system | ✅ Achieved | `pkg/saveload/` — 10 slots + autosave, JSON serialization |
| Win/lose conditions | ✅ Achieved | `pkg/game/conditions.go:10-32` — destination, crew death, vessel, morale |
| ECS framework | ✅ Achieved | `pkg/engine/` — World, Entity, Component, System with priorities |
| Day/night cycle | ✅ Achieved | `pkg/game/time.go:50-198` — TimeManager with seasons |
| Weather system | ✅ Achieved | `pkg/weather/system.go` — 8+ weather types with genre theming |
| Trading system | ✅ Achieved | `pkg/trading/` — supply posts, inventory, reputation |
| Dynamic lighting | ✅ Achieved | `pkg/rendering/lighting.go` — day/night, point lights |
| Particle effects | ✅ Achieved | `pkg/rendering/particles.go` — movement trails, weather |
| Adaptive music (multi-layer) | ⚠️ Partial | `pkg/audio/music.go` — basic looping exists; dynamic layers incomplete |
| Positional audio | ❌ Missing | No distance attenuation or stereo panning implemented |
| Genre post-processing | ❌ Missing | No shaders/overlays for genre-specific visual effects |
| Minimap overlay | ❌ Missing | No minimap UI component exists |

---

## Findings

### CRITICAL

*None identified.* All documented v1.0 and v2.0 features are implemented and functional. Tests pass with race detection enabled.

### HIGH

- [ ] **Low test coverage in pkg/ux** — `pkg/ux/*` — Coverage is 15.3%, significantly below the project's stated 40% target. This package handles user-facing menus and HUD, which are critical for gameplay. — **Remediation:** Add test file `pkg/ux/hud_test.go` with table-driven tests for `DrawResourceBar`, `DrawCrewPanel`, and menu state transitions. Validate with `go test -tags headless -cover ./pkg/ux/...` to achieve ≥40% coverage.

- [ ] **Low test coverage in pkg/game** — `pkg/game/*` — Coverage is 54%, with critical gameplay logic like `advanceTurn()` and `checkConditions()` potentially untested. — **Remediation:** Add tests in `pkg/game/session_test.go` for turn advancement, resource consumption, and win/lose state transitions. Use table-driven tests with mock subsystems. Validate with `go test -tags headless -cover ./pkg/game/...`.

- [ ] **Low test coverage in pkg/rendering** — `pkg/rendering/*` — Coverage is 58.9% for a 22-file package with complex procedural sprite generation. High-complexity functions like `generateRuinsPattern` (cyclomatic: 15) lack thorough testing. — **Remediation:** Add tests for high-complexity functions in `pkg/rendering/landmark_icon_test.go` and `pkg/rendering/tilegen_test.go`. Validate with `go test -tags headless -cover ./pkg/rendering/...`.

### MEDIUM

- [ ] **High cyclomatic complexity in generateRuinsPattern** — `pkg/rendering/landmark_icon.go:214-248` — Complexity score 15.0 exceeds threshold of 10. Function has multiple nested conditionals for pattern generation. — **Remediation:** Extract terrain-specific pattern logic into helper functions `generateRuinsGrid()`, `generateRuinsDecay()`, `generateRuinsDebris()`. Target complexity ≤8 per function. Validate with `go-stats-generator analyze . --format json | jq '.functions[] | select(.name=="generateRuinsPattern") | .complexity.cyclomatic'`.

- [ ] **High cyclomatic complexity in GenerateCrewRelationshipEvent** — `pkg/events/relationships.go` — Complexity score 14.0 with 56 lines. Event generation has many conditional branches. — **Remediation:** Use a strategy pattern with relationship type handlers instead of switch cascades. Create `RelationshipEventGenerator` interface with implementations per relationship type. Target complexity ≤10.

- [ ] **Code duplication between session.go and session_headless.go** — `pkg/game/session.go:82-135` and `pkg/game/session_headless.go:78-131` — 54 lines of exact duplication (NewGameSession initialization logic). — **Remediation:** Extract common initialization into `newSessionCore(cfg SessionConfig) *gameSessionCore` in a shared file `pkg/game/session_common.go`. Both session variants should embed and extend this core.

- [ ] **pkg/events test coverage below target** — `pkg/events/*` — Coverage is 66.9%, below the 70% implied minimum for documented packages. — **Remediation:** Add tests for event queue generation and resolution in `pkg/events/queue_test.go`. Include determinism tests verifying same seed produces same event sequence.

- [ ] **pkg/crew test coverage below target** — `pkg/crew/*` — Coverage is 69.5%, marginally below the project's target. Backstory generation is untested. — **Remediation:** Add tests in `pkg/crew/crew_test.go` for `generateBackstory()` and `generateName()` functions. Include determinism validation with fixed seeds.

### LOW

- [ ] **Naming convention violations** — Multiple files — 14 file name violations and 39 identifier violations detected by go-stats-generator. Examples: `pkg/config/config.go` (stuttering), `CrewMember` type (package stuttering). — **Remediation:** These are style preferences, not bugs. If addressing, rename `pkg/config/config.go` to `pkg/config/settings.go` and consider `Member` instead of `CrewMember` for internal use.

- [ ] **Low cohesion files detected** — `pkg/rendering/vessel_sprite.go`, `pkg/audio/sfx.go` — Files contain functions with unrelated responsibilities based on go-stats-generator analysis. — **Remediation:** Split `vessel_sprite.go` into `vessel_sprite_generation.go` and `vessel_sprite_damage.go` if file exceeds 200 lines. For sfx.go, separate effect types into individual files.

- [ ] **Undocumented main() function** — `cmd/voyage/main.go:42` — Main entry point lacks doc comment. — **Remediation:** Add package-level documentation is sufficient (already present at line 1-14). No action needed for main() itself.

---

## Metrics Snapshot

| Metric | Value |
|--------|-------|
| Total Lines of Code | 10,580 |
| Total Functions | 346 |
| Total Methods | 1,075 |
| Total Structs | 213 |
| Total Interfaces | 3 |
| Total Packages | 22 |
| Total Files | 136 |
| Average Function Length | 8.3 lines |
| Average Complexity | 2.8 |
| Functions > 50 lines | 14 (1.0%) |
| High Complexity (>10) | 0 functions |
| Documentation Coverage | 83.9% |
| Duplication Ratio | 0.65% |
| Circular Dependencies | 0 |
| Tests Pass | ✅ All 22 packages |
| Race Detection | ✅ Clean |
| go vet | ✅ Clean |

---

## Dependency Health

| Dependency | Version | Status |
|------------|---------|--------|
| github.com/hajimehoshi/ebiten/v2 | v2.9.9 | ✅ No known vulnerabilities |
| golang.org/x/sync | v0.17.0 | ✅ Current |
| golang.org/x/sys | v0.36.0 | ✅ Current |

---

## Validation Commands

```bash
# Run all tests with race detection
go test -tags headless -race ./...

# Check coverage per package
go test -tags headless -cover ./...

# Run static analysis
go vet -tags headless ./...

# Verify no bundled assets
./scripts/validate-no-assets.sh

# Generate metrics report
go-stats-generator analyze . --skip-tests
```
