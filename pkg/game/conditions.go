package game

import (
	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/procgen/world"
	"github.com/opd-ai/voyage/pkg/resources"
	"github.com/opd-ai/voyage/pkg/vessel"
)

// WinCondition represents the victory condition type.
type WinCondition int

const (
	// WinReachedDestination means the vessel reached the goal.
	WinReachedDestination WinCondition = iota
)

// LoseCondition represents the defeat condition type.
type LoseCondition int

const (
	// LoseNone means no loss condition has been met.
	LoseNone LoseCondition = iota
	// LoseAllCrewDead means all crew members died.
	LoseAllCrewDead
	// LoseVesselDestroyed means the vessel was destroyed.
	LoseVesselDestroyed
	// LoseMoraleZero means morale hit zero causing desertion.
	LoseMoraleZero
	// LoseStarvation means the party starved to death.
	LoseStarvation
)

// ConditionChecker evaluates win and loss conditions.
type ConditionChecker struct {
	vesselX, vesselY int
}

// NewConditionChecker creates a new condition checker.
func NewConditionChecker() *ConditionChecker {
	return &ConditionChecker{}
}

// SetVesselPosition updates the vessel's current position.
func (cc *ConditionChecker) SetVesselPosition(x, y int) {
	cc.vesselX = x
	cc.vesselY = y
}

// CheckWin checks if the win condition has been met.
func (cc *ConditionChecker) CheckWin(worldMap *world.WorldMap, party *crew.Party) (bool, WinCondition) {
	if worldMap == nil || party == nil {
		return false, WinReachedDestination
	}

	// Win: Reached destination with at least one living crew member
	dest := worldMap.Destination
	if cc.vesselX == dest.X && cc.vesselY == dest.Y {
		if party.LivingCount() > 0 {
			return true, WinReachedDestination
		}
	}

	return false, WinReachedDestination
}

// CheckLose checks if any loss condition has been met.
func (cc *ConditionChecker) CheckLose(party *crew.Party, v *vessel.Vessel, res *resources.Resources) (bool, LoseCondition) {
	// Check party conditions
	if party != nil && party.LivingCount() == 0 {
		return true, LoseAllCrewDead
	}

	// Check vessel condition
	if v != nil && v.IsDestroyed() {
		return true, LoseVesselDestroyed
	}

	// Check resource conditions
	if res != nil {
		if res.IsDepleted(resources.ResourceMorale) {
			return true, LoseMoraleZero
		}
		// Starvation requires food and water both depleted
		if res.IsDepleted(resources.ResourceFood) && res.IsDepleted(resources.ResourceWater) {
			return true, LoseStarvation
		}
	}

	return false, LoseNone
}

// WinConditionName returns a human-readable name for the win condition.
func WinConditionName(wc WinCondition) string {
	switch wc {
	case WinReachedDestination:
		return "Reached Destination"
	default:
		return "Victory"
	}
}

// LoseConditionName returns a human-readable name for the loss condition.
func LoseConditionName(lc LoseCondition) string {
	switch lc {
	case LoseAllCrewDead:
		return "All Crew Lost"
	case LoseVesselDestroyed:
		return "Vessel Destroyed"
	case LoseMoraleZero:
		return "Crew Deserted"
	case LoseStarvation:
		return "Starved"
	default:
		return "Unknown"
	}
}

// WinConditionDescription returns a detailed description.
func WinConditionDescription(wc WinCondition) string {
	switch wc {
	case WinReachedDestination:
		return "You have successfully completed your journey!"
	default:
		return "Congratulations on your victory!"
	}
}

// LoseConditionDescription returns a detailed description.
func LoseConditionDescription(lc LoseCondition) string {
	switch lc {
	case LoseAllCrewDead:
		return "Your entire crew has perished. The journey ends here."
	case LoseVesselDestroyed:
		return "Your vessel has been destroyed. There is no way to continue."
	case LoseMoraleZero:
		return "Morale has collapsed completely. Your crew has abandoned you."
	case LoseStarvation:
		return "Without food or water, your party could not survive."
	default:
		return "Your journey has come to an unfortunate end."
	}
}
