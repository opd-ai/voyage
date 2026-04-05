package procgen

import (
	"crypto/sha256"
	"encoding/json"
	"testing"

	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/events"
	"github.com/opd-ai/voyage/pkg/procgen/world"
)

// TestFullRunDeterminism verifies that the same seed produces byte-identical
// output across world maps, crew, and events. This validates the core promise
// of the project: deterministic procedural generation from a single seed.
func TestFullRunDeterminism(t *testing.T) {
	const masterSeed int64 = 12345
	const mapWidth = 30
	const mapHeight = 20
	const crewCount = 6
	const eventCount = 100

	// Generate first run
	state1 := generateFullState(masterSeed, mapWidth, mapHeight, crewCount, eventCount)
	hash1 := hashState(t, state1)

	// Generate second run with identical seed
	state2 := generateFullState(masterSeed, mapWidth, mapHeight, crewCount, eventCount)
	hash2 := hashState(t, state2)

	// Assert byte-equality
	if hash1 != hash2 {
		t.Errorf("Determinism violation: same seed produced different state hashes\n"+
			"Run 1: %x\nRun 2: %x", hash1, hash2)
	}
}

// TestDeterminismAcrossGenres verifies determinism holds for all five genres.
func TestDeterminismAcrossGenres(t *testing.T) {
	const masterSeed int64 = 54321

	for _, genre := range engine.AllGenres() {
		t.Run(string(genre), func(t *testing.T) {
			state1 := generateStateForGenre(masterSeed, genre)
			hash1 := hashState(t, state1)

			state2 := generateStateForGenre(masterSeed, genre)
			hash2 := hashState(t, state2)

			if hash1 != hash2 {
				t.Errorf("Determinism violation for genre %s: hashes differ", genre)
			}
		})
	}
}

// TestDifferentSeedsProduceDifferentOutput verifies that different seeds
// produce different output (non-trivial generation).
func TestDifferentSeedsProduceDifferentOutput(t *testing.T) {
	state1 := generateFullState(11111, 20, 15, 4, 50)
	hash1 := hashState(t, state1)

	state2 := generateFullState(22222, 20, 15, 4, 50)
	hash2 := hashState(t, state2)

	if hash1 == hash2 {
		t.Error("Different seeds should produce different output")
	}
}

// TestWorldMapDeterminism specifically tests world map generation determinism.
func TestWorldMapDeterminism(t *testing.T) {
	const masterSeed int64 = 99999
	genre := engine.GenreFantasy

	gen1 := world.NewGenerator(masterSeed, genre)
	worldMap1 := gen1.Generate(40, 30)

	gen2 := world.NewGenerator(masterSeed, genre)
	worldMap2 := gen2.Generate(40, 30)

	// Compare key world properties
	if worldMap1.Origin.X != worldMap2.Origin.X || worldMap1.Origin.Y != worldMap2.Origin.Y {
		t.Errorf("Origin mismatch: (%d,%d) vs (%d,%d)",
			worldMap1.Origin.X, worldMap1.Origin.Y,
			worldMap2.Origin.X, worldMap2.Origin.Y)
	}

	if worldMap1.Destination.X != worldMap2.Destination.X ||
		worldMap1.Destination.Y != worldMap2.Destination.Y {
		t.Errorf("Destination mismatch: (%d,%d) vs (%d,%d)",
			worldMap1.Destination.X, worldMap1.Destination.Y,
			worldMap2.Destination.X, worldMap2.Destination.Y)
	}

	// Compare all tiles
	for y := 0; y < worldMap1.Height; y++ {
		for x := 0; x < worldMap1.Width; x++ {
			tile1 := worldMap1.Tiles[y][x]
			tile2 := worldMap2.Tiles[y][x]

			if tile1.Terrain != tile2.Terrain {
				t.Errorf("Terrain mismatch at (%d,%d): %v vs %v",
					x, y, tile1.Terrain, tile2.Terrain)
			}
			if tile1.Biome != tile2.Biome {
				t.Errorf("Biome mismatch at (%d,%d): %v vs %v",
					x, y, tile1.Biome, tile2.Biome)
			}
		}
	}
}

// TestCrewDeterminism specifically tests crew generation determinism.
func TestCrewDeterminism(t *testing.T) {
	const masterSeed int64 = 77777
	genre := engine.GenreScifi

	gen1 := crew.NewGenerator(masterSeed, genre)
	gen2 := crew.NewGenerator(masterSeed, genre)

	for i := 0; i < 10; i++ {
		member1 := gen1.Generate()
		member2 := gen2.Generate()

		if member1.Name != member2.Name {
			t.Errorf("Crew %d name mismatch: %s vs %s", i, member1.Name, member2.Name)
		}
		if member1.Trait != member2.Trait {
			t.Errorf("Crew %d trait mismatch: %v vs %v", i, member1.Trait, member2.Trait)
		}
		if member1.Skill != member2.Skill {
			t.Errorf("Crew %d skill mismatch: %v vs %v", i, member1.Skill, member2.Skill)
		}
		if member1.Backstory.Origin != member2.Backstory.Origin {
			t.Errorf("Crew %d backstory origin mismatch", i)
		}
	}
}

// TestEventsDeterminism specifically tests event generation determinism.
func TestEventsDeterminism(t *testing.T) {
	const masterSeed int64 = 88888
	genre := engine.GenreHorror

	q1 := events.NewQueue(masterSeed, genre)
	q2 := events.NewQueue(masterSeed, genre)

	for turn := 0; turn < 50; turn++ {
		x, y := turn%10, turn/10

		event1 := q1.Generate(x, y, turn)
		event2 := q2.Generate(x, y, turn)

		if event1.Title != event2.Title {
			t.Errorf("Turn %d: event title mismatch: %s vs %s",
				turn, event1.Title, event2.Title)
		}
		if event1.Category != event2.Category {
			t.Errorf("Turn %d: event category mismatch: %v vs %v",
				turn, event1.Category, event2.Category)
		}
	}
}

// fullState captures the complete generated state for hashing.
type fullState struct {
	World  worldState   `json:"world"`
	Crew   []crewState  `json:"crew"`
	Events []eventState `json:"events"`
}

type worldState struct {
	OriginX     int    `json:"origin_x"`
	OriginY     int    `json:"origin_y"`
	DestX       int    `json:"dest_x"`
	DestY       int    `json:"dest_y"`
	TerrainHash string `json:"terrain_hash"`
	BiomeHash   string `json:"biome_hash"`
}

type crewState struct {
	Name   string `json:"name"`
	Trait  int    `json:"trait"`
	Skill  int    `json:"skill"`
	Health int    `json:"health"`
	Origin string `json:"origin"`
}

type eventState struct {
	Title    string `json:"title"`
	Category int    `json:"category"`
}

func generateFullState(seed int64, mapW, mapH, crewN, eventN int) *fullState {
	genre := engine.GenreFantasy

	// Generate world
	worldGen := world.NewGenerator(seed, genre)
	worldMap := worldGen.Generate(mapW, mapH)

	// Generate crew
	crewGen := crew.NewGenerator(seed, genre)
	crewMembers := make([]crewState, crewN)
	for i := 0; i < crewN; i++ {
		member := crewGen.Generate()
		crewMembers[i] = crewState{
			Name:   member.Name,
			Trait:  int(member.Trait),
			Skill:  int(member.Skill),
			Health: int(member.Health),
			Origin: member.Backstory.Origin,
		}
	}

	// Generate events
	eventQueue := events.NewQueue(seed, genre)
	eventStates := make([]eventState, eventN)
	for i := 0; i < eventN; i++ {
		x, y := i%mapW, (i/mapW)%mapH
		event := eventQueue.Generate(x, y, i)
		eventStates[i] = eventState{
			Title:    event.Title,
			Category: int(event.Category),
		}
	}

	// Create terrain and biome hashes
	terrainData := make([]byte, mapW*mapH)
	biomeData := make([]byte, mapW*mapH)
	for y := 0; y < mapH; y++ {
		for x := 0; x < mapW; x++ {
			tile := worldMap.Tiles[y][x]
			terrainData[y*mapW+x] = byte(tile.Terrain)
			biomeData[y*mapW+x] = byte(tile.Biome)
		}
	}
	terrainHash := sha256.Sum256(terrainData)
	biomeHash := sha256.Sum256(biomeData)

	return &fullState{
		World: worldState{
			OriginX:     worldMap.Origin.X,
			OriginY:     worldMap.Origin.Y,
			DestX:       worldMap.Destination.X,
			DestY:       worldMap.Destination.Y,
			TerrainHash: string(terrainHash[:16]),
			BiomeHash:   string(biomeHash[:16]),
		},
		Crew:   crewMembers,
		Events: eventStates,
	}
}

func generateStateForGenre(seed int64, genre engine.GenreID) *fullState {
	const mapW, mapH = 25, 18

	worldGen := world.NewGenerator(seed, genre)
	worldMap := worldGen.Generate(mapW, mapH)

	crewGen := crew.NewGenerator(seed, genre)
	crewMembers := make([]crewState, 4)
	for i := 0; i < 4; i++ {
		member := crewGen.Generate()
		crewMembers[i] = crewState{
			Name:   member.Name,
			Trait:  int(member.Trait),
			Skill:  int(member.Skill),
			Health: int(member.Health),
			Origin: member.Backstory.Origin,
		}
	}

	eventQueue := events.NewQueue(seed, genre)
	eventStates := make([]eventState, 20)
	for i := 0; i < 20; i++ {
		event := eventQueue.Generate(i%mapW, i/mapW, i)
		eventStates[i] = eventState{
			Title:    event.Title,
			Category: int(event.Category),
		}
	}

	terrainData := make([]byte, mapW*mapH)
	biomeData := make([]byte, mapW*mapH)
	for y := 0; y < mapH; y++ {
		for x := 0; x < mapW; x++ {
			tile := worldMap.Tiles[y][x]
			terrainData[y*mapW+x] = byte(tile.Terrain)
			biomeData[y*mapW+x] = byte(tile.Biome)
		}
	}
	terrainHash := sha256.Sum256(terrainData)
	biomeHash := sha256.Sum256(biomeData)

	return &fullState{
		World: worldState{
			OriginX:     worldMap.Origin.X,
			OriginY:     worldMap.Origin.Y,
			DestX:       worldMap.Destination.X,
			DestY:       worldMap.Destination.Y,
			TerrainHash: string(terrainHash[:16]),
			BiomeHash:   string(biomeHash[:16]),
		},
		Crew:   crewMembers,
		Events: eventStates,
	}
}

func hashState(t *testing.T, state *fullState) [32]byte {
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("Failed to marshal state: %v", err)
	}
	return sha256.Sum256(data)
}
