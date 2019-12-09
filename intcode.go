package aoc2019shared

import (
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
)

type operation int

const (
	add operation = iota + 1
	multiply
	input
	output
	jumpTrue
	jumpFalse
	lessThan
	equals
	relativeBaseOffset
	eof
)

func isValid(inst int64) bool {
	return inst >= int64(add) && inst < int64(eof)
}

// IntCodeInterpreter is a interpreter for the int code language defined in
// Advent Of Code 2019 day 2 and 5 (and more!)
type IntCodeInterpreter struct {
	name         string
	inst         []int64
	ip           int64
	RelativeBase int64
	Input        chan int64
	Output       chan int64
	LastOutput   *int64
}

// NewIntCodeInterpreter creates an int code interpreter with the
// given instructions.
func NewIntCodeInterpreter(name, input string) *IntCodeInterpreter {
	interpreter := IntCodeInterpreter{
		name:         name,
		inst:         parseInstructions(input),
		ip:           0,
		RelativeBase: 0,
		Input:        make(chan int64, 2),
		Output:       make(chan int64, 2),
	}

	return &interpreter
}

// Process runs the program in the IntCodeInterpreter's instructions. It returns
// the value in the 0 instruction at the end.
func (ici *IntCodeInterpreter) Process(wg *sync.WaitGroup) int64 {
	var param1Mode, param2Mode int64

	for {
		oper := ici.inst[ici.ip] % 100

		if !isValid(oper) {
			if wg != nil {
				wg.Done()
			}
			return ici.inst[0]
		}

		param1Mode = (ici.inst[ici.ip] / 100) % 10

		var p1, p2, p3 *int64
		if param1Mode == 1 {
			p1 = &ici.inst[ici.ip+1]
		} else {
			p1 = &ici.inst[ici.inst[ici.ip+1]]
		}

		if oper == 3 || oper == 4 || oper == 9 {
			param2Mode = 0
		} else {
			param2Mode = (ici.inst[ici.ip] / 1000) % 10

			if param2Mode == 1 {
				p2 = &ici.inst[ici.ip+2]
			} else {
				p2 = &ici.inst[ici.inst[ici.ip+2]]
			}

			if oper != 5 && oper != 6 {
				p3 = &ici.inst[ici.inst[ici.ip+3]]
			}
		}

		switch operation(oper) {
		case add:
			*p3 = *p1 + *p2
			ici.ip += 4
			break

		case multiply:
			*p3 = *p1 * *p2
			ici.ip += 4
			break

		case input:
			val := <-ici.Input

			// Parameters that an instruction writes to will never be immediate
			p1 = &ici.inst[ici.inst[ici.ip+1]]

			*p1 = val
			ici.ip += 2

		case output:
			ici.LastOutput = p1
			ici.Output <- *p1
			ici.ip += 2
			break

		case jumpTrue:
			if *p1 != 0 {
				ici.ip = *p2
			} else {
				ici.ip += 3
			}

		case jumpFalse:
			if *p1 == 0 {
				ici.ip = *p2
			} else {
				ici.ip += 3
			}
			break

		case lessThan:
			if *p1 < *p2 {
				*p3 = 1
			} else {
				*p3 = 0
			}
			ici.ip += 4
			break

		case equals:
			if *p1 == *p2 {
				*p3 = 1
			} else {
				*p3 = 0
			}
			ici.ip += 4
			break

		case relativeBaseOffset:
			ici.RelativeBase += *p1
			break
		}
	}
}

func parseInstructions(input string) []int64 {
	output := make([]int64, math.MaxUint16)

	split := strings.Split(input, ",")

	for i, strVal := range split {
		intVal, err := strconv.ParseInt(strVal, 10, 0)

		if err != nil {
			log.Fatal(err)
		}

		output[i] = intVal
	}

	return output
}
