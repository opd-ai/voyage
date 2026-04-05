# Implementation Gaps — 2026-04-05

This document identifies gaps between Voyage's stated goals and its current implementation.

---

## 1. Game Loop Integration

- **Stated Goal**: The README states "Core engine complete" with a playable journey from origin to destination. The ROADMAP marks all v1.0 items as complete, implying a functional game.

- **Current State**: The `cmd/voyage/main.go` entry point contains:
  ```go
  // TODO: Initialize game systems
  // - World map generation
  // - Rendering system
  // - Resource management
  // - Crew system
  // - Vessel system
  // - Event system
  // - Audio synthesis
  // - UI/HUD
  ```
  The `demo()` function only creates one entity and prints RNG values. No call to `ebiten.RunGame()` occurs. The game cannot be played.

- **Impact**: Users who install the game via `go install github.com/opd-ai/voyage/cmd/voyage@latest` will get a binary that prints text and exits. The 20+ packages of implemented functionality are unreachable from the main entry point.

- **Closing the Gap**:
  1. Complete system initialization in `main()` by instantiating each subsystem
  2. Create a `GameSession` orchestrator that connects world, crew, vessel, resources, and events
  3. Replace `demo()` with `g.Run()` using the existing `pkg/game.Game` type
  4. Add integration tests that verify the full game loop starts

---

## 2. Animated Sprite Generation (v3.0)

- **Stated Goal**: ROADMAP v3.0 lists:
  - Animated overworld tiles (flowing water, wind-swept grass, flickering fires)
  - Crew member portrait animation (idle breathing, hurt flinch, death fade)
  - Vessel damage states (pristine → worn → damaged → critical sprites)
  - Animated landmark icons (smoking ruins, blinking outpost lights)

- **Current State**: All v3.0 "Enhanced Sprite Generation" items are marked with `- [ ]` (unchecked) in the ROADMAP. The `pkg/rendering/tilegen.go` generates static procedural tiles using cellular automata but has no animation frame logic.

- **Impact**: The visual presentation lacks the "feel alive" quality promised by the ROADMAP. Static sprites reduce immersion compared to reference games like FTL.

- **Closing the Gap**:
  1. Add animation frame support to `TileGenerator`:
     ```go
     type AnimatedTile struct {
         Frames    []*ebiten.Image
         FrameTime float64
         Loop      bool
     }
     ```
  2. Implement animation presets per tile type (water cycles through 4 frames, fire flickers randomly)
  3. Add `Update()` method to advance frame counters
  4. Apply similar pattern to crew portraits and vessel sprites

---

## 3. Adaptive Multi-Layer Music

- **Stated Goal**: ROADMAP v3.0 specifies:
  - Dynamic layer system (peaceful travel, crisis, encounter, victory, death)
  - Biome-specific ambient music parameters
  - Smooth cross-fade between intensity states

- **Current State**: The `pkg/audio/music.go` generates a single looping ambient track per genre. There is no layer system, intensity detection, or cross-fade logic. The `GenerateLoop()` function produces a fixed arrangement.

- **Impact**: Music does not respond to gameplay state. A crisis event sounds the same as peaceful travel, reducing emotional feedback to the player.

- **Closing the Gap**:
  1. Define `MusicState` enum: `StateTravel`, `StateCrisis`, `StateEncounter`, `StateVictory`, `StateDeath`
  2. Generate separate tracks for each state, all quantized to the same BPM/bar length
  3. Add a `MusicMixer` that cross-fades between layers based on game state
  4. Hook `MusicMixer.SetState()` calls into event resolution and win/lose conditions

---

## 4. Positional Audio

- **Stated Goal**: ROADMAP v3.0 lists:
  - Distance attenuation for offscreen events
  - Left/right stereo panning for spatial awareness
  - Ambient loop per biome/region

- **Current State**: The `pkg/audio/sfx.go` generates mono sound effects with no position data. The `SFXGenerator.Generate()` method returns `[]float64` without stereo separation.

- **Impact**: Players cannot use audio cues to locate offscreen events, reducing tactical awareness.

- **Closing the Gap**:
  1. Add stereo output to SFX generation:
     ```go
     type StereoSample struct {
         Left, Right float64
     }
     func (g *SFXGenerator) GenerateStereo(sfxType SFXType, pan float64) []StereoSample
     ```
  2. Calculate pan from event position relative to vessel position
  3. Apply distance-based amplitude attenuation
  4. Add per-biome ambient loop selection in `MusicGenerator`

---

## 5. Genre Post-Processing Visual Effects

- **Stated Goal**: ROADMAP v3.0 specifies 5 genre-specific post-processing presets:
  - Fantasy: warm desaturated vignette, bloom on magic effects
  - Sci-fi: cool scanline overlay, chromatic aberration
  - Horror: desaturate + red-tint, film grain
  - Cyberpunk: neon bloom, CRT curvature, glitch artifacts
  - Post-apocalyptic: sepia wash, dust overlay, heavy vignette

- **Current State**: The `pkg/rendering/lighting.go` provides genre-specific color palettes and light colors, but there are no shader-based post-processing effects. The `Draw()` method applies tinting but no bloom, scanlines, vignette, or film grain.

- **Impact**: Each genre looks different in color but not in visual style. The distinct "feel" of each genre is diminished.

- **Closing the Gap**:
  1. Add Ebitengine shader support for post-processing
  2. Create a `PostProcessor` interface with `Apply(*ebiten.Image) *ebiten.Image`
  3. Implement genre-specific shaders (Kage shader language):
     ```go
     var horrorGrainShader = []byte(`
     package main
     func Fragment(...) vec4 {
         // Film grain + desaturation
     }
     `)
     ```
  4. Hook `PostProcessor.SetGenre()` into the render pipeline

---

## 6. Dynamic Minimap Overlay

- **Stated Goal**: ROADMAP v3.0 lists:
  - Always-visible corner minimap showing explored tiles
  - Icons for towns, ruins, hazards, destination
  - Minimap fades when navigation module is damaged
  - Genre-appropriate aesthetic

- **Current State**: The `pkg/ux/worldmap.go` renders the full map, not a minimap overlay. There is no corner widget or visibility-based fading.

- **Impact**: Players must switch to the full map view to check their position, interrupting gameplay flow.

- **Closing the Gap**:
  1. Add `MinimapWidget` to `pkg/ux/`:
     ```go
     type MinimapWidget struct {
         X, Y, Size   int
         ExploredMask [][]bool
         Opacity      float64
     }
     ```
  2. Render minimap in the corner during `StatePlaying`
  3. Apply genre-specific styling via `SetGenre()`
  4. Link opacity to vessel navigation module integrity

---

## 7. Crew Relationship Events

- **Stated Goal**: ROADMAP v2.0 marks as complete:
  - Crew relationship network (pairs that bicker or bond affect morale events)
  - Crew-specific events (personal crisis, milestone, sacrifice opportunity)

- **Current State**: The `pkg/crew/relationship.go` defines `RelationshipNetwork` with `AddBond()` and `AddRivalry()` methods. However, no code in `pkg/events/generator.go` queries the relationship network to create relationship-based events.

- **Impact**: The relationship system exists but has no gameplay effect. Crew members remain functionally independent despite the data structure supporting relationships.

- **Closing the Gap**:
  1. Add `GetStrongBonds()` and `GetRivalries()` methods to `RelationshipNetwork`
  2. In `EventGenerator`, periodically check for high-relationship pairs
  3. Generate events like "X and Y have a conflict" (rivalry) or "X saves Y" (bond)
  4. Add relationship delta outcomes to event choices

---

## 8. Packages Without Tests

- **Stated Goal**: CONTRIBUTING.md states "Aim for ≥40% coverage per package" and "Write or update tests as needed."

- **Current State**: Three packages have no test files:
  - `github.com/opd-ai/voyage/cmd/voyage`
  - `github.com/opd-ai/voyage/pkg/procgen/event`
  - `github.com/opd-ai/voyage/pkg/world`

- **Impact**: Changes to these packages cannot be validated automatically. The `pkg/world` package is particularly critical as it handles game state coordination.

- **Closing the Gap**:
  1. Add `cmd/voyage/main_test.go` with flag parsing tests
  2. Add `pkg/procgen/event/event_test.go` with generation tests
  3. Add `pkg/world/world_test.go` with state management tests
  4. Run `go test -cover ./...` and verify ≥40% per package

---

## Gap Priority Matrix

| Gap | Severity | Effort | Priority |
|-----|----------|--------|----------|
| Game Loop Integration | CRITICAL | Medium | P0 |
| Packages Without Tests | HIGH | Low | P1 |
| Crew Relationship Events | HIGH | Medium | P2 |
| Animated Sprite Generation | MEDIUM | High | P3 |
| Adaptive Multi-Layer Music | MEDIUM | High | P4 |
| Dynamic Minimap | MEDIUM | Medium | P5 |
| Positional Audio | LOW | Medium | P6 |
| Genre Post-Processing | LOW | High | P7 |

---

## Summary

The Voyage codebase contains **22 well-structured packages** implementing all claimed v1.0 and v2.0 features at the subsystem level. However, the **critical integration gap** in `cmd/voyage/main.go` means these subsystems are never connected into a playable game.

The v3.0 visual polish features (animated sprites, adaptive music, post-processing, minimap) are correctly marked as incomplete in the ROADMAP and should not be considered blockers.

**Recommendation**: Address Gap #1 (Game Loop Integration) as the highest priority. All other systems are ready to be wired together.

---

*Generated by functional audit on 2026-04-05*
