package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/dotMaro/AoC2019/utils"
)

const (
	desiredOutput = 19690720

	paramPositionMode  = 0
	paramImmediateMode = 1
	paramRelativeMode  = 2
)

func main() {
	input := utils.ReadFile("day9/input.txt")

	program := newProgram(input, 1)
	program.run()
	program.reset()
	program.inputInstruction = 2
	program.run()
}

type executor interface {
	// store value at address index.
	// Returns whether current instruction was overwritten.
	store(index int64, value int64) bool
	// stop execution.
	stop()
	// input returns input into the program.
	input() int
	// output a value.
	output(int64)
	// setPC sets the program counter to a value.
	setPC(int)
	adjustRelativeOffsetBase(offset int)
}

type program struct {
	memory           []int64
	initial          []int64 // what memory resets to when calling reset()
	inputInstruction int
	lastOutput       int64
	stopExecution    bool
	pc               int // program counter
	relativeBase     int
	operands         map[int]operand
}

func newProgram(input string, inputInstruction int) *program {
	split := strings.Split(input, ",")
	data := make([]int64, 40_000_000)
	for i, op := range split {
		var err error
		data[i], err = strconv.ParseInt(op, 10, 64)
		if err != nil {
			panic(err)
		}
	}
	dataCopy := make([]int64, len(data))
	copy(dataCopy, data)
	return &program{
		memory:           data,
		initial:          dataCopy,
		inputInstruction: inputInstruction,
		operands:         defaultOperands(),
	}
}

func defaultOperands() map[int]operand {
	return map[int]operand{
		1:  operand{2, 1, add},
		2:  operand{2, 1, multiply},
		3:  operand{0, 1, input},
		4:  operand{1, 0, output},
		5:  operand{2, 0, jumpIfTrue},
		6:  operand{2, 0, jumpIfFalse},
		7:  operand{2, 1, lessThan},
		8:  operand{2, 1, equals},
		9:  operand{1, 0, offsetRelativeBase},
		99: operand{0, 0, halt},
	}
}

// input returns input into the program.
func (p *program) input() int {
	return p.inputInstruction
}

func (p *program) output(o int64) {
	utils.Print("Output %d", o)
	p.lastOutput = o
}

func (p *program) setPC(pc int) {
	p.pc = pc
}

func (p *program) adjustRelativeOffsetBase(offset int) {
	p.relativeBase += offset
}

func (p *program) stop() {
	p.stopExecution = true
}

func (p *program) reset() {
	p.pc = 0
	p.relativeBase = 0
	copy(p.memory, p.initial)
}

func (p *program) inputParam(paramIndex int, mode int) int64 {
	var index int
	switch mode {
	case paramPositionMode:
		index = int(p.memory[p.pc+paramIndex])
	case paramImmediateMode:
		index = p.pc + paramIndex
	case paramRelativeMode:
		index = p.relativeBase + int(p.memory[p.pc+paramIndex])
	default:
		panic(fmt.Sprintf("unknown input param mode %d", mode))
	}

	return p.memory[index]
}

func (p *program) outputParam(paramIndex int, mode int) int64 {
	var index int64
	switch mode {
	case paramPositionMode:
		index = p.memory[p.pc+paramIndex]
	case paramRelativeMode:
		index = int64(p.relativeBase) + p.memory[p.pc+paramIndex]
	default:
		panic(fmt.Sprintf("unknown output param mode %d", mode))
	}

	return index
}

func (p *program) store(index int64, value int64) bool {
	p.memory[index] = value
	return index == int64(p.pc)
}

// run program.
func (p *program) run() {
	p.stopExecution = false
	for !p.stopExecution {
		p.step()
	}
}

func (p *program) step() {
	instruction := p.memory[p.pc]

	op, paramModes := instruction%100, instruction/100

	operand, foundOperand := p.operands[int(op)]
	if !foundOperand {
		panic(fmt.Sprintf("unknown op code %d at pos %d", op, p.pc))
	}

	inputParams := make([]int64, operand.inputParamCount)
	for i := 0; i < operand.inputParamCount; i++ {
		paramMode := int(paramModes) % int(math.Pow(10, float64(i+1))) / int(math.Pow(10, float64(i)))
		inputParams[i] = p.inputParam(i+1, paramMode)
	}
	outputParams := make([]int64, operand.outputParamCount)
	totalParams := operand.inputParamCount + operand.outputParamCount
	for i := operand.inputParamCount; i < totalParams; i++ {
		paramMode := int(paramModes) % int(math.Pow(10, float64(i+1))) / int(math.Pow(10, float64(i)))
		outputParams[i-operand.inputParamCount] = p.outputParam(i+1, paramMode)
	}

	p.pc += operand.opFunc(p, inputParams, outputParams)
}

type operand struct {
	inputParamCount  int
	outputParamCount int
	opFunc           operandFunc
}

type operandFunc func(s executor, inputParams []int64, outputParams []int64) (pcInc int)

func add(e executor, inputParams []int64, outputParams []int64) (pcInc int) {
	if !e.store(outputParams[0], inputParams[0]+inputParams[1]) {
		pcInc = 4
	}
	return
}

func multiply(e executor, inputParams []int64, outputParams []int64) (pcInc int) {
	if !e.store(outputParams[0], inputParams[0]*inputParams[1]) {
		pcInc = 4
	}
	return
}

func input(e executor, parinputParams []int64, outputParams []int64) (pcInc int) {
	if !e.store(outputParams[0], int64(e.input())) {
		pcInc = 2
	}
	return
}

func output(e executor, inputParams []int64, outputParams []int64) (pcInc int) {
	e.output(inputParams[0])
	return 2
}

func jumpIfTrue(e executor, inputParams []int64, outputParams []int64) (pcInc int) {
	if inputParams[0] != 0 {
		e.setPC(int(inputParams[1]))
	} else {
		pcInc = 3
	}
	return
}

func jumpIfFalse(e executor, inputParams []int64, outputParams []int64) (pcInc int) {
	if inputParams[0] == 0 {
		e.setPC(int(inputParams[1]))
	} else {
		pcInc = 3
	}
	return
}

func lessThan(e executor, inputParams []int64, outputParams []int64) (pcInc int) {
	var value int64
	if inputParams[0] < inputParams[1] {
		value = 1
	}
	e.store(outputParams[0], value)
	return 4
}

func equals(e executor, inputParams []int64, outputParams []int64) (pcInc int) {
	var value int64
	if inputParams[0] == inputParams[1] {
		value = 1
	}
	e.store(outputParams[0], value)
	return 4
}

func offsetRelativeBase(e executor, inputParams []int64, outputParams []int64) (pcInc int) {
	e.adjustRelativeOffsetBase(int(inputParams[0]))
	return 2
}

func halt(e executor, inputParams []int64, outputParams []int64) (pcInc int) {
	e.stop()
	return
}
