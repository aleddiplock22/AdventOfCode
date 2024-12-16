package main

import (
	"container/heap"
	"fmt"
)

func day16(part2 bool) Solution {
	example_filepath := GetExamplePath(16)
	input_filepath := GetInputPath(16)
	if !part2 {
		example_p1 := Part1_16(example_filepath)
		AssertExample("7036", example_p1, 1)
		AssertExample("11048", Part1_16("./inputs/16/example2.txt"), 1)
		return Solution{
			"16",
			example_p1,
			Part1_16(input_filepath),
		}
	} else {
		return Solution{
			"16",
			"example part 2",
			"input part 2",
		}
	}
}

type Position16 struct {
	y    int
	x    int
	dir  int
	cost int
}

// Heap Queue based on example on: https://pkg.go.dev/container/heap
// An Poisiton16Heap is a min-heap of Item7 structs.
type Position16Heap []Position16

func (h Position16Heap) Len() int           { return len(h) }
func (h Position16Heap) Less(i, j int) bool { return h[i].cost < h[j].cost }
func (h Position16Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *Position16Heap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(Position16))
}

func (h *Position16Heap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func Part1_16(filepath string) string {
	grid := readStringGrid(filepath)
	var sy, sx int
outer:
	for y, row := range grid {
		for x, char := range row {
			if char == "S" {
				sy, sx = y, x
				break outer
			}
		}
	}
	if sy == 0 || sx == 0 {
		panic("couldn't find start or end")
	}

	// STARTS FACING LEFT
	// following my classic 0 UP 1 RIGHT 2 DOWN 3 LEFT approach
	start := Position16{sy, sx, 3, 0} // y x dir cost
	queue := Position16Heap{start}
	heap.Init(&queue)
	var best_cost int
	seen := make(map[[3]int]bool) // y x dir
	for queue.Len() > 0 {
		item := heap.Pop(&queue).(Position16)
		if grid[item.y][item.x] == "E" {
			best_cost = item.cost
			break
		}
	dir_loop:
		for i := range 4 {
			// 0 1 2 3
			var ny, nx int
			var cost int
			switch i {
			case 0:
				// up
				ny, nx = item.y-1, item.x
				switch item.dir {
				case 0:
					cost = 1
				case 1:
					cost = 1001
				case 2:
					continue dir_loop
				case 3:
					cost = 1001
				}
			case 1:
				// right
				ny, nx = item.y, item.x+1
				switch item.dir {
				case 0:
					cost = 1001
				case 1:
					cost = 1
				case 2:
					cost = 1001
				case 3:
					continue dir_loop
				}
			case 2:
				// down
				ny, nx = item.y+1, item.x
				switch item.dir {
				case 0:
					continue dir_loop
				case 1:
					cost = 1001
				case 2:
					cost = 1
				case 3:
					cost = 1001
				}
			case 3:
				// left
				ny, nx = item.y, item.x-1
				switch item.dir {
				case 0:
					cost = 1001
				case 1:
					continue dir_loop
				case 2:
					cost = 1001
				case 3:
					cost = 1
				}
			default:
				panic("bad dir")
			}
			if ny < 0 || ny >= len(grid) || nx < 0 || nx >= len(grid[0]) || grid[ny][nx] == "#" {
				continue
			}
			_loc := [3]int{ny, nx, i}
			if _, have_seen := seen[_loc]; have_seen {
				continue
			}
			seen[_loc] = true
			queue = append(queue, Position16{ny, nx, i, item.cost + cost})
		}
	}
	return fmt.Sprintf("%d", best_cost)
}
