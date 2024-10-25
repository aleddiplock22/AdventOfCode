package main

import (
	"fmt"
	"strconv"
)

type Solution struct {
	day   string
	part1 string
	part2 string
}

func GetInputPath(day int) string {
	return fmt.Sprintf("./inputs/%v/input.txt", day)
}
func GetExamplePath(day int) string {
	return fmt.Sprintf("./inputs/%v/example.txt", day)
}

func GetAllSolutions(day int) []Solution {
	solutions := []Solution{}
	for range day {
		day_str := strconv.Itoa(day)
		if dayfunc, exists := dayMap[day_str]; exists {
			solutions = append(solutions, dayfunc())
		} else {
			panic(fmt.Sprintf("Day %v not in dayMap! Please update it!", day_str))
		}
	}
	return solutions
}

type SolutionFuncType func() Solution

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
	res := DayFunc()
	fmt.Printf("Day %v\n\tPart 1: %v\n\tPart 2: %v\n", res.day, res.part1, res.part2)
}
