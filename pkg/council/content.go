package council

import (
	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
)

// decisionDescriptions maps genre to decision type to descriptions.
var decisionDescriptions = map[engine.GenreID]map[DecisionType][]string{
	engine.GenreFantasy: {
		DecisionRoute: {"The road forks ahead. One path is treacherous but swift.", "Two trails diverge in the ancient wood."},
		DecisionCamp:  {"Night approaches. Do we press on or make camp?", "The weary party debates whether to rest."},
		DecisionTrade: {"A merchant offers a deal that seems too good.", "A trader presents an unusual proposition."},
		DecisionFight: {"Bandits block our path. We must decide quickly.", "An enemy patrol has spotted us."},
		DecisionSplit: {"The situation calls for some to scout ahead.", "We could cover more ground if we divide."},
	},
	engine.GenreScifi: {
		DecisionRoute: {"Two jump routes available: one through uncharted space.", "Navigation offers alternative courses."},
		DecisionCamp:  {"Power reserves low. Full stop or continue?", "The crew needs rest but time is critical."},
		DecisionTrade: {"An unusual contract has been proposed.", "A broker offers high-risk cargo."},
		DecisionFight: {"Hostile vessel detected. Engagement options?", "Pirates have locked weapons on us."},
		DecisionSplit: {"Away team deployment may be necessary.", "A shuttle mission could gather intel."},
	},
	engine.GenreHorror: {
		DecisionRoute: {"Two doors. One screams, one whispers.", "The path splits—neither feels safe."},
		DecisionCamp:  {"We can't keep running. But can we rest here?", "Exhaustion battles fear for control."},
		DecisionTrade: {"Something offers a bargain. Its smile is wrong.", "A stranger knows what we seek."},
		DecisionFight: {"It's found us. Run or fight?", "We can't outrun it forever."},
		DecisionSplit: {"We need to search faster. Split up?", "Cover more ground separately?"},
	},
	engine.GenreCyberpunk: {
		DecisionRoute: {"Two routes: corp sector or the sprawl.", "Matrix path or meatspace running."},
		DecisionCamp:  {"Lay low in this safehouse or keep moving?", "Heat's on. Bunker down or ghost?"},
		DecisionTrade: {"A fixer has a job. Smells like a setup.", "Contract's lucrative. Too lucrative."},
		DecisionFight: {"Corp security inbound. Engage or extract?", "Gangers want a piece. Pay or play?"},
		DecisionSplit: {"Need a distraction. Who stays behind?", "Solo run could get us in clean."},
	},
	engine.GenrePostapoc: {
		DecisionRoute: {"The safe road is long. The short one is lethal.", "Radiation or raiders—pick your poison."},
		DecisionCamp:  {"This ruin looks defensible. Rest or move?", "We need sleep but can't afford it."},
		DecisionTrade: {"Strangers want to trade. Trust them?", "They have supplies we need badly."},
		DecisionFight: {"Scavengers spotted us. Show strength?", "They look hungry. We look weak."},
		DecisionSplit: {"Scout the area or stay together?", "Need eyes on that ridge."},
	},
}

// riskyDescriptions maps genre to risky option descriptions.
var riskyDescriptions = map[engine.GenreID]map[DecisionType][]string{
	engine.GenreFantasy: {
		DecisionRoute: {"Take the mountain pass—faster but treacherous.", "Cut through the cursed woods."},
		DecisionCamp:  {"March through the night.", "Push on despite exhaustion."},
		DecisionTrade: {"Accept the deal.", "Take the merchant's offer."},
		DecisionFight: {"Stand and fight.", "Engage the enemy."},
		DecisionSplit: {"Divide the party.", "Send scouts ahead."},
	},
	engine.GenreScifi: {
		DecisionRoute: {"Plot course through the anomaly.", "Take the uncharted jump."},
		DecisionCamp:  {"Maintain full burn.", "Keep engines hot."},
		DecisionTrade: {"Accept the contract.", "Take the cargo."},
		DecisionFight: {"Target locked. Fire.", "Engage hostiles."},
		DecisionSplit: {"Launch the away team.", "Deploy the shuttle."},
	},
	engine.GenreHorror: {
		DecisionRoute: {"Through the darkness.", "Face what waits."},
		DecisionCamp:  {"Keep moving. No rest.", "Outrun the nightmare."},
		DecisionTrade: {"Accept the bargain.", "Take the offer."},
		DecisionFight: {"Fight it.", "Make a stand."},
		DecisionSplit: {"Split up to cover ground.", "Separate and search."},
	},
	engine.GenreCyberpunk: {
		DecisionRoute: {"Through corp territory.", "Hardline through the net."},
		DecisionCamp:  {"Keep running.", "Stay mobile."},
		DecisionTrade: {"Take the job.", "Sign the contract."},
		DecisionFight: {"Light them up.", "Go loud."},
		DecisionSplit: {"Solo run.", "One goes in alone."},
	},
	engine.GenrePostapoc: {
		DecisionRoute: {"Through the hot zone.", "Risk the rad storm."},
		DecisionCamp:  {"Keep moving.", "No rest for the wicked."},
		DecisionTrade: {"Deal with them.", "Trust the strangers."},
		DecisionFight: {"Show them we're not prey.", "Attack first."},
		DecisionSplit: {"Scout separately.", "Cover more ground."},
	},
}

// safeDescriptions maps genre to safe option descriptions.
var safeDescriptions = map[engine.GenreID]map[DecisionType][]string{
	engine.GenreFantasy: {
		DecisionRoute: {"Take the longer road.", "Follow the king's highway."},
		DecisionCamp:  {"Make camp and rest.", "Rest while we can."},
		DecisionTrade: {"Decline the offer.", "This deal seems wrong."},
		DecisionFight: {"Avoid confrontation.", "Find another way."},
		DecisionSplit: {"Stay together.", "Strength in numbers."},
	},
	engine.GenreScifi: {
		DecisionRoute: {"Use established routes.", "Follow the beacon path."},
		DecisionCamp:  {"Full stop for repairs.", "Power down and rest."},
		DecisionTrade: {"Decline the contract.", "Pass on this one."},
		DecisionFight: {"Evasive maneuvers.", "Run silent, run deep."},
		DecisionSplit: {"Keep the crew together.", "No away team."},
	},
	engine.GenreHorror: {
		DecisionRoute: {"The longer way.", "Avoid the darkness."},
		DecisionCamp:  {"Barricade and rest.", "We need to recover."},
		DecisionTrade: {"Refuse the bargain.", "Nothing is worth that."},
		DecisionFight: {"Run. Just run.", "Hide and survive."},
		DecisionSplit: {"Stay together.", "Never separate."},
	},
	engine.GenreCyberpunk: {
		DecisionRoute: {"Through the sprawl.", "The long way around."},
		DecisionCamp:  {"Go to ground.", "Safehouse protocol."},
		DecisionTrade: {"Walk away.", "Not worth the risk."},
		DecisionFight: {"Ghost out.", "Tactical retreat."},
		DecisionSplit: {"Team stays intact.", "No solo ops."},
	},
	engine.GenrePostapoc: {
		DecisionRoute: {"The long road.", "Avoid the hot zones."},
		DecisionCamp:  {"Fortify and rest.", "Secure the perimeter."},
		DecisionTrade: {"Keep our distance.", "We can't trust anyone."},
		DecisionFight: {"Stay hidden.", "Not worth the bullets."},
		DecisionSplit: {"Together or not at all.", "No splitting up."},
	},
}

// voteReasonings maps genre to trait to option to reasoning text.
var voteReasonings = map[engine.GenreID]map[crew.Trait]map[VoteOption][]string{
	engine.GenreFantasy: {
		crew.TraitBrave: {
			OptionRisky: {"Fortune favors the bold!", "We didn't come this far to turn back."},
			OptionSafe:  {"Even heroes need rest."},
		},
		crew.TraitCautious: {
			OptionRisky: {"If we must..."},
			OptionSafe:  {"Fools rush in. Let's be wise.", "The clever survive where the brave fall."},
		},
		crew.TraitOptimistic: {
			OptionRisky: {"I have a good feeling about this!", "Things will work out!"},
			OptionSafe:  {"Better to live and see tomorrow's sunrise."},
		},
		crew.TraitPessimistic: {
			OptionRisky: {"We're doomed either way."},
			OptionSafe:  {"Less ways to die this way.", "At least we'll live a little longer."},
		},
		crew.TraitGreedy: {
			OptionRisky: {"Think of the treasure!", "The reward outweighs the risk."},
			OptionSafe:  {"Dead men spend no gold."},
		},
		crew.TraitGenerous: {
			OptionRisky: {"For the good of all."},
			OptionSafe:  {"We must protect each other.", "Everyone's safety matters."},
		},
		crew.TraitStoic: {
			OptionRisky: {"Logic dictates action."},
			OptionSafe:  {"The rational choice is clear.", "Emotion clouds judgment."},
		},
		crew.TraitEmotional: {
			OptionRisky: {"My heart says yes!", "I feel it in my soul!"},
			OptionSafe:  {"I'm scared... but maybe that's okay."},
		},
		crew.TraitNavigator: {
			OptionRisky: {"I know these paths. Trust me.", "This way is faster."},
			OptionSafe:  {"The safe road is still a road."},
		},
		crew.TraitScavenger: {
			OptionRisky: {"More risk, more reward.", "I've found good things in bad places."},
			OptionSafe:  {"Sometimes survival is enough."},
		},
	},
	engine.GenreScifi: {
		crew.TraitBrave: {
			OptionRisky: {"Engage!", "We have the advantage."},
			OptionSafe:  {"Strategic withdrawal acknowledged."},
		},
		crew.TraitCautious: {
			OptionRisky: {"Running the calculations..."},
			OptionSafe:  {"Probability favors caution.", "Risk assessment says no."},
		},
		crew.TraitOptimistic: {
			OptionRisky: {"The odds could be worse!", "We've beaten worse."},
			OptionSafe:  {"Living to fight another day is winning."},
		},
		crew.TraitPessimistic: {
			OptionRisky: {"Might as well go out fighting."},
			OptionSafe:  {"Delay the inevitable.", "Every minute counts."},
		},
		crew.TraitGreedy: {
			OptionRisky: {"Credits don't earn themselves.", "Maximum profit potential."},
			OptionSafe:  {"Can't spend creds if you're vaporized."},
		},
		crew.TraitGenerous: {
			OptionRisky: {"For the crew."},
			OptionSafe:  {"Everyone makes it home.", "No casualties."},
		},
		crew.TraitStoic: {
			OptionRisky: {"Logic indicates this course."},
			OptionSafe:  {"Statistical analysis supports this.", "Vulcan logic."},
		},
		crew.TraitEmotional: {
			OptionRisky: {"I've got a feeling!", "Trust your gut!"},
			OptionSafe:  {"Something feels wrong about this."},
		},
		crew.TraitNavigator: {
			OptionRisky: {"My calculations show a path.", "Trust my navigation."},
			OptionSafe:  {"Established routes exist for a reason."},
		},
		crew.TraitScavenger: {
			OptionRisky: {"Salvage opportunities await.", "Debris fields hide treasures."},
			OptionSafe:  {"Dead crews don't salvage anything."},
		},
	},
	engine.GenreHorror: {
		crew.TraitBrave: {
			OptionRisky: {"We can't run forever.", "Face your fears."},
			OptionSafe:  {"Even I know when to retreat."},
		},
		crew.TraitCautious: {
			OptionRisky: {"If there's no other way..."},
			OptionSafe:  {"Stay quiet. Stay alive.", "Don't attract attention."},
		},
		crew.TraitOptimistic: {
			OptionRisky: {"Maybe it's not that bad?", "We'll be okay!"},
			OptionSafe:  {"Hope lies in survival."},
		},
		crew.TraitPessimistic: {
			OptionRisky: {"We're all going to die anyway."},
			OptionSafe:  {"At least this way we die slower.", "Every breath is borrowed."},
		},
		crew.TraitGreedy: {
			OptionRisky: {"Worth the risk.", "There might be something valuable."},
			OptionSafe:  {"Dead people own nothing."},
		},
		crew.TraitGenerous: {
			OptionRisky: {"I'll go first to protect you."},
			OptionSafe:  {"We protect each other.", "No one left behind."},
		},
		crew.TraitStoic: {
			OptionRisky: {"Fear is the mind-killer."},
			OptionSafe:  {"Rational decisions save lives.", "Panic kills."},
		},
		crew.TraitEmotional: {
			OptionRisky: {"I can't take the tension anymore!", "Just do something!"},
			OptionSafe:  {"I'm too scared. Please.", "I can't lose anyone else."},
		},
		crew.TraitNavigator: {
			OptionRisky: {"I know a shortcut.", "Trust my sense of direction."},
			OptionSafe:  {"The marked path is safer."},
		},
		crew.TraitScavenger: {
			OptionRisky: {"Could be supplies there.", "Nothing ventured..."},
			OptionSafe:  {"We have enough to survive."},
		},
	},
	engine.GenreCyberpunk: {
		crew.TraitBrave: {
			OptionRisky: {"Let's flatline some corpos.", "No guts, no glory."},
			OptionSafe:  {"Even I know when the odds are slag."},
		},
		crew.TraitCautious: {
			OptionRisky: {"Against my better judgment..."},
			OptionSafe:  {"Heat's too high.", "Risk assessment: negative."},
		},
		crew.TraitOptimistic: {
			OptionRisky: {"Our luck's gotta turn!", "This is our moment!"},
			OptionSafe:  {"Live to run another day."},
		},
		crew.TraitPessimistic: {
			OptionRisky: {"We're fragged either way."},
			OptionSafe:  {"Slower death is still death, but I'll take it."},
		},
		crew.TraitGreedy: {
			OptionRisky: {"Think of the eddies!", "Maximum payout."},
			OptionSafe:  {"Can't spend creds in the morgue."},
		},
		crew.TraitGenerous: {
			OptionRisky: {"For the team."},
			OptionSafe:  {"We all make it out.", "Nobody gets left behind."},
		},
		crew.TraitStoic: {
			OptionRisky: {"Cold calculation says go."},
			OptionSafe:  {"The math doesn't lie.", "Probability favors retreat."},
		},
		crew.TraitEmotional: {
			OptionRisky: {"Frag it, I'm in!", "Let's do this!"},
			OptionSafe:  {"I've got a bad feeling.", "Something's off."},
		},
		crew.TraitNavigator: {
			OptionRisky: {"I know the grid.", "Trust my routes."},
			OptionSafe:  {"Standard protocols exist for a reason."},
		},
		crew.TraitScavenger: {
			OptionRisky: {"Tech to salvage.", "Loot awaits."},
			OptionSafe:  {"Dead runners don't collect gear."},
		},
	},
	engine.GenrePostapoc: {
		crew.TraitBrave: {
			OptionRisky: {"We didn't survive this long to be cowards.", "Take the fight to them."},
			OptionSafe:  {"Even survivors know when to retreat."},
		},
		crew.TraitCautious: {
			OptionRisky: {"Only if absolutely necessary..."},
			OptionSafe:  {"Caution keeps us breathing.", "The wasteland punishes recklessness."},
		},
		crew.TraitOptimistic: {
			OptionRisky: {"It'll work out!", "The world's still got hope."},
			OptionSafe:  {"Tomorrow's another day if we live."},
		},
		crew.TraitPessimistic: {
			OptionRisky: {"Everything's already ruined anyway."},
			OptionSafe:  {"Why accelerate the end?", "The world's bad enough."},
		},
		crew.TraitGreedy: {
			OptionRisky: {"Think of the salvage!", "Resources waiting to be claimed."},
			OptionSafe:  {"Dead don't need supplies."},
		},
		crew.TraitGenerous: {
			OptionRisky: {"For our community."},
			OptionSafe:  {"Everyone lives. That's what matters.", "We protect our own."},
		},
		crew.TraitStoic: {
			OptionRisky: {"Survival requires risks."},
			OptionSafe:  {"Calculate, then act.", "Emotion is a luxury."},
		},
		crew.TraitEmotional: {
			OptionRisky: {"I can't stand waiting!", "Act now!"},
			OptionSafe:  {"Please, I'm tired of losing people."},
		},
		crew.TraitNavigator: {
			OptionRisky: {"I know these wastes.", "My route will work."},
			OptionSafe:  {"Known paths are safer."},
		},
		crew.TraitScavenger: {
			OptionRisky: {"Good picking ahead.", "Treasure in the ruins."},
			OptionSafe:  {"Can't scavenge if we're dead."},
		},
	},
}

// sceneOpenings maps genre to opening text.
var sceneOpenings = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"The party gathers around the crackling campfire.",
		"Beneath the stars, voices rise in debate.",
		"The campfire casts dancing shadows as the council begins.",
	},
	engine.GenreScifi: {
		"The crew assembles on the bridge.",
		"Holographic displays illuminate worried faces.",
		"The captain calls for a briefing.",
	},
	engine.GenreHorror: {
		"Huddled together, they argue in hushed, frantic whispers.",
		"Fear makes every voice sharp.",
		"The darkness presses in as they debate.",
	},
	engine.GenreCyberpunk: {
		"The team jacks in for a quick meet.",
		"AR overlays flicker as the runners talk biz.",
		"Smoke fills the room as the discussion begins.",
	},
	engine.GenrePostapoc: {
		"Around the guttering bonfire, the group debates.",
		"The wasteland wind howls as voices rise.",
		"Survivors gather to make a choice.",
	},
}

// discussionSnippets maps genre to generic discussion lines.
var discussionSnippets = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"Voices rise and fall like the flames.",
		"Old arguments surface briefly.",
		"Someone stokes the fire nervously.",
	},
	engine.GenreScifi: {
		"Data scrolls across screens as they talk.",
		"Someone pulls up a tactical overlay.",
		"The hum of the ship fills silences.",
	},
	engine.GenreHorror: {
		"Every shadow seems to move.",
		"Someone checks the doors again.",
		"A noise outside freezes everyone.",
	},
	engine.GenreCyberpunk: {
		"Someone pings the latest corp movements.",
		"Cred totals flash in ARs.",
		"The fixer's contact blinks online.",
	},
	engine.GenrePostapoc: {
		"The wind carries distant sounds.",
		"Someone counts their remaining rounds.",
		"Water rations are checked nervously.",
	},
}

// sceneClosings maps genre to closing text.
var sceneClosings = map[engine.GenreID][]string{
	engine.GenreFantasy: {
		"The decision is made. The fire crackles in agreement.",
		"And so the path is chosen.",
		"Fate awaits down the chosen road.",
	},
	engine.GenreScifi: {
		"The captain makes the call.",
		"Course laid in.",
		"The decision is logged.",
	},
	engine.GenreHorror: {
		"For better or worse, it's decided.",
		"No turning back now.",
		"Whatever happens, it's done.",
	},
	engine.GenreCyberpunk: {
		"The job is set.",
		"Time to run.",
		"Delta in three. Go.",
	},
	engine.GenrePostapoc: {
		"The wasteland waits for no one.",
		"Survival demands action.",
		"One way or another, they move.",
	},
}
