# Leaderboard System

Voyage includes an optional online leaderboard system for tracking and comparing high scores across players.

## Overview

The leaderboard client (`pkg/leaderboard/client.go`) supports:
- Score submission with automatic retry
- Offline caching when server is unavailable
- Querying by seed, genre, or global ranking
- Seed-based competitive play (everyone can race the same seed)

## Default Configuration

By default, the game uses a placeholder server URL:
```
https://api.voyage-game.example.com/leaderboard
```

**This URL is non-functional.** Online leaderboard features require deploying your own server.

## Self-Hosting Options

### Option 1: Simple REST Server

The leaderboard API expects these endpoints:

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/submit` | Submit a new entry (JSON body) |
| GET | `/query` | Query entries with filters |
| GET | `/health` | Health check endpoint |

#### Entry JSON Format

```json
{
  "player_name": "string",
  "seed": 12345,
  "genre": "fantasy",
  "score": 1000,
  "distance": 450,
  "turns": 100,
  "crew_survived": 3,
  "timestamp": "2024-01-01T12:00:00Z"
}
```

#### Query Parameters

- `seed` (int64): Filter by specific seed
- `genre` (string): Filter by genre (fantasy, scifi, horror, cyberpunk, postapoc)
- `limit` (int): Maximum entries to return
- `offset` (int): Pagination offset

### Option 2: Peer-to-Peer via Convoy

The convoy networking system (`pkg/convoy/`) can share leaderboard entries between connected players without a central server. This is currently experimental.

## Configuring Custom Server URL

To use a custom leaderboard server, modify the client configuration:

```go
import "github.com/opd-ai/voyage/pkg/leaderboard"

config := leaderboard.DefaultConfig()
config.ServerURL = "https://your-server.example.com/leaderboard"
client := leaderboard.NewClient(config)
```

## Offline Mode

When the server is unavailable, entries are stored locally in `~/.config/voyage/leaderboard.json` (or equivalent platform path). These entries sync automatically when connectivity is restored.

Local storage provides:
- Personal best tracking
- Seed replay recommendations
- Offline leaderboard viewing

## Running Without Online Features

The game works fully offline. Leaderboard features degrade gracefully:
- Score submission queues locally
- Queries return local data only
- No error dialogs or blocking behavior

## Example Server Implementation

A minimal Go+SQLite leaderboard server is planned for `cmd/leaderboard-server/`. Until then, any REST API implementing the endpoints above will work.

### Minimal Python Example

```python
from flask import Flask, request, jsonify
import sqlite3

app = Flask(__name__)

@app.route('/health')
def health():
    return jsonify({"status": "ok"})

@app.route('/submit', methods=['POST'])
def submit():
    data = request.json
    # Store in database...
    return jsonify({"success": True, "rank": 1})

@app.route('/query')
def query():
    seed = request.args.get('seed')
    genre = request.args.get('genre')
    limit = request.args.get('limit', 100, type=int)
    # Query database...
    return jsonify({"entries": []})

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)
```

## Security Considerations

- Validate all input on the server side
- Consider rate limiting submissions
- Use HTTPS in production
- Entries include timestamps for anti-cheat verification
- Seed verification can detect modified game states
