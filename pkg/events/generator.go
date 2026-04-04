package events

import "github.com/opd-ai/voyage/pkg/engine"

// EventTemplate is a reusable event structure.
type EventTemplate struct {
	Title       string
	Description string
	Choices     []ChoiceTemplate
}

// ChoiceTemplate is a reusable choice structure.
type ChoiceTemplate struct {
	Text    string
	Outcome EventOutcome
}

// Weather event templates by genre
var weatherTemplates = map[engine.GenreID][]EventTemplate{
	engine.GenreFantasy: {
		{
			Title:       "Storm Approaches",
			Description: "Dark clouds gather on the horizon. A fierce storm is approaching.",
			Choices: []ChoiceTemplate{
				{Text: "Press on through the storm", Outcome: EventOutcome{MoraleDelta: -5, TimeAdvance: 1}},
				{Text: "Find shelter and wait", Outcome: EventOutcome{FoodDelta: -5, WaterDelta: -3, TimeAdvance: 2}},
			},
		},
		{
			Title:       "Unseasonable Heat",
			Description: "The sun beats down mercilessly. Water supplies are being depleted faster than expected.",
			Choices: []ChoiceTemplate{
				{Text: "Ration water strictly", Outcome: EventOutcome{MoraleDelta: -8, WaterDelta: -5}},
				{Text: "Travel at night instead", Outcome: EventOutcome{FuelDelta: -3, TimeAdvance: 1}},
			},
		},
	},
	engine.GenreScifi: {
		{
			Title:       "Solar Flare Warning",
			Description: "Sensors detect an incoming solar flare. Radiation levels are rising.",
			Choices: []ChoiceTemplate{
				{Text: "Engage radiation shields", Outcome: EventOutcome{FuelDelta: -10}},
				{Text: "Find shelter behind a planetoid", Outcome: EventOutcome{TimeAdvance: 2}},
			},
		},
		{
			Title:       "Micro-Meteor Shower",
			Description: "A field of micro-meteors lies ahead.",
			Choices: []ChoiceTemplate{
				{Text: "Navigate through carefully", Outcome: EventOutcome{VesselDamage: 10, TimeAdvance: 1}},
				{Text: "Go around the field", Outcome: EventOutcome{FuelDelta: -15, TimeAdvance: 2}},
			},
		},
	},
	engine.GenreHorror: {
		{
			Title:       "Acid Rain",
			Description: "The clouds release a toxic downpour. The rain hisses on contact with metal.",
			Choices: []ChoiceTemplate{
				{Text: "Drive through it fast", Outcome: EventOutcome{VesselDamage: 15}},
				{Text: "Wait for it to pass", Outcome: EventOutcome{FoodDelta: -5, TimeAdvance: 2}},
			},
		},
	},
	engine.GenreCyberpunk: {
		{
			Title:       "Smog Alert",
			Description: "Toxic smog has descended on the district. Visibility is near zero.",
			Choices: []ChoiceTemplate{
				{Text: "Use IR sensors and push through", Outcome: EventOutcome{FuelDelta: -5, TimeAdvance: 1}},
				{Text: "Wait for the filters to clear", Outcome: EventOutcome{TimeAdvance: 2, MoraleDelta: -3}},
			},
		},
	},
	engine.GenrePostapoc: {
		{
			Title:       "Radiation Storm",
			Description: "A radioactive dust storm sweeps across the wasteland.",
			Choices: []ChoiceTemplate{
				{Text: "Seal the vehicle and wait", Outcome: EventOutcome{TimeAdvance: 3, FoodDelta: -5}},
				{Text: "Drive through with masks", Outcome: EventOutcome{CrewDamage: 10, MedicineDelta: -5}},
			},
		},
	},
}

// Encounter event templates by genre
var encounterTemplates = map[engine.GenreID][]EventTemplate{
	engine.GenreFantasy: {
		{
			Title:       "Travelers on the Road",
			Description: "A group of weary travelers approaches your caravan.",
			Choices: []ChoiceTemplate{
				{Text: "Share supplies and news", Outcome: EventOutcome{FoodDelta: -3, MoraleDelta: 5}},
				{Text: "Trade with them", Outcome: EventOutcome{CurrencyDelta: 10}},
				{Text: "Ignore them and move on", Outcome: EventOutcome{MoraleDelta: -2}},
			},
		},
		{
			Title:       "Bandits!",
			Description: "Armed figures emerge from hiding, demanding tribute.",
			Choices: []ChoiceTemplate{
				{Text: "Pay the toll", Outcome: EventOutcome{CurrencyDelta: -20}},
				{Text: "Fight them off", Outcome: EventOutcome{CrewDamage: 15, MoraleDelta: 5}},
				{Text: "Try to outrun them", Outcome: EventOutcome{FuelDelta: -10, VesselDamage: 5}},
			},
		},
	},
	engine.GenreScifi: {
		{
			Title:       "Distress Signal",
			Description: "Your sensors pick up a distress signal from a nearby ship.",
			Choices: []ChoiceTemplate{
				{Text: "Investigate and assist", Outcome: EventOutcome{FuelDelta: -5, MoraleDelta: 10}},
				{Text: "Ignore the signal", Outcome: EventOutcome{MoraleDelta: -5}},
				{Text: "Report to authorities", Outcome: EventOutcome{TimeAdvance: 1}},
			},
		},
	},
	engine.GenreHorror: {
		{
			Title:       "Survivors Ahead",
			Description: "A small group waves desperately from the roadside.",
			Choices: []ChoiceTemplate{
				{Text: "Stop and help them", Outcome: EventOutcome{FoodDelta: -5, MoraleDelta: 8}},
				{Text: "Drive past - it could be a trap", Outcome: EventOutcome{MoraleDelta: -5}},
			},
		},
	},
	engine.GenreCyberpunk: {
		{
			Title:       "Corporate Checkpoint",
			Description: "A megacorp security checkpoint blocks the road ahead.",
			Choices: []ChoiceTemplate{
				{Text: "Pay the bribe", Outcome: EventOutcome{CurrencyDelta: -30}},
				{Text: "Hack the systems", Outcome: EventOutcome{TimeAdvance: 1, MoraleDelta: 5}},
				{Text: "Find another route", Outcome: EventOutcome{FuelDelta: -15, TimeAdvance: 2}},
			},
		},
	},
	engine.GenrePostapoc: {
		{
			Title:       "Raider Gang",
			Description: "Armed raiders block the road, their vehicles forming a barricade.",
			Choices: []ChoiceTemplate{
				{Text: "Pay them off", Outcome: EventOutcome{CurrencyDelta: -15, FoodDelta: -5}},
				{Text: "Ram through the blockade", Outcome: EventOutcome{VesselDamage: 20, FuelDelta: -10}},
				{Text: "Negotiate passage", Outcome: EventOutcome{TimeAdvance: 1, MoraleDelta: -3}},
			},
		},
	},
}

// Discovery event templates
var discoveryTemplates = map[engine.GenreID][]EventTemplate{
	engine.GenreFantasy: {
		{
			Title:       "Abandoned Camp",
			Description: "You find an abandoned camp with supplies left behind.",
			Choices: []ChoiceTemplate{
				{Text: "Search thoroughly", Outcome: EventOutcome{FoodDelta: 10, WaterDelta: 5, TimeAdvance: 1}},
				{Text: "Take what's visible and go", Outcome: EventOutcome{FoodDelta: 5}},
			},
		},
		{
			Title:       "Hidden Spring",
			Description: "Your scout discovers a natural spring hidden among rocks.",
			Choices: []ChoiceTemplate{
				{Text: "Fill all water containers", Outcome: EventOutcome{WaterDelta: 20, TimeAdvance: 1}},
				{Text: "Take a quick drink and move on", Outcome: EventOutcome{WaterDelta: 5, MoraleDelta: 3}},
			},
		},
	},
	engine.GenreScifi: {
		{
			Title:       "Derelict Cargo Pod",
			Description: "Scanners detect an abandoned cargo pod floating nearby.",
			Choices: []ChoiceTemplate{
				{Text: "Dock and salvage", Outcome: EventOutcome{FuelDelta: 15, TimeAdvance: 1}},
				{Text: "Scan and move on", Outcome: EventOutcome{TimeAdvance: 0}},
			},
		},
	},
	engine.GenreHorror: {
		{
			Title:       "Abandoned Pharmacy",
			Description: "A looted pharmacy, but some supplies might remain.",
			Choices: []ChoiceTemplate{
				{Text: "Search inside", Outcome: EventOutcome{MedicineDelta: 15, TimeAdvance: 1, CrewDamage: 5}},
				{Text: "Too risky, move on", Outcome: EventOutcome{}},
			},
		},
	},
	engine.GenreCyberpunk: {
		{
			Title:       "Data Cache",
			Description: "An abandoned server node still has power.",
			Choices: []ChoiceTemplate{
				{Text: "Download and sell the data", Outcome: EventOutcome{CurrencyDelta: 25, TimeAdvance: 1}},
				{Text: "Too hot - leave it", Outcome: EventOutcome{}},
			},
		},
	},
	engine.GenrePostapoc: {
		{
			Title:       "Fuel Depot",
			Description: "Pre-war fuel storage tanks, partially intact.",
			Choices: []ChoiceTemplate{
				{Text: "Siphon what you can", Outcome: EventOutcome{FuelDelta: 20, TimeAdvance: 2}},
				{Text: "Check for traps first", Outcome: EventOutcome{FuelDelta: 10, TimeAdvance: 3}},
			},
		},
	},
}

// Hardship event templates
var hardshipTemplates = map[engine.GenreID][]EventTemplate{
	engine.GenreFantasy: {
		{
			Title:       "Fever Spreads",
			Description: "Several party members have fallen ill with a fever.",
			Choices: []ChoiceTemplate{
				{Text: "Use medicine to treat them", Outcome: EventOutcome{MedicineDelta: -10}},
				{Text: "Rest and hope for recovery", Outcome: EventOutcome{TimeAdvance: 2, CrewDamage: 10}},
			},
		},
		{
			Title:       "Wheel Broken",
			Description: "A wheel on the wagon has shattered.",
			Choices: []ChoiceTemplate{
				{Text: "Spend time repairing it", Outcome: EventOutcome{TimeAdvance: 2}},
				{Text: "Use spare parts", Outcome: EventOutcome{CurrencyDelta: -10}},
			},
		},
	},
	engine.GenreScifi: {
		{
			Title:       "Life Support Malfunction",
			Description: "The life support system is failing.",
			Choices: []ChoiceTemplate{
				{Text: "Emergency repairs", Outcome: EventOutcome{FuelDelta: -10, TimeAdvance: 1}},
				{Text: "Ration oxygen", Outcome: EventOutcome{MoraleDelta: -10, CrewDamage: 5}},
			},
		},
	},
	engine.GenreHorror: {
		{
			Title:       "Infection",
			Description: "A wound has become infected. Without treatment, it could be fatal.",
			Choices: []ChoiceTemplate{
				{Text: "Use precious medicine", Outcome: EventOutcome{MedicineDelta: -15}},
				{Text: "Cauterize the wound", Outcome: EventOutcome{CrewDamage: 15, MoraleDelta: -5}},
			},
		},
	},
	engine.GenreCyberpunk: {
		{
			Title:       "System Crash",
			Description: "Your vehicle's navigation systems have crashed.",
			Choices: []ChoiceTemplate{
				{Text: "Hire a tech to fix it", Outcome: EventOutcome{CurrencyDelta: -20}},
				{Text: "Navigate manually", Outcome: EventOutcome{TimeAdvance: 2, FuelDelta: -10}},
			},
		},
	},
	engine.GenrePostapoc: {
		{
			Title:       "Radiation Sickness",
			Description: "Exposure has made some crew members ill.",
			Choices: []ChoiceTemplate{
				{Text: "Administer anti-rad meds", Outcome: EventOutcome{MedicineDelta: -10}},
				{Text: "Rest and recover", Outcome: EventOutcome{TimeAdvance: 3, CrewDamage: 10}},
			},
		},
	},
}

// Windfall event templates
var windfallTemplates = map[engine.GenreID][]EventTemplate{
	engine.GenreFantasy: {
		{
			Title:       "Merchant Caravan",
			Description: "A wealthy merchant offers to trade at favorable rates.",
			Choices: []ChoiceTemplate{
				{Text: "Stock up on supplies", Outcome: EventOutcome{FoodDelta: 15, WaterDelta: 10, CurrencyDelta: -15}},
				{Text: "Sell excess cargo", Outcome: EventOutcome{CurrencyDelta: 25}},
				{Text: "Just exchange news", Outcome: EventOutcome{MoraleDelta: 5}},
			},
		},
	},
	engine.GenreScifi: {
		{
			Title:       "Favorable Trade Route",
			Description: "Updated charts reveal a more efficient route.",
			Choices: []ChoiceTemplate{
				{Text: "Take the new route", Outcome: EventOutcome{FuelDelta: 10, TimeAdvance: -1}},
				{Text: "Stick to the known path", Outcome: EventOutcome{}},
			},
		},
	},
	engine.GenreHorror: {
		{
			Title:       "Safe Haven",
			Description: "A fortified settlement offers sanctuary.",
			Choices: []ChoiceTemplate{
				{Text: "Rest and resupply", Outcome: EventOutcome{MoraleDelta: 15, FoodDelta: 10, TimeAdvance: 1}},
				{Text: "Trade and move on", Outcome: EventOutcome{CurrencyDelta: 10}},
			},
		},
	},
	engine.GenreCyberpunk: {
		{
			Title:       "Big Score",
			Description: "A contact tips you off about an easy job.",
			Choices: []ChoiceTemplate{
				{Text: "Take the job", Outcome: EventOutcome{CurrencyDelta: 40, TimeAdvance: 1}},
				{Text: "Too good to be true", Outcome: EventOutcome{MoraleDelta: -3}},
			},
		},
	},
	engine.GenrePostapoc: {
		{
			Title:       "Pre-War Cache",
			Description: "An untouched supply bunker!",
			Choices: []ChoiceTemplate{
				{Text: "Loot everything", Outcome: EventOutcome{FoodDelta: 20, MedicineDelta: 10, TimeAdvance: 2}},
				{Text: "Take what you need", Outcome: EventOutcome{FoodDelta: 10, MedicineDelta: 5, TimeAdvance: 1}},
			},
		},
	},
}

// HazardTemplates define genre-specific hazard events.
// These are thematic dangers unique to each setting.
var hazardTemplates = map[engine.GenreID][]EventTemplate{
	engine.GenreFantasy: {
		{
			Title:       "Magic Storm",
			Description: "Wild arcane energies swirl through the air. Lightning crackles with unnatural colors.",
			Choices: []ChoiceTemplate{
				{Text: "Shelter until it passes", Outcome: EventOutcome{TimeAdvance: 2, FoodDelta: -5}},
				{Text: "Use protective wards", Outcome: EventOutcome{MoraleDelta: -5}},
				{Text: "Push through quickly", Outcome: EventOutcome{CrewDamage: 15, VesselDamage: 10}},
			},
		},
		{
			Title:       "Cursed Grounds",
			Description: "An ancient battlefield. The dead do not rest easy here.",
			Choices: []ChoiceTemplate{
				{Text: "Offer tribute to the spirits", Outcome: EventOutcome{CurrencyDelta: -15, MoraleDelta: 5}},
				{Text: "Rush through at speed", Outcome: EventOutcome{FuelDelta: -10, MoraleDelta: -10}},
				{Text: "Take a longer detour", Outcome: EventOutcome{TimeAdvance: 3, FuelDelta: -5}},
			},
		},
	},
	engine.GenreScifi: {
		{
			Title:       "Asteroid Field",
			Description: "Dense asteroid clusters ahead. Navigation requires full attention.",
			Choices: []ChoiceTemplate{
				{Text: "Navigate through carefully", Outcome: EventOutcome{TimeAdvance: 2, FuelDelta: -5}},
				{Text: "Blast through with weapons", Outcome: EventOutcome{FuelDelta: -15, VesselDamage: 5}},
				{Text: "Chart a course around", Outcome: EventOutcome{TimeAdvance: 4, FuelDelta: -10}},
			},
		},
		{
			Title:       "Ion Storm",
			Description: "Electromagnetic interference blankets the region. Systems are flickering.",
			Choices: []ChoiceTemplate{
				{Text: "Power down non-essentials", Outcome: EventOutcome{TimeAdvance: 3}},
				{Text: "Boost shields and push through", Outcome: EventOutcome{FuelDelta: -20, VesselDamage: 10}},
			},
		},
	},
	engine.GenreHorror: {
		{
			Title:       "Zombie Horde",
			Description: "The shambling dead have gathered in massive numbers. They block the road ahead.",
			Choices: []ChoiceTemplate{
				{Text: "Ram through them", Outcome: EventOutcome{VesselDamage: 20, FuelDelta: -5}},
				{Text: "Draw them away with noise", Outcome: EventOutcome{TimeAdvance: 2, MoraleDelta: -5}},
				{Text: "Find another route", Outcome: EventOutcome{TimeAdvance: 4, FuelDelta: -15}},
			},
		},
		{
			Title:       "Infected Zone",
			Description: "Contamination levels spike. The air itself feels wrong.",
			Choices: []ChoiceTemplate{
				{Text: "Seal up and drive fast", Outcome: EventOutcome{FuelDelta: -10, TimeAdvance: 1}},
				{Text: "Use medical gear as protection", Outcome: EventOutcome{MedicineDelta: -10}},
				{Text: "Turn back", Outcome: EventOutcome{TimeAdvance: 3, MoraleDelta: -8}},
			},
		},
	},
	engine.GenreCyberpunk: {
		{
			Title:       "Netrunner Ambush",
			Description: "Your systems lock up. Someone's trying to jack your vehicle remotely.",
			Choices: []ChoiceTemplate{
				{Text: "Counter-hack them", Outcome: EventOutcome{TimeAdvance: 1, MoraleDelta: 5}},
				{Text: "Pay the ransom", Outcome: EventOutcome{CurrencyDelta: -30}},
				{Text: "Go manual and flee", Outcome: EventOutcome{VesselDamage: 10, FuelDelta: -10}},
			},
		},
		{
			Title:       "Corporate Drone Swarm",
			Description: "Security drones converge on your position. You've triggered something.",
			Choices: []ChoiceTemplate{
				{Text: "EMP burst", Outcome: EventOutcome{FuelDelta: -20, TimeAdvance: 1}},
				{Text: "Jam their signals", Outcome: EventOutcome{CurrencyDelta: -15, TimeAdvance: 2}},
				{Text: "Outrun them", Outcome: EventOutcome{FuelDelta: -15, VesselDamage: 15}},
			},
		},
	},
	engine.GenrePostapoc: {
		{
			Title:       "Radiation Storm",
			Description: "A wall of glowing dust sweeps across the wastes. Geiger counters scream.",
			Choices: []ChoiceTemplate{
				{Text: "Seal vehicle and wait", Outcome: EventOutcome{TimeAdvance: 4, FoodDelta: -5, WaterDelta: -5}},
				{Text: "Use rad meds and push through", Outcome: EventOutcome{MedicineDelta: -15, TimeAdvance: 1}},
				{Text: "Find shelter in ruins", Outcome: EventOutcome{TimeAdvance: 3, MoraleDelta: -5}},
			},
		},
		{
			Title:       "Mutant Swarm",
			Description: "Twisted creatures emerge from the wasteland, drawn by your engine noise.",
			Choices: []ChoiceTemplate{
				{Text: "Floor it and run", Outcome: EventOutcome{FuelDelta: -15, VesselDamage: 10}},
				{Text: "Stand and fight", Outcome: EventOutcome{CrewDamage: 15, TimeAdvance: 1, MoraleDelta: 5}},
				{Text: "Throw supplies to distract", Outcome: EventOutcome{FoodDelta: -10, WaterDelta: -5}},
			},
		},
	},
}

// CrewEventType identifies the type of crew-specific event.
type CrewEventType int

const (
	// CrewEventCrisis is a personal crisis requiring intervention.
	CrewEventCrisis CrewEventType = iota
	// CrewEventMilestone is a positive personal achievement.
	CrewEventMilestone
	// CrewEventSacrifice is an opportunity for heroic self-sacrifice.
	CrewEventSacrifice
)

// CrewEventTemplate includes a crew member placeholder.
type CrewEventTemplate struct {
	Type        CrewEventType
	Title       string // Use %s for crew member name
	Description string // Use %s for crew member name
	Choices     []ChoiceTemplate
}

// CrewEventTemplates define crew-specific events by genre.
var crewTemplates = map[engine.GenreID][]CrewEventTemplate{
	engine.GenreFantasy: {
		{
			Type:        CrewEventCrisis,
			Title:       "%s Falls Into Despair",
			Description: "%s is haunted by visions of home. Their spirit wavers.",
			Choices: []ChoiceTemplate{
				{Text: "Offer words of comfort", Outcome: EventOutcome{MoraleDelta: 10, TimeAdvance: 1}},
				{Text: "Give them space", Outcome: EventOutcome{MoraleDelta: -5}},
				{Text: "Share a drink and story", Outcome: EventOutcome{MoraleDelta: 15, WaterDelta: -2}},
			},
		},
		{
			Type:        CrewEventMilestone,
			Title:       "%s Discovers Their Purpose",
			Description: "%s has found inner strength. They move with new determination.",
			Choices: []ChoiceTemplate{
				{Text: "Celebrate their growth", Outcome: EventOutcome{MoraleDelta: 15}},
				{Text: "Ask them to share wisdom", Outcome: EventOutcome{MoraleDelta: 10, TimeAdvance: 1}},
			},
		},
		{
			Type:        CrewEventSacrifice,
			Title:       "%s Holds the Line",
			Description: "Danger approaches. %s steps forward to protect the group.",
			Choices: []ChoiceTemplate{
				{Text: "Accept their sacrifice", Outcome: EventOutcome{CrewDamage: 50, MoraleDelta: 20}},
				{Text: "Pull them back to safety", Outcome: EventOutcome{CrewDamage: 10, VesselDamage: 20, MoraleDelta: 5}},
			},
		},
	},
	engine.GenreScifi: {
		{
			Type:        CrewEventCrisis,
			Title:       "%s Questions the Mission",
			Description: "%s is showing signs of space fatigue. Their focus is slipping.",
			Choices: []ChoiceTemplate{
				{Text: "Run psychological protocols", Outcome: EventOutcome{MoraleDelta: 10, TimeAdvance: 1}},
				{Text: "Assign them light duties", Outcome: EventOutcome{MoraleDelta: 5, FuelDelta: -5}},
				{Text: "Ignore it - they'll adapt", Outcome: EventOutcome{MoraleDelta: -8}},
			},
		},
		{
			Type:        CrewEventMilestone,
			Title:       "%s Achieves Certification",
			Description: "%s has mastered a new skill. Their efficiency has improved.",
			Choices: []ChoiceTemplate{
				{Text: "Log it in their record", Outcome: EventOutcome{MoraleDelta: 10}},
				{Text: "Host a ceremony", Outcome: EventOutcome{MoraleDelta: 20, FoodDelta: -3}},
			},
		},
		{
			Type:        CrewEventSacrifice,
			Title:       "%s Volunteers for EVA",
			Description: "A critical repair is needed outside. %s volunteers for the dangerous task.",
			Choices: []ChoiceTemplate{
				{Text: "Accept the risk", Outcome: EventOutcome{CrewDamage: 40, VesselDamage: -20, MoraleDelta: 15}},
				{Text: "Find another way", Outcome: EventOutcome{FuelDelta: -20, TimeAdvance: 2}},
			},
		},
	},
	engine.GenreHorror: {
		{
			Type:        CrewEventCrisis,
			Title:       "%s Is Breaking Down",
			Description: "%s hasn't slept in days. The horrors are getting to them.",
			Choices: []ChoiceTemplate{
				{Text: "Sedate them", Outcome: EventOutcome{MedicineDelta: -5, MoraleDelta: 5}},
				{Text: "Talk them through it", Outcome: EventOutcome{MoraleDelta: 8, TimeAdvance: 1}},
				{Text: "They need to toughen up", Outcome: EventOutcome{MoraleDelta: -15}},
			},
		},
		{
			Type:        CrewEventMilestone,
			Title:       "%s Finds Resolve",
			Description: "%s has stared into the abyss and found their courage.",
			Choices: []ChoiceTemplate{
				{Text: "Acknowledge their strength", Outcome: EventOutcome{MoraleDelta: 15}},
				{Text: "Hope it lasts", Outcome: EventOutcome{MoraleDelta: 5}},
			},
		},
		{
			Type:        CrewEventSacrifice,
			Title:       "%s Creates a Distraction",
			Description: "They're closing in. %s offers to draw them away.",
			Choices: []ChoiceTemplate{
				{Text: "Let them do it", Outcome: EventOutcome{CrewDamage: 60, MoraleDelta: 25}},
				{Text: "No one gets left behind", Outcome: EventOutcome{CrewDamage: 20, FuelDelta: -15}},
			},
		},
	},
	engine.GenreCyberpunk: {
		{
			Type:        CrewEventCrisis,
			Title:       "%s Has Cyberpsychosis Symptoms",
			Description: "%s is twitching. The chrome is fighting their meat.",
			Choices: []ChoiceTemplate{
				{Text: "Get them to a ripperdoc", Outcome: EventOutcome{CurrencyDelta: -25, MoraleDelta: 10}},
				{Text: "Use suppressants", Outcome: EventOutcome{MedicineDelta: -10, MoraleDelta: 5}},
				{Text: "They'll flatline eventually anyway", Outcome: EventOutcome{MoraleDelta: -20}},
			},
		},
		{
			Type:        CrewEventMilestone,
			Title:       "%s Paid Off Their Debt",
			Description: "%s is finally free of their corporate obligations.",
			Choices: []ChoiceTemplate{
				{Text: "Throw them a party", Outcome: EventOutcome{MoraleDelta: 20, CurrencyDelta: -10}},
				{Text: "Good for them", Outcome: EventOutcome{MoraleDelta: 8}},
			},
		},
		{
			Type:        CrewEventSacrifice,
			Title:       "%s Jacks In",
			Description: "The ICE is brutal. %s offers to take the hit to break through.",
			Choices: []ChoiceTemplate{
				{Text: "Let them burn", Outcome: EventOutcome{CrewDamage: 45, CurrencyDelta: 30, MoraleDelta: 10}},
				{Text: "Find another angle", Outcome: EventOutcome{TimeAdvance: 2, CurrencyDelta: -15}},
			},
		},
	},
	engine.GenrePostapoc: {
		{
			Type:        CrewEventCrisis,
			Title:       "%s Remembers the Old World",
			Description: "%s found something from before. They can't stop crying.",
			Choices: []ChoiceTemplate{
				{Text: "Sit with them", Outcome: EventOutcome{MoraleDelta: 12, TimeAdvance: 1}},
				{Text: "Let them grieve alone", Outcome: EventOutcome{MoraleDelta: -3}},
				{Text: "Destroy the memento", Outcome: EventOutcome{MoraleDelta: -20}},
			},
		},
		{
			Type:        CrewEventMilestone,
			Title:       "%s Made Something Grow",
			Description: "%s coaxed life from the wasteland. A small plant, but it's hope.",
			Choices: []ChoiceTemplate{
				{Text: "Guard it carefully", Outcome: EventOutcome{MoraleDelta: 20}},
				{Text: "Eat it - calories are calories", Outcome: EventOutcome{FoodDelta: 2, MoraleDelta: -10}},
			},
		},
		{
			Type:        CrewEventSacrifice,
			Title:       "%s Stays Behind",
			Description: "Raiders are coming. %s will buy time.",
			Choices: []ChoiceTemplate{
				{Text: "Honor their choice", Outcome: EventOutcome{CrewDamage: 55, MoraleDelta: 25, FuelDelta: 10}},
				{Text: "Everyone runs together", Outcome: EventOutcome{CrewDamage: 15, VesselDamage: 25}},
			},
		},
	},
}

// GetCrewEventTemplates returns crew event templates for a genre.
func GetCrewEventTemplates(genre engine.GenreID) []CrewEventTemplate {
	return crewTemplates[genre]
}

// HazardVocabulary returns genre-specific hazard names for display.
func HazardVocabulary(genre engine.GenreID) []string {
	templates := hazardTemplates[genre]
	names := make([]string, len(templates))
	for i, t := range templates {
		names[i] = t.Title
	}
	return names
}

// AllEventTemplates returns all event templates for a given genre grouped by category.
func AllEventTemplates(genre engine.GenreID) map[EventCategory][]EventTemplate {
	return map[EventCategory][]EventTemplate{
		CategoryWeather:   weatherTemplates[genre],
		CategoryEncounter: encounterTemplates[genre],
		CategoryDiscovery: discoveryTemplates[genre],
		CategoryHardship:  hardshipTemplates[genre],
		CategoryWindfall:  windfallTemplates[genre],
		CategoryHazard:    hazardTemplates[genre],
		// CategoryCrew uses CrewEventTemplate, not EventTemplate
	}
}
