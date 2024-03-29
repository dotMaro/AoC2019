package main

import (
	"reflect"
	"testing"
)

func TestRunProgram(t *testing.T) {
	testCases := []struct {
		input string
		exp   []int
	}{
		{"1,0,0,0,99", []int{2, 0, 0, 0, 99}},
		{"2,3,0,3,99", []int{2, 3, 0, 6, 99}},
		{"2,4,4,5,99,0", []int{2, 4, 4, 5, 99, 9801}},
		{"1,1,1,4,99,5,6,0,99", []int{30, 1, 1, 4, 2, 5, 6, 0, 99}},
	}

	for _, c := range testCases {
		t.Logf("Input %v", c.input)
		p := newProgram(c.input)
		output := p.run()
		if !reflect.DeepEqual(p.memory, c.exp) {
			t.Errorf("Should have data\n%v, not\n%v on input %v", c.exp, p.memory, c.input)
		}
		if output != c.exp[0] {
			t.Errorf("Should return output %d, but returned %d", c.exp[0], output)
		}
	}
}

func TestResetProgram(t *testing.T) {
	p := program{
		memory:  []int{2, 5, 5, 5, 99},
		initial: []int{2, 0, 0, 0, 99},
	}
	p.reset()
	if !reflect.DeepEqual(p.memory, p.initial) {
		t.Errorf("Should reset to initial values %v, but was %v", p.initial, p.memory)
	}
	if p.pc != 0 {
		t.Errorf("PC should be 0, but was %d", p.pc)
	}
}

func TestProgramSetInputs(t *testing.T) {
	const (
		noun = 3
		verb = 5
	)
	p := newProgram("1,0,0,0,99")
	p.setInputs(noun, verb)
	if p.memory[1] != noun {
		t.Errorf("Should set noun to %d, but was %d", noun, p.memory[1])
	}
	if p.memory[2] != verb {
		t.Errorf("Should set verb to %d, but was %d", verb, p.memory[2])
	}
}
