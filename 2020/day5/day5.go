package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	const EXAMPLE_FILENAME = "example.txt"
	const INPUT_FILENAME = "input.txt"
	fmt.Println("---Day 5---")
	example_ans := Part1(EXAMPLE_FILENAME)
	fmt.Printf("[Example P1] Expected: 820 Got: %v\n", example_ans)
	part1_ans := Part1(INPUT_FILENAME)
	fmt.Printf("[Part 1] Answer: %v\n", part1_ans)

	part2_ans := Part2(INPUT_FILENAME)
	fmt.Printf("[Part 2] Answer: %v\n", part2_ans)
}

func findSeatID(seat string) int {
	// Rows 0-127
	// Cols 0-7
	lower_row_bound := 0.0
	upper_row_bound := 127.0
	lower_col_bound := 0.0
	upper_col_bound := 7.0

	instructions := []rune(seat)
	for i := 0; i < 7; {
		if instructions[i] == 'F' {
			upper_row_bound = upper_row_bound - math.Floor((upper_row_bound-lower_row_bound)/2) - 1
		} else if instructions[i] == 'B' {
			lower_row_bound = lower_row_bound + math.Floor((upper_row_bound-lower_row_bound)/2) + 1
		} else {
			panic("UNRECOGNISED INSTRUCTION RECIEVED!")
		}
		i++
	}
	for i := 7; i < 10; {
		if instructions[i] == 'L' {
			upper_col_bound = upper_col_bound - math.Floor((upper_col_bound-lower_col_bound)/2) - 1
		} else if instructions[i] == 'R' {
			lower_col_bound = lower_col_bound + math.Floor((upper_col_bound-lower_col_bound)/2) + 1
		} else {
			panic("UNRECOGNISED INSTRUCTION RECIEVED!")
		}
		i++
	}
	if (lower_row_bound != upper_row_bound) && (lower_col_bound != upper_col_bound) {
		fmt.Printf("LRB=%v URB=%v LCB=%v UCB=%v\n", lower_row_bound, upper_row_bound, lower_col_bound, upper_col_bound)
		panic("DONT HAVE PRECISE RANGES!")
	}

	// apparently can just convert to binary here. like replace F with 0 R with 1, then convert to decimal to get the number?
	return int(lower_row_bound)*8 + int(lower_col_bound)
}

func Part1(filename string) int {
	file, _ := os.Open(filename)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	highest_seat_id := 0
	var current_seat_id int
	for scanner.Scan() {
		current_seat_id = findSeatID(scanner.Text())
		highest_seat_id = max(highest_seat_id, current_seat_id)
	}
	return highest_seat_id
}

func Part2(filename string) int {
	file, _ := os.Open(filename)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var current_seat_id int
	var all_seats []int
	for scanner.Scan() {
		current_seat_id = findSeatID(scanner.Text())
		all_seats = append(all_seats, current_seat_id)
	}
	sort.SliceStable(all_seats, func(i int, j int) bool { return all_seats[i] < all_seats[j] })
	// contains := func(s []int, e int) bool {
	// 	for _, a := range s {
	// 		if a == e {
	// 			return true
	// 		}
	// 	}
	// 	return false
	// }
	// for i := 28; i < 842; {
	// 	if contains(all_seats, i) {
	// 		i++
	// 		continue
	// 	} else {
	// 		return i
	// 	}
	// }
	for i := 28; i < 842; i++ {
		if all_seats[i+1]-all_seats[i] != 1 {
			return all_seats[i] + 1
		}
	}
	// both methods quick enough, did think about this ^ but didn't go with it for some reason first lol
	panic("Didn't find!")
}
