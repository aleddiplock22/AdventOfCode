package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func day13(part2 bool) Solution {
	example_filepath := GetExamplePath(13)
	input_filepath := GetInputPath(13)
	if !part2 {
		example_p1 := Part1_13(example_filepath)
		AssertExample("480", example_p1, 1)
		return Solution{
			"13",
			example_p1,
			Part1_13(input_filepath),
		}
	} else {
		return Solution{
			"13",
			"N/A",
			Part2_13(input_filepath),
		}
	}
}

type ClawMachineButton struct {
	x int
	y int
}

type ClawMachine struct {
	A     ClawMachineButton
	B     ClawMachineButton
	Prize ClawMachineButton // yes this makes total sense
}

func parseClawMachineInputs(filepath string) []ClawMachine {
	s := readInput(filepath)
	definitions := strings.Split(s, "\r\n\r\n")
	var claw_machine_defs []ClawMachine
	for _, raw_def := range definitions {
		parts := strings.Split(raw_def, "\r\n")
		part_a, part_b, part_prize := parts[0], parts[1], parts[2]

		// Button A
		parts_a := strings.Split(part_a[12:], ", Y+")
		part_a_x, err := strconv.Atoi(parts_a[0])
		if err != nil {
			panic("trouble parsing part_a_x as int")
		}
		part_a_y, err := strconv.Atoi(parts_a[1])
		if err != nil {
			panic("trouble parsing part_a_y as int")
		}
		A := ClawMachineButton{
			part_a_x,
			part_a_y,
		}

		// Button B
		parts_b := strings.Split(part_b[12:], ", Y+")
		part_b_x, err := strconv.Atoi(parts_b[0])
		if err != nil {
			panic("trouble parsing part_b_x as int")
		}
		part_b_y, err := strconv.Atoi(parts_b[1])
		if err != nil {
			panic("trouble parsing part_b_y as int")
		}
		B := ClawMachineButton{
			part_b_x,
			part_b_y,
		}

		// Prize
		parts_prize := strings.Split(part_prize[9:], ", Y=")
		part_prize_x, err := strconv.Atoi(parts_prize[0])
		if err != nil {
			panic("trouble parsing part_prize_x as int")
		}
		part_prize_y, err := strconv.Atoi(parts_prize[1])
		if err != nil {
			panic("trouble parsing part_prize_y as int")
		}
		prize := ClawMachineButton{
			part_prize_x,
			part_prize_y,
		}
		claw_machine_defs = append(claw_machine_defs, ClawMachine{
			A,
			B,
			prize,
		})
	}
	return claw_machine_defs
}

func Part1_13(filepath string) string {
	claw_machines := parseClawMachineInputs(filepath)
	var total int
	for _, claw_machine := range claw_machines {
		total += SolveClawMachine(claw_machine, false)
	}

	return fmt.Sprintf("%d", total)
}

func Part2_13(filepath string) string {
	claw_machines := parseClawMachineInputs(filepath)
	var total int
	for _, claw_machine := range claw_machines {
		total += SolveClawMachine(claw_machine, true)
	}

	return fmt.Sprintf("%d", total)
}

func SolveClawMachine(claw_eqn ClawMachine, part2 bool) int {
	Xa := float64(claw_eqn.A.x)
	Ya := float64(claw_eqn.A.y)
	Xb := float64(claw_eqn.B.x)
	Yb := float64(claw_eqn.B.y)
	Xp := float64(claw_eqn.Prize.x)
	Yp := float64(claw_eqn.Prize.y)
	if part2 {
		Xp += 10000000000000
		Yp += 10000000000000
	}

	b := (Xa*Yp - Ya*Xp) / (Xa*Yb - Ya*Xb)
	a := (Xp - Xb*b) / Xa

	// check whole numbers
	if b != math.Trunc(b) || a != math.Trunc(a) {
		return 0
	}
	if part2 {
		return int(a*3 + b)
	}

	if a >= 0 && a <= 100 && b >= 0 && b <= 100 {
		return int(a*3 + b)
	} else {
		return 0
	}
}
