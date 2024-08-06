package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
)

func main() {
	const EXAMPLE1_FILEPATH = "example1.txt"
	const EXAMPLE2_FILEPATH = "example2.txt"
	const INPUT_FILEPATH = "input.txt"

	part1_ans := Part1(INPUT_FILEPATH)

	example1_p2 := Part2(EXAMPLE1_FILEPATH)
	example2_p2 := Part2(EXAMPLE2_FILEPATH)
	part2_ans := Part2(INPUT_FILEPATH)

	fmt.Println("---Day 10---")
	fmt.Printf("[Part 1]\n\tAnswer: %v\n", part1_ans)
	fmt.Printf("[Examples P2]\n\t(1) Expected: 8, Got: %v\n\t(2) Expected: 19208, Got: %v\n", example1_p2, example2_p2)
	fmt.Printf("[Part 2]\n\tAnswer: %v\n", part2_ans)
}

func parseInput(filepath string) (nums []int) {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var line string
	var num int
	for scanner.Scan() {
		line = scanner.Text()
		num, _ = strconv.Atoi(line)
		nums = append(nums, num)
	}
	return nums
}

func getDiffsCountsMap(nums []int) map[int]int {
	counts := make(map[int]int)
	ones, twos, threes := 0, 0, 0
	var diff int
	for i := 0; i < len(nums)-1; i++ {
		diff = nums[i+1] - nums[i]
		switch diff {
		case 1:
			ones++
		case 2:
			twos++
		case 3:
			threes++
		default:
			panic("Difference wasn't 1-3?")
		}
	}
	counts[1], counts[2], counts[3] = ones, twos, threes
	return counts
}

func setUpNums(nums []int) []int {
	slices.Sort(nums)
	nums = append(nums, nums[len(nums)-1]+3)
	nums = append([]int{0}, nums...)
	return nums
}

func Part1(filepath string) int {
	nums := setUpNums(parseInput(filepath))
	counts := getDiffsCountsMap(nums)
	return counts[1] * counts[3]
}

func computePart2(nums []int, index int, cache map[int]int) int {
	/*
		Think of a tree, can split into paths
		(
			i.e. the node can jump to the next node (by definition),
				the node after that (if within range),
				and one after that (similarly if in range)
		)
		Add to total every time we split, and see if we're at the end, if so only 1 way.
	*/

	if index == len(nums)-1 {
		return 1
	}

	cache_val, exists := cache[index]
	if exists {
		return cache_val
	}

	total := 0

	var next_idx int
	for i := 1; i <= 3; i++ {
		next_idx = index + i
		if next_idx <= len(nums)-1 {
			if nums[next_idx]-nums[index] <= 3 {
				total = total + computePart2(nums, next_idx, cache)
			}
		}
	}

	cache[index] = total
	return total
}

func Part2(filepath string) int {
	nums := setUpNums(parseInput(filepath))
	return computePart2(nums, 0, make(map[int]int))
}
