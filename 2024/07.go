package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

func day07(part2 bool) Solution {
	example_filepath := GetExamplePath(7)
	input_filepath := GetInputPath(7)
	if !part2 {
		example_p1 := Part1_07(example_filepath)
		AssertExample("3749", example_p1, 1)
		return Solution{
			"07",
			example_p1,
			Part1_07(input_filepath),
		}
	} else {
		example_p2 := Part2_07(example_filepath)
		AssertExample("11387", example_p2, 2)

		return Solution{
			"07",
			example_p2,
			Part2_07(input_filepath),
		}
	}
}

type Calibration struct {
	Result int
	Values []int
}

func readCalibrations(filepath string) (calibrations []Calibration) {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ": ")

		result_str := parts[0]
		result, err := strconv.Atoi(result_str)
		if err != nil {
			panic("trouble reading in result")
		}

		values_str := strings.Split(parts[1], " ")
		var values []int
		for _, val_str := range values_str {
			val, err := strconv.Atoi(val_str)
			if err != nil {
				panic("trouble reading in a value")
			}
			values = append(values, val)
		}
		calibrations = append(calibrations, Calibration{result, values})
	}
	return calibrations
}

type Item7 struct {
	total    int
	next_idx int
}

// Heap Queue based on example on: https://pkg.go.dev/container/heap
// An Item7Heap is a min-heap of Item7 structs.
type Item7Heap []Item7

func (h Item7Heap) Len() int           { return len(h) }
func (h Item7Heap) Less(i, j int) bool { return h[i].total < h[j].total }
func (h Item7Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *Item7Heap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(Item7))
}

func (h *Item7Heap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func doCalibrationCheck(calibration *Calibration) int {
	/*
		doing Kate's idea of working back from the answer... ofc!!!!
		should've considered other motivations for the left to right thing. ok
	*/

	// i.e. the first item in the heap queue is 'before' we start
	// traversing back through the calibration values
	// where our target is the result, and our next position is the final idx
	start_item := Item7{
		calibration.Result,
		len(calibration.Values) - 1,
	}

	item_heap := Item7Heap{start_item}
	heap.Init(&item_heap)
	for item_heap.Len() > 0 {
		item := heap.Pop(&item_heap).(Item7)
		if item.total == 0 {
			return calibration.Result
		}
		if item.next_idx < 0 {
			continue
		}
		val := calibration.Values[item.next_idx]

		current_value := item.total
		// 'multiply' op
		if current_value%val == 0 {
			heap.Push(&item_heap, Item7{current_value / val, item.next_idx - 1})
		}
		// 'addition' op
		heap.Push(&item_heap, Item7{current_value - val, item.next_idx - 1})
	}
	return 0
}

func doCalibrationCheckP2(calibration *Calibration) int {
	// similar to P1 but let's start from the front to avoid issues
	start_item := Item7{
		calibration.Result - calibration.Values[0],
		1,
	}

	item_heap := Item7Heap{start_item}
	heap.Init(&item_heap)

	for item_heap.Len() > 0 {
		item := heap.Pop(&item_heap).(Item7)

		// Sense Checks
		if item.total == 0 {
			panic("should've early returned no???")
		}
		if item.total > calibration.Result || item.total < 0 {
			// overshot or have negative value I guess
			panic("why invalid boys here")
		}
		if item.next_idx >= len(calibration.Values) {
			// exhausted the values, didn't reach answer, so skip
			panic("shouldn't be over last idx")
		}

		val := calibration.Values[item.next_idx]
		current_value := calibration.Result - item.total
		var diff int

		for _, op := range []string{"*", "+", "||"} {
			var updated_val int
			switch op {
			case "*":
				if current_value == 0 {
					panic("cant multiply by 0")
				}
				updated_val = current_value * val
			case "+":
				updated_val = current_value + val
			case "||":
				concatenated_str := strconv.Itoa(current_value) + strconv.Itoa(val) // "12"+"45"->"1245"
				updated_val, _ = strconv.Atoi(concatenated_str)                     // "1245" -> 1245
			default:
				panic("impossible op")
			}

			diff = calibration.Result - updated_val
			if diff == 0 && item.next_idx+1 == len(calibration.Values) {
				return calibration.Result
			}
			if diff == 0 && item.next_idx+1 < len(calibration.Values) {
				// ffs I didn't allow for the multiplying by 1 did I ....
				valid := true
				for _, val := range calibration.Values[item.next_idx+1:] {
					if val != 1 {
						valid = false
						break
					}
				}
				if valid {
					return calibration.Result
				}
			}
			if diff > 0 && (item.next_idx+1) < len(calibration.Values) {
				// haven't overshot & still in bounds
				heap.Push(&item_heap, Item7{diff, item.next_idx + 1})
			}
		}
	}
	return 0
}

func Solve(calibrations []Calibration, part2 bool) string {
	var total int
	var wg sync.WaitGroup
	results := make(chan int)

	for _, calibration := range calibrations {
		wg.Add(1)
		go func(cal *Calibration) {
			defer wg.Done()
			var res int
			if part2 {
				res = doCalibrationCheckP2(cal)
			} else {
				res = doCalibrationCheck(cal)
			}
			results <- res
		}(&calibration)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		total += result
	}

	return fmt.Sprintf("%d", total)
}
func Part1_07(filepath string) string {
	calibrations := readCalibrations(filepath)
	return Solve(calibrations, false)
}

func Part2_07(filepath string) string {
	calibrations := readCalibrations(filepath)
	return Solve(calibrations, true)
}
