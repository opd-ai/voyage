package convoy

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/opd-ai/voyage/pkg/engine"
)

// Convoy errors.
var (
	ErrConvoyNotFound    = errors.New("convoy not found")
	ErrConvoyFull        = errors.New("convoy is full")
	ErrPlayerNotInConvoy = errors.New("player not in convoy")
	ErrInvalidConvoyCode = errors.New("invalid convoy code")
	ErrConvoyStarted     = errors.New("convoy has already started")
)

// ConvoyID uniquely identifies a convoy session.
type ConvoyID string

// PlayerID uniquely identifies a player within a convoy.
type PlayerID string

// ConvoyState represents the current state of a convoy.
type ConvoyState int

const (
	// StateWaiting indicates convoy is waiting for players to join.
	StateWaiting ConvoyState = iota
	// StateStarted indicates all players have started their runs.
	StateStarted
	// StateCompleted indicates all players have finished.
	StateCompleted
)

// Convoy represents a shared-seed multiplayer session.
type Convoy struct {
	mu sync.RWMutex

	ID        ConvoyID       `json:"id"`
	Code      string         `json:"code"` // Short shareable code
	Seed      int64          `json:"seed"`
	Genre     engine.GenreID `json:"genre"`
	State     ConvoyState    `json:"state"`
	CreatedAt time.Time      `json:"createdAt"`
	StartedAt time.Time      `json:"startedAt,omitempty"`

	MaxPlayers int        `json:"maxPlayers"`
	Players    []*Player  `json:"players"`
	HostID     PlayerID   `json:"hostId"`
	Runs       []*RunData `json:"runs,omitempty"`
}

// Player represents a participant in a convoy.
type Player struct {
	ID          PlayerID  `json:"id"`
	Name        string    `json:"name"`
	JoinedAt    time.Time `json:"joinedAt"`
	IsHost      bool      `json:"isHost"`
	IsReady     bool      `json:"isReady"`
	HasFinished bool      `json:"hasFinished"`
}

// RunData tracks a player's run progress and results.
type RunData struct {
	PlayerID     PlayerID       `json:"playerId"`
	Genre        engine.GenreID `json:"genre"`
	StartedAt    time.Time      `json:"startedAt"`
	FinishedAt   time.Time      `json:"finishedAt,omitempty"`
	IsVictory    bool           `json:"isVictory"`
	Score        int            `json:"score"`
	DaysTraveled int            `json:"daysTraveled"`
	Survivors    int            `json:"survivors"`
	EventsSeen   int            `json:"eventsSeen"`
}

// NewConvoy creates a new convoy with the given seed and genre.
func NewConvoy(seed int64, genre engine.GenreID, maxPlayers int) *Convoy {
	if maxPlayers < 2 {
		maxPlayers = 2
	}
	if maxPlayers > 10 {
		maxPlayers = 10
	}

	return &Convoy{
		ID:         generateConvoyID(),
		Code:       generateConvoyCode(),
		Seed:       seed,
		Genre:      genre,
		State:      StateWaiting,
		CreatedAt:  time.Now().UTC(),
		MaxPlayers: maxPlayers,
		Players:    make([]*Player, 0),
		Runs:       make([]*RunData, 0),
	}
}

// generateConvoyID creates a unique convoy identifier.
func generateConvoyID() ConvoyID {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return ConvoyID(base64.RawURLEncoding.EncodeToString(b))
}

// generateConvoyCode creates a short, shareable convoy code.
func generateConvoyCode() string {
	// Generate a 6-character alphanumeric code
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

// AddPlayer adds a player to the convoy.
func (c *Convoy) AddPlayer(id PlayerID, name string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.State != StateWaiting {
		return ErrConvoyStarted
	}
	if len(c.Players) >= c.MaxPlayers {
		return ErrConvoyFull
	}

	// Check for duplicate player
	for _, p := range c.Players {
		if p.ID == id {
			return nil // Already in convoy
		}
	}

	isHost := len(c.Players) == 0
	player := &Player{
		ID:       id,
		Name:     name,
		JoinedAt: time.Now().UTC(),
		IsHost:   isHost,
		IsReady:  false,
	}
	c.Players = append(c.Players, player)

	if isHost {
		c.HostID = id
	}

	return nil
}

// RemovePlayer removes a player from the convoy.
func (c *Convoy) RemovePlayer(id PlayerID) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, p := range c.Players {
		if p.ID == id {
			c.Players = append(c.Players[:i], c.Players[i+1:]...)
			// Reassign host if needed
			if p.IsHost && len(c.Players) > 0 {
				c.Players[0].IsHost = true
				c.HostID = c.Players[0].ID
			}
			return nil
		}
	}
	return ErrPlayerNotInConvoy
}

// SetPlayerReady marks a player as ready to start.
func (c *Convoy) SetPlayerReady(id PlayerID, ready bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, p := range c.Players {
		if p.ID == id {
			p.IsReady = ready
			return nil
		}
	}
	return ErrPlayerNotInConvoy
}

// AllPlayersReady returns true if all players are ready.
func (c *Convoy) AllPlayersReady() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(c.Players) < 2 {
		return false
	}
	for _, p := range c.Players {
		if !p.IsReady {
			return false
		}
	}
	return true
}

// Start begins the convoy session.
func (c *Convoy) Start() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.State != StateWaiting {
		return ErrConvoyStarted
	}
	if len(c.Players) < 2 {
		return errors.New("need at least 2 players to start")
	}

	c.State = StateStarted
	c.StartedAt = time.Now().UTC()

	// Initialize runs for all players
	for _, p := range c.Players {
		c.Runs = append(c.Runs, &RunData{
			PlayerID:  p.ID,
			Genre:     c.Genre,
			StartedAt: c.StartedAt,
		})
	}

	return nil
}

// RecordRunResult records a player's completed run.
func (c *Convoy) RecordRunResult(playerID PlayerID, result *RunData) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.updatePlayerRun(playerID, result)
	c.markPlayerFinished(playerID)
	c.checkConvoyCompletion()

	return nil
}

// updatePlayerRun copies result data into the player's run record.
func (c *Convoy) updatePlayerRun(playerID PlayerID, result *RunData) {
	for _, run := range c.Runs {
		if run.PlayerID == playerID {
			run.FinishedAt = time.Now().UTC()
			run.IsVictory = result.IsVictory
			run.Score = result.Score
			run.DaysTraveled = result.DaysTraveled
			run.Survivors = result.Survivors
			run.EventsSeen = result.EventsSeen
			return
		}
	}
}

// markPlayerFinished sets HasFinished flag for a player.
func (c *Convoy) markPlayerFinished(playerID PlayerID) {
	for _, p := range c.Players {
		if p.ID == playerID {
			p.HasFinished = true
			return
		}
	}
}

// checkConvoyCompletion marks convoy complete if all players finished.
func (c *Convoy) checkConvoyCompletion() {
	for _, p := range c.Players {
		if !p.HasFinished {
			return
		}
	}
	c.State = StateCompleted
}

// GetPlayer returns a player by ID.
func (c *Convoy) GetPlayer(id PlayerID) *Player {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, p := range c.Players {
		if p.ID == id {
			return p
		}
	}
	return nil
}

// GetRun returns a player's run data.
func (c *Convoy) GetRun(playerID PlayerID) *RunData {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, run := range c.Runs {
		if run.PlayerID == playerID {
			return run
		}
	}
	return nil
}

// PlayerCount returns the number of players in the convoy.
func (c *Convoy) PlayerCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.Players)
}

// ReadyCount returns the number of ready players.
func (c *Convoy) ReadyCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	count := 0
	for _, p := range c.Players {
		if p.IsReady {
			count++
		}
	}
	return count
}

// FinishedCount returns the number of players who have finished.
func (c *Convoy) FinishedCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	count := 0
	for _, p := range c.Players {
		if p.HasFinished {
			count++
		}
	}
	return count
}

// IsComplete returns true if the convoy has completed.
func (c *Convoy) IsComplete() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.State == StateCompleted
}

// Marshal serializes the convoy to JSON.
func (c *Convoy) Marshal() ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return json.Marshal(c)
}

// UnmarshalConvoy deserializes a convoy from JSON.
func UnmarshalConvoy(data []byte) (*Convoy, error) {
	var convoy Convoy
	if err := json.Unmarshal(data, &convoy); err != nil {
		return nil, err
	}
	return &convoy, nil
}
