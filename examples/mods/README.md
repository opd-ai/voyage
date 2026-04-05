# Example Mods

This directory contains example mods for Voyage demonstrating how to extend the game with custom content.

## Available Examples

### steampunk-genre.json

A complete custom genre preset adding a Steampunk theme to Voyage. This mod demonstrates all the extension points available in the JSON modding system:

- **Custom Genre**: Full "steampunk" genre with Victorian-era clockwork aesthetics
- **Biome Names**: Clockwork Foundry, Brass Wasteland, Aether Heights, etc.
- **Resource Names**: Coal, Steam Cores, Brass Ingots, Aether Crystals
- **Faction Archetypes**: Gear Guild, Aether Pirates, Clockwork Union, Steam Barons
- **Vessel Types**: Steam Crawler, Aether Frigate, Clockwork Carriage
- **Crew Roles**: Gearwright, Aether Navigator, Steam Engineer
- **Custom Events**: 6 genre-specific events with multiple choices

## Using Mods

### Loading JSON Mods

```go
import "github.com/opd-ai/voyage/pkg/modding"

// Create a loader
loader := modding.NewLoader()

// Load a mod from file
mod, err := loader.LoadFromFile("steampunk-genre.json")
if err != nil {
    log.Fatal(err)
}

// Access custom events
events := loader.GetEventsForGenre("steampunk")

// Access custom genres
genres := loader.GetCustomGenres()
```

### Loading WASM Mods

For mods that need executable code, use the WASM loader with capability-based security:

```go
import "github.com/opd-ai/voyage/pkg/modding"

// Create a WASM loader with security config
config := modding.DefaultWASMConfig()
config.Capabilities = modding.CapReadEvents | modding.CapWriteEvents

wasmLoader := modding.NewWASMLoader(config)

// Load a WASM mod
mod, err := wasmLoader.LoadFromFile("my-mod.wasm", config.Capabilities)
if err != nil {
    log.Fatal(err)
}

// Initialize the mod
ctx := context.Background()
if err := mod.Initialize(ctx); err != nil {
    log.Fatal(err)
}
```

## Creating Your Own Mod

### Mod File Structure

```json
{
    "id": "my-mod-id",
    "name": "My Mod Name",
    "version": "1.0.0",
    "author": "Your Name",
    "description": "What this mod does",
    "genres": [...],
    "events": [...],
    "biomes": [...],
    "resources": [...],
    "factions": [...]
}
```

### Event Structure

Events have a category, genre, title, description, and choices:

```json
{
    "category": "encounter",
    "genre": "fantasy",
    "title": "Event Title",
    "description": "What happens in this event",
    "choices": [
        {
            "text": "Choice text shown to player",
            "outcome": {
                "description": "Result text",
                "food_delta": -10,
                "morale_delta": 5
            }
        }
    ],
    "weight": 1.0
}
```

### Valid Categories

- `weather` - Environmental conditions
- `encounter` - Meeting NPCs or enemies
- `discovery` - Finding locations or items
- `hardship` - Challenges and setbacks
- `windfall` - Fortunate events
- `hazard` - Dangerous situations
- `crew` - Crew-related events

### Outcome Fields

| Field | Type | Description |
|-------|------|-------------|
| `description` | string | Text shown after choice |
| `food_delta` | float | Change to food supply |
| `water_delta` | float | Change to water supply |
| `fuel_delta` | float | Change to fuel supply |
| `medicine_delta` | float | Change to medicine |
| `currency_delta` | float | Change to currency |
| `morale_delta` | float | Change to crew morale |
| `crew_damage` | float | Damage to crew health |
| `vessel_damage` | float | Damage to vessel |
| `time_advance` | int | Turns to skip |

### Choice Requirements

Choices can have requirements:

```json
{
    "text": "Use medical expertise",
    "outcome": {...},
    "require_skill": "medic",
    "require_min_resource": {
        "medicine": 10
    }
}
```

## WASM Mod Development

For executable mods, compile to WASM and export these functions:

### Required Exports

```
mod_get_id() -> (ptr, len)      // Return mod ID string
```

### Optional Exports

```
mod_get_name() -> (ptr, len)    // Return mod name
mod_get_version() -> (ptr, len) // Return version string
mod_get_author() -> (ptr, len)  // Return author name
mod_init()                      // Called once on load
mod_on_turn_start(turn)         // Called each turn
mod_on_event(ptr, len)          // Called on events
```

### Host Functions Available

```
voyage_read_events(ptr, len) -> bytes_written
voyage_add_event(ptr, len) -> success
voyage_read_genres(ptr, len) -> bytes_written
voyage_add_genre(ptr, len) -> success
voyage_read_resources(ptr, len) -> bytes_written
voyage_log(ptr, len)
voyage_has_capability(cap) -> bool
```

### Capability System

WASM mods run with limited permissions:

| Capability | Value | Description |
|------------|-------|-------------|
| `CapReadEvents` | 1 | Read existing events |
| `CapWriteEvents` | 2 | Add new events |
| `CapReadGenres` | 4 | Read genre configs |
| `CapWriteGenres` | 8 | Add custom genres |
| `CapReadResources` | 16 | Read resource data |
| `CapWriteResources` | 32 | Modify resources |
| `CapReadCrew` | 64 | Read crew data |
| `CapModifyCrew` | 128 | Modify crew |
| `CapTriggerEvents` | 256 | Trigger events |
| `CapAccessRNG` | 512 | Access game RNG |

## License

Example mods are provided under the MIT License.
