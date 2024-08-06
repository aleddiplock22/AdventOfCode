package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	const EXAMPLE_FILEPATH = "example.txt"
	const INPUT_FILEPATH = "input.txt"

	example_p1 := Part1(EXAMPLE_FILEPATH, 5)
	input_p1 := Part1(INPUT_FILEPATH, 25)
	example_p2 := Part2(EXAMPLE_FILEPATH, 5)
	input_p2 := Part2(INPUT_FILEPATH, 25)

	fmt.Println("---Day 9---")
	fmt.Printf("[Example P1] Expected: 127, Answer: %v\n", example_p1)
	fmt.Printf("[Part 1] Answer: %v\n", input_p1)
	fmt.Printf("[Example P2] Expected: 62, Answer: %v\n", example_p2)
	fmt.Printf("[Part 2] Answer: %v\n", input_p2)
}

func parseInput(filepath string) []int {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var numbers_output []int
	var line string
	var num int
	for scanner.Scan() {
		line = scanner.Text()
		num, _ = strconv.Atoi(line)
		numbers_output = append(numbers_output, num)
	}
	return numbers_output
}

func hasRelevantSum(candidate int, nums []int) bool {
	for i := 0; i < len(nums)-1; i++ {
		for j := i + 1; i != j && j < len(nums); j++ {
			if nums[i]+nums[j] == candidate {
				return true
			}
		}
	}
	return false
}

func findNonPatternFollower(nums []int, preamble int) (answer int) {
	for i := preamble; i < len(nums); i++ {
		if !hasRelevantSum(nums[i], nums[i-preamble:i]) {
			return nums[i]
		}
	}
	panic("Didn't find non pattern follower!")
}

func sumSlice(nums []int) (ans int) {
	for _, n := range nums {
		ans = ans + n
	}
	return ans
}

func findContiguousSetThatSumsToK(k int, nums []int) (set []int) {
	var candidate_set []int
	for i := 0; i < len(nums)-1; i++ {
		for j := i + 1; j < len(nums); j++ {
			candidate_set = nums[i : j+1]
			if sumSlice(candidate_set) == k {
				return candidate_set
			}
		}
	}
	panic("Couldn't find valid contiguous set!")
}

func getSumOfMinMaxInSlice(nums []int) int {
	minimum := math.MaxInt
	maximum := math.MinInt
	for _, n := range nums {
		minimum = min(minimum, n)
		maximum = max(maximum, n)
	}
	return minimum + maximum
}

func Part1(filepath string, preamble_length int) int {
	nums := parseInput(filepath)
	return findNonPatternFollower(nums, preamble_length)
}

func Part2(filepath string, preamble_length int) int {
	nums := parseInput(filepath)
	candidate := findNonPatternFollower(nums, preamble_length)
	contiguous_set := findContiguousSetThatSumsToK(candidate, nums)
	return getSumOfMinMaxInSlice(contiguous_set)
}
