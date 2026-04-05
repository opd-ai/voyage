package lore

import (
	"github.com/opd-ai/voyage/pkg/engine"
)

// InscriptionType identifies the type of world-map inscription.
type InscriptionType int

const (
	// TypeRuin represents ancient ruins with carved inscriptions.
	TypeRuin InscriptionType = iota
	// TypeGrave represents grave markers or memorials.
	TypeGrave
	// TypeSign represents warning signs or directions.
	TypeSign
	// TypeMonument represents monuments or statues.
	TypeMonument
	// TypeGraffiti represents recent markings or graffiti.
	TypeGraffiti
)

// AllInscriptionTypes returns all inscription types.
func AllInscriptionTypes() []InscriptionType {
	return []InscriptionType{
		TypeRuin,
		TypeGrave,
		TypeSign,
		TypeMonument,
		TypeGraffiti,
	}
}

// InscriptionTypeName returns the genre-appropriate name for an inscription type.
func InscriptionTypeName(t InscriptionType, genre engine.GenreID) string {
	names := inscriptionTypeNames[genre]
	if names == nil {
		names = inscriptionTypeNames[engine.GenreFantasy]
	}
	return names[t]
}

var inscriptionTypeNames = map[engine.GenreID]map[InscriptionType]string{
	engine.GenreFantasy: {
		TypeRuin:     "Ancient Ruin",
		TypeGrave:    "Grave Marker",
		TypeSign:     "Weathered Sign",
		TypeMonument: "Stone Monument",
		TypeGraffiti: "Carved Message",
	},
	engine.GenreScifi: {
		TypeRuin:     "Derelict Structure",
		TypeGrave:    "Memorial Plaque",
		TypeSign:     "Warning Beacon",
		TypeMonument: "Monument",
		TypeGraffiti: "Etched Message",
	},
	engine.GenreHorror: {
		TypeRuin:     "Abandoned Building",
		TypeGrave:    "Shallow Grave",
		TypeSign:     "Spray-Painted Warning",
		TypeMonument: "Memorial",
		TypeGraffiti: "Scrawled Note",
	},
	engine.GenreCyberpunk: {
		TypeRuin:     "Collapsed Structure",
		TypeGrave:    "Memorial Wall",
		TypeSign:     "Holo-Sign Fragment",
		TypeMonument: "Corporate Monument",
		TypeGraffiti: "Tagged Message",
	},
	engine.GenrePostapoc: {
		TypeRuin:     "Pre-War Ruins",
		TypeGrave:    "Mass Grave",
		TypeSign:     "Faded Road Sign",
		TypeMonument: "Crumbling Statue",
		TypeGraffiti: "Scratched Warning",
	},
}

// Inscription represents a world-map lore inscription.
type Inscription struct {
	ID         int
	Type       InscriptionType
	X          int
	Y          int
	Title      string
	Text       string
	Genre      engine.GenreID
	Discovered bool
}

// NewInscription creates a new inscription.
func NewInscription(id int, iType InscriptionType, x, y int, title, text string, genre engine.GenreID) *Inscription {
	return &Inscription{
		ID:    id,
		Type:  iType,
		X:     x,
		Y:     y,
		Title: title,
		Text:  text,
		Genre: genre,
	}
}

// SetGenre updates the inscription's genre.
func (i *Inscription) SetGenre(genre engine.GenreID) {
	i.Genre = genre
}

// Discover marks the inscription as discovered.
func (i *Inscription) Discover() {
	i.Discovered = true
}

// TypeDisplayName returns the genre-appropriate type name.
func (i *Inscription) TypeDisplayName() string {
	return InscriptionTypeName(i.Type, i.Genre)
}

// DiscoveryType identifies the type of abandoned discovery.
type DiscoveryType int

const (
	// DiscoveryVessel represents an abandoned vehicle or vessel.
	DiscoveryVessel DiscoveryType = iota
	// DiscoveryCamp represents an abandoned campsite.
	DiscoveryCamp
	// DiscoveryCache represents a hidden supply cache.
	DiscoveryCache
	// DiscoveryBody represents remains with personal effects.
	DiscoveryBody
)

// AllDiscoveryTypes returns all discovery types.
func AllDiscoveryTypes() []DiscoveryType {
	return []DiscoveryType{
		DiscoveryVessel,
		DiscoveryCamp,
		DiscoveryCache,
		DiscoveryBody,
	}
}

// DiscoveryTypeName returns the genre-appropriate name for a discovery type.
func DiscoveryTypeName(t DiscoveryType, genre engine.GenreID) string {
	names := discoveryTypeNames[genre]
	if names == nil {
		names = discoveryTypeNames[engine.GenreFantasy]
	}
	return names[t]
}

var discoveryTypeNames = map[engine.GenreID]map[DiscoveryType]string{
	engine.GenreFantasy: {
		DiscoveryVessel: "Abandoned Wagon",
		DiscoveryCamp:   "Deserted Camp",
		DiscoveryCache:  "Hidden Cache",
		DiscoveryBody:   "Fallen Traveler",
	},
	engine.GenreScifi: {
		DiscoveryVessel: "Derelict Ship",
		DiscoveryCamp:   "Abandoned Station",
		DiscoveryCache:  "Supply Cache",
		DiscoveryBody:   "Remains",
	},
	engine.GenreHorror: {
		DiscoveryVessel: "Wrecked Vehicle",
		DiscoveryCamp:   "Abandoned Shelter",
		DiscoveryCache:  "Hidden Supplies",
		DiscoveryBody:   "Corpse",
	},
	engine.GenreCyberpunk: {
		DiscoveryVessel: "Crashed Runner",
		DiscoveryCamp:   "Abandoned Squat",
		DiscoveryCache:  "Stash",
		DiscoveryBody:   "Flatlined Body",
	},
	engine.GenrePostapoc: {
		DiscoveryVessel: "Rusted Wreck",
		DiscoveryCamp:   "Abandoned Camp",
		DiscoveryCache:  "Buried Cache",
		DiscoveryBody:   "Skeleton",
	},
}

// DiscoveryItem represents loot from a discovery.
type DiscoveryItem struct {
	Name     string
	Quantity int
}

// Discovery represents an abandoned vessel or camp.
type Discovery struct {
	ID           int
	Type         DiscoveryType
	X            int
	Y            int
	Title        string
	VignetteText string
	Items        []DiscoveryItem
	Genre        engine.GenreID
	Discovered   bool
	Looted       bool
}

// NewDiscovery creates a new discovery.
func NewDiscovery(id int, dType DiscoveryType, x, y int, title, vignette string, genre engine.GenreID) *Discovery {
	return &Discovery{
		ID:           id,
		Type:         dType,
		X:            x,
		Y:            y,
		Title:        title,
		VignetteText: vignette,
		Items:        make([]DiscoveryItem, 0),
		Genre:        genre,
	}
}

// SetGenre updates the discovery's genre.
func (d *Discovery) SetGenre(genre engine.GenreID) {
	d.Genre = genre
}

// Discover marks the discovery as discovered.
func (d *Discovery) Discover() {
	d.Discovered = true
}

// Loot marks the discovery as looted.
func (d *Discovery) Loot() {
	d.Looted = true
}

// AddItem adds an item to the discovery.
func (d *Discovery) AddItem(name string, quantity int) {
	d.Items = append(d.Items, DiscoveryItem{Name: name, Quantity: quantity})
}

// TypeDisplayName returns the genre-appropriate type name.
func (d *Discovery) TypeDisplayName() string {
	return DiscoveryTypeName(d.Type, d.Genre)
}

// CodexCategory identifies the type of codex entry.
type CodexCategory int

const (
	// CodexHistory represents world history entries.
	CodexHistory CodexCategory = iota
	// CodexFaction represents faction biography entries.
	CodexFaction
	// CodexRoute represents route and location legends.
	CodexRoute
	// CodexCreature represents creature or enemy lore.
	CodexCreature
	// CodexTechnology represents technology or artifact lore.
	CodexTechnology
)

// AllCodexCategories returns all codex categories.
func AllCodexCategories() []CodexCategory {
	return []CodexCategory{
		CodexHistory,
		CodexFaction,
		CodexRoute,
		CodexCreature,
		CodexTechnology,
	}
}

// CodexCategoryName returns the display name for a category.
func CodexCategoryName(c CodexCategory) string {
	names := map[CodexCategory]string{
		CodexHistory:    "History",
		CodexFaction:    "Factions",
		CodexRoute:      "Routes & Locations",
		CodexCreature:   "Bestiary",
		CodexTechnology: "Technology",
	}
	if name, ok := names[c]; ok {
		return name
	}
	return "Unknown"
}

// CodexEntry represents a lore codex entry.
type CodexEntry struct {
	ID           string
	Category     CodexCategory
	Title        string
	Text         string
	Genre        engine.GenreID
	Unlocked     bool
	UnlockSource string // How it was unlocked (exploration, event, NPC)
}

// NewCodexEntry creates a new codex entry.
func NewCodexEntry(id string, cat CodexCategory, title, text string, genre engine.GenreID) *CodexEntry {
	return &CodexEntry{
		ID:       id,
		Category: cat,
		Title:    title,
		Text:     text,
		Genre:    genre,
	}
}

// SetGenre updates the entry's genre.
func (e *CodexEntry) SetGenre(genre engine.GenreID) {
	e.Genre = genre
}

// Unlock marks the entry as unlocked.
func (e *CodexEntry) Unlock(source string) {
	e.Unlocked = true
	e.UnlockSource = source
}

// Codex manages all lore entries.
type Codex struct {
	Entries map[string]*CodexEntry
	genre   engine.GenreID
}

// NewCodex creates a new codex.
func NewCodex(genre engine.GenreID) *Codex {
	return &Codex{
		Entries: make(map[string]*CodexEntry),
		genre:   genre,
	}
}

// SetGenre updates the codex's genre and all entries.
func (c *Codex) SetGenre(genre engine.GenreID) {
	c.genre = genre
	for _, entry := range c.Entries {
		entry.SetGenre(genre)
	}
}

// AddEntry adds an entry to the codex.
func (c *Codex) AddEntry(entry *CodexEntry) {
	c.Entries[entry.ID] = entry
}

// GetEntry returns an entry by ID.
func (c *Codex) GetEntry(id string) *CodexEntry {
	return c.Entries[id]
}

// UnlockEntry unlocks an entry by ID.
func (c *Codex) UnlockEntry(id, source string) bool {
	if entry := c.Entries[id]; entry != nil && !entry.Unlocked {
		entry.Unlock(source)
		return true
	}
	return false
}

// GetUnlockedEntries returns all unlocked entries.
func (c *Codex) GetUnlockedEntries() []*CodexEntry {
	var result []*CodexEntry
	for _, entry := range c.Entries {
		if entry.Unlocked {
			result = append(result, entry)
		}
	}
	return result
}

// GetEntriesByCategory returns all entries in a category.
func (c *Codex) GetEntriesByCategory(cat CodexCategory) []*CodexEntry {
	var result []*CodexEntry
	for _, entry := range c.Entries {
		if entry.Category == cat {
			result = append(result, entry)
		}
	}
	return result
}

// GetUnlockedByCategory returns unlocked entries in a category.
func (c *Codex) GetUnlockedByCategory(cat CodexCategory) []*CodexEntry {
	var result []*CodexEntry
	for _, entry := range c.Entries {
		if entry.Category == cat && entry.Unlocked {
			result = append(result, entry)
		}
	}
	return result
}

// UnlockedCount returns the number of unlocked entries.
func (c *Codex) UnlockedCount() int {
	count := 0
	for _, entry := range c.Entries {
		if entry.Unlocked {
			count++
		}
	}
	return count
}

// TotalCount returns the total number of entries.
func (c *Codex) TotalCount() int {
	return len(c.Entries)
}

// EnvironmentalManager manages all environmental storytelling elements.
type EnvironmentalManager struct {
	Inscriptions map[int]*Inscription
	Discoveries  map[int]*Discovery
	Codex        *Codex
	genre        engine.GenreID
}

// NewEnvironmentalManager creates a new environmental manager.
func NewEnvironmentalManager(genre engine.GenreID) *EnvironmentalManager {
	return &EnvironmentalManager{
		Inscriptions: make(map[int]*Inscription),
		Discoveries:  make(map[int]*Discovery),
		Codex:        NewCodex(genre),
		genre:        genre,
	}
}

// SetGenre updates the manager's genre and all managed elements.
func (m *EnvironmentalManager) SetGenre(genre engine.GenreID) {
	m.genre = genre
	for _, i := range m.Inscriptions {
		i.SetGenre(genre)
	}
	for _, d := range m.Discoveries {
		d.SetGenre(genre)
	}
	m.Codex.SetGenre(genre)
}

// AddInscription adds an inscription.
func (m *EnvironmentalManager) AddInscription(i *Inscription) {
	m.Inscriptions[i.ID] = i
}

// AddDiscovery adds a discovery.
func (m *EnvironmentalManager) AddDiscovery(d *Discovery) {
	m.Discoveries[d.ID] = d
}

// GetInscriptionAt returns the inscription at a position.
func (m *EnvironmentalManager) GetInscriptionAt(x, y int) *Inscription {
	for _, i := range m.Inscriptions {
		if i.X == x && i.Y == y {
			return i
		}
	}
	return nil
}

// GetDiscoveryAt returns the discovery at a position.
func (m *EnvironmentalManager) GetDiscoveryAt(x, y int) *Discovery {
	for _, d := range m.Discoveries {
		if d.X == x && d.Y == y {
			return d
		}
	}
	return nil
}

// DiscoverAt discovers any elements at a position.
func (m *EnvironmentalManager) DiscoverAt(x, y int) (inscription *Inscription, discovery *Discovery) {
	inscription = m.GetInscriptionAt(x, y)
	if inscription != nil {
		inscription.Discover()
	}
	discovery = m.GetDiscoveryAt(x, y)
	if discovery != nil {
		discovery.Discover()
	}
	return inscription, discovery
}
