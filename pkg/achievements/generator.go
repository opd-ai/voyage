package achievements

import (
	"fmt"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// Generator creates procedural achievements with genre-appropriate text
type Generator struct {
	gen   *seed.Generator
	genre engine.GenreID
}

// NewGenerator creates an achievement generator with the given seed and genre
func NewGenerator(masterSeed int64, genre engine.GenreID) *Generator {
	return &Generator{
		gen:   seed.NewGenerator(masterSeed, "achievements"),
		genre: genre,
	}
}

// SetGenre updates the generator's active genre
func (g *Generator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// GenerateAchievementTracker creates a complete tracker with 20+ achievements
func (g *Generator) GenerateAchievementTracker() *AchievementTracker {
	tracker := NewAchievementTracker(g.genre)

	// Survival achievements (5)
	tracker.AddAchievement(g.generateSurvivalAchievement("survive_10", TierBronze, 10))
	tracker.AddAchievement(g.generateSurvivalAchievement("survive_30", TierSilver, 30))
	tracker.AddAchievement(g.generateSurvivalAchievement("survive_100", TierGold, 100))
	tracker.AddAchievement(g.generateNoLossAchievement())
	tracker.AddAchievement(g.generateFullCrewAchievement())

	// Trade achievements (4)
	tracker.AddAchievement(g.generateTradeAchievement("trader", TierBronze, 10))
	tracker.AddAchievement(g.generateTradeAchievement("merchant", TierSilver, 50))
	tracker.AddAchievement(g.generateTycoonAchievement())
	tracker.AddAchievement(g.generateRegionalTraderAchievement())

	// Exploration achievements (5)
	tracker.AddAchievement(g.generateExplorerAchievement())
	tracker.AddAchievement(g.generateCartographerAchievement())
	tracker.AddAchievement(g.generateDiscovererAchievement())
	tracker.AddAchievement(g.generateLoreKeeperAchievement())
	tracker.AddAchievement(g.generateSecretFinderAchievement())

	// Combat achievements (3)
	tracker.AddAchievement(g.generateWarriorAchievement())
	tracker.AddAchievement(g.generateChampionAchievement())
	tracker.AddAchievement(g.generateFlawlessAchievement())

	// Social achievements (4)
	tracker.AddAchievement(g.generateDiplomatAchievement())
	tracker.AddAchievement(g.generateHelperAchievement())
	tracker.AddAchievement(g.generateQuestorAchievement())
	tracker.AddAchievement(g.generateRecruiterAchievement())

	// Special achievements (3)
	tracker.AddAchievement(g.generatePerfectRunAchievement())
	tracker.AddAchievement(g.generateCloseCallsAchievement())
	tracker.AddAchievement(g.generateCriticalMasterAchievement())

	return tracker
}

// Survival achievement generators

func (g *Generator) generateSurvivalAchievement(id string, tier AchievementTier, days int) *Achievement {
	names := map[engine.GenreID]map[int]string{
		engine.GenreFantasy: {
			10:  "First Steps",
			30:  "Seasoned Traveler",
			100: "Legendary Journey",
		},
		engine.GenreScifi: {
			10:  "Space Legs",
			30:  "Veteran Spacer",
			100: "Deep Space Pioneer",
		},
		engine.GenreHorror: {
			10:  "Still Breathing",
			30:  "Hardened Survivor",
			100: "Against All Odds",
		},
		engine.GenreCyberpunk: {
			10:  "Street Survivor",
			30:  "Chrome Veteran",
			100: "Living Legend",
		},
		engine.GenrePostapoc: {
			10:  "Not Dead Yet",
			30:  "Wasteland Hardened",
			100: "Last One Standing",
		},
	}

	descriptions := map[engine.GenreID]map[int]string{
		engine.GenreFantasy: {
			10:  "Survive the first 10 days of your journey through the realm.",
			30:  "Complete 30 days of travel across the lands.",
			100: "Endure 100 days of adventure and emerge victorious.",
		},
		engine.GenreScifi: {
			10:  "Survive 10 days in the void of space.",
			30:  "Complete 30 days of interstellar travel.",
			100: "Endure 100 days exploring the final frontier.",
		},
		engine.GenreHorror: {
			10:  "Survive the nightmare for 10 days.",
			30:  "Last 30 days against the horrors.",
			100: "Defy death for 100 days in the darkness.",
		},
		engine.GenreCyberpunk: {
			10:  "Survive 10 days in the neon jungle.",
			30:  "Last 30 days running the shadows.",
			100: "Outlive the system for 100 days.",
		},
		engine.GenrePostapoc: {
			10:  "Survive 10 days in the wasteland.",
			30:  "Endure 30 days after the fall.",
			100: "Persist for 100 days when the world ended.",
		},
	}

	name := names[g.genre][days]
	desc := descriptions[g.genre][days]

	return NewAchievement(id, name, desc, CategorySurvival, tier, days, g.genre)
}

func (g *Generator) generateNoLossAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Guardian Angel",
		engine.GenreScifi:     "Zero Casualties",
		engine.GenreHorror:    "No One Left Behind",
		engine.GenreCyberpunk: "Clean Run",
		engine.GenrePostapoc:  "Family Keeper",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Go 10 days without losing a single companion to harm.",
		engine.GenreScifi:     "Complete 10 days without any crew casualties.",
		engine.GenreHorror:    "Keep everyone alive for 10 consecutive days.",
		engine.GenreCyberpunk: "Run 10 days without losing a single teammate.",
		engine.GenrePostapoc:  "Protect your group for 10 days without a single loss.",
	}

	return NewAchievement("no_losses", names[g.genre], descriptions[g.genre], CategorySurvival, TierSilver, 10, g.genre)
}

func (g *Generator) generateFullCrewAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Heroes All",
		engine.GenreScifi:     "Full Complement",
		engine.GenreHorror:    "We All Made It",
		engine.GenreCyberpunk: "Crew Intact",
		engine.GenrePostapoc:  "Whole Family",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Complete the journey with your entire party surviving.",
		engine.GenreScifi:     "Arrive at destination with full crew complement.",
		engine.GenreHorror:    "Finish with everyone who started still alive.",
		engine.GenreCyberpunk: "End the run with your whole team intact.",
		engine.GenrePostapoc:  "Reach safety with everyone who started the journey.",
	}

	return NewAchievement("full_crew", names[g.genre], descriptions[g.genre], CategorySurvival, TierGold, 1, g.genre)
}

// Trade achievement generators

func (g *Generator) generateTradeAchievement(id string, tier AchievementTier, count int) *Achievement {
	names := map[engine.GenreID]map[int]string{
		engine.GenreFantasy: {
			10: "Budding Merchant",
			50: "Master Trader",
		},
		engine.GenreScifi: {
			10: "Licensed Trader",
			50: "Commerce Expert",
		},
		engine.GenreHorror: {
			10: "Resourceful",
			50: "Survival Expert",
		},
		engine.GenreCyberpunk: {
			10: "Small Timer",
			50: "Fixer",
		},
		engine.GenrePostapoc: {
			10: "Barterer",
			50: "Trade Baron",
		},
	}

	descriptions := map[engine.GenreID]map[int]string{
		engine.GenreFantasy: {
			10: "Complete 10 trades at marketplaces.",
			50: "Complete 50 trades across the realm.",
		},
		engine.GenreScifi: {
			10: "Complete 10 cargo transactions.",
			50: "Complete 50 interstellar trades.",
		},
		engine.GenreHorror: {
			10: "Trade 10 times with other survivors.",
			50: "Complete 50 desperate bargains.",
		},
		engine.GenreCyberpunk: {
			10: "Complete 10 deals in the shadows.",
			50: "Close 50 transactions on the street.",
		},
		engine.GenrePostapoc: {
			10: "Barter 10 times in the wasteland.",
			50: "Complete 50 trades in the ruins.",
		},
	}

	name := names[g.genre][count]
	desc := descriptions[g.genre][count]

	return NewAchievement(id, name, desc, CategoryTrade, tier, count, g.genre)
}

func (g *Generator) generateTycoonAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Golden Touch",
		engine.GenreScifi:     "Stellar Profits",
		engine.GenreHorror:    "Survival Capitalist",
		engine.GenreCyberpunk: "Corporate Level",
		engine.GenrePostapoc:  "Wasteland Tycoon",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Accumulate 10,000 gold through trading.",
		engine.GenreScifi:     "Earn 10,000 credits from commerce.",
		engine.GenreHorror:    "Gather 10,000 worth of supplies through trade.",
		engine.GenreCyberpunk: "Net 10,000 eurodollars from deals.",
		engine.GenrePostapoc:  "Accumulate 10,000 in trade value.",
	}

	return NewAchievement("tycoon", names[g.genre], descriptions[g.genre], CategoryTrade, TierGold, 10000, g.genre)
}

func (g *Generator) generateRegionalTraderAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Realm Connector",
		engine.GenreScifi:     "Sector Trader",
		engine.GenreHorror:    "Network Builder",
		engine.GenreCyberpunk: "District Dealer",
		engine.GenrePostapoc:  "Territory Trader",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Trade in 5 different regions of the realm.",
		engine.GenreScifi:     "Conduct commerce in 5 different sectors.",
		engine.GenreHorror:    "Establish trade with 5 different survivor groups.",
		engine.GenreCyberpunk: "Deal in 5 different districts.",
		engine.GenrePostapoc:  "Trade in 5 different wasteland territories.",
	}

	return NewAchievement("regional_trader", names[g.genre], descriptions[g.genre], CategoryTrade, TierSilver, 5, g.genre)
}

// Exploration achievement generators

func (g *Generator) generateExplorerAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Far Wanderer",
		engine.GenreScifi:     "Long Hauler",
		engine.GenreHorror:    "Distance Survivor",
		engine.GenreCyberpunk: "Urban Nomad",
		engine.GenrePostapoc:  "Road Warrior",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Travel a total of 1,000 leagues.",
		engine.GenreScifi:     "Cover 1,000 parsecs of space.",
		engine.GenreHorror:    "Traverse 1,000 miles through the nightmare.",
		engine.GenreCyberpunk: "Navigate 1,000 blocks of the megacity.",
		engine.GenrePostapoc:  "Cross 1,000 miles of wasteland.",
	}

	return NewAchievement("explorer", names[g.genre], descriptions[g.genre], CategoryExploration, TierSilver, 1000, g.genre)
}

func (g *Generator) generateCartographerAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Realm Mapper",
		engine.GenreScifi:     "Star Charter",
		engine.GenreHorror:    "Territory Mapper",
		engine.GenreCyberpunk: "Grid Navigator",
		engine.GenrePostapoc:  "Wasteland Mapper",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Visit 10 different regions of the world.",
		engine.GenreScifi:     "Chart 10 different star systems.",
		engine.GenreHorror:    "Map 10 different zones.",
		engine.GenreCyberpunk: "Navigate 10 different sectors.",
		engine.GenrePostapoc:  "Explore 10 different wasteland regions.",
	}

	return NewAchievement("cartographer", names[g.genre], descriptions[g.genre], CategoryExploration, TierSilver, 10, g.genre)
}

func (g *Generator) generateDiscovererAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Relic Hunter",
		engine.GenreScifi:     "Anomaly Finder",
		engine.GenreHorror:    "Clue Finder",
		engine.GenreCyberpunk: "Data Miner",
		engine.GenrePostapoc:  "Scavenger Elite",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Make 20 discoveries during your journey.",
		engine.GenreScifi:     "Log 20 new discoveries in the ship's computer.",
		engine.GenreHorror:    "Find 20 clues about what happened.",
		engine.GenreCyberpunk: "Uncover 20 pieces of hidden data.",
		engine.GenrePostapoc:  "Discover 20 useful finds in the ruins.",
	}

	return NewAchievement("discoverer", names[g.genre], descriptions[g.genre], CategoryExploration, TierSilver, 20, g.genre)
}

func (g *Generator) generateLoreKeeperAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Lore Master",
		engine.GenreScifi:     "Data Archivist",
		engine.GenreHorror:    "Truth Seeker",
		engine.GenreCyberpunk: "Info Broker",
		engine.GenrePostapoc:  "History Keeper",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Collect 30 pieces of lore and legend.",
		engine.GenreScifi:     "Archive 30 data entries about the universe.",
		engine.GenreHorror:    "Document 30 pieces of the dark truth.",
		engine.GenreCyberpunk: "Accumulate 30 pieces of valuable information.",
		engine.GenrePostapoc:  "Preserve 30 records of the old world.",
	}

	return NewAchievement("lore_keeper", names[g.genre], descriptions[g.genre], CategoryExploration, TierGold, 30, g.genre)
}

func (g *Generator) generateSecretFinderAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Secret Seeker",
		engine.GenreScifi:     "Hidden Protocol",
		engine.GenreHorror:    "Dark Secrets",
		engine.GenreCyberpunk: "Deep Net Diver",
		engine.GenrePostapoc:  "Hidden Cache",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Uncover 10 hidden secrets.",
		engine.GenreScifi:     "Decrypt 10 hidden data caches.",
		engine.GenreHorror:    "Discover 10 terrible secrets.",
		engine.GenreCyberpunk: "Hack 10 hidden systems.",
		engine.GenrePostapoc:  "Find 10 hidden caches.",
	}

	a := NewAchievement("secret_finder", names[g.genre], descriptions[g.genre], CategoryExploration, TierGold, 10, g.genre)
	a.Hidden = true
	return a
}

// Combat achievement generators

func (g *Generator) generateWarriorAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Monster Slayer",
		engine.GenreScifi:     "Combat Veteran",
		engine.GenreHorror:    "Creature Killer",
		engine.GenreCyberpunk: "High Body Count",
		engine.GenrePostapoc:  "Raider's Bane",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Defeat 50 enemies in combat.",
		engine.GenreScifi:     "Neutralize 50 hostile contacts.",
		engine.GenreHorror:    "Put down 50 creatures.",
		engine.GenreCyberpunk: "Flatline 50 enemies.",
		engine.GenrePostapoc:  "Eliminate 50 threats.",
	}

	return NewAchievement("warrior", names[g.genre], descriptions[g.genre], CategoryCombat, TierSilver, 50, g.genre)
}

func (g *Generator) generateChampionAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Battle Champion",
		engine.GenreScifi:     "Tactical Genius",
		engine.GenreHorror:    "Survival Fighter",
		engine.GenreCyberpunk: "Street Legend",
		engine.GenrePostapoc:  "Wasteland Champion",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Win 20 battles against hostile forces.",
		engine.GenreScifi:     "Achieve victory in 20 combat engagements.",
		engine.GenreHorror:    "Survive 20 hostile encounters.",
		engine.GenreCyberpunk: "Win 20 firefights.",
		engine.GenrePostapoc:  "Triumph in 20 wasteland battles.",
	}

	return NewAchievement("champion", names[g.genre], descriptions[g.genre], CategoryCombat, TierGold, 20, g.genre)
}

func (g *Generator) generateFlawlessAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Untouchable",
		engine.GenreScifi:     "Perfect Engagement",
		engine.GenreHorror:    "Unscathed",
		engine.GenreCyberpunk: "Clean Sweep",
		engine.GenrePostapoc:  "Flawless Victor",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Win 5 battles without any party member taking damage.",
		engine.GenreScifi:     "Complete 5 combat missions with zero casualties.",
		engine.GenreHorror:    "Survive 5 encounters without anyone getting hurt.",
		engine.GenreCyberpunk: "Execute 5 jobs without taking a single hit.",
		engine.GenrePostapoc:  "Win 5 fights without anyone getting wounded.",
	}

	return NewAchievement("flawless", names[g.genre], descriptions[g.genre], CategoryCombat, TierLegendary, 5, g.genre)
}

// Social achievement generators

func (g *Generator) generateDiplomatAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Realm Diplomat",
		engine.GenreScifi:     "Ambassador",
		engine.GenreHorror:    "Alliance Builder",
		engine.GenreCyberpunk: "Connected",
		engine.GenrePostapoc:  "Peacemaker",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Become allies with 3 different factions.",
		engine.GenreScifi:     "Establish positive relations with 3 factions.",
		engine.GenreHorror:    "Form alliances with 3 survivor groups.",
		engine.GenreCyberpunk: "Get in good with 3 different crews.",
		engine.GenrePostapoc:  "Ally with 3 different communities.",
	}

	return NewAchievement("diplomat", names[g.genre], descriptions[g.genre], CategorySocial, TierSilver, 3, g.genre)
}

func (g *Generator) generateHelperAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Helpful Hero",
		engine.GenreScifi:     "Good Samaritan",
		engine.GenreHorror:    "Lifesaver",
		engine.GenreCyberpunk: "Street Angel",
		engine.GenrePostapoc:  "Hope Giver",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Help 15 NPCs in need.",
		engine.GenreScifi:     "Assist 15 people during your voyage.",
		engine.GenreHorror:    "Save 15 people from danger.",
		engine.GenreCyberpunk: "Help out 15 people on the street.",
		engine.GenrePostapoc:  "Aid 15 survivors in the wasteland.",
	}

	return NewAchievement("helper", names[g.genre], descriptions[g.genre], CategorySocial, TierSilver, 15, g.genre)
}

func (g *Generator) generateQuestorAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Quest Champion",
		engine.GenreScifi:     "Mission Master",
		engine.GenreHorror:    "Objective Complete",
		engine.GenreCyberpunk: "Job Done",
		engine.GenrePostapoc:  "Task Master",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Complete 10 quests during your journey.",
		engine.GenreScifi:     "Successfully complete 10 missions.",
		engine.GenreHorror:    "Accomplish 10 survival objectives.",
		engine.GenreCyberpunk: "Finish 10 jobs.",
		engine.GenrePostapoc:  "Complete 10 tasks for the community.",
	}

	return NewAchievement("questor", names[g.genre], descriptions[g.genre], CategorySocial, TierGold, 10, g.genre)
}

func (g *Generator) generateRecruiterAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Band Builder",
		engine.GenreScifi:     "Crew Recruiter",
		engine.GenreHorror:    "Group Builder",
		engine.GenreCyberpunk: "Team Builder",
		engine.GenrePostapoc:  "Community Leader",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Recruit 5 companions to your party.",
		engine.GenreScifi:     "Add 5 crew members to your roster.",
		engine.GenreHorror:    "Bring 5 survivors into your group.",
		engine.GenreCyberpunk: "Recruit 5 runners to your team.",
		engine.GenrePostapoc:  "Welcome 5 survivors to your community.",
	}

	return NewAchievement("recruiter", names[g.genre], descriptions[g.genre], CategorySocial, TierSilver, 5, g.genre)
}

// Special achievement generators

func (g *Generator) generatePerfectRunAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Perfect Journey",
		engine.GenreScifi:     "Flawless Mission",
		engine.GenreHorror:    "Perfect Escape",
		engine.GenreCyberpunk: "Clean Run",
		engine.GenrePostapoc:  "Perfect Migration",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Complete a 30+ day journey with full party survival.",
		engine.GenreScifi:     "Finish a 30+ day mission with no crew losses.",
		engine.GenreHorror:    "Escape after 30+ days with everyone who started.",
		engine.GenreCyberpunk: "Complete a 30+ day run without losing anyone.",
		engine.GenrePostapoc:  "Survive 30+ days and reach safety with full group.",
	}

	a := NewAchievement("perfect_run", names[g.genre], descriptions[g.genre], CategorySpecial, TierLegendary, 1, g.genre)
	a.UnlockReward = fmt.Sprintf("Unlocks %s title", TierNameByGenre(TierLegendary, g.genre))
	return a
}

func (g *Generator) generateCloseCallsAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Fortune's Favor",
		engine.GenreScifi:     "Against All Odds",
		engine.GenreHorror:    "Death's Door",
		engine.GenreCyberpunk: "Nine Lives",
		engine.GenrePostapoc:  "Lucky Survivor",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Survive 5 near-death experiences.",
		engine.GenreScifi:     "Escape 5 situations with less than 10% survival odds.",
		engine.GenreHorror:    "Come back from 5 moments that should have killed you.",
		engine.GenreCyberpunk: "Flatline and recover 5 times.",
		engine.GenrePostapoc:  "Survive 5 situations that should have been fatal.",
	}

	a := NewAchievement("close_calls", names[g.genre], descriptions[g.genre], CategorySpecial, TierGold, 5, g.genre)
	a.Hidden = true
	return a
}

func (g *Generator) generateCriticalMasterAchievement() *Achievement {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Lucky Star",
		engine.GenreScifi:     "Probability Defier",
		engine.GenreHorror:    "Miracle Worker",
		engine.GenreCyberpunk: "Loaded Dice",
		engine.GenrePostapoc:  "Blessed",
	}

	descriptions := map[engine.GenreID]string{
		engine.GenreFantasy:   "Achieve 10 critical successes.",
		engine.GenreScifi:     "Log 10 improbable success events.",
		engine.GenreHorror:    "Experience 10 miraculous saves.",
		engine.GenreCyberpunk: "Hit 10 lucky breaks.",
		engine.GenrePostapoc:  "Get 10 blessed outcomes.",
	}

	return NewAchievement("critical_master", names[g.genre], descriptions[g.genre], CategorySpecial, TierSilver, 10, g.genre)
}
