package npc

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestAllNPCTypes(t *testing.T) {
	types := AllNPCTypes()
	if len(types) != 6 {
		t.Errorf("Expected 6 NPC types, got %d", len(types))
	}
}

func TestNPCTypeName(t *testing.T) {
	tests := []struct {
		npcType NPCType
		genre   engine.GenreID
		want    string
	}{
		{TypeTrader, engine.GenreFantasy, "Merchant"},
		{TypeTrader, engine.GenreScifi, "Trader"},
		{TypeBandit, engine.GenreHorror, "Raider"},
		{TypeScout, engine.GenreCyberpunk, "Netrunner"},
		{TypeRefugee, engine.GenrePostapoc, "Wanderer"},
	}

	for _, tt := range tests {
		got := NPCTypeName(tt.npcType, tt.genre)
		if got != tt.want {
			t.Errorf("NPCTypeName(%d, %s) = %q, want %q", tt.npcType, tt.genre, got, tt.want)
		}
	}
}

func TestAlignmentName(t *testing.T) {
	tests := []struct {
		alignment Alignment
		want      string
	}{
		{AlignmentHostile, "Hostile"},
		{AlignmentNeutral, "Neutral"},
		{AlignmentFriendly, "Friendly"},
	}

	for _, tt := range tests {
		got := AlignmentName(tt.alignment)
		if got != tt.want {
			t.Errorf("AlignmentName(%d) = %q, want %q", tt.alignment, got, tt.want)
		}
	}
}

func TestNewNPC(t *testing.T) {
	npc := NewNPC(1, "Test", TypeTrader, engine.GenreFantasy)

	if npc.ID != 1 {
		t.Errorf("Expected ID 1, got %d", npc.ID)
	}
	if npc.Name != "Test" {
		t.Errorf("Expected name 'Test', got %q", npc.Name)
	}
	if npc.NPCType != TypeTrader {
		t.Errorf("Expected TypeTrader, got %d", npc.NPCType)
	}
	if npc.Alignment != AlignmentNeutral {
		t.Errorf("Expected AlignmentNeutral, got %d", npc.Alignment)
	}
}

func TestNPCSetGenre(t *testing.T) {
	npc := NewNPC(1, "Test", TypeTrader, engine.GenreFantasy)
	npc.SetGenre(engine.GenreScifi)

	if npc.Genre != engine.GenreScifi {
		t.Errorf("Expected scifi, got %s", npc.Genre)
	}
}

func TestNPCIsHostile(t *testing.T) {
	npc := NewNPC(1, "Test", TypeBandit, engine.GenreFantasy)
	npc.Alignment = AlignmentHostile

	if !npc.IsHostile() {
		t.Error("Expected bandit to be hostile")
	}

	npc.Alignment = AlignmentNeutral
	if npc.IsHostile() {
		t.Error("Neutral NPC should not be hostile")
	}
}

func TestNPCCanTrade(t *testing.T) {
	trader := NewNPC(1, "Test", TypeTrader, engine.GenreFantasy)
	trader.Alignment = AlignmentNeutral

	if !trader.CanTrade() {
		t.Error("Neutral trader should be able to trade")
	}

	trader.Alignment = AlignmentHostile
	if trader.CanTrade() {
		t.Error("Hostile trader should not be able to trade")
	}

	bandit := NewNPC(2, "Test", TypeBandit, engine.GenreFantasy)
	bandit.Alignment = AlignmentNeutral
	if bandit.CanTrade() {
		t.Error("Bandit should not be able to trade")
	}
}

func TestGetDefaultAlignment(t *testing.T) {
	tests := []struct {
		npcType NPCType
		want    Alignment
	}{
		{TypeTrader, AlignmentNeutral},
		{TypeRefugee, AlignmentFriendly},
		{TypeBandit, AlignmentHostile},
	}

	for _, tt := range tests {
		got := GetDefaultAlignment(tt.npcType)
		if got != tt.want {
			t.Errorf("GetDefaultAlignment(%d) = %d, want %d", tt.npcType, got, tt.want)
		}
	}
}

func TestGenerator(t *testing.T) {
	gen := NewGenerator(12345, engine.GenreFantasy)

	npc := gen.Generate(TypeTrader)
	if npc == nil {
		t.Fatal("Expected non-nil NPC")
	}
	if npc.Name == "" {
		t.Error("NPC should have a name")
	}
	if npc.Description == "" {
		t.Error("NPC should have a description")
	}
	if len(npc.Dialogue) == 0 {
		t.Error("NPC should have dialogue")
	}
	if npc.NPCType != TypeTrader {
		t.Errorf("Expected TypeTrader, got %d", npc.NPCType)
	}
	if len(npc.TradeGoods) == 0 {
		t.Error("Trader should have trade goods")
	}
}

func TestGeneratorSetGenre(t *testing.T) {
	gen := NewGenerator(12345, engine.GenreFantasy)
	gen.SetGenre(engine.GenreScifi)

	npc := gen.Generate(TypeTrader)
	if npc.Genre != engine.GenreScifi {
		t.Errorf("Expected scifi genre, got %s", npc.Genre)
	}
}

func TestGeneratorRandom(t *testing.T) {
	gen := NewGenerator(12345, engine.GenreFantasy)

	npc := gen.GenerateRandom()
	if npc == nil {
		t.Fatal("Expected non-nil NPC")
	}
}

func TestGeneratorEncounter(t *testing.T) {
	gen := NewGenerator(12345, engine.GenreFantasy)

	npc := gen.GenerateEncounter()
	if npc == nil {
		t.Fatal("Expected non-nil NPC")
	}
}

func TestGeneratorDeterminism(t *testing.T) {
	gen1 := NewGenerator(42, engine.GenreFantasy)
	gen2 := NewGenerator(42, engine.GenreFantasy)

	npc1 := gen1.Generate(TypeTrader)
	npc2 := gen2.Generate(TypeTrader)

	if npc1.Name != npc2.Name {
		t.Errorf("Same seed should produce same name: %q vs %q", npc1.Name, npc2.Name)
	}
	if npc1.Alignment != npc2.Alignment {
		t.Errorf("Same seed should produce same alignment: %d vs %d", npc1.Alignment, npc2.Alignment)
	}
}

func TestAllGenresHaveContent(t *testing.T) {
	genres := engine.AllGenres()

	for _, genre := range genres {
		gen := NewGenerator(12345, genre)

		for _, npcType := range AllNPCTypes() {
			npc := gen.Generate(npcType)

			// Check type name exists
			name := NPCTypeName(npcType, genre)
			if name == "" {
				t.Errorf("Missing type name for genre %s, type %d", genre, npcType)
			}

			// Check NPC was generated properly
			if npc.Name == "" {
				t.Errorf("NPC missing name for genre %s, type %d", genre, npcType)
			}
			if npc.Description == "" {
				t.Errorf("NPC missing description for genre %s, type %d", genre, npcType)
			}
			if len(npc.Dialogue) == 0 {
				t.Errorf("NPC missing dialogue for genre %s, type %d", genre, npcType)
			}
		}
	}
}

func TestTradeGoodsHaveContent(t *testing.T) {
	genres := engine.AllGenres()

	for _, genre := range genres {
		gen := NewGenerator(12345, genre)
		npc := gen.Generate(TypeTrader)

		if len(npc.TradeGoods) == 0 {
			t.Errorf("Trader missing trade goods for genre %s", genre)
			continue
		}

		for _, good := range npc.TradeGoods {
			if good.Name == "" {
				t.Errorf("Trade good missing name for genre %s", genre)
			}
			if good.Price <= 0 {
				t.Errorf("Trade good has invalid price for genre %s", genre)
			}
			if good.Quantity <= 0 {
				t.Errorf("Trade good has invalid quantity for genre %s", genre)
			}
		}
	}
}
