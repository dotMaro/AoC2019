package main

import (
	"fmt"
	"github.com/dotMaro/AoC2019/utils"
	"strconv"
	"strings"
)

const desiredOutput = 19690720

func main() {
	input := utils.ReadFile("day2/input.txt")
	program := newProgram(input)
	// "before running the program, replace position 1 with the value 12 and replace position 2 with the value 2"
	program.setInputs(12, 2)
	output := program.run()
	utils.Print("Step 1 (input 12 and 2): %v", output)
	program.reset()

	noun, verb := program.findInputForOutput(desiredOutput)
	utils.Print("Step 2: Output is %d with noun %d verb %d (100 * noun + verb = %d)",
		desiredOutput, noun, verb, 100*noun+verb)
}

type program struct {
	memory  []int
	initial []int // what memory resets to when calling reset()
	pc      int   // program counter
	stop    bool  // running state
}

func newProgram(input string) program {
	split := strings.Split(input, ",")
	data := make([]int, len(split))
	for i, op := range split {
		var err error
		data[i], err = strconv.Atoi(op)
		if err != nil {
			panic(err)
		}
	}
	dataCopy := make([]int, len(data))
	copy(dataCopy, data)
	return program{memory: data, initial: dataCopy}
}

// findInputForOutput will try to find the input noun and verb
// to get the desired output in the interval 0-99.
// If none is found then both values will be returned as -1.
func (p *program) findInputForOutput(desiredOutput int) (noun, verb int) {
	for noun = 0; noun <= 99; noun++ {
		for verb = 0; verb <= 99; verb++ {
			p.setInputs(noun, verb)
			output := p.run()
			if output == desiredOutput {
				p.reset()
				return noun, verb
			}
			p.reset()
		}
	}
	return -1, -1
}

func (p *program) setInputs(noun, verb int) {
	p.memory[1] = noun
	p.memory[2] = verb
}

func (p *program) reset() {
	p.pc = 0
	copy(p.memory, p.initial)
}

// run program and return output.
func (p *program) run() int {
	p.stop = false
	for !p.stop {
		p.step()
	}
	return p.memory[0]
}

func (p *program) step() {
	op := p.memory[p.pc]
	switch op {
	case 1: // add
		p.memory[p.memory[p.pc+3]] = p.memory[p.memory[p.pc+1]] + p.memory[p.memory[p.pc+2]]
	case 2: // multiplication
		p.memory[p.memory[p.pc+3]] = p.memory[p.memory[p.pc+1]] * p.memory[p.memory[p.pc+2]]
	case 99: // halt program
		p.stop = true
		return
	default:
		panic(fmt.Sprintf("unknown op code %d at pos %d", op, p.pc))
	}

	p.pc += 4
}
