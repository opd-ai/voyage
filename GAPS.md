# Implementation Gaps — 2026-04-05

This document identifies gaps between the project's stated goals (per README and ROADMAP) and the current implementation state.

---

## Adaptive Multi-Layer Music System

- **Stated Goal**: ROADMAP v3.0 promises a dynamic layer system for music with states (peaceful travel, crisis, encounter, victory, death), biome-specific ambient parameters, smooth cross-fades between intensity states, and genre-specific instrument mapping.

- **Current State**: `pkg/audio/music.go` implements a basic looping ambient music generator with three layers (bass, pad, melody hints) and genre-specific parameters. However:
  - No dynamic state transitions between peaceful/crisis/encounter modes
  - No biome-aware ambient parameters
  - No cross-fade implementation between intensity states
  - Music parameters are static per loop, not reactive to gameplay

- **Impact**: Players experience the same audio intensity regardless of gameplay context. A zombie attack feels the same as peaceful travel, reducing immersion and audio feedback for game state.

- **Closing the Gap**: 
  1. Add `MusicState` enum in `pkg/audio/music.go` with values `Peaceful`, `Tense`, `Combat`, `Victory`, `Death`
  2. Implement `SetMusicState(state MusicState)` on `MusicGenerator` that modifies BPM, melody density, and waveform selection
  3. Add `CrossfadeTo(target []float64, duration float64)` for smooth transitions
  4. Wire music state changes to game events in `pkg/game/session.go`

---

## Positional Audio / Spatial SFX

- **Stated Goal**: ROADMAP v3.0 specifies distance attenuation for offscreen events, left/right stereo panning for spatial awareness, and ambient loops per biome/region.

- **Current State**: `pkg/audio/sfx.go` generates mono sound effects with no positional information. `pkg/audio/player.go` plays audio without spatial calculations.

- **Impact**: Players cannot perceive the direction or distance of events. An attack from the left sounds identical to one from the right, limiting tactical awareness in encounters.

- **Closing the Gap**:
  1. Add `SpatialAudioConfig` struct with `ListenerPos`, `SourcePos`, `MaxDistance`, `Falloff` fields
  2. Implement `ApplySpatialAudio(samples []float64, config SpatialAudioConfig) (left, right []float64)` in `pkg/audio/spatial.go`
  3. Modify `Player.PlaySFX()` to accept optional position parameter
  4. Calculate stereo pan based on angle from listener to source

---

## Genre Post-Processing Visual Effects

- **Stated Goal**: ROADMAP v3.0 promises genre-specific visual post-processing:
  - Fantasy: warm desaturated vignette, bloom on magic
  - Sci-fi: scanline overlay, chromatic aberration
  - Horror: desaturate + red-tint at low health, film grain
  - Cyberpunk: neon bloom, CRT curvature, glitch artifacts
  - Post-apocalyptic: sepia wash, dust overlay, heavy vignette

- **Current State**: `pkg/rendering/genre_overlay.go` implements basic palette swapping per genre but no shader-based post-processing effects. No bloom, scanlines, chromatic aberration, or film grain.

- **Impact**: All genres have different color palettes but feel visually similar in presentation style. The promised atmospheric differentiation is not achieved.

- **Closing the Gap**:
  1. Create `pkg/rendering/postprocess.go` with shader-like effects
  2. Implement `ApplyVignette(img *ebiten.Image, intensity float64)`
  3. Implement `ApplyScanlines(img *ebiten.Image, density, alpha float64)`
  4. Implement `ApplyFilmGrain(img *ebiten.Image, seed int64, intensity float64)`
  5. Implement `ApplyChromaticAberration(img *ebiten.Image, offset float64)`
  6. Add post-process pass in `Renderer.Draw()` based on current genre

---

## Dynamic Minimap Overlay

- **Stated Goal**: ROADMAP v3.0 specifies an always-visible corner minimap showing explored tiles, icons for landmarks, fading during crisis events, and genre-appropriate styling.

- **Current State**: No minimap implementation exists. Players rely solely on the main map view.

- **Impact**: Players must mentally track explored areas and cannot quickly orient to the destination or nearby points of interest without scrolling the main view.

- **Closing the Gap**:
  1. Create `pkg/ux/minimap.go` with `Minimap` struct tracking explored tiles
  2. Implement `Minimap.Draw(screen *ebiten.Image, worldMap *world.WorldMap, playerPos world.Point)`
  3. Add minimap to HUD in corner position (configurable)
  4. Implement fog overlay for unexplored tiles
  5. Add icons for origin, destination, supply posts, hazards
  6. Implement `SetGenre()` for minimap styling (parchment → hologram → torn atlas)

---

## pkg/ux Test Coverage

- **Stated Goal**: CONTRIBUTING.md states "Testing: Aim for ≥40% coverage per package"

- **Current State**: `pkg/ux` has 15.3% test coverage — the lowest in the project and significantly below the stated target.

- **Impact**: User interface code is largely untested, increasing risk of regressions in player-facing functionality. Menu state transitions, HUD rendering, and loadout screens may have latent bugs.

- **Closing the Gap**:
  1. Create `pkg/ux/hud_test.go` with tests for resource bar rendering
  2. Create `pkg/ux/menus_test.go` with state transition tests (MainMenu → Playing → Paused)
  3. Create `pkg/ux/loadout_test.go` with module selection validation
  4. Use headless rendering stubs for Ebiten-dependent tests
  5. Target: 45% coverage minimum

---

## pkg/game Test Coverage

- **Stated Goal**: CONTRIBUTING.md states "Testing: Aim for ≥40% coverage per package"

- **Current State**: `pkg/game` has 54% coverage, which meets the minimum but is low for the core gameplay package containing win/lose logic and turn advancement.

- **Impact**: Critical game loop logic (`advanceTurn`, `checkConditions`, `consumeResources`) may have edge cases that are untested, risking incorrect win/lose detection or resource depletion bugs.

- **Closing the Gap**:
  1. Add comprehensive tests for `checkConditions()` in `pkg/game/conditions_test.go`:
     - Test each LoseCondition scenario individually
     - Test WinCondition at destination with various party states
  2. Add tests for `advanceTurn()` resource consumption:
     - Verify correct depletion rates by crew size
     - Verify morale penalties when resources depleted
  3. Target: 70% coverage for core gameplay package

---

## Code Duplication in Session Files

- **Stated Goal**: ROADMAP states "No pre-written content" and the codebase shows emphasis on DRY principles with shared base implementations (e.g., `BaseSystem`).

- **Current State**: `pkg/game/session.go` lines 82-135 and `pkg/game/session_headless.go` lines 78-131 contain 54 lines of exact duplication — the largest clone in the codebase.

- **Impact**: Changes to session initialization require editing two files. Risk of drift between display and headless behavior. Maintenance burden is doubled for this critical code path.

- **Closing the Gap**:
  1. Create `pkg/game/session_core.go` with shared initialization logic
  2. Define `sessionCore` struct with common fields and `initCore(cfg SessionConfig) *sessionCore`
  3. Have both `GameSession` (display) and `GameSession` (headless) embed `sessionCore`
  4. Move duplicate NewGameSession logic into `initCore()`
  5. Validate with `go-stats-generator analyze . --format json | jq '.duplication'`

---

## v3.0 Features Marked Incomplete in ROADMAP

The following ROADMAP items are marked with `- [ ]` (incomplete) and represent the gap between stated v3.0 goals and current implementation:

| Feature | ROADMAP Line | Status |
|---------|--------------|--------|
| Dynamic music layer system | 279 | Not implemented |
| Biome-specific music parameters | 280 | Not implemented |
| Music cross-fade transitions | 281 | Not implemented |
| Genre instrument mapping | 282 | Partial (static params only) |
| Distance attenuation for SFX | 285 | Not implemented |
| Stereo panning for spatial audio | 286 | Not implemented |
| Biome-specific ambient loops | 287 | Not implemented |
| Fantasy post-processing | 292 | Not implemented |
| Sci-fi post-processing | 293 | Not implemented |
| Horror post-processing | 294 | Not implemented |
| Cyberpunk post-processing | 295 | Not implemented |
| Post-apocalyptic post-processing | 296 | Not implemented |
| Minimap overlay | 299-302 | Not implemented |

**Note**: All v1.0 and v2.0 features are complete and functional. The gaps above are v3.0 "Visual Polish" features.

---

## Summary

The project successfully delivers on its core v1.0 and v2.0 promises:
- ✅ 100% procedural content generation with no bundled assets
- ✅ Five genre themes with comprehensive `SetGenre()` support
- ✅ Complete gameplay loop with resource management, crew, vessel, and events
- ✅ Deterministic seed-based generation
- ✅ Full save/load system

The primary gaps are:
1. **v3.0 audio polish** — Adaptive multi-layer music and positional audio
2. **v3.0 visual polish** — Post-processing effects and minimap
3. **Test coverage** — `pkg/ux` at 15.3% (target 40%), `pkg/game` could be higher
4. **Code duplication** — Session initialization across display/headless variants

These gaps represent planned future work (v3.0) rather than missing core functionality.
