package vessel

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewSalvageManager(t *testing.T) {
	sm := NewSalvageManager(12345, engine.GenreFantasy)
	if sm == nil {
		t.Fatal("salvage manager should not be nil")
	}
}

func TestSalvageManagerSetGenre(t *testing.T) {
	sm := NewSalvageManager(12345, engine.GenreFantasy)
	sm.SetGenre(engine.GenreScifi)
	if sm.genre != engine.GenreScifi {
		t.Errorf("genre = %s, want scifi", sm.genre)
	}
}

func TestGenerateSite(t *testing.T) {
	sm := NewSalvageManager(12345, engine.GenreFantasy)

	site := sm.GenerateSite(1)

	if site.ID != 1 {
		t.Errorf("site ID = %d, want 1", site.ID)
	}
	if site.Salvaged {
		t.Error("new site should not be salvaged")
	}
	if site.Richness < 0.3 || site.Richness > 1.0 {
		t.Errorf("richness %f out of expected range [0.3, 1.0]", site.Richness)
	}
	if site.DangerLevel < 0 || site.DangerLevel > 0.6 {
		t.Errorf("danger level %f out of expected range [0, 0.6]", site.DangerLevel)
	}
}

func TestAttemptSalvage(t *testing.T) {
	sm := NewSalvageManager(12345, engine.GenreFantasy)
	site := sm.GenerateSite(1)

	result := sm.AttemptSalvage(site)

	if !result.Success {
		t.Error("salvage should succeed")
	}
	if len(result.Items) == 0 {
		t.Error("should get some items from salvage")
	}
	if !site.Salvaged {
		t.Error("site should be marked as salvaged")
	}
}

func TestAttemptSalvageAlreadySalvaged(t *testing.T) {
	sm := NewSalvageManager(12345, engine.GenreFantasy)
	site := sm.GenerateSite(1)
	site.Salvaged = true

	result := sm.AttemptSalvage(site)

	if result.Success {
		t.Error("salvage should fail for already salvaged site")
	}
	if result.Message != "This site has already been salvaged" {
		t.Errorf("unexpected message: %s", result.Message)
	}
}

func TestSalvageDeterminism(t *testing.T) {
	sm1 := NewSalvageManager(12345, engine.GenreFantasy)
	sm2 := NewSalvageManager(12345, engine.GenreFantasy)

	site1 := sm1.GenerateSite(1)
	site2 := sm2.GenerateSite(1)

	if site1.Type != site2.Type {
		t.Error("same seed should produce same site type")
	}
	if site1.Richness != site2.Richness {
		t.Error("same seed should produce same richness")
	}
}

func TestSalvageItems(t *testing.T) {
	sm := NewSalvageManager(12345, engine.GenreFantasy)
	site := &SalvageSite{
		ID:          1,
		Type:        SalvageWreck,
		Salvaged:    false,
		Richness:    1.0,
		DangerLevel: 0, // No danger
	}

	result := sm.AttemptSalvage(site)

	for _, item := range result.Items {
		if item.Name == "" {
			t.Error("item should have a name")
		}
		if item.Weight <= 0 {
			t.Error("item should have positive weight")
		}
		if item.Quantity <= 0 {
			t.Error("item should have positive quantity")
		}
	}
}

func TestSalvageTypeName(t *testing.T) {
	genres := engine.AllGenres()
	types := AllSalvageTypes()

	for _, g := range genres {
		for _, st := range types {
			name := SalvageTypeName(st, g)
			if name == "" {
				t.Errorf("missing name for genre=%s, type=%d", g, st)
			}
		}
	}
}

func TestAllSalvageTypes(t *testing.T) {
	types := AllSalvageTypes()
	if len(types) != 4 {
		t.Errorf("expected 4 salvage types, got %d", len(types))
	}
}

func TestAddSalvageToHold(t *testing.T) {
	hold := NewCargoHoldWithTier(1)

	items := []SalvageItem{
		{Name: "Scrap", Weight: 2, Volume: 2, Quantity: 5, Category: CargoRepair},
		{Name: "Food", Weight: 1, Volume: 1, Quantity: 10, Category: CargoSupplies},
	}

	added, failed := AddSalvageToHold(items, hold)

	if added != 2 {
		t.Errorf("added = %d, want 2", added)
	}
	if failed != 0 {
		t.Errorf("failed = %d, want 0", failed)
	}
	if hold.GetQuantity("Scrap") != 5 {
		t.Errorf("scrap quantity = %d, want 5", hold.GetQuantity("Scrap"))
	}
	if hold.GetQuantity("Food") != 10 {
		t.Errorf("food quantity = %d, want 10", hold.GetQuantity("Food"))
	}
}

func TestAddSalvageToHoldOverflow(t *testing.T) {
	hold := NewCargoHoldWithTier(1) // Small capacity

	// Try to add more than capacity
	items := []SalvageItem{
		{Name: "Heavy", Weight: 100, Volume: 100, Quantity: 1, Category: CargoRepair},
	}

	added, failed := AddSalvageToHold(items, hold)

	if added != 0 {
		t.Errorf("added = %d, want 0", added)
	}
	if failed != 1 {
		t.Errorf("failed = %d, want 1", failed)
	}
}

func TestSalvageGenreSpecificItems(t *testing.T) {
	genres := engine.AllGenres()

	for _, g := range genres {
		sm := NewSalvageManager(12345, g)
		site := &SalvageSite{
			ID:          1,
			Type:        SalvageWreck,
			Salvaged:    false,
			Richness:    1.0,
			DangerLevel: 0,
		}

		result := sm.AttemptSalvage(site)

		if !result.Success {
			t.Errorf("salvage should succeed for genre %s", g)
		}
		if len(result.Items) == 0 {
			t.Errorf("should get items for genre %s", g)
		}
	}
}

func TestSalvageDangerEncounter(t *testing.T) {
	// Use a seed that produces danger
	sm := NewSalvageManager(99999, engine.GenreHorror)
	site := &SalvageSite{
		ID:          1,
		Type:        SalvageWreck,
		Salvaged:    false,
		Richness:    0.5,
		DangerLevel: 1.0, // Guaranteed danger
	}

	result := sm.AttemptSalvage(site)

	if !result.Success {
		t.Error("salvage should still succeed with danger")
	}
	if !result.Danger {
		t.Error("should flag danger encounter")
	}
	if site.Salvaged != true {
		t.Error("site should still be marked salvaged")
	}
}

func TestSalvageMessages(t *testing.T) {
	genres := engine.AllGenres()

	for _, g := range genres {
		sm := NewSalvageManager(12345, g)

		// Test success messages
		for _, st := range AllSalvageTypes() {
			msg := sm.successMessage(st)
			if msg == "" {
				t.Errorf("missing success message for genre=%s, type=%d", g, st)
			}
		}

		// Test danger message
		msg := sm.dangerMessage()
		if msg == "" {
			t.Errorf("missing danger message for genre=%s", g)
		}
	}
}
