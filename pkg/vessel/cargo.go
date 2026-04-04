package vessel

// Cargo represents an item in the vessel's cargo hold.
type Cargo struct {
	ID       int
	Name     string
	Weight   int
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
	items    []*Cargo
	capacity int
	used     int
	nextID   int
}

// NewCargoHold creates a cargo hold with the given capacity.
func NewCargoHold(capacity int) *CargoHold {
	return &CargoHold{
		items:    make([]*Cargo, 0),
		capacity: capacity,
		used:     0,
		nextID:   1,
	}
}

// Capacity returns the total cargo capacity.
func (h *CargoHold) Capacity() int {
	return h.capacity
}

// Used returns the used cargo space.
func (h *CargoHold) Used() int {
	return h.used
}

// Free returns the remaining cargo space.
func (h *CargoHold) Free() int {
	return h.capacity - h.used
}

// Items returns all cargo items.
func (h *CargoHold) Items() []*Cargo {
	return h.items
}

// Add adds cargo to the hold. Returns true if successful.
func (h *CargoHold) Add(name string, weight, quantity int, cat CargoCategory) bool {
	totalWeight := weight * quantity
	if h.used+totalWeight > h.capacity {
		return false
	}

	// Check if we can stack with existing item
	for _, item := range h.items {
		if item.Name == name && item.Category == cat {
			item.Quantity += quantity
			h.used += totalWeight
			return true
		}
	}

	// Add new item
	h.items = append(h.items, &Cargo{
		ID:       h.nextID,
		Name:     name,
		Weight:   weight,
		Quantity: quantity,
		Category: cat,
	})
	h.nextID++
	h.used += totalWeight
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
			h.used -= item.Weight * quantity
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
	h.used = 0
}

// TotalWeight returns the total weight of all cargo.
func (h *CargoHold) TotalWeight() int {
	total := 0
	for _, item := range h.items {
		total += item.Weight * item.Quantity
	}
	return total
}
