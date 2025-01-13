package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func day24(part2 bool) Solution {
	example_filepath := GetExamplePath(24)
	input_filepath := GetInputPath(24)
	if !part2 {
		example_p1 := Part1_24(example_filepath)
		return Solution{
			"24",
			example_p1,
			Part1_24(input_filepath),
		}
	} else {
		return Solution{
			"24",
			"example part 2",
			"input part 2",
		}
	}
}

func doOperation(state1, state2 int, op string) int {
	switch op {
	case "XOR":
		if state1 == 1 && state2 == 0 || state1 == 0 && state2 == 1 {
			return 1
		} else {
			return 0
		}
	case "OR":
		if state1 == 1 || state2 == 1 {
			return 1
		} else {
			return 0
		}
	case "AND":
		if state1 == 1 && state2 == 1 {
			return 1
		} else {
			return 0
		}
	default:
		panic("Unrecognised op in doOperation!")
	}
}

func Part1_24(filepath string) string {
	raw_input := readInput(filepath)
	raw_parts := strings.Split(raw_input, "\r\n\r\n")
	initial_part, sequence_part := raw_parts[0], raw_parts[1]

	WiresState := make(map[string]int)
	for _, line := range strings.Split(initial_part, "\r\n") {
		parts := strings.Split(line, ": ")
		value, err := strconv.Atoi(parts[1])
		if err != nil {
			panic("Trouble parsing initial wire state as int.")
		}
		WiresState[parts[0]] = value
	}

	type Operation struct {
		wire1    string
		wire2    string
		op       string
		destWire string
	}

	var OpQueue []Operation

	handleOp := func(operation Operation) {
		state1, known1 := WiresState[operation.wire1]
		state2, known2 := WiresState[operation.wire2]
		if known1 && known2 {
			WiresState[operation.destWire] = doOperation(state1, state2, operation.op)
		} else {
			OpQueue = append(OpQueue, operation)
		}
	}

	var zSize int
	for _, operationLine := range strings.Split(sequence_part, "\r\n") {
		parts := strings.Split(operationLine, " ")
		wire1, op, wire2, destWire := parts[0], parts[1], parts[2], parts[len(parts)-1]
		if destWire[:1] == "z" {
			zVal, err := strconv.Atoi(destWire[1:])
			if err != nil {
				panic("trouble taking z val as int")
			}
			zSize = max(zVal, zSize)
		}
		handleOp(Operation{wire1, wire2, op, destWire})
	}
	for len(OpQueue) > 0 {
		operation := OpQueue[0]
		OpQueue = OpQueue[1:]
		handleOp(operation)
	}

	zBinary := make([]int, zSize+1)

	for wireName, wireState := range WiresState {
		// in Go we're guaranteed to go alphabetical right?
		if wireName[:1] == "z" {
			zIndex, err := strconv.Atoi(wireName[1:])
			if err != nil {
				panic("trouble parsing zIndex.")
			}
			zBinary[zIndex] = wireState
		}
	}
	var answer int
	for i, val := range zBinary {
		answer += val * int(math.Pow(2, float64(i)))
	}

	return fmt.Sprintf("%d", answer)
}
