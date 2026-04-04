package engine

// ComponentID uniquely identifies a component type.
type ComponentID string

// Component is the interface for all entity data containers.
// Components hold data but contain no logic.
type Component interface {
	// ID returns the unique identifier for this component type.
	ID() ComponentID
}

// ComponentRegistry maintains a mapping of component types.
type ComponentRegistry struct {
	factories map[ComponentID]func() Component
}

// NewComponentRegistry creates a new component registry.
func NewComponentRegistry() *ComponentRegistry {
	return &ComponentRegistry{
		factories: make(map[ComponentID]func() Component),
	}
}

// Register adds a component factory to the registry.
func (r *ComponentRegistry) Register(id ComponentID, factory func() Component) {
	r.factories[id] = factory
}

// Create instantiates a new component of the given type.
// Returns nil if the component type is not registered.
func (r *ComponentRegistry) Create(id ComponentID) Component {
	if factory, ok := r.factories[id]; ok {
		return factory()
	}
	return nil
}

// Has checks if a component type is registered.
func (r *ComponentRegistry) Has(id ComponentID) bool {
	_, ok := r.factories[id]
	return ok
}
