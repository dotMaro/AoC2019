package main

import "testing"

func TestHasAdjacentDigits(t *testing.T) {
	testCases := []struct {
		input int
		exp   bool
	}{
		{123, false},
		{1123, true},
	}

	for _, c := range testCases {
		got := hasDoubleDigits(c.input, false)
		if got != c.exp {
			t.Errorf("Should return %t, but got %t, for input %d", c.exp, got, c.input)
		}
	}
}
