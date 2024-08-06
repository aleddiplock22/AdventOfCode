package main

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
)

/*
Part 1 - Used goroutines as learning/pracitce

Part2 - Cool algorithm using prime number properties,
It apparently is an implementation of Chinese Remainder Theorem ?
Inspired by: https://www.youtube.com/watch?v=4_5mluiXF5I
*/

func main() {
	const EXAMPLE_FILEPATH = "example.txt"
	const INPUT_FILEPATH = "input.txt"

	fmt.Println("---Day 13---")

	example_p1 := Part1(EXAMPLE_FILEPATH)
	part1_ans := Part1(INPUT_FILEPATH)

	fmt.Printf("[Example P1] Expected: 295, Answer: %v\n", example_p1)
	fmt.Printf("[Part 1] Answer: %v\n", part1_ans)

	easy_example_to_brute_force := []any{17, "x", 13, 19}
	easy_example_answer := BruteForcePart2(easy_example_to_brute_force)
	if easy_example_answer != 3417 {
		panic("Even the brute force didn't work!!")
	} else {
		fmt.Println("Brute force worked on easy example!")
	}

	example_p2 := Part2(EXAMPLE_FILEPATH)
	part2_ans := Part2(INPUT_FILEPATH)

	fmt.Printf("[Example P2] Expected: 1068781, Answer: %v\n", example_p2)
	fmt.Printf("[Part 2] Answer: %v\n", part2_ans)
}

func getFileContent(filepath string) string {
	file, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("ERROR READING FILE!")
	}
	return string(file)
}

func parseInput(filepath string) (leave_time int, bus_ids []int) {
	file_content := getFileContent(filepath)
	parts := strings.Split(file_content, "\r\n")
	leave_string, buses_string := parts[0], parts[1]
	leave_time, err := strconv.Atoi(leave_string)
	if err != nil {
		panic("Couldn't parse leave time as integer!")
	}
	bus_id_strings := strings.Split(buses_string, ",")
	for _, bus_id := range bus_id_strings {
		bus_id_numeric, err := strconv.Atoi(bus_id)
		if err == nil {
			bus_ids = append(bus_ids, bus_id_numeric)
		}
	}

	return leave_time, bus_ids
}

func parseInputPart2(filepath string) (bus_ids []any) {
	buses_string := strings.Split(getFileContent(filepath), "\r\n")[1]
	bus_id_strings := strings.Split(buses_string, ",")
	for _, bus_id := range bus_id_strings {
		bus_id_numeric, err := strconv.Atoi(bus_id)
		if err == nil {
			bus_ids = append(bus_ids, int64(bus_id_numeric))
		} else {
			bus_ids = append(bus_ids, "x")
		}
	}

	return bus_ids
}

type wait_info struct {
	waiting_time int
	bus_id       int
}

func findMinimumViableOption(target int, bus_id int, ch chan wait_info, wg *sync.WaitGroup) {
	i := 0
	var tmp_check int
	// Little inefficient instead of doing bus_id - (target % bus_id) to experiment with goroutines :)
	for {
		i++
		tmp_check = bus_id * i
		if tmp_check >= target {
			defer wg.Done() // let the wait group know when function finishes (this is what defer does) that we've 'Done' another (not the entire group is done!)
			ch <- wait_info{
				tmp_check - target,
				bus_id,
			} // send our answer to the channel
			return
		}
	}
}

func Part1(filepath string) int {
	// get inputs
	leave_time, bus_ids := parseInput(filepath)

	// set up [pointer to] wait group & [buffered] channel for async goroutine stuff
	wait_times_wait_group := &sync.WaitGroup{}
	wait_times_channel := make(chan wait_info, len(bus_ids)) // buffer the size of max num we'll send

	for _, bus_id := range bus_ids {
		wait_times_wait_group.Add(1) // say we have to wait for additional wg.Done() before we can pass 'Wait'
		go findMinimumViableOption(leave_time, bus_id, wait_times_channel, wait_times_wait_group)
	}
	wait_times_wait_group.Wait()
	close(wait_times_channel) // done waiting so can close channel, not going to be sending anything else to it

	lowest_waiting_time := wait_info{
		math.MaxInt,
		0,
	}
	for waiting_time_info := range wait_times_channel { // loop over channel to find our best answer
		if waiting_time_info.waiting_time < lowest_waiting_time.waiting_time {
			lowest_waiting_time = waiting_time_info
		}
	}
	return lowest_waiting_time.waiting_time * lowest_waiting_time.bus_id
}

// Technically works! but lets not...
func BruteForcePart2(bus_ids []any) int {
	lowest_bus_id := math.MaxInt
	lowest_bus_idx := math.MaxInt
	for idx, bus_id := range bus_ids {
		if bus_id == "x" {
			continue
		}
		bus_id, is_int := bus_id.(int)
		if !is_int {
			panic("Expected bus_id to be an int here!")
		}
		if bus_id < lowest_bus_id {
			lowest_bus_id = bus_id
			lowest_bus_idx = idx
		}
	}

	var current_time int
	var position_diff int
	var happy bool
	i := 0
	for {
		i++
		happy = true
		current_time = lowest_bus_id * i
		for idx, bus_id := range bus_ids {
			if bus_id == "x" || bus_id == lowest_bus_id {
				continue
			}
			bus_id, is_int := bus_id.(int)
			if !is_int {
				panic("Expected bus_id to be an int here!")
			}
			position_diff = idx - lowest_bus_idx
			if !((current_time+position_diff)%bus_id == 0) {
				happy = false
				break
			}
		}
		if happy {
			first_bus_pos_diff := 0 - lowest_bus_idx
			return (current_time + first_bus_pos_diff)
		}
	}
}

func Part2(filepath string) int {
	bus_ids := parseInputPart2(filepath)
	for _, bus_id := range bus_ids {
		if bus_id == "x" {
			continue
		}
		bus_id, ok := bus_id.(int64)
		if !ok {
			continue
		}
		if !big.NewInt(bus_id).ProbablyPrime(0) {
			panic("Found a bus_id that was not prime! Cannot implement solution.")
		}
	}
	// Have now ensured our inputs are all prime! So LCM(a,b) = a*b

	initial_bus := bus_ids[0]                // assumining not 'x' !
	initial_bus_id, _ := initial_bus.(int64) // so compiler knows it's an int
	var step_size int64 = initial_bus_id
	var current_time int64 = 0
	/*
		Idea: Find the periodicity between elements that satisfy our needs one at a time
		So, if we have [3 7 11 5]
		firstly we take 3, so step_size=3
		and look at what's next, 7 [at index diff=1]
		start time at t=0
		t += step_size
		// t = 3
			> is (3 + 1) % 7 == 0 ?
			>no,
		t += step_size
		// t = 6
			> is (6 + 1) % 7 == 0 ?
			>yes! So we now combine the stepsize by doing
			step_size *= bus_id
			// step size = 21

		now we have a time where 3 & 7 are essentially 'in sync'
		by going up in stepsizes of lcm(3,7)=3*7, we will STAY in sync
		this greatly speeds up the search process

		so then continue for 11 [index 2] and so on...
	*/
	for idx, bus_id := range bus_ids {
		if bus_id == "x" || idx == 0 {
			continue
		}
		bus_id, ok := bus_id.(int64)
		if !ok {
			fmt.Println(idx, bus_id)
			panic("Expected integer bus_id!")
		}
		for {
			current_time += step_size
			if (current_time+int64(idx))%bus_id == 0 {
				step_size *= bus_id
				break
			}
		}
	}
	return int(current_time)
}
