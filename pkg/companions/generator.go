package companions

import (
	"fmt"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// Generator creates procedural companions
type Generator struct {
	gen         *seed.Generator
	genre       engine.GenreID
	companionID int
	eventID     int
}

// NewGenerator creates a companion generator with the given seed and genre
func NewGenerator(masterSeed int64, genre engine.GenreID) *Generator {
	return &Generator{
		gen:         seed.NewGenerator(masterSeed, "companions"),
		genre:       genre,
		companionID: 0,
		eventID:     0,
	}
}

// SetGenre updates the generator's active genre
func (g *Generator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// GenerateCompanion creates a new companion with the specified role
func (g *Generator) GenerateCompanion(role CompanionRole) *Companion {
	g.companionID++

	name := g.generateName()
	title := RoleName(role, g.genre)

	companion := NewCompanion(g.companionID, name, title, role, g.genre)
	companion.Backstory = g.generateBackstory(name, role)
	companion.SkillLevel = 1 + g.gen.Intn(3) // Start at 1-3

	// Add 2-3 personality traits
	traitCount := 2 + g.gen.Intn(2)
	traits := g.selectTraits(traitCount)
	for _, trait := range traits {
		companion.AddTrait(trait)
	}

	// Generate abilities based on role
	abilities := g.generateAbilities(role)
	for _, ability := range abilities {
		companion.AddAbility(ability)
	}

	// Initial stats based on traits
	g.applyTraitModifiers(companion)

	return companion
}

// GenerateRandomCompanion creates a companion with a random role
func (g *Generator) GenerateRandomCompanion() *Companion {
	roles := AllCompanionRoles()
	role := roles[g.gen.Intn(len(roles))]
	return g.GenerateCompanion(role)
}

// GenerateCompanionEvent creates a personality-driven event for a companion
func (g *Generator) GenerateCompanionEvent(companion *Companion) *CompanionEvent {
	g.eventID++

	// Pick a trait to base the event on
	if len(companion.Traits) == 0 {
		return nil
	}
	trait := companion.Traits[g.gen.Intn(len(companion.Traits))]

	title := g.generateEventTitle(trait)
	description := g.generateEventDescription(companion.Name, trait)
	dialogue := g.generateDialogue(companion.Name, trait)

	event := NewCompanionEvent(g.eventID, companion.ID, title, description, dialogue, g.genre)
	event.RequiredTrait = trait

	// Set effects based on trait
	g.applyEventEffects(event, trait)

	return event
}

func (g *Generator) generateName() string {
	names := map[engine.GenreID][]string{
		engine.GenreFantasy: {
			"Aldric", "Brynn", "Cedric", "Daria", "Elara", "Finn",
			"Gwendolyn", "Hadrian", "Isolde", "Jasper", "Kira", "Lysander",
			"Mira", "Nolan", "Ophelia", "Percival", "Quinn", "Rowan",
		},
		engine.GenreScifi: {
			"Zara", "Kai", "Nova", "Rex", "Luna", "Axel",
			"Cira", "Dash", "Echo", "Flux", "Gaia", "Hex",
			"Ion", "Jax", "Koda", "Lyra", "Max", "Nix",
		},
		engine.GenreHorror: {
			"Sarah", "Mike", "Jenny", "Tom", "Lisa", "Dan",
			"Amy", "Chris", "Beth", "Jack", "Kate", "Mark",
			"Nina", "Paul", "Rose", "Sam", "Tina", "Will",
		},
		engine.GenreCyberpunk: {
			"Zero", "Ghost", "Razor", "Neon", "Viper", "Chrome",
			"Pixel", "Glitch", "Surge", "Wire", "Byte", "Crash",
			"Flash", "Grid", "Hack", "Ice", "Jolt", "Link",
		},
		engine.GenrePostapoc: {
			"Rust", "Ash", "Duke", "Raven", "Storm", "Bolt",
			"Cinder", "Flint", "Grit", "Hawk", "Iris", "Knox",
			"Moss", "Nash", "Oak", "Pike", "Reed", "Sage",
		},
	}

	nameList := names[g.genre]
	return seed.Choice(g.gen, nameList)
}

func (g *Generator) generateBackstory(name string, role CompanionRole) string {
	backstories := map[engine.GenreID]map[CompanionRole][]string{
		engine.GenreFantasy: {
			RoleGuide: {
				fmt.Sprintf("%s studied the ancient arts in a tower that no longer exists. The secrets learned there guide the way forward.", name),
				fmt.Sprintf("Once an apprentice to a legendary wizard, %s now seeks to complete the unfinished quest of their master.", name),
			},
			RoleScout: {
				fmt.Sprintf("%s grew up in the wildlands, learning to read the forest like others read books.", name),
				fmt.Sprintf("A former hunter for the royal court, %s fell from favor and now walks the untamed paths.", name),
			},
			RoleMedic: {
				fmt.Sprintf("%s served in the temple healers until a crisis of faith sent them on a new journey.", name),
				fmt.Sprintf("Trained in both medicine and magic, %s heals wounds of body and spirit.", name),
			},
			RoleWarrior: {
				fmt.Sprintf("%s was a knight whose order fell to darkness. Honor remains, even without the banner.", name),
				fmt.Sprintf("Battle-scarred and weary, %s fights not for glory but for those who cannot fight.", name),
			},
			RoleTechnician: {
				fmt.Sprintf("%s creates wonders from gears and magic, blending the arcane with the mechanical.", name),
				fmt.Sprintf("Cast out by traditionalist mages, %s proved that innovation holds its own power.", name),
			},
			RoleLeader: {
				fmt.Sprintf("%s led a company of heroes until tragedy scattered them. Leadership endures.", name),
				fmt.Sprintf("Born to nobility, %s chose to earn respect through deeds, not birthright.", name),
			},
		},
		engine.GenreScifi: {
			RoleGuide: {
				fmt.Sprintf("%s is an AI construct that achieved sentience and chose to help organic lifeforms navigate the void.", name),
				fmt.Sprintf("After decades mapping uncharted systems, %s's navigation skills are legendary among spacers.", name),
			},
			RoleScout: {
				fmt.Sprintf("%s was enhanced for surveillance work until corporate warfare made them a liability. Now they watch for themselves.", name),
				fmt.Sprintf("Former military recon, %s sees threats before they materialize.", name),
			},
			RoleMedic: {
				fmt.Sprintf("%s practiced medicine on frontier stations where supplies were scarce and miracles necessary.", name),
				fmt.Sprintf("A medical license revoked for unauthorized procedures that saved lives. %s has no regrets.", name),
			},
			RoleWarrior: {
				fmt.Sprintf("%s served in the Colonial Marines until the Outer War ended. Peace never quite settled in.", name),
				fmt.Sprintf("Cybernetically enhanced for combat, %s struggles to remember life before the upgrades.", name),
			},
			RoleTechnician: {
				fmt.Sprintf("%s can make a hyperdrive from scrap metal and hope. Engine whisperer of the outer rim.", name),
				fmt.Sprintf("Trained in the orbital shipyards, %s left when corporations started cutting corners.", name),
			},
			RoleLeader: {
				fmt.Sprintf("%s commanded a cruiser until a command decision cost everything. Redemption lies ahead.", name),
				fmt.Sprintf("A natural leader who rose through the ranks on talent alone, %s now seeks a cause worth leading.", name),
			},
		},
		engine.GenreHorror: {
			RoleGuide: {
				fmt.Sprintf("%s has seen things that should not exist. That knowledge, however disturbing, might save lives.", name),
				fmt.Sprintf("Years of studying the occult prepared %s for when the impossible became real.", name),
			},
			RoleScout: {
				fmt.Sprintf("%s learned to move silently when the things that hunt in darkness came. Survival became instinct.", name),
				fmt.Sprintf("Once a park ranger, %s now scouts paths through territories no map shows.", name),
			},
			RoleMedic: {
				fmt.Sprintf("%s was working the ER when the outbreak started. Has seen too much to be shocked anymore.", name),
				fmt.Sprintf("Training to save lives didn't include dealing with the undead, but %s adapted.", name),
			},
			RoleWarrior: {
				fmt.Sprintf("%s survived when their group didn't. The guilt fuels a need to protect others.", name),
				fmt.Sprintf("Before the horror, %s was ordinary. Now, violence is a language they speak fluently.", name),
			},
			RoleTechnician: {
				fmt.Sprintf("%s keeps the generators running, the doors sealed, the lights on. Without that, darkness wins.", name),
				fmt.Sprintf("An engineering background meant nothing until %s realized machines could be the difference between life and death.", name),
			},
			RoleLeader: {
				fmt.Sprintf("%s held their community together when order collapsed. That responsibility never lifted.", name),
				fmt.Sprintf("When panic set in, %s stayed calm. Others looked to that calm and found hope.", name),
			},
		},
		engine.GenreCyberpunk: {
			RoleGuide: {
				fmt.Sprintf("%s runs the net like a second home, finding paths through corporate ICE that others miss.", name),
				fmt.Sprintf("Once a corporate decker, %s went rogue when they saw what the corps really do with data.", name),
			},
			RoleScout: {
				fmt.Sprintf("%s's reflex boosters and tactical implants make them see trouble before it arrives.", name),
				fmt.Sprintf("Born in the combat zone, %s learned to sense danger before learning to read.", name),
			},
			RoleMedic: {
				fmt.Sprintf("%s installs chrome and patches meat in a back-alley clinic. No questions, fair prices.", name),
				fmt.Sprintf("Medical school couldn't prepare %s for installing military-grade implants in teenagers.", name),
			},
			RoleWarrior: {
				fmt.Sprintf("%s is more chrome than meat, a weapon walking in human shape. But the humanity underneath remains.", name),
				fmt.Sprintf("Corporate assassin turned independent, %s uses skills meant to oppress to protect instead.", name),
			},
			RoleTechnician: {
				fmt.Sprintf("%s builds custom gear in a workshop that would make corps jealous. Street tech, but quality.", name),
				fmt.Sprintf("%s is an ex-corp engineer who left with some interesting blueprints in their head.", name),
			},
			RoleLeader: {
				fmt.Sprintf("%s knows people, knows connections, knows how to get things done in the shadows.", name),
				fmt.Sprintf("Running a crew taught %s that reputation is everything. Never break your word.", name),
			},
		},
		engine.GenrePostapoc: {
			RoleGuide: {
				fmt.Sprintf("%s remembers the old roads, the safe zones, the places where supplies might still exist.", name),
				fmt.Sprintf("Born after the fall, %s learned the new world's geography the hard way.", name),
			},
			RoleScout: {
				fmt.Sprintf("%s ranges ahead, finding water, shelter, and most importantly, avoiding the threats.", name),
				fmt.Sprintf("The wastelands claimed everyone %s loved. Now they walk alone by choice.", name),
			},
			RoleMedic: {
				fmt.Sprintf("%s treats radiation sickness, infection, and wounds with whatever can be scavenged.", name),
				fmt.Sprintf("Medical knowledge preserved through generations now lives in %s's careful hands.", name),
			},
			RoleWarrior: {
				fmt.Sprintf("%s has killed raiders, mutants, and worse. Each notch on the weapon has a story.", name),
				fmt.Sprintf("The strong prey on the weak. %s ensures their group is never the weak.", name),
			},
			RoleTechnician: {
				fmt.Sprintf("%s makes old tech work again, salvaging function from rust and ruin.", name),
				fmt.Sprintf("In a world where nothing new gets made, %s's ability to repair is priceless.", name),
			},
			RoleLeader: {
				fmt.Sprintf("%s united scattered survivors into something like community. That's rare post-fall.", name),
				fmt.Sprintf("When everything fell apart, %s found purpose in building something new.", name),
			},
		},
	}

	roleBackstories := backstories[g.genre][role]
	return seed.Choice(g.gen, roleBackstories)
}

func (g *Generator) selectTraits(count int) []PersonalityTrait {
	allTraits := AllPersonalityTraits()
	selected := make([]PersonalityTrait, 0, count)

	// Create copy to avoid modifying original
	available := make([]PersonalityTrait, len(allTraits))
	copy(available, allTraits)

	for i := 0; i < count && len(available) > 0; i++ {
		idx := g.gen.Intn(len(available))
		selected = append(selected, available[idx])
		// Remove selected trait
		available = append(available[:idx], available[idx+1:]...)
	}

	return selected
}

func (g *Generator) generateAbilities(role CompanionRole) []*Ability {
	abilities := make([]*Ability, 0, 2)

	abilityData := map[engine.GenreID]map[CompanionRole][][]string{
		engine.GenreFantasy: {
			RoleGuide:      {{"arcane_sight", "Arcane Sight", "Reveals hidden paths and magical dangers"}, {"teleport", "Short Teleport", "Instantly move the party a short distance"}},
			RoleScout:      {{"tracking", "Master Tracking", "Never lose a trail"}, {"ambush", "Set Ambush", "Prepare a devastating first strike"}},
			RoleMedic:      {{"heal", "Divine Healing", "Restore significant health"}, {"cure", "Cure Ailment", "Remove negative conditions"}},
			RoleWarrior:    {{"shield", "Shield Wall", "Block attacks for the party"}, {"rally", "Battle Cry", "Boost party combat strength"}},
			RoleTechnician: {{"repair", "Magical Repair", "Fix equipment instantly"}, {"construct", "Create Golem", "Build a temporary ally"}},
			RoleLeader:     {{"inspire", "Inspiring Presence", "Boost party morale and effectiveness"}, {"command", "Tactical Command", "Coordinate party actions"}},
		},
		engine.GenreScifi: {
			RoleGuide:      {{"nav_assist", "Navigation Assist", "Plot optimal course through hazards"}, {"scan", "Deep Scan", "Reveal hidden objects and threats"}},
			RoleScout:      {{"stealth", "Active Camouflage", "Become nearly invisible"}, {"recon", "Drone Recon", "Scout area remotely"}},
			RoleMedic:      {{"nano_heal", "Nanite Healing", "Deploy healing nanobots"}, {"stim", "Combat Stim", "Temporarily enhance abilities"}},
			RoleWarrior:    {{"suppression", "Suppressive Fire", "Pin down enemies"}, {"breach", "Tactical Breach", "Clear a room efficiently"}},
			RoleTechnician: {{"hack", "System Hack", "Override electronic systems"}, {"upgrade", "Field Upgrade", "Temporarily enhance equipment"}},
			RoleLeader:     {{"coordinate", "Combat Coordination", "Improve party effectiveness"}, {"morale", "Morale Boost", "Restore party will to fight"}},
		},
		engine.GenreHorror: {
			RoleGuide:      {{"ward", "Protective Ward", "Shield against supernatural threats"}, {"banish", "Banishing Ritual", "Drive back evil entities"}},
			RoleScout:      {{"silent_move", "Silent Movement", "Move without alerting threats"}, {"escape", "Escape Route", "Find a way out of danger"}},
			RoleMedic:      {{"first_aid", "Field Surgery", "Treat serious wounds"}, {"antidote", "Antitoxin", "Cure poison and infection"}},
			RoleWarrior:    {{"hold_line", "Hold the Line", "Protect the group from attack"}, {"berserker", "Desperate Strike", "Powerful attack when cornered"}},
			RoleTechnician: {{"fortify", "Quick Fortification", "Rapidly secure a location"}, {"trap", "Set Trap", "Create obstacles for pursuers"}},
			RoleLeader:     {{"calm", "Calm the Panic", "Reduce fear effects"}, {"organize", "Organized Retreat", "Safely withdraw from danger"}},
		},
		engine.GenreCyberpunk: {
			RoleGuide:      {{"netrun", "Deep Net Dive", "Access restricted data systems"}, {"decrypt", "Code Breaking", "Bypass security encryption"}},
			RoleScout:      {{"thermal", "Thermal Vision", "See through walls and darkness"}, {"counter", "Counter Surveillance", "Detect and disable tracking"}},
			RoleMedic:      {{"trauma", "Trauma Team", "Emergency medical intervention"}, {"implant", "Field Implant", "Install or repair cyberware"}},
			RoleWarrior:    {{"overdrive", "Combat Overdrive", "Enter enhanced combat state"}, {"takedown", "Silent Takedown", "Eliminate target quietly"}},
			RoleTechnician: {{"jury_rig", "Jury Rig", "Make broken tech work"}, {"emp", "EMP Burst", "Disable electronic systems"}},
			RoleLeader:     {{"contacts", "Call in Favor", "Get help from contacts"}, {"intimidate", "Street Cred", "Use reputation to resolve conflicts"}},
		},
		engine.GenrePostapoc: {
			RoleGuide:      {{"pathfind", "Wasteland Pathfinding", "Find safe routes through danger zones"}, {"weather", "Weather Sense", "Predict environmental hazards"}},
			RoleScout:      {{"scavenge", "Expert Scavenging", "Find hidden supplies"}, {"track", "Wasteland Tracking", "Follow or avoid trails"}},
			RoleMedic:      {{"rad_treat", "Radiation Treatment", "Reduce radiation damage"}, {"herbal", "Wasteland Medicine", "Create remedies from scavenged materials"}},
			RoleWarrior:    {{"intimidate", "Intimidating Presence", "Deter attackers"}, {"last_stand", "Last Stand", "Fight with desperate strength when wounded"}},
			RoleTechnician: {{"salvage", "Master Salvage", "Extract maximum value from junk"}, {"vehicle", "Vehicle Repair", "Keep transportation running"}},
			RoleLeader:     {{"ration", "Efficient Rationing", "Make supplies last longer"}, {"negotiate", "Wasteland Diplomacy", "Negotiate with hostile groups"}},
		},
	}

	roleAbilities := abilityData[g.genre][role]
	for i, abilityInfo := range roleAbilities {
		minSkill := 5 // First ability unlocks at skill 5
		if i > 0 {
			minSkill = 8 // Second ability unlocks at skill 8
		}

		ability := NewAbility(abilityInfo[0], abilityInfo[1], abilityInfo[2], AbilityActive, minSkill, g.genre)
		ability.Cooldown = 3 + i*2 // Later abilities have longer cooldowns
		abilities = append(abilities, ability)
	}

	return abilities
}

func (g *Generator) applyTraitModifiers(c *Companion) {
	for _, trait := range c.Traits {
		switch trait {
		case TraitBrave:
			c.Morale += 0.1
		case TraitCautious:
			c.Morale -= 0.05
		case TraitOptimistic:
			c.Morale += 0.1
		case TraitPessimistic:
			c.Morale -= 0.1
		case TraitLoyal:
			c.Loyalty += 0.2
		case TraitIndependent:
			c.Loyalty -= 0.1
		case TraitCompassionate:
			c.RelationshipWithPlayer += 0.1
		case TraitPragmatic:
			// No modifier
		}
	}

	// Clamp values
	c.AdjustMorale(0)
	c.AdjustLoyalty(0)
	c.AdjustRelationshipWithPlayer(0)
}

func (g *Generator) generateEventTitle(trait PersonalityTrait) string {
	titles := map[PersonalityTrait][]string{
		TraitBrave:         {"A Bold Stand", "Courage in the Face of Danger", "The Brave Heart"},
		TraitCautious:      {"A Word of Warning", "The Careful Path", "Prudent Counsel"},
		TraitOptimistic:    {"A Ray of Hope", "Looking Forward", "The Brightening"},
		TraitPessimistic:   {"Dark Thoughts", "The Weight of Worry", "Grim Reflections"},
		TraitLoyal:         {"Unwavering Loyalty", "Standing Together", "The Bond"},
		TraitIndependent:   {"Going It Alone", "Self-Reliance", "The Lone Path"},
		TraitCompassionate: {"An Act of Kindness", "Empathy's Call", "The Helping Hand"},
		TraitPragmatic:     {"A Practical Solution", "Hard Choices", "The Calculated Risk"},
	}

	return seed.Choice(g.gen, titles[trait])
}

func (g *Generator) generateEventDescription(name string, trait PersonalityTrait) string {
	descriptions := map[PersonalityTrait][]string{
		TraitBrave: {
			fmt.Sprintf("%s steps forward to face a threat that has others backing away.", name),
			fmt.Sprintf("When danger looms, %s's courage inspires the group.", name),
		},
		TraitCautious: {
			fmt.Sprintf("%s notices something the others missed - a warning sign.", name),
			fmt.Sprintf("Careful observation by %s reveals a hidden danger.", name),
		},
		TraitOptimistic: {
			fmt.Sprintf("%s finds something to be hopeful about, lifting spirits.", name),
			fmt.Sprintf("Even in darkness, %s sees a path to better days.", name),
		},
		TraitPessimistic: {
			fmt.Sprintf("%s's worries prove justified as trouble emerges.", name),
			fmt.Sprintf("The pessimism of %s feels vindicated by events.", name),
		},
		TraitLoyal: {
			fmt.Sprintf("%s's dedication to the group becomes clear in a moment of need.", name),
			fmt.Sprintf("When push comes to shove, %s's loyalty shines through.", name),
		},
		TraitIndependent: {
			fmt.Sprintf("%s handles a situation alone, without consulting the group.", name),
			fmt.Sprintf("Preferring to act alone, %s resolves an issue independently.", name),
		},
		TraitCompassionate: {
			fmt.Sprintf("%s shows kindness to someone in need, despite the risks.", name),
			fmt.Sprintf("Moved by empathy, %s reaches out to help.", name),
		},
		TraitPragmatic: {
			fmt.Sprintf("%s proposes a practical solution to an ongoing problem.", name),
			fmt.Sprintf("Setting aside emotion, %s addresses the situation logically.", name),
		},
	}

	return seed.Choice(g.gen, descriptions[trait])
}

func (g *Generator) generateDialogue(name string, trait PersonalityTrait) string {
	dialogue := map[engine.GenreID]map[PersonalityTrait][]string{
		engine.GenreFantasy: {
			TraitBrave:         {"\"Stand firm! We've faced worse and prevailed!\"", "\"I'll not let fear rule this day.\""},
			TraitCautious:      {"\"Wait - something feels wrong here.\"", "\"We should consider our options carefully.\""},
			TraitOptimistic:    {"\"The dawn always comes, my friends.\"", "\"I believe we'll find a way through this.\""},
			TraitPessimistic:   {"\"I told you this would happen.\"", "\"Prepare for the worst.\""},
			TraitLoyal:         {"\"I swore an oath, and I'll keep it.\"", "\"We're in this together, always.\""},
			TraitIndependent:   {"\"I'll handle this myself.\"", "\"Sometimes one must walk alone.\""},
			TraitCompassionate: {"\"We cannot abandon them to this fate.\"", "\"Mercy is never a weakness.\""},
			TraitPragmatic:     {"\"Let's focus on what we can control.\"", "\"Sentiment won't solve our problems.\""},
		},
		engine.GenreScifi: {
			TraitBrave:         {"\"We've got this. I've run worse odds.\"", "\"Fear is just data - process it and move on.\""},
			TraitCautious:      {"\"Sensors show anomalies. We should investigate.\"", "\"Let's not rush into an unknown situation.\""},
			TraitOptimistic:    {"\"The universe is vast, but so is hope.\"", "\"Statistics favor survival. Trust the numbers.\""},
			TraitPessimistic:   {"\"Probability of success: low.\"", "\"This mission was compromised from the start.\""},
			TraitLoyal:         {"\"My crew, my responsibility.\"", "\"We don't leave anyone behind.\""},
			TraitIndependent:   {"\"I work better alone.\"", "\"My methods, my results.\""},
			TraitCompassionate: {"\"They're sentient beings - we have to help.\"", "\"Protocol doesn't override conscience.\""},
			TraitPragmatic:     {"\"Let's optimize for survival.\"", "\"Resources are finite. We must choose wisely.\""},
		},
		engine.GenreHorror: {
			TraitBrave:         {"\"I'll check it out. Someone has to.\"", "\"We can't hide forever.\""},
			TraitCautious:      {"\"Did you hear that? Stay quiet.\"", "\"We need to find a safer route.\""},
			TraitOptimistic:    {"\"We've made it this far. We'll make it further.\"", "\"There has to be a safe zone somewhere.\""},
			TraitPessimistic:   {"\"We're just delaying the inevitable.\"", "\"How long until it finds us?\""},
			TraitLoyal:         {"\"I'm not leaving you behind.\"", "\"We survive together or not at all.\""},
			TraitIndependent:   {"\"I'm going to scout ahead alone.\"", "\"Sometimes being alone is safer.\""},
			TraitCompassionate: {"\"We can't just leave them here.\"", "\"They're still human, they deserve a chance.\""},
			TraitPragmatic:     {"\"We have to keep moving.\"", "\"Making noise will get us killed.\""},
		},
		engine.GenreCyberpunk: {
			TraitBrave:         {"\"Corp security? I've flatlined worse.\"", "\"Let them come. I'm ready.\""},
			TraitCautious:      {"\"There's too much ICE here. Something's off.\"", "\"This job doesn't smell right.\""},
			TraitOptimistic:    {"\"Big score waiting at the end of this.\"", "\"This could be our ticket out.\""},
			TraitPessimistic:   {"\"Corps always win. We're just buying time.\"", "\"Another day, another betrayal waiting.\""},
			TraitLoyal:         {"\"My crew or nobody.\"", "\"You don't sell out your people.\""},
			TraitIndependent:   {"\"I don't need backup for this.\"", "\"Works better when I solo.\""},
			TraitCompassionate: {"\"These people need help.\"", "\"Can't just ignore them because it's easier.\""},
			TraitPragmatic:     {"\"What's the profit margin on heroics?\"", "\"Let's focus on the objective.\""},
		},
		engine.GenrePostapoc: {
			TraitBrave:         {"\"Raiders? I've handled worse.\"", "\"We stand our ground.\""},
			TraitCautious:      {"\"That smoke on the horizon - we go around.\"", "\"Better to lose time than lives.\""},
			TraitOptimistic:    {"\"Heard there's clean water further west.\"", "\"We'll find a place to call home.\""},
			TraitPessimistic:   {"\"How long until the supplies run out?\"", "\"Nothing good lasts in this world.\""},
			TraitLoyal:         {"\"This group is all I have.\"", "\"We're family now, whether we like it or not.\""},
			TraitIndependent:   {"\"I move faster alone.\"", "\"Don't wait for me.\""},
			TraitCompassionate: {"\"We have enough to share.\"", "\"If we lose our humanity, what's the point?\""},
			TraitPragmatic:     {"\"One more mouth means less for everyone.\"", "\"Hard choices keep us alive.\""},
		},
	}

	return seed.Choice(g.gen, dialogue[g.genre][trait])
}

func (g *Generator) applyEventEffects(event *CompanionEvent, trait PersonalityTrait) {
	switch trait {
	case TraitBrave:
		event.MoraleChange = 0.1
		event.RelationshipChange = 0.05
	case TraitCautious:
		event.MoraleChange = 0.05
		event.SkillGain = 1
	case TraitOptimistic:
		event.MoraleChange = 0.15
	case TraitPessimistic:
		event.MoraleChange = -0.05
		event.LoyaltyChange = 0.05
	case TraitLoyal:
		event.LoyaltyChange = 0.15
		event.RelationshipChange = 0.1
	case TraitIndependent:
		event.SkillGain = 2
		event.LoyaltyChange = -0.05
	case TraitCompassionate:
		event.RelationshipChange = 0.15
		event.MoraleChange = 0.05
	case TraitPragmatic:
		event.SkillGain = 1
		event.MoraleChange = 0.05
	}
}

// GenerateParty creates a full party of companions
func (g *Generator) GenerateParty(size int) *CompanionManager {
	manager := NewCompanionManager(g.genre, size)

	// Generate diverse roles
	roles := AllCompanionRoles()
	for i := 0; i < size; i++ {
		role := roles[i%len(roles)]
		companion := g.GenerateCompanion(role)
		companion.JoinedAt = 0
		manager.AddCompanion(companion)

		// Generate events for each companion
		event := g.GenerateCompanionEvent(companion)
		if event != nil {
			manager.AddEvent(event)
		}
	}

	return manager
}
