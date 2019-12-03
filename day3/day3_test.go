package main

import "testing"

func TestVectorCollidesWith(t *testing.T) {
	v1 := vector{
		from: point{0, 0},
		to:   point{100, 0},
	}
	v2 := vector{
		from: point{50, 50},
		to:   point{50, -50},
	}
	got := v1.collidesWith(v2)
	if got.x != 50 || got.y != 0 {
		t.Errorf("%v and %v should intercept at (50, 0), but got %v", v1, v2, got)
	}
	v1 = vector{
		from: point{0, 0},
		to:   point{100, 0},
	}
	v2 = vector{
		from: point{-10, 50},
		to:   point{-20, 50},
	}
	got = v1.collidesWith(v2)
	if got.x != 0 || got.y != 0 {
		t.Errorf("%v and %v should not intercept, but got %v", v1, v2, got)
	}
	v1 = vector{
		from: point{0, 0},
		to:   point{0, 10},
	}
	v2 = vector{
		from: point{-3, 5},
		to:   point{3, 5},
	}
	got = v1.collidesWith(v2)
	if got.x != 0 || got.y != 5 {
		t.Errorf("%v and %v should intercept at (0, 5), but got %v", v1, v2, got)
	}
}
