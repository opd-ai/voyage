# Ebitengine Game Audit Report

## Executive Summary
- **Total Issues**: 52 (Critical: 6, High: 18, Medium: 18, Low: 10)
- **Performance Issues**: 14
- **Code Quality**: 8

## Critical Issues

### [C-001] ebiten.NewImage() Called Per Tile Per Frame in Renderer
- **Location**: `pkg/rendering/renderer.go:82-88`
- **Category**: Performance / Ebitengine-Specific
- **Description**: `DrawTile()` creates a new `ebiten.Image` for every tile, every frame. For a 15×15 viewport at 60 TPS, this produces ~13,500 image allocations per second, each requiring GPU-backed memory. These images are never disposed of explicitly and rely on GC finalization.
- **Impact**: Severe frame rate degradation, memory exhaustion, and eventual OOM crash during extended play sessions. This is the single largest performance bottleneck in the codebase.
- **Reproduction**: Run the game and observe memory growth via `runtime.ReadMemStats()`. Memory will increase continuously as the world map is rendered.
- **Suggested Fix**: Pre-allocate a tile cache (`map[int]*ebiten.Image`) keyed by tile type. Create each tile image once at initialization and reuse it via `DrawImage()` with translation options.

### [C-002] No Input Debouncing on Event Choice Keys
- **Location**: `pkg/game/session.go:113-120`
- **Category**: Input / Logic
- **Description**: Event choice keys (1-4) use `ebiten.IsKeyPressed()` which returns true every frame the key is held. If a player holds a number key, the same event choice fires every tick (60 times/sec), potentially resolving the same event repeatedly or consuming resources multiple times from a single keypress.
- **Impact**: Players can accidentally trigger event choices multiple times. Resources are consumed repeatedly. Events resolve with unintended outcomes.
- **Reproduction**: Open the game, trigger an event, and hold the "1" key. Observe that the choice resolves immediately and may fire multiple times if not gated by state.
- **Suggested Fix**: Replace `ebiten.IsKeyPressed()` with `inpututil.IsKeyJustPressed()` from `github.com/hajimehoshi/ebiten/v2/inpututil`, or implement a key-release detection pattern (already used in `handleDebugToggle` elsewhere in the same file).

### [C-003] Double Turn Advancement in Single Frame
- **Location**: `pkg/game/session.go:76-92, 115-120, 155-169`
- **Category**: Input / State
- **Description**: Movement input and event choice input are both processed in the same `Update()` frame without mutual exclusion. If a player presses a movement key and an event choice key simultaneously, `advanceTurn()` is called once from movement (line ~95) and potentially again from event resolution, causing two turns to elapse in one frame.
- **Impact**: Game timeline corrupted — resources consumed at double rate, events advance unexpectedly, and the turn counter skips values.
- **Reproduction**: While an event is active, press a movement key and an event choice key in the same frame. Observe that the turn counter increments by 2.
- **Suggested Fix**: Add an early return or flag after movement processing to prevent event handling in the same frame. Use a per-frame action budget (one action per `Update()` call).

### [C-004] Memory Leak in Pause Overlay
- **Location**: `pkg/game/game.go:250-252`
- **Category**: Performance / Ebitengine-Specific
- **Description**: `drawPauseOverlay()` allocates a full-screen `ebiten.NewImage(g.width, g.height)` every frame while paused. At 60 FPS, this creates 60 screen-sized GPU-backed images per second, none of which are explicitly disposed.
- **Impact**: Pausing the game for even 30 seconds produces ~1,800 orphaned images. Extended pausing causes OOM or severe GC pressure.
- **Reproduction**: Pause the game and monitor memory usage. It will climb continuously.
- **Suggested Fix**: Cache the overlay image as a field on the `Game` struct, create it once, and reuse it. Alternatively, use `screen.Fill()` with a semi-transparent color or apply a color matrix via `DrawImageOptions.ColorScale`.

### [C-005] Entity Recycling Without Component Cleanup
- **Location**: `pkg/engine/world.go:103-110, 161-171`
- **Category**: State
- **Description**: `DespawnImmediate()` adds entities to a pool for reuse but does not clear their component data. When a pooled entity is re-spawned via `Spawn()`, it retains stale components from its previous life. Systems holding references to despawned entities will access recycled data with undefined state.
- **Impact**: Old component data (position, health, damage, etc.) bleeds into newly spawned entities. Systems may process stale entity references, causing logic errors, visual glitches, and state corruption.
- **Reproduction**: Spawn an entity with components, despawn it, then spawn a new entity. Inspect the new entity's components — they will contain data from the previous entity.
- **Suggested Fix**: Clear all component maps for the entity in `DespawnImmediate()` before pooling. Alternatively, use entity generation counters to invalidate stale references.

### [C-006] Save Format Has No Migration Path
- **Location**: `pkg/saveload/save.go:145-155`
- **Category**: State / Logic
- **Description**: `Validate()` checks `Version < 1 || Version > CurrentVersion` but there is no migration logic. When `CurrentVersion` increments in a future update, all existing save files with older versions will fail validation with `ErrInvalidVersion` and become permanently unloadable.
- **Impact**: Players lose all save progress on any game update that changes the save format. No graceful degradation or migration.
- **Reproduction**: Increment `CurrentVersion` and attempt to load a save created with the previous version. The load will fail.
- **Suggested Fix**: Implement versioned migration handlers (e.g., `migrateV1toV2()`) and apply them sequentially during load.

## High Priority Issues

### [H-001] Fog Overlay Creates ebiten.NewImage Per Tile Per Frame
- **Location**: `pkg/ux/worldmap.go:120-127`
- **Category**: Performance / Rendering
- **Description**: `drawFogOverlay()` calls `ebiten.NewImage(tileSize, tileSize)` for every unexplored tile visible in the viewport, every frame. With a 15×15 viewport, this can be up to 225 allocations per frame.
- **Impact**: Significant memory pressure and frame drops when fog-of-war is extensive. Combined with C-001, worldmap rendering becomes the dominant performance bottleneck.
- **Suggested Fix**: Pre-allocate a single fog tile image and reuse it for all fog positions via translated `DrawImageOptions`.

### [H-002] Destination and Vessel Markers Create ebiten.NewImage Every Frame
- **Location**: `pkg/ux/worldmap.go:131, 160`
- **Category**: Performance / Rendering
- **Description**: Both `drawDestinationMarker()` and `drawVesselMarker()` create new `ebiten.Image` objects every frame.
- **Impact**: Two additional GPU image allocations per frame, contributing to the cumulative memory leak in the rendering pipeline.
- **Suggested Fix**: Cache marker images as struct fields; create once, reuse on every draw.

### [H-003] Panel Overlay Allocates Full-Screen Image Every Frame
- **Location**: `pkg/ux/panel.go:31, 44-46`
- **Category**: Performance / UI
- **Description**: `DrawOverlay()` creates `ebiten.NewImage(width, height)` (screen-sized) every call. `DrawCenteredPanel()` similarly allocates per frame. These are called whenever any menu, dialog, or popup is visible.
- **Impact**: Every visible panel causes a screen-sized GPU allocation per frame. Menus and dialogs become memory-intensive.
- **Suggested Fix**: Accept a pre-allocated image as parameter, or cache the overlay image.

### [H-004] Event Overlay Creates ebiten.NewImage Every Frame
- **Location**: `pkg/ux/events.go:96`
- **Category**: Performance / UI
- **Description**: `EventOverlay.Draw()` creates a new `ebiten.NewImage(overlayWidth, overlayHeight)` every frame while an event is displayed.
- **Impact**: Events displayed for multiple seconds cause hundreds of orphaned images. Combined with event choice issues (C-002), this compounds into severe performance degradation during event encounters.
- **Suggested Fix**: Cache the overlay image and only recreate when overlay dimensions change.

### [H-005] Pixel-by-Pixel Particle Rendering
- **Location**: `pkg/rendering/particles.go:346-361`
- **Category**: Performance / Rendering
- **Description**: `drawParticle()` uses nested loops with `screen.Set(x, y, color)` for every particle pixel. With up to 2,000 active particles, each potentially several pixels in size, this performs thousands of individual pixel writes per frame.
- **Impact**: Particle effects cause severe frame drops. Burst effects (explosions, weather) become unplayable on lower-end hardware.
- **Suggested Fix**: Use pre-rendered particle sprite images with `DrawImage()` and color scaling via `DrawImageOptions`. Batch particles of the same type.

### [H-006] Lighting Overlay Created Every Frame with Pixel-by-Pixel Drawing
- **Location**: `pkg/rendering/lighting.go:367-387, 427-435`
- **Category**: Performance / Rendering
- **Description**: `CreatePointLightOverlay()` and `CreateLightingOverlay()` allocate new images and perform per-pixel radial light calculations via `drawRadialLight()` (lines 403-418) every frame.
- **Impact**: Lighting system is the second largest performance bottleneck after tile rendering. Per-pixel operations on potentially large images cause consistent frame drops.
- **Suggested Fix**: Cache lighting overlays and only recalculate when light sources or player position change. Use shader-based lighting if available.

### [H-007] Landmark Icons Create ebiten.NewImage in Draw Functions
- **Location**: `pkg/rendering/landmark_icon.go:115, 221, 325, 390, 470, 523`
- **Category**: Performance / Rendering
- **Description**: Each landmark icon type (mountain, ruins, tower, etc.) creates a new `ebiten.NewImage()` every time it's drawn. Six separate draw functions all exhibit this pattern.
- **Impact**: Multiple landmark icons visible simultaneously multiply the allocation rate. Worldmap areas with many landmarks become performance hotspots.
- **Suggested Fix**: Generate landmark icons once at initialization and cache them by type/variant.

### [H-008] Genre Overlay Per-Pixel Image Processing
- **Location**: `pkg/rendering/genre_overlay.go:33-49, 111-153`
- **Category**: Performance / Rendering
- **Description**: `ApplyOverlay()` creates a new result image and iterates every pixel to apply saturation/brightness adjustments. This is called during genre transitions on potentially large images.
- **Impact**: Genre changes cause visible stutters as the overlay is recomputed. On large screens, the per-pixel loop becomes a bottleneck.
- **Suggested Fix**: Use Ebitengine's `ColorScale` or shader-based color adjustments instead of manual per-pixel processing.

### [H-009] Portrait Animation Creates Temporary Images Per Frame
- **Location**: `pkg/rendering/portrait.go:235-241, 266, 301`
- **Category**: Performance / Rendering
- **Description**: Breathing animation and death animation effects create temporary `ebiten.Image` objects per frame for visual effects.
- **Impact**: Crew portraits with active animations contribute to the cumulative memory leak. Multiple crew members animating simultaneously multiply the effect.
- **Suggested Fix**: Pre-allocate animation scratch images and reuse them.

### [H-010] Hardcoded Delta Time Ignores Actual TPS
- **Location**: `pkg/game/game.go:127`
- **Category**: Logic / Ebitengine-Specific
- **Description**: `g.world.Update(1.0 / 60.0)` passes a hardcoded delta time to the ECS world update. Ebitengine defaults to 60 TPS but this can change via `ebiten.SetTPS()` or vary under load.
- **Impact**: If TPS is changed or the game runs under heavy load with reduced tick rate, all physics and animation timing becomes incorrect — movements appear faster or slower than intended.
- **Suggested Fix**: Use `1.0 / float64(ebiten.TPS())` or `ebiten.ActualTPS()` for adaptive delta time.

### [H-011] Crew Relationship pairKey Truncates Large IDs
- **Location**: `pkg/crew/relationship.go:64-69`
- **Category**: Logic / State
- **Description**: `pairKey(a, b int)` converts integer crew IDs to runes via `string(rune(a))`. Rune values above Unicode max (1,114,111) are replaced with the Unicode replacement character (U+FFFD), causing distinct crew IDs to map to the same key. Even valid large IDs may collide if they share the same Unicode character representation.
- **Impact**: Relationship data corruption — different crew pairs silently overwrite each other's relationship strength. Morale modifiers become incorrect.
- **Suggested Fix**: Use `fmt.Sprintf("%d-%d", a, b)` or `strconv.Itoa(a) + "-" + strconv.Itoa(b)`.

### [H-012] Integer Overflow Risk in Cargo Weight Calculations
- **Location**: `pkg/vessel/cargo.go:182-187`
- **Category**: Logic
- **Description**: `CanAdd()` computes `totalWeight := weight * quantity` and `totalVolume := volume * quantity` using `int` multiplication with no overflow check. If both operands are large (e.g., weight=100000, quantity=100000), the result silently wraps to a negative value, causing the capacity check to incorrectly pass.
- **Impact**: Cargo hold can be overfilled beyond its limits. Potential exploit for carrying unlimited resources.
- **Suggested Fix**: Use `int64` for intermediate calculations, or validate that `weight * quantity` doesn't exceed `math.MaxInt` before computing.

### [H-013] Division by Zero in Audio SFX Generation
- **Location**: `pkg/audio/sfx.go:191`
- **Category**: Logic
- **Description**: `noteLen = samples / len(notes)` — if `samples < len(notes)`, `noteLen` becomes 0. Subsequent `noteIndex := i / noteLen` causes a division-by-zero panic.
- **Impact**: Game crash when generating success SFX with very short duration parameters.
- **Suggested Fix**: Add a guard: `if noteLen == 0 { noteLen = 1 }` or validate `samples >= len(notes)` before division.

### [H-014] Division by Zero in Encounter SFX
- **Location**: `pkg/encounters/sfx.go:104, 189`
- **Category**: Logic
- **Description**: Same pattern as H-013 — `noteLen` computed from integer division can be zero, causing panic in subsequent index calculations.
- **Impact**: Crash during encounter sound effect generation with edge-case parameters.
- **Suggested Fix**: Guard against zero `noteLen` before division.

### [H-015] Division by Zero in Encounter Resolution
- **Location**: `pkg/encounters/resolution.go:108, 134`
- **Category**: Logic
- **Description**: `progressRatio := enc.TotalProgress / float64(enc.MaxPhases)` and `total / float64(count)` have no zero-denomintor guards. If an encounter has 0 phases or 0 participants, these produce `+Inf` or `NaN`.
- **Impact**: NaN propagation corrupts encounter state and produces undefined behavior in subsequent calculations.
- **Suggested Fix**: Check denominators before division; return 0 or a sensible default for edge cases.

### [H-016] Music Crossfade Panic on Empty Sample Slice
- **Location**: `pkg/audio/music.go:159, 164-165`
- **Category**: Logic
- **Description**: `CrossfadeTo()` computes `fromIdx := i % len(fromSamples)`. If the source music generation failed and returned an empty slice, `len(fromSamples) == 0` causes a modulo-by-zero panic.
- **Impact**: Game crash during music transitions when audio generation fails silently.
- **Suggested Fix**: Check `len(fromSamples) > 0` before the crossfade loop; skip crossfade if source is empty.

### [H-017] Council Vote Division by Zero (Potential)
- **Location**: `pkg/council/council.go:100-112`
- **Category**: Logic
- **Description**: `calculateDissentPenalty()` divides by `totalVotes`. While line 101-103 has a zero-check that returns early, the function `resolveVote()` at line 80 does not check for zero total votes before determining the winning option, which could lead to a meaningless vote result.
- **Impact**: With zero total votes, the council decision defaults to `OptionSafe` regardless of context. While not a crash, it's a logic gap that could produce unexpected gameplay outcomes.
- **Suggested Fix**: Add a "no quorum" check in `resolveVote()` — if total votes is zero, return an abstention result or defer the decision.

### [H-018] Event Relationship Text Unsafe String Manipulation
- **Location**: `pkg/events/relationships.go:127-133`
- **Category**: Logic
- **Description**: `formatRelationshipText()` performs index-based string modification that can corrupt multi-byte UTF-8 sequences. The function modifies strings during iteration using byte-level indexing, which is unsafe for Unicode text.
- **Impact**: Name substitution produces garbled text for non-ASCII character names. Display corruption in relationship event descriptions.
- **Suggested Fix**: Use `strings.Replace()` or `strings.ReplaceAll()` instead of manual byte manipulation.

## Medium Priority Issues

### [M-001] ESC Key Pause Toggle Not Debounced
- **Location**: `pkg/game/session.go:237-239`
- **Category**: Input / State
- **Description**: Holding ESC causes the game state to rapidly alternate between `StatePaused` and `StatePlaying` every frame, as `IsKeyPressed()` is used without release detection.
- **Impact**: Pause screen flickers rapidly while ESC is held, making it unusable. Combined with C-004, each pause frame allocates a new overlay image.
- **Suggested Fix**: Use `inpututil.IsKeyJustPressed(ebiten.KeyEscape)`.

### [M-002] Foraged Tiles Map Grows Unbounded
- **Location**: `pkg/game/forage.go:129-130`
- **Category**: State / Performance
- **Description**: `fm.foragedTiles` map accumulates entries as the player forages but is never cleared. `ResetTile()` exists but is never called. Over a long game, this map grows proportional to the number of tiles visited.
- **Impact**: Memory leak proportional to game duration. For games with extensive exploration, this becomes significant.
- **Suggested Fix**: Periodically clear old entries (e.g., tiles foraged more than N turns ago) or call `ResetTile()` when tiles should regenerate.

### [M-003] Resource Consumption Return Values Ignored
- **Location**: `pkg/game/movement.go:82-83`
- **Category**: Logic
- **Description**: `res.Consume()` return values are discarded. If consumption fails (insufficient resources), the movement still proceeds as if resources were consumed, potentially allowing resources to go negative.
- **Impact**: Resources can desynchronize from actual consumption, allowing players to move without sufficient fuel/supplies.
- **Suggested Fix**: Check the return value and prevent movement if consumption fails.

### [M-004] Starvation Requires Both Food AND Water Depleted
- **Location**: `pkg/game/conditions.go:84-87`
- **Category**: Logic
- **Description**: The starvation loss condition checks `food <= 0 && water <= 0` (both depleted). Realistically, depletion of either resource alone should trigger crew health effects.
- **Impact**: Gameplay balance issue — crew can survive indefinitely on only water or only food.
- **Suggested Fix**: Change to `food <= 0 || water <= 0` for the starvation condition, or add separate conditions for each resource.

### [M-005] TimeOfDay Off-By-One in IsNight()
- **Location**: `pkg/game/time.go:85-86`
- **Category**: Logic
- **Description**: `IsNight()` returns true when `TimeOfDay() >= dayLength - 1`. With `dayLength = 4`, night triggers at time index 3. The `>=` operator means this is correct for "last turn of each day," but the comment and function name imply a longer nighttime period. If the intent is a single-turn night, `==` would be clearer.
- **Impact**: Night events may trigger at unexpected times if the intent was a range of night turns rather than a single turn.
- **Suggested Fix**: Clarify intent with documentation or change to `==` if only the last turn is night.

### [M-006] TurnsUntilNight Calculation Confusing
- **Location**: `pkg/game/time.go:107-113`
- **Category**: Logic
- **Description**: The `TurnsUntilNight()` calculation returns `dayLength - 1` when `remaining <= 0`, which represents the maximum time-of-day value, not the number of turns remaining. The semantics are unclear.
- **Impact**: Incorrect scheduling of night-dependent events and UI displays.
- **Suggested Fix**: Rewrite with clear variable names and add unit tests for edge cases.

### [M-007] Double Healing in CampRest
- **Location**: `pkg/game/rest.go:124-130`
- **Category**: Logic
- **Description**: `CampRest()` calls `Rest()` which heals crew members, then iterates the party again to apply `extraHeal`. Members are healed twice — once in `Rest()` and once in the `CampRest()` loop.
- **Impact**: Camp rest heals significantly more than intended, breaking game balance. Players who camp frequently gain a disproportionate advantage.
- **Suggested Fix**: Either remove the second healing loop in `CampRest()`, or have `Rest()` not heal when called from `CampRest()`.

### [M-008] Cargo Weight Tracking Mismatch on Stacking
- **Location**: `pkg/vessel/cargo.go:189-228`
- **Category**: State
- **Description**: When stacking cargo, `AddWithVolume()` matches items by name and category. If two items have the same name but different per-unit weights (e.g., from salvage with variance), the old item's weight is used for removal calculations but the new item's weight was used for addition. This creates phantom weight discrepancies.
- **Impact**: Inventory weight/volume desynchronizes from actual contents over time. Cargo hold may report incorrect capacity.
- **Suggested Fix**: Only stack items with matching weight/volume, or normalize weight on stack.

### [M-009] Status Effect Pointer Invalidation
- **Location**: `pkg/crew/status.go:191-198`
- **Category**: State
- **Description**: `GetEffect()` returns a pointer to a slice element. If the slice is later reallocated by `append()` in `AddEffect()`, the returned pointer becomes dangling. Callers modifying the effect through this pointer will write to freed memory.
- **Impact**: Status effect modifications silently lost; potential memory corruption in edge cases.
- **Suggested Fix**: Return a copy of the effect, or use indices instead of pointers.

### [M-010] Unchecked Map Access in GetResourceName
- **Location**: `pkg/resources/resources.go:159-166`
- **Category**: Logic
- **Description**: `GetResourceName()` accesses `names[rt]` without checking if the `ResourceType` key exists. If a new resource type is added but not all genre maps are updated, this returns an empty string.
- **Impact**: Missing resource names in UI; blank labels confuse players.
- **Suggested Fix**: Check for key existence and return a fallback name like the resource type's string representation.

### [M-011] Infinite Hue Wrapping Loop
- **Location**: `pkg/vessel/visual.go:157-166`
- **Category**: Performance / Logic
- **Description**: `wrapHue()` uses iterative addition/subtraction loops instead of modulo arithmetic. For extremely large or small input values, the loops iterate millions of times.
- **Impact**: Performance stall if hue values become very large due to accumulated floating-point drift.
- **Suggested Fix**: Use `math.Mod(h, 360)` with proper negative handling.

### [M-012] Particle MaxLife Division by Zero
- **Location**: `pkg/rendering/particles.go:159`
- **Category**: Logic
- **Description**: `p.Life -= dt / p.MaxLife` — if `MaxLife` is 0 (from uninitialized or edge-case particle creation), this produces `+Inf`, causing the particle's life to become `NaN` and never expire.
- **Impact**: Immortal particles accumulate, degrading performance and causing visual artifacts.
- **Suggested Fix**: Validate `MaxLife > 0` during particle creation; skip update if `MaxLife <= 0`.

### [M-013] Pathfinding reconstructPath Has No Cycle Guard
- **Location**: `pkg/procgen/world/pathfinding.go:119-130`
- **Category**: Logic
- **Description**: `reconstructPath()` follows the `cameFrom` map backwards with no maximum iteration limit. If the A* implementation has a bug that creates a cycle in the `cameFrom` map, this becomes an infinite loop.
- **Impact**: Potential game freeze during world generation pathfinding, though unlikely under normal A* operation.
- **Suggested Fix**: Add a maximum iteration count (e.g., `width * height`) as a safety guard.

### [M-014] NPC Dialogue Selection May Loop Excessively
- **Location**: `pkg/npc/generator.go:393-408`
- **Category**: Logic / Performance
- **Description**: Dialogue selection loops to find unique dialogues without duplicates. If `count > len(options)`, the loop can never satisfy the uniqueness constraint and continues indefinitely.
- **Impact**: Game freeze during NPC interaction if the dialogue count exceeds available options.
- **Suggested Fix**: Clamp `count` to `min(count, len(options))` before the loop.

### [M-015] Trading Reputation Can Enable Exploits
- **Location**: `pkg/trading/barter.go:142`
- **Category**: Logic
- **Description**: Price threshold formula `1.1 - (rep * 0.25)` produces values > 1.0 if reputation is negative. This means players with negative reputation face prices above 110% base, but the formula doesn't clamp the output, allowing extreme values.
- **Impact**: Unclamped price thresholds could produce unexpected trading behavior at extreme reputation values.
- **Suggested Fix**: Clamp the result to a reasonable range (e.g., `[0.5, 2.0]`).

### [M-016] Score Calculation Overflow Risk
- **Location**: `pkg/game/endscreen.go:53-88`
- **Category**: Logic
- **Description**: Score calculation uses `int` type multiplications without bounds checking. In extremely long games, intermediate products could overflow `int` (platform-dependent, 32-bit on some systems).
- **Impact**: Score displays as negative or wraps around on very long games.
- **Suggested Fix**: Use `int64` for score calculations or cap at a maximum score.

### [M-017] Forage Description Panic on Unknown Genre
- **Location**: `pkg/game/forage.go:264-301`
- **Category**: Logic
- **Description**: Description lookup maps (`nothingTexts`, `foodTexts`, etc.) are keyed by genre. If `fm.genre` is not present in these maps, `seed.Choice()` receives a nil slice, which may panic.
- **Impact**: Crash when foraging in an unsupported genre configuration.
- **Suggested Fix**: Add a default/fallback genre entry, or check for nil slice before calling `seed.Choice()`.

### [M-018] Save Slot Loading Performs Full Disk I/O Per Slot
- **Location**: `pkg/saveload/load.go:142-163`
- **Category**: Performance
- **Description**: `ListSlots()` calls `sm.Load(slot)` for every non-empty slot (up to 11 slots), performing full deserialization of each save file just to display slot summaries.
- **Impact**: Save/load menu is slow to open, especially on slower storage. Blocks the UI thread during loading.
- **Suggested Fix**: Store lightweight metadata (summary, timestamp) in separate files or a single index file.

## Low Priority Issues

### [L-001] Hardcoded Font Width for Text Centering
- **Location**: `pkg/ux/menus.go:189-190`
- **Category**: UI
- **Description**: Text centering uses a hardcoded character width of 7 pixels, assuming a fixed-width font. If the font changes, all menu text alignment breaks.
- **Impact**: Minor visual misalignment if font is changed.
- **Suggested Fix**: Calculate actual text width from the font face metrics.

### [L-002] Missing Capacity Hints in Slice Allocations
- **Location**: `pkg/crew/party.go` (Living()), `pkg/crew/crew.go`, and others
- **Category**: Performance
- **Description**: Many `make([]T, 0)` calls lack capacity hints even when the maximum size is known (e.g., `len(p.members)`). This causes unnecessary reallocations during `append()`.
- **Impact**: Minor GC pressure from extra allocations. Negligible for small slices but adds up across the codebase.
- **Suggested Fix**: Add capacity hints: `make([]*CrewMember, 0, len(p.members))`.

### [L-003] Magic Number in Movement Cost
- **Location**: `pkg/resources/consumption.go:62-66`
- **Category**: Code Quality
- **Description**: `baseCost := float64(terrainCost) * 2.0` uses an undocumented magic constant. The multiplier's gameplay purpose is unclear.
- **Impact**: Makes balancing difficult; maintainers must guess the intent.
- **Suggested Fix**: Extract to a named constant like `MovementCostMultiplier`.

### [L-004] Headless API Parameter Mismatches
- **Location**: `pkg/rendering/animation_headless.go:19`, `pkg/rendering/portrait_headless.go:31`, `pkg/rendering/landmark_icon_headless.go:30`
- **Category**: Code Quality
- **Description**: Headless build constructors (e.g., `NewAnimatedTile`) take different parameter types than their non-headless counterparts (e.g., `int` vs `[]*ebiten.Image`). While this works due to build tags, it creates an inconsistent API surface.
- **Impact**: Confusion for developers; code that works in one build mode may not compile in the other.
- **Suggested Fix**: Use an interface type or consistent parameter signatures across build tags.

### [L-005] Delete Confirmation Dialog Allocates Per Frame
- **Location**: `pkg/ux/slots.go:264`
- **Category**: Performance / UI
- **Description**: `drawDeleteConfirm()` creates a new `ebiten.NewImage()` every frame while the delete confirmation dialog is visible.
- **Impact**: Minor memory pressure during the brief confirmation dialog display.
- **Suggested Fix**: Cache the dialog image.

### [L-006] Crew SkillExpThreshold Accepts Negative Levels
- **Location**: `pkg/crew/crew.go:170-176`
- **Category**: Logic
- **Description**: Custom `pow()` function casts `exp` to `int` — negative values become 0 iterations, returning 1.0. `SkillExpThreshold()` could receive negative `currentLevel` and return `100.0 * 1.0 = 100` regardless.
- **Impact**: Minor — negative levels shouldn't occur in normal gameplay, but the function silently returns incorrect values instead of erroring.
- **Suggested Fix**: Add a `currentLevel >= 0` assertion.

### [L-007] Insignia SymbolScale Range Documentation Mismatch
- **Location**: `pkg/vessel/insignia.go:154-162`
- **Category**: Code Quality
- **Description**: `SymbolScale` is generated in range `[0.3, 0.7]` but documentation comment says `0.3-0.8`. Minor doc/code inconsistency.
- **Impact**: None on gameplay; developer confusion only.
- **Suggested Fix**: Align code and documentation.

### [L-008] Config Validation Missing Seed and Volume Checks
- **Location**: `pkg/config/settings.go:247-259`
- **Category**: Logic
- **Description**: `Validate()` checks screen dimensions, tile size, and master volume, but does not validate `MusicVolume`, `SFXVolume`, or the `Seed` value. A seed of 0 may cause non-deterministic behavior in generators.
- **Impact**: Invalid configuration accepted silently; reproducibility may be affected.
- **Suggested Fix**: Validate all volume fields are in `[0, 1]` and seed is non-zero (or document that 0 means "use system time").

### [L-009] Relationship Network Not Cleared on Regeneration
- **Location**: `pkg/crew/relationship.go:155-166`
- **Category**: State
- **Description**: `GenerateInitialRelationships()` only adds new relationships; it never clears existing ones. If called multiple times (e.g., on crew refresh), relationships accumulate.
- **Impact**: Duplicate relationship entries between the same crew pairs; minor morale calculation errors.
- **Suggested Fix**: Clear the network before generating initial relationships.

### [L-010] Unchecked InputConfig Action Cast
- **Location**: `pkg/config/persistence.go:94-107`
- **Category**: Logic
- **Description**: `Action(b.Action)` casts raw integers from JSON config to the `Action` enum without bounds validation. Invalid action IDs are silently accepted.
- **Impact**: Invalid keybindings from corrupted config files could cause undefined behavior.
- **Suggested Fix**: Validate that the action ID is within the valid enum range before casting.

## Performance Optimizations

### [P-001] Implement Tile Image Cache
- **Location**: `pkg/rendering/renderer.go:82-88`
- **Issue**: Every `DrawTile()` call creates a new `ebiten.Image`. This is the #1 performance bottleneck.
- **Optimization**: Create a `map[int]*ebiten.Image` tile cache. Populate it lazily or at initialization. Reuse cached images via `DrawImage()` with translated `DrawImageOptions`.
- **Expected Gain**: Eliminate ~13,500 allocations/sec (15×15 viewport × 60 TPS). Expected 50-80% frame time reduction in worldmap rendering.

### [P-002] Cache All UI Overlay Images
- **Location**: `pkg/ux/panel.go:31`, `pkg/ux/events.go:96`, `pkg/ux/worldmap.go:120,131,160`, `pkg/ux/slots.go:264`
- **Issue**: UI overlay images (fog, markers, panels, event displays, dialogs) are all created per frame.
- **Optimization**: Create each overlay image once and store it as a struct field. Recreate only when dimensions change (e.g., on resize).
- **Expected Gain**: Eliminate 200-500+ allocations/frame across UI layers. Major GC pressure reduction.

### [P-003] Replace Pixel-by-Pixel Drawing with Batched DrawImage
- **Location**: `pkg/rendering/particles.go:346-361`, `pkg/rendering/lighting.go:403-418`, `pkg/rendering/genre_overlay.go:116-152`
- **Issue**: Particles, lighting, and genre overlays use `Set(x, y, color)` loops — each call is a separate GPU operation.
- **Optimization**: Pre-render small sprite images for particles. Use Ebitengine's `ColorScale` and `DrawImage()` for lighting. Use shader-based overlays for genre effects.
- **Expected Gain**: 10-100x speedup for particle/lighting rendering depending on particle count and light source count.

### [P-004] Cache Landmark Icon Images
- **Location**: `pkg/rendering/landmark_icon.go:115, 221, 325, 390, 470, 523`
- **Issue**: Six different landmark icon draw functions each create new images every frame.
- **Optimization**: Generate all landmark icons at initialization and store in a `map[LandmarkType]*ebiten.Image` cache.
- **Expected Gain**: Eliminate 6× `ebiten.NewImage()` per visible landmark per frame.

### [P-005] Lazy Lighting Recalculation
- **Location**: `pkg/rendering/lighting.go:367-435`
- **Issue**: Full lighting overlay is recomputed every frame even when nothing has changed.
- **Optimization**: Track a "dirty" flag for the lighting system. Only recompute when player moves, time-of-day changes, or light sources are modified. Cache the computed overlay between frames.
- **Expected Gain**: Reduce lighting computation from 60/sec to 1-5/sec in typical gameplay (player moves ~1-2 times per second).

### [P-006] Object Pooling for Particles
- **Location**: `pkg/rendering/particles.go`
- **Issue**: Particles are created and garbage-collected continuously. Each burst can create hundreds of particle structs.
- **Optimization**: Implement a `sync.Pool` or ring-buffer pool for particle structs. Recycle expired particles instead of allocating new ones.
- **Expected Gain**: Reduce GC pressure during particle-heavy scenes by 50-90%.

### [P-007] Batch Save Slot Metadata Loading
- **Location**: `pkg/saveload/load.go:142-163`
- **Issue**: `ListSlots()` loads and deserializes every save file to display summaries.
- **Optimization**: Write a lightweight metadata sidecar file (`.meta.json`) alongside each save containing only the summary fields. Read only metadata files for `ListSlots()`.
- **Expected Gain**: 10-100x faster save menu opening, depending on save file sizes.

## Prioritized Recommendations

1. **Critical** — [C-001] Cache tile images in renderer — eliminates the largest source of memory leaks and frame drops
2. **Critical** — [C-002/C-003] Add input debouncing and mutual exclusion — prevents double-firing events and turn corruption
3. **Critical** — [C-004] Cache pause overlay image — stops memory leak during pause
4. **Critical** — [C-005] Clear entity components on despawn — prevents state corruption in ECS
5. **Critical** — [C-006] Implement save format migration — prevents player data loss on updates
6. **High** — [H-001/H-002/H-003/H-004] Cache all UI overlay images — bulk elimination of per-frame allocations
7. **High** — [H-005/H-006/H-008] Replace pixel-by-pixel rendering with batched draw calls
8. **High** — [H-010] Use dynamic delta time instead of hardcoded 1/60
9. **High** — [H-011] Fix relationship pairKey truncation
10. **High** — [H-012] Add integer overflow protection in cargo calculations
11. **High** — [H-013/H-014/H-015/H-016] Add zero-denominator guards across audio and encounter systems
12. **Performance** — [P-001/P-002] Tile and overlay caching (highest impact, lowest effort)
13. **Performance** — [P-003/P-006] Particle rendering and pooling overhaul
14. **Performance** — [P-005] Lazy lighting recalculation
15. **Medium** — [M-001 through M-018] Address remaining logic, state, and balance issues

## Audit Methodology

### Approach
This audit was conducted via static analysis of the complete Go source tree (~42,800 lines across 160+ files). All `.go` files were read and analyzed systematically by category:

1. **Core game loop**: `cmd/voyage/main.go`, `pkg/game/*.go` — focused on Update/Draw/Layout patterns, input handling, state transitions
2. **Rendering pipeline**: `pkg/rendering/*.go` — focused on image allocation patterns, draw order, camera transforms
3. **UX layer**: `pkg/ux/*.go` — focused on UI allocation patterns, input consumption, layout calculations
4. **Engine/ECS**: `pkg/engine/*.go` — focused on entity lifecycle, component management, system safety
5. **Game systems**: `pkg/vessel/*.go`, `pkg/crew/*.go`, `pkg/resources/*.go`, `pkg/events/*.go`, `pkg/encounters/*.go`, `pkg/trading/*.go`, `pkg/weather/*.go`, `pkg/audio/*.go`, `pkg/procgen/*.go`, `pkg/npc/*.go`, `pkg/council/*.go`, `pkg/destination/*.go`, `pkg/saveload/*.go`, `pkg/config/*.go` — focused on logic correctness, boundary conditions, state management
6. **Cross-cutting**: Verified exact line numbers for all critical and high-severity findings

### Key Patterns Searched
- `ebiten.NewImage(` in Draw-path functions (memory leak pattern)
- `IsKeyPressed(` without corresponding release detection (input debounce)
- Integer multiplication without overflow guards
- Division without zero-denominator checks
- Map/slice access without bounds validation
- State mutation in Draw() functions
- Hardcoded constants that should be configurable

### Assumptions
- Ebitengine TPS is 60 (default) unless explicitly configured otherwise
- The game targets desktop platforms (not mobile/web)
- Headless build mode is used for testing only
- `int` is 64-bit on the target platform (though 32-bit overflow risks are flagged)

### Limitations
- **Static analysis only** — no runtime profiling, no execution of the game
- **No build verification** — findings are based on code reading; some issues flagged as "Potential" may be mitigated by runtime conditions not visible in static analysis
- **Line numbers approximate** — exact line numbers were verified for critical and high-severity issues but may drift with future commits
- **No test coverage analysis** — existing test files were not evaluated for coverage adequacy
- **Ebitengine version assumptions** — analysis assumes Ebitengine v2 API; some API behaviors may differ across minor versions
