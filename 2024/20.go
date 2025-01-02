package main

import (
	"fmt"
)

func day20(part2 bool) Solution {
	example_filepath := GetExamplePath(20)
	// input_filepath := GetInputPath(20)
	if !part2 {
		example_p1 := Part1_20(example_filepath)
		return Solution{
			"20",
			example_p1,
			"", //Part1_20(input_filepath),
		}
	} else {
		fmt.Println("TRYING EXAMPLE P2")
		example_p2 := Part2_20(example_filepath)
		return Solution{
			"20",
			example_p2,
			"", //Part2_20(input_filepath),
		}
	}
}

func FindShortestPathInGridWithCheats(sr, sc, er, ec int, grid *[][]string, cheat_pos [2]int) int {
	queue := [][3]int{{sr, sc, 0}}
	seen := map[[2]int]bool{{sr, sc}: true}
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		r, c, steps := pos[0], pos[1], pos[2]
		for _, drdc := range FourSideDirs {
			dr, dc := drdc[0], drdc[1]
			nr := r + dr
			nc := c + dc
			_loc := [2]int{nr, nc}
			if nr < 0 || nr >= len(*grid) || nc < 0 || nc >= len((*grid)[0]) {
				continue
			}
			if (*grid)[nr][nc] == "#" && _loc != cheat_pos {
				continue
			}
			if nr == er && nc == ec {
				return steps + 1
			}
			if _, exists := seen[_loc]; exists {
				continue
			}
			seen[_loc] = true
			queue = append(queue, [3]int{nr, nc, steps + 1})
		}
	}
	return -1
}

func Part1_20(filepath string) string {
	grid := readStringGrid(filepath)
	var sr, sc, er, ec int
	for r, row := range grid {
		for c, char := range row {
			if char == "S" {
				sr, sc = r, c
			} else if char == "E" {
				er, ec = r, c
			}
		}
	}
	shortest_without_cheating := FindShortestPathInGridWithCheats(sr, sc, er, ec, &grid, [2]int{})

	var total int
	for r, row := range grid {
		for c := range row {
			steps := FindShortestPathInGridWithCheats(sr, sc, er, ec, &grid, [2]int{c, r})
			if shortest_without_cheating-steps >= 100 {
				total++
			}
		}
	}

	return fmt.Sprintf("%d", total)
}

type Point struct {
	r, c, steps int
	path        [][2]int
}

func FindShortestPathInGridDefault(sr, sc, er, ec int, grid *[][]string, cheat_pos [2]int) (int, [][2]int) {
	queue := []Point{{sr, sc, 0, [][2]int{{sr, sc}}}}
	seen := map[[2]int]bool{{sr, sc}: true}
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		r, c, steps, path := pos.r, pos.c, pos.steps, pos.path
		for _, drdc := range FourSideDirs {
			dr, dc := drdc[0], drdc[1]
			nr := r + dr
			nc := c + dc
			_loc := [2]int{nr, nc}
			var npath [][2]int
			npath = append(npath, path...)
			npath = append(npath, _loc)
			if nr < 0 || nr >= len(*grid) || nc < 0 || nc >= len((*grid)[0]) {
				continue
			}
			if (*grid)[nr][nc] == "#" && _loc != cheat_pos {
				continue
			}
			if nr == er && nc == ec {
				return steps + 1, npath
			}
			if _, exists := seen[_loc]; exists {
				continue
			}
			seen[_loc] = true
			queue = append(queue, Point{nr, nc, steps + 1, npath})
		}
	}
	return -1, [][2]int{}
}

func FindShortestPathsInGridWithLongCheats(sr, sc, er, ec int, grid *[][]string, cheat_start [2]int, target int) (results [][3]int) {
	queue := [][6]int{{sr, sc, 0, 20, -1, -1}} // start_r, start_c, steps, cheat_steps_remaining, cheat_end_r, cheat_end_c
	seen := map[[2]int]bool{{sr, sc}: true}
	seen2 := make(map[[2]int]bool)
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		r, c, steps, cheat_duration, cheat_end_r, cheat_end_c := pos[0], pos[1], pos[2], pos[3], pos[4], pos[5]
		for _, drdc := range FourSideDirs {
			dr, dc := drdc[0], drdc[1]
			nr := r + dr
			nc := c + dc
			_loc := [2]int{nr, nc}
			if nr < 0 || nr >= len(*grid) || nc < 0 || nc >= len((*grid)[0]) {
				continue
			}
			// if _, exists := seen[_loc]; exists {
			// 	continue
			// }

			n_cheat_end_r := cheat_end_r
			n_cheat_end_c := cheat_end_c
			n_cheat_duration := cheat_duration

			cheat_just_started := false
			if _loc == cheat_start {
				cheat_just_started = true
				n_cheat_duration--
			}

			if n_cheat_duration > 0 && n_cheat_duration < 20 {
				n_cheat_end_r, n_cheat_end_c = nr, nc
			}

			if (*grid)[nr][nc] == "#" {
				if !(cheat_duration > 0 && cheat_duration < 20) {
					continue
				}
			}
			if !cheat_just_started && cheat_duration < 20 {
				n_cheat_duration--
			}
			if nr == er && nc == ec {
				// fmt.Printf("END IN %v steps WITH %v cheat_duration left! (cheat_start=%v)\n", steps+1, n_cheat_duration, cheat_start)
				if _, exists := seen2[[2]int{n_cheat_end_r, n_cheat_end_c}]; exists {
					continue
				}
				seen2[[2]int{n_cheat_end_r, n_cheat_end_c}] = true
				if steps+1 <= target {
					fmt.Println("Found result:", steps+1, n_cheat_end_r, n_cheat_end_c)
					results = append(results, [3]int{steps + 1, n_cheat_end_r, n_cheat_end_c})
				} else {
					fmt.Println("returning")
					return results
				}
			} else {
				seen[_loc] = true
				queue = append(queue, [6]int{nr, nc, steps + 1, n_cheat_duration, n_cheat_end_r, n_cheat_end_c})
			}
		}
	}
	return results
}

func Part2_20(filepath string) string {
	grid := readStringGrid(filepath)
	var sr, sc, er, ec int
	for r, row := range grid {
		for c, char := range row {
			if char == "S" {
				sr, sc = r, c
			} else if char == "E" {
				er, ec = r, c
			}
		}
	}
	shortest_without_cheating, _ := FindShortestPathInGridDefault(sr, sc, er, ec, &grid, [2]int{})
	// got the path if we want it
	fmt.Println("SHORTEST PATH oRIGINAL", shortest_without_cheating)

	seen := make(map[[4]int]int)

	const TARGET = 34 // 84 - 50

	for r := range len(grid) - 1 {
		if r == 0 {
			continue
		}
		for c := range len(grid) - 1 {
			if c == 0 {
				continue
			}
			results := FindShortestPathsInGridWithLongCheats(sr, sc, er, ec, &grid, [2]int{c, r}, TARGET)
			for _, res := range results {
				_steps, _er, _ec := res[0], res[1], res[2]
				_val := [4]int{r, c, _er, _ec}
				if _, exists := seen[_val]; exists {
					continue
				}
				diff := shortest_without_cheating - _steps
				if diff >= 50 {
					fmt.Printf("(%v, %v) | (%v, %v) -> %v\n", r, c, _er, _ec, diff)
				}
				seen[_val] = diff
			}
		}
	}

	var total int
	for _, value := range seen {
		if value > 50 {
			total++
		}
	}

	return fmt.Sprintf("%d", total)
}
