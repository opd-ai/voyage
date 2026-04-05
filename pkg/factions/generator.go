package factions

import (
	"fmt"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// Generator creates procedural factions for a run.
type Generator struct {
	gen    *seed.Generator
	genre  engine.GenreID
	nextID FactionID
}

// NewGenerator creates a new faction generator.
func NewGenerator(masterSeed int64, genre engine.GenreID) *Generator {
	return &Generator{
		gen:    seed.NewGenerator(masterSeed, "factions"),
		genre:  genre,
		nextID: 1,
	}
}

// SetGenre updates the generator's genre.
func (g *Generator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// GenerateFactions creates 4-6 factions for a run.
func (g *Generator) GenerateFactions(mapWidth, mapHeight int) *FactionManager {
	manager := NewFactionManager(g.genre)

	count := 4 + g.gen.Intn(3) // 4-6 factions
	ideologies := g.selectUniqueIdeologies(count)

	for i := 0; i < count; i++ {
		faction := g.generateFaction(ideologies[i])
		g.assignTerritory(faction, mapWidth, mapHeight, manager)
		manager.AddFaction(faction)
	}

	g.generateRelationships(manager)
	return manager
}

// selectUniqueIdeologies picks count unique ideologies.
func (g *Generator) selectUniqueIdeologies(count int) []Ideology {
	all := AllIdeologies()
	g.gen.Shuffle(len(all), func(i, j int) {
		all[i], all[j] = all[j], all[i]
	})
	if count > len(all) {
		count = len(all)
	}
	return all[:count]
}

// generateFaction creates a single faction.
func (g *Generator) generateFaction(ideology Ideology) *Faction {
	name := g.generateName(ideology)
	faction := NewFaction(g.nextID, name, ideology, g.genre)
	g.nextID++

	faction.Description = g.generateDescription(ideology)
	return faction
}

// generateName creates a faction name based on genre and ideology.
func (g *Generator) generateName(ideology Ideology) string {
	prefixes := factionPrefixes[g.genre]
	if prefixes == nil {
		prefixes = factionPrefixes[engine.GenreFantasy]
	}

	suffixes := factionSuffixes[g.genre][ideology]
	if suffixes == nil {
		suffixes = factionSuffixes[engine.GenreFantasy][ideology]
	}

	prefix := seed.Choice(g.gen, prefixes)
	suffix := seed.Choice(g.gen, suffixes)
	return fmt.Sprintf("%s %s", prefix, suffix)
}

var factionPrefixes = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"Golden", "Silver", "Iron", "Crimson", "Azure",
		"Northern", "Southern", "Eastern", "Western", "Ancient",
		"Royal", "Imperial", "Noble", "Sacred", "Eternal",
	},
	engine.GenreScifi: {
		"Stellar", "Void", "Quantum", "Neo", "Omega",
		"Alpha", "Delta", "Sigma", "Prime", "United",
		"Galactic", "Colonial", "Solar", "Lunar", "Orbital",
	},
	engine.GenreHorror: {
		"Last", "Dead", "Lost", "Fallen", "Broken",
		"Forsaken", "Desperate", "Remnant", "Forgotten", "Cursed",
		"Dark", "Shadow", "Blood", "Bone", "Ash",
	},
	engine.GenreCyberpunk: {
		"Neon", "Chrome", "Digital", "Cyber", "Ghost",
		"Black", "Red", "Zero", "Apex", "Omega",
		"Grid", "Net", "Wire", "Steel", "Synth",
	},
	engine.GenrePostapoc: {
		"Rust", "Dust", "Ash", "Iron", "Steel",
		"New", "Lost", "Last", "Free", "United",
		"Wasteland", "Desert", "Rad", "Scrap", "Road",
	},
}

var factionSuffixes = map[engine.GenreID]map[Ideology][]string{
	engine.GenreFantasy: {
		IdeologyMerchant:    {"Trading Company", "Merchants Guild", "Commerce League"},
		IdeologyMilitary:    {"Knights", "Legion", "Guard", "Wardens"},
		IdeologyReligious:   {"Order", "Brotherhood", "Covenant", "Temple"},
		IdeologyCriminal:    {"Syndicate", "Shadow Guild", "Rogues"},
		IdeologyScientific:  {"Circle", "Academy", "Conclave"},
		IdeologySurvivalist: {"Rangers", "Scouts", "Pathfinders"},
	},
	engine.GenreScifi: {
		IdeologyMerchant:    {"Corporation", "Industries", "Consortium"},
		IdeologyMilitary:    {"Fleet", "Navy", "Defense Force"},
		IdeologyReligious:   {"Collective", "Unity", "Church"},
		IdeologyCriminal:    {"Pirates", "Cartel", "Black Market"},
		IdeologyScientific:  {"Research Division", "Labs", "Institute"},
		IdeologySurvivalist: {"Colony", "Settlers", "Pioneers"},
	},
	engine.GenreHorror: {
		IdeologyMerchant:    {"Traders", "Merchants", "Suppliers"},
		IdeologyMilitary:    {"Militia", "Defenders", "Watch"},
		IdeologyReligious:   {"Cult", "Believers", "Chosen"},
		IdeologyCriminal:    {"Raiders", "Marauders", "Scavengers"},
		IdeologyScientific:  {"Scientists", "Researchers", "Lab"},
		IdeologySurvivalist: {"Survivors", "Haven", "Refuge"},
	},
	engine.GenreCyberpunk: {
		IdeologyMerchant:    {"Corp", "Industries", "Holdings"},
		IdeologyMilitary:    {"Security", "PMC", "Enforcers"},
		IdeologyReligious:   {"Sect", "Brotherhood", "Believers"},
		IdeologyCriminal:    {"Gang", "Syndicate", "Cartel"},
		IdeologyScientific:  {"Labs", "Tech", "Systems"},
		IdeologySurvivalist: {"Runners", "Nomads", "Street Crew"},
	},
	engine.GenrePostapoc: {
		IdeologyMerchant:    {"Caravan", "Traders", "Market"},
		IdeologyMilitary:    {"Warband", "Militia", "Army"},
		IdeologyReligious:   {"Cult", "Children", "Followers"},
		IdeologyCriminal:    {"Raiders", "Horde", "Gang"},
		IdeologyScientific:  {"Vault", "Bunker", "Enclave"},
		IdeologySurvivalist: {"Settlement", "Town", "Community"},
	},
}

// generateDescription creates a procedural faction description.
func (g *Generator) generateDescription(ideology Ideology) string {
	descriptions := factionDescriptions[g.genre]
	if descriptions == nil {
		descriptions = factionDescriptions[engine.GenreFantasy]
	}
	ideologyDescs := descriptions[ideology]
	if len(ideologyDescs) == 0 {
		return "A powerful faction in these lands."
	}
	return seed.Choice(g.gen, ideologyDescs)
}

var factionDescriptions = map[engine.GenreID]map[Ideology][]string{
	engine.GenreFantasy: {
		IdeologyMerchant: {
			"Controls the major trade routes through the region.",
			"Wealthy merchants who deal in rare goods and secrets.",
			"A network of traders spanning the known world.",
		},
		IdeologyMilitary: {
			"Sworn defenders of the realm with ancient traditions.",
			"Battle-hardened warriors who protect the innocent.",
			"A martial order bound by honor and duty.",
		},
		IdeologyReligious: {
			"Devoted followers of an ancient faith.",
			"Keepers of sacred mysteries and holy relics.",
			"Zealous believers who spread their word by any means.",
		},
		IdeologyCriminal: {
			"Operates in the shadows, controlling the underworld.",
			"Thieves and rogues who answer to their own code.",
			"A network of smugglers and assassins.",
		},
		IdeologyScientific: {
			"Seekers of arcane knowledge and magical power.",
			"Scholars who unlock the secrets of the universe.",
			"Practitioners of forbidden arts and ancient lore.",
		},
		IdeologySurvivalist: {
			"Wanderers who know every hidden path and danger.",
			"Hardy folk who live off the land and its bounty.",
			"Guides who navigate the wilderness with ease.",
		},
	},
	engine.GenreScifi: {
		IdeologyMerchant: {
			"Dominates interstellar trade across multiple systems.",
			"Controls the flow of rare resources and technology.",
			"A corporate entity with holdings throughout the sector.",
		},
		IdeologyMilitary: {
			"A formidable naval force protecting their space.",
			"Veterans of countless battles against hostile forces.",
			"Elite soldiers with advanced combat technology.",
		},
		IdeologyReligious: {
			"Worshippers of technology and the machine spirit.",
			"Believers in transcendence through unity.",
			"A faith born in the void between stars.",
		},
		IdeologyCriminal: {
			"Pirates who prey on shipping lanes.",
			"Smugglers dealing in contraband and stolen tech.",
			"A criminal network spanning the outer colonies.",
		},
		IdeologyScientific: {
			"Researchers pushing the boundaries of known science.",
			"Innovators developing cutting-edge technology.",
			"Scientists studying the mysteries of the cosmos.",
		},
		IdeologySurvivalist: {
			"Colonists carving out a living on hostile worlds.",
			"Settlers adapting to the challenges of space.",
			"Pioneers establishing new communities in the void.",
		},
	},
	engine.GenreHorror: {
		IdeologyMerchant: {
			"Trades in the few resources left in this dying world.",
			"Barters goods between scattered survivor groups.",
			"Controls what little supply chain remains.",
		},
		IdeologyMilitary: {
			"Armed survivors defending their territory.",
			"The last remnants of organized military forces.",
			"Fighters who've learned to survive against all odds.",
		},
		IdeologyReligious: {
			"Believers who see meaning in the apocalypse.",
			"A cult that has found purpose in the darkness.",
			"Fanatics who worship the horror that has come.",
		},
		IdeologyCriminal: {
			"Ruthless survivors who take what they want.",
			"Predators who prey on the desperate.",
			"Those who've abandoned all morality to survive.",
		},
		IdeologyScientific: {
			"Researchers seeking a cure or explanation.",
			"Scientists who study the infected.",
			"Those who believe knowledge can save them.",
		},
		IdeologySurvivalist: {
			"Ordinary people banded together for survival.",
			"A community that has weathered the worst.",
			"Survivors who still believe in humanity.",
		},
	},
	engine.GenreCyberpunk: {
		IdeologyMerchant: {
			"A megacorp with fingers in every market.",
			"Controls vital infrastructure and services.",
			"Wealthy enough to own governments.",
		},
		IdeologyMilitary: {
			"Private military contractors for hire.",
			"Corporate security forces with military hardware.",
			"Veterans selling their skills to the highest bidder.",
		},
		IdeologyReligious: {
			"A new faith spreading through the connected.",
			"Believers who see divinity in the digital.",
			"Zealots who reject the machine or embrace it.",
		},
		IdeologyCriminal: {
			"Organized crime operating in the shadows.",
			"A gang that controls the streets.",
			"Fixers and smugglers who move illegal goods.",
		},
		IdeologyScientific: {
			"Hackers and researchers pushing boundaries.",
			"Tech developers creating the next innovation.",
			"Scientists unafraid of ethical constraints.",
		},
		IdeologySurvivalist: {
			"Street people surviving day by day.",
			"Nomads who live outside the system.",
			"Those who refuse to be controlled.",
		},
	},
	engine.GenrePostapoc: {
		IdeologyMerchant: {
			"Traders who brave the wastes to move goods.",
			"Controls the flow of water, fuel, and food.",
			"A caravan network connecting settlements.",
		},
		IdeologyMilitary: {
			"Warlords who rule through force.",
			"Armed bands protecting their territory.",
			"The strongest survive and dominate.",
		},
		IdeologyReligious: {
			"Worshippers of the atom and its glow.",
			"Believers who see the end as a new beginning.",
			"A cult that has found meaning in the ruins.",
		},
		IdeologyCriminal: {
			"Raiders who take what they want.",
			"Bandits preying on the weak.",
			"Those who've embraced the lawless new world.",
		},
		IdeologyScientific: {
			"Preservers of old world knowledge.",
			"Researchers trying to rebuild civilization.",
			"Vault dwellers with pre-war technology.",
		},
		IdeologySurvivalist: {
			"Settlers building a new community.",
			"People working together to survive.",
			"Those who believe in a better tomorrow.",
		},
	},
}

// assignTerritory assigns territory blocks to a faction.
func (g *Generator) assignTerritory(f *Faction, mapWidth, mapHeight int, manager *FactionManager) {
	numBlocks := 1 + g.gen.Intn(3) // 1-3 territory blocks

	for i := 0; i < numBlocks; i++ {
		// Try to find a non-overlapping position
		for attempts := 0; attempts < 10; attempts++ {
			x := g.gen.Range(mapWidth/6, mapWidth*5/6)
			y := g.gen.Range(mapHeight/6, mapHeight*5/6)
			radius := g.gen.Range(2, 5)

			// Check for overlap with existing territories
			overlaps := false
			for _, other := range manager.AllFactions() {
				if other.ControlsPosition(x, y) {
					overlaps = true
					break
				}
			}

			if !overlaps {
				f.AddTerritory(x, y, radius)
				break
			}
		}
	}
}

// generateRelationships creates the faction relationship matrix.
func (g *Generator) generateRelationships(manager *FactionManager) {
	factions := manager.AllFactions()

	for i := 0; i < len(factions); i++ {
		for j := i + 1; j < len(factions); j++ {
			rel := g.determineRelationship(factions[i], factions[j])
			factions[i].SetRelation(factions[j].ID, rel)
			factions[j].SetRelation(factions[i].ID, rel)
		}
	}
}

// determineRelationship calculates initial relationship between two factions.
func (g *Generator) determineRelationship(a, b *Faction) Relationship {
	// Similar ideologies tend to be friendlier
	if a.Ideology == b.Ideology {
		return g.weightedRelationship([]float64{0.3, 0.4, 0.2, 0.1, 0.0})
	}

	// Opposing ideologies (criminal vs military, etc.)
	if areOpposingIdeologies(a.Ideology, b.Ideology) {
		return g.weightedRelationship([]float64{0.0, 0.05, 0.15, 0.3, 0.5})
	}

	// Default random relationship
	return g.weightedRelationship([]float64{0.1, 0.2, 0.4, 0.2, 0.1})
}

// areOpposingIdeologies checks if two ideologies naturally conflict.
func areOpposingIdeologies(a, b Ideology) bool {
	oppositions := map[Ideology][]Ideology{
		IdeologyMerchant:    {IdeologyCriminal},
		IdeologyMilitary:    {IdeologyCriminal},
		IdeologyReligious:   {IdeologyScientific},
		IdeologyCriminal:    {IdeologyMerchant, IdeologyMilitary},
		IdeologyScientific:  {IdeologyReligious},
		IdeologySurvivalist: {},
	}

	for _, opp := range oppositions[a] {
		if opp == b {
			return true
		}
	}
	return false
}

// weightedRelationship selects a relationship based on weights.
func (g *Generator) weightedRelationship(weights []float64) Relationship {
	return seed.WeightedChoice(g.gen, AllRelationships(), weights)
}
