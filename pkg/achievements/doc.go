// Package achievements provides milestone tracking and achievement system for the voyage game.
//
// This package implements 20+ achievements that can be earned per run, displayed
// on end screens and in the Hall of Records, with procedurally generated descriptions.
//
// # Overview
//
// The achievement system tracks player milestones across a journey run and
// awards achievements based on performance, exploration, and survival. All
// achievement descriptions are generated from seed and genre context.
//
// # Core Components
//
//   - Achievement: Single milestone that can be earned
//   - AchievementTracker: Tracks progress toward all achievements
//   - AchievementCategory: Groups achievements by type (survival, trade, exploration, etc.)
//   - RunStatistics: Collected stats used to determine achievement eligibility
//
// # Features
//
//   - 20+ milestones tracked per run
//   - Categories: Survival, Trade, Exploration, Combat, Social, Special
//   - End-screen display of earned achievements
//   - Hall of Records integration with meta-progression
//   - All descriptions generated from seed/genre context
//
// # Genre Support
//
// All components implement engine.GenreSwitcher for genre-aware generation:
//
//   - Fantasy: Heroic deeds, legendary quests, noble titles
//   - Sci-Fi: Mission commendations, exploration logs, crew honors
//   - Horror: Survival milestones, escape records, sanity preservation
//   - Cyberpunk: Street cred, rep scores, legendary runs
//   - Post-Apocalyptic: Survival days, wasteland records, community building
//
// # Usage
//
//	g := achievements.NewGenerator(seed, engine.GenreFantasy)
//	tracker := g.GenerateAchievementTracker()
//
//	// Update stats during gameplay
//	tracker.Stats.DaysSurvived++
//
//	// Check for newly earned achievements
//	earned := tracker.CheckAchievements()
//	for _, a := range earned {
//	    // Display achievement notification
//	}
//
// # Achievement Categories
//
//   - Survival: Days survived, crew retention, health management
//   - Trade: Gold earned, trades completed, profit margins
//   - Exploration: Distance traveled, discoveries made, regions visited
//   - Combat: Enemies defeated, battles won without losses
//   - Social: Faction relations, NPCs helped, reputation earned
//   - Special: Unique challenges, perfect runs, genre-specific feats
package achievements
