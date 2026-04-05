package journey

import (
	"fmt"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// Generator creates procedural multi-leg journeys
type Generator struct {
	gen    *seed.Generator
	genre  engine.GenreID
	legID  int
	stopID int
}

// NewGenerator creates a journey generator with the given seed and genre
func NewGenerator(masterSeed int64, genre engine.GenreID) *Generator {
	return &Generator{
		gen:    seed.NewGenerator(masterSeed, "journey"),
		genre:  genre,
		legID:  0,
		stopID: 0,
	}
}

// SetGenre updates the generator's active genre
func (g *Generator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// GenerateCampaign creates a complete multi-leg campaign
func (g *Generator) GenerateCampaign(legCount int) *Campaign {
	if legCount < 2 {
		legCount = 2
	}
	if legCount > 4 {
		legCount = 4
	}

	name := g.generateCampaignName()
	desc := g.generateCampaignDescription(legCount)

	campaign := NewCampaign(
		fmt.Sprintf("campaign_%d", g.gen.Master()),
		name,
		desc,
		g.genre,
	)

	// Generate all legs with escalating difficulty
	prevDestination := g.generateOriginName()
	for i := 0; i < legCount; i++ {
		difficulty := g.difficultyForLeg(i, legCount)
		distance := g.distanceForLeg(i, legCount)
		destination := g.generateDestinationName(i)

		leg := g.GenerateLeg(prevDestination, destination, distance, difficulty)
		campaign.AddLeg(leg)

		// Generate stopover after each leg except the last
		if i < legCount-1 {
			stopover := g.GenerateStopover(leg.ID, destination)
			campaign.AddStopover(stopover)
		}

		prevDestination = destination
	}

	return campaign
}

// GenerateLeg creates a single journey leg
func (g *Generator) GenerateLeg(origin, destination string, distance int, difficulty DifficultyLevel) *Leg {
	g.legID++

	name := g.generateLegName(origin, destination)
	leg := NewLeg(LegID(g.legID), name, origin, destination, distance, difficulty, g.genre)

	leg.Description = g.generateLegDescription(distance, difficulty)
	leg.TerrainType = g.generateTerrainType()
	leg.Hazards = g.generateHazards(difficulty)

	return leg
}

// GenerateStopover creates an intermediate hub city
func (g *Generator) GenerateStopover(afterLeg LegID, locationName string) *Stopover {
	g.stopID++

	name := g.generateStopoverName(locationName)
	desc := g.generateStopoverDescription()

	stopover := NewStopover(g.stopID, name, desc, afterLeg, g.genre)

	// Add services based on random selection
	stopover.Services = g.generateServices()
	stopover.Features = g.generateStopoverFeatures()
	stopover.Inhabitants = g.generateInhabitants()

	// Price modifiers
	stopover.BuyPriceModifier = 0.9 + float64(g.gen.Intn(30))/100.0
	stopover.SellPriceModifier = 0.8 + float64(g.gen.Intn(40))/100.0

	return stopover
}

// difficultyForLeg calculates escalating difficulty
func (g *Generator) difficultyForLeg(legIndex, totalLegs int) DifficultyLevel {
	// Base difficulty scales with leg index
	if legIndex == 0 {
		return DifficultyEasy
	}
	if legIndex == totalLegs-1 {
		return DifficultyHard
	}
	return DifficultyNormal
}

// distanceForLeg calculates escalating distance
func (g *Generator) distanceForLeg(legIndex, totalLegs int) int {
	baseDistance := 100 + legIndex*50
	variance := g.gen.Intn(30)
	return baseDistance + variance
}

func (g *Generator) generateCampaignName() string {
	campaignNames := map[engine.GenreID][]string{
		engine.GenreFantasy: {
			"The Great Pilgrimage", "Journey to the Far Kingdoms", "Quest for the Eternal Flame",
			"The Dragon's Path", "Odyssey of the Sacred Relic", "March to the World's Edge",
		},
		engine.GenreScifi: {
			"Deep Space Expedition", "The Colony Run", "Voyage to Distant Stars",
			"Operation Frontier", "The Long Haul", "Interstellar Transit",
		},
		engine.GenreHorror: {
			"Flight from Darkness", "The Long Night", "Escape from the Damned Lands",
			"Journey Through Shadow", "The Desperate Exodus", "Path of the Survivors",
		},
		engine.GenreCyberpunk: {
			"The Grid Runner", "Megacity Exodus", "Data Trail", "The Smuggler's Route",
			"Corporate Extraction", "Underground Railroad",
		},
		engine.GenrePostapoc: {
			"The Long Walk", "Caravan of Hope", "Journey to New Eden",
			"The Wasteland Crossing", "Migration of the Last", "Road to Sanctuary",
		},
	}

	names := campaignNames[g.genre]
	return seed.Choice(g.gen, names)
}

func (g *Generator) generateCampaignDescription(legCount int) string {
	descriptions := map[engine.GenreID][]string{
		engine.GenreFantasy: {
			fmt.Sprintf("A perilous journey across %d realms, through ancient forests and over treacherous mountains.", legCount),
			fmt.Sprintf("An epic %d-part odyssey through lands both wondrous and deadly.", legCount),
		},
		engine.GenreScifi: {
			fmt.Sprintf("A multi-system voyage spanning %d jump points through uncharted space.", legCount),
			fmt.Sprintf("A %d-sector journey across the frontier of known space.", legCount),
		},
		engine.GenreHorror: {
			fmt.Sprintf("A desperate %d-stage flight from the horrors that have consumed the world.", legCount),
			fmt.Sprintf("A harrowing journey through %d zones of nightmare and despair.", legCount),
		},
		engine.GenreCyberpunk: {
			fmt.Sprintf("A %d-district run through the neon-lit underbelly of the megacities.", legCount),
			fmt.Sprintf("A dangerous %d-part extraction through corporate territory.", legCount),
		},
		engine.GenrePostapoc: {
			fmt.Sprintf("A %d-region trek across the blasted wastelands in search of hope.", legCount),
			fmt.Sprintf("A grueling %d-stage migration through the ruins of civilization.", legCount),
		},
	}

	descs := descriptions[g.genre]
	return seed.Choice(g.gen, descs)
}

func (g *Generator) generateOriginName() string {
	origins := map[engine.GenreID][]string{
		engine.GenreFantasy: {
			"Ironhold Keep", "Silverbrook Village", "The Shattered Tower",
			"Moonhaven", "Thornwick", "The Fallen Temple",
		},
		engine.GenreScifi: {
			"Station Terminus", "New Terra Colony", "Orbital Platform Zeta",
			"Outpost Vanguard", "The Drifter's Haven", "Mining Station Kappa",
		},
		engine.GenreHorror: {
			"The Last Refuge", "Barricade Town", "Sanctuary Falls",
			"The Old Bunker", "Haven's End", "The Forsaken Shelter",
		},
		engine.GenreCyberpunk: {
			"The Undercity", "Sector 7 Slums", "Neon District",
			"The Sprawl", "Dead Zone Alpha", "Freeport Terminal",
		},
		engine.GenrePostapoc: {
			"Rust Town", "The Bunker", "Survivor's Camp",
			"The Old Mall", "Scrap City", "The Underground",
		},
	}

	names := origins[g.genre]
	return seed.Choice(g.gen, names)
}

func (g *Generator) generateDestinationName(legIndex int) string {
	destinations := map[engine.GenreID][]string{
		engine.GenreFantasy: {
			"The Crystal City", "Dragonspire", "The Eternal Gardens",
			"Stormhold", "The Sacred Grove", "Goldenreach",
			"The Frozen Citadel", "Sunhaven", "The Emerald Isles",
		},
		engine.GenreScifi: {
			"New Eden Colony", "The Nexus Station", "Paradise Prime",
			"The Core Worlds", "Sanctuary Station", "The Fleet",
			"Hyperion Hub", "The Rim", "Frontier Station",
		},
		engine.GenreHorror: {
			"The Safe Zone", "Sanctuary", "The Fortified City",
			"Hope's Landing", "The Walled Town", "The Lighthouse",
			"The Island", "The Stronghold", "Dawn's Refuge",
		},
		engine.GenreCyberpunk: {
			"The Upper City", "Corporate Tower", "The Free Zone",
			"The Grid", "Neo Tokyo", "Data Haven",
			"The Enclave", "Neutral Ground", "The Bright Side",
		},
		engine.GenrePostapoc: {
			"New Eden", "The Green Zone", "Sanctuary",
			"The Settlement", "Safe Harbor", "The Oasis",
			"The Valley", "Hope City", "The Promised Land",
		},
	}

	names := destinations[g.genre]
	idx := g.gen.Intn(len(names))
	return names[idx]
}

func (g *Generator) generateLegName(origin, destination string) string {
	formats := map[engine.GenreID][]string{
		engine.GenreFantasy: {
			"The Road from %s to %s", "%s to %s: A Perilous Path",
			"Journey from %s toward %s", "The %s-%s Trail",
		},
		engine.GenreScifi: {
			"%s to %s Route", "Sector Transit: %s-%s",
			"Jump Route %s to %s", "%s-%s Corridor",
		},
		engine.GenreHorror: {
			"Escape from %s to %s", "The Dark Road: %s to %s",
			"Flight from %s toward %s", "%s to %s: Through the Nightmare",
		},
		engine.GenreCyberpunk: {
			"The %s-%s Run", "Extraction Route: %s to %s",
			"%s to %s Underground", "Grid Path: %s-%s",
		},
		engine.GenrePostapoc: {
			"The Wasteland: %s to %s", "%s to %s Migration",
			"The Long Road from %s to %s", "%s-%s Caravan Route",
		},
	}

	fmts := formats[g.genre]
	format := seed.Choice(g.gen, fmts)
	return fmt.Sprintf(format, origin, destination)
}

func (g *Generator) generateLegDescription(distance int, difficulty DifficultyLevel) string {
	difficultyDescs := map[DifficultyLevel]string{
		DifficultyEasy:    "a relatively safe passage",
		DifficultyNormal:  "a challenging trek",
		DifficultyHard:    "a dangerous journey",
		DifficultyExtreme: "a near-suicidal expedition",
	}

	templates := map[engine.GenreID][]string{
		engine.GenreFantasy: {
			"A %d-league journey through %s. The path promises %s.",
			"This %d-league stretch offers %s. Travelers speak of %s.",
		},
		engine.GenreScifi: {
			"A %d-parsec transit representing %s. Sensors indicate %s.",
			"This %d-parsec route is known as %s. Navigation charts show %s.",
		},
		engine.GenreHorror: {
			"A %d-mile trek through the darkness - %s. Survivors whisper of %s.",
			"This %d-mile journey is %s. Few return from what lies ahead.",
		},
		engine.GenreCyberpunk: {
			"A %d-block run across %s. Intel suggests %s.",
			"This %d-block stretch is %s. Street runners know the dangers.",
		},
		engine.GenrePostapoc: {
			"A %d-mile trek across the wastes - %s. Scouts report %s.",
			"This %d-mile journey is %s. Radiation and raiders await.",
		},
	}

	terrain := g.generateTerrainType()
	templates_genre := templates[g.genre]
	template := seed.Choice(g.gen, templates_genre)

	return fmt.Sprintf(template, distance, terrain, difficultyDescs[difficulty])
}

func (g *Generator) generateTerrainType() string {
	terrains := map[engine.GenreID][]string{
		engine.GenreFantasy: {
			"dense forest", "mountain pass", "open plains", "marshland",
			"desert wastes", "frozen tundra", "coastal cliffs", "dark caverns",
		},
		engine.GenreScifi: {
			"asteroid field", "nebula", "radiation belt", "debris field",
			"void space", "ion storm zone", "magnetic anomaly", "dark matter zone",
		},
		engine.GenreHorror: {
			"abandoned city", "haunted woods", "fog-shrouded valley", "infested zone",
			"corpse fields", "the red zone", "silent suburbs", "the darkness",
		},
		engine.GenreCyberpunk: {
			"industrial sector", "combat zone", "corporate territory", "no-man's land",
			"toxic district", "surveillance zone", "gang territory", "dead network zone",
		},
		engine.GenrePostapoc: {
			"irradiated zone", "dust storms", "urban ruins", "scorched earth",
			"toxic swamp", "glass desert", "the dead zone", "mutant territory",
		},
	}

	terrainList := terrains[g.genre]
	return seed.Choice(g.gen, terrainList)
}

func (g *Generator) generateHazards(difficulty DifficultyLevel) []string {
	hazardCount := int(difficulty) + 1
	if hazardCount > 4 {
		hazardCount = 4
	}

	hazards := map[engine.GenreID][]string{
		engine.GenreFantasy: {
			"bandit ambushes", "wild beast attacks", "magical storms",
			"cursed ground", "goblin raids", "dragon sightings",
			"enchanted traps", "spirit hauntings",
		},
		engine.GenreScifi: {
			"pirate attacks", "solar flares", "system failures",
			"hostile aliens", "gravity wells", "navigation hazards",
			"radiation bursts", "micrometeorite swarms",
		},
		engine.GenreHorror: {
			"undead hordes", "madness zones", "cult ambushes",
			"creature nests", "corruption spreading", "psychic attacks",
			"nightmare visions", "possession events",
		},
		engine.GenreCyberpunk: {
			"gang warfare", "corporate patrols", "security drones",
			"net attacks", "toxic spills", "emp zones",
			"killbot swarms", "sniper nests",
		},
		engine.GenrePostapoc: {
			"raider attacks", "radiation pockets", "mutant swarms",
			"dust storms", "resource scarcity", "cannibal bands",
			"equipment failure", "disease outbreaks",
		},
	}

	hazardList := hazards[g.genre]
	selected := make([]string, 0, hazardCount)

	for i := 0; i < hazardCount && i < len(hazardList); i++ {
		idx := g.gen.Intn(len(hazardList))
		selected = append(selected, hazardList[idx])
	}

	return selected
}

func (g *Generator) generateStopoverName(baseName string) string {
	prefixes := map[engine.GenreID][]string{
		engine.GenreFantasy: {
			"Fort", "Camp", "Haven", "Rest", "Way", "Trade",
		},
		engine.GenreScifi: {
			"Station", "Hub", "Port", "Dock", "Depot", "Node",
		},
		engine.GenreHorror: {
			"Refuge", "Bunker", "Camp", "Safe", "Hold", "Shelter",
		},
		engine.GenreCyberpunk: {
			"Hub", "Node", "Port", "Zone", "Den", "Nest",
		},
		engine.GenrePostapoc: {
			"Camp", "Post", "Haven", "Hold", "Stop", "Base",
		},
	}

	prefix := seed.Choice(g.gen, prefixes[g.genre])
	return fmt.Sprintf("%s %s", baseName, prefix)
}

func (g *Generator) generateStopoverDescription() string {
	descriptions := map[engine.GenreID][]string{
		engine.GenreFantasy: {
			"A fortified waystation offering respite to weary travelers.",
			"A bustling trade post where merchants gather from distant lands.",
			"A sacred refuge maintained by the temple order.",
		},
		engine.GenreScifi: {
			"An orbital station providing essential services to passing ships.",
			"A frontier outpost serving the needs of deep space travelers.",
			"A corporate-run facility offering repairs and resupply.",
		},
		engine.GenreHorror: {
			"A heavily fortified refuge against the horrors outside.",
			"A survivor camp where desperate souls find temporary safety.",
			"An underground shelter hidden from the things that hunt.",
		},
		engine.GenreCyberpunk: {
			"A neutral zone where all factions observe an uneasy truce.",
			"A fixer's haven offering services to those who can pay.",
			"An underground market operating beyond corporate reach.",
		},
		engine.GenrePostapoc: {
			"A survivor settlement built from the ruins of the old world.",
			"A fortified camp where travelers can trade and rest.",
			"A community of scavengers offering fair deals to passers-by.",
		},
	}

	return seed.Choice(g.gen, descriptions[g.genre])
}

func (g *Generator) generateServices() []StopoverService {
	// Every stopover has trading
	services := []StopoverService{ServiceTrading}

	// 80% chance of repairs
	if g.gen.Intn(10) < 8 {
		services = append(services, ServiceRepairs)
	}

	// 60% chance of recruitment
	if g.gen.Intn(10) < 6 {
		services = append(services, ServiceRecruitment)
	}

	// 40% chance of upgrades
	if g.gen.Intn(10) < 4 {
		services = append(services, ServiceUpgrades)
	}

	// 70% chance of information
	if g.gen.Intn(10) < 7 {
		services = append(services, ServiceInformation)
	}

	// 50% chance of healing
	if g.gen.Intn(10) < 5 {
		services = append(services, ServiceHealing)
	}

	return services
}

func (g *Generator) generateStopoverFeatures() []string {
	features := map[engine.GenreID][]string{
		engine.GenreFantasy: {
			"ancient well with healing waters", "shrine to the road god",
			"enchanted forge", "wandering minstrel", "fortune teller's tent",
			"exotic animal market", "mysterious locked vault",
		},
		engine.GenreScifi: {
			"advanced medical bay", "ship upgrade terminal",
			"alien artifact display", "bounty board", "news feed terminal",
			"smuggler's locker", "encrypted data cache",
		},
		engine.GenreHorror: {
			"reinforced panic room", "zombie disposal pit",
			"armored watchtower", "hidden exit tunnel", "communal fire pit",
			"makeshift shrine", "survivor memorial",
		},
		engine.GenreCyberpunk: {
			"black market terminal", "illegal upgrade clinic",
			"encrypted data broker", "underground fight ring", "synth bar",
			"neural dive booth", "weapons cache",
		},
		engine.GenrePostapoc: {
			"water purification system", "scrap metal forge",
			"community garden", "ammunition press", "radiation shelter",
			"vehicle repair bay", "radio tower",
		},
	}

	featureList := features[g.genre]
	count := 1 + g.gen.Intn(3)
	selected := make([]string, 0, count)

	for i := 0; i < count; i++ {
		idx := g.gen.Intn(len(featureList))
		selected = append(selected, featureList[idx])
	}

	return selected
}

func (g *Generator) generateInhabitants() string {
	inhabitants := map[engine.GenreID][]string{
		engine.GenreFantasy: {
			"a mix of travelers, merchants, and local militia",
			"devout pilgrims and temple guardians",
			"retired adventurers and their families",
			"suspicious traders and wandering sellswords",
		},
		engine.GenreScifi: {
			"station crew and transient spacers",
			"corporate employees and independent contractors",
			"refugees from frontier conflicts",
			"a diverse mix of species and cultures",
		},
		engine.GenreHorror: {
			"traumatized survivors and hardened fighters",
			"a tight-knit community of those who escaped",
			"paranoid watchmen and desperate refugees",
			"a mix of the hopeful and the broken",
		},
		engine.GenreCyberpunk: {
			"fixers, runners, and those who employ them",
			"corporate exiles and street-level operators",
			"hackers, mercs, and information brokers",
			"the dispossessed seeking a new start",
		},
		engine.GenrePostapoc: {
			"survivors banded together for protection",
			"traders, scavengers, and wanderers",
			"a community built on mutual aid",
			"hardened veterans and hopeful newcomers",
		},
	}

	return seed.Choice(g.gen, inhabitants[g.genre])
}

// GenerateCampaignWithGenreShifts creates a campaign where each leg may have a different genre
func (g *Generator) GenerateCampaignWithGenreShifts(legCount int, genres []engine.GenreID) *Campaign {
	campaign := g.GenerateCampaign(legCount)
	campaign.EnableGenreShifts()

	// Apply genre shifts to each leg
	for i, leg := range campaign.Legs {
		if i < len(genres) {
			leg.SetGenre(genres[i])
			campaign.LegGenres[i] = genres[i]
		}
	}

	// Update stopovers to match their preceding leg
	for _, stopover := range campaign.Stopovers {
		legIdx := int(stopover.AfterLeg) - 1
		if legIdx >= 0 && legIdx < len(campaign.LegGenres) {
			stopover.SetGenre(campaign.LegGenres[legIdx])
		}
	}

	return campaign
}
