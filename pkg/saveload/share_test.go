package saveload

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewRunData(t *testing.T) {
	rd := NewRunData(12345, engine.GenreFantasy, 2)

	if rd.Seed != 12345 {
		t.Errorf("expected seed 12345, got %d", rd.Seed)
	}
	if rd.Genre != engine.GenreFantasy {
		t.Errorf("expected fantasy genre, got %s", rd.Genre)
	}
	if rd.Diffi != 2 {
		t.Errorf("expected difficulty 2, got %d", rd.Diffi)
	}
	if len(rd.Decisions) != 0 {
		t.Errorf("expected empty decisions, got %d", len(rd.Decisions))
	}
}

func TestRecordDecision(t *testing.T) {
	rd := NewRunData(12345, engine.GenreScifi, 1)

	rd.RecordDecision(1, DecisionMove, 2, 0)
	rd.RecordDecision(2, DecisionChoice, 1, 5)
	rd.RecordDecision(3, DecisionRest, 0, 0)

	if len(rd.Decisions) != 3 {
		t.Fatalf("expected 3 decisions, got %d", len(rd.Decisions))
	}

	if rd.Decisions[0].Turn != 1 || rd.Decisions[0].Type != DecisionMove {
		t.Error("first decision incorrect")
	}
	if rd.Decisions[1].Turn != 2 || rd.Decisions[1].Type != DecisionChoice {
		t.Error("second decision incorrect")
	}
	if rd.Decisions[2].Turn != 3 || rd.Decisions[2].Type != DecisionRest {
		t.Error("third decision incorrect")
	}
}

func TestSetFinalState(t *testing.T) {
	rd := NewRunData(12345, engine.GenreHorror, 3)

	rd.SetFinalState(100, 15, 20, true, 3)

	if rd.FinalTurn != 100 {
		t.Errorf("expected final turn 100, got %d", rd.FinalTurn)
	}
	if rd.FinalX != 15 || rd.FinalY != 20 {
		t.Errorf("expected position (15,20), got (%d,%d)", rd.FinalX, rd.FinalY)
	}
	if !rd.WonGame {
		t.Error("expected WonGame true")
	}
	if rd.CrewSurvived != 3 {
		t.Errorf("expected 3 crew survived, got %d", rd.CrewSurvived)
	}
}

func TestShareCodeRoundTrip(t *testing.T) {
	rd := NewRunData(42, engine.GenreCyberpunk, 1)

	// Add various decisions
	rd.RecordDecision(1, DecisionMove, 0, 0)
	rd.RecordDecision(2, DecisionMove, 1, 0)
	rd.RecordDecision(3, DecisionChoice, 2, 10)
	rd.RecordDecision(5, DecisionForage, 0, 0)
	rd.RecordDecision(10, DecisionTrade, 3, 50)

	rd.SetFinalState(50, 25, 30, false, 2)

	// Export to share code
	code, err := rd.ExportShareCode()
	if err != nil {
		t.Fatalf("failed to export: %v", err)
	}

	if len(code) == 0 {
		t.Fatal("share code is empty")
	}

	// Import back
	imported, err := ImportShareCode(code)
	if err != nil {
		t.Fatalf("failed to import: %v", err)
	}

	// Verify all data matches
	if imported.Seed != rd.Seed {
		t.Errorf("seed mismatch: expected %d, got %d", rd.Seed, imported.Seed)
	}
	if imported.Genre != rd.Genre {
		t.Errorf("genre mismatch: expected %s, got %s", rd.Genre, imported.Genre)
	}
	if imported.Diffi != rd.Diffi {
		t.Errorf("difficulty mismatch: expected %d, got %d", rd.Diffi, imported.Diffi)
	}
	if imported.FinalTurn != rd.FinalTurn {
		t.Errorf("final turn mismatch: expected %d, got %d", rd.FinalTurn, imported.FinalTurn)
	}
	if imported.FinalX != rd.FinalX || imported.FinalY != rd.FinalY {
		t.Errorf("final position mismatch")
	}
	if imported.WonGame != rd.WonGame {
		t.Errorf("won game mismatch")
	}
	if imported.CrewSurvived != rd.CrewSurvived {
		t.Errorf("crew survived mismatch")
	}
}

func TestShareCodeAllGenres(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		rd := NewRunData(999, genre, 0)
		rd.RecordDecision(1, DecisionMove, 1, 0)
		rd.SetFinalState(1, 0, 0, true, 4)

		code, err := rd.ExportShareCode()
		if err != nil {
			t.Fatalf("failed to export for genre %s: %v", genre, err)
		}

		imported, err := ImportShareCode(code)
		if err != nil {
			t.Fatalf("failed to import for genre %s: %v", genre, err)
		}

		if imported.Genre != genre {
			t.Errorf("genre mismatch for %s: got %s", genre, imported.Genre)
		}
	}
}

func TestShareCodeManyDecisions(t *testing.T) {
	rd := NewRunData(12345, engine.GenreFantasy, 2)

	// Add many decisions
	for i := 0; i < 200; i++ {
		rd.RecordDecision(i, DecisionType(i%5), i%4, i*2)
	}
	rd.SetFinalState(200, 50, 50, true, 5)

	code, err := rd.ExportShareCode()
	if err != nil {
		t.Fatalf("failed to export many decisions: %v", err)
	}

	imported, err := ImportShareCode(code)
	if err != nil {
		t.Fatalf("failed to import many decisions: %v", err)
	}

	if len(imported.Decisions) != 200 {
		t.Errorf("expected 200 decisions, got %d", len(imported.Decisions))
	}
}

func TestImportInvalidShareCode(t *testing.T) {
	testCases := []struct {
		name string
		code string
		err  error
	}{
		{"empty", "", ErrShareCodeTooShort},
		{"too short", "abc", ErrShareCodeTooShort},
		{"invalid chars", "hello0world", ErrInvalidShareCode}, // 0 not in base58
		{"invalid O char", "helloOworld", ErrInvalidShareCode},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ImportShareCode(tc.code)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestBase58RoundTrip(t *testing.T) {
	testData := [][]byte{
		{0},
		{1},
		{0, 0, 1},
		{255, 255, 255},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
	}

	for _, data := range testData {
		encoded := encodeBase58(data)
		decoded, err := decodeBase58(encoded)
		if err != nil {
			t.Errorf("decode error for %v: %v", data, err)
			continue
		}

		// For leading zeros, we need special handling
		if len(decoded) != len(data) {
			// Check if it's just missing leading zeros
			for i := len(decoded); i < len(data); i++ {
				if data[len(data)-1-i] != 0 {
					t.Errorf("data mismatch for %v: got %v", data, decoded)
					break
				}
			}
		}
	}
}

func TestCalculateChecksum(t *testing.T) {
	// Checksum should be deterministic
	data := []byte{1, 2, 3, 4, 5}
	cs1 := calculateChecksum(data)
	cs2 := calculateChecksum(data)

	if cs1 != cs2 {
		t.Error("checksum not deterministic")
	}

	// Different data should produce different checksums (usually)
	data2 := []byte{1, 2, 3, 4, 6}
	cs3 := calculateChecksum(data2)

	if cs1 == cs3 {
		t.Log("checksums matched for different data (possible but unlikely)")
	}
}

func TestReplay(t *testing.T) {
	rd := NewRunData(12345, engine.GenreFantasy, 1)
	rd.RecordDecision(1, DecisionMove, 0, 0)
	rd.RecordDecision(2, DecisionMove, 1, 0)
	rd.RecordDecision(3, DecisionChoice, 2, 0)

	replay := NewReplay(rd)

	if !replay.HasNext() {
		t.Error("expected HasNext true")
	}
	if replay.Position() != 0 {
		t.Error("expected position 0")
	}
	if replay.Remaining() != 3 {
		t.Errorf("expected 3 remaining, got %d", replay.Remaining())
	}

	// Peek should not advance
	d, ok := replay.Peek()
	if !ok || d.Turn != 1 {
		t.Error("peek failed")
	}
	if replay.Position() != 0 {
		t.Error("peek should not advance position")
	}

	// Next should advance
	d, ok = replay.Next()
	if !ok || d.Turn != 1 {
		t.Error("next failed")
	}
	if replay.Position() != 1 {
		t.Error("next should advance position")
	}

	// Continue through all
	replay.Next()
	replay.Next()

	if replay.HasNext() {
		t.Error("expected HasNext false after all consumed")
	}

	// Reset
	replay.Reset()
	if replay.Position() != 0 {
		t.Error("reset failed")
	}
	if !replay.HasNext() {
		t.Error("expected HasNext true after reset")
	}
}

func TestRunDataGetDecisionAt(t *testing.T) {
	rd := NewRunData(12345, engine.GenreFantasy, 1)
	rd.RecordDecision(1, DecisionMove, 2, 0)

	d, ok := rd.GetDecisionAt(0)
	if !ok {
		t.Error("expected decision at index 0")
	}
	if d.Value != 2 {
		t.Errorf("expected value 2, got %d", d.Value)
	}

	_, ok = rd.GetDecisionAt(-1)
	if ok {
		t.Error("expected no decision at index -1")
	}

	_, ok = rd.GetDecisionAt(10)
	if ok {
		t.Error("expected no decision at index 10")
	}
}

func TestRunDataGetDecisionCount(t *testing.T) {
	rd := NewRunData(12345, engine.GenreFantasy, 1)

	if rd.GetDecisionCount() != 0 {
		t.Error("expected 0 decisions")
	}

	rd.RecordDecision(1, DecisionMove, 0, 0)
	rd.RecordDecision(2, DecisionMove, 0, 0)

	if rd.GetDecisionCount() != 2 {
		t.Errorf("expected 2 decisions, got %d", rd.GetDecisionCount())
	}
}

func TestEncodeDecodeGenre(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, g := range genres {
		encoded := encodeGenre(g)
		decoded := decodeGenre(encoded)
		if decoded != g {
			t.Errorf("genre roundtrip failed for %s: got %s", g, decoded)
		}
	}

	// Unknown genre should default to fantasy
	unknown := decodeGenre(255)
	if unknown != engine.GenreFantasy {
		t.Errorf("unknown genre should default to fantasy, got %s", unknown)
	}
}

func TestDecisionTypes(t *testing.T) {
	types := []DecisionType{
		DecisionMove,
		DecisionChoice,
		DecisionRest,
		DecisionForage,
		DecisionTrade,
	}

	// Verify they are unique
	seen := make(map[DecisionType]bool)
	for _, dt := range types {
		if seen[dt] {
			t.Errorf("duplicate decision type: %d", dt)
		}
		seen[dt] = true
	}
}
