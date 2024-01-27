package main

import "testing"

// TestAccumulator_Add tests the Add method of the Accumulator
func TestAccumulator_Add(t *testing.T) {
	acc := Accumulator{}

	// Test adding positive value
	acc.Add(5)
	if acc.GetValue() != 5 {
		t.Errorf("Expected value after adding 5 to be 5, got %d", acc.GetValue())
	}

	// Test adding negative value
	acc.Add(-3)
	if acc.GetValue() != 2 {
		t.Errorf("Expected value after adding -3 to be 2, got %d", acc.GetValue())
	}

	// Test adding zero
	acc.Add(0)
	if acc.GetValue() != 2 {
		t.Errorf("Expected value after adding 0 to remain 2, got %d", acc.GetValue())
	}
}

// TestAccumulator_GetValue tests the GetValue method of the Accumulator
func TestAccumulator_GetValue(t *testing.T) {
	acc := Accumulator{}

	// Initial value should be 0
	if acc.GetValue() != 0 {
		t.Errorf("Expected initial value to be 0, got %d", acc.GetValue())
	}

	// Change value and test again
	acc.Add(10)
	if acc.GetValue() != 10 {
		t.Errorf("Expected value after adding 10 to be 10, got %d", acc.GetValue())
	}
}
