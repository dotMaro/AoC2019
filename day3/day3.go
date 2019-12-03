package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dotMaro/AoC2019/utils"
)

func main() {
	input := utils.ReadFile("day3/input.txt")
	wireInput := strings.Split(input, "\n")
	w1 := newWire(wireInput[0])
	w2 := newWire(wireInput[1])
	var (
		closestInterceptDistance int
		closestInterceptWireLen  int
	)
	for _, p := range w1.interceptPoints(w2) {
		dist := p.distanceToOrigin()
		if closestInterceptDistance == 0 || dist < closestInterceptDistance {
			closestInterceptDistance = dist
		}
		if closestInterceptWireLen == 0 || p.wireLen < closestInterceptWireLen {
			closestInterceptWireLen = p.wireLen
		}
	}
	utils.Print("Part 1: Closest intercept point's Manhattan distance is %d", closestInterceptDistance)
	utils.Print("Part 2: Closest intercept point's wire length is %d", closestInterceptWireLen)
}

type direction int

const (
	up direction = iota
	down
	left
	right
)

// getDirection from byte.
// Panics on unknown direction.
func getDirection(d byte) direction {
	switch d {
	case 'R':
		return right
	case 'L':
		return left
	case 'U':
		return up
	case 'D':
		return down
	default:
		panic(fmt.Sprintf("unknown direction %v", d))
	}
}

type wire struct {
	// The points where the wire turns,
	// and where its segments split.
	points []point
}

func newWire(input string) wire {
	splitInput := strings.Split(input, ",")
	wire := wire{
		points: make([]point, 0, len(splitInput)),
	}
	for _, t := range splitInput {
		direction := getDirection(t[0])
		distance, err := strconv.Atoi(t[1:])
		if err != nil {
			panic(err)
		}
		wire.addSegment(direction, distance)
	}
	return wire
}

// addSegment to the wire.
func (w *wire) addSegment(dir direction, dist int) {
	var lastPoint point
	if len(w.points) != 0 {
		lastPoint = w.points[len(w.points)-1]
	}
	w.points = append(w.points, lastPoint.move(dir, dist))
}

// interceptPoints returns every point where the wire collides with wire o.
// The points' wireLen is the total wire length to get to that point (both wire combined).
func (w *wire) interceptPoints(o wire) []point {
	var interceptPoints []point
	for i := 1; i < len(w.points); i++ {
		v1 := segment{
			from: w.points[i-1],
			to:   w.points[i],
		}
		for u := 1; u < len(o.points); u++ {
			v2 := segment{
				from: o.points[u-1],
				to:   o.points[u],
			}
			intercept := v1.intercepts(v2)
			if intercept.x != 0 && intercept.y != 0 {
				// Calculate total wire length (both wires combined)
				intercept.wireLen = v1.from.wireLen + intercept.distanceToPoint(v1.from) +
					v2.from.wireLen + intercept.distanceToPoint(v2.from)
				interceptPoints = append(interceptPoints, intercept)
			}
		}
	}
	return interceptPoints
}

// point represents an edge of a wire segment.
// It can also be viewed as a corner of the wire.
type point struct {
	x, y    int // coordinates
	wireLen int // wire length
}

// move returns a new point that has the same properties as
// the point, but has moved a certain distance dist in direction dir.
func (p point) move(dir direction, dist int) point {
	var movedPoint point
	switch dir {
	case up:
		movedPoint = point{x: p.x, y: p.y + dist}
	case down:
		movedPoint = point{x: p.x, y: p.y - dist}
	case right:
		movedPoint = point{x: p.x + dist, y: p.y}
	case left:
		movedPoint = point{x: p.x - dist, y: p.y}
	}
	movedPoint.wireLen = p.wireLen + dist
	return movedPoint
}

// distanceToOrigin returns the Manhattan distance to origin (0, 0).
func (p point) distanceToOrigin() int {
	return p.distanceToPoint(point{x: 0, y: 0})
}

// distanceToPoint returns the Manhattan distance to point o.
func (p point) distanceToPoint(o point) int {
	return abs(abs(p.x)-abs(o.x)) + abs(abs(p.y)-abs(o.y))
}

// abs returns the absolute value of i.
func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

// segment represents a wire segment (which is always straight).
// It is a closed line segment between two points.
type segment struct {
	from, to point
}

// unchangingAxis returns what axis and the value of that axis.
// It assumes that exactly one axis is changing.
func (v segment) unchangingAxis() (val int, xAxis bool) {
	if v.from.x == v.to.x {
		return v.from.x, true
	}
	return v.from.y, false
}

// intercepts returns where the segment intercepts segment o.
// If there is no interception then (0, 0) will be returned. wirelen is not provided.
func (v segment) intercepts(o segment) point {
	// With the assumption that no interceptions occur when segments are
	// parallel, and that segments always move either horizontally or
	// vertically (not both), we can pretty easily check for interceptions.
	//
	// First find the values where interception could occur, and what axis for
	// both segments are changing. I.e. if the segments are horizontal
	// or vertical.
	a, axAxis := v.unchangingAxis()
	b, bxAxis := o.unchangingAxis()
	if axAxis == bxAxis {
		// We're assuming that they can't overlap
		// when they are parallel
		return point{}
	}

	// Check if the first value (x or y) is on the interval of the
	// same axis of the other segment. Do this for the other value (axis) too.
	var aCanCollide bool
	if axAxis {
		aCanCollide = inRange(a, o.from.x, o.to.x)
	} else {
		aCanCollide = inRange(a, o.from.y, o.to.y)
	}
	var bCanCollide bool
	if bxAxis {
		bCanCollide = inRange(b, v.from.x, v.to.x)
	} else {
		bCanCollide = inRange(b, v.from.y, v.to.y)
	}

	// If both axes are in range then they collide
	if aCanCollide && bCanCollide {
		// Check if a is an x- or y-value
		if axAxis {
			return point{x: a, y: b}
		}
		return point{x: b, y: a}
	}
	return point{x: 0, y: 0}
}

// inRange returns whether true if a >= val <= b
// or b >= val <= a.
func inRange(val, a, b int) bool {
	return val >= a && val <= b || val >= b && val <= a
}
