//go:build headless

package input

import (
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	m := NewManager()
	if m == nil {
		t.Fatal("expected manager to be created")
	}
	if !m.touchEnabled {
		t.Error("touch should be enabled by default")
	}
}

func TestManagerUpdate(t *testing.T) {
	m := NewManager()
	m.Update()

	state := m.State()
	if state.Direction != DirectionNone {
		t.Errorf("expected no direction, got %v", state.Direction)
	}
	if len(state.Actions) != 0 {
		t.Errorf("expected no actions, got %v", state.Actions)
	}
}

func TestManagerSimulateDirection(t *testing.T) {
	m := NewManager()

	m.SimulateDirection(DirectionUp)
	m.Update()

	if m.GetDirection() != DirectionUp {
		t.Errorf("expected up direction, got %v", m.GetDirection())
	}
}

func TestManagerSimulateAction(t *testing.T) {
	m := NewManager()

	m.SimulateAction(ActionConfirm)
	m.Update()

	if !m.JustConfirmed() {
		t.Error("expected confirm action")
	}
}

func TestManagerSimulateTap(t *testing.T) {
	m := NewManager()

	m.SimulateTap(100, 200)
	m.Update()

	pos := m.GetTapPosition()
	if pos == nil {
		t.Fatal("expected tap position")
	}
	if pos.X != 100 || pos.Y != 200 {
		t.Errorf("expected position (100, 200), got (%d, %d)", pos.X, pos.Y)
	}
	if !m.JustConfirmed() {
		t.Error("tap should also trigger confirm")
	}
}

func TestManagerJustCancelled(t *testing.T) {
	m := NewManager()

	m.SimulateAction(ActionCancel)
	m.Update()

	if !m.JustCancelled() {
		t.Error("expected cancel action")
	}
}

func TestManagerJustDebugToggled(t *testing.T) {
	m := NewManager()

	m.SimulateAction(ActionDebug)
	m.Update()

	if !m.JustDebugToggled() {
		t.Error("expected debug action")
	}
}

func TestManagerGetOptionPressed(t *testing.T) {
	tests := []struct {
		action Action
		want   int
	}{
		{ActionOption1, 1},
		{ActionOption2, 2},
		{ActionOption3, 3},
		{ActionOption4, 4},
		{ActionOption5, 5},
		{ActionOption6, 6},
		{ActionOption7, 7},
		{ActionOption8, 8},
		{ActionOption9, 9},
		{ActionConfirm, 0},
	}

	for _, tt := range tests {
		m := NewManager()
		m.SimulateAction(tt.action)
		m.Update()

		got := m.GetOptionPressed()
		if got != tt.want {
			t.Errorf("GetOptionPressed() for %v = %d, want %d", tt.action, got, tt.want)
		}
	}
}

func TestManagerClearSimulation(t *testing.T) {
	m := NewManager()

	m.SimulateDirection(DirectionLeft)
	m.SimulateAction(ActionConfirm)
	m.ClearSimulation()
	m.Update()

	if m.GetDirection() != DirectionNone {
		t.Error("expected no direction after clear")
	}
	if m.JustConfirmed() {
		t.Error("expected no confirm after clear")
	}
}

func TestManagerSetTouchEnabled(t *testing.T) {
	m := NewManager()

	m.SetTouchEnabled(false)
	if m.touchEnabled {
		t.Error("expected touch to be disabled")
	}

	m.SetTouchEnabled(true)
	if !m.touchEnabled {
		t.Error("expected touch to be enabled")
	}
}

func TestManagerSetKeyRepeatDelay(t *testing.T) {
	m := NewManager()

	m.SetKeyRepeatDelay(500 * time.Millisecond)
	if m.keyRepeatDelay != 500*time.Millisecond {
		t.Errorf("expected delay 500ms, got %v", m.keyRepeatDelay)
	}
}

func TestManagerSetKeyRepeatInterval(t *testing.T) {
	m := NewManager()

	m.SetKeyRepeatInterval(50 * time.Millisecond)
	if m.keyRepeatInterval != 50*time.Millisecond {
		t.Errorf("expected interval 50ms, got %v", m.keyRepeatInterval)
	}
}

func TestInputStateHasAction(t *testing.T) {
	state := InputState{
		Actions: []Action{ActionConfirm, ActionDebug},
	}

	if !state.HasAction(ActionConfirm) {
		t.Error("expected to have ActionConfirm")
	}
	if !state.HasAction(ActionDebug) {
		t.Error("expected to have ActionDebug")
	}
	if state.HasAction(ActionCancel) {
		t.Error("should not have ActionCancel")
	}
}

func TestDirectionString(t *testing.T) {
	tests := []struct {
		dir  Direction
		want string
	}{
		{DirectionNone, "none"},
		{DirectionUp, "up"},
		{DirectionDown, "down"},
		{DirectionLeft, "left"},
		{DirectionRight, "right"},
	}

	for _, tt := range tests {
		if got := tt.dir.String(); got != tt.want {
			t.Errorf("%v.String() = %q, want %q", tt.dir, got, tt.want)
		}
	}
}

func TestActionString(t *testing.T) {
	tests := []struct {
		action Action
		want   string
	}{
		{ActionNone, "none"},
		{ActionConfirm, "confirm"},
		{ActionCancel, "cancel"},
		{ActionDebug, "debug"},
		{ActionOption1, "option1"},
		{ActionOption9, "option9"},
	}

	for _, tt := range tests {
		if got := tt.action.String(); got != tt.want {
			t.Errorf("%v.String() = %q, want %q", tt.action, got, tt.want)
		}
	}
}

func TestManagerActionsCleared(t *testing.T) {
	m := NewManager()

	// First update with action
	m.SimulateAction(ActionConfirm)
	m.Update()
	if !m.JustConfirmed() {
		t.Error("expected confirm on first update")
	}

	// Second update should not have action
	m.Update()
	if m.JustConfirmed() {
		t.Error("action should be cleared on second update")
	}
}

func TestManagerKeyRepeatLogic(t *testing.T) {
	m := NewManager()
	m.SetKeyRepeatDelay(50 * time.Millisecond)
	m.SetKeyRepeatInterval(20 * time.Millisecond)

	// First press - should register
	m.SimulateDirection(DirectionUp)
	m.Update()
	if m.GetDirection() != DirectionUp {
		t.Error("first press should register direction")
	}

	// Immediate second update - should not repeat (delay not met)
	m.SimulateDirection(DirectionUp)
	m.Update()
	if m.GetDirection() != DirectionNone {
		t.Error("should not repeat before delay")
	}

	// Wait for delay and repeat
	time.Sleep(60 * time.Millisecond)
	m.SimulateDirection(DirectionUp)
	m.Update()
	if m.GetDirection() != DirectionUp {
		t.Error("should repeat after delay")
	}
}

func TestManagerDirectionChangeResets(t *testing.T) {
	m := NewManager()
	m.SetKeyRepeatDelay(50 * time.Millisecond)

	// Press up
	m.SimulateDirection(DirectionUp)
	m.Update()
	if m.GetDirection() != DirectionUp {
		t.Error("first press should register")
	}

	// Change to down - should register immediately
	m.SimulateDirection(DirectionDown)
	m.Update()
	if m.GetDirection() != DirectionDown {
		t.Error("direction change should register immediately")
	}
}
