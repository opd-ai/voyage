package journey

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// LegID uniquely identifies a journey leg
type LegID int

// DifficultyLevel represents the challenge rating of a leg
type DifficultyLevel int

const (
	DifficultyEasy DifficultyLevel = iota
	DifficultyNormal
	DifficultyHard
	DifficultyExtreme
)

// StopoverService represents a service available at a stopover city
type StopoverService int

const (
	ServiceTrading StopoverService = iota
	ServiceRepairs
	ServiceRecruitment
	ServiceUpgrades
	ServiceInformation
	ServiceHealing
)

// Leg represents a single journey segment in a multi-leg campaign
type Leg struct {
	ID          LegID
	Name        string
	Description string
	Genre       engine.GenreID

	// Origin and destination
	OriginName      string
	DestinationName string

	// Geography
	Distance    int      // Units of travel distance
	TerrainType string   // Primary terrain
	Hazards     []string // Environmental hazards

	// Difficulty scaling
	Difficulty       DifficultyLevel
	EnemyStrength    float64 // Multiplier for enemy encounters
	ResourceScarcity float64 // Multiplier for resource costs
	EventFrequency   float64 // Multiplier for event occurrence

	// State
	Started   bool
	Completed bool
	DaysTaken int
	Survivors int
}

// NewLeg creates a new journey leg
func NewLeg(id LegID, name, origin, destination string, distance int, difficulty DifficultyLevel, genre engine.GenreID) *Leg {
	// Calculate difficulty multipliers
	var strength, scarcity, frequency float64
	switch difficulty {
	case DifficultyEasy:
		strength, scarcity, frequency = 0.8, 0.9, 0.9
	case DifficultyNormal:
		strength, scarcity, frequency = 1.0, 1.0, 1.0
	case DifficultyHard:
		strength, scarcity, frequency = 1.3, 1.2, 1.2
	case DifficultyExtreme:
		strength, scarcity, frequency = 1.6, 1.5, 1.5
	}

	return &Leg{
		ID:               id,
		Name:             name,
		Genre:            genre,
		OriginName:       origin,
		DestinationName:  destination,
		Distance:         distance,
		Difficulty:       difficulty,
		EnemyStrength:    strength,
		ResourceScarcity: scarcity,
		EventFrequency:   frequency,
		Hazards:          make([]string, 0),
	}
}

// SetGenre updates the leg's genre
func (l *Leg) SetGenre(g engine.GenreID) {
	l.Genre = g
}

// Start marks the leg as started
func (l *Leg) Start() {
	l.Started = true
}

// Complete marks the leg as completed with statistics
func (l *Leg) Complete(days, survivors int) {
	l.Completed = true
	l.DaysTaken = days
	l.Survivors = survivors
}

// DifficultyName returns the human-readable difficulty name
func DifficultyName(d DifficultyLevel) string {
	names := map[DifficultyLevel]string{
		DifficultyEasy:    "Easy",
		DifficultyNormal:  "Normal",
		DifficultyHard:    "Hard",
		DifficultyExtreme: "Extreme",
	}
	if name, ok := names[d]; ok {
		return name
	}
	return "Unknown"
}

// AllDifficultyLevels returns all difficulty levels in order
func AllDifficultyLevels() []DifficultyLevel {
	return []DifficultyLevel{DifficultyEasy, DifficultyNormal, DifficultyHard, DifficultyExtreme}
}

// Stopover represents an intermediate hub city between journey legs
type Stopover struct {
	ID          int
	Name        string
	Description string
	Genre       engine.GenreID

	// Location between legs
	AfterLeg LegID

	// Available services
	Services []StopoverService

	// Trading modifiers
	BuyPriceModifier  float64 // Multiplier for buying resources
	SellPriceModifier float64 // Multiplier for selling resources

	// Special features
	Features    []string
	Inhabitants string // Description of who lives here

	// State
	Visited bool
}

// NewStopover creates a new stopover city
func NewStopover(id int, name, description string, afterLeg LegID, genre engine.GenreID) *Stopover {
	return &Stopover{
		ID:                id,
		Name:              name,
		Description:       description,
		AfterLeg:          afterLeg,
		Genre:             genre,
		Services:          make([]StopoverService, 0),
		BuyPriceModifier:  1.0,
		SellPriceModifier: 1.0,
		Features:          make([]string, 0),
	}
}

// SetGenre updates the stopover's genre
func (s *Stopover) SetGenre(g engine.GenreID) {
	s.Genre = g
}

// AddService adds a service to the stopover
func (s *Stopover) AddService(service StopoverService) {
	s.Services = append(s.Services, service)
}

// HasService checks if the stopover offers a specific service
func (s *Stopover) HasService(service StopoverService) bool {
	for _, svc := range s.Services {
		if svc == service {
			return true
		}
	}
	return false
}

// Visit marks the stopover as visited
func (s *Stopover) Visit() {
	s.Visited = true
}

// ServiceName returns the human-readable service name by genre
func ServiceName(service StopoverService, genre engine.GenreID) string {
	serviceNames := map[engine.GenreID]map[StopoverService]string{
		engine.GenreFantasy: {
			ServiceTrading:     "Marketplace",
			ServiceRepairs:     "Blacksmith",
			ServiceRecruitment: "Tavern",
			ServiceUpgrades:    "Enchanter",
			ServiceInformation: "Sage's Tower",
			ServiceHealing:     "Temple",
		},
		engine.GenreScifi: {
			ServiceTrading:     "Trade Hub",
			ServiceRepairs:     "Repair Bay",
			ServiceRecruitment: "Crew Exchange",
			ServiceUpgrades:    "Tech Lab",
			ServiceInformation: "Data Center",
			ServiceHealing:     "Medical Bay",
		},
		engine.GenreHorror: {
			ServiceTrading:     "Barter Post",
			ServiceRepairs:     "Workshop",
			ServiceRecruitment: "Refuge Hall",
			ServiceUpgrades:    "Tinkerer's Den",
			ServiceInformation: "Watch Tower",
			ServiceHealing:     "Infirmary",
		},
		engine.GenreCyberpunk: {
			ServiceTrading:     "Black Market",
			ServiceRepairs:     "Chop Shop",
			ServiceRecruitment: "Fixer's Joint",
			ServiceUpgrades:    "Upgrade Clinic",
			ServiceInformation: "Info Broker",
			ServiceHealing:     "Ripperdoc",
		},
		engine.GenrePostapoc: {
			ServiceTrading:     "Trade Shack",
			ServiceRepairs:     "Scrapyard",
			ServiceRecruitment: "Survivor Camp",
			ServiceUpgrades:    "Jury-Rig Station",
			ServiceInformation: "Lookout Post",
			ServiceHealing:     "Field Medic",
		},
	}

	if genreServices, ok := serviceNames[genre]; ok {
		if name, ok := genreServices[service]; ok {
			return name
		}
	}
	// Fallback
	return serviceNames[engine.GenreFantasy][service]
}

// AllStopoverServices returns all stopover service types
func AllStopoverServices() []StopoverService {
	return []StopoverService{
		ServiceTrading, ServiceRepairs, ServiceRecruitment,
		ServiceUpgrades, ServiceInformation, ServiceHealing,
	}
}

// CampaignState holds persistent state that carries between journey legs
type CampaignState struct {
	// Resources carried forward
	Gold       int
	Food       int
	Supplies   int
	Reputation int

	// Crew state
	CrewCount  int
	CrewHealth float64
	CrewMorale float64

	// Vessel state
	VesselHealth float64
	VesselCargo  int

	// Progress tracking
	TotalDistance     int
	TotalDays         int
	TotalDeaths       int
	EventsEncountered int

	// Unlocks and achievements
	Achievements   []string
	DiscoveredLore []string
}

// NewCampaignState creates a fresh campaign state with starting values
func NewCampaignState() *CampaignState {
	return &CampaignState{
		Gold:           100,
		Food:           50,
		Supplies:       30,
		Reputation:     0,
		CrewCount:      4,
		CrewHealth:     1.0,
		CrewMorale:     0.8,
		VesselHealth:   1.0,
		VesselCargo:    0,
		Achievements:   make([]string, 0),
		DiscoveredLore: make([]string, 0),
	}
}

// ApplyLegResults updates state based on completed leg
func (cs *CampaignState) ApplyLegResults(leg *Leg, survivors, daysUsed, goldEarned, foodUsed int) {
	cs.TotalDistance += leg.Distance
	cs.TotalDays += daysUsed
	cs.TotalDeaths += cs.CrewCount - survivors
	cs.CrewCount = survivors
	cs.Gold += goldEarned
	cs.Food -= foodUsed
	if cs.Food < 0 {
		cs.Food = 0
	}
}

// AddAchievement records a new achievement
func (cs *CampaignState) AddAchievement(achievement string) {
	cs.Achievements = append(cs.Achievements, achievement)
}

// AddLore records discovered lore
func (cs *CampaignState) AddLore(loreID string) {
	cs.DiscoveredLore = append(cs.DiscoveredLore, loreID)
}

// Campaign represents a complete multi-leg journey
type Campaign struct {
	ID          string
	Name        string
	Description string
	Genre       engine.GenreID

	// Journey structure
	Legs      []*Leg
	Stopovers []*Stopover

	// Persistent state
	State *CampaignState

	// Progress tracking
	CurrentLegIndex int
	StartedAt       int64 // Unix timestamp
	CompletedAt     int64

	// Genre variation
	AllowGenreShifts bool
	LegGenres        []engine.GenreID
}

// NewCampaign creates a new multi-leg campaign
func NewCampaign(id, name, description string, genre engine.GenreID) *Campaign {
	return &Campaign{
		ID:               id,
		Name:             name,
		Description:      description,
		Genre:            genre,
		Legs:             make([]*Leg, 0),
		Stopovers:        make([]*Stopover, 0),
		State:            NewCampaignState(),
		CurrentLegIndex:  0,
		AllowGenreShifts: false,
		LegGenres:        make([]engine.GenreID, 0),
	}
}

// SetGenre updates the campaign's genre (and optionally all components)
func (c *Campaign) SetGenre(g engine.GenreID) {
	c.Genre = g
	if !c.AllowGenreShifts {
		for _, leg := range c.Legs {
			leg.SetGenre(g)
		}
		for _, stopover := range c.Stopovers {
			stopover.SetGenre(g)
		}
	}
}

// AddLeg adds a journey leg to the campaign
func (c *Campaign) AddLeg(leg *Leg) {
	c.Legs = append(c.Legs, leg)
	c.LegGenres = append(c.LegGenres, leg.Genre)
}

// AddStopover adds a stopover city to the campaign
func (c *Campaign) AddStopover(stopover *Stopover) {
	c.Stopovers = append(c.Stopovers, stopover)
}

// CurrentLeg returns the current journey leg
func (c *Campaign) CurrentLeg() *Leg {
	if c.CurrentLegIndex < len(c.Legs) {
		return c.Legs[c.CurrentLegIndex]
	}
	return nil
}

// NextStopover returns the stopover after the current leg
func (c *Campaign) NextStopover() *Stopover {
	if c.CurrentLegIndex < len(c.Stopovers) {
		return c.Stopovers[c.CurrentLegIndex]
	}
	return nil
}

// CompleteLeg finishes the current leg and advances to the next
func (c *Campaign) CompleteLeg(days, survivors int) bool {
	leg := c.CurrentLeg()
	if leg == nil {
		return false
	}

	leg.Complete(days, survivors)
	c.State.ApplyLegResults(leg, survivors, days, 0, 0)
	c.CurrentLegIndex++

	return c.CurrentLegIndex < len(c.Legs)
}

// LegCount returns the number of legs in the campaign
func (c *Campaign) LegCount() int {
	return len(c.Legs)
}

// IsComplete returns true if all legs are completed
func (c *Campaign) IsComplete() bool {
	return c.CurrentLegIndex >= len(c.Legs)
}

// Progress returns the campaign completion percentage
func (c *Campaign) Progress() float64 {
	if len(c.Legs) == 0 {
		return 0
	}
	return float64(c.CurrentLegIndex) / float64(len(c.Legs)) * 100
}

// TotalDistance returns the sum of all leg distances
func (c *Campaign) TotalDistance() int {
	total := 0
	for _, leg := range c.Legs {
		total += leg.Distance
	}
	return total
}

// GetStopoverAfterLeg returns the stopover that comes after a specific leg
func (c *Campaign) GetStopoverAfterLeg(legID LegID) *Stopover {
	for _, s := range c.Stopovers {
		if s.AfterLeg == legID {
			return s
		}
	}
	return nil
}

// EnableGenreShifts allows different genres per leg
func (c *Campaign) EnableGenreShifts() {
	c.AllowGenreShifts = true
}
