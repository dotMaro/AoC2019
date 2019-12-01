package main

import (
	"testing"
)

func TestFuel(t *testing.T) {
	testCases := []struct {
		mass          int // input mass
		expNoFuelMass int // expected when not considering fuel mass
		expFuelMass   int // expected when considering fuel mass
	}{
		{12, 2, 2},
		{14, 2, 2},
		{1969, 654, 966},
		{100756, 33583, 50346},
	}

	for _, c := range testCases {
		got := fuel(c.mass, false)
		if got != c.expNoFuelMass {
			t.Errorf("Should return %d, not %d, for input %d when not considering fuel mass",
				c.expNoFuelMass, got, c.mass)
		}
		got = fuel(c.mass, true)
		if got != c.expFuelMass {
			t.Errorf("Should return %d, not %d, for input %d when considering fuel mass",
				c.expFuelMass, got, c.mass)
		}
	}
}

func TestCalculateFuelForModules(t *testing.T) {
	const input = "12\n14\n1969\n100756"
	got := calculateFuelForModules(input, false)
	exp := 2 + 2 + 654 + 33583
	if got != exp {
		t.Errorf("Should return %d, not %d, when not considering fuel mass", exp, got)
	}
	got = calculateFuelForModules(input, true)
	exp = 2 + 2 + 966 + 50346
	if got != exp {
		t.Errorf("Should return %d, not %d, when considering fuel mass", exp, got)
	}
}
