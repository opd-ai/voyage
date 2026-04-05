package input

import "time"

// Direction represents a movement direction.
type Direction int

const (
	// DirectionNone indicates no direction.
	DirectionNone Direction = iota
	// DirectionUp indicates upward movement.
	DirectionUp
	// DirectionDown indicates downward movement.
	DirectionDown
	// DirectionLeft indicates leftward movement.
	DirectionLeft
	// DirectionRight indicates rightward movement.
	DirectionRight
)

// Action represents a discrete user action.
type Action int

const (
	// ActionNone indicates no action.
	ActionNone Action = iota
	// ActionConfirm is triggered by Enter, Space, or tap.
	ActionConfirm
	// ActionCancel is triggered by Escape or back gesture.
	ActionCancel
	// ActionDebug toggles debug mode (F3).
	ActionDebug
	// ActionOption1 through ActionOption9 are number key options.
	ActionOption1
	ActionOption2
	ActionOption3
	ActionOption4
	ActionOption5
	ActionOption6
	ActionOption7
	ActionOption8
	ActionOption9
)

// TouchState tracks the state of a single touch point.
type TouchState struct {
	// ID is the unique identifier for this touch.
	ID int
	// StartX is the initial X position when touch began.
	StartX int
	// StartY is the initial Y position when touch began.
	StartY int
	// CurrentX is the current X position.
	CurrentX int
	// CurrentY is the current Y position.
	CurrentY int
	// StartTime is when the touch began.
	StartTime time.Time
	// IsActive indicates if this touch is still active.
	IsActive bool
}

// SwipeThreshold is the minimum distance (in pixels) for a swipe.
const SwipeThreshold = 30

// TapMaxDuration is the maximum time for a touch to be considered a tap.
const TapMaxDuration = 300 * time.Millisecond

// TapMaxDistance is the maximum movement (in pixels) for a tap.
const TapMaxDistance = 20

// SwipeResult represents the detected swipe direction.
type SwipeResult struct {
	// Direction is the detected swipe direction.
	Direction Direction
	// Distance is the swipe distance in pixels.
	Distance int
	// Duration is how long the swipe took.
	Duration time.Duration
}

// InputState holds the current frame's input state.
type InputState struct {
	// Direction is the current movement direction (keyboard or swipe).
	Direction Direction
	// Actions holds actions triggered this frame.
	Actions []Action
	// TapPosition holds the tap position if a tap was detected this frame.
	TapPosition *Position
	// ScreenSize holds the current screen dimensions.
	ScreenSize Size
}

// Position represents a screen coordinate.
type Position struct {
	X int
	Y int
}

// Size represents dimensions.
type Size struct {
	Width  int
	Height int
}

// HasAction checks if a specific action was triggered this frame.
func (s *InputState) HasAction(action Action) bool {
	for _, a := range s.Actions {
		if a == action {
			return true
		}
	}
	return false
}

// String returns a string representation of the direction.
func (d Direction) String() string {
	switch d {
	case DirectionUp:
		return "up"
	case DirectionDown:
		return "down"
	case DirectionLeft:
		return "left"
	case DirectionRight:
		return "right"
	default:
		return "none"
	}
}

// String returns a string representation of the action.
func (a Action) String() string {
	switch a {
	case ActionConfirm:
		return "confirm"
	case ActionCancel:
		return "cancel"
	case ActionDebug:
		return "debug"
	case ActionOption1:
		return "option1"
	case ActionOption2:
		return "option2"
	case ActionOption3:
		return "option3"
	case ActionOption4:
		return "option4"
	case ActionOption5:
		return "option5"
	case ActionOption6:
		return "option6"
	case ActionOption7:
		return "option7"
	case ActionOption8:
		return "option8"
	case ActionOption9:
		return "option9"
	default:
		return "none"
	}
}
