package convoy

import (
	"testing"
	"time"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewConvoy(t *testing.T) {
	c := NewConvoy(12345, engine.GenreFantasy, 4)

	if c.Seed != 12345 {
		t.Errorf("expected seed 12345, got %d", c.Seed)
	}
	if c.Genre != engine.GenreFantasy {
		t.Errorf("expected genre fantasy, got %s", c.Genre)
	}
	if c.MaxPlayers != 4 {
		t.Errorf("expected max players 4, got %d", c.MaxPlayers)
	}
	if c.State != StateWaiting {
		t.Errorf("expected state waiting, got %d", c.State)
	}
	if c.ID == "" {
		t.Error("convoy ID should not be empty")
	}
	if c.Code == "" {
		t.Error("convoy code should not be empty")
	}
	if len(c.Code) != 6 {
		t.Errorf("expected code length 6, got %d", len(c.Code))
	}
}

func TestConvoyMaxPlayersClamping(t *testing.T) {
	// Test minimum
	c1 := NewConvoy(1, engine.GenreScifi, 1)
	if c1.MaxPlayers != 2 {
		t.Errorf("expected max players clamped to 2, got %d", c1.MaxPlayers)
	}

	// Test maximum
	c2 := NewConvoy(1, engine.GenreScifi, 20)
	if c2.MaxPlayers != 10 {
		t.Errorf("expected max players clamped to 10, got %d", c2.MaxPlayers)
	}
}

func TestConvoyAddPlayer(t *testing.T) {
	c := NewConvoy(12345, engine.GenreHorror, 4)

	err := c.AddPlayer("p1", "Player One")
	if err != nil {
		t.Fatalf("AddPlayer() error = %v", err)
	}

	if c.PlayerCount() != 1 {
		t.Errorf("expected 1 player, got %d", c.PlayerCount())
	}

	// First player should be host
	p := c.GetPlayer("p1")
	if !p.IsHost {
		t.Error("first player should be host")
	}
	if c.HostID != "p1" {
		t.Errorf("expected host ID p1, got %s", c.HostID)
	}

	// Add second player
	err = c.AddPlayer("p2", "Player Two")
	if err != nil {
		t.Fatalf("AddPlayer() second player error = %v", err)
	}

	if c.PlayerCount() != 2 {
		t.Errorf("expected 2 players, got %d", c.PlayerCount())
	}

	// Second player should not be host
	p2 := c.GetPlayer("p2")
	if p2.IsHost {
		t.Error("second player should not be host")
	}
}

func TestConvoyAddPlayerDuplicate(t *testing.T) {
	c := NewConvoy(12345, engine.GenreFantasy, 4)

	_ = c.AddPlayer("p1", "Player One")
	err := c.AddPlayer("p1", "Player One Again")
	if err != nil {
		t.Error("adding duplicate player should not error")
	}
	if c.PlayerCount() != 1 {
		t.Errorf("expected 1 player after duplicate, got %d", c.PlayerCount())
	}
}

func TestConvoyAddPlayerFull(t *testing.T) {
	c := NewConvoy(12345, engine.GenreFantasy, 2)

	_ = c.AddPlayer("p1", "Player One")
	_ = c.AddPlayer("p2", "Player Two")
	err := c.AddPlayer("p3", "Player Three")

	if err != ErrConvoyFull {
		t.Errorf("expected ErrConvoyFull, got %v", err)
	}
}

func TestConvoyRemovePlayer(t *testing.T) {
	c := NewConvoy(12345, engine.GenreFantasy, 4)

	_ = c.AddPlayer("p1", "Player One")
	_ = c.AddPlayer("p2", "Player Two")

	err := c.RemovePlayer("p1")
	if err != nil {
		t.Fatalf("RemovePlayer() error = %v", err)
	}

	if c.PlayerCount() != 1 {
		t.Errorf("expected 1 player after removal, got %d", c.PlayerCount())
	}

	// p2 should now be host
	p2 := c.GetPlayer("p2")
	if !p2.IsHost {
		t.Error("p2 should be host after p1 removal")
	}
}

func TestConvoyRemovePlayerNotInConvoy(t *testing.T) {
	c := NewConvoy(12345, engine.GenreFantasy, 4)

	err := c.RemovePlayer("nonexistent")
	if err != ErrPlayerNotInConvoy {
		t.Errorf("expected ErrPlayerNotInConvoy, got %v", err)
	}
}

func TestConvoySetPlayerReady(t *testing.T) {
	c := NewConvoy(12345, engine.GenreFantasy, 4)
	_ = c.AddPlayer("p1", "Player One")

	err := c.SetPlayerReady("p1", true)
	if err != nil {
		t.Fatalf("SetPlayerReady() error = %v", err)
	}

	p := c.GetPlayer("p1")
	if !p.IsReady {
		t.Error("player should be ready")
	}

	err = c.SetPlayerReady("p1", false)
	if err != nil {
		t.Fatalf("SetPlayerReady() unset error = %v", err)
	}

	p = c.GetPlayer("p1")
	if p.IsReady {
		t.Error("player should not be ready")
	}
}

func TestConvoyAllPlayersReady(t *testing.T) {
	c := NewConvoy(12345, engine.GenreFantasy, 4)

	// Need at least 2 players
	_ = c.AddPlayer("p1", "Player One")
	if c.AllPlayersReady() {
		t.Error("should not be ready with only 1 player")
	}

	_ = c.AddPlayer("p2", "Player Two")
	if c.AllPlayersReady() {
		t.Error("should not be ready when players aren't ready")
	}

	_ = c.SetPlayerReady("p1", true)
	if c.AllPlayersReady() {
		t.Error("should not be ready when not all players are ready")
	}

	_ = c.SetPlayerReady("p2", true)
	if !c.AllPlayersReady() {
		t.Error("should be ready when all players are ready")
	}
}

func TestConvoyStart(t *testing.T) {
	c := NewConvoy(12345, engine.GenreFantasy, 4)

	// Can't start with less than 2 players
	_ = c.AddPlayer("p1", "Player One")
	err := c.Start()
	if err == nil {
		t.Error("should not start with less than 2 players")
	}

	_ = c.AddPlayer("p2", "Player Two")
	err = c.Start()
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if c.State != StateStarted {
		t.Errorf("expected state started, got %d", c.State)
	}
	if c.StartedAt.IsZero() {
		t.Error("StartedAt should be set")
	}
	if len(c.Runs) != 2 {
		t.Errorf("expected 2 runs, got %d", len(c.Runs))
	}

	// Can't start twice
	err = c.Start()
	if err != ErrConvoyStarted {
		t.Errorf("expected ErrConvoyStarted, got %v", err)
	}
}

func TestConvoyRecordRunResult(t *testing.T) {
	c := NewConvoy(12345, engine.GenreFantasy, 4)
	_ = c.AddPlayer("p1", "Player One")
	_ = c.AddPlayer("p2", "Player Two")
	_ = c.Start()

	result := &RunData{
		IsVictory:    true,
		Score:        1000,
		DaysTraveled: 30,
		Survivors:    4,
		EventsSeen:   15,
	}

	err := c.RecordRunResult("p1", result)
	if err != nil {
		t.Fatalf("RecordRunResult() error = %v", err)
	}

	p := c.GetPlayer("p1")
	if !p.HasFinished {
		t.Error("player should be marked as finished")
	}

	run := c.GetRun("p1")
	if run.Score != 1000 {
		t.Errorf("expected score 1000, got %d", run.Score)
	}

	// Convoy should still be in progress
	if c.State == StateCompleted {
		t.Error("convoy should not be complete yet")
	}

	// Second player finishes
	_ = c.RecordRunResult("p2", &RunData{
		IsVictory:    false,
		Score:        500,
		DaysTraveled: 20,
		Survivors:    0,
		EventsSeen:   10,
	})

	// Now convoy should be complete
	if c.State != StateCompleted {
		t.Errorf("expected state completed, got %d", c.State)
	}
}

func TestConvoyCounts(t *testing.T) {
	c := NewConvoy(12345, engine.GenreFantasy, 4)
	_ = c.AddPlayer("p1", "Player One")
	_ = c.AddPlayer("p2", "Player Two")

	if c.PlayerCount() != 2 {
		t.Errorf("expected player count 2, got %d", c.PlayerCount())
	}

	if c.ReadyCount() != 0 {
		t.Errorf("expected ready count 0, got %d", c.ReadyCount())
	}

	_ = c.SetPlayerReady("p1", true)
	if c.ReadyCount() != 1 {
		t.Errorf("expected ready count 1, got %d", c.ReadyCount())
	}

	_ = c.Start()

	if c.FinishedCount() != 0 {
		t.Errorf("expected finished count 0, got %d", c.FinishedCount())
	}

	_ = c.RecordRunResult("p1", &RunData{})
	if c.FinishedCount() != 1 {
		t.Errorf("expected finished count 1, got %d", c.FinishedCount())
	}
}

func TestConvoyMarshalUnmarshal(t *testing.T) {
	c := NewConvoy(12345, engine.GenreFantasy, 4)
	_ = c.AddPlayer("p1", "Player One")
	_ = c.AddPlayer("p2", "Player Two")

	data, err := c.Marshal()
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	restored, err := UnmarshalConvoy(data)
	if err != nil {
		t.Fatalf("UnmarshalConvoy() error = %v", err)
	}

	if restored.Seed != c.Seed {
		t.Errorf("seed mismatch: got %d, want %d", restored.Seed, c.Seed)
	}
	if len(restored.Players) != len(c.Players) {
		t.Errorf("player count mismatch: got %d, want %d", len(restored.Players), len(c.Players))
	}
}

func TestConvoyCompare(t *testing.T) {
	c := NewConvoy(12345, engine.GenreFantasy, 4)
	_ = c.AddPlayer("p1", "Player One")
	_ = c.AddPlayer("p2", "Player Two")
	_ = c.Start()

	// P1 wins with high score
	_ = c.RecordRunResult("p1", &RunData{
		IsVictory:    true,
		Score:        1000,
		DaysTraveled: 30,
		Survivors:    4,
		EventsSeen:   15,
	})

	// P2 loses with lower score
	_ = c.RecordRunResult("p2", &RunData{
		IsVictory:    false,
		Score:        500,
		DaysTraveled: 20,
		Survivors:    0,
		EventsSeen:   10,
	})

	result := c.Compare()

	if result.ConvoyID != c.ID {
		t.Errorf("expected convoy ID %s, got %s", c.ID, result.ConvoyID)
	}
	if len(result.Rankings) != 2 {
		t.Fatalf("expected 2 rankings, got %d", len(result.Rankings))
	}

	// First ranking should be P1 (higher score)
	if result.Rankings[0].PlayerID != "p1" {
		t.Errorf("expected first ranking to be p1, got %s", result.Rankings[0].PlayerID)
	}
	if result.Rankings[0].Rank != 1 {
		t.Errorf("expected rank 1, got %d", result.Rankings[0].Rank)
	}

	if result.WinnerID != "p1" {
		t.Errorf("expected winner p1, got %s", result.WinnerID)
	}
	if result.FastestID != "p1" {
		t.Errorf("expected fastest p1 (only victor), got %s", result.FastestID)
	}
}

func TestConvoyGetWinner(t *testing.T) {
	c := NewConvoy(12345, engine.GenreFantasy, 4)
	_ = c.AddPlayer("p1", "Player One")
	_ = c.AddPlayer("p2", "Player Two")
	_ = c.Start()

	// No one has finished
	winner, _ := c.GetWinner()
	if winner != nil {
		t.Error("expected no winner before anyone finishes")
	}

	// P1 loses
	_ = c.RecordRunResult("p1", &RunData{IsVictory: false, Score: 500})

	// P2 wins
	_ = c.RecordRunResult("p2", &RunData{IsVictory: true, Score: 1000})

	winner, run := c.GetWinner()
	if winner == nil {
		t.Fatal("expected a winner")
	}
	if winner.ID != "p2" {
		t.Errorf("expected winner p2, got %s", winner.ID)
	}
	if run.Score != 1000 {
		t.Errorf("expected winning score 1000, got %d", run.Score)
	}
}

func TestConvoyGetVictoryCount(t *testing.T) {
	c := NewConvoy(12345, engine.GenreFantasy, 4)
	_ = c.AddPlayer("p1", "Player One")
	_ = c.AddPlayer("p2", "Player Two")
	_ = c.Start()

	if c.GetVictoryCount() != 0 {
		t.Error("expected 0 victories initially")
	}

	_ = c.RecordRunResult("p1", &RunData{IsVictory: true})
	if c.GetVictoryCount() != 1 {
		t.Errorf("expected 1 victory, got %d", c.GetVictoryCount())
	}

	_ = c.RecordRunResult("p2", &RunData{IsVictory: true})
	if c.GetVictoryCount() != 2 {
		t.Errorf("expected 2 victories, got %d", c.GetVictoryCount())
	}
}

func TestManagerCreateAndGet(t *testing.T) {
	m := NewManager()
	defer m.Close()

	convoy, err := m.CreateConvoy(12345, "fantasy", 4, "host1", "Host Player")
	if err != nil {
		t.Fatalf("CreateConvoy() error = %v", err)
	}

	// Get by ID
	got, err := m.GetConvoy(convoy.ID)
	if err != nil {
		t.Fatalf("GetConvoy() error = %v", err)
	}
	if got.Seed != 12345 {
		t.Errorf("expected seed 12345, got %d", got.Seed)
	}

	// Get by code
	got, err = m.GetConvoyByCode(convoy.Code)
	if err != nil {
		t.Fatalf("GetConvoyByCode() error = %v", err)
	}
	if got.ID != convoy.ID {
		t.Errorf("expected convoy ID %s, got %s", convoy.ID, got.ID)
	}
}

func TestManagerJoinLeave(t *testing.T) {
	m := NewManager()
	defer m.Close()

	convoy, _ := m.CreateConvoy(12345, "fantasy", 4, "host1", "Host Player")

	// Join
	_, err := m.JoinConvoy(convoy.Code, "p2", "Player Two")
	if err != nil {
		t.Fatalf("JoinConvoy() error = %v", err)
	}

	if convoy.PlayerCount() != 2 {
		t.Errorf("expected 2 players, got %d", convoy.PlayerCount())
	}

	// Get player's convoy
	got, err := m.GetPlayerConvoy("p2")
	if err != nil {
		t.Fatalf("GetPlayerConvoy() error = %v", err)
	}
	if got.ID != convoy.ID {
		t.Error("player should be in the convoy")
	}

	// Leave
	err = m.LeaveConvoy("p2")
	if err != nil {
		t.Fatalf("LeaveConvoy() error = %v", err)
	}

	if convoy.PlayerCount() != 1 {
		t.Errorf("expected 1 player, got %d", convoy.PlayerCount())
	}
}

func TestManagerStartConvoy(t *testing.T) {
	m := NewManager()
	defer m.Close()

	convoy, _ := m.CreateConvoy(12345, "fantasy", 4, "host1", "Host Player")
	_, _ = m.JoinConvoy(convoy.Code, "p2", "Player Two")

	// Non-host can't start
	err := m.StartConvoy(convoy.ID, "p2")
	if err == nil {
		t.Error("non-host should not be able to start")
	}

	// Host can start
	err = m.StartConvoy(convoy.ID, "host1")
	if err != nil {
		t.Fatalf("StartConvoy() error = %v", err)
	}

	if convoy.State != StateStarted {
		t.Error("convoy should be started")
	}
}

func TestManagerRecordResult(t *testing.T) {
	m := NewManager()
	defer m.Close()

	convoy, _ := m.CreateConvoy(12345, "fantasy", 4, "host1", "Host Player")
	_, _ = m.JoinConvoy(convoy.Code, "p2", "Player Two")
	_ = m.StartConvoy(convoy.ID, "host1")

	result := &RunData{
		IsVictory:    true,
		Score:        1000,
		DaysTraveled: 30,
		Survivors:    4,
	}

	err := m.RecordResult("host1", result)
	if err != nil {
		t.Fatalf("RecordResult() error = %v", err)
	}

	run := convoy.GetRun("host1")
	if run.Score != 1000 {
		t.Errorf("expected score 1000, got %d", run.Score)
	}
}

func TestManagerListConvoys(t *testing.T) {
	m := NewManager()
	defer m.Close()

	_, _ = m.CreateConvoy(1, "fantasy", 4, "h1", "Host 1")
	c2, _ := m.CreateConvoy(2, "scifi", 4, "h2", "Host 2")

	// Start and complete c2
	_, _ = m.JoinConvoy(c2.Code, "p2", "Player 2")
	_ = m.StartConvoy(c2.ID, "h2")
	_ = m.RecordResult("h2", &RunData{})
	_ = m.RecordResult("p2", &RunData{})

	active := m.ListActive()
	if len(active) != 1 {
		t.Errorf("expected 1 active convoy, got %d", len(active))
	}

	waiting := m.ListWaiting()
	if len(waiting) != 1 {
		t.Errorf("expected 1 waiting convoy, got %d", len(waiting))
	}
}

func TestManagerCleanup(t *testing.T) {
	m := NewManager()
	defer m.Close()

	c, _ := m.CreateConvoy(1, "fantasy", 4, "h1", "Host 1")
	_, _ = m.JoinConvoy(c.Code, "p2", "Player 2")
	_ = m.StartConvoy(c.ID, "h1")
	_ = m.RecordResult("h1", &RunData{})
	_ = m.RecordResult("p2", &RunData{})

	// Manually set created time to old
	c.CreatedAt = time.Now().Add(-48 * time.Hour)

	removed := m.CleanupOld(24 * time.Hour)
	if removed != 1 {
		t.Errorf("expected 1 removed, got %d", removed)
	}

	if m.Count() != 0 {
		t.Errorf("expected 0 convoys, got %d", m.Count())
	}
}

func TestManagerExportImport(t *testing.T) {
	m := NewManager()
	defer m.Close()

	convoy, _ := m.CreateConvoy(12345, "fantasy", 4, "host1", "Host Player")
	_, _ = m.JoinConvoy(convoy.Code, "p2", "Player Two")

	data, err := m.ExportConvoy(convoy.ID)
	if err != nil {
		t.Fatalf("ExportConvoy() error = %v", err)
	}

	// Create new manager and import
	m2 := NewManager()
	defer m2.Close()

	imported, err := m2.ImportConvoy(data)
	if err != nil {
		t.Fatalf("ImportConvoy() error = %v", err)
	}

	if imported.Seed != convoy.Seed {
		t.Errorf("expected seed %d, got %d", convoy.Seed, imported.Seed)
	}
	if len(imported.Players) != len(convoy.Players) {
		t.Errorf("expected %d players, got %d", len(convoy.Players), len(imported.Players))
	}
}

func TestMessage(t *testing.T) {
	msg, err := NewMessage(MsgTypeJoin, "convoy1", "player1", map[string]string{"name": "Test"})
	if err != nil {
		t.Fatalf("NewMessage() error = %v", err)
	}

	if msg.Type != MsgTypeJoin {
		t.Errorf("expected type join, got %s", msg.Type)
	}

	data, err := msg.Marshal()
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	restored, err := UnmarshalMessage(data)
	if err != nil {
		t.Fatalf("UnmarshalMessage() error = %v", err)
	}

	if restored.Type != msg.Type {
		t.Errorf("type mismatch: got %s, want %s", restored.Type, msg.Type)
	}
	if restored.ConvoyID != msg.ConvoyID {
		t.Errorf("convoy ID mismatch")
	}
}

func TestManagerClosed(t *testing.T) {
	m := NewManager()
	m.Close()

	_, err := m.CreateConvoy(1, "fantasy", 4, "h1", "Host")
	if err != ErrManagerClosed {
		t.Errorf("expected ErrManagerClosed, got %v", err)
	}

	_, err = m.GetConvoy("any")
	if err != ErrManagerClosed {
		t.Errorf("expected ErrManagerClosed, got %v", err)
	}
}
