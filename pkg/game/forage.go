package game

import (
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
	"github.com/opd-ai/voyage/pkg/procgen/world"
	"github.com/opd-ai/voyage/pkg/resources"
)

// ForageOutcome represents the type of result from a forage attempt.
type ForageOutcome int

const (
	// ForageNothing indicates the forage attempt found nothing.
	ForageNothing ForageOutcome = iota
	// ForageFood indicates food was found.
	ForageFood
	// ForageWater indicates water was found.
	ForageWater
	// ForageFuel indicates fuel/supplies were found.
	ForageFuel
	// ForageMedicine indicates medicine was found.
	ForageMedicine
	// ForageCurrency indicates valuable items were found.
	ForageCurrency
	// ForageParts indicates repair materials were found.
	ForageParts
	// ForageEncounter indicates an encounter was triggered.
	ForageEncounter
)

// ForageResult contains the details of a forage attempt.
type ForageResult struct {
	Outcome      ForageOutcome
	Description  string
	FoodGain     float64
	WaterGain    float64
	FuelGain     float64
	MedsGain     float64
	CurrencyGain float64
	PartsGain    int
	TurnsSpent   int
	EncounterID  int // If Outcome is ForageEncounter, this is the event ID
}

// ForageManager handles gathering resources at wilderness and ruin tiles.
type ForageManager struct {
	gen          *seed.Generator
	genre        engine.GenreID
	foragedTiles map[string]int // Map of "x,y" -> times foraged
}

// NewForageManager creates a new forage manager.
func NewForageManager(masterSeed int64, genre engine.GenreID) *ForageManager {
	return &ForageManager{
		gen:          seed.NewGenerator(masterSeed, "forage"),
		genre:        genre,
		foragedTiles: make(map[string]int),
	}
}

// SetGenre changes the forage manager's genre for text generation.
func (fm *ForageManager) SetGenre(genre engine.GenreID) {
	fm.genre = genre
}

// Genre returns the current genre.
func (fm *ForageManager) Genre() engine.GenreID {
	return fm.genre
}

// CanForage checks if the player can forage at the given tile.
func (fm *ForageManager) CanForage(tile *world.Tile) bool {
	if tile == nil {
		return false
	}
	// Can forage at wilderness, ruin, or landmark tiles (except origin/destination)
	switch tile.Terrain {
	case world.TerrainForest, world.TerrainPlains, world.TerrainDesert,
		world.TerrainMountain, world.TerrainRuin:
		return true
	}
	// Check for landmark
	if tile.Landmark != nil {
		switch tile.Landmark.Type {
		case world.LandmarkRuins, world.LandmarkOutpost, world.LandmarkShrine:
			return true
		}
	}
	return false
}

// ActionName returns the genre-appropriate name for the forage action.
func (fm *ForageManager) ActionName() string {
	names := map[engine.GenreID]string{
		engine.GenreFantasy:   "Forage",
		engine.GenreScifi:     "Salvage",
		engine.GenreHorror:    "Scavenge",
		engine.GenreCyberpunk: "Jack Data",
		engine.GenrePostapoc:  "Strip Wreck",
	}
	if name, ok := names[fm.genre]; ok {
		return name
	}
	return "Forage"
}

// Forage attempts to gather resources at the given tile.
func (fm *ForageManager) Forage(tile *world.Tile, turn int) *ForageResult {
	if !fm.CanForage(tile) {
		return &ForageResult{
			Outcome:     ForageNothing,
			Description: "This location cannot be searched.",
			TurnsSpent:  0,
		}
	}

	// Create deterministic seed from position and turn
	posKey := fm.tileKey(tile.X, tile.Y)
	localSeed := fm.gen.Master() + int64(tile.X*1000+tile.Y*100+turn)
	localGen := seed.NewGenerator(localSeed, posKey)

	// Get diminishing returns factor
	timesForaged := fm.foragedTiles[posKey]
	yieldMod := fm.calculateYieldModifier(timesForaged)

	// Update forage count
	fm.foragedTiles[posKey]++

	// Determine outcome based on terrain and RNG
	outcome := fm.rollOutcome(localGen, tile)
	result := fm.generateResult(localGen, outcome, tile, yieldMod)
	result.TurnsSpent = 1

	return result
}

// tileKey creates a unique string key for a tile position.
func (fm *ForageManager) tileKey(x, y int) string {
	return string(rune(x)) + "," + string(rune(y))
}

// calculateYieldModifier returns a multiplier based on times foraged.
// Returns 1.0 for first forage, decreasing with each subsequent attempt.
func (fm *ForageManager) calculateYieldModifier(timesForaged int) float64 {
	if timesForaged == 0 {
		return 1.0
	}
	// Each subsequent forage yields 30% less, minimum 10%
	mod := 1.0
	for i := 0; i < timesForaged; i++ {
		mod *= 0.7
	}
	if mod < 0.1 {
		return 0.1
	}
	return mod
}

// rollOutcome determines what type of result the forage produces.
func (fm *ForageManager) rollOutcome(gen *seed.Generator, tile *world.Tile) ForageOutcome {
	weights := fm.getOutcomeWeights(tile)
	outcomes := []ForageOutcome{
		ForageNothing, ForageFood, ForageWater, ForageFuel,
		ForageMedicine, ForageCurrency, ForageParts, ForageEncounter,
	}
	return seed.WeightedChoice(gen, outcomes, weights)
}

// getOutcomeWeights returns probability weights based on terrain.
func (fm *ForageManager) getOutcomeWeights(tile *world.Tile) []float64 {
	// Base weights: nothing, food, water, fuel, meds, currency, parts, encounter
	switch tile.Terrain {
	case world.TerrainForest:
		return []float64{0.15, 0.35, 0.15, 0.05, 0.10, 0.05, 0.05, 0.10}
	case world.TerrainPlains:
		return []float64{0.25, 0.30, 0.10, 0.05, 0.05, 0.05, 0.05, 0.15}
	case world.TerrainDesert:
		return []float64{0.35, 0.10, 0.05, 0.10, 0.05, 0.10, 0.10, 0.15}
	case world.TerrainMountain:
		return []float64{0.30, 0.10, 0.15, 0.05, 0.10, 0.10, 0.10, 0.10}
	case world.TerrainRuin:
		return []float64{0.20, 0.05, 0.05, 0.15, 0.15, 0.15, 0.15, 0.10}
	default:
		return []float64{0.30, 0.20, 0.10, 0.05, 0.05, 0.05, 0.05, 0.20}
	}
}

// generateResult creates a ForageResult for the given outcome.
func (fm *ForageManager) generateResult(gen *seed.Generator, outcome ForageOutcome, tile *world.Tile, yieldMod float64) *ForageResult {
	result := &ForageResult{
		Outcome: outcome,
	}

	baseYield := 10 + float64(gen.Intn(16)) // 10-25 base yield
	yield := baseYield * yieldMod

	switch outcome {
	case ForageNothing:
		result.Description = fm.nothingDescription(gen)
	case ForageFood:
		result.FoodGain = yield
		result.Description = fm.foodDescription(gen, result.FoodGain)
	case ForageWater:
		result.WaterGain = yield
		result.Description = fm.waterDescription(gen, result.WaterGain)
	case ForageFuel:
		result.FuelGain = yield * 0.8 // Fuel slightly rarer
		result.Description = fm.fuelDescription(gen, result.FuelGain)
	case ForageMedicine:
		result.MedsGain = yield * 0.5 // Medicine is rare
		result.Description = fm.medsDescription(gen, result.MedsGain)
	case ForageCurrency:
		result.CurrencyGain = yield * 1.5 // Currency yields more
		result.Description = fm.currencyDescription(gen, result.CurrencyGain)
	case ForageParts:
		result.PartsGain = int(yield / 5) // 2-5 parts typical
		if result.PartsGain < 1 {
			result.PartsGain = 1
		}
		result.Description = fm.partsDescription(gen, result.PartsGain)
	case ForageEncounter:
		result.EncounterID = gen.Intn(1000) // Random encounter ID
		result.Description = fm.encounterDescription(gen)
	}

	return result
}

// ApplyResult applies a forage result to the player's resources.
func (fm *ForageManager) ApplyResult(result *ForageResult, res *resources.Resources) {
	if result.FoodGain > 0 {
		res.Add(resources.ResourceFood, result.FoodGain)
	}
	if result.WaterGain > 0 {
		res.Add(resources.ResourceWater, result.WaterGain)
	}
	if result.FuelGain > 0 {
		res.Add(resources.ResourceFuel, result.FuelGain)
	}
	if result.MedsGain > 0 {
		res.Add(resources.ResourceMedicine, result.MedsGain)
	}
	if result.CurrencyGain > 0 {
		res.Add(resources.ResourceCurrency, result.CurrencyGain)
	}
	// Parts would be added to cargo, but that's handled elsewhere
}

// GetForageCount returns how many times a tile has been foraged.
func (fm *ForageManager) GetForageCount(x, y int) int {
	return fm.foragedTiles[fm.tileKey(x, y)]
}

// ResetTile resets the forage count for a tile (e.g., after time passes).
func (fm *ForageManager) ResetTile(x, y int) {
	delete(fm.foragedTiles, fm.tileKey(x, y))
}

// Description generation methods
func (fm *ForageManager) nothingDescription(gen *seed.Generator) string {
	texts := nothingTexts[fm.genre]
	return seed.Choice(gen, texts)
}

func (fm *ForageManager) foodDescription(gen *seed.Generator, amount float64) string {
	texts := foodTexts[fm.genre]
	return seed.Choice(gen, texts)
}

func (fm *ForageManager) waterDescription(gen *seed.Generator, amount float64) string {
	texts := waterTexts[fm.genre]
	return seed.Choice(gen, texts)
}

func (fm *ForageManager) fuelDescription(gen *seed.Generator, amount float64) string {
	texts := fuelTexts[fm.genre]
	return seed.Choice(gen, texts)
}

func (fm *ForageManager) medsDescription(gen *seed.Generator, amount float64) string {
	texts := medsTexts[fm.genre]
	return seed.Choice(gen, texts)
}

func (fm *ForageManager) currencyDescription(gen *seed.Generator, amount float64) string {
	texts := currencyTexts[fm.genre]
	return seed.Choice(gen, texts)
}

func (fm *ForageManager) partsDescription(gen *seed.Generator, count int) string {
	texts := partsTexts[fm.genre]
	return seed.Choice(gen, texts)
}

func (fm *ForageManager) encounterDescription(gen *seed.Generator) string {
	texts := encounterTexts[fm.genre]
	return seed.Choice(gen, texts)
}

// Genre-specific text tables
var nothingTexts = map[engine.GenreID][]string{
	engine.GenreFantasy:   {"You search thoroughly but find nothing of value.", "The area has been picked clean.", "Your search yields only dust and disappointment."},
	engine.GenreScifi:     {"Scanners detect nothing salvageable.", "The area is depleted.", "No useful materials found."},
	engine.GenreHorror:    {"You find only decay and despair.", "Nothing here but death.", "The search turns up empty."},
	engine.GenreCyberpunk: {"The data nodes are fried.", "Nothing worth jacking here.", "The site's been cleaned out."},
	engine.GenrePostapoc:  {"Nothing but rust and bones.", "The scavengers got here first.", "Worthless junk, all of it."},
}

var foodTexts = map[engine.GenreID][]string{
	engine.GenreFantasy:   {"You gather edible berries and roots.", "Wild game provides fresh meat.", "Herbs and mushrooms fill your pack."},
	engine.GenreScifi:     {"Emergency ration cache discovered.", "Hydroponic supplies recovered.", "Preserved food stores found."},
	engine.GenreHorror:    {"Canned goods, still sealed.", "MREs in a supply closet.", "Non-perishables, enough to last."},
	engine.GenreCyberpunk: {"Nutrient paste stockpile.", "Soy-protein rations acquired.", "Synthetic food cache found."},
	engine.GenrePostapoc:  {"Pre-war cans in good condition.", "Dried goods in a bunker.", "Preserved provisions discovered."},
}

var waterTexts = map[engine.GenreID][]string{
	engine.GenreFantasy:   {"A clean spring provides fresh water.", "Rainwater collected in barrels.", "A hidden well yields pure water."},
	engine.GenreScifi:     {"Water recycler still functional.", "Ice deposits extracted.", "Condensation tanks operational."},
	engine.GenreHorror:    {"Bottled water in the basement.", "Water purification tablets found.", "A working tap, surprisingly."},
	engine.GenreCyberpunk: {"Filtered water reserves.", "Corporate-grade hydration packs.", "Clean water, worth its weight in creds."},
	engine.GenrePostapoc:  {"Uncontaminated water source.", "Purified water cache.", "Clean water, a rare find."},
}

var fuelTexts = map[engine.GenreID][]string{
	engine.GenreFantasy:   {"Your animals find good grazing.", "Oats and hay for the journey.", "Rest restores your stamina."},
	engine.GenreScifi:     {"Fuel cells in storage.", "Power core with charge remaining.", "Energy reserves recovered."},
	engine.GenreHorror:    {"Gas cans in a shed.", "Diesel in the tanks.", "Fuel siphoned from wrecks."},
	engine.GenreCyberpunk: {"Battery packs still charged.", "Power cells recovered.", "Energy cells, good capacity."},
	engine.GenrePostapoc:  {"Diesel in rusty drums.", "Fuel, though it smells bad.", "Gasoline, precious drops."},
}

var medsTexts = map[engine.GenreID][]string{
	engine.GenreFantasy:   {"Healing herbs and poultices.", "A healer's abandoned kit.", "Medicinal plants gathered."},
	engine.GenreScifi:     {"Medical supplies intact.", "Auto-injectors recovered.", "Med-bay supplies found."},
	engine.GenreHorror:    {"First aid kit, barely used.", "Antibiotics and bandages.", "Medical supplies, thank god."},
	engine.GenreCyberpunk: {"Stims and trauma patches.", "Med-tech supplies.", "Black market pharmaceuticals."},
	engine.GenrePostapoc:  {"Pre-war medicine cache.", "Surgical supplies found.", "Meds, worth killing for."},
}

var currencyTexts = map[engine.GenreID][]string{
	engine.GenreFantasy:   {"A hidden coin purse.", "Gems in a secret compartment.", "Trade goods of value."},
	engine.GenreScifi:     {"Credit chips recovered.", "Valuable data crystals.", "Trade commodities found."},
	engine.GenreHorror:    {"Valuables from the dead.", "A lockbox with jewelry.", "Cash and tradeable goods."},
	engine.GenreCyberpunk: {"Credsticks with balance.", "Encrypted data worth selling.", "Corporate scrip found."},
	engine.GenrePostapoc:  {"Bottle caps aplenty.", "Tradeable ammunition.", "Pre-war currency, still valued."},
}

var partsTexts = map[engine.GenreID][]string{
	engine.GenreFantasy:   {"Spare wagon parts.", "Tools and repair materials.", "Useful hardware salvaged."},
	engine.GenreScifi:     {"Hull plating recovered.", "Replacement components found.", "Repair drones salvaged."},
	engine.GenreHorror:    {"Auto parts in the garage.", "Repair supplies scavenged.", "Tools and spare parts."},
	engine.GenreCyberpunk: {"Chrome and circuitry.", "Replacement components.", "Tech parts, good condition."},
	engine.GenrePostapoc:  {"Scrap metal, good quality.", "Mechanical parts salvaged.", "Spare parts for days."},
}

var encounterTexts = map[engine.GenreID][]string{
	engine.GenreFantasy:   {"Something stirs in the shadows...", "You are not alone here!", "A figure emerges from hiding!"},
	engine.GenreScifi:     {"Movement detected!", "Life signs approaching!", "Contact! Unknown vessel!"},
	engine.GenreHorror:    {"Oh god, they're here!", "Something's coming!", "You've drawn attention!"},
	engine.GenreCyberpunk: {"You've been spotted!", "ICE detected, incoming!", "Security alert triggered!"},
	engine.GenrePostapoc:  {"Hostiles inbound!", "You've been followed!", "Company's coming!"},
}
