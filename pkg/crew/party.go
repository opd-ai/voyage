package crew

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// Party manages the crew roster.
type Party struct {
	genre    engine.GenreID
	members  []*CrewMember
	capacity int
}

// NewParty creates a new party with the given capacity.
func NewParty(genre engine.GenreID, capacity int) *Party {
	if capacity < 2 {
		capacity = 2
	}
	if capacity > 6 {
		capacity = 6
	}
	return &Party{
		genre:    genre,
		members:  make([]*CrewMember, 0, capacity),
		capacity: capacity,
	}
}

// SetGenre changes the party's genre.
func (p *Party) SetGenre(genre engine.GenreID) {
	p.genre = genre
}

// Genre returns the current genre.
func (p *Party) Genre() engine.GenreID {
	return p.genre
}

// Add adds a crew member to the party.
// Returns false if the party is full.
func (p *Party) Add(member *CrewMember) bool {
	if len(p.members) >= p.capacity {
		return false
	}
	p.members = append(p.members, member)
	return true
}

// Remove removes a crew member from the party.
func (p *Party) Remove(id int) bool {
	for i, m := range p.members {
		if m.ID == id {
			p.members = append(p.members[:i], p.members[i+1:]...)
			return true
		}
	}
	return false
}

// Get returns a crew member by ID.
func (p *Party) Get(id int) *CrewMember {
	for _, m := range p.members {
		if m.ID == id {
			return m
		}
	}
	return nil
}

// Members returns all crew members.
func (p *Party) Members() []*CrewMember {
	return p.members
}

// Living returns all living crew members.
func (p *Party) Living() []*CrewMember {
	living := make([]*CrewMember, 0)
	for _, m := range p.members {
		if m.IsAlive {
			living = append(living, m)
		}
	}
	return living
}

// Dead returns all dead crew members.
func (p *Party) Dead() []*CrewMember {
	dead := make([]*CrewMember, 0)
	for _, m := range p.members {
		if !m.IsAlive {
			dead = append(dead, m)
		}
	}
	return dead
}

// Count returns the total number of members.
func (p *Party) Count() int {
	return len(p.members)
}

// LivingCount returns the number of living members.
func (p *Party) LivingCount() int {
	count := 0
	for _, m := range p.members {
		if m.IsAlive {
			count++
		}
	}
	return count
}

// Capacity returns the party's maximum capacity.
func (p *Party) Capacity() int {
	return p.capacity
}

// IsFull returns true if the party is at capacity.
func (p *Party) IsFull() bool {
	return len(p.members) >= p.capacity
}

// IsEmpty returns true if the party has no living members.
func (p *Party) IsEmpty() bool {
	return p.LivingCount() == 0
}

// HasSkill returns true if any living member has the skill.
func (p *Party) HasSkill(skill Skill) bool {
	for _, m := range p.members {
		if m.IsAlive && m.Skill == skill {
			return true
		}
	}
	return false
}

// GetWithSkill returns the first living member with the skill.
func (p *Party) GetWithSkill(skill Skill) *CrewMember {
	for _, m := range p.members {
		if m.IsAlive && m.Skill == skill {
			return m
		}
	}
	return nil
}

// AverageHealth returns the average health of living members.
func (p *Party) AverageHealth() float64 {
	living := p.Living()
	if len(living) == 0 {
		return 0
	}
	total := 0.0
	for _, m := range living {
		total += m.Health
	}
	return total / float64(len(living))
}

// AdvanceDay updates all members for a new day.
func (p *Party) AdvanceDay() {
	for _, m := range p.members {
		if m.IsAlive {
			m.DaysWithParty++
		}
	}
}

// ApplyDamageToAll damages all living members.
func (p *Party) ApplyDamageToAll(amount float64) []string {
	deaths := make([]string, 0)
	for _, m := range p.members {
		if m.IsAlive && m.TakeDamage(amount) {
			deaths = append(deaths, m.Name)
		}
	}
	return deaths
}

// HealAll heals all living members.
func (p *Party) HealAll(amount float64) {
	for _, m := range p.members {
		if m.IsAlive {
			m.Heal(amount)
		}
	}
}
