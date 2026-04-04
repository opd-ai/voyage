package resources

import (
	"testing"

	"github.com/opd-ai/voyage/pkg/engine"
)

func TestNewResources(t *testing.T) {
	r := NewResources(engine.GenreFantasy)

	// Check initial levels
	if r.Get(ResourceFood) <= 0 {
		t.Error("Expected positive food level")
	}
	if r.Get(ResourceWater) <= 0 {
		t.Error("Expected positive water level")
	}

	// Check genre
	if r.Genre() != engine.GenreFantasy {
		t.Errorf("Expected fantasy genre, got %s", r.Genre())
	}
}

func TestResourceOperations(t *testing.T) {
	r := NewResources(engine.GenreFantasy)

	// Test Set and Get
	r.Set(ResourceFood, 50)
	if r.Get(ResourceFood) != 50 {
		t.Errorf("Expected food 50, got %f", r.Get(ResourceFood))
	}

	// Test clamping to max
	r.Set(ResourceFood, 200)
	if r.Get(ResourceFood) > r.GetMax(ResourceFood) {
		t.Error("Food should be clamped to max")
	}

	// Test clamping to zero
	r.Set(ResourceFood, -10)
	if r.Get(ResourceFood) < 0 {
		t.Error("Food should be clamped to zero")
	}

	// Test Add
	r.Set(ResourceFood, 50)
	r.Add(ResourceFood, 10)
	if r.Get(ResourceFood) != 60 {
		t.Errorf("Expected food 60, got %f", r.Get(ResourceFood))
	}

	// Test Consume
	r.Set(ResourceFood, 50)
	ok := r.Consume(ResourceFood, 30)
	if !ok {
		t.Error("Expected successful consume")
	}
	if r.Get(ResourceFood) != 20 {
		t.Errorf("Expected food 20, got %f", r.Get(ResourceFood))
	}

	// Test Consume failure
	ok = r.Consume(ResourceFood, 100)
	if ok {
		t.Error("Expected consume to fail")
	}
	if r.Get(ResourceFood) != 0 {
		t.Error("Food should be depleted after failed consume")
	}
}

func TestResourceStatus(t *testing.T) {
	r := NewResources(engine.GenreFantasy)

	// Normal
	r.Set(ResourceFood, 80)
	if r.GetStatus(ResourceFood) != StatusNormal {
		t.Error("Expected normal status at 80%")
	}

	// Low
	r.Set(ResourceFood, 20)
	if r.GetStatus(ResourceFood) != StatusLow {
		t.Error("Expected low status at 20%")
	}

	// Critical
	r.Set(ResourceFood, 5)
	if r.GetStatus(ResourceFood) != StatusCritical {
		t.Error("Expected critical status at 5%")
	}

	// Depleted
	r.Set(ResourceFood, 0)
	if r.GetStatus(ResourceFood) != StatusDepleted {
		t.Error("Expected depleted status at 0")
	}
}

func TestResourceNames(t *testing.T) {
	for _, genre := range engine.AllGenres() {
		for _, rt := range AllResourceTypes() {
			name := GetResourceName(rt, genre)
			if name == "" {
				t.Errorf("Empty name for resource %d, genre %s", rt, genre)
			}
		}
	}
}

func TestGenreSwitching(t *testing.T) {
	r := NewResources(engine.GenreFantasy)

	// Fantasy names
	if r.Name(ResourceCurrency) != "Gold" {
		t.Errorf("Expected 'Gold', got %s", r.Name(ResourceCurrency))
	}

	// Switch to scifi
	r.SetGenre(engine.GenreScifi)
	if r.Name(ResourceCurrency) != "Credits" {
		t.Errorf("Expected 'Credits', got %s", r.Name(ResourceCurrency))
	}
}

func TestConsumption(t *testing.T) {
	// Test daily consumption calculation
	consumption := CalculateDailyConsumption(ResourceFood, 4, 0)
	if consumption <= 0 {
		t.Error("Expected positive consumption")
	}

	// Test with desert terrain modifier for water
	normalWater := CalculateDailyConsumption(ResourceWater, 4, 0)
	desertWater := CalculateDailyConsumption(ResourceWater, 4, 3)
	if desertWater <= normalWater {
		t.Error("Desert should increase water consumption")
	}

	// Test movement cost
	cost := CalculateMovementCost(2, 1.0)
	if cost <= 0 {
		t.Error("Expected positive movement cost")
	}
}

func TestThresholdStatus(t *testing.T) {
	if !StatusLow.IsWarning() {
		t.Error("Low should be a warning")
	}
	if !StatusCritical.IsCritical() {
		t.Error("Critical should be critical")
	}
	if StatusNormal.IsWarning() {
		t.Error("Normal should not be a warning")
	}
}
