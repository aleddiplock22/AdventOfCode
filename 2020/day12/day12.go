package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	const EXAMPLE_FILEPATH = "example.txt"
	const INPUT_FILEPATH = "input.txt"
	fmt.Println("---Day 12---")

	example_p1 := Part1(EXAMPLE_FILEPATH)
	part1_ans := Part1(INPUT_FILEPATH)

	fmt.Printf("[Example P1] Expected: 25, Answer: %v\n", example_p1)
	fmt.Printf("[Part 1] Answer: %v\n", part1_ans)

	example_p2 := Part2(EXAMPLE_FILEPATH)
	part2_ans := Part2(INPUT_FILEPATH)

	fmt.Printf("[Example P2] Expected: 286, Answer: %v\n", example_p2)
	fmt.Printf("[Part 2] Answer: %v\n", part2_ans)
}

type Instruction struct {
	instruction_type rune
	value            int
}

type Vessel struct {
	current_direction rune
	x_pos             int // cartesian
	y_pos             int
	part2             bool // convenient way to not rewrite much code for p2 and reuse object not quite how they were designed !
}

func (v *Vessel) turn(instruction_input Instruction) {
	DIRECTIONS := []rune{'N', 'E', 'S', 'W'}
	direction_adjustment := func(adjustment int, dir rune) int {
		if dir == 'R' {
			return adjustment
		}
		return (4 - adjustment)
	}
	var representive_pos int
	switch v.current_direction {
	case 'N':
		representive_pos = 0
	case 'E':
		representive_pos = 1
	case 'S':
		representive_pos = 2
	case 'W':
		representive_pos = 3
	default:
		panic("What direction are we even turning??")
	}
	var adjustment int
	switch instruction_input.value {
	case 90:
		adjustment = 1
	case 180:
		adjustment = 2
	case 270:
		adjustment = 3
	default:
		panic("Unrecognised turning adjustment!")
	}

	if !v.part2 {
		representive_pos = (representive_pos + direction_adjustment(adjustment, instruction_input.instruction_type)) % 4
		new_direction := DIRECTIONS[representive_pos]
		v.current_direction = new_direction
		return
	}
	// part 2
	R_val := direction_adjustment(adjustment, instruction_input.instruction_type)
	var new_x int
	var new_y int
	switch R_val {
	case 1:
		// R90
		new_x = v.y_pos
		new_y = -1 * v.x_pos
	case 2:
		// R180
		new_x = -1 * v.x_pos
		new_y = -1 * v.y_pos
	case 3:
		// R270
		new_x = -1 * v.y_pos
		new_y = v.x_pos
	default:
		panic("Unknown R val?")
	}
	v.x_pos = new_x
	v.y_pos = new_y
}

func (v *Vessel) move(instruction_input Instruction) {
	switch instruction_input.instruction_type {
	// specific direction to move
	case 'N':
		v.y_pos = v.y_pos + instruction_input.value
	case 'E':
		v.x_pos = v.x_pos + instruction_input.value
	case 'S':
		v.y_pos = v.y_pos - instruction_input.value
	case 'W':
		v.x_pos = v.x_pos - instruction_input.value
	// forward in current direction
	case 'F':
		instruction_input.instruction_type = v.current_direction
		v.move(instruction_input)
	// turn
	case 'L', 'R':
		v.turn(instruction_input)
	default:
		panic("Didn't recognise instruction type provided!")
	}
}

func parseLine(line string) Instruction {
	chars := []rune(line)
	instructions_type := chars[0]
	instruction_value, _ := strconv.Atoi(string(chars[1:]))
	return Instruction{
		instructions_type,
		instruction_value,
	}
}

func Part1(filepath string) int {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	vessel := Vessel{
		current_direction: 'E',
		x_pos:             0,
		y_pos:             0,
		part2:             false,
	}
	for scanner.Scan() {
		line := scanner.Text()
		current_instruction := parseLine(line)
		// fmt.Printf("Vessel direction: %v, Vessel (x, y)=(%v, %v), NEXT INSTRUCTION: %v %v\n", string(vessel.current_direction), vessel.x_pos, vessel.y_pos, string(current_instruction.instruction_type), current_instruction.value)
		vessel.move(current_instruction)
	}

	return int(math.Abs(float64(vessel.x_pos)) + math.Abs(float64(vessel.y_pos)))
}

func Part2(filepath string) int {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	vessel := Vessel{
		current_direction: 'E',
		x_pos:             0,
		y_pos:             0,
		part2:             true,
	}
	waypoint := Vessel{
		current_direction: 'N',
		x_pos:             10,
		y_pos:             1,
		part2:             true,
	}
	for scanner.Scan() {
		line := scanner.Text()
		current_instruction := parseLine(line)
		// fmt.Printf("Vessel\n\tdirection: %v,\n\t(x, y)=(%v, %v)\nWaypoint\n\t(x, y)=(%v, %v)\nNEXT INSTRUCTION: %v %v\n", string(vessel.current_direction), vessel.x_pos, vessel.y_pos, waypoint.x_pos, waypoint.y_pos, string(current_instruction.instruction_type), current_instruction.value)
		if current_instruction.instruction_type == 'F' {
			// only case where the vessel itself moves, 'towards' the relative waypoint
			vessel.x_pos = vessel.x_pos + waypoint.x_pos*current_instruction.value
			vessel.y_pos = vessel.y_pos + waypoint.y_pos*current_instruction.value
		} else {
			waypoint.move(current_instruction)
		}
	}
	// fmt.Printf("Vessel\n\tdirection: %v,\n\t(x, y)=(%v, %v)\nWaypoint\n\t(x, y)=(%v, %v)\n", string(vessel.current_direction), vessel.x_pos, vessel.y_pos, waypoint.x_pos, waypoint.y_pos)

	return int(math.Abs(float64(vessel.x_pos)) + math.Abs(float64(vessel.y_pos)))
}
