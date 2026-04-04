// Package destination implements procedural destination generation and arrival
// sequences for Voyage. It provides multiple destination types per run, discovery
// events as players approach, and arrival ceremony sequences with genre-appropriate
// narrative text.
//
// # Destination Types
//
// The package supports five base destination archetypes that are re-skinned per genre:
//   - City: A populated settlement or outpost
//   - Sanctuary: A safe haven or refuge
//   - Treasure: A valuable cache or resource depot
//   - Escape: An exit point or evacuation craft
//   - Settlement: A frontier community or colony
//
// # Discovery Events
//
// As the party approaches a destination, discovery events reveal information:
//   - Distant sighting events at long range
//   - Signs of civilization/activity at medium range
//   - Approach events at close range
//   - Arrival events at the destination
//
// # Genre Support
//
// All destination names, descriptions, and narrative text adapt to the current genre
// through the SetGenre() method, which re-skins content while preserving mechanics.
package destination
