package world

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
)

// Tile represents a single map tile.
type Tile struct {
	X, Y        int
	Terrain     TerrainType
	Biome       BiomeType
	Explored    bool
	Landmark    *Landmark
	Connections []Point
}

// Point represents a 2D coordinate.
type Point struct {
	X, Y int
}

// Landmark represents a notable location on the map.
type Landmark struct {
	Type        LandmarkType
	Name        string
	Description string
}

// LandmarkType identifies the kind of landmark.
type LandmarkType int

const (
	// LandmarkTown is a settlement for trading.
	LandmarkTown LandmarkType = iota
	// LandmarkOutpost is a small rest stop.
	LandmarkOutpost
	// LandmarkRuins is an explorable location.
	LandmarkRuins
	// LandmarkShrine is a healing/rest location.
	LandmarkShrine
	// LandmarkOrigin is the starting point.
	LandmarkOrigin
	// LandmarkDestination is the goal.
	LandmarkDestination
)

// WorldMap holds the generated game world.
type WorldMap struct {
	Width       int
	Height      int
	Tiles       [][]*Tile
	Origin      Point
	Destination Point
	Genre       engine.GenreID
}

// Generator creates procedural world maps.
type Generator struct {
	gen   *seed.Generator
	genre engine.GenreID
}

// NewGenerator creates a new world generator.
func NewGenerator(masterSeed int64, genre engine.GenreID) *Generator {
	return &Generator{
		gen:   seed.NewGenerator(masterSeed, "world"),
		genre: genre,
	}
}

// SetGenre changes the generator's genre.
func (g *Generator) SetGenre(genre engine.GenreID) {
	g.genre = genre
}

// Generate creates a new world map.
func (g *Generator) Generate(width, height int) *WorldMap {
	w := &WorldMap{
		Width:  width,
		Height: height,
		Tiles:  make([][]*Tile, height),
		Genre:  g.genre,
	}

	// Initialize tiles
	for y := 0; y < height; y++ {
		w.Tiles[y] = make([]*Tile, width)
		for x := 0; x < width; x++ {
			w.Tiles[y][x] = &Tile{
				X:           x,
				Y:           y,
				Connections: make([]Point, 0),
			}
		}
	}

	// Assign biomes using Voronoi-like regions
	g.assignBiomes(w)

	// Assign terrain based on biomes
	g.assignTerrain(w)

	// Place origin and destination
	g.placeOriginDestination(w)

	// Generate path network
	g.generatePaths(w)

	// Place landmarks
	g.placeLandmarks(w)

	return w
}

// assignBiomes divides the map into biome regions.
func (g *Generator) assignBiomes(w *WorldMap) {
	centers := g.generateBiomeCenters(w)
	g.assignTilesToNearestCenter(w, centers)
}

// biomeCenter holds the position and biome type for a region.
type biomeCenter struct {
	x, y  int
	biome BiomeType
}

// generateBiomeCenters creates random biome region centers.
func (g *Generator) generateBiomeCenters(w *WorldMap) []biomeCenter {
	numRegions := 6 + g.gen.Intn(4)
	centers := make([]biomeCenter, numRegions)
	biomes := AllBiomeTypes()

	for i := 0; i < numRegions; i++ {
		centers[i] = biomeCenter{
			x:     g.gen.Intn(w.Width),
			y:     g.gen.Intn(w.Height),
			biome: seed.Choice(g.gen, biomes),
		}
	}
	return centers
}

// assignTilesToNearestCenter assigns each tile to the nearest biome center.
func (g *Generator) assignTilesToNearestCenter(w *WorldMap, centers []biomeCenter) {
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			w.Tiles[y][x].Biome = g.findNearestBiome(x, y, centers, w)
		}
	}
}

// findNearestBiome returns the biome type of the nearest center.
func (g *Generator) findNearestBiome(x, y int, centers []biomeCenter, w *WorldMap) BiomeType {
	minDist := w.Width * w.Height
	var nearestBiome BiomeType

	for _, c := range centers {
		dist := abs(x-c.x) + abs(y-c.y)
		if dist < minDist {
			minDist = dist
			nearestBiome = c.biome
		}
	}
	return nearestBiome
}

// assignTerrain fills in terrain based on biomes.
func (g *Generator) assignTerrain(w *WorldMap) {
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			tile := w.Tiles[y][x]
			tile.Terrain = SelectTerrain(g.gen, tile.Biome)
		}
	}
}

// placeOriginDestination sets start and end points.
func (g *Generator) placeOriginDestination(w *WorldMap) {
	// Origin on left side
	w.Origin = Point{
		X: 1 + g.gen.Intn(w.Width/4),
		Y: w.Height/3 + g.gen.Intn(w.Height/3),
	}

	// Destination on right side
	w.Destination = Point{
		X: w.Width - 2 - g.gen.Intn(w.Width/4),
		Y: w.Height/3 + g.gen.Intn(w.Height/3),
	}

	// Mark tiles
	w.Tiles[w.Origin.Y][w.Origin.X].Landmark = &Landmark{
		Type: LandmarkOrigin,
		Name: g.landmarkName(LandmarkOrigin),
	}
	w.Tiles[w.Origin.Y][w.Origin.X].Terrain = TerrainPlains

	w.Tiles[w.Destination.Y][w.Destination.X].Landmark = &Landmark{
		Type: LandmarkDestination,
		Name: g.landmarkName(LandmarkDestination),
	}
	w.Tiles[w.Destination.Y][w.Destination.X].Terrain = TerrainPlains
}

// generatePaths creates a connected path network.
func (g *Generator) generatePaths(w *WorldMap) {
	// Generate main path from origin to destination
	g.connectPoints(w, w.Origin, w.Destination)

	// Add some alternative paths
	numAltPaths := 2 + g.gen.Intn(3)
	for i := 0; i < numAltPaths; i++ {
		// Create waypoint
		midX := (w.Origin.X + w.Destination.X) / 2
		waypoint := Point{
			X: midX + g.gen.Range(-w.Width/4, w.Width/4),
			Y: g.gen.Intn(w.Height),
		}
		waypoint.X = clamp(waypoint.X, 1, w.Width-2)
		waypoint.Y = clamp(waypoint.Y, 1, w.Height-2)

		g.connectPoints(w, w.Origin, waypoint)
		g.connectPoints(w, waypoint, w.Destination)
	}
}

// connectPoints creates connections between two points.
func (g *Generator) connectPoints(w *WorldMap, from, to Point) {
	current := from
	for current.X != to.X || current.Y != to.Y {
		next := g.calculateNextStep(current, to, w)
		g.addBidirectionalConnection(w, current, next)
		current = next
	}
}

// calculateNextStep determines the next point in the path with optional deviation.
func (g *Generator) calculateNextStep(current, to Point, w *WorldMap) Point {
	dx := sign(to.X - current.X)
	dy := sign(to.Y - current.Y)

	if g.gen.Chance(0.3) {
		dx, dy = g.applyDeviation(dx, dy)
	}

	return Point{
		X: clamp(current.X+dx, 0, w.Width-1),
		Y: clamp(current.Y+dy, 0, w.Height-1),
	}
}

// applyDeviation randomly adjusts movement direction.
func (g *Generator) applyDeviation(dx, dy int) (int, int) {
	if g.gen.Chance(0.5) && dx != 0 {
		return dx, g.gen.Range(-1, 1)
	}
	if dy != 0 {
		return g.gen.Range(-1, 1), dy
	}
	return dx, dy
}

// addBidirectionalConnection adds a two-way connection between tiles.
func (g *Generator) addBidirectionalConnection(w *WorldMap, a, b Point) {
	tileA := w.Tiles[a.Y][a.X]
	tileB := w.Tiles[b.Y][b.X]

	if !containsPoint(tileA.Connections, b) {
		tileA.Connections = append(tileA.Connections, b)
	}
	if !containsPoint(tileB.Connections, a) {
		tileB.Connections = append(tileB.Connections, a)
	}
}

// placeLandmarks adds towns, outposts, and ruins.
func (g *Generator) placeLandmarks(w *WorldMap) {
	numTowns := 3 + g.gen.Intn(3)
	numOutposts := 4 + g.gen.Intn(4)
	numRuins := 2 + g.gen.Intn(3)
	numShrines := 1 + g.gen.Intn(2)

	g.placeLandmarkType(w, LandmarkTown, numTowns)
	g.placeLandmarkType(w, LandmarkOutpost, numOutposts)
	g.placeLandmarkType(w, LandmarkRuins, numRuins)
	g.placeLandmarkType(w, LandmarkShrine, numShrines)
}

// placeLandmarkType places landmarks of a specific type.
func (g *Generator) placeLandmarkType(w *WorldMap, lt LandmarkType, count int) {
	placed := 0
	attempts := 0
	maxAttempts := count * 20

	for placed < count && attempts < maxAttempts {
		attempts++
		x := g.gen.Intn(w.Width)
		y := g.gen.Intn(w.Height)

		tile := w.Tiles[y][x]
		if tile.Landmark != nil {
			continue
		}
		if len(tile.Connections) == 0 && lt != LandmarkRuins {
			continue
		}

		tile.Landmark = &Landmark{
			Type: lt,
			Name: g.landmarkName(lt),
		}
		placed++
	}
}

// landmarkName generates a name for a landmark.
func (g *Generator) landmarkName(lt LandmarkType) string {
	names := landmarkNames[g.genre][lt]
	if len(names) == 0 {
		return "Unknown"
	}
	prefix := seed.Choice(g.gen, names)
	suffix := seed.Choice(g.gen, landmarkSuffixes[g.genre])
	return prefix + " " + suffix
}

var landmarkNames = map[engine.GenreID]map[LandmarkType][]string{
	engine.GenreFantasy: {
		LandmarkTown:        {"Green", "Silver", "Golden", "Iron", "Stone"},
		LandmarkOutpost:     {"Lonely", "Far", "Last", "Northern", "Eastern"},
		LandmarkRuins:       {"Ancient", "Forgotten", "Cursed", "Haunted", "Lost"},
		LandmarkShrine:      {"Sacred", "Holy", "Blessed", "Divine", "Mystic"},
		LandmarkOrigin:      {"Home"},
		LandmarkDestination: {"Promised"},
	},
	engine.GenreScifi: {
		LandmarkTown:        {"Alpha", "Beta", "Gamma", "Delta", "Omega"},
		LandmarkOutpost:     {"Relay", "Beacon", "Signal", "Remote", "Frontier"},
		LandmarkRuins:       {"Derelict", "Abandoned", "Ancient", "Lost", "Dead"},
		LandmarkShrine:      {"Medical", "Repair", "Rest", "Recovery", "Safe"},
		LandmarkOrigin:      {"Launch"},
		LandmarkDestination: {"Target"},
	},
	engine.GenreHorror: {
		LandmarkTown:        {"Last", "Final", "Dead", "Silent", "Dark"},
		LandmarkOutpost:     {"Broken", "Crumbling", "Barricaded", "Hidden", "Safe"},
		LandmarkRuins:       {"Infested", "Overrun", "Fallen", "Doomed", "Cursed"},
		LandmarkShrine:      {"Survivor", "Hope", "Sanctuary", "Refuge", "Haven"},
		LandmarkOrigin:      {"Starting"},
		LandmarkDestination: {"Safe"},
	},
	engine.GenreCyberpunk: {
		LandmarkTown:        {"Neon", "Chrome", "Digital", "Cyber", "Tech"},
		LandmarkOutpost:     {"Node", "Hub", "Link", "Junction", "Gateway"},
		LandmarkRuins:       {"Crashed", "Corrupted", "Offline", "Dark", "Ghost"},
		LandmarkShrine:      {"Med", "Repair", "Upgrade", "Rest", "Clinic"},
		LandmarkOrigin:      {"Upload"},
		LandmarkDestination: {"Download"},
	},
	engine.GenrePostapoc: {
		LandmarkTown:        {"New", "Free", "Last", "Reclaimed", "Survivor"},
		LandmarkOutpost:     {"Rusty", "Patched", "Makeshift", "Hidden", "Watch"},
		LandmarkRuins:       {"Blasted", "Irradiated", "Collapsed", "Scavenged", "Toxic"},
		LandmarkShrine:      {"Clean", "Pure", "Med", "Rest", "Safe"},
		LandmarkOrigin:      {"Home"},
		LandmarkDestination: {"Promised"},
	},
}

var landmarkSuffixes = map[engine.GenreID][]string{
	engine.GenreFantasy:   {"Haven", "Hold", "Keep", "Valley", "Crossing", "Ford", "Gate"},
	engine.GenreScifi:     {"Station", "Base", "Outpost", "Array", "Platform", "Dock"},
	engine.GenreHorror:    {"Camp", "Base", "Zone", "Point", "Position", "Site"},
	engine.GenreCyberpunk: {"District", "Sector", "Block", "Zone", "Level", "Grid"},
	engine.GenrePostapoc:  {"Camp", "Settlement", "Bunker", "Post", "Site", "Zone"},
}

// Helper functions
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sign(x int) int {
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

func clamp(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func containsPoint(points []Point, p Point) bool {
	for _, pt := range points {
		if pt.X == p.X && pt.Y == p.Y {
			return true
		}
	}
	return false
}

// GetTile returns the tile at the given coordinates.
func (w *WorldMap) GetTile(x, y int) *Tile {
	if x < 0 || x >= w.Width || y < 0 || y >= w.Height {
		return nil
	}
	return w.Tiles[y][x]
}

// IsValidMove checks if movement to a tile is valid.
func (w *WorldMap) IsValidMove(from, to Point) bool {
	fromTile := w.GetTile(from.X, from.Y)
	if fromTile == nil {
		return false
	}
	return containsPoint(fromTile.Connections, to)
}
