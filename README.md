# Voyage

**100% Procedural Travel Simulator** — inspired by Oregon Trail, FTL, and Organ Trail

[![CI](https://github.com/opd-ai/voyage/actions/workflows/ci.yml/badge.svg)](https://github.com/opd-ai/voyage/actions/workflows/ci.yml)

## Overview

Voyage is a rogue-like travel simulator where every map, event, crew, vessel, audio, and narrative is procedurally generated from a single seed. The game supports five genre themes:

- **Fantasy** — Enchanted realm / Silk Road setting
- **Sci-fi** — Deep space / star lanes setting
- **Horror** — Zombie apocalypse wasteland
- **Cyberpunk** — Megacity sprawl
- **Post-apocalyptic** — Irradiated dust-bowl wastes

**All gameplay assets are generated at runtime** — no bundled images, audio, or pre-written content.

## Status

🚧 **Early Development** — Core engine and ECS framework implemented. Full gameplay coming soon.

### Implemented
- ✅ Go module and project structure
- ✅ ECS framework with GenreSwitcher interface
- ✅ Seed-based deterministic RNG
- ✅ Ebitengine rendering foundation
- ✅ Genre-switchable palettes
- ✅ CI/CD pipeline

### In Progress
- 🔄 Procedural world map generation
- 🔄 Resource management system
- 🔄 Crew/party system
- 🔄 Event system

## Prerequisites

- Go 1.22 or later
- OpenGL support (for Ebitengine rendering)

On Linux, you may need:
```bash
sudo apt-get install libgl1-mesa-dev xorg-dev
```

## Installation

```bash
go install github.com/opd-ai/voyage/cmd/voyage@latest
```

Or build from source:
```bash
git clone https://github.com/opd-ai/voyage.git
cd voyage
go build ./cmd/voyage
```

## Usage

```bash
# Start with random seed
./voyage

# Start with specific seed (for reproducible runs)
./voyage --seed 12345

# Start with specific genre
./voyage --genre scifi

# Show version
./voyage --version

# Show help
./voyage --help
```

## Controls

| Key | Action |
|-----|--------|
| Arrow keys | Move vessel |
| Enter/Space | Select/Confirm |
| Escape | Pause/Menu |
| 1-4 | Select event choices |
| F3 | Toggle debug info |

## Development

```bash
# Run tests
go test -race ./...

# Run benchmarks
go test -bench=. ./pkg/benchmark/...

# Check for issues
go vet ./...

# Validate no bundled assets
./scripts/validate-no-assets.sh
```

## Project Structure

```
cmd/voyage/         # Main entry point
pkg/
  engine/           # ECS framework with GenreSwitcher
  procgen/seed/     # Deterministic RNG
  procgen/world/    # World map generation
  rendering/        # Ebitengine rendering
  resources/        # Resource management
  crew/             # Party/crew system
  vessel/           # Transport system
  events/           # Event system
  audio/            # Audio synthesis
  ux/               # UI/HUD/Menus
  game/             # Game loop and state
  config/           # Configuration
  saveload/         # Save/load system
scripts/            # Validation and utility scripts
```

## Core Design Pillars

1. **Resource Attrition** — Food, water, fuel, medicine, morale, and currency deplete over time
2. **Party/Crew Mortality** — Procedurally generated crew members who can sicken, die, or desert
3. **Vessel Integrity** — Transport accumulates damage and requires repair
4. **Procedural Event Stream** — Grammar-based text generation with branching choices
5. **Route Choice with Consequence** — Multiple paths with varying risk/reward
6. **Fully Procedural World** — Every element generated from the master seed

## License

MIT License — see [LICENSE](LICENSE) for details.

## See Also

- [ROADMAP.md](ROADMAP.md) — Detailed feature planning and milestones
- [CONTRIBUTING.md](CONTRIBUTING.md) — Development guidelines
