# AUDIT — 2026-04-05

## Project Goals

Voyage is a **100% procedural travel simulator** inspired by Oregon Trail, FTL, and Organ Trail. The project claims:

1. **Fully procedural generation** — every map, event, crew, vessel, audio, and narrative generated from a single seed
2. **Five genre themes** — Fantasy, Sci-fi, Horror, Cyberpunk, Post-apocalyptic
3. **No bundled assets** — all visual, audio, and narrative content generated at runtime
4. **Core gameplay systems** — resource attrition, crew mortality, vessel integrity, procedural events, route choice
5. **GenreSwitcher interface** — all systems implement runtime genre switching
6. **Deterministic RNG** — same seed produces same game
7. **Ebitengine rendering** — graphical game using Ebitengine v2.9
8. **Save/Load system** — multiple slots with autosave
9. **Win/Lose conditions** — reach destination vs. all crew dead/vessel destroyed/morale collapse

### Stated Status (README)

The README claims "Early Development — Core engine complete. Full gameplay coming soon" with checkmarks on all v1.0 and v2.0 ROADMAP items.

---

## Goal-Achievement Summary

| Goal | Status | Evidence |
|------|--------|----------|
| ECS Framework with GenreSwitcher | ✅ Achieved | `pkg/engine/genre.go:22-24`, `pkg/engine/world.go:31-36` |
| Seed-based deterministic RNG | ✅ Achieved | `pkg/procgen/seed/seed.go:11-27`, SHA-256 hash derivation |
| Five genre support | ✅ Achieved | All 5 GenreIDs defined, all systems implement SetGenre |
| Procedural world map generation | ✅ Achieved | `pkg/procgen/world/generator.go:77-100` |
| Resource management (6-axis) | ✅ Achieved | `pkg/resources/resources.go:8-21` |
| Crew/party system | ✅ Achieved | `pkg/crew/crew.go`, traits, skills, backstory |
| Vessel/transport system | ✅ Achieved | `pkg/vessel/vessel.go`, modules, cargo, upgrades |
| Procedural event system | ✅ Achieved | `pkg/events/`, grammar-based text generation |
| Audio synthesis (waveforms, ADSR, SFX) | ✅ Achieved | `pkg/audio/waveforms.go`, `pkg/audio/music.go` |
| UI/HUD/Menus | ✅ Achieved | `pkg/ux/` — 14 files with genre theming |
| Win/lose conditions | ✅ Achieved | `pkg/game/conditions.go:51-91` |
| Save/load system | ✅ Achieved | `pkg/saveload/` — multiple slots, autosave |
| Configuration with CLI flags | ✅ Achieved | `cmd/voyage/main.go:30-34`, `pkg/config/` |
| No bundled assets | ✅ Achieved | No `.png`, `.jpg`, `.mp3`, `.wav`, `.ogg` files found |
| Weather system (8+ types) | ✅ Achieved | `pkg/weather/types.go`, genre-specific hazards |
| Trading and supply posts | ✅ Achieved | `pkg/trading/` — inventory, barter, reputation |
| Tactical encounters | ✅ Achieved | `pkg/encounters/` — roles, resolution phases |
| Crew council voting | ✅ Achieved | `pkg/council/` — trait-based voting |
| Foraging/scavenging | ✅ Achieved | `pkg/game/forage.go` — diminishing returns |
| Dynamic lighting | ✅ Achieved | `pkg/rendering/lighting.go` — day/night, genre presets |
| Particle effects | ✅ Achieved | `pkg/rendering/particles.go` — movement trails, weather |
| **Full game loop integration** | ⚠️ Partial | `cmd/voyage/main.go:80-91` — TODO comment, no Ebitengine.RunGame() |
| **Animated sprites** | ❌ Missing | ROADMAP v3.0 unchecked items at lines 272-276 |
| **Adaptive multi-layer music** | ❌ Missing | ROADMAP v3.0 unchecked items at lines 279-283 |
| **Positional audio** | ❌ Missing | ROADMAP v3.0 unchecked items at lines 285-288 |
| **Genre post-processing** | ❌ Missing | ROADMAP v3.0 unchecked items at lines 292-297 |
| **Dynamic minimap** | ❌ Missing | ROADMAP v3.0 unchecked items at lines 299-303 |

---

## Findings

### CRITICAL

- [ ] **Game loop not integrated** — `cmd/voyage/main.go:80-91` — The main entry point contains a TODO comment listing 8 uninitialized systems and never calls `ebiten.RunGame()`. The `demo()` function only spawns one entity and prints RNG values. This means the game cannot actually be played despite all subsystems being implemented.
  
  **Remediation:** Complete the `main()` function to initialize all game systems and call `g.Run()`:
  ```go
  // After line 79 in cmd/voyage/main.go, replace TODO block with:
  worldGen := world.NewGenerator(masterSeed, genre)
  worldMap := worldGen.Generate(50, 50)
  
  cfg := game.Config{
      Width:    800,
      Height:   600,
      TileSize: 16,
      Seed:     masterSeed,
      Genre:    genre,
  }
  g := game.NewGame(cfg)
  if err := g.Run(); err != nil {
      log.Fatal(err)
  }
  ```
  **Validation:** `go build ./cmd/voyage && ./voyage --seed 12345` should open a playable window.

### HIGH

- [ ] **Crew relationship network not used in gameplay** — `pkg/crew/relationship.go:45` — The `RelationshipNetwork` struct is defined with methods for bond/rivalry tracking, but no game system queries relationships to generate crew-specific events.

  **Remediation:** In `pkg/events/generator.go`, add relationship-based event generation:
  ```go
  func (g *Generator) GenerateCrewRelationshipEvent(network *crew.RelationshipNetwork) *Event {
      // Check for high bond/rivalry pairs and generate appropriate events
  }
  ```
  **Validation:** `go test -v ./pkg/events/... -run TestRelationshipEvents`

- [ ] **F3 debug toggle uses key press, not key release** — `pkg/game/game.go:143-145` — Using `IsKeyPressed` causes multiple toggles per frame. This is marked with `// BUG` comment at line 141.

  **Remediation:** Replace `IsKeyPressed` with state tracking:
  ```go
  var f3WasPressed bool
  if ebiten.IsKeyPressed(ebiten.KeyF3) {
      if !f3WasPressed {
          g.debugMode = !g.debugMode
      }
      f3WasPressed = true
  } else {
      f3WasPressed = false
  }
  ```
  **Validation:** `go build ./cmd/voyage && ./voyage` — press F3 once, debug overlay should toggle once.

- [ ] **NPC alignment variance function has high complexity** — `pkg/npc/generator.go:applyAlignmentVariance` — Cyclomatic complexity 11.9, highest in codebase. Complex nested logic increases bug risk.

  **Remediation:** Extract alignment calculation into a lookup table:
  ```go
  var alignmentVarianceTable = map[FactionType]map[NPCType]float64{...}
  func applyAlignmentVariance(npc *NPC, faction FactionType) {
      npc.Alignment += alignmentVarianceTable[faction][npc.Type]
  }
  ```
  **Validation:** `go-stats-generator analyze ./pkg/npc --format json | jq '.functions[] | select(.name=="applyAlignmentVariance") | .complexity'` should show cyclomatic < 10.

### MEDIUM

- [ ] **Packages with no test files** — `pkg/procgen/event/`, `pkg/world/` — Two packages have `[no test files]` per `go test ./...` output.

  **Remediation:** Add basic test files:
  ```go
  // pkg/procgen/event/event_test.go
  package event
  
  import "testing"
  
  func TestEventGeneration(t *testing.T) {
      // Test deterministic event generation
  }
  ```
  **Validation:** `go test -v ./pkg/procgen/event/... ./pkg/world/...`

- [ ] **Code duplication in UI drawing** — `pkg/ux/menus.go:164-178` and `pkg/ux/slots.go:153-167` — 15-line clone detected (0.36% duplication ratio).

  **Remediation:** Extract shared drawing logic to `pkg/ux/panel.go`:
  ```go
  func drawCenteredPanel(screen *ebiten.Image, x, y, w, h int, opts ...PanelOption) {...}
  ```
  **Validation:** `go-stats-generator analyze ./pkg/ux --sections duplication` should show 0 clone pairs.

- [ ] **Magic numbers in music generation** — `pkg/audio/music.go:66-77` — Hard-coded envelope values (0.5, 0.3, 0.6, 0.5) without explanation.

  **Remediation:** Define named constants:
  ```go
  const (
      bassBeatAttack  = 0.5
      bassBeatDecay   = 0.3
      bassBeatSustain = 0.6
      bassBeatRelease = 0.5
  )
  ```
  **Validation:** `grep -c "0\.[0-9]" pkg/audio/music.go` should decrease.

- [ ] **Tile key generation uses rune conversion** — `pkg/game/forage.go:140-141` — `string(rune(x))` is incorrect for values > 127 and produces invalid UTF-8 for negative values.

  **Remediation:** Use fmt.Sprintf:
  ```go
  func (fm *ForageManager) tileKey(x, y int) string {
      return fmt.Sprintf("%d,%d", x, y)
  }
  ```
  **Validation:** `go test -v ./pkg/game/... -run TestForage`

### LOW

- [ ] **Naming convention violations** — 28 identifier violations detected by go-stats-generator: `EnvelopeIdle` (acronym), `ConfigExists` (package stutter), etc.

  **Remediation:** Rename per Go conventions. Examples:
  - `ConfigExists` → `Exists` (called as `config.Exists()`)
  - `CrewMember` → `Member` (called as `crew.Member`)
  
  **Validation:** `go-stats-generator analyze . --sections naming` should show 0 violations.

- [ ] **Low cohesion in benchmark package** — `pkg/benchmark/` — 0.0 cohesion score, 1 file, 0 functions in main package.

  **Remediation:** Either remove the empty package or add benchmark functions:
  ```go
  package benchmark
  
  import "testing"
  
  func BenchmarkWorldGeneration(b *testing.B) {...}
  ```
  **Validation:** `go test -bench=. ./pkg/benchmark/...` should run benchmarks.

- [ ] **Feature envy in 142 methods** — Methods that access more data from other types than their own receiver, per go-stats-generator analysis.

  **Remediation:** Review high-impact methods and relocate to appropriate types. This is a code smell, not a bug.

  **Validation:** Manual review of flagged methods.

---

## Metrics Snapshot

| Metric | Value |
|--------|-------|
| Total Lines of Code | 8,902 |
| Total Functions | 307 |
| Total Methods | 923 |
| Total Structs | 181 |
| Total Interfaces | 3 |
| Total Packages | 22 |
| Total Files | 121 |
| Average Function Length | 8.0 lines |
| Average Complexity | 2.7 |
| Highest Complexity Function | `applyAlignmentVariance` (11.9) |
| Functions > 50 lines | 13 (1.1%) |
| Documentation Coverage | 84.2% |
| Package Coverage | 100.0% |
| Function Coverage | 85.7% |
| Type Coverage | 83.9% |
| Duplication Ratio | 0.36% |
| Clone Pairs | 7 |
| Circular Dependencies | 0 |
| Dead Code (unreferenced functions) | 9 |
| Tests Passing | ✅ All (20 packages) |
| Race Detector Issues | 0 |
| `go vet` Issues | 0 |

---

## Test Results

```
go test -tags headless -race ./...
?   	github.com/opd-ai/voyage/cmd/voyage	[no test files]
ok  	github.com/opd-ai/voyage/pkg/audio
ok  	github.com/opd-ai/voyage/pkg/benchmark
ok  	github.com/opd-ai/voyage/pkg/config
ok  	github.com/opd-ai/voyage/pkg/council
ok  	github.com/opd-ai/voyage/pkg/crew
ok  	github.com/opd-ai/voyage/pkg/destination
ok  	github.com/opd-ai/voyage/pkg/encounters
ok  	github.com/opd-ai/voyage/pkg/engine
ok  	github.com/opd-ai/voyage/pkg/events
ok  	github.com/opd-ai/voyage/pkg/game
ok  	github.com/opd-ai/voyage/pkg/npc
?   	github.com/opd-ai/voyage/pkg/procgen/event	[no test files]
ok  	github.com/opd-ai/voyage/pkg/procgen/seed
ok  	github.com/opd-ai/voyage/pkg/procgen/world
ok  	github.com/opd-ai/voyage/pkg/rendering
ok  	github.com/opd-ai/voyage/pkg/resources
ok  	github.com/opd-ai/voyage/pkg/saveload
ok  	github.com/opd-ai/voyage/pkg/trading
ok  	github.com/opd-ai/voyage/pkg/ux
ok  	github.com/opd-ai/voyage/pkg/vessel
ok  	github.com/opd-ai/voyage/pkg/weather
?   	github.com/opd-ai/voyage/pkg/world	[no test files]
```

---

## External Research Summary

- **No open GitHub issues** — The repository has 0 open issues as of 2026-04-05.
- **Ebitengine v2.9** — No known CVEs or security vulnerabilities. Current and maintained.
- **Dependencies** — All indirect dependencies (`golang.org/x/sync`, `golang.org/x/sys`, etc.) are standard Go modules with no known issues.

---

## Audit Methodology

1. **Phase 0**: Extracted 35+ claims from README.md and ROADMAP.md
2. **Phase 1**: Web search for GitHub issues and dependency vulnerabilities (none found)
3. **Phase 2**: `go-stats-generator analyze .` for baseline metrics
4. **Phase 3**: Package-by-package verification against stated goals
5. **Phase 4**: `go test -tags headless -race ./...` and `go vet -tags headless ./...`

---

*Generated by functional audit on 2026-04-05*
