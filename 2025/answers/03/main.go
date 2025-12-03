package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const TEST_FILEPATH string = "../../data/3/test.txt"
const INPUT_FILEPATH string = "../../data/3/input.txt"

func main() {
	fmt.Println("Advent of Code 2025 - Day 03")

	part1_test := Part1(TEST_FILEPATH)
	part1 := Part1(INPUT_FILEPATH)
	fmt.Printf("Part 1 (test): %v | Part 1 (answer): %v\n", part1_test, part1)

	part2_test := Part2(TEST_FILEPATH)
	part2 := Part2(INPUT_FILEPATH)
	fmt.Printf("Part 2 (test): %v | Part 2 (answer): %v\n", part2_test, part2)
}

func parseInput(filepath string) (joltages [][]int) {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		joltage_line := []int{}
		for _, s := range strings.Split(line, "") {
			num, err := strconv.Atoi(s)
			if err != nil {
				panic("trouble converting one of the joltage numbers!")
			}
			joltage_line = append(joltage_line, num)
		}
		joltages = append(joltages, joltage_line)
	}
	return joltages
}

func Part1(filepath string) int {
	joltages := parseInput(filepath)

	var total int
	for _, joltage_line := range joltages {
		lmax, rmax := 0, 0
		for i, joltage := range joltage_line {
			if joltage > lmax && i < len(joltage_line)-1 {
				lmax = joltage
				rmax = 0
				continue
			} else if joltage >= rmax {
				rmax = joltage
			}
		}
		total += lmax*10 + rmax
	}
	return total
}

func Part2(filepath string) int {
	joltages := parseInput(filepath)

	var total int
	for _, joltage_line := range joltages {
		var max_joltage [12]int
		for i, joltage := range joltage_line {
			for j := range 12 {
				if joltage > max_joltage[j] && len(joltage_line)-i >= 12-j {
					max_joltage[j] = joltage
					for k := j + 1; k < 12; k++ {
						max_joltage[k] = 0
					}
					break
				}
			}
		}
		var value int
		for i := range 12 {
			value += int(math.Pow(10, (11.0-float64(i)))) * max_joltage[i]
		}
		total += value
	}
	return total
}
