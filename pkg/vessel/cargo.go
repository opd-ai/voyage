package vessel

// Cargo represents an item in the vessel's cargo hold.
type Cargo struct {
	ID       int
	Name     string
	Weight   int
	Volume   int // Volume per unit
	Quantity int
	Category CargoCategory
}

// CargoCategory identifies the type of cargo.
type CargoCategory int

const (
	// CargoSupplies includes food, water, fuel.
	CargoSupplies CargoCategory = iota
	// CargoMedical includes medicine and healing items.
	CargoMedical
	// CargoRepair includes tools and repair materials.
	CargoRepair
	// CargoTrade includes valuable goods for trading.
	CargoTrade
	// CargoSpecial includes quest items and special cargo.
	CargoSpecial
)

// AllCargoCategories returns all cargo categories.
func AllCargoCategories() []CargoCategory {
	return []CargoCategory{
		CargoSupplies,
		CargoMedical,
		CargoRepair,
		CargoTrade,
		CargoSpecial,
	}
}

// CargoHold manages the vessel's cargo inventory.
type CargoHold struct {
	items       []*Cargo
	weightLimit int
	volumeLimit int
	usedWeight  int
	usedVolume  int
	nextID      int
	tier        int
}

// CargoHoldTierStats defines capacity per tier.
var CargoHoldTierStats = map[int]struct {
	WeightLimit int
	VolumeLimit int
}{
	1: {WeightLimit: 50, VolumeLimit: 40},
	2: {WeightLimit: 75, VolumeLimit: 60},
	3: {WeightLimit: 100, VolumeLimit: 80},
	4: {WeightLimit: 130, VolumeLimit: 105},
	5: {WeightLimit: 175, VolumeLimit: 140},
}

// NewCargoHold creates a cargo hold with the given capacity (tier 1).
func NewCargoHold(capacity int) *CargoHold {
	return &CargoHold{
		items:       make([]*Cargo, 0),
		weightLimit: capacity,
		volumeLimit: int(float64(capacity) * 0.8), // Default volume ~80% of weight
		usedWeight:  0,
		usedVolume:  0,
		nextID:      1,
		tier:        1,
	}
}

// NewCargoHoldWithTier creates a cargo hold at the specified tier.
func NewCargoHoldWithTier(tier int) *CargoHold {
	if tier < 1 {
		tier = 1
	}
	if tier > 5 {
		tier = 5
	}
	stats := CargoHoldTierStats[tier]
	return &CargoHold{
		items:       make([]*Cargo, 0),
		weightLimit: stats.WeightLimit,
		volumeLimit: stats.VolumeLimit,
		usedWeight:  0,
		usedVolume:  0,
		nextID:      1,
		tier:        tier,
	}
}

// Tier returns the cargo hold tier.
func (h *CargoHold) Tier() int {
	return h.tier
}

// SetTier changes the cargo hold tier, adjusting limits.
func (h *CargoHold) SetTier(tier int) {
	if tier < 1 {
		tier = 1
	}
	if tier > 5 {
		tier = 5
	}
	stats := CargoHoldTierStats[tier]
	h.tier = tier
	h.weightLimit = stats.WeightLimit
	h.volumeLimit = stats.VolumeLimit
}

// Capacity returns the total cargo weight capacity.
func (h *CargoHold) Capacity() int {
	return h.weightLimit
}

// WeightLimit returns the weight limit.
func (h *CargoHold) WeightLimit() int {
	return h.weightLimit
}

// VolumeLimit returns the volume limit.
func (h *CargoHold) VolumeLimit() int {
	return h.volumeLimit
}

// Used returns the used cargo weight (for backward compatibility).
func (h *CargoHold) Used() int {
	return h.usedWeight
}

// UsedWeight returns the used cargo weight.
func (h *CargoHold) UsedWeight() int {
	return h.usedWeight
}

// UsedVolume returns the used cargo volume.
func (h *CargoHold) UsedVolume() int {
	return h.usedVolume
}

// Free returns the remaining cargo weight space.
func (h *CargoHold) Free() int {
	return h.weightLimit - h.usedWeight
}

// FreeWeight returns the remaining weight capacity.
func (h *CargoHold) FreeWeight() int {
	return h.weightLimit - h.usedWeight
}

// FreeVolume returns the remaining volume capacity.
func (h *CargoHold) FreeVolume() int {
	return h.volumeLimit - h.usedVolume
}

// WeightRatio returns the used weight as a ratio [0, 1].
func (h *CargoHold) WeightRatio() float64 {
	if h.weightLimit <= 0 {
		return 0
	}
	return float64(h.usedWeight) / float64(h.weightLimit)
}

// VolumeRatio returns the used volume as a ratio [0, 1].
func (h *CargoHold) VolumeRatio() float64 {
	if h.volumeLimit <= 0 {
		return 0
	}
	return float64(h.usedVolume) / float64(h.volumeLimit)
}

// Items returns all cargo items.
func (h *CargoHold) Items() []*Cargo {
	return h.items
}

// CanAdd checks if cargo can be added without exceeding limits.
func (h *CargoHold) CanAdd(weight, volume, quantity int) bool {
	totalWeight := weight * quantity
	totalVolume := volume * quantity
	return h.usedWeight+totalWeight <= h.weightLimit &&
		h.usedVolume+totalVolume <= h.volumeLimit
}

// Add adds cargo to the hold. Returns true if successful.
func (h *CargoHold) Add(name string, weight, quantity int, cat CargoCategory) bool {
	return h.AddWithVolume(name, weight, weight, quantity, cat) // Default volume = weight
}

// AddWithVolume adds cargo with explicit volume to the hold.
func (h *CargoHold) AddWithVolume(name string, weight, volume, quantity int, cat CargoCategory) bool {
	totalWeight := weight * quantity
	totalVolume := volume * quantity
	if h.usedWeight+totalWeight > h.weightLimit {
		return false
	}
	if h.usedVolume+totalVolume > h.volumeLimit {
		return false
	}

	// Check if we can stack with existing item
	for _, item := range h.items {
		if item.Name == name && item.Category == cat {
			item.Quantity += quantity
			h.usedWeight += totalWeight
			h.usedVolume += totalVolume
			return true
		}
	}

	// Add new item
	h.items = append(h.items, &Cargo{
		ID:       h.nextID,
		Name:     name,
		Weight:   weight,
		Volume:   volume,
		Quantity: quantity,
		Category: cat,
	})
	h.nextID++
	h.usedWeight += totalWeight
	h.usedVolume += totalVolume
	return true
}

// Remove removes cargo from the hold. Returns true if successful.
func (h *CargoHold) Remove(name string, quantity int) bool {
	for i, item := range h.items {
		if item.Name == name {
			if item.Quantity < quantity {
				return false
			}
			item.Quantity -= quantity
			h.usedWeight -= item.Weight * quantity
			h.usedVolume -= item.Volume * quantity
			if item.Quantity <= 0 {
				h.items = append(h.items[:i], h.items[i+1:]...)
			}
			return true
		}
	}
	return false
}

// Has checks if the hold contains at least the given quantity.
func (h *CargoHold) Has(name string, quantity int) bool {
	for _, item := range h.items {
		if item.Name == name && item.Quantity >= quantity {
			return true
		}
	}
	return false
}

// GetQuantity returns the quantity of a specific item.
func (h *CargoHold) GetQuantity(name string) int {
	for _, item := range h.items {
		if item.Name == name {
			return item.Quantity
		}
	}
	return 0
}

// GetByCategory returns all cargo of a specific category.
func (h *CargoHold) GetByCategory(cat CargoCategory) []*Cargo {
	var result []*Cargo
	for _, item := range h.items {
		if item.Category == cat {
			result = append(result, item)
		}
	}
	return result
}

// Clear empties the cargo hold.
func (h *CargoHold) Clear() {
	h.items = make([]*Cargo, 0)
	h.usedWeight = 0
	h.usedVolume = 0
}

// TotalWeight returns the total weight of all cargo.
func (h *CargoHold) TotalWeight() int {
	total := 0
	for _, item := range h.items {
		total += item.Weight * item.Quantity
	}
	return total
}

// TotalVolume returns the total volume of all cargo.
func (h *CargoHold) TotalVolume() int {
	total := 0
	for _, item := range h.items {
		total += item.Volume * item.Quantity
	}
	return total
}
