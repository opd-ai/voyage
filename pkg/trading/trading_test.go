package trading

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/resources"
)

func TestNewSupplyPostGenerator(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreFantasy)
	if gen == nil {
		t.Fatal("expected non-nil generator")
	}
	if gen.Genre() != engine.GenreFantasy {
		t.Errorf("expected fantasy genre, got %s", gen.Genre())
	}
}

func TestSupplyPostGenerator_SetGenre(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreFantasy)
	gen.SetGenre(engine.GenreScifi)
	if gen.Genre() != engine.GenreScifi {
		t.Errorf("expected scifi genre, got %s", gen.Genre())
	}
}

func TestSupplyPostGenerator_Generate(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreFantasy)
	post := gen.Generate(10, 20, 1)

	if post == nil {
		t.Fatal("expected non-nil supply post")
	}
	if post.Name == "" {
		t.Error("expected non-empty name")
	}
	if post.Description == "" {
		t.Error("expected non-empty description")
	}
	if post.Position[0] != 10 || post.Position[1] != 20 {
		t.Errorf("expected position [10,20], got %v", post.Position)
	}
	if post.RegionID != 1 {
		t.Errorf("expected region ID 1, got %d", post.RegionID)
	}
	if post.Inventory == nil {
		t.Error("expected non-nil inventory")
	}
	if post.PriceModifier <= 0 {
		t.Error("expected positive price modifier")
	}
}

func TestSupplyPostGenerator_Determinism(t *testing.T) {
	gen1 := NewSupplyPostGenerator(42, engine.GenreHorror)
	gen2 := NewSupplyPostGenerator(42, engine.GenreHorror)

	post1 := gen1.Generate(5, 5, 1)
	post2 := gen2.Generate(5, 5, 1)

	if post1.Name != post2.Name {
		t.Errorf("same seed should produce same name: %s vs %s", post1.Name, post2.Name)
	}
	if post1.PostType != post2.PostType {
		t.Error("same seed should produce same post type")
	}
}

func TestSupplyPostGenerator_AllGenres(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		gen := NewSupplyPostGenerator(99999, genre)
		post := gen.Generate(0, 0, 0)

		if post.Genre != genre {
			t.Errorf("genre mismatch: expected %s, got %s", genre, post.Genre)
		}
		if post.Name == "" {
			t.Errorf("genre %s produced empty name", genre)
		}
	}
}

func TestSupplyPost_AdjustedPrice(t *testing.T) {
	post := &SupplyPost{
		PriceModifier: 1.0,
		Reputation:    0.5, // Neutral
	}

	// Test buying (price modifier applies)
	price := post.AdjustedPrice(100, false)
	if price != 100 {
		t.Errorf("expected 100, got %f", price)
	}

	// Test selling (60% of price)
	sellPrice := post.AdjustedPrice(100, true)
	if sellPrice != 60 {
		t.Errorf("expected 60 for selling, got %f", sellPrice)
	}

	// Test with high reputation
	post.Reputation = 1.0
	betterPrice := post.AdjustedPrice(100, false)
	if betterPrice >= 100 {
		t.Error("high reputation should give better (lower) buy prices")
	}
}

func TestSupplyPost_UpdateReputation(t *testing.T) {
	post := &SupplyPost{Reputation: 0.5}

	post.UpdateReputation(0.1)
	if post.Reputation != 0.6 {
		t.Errorf("expected 0.6, got %f", post.Reputation)
	}

	// Test clamping
	post.UpdateReputation(1.0)
	if post.Reputation != 1.0 {
		t.Errorf("reputation should clamp at 1.0, got %f", post.Reputation)
	}

	post.UpdateReputation(-2.0)
	if post.Reputation != 0.0 {
		t.Errorf("reputation should clamp at 0.0, got %f", post.Reputation)
	}
}

func TestSupplyPost_ReputationStatus(t *testing.T) {
	tests := []struct {
		rep    float64
		status string
	}{
		{0.9, "Honored"},
		{0.7, "Friendly"},
		{0.5, "Neutral"},
		{0.3, "Suspicious"},
		{0.1, "Hostile"},
	}

	for _, tc := range tests {
		post := &SupplyPost{Reputation: tc.rep}
		if post.ReputationStatus() != tc.status {
			t.Errorf("rep %f: expected %s, got %s", tc.rep, tc.status, post.ReputationStatus())
		}
	}
}

func TestPostTypeName_AllGenres(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		for _, postType := range AllPostTypes() {
			name := PostTypeName(postType, genre)
			if name == "" {
				t.Errorf("empty name for type %v in genre %s", postType, genre)
			}
		}
	}
}

func TestNewInventory(t *testing.T) {
	inv := NewInventory(engine.GenreFantasy)
	if inv == nil {
		t.Fatal("expected non-nil inventory")
	}
	if inv.ItemCount() != 0 {
		t.Error("expected empty inventory")
	}
}

func TestInventory_AddItem(t *testing.T) {
	inv := NewInventory(engine.GenreFantasy)
	item := &Item{Name: "Test Item", Category: CategoryFood, Quantity: 5}

	inv.AddItem(item)
	if inv.ItemCount() != 1 {
		t.Errorf("expected 1 item, got %d", inv.ItemCount())
	}

	// Adding same item should merge quantities
	item2 := &Item{Name: "Test Item", Category: CategoryFood, Quantity: 3}
	inv.AddItem(item2)
	if inv.ItemCount() != 1 {
		t.Errorf("expected 1 item after merge, got %d", inv.ItemCount())
	}
	if inv.Items[0].Quantity != 8 {
		t.Errorf("expected quantity 8, got %d", inv.Items[0].Quantity)
	}
}

func TestInventory_RemoveItem(t *testing.T) {
	inv := NewInventory(engine.GenreFantasy)
	inv.AddItem(&Item{Name: "Test Item", Category: CategoryFood, Quantity: 5})

	// Remove some
	if !inv.RemoveItem("Test Item", 3) {
		t.Error("expected successful removal")
	}
	if inv.Items[0].Quantity != 2 {
		t.Errorf("expected quantity 2, got %d", inv.Items[0].Quantity)
	}

	// Remove too many
	if inv.RemoveItem("Test Item", 5) {
		t.Error("expected failure removing more than available")
	}

	// Remove all remaining
	if !inv.RemoveItem("Test Item", 2) {
		t.Error("expected successful removal")
	}
	if inv.ItemCount() != 0 {
		t.Error("expected empty inventory after removing all")
	}
}

func TestInventory_GetItem(t *testing.T) {
	inv := NewInventory(engine.GenreFantasy)
	inv.AddItem(&Item{Name: "Test Item", Category: CategoryFood, Quantity: 5})

	item := inv.GetItem("Test Item")
	if item == nil {
		t.Error("expected to find item")
	}

	missing := inv.GetItem("Missing Item")
	if missing != nil {
		t.Error("expected nil for missing item")
	}
}

func TestInventory_GetByCategory(t *testing.T) {
	inv := NewInventory(engine.GenreFantasy)
	inv.AddItem(&Item{Name: "Food 1", Category: CategoryFood, Quantity: 1})
	inv.AddItem(&Item{Name: "Food 2", Category: CategoryFood, Quantity: 1})
	inv.AddItem(&Item{Name: "Water 1", Category: CategoryWater, Quantity: 1})

	foods := inv.GetByCategory(CategoryFood)
	if len(foods) != 2 {
		t.Errorf("expected 2 food items, got %d", len(foods))
	}

	waters := inv.GetByCategory(CategoryWater)
	if len(waters) != 1 {
		t.Errorf("expected 1 water item, got %d", len(waters))
	}
}

func TestInventory_TotalValue(t *testing.T) {
	inv := NewInventory(engine.GenreFantasy)
	inv.AddItem(&Item{Name: "Item 1", BasePrice: 10, Quantity: 2})
	inv.AddItem(&Item{Name: "Item 2", BasePrice: 5, Quantity: 3})

	total := inv.TotalValue()
	expected := 10.0*2 + 5.0*3
	if total != expected {
		t.Errorf("expected %f, got %f", expected, total)
	}
}

func TestNewItemGenerator(t *testing.T) {
	gen := NewItemGenerator(12345, engine.GenreFantasy)
	if gen == nil {
		t.Fatal("expected non-nil generator")
	}
}

func TestItemGenerator_Generate(t *testing.T) {
	gen := NewItemGenerator(12345, engine.GenreFantasy)
	item := gen.Generate(PostTypeMarket)

	if item == nil {
		t.Fatal("expected non-nil item")
	}
	if item.Name == "" {
		t.Error("expected non-empty name")
	}
	if item.Description == "" {
		t.Error("expected non-empty description")
	}
	if item.BasePrice <= 0 {
		t.Error("expected positive price")
	}
	if item.Quantity <= 0 {
		t.Error("expected positive quantity")
	}
	if item.Quality < 0 || item.Quality > 1 {
		t.Errorf("quality out of range: %f", item.Quality)
	}
}

func TestItemGenerator_Determinism(t *testing.T) {
	gen1 := NewItemGenerator(42, engine.GenreHorror)
	gen2 := NewItemGenerator(42, engine.GenreHorror)

	item1 := gen1.Generate(PostTypeMarket)
	item2 := gen2.Generate(PostTypeMarket)

	if item1.Name != item2.Name {
		t.Errorf("same seed should produce same name: %s vs %s", item1.Name, item2.Name)
	}
	if item1.Category != item2.Category {
		t.Error("same seed should produce same category")
	}
}

func TestCategoryName_AllGenres(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		for _, category := range AllItemCategories() {
			name := CategoryName(category, genre)
			if name == "" {
				t.Errorf("empty name for category %v in genre %s", category, genre)
			}
		}
	}
}

func TestAllPostTypes(t *testing.T) {
	types := AllPostTypes()
	if len(types) != 4 {
		t.Errorf("expected 4 post types, got %d", len(types))
	}
}

func TestAllItemCategories(t *testing.T) {
	categories := AllItemCategories()
	if len(categories) != 7 {
		t.Errorf("expected 7 categories, got %d", len(categories))
	}
}

// TradeInterface tests

func TestNewTradeInterface(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreFantasy)
	post := gen.Generate(0, 0, 0)
	playerRes := resources.NewResources(engine.GenreFantasy)
	playerRes.Set(resources.ResourceCurrency, 100)

	ti := NewTradeInterface(post, playerRes)
	if ti == nil {
		t.Fatal("expected non-nil trade interface")
	}
}

func TestTradeInterface_GetAvailableItems(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreFantasy)
	post := gen.Generate(0, 0, 0)
	playerRes := resources.NewResources(engine.GenreFantasy)

	ti := NewTradeInterface(post, playerRes)
	offers := ti.GetAvailableItems()

	if len(offers) == 0 {
		t.Error("expected some items available")
	}
	for _, offer := range offers {
		if offer.UnitPrice <= 0 {
			t.Error("expected positive price")
		}
		if offer.Available <= 0 {
			t.Error("expected positive availability")
		}
	}
}

func TestTradeInterface_Buy(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreFantasy)
	post := gen.Generate(0, 0, 0)
	playerRes := resources.NewResources(engine.GenreFantasy)
	playerRes.Set(resources.ResourceCurrency, 1000)

	ti := NewTradeInterface(post, playerRes)

	// Get first available item
	offers := ti.GetAvailableItems()
	if len(offers) == 0 {
		t.Skip("no items available to test")
	}
	itemName := offers[0].Item.Name

	// Test successful purchase
	result := ti.Buy(itemName, 1)
	if !result.Success {
		t.Errorf("expected successful purchase: %s", result.Message)
	}
	if result.TotalCost <= 0 {
		t.Error("expected positive cost")
	}

	// Verify currency was spent
	newCurrency := playerRes.Get(resources.ResourceCurrency)
	if newCurrency >= 1000 {
		t.Error("expected currency to decrease after purchase")
	}

	// Test purchase of non-existent item
	result = ti.Buy("Nonexistent Item", 1)
	if result.Success {
		t.Error("expected failure for non-existent item")
	}
}

func TestTradeInterface_Buy_InsufficientFunds(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreFantasy)
	post := gen.Generate(0, 0, 0)
	playerRes := resources.NewResources(engine.GenreFantasy)
	playerRes.Set(resources.ResourceCurrency, 0) // No money

	ti := NewTradeInterface(post, playerRes)

	offers := ti.GetAvailableItems()
	if len(offers) == 0 {
		t.Skip("no items available to test")
	}
	itemName := offers[0].Item.Name

	result := ti.Buy(itemName, 1)
	if result.Success {
		t.Error("expected failure due to insufficient funds")
	}
}

func TestTradeInterface_Sell(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreFantasy)
	post := gen.Generate(0, 0, 0)
	playerRes := resources.NewResources(engine.GenreFantasy)
	playerRes.Set(resources.ResourceCurrency, 0)

	ti := NewTradeInterface(post, playerRes)

	// Give player an item to sell
	playerInv := NewInventory(engine.GenreFantasy)
	playerInv.AddItem(&Item{
		Name:      "Test Item",
		BasePrice: 10,
		Quantity:  5,
		Quality:   0.8,
		Genre:     engine.GenreFantasy,
	})
	ti.SetPlayerInventory(playerInv)

	// Test successful sale
	result := ti.Sell("Test Item", 2)
	if !result.Success {
		t.Errorf("expected successful sale: %s", result.Message)
	}
	if result.TotalCost <= 0 {
		t.Error("expected positive value")
	}

	// Verify currency was gained
	newCurrency := playerRes.Get(resources.ResourceCurrency)
	if newCurrency <= 0 {
		t.Error("expected currency to increase after sale")
	}

	// Verify inventory decreased
	item := playerInv.GetItem("Test Item")
	if item.Quantity != 3 {
		t.Errorf("expected 3 remaining, got %d", item.Quantity)
	}
}

func TestTradeInterface_CanAfford(t *testing.T) {
	gen := NewSupplyPostGenerator(12345, engine.GenreFantasy)
	post := gen.Generate(0, 0, 0)
	playerRes := resources.NewResources(engine.GenreFantasy)

	ti := NewTradeInterface(post, playerRes)

	offers := ti.GetAvailableItems()
	if len(offers) == 0 {
		t.Skip("no items available to test")
	}
	itemName := offers[0].Item.Name

	// With no money
	playerRes.Set(resources.ResourceCurrency, 0)
	if ti.CanAfford(itemName, 1) {
		t.Error("expected cannot afford with no currency")
	}

	// With plenty of money
	playerRes.Set(resources.ResourceCurrency, 10000)
	if !ti.CanAfford(itemName, 1) {
		t.Error("expected can afford with lots of currency")
	}
}

func TestGetTradeVocab_AllGenres(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		vocab := GetTradeVocab(genre)
		if vocab == nil {
			t.Errorf("nil vocab for genre %s", genre)
			continue
		}
		if vocab.BuyAction == "" {
			t.Errorf("empty BuyAction for genre %s", genre)
		}
		if vocab.SellAction == "" {
			t.Errorf("empty SellAction for genre %s", genre)
		}
		if vocab.CurrencyName == "" {
			t.Errorf("empty CurrencyName for genre %s", genre)
		}
	}
}

func TestTradeOffer_TotalPrice(t *testing.T) {
	offer := &TradeOffer{
		UnitPrice: 10,
		Available: 5,
	}

	total := offer.TotalPrice(3)
	if total != 30 {
		t.Errorf("expected 30, got %f", total)
	}
}
