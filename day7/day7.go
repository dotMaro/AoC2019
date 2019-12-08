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
)

func main() {
	input := utils.ReadFile("day7/input.txt")

	utils.Print("Part 1: Output %d", findHighestOutput(input, 0, 4))
	output := findHighestOutput(input, 5, 9)
	utils.Print("Part 2: Output %v", output)
}

func findHighestOutput(input string, min, max int) int {
	var (
		highest int
		output  = 0
	)

	phaseSettings := make([]int, 0, max-min+1)
	for i := min; i <= max; i++ {
		phaseSettings = append(phaseSettings, i)
	}

	for p := make([]int, len(phaseSettings)); p[0] < len(p); nextPerm(p) {
		perm := getPerm(phaseSettings, p)
		p1, p2, p3, p4, p5 := perm[0], perm[1], perm[2], perm[3], perm[4]
		output = testPhases(input, p1, p2, p3, p4, p5)
		if output > highest {
			highest = output
		}
	}

	return highest
}

func nextPerm(p []int) {
	for i := len(p) - 1; i >= 0; i-- {
		if i == 0 || p[i] < len(p)-i-1 {
			p[i]++
			return
		}
		p[i] = 0
	}
}

func getPerm(orig, p []int) []int {
	result := append([]int{}, orig...)
	for i, v := range p {
		result[i], result[i+v] = result[i+v], result[i]
	}
	return result
}

func testPhases(input string, p1, p2, p3, p4, p5 int) int {
	signalInA := make(chan int, 100)
	signalAToB, _ := newProgram(input, signalInA).run()
	signalBToC, _ := newProgram(input, signalAToB).run()
	signalCToD, _ := newProgram(input, signalBToC).run()
	signalDToA, _ := newProgram(input, signalCToD).run()
	ampE := newProgram(input, signalDToA)
	signalOutE, stoppedSigE := ampE.run()

	// first signal is their phase settings
	signalInA <- p1
	signalAToB <- p2
	signalBToC <- p3
	signalCToD <- p4
	signalDToA <- p5

	// start off by sending signal 0 to amp A
	signalInA <- 0

	for {
		select {
		case s := <-signalOutE: // feedback loop
			signalInA <- s
		case <-stoppedSigE:
			return ampE.lastOutput
		}
	}
}

type executor interface {
	// store value at address that the param paramIndex points at.
	// Returns whether overwriting the current address's instruction.
	store(paramIndex int, value int) bool
	// stop execution.
	stop()
	// input returns input into the program.
	input() int
	// output a value.
	output(int)
	// setPC sets the program counter to a value.
	setPC(int)
}

type program struct {
	memory              []int
	initial             []int // what memory resets to when calling reset()
	signalIn, signalOut chan int
	lastOutput          int
	stopExecution       bool
	pc                  int // program counter
	operands            map[int]operand
}

func newProgram(input string, signalIn chan int) *program {
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
	return &program{
		memory:    data,
		initial:   dataCopy,
		signalIn:  signalIn,
		signalOut: make(chan int, 100),
		operands:  defaultOperands(),
	}
}

func defaultOperands() map[int]operand {
	return map[int]operand{
		1:  operand{3, add},
		2:  operand{3, multiply},
		3:  operand{1, input},
		4:  operand{1, output},
		5:  operand{3, jumpIfTrue},
		6:  operand{3, jumpIfFalse},
		7:  operand{3, lessThan},
		8:  operand{3, equals},
		99: operand{0, halt},
	}
}

// input returns input into the program.
func (p *program) input() int {
	return <-p.signalIn
}

func (p *program) output(o int) {
	p.signalOut <- o
	p.lastOutput = o
}

func (p *program) setPC(pc int) {
	p.pc = pc
}

func (p *program) stop() {
	p.stopExecution = true
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

// run program.
func (p *program) run() (signalOut chan int, stopped chan struct{}) {
	stopped = make(chan struct{})
	p.stopExecution = false
	go func(stopped chan struct{}) {
		for !p.stopExecution {
			p.step()
		}
		stopped <- struct{}{}
	}(stopped)
	return p.signalOut, stopped
}

func (p *program) step() {
	instruction := p.memory[p.pc]

	op, paramModes := instruction%100, instruction/100

	operand, foundOperand := p.operands[op]
	if !foundOperand {
		panic(fmt.Sprintf("unknown op code %d at pos %d", op, p.pc))
	}

	params := make([]int, operand.paramCount)
	for i := 0; i < operand.paramCount; i++ {
		paramMode := paramModes % int(math.Pow(10, float64(i+1))) / int(math.Pow(10, float64(i)))
		params[i] = p.param(i+1, paramMode)
	}

	p.pc += operand.opFunc(p, params)
}

type operand struct {
	paramCount int
	opFunc     operandFunc
}

type operandFunc func(s executor, params []int) (pcInc int)

func add(e executor, params []int) (pcInc int) {
	if !e.store(3, params[0]+params[1]) {
		pcInc = 4
	}
	return
}

func multiply(e executor, params []int) (pcInc int) {
	if !e.store(3, params[0]*params[1]) {
		pcInc = 4
	}
	return
}

func input(e executor, params []int) (pcInc int) {
	if !e.store(1, e.input()) {
		pcInc = 2
	}
	return
}

func output(e executor, params []int) (pcInc int) {
	e.output(params[0])
	return 2
}

func jumpIfTrue(e executor, params []int) (pcInc int) {
	if params[0] != 0 {
		e.setPC(params[1])
	} else {
		pcInc = 3
	}
	return
}

func jumpIfFalse(e executor, params []int) (pcInc int) {
	if params[0] == 0 {
		e.setPC(params[1])
	} else {
		pcInc = 3
	}
	return
}

func lessThan(e executor, params []int) (pcInc int) {
	value := 0
	if params[0] < params[1] {
		value = 1
	}
	e.store(params[2], value)
	return 4
}

func equals(e executor, params []int) (pcInc int) {
	value := 0
	if params[0] == params[1] {
		value = 1
	}
	e.store(params[2], value)
	return 4
}

func halt(e executor, params []int) (pcInc int) {
	e.stop()
	return
}
