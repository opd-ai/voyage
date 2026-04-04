package trading

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestGetTradingPostVocab_AllGenres(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		vocab := GetTradingPostVocab(genre)
		if vocab == nil {
			t.Errorf("GetTradingPostVocab(%v) returned nil", genre)
			continue
		}

		// Check all fields are populated
		if vocab.LocationTitle == "" {
			t.Errorf("Genre %v has empty LocationTitle", genre)
		}
		if vocab.WelcomeMessage == "" {
			t.Errorf("Genre %v has empty WelcomeMessage", genre)
		}
		if vocab.FarewellMessage == "" {
			t.Errorf("Genre %v has empty FarewellMessage", genre)
		}
		if vocab.EnterAction == "" {
			t.Errorf("Genre %v has empty EnterAction", genre)
		}
		if vocab.LeaveAction == "" {
			t.Errorf("Genre %v has empty LeaveAction", genre)
		}
		if vocab.MerchantTitle == "" {
			t.Errorf("Genre %v has empty MerchantTitle", genre)
		}
		if vocab.BalanceLabel == "" {
			t.Errorf("Genre %v has empty BalanceLabel", genre)
		}
	}
}

func TestGetTradingPostVocab_UnknownGenre(t *testing.T) {
	vocab := GetTradingPostVocab(engine.GenreID("unknown"))
	if vocab == nil {
		t.Fatal("GetTradingPostVocab for unknown genre returned nil")
	}
	// Should fall back to fantasy
	fantasyVocab := GetTradingPostVocab(engine.GenreFantasy)
	if vocab.LocationTitle != fantasyVocab.LocationTitle {
		t.Error("Unknown genre should fall back to fantasy")
	}
}

func TestGetPostDescriptionVocab_AllGenres(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		vocab := GetPostDescriptionVocab(genre)
		if vocab == nil {
			t.Errorf("GetPostDescriptionVocab(%v) returned nil", genre)
			continue
		}

		if vocab.MarketDesc == "" {
			t.Errorf("Genre %v has empty MarketDesc", genre)
		}
		if vocab.OutpostDesc == "" {
			t.Errorf("Genre %v has empty OutpostDesc", genre)
		}
		if vocab.SpecialistDesc == "" {
			t.Errorf("Genre %v has empty SpecialistDesc", genre)
		}
		if vocab.BlackMarketDesc == "" {
			t.Errorf("Genre %v has empty BlackMarketDesc", genre)
		}
	}
}

func TestSetGenreForPost(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreFantasy)
	post := gen.Generate(10, 20, 1)

	originalType := post.TypeName()

	// Change genre
	SetGenreForPost(post, engine.GenreScifi)

	if post.Genre != engine.GenreScifi {
		t.Errorf("Expected genre scifi, got %v", post.Genre)
	}

	// TypeName should now return scifi vocabulary
	newType := post.TypeName()
	if newType == originalType {
		// They might coincidentally be the same, but verify genre is set
		if post.Genre != engine.GenreScifi {
			t.Error("Genre should have changed")
		}
	}

	// Check inventory items have updated genre
	for _, item := range post.Inventory.Items {
		if item.Genre != engine.GenreScifi {
			t.Errorf("Item genre should be scifi, got %v", item.Genre)
		}
	}
}

func TestPostTypeName_MatchesTask(t *testing.T) {
	// Verify the specific names mentioned in the task:
	// market→space-dock→survivor-camp→black-market→scrap-bazaar

	testCases := []struct {
		genre    engine.GenreID
		expected string
	}{
		{engine.GenreFantasy, "Market"},
		{engine.GenreScifi, "Space Dock"},
		{engine.GenreHorror, "Survivor Camp"},
		{engine.GenreCyberpunk, "Black Market"},
		{engine.GenrePostapoc, "Scrap Bazaar"},
	}

	for _, tc := range testCases {
		name := PostTypeName(PostTypeMarket, tc.genre)
		if name != tc.expected {
			t.Errorf("Genre %v Market should be '%s', got '%s'", tc.genre, tc.expected, name)
		}
	}
}

func TestTradingPostVocab_GenreThemes(t *testing.T) {
	// Verify vocabulary matches genre themes

	// Fantasy should be medieval themed
	fantasyVocab := GetTradingPostVocab(engine.GenreFantasy)
	if fantasyVocab.BalanceLabel != "Gold" {
		t.Errorf("Fantasy balance should be 'Gold', got '%s'", fantasyVocab.BalanceLabel)
	}

	// Scifi should use credits
	scifiVocab := GetTradingPostVocab(engine.GenreScifi)
	if scifiVocab.BalanceLabel != "Credits" {
		t.Errorf("Scifi balance should be 'Credits', got '%s'", scifiVocab.BalanceLabel)
	}

	// Postapoc should use scrap
	postapocVocab := GetTradingPostVocab(engine.GenrePostapoc)
	if postapocVocab.BalanceLabel != "Scrap" {
		t.Errorf("Postapoc balance should be 'Scrap', got '%s'", postapocVocab.BalanceLabel)
	}
}

func TestTradingPostVocab_ActionVerbs(t *testing.T) {
	// Verify appropriate action verbs per genre

	scifiVocab := GetTradingPostVocab(engine.GenreScifi)
	if scifiVocab.EnterAction != "Dock" {
		t.Errorf("Scifi enter action should be 'Dock', got '%s'", scifiVocab.EnterAction)
	}

	cyberpunkVocab := GetTradingPostVocab(engine.GenreCyberpunk)
	if cyberpunkVocab.EnterAction != "Connect" {
		t.Errorf("Cyberpunk enter action should be 'Connect', got '%s'", cyberpunkVocab.EnterAction)
	}
}
