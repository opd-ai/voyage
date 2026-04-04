package destination

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// Destination represents a journey's end point with procedurally generated content.
type Destination struct {
	Type        DestinationType
	Name        string
	Description string
	Distance    int // Turns remaining to reach
	Phase       DiscoveryPhase
	Discovered  bool // Has the party discovered this destination?
	Reached     bool // Has the party arrived?
	genre       engine.GenreID
	seedGen     *seed.Generator
}

// Generator creates procedural destinations.
type Generator struct {
	genre   engine.GenreID
	seedGen *seed.Generator
}

// NewGenerator creates a destination generator for the given genre.
func NewGenerator(genre engine.GenreID) *Generator {
	return &Generator{
		genre:   genre,
		seedGen: seed.NewGenerator(0, "destination"),
	}
}

// SetGenre changes the generator's genre.
func (g *Generator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// Generate creates a specific destination type.
func (g *Generator) Generate(destType DestinationType, distance int) *Destination {
	dest := &Destination{
		Type:       destType,
		Distance:   distance,
		Phase:      Distant,
		Discovered: false,
		Reached:    false,
		genre:      g.genre,
		seedGen:    seed.NewGenerator(0, "destination"),
	}
	dest.generateContent()
	return dest
}

// GenerateRandom creates a random destination type.
func (g *Generator) GenerateRandom(distance int) *Destination {
	types := AllDestinationTypes()
	destType := seed.Choice(g.seedGen, types)
	return g.Generate(destType, distance)
}

// GenerateSet creates a set of destinations for a journey.
func (g *Generator) GenerateSet(count, maxDistance int) []*Destination {
	destinations := make([]*Destination, 0, count)
	types := AllDestinationTypes()

	for i := 0; i < count; i++ {
		destType := types[i%len(types)]
		distance := ((i + 1) * maxDistance) / count
		destinations = append(destinations, g.Generate(destType, distance))
	}
	return destinations
}

// generateContent fills in the procedural name and description.
func (d *Destination) generateContent() {
	d.Name = d.generateName()
	d.Description = d.generateDescription()
}

// generateName creates a procedural name.
func (d *Destination) generateName() string {
	baseName := destinationNames[d.genre][d.Type]
	prefixes := destinationPrefixes[d.genre]
	prefix := seed.Choice(d.seedGen, prefixes)
	return prefix + " " + baseName
}

// generateDescription creates a description.
func (d *Destination) generateDescription() string {
	return destinationDescriptions[d.genre][d.Type]
}

// SetGenre changes the destination's genre and regenerates content.
func (d *Destination) SetGenre(genre engine.GenreID) {
	d.genre = genre
	d.generateContent()
}

// AdvanceTurn updates the destination state as the party travels.
func (d *Destination) AdvanceTurn() *DiscoveryEvent {
	if d.Reached {
		return nil
	}

	d.Distance--
	if d.Distance < 0 {
		d.Distance = 0
	}

	// Update phase based on distance
	oldPhase := d.Phase
	switch {
	case d.Distance == 0:
		d.Phase = Arrival
		d.Reached = true
	case d.Distance <= 2:
		d.Phase = Approach
	case d.Distance <= 5:
		d.Phase = Signs
	default:
		d.Phase = Distant
	}

	// Generate event if phase changed
	if d.Phase != oldPhase {
		return d.generateDiscoveryEvent()
	}
	return nil
}

// GetArrivalText returns the arrival ceremony narrative.
func (d *Destination) GetArrivalText() string {
	if !d.Reached {
		return ""
	}
	return d.generateArrivalCeremony()
}

// destinationPrefixes for procedural name generation.
var destinationPrefixes = map[engine.GenreID][]string{
	engine.GenreFantasy:   {"Ancient", "Lost", "Sacred", "Dark", "High", "Old"},
	engine.GenreScifi:     {"Alpha", "Beta", "Nova", "Prime", "Deep", "Far"},
	engine.GenreHorror:    {"Cursed", "Forgotten", "Haunted", "Silent", "Lost", "Dark"},
	engine.GenreCyberpunk: {"Neo", "Shadow", "Ghost", "Black", "Red", "Free"},
	engine.GenrePostapoc:  {"Last", "Broken", "Lost", "Dead", "New", "Safe"},
}
