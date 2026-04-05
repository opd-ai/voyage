// Package companions provides advanced crew companion specializations for the voyage game.
//
// This package implements named companions with special abilities, genre-specific
// skills, and personality-driven events with procedurally generated dialogue.
//
// # Overview
//
// The companion system extends the basic crew system with named characters who
// have unique abilities, backstories, and personality traits. Companions can
// unlock special abilities at high skill levels and trigger unique events
// based on their personality.
//
// # Core Components
//
//   - Companion: Named crew member with abilities, personality, and backstory
//   - Ability: Special skill that companions can use during the journey
//   - PersonalityTrait: Character trait that affects events and dialogue
//   - CompanionManager: Tracks and manages all active companions
//
// # Features
//
//   - Named companions with procedurally generated names and backstories
//   - Special abilities unlocked at high skill levels
//   - Genre-skinned specializations (wizard, AI navigator, zombie handler, etc.)
//   - Personality-driven special events with procedural dialogue
//   - Relationship tracking between companions and the player
//
// # Genre Support
//
// All components implement engine.GenreSwitcher for genre-aware generation:
//
//   - Fantasy: Wizard Guide, Ranger Scout, Healer, Warrior, Bard
//   - Sci-Fi: AI Navigator, Engineer, Medic, Pilot, Science Officer
//   - Horror: Zombie Handler, Scout, Field Medic, Survivor, Occultist
//   - Cyberpunk: Netrunner, Techie, Street Samurai, Fixer, Ripperdoc
//   - Post-Apocalyptic: Rad-Doc, Mechanic, Scavenger, Scout, Leader
//
// # Usage
//
//	g := companions.NewGenerator(seed, engine.GenreFantasy)
//	companion := g.GenerateCompanion(companions.RoleGuide)
//
//	// Check if ability is available
//	if companion.CanUseAbility() {
//	    ability := companion.GetAbility()
//	    // Apply ability effect
//	}
//
// # Ability Unlocking
//
// Companions gain abilities through skill development:
//
//   - Base skill starts at 1-3 depending on background
//   - Skills increase through use and events
//   - Special abilities unlock at skill level 5+
//   - Master abilities unlock at skill level 8+
package companions
