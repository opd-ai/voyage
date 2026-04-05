package convoy

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// Network errors for high-latency scenarios.
var (
	ErrTimeout        = errors.New("operation timed out")
	ErrRetryExhausted = errors.New("retry attempts exhausted")
	ErrDisconnected   = errors.New("connection lost")
)

// LatencyConfig configures high-latency tolerant behavior.
type LatencyConfig struct {
	// MaxLatency is the maximum acceptable latency for operations.
	MaxLatency time.Duration

	// RetryCount is the number of retries for failed operations.
	RetryCount int

	// RetryBackoff is the initial backoff duration between retries.
	RetryBackoff time.Duration

	// MaxRetryBackoff is the maximum backoff duration.
	MaxRetryBackoff time.Duration

	// HeartbeatInterval is how often to send keep-alive messages.
	HeartbeatInterval time.Duration

	// DisconnectTimeout is how long to wait before considering disconnected.
	DisconnectTimeout time.Duration
}

// DefaultLatencyConfig returns a configuration suitable for high-latency
// networks like Tor (200-5000ms latency).
func DefaultLatencyConfig() *LatencyConfig {
	return &LatencyConfig{
		MaxLatency:        10 * time.Second,
		RetryCount:        3,
		RetryBackoff:      2 * time.Second,
		MaxRetryBackoff:   30 * time.Second,
		HeartbeatInterval: 30 * time.Second,
		DisconnectTimeout: 120 * time.Second,
	}
}

// LowLatencyConfig returns a configuration for low-latency networks.
func LowLatencyConfig() *LatencyConfig {
	return &LatencyConfig{
		MaxLatency:        2 * time.Second,
		RetryCount:        3,
		RetryBackoff:      500 * time.Millisecond,
		MaxRetryBackoff:   5 * time.Second,
		HeartbeatInterval: 10 * time.Second,
		DisconnectTimeout: 30 * time.Second,
	}
}

// MessageBuffer provides reliable message delivery with ordering guarantees.
type MessageBuffer struct {
	mu       sync.Mutex
	messages []*BufferedMessage
	nextSeq  int64
	ackSeq   int64
	config   *LatencyConfig
}

// BufferedMessage wraps a message with delivery metadata.
type BufferedMessage struct {
	Sequence  int64     `json:"seq"`
	Message   *Message  `json:"msg"`
	SentAt    time.Time `json:"sentAt"`
	Retries   int       `json:"retries"`
	Acked     bool      `json:"acked"`
	NextRetry time.Time `json:"nextRetry"`
}

// NewMessageBuffer creates a new message buffer with the given config.
func NewMessageBuffer(config *LatencyConfig) *MessageBuffer {
	if config == nil {
		config = DefaultLatencyConfig()
	}
	return &MessageBuffer{
		messages: make([]*BufferedMessage, 0),
		nextSeq:  1,
		ackSeq:   0,
		config:   config,
	}
}

// Enqueue adds a message to the buffer for reliable delivery.
func (mb *MessageBuffer) Enqueue(msg *Message) *BufferedMessage {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	buffered := &BufferedMessage{
		Sequence: mb.nextSeq,
		Message:  msg,
		SentAt:   time.Now(),
		Retries:  0,
		Acked:    false,
	}
	mb.nextSeq++
	mb.messages = append(mb.messages, buffered)

	return buffered
}

// Acknowledge marks a message as successfully delivered.
func (mb *MessageBuffer) Acknowledge(seq int64) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	for _, m := range mb.messages {
		if m.Sequence == seq {
			m.Acked = true
		}
	}

	// Update ack sequence
	if seq > mb.ackSeq {
		mb.ackSeq = seq
	}

	// Cleanup old acked messages
	mb.cleanup()
}

// GetPending returns messages that need to be sent or retried.
func (mb *MessageBuffer) GetPending() []*BufferedMessage {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	now := time.Now()
	pending := make([]*BufferedMessage, 0)

	for _, m := range mb.messages {
		if m.Acked {
			continue
		}
		// First send or ready for retry
		if m.Retries == 0 || now.After(m.NextRetry) {
			pending = append(pending, m)
		}
	}

	return pending
}

// MarkRetry marks a message for retry with exponential backoff.
func (mb *MessageBuffer) MarkRetry(seq int64) error {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	for _, m := range mb.messages {
		if m.Sequence == seq && !m.Acked {
			m.Retries++
			if m.Retries > mb.config.RetryCount {
				return ErrRetryExhausted
			}
			// Exponential backoff
			backoff := mb.config.RetryBackoff * time.Duration(1<<uint(m.Retries-1))
			if backoff > mb.config.MaxRetryBackoff {
				backoff = mb.config.MaxRetryBackoff
			}
			m.NextRetry = time.Now().Add(backoff)
			return nil
		}
	}
	return nil
}

// cleanup removes old acknowledged messages.
func (mb *MessageBuffer) cleanup() {
	var filtered []*BufferedMessage
	for _, m := range mb.messages {
		// Keep unacked messages and recently acked ones
		if !m.Acked || time.Since(m.SentAt) < time.Minute {
			filtered = append(filtered, m)
		}
	}
	mb.messages = filtered
}

// UnackedCount returns the number of unacknowledged messages.
func (mb *MessageBuffer) UnackedCount() int {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	count := 0
	for _, m := range mb.messages {
		if !m.Acked {
			count++
		}
	}
	return count
}

// PeerState tracks the connection state of a peer.
type PeerState struct {
	PlayerID     PlayerID      `json:"playerId"`
	LastSeen     time.Time     `json:"lastSeen"`
	Latency      time.Duration `json:"latency"`
	IsConnected  bool          `json:"isConnected"`
	MessagesSent int64         `json:"messagesSent"`
	MessagesRecv int64         `json:"messagesRecv"`
}

// ConnectionTracker monitors peer connections with high-latency tolerance.
type ConnectionTracker struct {
	mu     sync.RWMutex
	peers  map[PlayerID]*PeerState
	config *LatencyConfig
}

// NewConnectionTracker creates a new connection tracker.
func NewConnectionTracker(config *LatencyConfig) *ConnectionTracker {
	if config == nil {
		config = DefaultLatencyConfig()
	}
	return &ConnectionTracker{
		peers:  make(map[PlayerID]*PeerState),
		config: config,
	}
}

// RegisterPeer adds a peer to tracking.
func (ct *ConnectionTracker) RegisterPeer(playerID PlayerID) {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	ct.peers[playerID] = &PeerState{
		PlayerID:    playerID,
		LastSeen:    time.Now(),
		IsConnected: true,
	}
}

// RemovePeer removes a peer from tracking.
func (ct *ConnectionTracker) RemovePeer(playerID PlayerID) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	delete(ct.peers, playerID)
}

// RecordActivity updates the last seen time for a peer.
func (ct *ConnectionTracker) RecordActivity(playerID PlayerID) {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	if peer, ok := ct.peers[playerID]; ok {
		peer.LastSeen = time.Now()
		peer.IsConnected = true
		peer.MessagesRecv++
	}
}

// RecordLatency records a latency measurement for a peer.
func (ct *ConnectionTracker) RecordLatency(playerID PlayerID, latency time.Duration) {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	if peer, ok := ct.peers[playerID]; ok {
		// Simple moving average
		if peer.Latency == 0 {
			peer.Latency = latency
		} else {
			peer.Latency = (peer.Latency*3 + latency) / 4
		}
	}
}

// GetPeer returns the state of a specific peer.
func (ct *ConnectionTracker) GetPeer(playerID PlayerID) *PeerState {
	ct.mu.RLock()
	defer ct.mu.RUnlock()
	return ct.peers[playerID]
}

// GetConnectedPeers returns all connected peers.
func (ct *ConnectionTracker) GetConnectedPeers() []*PeerState {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	var result []*PeerState
	for _, peer := range ct.peers {
		if peer.IsConnected {
			result = append(result, peer)
		}
	}
	return result
}

// CheckTimeouts checks for timed-out peers.
func (ct *ConnectionTracker) CheckTimeouts() []PlayerID {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	var timedOut []PlayerID
	now := time.Now()

	for id, peer := range ct.peers {
		if peer.IsConnected && now.Sub(peer.LastSeen) > ct.config.DisconnectTimeout {
			peer.IsConnected = false
			timedOut = append(timedOut, id)
		}
	}

	return timedOut
}

// AverageLatency returns the average latency across all connected peers.
func (ct *ConnectionTracker) AverageLatency() time.Duration {
	ct.mu.RLock()
	defer ct.mu.RUnlock()

	var total time.Duration
	var count int

	for _, peer := range ct.peers {
		if peer.IsConnected && peer.Latency > 0 {
			total += peer.Latency
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return total / time.Duration(count)
}

// SyncState represents the synchronization state for high-latency scenarios.
type SyncState struct {
	ConvoyID     ConvoyID                      `json:"convoyId"`
	LastSync     time.Time                     `json:"lastSync"`
	StateHash    string                        `json:"stateHash"`
	Version      int64                         `json:"version"`
	PlayerStates map[PlayerID]*PlayerSyncState `json:"playerStates"`
}

// PlayerSyncState represents a player's synced state.
type PlayerSyncState struct {
	PlayerID   PlayerID  `json:"playerId"`
	Turn       int       `json:"turn"`
	Position   [2]int    `json:"position"`
	IsFinished bool      `json:"isFinished"`
	LastUpdate time.Time `json:"lastUpdate"`
}

// NewSyncState creates a new sync state for a convoy.
func NewSyncState(convoyID ConvoyID) *SyncState {
	return &SyncState{
		ConvoyID:     convoyID,
		LastSync:     time.Now(),
		Version:      1,
		PlayerStates: make(map[PlayerID]*PlayerSyncState),
	}
}

// UpdatePlayer updates a player's sync state.
func (ss *SyncState) UpdatePlayer(playerID PlayerID, turn, x, y int, finished bool) {
	ss.PlayerStates[playerID] = &PlayerSyncState{
		PlayerID:   playerID,
		Turn:       turn,
		Position:   [2]int{x, y},
		IsFinished: finished,
		LastUpdate: time.Now(),
	}
	ss.Version++
	ss.LastSync = time.Now()
}

// Marshal serializes the sync state.
func (ss *SyncState) Marshal() ([]byte, error) {
	return json.Marshal(ss)
}

// UnmarshalSyncState deserializes a sync state.
func UnmarshalSyncState(data []byte) (*SyncState, error) {
	var ss SyncState
	if err := json.Unmarshal(data, &ss); err != nil {
		return nil, err
	}
	return &ss, nil
}

// IsStale checks if a player's state is stale given the latency config.
func (ss *SyncState) IsStale(playerID PlayerID, config *LatencyConfig) bool {
	ps, ok := ss.PlayerStates[playerID]
	if !ok {
		return true
	}
	return time.Since(ps.LastUpdate) > config.DisconnectTimeout
}
