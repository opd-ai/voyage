package narrative

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// StoryAct represents a phase of the three-act structure.
type StoryAct int

const (
	// ActDeparture is the opening crisis that starts the journey.
	ActDeparture StoryAct = iota
	// ActMidJourney is the revelation that changes the stakes.
	ActMidJourney
	// ActArrival is the twist at the destination.
	ActArrival
)

// AllStoryActs returns all story acts.
func AllStoryActs() []StoryAct {
	return []StoryAct{ActDeparture, ActMidJourney, ActArrival}
}

// ActName returns the display name for an act.
func ActName(act StoryAct) string {
	names := map[StoryAct]string{
		ActDeparture:  "Departure",
		ActMidJourney: "Mid-Journey",
		ActArrival:    "Arrival",
	}
	if name, ok := names[act]; ok {
		return name
	}
	return "Unknown"
}

// RecurringNPCRole represents the NPC's relationship to the party.
type RecurringNPCRole int

const (
	// RoleFriend actively helps the party.
	RoleFriend RecurringNPCRole = iota
	// RoleNemesis opposes the party.
	RoleNemesis
	// RoleAmbiguous has unclear motives.
	RoleAmbiguous
)

// AllNPCRoles returns all NPC roles.
func AllNPCRoles() []RecurringNPCRole {
	return []RecurringNPCRole{RoleFriend, RoleNemesis, RoleAmbiguous}
}

// RoleName returns the display name for a role.
func RoleName(role RecurringNPCRole) string {
	names := map[RecurringNPCRole]string{
		RoleFriend:    "Friend",
		RoleNemesis:   "Nemesis",
		RoleAmbiguous: "Ambiguous",
	}
	if name, ok := names[role]; ok {
		return name
	}
	return "Unknown"
}

// StoryBeat represents a narrative moment in the journey.
type StoryBeat struct {
	Act         StoryAct
	Title       string
	Description string
	Triggered   bool
	TriggerTurn int // 0 = start, positive = mid/arrival trigger
}

// RecurringNPC represents a named character that reappears.
type RecurringNPC struct {
	Name        string
	Role        RecurringNPCRole
	Description string
	Genre       engine.GenreID
	Appearances int
	Dialogues   []string
}

// NewRecurringNPC creates a new recurring NPC.
func NewRecurringNPC(name string, role RecurringNPCRole, desc string, genre engine.GenreID) *RecurringNPC {
	return &RecurringNPC{
		Name:        name,
		Role:        role,
		Description: desc,
		Genre:       genre,
		Dialogues:   make([]string, 0),
	}
}

// SetGenre updates the NPC's genre.
func (n *RecurringNPC) SetGenre(genre engine.GenreID) {
	n.Genre = genre
}

// AddDialogue adds a dialogue line.
func (n *RecurringNPC) AddDialogue(line string) {
	n.Dialogues = append(n.Dialogues, line)
}

// RecordAppearance increments the appearance counter.
func (n *RecurringNPC) RecordAppearance() {
	n.Appearances++
}

// CrewBackstory links a crew member's past to the narrative.
type CrewBackstory struct {
	CrewID          int
	CrewName        string
	BackstoryHook   string // Brief teaser revealed at start
	FullBackstory   string // Full story revealed mid-journey
	DestinationLink string // Connection to the destination
	Revealed        bool
}

// NewCrewBackstory creates a new crew backstory.
func NewCrewBackstory(crewID int, name, hook, full, link string) *CrewBackstory {
	return &CrewBackstory{
		CrewID:          crewID,
		CrewName:        name,
		BackstoryHook:   hook,
		FullBackstory:   full,
		DestinationLink: link,
	}
}

// Reveal marks the backstory as revealed.
func (b *CrewBackstory) Reveal() {
	b.Revealed = true
}

// StoryArc contains all narrative elements for a run.
type StoryArc struct {
	Genre           engine.GenreID
	Beats           []*StoryBeat
	RecurringNPC    *RecurringNPC
	CrewBackstories []*CrewBackstory
	CurrentAct      StoryAct
}

// NewStoryArc creates a new story arc.
func NewStoryArc(genre engine.GenreID) *StoryArc {
	return &StoryArc{
		Genre:           genre,
		Beats:           make([]*StoryBeat, 0),
		CrewBackstories: make([]*CrewBackstory, 0),
		CurrentAct:      ActDeparture,
	}
}

// SetGenre updates the story arc's genre.
func (s *StoryArc) SetGenre(genre engine.GenreID) {
	s.Genre = genre
	if s.RecurringNPC != nil {
		s.RecurringNPC.SetGenre(genre)
	}
}

// AddBeat adds a story beat.
func (s *StoryArc) AddBeat(beat *StoryBeat) {
	s.Beats = append(s.Beats, beat)
}

// AddCrewBackstory adds a crew backstory.
func (s *StoryArc) AddCrewBackstory(backstory *CrewBackstory) {
	s.CrewBackstories = append(s.CrewBackstories, backstory)
}

// GetBeatForAct returns the story beat for a given act.
func (s *StoryArc) GetBeatForAct(act StoryAct) *StoryBeat {
	for _, beat := range s.Beats {
		if beat.Act == act {
			return beat
		}
	}
	return nil
}

// TriggerBeat marks a beat as triggered.
func (s *StoryArc) TriggerBeat(act StoryAct, turn int) bool {
	if beat := s.GetBeatForAct(act); beat != nil && !beat.Triggered {
		beat.Triggered = true
		beat.TriggerTurn = turn
		s.CurrentAct = act
		return true
	}
	return false
}

// GetCrewBackstory returns the backstory for a crew member.
func (s *StoryArc) GetCrewBackstory(crewID int) *CrewBackstory {
	for _, bs := range s.CrewBackstories {
		if bs.CrewID == crewID {
			return bs
		}
	}
	return nil
}

// AdvanceAct advances to the next act.
func (s *StoryArc) AdvanceAct() bool {
	switch s.CurrentAct {
	case ActDeparture:
		s.CurrentAct = ActMidJourney
		return true
	case ActMidJourney:
		s.CurrentAct = ActArrival
		return true
	default:
		return false
	}
}

// IsComplete returns true if all acts have been triggered.
func (s *StoryArc) IsComplete() bool {
	for _, beat := range s.Beats {
		if !beat.Triggered {
			return false
		}
	}
	return len(s.Beats) > 0
}

// GetActiveBeats returns all triggered beats.
func (s *StoryArc) GetActiveBeats() []*StoryBeat {
	var result []*StoryBeat
	for _, beat := range s.Beats {
		if beat.Triggered {
			result = append(result, beat)
		}
	}
	return result
}
