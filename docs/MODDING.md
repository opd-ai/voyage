# Modding Guide

Voyage supports content mods via JSON files and WASM extensions. This guide covers both systems.

## Quick Start

1. Create a JSON file (e.g., `my-mod.json`) in the `mods/` directory
2. Define your mod with events, genres, or other content
3. Launch Voyage — mods load automatically

## JSON Mod Format

Mods are defined in JSON with the following structure:

```json
{
    "id": "my-mod-id",
    "name": "My Custom Mod",
    "version": "1.0.0",
    "author": "Your Name",
    "description": "Description of what this mod adds",
    "events": [...],
    "genres": [...],
    "biomes": [...],
    "resources": [...],
    "factions": [...]
}
```

### Required Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier (lowercase, no spaces) |
| `name` | string | Human-readable name |
| `version` | string | Semantic version (e.g., "1.0.0") |

### Optional Fields

| Field | Type | Description |
|-------|------|-------------|
| `author` | string | Mod creator's name |
| `description` | string | What the mod adds |
| `events` | array | Custom event definitions |
| `genres` | array | Custom genre definitions |
| `biomes` | array | Additional biomes for existing genres |
| `resources` | array | Custom resource types |
| `factions` | array | Custom faction definitions |

## Custom Events

Events are the primary way to add content. Each event has a category, narrative text, and choices.

### Event Categories

| Category | Description |
|----------|-------------|
| `weather` | Environmental conditions |
| `encounter` | Meeting NPCs or creatures |
| `discovery` | Finding locations or items |
| `hardship` | Equipment failures, shortages |
| `windfall` | Lucky finds, gifts |
| `hazard` | Dangerous terrain or situations |
| `crew` | Party member events |

### Event Definition

```json
{
    "category": "encounter",
    "genre": "fantasy",
    "title": "Wandering Merchant",
    "description": "A merchant approaches with exotic wares from distant lands.",
    "choices": [
        {
            "text": "Trade with them",
            "outcome": {
                "description": "You browse their goods and make a purchase.",
                "currency_delta": -10,
                "morale_delta": 5
            }
        },
        {
            "text": "Decline politely",
            "outcome": {
                "description": "The merchant tips their hat and continues on."
            }
        }
    ],
    "weight": 1.0,
    "min_turn": 0,
    "max_turn": 0,
    "requires_biome": []
}
```

### Event Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `category` | string | Yes | Event type (see categories above) |
| `genre` | string | No | Genre ID or empty for all genres |
| `title` | string | Yes | Event headline |
| `description` | string | Yes | Narrative text shown to player |
| `choices` | array | Yes | At least one choice required |
| `weight` | float | No | Selection frequency (default: 1.0) |
| `min_turn` | int | No | Earliest turn this can appear (default: 0) |
| `max_turn` | int | No | Latest turn (0 = no limit) |
| `requires_biome` | array | No | Limit to specific biome types |

### Choice Definition

```json
{
    "text": "Choice text shown to player",
    "outcome": {
        "description": "Result description",
        "food_delta": 10,
        "water_delta": 5,
        "fuel_delta": -10,
        "medicine_delta": 0,
        "currency_delta": -20,
        "morale_delta": 5,
        "crew_damage": 0,
        "vessel_damage": 10,
        "time_advance": 1
    },
    "require_skill": "navigator",
    "require_min_resource": {
        "currency": 20
    }
}
```

### Outcome Fields

All outcome fields are optional. Positive values indicate gains, negative indicate losses.

| Field | Type | Description |
|-------|------|-------------|
| `description` | string | Flavor text shown after selection |
| `food_delta` | float | Food change |
| `water_delta` | float | Water change |
| `fuel_delta` | float | Fuel/stamina change |
| `medicine_delta` | float | Medicine change |
| `currency_delta` | float | Currency change |
| `morale_delta` | float | Morale change (-100 to 100) |
| `crew_damage` | float | Damage to crew members (0-100) |
| `vessel_damage` | float | Damage to vessel (0-100, negative = repair) |
| `time_advance` | int | Turns to skip |

### Choice Conditions

| Field | Type | Description |
|-------|------|-------------|
| `require_skill` | string | Only show if crew has this skill |
| `require_min_resource` | object | Resource minimums to show choice |

## Custom Genres

Define entirely new genre themes with custom vocabulary:

```json
{
    "id": "steampunk",
    "name": "Steampunk",
    "description": "Victorian-era clockwork and steam technology.",
    "biomes": ["Clockwork Foundry", "Brass Wasteland", "Aether Heights"],
    "resources": ["Coal", "Steam Cores", "Brass Ingots"],
    "factions": ["Gear Guild", "Aether Pirates"],
    "vessel_types": ["Steam Crawler", "Aether Frigate"],
    "crew_roles": ["Gearwright", "Steam Engineer"],
    "category_names": {
        "weather": "atmospheric disturbance",
        "encounter": "clockwork confrontation",
        "discovery": "mechanical marvel"
    }
}
```

### Genre Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | Yes | Internal identifier |
| `name` | string | Yes | Display name |
| `description` | string | No | Theme description |
| `biomes` | array | No | Location type names |
| `resources` | array | No | Resource display names |
| `factions` | array | No | Faction names |
| `vessel_types` | array | No | Transport names |
| `crew_roles` | array | No | Crew job titles |
| `category_names` | object | No | Event category display names |

## Additional Biomes

Add biomes to existing genres:

```json
{
    "biomes": [
        {
            "genre": "fantasy",
            "names": ["Enchanted Grove", "Crystal Caverns"]
        }
    ]
}
```

## Custom Resources

Define new tradeable resources:

```json
{
    "id": "aether_crystals",
    "name": "Aether Crystals",
    "genre": "steampunk",
    "description": "Crystallized energy for powering devices.",
    "max_stack": 100,
    "tradeable": true
}
```

## Custom Factions

Add new factions:

```json
{
    "id": "gear_guild",
    "name": "Gear Guild",
    "genre": "steampunk",
    "description": "A union of engineers and mechanics.",
    "hostile": false
}
```

## WASM Mods (Advanced)

For complex logic, mods can be written in any language that compiles to WebAssembly.

### Capabilities

WASM mods run in a sandbox with explicit permissions:

| Capability | Description |
|------------|-------------|
| `read_events` | Read existing event data |
| `write_events` | Add new events |
| `read_genres` | Read genre configurations |
| `write_genres` | Add custom genres |
| `read_resources` | Read resource values |
| `write_resources` | Modify resource values |
| `read_crew` | Read crew member data |
| `modify_crew` | Modify crew members |
| `trigger_events` | Trigger custom events |
| `access_rng` | Access game's RNG |

### Required Exports

WASM mods must export these functions:

```c
// Required: Returns mod ID as (ptr, len)
int64_t mod_get_id(void);

// Optional: Returns mod name as (ptr, len)
int64_t mod_get_name(void);

// Optional: Returns version as (ptr, len)
int64_t mod_get_version(void);

// Optional: Returns author as (ptr, len)
int64_t mod_get_author(void);

// Optional: Called when mod is loaded
void mod_init(void);

// Optional: Called at start of each turn
void mod_on_turn_start(int32_t turn);

// Optional: Called when an event occurs
void mod_on_event(int32_t category_ptr, int32_t category_len);
```

### Host Functions

WASM mods can call these host functions:

```c
// Read events as JSON (returns bytes written)
uint32_t voyage_read_events(uint32_t ptr, uint32_t max_len);

// Add an event from JSON (returns 1 on success)
uint32_t voyage_add_event(uint32_t ptr, uint32_t len);

// Read genres as JSON
uint32_t voyage_read_genres(uint32_t ptr, uint32_t max_len);

// Add a genre from JSON
uint32_t voyage_add_genre(uint32_t ptr, uint32_t len);

// Read resources as JSON
uint32_t voyage_read_resources(uint32_t ptr, uint32_t max_len);

// Log a message (for debugging)
void voyage_log(uint32_t ptr, uint32_t len);

// Check if capability is granted (returns 1 or 0)
uint32_t voyage_has_capability(uint32_t cap);
```

## Examples

See `examples/mods/` for complete working examples:

- `steampunk-genre.json` — Full custom genre with events, biomes, resources, and factions

## Validation

Mods are validated on load. Common errors:

| Error | Cause |
|-------|-------|
| `missing id` | Required `id` field not set |
| `missing name` | Required `name` field not set |
| `missing version` | Required `version` field not set |
| `missing title` | Event has no `title` |
| `missing description` | Event has no `description` |
| `no choices` | Event has empty `choices` array |
| `invalid category` | Event `category` not in valid list |
| `choice text required` | A choice has empty `text` |

## Best Practices

1. **Use descriptive IDs** — `my-steampunk-events` not `mod1`
2. **Set appropriate weights** — Higher weight = more frequent selection
3. **Balance outcomes** — Powerful rewards should have risks
4. **Write immersive text** — Match the genre's tone and vocabulary
5. **Test thoroughly** — Load your mod and play through events
6. **Use `require_skill`** — Add meaningful choices for skilled crew
7. **Version semantically** — Increment versions when updating

## Loading Mods

Mods are loaded from the `mods/` directory in the game's data path:

- **Linux**: `~/.local/share/voyage/mods/`
- **macOS**: `~/Library/Application Support/voyage/mods/`
- **Windows**: `%APPDATA%/voyage/mods/`

Alternatively, specify a custom mods directory with the `--mods-dir` flag.

## API Reference

For full API details, see the `pkg/modding/` package documentation:

```go
import "github.com/opd-ai/voyage/pkg/modding"

// JSON mod loading
loader := modding.NewLoader()
mod, err := loader.LoadFromFile("my-mod.json")

// WASM mod loading
wasmLoader := modding.NewWASMLoader(modding.DefaultWASMConfig())
wasmMod, err := wasmLoader.LoadFromFile("my-mod.wasm", modding.MinimalCapabilities)
```
