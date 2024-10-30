package main

import (
	"fmt"
	"strconv"
)

const P1_EXAMPLE = "P1_EXAMPLE"
const P2_EXAMPLE = "P2_EXAMPLE"
const P1_INPUT = "P1_INPUT"
const P2_INPUT = "P2_INPUT"

type Solution struct {
	day         string
	example_ans string
	input_ans   string
}

func GetInputPath(day int) string {
	return fmt.Sprintf("./inputs/%v/input.txt", day)
}
func GetExamplePath(day int) string {
	return fmt.Sprintf("./inputs/%v/example.txt", day)
}

func GetAllSolutions(day int) [][2]Solution {
	solutions := [][2]Solution{}
	for i := range day {
		day_str := strconv.Itoa(i + 1)
		if dayfunc, exists := dayMap[day_str]; exists {
			solutions_p1 := dayfunc(false)
			solutions_p2 := dayfunc(true)
			solutions = append(solutions, [2]Solution{solutions_p1, solutions_p2})
		} else {
			panic(fmt.Sprintf("Day %v not in dayMap! Please update it!", day_str))
		}
	}
	return solutions
}

type SolutionFuncType func(bool) Solution

var dayMap = map[string]SolutionFuncType{
	"1":  day01,
	"2":  day02,
	"3":  day03,
	"4":  day04,
	"5":  day05,
	"6":  day06,
	"7":  day07,
	"8":  day08,
	"9":  day09,
	"10": day10,
	"11": day11,
	"12": day12,
	"13": day13,
	"14": day14,
	"15": day15,
	"16": day16,
	"17": day17,
	"18": day18,
	"19": day19,
	"20": day20,
	"21": day21,
	"22": day22,
	"23": day23,
	"24": day24,
	"25": day25,
}

func SingleDayDevelopment(DayFunc SolutionFuncType) {
	res_p1 := DayFunc(false)
	res_p2 := DayFunc(true)
	fmt.Printf("Day %v\n\tPart 1: %v\n\tPart 2: %v\n", res_p1.day, res_p1.input_ans, res_p2.input_ans)
}
