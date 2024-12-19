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
		// DP approach
		example_p2_dp := Part2_19_dp(example_filepath)
		AssertExample("16", example_p2_dp, 2)
		// Recursion with cache
		example_p2_r := Part2_19_recursive(example_filepath)
		AssertExample("16", example_p2_r, 2)

		return Solution{
			"19",
			example_p2_r,
			Part2_19_recursive(input_filepath), // my recursion wins!!!
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

func Part2_19_dp(filepath string) string {
	options, targets := ParseDay19(filepath)

	var count int
	for _, target := range targets {
		count += NumViableTowelDP(target, options)
	}

	return fmt.Sprintf("%d", count)
}

func Part2_19_recursive(filepath string) string {
	options, targets := ParseDay19(filepath)
	cache := make(map[string]int)

	var count int
	for _, target := range targets {
		count += NumViableTowelComboRecurisve(target, options, cache)
	}

	return fmt.Sprintf("%d", count)
}

// I was quite pleased with this but it's crazy slow lol
func NumViableTowelComboRecurisve(target string, options []string, cache map[string]int) int {
	if answer, exists := cache[target]; exists {
		return answer
	}
	if len(target) == 0 {
		// succeed
		return 1
	}
	total := 0
	for _, option := range options {
		n := len(option)
		if len(target) < n {
			continue
		}
		sub_target := target[:n]

		if sub_target == option {
			ans := NumViableTowelComboRecurisve(target[n:], options, cache)
			cache[target[n:]] = ans
			total += ans
		}
	}
	return total
}

// this got spoiled a bit when I was asking AI how to optimise my recursive one, but I've learned and understood!!
// similar to leetcode problems I've done tbf, should've spotted it
func NumViableTowelDP(target string, options []string) int {
	valid_combinations := make([]int, len(target)+1)
	// valid_combinations[i] represents number of valid combinations for substring target[:i]
	valid_combinations[0] = 1 // Empty string has one valid combination

	for i := 1; i <= len(target); i++ {
		for _, option := range options {
			L := len(option)
			if i < L {
				continue
			}
			// i >= L just guarantees we can slice
			substring_window := target[i-L : i]
			// target[i-len(option):i] is a window into target from just before i up to it
			if substring_window == option {
				// match -> add the number of combinations that were possible before this match
				// (valid_combinations[i-len(option)]) to the current position (valid_combinations[i]).
				valid_combinations[i] += valid_combinations[i-L]
			}
		}
	}

	// thus valid_combinations[len(target)] gives us the total valid combinations for the entire substring,
	// since we have this cumulative sum thing going on
	return valid_combinations[len(target)]
}
