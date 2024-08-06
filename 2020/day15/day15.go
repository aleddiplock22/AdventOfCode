package main

import (
	"fmt"
	"strconv"
	"strings"
)

const INPUT = "19,0,5,1,10,13"
const EXAMPLE_1 = "0,3,6"

func main() {
	/*
		Not super optimised but didn't have to adjust for part2 - win!
	*/
	fmt.Println("---Day 15---")

	example_1_p1 := Part1(EXAMPLE_1)
	part1_ans := Part1(INPUT)

	fmt.Printf("[Example 1 P1] Expected: 436, Answer: %v\n", example_1_p1)
	fmt.Printf("[Part 1] Answer: %v\n", part1_ans)

	example_1_p2 := Part2(EXAMPLE_1)
	part2_ans := Part2(INPUT)

	fmt.Printf("[Example 1 P2] Expected: 175594, Answer: %v\n", example_1_p2)
	fmt.Printf("[Part 2] Answer: %v\n", part2_ans)
}

func getStartingNums(input string) []int {
	var starting_nums []int
	for _, char := range strings.Split(input, ",") {
		num, err := strconv.Atoi(char)
		if err != nil {
			panic("Tried to parse non-number!")
		}
		starting_nums = append(starting_nums, num)
	}
	return starting_nums
}

func findNthNumberInSequence(N int, starting_nums []int) int {
	turns := make(map[int]int, N)
	num_mentions := make(map[int][]int, N)
	for i, turn := range starting_nums {
		turns[i+1] = turn
		num_mentions[turn] = append(num_mentions[turn], i+1)
	}

	var current_turn int

	for turn_n := len(starting_nums) + 1; turn_n <= N; turn_n++ {
		previous_num_said := turns[turn_n-1]
		mentions_of_previous_num := num_mentions[previous_num_said]

		if len(mentions_of_previous_num) == 1 {
			// first time spoken, so we say: 0
			current_turn = 0
		} else {
			current_turn = mentions_of_previous_num[len(mentions_of_previous_num)-1] - mentions_of_previous_num[len(mentions_of_previous_num)-2]
		}

		turns[turn_n] = current_turn
		num_mentions[current_turn] = append(num_mentions[current_turn], turn_n)
	}

	return turns[N]
}

func Part1(input string) int {
	starting_nums := getStartingNums(input)
	n_2020 := findNthNumberInSequence(2020, starting_nums)

	return n_2020
}

func Part2(input string) int {
	starting_nums := getStartingNums(input)
	n_30000000 := findNthNumberInSequence(30000000, starting_nums)

	return n_30000000
}
