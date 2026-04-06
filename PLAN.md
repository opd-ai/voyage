# Game Playability Issues

## Critical Blockers (prevents launch/play)

- [ ] **Event category/weight mismatch causes index panic risk** | Location: `pkg/events/queue.go:47-48` | Fix: `AllEventCategories()` returns 7 categories (Weather, Encounter, Discovery, Hardship, Windfall, Hazard, Crew) but `weights` slice has only 5 entries `{0.15, 0.25, 0.20, 0.25, 0.15}`. `seed.WeightedChoice` at `pkg/procgen/seed/seed.go:112-127` iterates `weights` but indexes into `choices` — when `len(choices) != len(weights)`, the weight distribution is wrong and always falls back to `CategoryCrew` (index 6). Fix by either adding 2 more weights for Hazard and Crew categories, or limiting `AllEventCategories()` to the 5 that have templates.

- [ ] **Small map dimensions cause panic in `Intn(0)`** | Location: `pkg/procgen/world/generator.go:182-189` | Fix: `placeOriginDestination()` calls `g.gen.Intn(w.Width/4)` and `g.gen.Intn(w.Height/3)`. For maps where Width<4 or Height<3, integer division yields 0, and `Intn(0)` panics. While default config uses 50×50, no validation prevents smaller maps. Add `if w.Width < 8 { w.Width = 8 }; if w.Height < 6 { w.Height = 6 }` guard, or use `max(1, w.Width/4)`.

- [ ] **ESC key rapid-toggles between Playing/Paused in Game struct** | Location: `pkg/game/game.go:127,137` | Fix: `handlePlayingInput()` and `handlePausedInput()` both use `ebiten.IsKeyPressed(ebiten.KeyEscape)` which fires every frame while held. This causes uncontrollable rapid state toggling. Change to `inpututil.IsKeyJustPressed(ebiten.KeyEscape)` as already done in `session.go:69,173`. Note: This only affects the standalone `Game` path, not the `GameSession` path used by `cmd/voyage/main.go`.

## High Priority (severely degrades gameplay)

- [ ] **Player and destination markers are invisible** | Location: `pkg/game/session.go:258,261` | Fix: `DrawTile` is called with tile types `10` (player marker) and `11` (destination marker), but all palettes only define 4 tile colors (indices 0–3). `GetTileColor()` at `pkg/rendering/renderer_core.go:141-146` falls back to `Background` color, making markers indistinguishable from background. Add dedicated player/destination colors to the `Palette` struct or extend `TileColors` arrays to include indices for markers.

- [ ] **`formatRelationshipText` corrupts string during in-place mutation** | Location: `pkg/events/relationships.go:124-135` | Fix: Loop iterates `i < len(result)-1` but mutates `result` string length mid-iteration when replacing `%A`/`%B` with names. After replacing `%A`, the loop index `i` may point to the wrong position, skipping `%B` or accessing out-of-bounds. Replace with `strings.ReplaceAll(result, "%A", nameA)` followed by `strings.ReplaceAll(result, "%B", nameB)`.

- [ ] **CategoryHazard and CategoryCrew missing from `categoryTemplates`** | Location: `pkg/events/queue.go:64-70` | Fix: The map only contains 5 of 7 categories. When Hazard or Crew is selected by `WeightedChoice`, `generateForCategory()` at line 76-79 falls back to `CategoryHardship` templates. Add Hazard and Crew entries to `categoryTemplates` to enable all event types.

- [ ] **`Resolve()` returns nil outcome without caller guard** | Location: `pkg/events/queue.go:100-114`, callers at `pkg/game/session.go:140`, `pkg/game/session_common.go:74` | Fix: `Resolve()` returns nil for invalid eventID or choiceID. Session callers at `session.go:140-141` check for nil before calling `applyOutcome`, so immediate crash is prevented. However, `pkg/events/resolution.go`'s `Apply()` dereferences the outcome pointer without nil check — if used elsewhere, it panics. Add nil guard in `Apply()`.

- [ ] **Audio methods return nil when muted, risking caller panics** | Location: `pkg/audio/player.go:74,113,141` | Fix: `PlaySFX()`, `GenerateAmbientMusic()`, and `CrossfadeMusicTo()` return nil when `p.muted == true`. Any caller iterating over the returned slice/struct without nil-checking will panic. Return empty `[]float64{}` or empty `*AmbientLoop{}` instead of nil.

## Medium Priority (playable but broken)

- [ ] **Movement input uses `IsKeyPressed` causing multi-tile-per-frame movement** | Location: `pkg/game/session.go:99-112` | Fix: `getMovementInput()` uses `ebiten.IsKeyPressed()` which fires every frame while held, moving the player every tick (60 tiles/sec at 60 TPS). Use `inpututil.IsKeyJustPressed()` for single-tile movement per keypress, or implement a movement cooldown/repeat delay.

- [ ] **Menu input uses `IsKeyPressed` causing instant state skip** | Location: `pkg/game/session.go:61-63` | Fix: `handleMenuInput()` uses `IsKeyPressed(KeyEnter)` — if the player presses Enter on the menu, the game may process multiple frames worth of input. Use `inpututil.IsKeyJustPressed` for clean single-press detection.

- [ ] **`Consume()` accepts negative amounts, enabling resource exploitation** | Location: `pkg/resources/resources.go:130-139` | Fix: Calling `Consume(type, -100)` effectively adds resources since the subtraction of a negative increases the value. Add `if amount < 0 { return false }` guard.

- [ ] **`Vessel.TakeDamage()` accepts negative amounts, healing without limit** | Location: `pkg/vessel/vessel.go:122-130` | Fix: Calling `TakeDamage(-50)` increases integrity beyond maxIntegrity with no clamping. Add `if amount <= 0 { return false }` guard.

- [ ] **`Party.Add()` accepts nil crew members** | Location: `pkg/crew/party.go:41-47` | Fix: No nil check on `member` parameter. If nil is added, `LivingCount()` will panic when iterating and accessing `m.IsAlive`. Add `if member == nil { return false }` at the top of the function.

- [ ] **`CalculateMoveCost` divides by vessel `Speed()` without zero-check** | Location: `pkg/game/movement.go:35` | Fix: `fuelCost := mm.baseFuelCost * float64(terrain.MovementCost) / v.Speed()` produces `+Inf` if speed is 0 (possible with destroyed vessel). Not currently called in active session path but will produce `NaN`/`Inf` if integrated. Add `if v.Speed() == 0 { return math.Inf(1), timeCost }` guard or prevent the call.

- [ ] **`DifficultyName` and `ActionName` return empty string for invalid values** | Location: `pkg/config/settings.go:186-195,56-72` | Fix: Map lookup returns zero-value (empty string) for undefined constants. Return "Unknown" as default fallback.

- [ ] **`SetMusicState` silently ignores invalid states** | Location: `pkg/audio/music.go:96-113` | Fix: No default case in switch — passing an invalid `MusicState` value keeps previous BPM, creating audio state inconsistency. Add a default case that sets a fallback BPM or logs a warning.

## Resolution Order

1. Event category/weight mismatch (`queue.go:47-48`) — prevents correct event generation, most likely crash vector
2. Player/destination markers invisible (`session.go:258,261`, `renderer_core.go`) — game is unnavigable without markers
3. `formatRelationshipText` string corruption (`relationships.go:124-135`) — corrupts event text display
4. Small map Intn(0) panic (`generator.go:182-189`) — crashes with certain configs
5. Missing Hazard/Crew category templates (`queue.go:64-70`) — complements fix #1
6. Movement multi-tile-per-frame (`session.go:99-112`) — movement is uncontrollable
7. Menu instant-skip (`session.go:61-63`) — menu is barely usable
8. Nil return from audio when muted (`player.go:74,113,141`) — crashes if audio is muted
9. `Resolve()` nil outcome without guard (`queue.go:100-114`) — potential crash on edge case
10. `Party.Add()` nil member (`party.go:41-47`) — crash if code path produces nil member
11. Negative `Consume`/`TakeDamage` exploits (`resources.go:130`, `vessel.go:122`) — game balance issues
12. ESC rapid toggle in Game struct (`game.go:127,137`) — only affects non-primary code path
13. Division by zero in `CalculateMoveCost` (`movement.go:35`) — only affects unused code path
14. `DifficultyName`/`ActionName` empty returns (`settings.go`) — cosmetic issue
15. `SetMusicState` invalid state handling (`music.go:96`) — minor audio inconsistency

## Notes

### Architecture Observations
- The codebase has two parallel game implementations: `Game` (game.go) and `GameSession` (session.go/session_common.go). The `cmd/voyage/main.go` entry point uses `GameSession`, making `Game` effectively dead code. Consider removing `Game` or unifying the two.
- The `MovementManager`, `RestManager`, `ForageManager`, and `TimeManager` are well-designed subsystems but are not integrated into the `GameSession` loop — movement in session.go is bare coordinate manipulation without fuel costs, terrain effects, or rest mechanics. Integrating these would significantly improve gameplay depth.
- The headless/non-headless build tag strategy is well-implemented for testing but means the test suite can only validate logic, not rendering or input handling.
- All 37 packages compile and pass tests under `headless` tag. Non-headless build requires X11 dev headers (expected for Ebitengine on Linux).

### Technical Debt
- Event system uses 7 categories but only 5 have templates and weights — this asymmetry suggests incomplete feature work.
- `seed.WeightedChoice` assumes `len(choices) == len(weights)` but does not validate this, making it fragile against caller mistakes.
- String formatting in `formatRelationshipText` should use Go's `strings.Replacer` or `text/template` instead of manual byte-by-byte replacement.
- Resource consumption in the session loop (`consumeResources`) does not integrate with `TimeManager` seasons, `MovementManager` terrain costs, or difficulty settings — these subsystems exist but are disconnected.
- The `endscreen.go` `EndStats` struct is never populated during gameplay — the session transitions to `StateGameOver` but doesn't create or display end-game statistics.
