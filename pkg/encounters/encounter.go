package encounters

import (
	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// Encounter represents a tactical encounter in progress.
type Encounter struct {
	ID          int
	Type        EncounterType
	Title       string
	Description string
	Genre       engine.GenreID
	Difficulty  float64 // 0.0-1.0, affects outcome chances

	// Resolution state
	State        EncounterState
	CurrentPhase int
	MaxPhases    int
	IsPaused     bool

	// Crew assignments
	Assignments map[EncounterRole]int // Role -> CrewMember ID

	// Required roles (some encounters need specific roles filled)
	RequiredRoles []EncounterRole
	OptimalRoles  []EncounterRole // Roles that provide bonuses

	// Accumulated stats during resolution
	TotalDamage   float64
	TotalProgress float64
	TurnsElapsed  int
}

// EncounterState tracks the current state of encounter resolution.
type EncounterState int

const (
	// StatePending means the encounter hasn't started yet.
	StatePending EncounterState = iota
	// StateAssignment means crew are being assigned to roles.
	StateAssignment
	// StateResolution means the encounter is being resolved.
	StateResolution
	// StateComplete means the encounter has finished.
	StateComplete
)

// NewEncounter creates a new encounter.
func NewEncounter(id int, encType EncounterType, genre engine.GenreID) *Encounter {
	return &Encounter{
		ID:          id,
		Type:        encType,
		Genre:       genre,
		Difficulty:  0.5,
		State:       StatePending,
		MaxPhases:   3,
		Assignments: make(map[EncounterRole]int),
	}
}

// SetGenre updates the encounter's genre theme.
func (e *Encounter) SetGenre(genre engine.GenreID) {
	e.Genre = genre
}

// AssignCrew assigns a crew member to a role.
func (e *Encounter) AssignCrew(role EncounterRole, memberID int) {
	e.Assignments[role] = memberID
}

// UnassignCrew removes a crew member from a role.
func (e *Encounter) UnassignCrew(role EncounterRole) {
	delete(e.Assignments, role)
}

// GetAssignment returns the crew member ID assigned to a role.
func (e *Encounter) GetAssignment(role EncounterRole) (int, bool) {
	id, ok := e.Assignments[role]
	return id, ok
}

// IsRoleFilled checks if a role has a crew member assigned.
func (e *Encounter) IsRoleFilled(role EncounterRole) bool {
	_, ok := e.Assignments[role]
	return ok
}

// RequiredRolesFilled checks if all required roles have crew assigned.
func (e *Encounter) RequiredRolesFilled() bool {
	for _, role := range e.RequiredRoles {
		if !e.IsRoleFilled(role) {
			return false
		}
	}
	return true
}

// Start begins encounter resolution.
func (e *Encounter) Start() bool {
	if e.State != StatePending && e.State != StateAssignment {
		return false
	}
	if !e.RequiredRolesFilled() {
		return false
	}
	e.State = StateResolution
	e.CurrentPhase = 0
	return true
}

// Pause pauses encounter resolution.
func (e *Encounter) Pause() {
	if e.State == StateResolution {
		e.IsPaused = true
	}
}

// Resume resumes encounter resolution.
func (e *Encounter) Resume() {
	if e.State == StateResolution {
		e.IsPaused = false
	}
}

// Generator creates procedural encounters.
type Generator struct {
	gen    *seed.Generator
	genre  engine.GenreID
	nextID int
}

// NewGenerator creates a new encounter generator.
func NewGenerator(masterSeed int64, genre engine.GenreID) *Generator {
	return &Generator{
		gen:    seed.NewGenerator(masterSeed, "encounters"),
		genre:  genre,
		nextID: 1,
	}
}

// SetGenre updates the generator's genre.
func (g *Generator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// Generate creates a new encounter of the specified type.
func (g *Generator) Generate(encType EncounterType) *Encounter {
	enc := NewEncounter(g.nextID, encType, g.genre)
	g.nextID++

	g.setupEncounterByType(enc)

	return enc
}

// GenerateRandom creates a random encounter type.
func (g *Generator) GenerateRandom() *Encounter {
	encType := seed.Choice(g.gen, AllEncounterTypes())
	return g.Generate(encType)
}

func (g *Generator) setupEncounterByType(enc *Encounter) {
	templates := encounterTemplates[g.genre]
	if templates == nil {
		templates = encounterTemplates[engine.GenreFantasy]
	}

	typeTemplates := templates[enc.Type]
	if len(typeTemplates) == 0 {
		enc.Title = "Unknown Encounter"
		enc.Description = "Something unexpected happens."
		return
	}

	template := seed.Choice(g.gen, typeTemplates)
	enc.Title = template.Title
	enc.Description = template.Description
	enc.RequiredRoles = template.RequiredRoles
	enc.OptimalRoles = template.OptimalRoles
	enc.MaxPhases = template.Phases
	enc.Difficulty = 0.3 + g.gen.Float64()*0.5 // 0.3-0.8 range
}

type encounterTemplate struct {
	Title         string
	Description   string
	RequiredRoles []EncounterRole
	OptimalRoles  []EncounterRole
	Phases        int
}

var encounterTemplates = map[engine.GenreID]map[EncounterType][]encounterTemplate{
	engine.GenreFantasy: {
		TypeAmbush: {
			{
				Title:         "Bandit Ambush",
				Description:   "A group of bandits emerges from the treeline, weapons drawn.",
				RequiredRoles: []EncounterRole{RoleFighter},
				OptimalRoles:  []EncounterRole{RoleFighter, RoleMedic},
				Phases:        3,
			},
			{
				Title:         "Monster Attack",
				Description:   "A fearsome beast blocks your path, snarling with hunger.",
				RequiredRoles: []EncounterRole{RoleFighter},
				OptimalRoles:  []EncounterRole{RoleFighter, RoleNegotiator},
				Phases:        4,
			},
		},
		TypeNegotiation: {
			{
				Title:         "Toll Collectors",
				Description:   "Armed guards demand payment to pass through their territory.",
				RequiredRoles: []EncounterRole{RoleNegotiator},
				OptimalRoles:  []EncounterRole{RoleNegotiator, RoleFighter},
				Phases:        2,
			},
			{
				Title:         "Merchant Caravan",
				Description:   "A wealthy merchant offers to trade goods.",
				RequiredRoles: []EncounterRole{},
				OptimalRoles:  []EncounterRole{RoleNegotiator},
				Phases:        2,
			},
		},
		TypeRace: {
			{
				Title:         "Pursuit",
				Description:   "Enemies give chase! You must outrun them.",
				RequiredRoles: []EncounterRole{RoleEngineer},
				OptimalRoles:  []EncounterRole{RoleEngineer, RoleFighter},
				Phases:        3,
			},
		},
		TypeCrisis: {
			{
				Title:         "Bridge Collapse",
				Description:   "The bridge ahead is crumbling. Time is short.",
				RequiredRoles: []EncounterRole{RoleEngineer},
				OptimalRoles:  []EncounterRole{RoleEngineer, RoleMedic},
				Phases:        2,
			},
			{
				Title:         "Plague Outbreak",
				Description:   "A sickness spreads through the caravan.",
				RequiredRoles: []EncounterRole{RoleMedic},
				OptimalRoles:  []EncounterRole{RoleMedic},
				Phases:        3,
			},
		},
		TypePuzzle: {
			{
				Title:         "Ancient Lock",
				Description:   "A sealed door blocks access to valuable supplies.",
				RequiredRoles: []EncounterRole{},
				OptimalRoles:  []EncounterRole{RoleEngineer},
				Phases:        2,
			},
		},
	},
	engine.GenreScifi: {
		TypeAmbush: {
			{
				Title:         "Pirate Attack",
				Description:   "A hostile vessel emerges from an asteroid's shadow.",
				RequiredRoles: []EncounterRole{RoleFighter},
				OptimalRoles:  []EncounterRole{RoleFighter, RoleEngineer},
				Phases:        4,
			},
			{
				Title:         "Alien Hostiles",
				Description:   "Unknown life forms breach the hull.",
				RequiredRoles: []EncounterRole{RoleFighter},
				OptimalRoles:  []EncounterRole{RoleFighter, RoleMedic},
				Phases:        3,
			},
		},
		TypeNegotiation: {
			{
				Title:         "Station Docking Request",
				Description:   "The station demands credentials before granting access.",
				RequiredRoles: []EncounterRole{RoleNegotiator},
				OptimalRoles:  []EncounterRole{RoleNegotiator},
				Phases:        2,
			},
		},
		TypeRace: {
			{
				Title:         "Pursuit Vector",
				Description:   "Hostile ships are closing fast. Escape is the only option.",
				RequiredRoles: []EncounterRole{RoleEngineer},
				OptimalRoles:  []EncounterRole{RoleEngineer, RoleFighter},
				Phases:        3,
			},
		},
		TypeCrisis: {
			{
				Title:         "Hull Breach",
				Description:   "Emergency! The hull is compromised.",
				RequiredRoles: []EncounterRole{RoleEngineer},
				OptimalRoles:  []EncounterRole{RoleEngineer, RoleMedic},
				Phases:        2,
			},
		},
		TypePuzzle: {
			{
				Title:         "Encrypted Signal",
				Description:   "A signal contains valuable data, but it's encrypted.",
				RequiredRoles: []EncounterRole{},
				OptimalRoles:  []EncounterRole{RoleEngineer},
				Phases:        2,
			},
		},
	},
	engine.GenreHorror: {
		TypeAmbush: {
			{
				Title:         "Horde Attack",
				Description:   "They emerge from everywhere, moaning and hungry.",
				RequiredRoles: []EncounterRole{RoleFighter},
				OptimalRoles:  []EncounterRole{RoleFighter, RoleMedic},
				Phases:        4,
			},
		},
		TypeNegotiation: {
			{
				Title:         "Survivor Standoff",
				Description:   "Other survivors aim weapons at you, terrified.",
				RequiredRoles: []EncounterRole{RoleNegotiator},
				OptimalRoles:  []EncounterRole{RoleNegotiator, RoleMedic},
				Phases:        2,
			},
		},
		TypeRace: {
			{
				Title:         "Run!",
				Description:   "A massive horde approaches. There's no fighting this.",
				RequiredRoles: []EncounterRole{},
				OptimalRoles:  []EncounterRole{RoleEngineer},
				Phases:        3,
			},
		},
		TypeCrisis: {
			{
				Title:         "Infection",
				Description:   "Someone was bitten. Decisions must be made.",
				RequiredRoles: []EncounterRole{RoleMedic},
				OptimalRoles:  []EncounterRole{RoleMedic, RoleNegotiator},
				Phases:        2,
			},
		},
		TypePuzzle: {
			{
				Title:         "Safe Room",
				Description:   "The door is locked. Supplies are visible inside.",
				RequiredRoles: []EncounterRole{},
				OptimalRoles:  []EncounterRole{RoleEngineer},
				Phases:        2,
			},
		},
	},
	engine.GenreCyberpunk: {
		TypeAmbush: {
			{
				Title:         "Gang Ambush",
				Description:   "A local gang has decided you're worth robbing.",
				RequiredRoles: []EncounterRole{RoleFighter},
				OptimalRoles:  []EncounterRole{RoleFighter, RoleNegotiator},
				Phases:        3,
			},
		},
		TypeNegotiation: {
			{
				Title:         "Fixer Meeting",
				Description:   "Your contact wants to renegotiate terms.",
				RequiredRoles: []EncounterRole{RoleNegotiator},
				OptimalRoles:  []EncounterRole{RoleNegotiator, RoleFighter},
				Phases:        2,
			},
		},
		TypeRace: {
			{
				Title:         "Corporate Pursuit",
				Description:   "Corp security is on your tail. Time to floor it.",
				RequiredRoles: []EncounterRole{RoleEngineer},
				OptimalRoles:  []EncounterRole{RoleEngineer, RoleFighter},
				Phases:        3,
			},
		},
		TypeCrisis: {
			{
				Title:         "Cyberware Malfunction",
				Description:   "Someone's chrome is glitching dangerously.",
				RequiredRoles: []EncounterRole{RoleEngineer},
				OptimalRoles:  []EncounterRole{RoleEngineer, RoleMedic},
				Phases:        2,
			},
		},
		TypePuzzle: {
			{
				Title:         "ICE Breach",
				Description:   "Black ICE protects valuable data. Crack it.",
				RequiredRoles: []EncounterRole{RoleEngineer},
				OptimalRoles:  []EncounterRole{RoleEngineer},
				Phases:        3,
			},
		},
	},
	engine.GenrePostapoc: {
		TypeAmbush: {
			{
				Title:         "Raider Attack",
				Description:   "Raiders emerge from the ruins, hungry for supplies.",
				RequiredRoles: []EncounterRole{RoleFighter},
				OptimalRoles:  []EncounterRole{RoleFighter, RoleMedic},
				Phases:        3,
			},
		},
		TypeNegotiation: {
			{
				Title:         "Settlement Gate",
				Description:   "The guards demand proof you're not sick.",
				RequiredRoles: []EncounterRole{RoleNegotiator},
				OptimalRoles:  []EncounterRole{RoleNegotiator, RoleMedic},
				Phases:        2,
			},
		},
		TypeRace: {
			{
				Title:         "Rad Storm",
				Description:   "A radiation storm approaches. Find shelter fast.",
				RequiredRoles: []EncounterRole{},
				OptimalRoles:  []EncounterRole{RoleEngineer},
				Phases:        2,
			},
		},
		TypeCrisis: {
			{
				Title:         "Radiation Sickness",
				Description:   "Someone absorbed too much rads.",
				RequiredRoles: []EncounterRole{RoleMedic},
				OptimalRoles:  []EncounterRole{RoleMedic},
				Phases:        3,
			},
		},
		TypePuzzle: {
			{
				Title:         "Pre-War Cache",
				Description:   "A sealed bunker door. Something valuable inside.",
				RequiredRoles: []EncounterRole{},
				OptimalRoles:  []EncounterRole{RoleEngineer},
				Phases:        2,
			},
		},
	},
}

// CalculateRoleEffectiveness calculates how effective a crew member is in a role.
func CalculateRoleEffectiveness(member *crew.CrewMember, role EncounterRole) float64 {
	if member == nil || !member.IsAlive {
		return 0
	}

	base := 0.5 // Base effectiveness

	// Skill bonus
	skillBonus := getSkillBonusForRole(member.Skill, role)
	base += skillBonus

	// Trait bonus
	traitBonus := getTraitBonusForRole(member.Trait, role)
	base += traitBonus

	// Health penalty
	healthMod := member.HealthRatio()
	base *= healthMod

	// Skill level bonus
	base *= member.SkillEffectiveness()

	return clamp(base, 0, 2.0)
}

func getSkillBonusForRole(skill crew.Skill, role EncounterRole) float64 {
	bonuses := map[EncounterRole]map[crew.Skill]float64{
		RoleFighter: {
			crew.SkillWarrior: 0.4,
			crew.SkillLeader:  0.2,
		},
		RoleMedic: {
			crew.SkillMedic:  0.4,
			crew.SkillLeader: 0.1,
		},
		RoleEngineer: {
			crew.SkillMechanic: 0.4,
			crew.SkillScout:    0.2,
		},
		RoleNegotiator: {
			crew.SkillTrader: 0.4,
			crew.SkillLeader: 0.3,
		},
	}

	if roleBonuses, ok := bonuses[role]; ok {
		if bonus, ok := roleBonuses[skill]; ok {
			return bonus
		}
	}
	return 0
}

func getTraitBonusForRole(trait crew.Trait, role EncounterRole) float64 {
	bonuses := map[EncounterRole]map[crew.Trait]float64{
		RoleFighter: {
			crew.TraitBrave: 0.2,
			crew.TraitStoic: 0.1,
		},
		RoleMedic: {
			crew.TraitOptimistic: 0.1,
			crew.TraitGenerous:   0.1,
		},
		RoleEngineer: {
			crew.TraitCautious:  0.1,
			crew.TraitScavenger: 0.2,
		},
		RoleNegotiator: {
			crew.TraitGenerous:   0.2,
			crew.TraitOptimistic: 0.1,
		},
	}

	if roleBonuses, ok := bonuses[role]; ok {
		if bonus, ok := roleBonuses[trait]; ok {
			return bonus
		}
	}
	return 0
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
