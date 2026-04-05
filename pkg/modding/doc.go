// Package modding provides a JSON-based extension system for Voyage.
//
// This package allows users to add custom content without modifying game code:
//   - Custom event tables: Add new events for any genre
//   - Genre presets: Define new themes with custom names
//   - Mod loading: Load mods from JSON files at runtime
//
// # Mod File Format
//
// Mods are defined in JSON files with the following structure:
//
//	{
//	    "id": "my-mod",
//	    "name": "My Custom Mod",
//	    "version": "1.0.0",
//	    "author": "ModAuthor",
//	    "events": [...],
//	    "genres": [...]
//	}
//
// # Custom Events
//
// Events can be added to existing genres:
//
//	{
//	    "category": "encounter",
//	    "genre": "fantasy",
//	    "title": "Wandering Merchant",
//	    "description": "A merchant approaches with exotic wares.",
//	    "choices": [
//	        {
//	            "text": "Trade with them",
//	            "outcome": {"currency_delta": -10, "morale_delta": 5}
//	        },
//	        {
//	            "text": "Decline politely",
//	            "outcome": {}
//	        }
//	    ]
//	}
//
// # Custom Genres
//
// New genres can be defined with custom vocabulary:
//
//	{
//	    "id": "steampunk",
//	    "name": "Steampunk",
//	    "biomes": ["Clockwork City", "Aether Plains"],
//	    "resources": ["Coal", "Steam Cores"],
//	    "factions": ["Gear Guild", "Aether Pirates"]
//	}
package modding
