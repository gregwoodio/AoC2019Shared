package aoc2019shared

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

type testData struct {
	input          string
	expectedZero   int
	expectedOutput int
	userInput      string
}

func TestProcessIntCode(t *testing.T) {
	testDatas := []testData{
		testData{
			input:          "3,1,4,1,99",
			expectedOutput: 42,
			expectedZero:   3,
			userInput:      "42\n",
		},
		// Day 5 tests (less than and equal operators)
		testData{
			input:          "3,9,8,9,10,9,4,9,99,-1,8",
			expectedOutput: 1,
			expectedZero:   3,
			userInput:      "8\n",
		},
		testData{
			input:          "3,9,8,9,10,9,4,9,99,-1,8",
			expectedOutput: 0,
			expectedZero:   3,
			userInput:      "7\n",
		},
		testData{
			input:          "3,9,7,9,10,9,4,9,99,-1,8",
			expectedOutput: 0,
			expectedZero:   3,
			userInput:      "8\n",
		},
		testData{
			input:          "3,9,7,9,10,9,4,9,99,-1,8",
			expectedOutput: 1,
			expectedZero:   3,
			userInput:      "7\n",
		},
		testData{
			input:          "3,3,1108,-1,8,3,4,3,99",
			expectedOutput: 1,
			expectedZero:   3,
			userInput:      "8\n",
		},
		testData{
			input:          "3,3,1108,-1,8,3,4,3,99",
			expectedOutput: 0,
			expectedZero:   3,
			userInput:      "7\n",
		},
		testData{
			input:          "3,3,1107,-1,8,3,4,3,99",
			expectedOutput: 0,
			expectedZero:   3,
			userInput:      "8\n",
		},
		testData{
			input:          "3,3,1107,-1,8,3,4,3,99",
			expectedOutput: 1,
			expectedZero:   3,
			userInput:      "7\n",
		},
		// Day 5 tests (jump operators)
		testData{
			input:          "3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9",
			expectedOutput: 0,
			expectedZero:   3,
			userInput:      "0\n",
		},
		testData{
			input:          "3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9",
			expectedOutput: 1,
			expectedZero:   3,
			userInput:      "999\n",
		},
		testData{
			input:          "3,3,1105,-1,9,1101,0,0,12,4,12,99,1",
			expectedOutput: 0,
			expectedZero:   3,
			userInput:      "0\n",
		},
		testData{
			input:          "3,3,1105,-1,9,1101,0,0,12,4,12,99,1",
			expectedOutput: 1,
			expectedZero:   3,
			userInput:      "999\n",
		},
	}

	var mockInput bytes.Buffer
	var mockOutput bytes.Buffer

	for _, td := range testDatas {
		ici := NewIntCodeInterpreter(td.input)

		mockInput.Write([]byte(td.userInput))

		zero := ici.Process(&mockInput, &mockOutput)

		if zero != td.expectedZero {
			t.Errorf("Expected %d but was %d\n", td.expectedZero, zero)
		}

		out := string(mockOutput.Bytes())
		if out != "" {
			actual, _ := strconv.Atoi(out[:strings.Index(out, "\n")])

			if actual != td.expectedOutput {
				t.Errorf("Expected %d but was %d\n", td.expectedOutput, actual)
			}
		}

		mockOutput.Reset()
	}
}
