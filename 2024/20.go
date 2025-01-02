package main

import (
	"fmt"
	"math"
)

func day20(part2 bool) Solution {
	example_filepath := GetExamplePath(20)
	input_filepath := GetInputPath(20)
	if !part2 {
		example_p1 := Part1_20(example_filepath)
		return Solution{
			"20",
			example_p1,
			Part1_20(input_filepath), // could and should optimise p1 using p2 method but 2 instead of 20, but alas
		}
	} else {
		example_p2 := Part2_20(example_filepath, 50)
		AssertExample("285", example_p2, 2)
		return Solution{
			"20",
			example_p2,
			Part2_20(input_filepath, 100),
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

func Part2_20(filepath string, save_goal int) string {
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
	shortest_without_cheating, path := FindShortestPathInGridDefault(sr, sc, er, ec, &grid, [2]int{})

	distOnPathToEndMap := make(map[[2]int]int)
	for steps, point := range path {
		distOnPathToEndMap[point] = shortest_without_cheating - steps
	}

	const CHEAT_RANGE = 20
	pointsInRange := func(r1, c1 int) [][3]int {
		var points_in_range [][3]int
		for i := -CHEAT_RANGE; i <= CHEAT_RANGE; i++ {
			for j := -CHEAT_RANGE; j <= CHEAT_RANGE; j++ {
				mDist := math.Abs(float64(i)) + math.Abs(float64(j))
				if mDist <= float64(CHEAT_RANGE) {
					points_in_range = append(points_in_range, [3]int{r1 + i, c1 + j, int(mDist)})
				}
			}
		}
		return points_in_range
	}

	// for each bit of path, find all the points we can reach with cheats,
	// if that's a point on the path, we can add current steps to steps to end

	path_length := 0
	var total_saved_by_at_least_goal int
	for _, point := range path {
		sr, sc := point[0], point[1]
		points_in_range := pointsInRange(sr, sc)
		for _, point := range points_in_range {
			_point, cheat_dist := [2]int{point[0], point[1]}, point[2]
			if dist_to_end, exists := distOnPathToEndMap[_point]; exists {
				diff := shortest_without_cheating - (path_length + cheat_dist + dist_to_end)
				if diff >= save_goal {
					total_saved_by_at_least_goal++
				}
			}
		}
		path_length++
	}

	return fmt.Sprintf("%d", total_saved_by_at_least_goal)
}
