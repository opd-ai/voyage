package engine

import "sort"

// World manages entities and systems in the ECS framework.
type World struct {
	entities    map[EntityID]*Entity
	systems     []System
	genre       GenreID
	registry    *ComponentRegistry
	entityPool  []*Entity
	toSpawn     []*Entity
	toDespawn   []EntityID
	initialized bool
}

// NewWorld creates a new world with the given component registry.
func NewWorld(registry *ComponentRegistry) *World {
	return &World{
		entities:   make(map[EntityID]*Entity),
		systems:    make([]System, 0),
		genre:      GenreFantasy,
		registry:   registry,
		entityPool: make([]*Entity, 0),
		toSpawn:    make([]*Entity, 0),
		toDespawn:  make([]EntityID, 0),
	}
}

// SetGenre changes the genre for the world and all registered systems.
func (w *World) SetGenre(genreID GenreID) {
	w.genre = genreID
	for _, sys := range w.systems {
		sys.SetGenre(genreID)
	}
}

// Genre returns the current genre.
func (w *World) Genre() GenreID {
	return w.genre
}

// AddSystem registers a system with the world.
// Systems are automatically sorted by priority.
func (w *World) AddSystem(sys System) {
	sys.SetGenre(w.genre)
	w.systems = append(w.systems, sys)
	sort.Slice(w.systems, func(i, j int) bool {
		return w.systems[i].Priority() < w.systems[j].Priority()
	})
}

// RemoveSystem unregisters a system from the world.
func (w *World) RemoveSystem(sys System) {
	for i, s := range w.systems {
		if s == sys {
			w.systems = append(w.systems[:i], w.systems[i+1:]...)
			return
		}
	}
}

// Systems returns all registered systems.
func (w *World) Systems() []System {
	return w.systems
}

// Spawn creates a new entity and adds it to the world.
// The entity is added at the start of the next Update cycle.
func (w *World) Spawn() *Entity {
	var e *Entity
	if len(w.entityPool) > 0 {
		e = w.entityPool[len(w.entityPool)-1]
		w.entityPool = w.entityPool[:len(w.entityPool)-1]
		e.active = true
	} else {
		e = NewEntity()
	}
	w.toSpawn = append(w.toSpawn, e)
	return e
}

// SpawnImmediate creates a new entity and adds it to the world immediately.
func (w *World) SpawnImmediate() *Entity {
	var e *Entity
	if len(w.entityPool) > 0 {
		e = w.entityPool[len(w.entityPool)-1]
		w.entityPool = w.entityPool[:len(w.entityPool)-1]
		e.active = true
	} else {
		e = NewEntity()
	}
	w.entities[e.ID()] = e
	return e
}

// Despawn marks an entity for removal.
// The entity is removed at the start of the next Update cycle.
func (w *World) Despawn(id EntityID) {
	w.toDespawn = append(w.toDespawn, id)
}

// DespawnImmediate removes an entity from the world immediately.
// Clears all components and tags before pooling to prevent state corruption (C-005).
func (w *World) DespawnImmediate(id EntityID) {
	if e, ok := w.entities[id]; ok {
		e.active = false
		e.Clear() // Clear components and tags before pooling (C-005)
		delete(w.entities, id)
		w.entityPool = append(w.entityPool, e)
	}
}

// Entity retrieves an entity by ID.
// Returns nil if the entity does not exist.
func (w *World) Entity(id EntityID) *Entity {
	return w.entities[id]
}

// Entities returns all active entities in the world.
func (w *World) Entities() []*Entity {
	result := make([]*Entity, 0, len(w.entities))
	for _, e := range w.entities {
		if e.IsActive() {
			result = append(result, e)
		}
	}
	return result
}

// EntitiesWith returns all entities that have all specified components.
func (w *World) EntitiesWith(ids ...ComponentID) []*Entity {
	result := make([]*Entity, 0)
	for _, e := range w.entities {
		if e.IsActive() && e.HasAll(ids...) {
			result = append(result, e)
		}
	}
	return result
}

// EntitiesWithTag returns all entities that have the specified tag.
func (w *World) EntitiesWithTag(tag string) []*Entity {
	result := make([]*Entity, 0)
	for _, e := range w.entities {
		if e.IsActive() && e.HasTag(tag) {
			result = append(result, e)
		}
	}
	return result
}

// Update runs all systems for one frame/tick.
func (w *World) Update(dt float64) {
	w.processSpawnDespawn()

	for _, sys := range w.systems {
		sys.Update(w, dt)
	}
}

// processSpawnDespawn handles deferred entity spawning and despawning.
func (w *World) processSpawnDespawn() {
	for _, e := range w.toSpawn {
		w.entities[e.ID()] = e
	}
	w.toSpawn = w.toSpawn[:0]

	for _, id := range w.toDespawn {
		w.DespawnImmediate(id)
	}
	w.toDespawn = w.toDespawn[:0]
}

// EntityCount returns the number of active entities.
func (w *World) EntityCount() int {
	return len(w.entities)
}

// Registry returns the component registry.
func (w *World) Registry() *ComponentRegistry {
	return w.registry
}
