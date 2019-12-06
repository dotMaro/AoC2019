package main

import (
	"strings"

	"github.com/dotMaro/AoC2019/utils"
)

const (
	com   = "COM"
	you   = "YOU"
	santa = "SAN"
)

func main() {
	input := utils.ReadFile("day6/input.txt")

	orbitMap := parseInput(input)
	com := orbitMap[com]
	utils.Print("Part 1: %d", com.totalOrbits())
	you := orbitMap[you]
	utils.Print("Part 2: %d", you.findOrbitalTransfersTo(santa)-2) // -2 because YOU and SAN are treated specially
}

type orbitMap map[string]*spaceObject

func parseInput(input string) orbitMap {
	orbitMap := make(map[string]*spaceObject)
	for _, o := range strings.Split(input, "\n") {
		objects := strings.Split(o, ")")
		orbitedName, orbiterName := objects[0], objects[1]

		orbiter, orbiterExists := orbitMap[orbiterName]
		if !orbiterExists {
			orbiter = &spaceObject{
				name: orbiterName,
			}
			orbitMap[orbiterName] = orbiter
		}

		if orbited, alreadyExists := orbitMap[orbitedName]; alreadyExists {
			orbiter.orbiting = orbited
			orbited.addOrbiter(orbiter)
		} else {
			orbited = &spaceObject{
				name:     orbitedName,
				orbiters: []*spaceObject{orbiter},
			}
			orbiter.orbiting = orbited
			orbitMap[orbitedName] = orbited
		}
	}
	return orbitMap
}

type spaceObject struct {
	name     string
	orbiting *spaceObject
	orbiters []*spaceObject
}

func (o *spaceObject) addOrbiter(orbiter *spaceObject) {
	if o.orbiters == nil {
		o.orbiters = []*spaceObject{orbiter}
	} else {
		o.orbiters = append(o.orbiters, orbiter)
	}
}

func (o *spaceObject) findOrbitalTransfersTo(find string) int {
	return o.findOrbitalTransfersToRecursive(find, o.name, 0)
}

func (o *spaceObject) findOrbitalTransfersToRecursive(findName, justVisitedName string, depth int) int {
	if o.name == findName {
		return depth
	}
	shortest := 999999
	if o.orbiting != nil && o.orbiting.name != justVisitedName {
		shortest = o.orbiting.findOrbitalTransfersToRecursive(findName, o.name, depth+1)
	}

	for _, orbiter := range o.orbiters {
		if orbiter.name != justVisitedName {
			transfers := orbiter.findOrbitalTransfersToRecursive(findName, o.name, depth+1)
			if transfers < shortest {
				shortest = transfers
			}
		}
	}

	return shortest
}

func (o *spaceObject) totalOrbits() int {
	return o.totalOrbitsRecursive(0)
}

func (o *spaceObject) totalOrbitsRecursive(depth int) int {
	if o.orbiters == nil {
		return depth
	}

	var orbits int
	for _, orbiter := range o.orbiters {
		orbits += orbiter.totalOrbitsRecursive(depth + 1)
	}
	return orbits + depth
}
