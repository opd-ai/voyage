package crew

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// CrewMember represents an individual party member.
type CrewMember struct {
	ID          int
	Name        string
	Health      float64
	MaxHealth   float64
	Trait       Trait
	Skill       Skill
	IsAlive     bool
	DaysWithParty int
}

// Trait represents a character personality trait.
type Trait int

const (
	TraitBrave Trait = iota
	TraitCautious
	TraitOptimistic
	TraitPessimistic
	TraitGreedy
	TraitGenerous
	TraitStoic
	TraitEmotional
)

// Skill represents a character's primary skill.
type Skill int

const (
	SkillNone Skill = iota
	SkillMedic
	SkillMechanic
	SkillScout
	SkillTrader
	SkillWarrior
	SkillLeader
)

// NewCrewMember creates a new crew member with the given attributes.
func NewCrewMember(id int, name string, trait Trait, skill Skill) *CrewMember {
	return &CrewMember{
		ID:          id,
		Name:        name,
		Health:      100,
		MaxHealth:   100,
		Trait:       trait,
		Skill:       skill,
		IsAlive:     true,
		DaysWithParty: 0,
	}
}

// TakeDamage reduces health by the given amount.
// Returns true if the crew member died.
func (c *CrewMember) TakeDamage(amount float64) bool {
	c.Health -= amount
	if c.Health <= 0 {
		c.Health = 0
		c.IsAlive = false
		return true
	}
	return false
}

// Heal increases health by the given amount.
func (c *CrewMember) Heal(amount float64) {
	c.Health += amount
	if c.Health > c.MaxHealth {
		c.Health = c.MaxHealth
	}
}

// HealthRatio returns health as a ratio [0, 1].
func (c *CrewMember) HealthRatio() float64 {
	if c.MaxHealth <= 0 {
		return 0
	}
	return c.Health / c.MaxHealth
}

// TraitName returns the name of the trait for the given genre.
func TraitName(t Trait, genre engine.GenreID) string {
	names, ok := traitNames[genre]
	if !ok {
		names = traitNames[engine.GenreFantasy]
	}
	return names[t]
}

// SkillName returns the name of the skill for the given genre.
func SkillName(s Skill, genre engine.GenreID) string {
	names, ok := skillNames[genre]
	if !ok {
		names = skillNames[engine.GenreFantasy]
	}
	return names[s]
}

var traitNames = map[engine.GenreID]map[Trait]string{
	engine.GenreFantasy: {
		TraitBrave:       "Brave",
		TraitCautious:    "Cautious",
		TraitOptimistic:  "Hopeful",
		TraitPessimistic: "Gloomy",
		TraitGreedy:      "Greedy",
		TraitGenerous:    "Generous",
		TraitStoic:       "Stoic",
		TraitEmotional:   "Passionate",
	},
	engine.GenreScifi: {
		TraitBrave:       "Fearless",
		TraitCautious:    "Calculated",
		TraitOptimistic:  "Optimistic",
		TraitPessimistic: "Cynical",
		TraitGreedy:      "Acquisitive",
		TraitGenerous:    "Altruistic",
		TraitStoic:       "Logical",
		TraitEmotional:   "Empathic",
	},
	engine.GenreHorror: {
		TraitBrave:       "Fearless",
		TraitCautious:    "Paranoid",
		TraitOptimistic:  "Delusional",
		TraitPessimistic: "Fatalist",
		TraitGreedy:      "Hoarder",
		TraitGenerous:    "Selfless",
		TraitStoic:       "Hardened",
		TraitEmotional:   "Unstable",
	},
	engine.GenreCyberpunk: {
		TraitBrave:       "Reckless",
		TraitCautious:    "Street-smart",
		TraitOptimistic:  "Dreamer",
		TraitPessimistic: "Nihilist",
		TraitGreedy:      "Corporate",
		TraitGenerous:    "Anarchist",
		TraitStoic:       "Chrome-cold",
		TraitEmotional:   "Wire-hot",
	},
	engine.GenrePostapoc: {
		TraitBrave:       "Survivor",
		TraitCautious:    "Wary",
		TraitOptimistic:  "Believer",
		TraitPessimistic: "Doom-sayer",
		TraitGreedy:      "Scavenger",
		TraitGenerous:    "Sharer",
		TraitStoic:       "Weathered",
		TraitEmotional:   "Broken",
	},
}

var skillNames = map[engine.GenreID]map[Skill]string{
	engine.GenreFantasy: {
		SkillNone:     "Peasant",
		SkillMedic:    "Healer",
		SkillMechanic: "Craftsman",
		SkillScout:    "Ranger",
		SkillTrader:   "Merchant",
		SkillWarrior:  "Knight",
		SkillLeader:   "Noble",
	},
	engine.GenreScifi: {
		SkillNone:     "Civilian",
		SkillMedic:    "Medic",
		SkillMechanic: "Engineer",
		SkillScout:    "Navigator",
		SkillTrader:   "Merchant",
		SkillWarrior:  "Marine",
		SkillLeader:   "Captain",
	},
	engine.GenreHorror: {
		SkillNone:     "Survivor",
		SkillMedic:    "Nurse",
		SkillMechanic: "Mechanic",
		SkillScout:    "Scout",
		SkillTrader:   "Scavenger",
		SkillWarrior:  "Fighter",
		SkillLeader:   "Leader",
	},
	engine.GenreCyberpunk: {
		SkillNone:     "Street Rat",
		SkillMedic:    "Ripperdoc",
		SkillMechanic: "Tech",
		SkillScout:    "Netrunner",
		SkillTrader:   "Fixer",
		SkillWarrior:  "Solo",
		SkillLeader:   "Face",
	},
	engine.GenrePostapoc: {
		SkillNone:     "Wanderer",
		SkillMedic:    "Medic",
		SkillMechanic: "Wrench",
		SkillScout:    "Tracker",
		SkillTrader:   "Trader",
		SkillWarrior:  "Raider",
		SkillLeader:   "Chief",
	},
}

// AllTraits returns all trait types.
func AllTraits() []Trait {
	return []Trait{
		TraitBrave,
		TraitCautious,
		TraitOptimistic,
		TraitPessimistic,
		TraitGreedy,
		TraitGenerous,
		TraitStoic,
		TraitEmotional,
	}
}

// AllSkills returns all skill types.
func AllSkills() []Skill {
	return []Skill{
		SkillNone,
		SkillMedic,
		SkillMechanic,
		SkillScout,
		SkillTrader,
		SkillWarrior,
		SkillLeader,
	}
}

// Generator creates procedural crew members.
type Generator struct {
	gen   *seed.Generator
	genre engine.GenreID
	nextID int
}

// NewGenerator creates a new crew generator.
func NewGenerator(masterSeed int64, genre engine.GenreID) *Generator {
	return &Generator{
		gen:   seed.NewGenerator(masterSeed, "crew"),
		genre: genre,
		nextID: 1,
	}
}

// SetGenre changes the generator's genre.
func (g *Generator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// Generate creates a new crew member.
func (g *Generator) Generate() *CrewMember {
	name := g.generateName()
	trait := seed.Choice(g.gen, AllTraits())
	
	// Weight skills - most are None, specialists are rarer
	skillWeights := []float64{0.4, 0.12, 0.12, 0.12, 0.08, 0.08, 0.08}
	skill := seed.WeightedChoice(g.gen, AllSkills(), skillWeights)
	
	member := NewCrewMember(g.nextID, name, trait, skill)
	g.nextID++
	
	// Vary starting health slightly
	member.Health = 80 + float64(g.gen.Intn(21))
	
	return member
}

// generateName creates a procedural name based on genre.
func (g *Generator) generateName() string {
	firstNames := firstNamesByGenre[g.genre]
	lastNames := lastNamesByGenre[g.genre]
	
	first := seed.Choice(g.gen, firstNames)
	last := seed.Choice(g.gen, lastNames)
	
	return first + " " + last
}

var firstNamesByGenre = map[engine.GenreID][]string{
	engine.GenreFantasy:   {"Aldric", "Brynn", "Cedric", "Dara", "Elara", "Finn", "Gwen", "Hector", "Ivy", "Jareth", "Kira", "Liam"},
	engine.GenreScifi:     {"Astra", "Beck", "Cade", "Dex", "Echo", "Flux", "Gaia", "Hex", "Ion", "Jax", "Kira", "Luna"},
	engine.GenreHorror:    {"Alex", "Blake", "Casey", "Dana", "Eli", "Frank", "Grace", "Hunter", "Isaac", "Jamie", "Kelly", "Lee"},
	engine.GenreCyberpunk: {"Blade", "Chrome", "Dash", "Edge", "Flash", "Ghost", "Hack", "Ice", "Jack", "Knife", "Link", "Mox"},
	engine.GenrePostapoc:  {"Ash", "Blaze", "Crow", "Dust", "Echo", "Flint", "Grit", "Haze", "Iron", "Junk", "Knox", "Rust"},
}

var lastNamesByGenre = map[engine.GenreID][]string{
	engine.GenreFantasy:   {"Blackwood", "Brightblade", "Dawnseeker", "Evershade", "Frostborn", "Goldleaf", "Ironhelm", "Moonshadow"},
	engine.GenreScifi:     {"Nova", "Stellar", "Quantum", "Vector", "Nexus", "Cosmos", "Orbital", "Prime"},
	engine.GenreHorror:    {"Walker", "Hunter", "Cross", "Graves", "Stone", "Black", "Cold", "Sharp"},
	engine.GenreCyberpunk: {"Zero", "One", "Two", "Three", "Four", "Five", "Six", "Seven"},
	engine.GenrePostapoc:  {"Wasteland", "Sandstorm", "Ironsides", "Scorched", "Razorback", "Deadzone", "Thunderhead", "Firestorm"},
}
