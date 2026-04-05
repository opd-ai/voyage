package metaprog

import (
	"testing"
	"time"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewUnlockLog(t *testing.T) {
	log := NewUnlockLog()

	if log.TotalRuns != 0 {
		t.Error("initial total runs should be 0")
	}
	if log.TotalWins != 0 {
		t.Error("initial total wins should be 0")
	}
	if len(log.CrewArchetypes) == 0 {
		t.Error("crew archetypes should be initialized")
	}
	if len(log.VesselConfigs) == 0 {
		t.Error("vessel configs should be initialized")
	}
}

func TestUnlockLogRecordEvent(t *testing.T) {
	log := NewUnlockLog()

	log.RecordEvent("weather_storm")
	if !log.HasSeenEvent("weather_storm") {
		t.Error("event should be recorded")
	}

	// Recording same event shouldn't increase count
	log.RecordEvent("weather_storm")
	if log.GetSeenEventsCount() != 1 {
		t.Error("duplicate events shouldn't increase count")
	}
}

func TestUnlockLogRecordDestination(t *testing.T) {
	log := NewUnlockLog()

	log.RecordDestination("legendary_city")
	if !log.HasReachedDestination("legendary_city") {
		t.Error("destination should be recorded")
	}

	log.RecordDestination("legendary_city")
	if log.GetReachedDestinationsCount() != 1 {
		t.Error("duplicate destinations shouldn't increase count")
	}
}

func TestUnlockLogRecordRun(t *testing.T) {
	log := NewUnlockLog()

	summary := &RunSummary{
		Seed:          12345,
		Genre:         engine.GenreFantasy,
		DaysTraveled:  50,
		CrewSurvivors: 3,
		FinalScore:    1500,
		EventsSeen:    []string{"weather_storm", "encounter_bandit"},
		DestinationID: "legendary_city",
		CompletedAt:   time.Now(),
		Victory:       true,
	}

	log.RecordRun(summary)

	if log.TotalRuns != 1 {
		t.Error("total runs should be 1")
	}
	if log.TotalWins != 1 {
		t.Error("total wins should be 1")
	}
	if !log.HasSeenEvent("weather_storm") {
		t.Error("events from run should be recorded")
	}
	if !log.HasReachedDestination("legendary_city") {
		t.Error("destination from victorious run should be recorded")
	}
}

func TestUnlockLogUnlocks(t *testing.T) {
	log := NewUnlockLog()

	// Initially no unlocks
	if len(log.GetUnlockedCrewArchetypes()) != 0 {
		t.Error("no crew archetypes should be unlocked initially")
	}

	// Simulate high progress
	for i := 0; i < 10; i++ {
		summary := &RunSummary{
			Seed:          int64(i),
			Genre:         engine.GenreFantasy,
			DaysTraveled:  50,
			CrewSurvivors: 3,
			FinalScore:    2000,
			EventsSeen:    []string{"event_" + string(rune('a'+i))},
			CompletedAt:   time.Now(),
			Victory:       true,
		}
		log.RecordRun(summary)
	}

	// Should have some unlocks now
	unlocked := log.GetUnlockedCrewArchetypes()
	if len(unlocked) == 0 {
		t.Error("should have some crew archetypes unlocked after progress")
	}
}

func TestHallOfRecords(t *testing.T) {
	hall := NewHallOfRecords()

	// Check all genres initialized
	for _, genre := range engine.AllGenres() {
		entry := hall.GetEntry(genre)
		if entry == nil {
			t.Errorf("genre %s should have entry", genre)
		}
	}
}

func TestHallOfRecordsRecordRun(t *testing.T) {
	hall := NewHallOfRecords()

	summary := &RunSummary{
		Seed:          12345,
		Genre:         engine.GenreFantasy,
		DaysTraveled:  50,
		CrewSurvivors: 3,
		FinalScore:    1500,
		CompletedAt:   time.Now(),
		Victory:       true,
	}

	hall.RecordRun(summary)

	entry := hall.GetEntry(engine.GenreFantasy)
	if entry.TotalRuns != 1 {
		t.Error("total runs should be 1")
	}
	if entry.TotalWins != 1 {
		t.Error("total wins should be 1")
	}
	if entry.BestScore != 1500 {
		t.Errorf("best score should be 1500, got %d", entry.BestScore)
	}
	if entry.BestSummary != summary {
		t.Error("best summary should be set")
	}
}

func TestHallOfRecordsBestRun(t *testing.T) {
	hall := NewHallOfRecords()

	// First run
	summary1 := &RunSummary{
		Genre:      engine.GenreScifi,
		FinalScore: 1000,
		Victory:    true,
	}
	hall.RecordRun(summary1)

	// Better run
	summary2 := &RunSummary{
		Genre:      engine.GenreScifi,
		FinalScore: 2000,
		Victory:    true,
	}
	hall.RecordRun(summary2)

	// Worse run
	summary3 := &RunSummary{
		Genre:      engine.GenreScifi,
		FinalScore: 500,
		Victory:    false,
	}
	hall.RecordRun(summary3)

	best := hall.GetBestRun(engine.GenreScifi)
	if best != summary2 {
		t.Error("best run should be the highest scoring run")
	}

	entry := hall.GetEntry(engine.GenreScifi)
	if entry.TotalRuns != 3 {
		t.Error("total runs should be 3")
	}
	if entry.TotalWins != 2 {
		t.Error("total wins should be 2")
	}
}

func TestHallOfRecordsTotals(t *testing.T) {
	hall := NewHallOfRecords()

	// Add runs to different genres
	genres := engine.AllGenres()
	for _, genre := range genres {
		summary := &RunSummary{
			Genre:      genre,
			FinalScore: 1000,
			Victory:    true,
		}
		hall.RecordRun(summary)
	}

	if hall.GetTotalRuns() != len(genres) {
		t.Errorf("total runs should be %d, got %d", len(genres), hall.GetTotalRuns())
	}
	if hall.GetTotalWins() != len(genres) {
		t.Errorf("total wins should be %d, got %d", len(genres), hall.GetTotalWins())
	}
}

func TestMetaProgress(t *testing.T) {
	meta := NewMetaProgress()

	if meta.UnlockLog == nil {
		t.Error("unlock log should be initialized")
	}
	if meta.HallOfRecords == nil {
		t.Error("hall of records should be initialized")
	}
}

func TestMetaProgressRecordRun(t *testing.T) {
	meta := NewMetaProgress()

	summary := &RunSummary{
		Seed:          12345,
		Genre:         engine.GenreHorror,
		DaysTraveled:  30,
		CrewSurvivors: 2,
		FinalScore:    1200,
		EventsSeen:    []string{"zombie_horde", "supply_find"},
		DestinationID: "safe_zone",
		CompletedAt:   time.Now(),
		Victory:       true,
	}

	meta.RecordRun(summary)

	// Check unlock log
	if meta.UnlockLog.TotalRuns != 1 {
		t.Error("unlock log should have 1 run")
	}
	if !meta.UnlockLog.HasSeenEvent("zombie_horde") {
		t.Error("event should be recorded in unlock log")
	}

	// Check hall of records
	entry := meta.HallOfRecords.GetEntry(engine.GenreHorror)
	if entry.TotalRuns != 1 {
		t.Error("hall of records should have 1 run for horror")
	}
	if entry.BestScore != 1200 {
		t.Error("hall of records should have correct best score")
	}
}

func TestUnlockCategories(t *testing.T) {
	log := NewUnlockLog()

	// Check crew archetypes have correct category
	for _, unlock := range log.CrewArchetypes {
		if unlock.Category != CategoryCrewArchetype {
			t.Errorf("crew archetype %s has wrong category", unlock.ID)
		}
	}

	// Check vessel configs have correct category
	for _, unlock := range log.VesselConfigs {
		if unlock.Category != CategoryVesselConfig {
			t.Errorf("vessel config %s has wrong category", unlock.ID)
		}
	}
}

func TestLockedUnlocks(t *testing.T) {
	log := NewUnlockLog()

	// Initially all should be locked
	lockedCrew := log.GetLockedCrewArchetypes()
	if len(lockedCrew) == 0 {
		t.Error("initially all crew archetypes should be locked")
	}

	lockedVessel := log.GetLockedVesselConfigs()
	if len(lockedVessel) == 0 {
		t.Error("initially all vessel configs should be locked")
	}

	// Total should equal unlocked + locked
	totalCrew := len(log.CrewArchetypes)
	if len(lockedCrew)+len(log.GetUnlockedCrewArchetypes()) != totalCrew {
		t.Error("locked + unlocked should equal total")
	}
}
