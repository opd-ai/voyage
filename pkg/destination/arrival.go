package destination

import (
	"strings"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// ArrivalCeremony contains the full arrival narrative.
type ArrivalCeremony struct {
	Opening    string
	Narrative  string
	Reflection string
	Closing    string
}

// generateArrivalCeremony creates the full arrival text.
func (d *Destination) generateArrivalCeremony() string {
	ceremony := d.buildArrivalCeremony()
	var sb strings.Builder
	sb.WriteString(ceremony.Opening)
	sb.WriteString("\n\n")
	sb.WriteString(ceremony.Narrative)
	sb.WriteString("\n\n")
	sb.WriteString(ceremony.Reflection)
	sb.WriteString("\n\n")
	sb.WriteString(ceremony.Closing)
	return sb.String()
}

// buildArrivalCeremony constructs the ceremony components.
func (d *Destination) buildArrivalCeremony() *ArrivalCeremony {
	return &ArrivalCeremony{
		Opening:    d.generateOpening(),
		Narrative:  d.generateNarrative(),
		Reflection: d.generateReflection(),
		Closing:    d.generateClosing(),
	}
}

// generateOpening creates the opening line of the ceremony.
func (d *Destination) generateOpening() string {
	openings := arrivalOpenings[d.genre]
	return seed.Choice(d.seedGen, openings)
}

// generateNarrative creates the main narrative text.
func (d *Destination) generateNarrative() string {
	narratives := arrivalNarratives[d.genre][d.Type]
	return seed.Choice(d.seedGen, narratives)
}

// generateReflection creates a reflective moment.
func (d *Destination) generateReflection() string {
	reflections := arrivalReflections[d.genre]
	return seed.Choice(d.seedGen, reflections)
}

// generateClosing creates the closing statement.
func (d *Destination) generateClosing() string {
	closings := arrivalClosings[d.genre]
	return seed.Choice(d.seedGen, closings)
}

// arrivalOpenings maps genre to opening statements.
var arrivalOpenings = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"The weary travelers look upon their destination at last.",
		"After countless leagues, the journey's end stands before them.",
		"Songs will be sung of this moment for generations to come.",
	},
	engine.GenreScifi: {
		"The ship's computer announces: Destination reached.",
		"Through the viewport, their goal fills the screen.",
		"After light-years of travel, the voyage concludes.",
	},
	engine.GenreHorror: {
		"Through the mist, the destination materializes.",
		"The long nightmare may finally be ending.",
		"Hope and dread mingle as they arrive.",
	},
	engine.GenreCyberpunk: {
		"The neon glow of safety washes over them.",
		"Connection established. They've made it.",
		"The grid welcomes them with flickering lights.",
	},
	engine.GenrePostapoc: {
		"Against all odds, they stand at the gates of salvation.",
		"The wasteland journey ends here.",
		"Survivors no more—they've become something more.",
	},
}

// arrivalNarratives maps genre to destination type to narrative texts.
var arrivalNarratives = map[engine.GenreID]map[DestinationType][]string{
	engine.GenreFantasy: {
		City: {
			"The great city gates swing wide, trumpets heralding the travelers' arrival. Crowds gather to witness those who braved the perils of the wild.",
			"Spires and towers reach for the clouds as the party passes beneath the ancient archway. Civilization embraces them once more.",
		},
		Sanctuary: {
			"The sacred grove welcomes them with gentle birdsong and dappled light. Priests emerge to offer blessings and healing.",
			"Ancient trees bow in the breeze as if greeting old friends. The grove's magic soothes their weary spirits.",
		},
		Treasure: {
			"Mountains of gold and gems glitter in the torchlight. The dragon's hoard is theirs at last.",
			"Legendary artifacts line the walls. Their quest for treasure has reached its triumphant conclusion.",
		},
		Escape: {
			"The portal shimmers with otherworldly light, offering passage to distant lands and new beginnings.",
			"Ancient runes flare to life as the gateway opens. Their old life fades as a new one beckons.",
		},
		Settlement: {
			"Simple cottages and friendly faces greet the weary travelers. The village opens its arms to them.",
			"Smoke rises from hearths, children play in the streets. Peace at last.",
		},
	},
	engine.GenreScifi: {
		City: {
			"The station's massive docking bay hums with activity. Thousands of lives intersect in this orbital metropolis.",
			"Viewports reveal the station's majesty—a testament to human achievement among the stars.",
		},
		Sanctuary: {
			"The orbital haven's environmental systems provide perfect comfort. Medical bays stand ready to heal.",
			"Gentle artificial gravity welcomes them. Here, at last, they can truly rest.",
		},
		Treasure: {
			"The derelict's cargo hold contains wealth beyond measure—rare elements, advanced tech, priceless data.",
			"Salvage drones begin cataloging the haul. This find will change everything.",
		},
		Escape: {
			"The jump gate's energy field beckons, promising escape to civilized space.",
			"Coordinates locked. One jump and they'll be light-years from danger.",
		},
		Settlement: {
			"The colony dome shimmers in the alien sun. A new world, a new home.",
			"Settlers wave from the hydroponics bay. Humanity persists, even here.",
		},
	},
	engine.GenreHorror: {
		City: {
			"The abandoned town is silent, but perhaps that silence means safety from what hunts in the dark.",
			"Boarded windows and empty streets. But at least the horrors have not followed them here.",
		},
		Sanctuary: {
			"The church's blessed walls hold firm against the darkness. Holy symbols ward the doorways.",
			"Within these consecrated walls, the nightmares cannot reach them. For now.",
		},
		Treasure: {
			"The crypt's relics pulse with power—dangerous, but perhaps enough to fight back.",
			"Ancient artifacts of light and darkness. Tools to survive what comes next.",
		},
		Escape: {
			"The bridge stretches into fog, but beyond it lies freedom from this cursed land.",
			"One crossing. One chance. Leave the horrors behind forever.",
		},
		Settlement: {
			"Other survivors huddle together. Safety in numbers against the night.",
			"Haunted eyes meet theirs—fellow victims who have endured and survived.",
		},
	},
	engine.GenreCyberpunk: {
		City: {
			"The megacity sector hums with data and commerce. Anonymous among millions, finally safe.",
			"Neon reflections paint their faces. The corporate web stretches everywhere, but so do the shadows.",
		},
		Sanctuary: {
			"The safehouse's faraday cage blocks all signals. Off-grid, off-map, finally free.",
			"No cameras, no trackers, no corporate eyes. Just peace.",
		},
		Treasure: {
			"The data vault's contents scroll across screens—corporate secrets worth millions.",
			"Encrypted files unlock one by one. Leverage against the megacorps at last.",
		},
		Escape: {
			"The black market port bustles with smugglers. New identities, new lives, just credits away.",
			"Forged papers in hand, they prepare to disappear into the sprawl.",
		},
		Settlement: {
			"The free zone operates beyond corporate law. Rough but honest, a place to truly live.",
			"Hackers, fixers, and dreamers share space here. The underground welcomes its own.",
		},
	},
	engine.GenrePostapoc: {
		City: {
			"The survivor camp's walls are built from the bones of the old world. Strong. Defiant.",
			"Fires burn in the camp. Life persists despite everything that tried to end it.",
		},
		Sanctuary: {
			"The bunker's blast doors seal behind them. Pre-war technology still functions here.",
			"Clean air, clean water, stored supplies. The bunker preserved a piece of the old world.",
		},
		Treasure: {
			"The cache holds medicine, ammunition, fuel—currency of the new world.",
			"Pre-war supplies in perfect condition. Wealth beyond measure in the wasteland.",
		},
		Escape: {
			"The evacuation point is still staffed. Against all odds, rescue arrived.",
			"Helicopters wait on the pad. The nightmare ends here.",
		},
		Settlement: {
			"The reclaimed zone shows signs of real recovery—crops, buildings, hope.",
			"From ashes, a new community rises. They've found a place to call home.",
		},
	},
}

// arrivalReflections maps genre to reflective statements.
var arrivalReflections = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"They think of companions lost and trials overcome. Every step was worth it.",
		"The journey changed them in ways they are only beginning to understand.",
		"Heroes are not born—they are forged by the path they walk.",
	},
	engine.GenreScifi: {
		"The ship's log records the journey's end. Data will be analyzed for years.",
		"Light-years traveled, anomalies cataloged, lives transformed.",
		"Space changes everyone who truly faces its vastness.",
	},
	engine.GenreHorror: {
		"They survived. But will the memories ever fade?",
		"Some scars are invisible. Some nightmares never truly end.",
		"Safety is relative. They know now what lurks in the dark.",
	},
	engine.GenreCyberpunk: {
		"The data is secured, the mission complete. But the game never truly ends.",
		"They've become ghosts in the machine, invisible but powerful.",
		"In the neon shadows, a new legend takes shape.",
	},
	engine.GenrePostapoc: {
		"They look back at the wasteland they crossed. Never again.",
		"Survival wasn't enough—they had to find something worth living for.",
		"The old world is dead. But they carry its memory forward.",
	},
}

// arrivalClosings maps genre to closing statements.
var arrivalClosings = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"And so the journey ends, but the story continues...",
		"The final chapter of this tale closes. New adventures await.",
		"Thus concludes the voyage. May fortune favor future travels.",
	},
	engine.GenreScifi: {
		"Mission complete. Standby for next assignment.",
		"The voyage ends. The exploration continues.",
		"Destination reached. What wonders await?",
	},
	engine.GenreHorror: {
		"For now, they rest. But evil never truly sleeps.",
		"The story pauses, but the darkness remains.",
		"Survival is its own reward. For now.",
	},
	engine.GenreCyberpunk: {
		"Connection terminated. Until the next run.",
		"The job is done. Time to disappear.",
		"End of line. But the network never forgets.",
	},
	engine.GenrePostapoc: {
		"They made it. But the wasteland always waits.",
		"One journey ends. The struggle to rebuild begins.",
		"Survivors. Pioneers. The future starts here.",
	},
}
