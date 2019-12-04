package main

import "github.com/dotMaro/AoC2019/utils"

const (
	lowerBound = 272091
	upperBound = 815432
)

func main() {
	var validCount, validStrictlyTwoCount int
	for i := lowerBound; i <= upperBound; i++ {
		if validPassword(i, false) {
			validCount++
		}
		if validPassword(i, true) {
			validStrictlyTwoCount++
		}
	}

	utils.Print("Part 1: There are %d valid passwords", validCount)
	utils.Print("Part 2: There are %d valid passwords", validStrictlyTwoCount)
}

// validPassword is true if i has two adjacent digits and no digits are
// decreasing. If strictlyTwo is true then adjacent chains of digits that
// are longer than two will not be seen as valid.
func validPassword(i int, strictlyTwo bool) bool {
	return hasDoubleDigits(i, strictlyTwo) && nonDecreasingDigits(i)
}

func hasDoubleDigits(i int, strictlyTwo bool) bool {
	var (
		prevDigit int
		// The digit that is adjacent.
		// If zero then no adjacent digit has been found.
		adjDigit int
		// If the amount of adjacent digits are more than two.
		overTwo bool
	)
	for i != 0 {
		rightMost := i % 10
		if rightMost == prevDigit {
			if !strictlyTwo {
				return true
			}
			if adjDigit == rightMost {
				// longer repetition than two
				overTwo = true
			} else if adjDigit == 0 || overTwo {
				// only overwrite if there has not been another successful double
				adjDigit = rightMost
				overTwo = false
			} else { // there has been a successful double
				return true
			}
		}
		prevDigit = rightMost
		i /= 10
	}

	return adjDigit != 0 && !overTwo
}

func nonDecreasingDigits(i int) bool {
	prevDigit := 10
	for i != 0 {
		rightMost := i % 10
		if rightMost > prevDigit {
			return false
		}
		prevDigit = rightMost
		i /= 10
	}

	return true
}
