package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	const EXAMPLE_FILEPATH = "example.txt"
	const INPUT_FILEPATH = "input.txt"

	example_p1 := Part1(EXAMPLE_FILEPATH)
	input_p1 := Part1(INPUT_FILEPATH)

	example_p2 := Part2(EXAMPLE_FILEPATH)
	input_p2 := Part2(INPUT_FILEPATH)

	fmt.Println("---Day 8---")
	fmt.Printf("[Example P1] Expected: 5, Answer: %v\n", example_p1)
	fmt.Printf("[Part 1] Answer: %v\n", input_p1)
	fmt.Printf("[Example P1] Expected: 8, Answer: %v\n", example_p2)
	fmt.Printf("[Part 2] Answer: %v\n", input_p2)
}

const NOP = "nop"
const ACC = "acc"
const JMP = "jmp"

func readInput(filepath string) string {
	file, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("ERROR READING FILE!")
	}
	file_content := string(file)
	return file_content
}

type instruction struct {
	operation string
	value     int
}

func getInstruction(instruction_string string) instruction {
	parts := strings.Split(instruction_string, " ")
	lhs := parts[0]
	rhs := parts[1]
	val, _ := strconv.Atoi(rhs[1:])
	if rhs[0] == '+' {
		return instruction{
			operation: lhs,
			value:     val,
		}
	} else {
		return instruction{
			operation: lhs,
			value:     (-1 * val),
		}
	}
}

func parseStringInstructions(instructions []string) (formatted_instructions []instruction) {
	for _, string_instruction := range instructions {
		formatted_instructions = append(formatted_instructions, getInstruction(string_instruction))
	}
	return formatted_instructions
}

func getAccumalatorValue(instructions []instruction) (accumulator int, early bool) {
	var current_instruction instruction
	var cache []int
	for i := 0; i < len(instructions); {
		for _, seen_i := range cache {
			if i == seen_i {
				return accumulator, true
			}
		}
		cache = append(cache, i)
		current_instruction = instructions[i]
		switch current_instruction.operation {
		case NOP:
			i++
		case JMP:
			i = i + current_instruction.value
		case ACC:
			accumulator = accumulator + current_instruction.value
			i++
		default:
			panic("UNRECOGNISED OPERATION!")
		}
	}
	return accumulator, false
}

func getFormattedInstructionsFromInput(filepath string) (formatted_instructions []instruction) {
	instructions := strings.Split(readInput(filepath), "\r\n")
	formatted_instructions = parseStringInstructions(instructions)
	return formatted_instructions
}

func Part1(filepath string) int {
	formatted_instructions := getFormattedInstructionsFromInput(filepath)
	accumulator_value, _ := getAccumalatorValue(formatted_instructions)
	return accumulator_value
}

func flipInstructionAtIdx(instructions []instruction, index int) []instruction {
	new_instructions := make([]instruction, len(instructions))
	_ = copy(new_instructions, instructions)
	tmp := new_instructions[index]
	switch tmp.operation {
	case NOP:
		tmp.operation = JMP
	case JMP:
		tmp.operation = NOP
	}
	new_instructions[index] = tmp
	return new_instructions
}

func Part2(filepath string) int {
	formatted_instructions := getFormattedInstructionsFromInput(filepath)
	var new_instructions []instruction
	for i := 0; i < len(formatted_instructions); i++ {
		new_instructions = flipInstructionAtIdx(formatted_instructions, i)
		accumulator_value, early := getAccumalatorValue(new_instructions)
		if early {
			continue
		} else {
			return accumulator_value
		}
	}
	panic("Didn't find a non-early return accumulator value!")
}
