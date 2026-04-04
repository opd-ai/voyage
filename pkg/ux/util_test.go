package ux

import (
	"testing"
)

func TestSplitWords(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"hello world", []string{"hello", "world"}},
		{"one", []string{"one"}},
		{"", nil},
		{"  spaces  between  ", []string{"spaces", "between"}},
		{"line1\nline2", []string{"line1", "line2"}},
		{"mixed spaces\nand\nnewlines", []string{"mixed", "spaces", "and", "newlines"}},
	}

	for _, tt := range tests {
		result := splitWords(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("splitWords(%q) = %v, want %v", tt.input, result, tt.expected)
			continue
		}
		for i := range result {
			if result[i] != tt.expected[i] {
				t.Errorf("splitWords(%q)[%d] = %q, want %q", tt.input, i, result[i], tt.expected[i])
			}
		}
	}
}

func TestGameStats(t *testing.T) {
	stats := GameStats{
		DaysTraveled:     30,
		DistanceTraveled: 150,
		CrewLost:         2,
		EventsResolved:   25,
		Victory:          true,
	}

	if stats.DaysTraveled != 30 {
		t.Errorf("expected DaysTraveled 30, got %d", stats.DaysTraveled)
	}
	if !stats.Victory {
		t.Error("expected Victory true")
	}
}
