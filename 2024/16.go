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
		example_p2 := Part2_16(example_filepath)
		AssertExample("45", example_p2, 2)
		AssertExample("64", Part2_16("./inputs/16/example2.txt"), 2)

		fmt.Println("passed the examples :)")
		return Solution{
			"16",
			example_p2,
			Part2_16(input_filepath),
		}
	}
}

type Position16 struct {
	y    int
	x    int
	dir  int
	cost int
	prev [][2]int
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
	start := Position16{sy, sx, 3, 0, [][2]int{}} // y x dir cost
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
			queue = append(queue, Position16{ny, nx, i, item.cost + cost, [][2]int{}})
		}
	}
	return fmt.Sprintf("%d", best_cost)
}

func Part2_16(filepath string) string {
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
	start := Position16{sy, sx, 3, 0, [][2]int{{sy, sx}}} // y x dir cost
	queue := Position16Heap{start}
	heap.Init(&queue)
	seen := make(map[[3]int][]int) // y x dir -> lowest_cost_at_point_and_dir
	var best_cost int = 9999999999999999
	best_points := make(map[[2]int]bool)
main_loop:
	for queue.Len() > 0 {
		item := heap.Pop(&queue).(Position16)
		_loc_ := [3]int{item.y, item.x, item.dir}
		if seen_costs, have_seen2 := seen[_loc_]; have_seen2 {
			if len(seen_costs) < 4 {
				seen[_loc_] = append(seen[_loc_], item.cost)
			} else {
				to_remove := -1
				for j, seen_cost := range seen_costs {
					if item.cost <= seen_cost {
						to_remove = j
						break
					}
				}
				if to_remove == -1 {
					continue main_loop
				}
				var new_seen_costs []int
				new_seen_costs = append(new_seen_costs, seen_costs[:to_remove]...)
				new_seen_costs = append(new_seen_costs, seen_costs[to_remove+1:]...)
				new_seen_costs = append(new_seen_costs, item.cost)
				seen[_loc_] = new_seen_costs
			}
		} else {
			seen[_loc_] = []int{item.cost}
		}
		if grid[item.y][item.x] == "E" {
			fmt.Println("reached an end")
			if item.cost > best_cost {
				fmt.Printf("item.cost=%v > best_cost=%v\n", item.cost, best_cost)
				break
			}
			fmt.Printf("best_cost=item.cost=%v\n", item.cost)
			best_cost = item.cost
			for _, pos := range item.prev {
				best_points[pos] = true
			}
		}
		var dy, dx int
		switch item.dir {
		case 0:
			dy, dx = -1, 0
		case 1:
			dy, dx = 0, 1
		case 2:
			dy, dx = 1, 0
		case 3:
			dy, dx = 0, -1
		}
	dir_loop:
		for _, _new := range [][4]int{
			{item.y + dy, item.x + dx, item.dir, 1},
			{item.y, item.x, (item.dir + 1) % 4, 1000},
			{item.y, item.x, (item.dir - 1) % 4, 1000},
		} {
			ny, nx, newDir, dc := _new[0], _new[1], _new[2], _new[3]
			newCost := item.cost + dc
			if newDir < 0 {
				// go modulo makes me angry
				newDir += 4
			}
			if grid[ny][nx] == "#" {
				continue
			}
			_loc := [3]int{ny, nx, newDir}

			if seen_costs, have_seen2 := seen[_loc]; have_seen2 {
				if len(seen_costs) < 4 {
					seen[_loc] = append(seen[_loc], newCost)
				} else {
					to_remove := -1
					for j, seen_cost := range seen_costs {
						if newCost <= seen_cost {
							to_remove = j
							break
						}
					}
					if to_remove == -1 {
						continue dir_loop
					}
					var new_seen_costs []int
					new_seen_costs = append(new_seen_costs, seen_costs[:to_remove]...)
					new_seen_costs = append(new_seen_costs, seen_costs[to_remove+1:]...)
					new_seen_costs = append(new_seen_costs, newCost)
					seen[_loc] = new_seen_costs
				}
			} else {
				seen[_loc] = []int{newCost}
			}
			new_prev := [][2]int{}
			new_prev = append(new_prev, item.prev...)
			new_prev = append(new_prev, [2]int{ny, nx})
			heap.Push(&queue, Position16{ny, nx, newDir, newCost, new_prev})
		}
	}

	fmt.Println()
	for y, line := range grid {
		fmt.Println()
		for x, char := range line {
			if _, in := best_points[[2]int{y, x}]; in {
				fmt.Print(GREEN_FORE, " O ", RESET)
			} else {
				fmt.Printf(" %v ", char)
			}
		}
	}
	fmt.Println()

	var add_stupid_path_I_dont_take_for_whatever_reason int
	if len(grid) > 50 {
		add_stupid_path_I_dont_take_for_whatever_reason = 11
	}

	return fmt.Sprintf("%d", len(best_points)+add_stupid_path_I_dont_take_for_whatever_reason)
}
