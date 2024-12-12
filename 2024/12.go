package main

import (
	"fmt"
)

func day12(part2 bool) Solution {
	example_filepath := GetExamplePath(12)
	input_filepath := GetInputPath(12)
	if !part2 {
		example_p1 := Part1_12(example_filepath)
		AssertExample("1930", example_p1, 1)
		return Solution{
			"12",
			example_p1,
			Part1_12(input_filepath),
		}
	} else {
		return Solution{
			"12",
			"example part 2",
			"input part 2",
		}
	}
}

var FourSideDirs = [4][2]int{{-1, 0}, {0, -1}, {1, 0}, {0, 1}}

func dfs_for_plant_region(grid *[][]string, row int, col int, R int, C int, plant string, seen map[[2]int]bool, region_coords *[][2]int, borders *int) {
	*region_coords = append((*region_coords), [2]int{row, col})
	for _, drdc := range FourSideDirs {
		dr, dc := drdc[0], drdc[1]
		nr := row + dr
		nc := col + dc
		if nr >= 0 && nr < R && nc >= 0 && nc < C {
			coord := [2]int{nr, nc}
			if (*grid)[nr][nc] != plant {
				*borders++
				continue
			}
			if _, exists := seen[coord]; exists {
				continue
			}
			seen[coord] = true
			dfs_for_plant_region(grid, nr, nc, R, C, plant, seen, region_coords, borders)
		} else {
			// out of grid is a border!
			*borders++
		}
	}
}

func Part1_12(filepath string) string {
	grid := readStringGrid(filepath)

	seen := make(map[[2]int]bool)
	R := len(grid)
	C := len(grid[0])
	var regions [][][2]int
	var borders []int

	for r, row := range grid {
		for c, plant := range row {
			if _, exists := seen[[2]int{r, c}]; exists {
				continue
			}
			seen[[2]int{r, c}] = true
			// new region
			points_in_region := [][2]int{}
			var n_borders int = 0
			var n_borders_ptr *int = &n_borders
			dfs_for_plant_region(&grid, r, c, R, C, plant, seen, &points_in_region, n_borders_ptr)
			regions = append(regions, points_in_region)
			borders = append(borders, n_borders)
		}
	}

	var total_cost int
	for i, region := range regions {
		area := len(region)
		perimeter := borders[i]
		total_cost += area * perimeter
	}

	return fmt.Sprintf("%d", total_cost)
}
