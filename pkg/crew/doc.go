// Package crew provides the party/crew management system for Voyage.
//
// Features:
//   - Party entity with 2-6 crew member slots
//   - Procedurally generated crew members with names, traits, and skills
//   - Individual health tracking per crew member
//   - Crew mortality from starvation, disease, and injury
//   - Genre-switchable names and portrait palettes
//
// All crew generation is deterministic based on the master seed.
package crew
