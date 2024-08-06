package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func main() {
	const EXAMPLE = "389125467"
	const INPUT = "624397158"

	fmt.Println("---Day 23---")

	fmt.Println("[Example P1] Expected: 67384529, Answer:", Part1(EXAMPLE))
	fmt.Println("[Part 1] Answer:", Part1(INPUT))

	// unfortunately it's not a 'find a quick cycle' kind of problem and this isn't going to work out for part 2.
	// Going to implement Jonathon Paulson's idea of a list that just updates which 'next' value a given point is pointing at https://github.com/jonathanpaulson/AdventOfCode/blob/master/2020/23.py

	// THIS IS A LINKED LIST, where each element points to to the value at the 'next' element despite no 'physical' proximity in memory or actual ordering.

	// JP's idea:
	fmt.Println("PART 2 ANSWER:", Part2_NextInListApproach(INPUT))

}

type Circle struct {
	numbers         []int
	current_cup_idx int
	current_cup_val int
}

func (circle Circle) FindIdxOfX(X int) int {
	idx := slices.Index(circle.numbers, X)
	if idx == -1 {
		panic("Didn't find X!")
	}
	return idx
}

func (circle Circle) NextThreeNumbersFromCurrent() (next_three []int) {
	max_idx := len(circle.numbers) - 1
	next := circle.current_cup_idx
	var next_three_idx []int
	for j := 0; j < 3; j++ {
		if next+1 <= max_idx {
			next += 1
		} else {
			next = 0
		}

		next_three_idx = append(next_three_idx, next)
	}
	for _, idx := range next_three_idx {
		next_three = append(next_three, circle.numbers[idx])
	}
	return next_three
}

func (circle *Circle) PlaceThreeDownAfterVal(destination_val int, picked_up []int, still_in_circle []int) {
	destination_idx := slices.Index(still_in_circle, destination_val)
	new_nums := make([]int, len(still_in_circle)+len(picked_up))

	copy(new_nums[:destination_idx+1], still_in_circle[:destination_idx+1])
	copy(new_nums[destination_idx+1:], picked_up)
	copy(new_nums[destination_idx+1+len(picked_up):], still_in_circle[destination_idx+1:])

	circle.numbers = new_nums
}

func (circle Circle) GetLabelAfterOne() string {
	one_idx := circle.FindIdxOfX(1)
	new_num_str := ""
	for _, num := range circle.numbers[one_idx+1:] {
		new_num_str += strconv.Itoa(num)
	}
	for _, num := range circle.numbers[:one_idx] {
		new_num_str += strconv.Itoa(num)
	}
	return new_num_str
}

func (circle Circle) MultiplyTwoCupsAfterCupOne() int {
	idx := circle.FindIdxOfX(1)
	if idx+1 < len(circle.numbers) {
		idx += 1
	} else {
		idx = 0
	}
	val1 := circle.numbers[idx]
	if idx+1 < len(circle.numbers) {
		idx += 1
	} else {
		idx = 0
	}
	val2 := circle.numbers[idx]
	return val1 * val2
}

func (circle *Circle) Move() {
	findDestinationVal := func(starting_value int, circle_vals []int) int {
		minimum := slices.Min(circle_vals)
		destination := starting_value - 1
		for {
			if slices.Contains(circle_vals, destination) {
				return destination
			} else {
				if destination < minimum {
					return slices.Max(circle_vals)
				}
				destination--
			}
		}
	}

	current_card_val := circle.numbers[circle.current_cup_idx]
	picked_up := []int{}
	still_in_circle := []int{}

	next_three_nums := circle.NextThreeNumbersFromCurrent()

	for _, val := range circle.numbers[circle.current_cup_idx:] {
		if slices.Contains(next_three_nums, val) {
			picked_up = append(picked_up, val)
		} else {
			still_in_circle = append(still_in_circle, val)
		}
	}
	for _, val := range circle.numbers[:circle.current_cup_idx] {
		if slices.Contains(next_three_nums, val) {
			picked_up = append(picked_up, val)
		} else {
			still_in_circle = append(still_in_circle, val)
		}
	}

	destination_val := findDestinationVal(current_card_val, still_in_circle)
	// put picked up down next to next_val
	circle.PlaceThreeDownAfterVal(destination_val, picked_up, still_in_circle) // BOTTLENECK

	old_current_idx := circle.FindIdxOfX(circle.current_cup_val)
	var new_current_idx int
	if old_current_idx+1 <= len(circle.numbers)-1 {
		new_current_idx = old_current_idx + 1
	} else {
		new_current_idx = 0
	}
	circle.current_cup_idx = new_current_idx
	circle.current_cup_val = circle.numbers[new_current_idx]
}

func (circle *Circle) ExtendToMillion() {
	new_nums := make([]int, 1000000)
	max_val := slices.Max(circle.numbers)
	copy(new_nums, circle.numbers) // only copies as many elements as are in the smaller arr
	for i := max_val + 1; i <= 1000000; i++ {
		circle.numbers = append(circle.numbers, i)
	}
}

func parseInput(input string) Circle {
	values := []int{}
	for _, s := range strings.Split(input, "") {
		value, err := strconv.Atoi(s)
		if err != nil {
			panic("Couldn't parse input num")
		}
		values = append(values, value)
	}
	return Circle{
		values,
		0,
		values[0],
	}
}

// func Hash(slice []int) uint32 {
// 	h := fnv.New32a()

// 	for _, i := range slice {
// 		// Convert int to byte slice
// 		bs := []byte(fmt.Sprintf("%d", i))
// 		h.Write(bs)
// 	}
// 	return h.Sum32()
// }

func KeyifyArr(slice []int) string { return fmt.Sprintf("%q", slice) }

func Part1(input string) string {
	circle := parseInput(input)
	for i := 0; i < 100; i++ {
		circle.Move()
	}
	return circle.GetLabelAfterOne()
}

// UNUSED FIRST ATTEMPT AT PART 2:
func Part2(input string) int {
	// need to set up with a million cups
	// and do ten million moves

	// TODO:
	// 	- gonna have to find a cycle here!
	// 	- going to have eliminate as many loops as possible I think! - Done: by using more copy operations & pre-allocating and also using standard lib funcs like Slices.Contains and Slices.Index
	fmt.Println()

	seen := make(map[string]bool)

	circle := parseInput(input)
	circle.ExtendToMillion()
	for i := 0; i < 10000000; i++ { // break when we find a cycle and then do modulo and just have to do a few steps was the plan
		// if i%100 == 0 {
		// 	fmt.Println(i, "/", 10000)
		// }
		_, already_seen := seen[KeyifyArr(circle.numbers)]
		if already_seen {
			fmt.Println(i)
			break
		} else {
			seen[KeyifyArr(circle.numbers)] = true
		}
		circle.Move()
	}
	// Multiply the two values after Cup '1'
	return circle.MultiplyTwoCupsAfterCupOne()
}

func getNextsOfSizeN(n int, X []int) []int {
	nexts := make([]int, n+1)

	for i := 0; i < len(X); i++ {
		nexts[X[i]] = X[(i+1)%len(X)] // modulo just creates nice looping effect s.t. last element maps to 0 here, rest remain i + 1
	}
	nexts[X[len(X)-1]] = len(X) + 1
	for i := len(X) + 1; i <= n; i++ {
		nexts[i] = i + 1
	}
	nexts[len(nexts)-1] = X[0]
	return nexts
}

func Part2_NextInListApproach(input string) int {
	const MILLION = 1000000
	const TEN_MILLION = 10000000
	n := MILLION
	var starts []int

	// SETUP
	for _, s := range strings.Split(input, "") {
		val, _ := strconv.Atoi(s)
		starts = append(starts, val)
	}
	nexts := getNextsOfSizeN(n, starts)

	// 'Moves'!
	current := starts[0]
	for i := 0; i < TEN_MILLION; i++ {
		pickup := nexts[current]                     // value of first card to pick up
		nexts[current] = nexts[nexts[nexts[pickup]]] // 'next' pointer of current becomes the card 'next' from third pickup

		var dest int
		if current == 1 {
			dest = n
		} else {
			dest = current - 1 // last elem idx of nexts
		}
		for dest == pickup || dest == nexts[pickup] || dest == nexts[nexts[pickup]] {
			// while destination value is one of the three pick up cards
			if dest == 1 {
				dest = n
			} else {
				dest -= 1
			} // typical decrimenting to find destination value, else loops back to max (always n)
		}
		nexts[nexts[nexts[pickup]]] = nexts[dest] // i.e. the 'next' value after third pickup card is the next value after destination
		nexts[dest] = pickup                      // point destination to first pickup card
		current = nexts[current]                  // current moves on one
	}

	// Answer
	return nexts[1] * nexts[nexts[1]] // i.e. multiply value after 1 by the value after the value after 1 :)
}
