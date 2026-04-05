// Package council implements crew voting for critical decisions in Voyage.
// When facing dangerous shortcuts, costly detours, or other major route choices,
// crew members vote based on their personality traits.
//
// # Voting Mechanics
//
// Each crew member votes based on their dominant trait:
//   - Brave traits vote for risky but faster options
//   - Cautious traits vote for safer options
//   - Greedy traits vote for profitable options
//   - Other traits have nuanced voting patterns
//
// # Player Interaction
//
// The player may:
//   - Follow the crew's decision (no morale change)
//   - Overrule the vote (morale penalty proportional to dissent)
//   - Receive a morale bonus for unanimous agreement
//
// # Genre Support
//
// Council scenes are re-skinned per genre through SetGenre():
//   - Fantasy: Campfire debate
//   - Scifi: Bridge briefing
//   - Horror: Group argument
//   - Cyberpunk: Exec meeting
//   - Postapoc: Bonfire council
package council
