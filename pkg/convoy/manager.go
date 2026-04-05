package convoy

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/opd-ai/voyage/pkg/engine"
)

// Manager errors.
var (
	ErrManagerClosed = errors.New("convoy manager is closed")
)

// Manager manages multiple convoy sessions.
type Manager struct {
	mu       sync.RWMutex
	convoys  map[ConvoyID]*Convoy
	byCode   map[string]*Convoy
	byPlayer map[PlayerID]*Convoy
	closed   bool
}

// NewManager creates a new convoy manager.
func NewManager() *Manager {
	return &Manager{
		convoys:  make(map[ConvoyID]*Convoy),
		byCode:   make(map[string]*Convoy),
		byPlayer: make(map[PlayerID]*Convoy),
	}
}

// CreateConvoy creates a new convoy and registers it.
func (m *Manager) CreateConvoy(seed int64, genre string, maxPlayers int, hostID PlayerID, hostName string) (*Convoy, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil, ErrManagerClosed
	}

	convoy := NewConvoy(seed, engine.GenreID(genre), maxPlayers)

	// Add the host as the first player
	if err := convoy.AddPlayer(hostID, hostName); err != nil {
		return nil, err
	}

	m.convoys[convoy.ID] = convoy
	m.byCode[convoy.Code] = convoy
	m.byPlayer[hostID] = convoy

	return convoy, nil
}

// lookupConvoy is a helper that performs a locked lookup with closed-check.
func (m *Manager) lookupConvoy(lookup func() (*Convoy, bool)) (*Convoy, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return nil, ErrManagerClosed
	}

	convoy, ok := lookup()
	if !ok {
		return nil, ErrConvoyNotFound
	}
	return convoy, nil
}

// GetConvoy returns a convoy by ID.
func (m *Manager) GetConvoy(id ConvoyID) (*Convoy, error) {
	return m.lookupConvoy(func() (*Convoy, bool) {
		c, ok := m.convoys[id]
		return c, ok
	})
}

// GetConvoyByCode returns a convoy by its shareable code.
func (m *Manager) GetConvoyByCode(code string) (*Convoy, error) {
	return m.lookupConvoy(func() (*Convoy, bool) {
		c, ok := m.byCode[code]
		return c, ok
	})
}

// GetPlayerConvoy returns the convoy a player is currently in.
func (m *Manager) GetPlayerConvoy(playerID PlayerID) (*Convoy, error) {
	return m.lookupConvoy(func() (*Convoy, bool) {
		c, ok := m.byPlayer[playerID]
		return c, ok
	})
}

// JoinConvoy adds a player to a convoy by code.
func (m *Manager) JoinConvoy(code string, playerID PlayerID, playerName string) (*Convoy, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil, ErrManagerClosed
	}

	convoy, ok := m.byCode[code]
	if !ok {
		return nil, ErrConvoyNotFound
	}

	if err := convoy.AddPlayer(playerID, playerName); err != nil {
		return nil, err
	}

	m.byPlayer[playerID] = convoy
	return convoy, nil
}

// LeaveConvoy removes a player from their convoy.
func (m *Manager) LeaveConvoy(playerID PlayerID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return ErrManagerClosed
	}

	convoy, ok := m.byPlayer[playerID]
	if !ok {
		return ErrPlayerNotInConvoy
	}

	if err := convoy.RemovePlayer(playerID); err != nil {
		return err
	}

	delete(m.byPlayer, playerID)

	// Clean up empty convoys
	if convoy.PlayerCount() == 0 {
		delete(m.convoys, convoy.ID)
		delete(m.byCode, convoy.Code)
	}

	return nil
}

// StartConvoy starts a convoy (only host can do this).
func (m *Manager) StartConvoy(convoyID ConvoyID, playerID PlayerID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return ErrManagerClosed
	}

	convoy, ok := m.convoys[convoyID]
	if !ok {
		return ErrConvoyNotFound
	}

	// Verify player is host
	if convoy.HostID != playerID {
		return errors.New("only host can start the convoy")
	}

	return convoy.Start()
}

// RecordResult records a player's run result.
func (m *Manager) RecordResult(playerID PlayerID, result *RunData) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return ErrManagerClosed
	}

	convoy, ok := m.byPlayer[playerID]
	if !ok {
		return ErrPlayerNotInConvoy
	}

	return convoy.RecordRunResult(playerID, result)
}

// ListActive returns all active (non-completed) convoys.
func (m *Manager) ListActive() []*Convoy {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*Convoy, 0)
	for _, c := range m.convoys {
		if c.State != StateCompleted {
			result = append(result, c)
		}
	}
	return result
}

// ListWaiting returns convoys still waiting for players.
func (m *Manager) ListWaiting() []*Convoy {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*Convoy, 0)
	for _, c := range m.convoys {
		if c.State == StateWaiting {
			result = append(result, c)
		}
	}
	return result
}

// CleanupOld removes convoys older than the given duration.
func (m *Manager) CleanupOld(maxAge time.Duration) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	cutoff := time.Now().Add(-maxAge)
	removed := 0

	for id, convoy := range m.convoys {
		if convoy.CreatedAt.Before(cutoff) && convoy.State == StateCompleted {
			// Remove player mappings
			for _, p := range convoy.Players {
				delete(m.byPlayer, p.ID)
			}
			delete(m.byCode, convoy.Code)
			delete(m.convoys, id)
			removed++
		}
	}

	return removed
}

// Count returns the total number of convoys.
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.convoys)
}

// Close shuts down the manager.
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.closed = true
	m.convoys = nil
	m.byCode = nil
	m.byPlayer = nil

	return nil
}

// ExportConvoy exports a convoy to JSON.
func (m *Manager) ExportConvoy(id ConvoyID) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	convoy, ok := m.convoys[id]
	if !ok {
		return nil, ErrConvoyNotFound
	}

	return convoy.Marshal()
}

// ImportConvoy imports a convoy from JSON.
func (m *Manager) ImportConvoy(data []byte) (*Convoy, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil, ErrManagerClosed
	}

	convoy, err := UnmarshalConvoy(data)
	if err != nil {
		return nil, err
	}

	m.convoys[convoy.ID] = convoy
	m.byCode[convoy.Code] = convoy

	for _, p := range convoy.Players {
		m.byPlayer[p.ID] = convoy
	}

	return convoy, nil
}

// Message represents a convoy message for synchronization.
type Message struct {
	Type      string          `json:"type"`
	ConvoyID  ConvoyID        `json:"convoyId,omitempty"`
	PlayerID  PlayerID        `json:"playerId,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

// Message types for convoy synchronization.
const (
	MsgTypeJoin     = "join"
	MsgTypeLeave    = "leave"
	MsgTypeReady    = "ready"
	MsgTypeStart    = "start"
	MsgTypeProgress = "progress"
	MsgTypeFinish   = "finish"
	MsgTypeSync     = "sync"
)

// NewMessage creates a new convoy message.
func NewMessage(msgType string, convoyID ConvoyID, playerID PlayerID, payload interface{}) (*Message, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Message{
		Type:      msgType,
		ConvoyID:  convoyID,
		PlayerID:  playerID,
		Timestamp: time.Now().UTC(),
		Payload:   payloadBytes,
	}, nil
}

// Marshal serializes a message to JSON.
func (msg *Message) Marshal() ([]byte, error) {
	return json.Marshal(msg)
}

// UnmarshalMessage deserializes a message from JSON.
func UnmarshalMessage(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
