package main

import (
	"strconv"
	"strings"

	"github.com/dotMaro/AoC2019/utils"
)

func main() {
	input := utils.ReadFile("day1/input.txt")
	utils.Print("Fuel required (ignoring the fuel's mass): %s", calculateFuelForModules(input, false))
	utils.Print("Fuel required (considering the fuel's mass): %s", calculateFuelForModules(input, true))
}

func calculateFuelForModules(input string, considerFuelMass bool) int {
	var fuelRequired int
	for _, massString := range strings.Split(input, "\n") {
		mass, err := strconv.Atoi(massString)
		if err != nil {
			panic(err)
		}
		fuelRequired += fuel(mass, considerFuelMass)
	}
	return fuelRequired
}

func fuel(mass int, considerFuelMass bool) int {
	if mass <= 0 {
		return 0
	}
	fuelReq := mass/3 - 2 // int division so it will be floored anyway
	if considerFuelMass {
		fuelForFuel := fuel(fuelReq, true)
		if fuelForFuel > 0 {
			fuelReq += fuelForFuel
		}
	}
	return fuelReq
}
