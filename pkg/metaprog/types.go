package metaprog

import (
	"crypto/sha256"
	"encoding/binary"
	"time"

	"github.com/opd-ai/voyage/pkg/engine"
)

// UnlockID uniquely identifies an unlockable item.
type UnlockID string

// UnlockCategory categorizes unlockable content.
type UnlockCategory int

const (
	// CategoryEvent represents event types that have been seen.
	CategoryEvent UnlockCategory = iota
	// CategoryDestination represents destinations that have been reached.
	CategoryDestination
	// CategoryCrewArchetype represents crew starting archetypes.
	CategoryCrewArchetype
	// CategoryVesselConfig represents vessel starting configurations.
	CategoryVesselConfig
)

// Unlock represents an unlockable item.
type Unlock struct {
	ID          UnlockID
	Category    UnlockCategory
	Name        string
	Description string
	Unlocked    bool
	UnlockHash  int64 // Hash threshold for unlocking
}

// NewUnlock creates a new unlock definition.
func NewUnlock(id UnlockID, cat UnlockCategory, name, desc string, hash int64) *Unlock {
	return &Unlock{
		ID:          id,
		Category:    cat,
		Name:        name,
		Description: desc,
		Unlocked:    false,
		UnlockHash:  hash,
	}
}

// RunSummary captures the results of a completed run.
type RunSummary struct {
	Seed          int64
	Genre         engine.GenreID
	DaysTraveled  int
	CrewSurvivors int
	FinalScore    int
	EventsSeen    []string
	DestinationID string
	CompletedAt   time.Time
	Victory       bool
}

// HallOfRecordsEntry stores the best run for a genre.
type HallOfRecordsEntry struct {
	Genre       engine.GenreID
	BestScore   int
	BestSummary *RunSummary
	TotalRuns   int
	TotalWins   int
}

// UnlockLog tracks all discoveries and unlocks.
type UnlockLog struct {
	EventsSeen          map[string]bool
	DestinationsReached map[string]bool
	TotalRuns           int
	TotalWins           int
	CumulativeHash      int64
	CrewArchetypes      map[UnlockID]*Unlock
	VesselConfigs       map[UnlockID]*Unlock
}

// NewUnlockLog creates a new unlock log.
func NewUnlockLog() *UnlockLog {
	log := &UnlockLog{
		EventsSeen:          make(map[string]bool),
		DestinationsReached: make(map[string]bool),
		CrewArchetypes:      make(map[UnlockID]*Unlock),
		VesselConfigs:       make(map[UnlockID]*Unlock),
	}
	log.initializeUnlocks()
	return log
}

// initializeUnlocks sets up the default unlockable content.
func (l *UnlockLog) initializeUnlocks() {
	// Crew archetypes - unlock thresholds based on cumulative hash
	crewUnlocks := []struct {
		id   UnlockID
		name string
		desc string
		hash int64
	}{
		{"veteran_leader", "Veteran Leader", "A seasoned leader with bonus morale", 1000},
		{"skilled_medic", "Skilled Medic", "An experienced healer with better healing", 2000},
		{"master_mechanic", "Master Mechanic", "Expert at repairs and maintenance", 3000},
		{"lucky_scavenger", "Lucky Scavenger", "Finds more resources when foraging", 5000},
		{"hardened_survivor", "Hardened Survivor", "Extra health and resistance", 7500},
		{"wise_navigator", "Wise Navigator", "Reveals more of the map", 10000},
	}

	for _, cu := range crewUnlocks {
		l.CrewArchetypes[cu.id] = NewUnlock(cu.id, CategoryCrewArchetype, cu.name, cu.desc, cu.hash)
	}

	// Vessel configurations
	vesselUnlocks := []struct {
		id   UnlockID
		name string
		desc string
		hash int64
	}{
		{"cargo_hauler", "Cargo Hauler", "Extra cargo capacity, slower speed", 1500},
		{"fast_runner", "Fast Runner", "Higher speed, less cargo space", 2500},
		{"armored_transport", "Armored Transport", "More hull integrity, slower", 4000},
		{"balanced_cruiser", "Balanced Cruiser", "Well-rounded starting stats", 6000},
		{"medical_frigate", "Medical Frigate", "Bonus medicine storage and healing", 8000},
		{"explorer_vessel", "Explorer Vessel", "Extended range and visibility", 12000},
	}

	for _, vu := range vesselUnlocks {
		l.VesselConfigs[vu.id] = NewUnlock(vu.id, CategoryVesselConfig, vu.name, vu.desc, vu.hash)
	}
}

// RecordEvent marks an event type as seen.
func (l *UnlockLog) RecordEvent(eventType string) {
	if !l.EventsSeen[eventType] {
		l.EventsSeen[eventType] = true
		l.updateCumulativeHash(eventType)
	}
}

// RecordDestination marks a destination as reached.
func (l *UnlockLog) RecordDestination(destID string) {
	if !l.DestinationsReached[destID] {
		l.DestinationsReached[destID] = true
		l.updateCumulativeHash(destID)
	}
}

// RecordRun adds a completed run to the log.
func (l *UnlockLog) RecordRun(summary *RunSummary) {
	l.TotalRuns++
	if summary.Victory {
		l.TotalWins++
	}

	// Record all events seen in this run
	for _, event := range summary.EventsSeen {
		l.RecordEvent(event)
	}

	// Record destination if reached
	if summary.Victory && summary.DestinationID != "" {
		l.RecordDestination(summary.DestinationID)
	}

	// Update cumulative hash with run results
	l.updateCumulativeHashWithRun(summary)
	l.checkUnlocks()
}

// updateCumulativeHash updates the hash with new discovery.
func (l *UnlockLog) updateCumulativeHash(data string) {
	h := sha256.New()
	seedBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(seedBytes, uint64(l.CumulativeHash))
	h.Write(seedBytes)
	h.Write([]byte(data))
	sum := h.Sum(nil)
	l.CumulativeHash = int64(binary.LittleEndian.Uint64(sum[:8]))
}

// updateCumulativeHashWithRun updates the hash based on run results.
func (l *UnlockLog) updateCumulativeHashWithRun(summary *RunSummary) {
	// Add points based on run performance
	points := int64(summary.FinalScore)
	if summary.Victory {
		points += 500
	}
	points += int64(summary.CrewSurvivors * 100)
	points += int64(len(summary.EventsSeen) * 10)

	l.CumulativeHash += points
	if l.CumulativeHash < 0 {
		l.CumulativeHash = -l.CumulativeHash
	}
}

// checkUnlocks evaluates which items should be unlocked.
func (l *UnlockLog) checkUnlocks() {
	for _, unlock := range l.CrewArchetypes {
		if !unlock.Unlocked && l.CumulativeHash >= unlock.UnlockHash {
			unlock.Unlocked = true
		}
	}
	for _, unlock := range l.VesselConfigs {
		if !unlock.Unlocked && l.CumulativeHash >= unlock.UnlockHash {
			unlock.Unlocked = true
		}
	}
}

// GetUnlockedCrewArchetypes returns all unlocked crew archetypes.
func (l *UnlockLog) GetUnlockedCrewArchetypes() []*Unlock {
	var result []*Unlock
	for _, u := range l.CrewArchetypes {
		if u.Unlocked {
			result = append(result, u)
		}
	}
	return result
}

// GetUnlockedVesselConfigs returns all unlocked vessel configurations.
func (l *UnlockLog) GetUnlockedVesselConfigs() []*Unlock {
	var result []*Unlock
	for _, u := range l.VesselConfigs {
		if u.Unlocked {
			result = append(result, u)
		}
	}
	return result
}

// GetLockedCrewArchetypes returns all locked crew archetypes.
func (l *UnlockLog) GetLockedCrewArchetypes() []*Unlock {
	var result []*Unlock
	for _, u := range l.CrewArchetypes {
		if !u.Unlocked {
			result = append(result, u)
		}
	}
	return result
}

// GetLockedVesselConfigs returns all locked vessel configurations.
func (l *UnlockLog) GetLockedVesselConfigs() []*Unlock {
	var result []*Unlock
	for _, u := range l.VesselConfigs {
		if !u.Unlocked {
			result = append(result, u)
		}
	}
	return result
}

// HasSeenEvent returns true if the event type has been seen.
func (l *UnlockLog) HasSeenEvent(eventType string) bool {
	return l.EventsSeen[eventType]
}

// HasReachedDestination returns true if the destination was reached.
func (l *UnlockLog) HasReachedDestination(destID string) bool {
	return l.DestinationsReached[destID]
}

// GetSeenEventsCount returns the number of unique events seen.
func (l *UnlockLog) GetSeenEventsCount() int {
	return len(l.EventsSeen)
}

// GetReachedDestinationsCount returns the number of destinations reached.
func (l *UnlockLog) GetReachedDestinationsCount() int {
	return len(l.DestinationsReached)
}

// HallOfRecords tracks best runs per genre.
type HallOfRecords struct {
	Entries map[engine.GenreID]*HallOfRecordsEntry
}

// NewHallOfRecords creates a new hall of records.
func NewHallOfRecords() *HallOfRecords {
	hall := &HallOfRecords{
		Entries: make(map[engine.GenreID]*HallOfRecordsEntry),
	}
	// Initialize entries for all genres
	for _, genre := range engine.AllGenres() {
		hall.Entries[genre] = &HallOfRecordsEntry{
			Genre: genre,
		}
	}
	return hall
}

// RecordRun records a run result, updating if it's a new best.
func (h *HallOfRecords) RecordRun(summary *RunSummary) {
	entry := h.Entries[summary.Genre]
	if entry == nil {
		entry = &HallOfRecordsEntry{Genre: summary.Genre}
		h.Entries[summary.Genre] = entry
	}

	entry.TotalRuns++
	if summary.Victory {
		entry.TotalWins++
	}

	// Update best score if this is higher
	if summary.FinalScore > entry.BestScore {
		entry.BestScore = summary.FinalScore
		entry.BestSummary = summary
	}
}

// GetEntry returns the hall of records entry for a genre.
func (h *HallOfRecords) GetEntry(genre engine.GenreID) *HallOfRecordsEntry {
	return h.Entries[genre]
}

// GetBestRun returns the best run summary for a genre.
func (h *HallOfRecords) GetBestRun(genre engine.GenreID) *RunSummary {
	if entry := h.Entries[genre]; entry != nil {
		return entry.BestSummary
	}
	return nil
}

// GetTotalRuns returns total runs across all genres.
func (h *HallOfRecords) GetTotalRuns() int {
	total := 0
	for _, entry := range h.Entries {
		total += entry.TotalRuns
	}
	return total
}

// GetTotalWins returns total wins across all genres.
func (h *HallOfRecords) GetTotalWins() int {
	total := 0
	for _, entry := range h.Entries {
		total += entry.TotalWins
	}
	return total
}

// MetaProgress combines unlock log and hall of records.
type MetaProgress struct {
	UnlockLog     *UnlockLog
	HallOfRecords *HallOfRecords
}

// NewMetaProgress creates a new meta-progression tracker.
func NewMetaProgress() *MetaProgress {
	return &MetaProgress{
		UnlockLog:     NewUnlockLog(),
		HallOfRecords: NewHallOfRecords(),
	}
}

// RecordRun records a completed run in both unlock log and hall of records.
func (m *MetaProgress) RecordRun(summary *RunSummary) {
	m.UnlockLog.RecordRun(summary)
	m.HallOfRecords.RecordRun(summary)
}
