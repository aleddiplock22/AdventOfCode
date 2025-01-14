package main

import (
	"fmt"
	"strings"
)

func day25(part2 bool) Solution {
	example_filepath := GetExamplePath(25)
	input_filepath := GetInputPath(25)
	if !part2 {
		example_p1 := Part1_25(example_filepath)
		return Solution{
			"25",
			example_p1,
			Part1_25(input_filepath),
		}
	} else {
		return Solution{
			"25",
			"",
			"Just do the other days!",
		}
	}
}

type FivePinCombination struct {
	heights [5]int
}

func parseKeysAndLocks(filepath string) (keys, locks []FivePinCombination) {
	rawInput := readInput(filepath)
	for _, rawCombination := range strings.Split(rawInput, "\r\n\r\n") {
		parts := strings.Split(rawCombination, "\r\n")
		if parts[0] == "#####" {
			// lock
			var _lock_heights [5]int
			for _, line := range parts[1:] {
				for i, char := range strings.Split(line, "") {
					if char == "#" {
						_lock_heights[i]++
					}
				}
			}
			locks = append(locks, FivePinCombination{_lock_heights})
		} else if parts[len(parts)-1] == "#####" {
			var _key_heights [5]int
			for i := len(parts) - 2; i >= 0; i-- {
				for j, char := range strings.Split(parts[i], "") {
					if char == "#" {
						_key_heights[j]++
					}
				}
			}
			keys = append(keys, FivePinCombination{_key_heights})
		}
	}
	return keys, locks
}

func Part1_25(filepath string) string {
	keys, locks := parseKeysAndLocks(filepath)
	var total_valid int
	for _, key := range keys {
	lock_loop:
		for _, lock := range locks {
			for i := 0; i < 5; i++ {
				if key.heights[i]+lock.heights[i] > 5 {
					continue lock_loop
				}
			}
			total_valid++
		}
	}
	return fmt.Sprintf("%d", total_valid)
}
