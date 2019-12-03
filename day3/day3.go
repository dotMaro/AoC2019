package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dotMaro/AoC2019/utils"
)

/*
	There are a lot of improvements that can be made here, but I feel
	I've spent enough time on this day's tasks so I'm just going to leave it
	as it is. At least it works!
*/

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
		dist := p.distance()
		if closestInterceptDistance == 0 || dist < closestInterceptDistance {
			closestInterceptDistance = dist
		}
		if closestInterceptWireLen == 0 || p.wireLen < closestInterceptWireLen {
			closestInterceptWireLen = p.wireLen
		}
	}
	utils.Print("Part 1: Closest intercept's distance is %d", closestInterceptDistance)
	utils.Print("Part 2: Closest intercept's wire length is %d", closestInterceptWireLen)
}

type direction int

const (
	up direction = iota
	down
	left
	right
)

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
		wire.addTraversal(direction, distance)
	}
	return wire
}

func (w *wire) addTraversal(dir direction, dist int) {
	var lastPoint, newPoint point
	if len(w.points) != 0 {
		lastPoint = w.points[len(w.points)-1]
	}

	switch dir {
	case up:
		newPoint = point{x: lastPoint.x, y: lastPoint.y + dist}
	case down:
		newPoint = point{x: lastPoint.x, y: lastPoint.y - dist}
	case right:
		newPoint = point{x: lastPoint.x + dist, y: lastPoint.y}
	case left:
		newPoint = point{x: lastPoint.x - dist, y: lastPoint.y}
	}
	newPoint.wireLen = lastPoint.wireLen + dist

	w.points = append(w.points, newPoint)
}

func (w *wire) interceptPoints(o wire) []point {
	var interceptPoints []point
	for i := 1; i < len(w.points); i++ {
		v1 := vector{
			from: w.points[i-1],
			to:   w.points[i],
		}
		for u := 1; u < len(o.points); u++ {
			v2 := vector{
				from: o.points[u-1],
				to:   o.points[u],
			}
			intercept := v1.collidesWith(v2)
			if intercept.x != 0 && intercept.y != 0 {
				intercept.wireLen = v1.from.wireLen + intercept.distanceToPoint(v1.from) +
					v2.from.wireLen + intercept.distanceToPoint(v2.from)
				interceptPoints = append(interceptPoints, intercept)
			}
		}
	}
	return interceptPoints
}

type point struct {
	x, y    int
	wireLen int
}

func (p point) distance() int {
	return p.distanceToPoint(point{x: 0, y: 0})
}

func (p point) distanceToPoint(o point) int {
	return abs(abs(p.x)-abs(o.x)) + abs(abs(p.y)-abs(o.y))
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

type vector struct {
	from, to point
}

func (v vector) unchangingAxis() (val int, xAxis bool) {
	if v.from.x == v.to.x {
		return v.from.x, true
	}
	// assuming y is unchanging
	return v.from.y, false
}

// collidesWith returns where the vector collides with vector o.
// If there is no collision 0,0 will be returned. wirelen is not provided.
func (v vector) collidesWith(o vector) point {
	/*
		a = Find unchanging value in v1
		b = Find unchanging value in v2

		is a within v2's values in same axis (x or y)?
		is b within v1's values in same axis (x or y)?

		if both yes:
			intersect at (a, b)
	*/
	a, axAxis := v.unchangingAxis()
	b, bxAxis := o.unchangingAxis()
	// axAxis should never be equal to bxAxis

	var aCanCollide bool
	if axAxis {
		aCanCollide = a >= o.from.x && a <= o.to.x ||
			a <= o.from.x && a >= o.to.x
	} else {
		aCanCollide = a >= o.from.y && a <= o.to.y ||
			a <= o.from.y && a >= o.to.y
	}

	var bCanCollide bool
	if bxAxis {
		bCanCollide = b >= v.from.x && b <= v.to.x ||
			b <= v.from.x && b >= v.to.x
	} else {
		bCanCollide = b >= v.from.y && b <= v.to.y ||
			b <= v.from.y && b >= v.to.y
	}

	if aCanCollide && bCanCollide {
		if axAxis {
			return point{x: a, y: b}
		}
		return point{x: b, y: a}
	}
	return point{x: 0, y: 0}
}
