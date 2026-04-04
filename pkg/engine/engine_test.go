package engine

import (
	"testing"
)

func TestGenreIDConstants(t *testing.T) {
	genres := AllGenres()
	if len(genres) != 5 {
		t.Errorf("Expected 5 genres, got %d", len(genres))
	}

	expected := []GenreID{GenreFantasy, GenreScifi, GenreHorror, GenreCyberpunk, GenrePostapoc}
	for i, g := range expected {
		if genres[i] != g {
			t.Errorf("Genre %d: expected %s, got %s", i, g, genres[i])
		}
	}
}

func TestIsValidGenre(t *testing.T) {
	tests := []struct {
		input string
		valid bool
	}{
		{"fantasy", true},
		{"scifi", true},
		{"horror", true},
		{"cyberpunk", true},
		{"postapoc", true},
		{"invalid", false},
		{"", false},
		{"Fantasy", false},
	}

	for _, tt := range tests {
		got := IsValidGenre(tt.input)
		if got != tt.valid {
			t.Errorf("IsValidGenre(%q) = %v, want %v", tt.input, got, tt.valid)
		}
	}
}

type mockComponent struct {
	id ComponentID
}

func (c *mockComponent) ID() ComponentID { return c.id }

func TestComponentRegistry(t *testing.T) {
	r := NewComponentRegistry()

	posID := ComponentID("position")
	r.Register(posID, func() Component {
		return &mockComponent{id: posID}
	})

	if !r.Has(posID) {
		t.Error("Expected registry to have position component")
	}

	c := r.Create(posID)
	if c == nil {
		t.Error("Expected component to be created")
	}
	if c.ID() != posID {
		t.Errorf("Expected component ID %s, got %s", posID, c.ID())
	}

	// Test non-existent component
	if r.Has(ComponentID("nonexistent")) {
		t.Error("Expected registry to not have nonexistent component")
	}
	if r.Create(ComponentID("nonexistent")) != nil {
		t.Error("Expected nil for nonexistent component")
	}
}

func TestEntity(t *testing.T) {
	e := NewEntity()

	if e.ID() == 0 {
		t.Error("Expected non-zero entity ID")
	}
	if !e.IsActive() {
		t.Error("Expected new entity to be active")
	}

	// Test component operations
	posID := ComponentID("position")
	pos := &mockComponent{id: posID}
	e.Add(pos)

	if !e.Has(posID) {
		t.Error("Expected entity to have position component")
	}
	if e.Get(posID) != pos {
		t.Error("Expected to get same component back")
	}

	velID := ComponentID("velocity")
	vel := &mockComponent{id: velID}
	e.Add(vel)

	if !e.HasAll(posID, velID) {
		t.Error("Expected entity to have both components")
	}

	components := e.Components()
	if len(components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(components))
	}

	e.Remove(posID)
	if e.Has(posID) {
		t.Error("Expected position component to be removed")
	}

	// Test tags
	e.AddTag("player")
	if !e.HasTag("player") {
		t.Error("Expected entity to have player tag")
	}
	e.RemoveTag("player")
	if e.HasTag("player") {
		t.Error("Expected player tag to be removed")
	}

	// Test active state
	e.SetActive(false)
	if e.IsActive() {
		t.Error("Expected entity to be inactive")
	}
}

func TestEntityUniqueIDs(t *testing.T) {
	e1 := NewEntity()
	e2 := NewEntity()
	e3 := NewEntity()

	if e1.ID() == e2.ID() || e2.ID() == e3.ID() || e1.ID() == e3.ID() {
		t.Error("Expected unique entity IDs")
	}
}

type mockSystem struct {
	BaseSystem
	updateCount int
}

func (s *mockSystem) Update(world *World, dt float64) {
	s.updateCount++
}

func TestWorld(t *testing.T) {
	r := NewComponentRegistry()
	w := NewWorld(r)

	if w.Genre() != GenreFantasy {
		t.Errorf("Expected default genre fantasy, got %s", w.Genre())
	}

	// Test genre switching
	w.SetGenre(GenreScifi)
	if w.Genre() != GenreScifi {
		t.Errorf("Expected genre scifi, got %s", w.Genre())
	}

	// Test entity spawning
	e := w.SpawnImmediate()
	if e == nil {
		t.Error("Expected entity to be created")
	}
	if w.EntityCount() != 1 {
		t.Errorf("Expected 1 entity, got %d", w.EntityCount())
	}
	if w.Entity(e.ID()) != e {
		t.Error("Expected to retrieve same entity")
	}

	// Test deferred spawning
	e2 := w.Spawn()
	if w.EntityCount() != 1 {
		t.Error("Expected deferred spawn to not add immediately")
	}
	w.Update(0)
	if w.EntityCount() != 2 {
		t.Errorf("Expected 2 entities after update, got %d", w.EntityCount())
	}

	// Test despawning
	w.DespawnImmediate(e.ID())
	if w.EntityCount() != 1 {
		t.Errorf("Expected 1 entity after despawn, got %d", w.EntityCount())
	}

	w.Despawn(e2.ID())
	if w.EntityCount() != 1 {
		t.Error("Expected deferred despawn to not remove immediately")
	}
	w.Update(0)
	if w.EntityCount() != 0 {
		t.Errorf("Expected 0 entities after update, got %d", w.EntityCount())
	}
}

func TestWorldSystems(t *testing.T) {
	r := NewComponentRegistry()
	w := NewWorld(r)

	sys1 := &mockSystem{BaseSystem: NewBaseSystem(PriorityLogic)}
	sys2 := &mockSystem{BaseSystem: NewBaseSystem(PriorityRender)}
	sys3 := &mockSystem{BaseSystem: NewBaseSystem(PriorityInput)}

	w.AddSystem(sys1)
	w.AddSystem(sys2)
	w.AddSystem(sys3)

	// Test priority ordering
	systems := w.Systems()
	if systems[0].Priority() != PriorityInput {
		t.Error("Expected input system first")
	}
	if systems[1].Priority() != PriorityLogic {
		t.Error("Expected logic system second")
	}
	if systems[2].Priority() != PriorityRender {
		t.Error("Expected render system last")
	}

	// Test update
	w.Update(1.0 / 60.0)
	if sys1.updateCount != 1 || sys2.updateCount != 1 || sys3.updateCount != 1 {
		t.Error("Expected all systems to be updated once")
	}

	// Test genre propagation
	w.SetGenre(GenreHorror)
	if sys1.Genre() != GenreHorror {
		t.Error("Expected genre to propagate to systems")
	}

	// Test system removal
	w.RemoveSystem(sys2)
	if len(w.Systems()) != 2 {
		t.Errorf("Expected 2 systems after removal, got %d", len(w.Systems()))
	}
}

func TestWorldEntitiesWith(t *testing.T) {
	r := NewComponentRegistry()
	w := NewWorld(r)

	posID := ComponentID("position")
	velID := ComponentID("velocity")

	e1 := w.SpawnImmediate()
	e1.Add(&mockComponent{id: posID})
	e1.Add(&mockComponent{id: velID})

	e2 := w.SpawnImmediate()
	e2.Add(&mockComponent{id: posID})

	e3 := w.SpawnImmediate()
	e3.Add(&mockComponent{id: velID})

	withPos := w.EntitiesWith(posID)
	if len(withPos) != 2 {
		t.Errorf("Expected 2 entities with position, got %d", len(withPos))
	}

	withBoth := w.EntitiesWith(posID, velID)
	if len(withBoth) != 1 {
		t.Errorf("Expected 1 entity with both, got %d", len(withBoth))
	}

	// Test tag filtering
	e1.AddTag("player")
	tagged := w.EntitiesWithTag("player")
	if len(tagged) != 1 {
		t.Errorf("Expected 1 tagged entity, got %d", len(tagged))
	}
}

func TestBaseSystem(t *testing.T) {
	posID := ComponentID("position")
	velID := ComponentID("velocity")

	bs := NewBaseSystem(PriorityLogic, posID, velID)

	if bs.Priority() != PriorityLogic {
		t.Errorf("Expected priority %d, got %d", PriorityLogic, bs.Priority())
	}

	required := bs.RequiredComponents()
	if len(required) != 2 {
		t.Errorf("Expected 2 required components, got %d", len(required))
	}

	bs.SetGenre(GenreCyberpunk)
	if bs.Genre() != GenreCyberpunk {
		t.Errorf("Expected genre cyberpunk, got %s", bs.Genre())
	}
}
