// Package narrative provides procedural story arc generation for Voyage.
//
// This package generates deterministic narrative content from seed including:
//   - Three-act structure (departure crisis → mid-journey revelation → arrival twist)
//   - Named recurring NPCs (friend, nemesis, or ambiguous figure)
//   - Crew backstory events that connect to the destination
//
// # Story Arc
//
// Each run has a procedurally generated story arc with:
//   - Departure Crisis: An event that sets the journey in motion
//   - Mid-Journey Revelation: A discovery that changes the stakes
//   - Arrival Twist: A final reveal at the destination
//
// # Recurring NPCs
//
// A named NPC reappears throughout the journey with consistent
// characterization determined by the seed. This NPC may be:
//   - A friend who helps the party
//   - A nemesis who opposes them
//   - An ambiguous figure with unclear motives
//
// # Genre Support
//
// All narrative elements implement SetGenre() to adjust vocabulary
// and tone to match the current genre theme.
package narrative
