package encounters

import (
	"github.com/opd-ai/voyage/pkg/crew"
	"github.com/opd-ai/voyage/pkg/engine"
	"github.com/opd-ai/voyage/pkg/procgen/seed"
	"github.com/opd-ai/voyage/pkg/resources"
	"github.com/opd-ai/voyage/pkg/vessel"
)

// Resolver handles tactical encounter resolution.
type Resolver struct {
	gen   *seed.Generator
	genre engine.GenreID
}

// NewResolver creates a new encounter resolver.
func NewResolver(masterSeed int64, genre engine.GenreID) *Resolver {
	return &Resolver{
		gen:   seed.NewGenerator(masterSeed, "encounter_resolver"),
		genre: genre,
	}
}

// SetGenre updates the resolver's genre.
func (r *Resolver) SetGenre(genre engine.GenreID) {
	r.genre = genre
}

// ResolvePhase resolves a single phase of an encounter.
func (r *Resolver) ResolvePhase(enc *Encounter, party *crew.Party) PhaseResult {
	if enc.State != StateResolution || enc.IsPaused {
		return PhaseResult{Success: false, Message: "Encounter not in resolution state"}
	}

	result := PhaseResult{
		PhaseNumber: enc.CurrentPhase + 1,
		TotalPhases: enc.MaxPhases,
	}

	// Calculate team effectiveness
	teamEff := r.calculateTeamEffectiveness(enc, party)

	// Roll for success
	difficultyMod := enc.Difficulty
	successChance := teamEff * (1 - difficultyMod*0.5)
	roll := r.gen.Float64()

	result.Success = roll < successChance
	result.EffectivenessScore = teamEff

	// Apply phase outcome
	if result.Success {
		enc.TotalProgress += teamEff * 0.5
		result.ProgressGained = teamEff * 0.5
		result.Message = r.getSuccessMessage(enc.Type)
	} else {
		enc.TotalDamage += difficultyMod * 10
		result.DamageTaken = difficultyMod * 10
		result.Message = r.getFailureMessage(enc.Type)
	}

	// Award skill experience
	result.ExpGains = r.calculateExpGains(enc, party, result.Success)

	enc.CurrentPhase++
	enc.TurnsElapsed++

	// Check if encounter is complete
	if enc.CurrentPhase >= enc.MaxPhases {
		enc.State = StateComplete
	}

	return result
}

// PhaseResult contains the result of a single phase resolution.
type PhaseResult struct {
	PhaseNumber        int
	TotalPhases        int
	Success            bool
	Message            string
	EffectivenessScore float64
	ProgressGained     float64
	DamageTaken        float64
	ExpGains           map[int]float64
}

// ResolveComplete finalizes an encounter and returns the full result.
func (r *Resolver) ResolveComplete(enc *Encounter, party *crew.Party) *EncounterResult {
	// Determine final outcome based on accumulated progress vs damage
	outcome := r.determineOutcome(enc)

	result := NewEncounterResult(outcome)
	result.TurnsElapsed = enc.TurnsElapsed
	result.Description = r.getOutcomeDescription(enc.Type, outcome)

	// Apply outcome-based rewards/penalties
	r.applyOutcomeEffects(result, enc, outcome)

	// Calculate crew damage distribution
	r.distributeCrewDamage(result, enc, party)

	return result
}

func (r *Resolver) determineOutcome(enc *Encounter) EncounterOutcome {
	// Guard against zero MaxPhases (H-015)
	if enc.MaxPhases == 0 {
		return OutcomeRetreat
	}
	progressRatio := enc.TotalProgress / float64(enc.MaxPhases)
	damageRatio := enc.TotalDamage / 100.0

	score := progressRatio - damageRatio

	switch {
	case score >= 0.6:
		return OutcomeVictory
	case score >= 0.2:
		return OutcomePartialSuccess
	case score >= -0.2:
		return OutcomeRetreat
	default:
		return OutcomeDefeat
	}
}

func (r *Resolver) calculateTeamEffectiveness(enc *Encounter, party *crew.Party) float64 {
	if party == nil || len(party.Living()) == 0 {
		return 0.1 // Minimum effectiveness
	}

	total, count := r.sumAssignmentEffectiveness(enc, party)
	if count == 0 {
		return 0.3 // Base effectiveness with no assignments
	}
	return total / float64(count)
}

// sumAssignmentEffectiveness calculates total effectiveness of all assigned crew.
func (r *Resolver) sumAssignmentEffectiveness(enc *Encounter, party *crew.Party) (float64, int) {
	total := 0.0
	count := 0
	for role, memberID := range enc.Assignments {
		member := party.Get(memberID)
		if member != nil && member.IsAlive {
			eff := r.calculateMemberEffectiveness(member, role, enc.OptimalRoles)
			total += eff
			count++
		}
	}
	return total, count
}

// calculateMemberEffectiveness calculates effectiveness for a single crew member in a role.
func (r *Resolver) calculateMemberEffectiveness(member *crew.CrewMember, role EncounterRole, optimalRoles []EncounterRole) float64 {
	eff := CalculateRoleEffectiveness(member, role)
	if isOptimalRole(role, optimalRoles) {
		eff *= 1.2
	}
	return eff
}

// isOptimalRole checks if a role is in the list of optimal roles.
func isOptimalRole(role EncounterRole, optimalRoles []EncounterRole) bool {
	for _, optRole := range optimalRoles {
		if role == optRole {
			return true
		}
	}
	return false
}

func (r *Resolver) calculateExpGains(enc *Encounter, party *crew.Party, success bool) map[int]float64 {
	gains := make(map[int]float64)
	baseExp := 10.0
	if success {
		baseExp = 20.0
	}

	for _, memberID := range enc.Assignments {
		member := party.Get(memberID)
		if member != nil && member.IsAlive {
			gains[memberID] = baseExp * (1 + enc.Difficulty)
		}
	}

	return gains
}

func (r *Resolver) applyOutcomeEffects(result *EncounterResult, enc *Encounter, outcome EncounterOutcome) {
	baseMod := enc.Difficulty

	switch outcome {
	case OutcomeVictory:
		result.MoraleDelta = 10 + baseMod*5
		result.CurrencyDelta = 20 * (1 + baseMod)
		result.FoodDelta = 5 * r.gen.Float64()
	case OutcomePartialSuccess:
		result.MoraleDelta = 2
		result.CurrencyDelta = 5
		result.VesselDamage = 5 * baseMod
	case OutcomeRetreat:
		result.MoraleDelta = -5
		result.FuelDelta = -10 * baseMod
		result.VesselDamage = 10 * baseMod
	case OutcomeDefeat:
		result.MoraleDelta = -15
		result.FoodDelta = -10 * baseMod
		result.CurrencyDelta = -20 * baseMod
		result.VesselDamage = 25 * baseMod
	}
}

func (r *Resolver) distributeCrewDamage(result *EncounterResult, enc *Encounter, party *crew.Party) {
	if enc.TotalDamage <= 0 {
		return
	}

	// Distribute damage among assigned crew
	assignedCount := len(enc.Assignments)
	if assignedCount == 0 {
		return
	}

	damagePerCrew := enc.TotalDamage / float64(assignedCount)
	for _, memberID := range enc.Assignments {
		// Add some variance
		variance := 0.8 + r.gen.Float64()*0.4 // 0.8-1.2
		result.CrewDamage[memberID] = damagePerCrew * variance
	}
}

func (r *Resolver) getSuccessMessage(encType EncounterType) string {
	messages := phaseSuccessMessages[r.genre]
	if messages == nil {
		messages = phaseSuccessMessages[engine.GenreFantasy]
	}
	typeMessages := messages[encType]
	if len(typeMessages) == 0 {
		return "The phase succeeds."
	}
	return seed.Choice(r.gen, typeMessages)
}

func (r *Resolver) getFailureMessage(encType EncounterType) string {
	messages := phaseFailureMessages[r.genre]
	if messages == nil {
		messages = phaseFailureMessages[engine.GenreFantasy]
	}
	typeMessages := messages[encType]
	if len(typeMessages) == 0 {
		return "The phase fails."
	}
	return seed.Choice(r.gen, typeMessages)
}

func (r *Resolver) getOutcomeDescription(encType EncounterType, outcome EncounterOutcome) string {
	descriptions := outcomeDescriptions[r.genre]
	if descriptions == nil {
		descriptions = outcomeDescriptions[engine.GenreFantasy]
	}
	outcomeDescs := descriptions[outcome]
	if len(outcomeDescs) == 0 {
		return "The encounter concludes."
	}
	return seed.Choice(r.gen, outcomeDescs)
}

var phaseSuccessMessages = map[engine.GenreID]map[EncounterType][]string{
	engine.GenreFantasy: {
		TypeAmbush:      {"Your defenders hold the line.", "The enemy falters."},
		TypeNegotiation: {"Your words find their mark.", "Understanding is reached."},
		TypeRace:        {"You gain ground.", "The pursuers fall behind."},
		TypeCrisis:      {"The situation stabilizes.", "Quick thinking saves the day."},
		TypePuzzle:      {"The pieces fall into place.", "Insight strikes."},
	},
	engine.GenreScifi: {
		TypeAmbush:      {"Hostile fire suppressed.", "Shields holding."},
		TypeNegotiation: {"Communication established.", "Terms accepted."},
		TypeRace:        {"Gaining distance.", "Pursuit falling behind."},
		TypeCrisis:      {"Systems stabilizing.", "Emergency contained."},
		TypePuzzle:      {"Data decoded.", "Pattern recognized."},
	},
	engine.GenreHorror: {
		TypeAmbush:      {"You push them back.", "A moment to breathe."},
		TypeNegotiation: {"They lower their weapons.", "Trust, for now."},
		TypeRace:        {"You're pulling ahead.", "The sounds grow distant."},
		TypeCrisis:      {"Under control. Barely.", "You act in time."},
		TypePuzzle:      {"It clicks.", "The way becomes clear."},
	},
	engine.GenreCyberpunk: {
		TypeAmbush:      {"Hostiles down.", "You've got the edge."},
		TypeNegotiation: {"Deal sweetened.", "They're buying it."},
		TypeRace:        {"Losing them.", "Clean getaway in sight."},
		TypeCrisis:      {"Situation contained.", "Crisis averted."},
		TypePuzzle:      {"ICE breached.", "Access granted."},
	},
	engine.GenrePostapoc: {
		TypeAmbush:      {"Raiders retreat.", "You hold your ground."},
		TypeNegotiation: {"They see reason.", "An accord is struck."},
		TypeRace:        {"You're outrunning it.", "Safety draws near."},
		TypeCrisis:      {"Disaster averted.", "You stabilize things."},
		TypePuzzle:      {"The lock gives.", "Pre-war tech yields."},
	},
}

var phaseFailureMessages = map[engine.GenreID]map[EncounterType][]string{
	engine.GenreFantasy: {
		TypeAmbush:      {"Your line wavers.", "The enemy presses forward."},
		TypeNegotiation: {"Your words fall flat.", "Tensions rise."},
		TypeRace:        {"They gain on you.", "Your escape narrows."},
		TypeCrisis:      {"Things get worse.", "The situation deteriorates."},
		TypePuzzle:      {"The solution eludes you.", "Time slips away."},
	},
	engine.GenreScifi: {
		TypeAmbush:      {"Shields taking damage.", "Hostiles advancing."},
		TypeNegotiation: {"Signal rejected.", "Negotiations stall."},
		TypeRace:        {"They're closing.", "Engine strain detected."},
		TypeCrisis:      {"Systems critical.", "Emergency worsening."},
		TypePuzzle:      {"Access denied.", "Encryption holds."},
	},
	engine.GenreHorror: {
		TypeAmbush:      {"They break through.", "Too many of them."},
		TypeNegotiation: {"They don't believe you.", "Fear takes hold."},
		TypeRace:        {"They're catching up.", "You can hear them."},
		TypeCrisis:      {"It's spreading.", "Things are getting worse."},
		TypePuzzle:      {"Still locked.", "Time's running out."},
	},
	engine.GenreCyberpunk: {
		TypeAmbush:      {"Taking fire.", "They've got numbers."},
		TypeNegotiation: {"They're not buying.", "Deal's going south."},
		TypeRace:        {"They're on you.", "Can't shake them."},
		TypeCrisis:      {"Going critical.", "Situation escalating."},
		TypePuzzle:      {"Locked out.", "ICE holds."},
	},
	engine.GenrePostapoc: {
		TypeAmbush:      {"They push forward.", "You're outgunned."},
		TypeNegotiation: {"They're hostile.", "Words aren't working."},
		TypeRace:        {"It's gaining.", "No time left."},
		TypeCrisis:      {"It's getting bad.", "You're losing control."},
		TypePuzzle:      {"Still sealed.", "Old tech resists."},
	},
}

var outcomeDescriptions = map[engine.GenreID]map[EncounterOutcome][]string{
	engine.GenreFantasy: {
		OutcomeVictory:        {"A glorious victory!", "The day is won."},
		OutcomePartialSuccess: {"Victory, but at a cost.", "You prevail, barely."},
		OutcomeRetreat:        {"You escape with your lives.", "A tactical retreat."},
		OutcomeDefeat:         {"A bitter defeat.", "All is lost."},
	},
	engine.GenreScifi: {
		OutcomeVictory:        {"Mission accomplished.", "Objectives achieved."},
		OutcomePartialSuccess: {"Partial mission success.", "Objectives partially met."},
		OutcomeRetreat:        {"Emergency retreat completed.", "Disengaged safely."},
		OutcomeDefeat:         {"Mission failed.", "Critical failure."},
	},
	engine.GenreHorror: {
		OutcomeVictory:        {"You survived.", "Against all odds, you made it."},
		OutcomePartialSuccess: {"Survival has a price.", "You made it, barely."},
		OutcomeRetreat:        {"You got away.", "Escape, for now."},
		OutcomeDefeat:         {"The horror takes its toll.", "You couldn't stop it."},
	},
	engine.GenreCyberpunk: {
		OutcomeVictory:        {"Clean run, choom.", "Mission complete."},
		OutcomePartialSuccess: {"Got messy, but you're out.", "Job done, mostly."},
		OutcomeRetreat:        {"Had to bail.", "You're out, job's not."},
		OutcomeDefeat:         {"Flatlined.", "Run failed."},
	},
	engine.GenrePostapoc: {
		OutcomeVictory:        {"You survived another day.", "Victory in the wastes."},
		OutcomePartialSuccess: {"Survival isn't pretty.", "You made it through."},
		OutcomeRetreat:        {"You got away.", "Escape, at least."},
		OutcomeDefeat:         {"The wasteland claims its due.", "You couldn't make it."},
	},
}

// ApplyResult applies encounter results to game state.
func (r *Resolver) ApplyResult(result *EncounterResult, res *resources.Resources, party *crew.Party, v *vessel.Vessel) []string {
	applyResourceDeltas(result, res)
	applyVesselDamage(result, v)
	deaths := applyCrewDamage(result, party)
	applyCrewExpGains(result, party)
	return deaths
}

// applyResourceDeltas applies all resource changes from an encounter result.
func applyResourceDeltas(result *EncounterResult, res *resources.Resources) {
	deltas := map[resources.ResourceType]float64{
		resources.ResourceFood:     result.FoodDelta,
		resources.ResourceWater:    result.WaterDelta,
		resources.ResourceFuel:     result.FuelDelta,
		resources.ResourceMedicine: result.MedicineDelta,
		resources.ResourceMorale:   result.MoraleDelta,
		resources.ResourceCurrency: result.CurrencyDelta,
	}
	for resType, delta := range deltas {
		if delta != 0 {
			res.Add(resType, delta)
		}
	}
}

// applyVesselDamage applies vessel damage from an encounter result.
func applyVesselDamage(result *EncounterResult, v *vessel.Vessel) {
	if result.VesselDamage > 0 {
		v.TakeDamage(result.VesselDamage)
	}
}

// applyCrewDamage applies damage to crew members and returns names of any deaths.
func applyCrewDamage(result *EncounterResult, party *crew.Party) []string {
	var deaths []string
	for memberID, damage := range result.CrewDamage {
		member := party.Get(memberID)
		if member != nil && member.TakeDamage(damage) {
			deaths = append(deaths, member.Name)
		}
	}
	return deaths
}

// applyCrewExpGains applies skill experience gains to crew members.
func applyCrewExpGains(result *EncounterResult, party *crew.Party) {
	for memberID, exp := range result.SkillExpGains {
		member := party.Get(memberID)
		if member != nil {
			member.GainSkillExp(exp)
		}
	}
}
