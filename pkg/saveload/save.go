package saveload

import (
	"encoding/json"
	"time"

	"github.com/opd-ai/voyage/pkg/engine"
)

// SaveData represents the complete game state for serialization.
type SaveData struct {
	// Metadata
	Version    int       `json:"version"`
	Slot       int       `json:"slot"`
	SavedAt    time.Time `json:"savedAt"`
	PlayTime   int64     `json:"playTime"` // seconds
	IsAutosave bool      `json:"isAutosave"`

	// Core game state
	MasterSeed int64          `json:"masterSeed"`
	Genre      engine.GenreID `json:"genre"`
	Turn       int            `json:"turn"`
	Day        int            `json:"day"`

	// Vessel state
	Vessel VesselState `json:"vessel"`

	// Party state
	Party PartyState `json:"party"`

	// Resource state
	Resources ResourceState `json:"resources"`

	// World state
	World WorldState `json:"world"`

	// Event state
	Events EventState `json:"events"`

	// Game stats
	Stats StatsState `json:"stats"`
}

// VesselState holds vessel serialization data.
type VesselState struct {
	Type      int         `json:"type"`
	X         int         `json:"x"`
	Y         int         `json:"y"`
	Integrity float64     `json:"integrity"`
	Cargo     []CargoItem `json:"cargo"`
}

// CargoItem represents an item in cargo.
type CargoItem struct {
	Type     int     `json:"type"`
	Quantity float64 `json:"quantity"`
}

// PartyState holds party serialization data.
type PartyState struct {
	Capacity int         `json:"capacity"`
	Members  []CrewState `json:"members"`
}

// CrewState holds individual crew member data.
type CrewState struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Health        float64 `json:"health"`
	MaxHealth     float64 `json:"maxHealth"`
	Trait         int     `json:"trait"`
	Skill         int     `json:"skill"`
	IsAlive       bool    `json:"isAlive"`
	DaysWithParty int     `json:"daysWithParty"`
}

// ResourceState holds resource serialization data.
type ResourceState struct {
	Levels    map[int]float64 `json:"levels"`
	MaxLevels map[int]float64 `json:"maxLevels"`
}

// WorldState holds world map serialization data.
type WorldState struct {
	Width        int      `json:"width"`
	Height       int      `json:"height"`
	OriginX      int      `json:"originX"`
	OriginY      int      `json:"originY"`
	DestinationX int      `json:"destinationX"`
	DestinationY int      `json:"destinationY"`
	ExploredMap  [][]bool `json:"exploredMap"`
}

// EventState holds event queue serialization data.
type EventState struct {
	NextID        int   `json:"nextId"`
	ResolvedCount int   `json:"resolvedCount"`
	PendingIDs    []int `json:"pendingIds"`
}

// StatsState holds game statistics.
type StatsState struct {
	TilesExplored    int     `json:"tilesExplored"`
	DistanceTraveled int     `json:"distanceTraveled"`
	EventsResolved   int     `json:"eventsResolved"`
	CrewLost         int     `json:"crewLost"`
	FoodConsumed     float64 `json:"foodConsumed"`
	WaterConsumed    float64 `json:"waterConsumed"`
	FuelConsumed     float64 `json:"fuelConsumed"`
}

// CurrentVersion is the save file format version.
const CurrentVersion = 1

// NewSaveData creates a new SaveData with default metadata.
func NewSaveData(slot int, seed int64, genre engine.GenreID) *SaveData {
	return &SaveData{
		Version:    CurrentVersion,
		Slot:       slot,
		SavedAt:    time.Now(),
		MasterSeed: seed,
		Genre:      genre,
		Resources: ResourceState{
			Levels:    make(map[int]float64),
			MaxLevels: make(map[int]float64),
		},
	}
}

// Marshal serializes the save data to JSON.
func (sd *SaveData) Marshal() ([]byte, error) {
	return json.MarshalIndent(sd, "", "  ")
}

// Unmarshal deserializes JSON data into SaveData.
func Unmarshal(data []byte) (*SaveData, error) {
	sd := &SaveData{}
	if err := json.Unmarshal(data, sd); err != nil {
		return nil, err
	}
	return sd, nil
}

// Validate checks if the save data is valid.
func (sd *SaveData) Validate() error {
	if sd.Version < 1 || sd.Version > CurrentVersion {
		return ErrInvalidVersion
	}
	if sd.MasterSeed == 0 {
		return ErrInvalidSeed
	}
	if sd.Slot < 0 || sd.Slot > MaxSlots {
		return ErrInvalidSlot
	}
	return nil
}

// GetSummary returns a brief summary of the save.
func (sd *SaveData) GetSummary() SaveSummary {
	return SaveSummary{
		Slot:       sd.Slot,
		SavedAt:    sd.SavedAt,
		PlayTime:   time.Duration(sd.PlayTime) * time.Second,
		Genre:      sd.Genre,
		Turn:       sd.Turn,
		Day:        sd.Day,
		CrewCount:  countLivingCrew(sd.Party.Members),
		IsAutosave: sd.IsAutosave,
	}
}

// SaveSummary provides a brief overview of a save file.
type SaveSummary struct {
	Slot       int
	SavedAt    time.Time
	PlayTime   time.Duration
	Genre      engine.GenreID
	Turn       int
	Day        int
	CrewCount  int
	IsAutosave bool
}

// countLivingCrew counts living crew members.
func countLivingCrew(members []CrewState) int {
	count := 0
	for _, m := range members {
		if m.IsAlive {
			count++
		}
	}
	return count
}
