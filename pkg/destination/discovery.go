package destination

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// DiscoveryEvent represents an event triggered when approaching a destination.
type DiscoveryEvent struct {
	Phase       DiscoveryPhase
	Title       string
	Description string
	Destination *Destination
}

// generateDiscoveryEvent creates an event for the current phase.
func (d *Destination) generateDiscoveryEvent() *DiscoveryEvent {
	event := &DiscoveryEvent{
		Phase:       d.Phase,
		Destination: d,
	}
	event.Title = d.getEventTitle()
	event.Description = d.getEventDescription()
	return event
}

// getEventTitle returns a phase-appropriate title.
// Falls back to fantasy genre if the current genre is not found (L-013).
func (d *Destination) getEventTitle() string {
	genreTitles, ok := discoveryTitles[d.genre]
	if !ok {
		genreTitles = discoveryTitles[engine.GenreFantasy]
	}
	titles, ok := genreTitles[d.Phase]
	if !ok || len(titles) == 0 {
		return "Discovery"
	}
	return seed.Choice(d.seedGen, titles)
}

// getEventDescription returns a phase-appropriate description.
// Falls back to fantasy genre if the current genre is not found (L-013).
func (d *Destination) getEventDescription() string {
	genreDescs, ok := discoveryDescriptions[d.genre]
	if !ok {
		genreDescs = discoveryDescriptions[engine.GenreFantasy]
	}
	descriptions, ok := genreDescs[d.Phase]
	if !ok || len(descriptions) == 0 {
		return "Something significant has occurred."
	}
	return seed.Choice(d.seedGen, descriptions)
}

// discoveryTitles maps genre to phase to possible titles.
var discoveryTitles = map[engine.GenreID]map[DiscoveryPhase][]string{
	engine.GenreFantasy: {
		Distant:  {"A Distant Glimmer", "Smoke on the Horizon", "Tales Confirmed"},
		Signs:    {"Waymarkers Found", "Travelers' Rest", "Signs of Life"},
		Approach: {"The Gates Appear", "Journey's End in Sight", "Final Stretch"},
		Arrival:  {"Arrival at Last", "Journey's End", "We Have Arrived"},
	},
	engine.GenreScifi: {
		Distant:  {"Long-Range Contact", "Sensor Anomaly", "Signal Detected"},
		Signs:    {"Beacon Located", "Transmission Received", "Station Signature"},
		Approach: {"Docking Clearance", "Final Approach", "Coordinates Locked"},
		Arrival:  {"Docking Complete", "Arrival Confirmed", "Mission Success"},
	},
	engine.GenreHorror: {
		Distant:  {"Lights in the Fog", "Distant Bells", "A Bad Omen"},
		Signs:    {"Warning Signs", "Abandoned Camps", "Fresh Graves"},
		Approach: {"The Mist Parts", "Shadows Deepen", "No Turning Back"},
		Arrival:  {"The End Awaits", "Sanctuary at Last", "Arrival"},
	},
	engine.GenreCyberpunk: {
		Distant:  {"Grid Ping", "Node Located", "Chatter Intercepted"},
		Signs:    {"Firewall Detected", "Signal Strength Rising", "Cache Found"},
		Approach: {"Perimeter Breach", "Access Point Located", "Final Push"},
		Arrival:  {"Connection Established", "We're In", "Mission Complete"},
	},
	engine.GenrePostapoc: {
		Distant:  {"Smoke Sighted", "Radio Static", "Hope on the Wind"},
		Signs:    {"Markers Found", "Recent Tracks", "Signs of Survivors"},
		Approach: {"Walls Visible", "Guards Spotted", "Almost There"},
		Arrival:  {"Safe at Last", "We Made It", "Sanctuary Found"},
	},
}

// discoveryDescriptions maps genre to phase to possible descriptions.
var discoveryDescriptions = map[engine.GenreID]map[DiscoveryPhase][]string{
	engine.GenreFantasy: {
		Distant: {
			"A faint shimmer catches the eye on the distant horizon.",
			"Campfire tales speak of this place. It seems they were true.",
			"Birds fly toward the same point. Something draws them.",
		},
		Signs: {
			"Ancient waymarkers line the path, carved with protective runes.",
			"A traveler's rest, recently used, shows others walk this path.",
			"The trees thin and farmland appears. Civilization is near.",
		},
		Approach: {
			"The destination reveals itself fully. Your journey nears its end.",
			"Guards patrol the outer reaches. You have been seen.",
			"The road widens into a proper thoroughfare. Almost there.",
		},
		Arrival: {
			"At last! The weary travelers have reached their destination.",
			"The gates open to welcome the road-worn party.",
			"Your epic journey comes to a triumphant close.",
		},
	},
	engine.GenreScifi: {
		Distant: {
			"Long-range sensors detect a mass anomaly in the target sector.",
			"A faint electromagnetic signature matches our destination.",
			"Subspace echoes suggest a large structure ahead.",
		},
		Signs: {
			"Navigation beacon locked. Automated transmission confirms identity.",
			"Traffic control frequencies detected. We're on the right path.",
			"Passive scans show energy signatures consistent with habitation.",
		},
		Approach: {
			"Visual confirmation. The station fills the viewport.",
			"Docking request transmitted. Awaiting clearance.",
			"Final approach vector calculated. Destination reached.",
		},
		Arrival: {
			"Docking clamps engaged. The journey is complete.",
			"Pressure equalized. You may disembark.",
			"Mission parameters satisfied. Welcome to your destination.",
		},
	},
	engine.GenreHorror: {
		Distant: {
			"Faint lights flicker through the fog. Something waits ahead.",
			"Distant bells toll a funeral dirge. Or is it a warning?",
			"The air grows cold. Your destination feels... wrong.",
		},
		Signs: {
			"Abandoned camps with supplies still laid out. Where did they go?",
			"Fresh graves line the road. Many did not make it this far.",
			"Warning signs in desperate handwriting. 'Turn back.'",
		},
		Approach: {
			"The mist parts reluctantly, revealing your destination.",
			"Shadows seem to follow your every step. Almost there.",
			"The point of no return. Whatever awaits, you must face it.",
		},
		Arrival: {
			"You've arrived. But are you safe, or merely trapped?",
			"The destination stands before you. May it offer the sanctuary you seek.",
			"At last. Whether this is salvation or doom remains to be seen.",
		},
	},
	engine.GenreCyberpunk: {
		Distant: {
			"A ping on the underground mesh. Someone knows we're coming.",
			"Grid traffic increases. The node is real.",
			"Intercepted chatter mentions a safehouse matching our target.",
		},
		Signs: {
			"Corporate firewalls detected. We're getting close to something big.",
			"Dead drops along the route. Someone left breadcrumbs.",
			"Signal strength increasing. The cache is near.",
		},
		Approach: {
			"Visual on the perimeter. Security looks manageable.",
			"Access point located. One more push and we're in.",
			"Facial recognition cameras ahead. Stay alert.",
		},
		Arrival: {
			"Connection established. We made it through the corp-net.",
			"The safehouse door slides open. Welcome to the underground.",
			"Mission complete. Time to collect and get ghost.",
		},
	},
	engine.GenrePostapoc: {
		Distant: {
			"Smoke rises on the horizon. Could be survivors. Could be raiders.",
			"Static breaks on the radio. Someone is broadcasting.",
			"The wasteland stretches on, but hope glimmers ahead.",
		},
		Signs: {
			"Painted markers on the road. Someone wants travelers to find this place.",
			"Recent tracks in the dust. Others have passed this way.",
			"An old supply cache, recently restocked. Survivors are near.",
		},
		Approach: {
			"Walls of scrap and concrete rise ahead. A settlement.",
			"Guards on the perimeter. They've seen us.",
			"The final stretch of blasted earth. Almost safe.",
		},
		Arrival: {
			"The gates creak open. Safety at last.",
			"Against all odds, you've reached sanctuary.",
			"The wasteland journey ends. A new chapter begins.",
		},
	},
}
