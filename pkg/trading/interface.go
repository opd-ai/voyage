package trading

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/resources"
)

// TradeAction represents a buy or sell action.
type TradeAction int

const (
	// ActionBuy represents purchasing from the supply post.
	ActionBuy TradeAction = iota
	// ActionSell represents selling to the supply post.
	ActionSell
)

// TradeResult represents the outcome of a trade transaction.
type TradeResult struct {
	Success   bool
	Action    TradeAction
	ItemName  string
	Quantity  int
	UnitPrice float64
	TotalCost float64
	Message   string
}

// TradeInterface manages buying and selling at a supply post.
type TradeInterface struct {
	post            *SupplyPost
	playerResources *resources.Resources
	playerInventory *Inventory
	genre           engine.GenreID
}

// NewTradeInterface creates a new trading interface for a supply post.
func NewTradeInterface(post *SupplyPost, playerResources *resources.Resources) *TradeInterface {
	return &TradeInterface{
		post:            post,
		playerResources: playerResources,
		playerInventory: NewInventory(post.Genre),
		genre:           post.Genre,
	}
}

// SetPlayerInventory sets the player's inventory for trading.
func (ti *TradeInterface) SetPlayerInventory(inv *Inventory) {
	ti.playerInventory = inv
}

// GetAvailableItems returns items available for purchase at this post.
func (ti *TradeInterface) GetAvailableItems() []*TradeOffer {
	offers := make([]*TradeOffer, 0)
	for _, item := range ti.post.Inventory.Items {
		if item.Quantity <= 0 {
			continue
		}
		price := ti.post.AdjustedPrice(item.BasePrice, false)
		offers = append(offers, &TradeOffer{
			Item:      item,
			Action:    ActionBuy,
			UnitPrice: price,
			Available: item.Quantity,
		})
	}
	return offers
}

// GetSellableItems returns items the player can sell at this post.
func (ti *TradeInterface) GetSellableItems() []*TradeOffer {
	if ti.playerInventory == nil {
		return nil
	}
	offers := make([]*TradeOffer, 0)
	for _, item := range ti.playerInventory.Items {
		if item.Quantity <= 0 {
			continue
		}
		price := ti.post.AdjustedPrice(item.BasePrice, true)
		offers = append(offers, &TradeOffer{
			Item:      item,
			Action:    ActionSell,
			UnitPrice: price,
			Available: item.Quantity,
		})
	}
	return offers
}

// TradeOffer represents an item available for trade.
type TradeOffer struct {
	Item      *Item
	Action    TradeAction
	UnitPrice float64
	Available int
}

// TotalPrice returns the total price for a given quantity.
func (to *TradeOffer) TotalPrice(quantity int) float64 {
	return to.UnitPrice * float64(quantity)
}

// Buy attempts to purchase an item from the supply post.
func (ti *TradeInterface) Buy(itemName string, quantity int) TradeResult {
	// Find the item
	item := ti.post.Inventory.GetItem(itemName)
	if item == nil {
		return TradeResult{
			Success:  false,
			Action:   ActionBuy,
			ItemName: itemName,
			Message:  "Item not available",
		}
	}

	// Check quantity
	if item.Quantity < quantity {
		return TradeResult{
			Success:  false,
			Action:   ActionBuy,
			ItemName: itemName,
			Quantity: quantity,
			Message:  "Insufficient stock",
		}
	}

	// Calculate price
	unitPrice := ti.post.AdjustedPrice(item.BasePrice, false)
	totalCost := unitPrice * float64(quantity)

	// Check player currency
	currency := ti.playerResources.Get(resources.ResourceCurrency)
	if currency < totalCost {
		return TradeResult{
			Success:   false,
			Action:    ActionBuy,
			ItemName:  itemName,
			Quantity:  quantity,
			UnitPrice: unitPrice,
			TotalCost: totalCost,
			Message:   "Insufficient currency",
		}
	}

	// Execute transaction
	ti.playerResources.Consume(resources.ResourceCurrency, totalCost)
	ti.post.Inventory.RemoveItem(itemName, quantity)

	// Add to player inventory (create copy to avoid aliasing issues)
	boughtItem := &Item{
		Name:        item.Name,
		Description: item.Description,
		Category:    item.Category,
		BasePrice:   item.BasePrice,
		Quantity:    quantity,
		Quality:     item.Quality,
		Genre:       item.Genre,
	}
	ti.playerInventory.AddItem(boughtItem)

	// Apply resource effect if applicable
	ti.applyResourceEffect(item.Category, quantity, item.Quality)

	// Small reputation boost for trading
	ti.post.UpdateReputation(0.01)

	return TradeResult{
		Success:   true,
		Action:    ActionBuy,
		ItemName:  itemName,
		Quantity:  quantity,
		UnitPrice: unitPrice,
		TotalCost: totalCost,
		Message:   ti.buyMessage(itemName, quantity),
	}
}

// Sell attempts to sell an item to the supply post.
func (ti *TradeInterface) Sell(itemName string, quantity int) TradeResult {
	if ti.playerInventory == nil {
		return TradeResult{
			Success:  false,
			Action:   ActionSell,
			ItemName: itemName,
			Message:  "No player inventory",
		}
	}

	// Find the item
	item := ti.playerInventory.GetItem(itemName)
	if item == nil {
		return TradeResult{
			Success:  false,
			Action:   ActionSell,
			ItemName: itemName,
			Message:  "Item not in inventory",
		}
	}

	// Check quantity
	if item.Quantity < quantity {
		return TradeResult{
			Success:  false,
			Action:   ActionSell,
			ItemName: itemName,
			Quantity: quantity,
			Message:  "Insufficient quantity to sell",
		}
	}

	// Calculate price
	unitPrice := ti.post.AdjustedPrice(item.BasePrice, true)
	totalValue := unitPrice * float64(quantity)

	// Execute transaction
	ti.playerInventory.RemoveItem(itemName, quantity)
	ti.playerResources.Add(resources.ResourceCurrency, totalValue)

	// Add to post inventory
	soldItem := &Item{
		Name:        item.Name,
		Description: item.Description,
		Category:    item.Category,
		BasePrice:   item.BasePrice,
		Quantity:    quantity,
		Quality:     item.Quality,
		Genre:       item.Genre,
	}
	ti.post.Inventory.AddItem(soldItem)

	// Small reputation boost for trading
	ti.post.UpdateReputation(0.01)

	return TradeResult{
		Success:   true,
		Action:    ActionSell,
		ItemName:  itemName,
		Quantity:  quantity,
		UnitPrice: unitPrice,
		TotalCost: totalValue,
		Message:   ti.sellMessage(itemName, quantity),
	}
}

// applyResourceEffect adds resources based on item category.
func (ti *TradeInterface) applyResourceEffect(category ItemCategory, quantity int, quality float64) {
	amount := float64(quantity) * quality * 10

	switch category {
	case CategoryFood:
		ti.playerResources.Add(resources.ResourceFood, amount)
	case CategoryWater:
		ti.playerResources.Add(resources.ResourceWater, amount)
	case CategoryFuel:
		ti.playerResources.Add(resources.ResourceFuel, amount)
	case CategoryMedicine:
		ti.playerResources.Add(resources.ResourceMedicine, amount)
	}
}

// buyMessage generates a genre-appropriate purchase message.
func (ti *TradeInterface) buyMessage(itemName string, quantity int) string {
	messages := map[engine.GenreID]string{
		engine.GenreFantasy:   "You have acquired %d %s",
		engine.GenreScifi:     "Transaction complete: %d %s transferred",
		engine.GenreHorror:    "You grabbed %d %s",
		engine.GenreCyberpunk: "Download complete: %d %s",
		engine.GenrePostapoc:  "You scored %d %s",
	}
	template := messages[ti.genre]
	if template == "" {
		template = "Purchased %d %s"
	}
	return formatQuantityMessage(template, quantity, itemName)
}

// sellMessage generates a genre-appropriate sell message.
func (ti *TradeInterface) sellMessage(itemName string, quantity int) string {
	messages := map[engine.GenreID]string{
		engine.GenreFantasy:   "You sold %d %s",
		engine.GenreScifi:     "Transfer complete: %d %s sold",
		engine.GenreHorror:    "You traded away %d %s",
		engine.GenreCyberpunk: "Upload complete: %d %s sold",
		engine.GenrePostapoc:  "You offloaded %d %s",
	}
	template := messages[ti.genre]
	if template == "" {
		template = "Sold %d %s"
	}
	return formatQuantityMessage(template, quantity, itemName)
}

// formatQuantityMessage formats a message with quantity and item name.
func formatQuantityMessage(template string, quantity int, itemName string) string {
	// Simple placeholder replacement
	result := template
	quantityPlaced := false
	for i := 0; i < len(result)-1; i++ {
		if result[i] == '%' && result[i+1] == 'd' && !quantityPlaced {
			// Replace %d with quantity
			prefix := result[:i]
			suffix := result[i+2:]
			result = prefix + intToString(quantity) + suffix
			quantityPlaced = true
			continue
		}
		if result[i] == '%' && result[i+1] == 's' {
			// Replace %s with item name
			prefix := result[:i]
			suffix := result[i+2:]
			result = prefix + itemName + suffix
			break
		}
	}
	return result
}

// intToString converts an int to string without importing strconv.
func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	negative := n < 0
	if negative {
		n = -n
	}
	digits := make([]byte, 0, 10)
	for n > 0 {
		digits = append(digits, byte('0'+n%10))
		n /= 10
	}
	// Reverse
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	if negative {
		return "-" + string(digits)
	}
	return string(digits)
}

// CanAfford checks if the player can afford a purchase.
func (ti *TradeInterface) CanAfford(itemName string, quantity int) bool {
	item := ti.post.Inventory.GetItem(itemName)
	if item == nil {
		return false
	}
	totalCost := ti.post.AdjustedPrice(item.BasePrice, false) * float64(quantity)
	return ti.playerResources.Get(resources.ResourceCurrency) >= totalCost
}

// GetPlayerCurrency returns the player's current currency.
func (ti *TradeInterface) GetPlayerCurrency() float64 {
	return ti.playerResources.Get(resources.ResourceCurrency)
}

// GetPostInventoryValue returns the total value of items at the post.
func (ti *TradeInterface) GetPostInventoryValue() float64 {
	return ti.post.Inventory.TotalValue()
}

// GetPlayerInventoryValue returns the total sell value of player's items.
func (ti *TradeInterface) GetPlayerInventoryValue() float64 {
	if ti.playerInventory == nil {
		return 0
	}
	total := 0.0
	for _, item := range ti.playerInventory.Items {
		price := ti.post.AdjustedPrice(item.BasePrice, true)
		total += price * float64(item.Quantity)
	}
	return total
}

// TradeVocab holds genre-specific trading vocabulary.
type TradeVocab struct {
	BuyAction            string
	SellAction           string
	CurrencyName         string
	InsufficientCurrency string
	InsufficientStock    string
	TradeComplete        string
	NoStock              string
	BrowsePrompt         string
	SellPrompt           string
}

// GetTradeVocab returns genre-specific trading vocabulary.
func GetTradeVocab(genre engine.GenreID) *TradeVocab {
	vocab, ok := tradeVocabs[genre]
	if !ok {
		return tradeVocabs[engine.GenreFantasy]
	}
	return vocab
}

var tradeVocabs = map[engine.GenreID]*TradeVocab{
	engine.GenreFantasy: {
		BuyAction:            "Purchase",
		SellAction:           "Sell",
		CurrencyName:         "Gold",
		InsufficientCurrency: "You haven't enough gold",
		InsufficientStock:    "The merchant hasn't enough stock",
		TradeComplete:        "Trade complete",
		NoStock:              "Nothing for sale",
		BrowsePrompt:         "What would you like to buy?",
		SellPrompt:           "What would you like to sell?",
	},
	engine.GenreScifi: {
		BuyAction:            "Purchase",
		SellAction:           "Sell",
		CurrencyName:         "Credits",
		InsufficientCurrency: "Insufficient credits",
		InsufficientStock:    "Item out of stock",
		TradeComplete:        "Transaction complete",
		NoStock:              "No inventory available",
		BrowsePrompt:         "Browse available inventory",
		SellPrompt:           "Select items to liquidate",
	},
	engine.GenreHorror: {
		BuyAction:            "Take",
		SellAction:           "Trade",
		CurrencyName:         "Supplies",
		InsufficientCurrency: "Not enough to trade",
		InsufficientStock:    "They don't have enough",
		TradeComplete:        "Deal done",
		NoStock:              "Nothing to trade",
		BrowsePrompt:         "What do you need?",
		SellPrompt:           "What can you spare?",
	},
	engine.GenreCyberpunk: {
		BuyAction:            "Download",
		SellAction:           "Upload",
		CurrencyName:         "Creds",
		InsufficientCurrency: "Insufficient creds",
		InsufficientStock:    "Item unavailable",
		TradeComplete:        "Transaction logged",
		NoStock:              "Inventory empty",
		BrowsePrompt:         "Browse the catalog",
		SellPrompt:           "Liquidate assets",
	},
	engine.GenrePostapoc: {
		BuyAction:            "Barter",
		SellAction:           "Trade",
		CurrencyName:         "Scrap",
		InsufficientCurrency: "Not enough scrap",
		InsufficientStock:    "They're out",
		TradeComplete:        "Fair trade",
		NoStock:              "Shelves are empty",
		BrowsePrompt:         "What've they got?",
		SellPrompt:           "What can you give up?",
	},
}
