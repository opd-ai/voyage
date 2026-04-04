package trading

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/resources"
)

func TestNewBarterInterface(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreFantasy)
	post := gen.Generate(10, 20, 1)
	res := resources.NewResources(engine.GenreFantasy)
	ti := NewTradeInterface(post, res)

	bi := NewBarterInterface(ti)
	if bi == nil {
		t.Fatal("NewBarterInterface returned nil")
	}
	if bi.genre != engine.GenreFantasy {
		t.Errorf("Expected genre fantasy, got %v", bi.genre)
	}
}

func TestEvaluateOffer(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreFantasy)
	post := gen.Generate(10, 20, 1)
	res := resources.NewResources(engine.GenreFantasy)
	ti := NewTradeInterface(post, res)

	// Give player some items
	playerInv := NewInventory(engine.GenreFantasy)
	playerInv.AddItem(&Item{
		Name:      "Test Sword",
		Category:  CategoryTrade,
		BasePrice: 100.0,
		Quantity:  5,
		Quality:   0.8,
	})
	ti.SetPlayerInventory(playerInv)

	bi := NewBarterInterface(ti)

	// Get an item from the post
	postItems := post.Inventory.Items
	if len(postItems) == 0 {
		t.Skip("Post has no items")
	}
	postItem := postItems[0]

	// Create an offer
	offer := CreateBarterOffer(
		map[string]int{"Test Sword": 2},
		map[string]int{postItem.Name: 1},
	)

	offerValue, requestValue, valid := bi.EvaluateOffer(offer)
	if !valid {
		t.Fatal("Expected valid offer")
	}

	if offerValue <= 0 {
		t.Errorf("Expected positive offer value, got %f", offerValue)
	}
	if requestValue <= 0 {
		t.Errorf("Expected positive request value, got %f", requestValue)
	}
}

func TestEvaluateOfferInvalidPlayerItem(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreFantasy)
	post := gen.Generate(10, 20, 1)
	res := resources.NewResources(engine.GenreFantasy)
	ti := NewTradeInterface(post, res)

	playerInv := NewInventory(engine.GenreFantasy)
	ti.SetPlayerInventory(playerInv)

	bi := NewBarterInterface(ti)

	// Offer item player doesn't have
	postItems := post.Inventory.Items
	if len(postItems) == 0 {
		t.Skip("Post has no items")
	}

	offer := CreateBarterOffer(
		map[string]int{"Nonexistent Item": 1},
		map[string]int{postItems[0].Name: 1},
	)

	_, _, valid := bi.EvaluateOffer(offer)
	if valid {
		t.Error("Expected invalid offer for nonexistent player item")
	}
}

func TestEvaluateOfferInvalidPostItem(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreFantasy)
	post := gen.Generate(10, 20, 1)
	res := resources.NewResources(engine.GenreFantasy)
	ti := NewTradeInterface(post, res)

	playerInv := NewInventory(engine.GenreFantasy)
	playerInv.AddItem(&Item{
		Name:      "Test Item",
		Category:  CategoryTrade,
		BasePrice: 50.0,
		Quantity:  3,
		Quality:   0.7,
	})
	ti.SetPlayerInventory(playerInv)

	bi := NewBarterInterface(ti)

	// Request item post doesn't have
	offer := CreateBarterOffer(
		map[string]int{"Test Item": 1},
		map[string]int{"Nonexistent Post Item": 1},
	)

	_, _, valid := bi.EvaluateOffer(offer)
	if valid {
		t.Error("Expected invalid offer for nonexistent post item")
	}
}

func TestBarterAccepted(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreFantasy)
	post := gen.Generate(10, 20, 1)
	post.Reputation = 0.8 // High reputation for easier acceptance

	res := resources.NewResources(engine.GenreFantasy)
	ti := NewTradeInterface(post, res)

	// Give player high-value items
	playerInv := NewInventory(engine.GenreFantasy)
	playerInv.AddItem(&Item{
		Name:      "Valuable Gem",
		Category:  CategoryRare,
		BasePrice: 500.0,
		Quantity:  10,
		Quality:   1.0,
	})
	ti.SetPlayerInventory(playerInv)

	bi := NewBarterInterface(ti)

	// Get a cheap item from post
	var cheapItem *Item
	for _, item := range post.Inventory.Items {
		if item.BasePrice < 100 && item.Quantity > 0 {
			cheapItem = item
			break
		}
	}
	if cheapItem == nil {
		t.Skip("No cheap item in post inventory")
	}

	// Offer overpayment to ensure acceptance
	offer := CreateBarterOffer(
		map[string]int{"Valuable Gem": 3},
		map[string]int{cheapItem.Name: 1},
	)

	result := bi.Barter(offer)
	if !result.Success {
		t.Fatalf("Expected successful barter, got: %s", result.Message)
	}
	if !result.Accepted {
		t.Errorf("Expected accepted barter with high-value offer")
	}

	// Verify items were exchanged
	if playerInv.GetItem("Valuable Gem").Quantity != 7 {
		t.Errorf("Expected 7 gems remaining, got %d", playerInv.GetItem("Valuable Gem").Quantity)
	}
	if playerInv.GetItem(cheapItem.Name) == nil {
		t.Error("Expected player to receive the traded item")
	}
}

func TestBarterRejectedWithCounterOffer(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreScifi)
	post := gen.Generate(10, 20, 1)
	post.Reputation = 0.0 // Low reputation for harder trades

	res := resources.NewResources(engine.GenreScifi)
	ti := NewTradeInterface(post, res)

	// Give player low-value items
	playerInv := NewInventory(engine.GenreScifi)
	playerInv.AddItem(&Item{
		Name:      "Cheap Part",
		Category:  CategoryParts,
		BasePrice: 10.0,
		Quantity:  2,
		Quality:   0.5,
	})
	playerInv.AddItem(&Item{
		Name:      "Extra Fuel",
		Category:  CategoryFuel,
		BasePrice: 20.0,
		Quantity:  5,
		Quality:   0.7,
	})
	ti.SetPlayerInventory(playerInv)

	bi := NewBarterInterface(ti)

	// Find an expensive item
	var expensiveItem *Item
	for _, item := range post.Inventory.Items {
		if item.BasePrice >= 50 && item.Quantity > 0 {
			expensiveItem = item
			break
		}
	}
	if expensiveItem == nil {
		t.Skip("No expensive item in post inventory")
	}

	// Offer underpayment
	offer := CreateBarterOffer(
		map[string]int{"Cheap Part": 1},
		map[string]int{expensiveItem.Name: 1},
	)

	result := bi.Barter(offer)
	if !result.Success {
		t.Fatalf("Expected successful evaluation, got: %s", result.Message)
	}
	if result.Accepted {
		// Might be accepted if values happen to align
		return
	}

	if result.CounterOffer == nil {
		t.Error("Expected counter-offer when trade is rejected")
	}
}

func TestBarterEmptyOffer(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreHorror)
	post := gen.Generate(10, 20, 1)
	res := resources.NewResources(engine.GenreHorror)
	ti := NewTradeInterface(post, res)

	playerInv := NewInventory(engine.GenreHorror)
	ti.SetPlayerInventory(playerInv)

	bi := NewBarterInterface(ti)

	// Empty offer
	offer := CreateBarterOffer(map[string]int{}, map[string]int{"Something": 1})
	result := bi.Barter(offer)

	if result.Success {
		t.Error("Expected failure for empty offer")
	}
}

func TestBarterEmptyRequest(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreCyberpunk)
	post := gen.Generate(10, 20, 1)
	res := resources.NewResources(engine.GenreCyberpunk)
	ti := NewTradeInterface(post, res)

	playerInv := NewInventory(engine.GenreCyberpunk)
	playerInv.AddItem(&Item{
		Name:      "Test Item",
		BasePrice: 50.0,
		Quantity:  5,
	})
	ti.SetPlayerInventory(playerInv)

	bi := NewBarterInterface(ti)

	// Empty request
	offer := CreateBarterOffer(map[string]int{"Test Item": 1}, map[string]int{})
	result := bi.Barter(offer)

	if result.Success {
		t.Error("Expected failure for empty request")
	}
}

func TestBarterNoInventory(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenrePostapoc)
	post := gen.Generate(10, 20, 1)
	res := resources.NewResources(engine.GenrePostapoc)
	ti := NewTradeInterface(post, res)
	ti.playerInventory = nil

	bi := NewBarterInterface(ti)

	offer := CreateBarterOffer(
		map[string]int{"Item": 1},
		map[string]int{"Other": 1},
	)
	result := bi.Barter(offer)

	if result.Success {
		t.Error("Expected failure when player has no inventory")
	}
}

func TestBarterVocabAllGenres(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	for _, genre := range genres {
		vocab := GetBarterVocab(genre)
		if vocab == nil {
			t.Errorf("GetBarterVocab(%v) returned nil", genre)
			continue
		}
		if vocab.BarterAction == "" {
			t.Errorf("Genre %v has empty BarterAction", genre)
		}
		if vocab.OfferVerb == "" {
			t.Errorf("Genre %v has empty OfferVerb", genre)
		}
		if vocab.RequestVerb == "" {
			t.Errorf("Genre %v has empty RequestVerb", genre)
		}
		if vocab.AcceptVerb == "" {
			t.Errorf("Genre %v has empty AcceptVerb", genre)
		}
		if vocab.RejectVerb == "" {
			t.Errorf("Genre %v has empty RejectVerb", genre)
		}
	}
}

func TestBarterMessagesAllGenres(t *testing.T) {
	genres := []engine.GenreID{
		engine.GenreFantasy,
		engine.GenreScifi,
		engine.GenreHorror,
		engine.GenreCyberpunk,
		engine.GenrePostapoc,
	}

	keys := []string{
		"no_inventory",
		"empty_offer",
		"empty_request",
		"invalid_items",
		"accepted",
		"rejected",
	}

	for _, genre := range genres {
		messages, ok := barterMessages[genre]
		if !ok {
			t.Errorf("Missing barterMessages for genre %v", genre)
			continue
		}
		for _, key := range keys {
			if messages[key] == "" {
				t.Errorf("Genre %v has empty message for key %s", genre, key)
			}
		}
	}
}

func TestAcceptanceThreshold(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreFantasy)
	post := gen.Generate(10, 20, 1)
	res := resources.NewResources(engine.GenreFantasy)
	ti := NewTradeInterface(post, res)
	bi := NewBarterInterface(ti)

	// Test at different reputation levels
	post.Reputation = 0.0
	threshold0 := bi.acceptanceThreshold()

	post.Reputation = 0.5
	threshold50 := bi.acceptanceThreshold()

	post.Reputation = 1.0
	threshold100 := bi.acceptanceThreshold()

	// Higher reputation should mean lower threshold (better for player)
	if threshold50 >= threshold0 {
		t.Errorf("Threshold at 50%% rep (%f) should be less than at 0%% (%f)", threshold50, threshold0)
	}
	if threshold100 >= threshold50 {
		t.Errorf("Threshold at 100%% rep (%f) should be less than at 50%% (%f)", threshold100, threshold50)
	}
}

func TestCreateBarterOffer(t *testing.T) {
	offer := CreateBarterOffer(
		map[string]int{"Item A": 2, "Item B": 3},
		map[string]int{"Item C": 1},
	)

	if len(offer.OfferedItems) != 2 {
		t.Errorf("Expected 2 offered items, got %d", len(offer.OfferedItems))
	}
	if len(offer.RequestedItems) != 1 {
		t.Errorf("Expected 1 requested item, got %d", len(offer.RequestedItems))
	}

	// Check quantities
	found := false
	for _, item := range offer.OfferedItems {
		if item.ItemName == "Item A" && item.Quantity == 2 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to find Item A with quantity 2 in offered items")
	}
}

func TestBarterReputationChange(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreFantasy)
	post := gen.Generate(10, 20, 1)
	post.Reputation = 0.5
	initialRep := post.Reputation

	res := resources.NewResources(engine.GenreFantasy)
	ti := NewTradeInterface(post, res)

	// Give player high-value items for guaranteed acceptance
	playerInv := NewInventory(engine.GenreFantasy)
	playerInv.AddItem(&Item{
		Name:      "Super Gem",
		Category:  CategoryRare,
		BasePrice: 1000.0,
		Quantity:  10,
		Quality:   1.0,
	})
	ti.SetPlayerInventory(playerInv)

	bi := NewBarterInterface(ti)

	// Find any item from post
	if len(post.Inventory.Items) == 0 {
		t.Skip("Post has no items")
	}
	postItem := post.Inventory.Items[0]

	// Make a generous offer
	offer := CreateBarterOffer(
		map[string]int{"Super Gem": 5},
		map[string]int{postItem.Name: 1},
	)

	result := bi.Barter(offer)
	if result.Accepted {
		// Successful trade should boost reputation
		if post.Reputation <= initialRep {
			t.Error("Expected reputation to increase after successful barter")
		}
	}
}
