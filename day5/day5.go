package main

import (
	"fmt"
	"github.com/dotMaro/AoC2019/utils"
	"strconv"
	"strings"
)

const (
	desiredOutput = 19690720

	paramPositionMode  = 0
	paramImmediateMode = 1
)

func main() {
	input := utils.ReadFile("day5/input.txt")
	program := newProgram(input, 2)
	// program.setInputs(12, 2)
	output := program.run()
	utils.Print("Step 1: %v", output)
	program.reset()

	program.setInputInstruction(5)
	output = program.run()
	utils.Print("Step 2: %v", output)
}

type storer interface {
	// store value at paramIndex
	store(paramIndex int, value int) bool
}

type program struct {
	memory  []int
	initial []int // what memory resets to when calling reset()
	input   int   // input instruction
	pc      int   // program counter
}

func newProgram(input string, inputInstruction int) program {
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
	return program{
		memory:  data,
		initial: dataCopy,
		input:   inputInstruction,
	}
}

func (p *program) setInputInstruction(i int) {
	p.input = i
}

func (p *program) reset() {
	p.pc = 0
	copy(p.memory, p.initial)
}

func (p *program) param(paramIndex int, mode int) int {
	if mode == paramImmediateMode {
		return p.memory[p.pc+paramIndex]
	}

	// position mode
	return p.memory[p.memory[p.pc+paramIndex]]
}

func (p *program) store(paramIndex int, value int) bool {
	index := p.memory[p.pc+paramIndex]
	p.memory[index] = value
	return index == p.pc
}

// run program and return output.
func (p *program) run() int {
	var res int
	for res == 0 {
		res = p.step()
	}
	return res
}

func (p *program) step() int {
	instruction := p.memory[p.pc]

	op, paramModes := instruction%100, instruction/100
	param1Mode, param2Mode := paramModes%10, paramModes%100/10
	// param3Mode := paramModes/100 // never used as so far the third param is always an address to store something in

	switch op {
	case 1: // add
		if !p.store(3, p.param(1, param1Mode)+p.param(2, param2Mode)) {
			p.pc += 4
		}
	case 2: // multiplication
		if !p.store(3, p.param(1, param1Mode)*p.param(2, param2Mode)) {
			p.pc += 4
		}
	case 3: // input
		if !p.store(1, p.input) {
			p.pc += 2
		}
	case 4: // output
		output := p.param(1, param1Mode)
		pcAfter := p.pc + 2
		if pcAfter < len(p.memory) && p.memory[pcAfter] == 99 {
			return output
		}
		if output != 0 {
			utils.Print("pc %d: output should have been 0 but was %d", p.pc, output)
			return -1
		}
		p.pc = pcAfter
	case 5: // jump if true
		if p.param(1, param1Mode) != 0 {
			p.pc = p.param(2, param2Mode)
		} else {
			p.pc += 3
		}
	case 6: // jump if false
		if p.param(1, param1Mode) == 0 {
			p.pc = p.param(2, param2Mode)
		} else {
			p.pc += 3
		}
	case 7: // less than
		isLess := p.param(1, param1Mode) < p.param(2, param2Mode)
		var store int
		if isLess {
			store = 1
		}
		if !p.store(3, store) {
			p.pc += 4
		}
	case 8: // equals
		equals := p.param(1, param1Mode) == p.param(2, param2Mode)
		var store int
		if equals {
			store = 1
		}
		if !p.store(3, store) {
			p.pc += 4
		}
	case 99: // halt program
		return -1
	default:
		panic(fmt.Sprintf("unknown op code %d at pos %d", op, p.pc))
	}

	return 0
}

type operand func(s storer) (pcInc int, output int)

func add(s storer) (pcInc int, output int) {
	// p.memory[p.memory[p.pc+3]] = p.memory[p.memory[p.pc+1]] + p.memory[p.memory[p.pc+2]]
	return 4, 0
}

func multiply(p *program) (pcInc int, output int) {
	p.memory[p.memory[p.pc+3]] = p.memory[p.memory[p.pc+1]] * p.memory[p.memory[p.pc+2]]
	return 4, 0
}

func halt(p *program) (pcInc int, output int) {
	return 0, -1
}
