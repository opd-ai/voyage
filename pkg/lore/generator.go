package lore

import (
	"fmt"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// Generator creates procedural lore content.
type Generator struct {
	gen        *seed.Generator
	genre      engine.GenreID
	nextInscID int
	nextDiscID int
}

// NewGenerator creates a new lore generator.
func NewGenerator(masterSeed int64, genre engine.GenreID) *Generator {
	return &Generator{
		gen:        seed.NewGenerator(masterSeed, "lore"),
		genre:      genre,
		nextInscID: 1,
		nextDiscID: 1,
	}
}

// SetGenre updates the generator's genre.
func (g *Generator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// GenerateInscription creates a procedural inscription.
func (g *Generator) GenerateInscription(x, y int) *Inscription {
	iType := seed.Choice(g.gen, AllInscriptionTypes())
	title := g.generateInscriptionTitle(iType)
	text := g.generateInscriptionText(iType)

	insc := NewInscription(g.nextInscID, iType, x, y, title, text, g.genre)
	g.nextInscID++
	return insc
}

// generateInscriptionTitle creates a title for an inscription.
func (g *Generator) generateInscriptionTitle(iType InscriptionType) string {
	titles := inscriptionTitles[g.genre]
	if titles == nil {
		titles = inscriptionTitles[engine.GenreFantasy]
	}
	typeTitles := titles[iType]
	if len(typeTitles) == 0 {
		return "Unknown Inscription"
	}
	return seed.Choice(g.gen, typeTitles)
}

var inscriptionTitles = map[engine.GenreID]map[InscriptionType][]string{
	engine.GenreFantasy: {
		TypeRuin:     {"Ancient Temple", "Forgotten Shrine", "Crumbling Tower", "Lost Archive"},
		TypeGrave:    {"Unmarked Grave", "Hero's Rest", "Mass Burial", "Noble's Tomb"},
		TypeSign:     {"Crossroads Marker", "Warning Post", "Mile Stone", "Direction Stone"},
		TypeMonument: {"Victory Column", "Memorial Obelisk", "Founder's Statue", "Sacred Stone"},
		TypeGraffiti: {"Scratched Warning", "Carved Initials", "Desperate Message", "Pilgrim's Mark"},
	},
	engine.GenreScifi: {
		TypeRuin:     {"Research Station", "Colony Remnant", "Derelict Outpost", "Crashed Ship"},
		TypeGrave:    {"Memorial Site", "Crash Memorial", "Pioneer Graves", "Colonist Markers"},
		TypeSign:     {"Navigation Beacon", "Warning Marker", "Landing Zone", "Hazard Signal"},
		TypeMonument: {"First Landing", "Founder Memorial", "Unity Pillar", "Achievement Monument"},
		TypeGraffiti: {"Distress Marker", "Survival Note", "Crew Tags", "Final Message"},
	},
	engine.GenreHorror: {
		TypeRuin:     {"Abandoned Hospital", "Burnt Church", "Collapsed Shelter", "Overrun Camp"},
		TypeGrave:    {"Fresh Graves", "Body Pit", "Memorial Wall", "Lost Names"},
		TypeSign:     {"Danger Warning", "Safe Zone Arrow", "Contamination Notice", "Survivor Message"},
		TypeMonument: {"Before the Fall", "Memorial Garden", "Lost District", "Remembrance"},
		TypeGraffiti: {"Spray-Painted Warning", "Scratched Note", "Blood Message", "Last Words"},
	},
	engine.GenreCyberpunk: {
		TypeRuin:     {"Corp Tower Ruins", "Collapsed Arcology", "Burned District", "Abandoned Lab"},
		TypeGrave:    {"Memorial Wall", "Gang Tags", "Lost Citizens", "Riot Memorial"},
		TypeSign:     {"Hazard Warning", "Zone Border", "Security Notice", "Corp Territory"},
		TypeMonument: {"Founder's Plaza", "Corporate Pride", "Unity Monument", "Progress Tower"},
		TypeGraffiti: {"Gang Tags", "Resistance Message", "Hacker Mark", "Street Art"},
	},
	engine.GenrePostapoc: {
		TypeRuin:     {"Bombed Building", "Vault Entrance", "Factory Ruins", "Highway Wreck"},
		TypeGrave:    {"Mass Grave", "Family Plot", "Soldier's Rest", "Unknown Dead"},
		TypeSign:     {"Radiation Warning", "Safe Water", "Raider Territory", "Route Marker"},
		TypeMonument: {"Pre-War Memorial", "Founder Statue", "Old World Sign", "Town Marker"},
		TypeGraffiti: {"Scratched Warning", "Survivor Count", "Supply Note", "Direction Arrow"},
	},
}

// generateInscriptionText creates text for an inscription.
func (g *Generator) generateInscriptionText(iType InscriptionType) string {
	texts := inscriptionTexts[g.genre]
	if texts == nil {
		texts = inscriptionTexts[engine.GenreFantasy]
	}
	typeTexts := texts[iType]
	if len(typeTexts) == 0 {
		return "The inscription is too weathered to read."
	}
	return seed.Choice(g.gen, typeTexts)
}

var inscriptionTexts = map[engine.GenreID]map[InscriptionType][]string{
	engine.GenreFantasy: {
		TypeRuin: {
			"Ancient symbols speak of a kingdom that fell to pride.",
			"These stones once housed a library of forbidden knowledge.",
			"The tower's collapse came when the mages reached too far.",
			"Here stood the sanctuary, before the darkness came.",
		},
		TypeGrave: {
			"Here lies one who traveled far but found only this.",
			"In memory of those who gave everything for the realm.",
			"Names worn away by time, but not forgotten.",
			"The fallen of the last battle rest here.",
		},
		TypeSign: {
			"Turn back. The path ahead holds only death.",
			"Three days to the next water. Travel wisely.",
			"The mountain pass is treacherous in winter.",
			"Safe haven lies to the east, if you can reach it.",
		},
		TypeMonument: {
			"Erected in honor of those who founded this land.",
			"Let this stone remind all of the price of freedom.",
			"Victory came at great cost. Remember the fallen.",
			"The ancients blessed this place. Treat it with respect.",
		},
		TypeGraffiti: {
			"Don't trust the ones who smile. -A Friend",
			"We passed this way. May others find hope. -The Lost",
			"Water is poisoned. Find the stream instead.",
			"Three of us left. Only I remain. Keep moving.",
		},
	},
	engine.GenreScifi: {
		TypeRuin: {
			"This station's logs speak of an experiment gone wrong.",
			"Emergency evacuation initiated. No further records.",
			"Research team vanished on day 47. Investigation inconclusive.",
			"Hull breach in sector 7. All hands lost.",
		},
		TypeGrave: {
			"First colonists of Kepler. They dreamed of a new home.",
			"Ship's crew. They gave their lives so others might live.",
			"Unknown casualties. May the void carry them home.",
			"Pioneer team Alpha. They found what we were looking for.",
		},
		TypeSign: {
			"HAZARD: Radiation levels exceed safe limits.",
			"Navigation hazard. Proceed with caution.",
			"Emergency beacon. Distress frequency 121.5.",
			"Quarantine zone. No unauthorized entry.",
		},
		TypeMonument: {
			"On this site, humanity first touched alien soil.",
			"In memory of the Unity mission. May their sacrifice inspire.",
			"Progress demands sacrifice. Remember those who gave all.",
			"The founders who made this possible are honored here.",
		},
		TypeGraffiti: {
			"They lied to us about what's out here. -Crew of Icarus",
			"If you find this, the corporation cannot be trusted.",
			"Supplies hidden at coordinates: CORRUPTED DATA",
			"We're not alone. Trust no one who says otherwise.",
		},
	},
	engine.GenreHorror: {
		TypeRuin: {
			"This was where they first appeared. No one believed us.",
			"Patient zero was found here. By then it was too late.",
			"The barricades held for three days. Then they adapted.",
			"We thought this place was safe. We were wrong.",
		},
		TypeGrave: {
			"Names of those who didn't make it. Remember them.",
			"We buried them quickly. Too quickly to mourn properly.",
			"The lost. The turned. The forgotten.",
			"My family. I'm sorry I couldn't save you.",
		},
		TypeSign: {
			"WARNING: Infected zone. Do not enter.",
			"Safe zone overrun. Find alternate route.",
			"Supplies depleted. Move on.",
			"They come at night. Don't stop moving.",
		},
		TypeMonument: {
			"Before the outbreak. When we still had hope.",
			"The old world is gone. This is what remains.",
			"In memory of what we lost.",
			"Never forget. Never let it happen again.",
		},
		TypeGraffiti: {
			"DON'T SLEEP. THEY KNOW WHEN YOU SLEEP.",
			"Bit on day 3. Still human on day 5. Why?",
			"The cure is real. Find Dr. Moran.",
			"If you read this, there's still hope. Keep going.",
		},
	},
	engine.GenreCyberpunk: {
		TypeRuin: {
			"Zaibatsu Corp headquarters. Fell when the riots came.",
			"Data center compromised. All records purged.",
			"Security breach in the old network. Stay offline.",
			"The AI that ran this place is still active. Be careful.",
		},
		TypeGrave: {
			"Runners who didn't make it. The street remembers.",
			"Zeroed by the corps. Their names live on.",
			"Victims of the corporate wars. Collateral damage.",
			"The forgotten. The discarded. The erased.",
		},
		TypeSign: {
			"TOXIC ZONE. Filtration required.",
			"Gang territory. Proceed at own risk.",
			"Corporate surveillance active. Mind your data.",
			"Safe house beyond. Knock three times.",
		},
		TypeMonument: {
			"Progress through technology. Corp propaganda from before.",
			"The founders. Before they became what they are now.",
			"Unity through submission. The old lie.",
			"Built on the backs of those they discarded.",
		},
		TypeGraffiti: {
			"The corps own everything. But they don't own us.",
			"Data cache nearby. Use code: NIGHTFALL",
			"AI is watching. Act natural.",
			"Free zone ahead. If you can make it.",
		},
	},
	engine.GenrePostapoc: {
		TypeRuin: {
			"This was a city once. Now it's just bones.",
			"The vault beneath held survivors. Empty now.",
			"Factory still runs. Don't know why. Don't care to find out.",
			"Hospital. Picked clean years ago. Not worth the risk.",
		},
		TypeGrave: {
			"We buried them as best we could. The land will remember.",
			"The first generation. They survived the bombs but not the after.",
			"Here lie those who tried to rebuild. They deserve better than this.",
			"Unknown remains. Too many to count. May they find peace.",
		},
		TypeSign: {
			"RADIATION HOT ZONE. Detour recommended.",
			"Clean water 10 miles east. Worth the trip.",
			"Raider territory. Turn back or prepare to fight.",
			"Settlement ahead. Peaceful. Trade welcome.",
		},
		TypeMonument: {
			"The old world fell here. A new one rose from the ashes.",
			"In memory of what was. A warning of what could be again.",
			"The founders of this settlement. May their legacy endure.",
			"Before the bombs. When the world still made sense.",
		},
		TypeGraffiti: {
			"Water here is good. Tested it myself. -Wanderer",
			"Raiders hit this area every full moon. Plan accordingly.",
			"Buried supplies 20 paces north. Share with others.",
			"We made it. You can too. Don't give up.",
		},
	},
}

// GenerateDiscovery creates a procedural discovery.
func (g *Generator) GenerateDiscovery(x, y int) *Discovery {
	dType := seed.Choice(g.gen, AllDiscoveryTypes())
	title := g.generateDiscoveryTitle(dType)
	vignette := g.generateDiscoveryVignette(dType)

	disc := NewDiscovery(g.nextDiscID, dType, x, y, title, vignette, g.genre)
	g.nextDiscID++

	// Add items
	items := g.generateDiscoveryItems(dType)
	for _, item := range items {
		disc.AddItem(item.Name, item.Quantity)
	}

	return disc
}

// generateDiscoveryTitle creates a title for a discovery.
func (g *Generator) generateDiscoveryTitle(dType DiscoveryType) string {
	titles := discoveryTitles[g.genre]
	if titles == nil {
		titles = discoveryTitles[engine.GenreFantasy]
	}
	typeTitles := titles[dType]
	if len(typeTitles) == 0 {
		return "Unknown Discovery"
	}
	return seed.Choice(g.gen, typeTitles)
}

var discoveryTitles = map[engine.GenreID]map[DiscoveryType][]string{
	engine.GenreFantasy: {
		DiscoveryVessel: {"Overturned Cart", "Broken Wagon", "Merchant's Carriage", "Noble's Coach"},
		DiscoveryCamp:   {"Cold Campfire", "Abandoned Tent", "Ranger's Hideout", "Pilgrim's Rest"},
		DiscoveryCache:  {"Buried Chest", "Hidden Stash", "Smuggler's Drop", "Emergency Supplies"},
		DiscoveryBody:   {"Fallen Knight", "Lost Merchant", "Unknown Traveler", "Wounded Scout"},
	},
	engine.GenreScifi: {
		DiscoveryVessel: {"Crashed Shuttle", "Derelict Pod", "Mining Rig", "Scout Ship"},
		DiscoveryCamp:   {"Research Camp", "Survival Shelter", "Emergency Habitat", "Field Station"},
		DiscoveryCache:  {"Supply Drop", "Equipment Cache", "Emergency Kit", "Hidden Stash"},
		DiscoveryBody:   {"Suited Figure", "Crew Remains", "Unknown Colonist", "Survey Team"},
	},
	engine.GenreHorror: {
		DiscoveryVessel: {"Crashed Car", "Abandoned Van", "Wrecked Truck", "Police Cruiser"},
		DiscoveryCamp:   {"Makeshift Shelter", "Barricaded Room", "Survivor Camp", "Safe House"},
		DiscoveryCache:  {"Hidden Supplies", "Emergency Stash", "Weapon Cache", "Medical Kit"},
		DiscoveryBody:   {"Fresh Corpse", "Suicide Victim", "Turned Victim", "Soldier's Body"},
	},
	engine.GenreCyberpunk: {
		DiscoveryVessel: {"Crashed Aerodyne", "Abandoned Runner", "Corporate Car", "Delivery Drone"},
		DiscoveryCamp:   {"Squatter's Hideout", "Runner's Den", "Tech Lab", "Safe House"},
		DiscoveryCache:  {"Data Stash", "Weapon Cache", "Medical Drop", "Credit Chip"},
		DiscoveryBody:   {"Flatlined Runner", "Corp Victim", "Street Casualty", "Gang Member"},
	},
	engine.GenrePostapoc: {
		DiscoveryVessel: {"Rusted Truck", "Military Vehicle", "Trader's Cart", "Scavenger's Ride"},
		DiscoveryCamp:   {"Survivor Camp", "Trading Post", "Raider Hideout", "Shelter"},
		DiscoveryCache:  {"Buried Supplies", "Pre-War Stash", "Weapon Cache", "Water Store"},
		DiscoveryBody:   {"Skeleton", "Fresh Victim", "Raider Corpse", "Wanderer's Remains"},
	},
}

// generateDiscoveryVignette creates vignette text for a discovery.
func (g *Generator) generateDiscoveryVignette(dType DiscoveryType) string {
	vignettes := discoveryVignettes[g.genre]
	if vignettes == nil {
		vignettes = discoveryVignettes[engine.GenreFantasy]
	}
	typeVignettes := vignettes[dType]
	if len(typeVignettes) == 0 {
		return "You find evidence of those who came before."
	}
	return seed.Choice(g.gen, typeVignettes)
}

var discoveryVignettes = map[engine.GenreID]map[DiscoveryType][]string{
	engine.GenreFantasy: {
		DiscoveryVessel: {
			"The wagon lies on its side, one wheel still slowly turning. Whatever attacked its owners did so quickly.",
			"A merchant's cart, its contents scattered. A journal lies open, the last entry hastily scrawled.",
			"The carriage bears a noble crest, but its occupants are long gone. Something valuable might remain.",
		},
		DiscoveryCamp: {
			"The campfire is cold, but the ashes are recent. Whoever was here left in a hurry.",
			"A tent flaps in the wind, its occupant absent. Personal effects suggest they meant to return.",
			"Signs of a struggle. Blood on the ground. But no bodies. What happened here?",
		},
		DiscoveryCache: {
			"Beneath a marked stone, you find a hidden cache. Someone prepared for emergencies.",
			"A hollowed tree reveals supplies wrapped in oilcloth. A smuggler's drop, perhaps.",
			"The chest is locked, but the lock is rusted. Inside, supplies that could save lives.",
		},
		DiscoveryBody: {
			"A fallen traveler lies here, wounds telling a story of violence. Their pack remains.",
			"The knight's armor is dented but salvageable. They died fighting something terrible.",
			"A peaceful expression. No wounds. Disease, perhaps, or exhaustion. Their belongings remain.",
		},
	},
	engine.GenreScifi: {
		DiscoveryVessel: {
			"The shuttle's hull is breached, but the cargo hold appears intact. Emergency lights still flicker.",
			"A mining rig, abandoned mid-operation. The crew's personal effects remain at their stations.",
			"The escape pod made a hard landing. Its occupant didn't survive, but their supplies did.",
		},
		DiscoveryCamp: {
			"Research equipment lies scattered. Data pads contain partial logs. Something interrupted their work.",
			"The habitat's life support still runs on battery backup. Empty, but recently occupied.",
			"An emergency shelter, properly deployed. Inside, signs of extended habitation.",
		},
		DiscoveryCache: {
			"A supply drop, marked with corporate insignia. Standard emergency kit.",
			"Someone hid equipment here deliberately. The caching was methodical.",
			"Emergency supplies in a sealed container. Standard colony survival package.",
		},
		DiscoveryBody: {
			"The suit's occupant died from exposure. Their beacon was never activated.",
			"A survey team member, based on insignia. Their data recorder might still work.",
			"Someone in a corporate uniform. Whatever they were running from caught up.",
		},
	},
	engine.GenreHorror: {
		DiscoveryVessel: {
			"The car's doors are open. Blood smears the seats. The keys are still in the ignition.",
			"A crashed ambulance, supplies scattered. Someone tried to flee with medical equipment.",
			"The van's back doors are torn open from inside. Whatever was locked in got out.",
		},
		DiscoveryCamp: {
			"The barricade failed. Boards torn away. Inside, signs of last stand.",
			"A survivor camp, recently abandoned. Food still warm. They heard something coming.",
			"The shelter is intact, but its occupant chose to leave. A note explains why.",
		},
		DiscoveryCache: {
			"Hidden under floorboards, a cache of supplies. Someone planned for the worst.",
			"A weapon stash, carefully concealed. Not used. The owner didn't make it back.",
			"Medical supplies in a hidden compartment. More valuable than gold now.",
		},
		DiscoveryBody: {
			"They were bitten but didn't turn. A bullet in their temple. Their choice.",
			"A soldier, still in uniform. They went down fighting. Their ammunition is spent.",
			"No marks. No wounds. Sometimes people just... give up. Their supplies remain.",
		},
	},
	engine.GenreCyberpunk: {
		DiscoveryVessel: {
			"The aerodyne's crash was recent. Corporate markings are still visible. So is the damage.",
			"A runner's vehicle, boosted but abandoned. The previous owner left in a hurry.",
			"Delivery drone, grounded by EMP. Its cargo is still secured.",
		},
		DiscoveryCamp: {
			"A squatter's den, recently vacated. Someone was watching the building across the street.",
			"Runner safe house. The tech here is outdated but functional. No one's coming back for it.",
			"A makeshift tech lab. Someone was working on something big. The project abandoned.",
		},
		DiscoveryCache: {
			"Data stash hidden behind a false wall. Corporate secrets, maybe. Worth checking.",
			"Weapon cache under the floorboards. Street level gear, but functional.",
			"Credit chips in a hidden safe. The previous owner isn't coming back for them.",
		},
		DiscoveryBody: {
			"A runner, flatlined. Neural burnout from the look of it. Their deck might be salvageable.",
			"Corporate uniform, bullet holes in the back. The street claimed another one.",
			"No ID chip. No records. A ghost in life and death. Their gear remains.",
		},
	},
	engine.GenrePostapoc: {
		DiscoveryVessel: {
			"The truck's engine is salvageable. Its previous owners weren't so lucky.",
			"A military vehicle, pre-war markings. The cargo hold has been picked clean, but...",
			"Trader's cart, overturned by raiders. They took the goods but missed the hidden compartment.",
		},
		DiscoveryCamp: {
			"A survivor camp, abandoned when the water ran out. Some supplies remain.",
			"Raider hideout, recently cleared. Someone else got here first. Maybe they missed something.",
			"A trading post that didn't last. Disease, judging by the makeshift quarantine.",
		},
		DiscoveryCache: {
			"Pre-war supplies buried under marked stones. Someone's emergency plan.",
			"A water cache, still sealed. More valuable than anything else out here.",
			"Weapon stash hidden in a collapsed building. The original owner didn't make it back.",
		},
		DiscoveryBody: {
			"A wanderer who didn't make it. Dehydration. Their supplies ran out miles ago.",
			"Raider corpse, killed by their own kind. Fighting over scraps.",
			"A settler, died protecting their cache. The cache is still intact.",
		},
	},
}

// generateDiscoveryItems creates items for a discovery.
func (g *Generator) generateDiscoveryItems(dType DiscoveryType) []DiscoveryItem {
	itemLists := discoveryItems[g.genre]
	if itemLists == nil {
		itemLists = discoveryItems[engine.GenreFantasy]
	}
	typeItems := itemLists[dType]
	if len(typeItems) == 0 {
		return nil
	}

	// Pick 1-3 items
	count := 1 + g.gen.Intn(3)
	if count > len(typeItems) {
		count = len(typeItems)
	}

	g.gen.Shuffle(len(typeItems), func(i, j int) {
		typeItems[i], typeItems[j] = typeItems[j], typeItems[i]
	})

	result := make([]DiscoveryItem, 0, count)
	for i := 0; i < count; i++ {
		item := typeItems[i]
		item.Quantity = 1 + g.gen.Intn(item.Quantity)
		result = append(result, item)
	}
	return result
}

var discoveryItems = map[engine.GenreID]map[DiscoveryType][]DiscoveryItem{
	engine.GenreFantasy: {
		DiscoveryVessel: {{Name: "Gold Coins", Quantity: 20}, {Name: "Rations", Quantity: 5}, {Name: "Cloth", Quantity: 3}},
		DiscoveryCamp:   {{Name: "Bedroll", Quantity: 1}, {Name: "Tinderbox", Quantity: 1}, {Name: "Dried Meat", Quantity: 3}},
		DiscoveryCache:  {{Name: "Gold Coins", Quantity: 50}, {Name: "Healing Herbs", Quantity: 5}, {Name: "Rope", Quantity: 2}},
		DiscoveryBody:   {{Name: "Personal Effects", Quantity: 1}, {Name: "Coin Purse", Quantity: 10}, {Name: "Dagger", Quantity: 1}},
	},
	engine.GenreScifi: {
		DiscoveryVessel: {{Name: "Fuel Cells", Quantity: 5}, {Name: "Ration Packs", Quantity: 10}, {Name: "Repair Kit", Quantity: 1}},
		DiscoveryCamp:   {{Name: "Data Pad", Quantity: 1}, {Name: "Med-Gel", Quantity: 3}, {Name: "O2 Tank", Quantity: 2}},
		DiscoveryCache:  {{Name: "Credits", Quantity: 100}, {Name: "Stim Pack", Quantity: 3}, {Name: "Tool Kit", Quantity: 1}},
		DiscoveryBody:   {{Name: "ID Tags", Quantity: 1}, {Name: "Personal Device", Quantity: 1}, {Name: "Credits", Quantity: 50}},
	},
	engine.GenreHorror: {
		DiscoveryVessel: {{Name: "Gasoline", Quantity: 5}, {Name: "Canned Food", Quantity: 5}, {Name: "Batteries", Quantity: 3}},
		DiscoveryCamp:   {{Name: "First Aid Kit", Quantity: 1}, {Name: "Bottled Water", Quantity: 3}, {Name: "Ammunition", Quantity: 10}},
		DiscoveryCache:  {{Name: "Antibiotics", Quantity: 2}, {Name: "Ammunition", Quantity: 20}, {Name: "Canned Food", Quantity: 10}},
		DiscoveryBody:   {{Name: "Personal Effects", Quantity: 1}, {Name: "Supplies", Quantity: 3}, {Name: "Weapon", Quantity: 1}},
	},
	engine.GenreCyberpunk: {
		DiscoveryVessel: {{Name: "Battery Pack", Quantity: 3}, {Name: "Synth-Food", Quantity: 5}, {Name: "Cred Chip", Quantity: 50}},
		DiscoveryCamp:   {{Name: "Data Shard", Quantity: 1}, {Name: "Stim", Quantity: 2}, {Name: "Tech Components", Quantity: 3}},
		DiscoveryCache:  {{Name: "Cred Chips", Quantity: 200}, {Name: "Cyberware Patch", Quantity: 1}, {Name: "ICE Breaker", Quantity: 1}},
		DiscoveryBody:   {{Name: "ID Chip", Quantity: 1}, {Name: "Cred Chip", Quantity: 75}, {Name: "Data Shard", Quantity: 1}},
	},
	engine.GenrePostapoc: {
		DiscoveryVessel: {{Name: "Fuel Can", Quantity: 3}, {Name: "MRE", Quantity: 5}, {Name: "Scrap Metal", Quantity: 5}},
		DiscoveryCamp:   {{Name: "Clean Water", Quantity: 5}, {Name: "Rad-Away", Quantity: 2}, {Name: "Ammo Box", Quantity: 1}},
		DiscoveryCache:  {{Name: "Caps", Quantity: 100}, {Name: "Med Kit", Quantity: 2}, {Name: "Pre-War Tech", Quantity: 1}},
		DiscoveryBody:   {{Name: "Personal Effects", Quantity: 1}, {Name: "Caps", Quantity: 30}, {Name: "Supplies", Quantity: 3}},
	},
}

// GenerateCodexEntry creates a procedural codex entry.
func (g *Generator) GenerateCodexEntry(cat CodexCategory, uniqueID string) *CodexEntry {
	title := g.generateCodexTitle(cat)
	text := g.generateCodexText(cat)

	return NewCodexEntry(uniqueID, cat, title, text, g.genre)
}

// generateCodexTitle creates a title for a codex entry.
func (g *Generator) generateCodexTitle(cat CodexCategory) string {
	titles := codexTitles[g.genre]
	if titles == nil {
		titles = codexTitles[engine.GenreFantasy]
	}
	catTitles := titles[cat]
	if len(catTitles) == 0 {
		return "Unknown Entry"
	}
	return seed.Choice(g.gen, catTitles)
}

var codexTitles = map[engine.GenreID]map[CodexCategory][]string{
	engine.GenreFantasy: {
		CodexHistory:    {"The Age of Kingdoms", "The Great War", "Rise of Magic", "Fall of the Ancients"},
		CodexFaction:    {"The Order of Light", "Shadow Guild", "Merchant League", "Wilderness Rangers"},
		CodexRoute:      {"The King's Road", "Mountain Pass", "Forest Trail", "River Crossing"},
		CodexCreature:   {"Dragons of Old", "Forest Spirits", "Undead Horrors", "Beast Tribes"},
		CodexTechnology: {"Runecraft", "Enchantments", "Alchemy", "Siege Weapons"},
	},
	engine.GenreScifi: {
		CodexHistory:    {"First Contact", "Colony Wars", "Corporate Expansion", "The Singularity"},
		CodexFaction:    {"United Colonies", "Pirate Clans", "Corporate Alliance", "Research Council"},
		CodexRoute:      {"Trade Lanes", "Void Routes", "Jump Points", "Patrol Corridors"},
		CodexCreature:   {"Alien Species", "Void Entities", "Mutant Strains", "AI Constructs"},
		CodexTechnology: {"FTL Drives", "Energy Weapons", "Neural Links", "Terraforming"},
	},
	engine.GenreHorror: {
		CodexHistory:    {"Day Zero", "The Spread", "Government Response", "Society's Fall"},
		CodexFaction:    {"Military Remnants", "Survivor Groups", "Raider Gangs", "Research Teams"},
		CodexRoute:      {"Safe Corridors", "Danger Zones", "Evacuation Routes", "Supply Lines"},
		CodexCreature:   {"The Infected", "Mutations", "Pack Behavior", "Evolution Stages"},
		CodexTechnology: {"Cure Research", "Weapons", "Fortifications", "Communication"},
	},
	engine.GenreCyberpunk: {
		CodexHistory:    {"Corporate Rise", "The Riots", "Net Wars", "Social Collapse"},
		CodexFaction:    {"Megacorps", "Street Gangs", "Fixers", "Resistance"},
		CodexRoute:      {"Corporate Zones", "Combat Zones", "Free Zones", "Underground"},
		CodexCreature:   {"Cyborgs", "AI Entities", "Mutants", "Enhanced Humans"},
		CodexTechnology: {"Cyberware", "Netrunning", "Weapons", "Vehicles"},
	},
	engine.GenrePostapoc: {
		CodexHistory:    {"The Old World", "The Bombs", "First Years", "New Societies"},
		CodexFaction:    {"Vault Dwellers", "Raider Hordes", "Trader Guilds", "Settlements"},
		CodexRoute:      {"Trade Routes", "Danger Zones", "Water Sources", "Radiation Maps"},
		CodexCreature:   {"Mutants", "Ferals", "Rad-Beasts", "Pre-War Creatures"},
		CodexTechnology: {"Pre-War Tech", "Jury-Rigging", "Weapons", "Vehicles"},
	},
}

// generateCodexText creates text for a codex entry.
func (g *Generator) generateCodexText(cat CodexCategory) string {
	texts := codexTexts[g.genre]
	if texts == nil {
		texts = codexTexts[engine.GenreFantasy]
	}
	catTexts := texts[cat]
	if len(catTexts) == 0 {
		return "Information on this topic is incomplete."
	}
	return seed.Choice(g.gen, catTexts)
}

var codexTexts = map[engine.GenreID]map[CodexCategory][]string{
	engine.GenreFantasy: {
		CodexHistory: {
			"The old kingdoms rose and fell before the current age. Their ruins dot the landscape.",
			"Magic once flowed freely through the land. Now it is rare and dangerous.",
			"The great war changed everything. Cities fell, and the survivors scattered.",
		},
		CodexFaction: {
			"This faction controls territory through a combination of military might and trade.",
			"Known for their strict code of honor, they are respected and feared in equal measure.",
			"Operating in the shadows, their influence extends further than most realize.",
		},
		CodexRoute: {
			"This route is well-traveled but not without danger. Bandits are known to operate here.",
			"An ancient path, used by traders for generations. Many landmarks guide the way.",
			"Treacherous terrain makes this route difficult. Only the experienced attempt it.",
		},
		CodexCreature: {
			"These creatures are rarely seen but always feared. Encounters are often fatal.",
			"Once common, they now lurk in the wild places. Approach with extreme caution.",
			"Legend speaks of their origins in dark magic. Truth and myth are hard to separate.",
		},
		CodexTechnology: {
			"This craft has been passed down through generations. Its secrets are closely guarded.",
			"Ancient techniques, rediscovered in recent years. Understanding remains incomplete.",
			"Powerful but dangerous. Those who misuse it rarely survive the consequences.",
		},
	},
	engine.GenreScifi: {
		CodexHistory: {
			"Humanity spread across the stars, but unity proved elusive. Conflicts shaped the sector.",
			"The corporations rose to fill the power vacuum. Their influence is everywhere.",
			"First contact changed everything. We were not alone, and they knew it before we did.",
		},
		CodexFaction: {
			"Corporate interests drive their every action. Profit is the only morality.",
			"Military discipline holds them together. Loyalty is demanded and returned.",
			"Operating beyond the law, they provide services others won't. For a price.",
		},
		CodexRoute: {
			"This lane is heavily trafficked. Patrols are common, but so are pirates.",
			"A dangerous route, but the shortest distance between key systems.",
			"Uncharted territory. Navigation data is incomplete. Expect surprises.",
		},
		CodexCreature: {
			"Alien biology is poorly understood. Contact protocols exist but are often inadequate.",
			"Artificial life forms operate by their own logic. Predicting their actions is difficult.",
			"Mutations from colonial experiments. Not natural, but very real threats.",
		},
		CodexTechnology: {
			"Cutting-edge technology, reverse-engineered from alien artifacts.",
			"Military-grade equipment, restricted to authorized personnel.",
			"Prototype systems with unknown reliability. Use at your own risk.",
		},
	},
	engine.GenreHorror: {
		CodexHistory: {
			"It started without warning. Within weeks, civilization collapsed.",
			"The government's response was too slow. By the time they acted, it was too late.",
			"Survivors banded together, but trust was in short supply. Many didn't make it.",
		},
		CodexFaction: {
			"Military remnants maintain order in their territory. Their methods are harsh.",
			"Survivors who work together. Community is their strength, but resources are scarce.",
			"Those who prey on others. In the new world, some chose to become predators.",
		},
		CodexRoute: {
			"Known to be relatively safe. Survivors have cleared and maintained it.",
			"Dangerous territory. The infected concentration here is high.",
			"Evacuation route from the early days. Now abandoned and overgrown.",
		},
		CodexCreature: {
			"The infected display disturbing behaviors. They're changing, adapting.",
			"Something in the infection is rewriting them. What they're becoming is unclear.",
			"Not all infected are the same. Some mutations are far more dangerous.",
		},
		CodexTechnology: {
			"Medical research continues in scattered facilities. Progress is slow.",
			"Improvised weapons and fortifications. Survival demands adaptation.",
			"Pre-outbreak technology. Still useful, but parts are getting scarce.",
		},
	},
	engine.GenreCyberpunk: {
		CodexHistory: {
			"The corporations filled the void when governments failed. Now they are the government.",
			"The riots changed everything. The old social contract burned with the buildings.",
			"The net was supposed to connect us. Instead, it became another battlefield.",
		},
		CodexFaction: {
			"Corporate power is absolute in their zones. Resistance is futile without inside help.",
			"Street gangs control the margins. Loyalty is earned, not demanded.",
			"Fixers connect all the pieces. Information is their currency.",
		},
		CodexRoute: {
			"Corporate zones are safe but monitored. Every movement is tracked.",
			"Combat zones are lawless. Corporate security doesn't enter without heavy support.",
			"Underground routes bypass checkpoints. Knowledge of them is valuable.",
		},
		CodexCreature: {
			"Cyborgs blur the line between human and machine. Some embrace it, others resist.",
			"AI constructs operate in the net and sometimes beyond. Their goals are their own.",
			"Corporate experiments created things that were never meant to exist.",
		},
		CodexTechnology: {
			"Cyberware ranges from basic to military-grade. Installation is risky without proper clinics.",
			"Netrunning is an art and a science. ICE will kill the unwary.",
			"Street tech is improvised but effective. Chrome from parts.",
		},
	},
	engine.GenrePostapoc: {
		CodexHistory: {
			"The bombs fell without warning. What survived the blast faced worse in the aftermath.",
			"The first years were the hardest. Those who made it through became the new humanity.",
			"Civilization is rebuilding, slowly. Old knowledge is precious and rare.",
		},
		CodexFaction: {
			"Vault dwellers possess pre-war knowledge. Their isolation preserved what others lost.",
			"Raiders take what they want. Might makes right in the wastes.",
			"Traders connect scattered settlements. The lifeblood of the new economy.",
		},
		CodexRoute: {
			"Trade routes are maintained by common agreement. Attacking traders hurts everyone.",
			"Radiation hot zones must be avoided or crossed quickly with proper protection.",
			"Water sources determine where settlements can exist. Control water, control the region.",
		},
		CodexCreature: {
			"Radiation changed everything it touched. The creatures that emerged are dangerous.",
			"Feral humans who've lost their minds to radiation or trauma. Unpredictable.",
			"Pre-war animals, mutated beyond recognition. Some are useful, most are deadly.",
		},
		CodexTechnology: {
			"Pre-war tech is priceless. Working examples are worth dying for—or killing for.",
			"Jury-rigged solutions from salvaged parts. Not pretty, but functional.",
			"Water purification is the most important technology. Everything else is secondary.",
		},
	},
}

// GenerateEnvironment creates a full environmental manager with content.
func (g *Generator) GenerateEnvironment(mapWidth, mapHeight, numInscriptions, numDiscoveries int) *EnvironmentalManager {
	manager := NewEnvironmentalManager(g.genre)

	// Generate inscriptions at random positions
	for i := 0; i < numInscriptions; i++ {
		x := g.gen.Range(1, mapWidth-1)
		y := g.gen.Range(1, mapHeight-1)
		manager.AddInscription(g.GenerateInscription(x, y))
	}

	// Generate discoveries at random positions
	for i := 0; i < numDiscoveries; i++ {
		x := g.gen.Range(1, mapWidth-1)
		y := g.gen.Range(1, mapHeight-1)
		manager.AddDiscovery(g.GenerateDiscovery(x, y))
	}

	// Generate codex entries for each category
	for _, cat := range AllCodexCategories() {
		for i := 0; i < 3; i++ { // 3 entries per category
			id := fmt.Sprintf("%d_%d", cat, i)
			manager.Codex.AddEntry(g.GenerateCodexEntry(cat, id))
		}
	}

	return manager
}
