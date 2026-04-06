# Game Repair Plan

## Critical Issues (Blocks Playability)

### Compilation Failures (Game Cannot Build)

- [ ] **Undefined `InputAction` type in input manager** - Location: `pkg/input/manager.go:153` - Impact: Non-headless build fails to compile. The `hasAction` method uses type `InputAction` but the actual type defined in `pkg/input/action.go:22` is `Action`. This blocks all non-headless builds entirely.

- [ ] **Double-pointer type mismatch in event snapshot** - Location: `pkg/game/session.go:48` - Impact: Non-headless build fails to compile. `s.eventQueue.Pending()` returns `[]*Event`, so `pending[0]` is already `*Event`. Taking `&e` creates `**Event` which cannot be assigned to `*events.Event`. Should be `s.currentEventSnapshot = pending[0]`.

### Event System Broken (Core Gameplay Non-Functional)

- [ ] **Event choice selection off-by-one error** - Location: `pkg/game/session.go:194-195` - Impact: Players cannot correctly resolve events. Choice IDs start at 1 (set in `pkg/events/event.go:87` as `len(e.Choices)+1`), but the code passes the loop index `i` (0-based) to `resolveEvent()`. Key 1 → choiceID 0 → `GetChoice(0)` returns nil (no match), Key 2 → choiceID 1 → selects choice 1 instead of choice 2, etc. First choice is always unreachable; all others are shifted. Depends on: compilation fixes above.

### Input Handling Broken

- [ ] **Menu uses `IsKeyPressed` instead of `IsKeyJustPressed`** - Location: `pkg/game/game.go:130`, `pkg/game/game.go:154` - Impact: `handleMenuInput()` and `handleGameOverInput()` use `ebiten.IsKeyPressed` (continuous detection) instead of `inpututil.IsKeyJustPressed` (edge detection). Pressing Enter/Space fires every frame, causing rapid unintended state transitions. Menu → Playing transition happens on the same frame the key is held, and GameOver → Menu fires repeatedly.

- [ ] **Direction state not cleared between key repeat intervals** - Location: `pkg/input/manager.go:123-130` - Impact: When a direction key is held and the repeat interval hasn't elapsed, `currentState.Direction` retains the previous frame's value instead of being reset to `DirectionNone`. This causes unintended movement between repeat ticks. The headless version (`manager_headless.go`) correctly handles this, confirming the bug.

## High Priority (Degrades Experience)

### UI Rendering Issues

- [ ] **Save/load slots invisible on small panels** - Location: `pkg/ux/slots.go:158` - Impact: `visibleSlots = (panelHeight - 100) / 50` produces 0 when `panelHeight < 150`. When `visibleSlots == 0`, `startIdx` calculation becomes `selectedIndex + 1`, skipping the selected slot entirely. Save/load screen becomes unusable on smaller screen sizes.

- [ ] **Event overlay text overflow for long words** - Location: `pkg/ux/events.go:201-207` - Impact: `addWordToLine()` places words longer than `maxWidth` on a single line without truncation or splitting. Long procedurally-generated words overflow the event overlay boundaries, making event text unreadable. The word is returned as-is on a new line regardless of length.

- [ ] **Pause overlay not recreated on window resize** - Location: `pkg/game/game.go:264-269` - Impact: The pause overlay image is lazily created at `g.width × g.height` and cached. If the window resizes, the overlay retains its original dimensions, leaving uncovered screen areas or displaying incorrectly.

- [ ] **Minimap renders off-screen on narrow windows** - Location: `pkg/ux/minimap.go:103` - Impact: Position calculated as `screenW - m.width - 10` with no bounds check. If `screenW < m.width + 10`, the minimap renders at negative X coordinates (off-screen to the left).

### Game Logic Issues

- [ ] **Event generation after game over condition** - Location: `pkg/game/session.go:222-237` - Impact: In `advanceTurn()`, `consumeResources()` and `maybeGenerateEvent()` execute before `checkConditions()`. If resources deplete to a loss condition, a new event can still be generated and queued. This event then displays over the game-over screen or persists into a new game.

- [ ] **Negative time advance not clamped** - Location: `pkg/game/session.go:211-215` - Impact: `outcome.TimeAdvance` is only clamped against `maxTimeAdvance` (upper bound). A negative value makes the `for` loop condition immediately false, silently skipping turn advancement. No error feedback.

- [ ] **Stacking morale penalties for simultaneous resource depletion** - Location: `pkg/game/session_common.go:120-125` - Impact: When both food AND water are depleted in the same turn, morale receives -5 AND -8 (total -13 per turn). Both penalties apply independently with no cap, causing morale to collapse extremely quickly, creating an unfair difficulty spike that prevents recovery.

### Input Architecture Issues

- [ ] **Input Manager abstraction layer entirely unused** - Location: `pkg/input/` (entire package) vs `pkg/game/game.go`, `pkg/game/session.go` - Impact: The `pkg/input` package implements a full input Manager with touch, swipe, key repeat, and action deduplication, but the game code (`pkg/game/`) directly calls raw Ebiten API functions (`ebiten.IsKeyPressed`, `inpututil.IsKeyJustPressed`) instead of using the Manager. This creates inconsistent input behavior across the codebase and renders the input package dead code.

- [ ] **Missing touch action deduplication** - Location: `pkg/input/manager.go:243-257` - Impact: The `handleTouchEnd` path appends `ActionConfirm` without checking for duplicates (no `HasAction()` check), unlike the mouse click path at line 307 which correctly deduplicates. On touch devices, simultaneous touch + keyboard could trigger duplicate confirm actions.

## Medium Priority (Polish/Optimization)

### Performance Issues

- [ ] **Per-pixel operations in post-processing effects** - Location: `pkg/rendering/postprocess.go:193,341` - Impact: `ApplyScanlines()` and `ApplySepia()` use `img.At(x,y)` to read every pixel individually, and `result.Set(x,y,...)` to write. These are extremely expensive Ebiten operations in tight loops. If post-processing runs per frame, this causes severe frame drops at higher resolutions.

- [ ] **Per-pixel Set in vignette generation** - Location: `pkg/rendering/postprocess.go:161-175` - Impact: `generateVignetteCache()` uses `Set(x, y, ...)` in a double loop for every pixel. While cached (not per-frame), the first-time generation causes a visible stutter.

- [ ] **Per-pixel Set in radial light rendering** - Location: `pkg/rendering/lighting.go:414` - Impact: `drawRadialLight()` uses per-pixel `Set()` operations in a double loop. Causes frame stutters when point lights are updated.

- [ ] **Per-pixel marker drawing on world map** - Location: `pkg/ux/worldmap.go:149-159,184-189` - Impact: Destination and vessel markers are created using per-pixel `Set()` for line drawing. Markers recreated on first use with slow operations.

- [ ] **Mass image creation in PostProcessor.Apply()** - Location: `pkg/rendering/postprocess.go:58,136,221,244,286,337` - Impact: Multiple `ebiten.NewImage()` calls per post-processing pass. If called every frame, creates 5+ images per frame, causing GPU memory churn and GC pressure.

- [ ] **Per-pixel animation frame generation** - Location: `pkg/rendering/animation.go:103-120,134-148,186-196` - Impact: Animated tiles (water, grass, fire) are generated pixel-by-pixel using `img.Set()`. Creates 12 images (3 types × 4 frames) during initialization with very slow per-pixel operations.

### UI Polish Issues

- [ ] **Cargo screen scroll indicator shows negative values** - Location: `pkg/ux/cargo.go:171` - Impact: When `len(items) <= maxVisible` (8), the denominator `len(items)-maxVisible+1` becomes zero or negative, displaying misleading indicators like `[1/-3]` or `[1/0]`.

- [ ] **Scrollbar thumb can extend past track bounds** - Location: `pkg/ux/leaderboard.go:360-367` - Impact: No bounds checking on `thumbPos + thumbSize` means the scrollbar thumb can visually extend past the track area at high scroll offsets. Minor visual glitch.

- [ ] **Division by zero risk in vignette generation** - Location: `pkg/rendering/postprocess.go:149-176` - Impact: If screen width or height is 0, `maxDist` and `maxDistSq` become 0, causing division by zero at the distance calculation. Unlikely in normal operation but possible during initialization or window minimization.

### Code Quality Issues

- [ ] **No `Image.Dispose()` calls anywhere in codebase** - Location: `pkg/rendering/` (all files) - Impact: Ebiten images are never explicitly disposed. For long-running sessions or genre switches that recreate images, GPU memory accumulates without cleanup.

- [ ] **Redundant F3 key-release detection** - Location: `pkg/game/game.go:162-169`, `pkg/game/session.go:94-102` - Impact: Manual key-release tracking (`f3WasPressed` flag) reimplements what `inpututil.IsKeyJustPressed(ebiten.KeyF3)` already provides. Works correctly but adds unnecessary state and complexity.

- [ ] **Event queue resolved list grows unbounded** - Location: `pkg/events/queue.go:148-163` - Impact: `q.resolved` slice accumulates all resolved events forever. Over very long game sessions (1000+ turns), this creates steadily growing memory usage with no pruning mechanism.

## Resolution Notes

### Build Dependencies
1. **Fix compilation errors first** (items 1-2) — nothing else can be tested until the non-headless build compiles.
2. **Event choice off-by-one** (item 3) depends on compilation fixes and should be fixed immediately after.
3. **Input fixes** (items 4-5) can be done in parallel with event fixes.

### Fix Ordering Recommendations
- **Phase 1**: Compilation fixes (`InputAction` → `Action`, `&e` → `pending[0]`)
- **Phase 2**: Event system fix (choice ID off-by-one `i` → `i+1`)
- **Phase 3**: Input edge detection fixes (`IsKeyPressed` → `IsKeyJustPressed` in menus)
- **Phase 4**: UI rendering fixes (slots visibility, text overflow, pause overlay)
- **Phase 5**: Game logic fixes (turn ordering, time advance clamping, morale stacking)
- **Phase 6**: Performance optimization (per-pixel operations, image caching)

### Testing Strategy
- The headless build (`go build -tags headless ./...`) passes all tests, confirming headless code paths are clean
- Non-headless build requires X11 headers and fails due to the two compilation errors identified above
- After compilation fixes, run `go test -tags headless -race ./...` to verify no regressions
- UI rendering issues require manual visual testing or screenshot comparison

### Key Architectural Observation
The `pkg/input` Manager provides a well-designed abstraction layer with touch support, key repeat, and action deduplication. However, it is completely unused by the game code in `pkg/game/`. The game directly calls raw Ebiten APIs, bypassing all the Manager's features. A significant improvement would be to wire the Input Manager into the game loop, which would fix several input issues simultaneously and unify input handling.
