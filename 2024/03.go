package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func day03(part2 bool) Solution {
	example_p1 := Part1_03(GetExamplePath(3))
	AssertExample("161", example_p1, 1)

	example_p2 := Part2_03("./inputs/3/example2.txt")
	AssertExample("48", example_p2, 2)

	input := GetInputPath(3)
	if !part2 {
		return Solution{
			"03",
			example_p1,
			Part1_03(input),
		}
	} else {
		return Solution{
			"03",
			example_p2,
			Part2_03(input),
		}
	}
}

func DoMulOperation(mul_str string) int {
	s := mul_str[:len(mul_str)-1]
	parts := strings.Split(s, ",")
	left, right := parts[0], parts[1]
	l, err := strconv.Atoi(left)
	if err != nil {
		panic("trouble parsing left to int")
	}
	r, err := strconv.Atoi(right)
	if err != nil {
		panic("trouble parsing right to int")
	}
	return l * r
}

func Part1_03(filepath string) string {
	re := regexp.MustCompile(`mul\(\d{1,3}\,\d{1,3}\)`)
	matches := re.FindAllSubmatch([]byte(readInput(filepath)), -1)
	var total int
	for _, match := range matches {
		s := string(match[0])[4:]
		total += DoMulOperation(s)
	}
	return fmt.Sprintf("%v", total)
}

func Part2_03(filepath string) string {
	input := readInput(filepath)

	re := regexp.MustCompile(`mul\(\d{1,3}\,\d{1,3}\)`)
	matches := re.FindAllStringSubmatchIndex(input, -1)

	dont_re := regexp.MustCompile(`don\'t\(\)`)
	dont_matches := dont_re.FindAllStringSubmatchIndex(input, -1)

	do_re := regexp.MustCompile(`do\(\)`)
	do_matches := do_re.FindAllStringSubmatchIndex(input, -1)

	var total int
	var DONT_idx int
	for _, match := range matches {
		s := input[match[0]:match[1]][4:]
		valid := true
		for j, dont_range := range dont_matches {
			if dont_range[1] <= match[0] {
				DONT_idx = j
				valid = false
			} else {
				break
			}
		}
		if !valid {
			dont_upper_bound := dont_matches[DONT_idx][1]
			for _, do_range := range do_matches {
				if do_range[0] >= dont_upper_bound && do_range[1] <= match[0] {
					valid = true
					break
				}
			}
		}
		if valid {
			total += DoMulOperation(s)
		}
	}
	return fmt.Sprintf("%v", total)
}
