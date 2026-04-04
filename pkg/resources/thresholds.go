package resources

// ThresholdStatus indicates the warning level of a resource.
type ThresholdStatus int

const (
	// StatusNormal indicates the resource is at healthy levels.
	StatusNormal ThresholdStatus = iota
	// StatusLow indicates the resource is getting low.
	StatusLow
	// StatusCritical indicates the resource is critically low.
	StatusCritical
	// StatusDepleted indicates the resource is empty.
	StatusDepleted
)

// Threshold defines warning levels for a resource.
type Threshold struct {
	Low      float64 // ratio below which "low" warning triggers
	Critical float64 // ratio below which "critical" warning triggers
}

// DefaultThresholds returns the default warning thresholds.
func DefaultThresholds() map[ResourceType]Threshold {
	return map[ResourceType]Threshold{
		ResourceFood:     {Low: 0.30, Critical: 0.10},
		ResourceWater:    {Low: 0.30, Critical: 0.10},
		ResourceFuel:     {Low: 0.25, Critical: 0.08},
		ResourceMedicine: {Low: 0.40, Critical: 0.15},
		ResourceMorale:   {Low: 0.30, Critical: 0.10},
		ResourceCurrency: {Low: 0.20, Critical: 0.05},
	}
}

// GetThresholdStatus returns the status for a resource at the given ratio.
func GetThresholdStatus(rt ResourceType, ratio float64) ThresholdStatus {
	if ratio <= 0 {
		return StatusDepleted
	}

	thresholds := DefaultThresholds()
	t, ok := thresholds[rt]
	if !ok {
		return StatusNormal
	}

	if ratio <= t.Critical {
		return StatusCritical
	}
	if ratio <= t.Low {
		return StatusLow
	}
	return StatusNormal
}

// StatusString returns a human-readable status string.
func (s ThresholdStatus) String() string {
	switch s {
	case StatusNormal:
		return "Normal"
	case StatusLow:
		return "Low"
	case StatusCritical:
		return "Critical"
	case StatusDepleted:
		return "Depleted"
	default:
		return "Unknown"
	}
}

// IsWarning returns true if the status requires attention.
func (s ThresholdStatus) IsWarning() bool {
	return s >= StatusLow
}

// IsCritical returns true if the status is critical or depleted.
func (s ThresholdStatus) IsCritical() bool {
	return s >= StatusCritical
}
