package lore

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewInscription(t *testing.T) {
	insc := NewInscription(1, TypeRuin, 10, 20, "Test Ruin", "Ancient text", engine.GenreFantasy)

	if insc.ID != 1 {
		t.Errorf("expected ID 1, got %d", insc.ID)
	}
	if insc.Type != TypeRuin {
		t.Error("type should be Ruin")
	}
	if insc.X != 10 || insc.Y != 20 {
		t.Error("position incorrect")
	}
	if insc.Discovered {
		t.Error("should not be discovered initially")
	}
}

func TestInscriptionDiscover(t *testing.T) {
	insc := NewInscription(1, TypeGrave, 0, 0, "Test", "Text", engine.GenreFantasy)
	insc.Discover()

	if !insc.Discovered {
		t.Error("should be discovered after Discover()")
	}
}

func TestInscriptionSetGenre(t *testing.T) {
	insc := NewInscription(1, TypeSign, 0, 0, "Test", "Text", engine.GenreFantasy)
	insc.SetGenre(engine.GenreScifi)

	if insc.Genre != engine.GenreScifi {
		t.Error("genre should be updated")
	}
}

func TestNewDiscovery(t *testing.T) {
	disc := NewDiscovery(1, DiscoveryVessel, 10, 20, "Test Ship", "Vignette", engine.GenreScifi)

	if disc.ID != 1 {
		t.Errorf("expected ID 1, got %d", disc.ID)
	}
	if disc.Type != DiscoveryVessel {
		t.Error("type should be Vessel")
	}
	if disc.Discovered || disc.Looted {
		t.Error("should not be discovered or looted initially")
	}
}

func TestDiscoveryItems(t *testing.T) {
	disc := NewDiscovery(1, DiscoveryCamp, 0, 0, "Test", "Vignette", engine.GenreFantasy)
	disc.AddItem("Gold", 10)
	disc.AddItem("Food", 5)

	if len(disc.Items) != 2 {
		t.Errorf("expected 2 items, got %d", len(disc.Items))
	}
	if disc.Items[0].Name != "Gold" || disc.Items[0].Quantity != 10 {
		t.Error("first item incorrect")
	}
}

func TestDiscoveryLoot(t *testing.T) {
	disc := NewDiscovery(1, DiscoveryCache, 0, 0, "Test", "Vignette", engine.GenreFantasy)
	disc.Discover()
	disc.Loot()

	if !disc.Discovered {
		t.Error("should be discovered")
	}
	if !disc.Looted {
		t.Error("should be looted")
	}
}

func TestCodexEntry(t *testing.T) {
	entry := NewCodexEntry("test_1", CodexHistory, "Test Title", "Test Text", engine.GenreFantasy)

	if entry.ID != "test_1" {
		t.Error("ID should match")
	}
	if entry.Category != CodexHistory {
		t.Error("category should be History")
	}
	if entry.Unlocked {
		t.Error("should not be unlocked initially")
	}
}

func TestCodexEntryUnlock(t *testing.T) {
	entry := NewCodexEntry("test_1", CodexFaction, "Test", "Text", engine.GenreFantasy)
	entry.Unlock("exploration")

	if !entry.Unlocked {
		t.Error("should be unlocked")
	}
	if entry.UnlockSource != "exploration" {
		t.Error("unlock source should be set")
	}
}

func TestCodex(t *testing.T) {
	codex := NewCodex(engine.GenreFantasy)

	entry1 := NewCodexEntry("h1", CodexHistory, "History 1", "Text", engine.GenreFantasy)
	entry2 := NewCodexEntry("f1", CodexFaction, "Faction 1", "Text", engine.GenreFantasy)

	codex.AddEntry(entry1)
	codex.AddEntry(entry2)

	if codex.GetEntry("h1") != entry1 {
		t.Error("should retrieve entry h1")
	}
	if codex.TotalCount() != 2 {
		t.Error("total count should be 2")
	}
	if codex.UnlockedCount() != 0 {
		t.Error("unlocked count should be 0")
	}
}

func TestCodexUnlock(t *testing.T) {
	codex := NewCodex(engine.GenreFantasy)
	entry := NewCodexEntry("test", CodexRoute, "Test", "Text", engine.GenreFantasy)
	codex.AddEntry(entry)

	if !codex.UnlockEntry("test", "event") {
		t.Error("should successfully unlock entry")
	}
	if codex.UnlockedCount() != 1 {
		t.Error("unlocked count should be 1")
	}

	// Can't unlock again
	if codex.UnlockEntry("test", "event") {
		t.Error("should not unlock already unlocked entry")
	}
}

func TestCodexByCategory(t *testing.T) {
	codex := NewCodex(engine.GenreFantasy)

	codex.AddEntry(NewCodexEntry("h1", CodexHistory, "H1", "T", engine.GenreFantasy))
	codex.AddEntry(NewCodexEntry("h2", CodexHistory, "H2", "T", engine.GenreFantasy))
	codex.AddEntry(NewCodexEntry("f1", CodexFaction, "F1", "T", engine.GenreFantasy))

	history := codex.GetEntriesByCategory(CodexHistory)
	if len(history) != 2 {
		t.Errorf("expected 2 history entries, got %d", len(history))
	}

	factions := codex.GetEntriesByCategory(CodexFaction)
	if len(factions) != 1 {
		t.Errorf("expected 1 faction entry, got %d", len(factions))
	}
}

func TestEnvironmentalManager(t *testing.T) {
	manager := NewEnvironmentalManager(engine.GenreFantasy)

	insc := NewInscription(1, TypeRuin, 10, 10, "Test", "Text", engine.GenreFantasy)
	disc := NewDiscovery(1, DiscoveryVessel, 20, 20, "Test", "Vignette", engine.GenreFantasy)

	manager.AddInscription(insc)
	manager.AddDiscovery(disc)

	foundInsc := manager.GetInscriptionAt(10, 10)
	if foundInsc != insc {
		t.Error("should find inscription at position")
	}

	foundDisc := manager.GetDiscoveryAt(20, 20)
	if foundDisc != disc {
		t.Error("should find discovery at position")
	}
}

func TestEnvironmentalManagerDiscoverAt(t *testing.T) {
	manager := NewEnvironmentalManager(engine.GenreFantasy)

	insc := NewInscription(1, TypeSign, 10, 10, "Test", "Text", engine.GenreFantasy)
	disc := NewDiscovery(1, DiscoveryCamp, 10, 10, "Test", "Vignette", engine.GenreFantasy)

	manager.AddInscription(insc)
	manager.AddDiscovery(disc)

	foundInsc, foundDisc := manager.DiscoverAt(10, 10)

	if foundInsc == nil || !foundInsc.Discovered {
		t.Error("inscription should be discovered")
	}
	if foundDisc == nil || !foundDisc.Discovered {
		t.Error("discovery should be discovered")
	}
}

func TestGenerator(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)

	insc := g.GenerateInscription(10, 20)
	if insc.Title == "" {
		t.Error("inscription should have title")
	}
	if insc.Text == "" {
		t.Error("inscription should have text")
	}

	disc := g.GenerateDiscovery(30, 40)
	if disc.Title == "" {
		t.Error("discovery should have title")
	}
	if disc.VignetteText == "" {
		t.Error("discovery should have vignette")
	}
	if len(disc.Items) == 0 {
		t.Error("discovery should have items")
	}
}

func TestGeneratorDeterminism(t *testing.T) {
	g1 := NewGenerator(12345, engine.GenreFantasy)
	g2 := NewGenerator(12345, engine.GenreFantasy)

	insc1 := g1.GenerateInscription(10, 10)
	insc2 := g2.GenerateInscription(10, 10)

	if insc1.Title != insc2.Title {
		t.Error("same seed should produce same titles")
	}
	if insc1.Text != insc2.Text {
		t.Error("same seed should produce same text")
	}
}

func TestGeneratorSetGenre(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)
	g.SetGenre(engine.GenreCyberpunk)

	insc := g.GenerateInscription(10, 10)
	if insc.Genre != engine.GenreCyberpunk {
		t.Error("inscription should have cyberpunk genre")
	}
}

func TestAllGenreLoreGeneration(t *testing.T) {
	genres := engine.AllGenres()

	for _, genre := range genres {
		g := NewGenerator(12345, genre)

		insc := g.GenerateInscription(10, 10)
		if insc.Title == "" {
			t.Errorf("genre %s: inscription should have title", genre)
		}

		disc := g.GenerateDiscovery(20, 20)
		if disc.Title == "" {
			t.Errorf("genre %s: discovery should have title", genre)
		}

		codex := g.GenerateCodexEntry(CodexHistory, "test")
		if codex.Title == "" {
			t.Errorf("genre %s: codex should have title", genre)
		}
	}
}

func TestGenerateEnvironment(t *testing.T) {
	g := NewGenerator(12345, engine.GenreFantasy)

	manager := g.GenerateEnvironment(100, 100, 10, 5)

	if len(manager.Inscriptions) != 10 {
		t.Errorf("expected 10 inscriptions, got %d", len(manager.Inscriptions))
	}
	if len(manager.Discoveries) != 5 {
		t.Errorf("expected 5 discoveries, got %d", len(manager.Discoveries))
	}
	if manager.Codex.TotalCount() == 0 {
		t.Error("codex should have entries")
	}
}

func TestInscriptionTypeNames(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		for _, iType := range AllInscriptionTypes() {
			name := InscriptionTypeName(iType, genre)
			if name == "" {
				t.Errorf("inscription type %v genre %s should have name", iType, genre)
			}
		}
	}
}

func TestDiscoveryTypeNames(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		for _, dType := range AllDiscoveryTypes() {
			name := DiscoveryTypeName(dType, genre)
			if name == "" {
				t.Errorf("discovery type %v genre %s should have name", dType, genre)
			}
		}
	}
}

func TestCodexCategoryNames(t *testing.T) {
	for _, cat := range AllCodexCategories() {
		name := CodexCategoryName(cat)
		if name == "" || name == "Unknown" {
			t.Errorf("codex category %v should have name", cat)
		}
	}
}

func TestEnvironmentalManagerSetGenre(t *testing.T) {
	manager := NewEnvironmentalManager(engine.GenreFantasy)
	insc := NewInscription(1, TypeRuin, 10, 10, "Test", "Text", engine.GenreFantasy)
	disc := NewDiscovery(1, DiscoveryVessel, 20, 20, "Test", "Vignette", engine.GenreFantasy)
	manager.AddInscription(insc)
	manager.AddDiscovery(disc)

	manager.SetGenre(engine.GenreScifi)

	if insc.Genre != engine.GenreScifi {
		t.Error("inscription genre should be updated")
	}
	if disc.Genre != engine.GenreScifi {
		t.Error("discovery genre should be updated")
	}
}
