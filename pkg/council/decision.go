package council

import "github.com/opd-ai/voyage/pkg/engine"

// DecisionType represents the category of decision being voted on.
type DecisionType int

const (
	DecisionRoute DecisionType = iota // Shortcut vs safe path
	DecisionCamp                      // Rest vs push on
	DecisionTrade                     // Accept deal vs decline
	DecisionFight                     // Engage vs evade
	DecisionSplit                     // Divide party vs stay together
)

// String returns the base type name.
func (d DecisionType) String() string {
	return [...]string{"Route", "Camp", "Trade", "Fight", "Split"}[d]
}

// AllDecisionTypes returns all decision types.
func AllDecisionTypes() []DecisionType {
	return []DecisionType{DecisionRoute, DecisionCamp, DecisionTrade, DecisionFight, DecisionSplit}
}

// VoteOption represents one side of a binary decision.
type VoteOption int

const (
	OptionRisky VoteOption = iota // The riskier but potentially faster/better option
	OptionSafe                    // The safer but potentially slower/costlier option
)

// String returns the option name.
func (v VoteOption) String() string {
	return [...]string{"Risky", "Safe"}[v]
}

// Vote represents a single crew member's vote.
type Vote struct {
	CrewID     int
	CrewName   string
	Option     VoteOption
	Reasoning  string
	Confidence float64 // 0.0-1.0, how strongly they feel
}

// VoteResult represents the outcome of a council vote.
type VoteResult struct {
	Decision     DecisionType
	RiskyVotes   int
	SafeVotes    int
	Unanimous    bool
	PlayerChoice VoteOption
	Overruled    bool
	MoraleDelta  float64
	Votes        []Vote
}

// CouncilScene contains all text for the voting scene.
type CouncilScene struct {
	SceneName   string
	Opening     string
	RiskyOption string
	SafeOption  string
	Discussion  []string
	Closing     string
}

// decisionNames maps genre to decision type names.
var decisionNames = map[engine.GenreID]map[DecisionType]string{
	engine.GenreFantasy: {
		DecisionRoute: "Path Choice",
		DecisionCamp:  "Rest or March",
		DecisionTrade: "Merchant's Offer",
		DecisionFight: "Stand or Flee",
		DecisionSplit: "Divide the Party",
	},
	engine.GenreScifi: {
		DecisionRoute: "Course Correction",
		DecisionCamp:  "Power Down or Continue",
		DecisionTrade: "Contract Negotiation",
		DecisionFight: "Engage or Evade",
		DecisionSplit: "Deploy Away Team",
	},
	engine.GenreHorror: {
		DecisionRoute: "Which Way",
		DecisionCamp:  "Hide or Run",
		DecisionTrade: "Dark Bargain",
		DecisionFight: "Fight or Flight",
		DecisionSplit: "Split Up",
	},
	engine.GenreCyberpunk: {
		DecisionRoute: "Route Selection",
		DecisionCamp:  "Lay Low or Move",
		DecisionTrade: "The Deal",
		DecisionFight: "Assault or Ghost",
		DecisionSplit: "Solo Mission",
	},
	engine.GenrePostapoc: {
		DecisionRoute: "The Fork",
		DecisionCamp:  "Rest or Push",
		DecisionTrade: "Barter Terms",
		DecisionFight: "Engage or Evade",
		DecisionSplit: "Scout Ahead",
	},
}

// sceneNames maps genre to council scene names.
var sceneNames = map[engine.GenreID]string{
	engine.GenreFantasy:   "Campfire Debate",
	engine.GenreScifi:     "Bridge Briefing",
	engine.GenreHorror:    "Desperate Argument",
	engine.GenreCyberpunk: "Exec Meeting",
	engine.GenrePostapoc:  "Bonfire Council",
}
