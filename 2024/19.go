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
		return Solution{
			"19",
			"example part 2",
			"input part 2",
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
	// fmt.Println("LOOKING FOR: ", target)

	pointer := 0
	tried := make([][]int, len(target))
	for j := range tried {
		tried[j] = make([]int, len(options))
	}

	for {
		if pointer == len(target) {
			// fmt.Println("FOUND TARGET!", target)
			return true
		}
		// fmt.Println(pointer, target[pointer:])
		sub_found := false
		if RowFilledWithOnes(tried, pointer) {
			return false
		}
		for j, option := range options {
			// fmt.Println("option:", option)
			if tried[pointer][j] == 1 {
				// fmt.Println("already tried")
				continue
			}
			n := len(option)
			if pointer+n > len(target) {
				tried[pointer][j] = 1
				continue
			}
			sub_target := target[pointer : pointer+n]
			tried[pointer][j] = 1

			// fmt.Println("sub:", sub_target, "option", option)
			if sub_target == option {
				// fmt.Println("sub == option!")
				pointer += n
				sub_found = true
				break
			}
		}
		if !sub_found {
			// fmt.Println("not sub found, tried=", tried)
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
