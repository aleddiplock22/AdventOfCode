package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const P1_EXAMPLE = "P1_EXAMPLE"
const P2_EXAMPLE = "P2_EXAMPLE"
const P1_INPUT = "P1_INPUT"
const P2_INPUT = "P2_INPUT"

const RED_FORE = "\x1b[31m"
const BOLD = "\x1b[1m"
const GREEN_FORE = "\x1b[32m"
const RESET = "\x1b[39m"
const RESET_STYLE = "\x1b[0m"

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
	fmt.Println(BOLD, RED_FORE, "---------------------------------", RESET)
	fmt.Println(GREEN_FORE, "  🎅🎄❄️  AOC 2024 - Day", res_p1.day, "❄️ 🎄🎅", RESET)
	fmt.Println(RED_FORE, " ---------------------------------", RESET)
	fmt.Printf("%v%vPart 1:%v\n\tExample: %v%v\n\t%vInput: %v%v\n%vPart 2:%v\n\tExample: %v%v\n\t%vInput: %v%v\n%v%v", GREEN_FORE, BOLD, RED_FORE, GREEN_FORE, res_p1.example_ans, RED_FORE, GREEN_FORE, res_p1.input_ans, GREEN_FORE, RED_FORE, GREEN_FORE, res_p2.example_ans, RED_FORE, GREEN_FORE, res_p2.input_ans, RESET, RESET_STYLE)
}

func RunVisualisation(day int) {
	if day != 10 {
		panic("I don't have visualisations for other days yet!!!!")
	}
	DoDay10Visualisation()
}

func AssertExample(expected string, result string, part int) {
	if expected != result {
		panic(fmt.Sprintf("Failed example p%v. Expected: %v, Got: %v\n", part, expected, result))
	}
}

func readInput(filepath string) string {
	// Get full file as string
	file, err := os.ReadFile(filepath)
	if err != nil {
		panic("ERROR READING FILE!")
	}
	file_content := string(file)
	return file_content
}

func readStringGrid(filepath string) [][]string {
	// Get grid input as an grid[col][row] of strings
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	grid := [][]string{}
	for scanner.Scan() {
		line := scanner.Text()
		chars := strings.Split(line, "")
		grid = append(grid, chars)
	}
	return grid
}

func readIntGrid(filepath string) [][]int {
	// Get grid input as an grid[col][row] of strings
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	grid := [][]int{}
	for scanner.Scan() {
		line := scanner.Text()
		chars := strings.Split(line, "")
		var row []int
		for _, char := range chars {
			digit, err := strconv.Atoi(char)
			if err != nil {
				panic("trouble parsing grid value as digit")
			}
			row = append(row, digit)
		}
		grid = append(grid, row)
	}
	return grid
}
