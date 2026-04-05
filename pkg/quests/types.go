package quests

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// QuestID uniquely identifies a quest.
type QuestID int

// QuestType identifies the type of quest.
type QuestType int

const (
	// TypeDelivery transports goods between locations.
	TypeDelivery QuestType = iota
	// TypeRescue saves stranded people.
	TypeRescue
	// TypeRetrieve finds and returns an item.
	TypeRetrieve
	// TypeExplore maps or investigates a location.
	TypeExplore
	// TypeEliminate deals with a threat.
	TypeEliminate
)

// AllQuestTypes returns all quest types.
func AllQuestTypes() []QuestType {
	return []QuestType{
		TypeDelivery,
		TypeRescue,
		TypeRetrieve,
		TypeExplore,
		TypeEliminate,
	}
}

// QuestTypeName returns the genre-appropriate name for a quest type.
func QuestTypeName(t QuestType, genre engine.GenreID) string {
	names := questTypeNames[genre]
	if names == nil {
		names = questTypeNames[engine.GenreFantasy]
	}
	return names[t]
}

var questTypeNames = map[engine.GenreID]map[QuestType]string{
	engine.GenreFantasy: {
		TypeDelivery:  "Courier Mission",
		TypeRescue:    "Rescue Quest",
		TypeRetrieve:  "Artifact Hunt",
		TypeExplore:   "Exploration",
		TypeEliminate: "Monster Slaying",
	},
	engine.GenreScifi: {
		TypeDelivery:  "Cargo Run",
		TypeRescue:    "Search and Rescue",
		TypeRetrieve:  "Salvage Operation",
		TypeExplore:   "Survey Mission",
		TypeEliminate: "Threat Elimination",
	},
	engine.GenreHorror: {
		TypeDelivery:  "Supply Run",
		TypeRescue:    "Survivor Rescue",
		TypeRetrieve:  "Scavenger Hunt",
		TypeExplore:   "Recon Mission",
		TypeEliminate: "Clear the Area",
	},
	engine.GenreCyberpunk: {
		TypeDelivery:  "Data Courier",
		TypeRescue:    "Extraction Job",
		TypeRetrieve:  "Acquisition",
		TypeExplore:   "Recon",
		TypeEliminate: "Wetwork",
	},
	engine.GenrePostapoc: {
		TypeDelivery:  "Caravan Run",
		TypeRescue:    "Rescue Op",
		TypeRetrieve:  "Salvage Run",
		TypeExplore:   "Scouting",
		TypeEliminate: "Clear Out",
	},
}

// QuestStatus tracks quest progress.
type QuestStatus int

const (
	// StatusAvailable means the quest can be accepted.
	StatusAvailable QuestStatus = iota
	// StatusActive means the quest is in progress.
	StatusActive
	// StatusCompleted means the quest was finished successfully.
	StatusCompleted
	// StatusFailed means the quest failed (time expired, etc.).
	StatusFailed
	// StatusDeclined means the player declined the quest.
	StatusDeclined
)

// StatusName returns the display name for a status.
func StatusName(s QuestStatus) string {
	names := map[QuestStatus]string{
		StatusAvailable: "Available",
		StatusActive:    "Active",
		StatusCompleted: "Completed",
		StatusFailed:    "Failed",
		StatusDeclined:  "Declined",
	}
	if name, ok := names[s]; ok {
		return name
	}
	return "Unknown"
}

// QuestReward represents the rewards for completing a quest.
type QuestReward struct {
	Currency    float64
	Food        float64
	Water       float64
	Fuel        float64
	Medicine    float64
	Morale      float64
	Reputation  int // Faction reputation if applicable
	FactionID   int // 0 = no faction
	SpecialItem string
}

// QuestObjective represents a quest goal.
type QuestObjective struct {
	Description string
	TargetX     int
	TargetY     int
	TargetName  string
	Completed   bool
}

// Quest represents a procedurally generated quest.
type Quest struct {
	ID          QuestID
	Type        QuestType
	Title       string
	Description string
	Genre       engine.GenreID
	Status      QuestStatus
	Objectives  []QuestObjective
	Reward      QuestReward
	TimeLimit   int // Turns remaining, 0 = no limit
	GiverName   string
	OriginX     int
	OriginY     int
}

// NewQuest creates a new quest.
func NewQuest(id QuestID, qType QuestType, title, description string, genre engine.GenreID) *Quest {
	return &Quest{
		ID:          id,
		Type:        qType,
		Title:       title,
		Description: description,
		Genre:       genre,
		Status:      StatusAvailable,
		Objectives:  make([]QuestObjective, 0),
	}
}

// SetGenre updates the quest's genre.
func (q *Quest) SetGenre(genre engine.GenreID) {
	q.Genre = genre
}

// TypeDisplayName returns the genre-appropriate type name.
func (q *Quest) TypeDisplayName() string {
	return QuestTypeName(q.Type, q.Genre)
}

// AddObjective adds an objective to the quest.
func (q *Quest) AddObjective(desc string, x, y int, targetName string) {
	q.Objectives = append(q.Objectives, QuestObjective{
		Description: desc,
		TargetX:     x,
		TargetY:     y,
		TargetName:  targetName,
	})
}

// CompleteObjective marks an objective as complete.
func (q *Quest) CompleteObjective(index int) {
	if index >= 0 && index < len(q.Objectives) {
		q.Objectives[index].Completed = true
	}
}

// IsComplete returns true if all objectives are completed.
func (q *Quest) IsComplete() bool {
	for _, obj := range q.Objectives {
		if !obj.Completed {
			return false
		}
	}
	return len(q.Objectives) > 0
}

// Accept marks the quest as active.
func (q *Quest) Accept() {
	if q.Status == StatusAvailable {
		q.Status = StatusActive
	}
}

// Decline marks the quest as declined.
func (q *Quest) Decline() {
	if q.Status == StatusAvailable {
		q.Status = StatusDeclined
	}
}

// Complete marks the quest as completed.
func (q *Quest) Complete() {
	if q.Status == StatusActive {
		q.Status = StatusCompleted
	}
}

// Fail marks the quest as failed.
func (q *Quest) Fail() {
	if q.Status == StatusActive {
		q.Status = StatusFailed
	}
}

// AdvanceTime decrements the time limit and fails if expired.
func (q *Quest) AdvanceTime(turns int) {
	if q.Status != StatusActive || q.TimeLimit <= 0 {
		return
	}
	q.TimeLimit -= turns
	if q.TimeLimit <= 0 {
		q.TimeLimit = 0
		q.Fail()
	}
}

// QuestTracker manages all quests in a run.
type QuestTracker struct {
	Quests      map[QuestID]*Quest
	ActiveLimit int
	genre       engine.GenreID
}

// NewQuestTracker creates a new quest tracker.
func NewQuestTracker(genre engine.GenreID) *QuestTracker {
	return &QuestTracker{
		Quests:      make(map[QuestID]*Quest),
		ActiveLimit: 3, // Max 3 active quests by default
		genre:       genre,
	}
}

// SetGenre updates the tracker's genre and all tracked quests.
func (t *QuestTracker) SetGenre(genre engine.GenreID) {
	t.genre = genre
	for _, q := range t.Quests {
		q.SetGenre(genre)
	}
}

// AddQuest adds a quest to the tracker.
func (t *QuestTracker) AddQuest(q *Quest) {
	t.Quests[q.ID] = q
}

// GetQuest returns a quest by ID.
func (t *QuestTracker) GetQuest(id QuestID) *Quest {
	return t.Quests[id]
}

// AcceptQuest accepts a quest if under the active limit.
func (t *QuestTracker) AcceptQuest(id QuestID) bool {
	if t.ActiveQuestCount() >= t.ActiveLimit {
		return false
	}
	if q := t.Quests[id]; q != nil {
		q.Accept()
		return true
	}
	return false
}

// ActiveQuestCount returns the number of active quests.
func (t *QuestTracker) ActiveQuestCount() int {
	count := 0
	for _, q := range t.Quests {
		if q.Status == StatusActive {
			count++
		}
	}
	return count
}

// ActiveQuests returns all active quests.
func (t *QuestTracker) ActiveQuests() []*Quest {
	var result []*Quest
	for _, q := range t.Quests {
		if q.Status == StatusActive {
			result = append(result, q)
		}
	}
	return result
}

// AvailableQuests returns all available quests.
func (t *QuestTracker) AvailableQuests() []*Quest {
	var result []*Quest
	for _, q := range t.Quests {
		if q.Status == StatusAvailable {
			result = append(result, q)
		}
	}
	return result
}

// CompletedQuests returns all completed quests.
func (t *QuestTracker) CompletedQuests() []*Quest {
	var result []*Quest
	for _, q := range t.Quests {
		if q.Status == StatusCompleted {
			result = append(result, q)
		}
	}
	return result
}

// AdvanceAllTime advances time for all active quests.
func (t *QuestTracker) AdvanceAllTime(turns int) {
	for _, q := range t.Quests {
		q.AdvanceTime(turns)
	}
}

// CheckObjectivesAt checks if any objectives are at the given position.
func (t *QuestTracker) CheckObjectivesAt(x, y int) []*Quest {
	var result []*Quest
	for _, q := range t.Quests {
		if q.Status != StatusActive {
			continue
		}
		for i := range q.Objectives {
			if !q.Objectives[i].Completed && q.Objectives[i].TargetX == x && q.Objectives[i].TargetY == y {
				result = append(result, q)
				break
			}
		}
	}
	return result
}
