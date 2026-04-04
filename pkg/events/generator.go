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
