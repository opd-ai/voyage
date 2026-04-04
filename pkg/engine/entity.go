package engine

import "sync/atomic"

// EntityID uniquely identifies an entity within a world.
type EntityID uint64

// entityCounter provides unique IDs for entities.
var entityCounter uint64

// nextEntityID generates a new unique entity ID.
func nextEntityID() EntityID {
	return EntityID(atomic.AddUint64(&entityCounter, 1))
}

// Entity is a container holding a set of components.
// Entities are pure data containers with no behavior.
type Entity struct {
	id         EntityID
	components map[ComponentID]Component
	active     bool
	tags       map[string]struct{}
}

// NewEntity creates a new entity with a unique ID.
func NewEntity() *Entity {
	return &Entity{
		id:         nextEntityID(),
		components: make(map[ComponentID]Component),
		active:     true,
		tags:       make(map[string]struct{}),
	}
}

// ID returns the entity's unique identifier.
func (e *Entity) ID() EntityID {
	return e.id
}

// Add attaches a component to the entity.
// If a component of the same type exists, it is replaced.
func (e *Entity) Add(c Component) {
	e.components[c.ID()] = c
}

// Remove detaches a component from the entity.
func (e *Entity) Remove(id ComponentID) {
	delete(e.components, id)
}

// Get retrieves a component by its ID.
// Returns nil if the component is not attached.
func (e *Entity) Get(id ComponentID) Component {
	return e.components[id]
}

// Has checks if the entity has a component of the given type.
func (e *Entity) Has(id ComponentID) bool {
	_, ok := e.components[id]
	return ok
}

// HasAll checks if the entity has all the specified component types.
func (e *Entity) HasAll(ids ...ComponentID) bool {
	for _, id := range ids {
		if !e.Has(id) {
			return false
		}
	}
	return true
}

// Components returns all components attached to the entity.
func (e *Entity) Components() []Component {
	result := make([]Component, 0, len(e.components))
	for _, c := range e.components {
		result = append(result, c)
	}
	return result
}

// IsActive returns whether the entity is active.
func (e *Entity) IsActive() bool {
	return e.active
}

// SetActive sets the entity's active state.
func (e *Entity) SetActive(active bool) {
	e.active = active
}

// AddTag adds a tag to the entity.
func (e *Entity) AddTag(tag string) {
	e.tags[tag] = struct{}{}
}

// RemoveTag removes a tag from the entity.
func (e *Entity) RemoveTag(tag string) {
	delete(e.tags, tag)
}

// HasTag checks if the entity has the specified tag.
func (e *Entity) HasTag(tag string) bool {
	_, ok := e.tags[tag]
	return ok
}
