package trading

import "github.com/opd-ai/voyage/pkg/engine"

// TradingPostVocab holds genre-specific trading post UI vocabulary.
type TradingPostVocab struct {
	// Location names
	LocationTitle   string
	WelcomeMessage  string
	FarewellMessage string

	// Actions
	EnterAction  string
	LeaveAction  string
	BrowseAction string
	TradeAction  string
	BarterAction string

	// Merchant
	MerchantTitle   string
	GreetingMessage string

	// Status
	OpenStatus       string
	ClosedStatus     string
	RestrictedStatus string

	// Inventory
	InventoryLabel string
	StockLabel     string
	PriceLabel     string
	QualityLabel   string

	// Transaction
	BuyLabel     string
	SellLabel    string
	TotalLabel   string
	BalanceLabel string

	// Special offers
	DealLabel    string
	RareLabel    string
	LimitedLabel string
}

// GetTradingPostVocab returns genre-specific trading post vocabulary.
func GetTradingPostVocab(genre engine.GenreID) *TradingPostVocab {
	vocab, ok := tradingPostVocabs[genre]
	if !ok {
		return tradingPostVocabs[engine.GenreFantasy]
	}
	return vocab
}

var tradingPostVocabs = map[engine.GenreID]*TradingPostVocab{
	engine.GenreFantasy: {
		LocationTitle:    "Trading Market",
		WelcomeMessage:   "Welcome, traveler!",
		FarewellMessage:  "Safe travels!",
		EnterAction:      "Enter",
		LeaveAction:      "Depart",
		BrowseAction:     "Browse Wares",
		TradeAction:      "Trade",
		BarterAction:     "Barter",
		MerchantTitle:    "Merchant",
		GreetingMessage:  "What can I interest you in today?",
		OpenStatus:       "Open for Trade",
		ClosedStatus:     "Closed",
		RestrictedStatus: "Restricted",
		InventoryLabel:   "Wares",
		StockLabel:       "Stock",
		PriceLabel:       "Price",
		QualityLabel:     "Quality",
		BuyLabel:         "Purchase",
		SellLabel:        "Sell",
		TotalLabel:       "Total",
		BalanceLabel:     "Gold",
		DealLabel:        "Special Offer",
		RareLabel:        "Rare",
		LimitedLabel:     "Limited Stock",
	},
	engine.GenreScifi: {
		LocationTitle:    "Space Dock Commerce Hub",
		WelcomeMessage:   "Docking complete. Welcome aboard.",
		FarewellMessage:  "Safe journey through the void.",
		EnterAction:      "Dock",
		LeaveAction:      "Undock",
		BrowseAction:     "Access Catalog",
		TradeAction:      "Transact",
		BarterAction:     "Exchange Protocol",
		MerchantTitle:    "Vendor",
		GreetingMessage:  "Browse our inventory database.",
		OpenStatus:       "Systems Online",
		ClosedStatus:     "Systems Offline",
		RestrictedStatus: "Access Restricted",
		InventoryLabel:   "Inventory",
		StockLabel:       "Units",
		PriceLabel:       "Credits",
		QualityLabel:     "Grade",
		BuyLabel:         "Purchase",
		SellLabel:        "Liquidate",
		TotalLabel:       "Total",
		BalanceLabel:     "Credits",
		DealLabel:        "Priority Offer",
		RareLabel:        "Premium",
		LimitedLabel:     "Low Stock",
	},
	engine.GenreHorror: {
		LocationTitle:    "Survivor Camp",
		WelcomeMessage:   "You made it. Come in, quick.",
		FarewellMessage:  "Be careful out there.",
		EnterAction:      "Enter",
		LeaveAction:      "Leave",
		BrowseAction:     "See What's Available",
		TradeAction:      "Trade",
		BarterAction:     "Barter",
		MerchantTitle:    "Trader",
		GreetingMessage:  "Take what you need. We all need to survive.",
		OpenStatus:       "Trading",
		ClosedStatus:     "Closed - Too Dangerous",
		RestrictedStatus: "Not Trusted",
		InventoryLabel:   "Supplies",
		StockLabel:       "Left",
		PriceLabel:       "Cost",
		QualityLabel:     "Condition",
		BuyLabel:         "Take",
		SellLabel:        "Give",
		TotalLabel:       "Total",
		BalanceLabel:     "Supplies",
		DealLabel:        "Urgent Need",
		RareLabel:        "Rare Find",
		LimitedLabel:     "Almost Gone",
	},
	engine.GenreCyberpunk: {
		LocationTitle:    "Black Market Terminal",
		WelcomeMessage:   "Connection established. Welcome, user.",
		FarewellMessage:  "Connection terminated.",
		EnterAction:      "Connect",
		LeaveAction:      "Disconnect",
		BrowseAction:     "Browse Listings",
		TradeAction:      "Execute Transaction",
		BarterAction:     "Exchange Protocol",
		MerchantTitle:    "Fixer",
		GreetingMessage:  "What are you looking for?",
		OpenStatus:       "Online",
		ClosedStatus:     "Offline",
		RestrictedStatus: "Access Denied",
		InventoryLabel:   "Listings",
		StockLabel:       "Available",
		PriceLabel:       "Creds",
		QualityLabel:     "Rating",
		BuyLabel:         "Download",
		SellLabel:        "Upload",
		TotalLabel:       "Total",
		BalanceLabel:     "Creds",
		DealLabel:        "Hot Deal",
		RareLabel:        "Illegal",
		LimitedLabel:     "Limited",
	},
	engine.GenrePostapoc: {
		LocationTitle:    "Scrap Bazaar",
		WelcomeMessage:   "Welcome to the bazaar, stranger.",
		FarewellMessage:  "Don't get killed out there.",
		EnterAction:      "Enter",
		LeaveAction:      "Leave",
		BrowseAction:     "Look Around",
		TradeAction:      "Trade",
		BarterAction:     "Swap",
		MerchantTitle:    "Trader",
		GreetingMessage:  "Got some good stuff today.",
		OpenStatus:       "Open",
		ClosedStatus:     "Closed",
		RestrictedStatus: "Stay Out",
		InventoryLabel:   "Goods",
		StockLabel:       "Count",
		PriceLabel:       "Scrap",
		QualityLabel:     "Shape",
		BuyLabel:         "Buy",
		SellLabel:        "Sell",
		TotalLabel:       "Total",
		BalanceLabel:     "Scrap",
		DealLabel:        "Good Deal",
		RareLabel:        "Pre-War",
		LimitedLabel:     "Running Low",
	},
}

// PostDescriptionVocab provides genre-specific post description templates.
type PostDescriptionVocab struct {
	MarketDesc      string
	OutpostDesc     string
	SpecialistDesc  string
	BlackMarketDesc string
}

// GetPostDescriptionVocab returns genre-specific post description vocabulary.
func GetPostDescriptionVocab(genre engine.GenreID) *PostDescriptionVocab {
	vocab, ok := postDescriptionVocabs[genre]
	if !ok {
		return postDescriptionVocabs[engine.GenreFantasy]
	}
	return vocab
}

var postDescriptionVocabs = map[engine.GenreID]*PostDescriptionVocab{
	engine.GenreFantasy: {
		MarketDesc:      "A bustling marketplace filled with traders from distant lands",
		OutpostDesc:     "A frontier trading post serving weary travelers",
		SpecialistDesc:  "A guild shop specializing in crafted wares",
		BlackMarketDesc: "A shadowy den where forbidden goods change hands",
	},
	engine.GenreScifi: {
		MarketDesc:      "A commercial hub at the center of interstellar trade routes",
		OutpostDesc:     "A relay station offering basic supplies to passing ships",
		SpecialistDesc:  "A tech bay stocked with specialized equipment",
		BlackMarketDesc: "An off-grid market dealing in restricted technology",
	},
	engine.GenreHorror: {
		MarketDesc:      "A survivor camp where the desperate trade what they can",
		OutpostDesc:     "A fortified safe house offering refuge and supplies",
		SpecialistDesc:  "A specialist with rare knowledge and rarer goods",
		BlackMarketDesc: "A place where dangerous items find new owners",
	},
	engine.GenreCyberpunk: {
		MarketDesc:      "A corporate-controlled marketplace with monitored transactions",
		OutpostDesc:     "A data node offering network access and basic supplies",
		SpecialistDesc:  "A mod shop for those seeking enhancement",
		BlackMarketDesc: "An encrypted market where anonymity is currency",
	},
	engine.GenrePostapoc: {
		MarketDesc:      "A scrap bazaar where salvagers trade their finds",
		OutpostDesc:     "A fortified post offering shelter and trade",
		SpecialistDesc:  "A workshop where skilled hands repair the old world's tech",
		BlackMarketDesc: "An underground market dealing in dangerous goods",
	},
}

// SetGenreForPost updates a supply post's vocabulary when genre changes.
// This allows dynamic re-skinning of existing posts.
func SetGenreForPost(post *SupplyPost, newGenre engine.GenreID) {
	oldGenre := post.Genre
	post.Genre = newGenre

	// Update inventory item genre
	for _, item := range post.Inventory.Items {
		item.Genre = newGenre
	}

	// The name and description were generated with the old genre,
	// but TypeName() will automatically use the new genre vocabulary.
	// For full re-skinning, regenerate name components if needed.
	_ = oldGenre // May be used for logging or transition effects
}
