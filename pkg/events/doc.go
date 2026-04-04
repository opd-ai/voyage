// Package events provides the procedural event system for Voyage.
//
// Features:
//   - Seeded event queue based on position and turn
//   - Event categories: weather, encounter, discovery, hardship, windfall
//   - Choice-based event resolution with 2-4 options
//   - Grammar-based text generation from templates
//   - Outcome application affecting resources, crew, and vessel
//   - Genre-switchable event vocabulary
//
// All event text is procedurally generated - no pre-written content.
package events
