package npc

import (
	"fmt"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// Generator creates procedural NPCs.
type Generator struct {
	gen    *seed.Generator
	genre  engine.GenreID
	nextID int
}

// NewGenerator creates a new NPC generator.
func NewGenerator(masterSeed int64, genre engine.GenreID) *Generator {
	return &Generator{
		gen:    seed.NewGenerator(masterSeed, "npc"),
		genre:  genre,
		nextID: 1,
	}
}

// SetGenre updates the generator's genre.
func (g *Generator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// Generate creates a new NPC of the specified type.
func (g *Generator) Generate(npcType NPCType) *NPC {
	name := g.generateName(npcType)
	npc := NewNPC(g.nextID, name, npcType, g.genre)
	g.nextID++

	npc.Alignment = g.determineAlignment(npcType)
	npc.Description = g.generateDescription(npcType)
	npc.Dialogue = g.generateDialogue(npcType, npc.Alignment)

	if npcType == TypeTrader {
		npc.TradeGoods = g.generateTradeGoods()
	}

	return npc
}

// GenerateRandom creates a random NPC type.
func (g *Generator) GenerateRandom() *NPC {
	npcType := seed.Choice(g.gen, AllNPCTypes())
	return g.Generate(npcType)
}

// GenerateEncounter creates an NPC suitable for an encounter.
func (g *Generator) GenerateEncounter() *NPC {
	// Weight hostile types lower
	weights := []float64{0.25, 0.15, 0.15, 0.2, 0.15, 0.1}
	npcType := seed.WeightedChoice(g.gen, AllNPCTypes(), weights)
	return g.Generate(npcType)
}

func (g *Generator) generateName(npcType NPCType) string {
	firstNames := firstNamesByGenre[g.genre]
	if firstNames == nil {
		firstNames = firstNamesByGenre[engine.GenreFantasy]
	}

	titlesByType := titlesByGenreType[g.genre]
	if titlesByType == nil {
		titlesByType = titlesByGenreType[engine.GenreFantasy]
	}
	titles := titlesByType[npcType]

	firstName := seed.Choice(g.gen, firstNames)

	// Some NPCs get titles
	if len(titles) > 0 && g.gen.Float64() < 0.4 {
		title := seed.Choice(g.gen, titles)
		return fmt.Sprintf("%s %s", title, firstName)
	}

	return firstName
}

var firstNamesByGenre = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"Aldric", "Brynn", "Cedric", "Dara", "Elara", "Finn",
		"Gwen", "Hector", "Ivy", "Jareth", "Kira", "Liam",
		"Mira", "Nolan", "Orla", "Piers", "Quinn", "Rhys",
	},
	engine.GenreScifi: {
		"Astra", "Beck", "Cade", "Dex", "Echo", "Flux",
		"Gaia", "Hex", "Ion", "Jax", "Kira", "Luna",
		"Nova", "Orion", "Pulse", "Quasar", "Rex", "Siren",
	},
	engine.GenreHorror: {
		"Alex", "Blake", "Casey", "Dana", "Eli", "Frank",
		"Grace", "Hunter", "Isaac", "Jamie", "Kelly", "Lee",
		"Morgan", "Nick", "Olive", "Pat", "Quinn", "Riley",
	},
	engine.GenreCyberpunk: {
		"Blade", "Chrome", "Dash", "Edge", "Flash", "Ghost",
		"Hack", "Ice", "Jack", "Knife", "Link", "Mox",
		"Neon", "Oxide", "Pulse", "Razor", "Spike", "Volt",
	},
	engine.GenrePostapoc: {
		"Ash", "Blaze", "Crow", "Dust", "Echo", "Flint",
		"Grit", "Haze", "Iron", "Junk", "Knox", "Rust",
		"Sage", "Tank", "Vex", "Wire", "Zeke", "Thorn",
	},
}

var titlesByGenreType = map[engine.GenreID]map[NPCType][]string{
	engine.GenreFantasy: {
		TypeTrader:   {"Honest", "Traveling", "Master"},
		TypeRefugee:  {"Poor", "Lost", "Weary"},
		TypeBandit:   {"Black", "Scarred", "One-Eyed"},
		TypeTraveler: {"Wandering", "Curious", "Old"},
		TypeScout:    {"Swift", "Silent", "Far-Seeing"},
		TypeGuard:    {"Stern", "Vigilant", "Armored"},
	},
	engine.GenreScifi: {
		TypeTrader:   {"Licensed", "Void", "Deep-Space"},
		TypeRefugee:  {"Displaced", "Colonial", "Stranded"},
		TypeBandit:   {"Void", "Renegade", "Cutthroat"},
		TypeTraveler: {"Star", "Free", "Rim"},
		TypeScout:    {"Forward", "Long-Range", "Recon"},
		TypeGuard:    {"Station", "Corporate", "Military"},
	},
	engine.GenreHorror: {
		TypeTrader:   {"Cautious", "Armed", "Nervous"},
		TypeRefugee:  {"Terrified", "Bitten", "Desperate"},
		TypeBandit:   {"Ruthless", "Mad", "Hungry"},
		TypeTraveler: {"Lone", "Silent", "Paranoid"},
		TypeScout:    {"Sharp-Eyed", "Quick", "Wary"},
		TypeGuard:    {"Tired", "Grizzled", "Armed"},
	},
	engine.GenreCyberpunk: {
		TypeTrader:   {"Street", "Licensed", "Black-Market"},
		TypeRefugee:  {"Zeroed", "Burned", "Corporate"},
		TypeBandit:   {"Chrome", "Wired", "Boosted"},
		TypeTraveler: {"Data", "Grid", "Street"},
		TypeScout:    {"Net", "Ghost", "Ice"},
		TypeGuard:    {"Corpo", "Private", "Armored"},
	},
	engine.GenrePostapoc: {
		TypeTrader:   {"Water", "Wandering", "Vault"},
		TypeRefugee:  {"Irradiated", "Lost", "Desperate"},
		TypeBandit:   {"Blood", "War", "Cannibal"},
		TypeTraveler: {"Wasteland", "Road", "Desert"},
		TypeScout:    {"Recon", "Path", "Trail"},
		TypeGuard:    {"Gate", "Wall", "Perimeter"},
	},
}

func (g *Generator) determineAlignment(npcType NPCType) Alignment {
	defaultAlign := GetDefaultAlignment(npcType)

	// Add some variance
	roll := g.gen.Float64()
	switch defaultAlign {
	case AlignmentHostile:
		if roll < 0.15 {
			return AlignmentSuspicious // Sometimes bandits can be reasoned with
		}
	case AlignmentNeutral:
		if roll < 0.2 {
			return AlignmentFriendly
		} else if roll < 0.3 {
			return AlignmentSuspicious
		}
	case AlignmentFriendly:
		if roll < 0.1 {
			return AlignmentAllied
		}
	case AlignmentSuspicious:
		if roll < 0.2 {
			return AlignmentNeutral
		} else if roll < 0.1 {
			return AlignmentHostile
		}
	}

	return defaultAlign
}

func (g *Generator) generateDescription(npcType NPCType) string {
	descriptions := descriptionsByGenreType[g.genre]
	if descriptions == nil {
		descriptions = descriptionsByGenreType[engine.GenreFantasy]
	}
	typeDescs := descriptions[npcType]
	if len(typeDescs) == 0 {
		return "A stranger on the road."
	}
	return seed.Choice(g.gen, typeDescs)
}

var descriptionsByGenreType = map[engine.GenreID]map[NPCType][]string{
	engine.GenreFantasy: {
		TypeTrader: {
			"A weathered merchant with a cart full of goods.",
			"A traveling peddler with jingling packs.",
			"A shrewd-looking trader counting coins.",
		},
		TypeRefugee: {
			"A tired family fleeing some distant trouble.",
			"A lone wanderer with haunted eyes.",
			"A group of displaced villagers.",
		},
		TypeBandit: {
			"A rough-looking band blocking the road.",
			"Masked figures emerge from the shadows.",
			"Armed desperados demand your attention.",
		},
		TypeTraveler: {
			"A pilgrim heading to distant shrines.",
			"A wandering bard seeking stories.",
			"A curious explorer studying the land.",
		},
		TypeScout: {
			"A keen-eyed ranger surveys the path ahead.",
			"A skilled tracker offers knowledge of the route.",
			"A swift messenger pauses briefly.",
		},
		TypeGuard: {
			"Armored sentries block the way.",
			"Vigilant guards patrol the area.",
			"A checkpoint manned by stern soldiers.",
		},
	},
	engine.GenreScifi: {
		TypeTrader: {
			"A ship-board merchant with cargo pods.",
			"A licensed vendor with holographic wares.",
			"A void trader offering exotic goods.",
		},
		TypeRefugee: {
			"Colonists fleeing a dead world.",
			"A stranded crew seeking passage.",
			"Evacuees from a station disaster.",
		},
		TypeBandit: {
			"Pirates hail on an intercept course.",
			"Renegade raiders arm their weapons.",
			"Void wolves demand tribute.",
		},
		TypeTraveler: {
			"A free spacer drifting between ports.",
			"An explorer charting unknown routes.",
			"A courier carrying sealed packages.",
		},
		TypeScout: {
			"A reconnaissance vessel shares sensor data.",
			"A forward scout reports conditions.",
			"A pathfinder offers navigation intel.",
		},
		TypeGuard: {
			"A patrol ship demands identification.",
			"Corporate security scans your vessel.",
			"Military pickets challenge your approach.",
		},
	},
	engine.GenreHorror: {
		TypeTrader: {
			"A survivor willing to trade supplies.",
			"A scavenger with salvaged goods.",
			"Someone who found things worth trading.",
		},
		TypeRefugee: {
			"Terrified survivors seeking safety.",
			"A family running from something terrible.",
			"Desperate people with nowhere to go.",
		},
		TypeBandit: {
			"Dangerous people who've lost their humanity.",
			"Hungry raiders with dead eyes.",
			"Desperate killers who want what you have.",
		},
		TypeTraveler: {
			"A lone survivor, cautious but not hostile.",
			"Someone trying to find somewhere safe.",
			"A wanderer who's seen too much.",
		},
		TypeScout: {
			"Someone who knows the safe routes.",
			"A survivor who scouts ahead.",
			"A watcher who knows what's coming.",
		},
		TypeGuard: {
			"Armed survivors protecting their people.",
			"A checkpoint with suspicious guards.",
			"Defenders who shoot first and ask later.",
		},
	},
	engine.GenreCyberpunk: {
		TypeTrader: {
			"A street fixer with connections.",
			"A black market dealer with rare tech.",
			"A corporate vendor with licensed goods.",
		},
		TypeRefugee: {
			"Corporate refugees with burned identities.",
			"Zeroed citizens with nothing left.",
			"People fleeing gang territory.",
		},
		TypeBandit: {
			"Chrome-heavy gangers looking for prey.",
			"Boosted thugs blocking the way.",
			"Corporate muscle gone rogue.",
		},
		TypeTraveler: {
			"A data courier moving between sectors.",
			"A nomad passing through the zone.",
			"A wandering tech looking for work.",
		},
		TypeScout: {
			"A netrunner with local intel.",
			"A street kid who knows the territory.",
			"A drone operator with eyes everywhere.",
		},
		TypeGuard: {
			"Corporate security in full tactical gear.",
			"Private military contractors on patrol.",
			"Gang enforcers guarding their turf.",
		},
	},
	engine.GenrePostapoc: {
		TypeTrader: {
			"A water merchant with precious cargo.",
			"A wandering trader with salvage.",
			"A vault dweller with pre-war goods.",
		},
		TypeRefugee: {
			"Irradiated survivors seeking clean land.",
			"People fleeing raider territory.",
			"A family looking for somewhere safe.",
		},
		TypeBandit: {
			"Raiders in salvaged armor.",
			"War boys howling for blood.",
			"Desperate killers who've gone feral.",
		},
		TypeTraveler: {
			"A wasteland wanderer seeking something.",
			"A lone drifter walking the roads.",
			"Someone following old pre-war maps.",
		},
		TypeScout: {
			"A tracker who knows the safe paths.",
			"A scout from a nearby settlement.",
			"Someone who's mapped the radiation zones.",
		},
		TypeGuard: {
			"Settlement militia on watch.",
			"Gate guards with makeshift weapons.",
			"Perimeter patrol from a fortified camp.",
		},
	},
}

func (g *Generator) generateDialogue(npcType NPCType, alignment Alignment) []string {
	dialogues := dialogueByGenreAlignment[g.genre]
	if dialogues == nil {
		dialogues = dialogueByGenreAlignment[engine.GenreFantasy]
	}
	alignDialogues := dialogues[alignment]
	if len(alignDialogues) == 0 {
		return []string{"..."}
	}

	// Pick 2-3 dialogue options
	count := 2 + g.gen.Intn(2)
	if count > len(alignDialogues) {
		count = len(alignDialogues)
	}

	result := make([]string, 0, count)
	used := make(map[int]bool)

	for len(result) < count {
		idx := g.gen.Intn(len(alignDialogues))
		if !used[idx] {
			used[idx] = true
			result = append(result, alignDialogues[idx])
		}
	}

	return result
}

var dialogueByGenreAlignment = map[engine.GenreID]map[Alignment][]string{
	engine.GenreFantasy: {
		AlignmentHostile: {
			"Your coin or your life!",
			"No one passes without paying tribute.",
			"Take them!",
		},
		AlignmentSuspicious: {
			"State your business, stranger.",
			"Keep your hands where I can see them.",
			"What brings you to these parts?",
		},
		AlignmentNeutral: {
			"Greetings, traveler.",
			"The road is long, is it not?",
			"Safe travels to you.",
		},
		AlignmentFriendly: {
			"Well met, friend!",
			"It's good to see friendly faces.",
			"How can I help you?",
		},
		AlignmentAllied: {
			"We stand together in these dark times.",
			"Call on me if you need aid.",
			"Your cause is my cause.",
		},
	},
	engine.GenreScifi: {
		AlignmentHostile: {
			"Power down and prepare to be boarded.",
			"Your cargo is now ours.",
			"Surrender or be destroyed.",
		},
		AlignmentSuspicious: {
			"Transmit your credentials.",
			"State your origin and destination.",
			"You're in restricted space.",
		},
		AlignmentNeutral: {
			"Hailing frequencies open.",
			"Safe travels through the void.",
			"Clear skies, spacer.",
		},
		AlignmentFriendly: {
			"Good to see another friendly ship.",
			"Need any assistance out here?",
			"The black is lonely - good to have company.",
		},
		AlignmentAllied: {
			"We've got your back out here.",
			"Our frequencies are always open to you.",
			"Together we're stronger.",
		},
	},
	engine.GenreHorror: {
		AlignmentHostile: {
			"We need what you've got.",
			"Nothing personal. Survival.",
			"Don't make this harder than it needs to be.",
		},
		AlignmentSuspicious: {
			"Are you bitten? Show me your arms.",
			"How many are you? Where did you come from?",
			"Stay back. I mean it.",
		},
		AlignmentNeutral: {
			"Just passing through. Don't want trouble.",
			"Seen anything out there?",
			"Stay safe.",
		},
		AlignmentFriendly: {
			"Thank god, other survivors.",
			"We have to stick together.",
			"There's safety in numbers.",
		},
		AlignmentAllied: {
			"We're in this together.",
			"I've got your back.",
			"Till the end.",
		},
	},
	engine.GenreCyberpunk: {
		AlignmentHostile: {
			"Wrong sector, choom.",
			"Transfer your creds. Now.",
			"Flatline or cooperate.",
		},
		AlignmentSuspicious: {
			"Who sent you?",
			"You corpo or street?",
			"What's your angle?",
		},
		AlignmentNeutral: {
			"Just biz, nothing personal.",
			"Keep it professional.",
			"We done here?",
		},
		AlignmentFriendly: {
			"Need a fixer? I'm your choom.",
			"Got some good deals for you.",
			"Always good to meet new faces.",
		},
		AlignmentAllied: {
			"You're one of us now.",
			"We ride together.",
			"Your enemies are my enemies.",
		},
	},
	engine.GenrePostapoc: {
		AlignmentHostile: {
			"Everything you've got. Now.",
			"Shiny and chrome!",
			"You're dead, smoothskin.",
		},
		AlignmentSuspicious: {
			"You carrying the sickness?",
			"Where'd you get those supplies?",
			"What do you want?",
		},
		AlignmentNeutral: {
			"Just trying to survive.",
			"The wastes are hard on everyone.",
			"Good luck out there.",
		},
		AlignmentFriendly: {
			"You look like good people.",
			"We help each other out here.",
			"Come, rest by our fire.",
		},
		AlignmentAllied: {
			"You saved us. We won't forget.",
			"Our camp is your camp.",
			"Together we rebuild.",
		},
	},
}

func (g *Generator) generateTradeGoods() []TradeGood {
	goods := tradeGoodsByGenre[g.genre]
	if goods == nil {
		goods = tradeGoodsByGenre[engine.GenreFantasy]
	}

	// Generate 3-5 trade goods
	count := 3 + g.gen.Intn(3)
	if count > len(goods) {
		count = len(goods)
	}

	result := make([]TradeGood, 0, count)
	used := make(map[int]bool)

	for len(result) < count {
		idx := g.gen.Intn(len(goods))
		if !used[idx] {
			used[idx] = true
			good := goods[idx]
			// Vary quantity and price
			good.Quantity = 1 + g.gen.Intn(10)
			good.Price = good.Price * (0.8 + g.gen.Float64()*0.4)
			result = append(result, good)
		}
	}

	return result
}

var tradeGoodsByGenre = map[engine.GenreID][]TradeGood{
	engine.GenreFantasy: {
		{Name: "Bread", Price: 5},
		{Name: "Water Skin", Price: 8},
		{Name: "Healing Herbs", Price: 15},
		{Name: "Rope", Price: 10},
		{Name: "Torches", Price: 3},
		{Name: "Iron Rations", Price: 12},
		{Name: "Blanket", Price: 7},
		{Name: "Lantern Oil", Price: 6},
	},
	engine.GenreScifi: {
		{Name: "Ration Packs", Price: 10},
		{Name: "Water Purifier", Price: 25},
		{Name: "Med-Gel", Price: 20},
		{Name: "Fuel Cells", Price: 30},
		{Name: "O2 Tanks", Price: 15},
		{Name: "Repair Kit", Price: 35},
		{Name: "Stim Pack", Price: 40},
		{Name: "Navigation Data", Price: 50},
	},
	engine.GenreHorror: {
		{Name: "Canned Food", Price: 15},
		{Name: "Bottled Water", Price: 20},
		{Name: "First Aid Kit", Price: 30},
		{Name: "Gasoline", Price: 25},
		{Name: "Ammunition", Price: 35},
		{Name: "Batteries", Price: 10},
		{Name: "Antibiotics", Price: 50},
		{Name: "Rope", Price: 8},
	},
	engine.GenreCyberpunk: {
		{Name: "Synth-Food", Price: 10},
		{Name: "Purified Water", Price: 15},
		{Name: "Trauma Kit", Price: 40},
		{Name: "Battery Pack", Price: 20},
		{Name: "Stim", Price: 25},
		{Name: "ICE Breaker", Price: 75},
		{Name: "Cyberware Patch", Price: 60},
		{Name: "Cred Chip", Price: 100},
	},
	engine.GenrePostapoc: {
		{Name: "MRE", Price: 15},
		{Name: "Clean Water", Price: 25},
		{Name: "Rad-Away", Price: 40},
		{Name: "Fuel Can", Price: 30},
		{Name: "Ammo Box", Price: 35},
		{Name: "Med Kit", Price: 45},
		{Name: "Scrap Metal", Price: 10},
		{Name: "Pre-War Tech", Price: 80},
	},
}
