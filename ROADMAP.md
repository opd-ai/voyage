# VOYAGE — Procedural Travel Simulator

**Gameplay style**: Rogue-like travel simulator — top-down 2D overworld navigation with turn-based time progression, systemic resource management, procedural event resolution, and crew/vessel stewardship across a procedurally generated world.

**Inspirations**: Oregon Trail (resource attrition, journey events, permadeath party members, branching decisions), FTL: Faster Than Light (vessel system damage, crew roles, tactical encounter pausing, multi-room management), Organ Trail (survival horror tone, dark humor, vehicle maintenance, desperate choices).

**Vision**: Ship an infinitely replayable, fully procedural travel simulator — every map, route, encounter, character, crew member, event, piece of lore, and piece of dialogue generated from a single seed — achieving feature parity with the reference-complete **venture** roguelike across all five setting genres. **All gameplay assets — including audio, visual, and narrative/story-driven components — must be procedurally generated at runtime using deterministic algorithms. No pre-rendered images (`.png`, `.jpg`, `.svg`, `.gif`), bundled audio files (`.mp3`, `.wav`, `.ogg`), or static narrative content (hardcoded dialogue, pre-written event scripts, fixed story arcs, embedded text assets) are permitted in the project.**

---

## Genre Support

Every system must implement the `GenreSwitcher` interface to switch thematic presentation at runtime. **This interface is not yet implemented** — creating it is the first task in v1.0 ECS Framework.

```go
type GenreID string

const (
    GenreIDFantasy   GenreID = "fantasy"
    GenreIDScifi     GenreID = "scifi"
    GenreIDHorror    GenreID = "horror"
    GenreIDCyberpunk GenreID = "cyberpunk"
    GenreIDPostapoc  GenreID = "postapoc"
)

type GenreSwitcher interface {
    SetGenre(genreID GenreID)
}
```

This applies to all ECS Systems (renderer, audio, AI, event-generator, HUD, narrative). Components and Entities hold genre-tagged data; only Systems are required to implement `SetGenre()`.

| Genre ID    | Setting                     | Vessel / Transport              | Destination                          | Signature Hazards                                           |
|-------------|-----------------------------|---------------------------------|--------------------------------------|-------------------------------------------------------------|
| `fantasy`   | Enchanted realm / Silk Road | Horse-drawn wagon caravan       | Legendary city / lost dungeon        | Magical storms, monster ambushes, cursed passes, famine     |
| `scifi`     | Deep space / star lanes     | FTL-capable spacecraft          | Distant space station / colony world | Asteroid fields, alien encounters, hull breaches, system failures |
| `horror`    | Zombie apocalypse wasteland | Armored vehicle convoy          | Fortified safe zone                  | Zombie hordes, raider gangs, disease outbreaks, fuel shortage |
| `cyberpunk` | Megacity sprawl / corridors | Armored runner vehicle / drone  | Corporate enclave / free zone        | Corporate checkpoints, netrunner ambushes, gang territory    |
| `postapoc`  | Irradiated dust-bowl wastes | Jury-rigged diesel transport    | New settlement / promised land       | Radiation storms, mutant packs, scavenger warbands           |

---

## Core Design Pillars

1. **Resource Attrition** — food, water, fuel/stamina, medicine, morale, and currency all deplete over time. Every decision has a resource cost.
2. **Party/Crew Mortality** — crew members are procedurally generated individuals with names, traits, skills, and health. They can sicken, die, desert, or grow as the journey progresses.
3. **Vessel Integrity** — the transport accumulates wear, can be damaged in encounters, and must be repaired with scavenged parts. A destroyed vessel ends the run.
4. **Procedural Event Stream** — the journey is driven by a continuous stream of seeded events: weather, encounters, discoveries, moral dilemmas, and crisis moments.
5. **Route Choice with Consequence** — the world map offers multiple branching paths with varying distance, terrain, hazard, and reward profiles. Faster routes are more dangerous.
6. **Fully Procedural World** — every map tile, landmark name, NPC, trade post, event text, lore entry, and piece of audio is generated at runtime from the master seed.

---

## Phased Milestones

### v1.0 — Core Engine + Playable Single-Player Journey

*Goal: ECS scaffold, seed-based PCG, rendering, and a fully playable journey from origin to destination in one genre (`fantasy` baseline).*

#### ECS Framework
- [ ] Component / Entity / System interfaces (`SetGenre(genreID GenreID)` required on every **System**; see interface definition above)
- [ ] System execution ordering and dependency graph
- [ ] Entity lifecycle management (spawn, despawn, pooling)

#### Seed-Based Deterministic RNG
- [ ] Master seed → subsystem seed derivation (`HashSeed` via SHA-256)
- [ ] Per-subsystem isolated `math/rand` sources
- [ ] Determinism test suite (same seed → same game)

#### Input System
- [ ] Keyboard / gamepad mapping
- [ ] Rebindable controls (stored in config)
- [ ] Modal input handling (overworld navigation vs. event resolution vs. menus)

#### Procedural World Map Generation
- [ ] Voronoi / grid-based overworld with regions and biomes
- [ ] Origin → destination placement with guaranteed solvable path
- [ ] Waypoint and landmark seeding (towns, outposts, ruins, wilderness)
- [ ] Terrain type assignment (plains, forest, mountain, desert, river, ocean, ruin)
- [ ] Branching path network with risk/reward tradeoffs (short=dangerous, long=safer)
- [ ] `SetGenre()` on map generator to swap biome vocabulary (forest→void nebula, mountain→asteroid belt)

#### Overworld Rendering
- [ ] Ebiten tile renderer for the world map
- [ ] Procedural tile sprite generation (cellular automata + palette)
- [ ] Fog-of-war / unexplored region masking
- [ ] Player vessel token (procedurally generated sprite)
- [ ] Landmark icons (procedurally generated per type and genre)
- [ ] `SetGenre()` on renderer to swap palette and tile-theme presets

#### Time Progression System
- [ ] Turn-based day/night cycle
- [ ] Movement costs per terrain type (mountains cost more turns than plains)
- [ ] Rest mechanic (spend turns stationary to recover morale/health)
- [ ] Seasonal time tracking (affects hazard frequency and resource costs)

#### Resource Management — Core Six
- [ ] Food (depletes daily; crew starves if empty)
- [ ] Water (depletes daily; faster in desert/hot biomes)
- [ ] Fuel / Stamina (depletes per movement; vessel stops if empty)
- [ ] Medicine (consumed on injury/disease events; death without it)
- [ ] Morale (falls on hardship, rises on rest/success; crew desert at zero)
- [ ] Currency / Trade Goods (used at supply points)
- [ ] Resource HUD with warning thresholds
- [ ] `SetGenre()` renames resources (food→rations→biomass→credit-chips→scrap)

#### Party / Crew System — Foundation
- [ ] Party entity with 2–6 crew member slots
- [ ] Procedurally generated crew member (name, portrait-sprite, trait, skill)
- [ ] Individual health tracking per crew member
- [ ] Crew mortality (starvation, disease, injury)
- [ ] `SetGenre()` re-skins crew names and portrait palette (medieval → alien → survivor → street-punk → wastelander)

#### Vessel / Transport System — Foundation
- [ ] Vessel entity with hull integrity, speed, and cargo capacity stats
- [ ] Cargo inventory (items the party carries)
- [ ] Basic breakdown events (random chance per turn based on vessel condition)
- [ ] Repair mechanic (spend materials to restore integrity)
- [ ] `SetGenre()` swaps vessel type vocabulary (wagon → spacecraft → car → runner-rig → diesel-hauler)

#### Vessel Customization — Foundation
- [ ] Procedurally generated vessel name (seed-derived; player may rename)
- [ ] Starting loadout selection (3 procedurally generated preset configurations: balanced, fast/light, slow/heavy)
- [ ] Visual variant selection (3 procedurally generated hull skins per genre; no bundled sprites)

#### Procedural Event System — Core
- [ ] Event queue seeded from master seed + current map position
- [ ] Event categories: weather, encounter, discovery, hardship, windfall
- [ ] Choice-based event resolution (present 2–4 options with different resource costs/gains)
- [ ] Outcome application (apply resource deltas, crew health changes, vessel damage)
- [ ] All event text procedurally generated at runtime from grammar templates driven by seed — no pre-authored event scripts
- [ ] `SetGenre()` re-skins event vocabulary and flavor text generation parameters

#### Audio — Waveform Synthesis & SFX
- [ ] Sine / square / sawtooth / triangle / noise waveforms
- [ ] ADSR envelope system
- [ ] SFX generation (travel movement, event fanfare, crisis alarm, success jingle, death toll)
- [ ] Ambient travel music (looping procedural composition)
- [ ] `SetGenre()` on audio to select thematic instrument/timbre presets

#### UI / HUD / Menus
- [ ] World map screen with vessel position, explored tiles, and route overlay
- [ ] Resource panel (six resources with bar indicators)
- [ ] Crew roster panel (names, health, morale per member)
- [ ] Event overlay (text, choices, outcome display)
- [ ] Main menu, pause menu, options screen
- [ ] Genre-themed UI skin switchable via `SetGenre()`

#### Save / Load
- [ ] Multiple save slots with autosave on turn advance
- [ ] Slot selection screen
- [ ] Seed embedded in save for reproducibility

#### Config / Settings
- [ ] Resolution, volume, key bindings persisted to disk
- [ ] CLI flags (`--seed`, `--genre`, `--difficulty`)

#### Win / Lose Conditions
- [ ] Win: vessel reaches destination tile with ≥1 living crew member
- [ ] Lose: vessel destroyed, or all crew dead, or morale hits zero and full mutiny
- [ ] End-screen with run summary (days traveled, crew lost, events survived, score)

#### Foraging and Scavenging
- [ ] Spend turns at wilderness, ruin, or landmark tiles to attempt a gather action
- [ ] Outcome table seeded from position + turn count (find food, find parts, find nothing, trigger encounter)
- [ ] Diminishing returns per tile (repeated foraging same tile yields less)
- [ ] `SetGenre()` re-skins gather action (forage → salvage → scavenge → jack data → strip wreck)

---

### v2.0 — Full Journey Loop (All 5 Genres, Crew Depth, Vessel Upgrades, Trading, Tactical Encounters)

*Goal: Complete the gameplay loop — deep crew management, vessel upgrades, trading economy, tactical encounter resolution, and all five genre skins.*

#### All 5 Genres — Full Integration
- [ ] `SetGenre()` implemented on every system (renderer, audio, map, event, HUD, narrative, crew, vessel)
- [ ] Genre selection at game start (or seed-derived genre) using `GenreID` constants
- [ ] Per-genre biome / tile / palette / SFX / music generation parameter presets (configuration values that drive procedural generation — not bundled asset files)
- [ ] Per-genre hazard vocabulary (magic storms, asteroid fields, zombie hordes, netrunner ambushes, radiation storms)

#### Crew Depth — Traits, Skills, and Relationships
- [ ] Trait system (brave, cautious, medic, mechanic, navigator, scavenger, etc.)
- [ ] Skill system (skills improve with use — experienced medic heals more effectively)
- [ ] Crew relationship network (pairs that bicker or bond affect morale events)
- [ ] Crew-specific events (personal crisis, milestone, sacrifice opportunity)
- [ ] Procedurally generated crew backstory surfaced in crew detail screen — no pre-written character bios
- [ ] `SetGenre()` re-skins trait/skill names (medic→biomancer→doc→netdoc→chem-doc)

#### Status Effects — Crew
- [ ] Disease (spreads between crew; slows recovery; fatal without medicine)
- [ ] Injury (reduces action effectiveness; requires rest + medicine)
- [ ] Exhaustion (from overtravel; reduces skill effectiveness)
- [ ] Despair (low morale debuff; increases desertion chance)
- [ ] Genre-specific afflictions (cursed → irradiated → infected → glitched → mutated)
- [ ] `SetGenre()` renames and recolours status icons

#### Vessel Systems — Depth and Upgrades
- [ ] Modular vessel system: engine, cargo hold, medical bay, navigation, defense
- [ ] Per-system integrity tracking (damaged navigation = harder routing)
- [ ] Upgrade system (spend currency at supply points to improve modules)
- [ ] Cargo management screen (weight/volume limits per cargo hold tier)
- [ ] Salvage mechanic (strip fallen vessels/wrecks for parts)
- [ ] `SetGenre()` re-skins vessel modules (stable→engine room→engine bay→core systems→reactor)

#### Vessel Customization — Full System
- [ ] Custom module loadout screen before departure (swap starting module tier per slot)
- [ ] Upgrade path branching: each module offers speed, cargo, or defense specialization tracks
- [ ] Vessel insignia / livery selection (procedurally generated emblems; no bundled art)
- [ ] Vessel insurance mechanic: pay currency to protect one module from a catastrophic breakdown
- [ ] `SetGenre()` re-skins customization screen and module upgrade vocabulary

#### Trading and Supply Points
- [ ] Procedurally generated supply posts / towns at waypoints
- [ ] Dynamic inventory seeded from map region and genre
- [ ] Buy/sell interface with supply/demand pricing
- [ ] All item names and descriptions procedurally generated from seed — no embedded item text
- [ ] Bartering option (trade goods instead of currency)
- [ ] Town reputation track (friendly towns offer better prices; hostile towns may attack)
- [ ] `SetGenre()` re-skins trading posts (market→space-dock→survivor-camp→black-market→scrap-bazaar)

#### Tactical Encounter Resolution
- [ ] Encounter types: ambush, negotiation, race/chase, crisis management, puzzle
- [ ] Pausable real-time or turn-based resolution phase (FTL-style)
- [ ] Crew assignment to encounter roles (fighter, medic, engineer, negotiator)
- [ ] Outcome branches: victory, partial success, retreat, defeat
- [ ] `SetGenre()` re-skins encounter imagery and sound design

#### Weather and Environmental Hazards
- [ ] Weather system: 8+ types (storm, blizzard, heatwave, flood, fog, meteor shower, dust storm, acid rain)
- [ ] Weather affects movement cost, resource consumption, visibility, and crew health
- [ ] Terrain hazards: mountain passes (injury risk), river crossings (fuel cost), desert (water crisis), ruin (random loot + danger)
- [ ] Genre-appropriate hazard subset per theme via `SetGenre()`

#### Procedural NPC Generation
- [ ] Wandering NPC encounters (traders, refugees, bandits, lost travelers)
- [ ] Faction-affiliated NPCs with alignment tags
- [ ] All NPC names, dialogue, and descriptions procedurally generated — no pre-authored NPC text
- [ ] `SetGenre()` re-skins NPC archetypes (merchant→trader→survivor→fixer→scavenger)

#### Destination Depth
- [ ] Multiple destination types seeded per run (city, sanctuary, treasure vault, escape craft, settlement)
- [ ] Destination discovery events as the party approaches
- [ ] Arrival ceremony sequence with procedurally generated narrative payoff text
- [ ] `SetGenre()` re-skins destination type and arrival text vocabulary

#### Crew Council
- [ ] Critical route-choice decisions (dangerous shortcut, costly detour) trigger a crew vote
- [ ] Each crew member votes based on their dominant trait (brave votes for risk, cautious votes against)
- [ ] Player may overrule the vote; doing so applies a morale penalty proportional to dissent
- [ ] Unanimous votes in the player's favor grant a small morale bonus
- [ ] `SetGenre()` re-skins the council scene (campfire debate → bridge briefing → group argument → exec meeting → bonfire council)

---

### v3.0 — Visual Polish (Lighting, Particles, Weather Visuals, Enhanced Sprites, Adaptive Audio)

*Goal: Make the procedurally generated world feel alive and distinct per genre.*

#### Dynamic Lighting
- [ ] Day/night cycle lighting on overworld (dawn→day→dusk→night transitions)
- [ ] Point lights at towns, campfires, vessel lanterns
- [ ] Darkness penalty at night (reduced visibility unless torches/power used)
- [ ] Genre presets via `SetGenre()` (warm campfire glow for `fantasy`, blue-white hull lights for `scifi`, emergency red for `horror`, neon spillover for `cyberpunk`, dim salvage lanterns for `postapoc`)

#### Particle Effects
- [ ] Movement trail (dust clouds, thruster exhaust, tire tracks)
- [ ] Weather particles (rain, snow, sand, embers, ash)
- [ ] Event flash effects (combat sparks, healing glow, disaster explosion)
- [ ] Genre-specific particle themes via `SetGenre()`

#### Enhanced Sprite Generation
- [ ] Animated overworld tiles (flowing water, wind-swept grass, flickering fires)
- [ ] Crew member portrait animation (idle breathing, hurt flinch, death fade)
- [ ] Vessel damage states (pristine → worn → damaged → critical sprites)
- [ ] Animated landmark icons (smoking ruins, blinking outpost lights)
- [ ] Genre palette overlays via `SetGenre()`

#### Music — Adaptive Multi-Layer
- [ ] Dynamic layer system (peaceful travel, crisis, encounter, victory, death)
- [ ] Biome-specific ambient music parameters (forest calm → asteroid tension → wastes drone)
- [ ] Smooth cross-fade between intensity states
- [ ] Genre instrument mapping via `SetGenre()` (lute/harp → synthesizer pad → distorted bass → glitch-synth → industrial grind) — all instruments procedurally synthesized, not sampled from bundled audio

#### Audio — Positional SFX
- [ ] Distance attenuation for offscreen events
- [ ] Left/right stereo panning for spatial awareness
- [ ] Ambient loop per biome/region (wind, space hum, groaning metal, city noise, silence)

All SFX and music are procedurally synthesized at runtime — no pre-recorded or bundled audio files.

#### Genre Post-Processing Presets
- [ ] `fantasy` — warm desaturated vignette, bloom on magic effects
- [ ] `scifi` — cool scanline overlay, chromatic aberration at screen edges
- [ ] `horror` — desaturate + red-tint at low crew health, film grain
- [ ] `cyberpunk` — neon bloom, CRT curvature, glitch artifacts on hacks
- [ ] `postapoc` — sepia wash, dust overlay, heavy vignette

#### Dynamic Minimap Overlay
- [ ] Always-visible procedurally rendered corner minimap showing explored tiles and current position
- [ ] Icons for towns, ruins, hazards, and the destination (revealed as explored)
- [ ] Minimap fades or dims in crisis events (damaged navigation module reduces fidelity)
- [ ] `SetGenre()` applies genre-appropriate minimap aesthetic (parchment map → holographic display → torn atlas → AR overlay → scratched road atlas)

---

### v4.0 — Depth Expansion (Factions, Quests, Meta-Progression, Advanced Narrative, Multi-Leg Journeys)

*Goal: Deepen the simulator loop — richer political landscape, procedural quest objectives, cross-run persistence, and narrative payoff.*

#### Faction System
- [ ] 4–6 procedurally generated factions per run (seeded names, ideologies, territory)
- [ ] Faction relationship matrix (allied, neutral, hostile) affected by player choices
- [ ] Faction-controlled territory blocks on the overworld (require safe passage or conflict)
- [ ] Reputation track per faction (favors, betrayals, bribes shift standing)
- [ ] Genre-mapped factions (guild/duchy/cult → corp/colony/pirate → gang/survivor-band/military remnant)

#### Quest / Objective System
- [ ] Primary objective: reach destination (always)
- [ ] Procedurally generated side quests: deliver parcel, rescue stranded crew, retrieve artifact
- [ ] Quest board at supply points (accept/decline optional missions)
- [ ] All quest text, objectives, and flavor procedurally generated from seed — no pre-authored quest scripts
- [ ] Objective tracker in HUD
- [ ] `SetGenre()` re-flavors quest vocabulary

#### Meta-Progression Between Runs
- [ ] Unlock log: persistent record of event types and destinations seen
- [ ] Unlockable starting crew archetypes (seeded by cumulative game-state hash, not random)
- [ ] Unlockable vessel starting configurations
- [ ] Hall of Records: best run summary per genre (days, crew survivors, score)

#### Advanced Narrative — Procedural Story Arc
- [ ] Three-act structure derived from seed (departure crisis → mid-journey revelation → arrival twist)
- [ ] Named recurring NPC that reappears across the journey (friend, nemesis, or ambiguous figure)
- [ ] Crew backstory events that surface mid-journey and connect to the destination
- [ ] All narrative text, character arcs, and story beats generated deterministically from seed — no pre-authored story content

#### Environmental Storytelling
- [ ] Procedurally generated world-map lore inscriptions (ruins with descriptions, grave markers, burnt signs)
- [ ] Abandoned vessel/camp discoveries with item inventories and procedurally generated vignette text
- [ ] All environmental text generated algorithmically — no pre-authored flavor text

#### Lore Codex
- [ ] In-game codex screen with discovered lore entries (world history, faction bios, route legends)
- [ ] All lore texts procedurally generated per genre from seed — no embedded text assets
- [ ] Unlock via exploration (ruins, events, NPC conversations)

#### Multi-Leg Journey Support
- [ ] Campaign mode: chain 2–4 journey legs with state persisting between legs
- [ ] Intermediate stopover city as hub between legs (buy, upgrade, recruit)
- [ ] Escalating difficulty per leg (longer distances, harsher terrain, stronger factions)
- [ ] `SetGenre()` applied per leg (option: genre shifts between legs for narrative variety)

#### Companion Specializations (Advanced Crew)
- [ ] Named companions unlock special abilities at high skill levels
- [ ] Companion ability: genre-skinned (wizard guide, AI navigator, zombie handler, netrunner, rad-doc)
- [ ] Companion-driven special events that depend on their personality and backstory — all dialogue procedurally generated

#### Achievement System
- [ ] 20+ milestones tracked per run (survived X days, traded in every region, lost no crew, etc.)
- [ ] Achievements displayed on end-screen and in main menu Hall of Records
- [ ] All achievement descriptions generated from seed/genre context

#### Trade Route Dynamics
- [ ] Regional supply and demand model: goods sold in a region become cheaper; scarcities drive prices up
- [ ] Demand shifts propagate along procedurally generated trade routes (selling food in one town raises its price there and lowers it along the supply chain)
- [ ] Price history display at supply points (sparkline of recent trade activity for that good)
- [ ] Speculation mechanic: buy cheap goods early in the journey to sell dear at the destination
- [ ] `SetGenre()` re-skins trade goods and economic vocabulary (grain/spices → fuel cells/ore → medical supplies/ammo → data chips/access codes → scrap/water)

---

### v5.0 — Online / Social / Platform Expansion

*Goal: Shared runs, leaderboards, async convoy play, WASM browser build, modding hooks.*

#### Per-Seed Leaderboards
- [ ] On run completion, submit (seed, genre, score, days, survivors, timestamp) to leaderboard
- [ ] Global leaderboard screen with filter by genre and seed
- [ ] Replay seed to attempt same world as top-score run

#### Async Convoy Mode
- [ ] Shared-seed co-op: multiple players run same world seed simultaneously
- [ ] Async event resolution: each player's run can diverge from shared events
- [ ] End-of-run comparison screen (who survived best, who reached destination first)
- [ ] High-latency tolerant design (Tor / onion-service friendly, 200–5000ms latency)

#### WebAssembly Build
- [ ] `make build-wasm` target with Ebitengine WASM output
- [ ] Browser-playable version deployed via GitHub Pages
- [ ] Touch input support for mobile browsers

#### Modding System
- [ ] JSON-based event grammar extension point (add custom event tables without code changes)
- [ ] WASM-sandboxed mod loader (capability-based security)
- [ ] Example mod: custom genre preset with new biome names, resource names, faction archetypes

#### Mobile Support
- [ ] Android APK build via `gomobile`
- [ ] iOS build via `gomobile`
- [ ] Touch controls (tap to move, swipe to scroll map, tap to select options)

#### Run Sharing
- [ ] Export any completed run as a compact shareable code (seed + genre + decision sequence encoded as base58 string)
- [ ] Import a run code to replay another player's exact route and choices as a ghost overlay
- [ ] Ghost mode: player plays the same seed while a translucent ghost vessel shows the shared run's path and timing
- [ ] Share codes copyable from the end-screen with a single button press

---

## Architecture Overview

```
cmd/voyage/          — Ebitengine game loop entry point
pkg/engine/          — ECS core (World, Entity, Component, System, GenreSwitcher)
pkg/procgen/         — Procedural generators
  world/             — Overworld map (Voronoi, terrain, landmark placement)
  event/             — Event stream (grammar-based text, choice trees, outcomes)
  crew/              — Crew member generation (name, trait, skill, portrait)
  vessel/            — Vessel stats, module layout, cargo generation
  npc/               — Wandering NPC generation
  narrative/         — Story arc, lore codex, environmental text
  genre/             — Genre preset registry and GenreSwitcher dispatch
pkg/rendering/       — Sprite generation, tile renderer, particle system, post-processing
pkg/audio/           — PCM synthesis, SFX generation, adaptive music layers
pkg/game/            — Game struct, state machine, system integration
pkg/resources/       — Resource management (six-axis attrition model)
pkg/crew/            — Crew entity management, status effects, relationships
pkg/vessel/          — Vessel entity management, module damage, cargo
pkg/events/          — Event queue, resolution, outcome application
pkg/world/           — Overworld state, fog-of-war, route management
pkg/trading/         — Supply point inventory, buy/sell, bartering
pkg/factions/        — Faction state, reputation, territory
pkg/saveload/        — Save slots, serialization, seed embedding
pkg/config/          — Configuration, CLI flags, settings persistence
pkg/ux/              — HUD, menus, event overlay, codex screen
pkg/benchmark/       — Performance benchmarks
pkg/audit/           — Feature audit tooling
```

---

## Procedural Generation Constraints

All content generators must satisfy:

1. **Determinism**: Given the same seed and genre, every run produces an identical world. Verified by determinism test suite (to be implemented in v1.0).
2. **No Bundled Assets**: Zero embedded images, audio files, or pre-authored text. Confirmed by `scripts/validate-no-assets.sh` (to be implemented in v1.0).
3. **Single Binary**: `go build ./cmd/voyage` produces one self-contained executable.
4. **GenreSwitcher Compliance**: Every ECS System implements `SetGenre(genreID GenreID)`. Verified by interface conformance tests (to be implemented in v1.0).
5. **Seed Isolation**: Each subsystem derives its own `math/rand` source from the master seed via `HashSeed`. Cross-subsystem seeding is prohibited.

---

## Success Criteria

| Milestone | Target | Measurement |
|-----------|--------|-------------|
| No bundled assets | 0 `.png`/`.mp3`/`.ogg`/hardcoded text | `scripts/validate-no-assets.sh` *(to be implemented)* |
| Deterministic runs | Same seed produces identical world | Determinism test suite *(to be implemented in v1.0)* |
| Single binary | `go build ./cmd/voyage` succeeds | CI build job |
| All 5 genres playable | Each genre passes smoke-test | `go test ./pkg/procgen/genre/...` |
| GenreSwitcher compliance | All Systems implement interface | Interface conformance test *(to be implemented in v1.0)* |
| Test coverage | ≥40% per package | `go test -cover ./pkg/...` |
| `go vet` clean | 0 errors | `go vet ./...` |
| 60 FPS on target hardware | Benchmark ≥60 FPS | `pkg/benchmark/fps/` *(to be implemented)* |
| `<500MB` client memory | Benchmark <500 MB heap | `pkg/benchmark/memory/` *(to be implemented)* |
