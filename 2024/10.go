package main

import (
	"fmt"
	"sync"
)

func day10(part2 bool) Solution {
	example_filepath := GetExamplePath(10)
	input_filepath := GetInputPath(10)
	if !part2 {
		example_p1 := Part1_10(example_filepath)
		AssertExample("36", example_p1, 1)
		return Solution{
			"10",
			example_p1,
			Part1_10(input_filepath),
		}
	} else {
		example_p2 := Part2_10(example_filepath)
		AssertExample("81", example_p2, 2)
		return Solution{
			"10",
			example_p2,
			Part2_10(input_filepath),
		}
	}
}

func calculateTrailheadScore(grid *[][]int, sr, sc int, part2 bool) (score int) {
	R := len(*grid)
	C := len((*grid)[0])

	// point = [pos_r, pos_x, len_path]
	seen := make(map[[3]int]bool)
	queue := [][3]int{{sr, sc, 0}}
	for len(queue) > 0 {
		// pop
		tile := queue[0]
		queue = queue[1:]

		r, c, l := tile[0], tile[1], tile[2]
		height := (*grid)[r][c]
		// reached an end
		if height == 9 {
			score++
			continue
		}

		// find next step
		nl := l + 1
		for _, drdc := range [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
			dr, dc := drdc[0], drdc[1]
			nr := r + dr
			nc := c + dc
			if nr >= 0 && nr < R && nc >= 0 && nc < C {
				new_height := (*grid)[nr][nc]
				if new_height-height != 1 {
					continue
				}
				new_point := [3]int{nr, nc, nl}
				if !part2 {
					if already_travelled := seen[new_point]; already_travelled {
						continue
					}
					seen[new_point] = true
				}
				queue = append(queue, new_point)
			}
		}

	}
	return score
}

func calculateTrailheadScoreRecursively(grid *[][]int, sr, sc int) (score int) {
	// this is nice I guess. I can't think how to cleanly do p1 with recursion though
	// would have to keep some set of the '9's we reached, i.e. their coords, and count them
	// but tricky to combine a p1 and p2 function in that case. I'll leave it.
	R := len(*grid)
	C := len((*grid)[0])

	height := (*grid)[sr][sc]
	// reached an end
	if height == 9 {
		return 1
	}

	// find next step
	for _, drdc := range [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
		dr, dc := drdc[0], drdc[1]
		nr := sr + dr
		nc := sc + dc
		if nr >= 0 && nr < R && nc >= 0 && nc < C {
			new_height := (*grid)[nr][nc]
			if new_height-height != 1 {
				continue
			}
			score += calculateTrailheadScoreRecursively(grid, nr, nc)
		}
	}
	return score
}

func SolvePathCountingTask(filepath string, part2 bool, use_recursive_solution bool) string {
	grid := readIntGrid(filepath) // grid[r][c]

	var wg sync.WaitGroup
	scoreChan := make(chan int)
	for r, row := range grid {
		for c, value := range row {
			if value == 0 {
				// start path here!
				wg.Add(1)
				go func(sr, sc int) {
					defer wg.Done()
					if use_recursive_solution {
						if part2 {
							scoreChan <- calculateTrailheadScoreRecursively(&grid, sr, sc)
						}
					} else {
						scoreChan <- calculateTrailheadScore(&grid, sr, sc, part2)
					}
				}(r, c)
			}
		}
	}

	go func() {
		wg.Wait()
		close(scoreChan)
	}()

	var total int
	for score := range scoreChan {
		total += score
	}

	return fmt.Sprintf("%d", total)
}

func Part1_10(filepath string) string {
	return SolvePathCountingTask(filepath, false, false)
}

func Part2_10(filepath string) string {
	return SolvePathCountingTask(filepath, true, true)
}
