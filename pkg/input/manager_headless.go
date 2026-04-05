//go:build headless

package input

import "time"

// Manager handles all input from keyboard, mouse, and touch.
// This is the headless version for testing.
type Manager struct {
	currentState        InputState
	simulatedDirection  Direction
	simulatedActions    []Action
	simulatedTap        *Position
	keyRepeatDelay      time.Duration
	keyRepeatInterval   time.Duration
	touchEnabled        bool
	lastDirection       Direction
	directionHeldSince  time.Time
	lastDirectionRepeat time.Time
}

// NewManager creates a new input manager.
func NewManager() *Manager {
	return &Manager{
		keyRepeatDelay:    400 * time.Millisecond,
		keyRepeatInterval: 100 * time.Millisecond,
		touchEnabled:      true,
	}
}

// SetKeyRepeatDelay sets the initial delay before key repeat starts.
func (m *Manager) SetKeyRepeatDelay(d time.Duration) {
	m.keyRepeatDelay = d
}

// SetKeyRepeatInterval sets the interval between repeated inputs.
func (m *Manager) SetKeyRepeatInterval(d time.Duration) {
	m.keyRepeatInterval = d
}

// SetTouchEnabled enables or disables touch input processing.
func (m *Manager) SetTouchEnabled(enabled bool) {
	m.touchEnabled = enabled
}

// Update processes input for the current frame.
func (m *Manager) Update() {
	m.currentState = InputState{
		Direction:   m.simulatedDirection,
		Actions:     m.simulatedActions,
		TapPosition: m.simulatedTap,
		ScreenSize:  Size{Width: 800, Height: 600},
	}

	// Handle direction with key repeat logic
	if m.simulatedDirection != DirectionNone {
		now := time.Now()
		if m.simulatedDirection != m.lastDirection {
			m.lastDirection = m.simulatedDirection
			m.directionHeldSince = now
			m.lastDirectionRepeat = now
			m.currentState.Direction = m.simulatedDirection
		} else {
			heldDuration := now.Sub(m.directionHeldSince)
			if heldDuration >= m.keyRepeatDelay {
				sinceLast := now.Sub(m.lastDirectionRepeat)
				if sinceLast >= m.keyRepeatInterval {
					m.lastDirectionRepeat = now
					m.currentState.Direction = m.simulatedDirection
				} else {
					m.currentState.Direction = DirectionNone
				}
			} else {
				m.currentState.Direction = DirectionNone
			}
		}
	} else {
		m.lastDirection = DirectionNone
	}

	// Clear simulated state after processing
	m.simulatedActions = nil
	m.simulatedTap = nil
}

// State returns the current input state.
func (m *Manager) State() InputState {
	return m.currentState
}

// JustConfirmed returns true if confirm was triggered this frame.
func (m *Manager) JustConfirmed() bool {
	return m.currentState.HasAction(ActionConfirm)
}

// JustCancelled returns true if cancel was triggered this frame.
func (m *Manager) JustCancelled() bool {
	return m.currentState.HasAction(ActionCancel)
}

// JustDebugToggled returns true if debug was triggered this frame.
func (m *Manager) JustDebugToggled() bool {
	return m.currentState.HasAction(ActionDebug)
}

// GetDirection returns the current movement direction.
func (m *Manager) GetDirection() Direction {
	return m.currentState.Direction
}

// GetTapPosition returns the tap/click position if one occurred this frame.
func (m *Manager) GetTapPosition() *Position {
	return m.currentState.TapPosition
}

// GetOptionPressed returns the option number (1-9) if pressed, or 0 if none.
func (m *Manager) GetOptionPressed() int {
	for _, action := range m.currentState.Actions {
		switch action {
		case ActionOption1:
			return 1
		case ActionOption2:
			return 2
		case ActionOption3:
			return 3
		case ActionOption4:
			return 4
		case ActionOption5:
			return 5
		case ActionOption6:
			return 6
		case ActionOption7:
			return 7
		case ActionOption8:
			return 8
		case ActionOption9:
			return 9
		}
	}
	return 0
}

// SimulateDirection simulates a direction input for testing.
func (m *Manager) SimulateDirection(dir Direction) {
	m.simulatedDirection = dir
}

// SimulateAction simulates an action input for testing.
func (m *Manager) SimulateAction(action Action) {
	m.simulatedActions = append(m.simulatedActions, action)
}

// SimulateTap simulates a tap at the given position for testing.
func (m *Manager) SimulateTap(x, y int) {
	m.simulatedTap = &Position{X: x, Y: y}
	m.simulatedActions = append(m.simulatedActions, ActionConfirm)
}

// ClearSimulation clears all simulated inputs.
func (m *Manager) ClearSimulation() {
	m.simulatedDirection = DirectionNone
	m.simulatedActions = nil
	m.simulatedTap = nil
}
