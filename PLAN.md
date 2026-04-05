# Implementation Plan: v3.0 Visual Polish Milestone

## Project Context
- **What it does**: Voyage is a 100% procedural rogue-like travel simulator where every map, event, crew, vessel, audio, and narrative is generated from a single seed.
- **Current goal**: Complete v3.0 Visual Polish — adaptive audio, post-processing effects, and minimap overlay
- **Estimated Scope**: Medium (12 incomplete features, 2 functions at complexity threshold, 0.65% duplication)

## Goal-Achievement Status

| Stated Goal | Current Status | This Plan Addresses |
|-------------|----------------|---------------------|
| 100% procedural generation | ✅ Achieved | No |
| Five genre themes with GenreSwitcher | ✅ Achieved | No |
| Seed-based determinism | ✅ Achieved | No |
| Core gameplay loop (v1.0) | ✅ Achieved | No |
| Crew depth & vessel upgrades (v2.0) | ✅ Achieved | No |
| Dynamic lighting & particles | ✅ Achieved | No |
| Adaptive multi-layer music | ❌ Not implemented | Yes |
| Positional audio/spatial SFX | ❌ Not implemented | Yes |
| Genre post-processing effects | ❌ Not implemented | Yes |
| Dynamic minimap overlay | ❌ Not implemented | Yes |
| Test coverage ≥40% per package | ⚠️ pkg/ux at 15.3% | Yes |
| Code duplication reduction | ⚠️ 54-line clone in session files | Yes |

## Metrics Summary

- **Complexity hotspots on goal-critical paths**: 2 functions at threshold (cyclomatic = 10)
  - `GenerateCrewRelationshipEvent` in `pkg/events/relationships.go:33`
  - `generateRuinsPattern` in `pkg/rendering/landmark_icon.go:125`
- **Duplication ratio**: 0.65% (176 duplicated lines, largest clone 54 lines)
- **Doc coverage**: 83.9% overall
- **Test coverage gaps**: 
  - `pkg/ux`: 15.3% (target 40%)
  - `pkg/game`: 54.0% (target 70% for core loop)
  - `pkg/rendering`: 58.9%
- **Package coupling**: Clean — 0 circular dependencies

## Implementation Steps

### Step 1: Extract Shared Session Initialization Logic ✅ COMPLETE

- **Deliverable**: New file `pkg/game/session_core.go` containing shared `initSessionCore()` function; refactored `session.go` and `session_headless.go` to embed/call shared logic
- **Dependencies**: None
- **Goal Impact**: Reduces largest code clone (54 lines) and enables safer parallel changes to session variants
- **Acceptance**: Duplication ratio drops below 0.5%; largest clone < 30 lines
- **Status**: ✅ Already implemented in `pkg/game/session_common.go` with `initializeSession()` function. Duplication ratio: 0.45%, largest clone: 27 lines.
- **Validation**: 
  ```bash
  go-stats-generator analyze . --skip-tests --format json --sections duplication 2>/dev/null | jq '.duplication | {ratio: .duplication_ratio, largest: .largest_clone_size}'
  ```

### Step 2: Implement Music State Machine for Adaptive Audio

- **Deliverable**: 
  - Add `MusicState` enum (`Peaceful`, `Tense`, `Combat`, `Victory`, `Death`) in `pkg/audio/music.go`
  - Add `SetMusicState(state MusicState)` method on `MusicGenerator`
  - Modify BPM, melody density, and waveform selection per state
- **Dependencies**: None
- **Goal Impact**: Addresses ROADMAP v3.0 "Dynamic layer system" requirement
- **Acceptance**: `MusicGenerator` responds to state changes with audibly different output; test verifies state-dependent parameter changes
- **Validation**: 
  ```bash
  go test -tags headless -v ./pkg/audio/... -run TestMusicState
  ```

### Step 3: Implement Audio Cross-Fade Transitions

- **Deliverable**: 
  - Add `CrossfadeTo(targetState MusicState, durationMs int)` method in `pkg/audio/music.go`
  - Implement linear interpolation between current and target audio buffers
- **Dependencies**: Step 2 (MusicState enum must exist)
- **Goal Impact**: Addresses ROADMAP v3.0 "Smooth cross-fade between intensity states"
- **Acceptance**: Cross-fade produces gradual transition over specified duration without audio artifacts
- **Validation**: 
  ```bash
  go test -tags headless -v ./pkg/audio/... -run TestCrossfade
  ```

### Step 4: Implement Spatial Audio System

- **Deliverable**: 
  - Create `pkg/audio/spatial.go` with `SpatialAudioConfig` struct
  - Implement `ApplySpatialAudio(samples []float64, config SpatialAudioConfig) (left, right []float64)`
  - Add distance attenuation and stereo panning based on listener/source positions
- **Dependencies**: None
- **Goal Impact**: Addresses ROADMAP v3.0 "Distance attenuation" and "Stereo panning"
- **Acceptance**: Audio samples correctly attenuate with distance; stereo field reflects source position
- **Validation**: 
  ```bash
  go test -tags headless -v ./pkg/audio/... -run TestSpatialAudio
  ```

### Step 5: Wire Music State to Game Events

- **Deliverable**: 
  - Add `SetMusicState()` calls in `pkg/game/session.go` during:
    - Normal travel → `Peaceful`
    - Encounter start → `Combat` or `Tense`
    - Win condition → `Victory`
    - Lose condition → `Death`
- **Dependencies**: Steps 2, 3 (MusicState and cross-fade must exist)
- **Goal Impact**: Connects adaptive audio to actual gameplay states for immersive feedback
- **Acceptance**: Music transitions during gameplay events are audible and contextually appropriate
- **Validation**: 
  ```bash
  go test -tags headless -v ./pkg/game/... -run TestMusicStateTransitions
  ```

### Step 6: Create Post-Processing Effects Foundation

- **Deliverable**: 
  - Create `pkg/rendering/postprocess.go` with `PostProcessor` struct
  - Implement `ApplyVignette(img *ebiten.Image, intensity float64) *ebiten.Image`
  - Use Ebitengine's offscreen render target approach for effect chaining
- **Dependencies**: None
- **Goal Impact**: Foundation for all genre-specific visual effects (ROADMAP v3.0)
- **Acceptance**: Vignette effect visibly darkens screen edges at intensity > 0.5
- **Validation**: 
  ```bash
  go test -tags headless -v ./pkg/rendering/... -run TestVignette
  ```

### Step 7: Implement Genre-Specific Post-Processing Effects

- **Deliverable**: 
  - Add `ApplyScanlines(img, density, alpha float64)` for sci-fi
  - Add `ApplyFilmGrain(img, seed int64, intensity float64)` for horror
  - Add `ApplyChromaticAberration(img, offset float64)` for cyberpunk
  - Add `ApplySepia(img, intensity float64)` for post-apocalyptic
  - Wire `SetGenre()` to apply appropriate effect chain
- **Dependencies**: Step 6 (PostProcessor foundation)
- **Goal Impact**: Addresses all 5 genre post-processing requirements from ROADMAP v3.0
- **Acceptance**: Each genre displays distinct visual style; effects are procedural (no bundled textures)
- **Validation**: 
  ```bash
  go test -tags headless -v ./pkg/rendering/... -run TestGenrePostProcessing
  ```

### Step 8: Implement Dynamic Minimap Component

- **Deliverable**: 
  - Create `pkg/ux/minimap.go` with `Minimap` struct
  - Implement `Draw(screen *ebiten.Image, worldMap, playerPos)` rendering explored tiles
  - Add fog overlay for unexplored areas
  - Add icons for origin, destination, supply posts, hazards
- **Dependencies**: None
- **Goal Impact**: Addresses ROADMAP v3.0 "Always-visible minimap" requirement
- **Acceptance**: Minimap renders in corner; explored tiles visible; player position accurate
- **Validation**: 
  ```bash
  go test -tags headless -v ./pkg/ux/... -run TestMinimap
  ```

### Step 9: Add Minimap Genre Theming and HUD Integration

- **Deliverable**: 
  - Add `SetGenre()` on Minimap for aesthetic variants (parchment, hologram, torn atlas, AR overlay, scratched road atlas)
  - Integrate minimap into HUD draw loop
  - Add minimap fade during crisis events
- **Dependencies**: Step 8 (Minimap component)
- **Goal Impact**: Completes minimap feature with genre theming per ROADMAP v3.0
- **Acceptance**: Minimap appearance changes with genre; dims appropriately during encounters
- **Validation**: 
  ```bash
  go test -tags headless -v ./pkg/ux/... -run TestMinimapGenre
  ```

### Step 10: Increase pkg/ux Test Coverage to 40%

- **Deliverable**: 
  - Create `pkg/ux/hud_test.go` with tests for `DrawResourceBar`, `DrawCrewPanel`
  - Create `pkg/ux/menus_test.go` with state transition tests (MainMenu → Playing → Paused)
  - Add table-driven tests for menu option navigation
- **Dependencies**: Steps 8, 9 (Minimap tests contribute to coverage)
- **Goal Impact**: Addresses CONTRIBUTING.md test coverage requirement (≥40%)
- **Acceptance**: `go test -cover ./pkg/ux/...` reports ≥40% coverage
- **Validation**: 
  ```bash
  go test -tags headless -cover ./pkg/ux/... | grep coverage
  ```

### Step 11: Increase pkg/game Test Coverage to 70%

- **Deliverable**: 
  - Add `pkg/game/conditions_test.go` with tests for each `LoseCondition` and `WinCondition`
  - Add `pkg/game/turn_test.go` with tests for `advanceTurn()` resource consumption
  - Test edge cases: zero resources, full crew death, morale depletion
- **Dependencies**: Step 1 (refactored session core simplifies test setup)
- **Goal Impact**: Ensures core gameplay logic is thoroughly tested for correctness
- **Acceptance**: `go test -cover ./pkg/game/...` reports ≥70% coverage
- **Validation**: 
  ```bash
  go test -tags headless -cover ./pkg/game/... | grep coverage
  ```

### Step 12: Add Biome-Specific Ambient Audio Loops

- **Deliverable**: 
  - Add `AmbientLoopType` enum in `pkg/audio/ambient.go`
  - Implement procedural ambient generation (wind, space hum, groaning metal, city noise)
  - Add `SetBiome(biomeType)` for ambient audio selection
  - Wire biome changes to ambient audio in world traversal
- **Dependencies**: Step 4 (spatial audio system for positioning ambient sources)
- **Goal Impact**: Addresses ROADMAP v3.0 "Ambient loop per biome/region"
- **Acceptance**: Different biomes produce distinct ambient soundscapes
- **Validation**: 
  ```bash
  go test -tags headless -v ./pkg/audio/... -run TestAmbientBiome
  ```

## Dependency Graph

```
Step 1 (Session Core)
    └── Step 11 (Game Tests)

Step 2 (Music State)
    └── Step 3 (Cross-fade)
        └── Step 5 (Wire to Game)

Step 4 (Spatial Audio)
    └── Step 12 (Ambient Loops)

Step 6 (PostProcess Foundation)
    └── Step 7 (Genre Effects)

Step 8 (Minimap)
    └── Step 9 (Minimap Genre)
        └── Step 10 (UX Tests)
```

## Recommended Execution Order

1. **Step 1** — Session core extraction (unblocks Step 11, reduces maintenance burden)
2. **Steps 2, 4, 6, 8** — Independent foundations (can be parallelized)
3. **Steps 3, 7, 9** — Feature completions (depend on foundations)
4. **Steps 5, 12** — Game integration (wire features to gameplay)
5. **Steps 10, 11** — Test coverage (run last to include new code)

## Notes

- All audio features must use procedural synthesis — no bundled samples per project constraint
- Post-processing effects use Ebitengine's Kage shader language for cross-platform compatibility
- Test commands use `-tags headless` for CI/server environments without display
- Complexity of `generateRuinsPattern` (cyclomatic 10) is at threshold but not blocking; consider refactoring after v3.0 if pattern grows
