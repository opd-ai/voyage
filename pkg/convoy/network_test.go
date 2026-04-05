package convoy

import (
	"testing"
	"time"
)

func TestDefaultLatencyConfig(t *testing.T) {
	config := DefaultLatencyConfig()

	if config.MaxLatency != 10*time.Second {
		t.Errorf("expected max latency 10s, got %v", config.MaxLatency)
	}
	if config.RetryCount != 3 {
		t.Errorf("expected retry count 3, got %d", config.RetryCount)
	}
	if config.DisconnectTimeout != 120*time.Second {
		t.Errorf("expected disconnect timeout 120s, got %v", config.DisconnectTimeout)
	}
}

func TestLowLatencyConfig(t *testing.T) {
	config := LowLatencyConfig()

	if config.MaxLatency != 2*time.Second {
		t.Errorf("expected max latency 2s, got %v", config.MaxLatency)
	}
	if config.DisconnectTimeout != 30*time.Second {
		t.Errorf("expected disconnect timeout 30s, got %v", config.DisconnectTimeout)
	}
}

func TestMessageBufferEnqueue(t *testing.T) {
	mb := NewMessageBuffer(nil)

	msg := &Message{Type: "test", ConvoyID: "c1"}
	buffered := mb.Enqueue(msg)

	if buffered.Sequence != 1 {
		t.Errorf("expected sequence 1, got %d", buffered.Sequence)
	}
	if buffered.Acked {
		t.Error("message should not be acked initially")
	}

	msg2 := &Message{Type: "test2", ConvoyID: "c1"}
	buffered2 := mb.Enqueue(msg2)

	if buffered2.Sequence != 2 {
		t.Errorf("expected sequence 2, got %d", buffered2.Sequence)
	}
}

func TestMessageBufferAcknowledge(t *testing.T) {
	mb := NewMessageBuffer(nil)

	msg := &Message{Type: "test"}
	buffered := mb.Enqueue(msg)

	if mb.UnackedCount() != 1 {
		t.Errorf("expected 1 unacked, got %d", mb.UnackedCount())
	}

	mb.Acknowledge(buffered.Sequence)

	// After acknowledgment, message is acked but not yet cleaned up
	// UnackedCount should decrease
	pending := mb.GetPending()
	if len(pending) != 0 {
		t.Errorf("expected 0 pending after ack, got %d", len(pending))
	}
}

func TestMessageBufferGetPending(t *testing.T) {
	mb := NewMessageBuffer(nil)

	msg1 := &Message{Type: "test1"}
	msg2 := &Message{Type: "test2"}
	mb.Enqueue(msg1)
	mb.Enqueue(msg2)

	pending := mb.GetPending()
	if len(pending) != 2 {
		t.Errorf("expected 2 pending, got %d", len(pending))
	}

	// Ack first message
	mb.Acknowledge(1)
	pending = mb.GetPending()
	if len(pending) != 1 {
		t.Errorf("expected 1 pending after ack, got %d", len(pending))
	}
}

func TestMessageBufferRetry(t *testing.T) {
	config := &LatencyConfig{
		RetryCount:      2,
		RetryBackoff:    100 * time.Millisecond,
		MaxRetryBackoff: 500 * time.Millisecond,
	}
	mb := NewMessageBuffer(config)

	msg := &Message{Type: "test"}
	buffered := mb.Enqueue(msg)

	// First retry
	err := mb.MarkRetry(buffered.Sequence)
	if err != nil {
		t.Fatalf("MarkRetry() error = %v", err)
	}

	// Check backoff is set
	if buffered.NextRetry.IsZero() {
		t.Error("NextRetry should be set after retry")
	}
	if buffered.Retries != 1 {
		t.Errorf("expected 1 retry, got %d", buffered.Retries)
	}

	// Second retry
	err = mb.MarkRetry(buffered.Sequence)
	if err != nil {
		t.Fatalf("MarkRetry() second error = %v", err)
	}

	// Third retry should exhaust
	err = mb.MarkRetry(buffered.Sequence)
	if err != ErrRetryExhausted {
		t.Errorf("expected ErrRetryExhausted, got %v", err)
	}
}

func TestConnectionTrackerRegisterPeer(t *testing.T) {
	ct := NewConnectionTracker(nil)

	ct.RegisterPeer("player1")

	peer := ct.GetPeer("player1")
	if peer == nil {
		t.Fatal("expected peer to exist")
	}
	if !peer.IsConnected {
		t.Error("peer should be connected")
	}
	if peer.PlayerID != "player1" {
		t.Errorf("expected player ID player1, got %s", peer.PlayerID)
	}
}

func TestConnectionTrackerRecordActivity(t *testing.T) {
	ct := NewConnectionTracker(nil)
	ct.RegisterPeer("player1")

	initialSeen := ct.GetPeer("player1").LastSeen

	time.Sleep(10 * time.Millisecond)
	ct.RecordActivity("player1")

	peer := ct.GetPeer("player1")
	if !peer.LastSeen.After(initialSeen) {
		t.Error("LastSeen should be updated")
	}
	if peer.MessagesRecv != 1 {
		t.Errorf("expected 1 message received, got %d", peer.MessagesRecv)
	}
}

func TestConnectionTrackerRecordLatency(t *testing.T) {
	ct := NewConnectionTracker(nil)
	ct.RegisterPeer("player1")

	ct.RecordLatency("player1", 100*time.Millisecond)
	peer := ct.GetPeer("player1")
	if peer.Latency != 100*time.Millisecond {
		t.Errorf("expected latency 100ms, got %v", peer.Latency)
	}

	// Moving average
	ct.RecordLatency("player1", 200*time.Millisecond)
	peer = ct.GetPeer("player1")
	// (100*3 + 200) / 4 = 125
	if peer.Latency != 125*time.Millisecond {
		t.Errorf("expected latency 125ms (moving average), got %v", peer.Latency)
	}
}

func TestConnectionTrackerGetConnectedPeers(t *testing.T) {
	ct := NewConnectionTracker(nil)
	ct.RegisterPeer("player1")
	ct.RegisterPeer("player2")

	peers := ct.GetConnectedPeers()
	if len(peers) != 2 {
		t.Errorf("expected 2 connected peers, got %d", len(peers))
	}
}

func TestConnectionTrackerCheckTimeouts(t *testing.T) {
	config := &LatencyConfig{
		DisconnectTimeout: 50 * time.Millisecond,
	}
	ct := NewConnectionTracker(config)
	ct.RegisterPeer("player1")

	// No timeouts yet
	timedOut := ct.CheckTimeouts()
	if len(timedOut) != 0 {
		t.Errorf("expected 0 timeouts, got %d", len(timedOut))
	}

	// Wait for timeout
	time.Sleep(60 * time.Millisecond)

	timedOut = ct.CheckTimeouts()
	if len(timedOut) != 1 {
		t.Errorf("expected 1 timeout, got %d", len(timedOut))
	}
	if timedOut[0] != "player1" {
		t.Errorf("expected player1 to timeout, got %s", timedOut[0])
	}

	// Peer should be marked disconnected
	peer := ct.GetPeer("player1")
	if peer.IsConnected {
		t.Error("peer should be marked disconnected")
	}
}

func TestConnectionTrackerAverageLatency(t *testing.T) {
	ct := NewConnectionTracker(nil)
	ct.RegisterPeer("player1")
	ct.RegisterPeer("player2")

	ct.RecordLatency("player1", 100*time.Millisecond)
	ct.RecordLatency("player2", 200*time.Millisecond)

	avg := ct.AverageLatency()
	if avg != 150*time.Millisecond {
		t.Errorf("expected average latency 150ms, got %v", avg)
	}
}

func TestConnectionTrackerAverageLatencyNoConnections(t *testing.T) {
	ct := NewConnectionTracker(nil)

	avg := ct.AverageLatency()
	if avg != 0 {
		t.Errorf("expected 0 latency for no connections, got %v", avg)
	}
}

func TestConnectionTrackerRemovePeer(t *testing.T) {
	ct := NewConnectionTracker(nil)
	ct.RegisterPeer("player1")
	ct.RemovePeer("player1")

	peer := ct.GetPeer("player1")
	if peer != nil {
		t.Error("expected peer to be removed")
	}
}

func TestSyncStateUpdatePlayer(t *testing.T) {
	ss := NewSyncState("convoy1")

	ss.UpdatePlayer("player1", 10, 5, 3, false)

	ps := ss.PlayerStates["player1"]
	if ps == nil {
		t.Fatal("expected player state to exist")
	}
	if ps.Turn != 10 {
		t.Errorf("expected turn 10, got %d", ps.Turn)
	}
	if ps.Position[0] != 5 || ps.Position[1] != 3 {
		t.Errorf("expected position [5, 3], got %v", ps.Position)
	}
	if ps.IsFinished {
		t.Error("player should not be finished")
	}

	if ss.Version != 2 {
		t.Errorf("expected version 2, got %d", ss.Version)
	}
}

func TestSyncStateIsStale(t *testing.T) {
	config := &LatencyConfig{
		DisconnectTimeout: 50 * time.Millisecond,
	}
	ss := NewSyncState("convoy1")

	// Non-existent player is stale
	if !ss.IsStale("player1", config) {
		t.Error("non-existent player should be stale")
	}

	ss.UpdatePlayer("player1", 1, 0, 0, false)

	// Recently updated is not stale
	if ss.IsStale("player1", config) {
		t.Error("recently updated player should not be stale")
	}

	// Wait for staleness
	time.Sleep(60 * time.Millisecond)
	if !ss.IsStale("player1", config) {
		t.Error("old player state should be stale")
	}
}

func TestSyncStateMarshalUnmarshal(t *testing.T) {
	ss := NewSyncState("convoy1")
	ss.UpdatePlayer("player1", 10, 5, 3, true)

	data, err := ss.Marshal()
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	restored, err := UnmarshalSyncState(data)
	if err != nil {
		t.Fatalf("UnmarshalSyncState() error = %v", err)
	}

	if restored.ConvoyID != ss.ConvoyID {
		t.Errorf("convoy ID mismatch")
	}
	if restored.Version != ss.Version {
		t.Errorf("version mismatch: got %d, want %d", restored.Version, ss.Version)
	}

	ps := restored.PlayerStates["player1"]
	if ps == nil {
		t.Fatal("expected player state to exist")
	}
	if ps.Turn != 10 {
		t.Errorf("expected turn 10, got %d", ps.Turn)
	}
}

func TestMessageBufferNilConfig(t *testing.T) {
	mb := NewMessageBuffer(nil)

	// Should use default config
	msg := &Message{Type: "test"}
	buffered := mb.Enqueue(msg)

	// Should work without panicking
	_ = mb.MarkRetry(buffered.Sequence)
	mb.Acknowledge(buffered.Sequence)
}

func TestConnectionTrackerNilConfig(t *testing.T) {
	ct := NewConnectionTracker(nil)

	// Should use default config
	ct.RegisterPeer("player1")
	peer := ct.GetPeer("player1")
	if peer == nil {
		t.Error("peer should exist")
	}
}
