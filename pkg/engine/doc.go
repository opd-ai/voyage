// Package engine provides the Entity-Component-System (ECS) framework for Voyage.
//
// The ECS architecture enables genre-switchable game systems where every System
// implements the GenreSwitcher interface to swap thematic presentation at runtime.
//
// Core types:
//   - GenreID: Identifies one of five supported genre themes
//   - GenreSwitcher: Interface for genre-aware systems
//   - Component: Interface for entity data containers
//   - Entity: Container holding a set of components
//   - System: Interface for game logic processing entities
//   - World: Manager for entities and systems
package engine
