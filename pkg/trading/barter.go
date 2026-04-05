package trading

import "github.com/opd-ai/voyage/pkg/engine"

// BarterOffer represents a proposed goods-for-goods trade.
type BarterOffer struct {
	OfferedItems   []*BarterItem // Items the player is offering
	RequestedItems []*BarterItem // Items the player wants
}

// BarterItem represents an item and quantity in a barter offer.
type BarterItem struct {
	ItemName string
	Quantity int
}

// BarterResult represents the outcome of a barter attempt.
type BarterResult struct {
	Success      bool
	Accepted     bool
	OfferValue   float64
	RequestValue float64
	CounterOffer *BarterOffer // Non-nil if merchant makes counter-offer
	Message      string
}

// BarterInterface handles goods-for-goods trading.
type BarterInterface struct {
	tradeInterface *TradeInterface
	genre          engine.GenreID
}

// NewBarterInterface creates a barter interface from an existing trade interface.
func NewBarterInterface(ti *TradeInterface) *BarterInterface {
	return &BarterInterface{
		tradeInterface: ti,
		genre:          ti.genre,
	}
}

// EvaluateOffer calculates the value of a barter offer from the merchant's perspective.
func (bi *BarterInterface) EvaluateOffer(offer *BarterOffer) (offerValue, requestValue float64, valid bool) {
	ti := bi.tradeInterface

	// Calculate value of items being offered (what player gives)
	for _, item := range offer.OfferedItems {
		playerItem := ti.playerInventory.GetItem(item.ItemName)
		if playerItem == nil || playerItem.Quantity < item.Quantity {
			return 0, 0, false // Invalid offer
		}
		// Value at sell price (what merchant would pay)
		unitPrice := ti.post.AdjustedPrice(playerItem.BasePrice, true)
		offerValue += unitPrice * float64(item.Quantity)
	}

	// Calculate value of items being requested (what player wants)
	for _, item := range offer.RequestedItems {
		postItem := ti.post.Inventory.GetItem(item.ItemName)
		if postItem == nil || postItem.Quantity < item.Quantity {
			return 0, 0, false // Invalid request
		}
		// Value at buy price (what player would pay)
		unitPrice := ti.post.AdjustedPrice(postItem.BasePrice, false)
		requestValue += unitPrice * float64(item.Quantity)
	}

	return offerValue, requestValue, true
}

// Barter attempts to execute a goods-for-goods trade.
func (bi *BarterInterface) Barter(offer *BarterOffer) BarterResult {
	if err := bi.validateBarterOffer(offer); err != nil {
		return BarterResult{Success: false, Message: err.Error()}
	}

	offerValue, requestValue, valid := bi.EvaluateOffer(offer)
	if !valid {
		return BarterResult{Success: false, Message: bi.getMessage("invalid_items")}
	}

	return bi.processBarterOffer(offer, offerValue, requestValue)
}

// validateBarterOffer checks if the barter offer is valid.
func (bi *BarterInterface) validateBarterOffer(offer *BarterOffer) error {
	if bi.tradeInterface.playerInventory == nil {
		return barterValidationError(bi.getMessage("no_inventory"))
	}
	if len(offer.OfferedItems) == 0 {
		return barterValidationError(bi.getMessage("empty_offer"))
	}
	if len(offer.RequestedItems) == 0 {
		return barterValidationError(bi.getMessage("empty_request"))
	}
	return nil
}

// barterValidationError is a simple error type for barter validation.
type barterValidationError string

func (e barterValidationError) Error() string { return string(e) }

// processBarterOffer processes a validated barter offer.
func (bi *BarterInterface) processBarterOffer(offer *BarterOffer, offerValue, requestValue float64) BarterResult {
	result := BarterResult{
		Success:      true,
		OfferValue:   offerValue,
		RequestValue: requestValue,
	}

	threshold := bi.acceptanceThreshold()
	if offerValue >= requestValue*threshold {
		bi.acceptBarter(&result, offer)
	} else {
		bi.rejectBarter(&result, offer, offerValue, requestValue)
	}

	return result
}

// acceptBarter handles an accepted barter offer.
func (bi *BarterInterface) acceptBarter(result *BarterResult, offer *BarterOffer) {
	result.Accepted = true
	bi.executeBarter(offer)
	result.Message = bi.getMessage("accepted")
	bi.tradeInterface.post.UpdateReputation(0.02)
}

// rejectBarter handles a rejected barter offer with counter-offer.
func (bi *BarterInterface) rejectBarter(result *BarterResult, offer *BarterOffer, offerValue, requestValue float64) {
	result.Accepted = false
	result.CounterOffer = bi.generateCounterOffer(offer, offerValue, requestValue)
	result.Message = bi.getMessage("rejected")
	bi.tradeInterface.post.UpdateReputation(-0.005)
}

// acceptanceThreshold returns the value ratio needed for acceptance.
// Lower values = merchant accepts worse deals (better for player).
func (bi *BarterInterface) acceptanceThreshold() float64 {
	rep := bi.tradeInterface.post.Reputation
	// Reputation 0 = 1.1 threshold (merchant wants 10% more value)
	// Reputation 1 = 0.85 threshold (merchant accepts 15% less value)
	return 1.1 - (rep * 0.25)
}

// executeBarter performs the actual item exchange.
func (bi *BarterInterface) executeBarter(offer *BarterOffer) {
	ti := bi.tradeInterface

	// Remove offered items from player, add to post
	for _, item := range offer.OfferedItems {
		playerItem := ti.playerInventory.GetItem(item.ItemName)
		if playerItem == nil {
			continue
		}

		ti.playerInventory.RemoveItem(item.ItemName, item.Quantity)

		tradedItem := &Item{
			Name:        playerItem.Name,
			Description: playerItem.Description,
			Category:    playerItem.Category,
			BasePrice:   playerItem.BasePrice,
			Quantity:    item.Quantity,
			Quality:     playerItem.Quality,
			Genre:       playerItem.Genre,
		}
		ti.post.Inventory.AddItem(tradedItem)
	}

	// Remove requested items from post, add to player
	for _, item := range offer.RequestedItems {
		postItem := ti.post.Inventory.GetItem(item.ItemName)
		if postItem == nil {
			continue
		}

		ti.post.Inventory.RemoveItem(item.ItemName, item.Quantity)

		receivedItem := &Item{
			Name:        postItem.Name,
			Description: postItem.Description,
			Category:    postItem.Category,
			BasePrice:   postItem.BasePrice,
			Quantity:    item.Quantity,
			Quality:     postItem.Quality,
			Genre:       postItem.Genre,
		}
		ti.playerInventory.AddItem(receivedItem)

		// Apply resource effects
		ti.applyResourceEffect(postItem.Category, item.Quantity, postItem.Quality)
	}
}

// generateCounterOffer creates a counter-offer when original is rejected.
func (bi *BarterInterface) generateCounterOffer(original *BarterOffer, offerValue, requestValue float64) *BarterOffer {
	deficit := bi.calculateDeficit(offerValue, requestValue)
	if deficit <= 0 {
		return nil
	}

	counter := bi.copyOriginalOffer(original)
	bi.addItemsToMeetDeficit(counter, original, &deficit)
	return counter
}

// calculateDeficit returns how much more value the merchant needs.
func (bi *BarterInterface) calculateDeficit(offerValue, requestValue float64) float64 {
	threshold := bi.acceptanceThreshold()
	needed := requestValue * threshold
	return needed - offerValue
}

// copyOriginalOffer creates a deep copy of the original barter offer.
func (bi *BarterInterface) copyOriginalOffer(original *BarterOffer) *BarterOffer {
	counter := &BarterOffer{
		OfferedItems:   make([]*BarterItem, len(original.OfferedItems)),
		RequestedItems: make([]*BarterItem, len(original.RequestedItems)),
	}
	for i, item := range original.OfferedItems {
		counter.OfferedItems[i] = &BarterItem{
			ItemName: item.ItemName,
			Quantity: item.Quantity,
		}
	}
	for i, item := range original.RequestedItems {
		counter.RequestedItems[i] = &BarterItem{
			ItemName: item.ItemName,
			Quantity: item.Quantity,
		}
	}
	return counter
}

// addItemsToMeetDeficit finds additional items from player inventory to cover the deficit.
func (bi *BarterInterface) addItemsToMeetDeficit(counter, original *BarterOffer, deficit *float64) {
	ti := bi.tradeInterface
	for _, item := range ti.playerInventory.Items {
		if *deficit <= 0 {
			break
		}
		if bi.shouldSkipItem(item, original.OfferedItems) {
			continue
		}
		bi.addItemToCounter(counter, item, deficit)
	}
}

// shouldSkipItem returns true if the item should be skipped during counter-offer generation.
func (bi *BarterInterface) shouldSkipItem(item *Item, offeredItems []*BarterItem) bool {
	return bi.isInBarterOffer(item.Name, offeredItems) || item.Quantity <= 0
}

// addItemToCounter adds an item to the counter-offer if it has positive value.
func (bi *BarterInterface) addItemToCounter(counter *BarterOffer, item *Item, deficit *float64) {
	ti := bi.tradeInterface
	unitPrice := ti.post.AdjustedPrice(item.BasePrice, true)
	if unitPrice <= 0 {
		return
	}

	quantityNeeded := bi.calculateQuantityNeeded(*deficit, unitPrice, item.Quantity)
	if quantityNeeded > 0 {
		counter.OfferedItems = append(counter.OfferedItems, &BarterItem{
			ItemName: item.Name,
			Quantity: quantityNeeded,
		})
		*deficit -= unitPrice * float64(quantityNeeded)
	}
}

// calculateQuantityNeeded determines how many items are needed to cover a deficit.
func (bi *BarterInterface) calculateQuantityNeeded(deficit, unitPrice float64, available int) int {
	quantityNeeded := int(deficit/unitPrice) + 1
	if quantityNeeded > available {
		quantityNeeded = available
	}
	return quantityNeeded
}

// isInBarterOffer checks if an item is already in a barter list.
func (bi *BarterInterface) isInBarterOffer(itemName string, items []*BarterItem) bool {
	for _, item := range items {
		if item.ItemName == itemName {
			return true
		}
	}
	return false
}

// getMessage returns a genre-appropriate barter message.
func (bi *BarterInterface) getMessage(key string) string {
	messages, ok := barterMessages[bi.genre]
	if !ok {
		messages = barterMessages[engine.GenreFantasy]
	}
	msg, ok := messages[key]
	if !ok {
		return key
	}
	return msg
}

// BarterVocab holds genre-specific bartering vocabulary.
type BarterVocab struct {
	BarterAction string
	OfferVerb    string
	RequestVerb  string
	AcceptVerb   string
	RejectVerb   string
	CounterVerb  string
	ValueLabel   string
}

// GetBarterVocab returns genre-specific bartering vocabulary.
func GetBarterVocab(genre engine.GenreID) *BarterVocab {
	vocab, ok := barterVocabs[genre]
	if !ok {
		return barterVocabs[engine.GenreFantasy]
	}
	return vocab
}

var barterVocabs = map[engine.GenreID]*BarterVocab{
	engine.GenreFantasy: {
		BarterAction: "Barter",
		OfferVerb:    "Offer",
		RequestVerb:  "Request",
		AcceptVerb:   "Accept",
		RejectVerb:   "Decline",
		CounterVerb:  "Counter-offer",
		ValueLabel:   "Worth",
	},
	engine.GenreScifi: {
		BarterAction: "Trade Protocol",
		OfferVerb:    "Transfer",
		RequestVerb:  "Requisition",
		AcceptVerb:   "Confirm",
		RejectVerb:   "Deny",
		CounterVerb:  "Alternative proposal",
		ValueLabel:   "Value",
	},
	engine.GenreHorror: {
		BarterAction: "Exchange",
		OfferVerb:    "Give",
		RequestVerb:  "Take",
		AcceptVerb:   "Deal",
		RejectVerb:   "No deal",
		CounterVerb:  "New terms",
		ValueLabel:   "Worth",
	},
	engine.GenreCyberpunk: {
		BarterAction: "Exchange Protocol",
		OfferVerb:    "Upload",
		RequestVerb:  "Download",
		AcceptVerb:   "Execute",
		RejectVerb:   "Abort",
		CounterVerb:  "Renegotiate",
		ValueLabel:   "Market value",
	},
	engine.GenrePostapoc: {
		BarterAction: "Swap",
		OfferVerb:    "Trade",
		RequestVerb:  "Get",
		AcceptVerb:   "Deal",
		RejectVerb:   "Walk away",
		CounterVerb:  "Counter",
		ValueLabel:   "Worth",
	},
}

var barterMessages = map[engine.GenreID]map[string]string{
	engine.GenreFantasy: {
		"no_inventory":  "You have nothing to trade",
		"empty_offer":   "You must offer something",
		"empty_request": "You must request something",
		"invalid_items": "Invalid items in the offer",
		"accepted":      "The merchant accepts your offer",
		"rejected":      "The merchant refuses your offer",
	},
	engine.GenreScifi: {
		"no_inventory":  "Inventory empty",
		"empty_offer":   "Offer manifest empty",
		"empty_request": "Request manifest empty",
		"invalid_items": "Invalid items detected",
		"accepted":      "Trade protocol executed",
		"rejected":      "Trade protocol rejected",
	},
	engine.GenreHorror: {
		"no_inventory":  "You have nothing",
		"empty_offer":   "Offer something",
		"empty_request": "What do you want?",
		"invalid_items": "Can't trade those",
		"accepted":      "They nod slowly",
		"rejected":      "They shake their head",
	},
	engine.GenreCyberpunk: {
		"no_inventory":  "Asset inventory: null",
		"empty_offer":   "Upload queue empty",
		"empty_request": "Download queue empty",
		"invalid_items": "Asset verification failed",
		"accepted":      "Transaction confirmed",
		"rejected":      "Transaction denied",
	},
	engine.GenrePostapoc: {
		"no_inventory":  "Got nothing to trade",
		"empty_offer":   "Gotta offer something",
		"empty_request": "What do you want?",
		"invalid_items": "Bad deal",
		"accepted":      "Fair enough",
		"rejected":      "No way",
	},
}

// CreateBarterOffer is a helper to create a barter offer.
func CreateBarterOffer(offered, requested map[string]int) *BarterOffer {
	offer := &BarterOffer{
		OfferedItems:   make([]*BarterItem, 0, len(offered)),
		RequestedItems: make([]*BarterItem, 0, len(requested)),
	}

	for name, qty := range offered {
		offer.OfferedItems = append(offer.OfferedItems, &BarterItem{
			ItemName: name,
			Quantity: qty,
		})
	}

	for name, qty := range requested {
		offer.RequestedItems = append(offer.RequestedItems, &BarterItem{
			ItemName: name,
			Quantity: qty,
		})
	}

	return offer
}
