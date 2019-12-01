package main

import (
	"io/ioutil"
	"math"
	"strconv"
	"strings"

	"github.com/dotMaro/AoC2019/utils"
)

func main() {
	inputFile := utils.OpenFile("day1/input.txt")
	defer inputFile.Close()
	inputBytes, err := ioutil.ReadAll(inputFile)
	if err != nil {
		panic(err)
	}
	input := string(inputBytes)

	utils.Print("Fuel required (ignoring the fuel's mass): %v", calculateFuelRequired(input, false))
	utils.Print("Fuel required (considering the fuel's mass): %v", calculateFuelRequired(input, true))
}

func calculateFuelRequired(input string, considerFuelMass bool) int {
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
	fuelReq := int(math.Floor(float64(mass)/3.0)) - 2
	if considerFuelMass {
		fuelForFuel := fuel(fuelReq, true)
		if fuelForFuel > 0 {
			fuelReq += fuelForFuel
		}
	}
	return fuelReq
}
