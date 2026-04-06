package game

// TutorialPhase represents the current phase of the tutorial.
type TutorialPhase int

const (
	// TutorialWelcome is shown at the start of the game.
	TutorialWelcome TutorialPhase = iota
	// TutorialMovement teaches the player to move.
	TutorialMovement
	// TutorialResources introduces resource management.
	TutorialResources
	// TutorialEvents introduces the event system.
	TutorialEvents
	// TutorialComplete means the tutorial is finished.
	TutorialComplete
)

// tutorialPhaseThresholds defines turn thresholds for each phase transition.
var tutorialPhaseThresholds = map[TutorialPhase]int{
	TutorialWelcome:   0, // Shown immediately
	TutorialMovement:  1, // After first move
	TutorialResources: 5, // After a few turns
	TutorialEvents:    0, // Triggered by first event, not by turn count
}

// TutorialManager tracks tutorial progress and generates contextual hints.
type TutorialManager struct {
	phase          TutorialPhase
	turnsSinceMove int
	hasSeenEvent   bool
	dismissed      bool // Player dismissed the current hint
}

// NewTutorialManager creates a new tutorial manager starting at the welcome phase.
func NewTutorialManager() *TutorialManager {
	return &TutorialManager{
		phase: TutorialWelcome,
	}
}

// Phase returns the current tutorial phase.
func (tm *TutorialManager) Phase() TutorialPhase {
	return tm.phase
}

// IsComplete returns true if the tutorial is finished.
func (tm *TutorialManager) IsComplete() bool {
	return tm.phase >= TutorialComplete
}

// IsEarlyGame returns true if the player is still in the first few turns
// where events should be suppressed for orientation.
func (tm *TutorialManager) IsEarlyGame(turn int) bool {
	return turn < 3
}

// Dismiss hides the current hint until the phase advances.
func (tm *TutorialManager) Dismiss() {
	tm.dismissed = true
}

// OnMove is called when the player moves, advancing tutorial state.
func (tm *TutorialManager) OnMove() {
	tm.turnsSinceMove = 0
	if tm.phase == TutorialWelcome {
		tm.phase = TutorialMovement
		tm.dismissed = false
	}
}

// OnTurnAdvance is called each turn to track progress.
func (tm *TutorialManager) OnTurnAdvance(turn int) {
	tm.turnsSinceMove++
	if tm.phase == TutorialMovement && turn >= tutorialPhaseThresholds[TutorialResources] {
		tm.phase = TutorialResources
		tm.dismissed = false
	}
}

// OnEventSeen is called when the player encounters their first event.
func (tm *TutorialManager) OnEventSeen() {
	if tm.hasSeenEvent {
		return
	}
	tm.hasSeenEvent = true
	if tm.phase == TutorialResources || tm.phase == TutorialMovement {
		tm.phase = TutorialEvents
		tm.dismissed = false
	}
}

// OnEventResolved is called when the player resolves an event.
func (tm *TutorialManager) OnEventResolved() {
	if tm.phase == TutorialEvents {
		tm.phase = TutorialComplete
		tm.dismissed = false
	}
}

// ShouldShowHint returns true if a tutorial hint should be displayed.
func (tm *TutorialManager) ShouldShowHint() bool {
	return !tm.IsComplete() && !tm.dismissed
}

// GetHintText returns the current tutorial hint text.
func (tm *TutorialManager) GetHintText() string {
	switch tm.phase {
	case TutorialWelcome:
		return "Use ARROW KEYS to move your vessel toward the destination."
	case TutorialMovement:
		return "Each move costs resources. Watch your Food, Water, and Fuel."
	case TutorialResources:
		return "Keep Morale above zero! Low Food or Water drains Morale quickly."
	case TutorialEvents:
		return "Press 1-4 to choose during events. Choose wisely!"
	default:
		return ""
	}
}

// GetObjectiveText returns a description of the current objective.
func GetObjectiveText() string {
	return "Reach the destination with at least one crew member alive."
}

// GetControlsText returns a summary of game controls.
func GetControlsText() string {
	return "ARROW KEYS: Move | ESC: Pause | 1-4: Event Choices | R: Rest | F3: Debug"
}

// GetLoseReasonTip returns a tip based on the loss condition.
func GetLoseReasonTip(lc LoseCondition) string {
	switch lc {
	case LoseAllCrewDead:
		return "Tip: Use Medicine and rest often to keep your crew healthy."
	case LoseVesselDestroyed:
		return "Tip: Avoid hazardous events and choose cautious options to protect your vessel."
	case LoseMoraleZero:
		return "Tip: Keep Food and Water stocked to prevent morale collapse."
	case LoseStarvation:
		return "Tip: Forage for supplies regularly and manage consumption carefully."
	default:
		return "Tip: Balance resource use and take fewer risks on your next journey."
	}
}

// GetResourceDescription returns a brief description of what a resource does.
func GetResourceDescription(name string) string {
	descriptions := map[string]string{
		"Food":     "Consumed each turn per crew member. Depletion drops Morale.",
		"Water":    "Consumed each turn per crew member. Depletion drops Morale fast.",
		"Fuel":     "Consumed each move based on vessel speed. Without it, you stop.",
		"Medicine": "Used to heal injured crew during events.",
		"Morale":   "Drops when resources deplete. At zero, your crew deserts you.",
		"Currency": "Trade for supplies at destinations. Earned through events.",
	}
	if desc, ok := descriptions[name]; ok {
		return desc
	}
	return ""
}

// directionArrow returns a cardinal direction string from a delta vector.
func directionArrow(dx, dy int) string {
	if dx == 0 && dy == 0 {
		return "HERE!"
	}

	var dir string
	if dy < 0 {
		dir = "N"
	} else if dy > 0 {
		dir = "S"
	}
	if dx > 0 {
		dir += "E"
	} else if dx < 0 {
		dir += "W"
	}
	return dir
}
