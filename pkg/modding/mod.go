package modding

import (
	"encoding/json"
	"time"
)

// Mod represents a loadable game modification.
type Mod struct {
	// ID is a unique identifier for this mod (e.g., "my-custom-mod").
	ID string `json:"id"`
	// Name is the human-readable name.
	Name string `json:"name"`
	// Version follows semantic versioning (e.g., "1.0.0").
	Version string `json:"version"`
	// Author is the mod creator's name.
	Author string `json:"author"`
	// Description explains what the mod adds.
	Description string `json:"description,omitempty"`
	// Events contains custom event definitions.
	Events []EventDef `json:"events,omitempty"`
	// Genres contains custom genre definitions.
	Genres []GenreDef `json:"genres,omitempty"`
	// Biomes contains custom biome names for existing genres.
	Biomes []BiomeDef `json:"biomes,omitempty"`
	// Resources contains custom resource definitions.
	Resources []ResourceDef `json:"resources,omitempty"`
	// Factions contains custom faction definitions.
	Factions []FactionDef `json:"factions,omitempty"`
	// LoadedAt is when the mod was loaded (set at runtime).
	LoadedAt time.Time `json:"-"`
	// Enabled indicates if the mod is currently active.
	Enabled bool `json:"-"`
}

// EventDef defines a custom event in JSON.
type EventDef struct {
	// Category is the event type: weather, encounter, discovery, etc.
	Category string `json:"category"`
	// Genre is which genre this event belongs to (fantasy, scifi, etc.).
	Genre string `json:"genre"`
	// Title is the event headline.
	Title string `json:"title"`
	// Description is the event narrative text.
	Description string `json:"description"`
	// Choices are the available player options.
	Choices []ChoiceDef `json:"choices"`
	// Weight affects how often this event is selected (default 1.0).
	Weight float64 `json:"weight,omitempty"`
	// MinTurn is the earliest turn this event can appear (default 0).
	MinTurn int `json:"min_turn,omitempty"`
	// MaxTurn is the latest turn this event can appear (0 = no limit).
	MaxTurn int `json:"max_turn,omitempty"`
	// RequiresBiome limits this event to specific biomes.
	RequiresBiome []string `json:"requires_biome,omitempty"`
}

// ChoiceDef defines a choice option in an event.
type ChoiceDef struct {
	// Text is what the player sees.
	Text string `json:"text"`
	// Outcome defines the effects of this choice.
	Outcome OutcomeDef `json:"outcome"`
	// RequireSkill makes this choice only available with a skill.
	RequireSkill string `json:"require_skill,omitempty"`
	// RequireMinResource only shows this choice with enough resources.
	RequireMinResource map[string]float64 `json:"require_min_resource,omitempty"`
}

// OutcomeDef defines the effects of a choice.
type OutcomeDef struct {
	// Description is flavor text shown after selection.
	Description string `json:"description,omitempty"`
	// Resource deltas (positive = gain, negative = loss).
	FoodDelta     float64 `json:"food_delta,omitempty"`
	WaterDelta    float64 `json:"water_delta,omitempty"`
	FuelDelta     float64 `json:"fuel_delta,omitempty"`
	MedicineDelta float64 `json:"medicine_delta,omitempty"`
	CurrencyDelta float64 `json:"currency_delta,omitempty"`
	// Morale change (-100 to 100).
	MoraleDelta float64 `json:"morale_delta,omitempty"`
	// Damage (0 to 100).
	CrewDamage   float64 `json:"crew_damage,omitempty"`
	VesselDamage float64 `json:"vessel_damage,omitempty"`
	// TimeAdvance is turns to skip (can simulate delays).
	TimeAdvance int `json:"time_advance,omitempty"`
}

// GenreDef defines a custom genre.
type GenreDef struct {
	// ID is the internal identifier (e.g., "steampunk").
	ID string `json:"id"`
	// Name is the display name (e.g., "Steampunk").
	Name string `json:"name"`
	// Description explains the genre theme.
	Description string `json:"description,omitempty"`
	// Biomes are location names for this genre.
	Biomes []string `json:"biomes,omitempty"`
	// Resources are resource names for this genre.
	Resources []string `json:"resources,omitempty"`
	// Factions are group names for this genre.
	Factions []string `json:"factions,omitempty"`
	// VesselTypes are vehicle names for this genre.
	VesselTypes []string `json:"vessel_types,omitempty"`
	// CrewRoles are job names for crew members.
	CrewRoles []string `json:"crew_roles,omitempty"`
	// CategoryNames maps event categories to genre-specific names.
	CategoryNames map[string]string `json:"category_names,omitempty"`
}

// BiomeDef adds biomes to an existing genre.
type BiomeDef struct {
	// Genre to add these biomes to.
	Genre string `json:"genre"`
	// Names are the biome names to add.
	Names []string `json:"names"`
}

// ResourceDef defines a custom resource type.
type ResourceDef struct {
	// ID is the internal identifier.
	ID string `json:"id"`
	// Name is the display name.
	Name string `json:"name"`
	// Genre limits this resource to a specific genre (empty = all).
	Genre string `json:"genre,omitempty"`
	// Description explains the resource.
	Description string `json:"description,omitempty"`
	// MaxStack is the maximum amount (0 = unlimited).
	MaxStack float64 `json:"max_stack,omitempty"`
	// Tradeable indicates if this can be bought/sold.
	Tradeable bool `json:"tradeable,omitempty"`
}

// FactionDef defines a custom faction.
type FactionDef struct {
	// ID is the internal identifier.
	ID string `json:"id"`
	// Name is the display name.
	Name string `json:"name"`
	// Genre limits this faction to a specific genre (empty = all).
	Genre string `json:"genre,omitempty"`
	// Description explains the faction.
	Description string `json:"description,omitempty"`
	// Hostile indicates default stance toward player.
	Hostile bool `json:"hostile,omitempty"`
}

// Validate checks if the mod definition is valid.
func (m *Mod) Validate() error {
	if err := m.validateRequiredFields(); err != nil {
		return err
	}
	if err := m.validateEvents(); err != nil {
		return err
	}
	return m.validateGenres()
}

// validateRequiredFields checks that mod has required ID, name, and version.
func (m *Mod) validateRequiredFields() error {
	if m.ID == "" {
		return ErrMissingID
	}
	if m.Name == "" {
		return ErrMissingName
	}
	if m.Version == "" {
		return ErrMissingVersion
	}
	return nil
}

// validateEvents checks all events in the mod.
func (m *Mod) validateEvents() error {
	for i, e := range m.Events {
		if err := e.Validate(); err != nil {
			return &ValidationError{Field: "events", Index: i, Message: err.Error()}
		}
	}
	return nil
}

// validateGenres checks all genres in the mod.
func (m *Mod) validateGenres() error {
	for i, g := range m.Genres {
		if err := g.Validate(); err != nil {
			return &ValidationError{Field: "genres", Index: i, Message: err.Error()}
		}
	}
	return nil
}

// Validate checks if the event definition is valid.
func (e *EventDef) Validate() error {
	if err := e.validateRequiredFields(); err != nil {
		return err
	}
	return e.validateChoices()
}

// validateRequiredFields checks event's required text fields and category.
func (e *EventDef) validateRequiredFields() error {
	if e.Title == "" {
		return ErrMissingTitle
	}
	if e.Description == "" {
		return ErrMissingDescription
	}
	if len(e.Choices) == 0 {
		return ErrNoChoices
	}
	if e.Category == "" {
		return ErrMissingCategory
	}
	if !isValidCategory(e.Category) {
		return ErrInvalidCategory
	}
	return nil
}

// validateChoices checks that all choices have text.
func (e *EventDef) validateChoices() error {
	for i, c := range e.Choices {
		if c.Text == "" {
			return &ValidationError{Field: "choices", Index: i, Message: "choice text is required"}
		}
	}
	return nil
}

// Validate checks if the genre definition is valid.
func (g *GenreDef) Validate() error {
	if g.ID == "" {
		return ErrMissingID
	}
	if g.Name == "" {
		return ErrMissingName
	}
	return nil
}

// MarshalJSON implements custom JSON marshaling.
func (m *Mod) MarshalJSON() ([]byte, error) {
	type Alias Mod
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// isValidCategory checks if a category string is valid.
func isValidCategory(cat string) bool {
	validCategories := []string{
		"weather", "encounter", "discovery",
		"hardship", "windfall", "hazard", "crew",
	}
	for _, v := range validCategories {
		if cat == v {
			return true
		}
	}
	return false
}

// ValidCategories returns all valid event category names.
func ValidCategories() []string {
	return []string{
		"weather", "encounter", "discovery",
		"hardship", "windfall", "hazard", "crew",
	}
}

// ValidGenres returns the built-in genre IDs.
func ValidGenres() []string {
	return []string{
		"fantasy", "scifi", "horror", "cyberpunk", "postapoc",
	}
}
