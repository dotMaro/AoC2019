package main

import (
	"io/ioutil"

	"github.com/dotMaro/AoC2019/utils"
)

func main() {
	input := utils.ReadFile("day1/input.txt")
	defer input.Close()
	content, err := ioutil.ReadAll(input)
	if err != nil {
		panic(err)
	}
	utils.Print("Input: %v", string(content))
}
