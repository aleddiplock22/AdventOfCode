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
		return Solution{
			"07",
			"example part 2",
			"input part 2",
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

		current_target := item.total
		// 'multiply' op
		if current_target%val == 0 {
			heap.Push(&item_heap, Item7{current_target / val, item.next_idx - 1})
		}
		// 'addition' op
		heap.Push(&item_heap, Item7{current_target - val, item.next_idx - 1})
	}
	return 0
}

func Part1_07(filepath string) string {
	calibrations := readCalibrations(filepath)
	var total int

	// operators_map := make(map[[2]int][]string)
	// a map which does [2]int{num_multiply, num_add} -> []string{...the produced operators...}
	// so we avoid repeating work
	// operator_permutations_map := make(map[[2]int][][]string)
	// similar concept to above

	// god this was slow, lets try go routines ffs
	var wg sync.WaitGroup
	results := make(chan int)

	for _, calibration := range calibrations {
		wg.Add(1)
		go func(cal *Calibration) {
			defer wg.Done()
			res := doCalibrationCheck(cal)
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
