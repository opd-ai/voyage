package companions

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// CompanionRole represents the functional role of a companion
type CompanionRole int

const (
	RoleGuide CompanionRole = iota
	RoleScout
	RoleMedic
	RoleWarrior
	RoleTechnician
	RoleLeader
)

// PersonalityTrait represents a character personality trait
type PersonalityTrait int

const (
	TraitBrave PersonalityTrait = iota
	TraitCautious
	TraitOptimistic
	TraitPessimistic
	TraitLoyal
	TraitIndependent
	TraitCompassionate
	TraitPragmatic
)

// AbilityType represents the category of special ability
type AbilityType int

const (
	AbilityPassive  AbilityType = iota // Always active
	AbilityActive                      // Must be triggered
	AbilityReactive                    // Triggers on condition
)

// Ability represents a special skill a companion can use
type Ability struct {
	ID          string
	Name        string
	Description string
	Type        AbilityType
	Genre       engine.GenreID

	// Requirements
	MinSkillLevel int
	Cooldown      int // Turns between uses (0 = no cooldown)

	// Effects
	EffectStrength float64 // Multiplier for ability effect
	TargetsSelf    bool
	TargetsAlly    bool
	TargetsEnemy   bool

	// State
	CurrentCooldown int
	Unlocked        bool
}

// NewAbility creates a new ability
func NewAbility(id, name, description string, abilityType AbilityType, minSkill int, genre engine.GenreID) *Ability {
	return &Ability{
		ID:             id,
		Name:           name,
		Description:    description,
		Type:           abilityType,
		MinSkillLevel:  minSkill,
		EffectStrength: 1.0,
		Genre:          genre,
	}
}

// SetGenre updates the ability's genre
func (a *Ability) SetGenre(g engine.GenreID) {
	a.Genre = g
}

// IsReady returns true if the ability can be used
func (a *Ability) IsReady() bool {
	return a.Unlocked && a.CurrentCooldown == 0
}

// Use activates the ability and starts cooldown
func (a *Ability) Use() bool {
	if !a.IsReady() {
		return false
	}
	a.CurrentCooldown = a.Cooldown
	return true
}

// Tick reduces cooldown by one turn
func (a *Ability) Tick() {
	if a.CurrentCooldown > 0 {
		a.CurrentCooldown--
	}
}

// Companion represents a named crew member with special abilities
type Companion struct {
	ID        int
	Name      string
	Title     string // Genre-appropriate title (Wizard, AI, Handler, etc.)
	Backstory string
	Genre     engine.GenreID

	// Role and skills
	Role          CompanionRole
	SkillLevel    int
	Experience    int
	MaxSkillLevel int

	// Personality
	Traits  []PersonalityTrait
	Morale  float64 // 0.0 to 1.0
	Loyalty float64 // 0.0 to 1.0

	// Abilities
	Abilities []*Ability

	// State
	Health   float64
	IsActive bool
	JoinedAt int // Day joined the party

	// Relationships
	RelationshipWithPlayer float64         // -1.0 to 1.0
	RelationshipWithCrew   map[int]float64 // Other companion ID -> relationship
}

// NewCompanion creates a new companion
func NewCompanion(id int, name, title string, role CompanionRole, genre engine.GenreID) *Companion {
	return &Companion{
		ID:                     id,
		Name:                   name,
		Title:                  title,
		Role:                   role,
		Genre:                  genre,
		SkillLevel:             1,
		MaxSkillLevel:          10,
		Traits:                 make([]PersonalityTrait, 0),
		Morale:                 0.7,
		Loyalty:                0.5,
		Abilities:              make([]*Ability, 0),
		Health:                 1.0,
		IsActive:               true,
		RelationshipWithPlayer: 0.0,
		RelationshipWithCrew:   make(map[int]float64),
	}
}

// SetGenre updates the companion's genre and all abilities
func (c *Companion) SetGenre(g engine.GenreID) {
	c.Genre = g
	for _, ability := range c.Abilities {
		ability.SetGenre(g)
	}
}

// AddTrait adds a personality trait
func (c *Companion) AddTrait(trait PersonalityTrait) {
	c.Traits = append(c.Traits, trait)
}

// HasTrait checks if companion has a specific trait
func (c *Companion) HasTrait(trait PersonalityTrait) bool {
	for _, t := range c.Traits {
		if t == trait {
			return true
		}
	}
	return false
}

// AddAbility adds a special ability
func (c *Companion) AddAbility(ability *Ability) {
	c.Abilities = append(c.Abilities, ability)
	// Check if ability should be unlocked
	if c.SkillLevel >= ability.MinSkillLevel {
		ability.Unlocked = true
	}
}

// GetAbility returns the companion's primary ability (first unlocked one)
func (c *Companion) GetAbility() *Ability {
	for _, ability := range c.Abilities {
		if ability.Unlocked {
			return ability
		}
	}
	return nil
}

// CanUseAbility returns true if companion has an ability ready to use
func (c *Companion) CanUseAbility() bool {
	for _, ability := range c.Abilities {
		if ability.IsReady() {
			return true
		}
	}
	return false
}

// GainExperience adds experience and potentially levels up skill
func (c *Companion) GainExperience(amount int) bool {
	c.Experience += amount
	leveledUp := false

	// Level up at 100 XP per level
	for c.Experience >= 100 && c.SkillLevel < c.MaxSkillLevel {
		c.Experience -= 100
		c.SkillLevel++
		leveledUp = true

		// Check for ability unlocks
		for _, ability := range c.Abilities {
			if !ability.Unlocked && c.SkillLevel >= ability.MinSkillLevel {
				ability.Unlocked = true
			}
		}
	}

	return leveledUp
}

// AdjustMorale changes morale within bounds
func (c *Companion) AdjustMorale(delta float64) {
	c.Morale += delta
	if c.Morale < 0 {
		c.Morale = 0
	}
	if c.Morale > 1 {
		c.Morale = 1
	}
}

// AdjustLoyalty changes loyalty within bounds
func (c *Companion) AdjustLoyalty(delta float64) {
	c.Loyalty += delta
	if c.Loyalty < 0 {
		c.Loyalty = 0
	}
	if c.Loyalty > 1 {
		c.Loyalty = 1
	}
}

// AdjustRelationshipWithPlayer changes relationship with player
func (c *Companion) AdjustRelationshipWithPlayer(delta float64) {
	c.RelationshipWithPlayer += delta
	if c.RelationshipWithPlayer < -1 {
		c.RelationshipWithPlayer = -1
	}
	if c.RelationshipWithPlayer > 1 {
		c.RelationshipWithPlayer = 1
	}
}

// Tick processes one turn for the companion (cooldowns, etc.)
func (c *Companion) Tick() {
	for _, ability := range c.Abilities {
		ability.Tick()
	}
}

// GetUnlockedAbilities returns all unlocked abilities
func (c *Companion) GetUnlockedAbilities() []*Ability {
	unlocked := make([]*Ability, 0)
	for _, ability := range c.Abilities {
		if ability.Unlocked {
			unlocked = append(unlocked, ability)
		}
	}
	return unlocked
}

// CompanionEvent represents a personality-driven special event
type CompanionEvent struct {
	ID          int
	CompanionID int
	Title       string
	Description string
	Dialogue    string
	Genre       engine.GenreID

	// Trigger conditions
	RequiredTrait   PersonalityTrait
	MinMorale       float64
	MaxMorale       float64
	MinLoyalty      float64
	MinRelationship float64

	// Effects
	MoraleChange       float64
	LoyaltyChange      float64
	RelationshipChange float64
	SkillGain          int

	// State
	Triggered bool
}

// NewCompanionEvent creates a new companion event
func NewCompanionEvent(id, companionID int, title, description, dialogue string, genre engine.GenreID) *CompanionEvent {
	return &CompanionEvent{
		ID:              id,
		CompanionID:     companionID,
		Title:           title,
		Description:     description,
		Dialogue:        dialogue,
		Genre:           genre,
		MinMorale:       0.0,
		MaxMorale:       1.0,
		MinLoyalty:      0.0,
		MinRelationship: -1.0,
	}
}

// SetGenre updates the event's genre
func (e *CompanionEvent) SetGenre(g engine.GenreID) {
	e.Genre = g
}

// CanTrigger checks if the event can fire for the given companion
func (e *CompanionEvent) CanTrigger(c *Companion) bool {
	if e.Triggered {
		return false
	}
	if c.ID != e.CompanionID {
		return false
	}
	if !c.HasTrait(e.RequiredTrait) {
		return false
	}
	if c.Morale < e.MinMorale || c.Morale > e.MaxMorale {
		return false
	}
	if c.Loyalty < e.MinLoyalty {
		return false
	}
	if c.RelationshipWithPlayer < e.MinRelationship {
		return false
	}
	return true
}

// Trigger marks the event as triggered and applies effects
func (e *CompanionEvent) Trigger(c *Companion) {
	e.Triggered = true
	c.AdjustMorale(e.MoraleChange)
	c.AdjustLoyalty(e.LoyaltyChange)
	c.AdjustRelationshipWithPlayer(e.RelationshipChange)
	if e.SkillGain > 0 {
		c.GainExperience(e.SkillGain * 25) // Convert skill gain to XP
	}
}

// CompanionManager tracks and manages all active companions
type CompanionManager struct {
	Companions []*Companion
	Events     []*CompanionEvent
	Genre      engine.GenreID
	MaxSize    int
}

// NewCompanionManager creates a new companion manager
func NewCompanionManager(genre engine.GenreID, maxSize int) *CompanionManager {
	return &CompanionManager{
		Companions: make([]*Companion, 0),
		Events:     make([]*CompanionEvent, 0),
		Genre:      genre,
		MaxSize:    maxSize,
	}
}

// SetGenre updates all companions and events to new genre
func (m *CompanionManager) SetGenre(g engine.GenreID) {
	m.Genre = g
	for _, c := range m.Companions {
		c.SetGenre(g)
	}
	for _, e := range m.Events {
		e.SetGenre(g)
	}
}

// AddCompanion adds a companion to the party
func (m *CompanionManager) AddCompanion(c *Companion) bool {
	if len(m.Companions) >= m.MaxSize {
		return false
	}
	m.Companions = append(m.Companions, c)
	return true
}

// RemoveCompanion removes a companion from the party
func (m *CompanionManager) RemoveCompanion(id int) bool {
	for i, c := range m.Companions {
		if c.ID == id {
			m.Companions = append(m.Companions[:i], m.Companions[i+1:]...)
			return true
		}
	}
	return false
}

// GetCompanion retrieves a companion by ID
func (m *CompanionManager) GetCompanion(id int) *Companion {
	for _, c := range m.Companions {
		if c.ID == id {
			return c
		}
	}
	return nil
}

// GetCompanionByRole retrieves the first companion with a specific role
func (m *CompanionManager) GetCompanionByRole(role CompanionRole) *Companion {
	for _, c := range m.Companions {
		if c.Role == role && c.IsActive {
			return c
		}
	}
	return nil
}

// AddEvent registers an event for triggering
func (m *CompanionManager) AddEvent(e *CompanionEvent) {
	m.Events = append(m.Events, e)
}

// CheckEvents returns events that can trigger for active companions
func (m *CompanionManager) CheckEvents() []*CompanionEvent {
	triggerable := make([]*CompanionEvent, 0)
	for _, e := range m.Events {
		for _, c := range m.Companions {
			if e.CanTrigger(c) {
				triggerable = append(triggerable, e)
				break
			}
		}
	}
	return triggerable
}

// Tick processes one turn for all companions
func (m *CompanionManager) Tick() {
	for _, c := range m.Companions {
		if c.IsActive {
			c.Tick()
		}
	}
}

// ActiveCount returns the number of active companions
func (m *CompanionManager) ActiveCount() int {
	count := 0
	for _, c := range m.Companions {
		if c.IsActive {
			count++
		}
	}
	return count
}

// TotalSkillLevel returns sum of all active companion skill levels
func (m *CompanionManager) TotalSkillLevel() int {
	total := 0
	for _, c := range m.Companions {
		if c.IsActive {
			total += c.SkillLevel
		}
	}
	return total
}

// AverageMorale returns average morale of active companions
func (m *CompanionManager) AverageMorale() float64 {
	if len(m.Companions) == 0 {
		return 0
	}
	total := 0.0
	count := 0
	for _, c := range m.Companions {
		if c.IsActive {
			total += c.Morale
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

// RoleName returns the genre-appropriate name for a role
func RoleName(role CompanionRole, genre engine.GenreID) string {
	roleNames := map[engine.GenreID]map[CompanionRole]string{
		engine.GenreFantasy: {
			RoleGuide:      "Wizard Guide",
			RoleScout:      "Ranger Scout",
			RoleMedic:      "Healer",
			RoleWarrior:    "Knight",
			RoleTechnician: "Artificer",
			RoleLeader:     "Champion",
		},
		engine.GenreScifi: {
			RoleGuide:      "AI Navigator",
			RoleScout:      "Sensor Officer",
			RoleMedic:      "Ship's Doctor",
			RoleWarrior:    "Security Chief",
			RoleTechnician: "Chief Engineer",
			RoleLeader:     "First Officer",
		},
		engine.GenreHorror: {
			RoleGuide:      "Occultist",
			RoleScout:      "Lookout",
			RoleMedic:      "Field Medic",
			RoleWarrior:    "Survivor",
			RoleTechnician: "Mechanic",
			RoleLeader:     "Group Leader",
		},
		engine.GenreCyberpunk: {
			RoleGuide:      "Netrunner",
			RoleScout:      "Street Samurai",
			RoleMedic:      "Ripperdoc",
			RoleWarrior:    "Solo",
			RoleTechnician: "Techie",
			RoleLeader:     "Fixer",
		},
		engine.GenrePostapoc: {
			RoleGuide:      "Wasteland Guide",
			RoleScout:      "Scout",
			RoleMedic:      "Rad-Doc",
			RoleWarrior:    "Enforcer",
			RoleTechnician: "Scrapper",
			RoleLeader:     "Warlord",
		},
	}

	if genreRoles, ok := roleNames[genre]; ok {
		if name, ok := genreRoles[role]; ok {
			return name
		}
	}
	return roleNames[engine.GenreFantasy][role]
}

// TraitName returns the human-readable name for a personality trait
func TraitName(trait PersonalityTrait) string {
	names := map[PersonalityTrait]string{
		TraitBrave:         "Brave",
		TraitCautious:      "Cautious",
		TraitOptimistic:    "Optimistic",
		TraitPessimistic:   "Pessimistic",
		TraitLoyal:         "Loyal",
		TraitIndependent:   "Independent",
		TraitCompassionate: "Compassionate",
		TraitPragmatic:     "Pragmatic",
	}
	if name, ok := names[trait]; ok {
		return name
	}
	return "Unknown"
}

// AllCompanionRoles returns all companion roles
func AllCompanionRoles() []CompanionRole {
	return []CompanionRole{
		RoleGuide, RoleScout, RoleMedic, RoleWarrior, RoleTechnician, RoleLeader,
	}
}

// AllPersonalityTraits returns all personality traits
func AllPersonalityTraits() []PersonalityTrait {
	return []PersonalityTrait{
		TraitBrave, TraitCautious, TraitOptimistic, TraitPessimistic,
		TraitLoyal, TraitIndependent, TraitCompassionate, TraitPragmatic,
	}
}
