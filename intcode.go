package aoc2019shared

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
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
	eof
)

func isValid(inst int) bool {
	return inst >= int(add) && inst < int(eof)
}

// IntCodeInterpreter is a interpreter for the int code language defined in
// Advent Of Code 2019 day 2 and 5 (and more!)
type IntCodeInterpreter struct {
	inst []int
}

// NewIntCodeInterpreter creates an int code interpreter with the
// given instructions.
func NewIntCodeInterpreter(input string) *IntCodeInterpreter {
	interpreter := IntCodeInterpreter{
		inst: parseInstructions(input),
	}

	return &interpreter
}

// Process runs the program in the IntCodeInterpreter's instructions. It returns
// the value in the 0 instruction at the end. Predefined inputs can be provided
// into the reader for testing, or provide os.Stdin and os.Stdout.
func (ici IntCodeInterpreter) Process(r io.Reader, w io.Writer) int {
	ip := 0
	var isParam1Immediate, isParam2Immediate bool
	reader := bufio.NewReader(r)

	for {
		oper := ici.inst[ip] % 10

		if !isValid(oper) {
			return ici.inst[0]
		}

		isParam1Immediate = (ici.inst[ip]/100)%10 == 1

		var p1, p2, p3 *int
		if isParam1Immediate {
			p1 = &ici.inst[ip+1]
		} else {
			p1 = &ici.inst[ici.inst[ip+1]]
		}

		if oper == 3 || oper == 4 {
			isParam2Immediate = false
		} else {
			isParam2Immediate = (ici.inst[ip]/1000)%10 == 1

			if isParam2Immediate {
				p2 = &ici.inst[ip+2]
			} else {
				p2 = &ici.inst[ici.inst[ip+2]]
			}

			if oper != 5 && oper != 6 {
				p3 = &ici.inst[ici.inst[ip+3]]
			}
		}

		switch operation(oper) {
		case add:
			*p3 = *p1 + *p2
			ip += 4
			break

		case multiply:
			*p3 = *p1 * *p2
			ip += 4
			break

		case input:
			fmt.Println("Input integer: ")
			text, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}

			text = text[0:strings.Index(text, "\n")]
			val, err := strconv.Atoi(text)
			if err != nil {
				log.Fatal(err)
			}

			// Parameters that an instruction writes to will never be immediate
			p1 = &ici.inst[ici.inst[ip+1]]

			*p1 = val
			ip += 2

		case output:
			fmt.Fprintf(w, "%d\n", *p1)
			ip += 2
			break

		case jumpTrue:
			if *p1 != 0 {
				ip = *p2
			} else {
				ip += 3
			}

		case jumpFalse:
			if *p1 == 0 {
				ip = *p2
			} else {
				ip += 3
			}
			break

		case lessThan:
			if *p1 < *p2 {
				*p3 = 1
			} else {
				*p3 = 0
			}
			ip += 4
			break

		case equals:
			if *p1 == *p2 {
				*p3 = 1
			} else {
				*p3 = 0
			}
			ip += 4
			break
		}
	}
}

func parseInstructions(input string) []int {
	output := []int{}

	split := strings.Split(input, ",")

	for _, strVal := range split {
		intVal, err := strconv.Atoi(strVal)

		if err != nil {
			log.Fatal(err)
		}

		output = append(output, intVal)
	}

	return output
}
