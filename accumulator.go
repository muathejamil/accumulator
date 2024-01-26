package main

// Accumulator struct
type Accumulator struct {
	value int
}

// Add method
func (a *Accumulator) Add(val int) {
	a.value += val
}

// GetValue Get value
func (a *Accumulator) GetValue() int {
	return a.value
}
