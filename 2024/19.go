package main

import (
	"fmt"
	"strings"
)

func day19(part2 bool) Solution {
	example_filepath := GetExamplePath(19)
	input_filepath := GetInputPath(19)
	if !part2 {
		example_p1 := Part1_19(example_filepath)
		AssertExample("6", example_p1, 1)
		return Solution{
			"19",
			example_p1,
			Part1_19(input_filepath),
		}
	} else {
		example_p2 := Part2_19(example_filepath)
		return Solution{
			"19",
			example_p2,
			Part2_19(input_filepath),
		}
	}
}

func ParseDay19(filepath string) (options []string, targets []string) {
	raw_input := readInput(filepath)
	parts := strings.Split(raw_input, "\r\n\r\n")
	options = append(options, strings.Split(parts[0], ", ")...)
	targets = append(targets, strings.Split(parts[1], "\r\n")...)
	return options, targets
}

func IsViableTowelCombo(target string, options []string) bool {
	pointer := 0
	tried := make([][]int, len(target))
	for i := range tried {
		tried[i] = make([]int, len(options))
	}

	var previous_pointer int
	var previous_j int

	for {
		if pointer == len(target) {
			return true
		}
		sub_found := false
		if RowFilledWithOnes(tried, pointer) {
			return false
		}
		for j, option := range options {
			if tried[pointer][j] == 1 {
				continue
			}
			n := len(option)
			if pointer+n > len(target) {
				tried[pointer][j] = 1
				continue
			}
			sub_target := target[pointer : pointer+n]

			if sub_target == option {
				sub_found = true
				previous_j = j
				previous_pointer = pointer
				pointer += n
				break
			}
		}
		if !sub_found {
			tried[previous_pointer][previous_j] = 1 // this combo caused issue, so dont allow it!
			if pointer == 0 {
				return false
			}
			pointer = 0
		}
	}
}

func RowFilledWithOnes(arr [][]int, i int) bool {
	for j := range arr[i] {
		if arr[i][j] != 1 {
			return false
		}
	}
	return true
}

func Part1_19(filepath string) string {
	options, targets := ParseDay19(filepath)

	var count int
	for _, target := range targets {
		if IsViableTowelCombo(target, options) {
			count++
		}
	}

	return fmt.Sprintf("%d", count)
}

func Part2_19(filepath string) string {
	options, targets := ParseDay19(filepath)

	var count int
	for _, target := range targets {
		count += NumViableTowelComboRecurisve(target, options)
		fmt.Printf("did target %v, now count=%v.\n", target, count)
	}

	return fmt.Sprintf("%d", count)
}

func NumViableTowelComboRecurisve(target string, options []string) int {
	/*
		TODO: UNDERSTAND, write pen and paper, then reimplement without looking,
		then allowed to use this!
	*/

	dp := make([]int, len(target)+1)
	// dp[i] represents number of valid combinations for substring target[:i]
	dp[0] = 1 // Empty string has one valid combination

	for i := 1; i <= len(target); i++ {
		for _, option := range options {
			if i >= len(option) && target[i-len(option):i] == option {
				// match -> add the number of combinations that were possible before this match
				// (dp[i-len(option)]) to the current position (dp[i]).
				dp[i] += dp[i-len(option)]
			}
		}
	}

	return dp[len(target)]
}
