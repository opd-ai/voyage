package narrative

import (
	"fmt"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// Generator creates procedural narrative content.
type Generator struct {
	gen   *seed.Generator
	genre engine.GenreID
}

// NewGenerator creates a new narrative generator.
func NewGenerator(masterSeed int64, genre engine.GenreID) *Generator {
	return &Generator{
		gen:   seed.NewGenerator(masterSeed, "narrative"),
		genre: genre,
	}
}

// SetGenre updates the generator's genre.
func (g *Generator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// GenerateStoryArc creates a complete story arc for a run.
func (g *Generator) GenerateStoryArc(crewNames []string, destName string) *StoryArc {
	arc := NewStoryArc(g.genre)

	// Generate three-act structure
	arc.AddBeat(g.generateDepartureBeat())
	arc.AddBeat(g.generateMidJourneyBeat())
	arc.AddBeat(g.generateArrivalBeat(destName))

	// Generate recurring NPC
	arc.RecurringNPC = g.generateRecurringNPC()

	// Generate crew backstories
	for i, name := range crewNames {
		if g.gen.Chance(0.5) { // 50% chance each crew member has a backstory
			arc.AddCrewBackstory(g.generateCrewBackstory(i, name, destName))
		}
	}

	return arc
}

// generateDepartureBeat creates the opening crisis.
func (g *Generator) generateDepartureBeat() *StoryBeat {
	titles := departureTitles[g.genre]
	if titles == nil {
		titles = departureTitles[engine.GenreFantasy]
	}
	descriptions := departureDescriptions[g.genre]
	if descriptions == nil {
		descriptions = departureDescriptions[engine.GenreFantasy]
	}

	return &StoryBeat{
		Act:         ActDeparture,
		Title:       seed.Choice(g.gen, titles),
		Description: seed.Choice(g.gen, descriptions),
	}
}

var departureTitles = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"The Burning Village",
		"Dark Portents",
		"The Royal Decree",
		"Curse of the Ancients",
	},
	engine.GenreScifi: {
		"Distress Signal",
		"Station Compromised",
		"The Last Transmission",
		"Emergency Evacuation",
	},
	engine.GenreHorror: {
		"Day Zero",
		"The Outbreak",
		"First Blood",
		"They're Here",
	},
	engine.GenreCyberpunk: {
		"The Burn Notice",
		"Corporate Betrayal",
		"Zeroed Out",
		"The Setup",
	},
	engine.GenrePostapoc: {
		"The Bombs Fell",
		"Water Wars Begin",
		"Exodus",
		"The Horde Approaches",
	},
}

var departureDescriptions = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"Your home is destroyed. Only flight can save the survivors.",
		"A prophet foretells doom unless you reach the sacred city.",
		"The king demands you deliver a vital message, or face execution.",
		"An ancient curse awakens. Only the distant temple holds the cure.",
	},
	engine.GenreScifi: {
		"The station's reactor is failing. You must evacuate immediately.",
		"A distress call from deep space. Someone needs rescue.",
		"The last transmission spoke of something terrible. You must investigate.",
		"Emergency protocols activated. Get your crew to safety.",
	},
	engine.GenreHorror: {
		"It started without warning. Now you must run or die.",
		"The infection spreads fast. Find the safe zone before it's too late.",
		"They came in the night. Only a few of you made it out.",
		"The screaming started at dawn. By noon, the city was lost.",
	},
	engine.GenreCyberpunk: {
		"Your identity was erased. Now you're running from everyone.",
		"The corporation turned on you. Time to disappear.",
		"Your accounts are frozen, your face is on every screen. Run.",
		"Someone set you up. Now you're public enemy number one.",
	},
	engine.GenrePostapoc: {
		"The bombs fell last week. Now the radiation forces you to move.",
		"The water dried up. Your settlement must find a new source.",
		"Raiders destroyed your home. The survivors must find sanctuary.",
		"A massive horde approaches. Stay and die, or flee and live.",
	},
}

// generateMidJourneyBeat creates the revelation.
func (g *Generator) generateMidJourneyBeat() *StoryBeat {
	titles := midJourneyTitles[g.genre]
	if titles == nil {
		titles = midJourneyTitles[engine.GenreFantasy]
	}
	descriptions := midJourneyDescriptions[g.genre]
	if descriptions == nil {
		descriptions = midJourneyDescriptions[engine.GenreFantasy]
	}

	return &StoryBeat{
		Act:         ActMidJourney,
		Title:       seed.Choice(g.gen, titles),
		Description: seed.Choice(g.gen, descriptions),
	}
}

var midJourneyTitles = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"The Hidden Truth",
		"Prophecy Revealed",
		"Secrets of the Past",
		"The True Enemy",
	},
	engine.GenreScifi: {
		"Data Decrypted",
		"The Signal's Source",
		"Hidden Agenda",
		"Conspiracy Unveiled",
	},
	engine.GenreHorror: {
		"Patient Zero",
		"The Source",
		"What They Became",
		"No Escape",
	},
	engine.GenreCyberpunk: {
		"The Real Target",
		"Deep Cover",
		"Memory Fragment",
		"The Connection",
	},
	engine.GenrePostapoc: {
		"Before the Fall",
		"The Vault's Secret",
		"What Really Happened",
		"The Last Message",
	},
}

var midJourneyDescriptions = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"You discover the curse was no accident—someone planned this.",
		"An ancient prophecy names one of your crew as the chosen one.",
		"The ruins reveal a forgotten history that changes everything.",
		"Your true enemy is not who you thought. The real threat awaits.",
	},
	engine.GenreScifi: {
		"The encrypted files reveal a corporate conspiracy reaching high.",
		"The signal wasn't a distress call—it was a warning.",
		"Someone on the inside has been feeding information to the enemy.",
		"The conspiracy goes deeper than anyone imagined.",
	},
	engine.GenreHorror: {
		"You find records of the first infection. This was engineered.",
		"The cure exists—but getting it means going back.",
		"They're not just sick. They're changing into something else.",
		"The safe zone isn't safe. It never was.",
	},
	engine.GenreCyberpunk: {
		"The burn wasn't random. You were the real target all along.",
		"A hidden implant reveals memories that aren't yours.",
		"Your fixer was playing both sides. The question is why.",
		"The data you're carrying is worth more than you knew.",
	},
	engine.GenrePostapoc: {
		"The bombs weren't an accident. Someone survived to plan this.",
		"A pre-war vault contains the truth about the apocalypse.",
		"The radiation isn't natural. It's being maintained deliberately.",
		"Someone is controlling the raiders. Orchestrating the chaos.",
	},
}

// generateArrivalBeat creates the final twist.
func (g *Generator) generateArrivalBeat(destName string) *StoryBeat {
	titles := arrivalTitles[g.genre]
	if titles == nil {
		titles = arrivalTitles[engine.GenreFantasy]
	}
	descriptions := arrivalDescriptions[g.genre]
	if descriptions == nil {
		descriptions = arrivalDescriptions[engine.GenreFantasy]
	}

	desc := seed.Choice(g.gen, descriptions)
	desc = fmt.Sprintf(desc, destName)

	return &StoryBeat{
		Act:         ActArrival,
		Title:       seed.Choice(g.gen, titles),
		Description: desc,
	}
}

var arrivalTitles = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"Journey's End",
		"The Final Gate",
		"Destiny Awaits",
		"The Last Trial",
	},
	engine.GenreScifi: {
		"Final Approach",
		"Destination Reached",
		"End of the Line",
		"The Last Jump",
	},
	engine.GenreHorror: {
		"Sanctuary",
		"The Last Stand",
		"Safe at Last?",
		"End of Days",
	},
	engine.GenreCyberpunk: {
		"The Score",
		"Endgame",
		"Final Run",
		"Payoff",
	},
	engine.GenrePostapoc: {
		"The Promised Land",
		"New Beginning",
		"Home at Last",
		"The Settlement",
	},
}

var arrivalDescriptions = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"As you approach %s, you realize the legends were true—but not in the way you expected.",
		"The gates of %s open, but what lies within challenges everything you believed.",
		"%s stands before you. The final trial is not what you imagined.",
		"You've reached %s, but the journey's end is only the beginning.",
	},
	engine.GenreScifi: {
		"The scanners show %s dead ahead, but something about the readings is wrong.",
		"Docking at %s, you realize this destination holds more than shelter.",
		"%s appears on the viewport. The signal that led you here was only half the story.",
		"You've arrived at %s, but the welcome you receive is unexpected.",
	},
	engine.GenreHorror: {
		"The walls of %s come into view. But are you really safe here?",
		"You've made it to %s. The nightmare should be over—shouldn't it?",
		"%s stands before you. Those who survived the journey now face one final truth.",
		"At last, %s. But something feels wrong about this sanctuary.",
	},
	engine.GenreCyberpunk: {
		"The border checkpoint to %s lies ahead. Time to see if the price was worth it.",
		"You've reached %s. Now to collect what you're owed—if you can trust anyone.",
		"%s sprawls before you. But freedom comes with its own cost.",
		"The gates of %s open. The final play begins now.",
	},
	engine.GenrePostapoc: {
		"The walls of %s rise from the wasteland. Is this truly the promised land?",
		"You've crossed the wastes to reach %s. But is it everything you hoped?",
		"%s stands in the distance. What you find there will change everything.",
		"At last, %s. The journey ends, but a new chapter begins.",
	},
}

// generateRecurringNPC creates a named character that reappears.
func (g *Generator) generateRecurringNPC() *RecurringNPC {
	role := seed.Choice(g.gen, AllNPCRoles())
	name := g.generateNPCName()
	desc := g.generateNPCDescription(role)

	npc := NewRecurringNPC(name, role, desc, g.genre)
	npc.Dialogues = g.generateNPCDialogues(role)

	return npc
}

// generateNPCName creates a name for the recurring NPC.
func (g *Generator) generateNPCName() string {
	names := npcNames[g.genre]
	if names == nil {
		names = npcNames[engine.GenreFantasy]
	}
	return seed.Choice(g.gen, names)
}

var npcNames = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"Mira the Wanderer", "Lord Varek", "The Gray Stranger",
		"Sera of the Woods", "Old Thorne", "The Hooded Figure",
	},
	engine.GenreScifi: {
		"Agent Nine", "The Broker", "Commander Vex",
		"Doctor Omega", "The Pilot", "Zero",
	},
	engine.GenreHorror: {
		"The Doctor", "Murphy", "The Priest",
		"Silent Sam", "The Watcher", "Gracie",
	},
	engine.GenreCyberpunk: {
		"Mr. White", "Ghost", "The Fixer",
		"Neon", "Cipher", "The Broker",
	},
	engine.GenrePostapoc: {
		"The Wanderer", "Old Pete", "Dust Devil",
		"The Prophet", "Ironside", "The Mechanic",
	},
}

// generateNPCDescription creates a description based on role.
func (g *Generator) generateNPCDescription(role RecurringNPCRole) string {
	descriptions := npcDescriptions[g.genre]
	if descriptions == nil {
		descriptions = npcDescriptions[engine.GenreFantasy]
	}
	roleDescs := descriptions[role]
	if len(roleDescs) == 0 {
		return "A mysterious figure whose motives are unclear."
	}
	return seed.Choice(g.gen, roleDescs)
}

var npcDescriptions = map[engine.GenreID]map[RecurringNPCRole][]string{
	engine.GenreFantasy: {
		RoleFriend: {
			"A wise traveler who offers guidance and aid.",
			"A former knight who has taken an interest in your quest.",
			"A mysterious benefactor who appears when needed most.",
		},
		RoleNemesis: {
			"A dark figure who seems to follow your every step.",
			"An old rival who seeks to claim what you're after.",
			"A villain whose goals oppose everything you stand for.",
		},
		RoleAmbiguous: {
			"A merchant who seems to know too much.",
			"A stranger whose help always comes with a price.",
			"Someone who appears at crossroads with cryptic advice.",
		},
	},
	engine.GenreScifi: {
		RoleFriend: {
			"A rogue operative who shares valuable intelligence.",
			"A smuggler who owes you a favor—and always pays debts.",
			"A scientist with knowledge that could save everyone.",
		},
		RoleNemesis: {
			"A corporate agent who wants what you're carrying.",
			"A hunter who has been tracking you across space.",
			"A zealot who believes you must be stopped at any cost.",
		},
		RoleAmbiguous: {
			"A broker who deals in information—and secrets.",
			"An AI construct with its own mysterious agenda.",
			"A diplomat whose loyalties shift with the wind.",
		},
	},
	engine.GenreHorror: {
		RoleFriend: {
			"A survivor who has knowledge of safe routes.",
			"A former scientist who might know about a cure.",
			"Someone who lost everything but keeps helping others.",
		},
		RoleNemesis: {
			"A survivor who believes you're the cause of this.",
			"A raider leader who wants your supplies—and your lives.",
			"Someone who was changed by the outbreak in dark ways.",
		},
		RoleAmbiguous: {
			"A doctor who's been exposed but hasn't turned yet.",
			"A soldier whose orders might not align with survival.",
			"Someone who knows more about the outbreak than they say.",
		},
	},
	engine.GenreCyberpunk: {
		RoleFriend: {
			"A fixer who's been in your corner since the beginning.",
			"A netrunner who provides intel and support.",
			"A street doc who patches you up, no questions asked.",
		},
		RoleNemesis: {
			"A corporate agent who was paid to end you.",
			"A rival runner who wants what you're carrying.",
			"Someone you wronged who's come to collect.",
		},
		RoleAmbiguous: {
			"A fixer who plays all sides of every deal.",
			"An AI that contacts you with unclear motives.",
			"A corpo exec who might be helping—or using you.",
		},
	},
	engine.GenrePostapoc: {
		RoleFriend: {
			"A wanderer who knows the safe paths through the wastes.",
			"A mechanic who keeps your transport running.",
			"A trader who always has what you need—at fair prices.",
		},
		RoleNemesis: {
			"A warlord who claims the territory you must cross.",
			"A raider chief who's marked you for death.",
			"Someone who blames you for their losses.",
		},
		RoleAmbiguous: {
			"A prophet who speaks in riddles about the future.",
			"A vault dweller with pre-war knowledge and secrets.",
			"A trader whose deals always have hidden terms.",
		},
	},
}

// generateNPCDialogues creates dialogue lines for the NPC.
func (g *Generator) generateNPCDialogues(role RecurringNPCRole) []string {
	dialogues := npcDialogues[g.genre]
	if dialogues == nil {
		dialogues = npcDialogues[engine.GenreFantasy]
	}
	roleDialogues := dialogues[role]
	if len(roleDialogues) == 0 {
		return []string{"..."}
	}

	// Select 3-4 dialogues
	count := 3 + g.gen.Intn(2)
	if count > len(roleDialogues) {
		count = len(roleDialogues)
	}

	result := make([]string, 0, count)
	g.gen.Shuffle(len(roleDialogues), func(i, j int) {
		roleDialogues[i], roleDialogues[j] = roleDialogues[j], roleDialogues[i]
	})
	return append(result, roleDialogues[:count]...)
}

var npcDialogues = map[engine.GenreID]map[RecurringNPCRole][]string{
	engine.GenreFantasy: {
		RoleFriend: {
			"We meet again, friends. The road ahead is perilous.",
			"I've been watching your progress. You're doing well.",
			"Take this. You'll need it more than I.",
			"The path you seek lies to the north. Trust me.",
		},
		RoleNemesis: {
			"You think you can outrun fate? I am fate.",
			"Everything you've worked for will be mine.",
			"We meet again. This time, you won't escape.",
			"Your journey ends here, one way or another.",
		},
		RoleAmbiguous: {
			"Interesting. You've made it this far.",
			"I could help you—for a price.",
			"The truth isn't what you think it is.",
			"Choose wisely. Not all paths lead forward.",
		},
	},
	engine.GenreScifi: {
		RoleFriend: {
			"I've got intel that might help you.",
			"The coordinates I'm sending should get you through.",
			"I owed you one. Now we're even.",
			"Stay alive out there. We need people like you.",
		},
		RoleNemesis: {
			"Your contract has been terminated—permanently.",
			"Did you really think you could escape?",
			"The corporation thanks you for your service.",
			"Nothing personal. Just business.",
		},
		RoleAmbiguous: {
			"Information has a price. What are you offering?",
			"I know things that could change everything.",
			"Trust is a luxury in the void.",
			"My motives are my own. Shall we deal or not?",
		},
	},
	engine.GenreHorror: {
		RoleFriend: {
			"Thank god, you're still alive.",
			"I found a route that's mostly clear.",
			"We have to stick together. It's the only way.",
			"I won't leave anyone behind.",
		},
		RoleNemesis: {
			"You brought this on all of us.",
			"Survival of the fittest. And you're not fit.",
			"Your supplies. Your weapons. Hand them over.",
			"In this world, only the ruthless survive.",
		},
		RoleAmbiguous: {
			"I've seen what happens. You should know.",
			"I'm not sure I trust anyone anymore.",
			"The infection... there's something they didn't tell us.",
			"I can help you, but you won't like what I know.",
		},
	},
	engine.GenreCyberpunk: {
		RoleFriend: {
			"Got your back, choom. Always.",
			"Here's what you need. Don't ask where I got it.",
			"The net's saying things about you. Watch yourself.",
			"When this is over, drinks are on me.",
		},
		RoleNemesis: {
			"You cost me everything. Time to return the favor.",
			"The bounty on your head just went up.",
			"Nothing personal. Cred is cred.",
			"End of the line, runner.",
		},
		RoleAmbiguous: {
			"I could sell you out. But where's the profit?",
			"Everyone's playing everyone. Including me.",
			"What I know is worth more than your life.",
			"Let's make a deal. We both win.",
		},
	},
	engine.GenrePostapoc: {
		RoleFriend: {
			"The wasteland's hard on everyone. Let me help.",
			"I know where there's clean water. Follow me.",
			"Your people? They're good folks. Worth saving.",
			"In the old world, I was nobody. Out here, I matter.",
		},
		RoleNemesis: {
			"This territory is mine. You don't belong here.",
			"The strong take. The weak give. Which are you?",
			"Water, food, fuel—everything has a price in blood.",
			"I've killed for less than what you're carrying.",
		},
		RoleAmbiguous: {
			"The old world is gone. Accept it.",
			"I've seen things in the wastes. Things that shouldn't be.",
			"Trust is a pre-war luxury.",
			"I'll trade knowledge for supplies. Fair deal?",
		},
	},
}

// generateCrewBackstory creates a backstory connecting crew to destination.
func (g *Generator) generateCrewBackstory(crewID int, name, destName string) *CrewBackstory {
	hooks := backstoryHooks[g.genre]
	if hooks == nil {
		hooks = backstoryHooks[engine.GenreFantasy]
	}
	fulls := backstoryFulls[g.genre]
	if fulls == nil {
		fulls = backstoryFulls[engine.GenreFantasy]
	}
	links := backstoryLinks[g.genre]
	if links == nil {
		links = backstoryLinks[engine.GenreFantasy]
	}

	hook := fmt.Sprintf(seed.Choice(g.gen, hooks), name)
	full := fmt.Sprintf(seed.Choice(g.gen, fulls), name)
	link := fmt.Sprintf(seed.Choice(g.gen, links), name, destName)

	return NewCrewBackstory(crewID, name, hook, full, link)
}

var backstoryHooks = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"%s seems troubled by memories of the past.",
		"%s carries an old keepsake but won't discuss it.",
		"%s knows more about the destination than they let on.",
	},
	engine.GenreScifi: {
		"%s receives encrypted messages they don't share.",
		"%s's background check has strange gaps.",
		"%s seems familiar with the destination coordinates.",
	},
	engine.GenreHorror: {
		"%s has nightmares about something they won't discuss.",
		"%s was in the area when the outbreak started.",
		"%s lost someone and won't say how.",
	},
	engine.GenreCyberpunk: {
		"%s has implants that don't match their story.",
		"%s's ID doesn't quite add up.",
		"%s knows people at the destination—maybe too well.",
	},
	engine.GenrePostapoc: {
		"%s has scars they won't explain.",
		"%s knows the old world better than most.",
		"%s has been to the destination before.",
	},
}

var backstoryFulls = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"%s reveals they once lived in the city you're seeking.",
		"%s confesses to fleeing a terrible crime years ago.",
		"%s admits their family was involved in the curse's origins.",
	},
	engine.GenreScifi: {
		"%s was part of the original research team at the station.",
		"%s worked for the corporation now hunting you.",
		"%s's family was affected by the same conspiracy.",
	},
	engine.GenreHorror: {
		"%s was bitten but somehow didn't turn.",
		"%s worked at the facility where it all started.",
		"%s's loved one might still be alive—infected.",
	},
	engine.GenreCyberpunk: {
		"%s used to work for the people now hunting you.",
		"%s has memories that were supposed to be wiped.",
		"%s was the original target—you were just collateral.",
	},
	engine.GenrePostapoc: {
		"%s survived the bombs in a vault near the destination.",
		"%s's family founded the settlement you're seeking.",
		"%s knows what really caused the apocalypse.",
	},
}

var backstoryLinks = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"%s's family has an ancient connection to %s.",
		"%s was exiled from %s years ago.",
		"%s seeks redemption that can only be found in %s.",
	},
	engine.GenreScifi: {
		"%s was born on %s before the disaster.",
		"%s's research was stolen and taken to %s.",
		"%s's missing sibling was last seen heading to %s.",
	},
	engine.GenreHorror: {
		"%s's family might still be alive in %s.",
		"%s helped build %s before everything fell apart.",
		"%s has unfinished business with someone in %s.",
	},
	engine.GenreCyberpunk: {
		"%s has a contact in %s who might help—or betray.",
		"%s's real identity is known in %s.",
		"%s stole something from %s and fears recognition.",
	},
	engine.GenrePostapoc: {
		"%s's ancestors built the walls of %s.",
		"%s left %s years ago and never thought they'd return.",
		"%s knows secrets about %s that could change everything.",
	},
}
