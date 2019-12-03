package main

import "testing"

func TestVectorCollidesWith(t *testing.T) {
	v1 := segment{
		from: point{x: 0, y: 0},
		to:   point{x: 100, y: 0},
	}
	v2 := segment{
		from: point{x: 50, y: 50},
		to:   point{x: 50, y: -50},
	}
	got := v1.collidesWith(v2)
	if got.x != 50 || got.y != 0 {
		t.Errorf("%v and %v should intercept at (50, 0), but got %v", v1, v2, got)
	}
	v1 = segment{
		from: point{x: 0, y: 0},
		to:   point{x: 100, y: 0},
	}
	v2 = segment{
		from: point{x: -10, y: 50},
		to:   point{x: -20, y: 50},
	}
	got = v1.collidesWith(v2)
	if got.x != 0 || got.y != 0 {
		t.Errorf("%v and %v should not intercept, but got %v", v1, v2, got)
	}
	v1 = segment{
		from: point{x: 0, y: 0},
		to:   point{x: 0, y: 10},
	}
	v2 = segment{
		from: point{x: -3, y: 5},
		to:   point{x: 3, y: 5},
	}
	got = v1.collidesWith(v2)
	if got.x != 0 || got.y != 5 {
		t.Errorf("%v and %v should intercept at (0, 5), but got %v", v1, v2, got)
	}
}
