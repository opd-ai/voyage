package rendering

import (
	"testing"
)

func TestClamp(t *testing.T) {
	tests := []struct {
		v, min, max, expected int
	}{
		{5, 0, 10, 5},
		{-5, 0, 10, 0},
		{15, 0, 10, 10},
		{0, 0, 10, 0},
		{10, 0, 10, 10},
	}

	for _, tt := range tests {
		result := clamp(tt.v, tt.min, tt.max)
		if result != tt.expected {
			t.Errorf("clamp(%d, %d, %d) = %d, want %d",
				tt.v, tt.min, tt.max, result, tt.expected)
		}
	}
}
