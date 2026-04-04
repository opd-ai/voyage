// Package saveload provides game state persistence for Voyage.
//
// Features:
//   - Multiple save slots with autosave on turn advance
//   - Full game state serialization (JSON format)
//   - Seed embedded in save for reproducibility verification
//   - Load state restoration with validation
package saveload
