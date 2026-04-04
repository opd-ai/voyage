package engine

// SystemPriority determines the execution order of systems.
// Lower values execute first.
type SystemPriority int

const (
	// PriorityInput handles input processing first.
	PriorityInput SystemPriority = 100
	// PriorityLogic handles game logic after input.
	PriorityLogic SystemPriority = 200
	// PriorityPhysics handles physics calculations.
	PriorityPhysics SystemPriority = 300
	// PriorityAudio handles audio processing.
	PriorityAudio SystemPriority = 400
	// PriorityRender handles rendering last.
	PriorityRender SystemPriority = 500
)

// System processes entities with specific component configurations.
// All systems must implement GenreSwitcher to support genre switching.
type System interface {
	GenreSwitcher

	// Update processes entities for one frame/tick.
	// The world provides access to entities and other systems.
	Update(world *World, dt float64)

	// Priority returns the system's execution priority.
	// Lower values execute first.
	Priority() SystemPriority

	// RequiredComponents returns the component types this system requires.
	// Only entities with all required components will be processed.
	RequiredComponents() []ComponentID
}

// BaseSystem provides a default implementation of common System methods.
// Embed this in concrete systems and override as needed.
type BaseSystem struct {
	genre    GenreID
	priority SystemPriority
	required []ComponentID
}

// NewBaseSystem creates a new base system with the given priority.
func NewBaseSystem(priority SystemPriority, required ...ComponentID) BaseSystem {
	return BaseSystem{
		genre:    GenreFantasy,
		priority: priority,
		required: required,
	}
}

// SetGenre implements GenreSwitcher.
func (s *BaseSystem) SetGenre(genreID GenreID) {
	s.genre = genreID
}

// Genre returns the current genre.
func (s *BaseSystem) Genre() GenreID {
	return s.genre
}

// Priority returns the system's execution priority.
func (s *BaseSystem) Priority() SystemPriority {
	return s.priority
}

// RequiredComponents returns the component types this system requires.
func (s *BaseSystem) RequiredComponents() []ComponentID {
	return s.required
}
