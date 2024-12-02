package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func day02(part2 bool) Solution {
	// trying some fun calling python to get the answer from fun dodgy one liners
	cmd := exec.Command("python", "./python_solutions/day2.py")
	output, err := cmd.CombinedOutput()
	if err != nil {
		panic("Failed running Python day2.py script!")
	}
	parts := strings.Split(string(output), ",")
	p1, p2 := parts[0], parts[1]
	if !part2 {
		return Solution{
			"02",
			"4",
			p1,
		}
	} else {
		return Solution{
			"02",
			"5",
			p2,
		}
	}

	// below is my actual Go implementation etc. etc.

	part1_example := Part1(GetExamplePath(2))
	AssertExample("2", part1_example, 1)

	part2_example := Part2(GetExamplePath(2))
	AssertExample("4", part2_example, 2)

	if !part2 {
		return Solution{
			"02",
			part1_example,
			Part1(GetInputPath(2)),
		}
	} else {
		return Solution{
			"02",
			part2_example,
			Part2(GetInputPath(2)),
		}
	}
}

func parseInput(filepath string) (reports [][]int) {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		digits_strings := strings.Split(line, " ")
		reports = append(reports, make([]int, len(digits_strings)))
		j := 0
		for _, digit_str := range digits_strings {
			digit, err := strconv.Atoi(digit_str)
			if err != nil {
				panic("Trouble parsing digit!")
			}
			reports[i][j] = digit
			j++
		}
		i++
	}
	return reports
}

func Part1(filepath string) string {
	return Solve(filepath, false)
}

func Part2(filepath string) string {
	return Solve(filepath, true)
}

func CheckNums(nums []int) bool {
	var increasing bool

	curr := nums[0]
	first := true
	for _, num := range nums[1:] {
		if first {
			if num > curr {
				increasing = true
			} else if num < curr {
				increasing = false
			} else {
				return false
			}
			first = false
		}

		if increasing {
			if num < curr || num-curr < 1 || num-curr > 3 {
				return false
			}
		} else {
			if num > curr || curr-num < 1 || curr-num > 3 {
				return false
			}
		}
		curr = num
	}
	return true
}

func getCombos(nums []int) [][]int {
	combos := [][]int{}
	for i := range nums {
		new_nums := []int{}
		new_nums = append(new_nums, nums[:i]...)
		new_nums = append(new_nums, nums[i+1:]...)
		combos = append(combos, new_nums)
	}
	return combos
}

func Solve(filepath string, second_chance bool) string {
	input := parseInput(filepath)

	var total_safe int
	for _, nums := range input {
		if second_chance {
			for _, combo := range getCombos(nums) {
				if CheckNums(combo) {
					total_safe++
					break
				}
			}
		} else {
			if CheckNums(nums) {
				total_safe++
			}
		}
	}

	return fmt.Sprintf("%v", total_safe)
}
