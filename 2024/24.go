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
		AssertExample("2024", example_p1, 1)
		return Solution{
			"24",
			example_p1,
			Part1_24(input_filepath),
		}
	} else {
		example_p2 := Part2_24("./inputs/24/input.txt")
		return Solution{
			"24",
			example_p2,
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

func ParseInput24(filepath string) (WiresState map[string]int, sequence_part string) {
	raw_input := readInput(filepath)
	raw_parts := strings.Split(raw_input, "\r\n\r\n")
	initial_part, sequence_part := raw_parts[0], raw_parts[1]

	WiresState = make(map[string]int)
	for _, line := range strings.Split(initial_part, "\r\n") {
		parts := strings.Split(line, ": ")
		value, err := strconv.Atoi(parts[1])
		if err != nil {
			panic("Trouble parsing initial wire state as int.")
		}
		WiresState[parts[0]] = value
	}

	return WiresState, sequence_part
}

type Operation struct {
	wire1    string
	wire2    string
	op       string
	destWire string
}

func HandleOperationWithQueue(WiresState map[string]int, operation Operation, OpQueue *[]Operation) {
	state1, known1 := WiresState[operation.wire1]
	state2, known2 := WiresState[operation.wire2]
	if known1 && known2 {
		WiresState[operation.destWire] = doOperation(state1, state2, operation.op)
	} else {
		*OpQueue = append(*OpQueue, operation)
	}
}

func Part1_24(filepath string) string {
	WiresState, sequence_part := ParseInput24(filepath)

	var OpQueue []Operation

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
		HandleOperationWithQueue(WiresState, Operation{wire1, wire2, op, destWire}, &OpQueue)
	}

	for len(OpQueue) > 0 {
		operation := OpQueue[0]
		OpQueue = OpQueue[1:]
		HandleOperationWithQueue(WiresState, operation, &OpQueue)
	}

	zBinary := make([]int, zSize+1)

	for wireName, wireState := range WiresState {
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

func Part2_24(filepath string) string {
	WiresState, sequence_part := ParseInput24(filepath)

	var OpQueue []Operation
	var xSize, ySize, zSize int
	for wireName := range WiresState {
		switch wireName[:1] {
		case "x":
			xVal, err := strconv.Atoi(wireName[1:])
			if err != nil {
				panic("trouble taking z val as int")
			}
			xSize = max(xVal, xSize)
		case "y":
			yVal, err := strconv.Atoi(wireName[1:])
			if err != nil {
				panic("trouble taking z val as int")
			}
			ySize = max(yVal, ySize)
		}
	}

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
		OpQueue = append(OpQueue, Operation{wire1, wire2, op, destWire})
	}

	Simulate := func(opQueue *[]Operation, wiresState map[string]int) (WiresState map[string]int) {
		WiresState = make(map[string]int)
		for k, v := range wiresState {
			WiresState[k] = v // copying each entry
		}
		OpQueue := []Operation{}
		OpQueue = append(OpQueue, (*opQueue)...) // copy contents

		for len(OpQueue) > 0 {
			operation := OpQueue[0]
			OpQueue = OpQueue[1:]
			HandleOperationWithQueue(WiresState, operation, &OpQueue)
		}
		return WiresState
	}

	CheckAnswer := func(wiresState map[string]int) bool {
		xBinary, yBinary, zBinary := make([]int, xSize+1), make([]int, ySize+1), make([]int, zSize+1)
		for wireName, wireState := range WiresState {
			index, err := strconv.Atoi(wireName[1:])
			if err != nil {
				continue
				// error handling, that's right
			}
			switch wireName[:1] {
			case "x":
				xBinary[index] = wireState
			case "y":
				yBinary[index] = wireState
			case "z":
				zBinary[index] = wireState
			}
		}
		var xAnswer, yAnswer, zAnswer int

		for i, val := range xBinary {
			xAnswer += val * int(math.Pow(2, float64(i)))
		}
		for i, val := range yBinary {
			yAnswer += val * int(math.Pow(2, float64(i)))
		}
		for i, val := range zBinary {
			zAnswer += val * int(math.Pow(2, float64(i)))
		}

		fmt.Printf("%v+%v =? %v, hmm: %v\n", xAnswer, yAnswer, zAnswer, xAnswer+yAnswer)
		return xAnswer+yAnswer == zAnswer
	}

	// so we have to somehow figure out which four wires to flip to get to this answer...
	// I tried a function to generate all combinations of four pairs of indexes and it was insanely large and slow,
	// so dont think we'll be brute forcing...!
	// above have setup 'Simulate' and 'Check' functions so if we have a systematic way of flipping, seeing result, then honing in like a binary search idk that might help...?
	// might need a map of the destWire name to index of that in original opQueue

	/*
		// example manipulate pairs to be swapped!

		for i, op := range OpQueue {
			if op.destWire == "z05" {
				OpQueue[i].destWire = "z00"
			} else if op.destWire == "z00" {
				OpQueue[i].destWire = "z05"
			} else if op.destWire == "z02" {
				OpQueue[i].destWire = "z01"
			} else if op.destWire == "z01" {
				OpQueue[i].destWire = "z02"
			}
		}
	*/

	WiresState = Simulate(&OpQueue, WiresState)
	CheckAnswer(WiresState)

	return fmt.Sprintf("%d", 0)
}
