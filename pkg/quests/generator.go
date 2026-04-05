package quests

import (
	"fmt"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// Generator creates procedural quests.
type Generator struct {
	gen    *seed.Generator
	genre  engine.GenreID
	nextID QuestID
}

// NewGenerator creates a new quest generator.
func NewGenerator(masterSeed int64, genre engine.GenreID) *Generator {
	return &Generator{
		gen:    seed.NewGenerator(masterSeed, "quests"),
		genre:  genre,
		nextID: 1,
	}
}

// SetGenre updates the generator's genre.
func (g *Generator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// GenerateQuestBoard creates 2-4 quests for a supply point.
func (g *Generator) GenerateQuestBoard(supplyX, supplyY, mapWidth, mapHeight int) []*Quest {
	count := 2 + g.gen.Intn(3) // 2-4 quests
	quests := make([]*Quest, 0, count)

	for i := 0; i < count; i++ {
		qType := seed.Choice(g.gen, AllQuestTypes())
		quest := g.generateQuest(qType, supplyX, supplyY, mapWidth, mapHeight)
		quests = append(quests, quest)
	}

	return quests
}

// generateQuest creates a single quest of the specified type.
func (g *Generator) generateQuest(qType QuestType, originX, originY, mapWidth, mapHeight int) *Quest {
	title := g.generateTitle(qType)
	description := g.generateDescription(qType)

	quest := NewQuest(g.nextID, qType, title, description, g.genre)
	g.nextID++

	quest.OriginX = originX
	quest.OriginY = originY
	quest.GiverName = g.generateGiverName()
	quest.TimeLimit = g.generateTimeLimit(qType)
	quest.Reward = g.generateReward(qType)

	g.addObjectives(quest, qType, originX, originY, mapWidth, mapHeight)

	return quest
}

// generateTitle creates a quest title.
func (g *Generator) generateTitle(qType QuestType) string {
	titles := questTitles[g.genre]
	if titles == nil {
		titles = questTitles[engine.GenreFantasy]
	}
	typeTitles := titles[qType]
	if len(typeTitles) == 0 {
		return "Unknown Quest"
	}
	return seed.Choice(g.gen, typeTitles)
}

var questTitles = map[engine.GenreID]map[QuestType][]string{
	engine.GenreFantasy: {
		TypeDelivery: {
			"The Royal Dispatch",
			"Merchant's Request",
			"Urgent Scroll Delivery",
			"Package for the Guild",
		},
		TypeRescue: {
			"Lost in the Wilderness",
			"Kidnapped by Bandits",
			"Stranded Travelers",
			"Missing Pilgrims",
		},
		TypeRetrieve: {
			"The Lost Heirloom",
			"Stolen Artifact",
			"Ancient Relic",
			"The Precious Gem",
		},
		TypeExplore: {
			"Map the Unknown",
			"Survey the Ruins",
			"Scout Ahead",
			"The Hidden Path",
		},
		TypeEliminate: {
			"Beast Hunt",
			"Bandit Trouble",
			"Monster in the Hills",
			"Clear the Road",
		},
	},
	engine.GenreScifi: {
		TypeDelivery: {
			"Priority Cargo",
			"Data Transfer",
			"Medical Supplies Run",
			"Corporate Package",
		},
		TypeRescue: {
			"Distress Signal",
			"Stranded Crew",
			"Escape Pod Recovery",
			"Colony Evacuation",
		},
		TypeRetrieve: {
			"Salvage Operation",
			"Lost Probe Recovery",
			"Stolen Tech",
			"Derelict Artifacts",
		},
		TypeExplore: {
			"Survey Mission",
			"System Scan",
			"Anomaly Investigation",
			"Uncharted Sector",
		},
		TypeEliminate: {
			"Pirate Elimination",
			"Hostile Neutralization",
			"Threat Assessment",
			"Clear the Lane",
		},
	},
	engine.GenreHorror: {
		TypeDelivery: {
			"Medical Run",
			"Supply Drop",
			"Food Delivery",
			"Emergency Supplies",
		},
		TypeRescue: {
			"Survivor Signal",
			"Trapped Civilians",
			"Lost Patrol",
			"Family Reunion",
		},
		TypeRetrieve: {
			"Vital Equipment",
			"The Cure Components",
			"Weapons Cache",
			"Research Data",
		},
		TypeExplore: {
			"Area Recon",
			"Safe Route",
			"The Dead Zone",
			"Unknown Territory",
		},
		TypeEliminate: {
			"Clear the Nest",
			"Horde Thinning",
			"Raider Problem",
			"The Pack Leader",
		},
	},
	engine.GenreCyberpunk: {
		TypeDelivery: {
			"Data Courier",
			"Package Run",
			"Hot Delivery",
			"Priority Transfer",
		},
		TypeRescue: {
			"Extraction Op",
			"Asset Recovery",
			"Hostage Situation",
			"Burned Agent",
		},
		TypeRetrieve: {
			"Acquisition Job",
			"Tech Heist",
			"Data Grab",
			"The Prototype",
		},
		TypeExplore: {
			"Recon Work",
			"Intel Gathering",
			"Area Surveillance",
			"Net Dive",
		},
		TypeEliminate: {
			"Contract Work",
			"Gang Problem",
			"Corporate Cleanup",
			"Wetwork",
		},
	},
	engine.GenrePostapoc: {
		TypeDelivery: {
			"Water Run",
			"Supply Caravan",
			"Medicine Delivery",
			"Food Drop",
		},
		TypeRescue: {
			"Survivor Rescue",
			"Captured Settlers",
			"Stranded Scout",
			"Lost Child",
		},
		TypeRetrieve: {
			"Scavenger Hunt",
			"Pre-War Tech",
			"The Water Chip",
			"Vital Parts",
		},
		TypeExplore: {
			"Wasteland Survey",
			"Safe Passage",
			"Radiation Mapping",
			"The Old Road",
		},
		TypeEliminate: {
			"Raider Camp",
			"Mutant Nest",
			"Bandit Problem",
			"Clear the Way",
		},
	},
}

// generateDescription creates a quest description.
func (g *Generator) generateDescription(qType QuestType) string {
	descriptions := questDescriptions[g.genre]
	if descriptions == nil {
		descriptions = questDescriptions[engine.GenreFantasy]
	}
	typeDescs := descriptions[qType]
	if len(typeDescs) == 0 {
		return "Complete this task for a reward."
	}
	return seed.Choice(g.gen, typeDescs)
}

var questDescriptions = map[engine.GenreID]map[QuestType][]string{
	engine.GenreFantasy: {
		TypeDelivery: {
			"A merchant needs urgent goods delivered to a distant settlement.",
			"Royal correspondence must reach its destination without delay.",
			"Precious cargo awaits transport to waiting customers.",
		},
		TypeRescue: {
			"Travelers have gone missing on a dangerous route. Find them.",
			"A family was taken by bandits. Bring them home safely.",
			"Lost pilgrims need guidance back to civilization.",
		},
		TypeRetrieve: {
			"A family heirloom was stolen. Retrieve it from the thieves.",
			"An ancient artifact lies in dangerous ruins. Recover it.",
			"A valuable relic must be returned to its rightful owners.",
		},
		TypeExplore: {
			"Chart an unknown path through treacherous territory.",
			"Survey old ruins for signs of danger or treasure.",
			"Scout the road ahead and report what you find.",
		},
		TypeEliminate: {
			"A dangerous beast terrorizes travelers. Hunt it down.",
			"Bandits have made the roads unsafe. Deal with them.",
			"A monster nests too close to civilization. Remove the threat.",
		},
	},
	engine.GenreScifi: {
		TypeDelivery: {
			"High-priority cargo requires secure transport between stations.",
			"Sensitive data must be physically transferred. No network trust.",
			"Medical supplies are urgently needed at a remote colony.",
		},
		TypeRescue: {
			"A distress signal has been received. Investigate and rescue.",
			"A crew is stranded after their ship was disabled. Extract them.",
			"Escape pods were jettisoned during an attack. Recover survivors.",
		},
		TypeRetrieve: {
			"A derelict ship contains valuable salvage. Secure it.",
			"A stolen prototype must be recovered before it's analyzed.",
			"A lost probe contains critical data. Find and retrieve it.",
		},
		TypeExplore: {
			"An uncharted sector requires mapping for navigation data.",
			"An anomaly has been detected. Investigate and report.",
			"Survey a new system for colonization potential.",
		},
		TypeEliminate: {
			"Pirates have been attacking shipping lanes. Neutralize them.",
			"A hostile force threatens a key route. Clear it.",
			"A threat to colonial security must be eliminated.",
		},
	},
	engine.GenreHorror: {
		TypeDelivery: {
			"A settlement desperately needs medical supplies. Get them there.",
			"Food is running low. A supply run could save lives.",
			"Emergency equipment must reach survivors before it's too late.",
		},
		TypeRescue: {
			"A distress signal was received. Someone needs help.",
			"Civilians are trapped in a dangerous zone. Extract them.",
			"A patrol never returned. Find out what happened.",
		},
		TypeRetrieve: {
			"Vital equipment was lost during an attack. Find it.",
			"Research data could help understand the threat. Retrieve it.",
			"A weapons cache could make all the difference. Secure it.",
		},
		TypeExplore: {
			"We need to know if this route is safe. Scout ahead.",
			"Something is in that area. Find out what.",
			"Map a path through the dead zone if possible.",
		},
		TypeEliminate: {
			"A nest is too close for comfort. Clear it out.",
			"The horde is growing. We need to thin their numbers.",
			"Raiders are preying on survivors. Stop them.",
		},
	},
	engine.GenreCyberpunk: {
		TypeDelivery: {
			"Hot data needs to move, and the net isn't safe. Carry it.",
			"A package needs to get there, no questions asked.",
			"Priority transfer, double rate. Interested?",
		},
		TypeRescue: {
			"An asset got burned. We need them extracted, fast.",
			"Hostage situation in corpo territory. Get them out.",
			"Someone important is being held. Make them not held.",
		},
		TypeRetrieve: {
			"A prototype went missing. Acquire it.",
			"Data was stolen from us. Steal it back.",
			"Something valuable is poorly guarded. Liberate it.",
		},
		TypeExplore: {
			"We need eyes on a location. Recon work.",
			"Intel gathering. Observe and report.",
			"Net dive required. Find what's hidden.",
		},
		TypeEliminate: {
			"Someone needs to stop breathing. Contract work.",
			"A gang moved into our territory. Fix that.",
			"Corporate cleanup. Make problems disappear.",
		},
	},
	engine.GenrePostapoc: {
		TypeDelivery: {
			"Clean water is worth its weight in gold. Deliver it safely.",
			"A settlement needs supplies or they won't last the winter.",
			"Medicine could save lives. Get it where it needs to go.",
		},
		TypeRescue: {
			"Survivors were spotted in hostile territory. Bring them back.",
			"Raiders took prisoners. Time to rescue them.",
			"A scout went missing. Find out if they're still alive.",
		},
		TypeRetrieve: {
			"Pre-war tech was spotted in a ruin. Salvage it.",
			"We need parts to fix the water purifier. Find them.",
			"An old cache might have what we need. Retrieve it.",
		},
		TypeExplore: {
			"The wasteland is dangerous. Map a safe route.",
			"Radiation levels need to be checked. Survey the area.",
			"We don't know what's out there. Scout and report.",
		},
		TypeEliminate: {
			"A raider camp is too close. Clear them out.",
			"Mutants are breeding nearby. Deal with it.",
			"Bandits are hitting caravans. Make them stop.",
		},
	},
}

// generateGiverName creates a procedural quest giver name.
func (g *Generator) generateGiverName() string {
	titles := giverTitles[g.genre]
	if titles == nil {
		titles = giverTitles[engine.GenreFantasy]
	}
	names := giverNames[g.genre]
	if names == nil {
		names = giverNames[engine.GenreFantasy]
	}

	if g.gen.Chance(0.4) {
		title := seed.Choice(g.gen, titles)
		name := seed.Choice(g.gen, names)
		return fmt.Sprintf("%s %s", title, name)
	}
	return seed.Choice(g.gen, names)
}

var giverTitles = map[engine.GenreID][]string{
	engine.GenreFantasy:   {"Lord", "Lady", "Master", "Elder", "Captain"},
	engine.GenreScifi:     {"Commander", "Director", "Chief", "Admiral", "Agent"},
	engine.GenreHorror:    {"Doctor", "Chief", "Leader", "Captain", "Elder"},
	engine.GenreCyberpunk: {"Mr.", "Ms.", "Boss", "Fixer", "Agent"},
	engine.GenrePostapoc:  {"Chief", "Elder", "Boss", "Leader", "Captain"},
}

var giverNames = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"Aldric", "Brynn", "Cedric", "Dara", "Elara",
		"Finn", "Gwen", "Hector", "Ivy", "Jareth",
	},
	engine.GenreScifi: {
		"Chen", "Rodriguez", "Okonkwo", "Petrov", "Singh",
		"Nakamura", "Schmidt", "Al-Rashid", "O'Brien", "Kim",
	},
	engine.GenreHorror: {
		"Miller", "Johnson", "Williams", "Brown", "Davis",
		"Garcia", "Wilson", "Anderson", "Thomas", "Lee",
	},
	engine.GenreCyberpunk: {
		"V", "Rogue", "Spider", "Ghost", "Razor",
		"Chrome", "Null", "Cipher", "Delta", "Smoke",
	},
	engine.GenrePostapoc: {
		"Dust", "Ash", "Crow", "Flint", "Sage",
		"Iron", "Wire", "Rust", "Blaze", "Grit",
	},
}

// generateTimeLimit creates a time limit based on quest type.
func (g *Generator) generateTimeLimit(qType QuestType) int {
	baseLimits := map[QuestType]int{
		TypeDelivery:  20 + g.gen.Intn(10),
		TypeRescue:    15 + g.gen.Intn(10),
		TypeRetrieve:  30 + g.gen.Intn(15),
		TypeExplore:   25 + g.gen.Intn(10),
		TypeEliminate: 20 + g.gen.Intn(10),
	}
	if limit, ok := baseLimits[qType]; ok {
		return limit
	}
	return 20
}

// generateReward creates quest rewards.
func (g *Generator) generateReward(qType QuestType) QuestReward {
	baseRewards := questBaseRewards[qType]
	if baseRewards.Currency == 0 {
		baseRewards = QuestReward{Currency: 50}
	}

	// Apply random variance
	variance := 0.8 + g.gen.Float64()*0.4
	return QuestReward{
		Currency:   baseRewards.Currency * variance,
		Food:       baseRewards.Food * variance,
		Water:      baseRewards.Water * variance,
		Fuel:       baseRewards.Fuel * variance,
		Medicine:   baseRewards.Medicine * variance,
		Morale:     baseRewards.Morale * variance,
		Reputation: baseRewards.Reputation,
	}
}

var questBaseRewards = map[QuestType]QuestReward{
	TypeDelivery: {
		Currency: 100,
		Morale:   5,
	},
	TypeRescue: {
		Currency:   75,
		Morale:     10,
		Reputation: 5,
	},
	TypeRetrieve: {
		Currency: 150,
		Morale:   5,
	},
	TypeExplore: {
		Currency: 50,
		Food:     10,
		Fuel:     10,
	},
	TypeEliminate: {
		Currency:   125,
		Morale:     5,
		Reputation: 3,
	},
}

// addObjectives adds appropriate objectives based on quest type.
func (g *Generator) addObjectives(quest *Quest, qType QuestType, originX, originY, mapWidth, mapHeight int) {
	targetX, targetY := g.generateTargetPosition(originX, originY, mapWidth, mapHeight)
	targetName := g.generateTargetName(qType)

	switch qType {
	case TypeDelivery:
		quest.AddObjective(
			fmt.Sprintf("Deliver package to %s", targetName),
			targetX, targetY, targetName,
		)
	case TypeRescue:
		quest.AddObjective(
			fmt.Sprintf("Find survivors at %s", targetName),
			targetX, targetY, targetName,
		)
		quest.AddObjective(
			"Return to safety",
			originX, originY, "Origin",
		)
	case TypeRetrieve:
		quest.AddObjective(
			fmt.Sprintf("Locate the artifact at %s", targetName),
			targetX, targetY, targetName,
		)
		quest.AddObjective(
			"Return with the artifact",
			originX, originY, "Origin",
		)
	case TypeExplore:
		quest.AddObjective(
			fmt.Sprintf("Survey %s", targetName),
			targetX, targetY, targetName,
		)
	case TypeEliminate:
		quest.AddObjective(
			fmt.Sprintf("Eliminate threat at %s", targetName),
			targetX, targetY, targetName,
		)
	}
}

// generateTargetPosition creates a valid target position.
func (g *Generator) generateTargetPosition(originX, originY, mapWidth, mapHeight int) (int, int) {
	// Generate a target 10-30 tiles away from origin
	distance := 10 + g.gen.Intn(21)
	angle := g.gen.Float64() * 6.28318 // 2*PI

	dx := int(float64(distance) * cos(angle))
	dy := int(float64(distance) * sin(angle))

	targetX := originX + dx
	targetY := originY + dy

	// Clamp to map bounds
	if targetX < 1 {
		targetX = 1
	}
	if targetX >= mapWidth-1 {
		targetX = mapWidth - 2
	}
	if targetY < 1 {
		targetY = 1
	}
	if targetY >= mapHeight-1 {
		targetY = mapHeight - 2
	}

	return targetX, targetY
}

// cos approximation
func cos(x float64) float64 {
	// Simple Taylor series approximation
	x2 := x * x
	return 1 - x2/2 + x2*x2/24
}

// sin approximation
func sin(x float64) float64 {
	x2 := x * x
	return x - x*x2/6 + x*x2*x2/120
}

// generateTargetName creates a procedural target location name.
func (g *Generator) generateTargetName(qType QuestType) string {
	locations := targetLocations[g.genre]
	if locations == nil {
		locations = targetLocations[engine.GenreFantasy]
	}
	typeLocations := locations[qType]
	if len(typeLocations) == 0 {
		return "the target location"
	}
	return seed.Choice(g.gen, typeLocations)
}

var targetLocations = map[engine.GenreID]map[QuestType][]string{
	engine.GenreFantasy: {
		TypeDelivery:  {"Riverside Market", "Mountain Lodge", "Forest Waypoint", "Coastal Village"},
		TypeRescue:    {"the Dark Woods", "Bandit Cave", "Mountain Pass", "Ruined Tower"},
		TypeRetrieve:  {"Ancient Ruins", "Forgotten Temple", "Dragon's Lair", "Cursed Tomb"},
		TypeExplore:   {"Uncharted Forest", "Mountain Peak", "Desert Oasis", "Hidden Valley"},
		TypeEliminate: {"Goblin Warren", "Bandit Camp", "Monster Lair", "Cursed Grounds"},
	},
	engine.GenreScifi: {
		TypeDelivery:  {"Station Alpha", "Asteroid Outpost", "Colony Hub", "Orbital Platform"},
		TypeRescue:    {"Disabled Vessel", "Crashed Ship", "Remote Station", "Escape Coordinates"},
		TypeRetrieve:  {"Derelict Hulk", "Abandoned Lab", "Crashed Probe", "Anomaly Zone"},
		TypeExplore:   {"Uncharted System", "Anomaly Sector", "New Planet", "Asteroid Field"},
		TypeEliminate: {"Pirate Base", "Hostile Station", "Threat Vector", "Enemy Territory"},
	},
	engine.GenreHorror: {
		TypeDelivery:  {"Safe House", "Survivor Camp", "Medical Station", "Supply Depot"},
		TypeRescue:    {"Overrun Building", "Quarantine Zone", "The Dead Zone", "Collapsed Structure"},
		TypeRetrieve:  {"Hospital Ruins", "Research Lab", "Military Bunker", "Abandoned Store"},
		TypeExplore:   {"The Dead Sector", "Unknown Territory", "Contaminated Zone", "Silent Streets"},
		TypeEliminate: {"Infected Nest", "Raider Camp", "Horde Gathering", "The Den"},
	},
	engine.GenreCyberpunk: {
		TypeDelivery:  {"Drop Point", "Secure Facility", "Corporate Tower", "Underground Hub"},
		TypeRescue:    {"Black Site", "Gang Territory", "Corporate Prison", "The Extraction Point"},
		TypeRetrieve:  {"Research Lab", "Data Center", "Secure Storage", "The Vault"},
		TypeExplore:   {"Restricted Zone", "Net Node", "Corporate Sector", "Underground"},
		TypeEliminate: {"Gang Hideout", "Corporate Office", "Black Site", "The Target Location"},
	},
	engine.GenrePostapoc: {
		TypeDelivery:  {"Water Station", "Trading Post", "Survivor Camp", "Fortified Town"},
		TypeRescue:    {"Raider Territory", "The Wastes", "Collapsed Building", "Danger Zone"},
		TypeRetrieve:  {"Pre-War Ruins", "Vault Entrance", "Military Base", "Factory Ruins"},
		TypeExplore:   {"Radiation Zone", "Unknown Territory", "The Badlands", "Abandoned City"},
		TypeEliminate: {"Raider Camp", "Mutant Nest", "Bandit Hideout", "The Stronghold"},
	},
}

// GeneratePrimaryObjective creates the main "reach destination" quest.
func (g *Generator) GeneratePrimaryObjective(destX, destY int, destName string) *Quest {
	title := primaryTitles[g.genre]
	if title == "" {
		title = "Journey's End"
	}

	desc := primaryDescriptions[g.genre]
	if desc == "" {
		desc = "Reach your destination safely."
	}

	quest := NewQuest(0, TypeDelivery, title, desc, g.genre) // ID 0 for primary
	quest.Status = StatusActive                              // Always active
	quest.TimeLimit = 0                                      // No time limit
	quest.AddObjective(
		fmt.Sprintf("Reach %s", destName),
		destX, destY, destName,
	)

	// Primary quest has no standard rewards - winning is the reward
	quest.Reward = QuestReward{}

	return quest
}

var primaryTitles = map[engine.GenreID]string{
	engine.GenreFantasy:   "The Long Journey",
	engine.GenreScifi:     "Final Destination",
	engine.GenreHorror:    "Find Safety",
	engine.GenreCyberpunk: "The Big Score",
	engine.GenrePostapoc:  "The Promised Land",
}

var primaryDescriptions = map[engine.GenreID]string{
	engine.GenreFantasy:   "Complete your journey to the legendary city.",
	engine.GenreScifi:     "Navigate through space to reach your destination.",
	engine.GenreHorror:    "Survive the nightmare and reach the safe zone.",
	engine.GenreCyberpunk: "Make it to the free zone and start a new life.",
	engine.GenrePostapoc:  "Cross the wastes to find the promised settlement.",
}
