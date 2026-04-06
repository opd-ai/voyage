//go:build !headless

package input

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Manager handles all input from keyboard, mouse, and touch.
type Manager struct {
	// touches tracks active touch points.
	touches map[ebiten.TouchID]*TouchState
	// lastDirection holds the last detected direction for key repeat.
	lastDirection Direction
	// directionHeldSince tracks when direction key was first pressed.
	directionHeldSince time.Time
	// lastDirectionRepeat tracks last repeat time.
	lastDirectionRepeat time.Time
	// currentState holds the computed state for this frame.
	currentState InputState
	// keyRepeatDelay is the initial delay before key repeat starts.
	keyRepeatDelay time.Duration
	// keyRepeatInterval is the interval between repeated inputs.
	keyRepeatInterval time.Duration
	// prevTouchIDs tracks touch IDs from previous frame.
	prevTouchIDs []ebiten.TouchID
	// touchEnabled allows disabling touch input.
	touchEnabled bool
}

// NewManager creates a new input manager.
func NewManager() *Manager {
	return &Manager{
		touches:           make(map[ebiten.TouchID]*TouchState),
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

// Reset clears all input state for scene transitions (H-008).
// Call this when transitioning between scenes to prevent phantom inputs.
func (m *Manager) Reset() {
	m.lastDirection = DirectionNone
	m.directionHeldSince = time.Time{}
	m.lastDirectionRepeat = time.Time{}
	m.currentState = InputState{
		Direction: DirectionNone,
		Actions:   nil,
	}
}

// Update processes input for the current frame.
// Call this once per frame before querying input state.
func (m *Manager) Update() {
	m.currentState = InputState{
		Direction: DirectionNone,
		Actions:   nil,
	}

	// Get screen size
	w, h := ebiten.WindowSize()
	m.currentState.ScreenSize = Size{Width: w, Height: h}

	// Process keyboard input
	m.processKeyboard()

	// Process touch input
	if m.touchEnabled {
		m.processTouch()
	}

	// Process mouse input (for desktop browser)
	m.processMouse()
}

// processKeyboard handles keyboard input.
func (m *Manager) processKeyboard() {
	m.processDirectionInput()
	m.processActionKeys()
	m.processOptionKeys()
}

// processDirectionInput handles arrow/WASD key input with repeat.
func (m *Manager) processDirectionInput() {
	dir := m.getKeyboardDirection()
	if dir != DirectionNone {
		m.handleDirectionPressed(dir)
	} else {
		m.lastDirection = DirectionNone
	}
}

// handleDirectionPressed processes a direction key being held.
func (m *Manager) handleDirectionPressed(dir Direction) {
	now := time.Now()
	if dir != m.lastDirection {
		m.lastDirection = dir
		m.directionHeldSince = now
		m.lastDirectionRepeat = now
		m.currentState.Direction = dir
		return
	}
	// Same direction held - check for repeat
	heldDuration := now.Sub(m.directionHeldSince)
	if heldDuration >= m.keyRepeatDelay {
		if now.Sub(m.lastDirectionRepeat) >= m.keyRepeatInterval {
			m.lastDirectionRepeat = now
			m.currentState.Direction = dir
		}
	}
}

// processActionKeys handles confirm, cancel, and debug keys.
// Uses separate conditions to prevent duplicate ActionConfirm (M-004).
func (m *Manager) processActionKeys() {
	// Check for confirm action, only add once even if both keys pressed (M-004)
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
		inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		// Only add if not already present
		if !m.hasAction(ActionConfirm) {
			m.currentState.Actions = append(m.currentState.Actions, ActionConfirm)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		m.currentState.Actions = append(m.currentState.Actions, ActionCancel)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		m.currentState.Actions = append(m.currentState.Actions, ActionDebug)
	}
}

// hasAction checks if an action is already in the current state (M-004).
func (m *Manager) hasAction(action InputAction) bool {
	for _, a := range m.currentState.Actions {
		if a == action {
			return true
		}
	}
	return false
}

// processOptionKeys handles number keys 1-9 for options.
func (m *Manager) processOptionKeys() {
	optionKeys := []ebiten.Key{
		ebiten.Key1, ebiten.Key2, ebiten.Key3,
		ebiten.Key4, ebiten.Key5, ebiten.Key6,
		ebiten.Key7, ebiten.Key8, ebiten.Key9,
	}
	optionActions := []Action{
		ActionOption1, ActionOption2, ActionOption3,
		ActionOption4, ActionOption5, ActionOption6,
		ActionOption7, ActionOption8, ActionOption9,
	}
	for i, key := range optionKeys {
		if inpututil.IsKeyJustPressed(key) {
			m.currentState.Actions = append(m.currentState.Actions, optionActions[i])
		}
	}
}

// getKeyboardDirection returns the current direction from keyboard.
func (m *Manager) getKeyboardDirection() Direction {
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		return DirectionUp
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		return DirectionDown
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		return DirectionLeft
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		return DirectionRight
	}
	return DirectionNone
}

// processTouch handles touch input.
func (m *Manager) processTouch() {
	touchIDs := ebiten.AppendTouchIDs(nil)
	m.trackNewTouches(touchIDs)
	m.handleEndedTouches(touchIDs)
	m.prevTouchIDs = touchIDs
}

// trackNewTouches registers new touch points and updates existing ones.
func (m *Manager) trackNewTouches(touchIDs []ebiten.TouchID) {
	for _, id := range touchIDs {
		x, y := ebiten.TouchPosition(id)
		if _, exists := m.touches[id]; !exists {
			m.touches[id] = &TouchState{
				ID:        int(id),
				StartX:    x,
				StartY:    y,
				CurrentX:  x,
				CurrentY:  y,
				StartTime: time.Now(),
				IsActive:  true,
			}
		} else {
			m.touches[id].CurrentX = x
			m.touches[id].CurrentY = y
		}
	}
}

// handleEndedTouches processes touches that ended this frame.
func (m *Manager) handleEndedTouches(currentIDs []ebiten.TouchID) {
	currentIDSet := make(map[ebiten.TouchID]bool)
	for _, id := range currentIDs {
		currentIDSet[id] = true
	}
	for _, prevID := range m.prevTouchIDs {
		if !currentIDSet[prevID] {
			if touch, exists := m.touches[prevID]; exists {
				m.handleTouchEnd(touch)
				delete(m.touches, prevID)
			}
		}
	}
}

// handleTouchEnd processes a touch that just ended.
func (m *Manager) handleTouchEnd(touch *TouchState) {
	duration := time.Since(touch.StartTime)
	dx := touch.CurrentX - touch.StartX
	dy := touch.CurrentY - touch.StartY
	distance := int(math.Sqrt(float64(dx*dx + dy*dy)))

	// Check if it's a tap
	if duration < TapMaxDuration && distance < TapMaxDistance {
		m.currentState.Actions = append(m.currentState.Actions, ActionConfirm)
		m.currentState.TapPosition = &Position{
			X: touch.CurrentX,
			Y: touch.CurrentY,
		}
		return
	}

	// Check if it's a swipe
	if distance >= SwipeThreshold {
		swipe := m.detectSwipeDirection(dx, dy, distance, duration)
		if swipe.Direction != DirectionNone {
			m.currentState.Direction = swipe.Direction
		}
	}
}

// detectSwipeDirection determines swipe direction from deltas.
func (m *Manager) detectSwipeDirection(dx, dy, distance int, duration time.Duration) SwipeResult {
	result := SwipeResult{
		Direction: DirectionNone,
		Distance:  distance,
		Duration:  duration,
	}

	// Determine primary axis
	absDX := abs(dx)
	absDY := abs(dy)

	if absDX > absDY {
		// Horizontal swipe
		if dx > 0 {
			result.Direction = DirectionRight
		} else {
			result.Direction = DirectionLeft
		}
	} else {
		// Vertical swipe
		if dy > 0 {
			result.Direction = DirectionDown
		} else {
			result.Direction = DirectionUp
		}
	}

	return result
}

// processMouse handles mouse input for desktop.
func (m *Manager) processMouse() {
	// Left click as confirm
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		m.currentState.TapPosition = &Position{X: x, Y: y}
		// Only add confirm if we don't already have a tap from touch
		if !m.currentState.HasAction(ActionConfirm) {
			m.currentState.Actions = append(m.currentState.Actions, ActionConfirm)
		}
	}

	// Right click as cancel
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		if !m.currentState.HasAction(ActionCancel) {
			m.currentState.Actions = append(m.currentState.Actions, ActionCancel)
		}
	}
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

// abs returns the absolute value of an int.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
