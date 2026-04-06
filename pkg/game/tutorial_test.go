//go:build headless

package game

import (
	"testing"
)

func TestTutorialPhaseProgression(t *testing.T) {
	tm := NewTutorialManager()

	// Starts at welcome phase
	if tm.Phase() != TutorialWelcome {
		t.Errorf("expected TutorialWelcome, got %d", tm.Phase())
	}
	if tm.IsComplete() {
		t.Error("tutorial should not be complete at start")
	}

	// Should show hints initially
	if !tm.ShouldShowHint() {
		t.Error("should show hint at welcome phase")
	}

	// Welcome hint text should mention arrow keys and WASD
	hint := tm.GetHintText()
	if hint == "" {
		t.Error("hint text should not be empty at welcome phase")
	}

	// Moving advances to movement phase
	tm.OnMove()
	if tm.Phase() != TutorialMovement {
		t.Errorf("expected TutorialMovement after move, got %d", tm.Phase())
	}

	// Advancing turns advances to resource phase
	for i := 1; i <= 5; i++ {
		tm.OnTurnAdvance(i)
	}
	if tm.Phase() != TutorialResources {
		t.Errorf("expected TutorialResources after 5 turns, got %d", tm.Phase())
	}

	// Seeing an event advances to events phase
	tm.OnEventSeen()
	if tm.Phase() != TutorialEvents {
		t.Errorf("expected TutorialEvents after seeing event, got %d", tm.Phase())
	}

	// Resolving an event completes the tutorial
	tm.OnEventResolved()
	if tm.Phase() != TutorialComplete {
		t.Errorf("expected TutorialComplete after resolving event, got %d", tm.Phase())
	}
	if !tm.IsComplete() {
		t.Error("tutorial should be complete")
	}
	if tm.ShouldShowHint() {
		t.Error("should not show hint when tutorial is complete")
	}
}

func TestTutorialEarlyGame(t *testing.T) {
	tm := NewTutorialManager()

	if !tm.IsEarlyGame(0) {
		t.Error("turn 0 should be early game")
	}
	if !tm.IsEarlyGame(2) {
		t.Error("turn 2 should be early game")
	}
	if tm.IsEarlyGame(3) {
		t.Error("turn 3 should not be early game")
	}
}

func TestTutorialSkipOnQuickEvent(t *testing.T) {
	tm := NewTutorialManager()

	// If an event fires before resource phase, skip directly to events
	tm.OnMove()
	tm.OnEventSeen()
	if tm.Phase() != TutorialEvents {
		t.Errorf("expected TutorialEvents, got %d", tm.Phase())
	}
}

func TestGetObjectiveText(t *testing.T) {
	text := GetObjectiveText()
	if text == "" {
		t.Error("objective text should not be empty")
	}
}

func TestGetControlsText(t *testing.T) {
	text := GetControlsText()
	if text == "" {
		t.Error("controls text should not be empty")
	}
}

func TestGetLoseReasonTip(t *testing.T) {
	conditions := []LoseCondition{
		LoseAllCrewDead,
		LoseVesselDestroyed,
		LoseMoraleZero,
		LoseStarvation,
		LoseNone,
	}
	for _, lc := range conditions {
		tip := GetLoseReasonTip(lc)
		if tip == "" {
			t.Errorf("tip should not be empty for condition %d", lc)
		}
	}
}

func TestGetResourceDescription(t *testing.T) {
	names := []string{"Food", "Water", "Fuel", "Medicine", "Morale", "Currency"}
	for _, name := range names {
		desc := GetResourceDescription(name)
		if desc == "" {
			t.Errorf("description should not be empty for resource %s", name)
		}
	}

	// Unknown resource should return empty
	if desc := GetResourceDescription("Unknown"); desc != "" {
		t.Errorf("expected empty description for unknown resource, got %s", desc)
	}
}

func TestDirectionArrow(t *testing.T) {
	tests := []struct {
		dx, dy int
		want   string
	}{
		{0, 0, "HERE!"},
		{0, -5, "N"},
		{0, 5, "S"},
		{5, 0, "E"},
		{-5, 0, "W"},
		{5, -5, "NE"},
		{-5, -5, "NW"},
		{5, 5, "SE"},
		{-5, 5, "SW"},
	}

	for _, tt := range tests {
		got := directionArrow(tt.dx, tt.dy)
		if got != tt.want {
			t.Errorf("directionArrow(%d, %d) = %s, want %s", tt.dx, tt.dy, got, tt.want)
		}
	}
}

func TestEarlyGameEventSuppression(t *testing.T) {
	cfg := DefaultSessionConfig()
	cfg.Seed = 42
	session := NewGameSession(cfg)

	// Tutorial should be initialized
	if session.Tutorial() == nil {
		t.Fatal("tutorial should be initialized")
	}

	// At turn 0, early game suppression should be active
	if !session.Tutorial().IsEarlyGame(0) {
		t.Error("turn 0 should be early game")
	}
}
