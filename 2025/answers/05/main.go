package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const TEST_FILEPATH string = "../../data/5/test.txt"
const INPUT_FILEPATH string = "../../data/5/input.txt"

func main() {
	fmt.Println("🎄 Advent Of Code 2025 | Day 5 🎅")

	Part1Test := Part1(TEST_FILEPATH)
	Part1 := Part1(INPUT_FILEPATH)
	fmt.Printf("[Part 1] Test: %v | Real: %v\n", Part1Test, Part1)

	Part2Test := Part2(TEST_FILEPATH)
	Part2 := Part2(INPUT_FILEPATH)
	fmt.Printf("[Part 2] Test: %v | Real: %v\n", Part2Test, Part2)
}

func parseInput(filepath string) (input [2][]string) {
	file, _ := os.ReadFile(filepath)
	fileContent := string(file)
	parts := strings.Split(fileContent, "\r\n\r\n")

	input[0] = strings.Split(parts[0], "\r\n")
	input[1] = strings.Split(parts[1], "\r\n")

	return input
}

func getValidIdRanges(rawRanges []string) (valid [][2]int) {
	for _, idRange := range rawRanges {
		nums := strings.Split(idRange, "-")
		num1, _ := strconv.Atoi(nums[0])
		num2, _ := strconv.Atoi(nums[1])
		added := false
		for i := range valid {
			l, r := valid[i][0], valid[i][1]
			if num1 < l && num1 <= r && num2 <= r && num2 >= l {
				valid[i][0] = num1
				added = true
			}
			if num2 >= l && num1 >= l && num1 <= r && num2 > r {
				valid[i][1] = num2
				added = true
			}
			if added {
				break
			}
		}
		if !added {
			valid = append(valid, [2]int{num1, num2})
		}
	}
	return valid
}

func Part1(filepath string) (answer int) {
	input := parseInput(filepath)
	valid := getValidIdRanges(input[0])

	for _, id := range input[1] {
		num, _ := strconv.Atoi(id)
		for _, validRange := range valid {
			l, r := validRange[0], validRange[1]
			if l <= num && num <= r {
				answer++
				break
			}
		}
	}
	return answer
}

func Part2(filepath string) (answer int) {
	input := parseInput(filepath)
	valid := getValidIdRanges(input[0])

	// Make the ranges exclusive of each other
	changed := true
	for changed {
		changed = false

		sort.Slice(valid, func(i, j int) bool {
			return valid[i][0] < valid[j][0]
		})

		// fmt.Println(valid)

	outer:
		for i, validRange := range valid {
			l, r := validRange[0], validRange[1]
			for j := range len(valid) {
				if i == j {
					continue
				}
				lj, rj := valid[j][0], valid[j][1]
				if lj > l && lj <= r && rj > r {
					valid[i][1] = lj - 1
					changed = true
					break outer
				}
				if rj >= l && rj < r && lj < l {
					valid[i][0] = rj + 1
					changed = true
					break outer
				}

				if l >= lj && l <= rj && r >= lj && r <= rj {
					// remove entry at index i as it is contained within another
					changed = true
					if i == len(valid)-1 {
						valid = valid[:i]
					} else {
						valid = append(valid[:i], valid[i+1:]...)
					}
					break outer
				}
			}
		}
	}

	for _, validRange := range valid {
		l, r := validRange[0], validRange[1]
		answer += r - l + 1
	}

	return answer
}
