package main

import (
	"fmt"
	"math"
	"slices"
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
		example_p2 := Part2_12(example_filepath)
		AssertExample("1206", example_p2, 2)
		return Solution{
			"12",
			example_p2,
			Part2_12(input_filepath),
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

func FindRegionsAndNumBordersForEachPlant(grid [][]string) (regions [][][2]int, borders []int) {
	seen := make(map[[2]int]bool)
	R := len(grid)
	C := len(grid[0])

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
	return regions, borders
}

func Part1_12(filepath string) string {
	grid := readStringGrid(filepath)
	regions, borders := FindRegionsAndNumBordersForEachPlant(grid)

	var total_cost int
	for i, region := range regions {
		area := len(region)
		perimeter := borders[i]
		total_cost += area * perimeter
	}

	return fmt.Sprintf("%d", total_cost)
}

type SideOfPlantRegion struct {
	n_or_c_coords []int
	side_dir      int // 0 UP 1 RIGHT 2 DOWN 3 LEFT
}

func CalculateNumberOfSidesFromRegion(region [][2]int) int {
	side_rows := make(map[int][]*SideOfPlantRegion)
	side_cols := make(map[int][]*SideOfPlantRegion)

	for _, coord := range region {
		r, c := coord[0], coord[1]
		// cycle over UP RIGHT DOWN LEFT, corresponding to i = 0,1,2,3
		for i, drdc := range [][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} {
			dr, dc := drdc[0], drdc[1]
			nr := r + dr
			nc := c + dc

			var _map *map[int][]*SideOfPlantRegion
			var _nRorC *int // ENTRY POINT TO MAP nr for rows, nc for cols
			var _nCorR *int // POINT in 'side'
			if i%2 == 0 {
				// up down
				_map = &side_rows
				_nRorC = &nr
				_nCorR = &nc
			} else {
				// right left
				_map = &side_cols
				_nRorC = &nc
				_nCorR = &nr
			}
			coord := [2]int{nr, nc}
			if slices.Contains(region, coord) {
				// same plant neighbour! We ignore this.
				continue
			}
			// unfriendly plant neighbour
			if sides_in_row, exists := (*_map)[*_nRorC]; exists {
				// check sides_in_row
				adding := true
				for _, side := range sides_in_row {
					if i != side.side_dir {
						continue
					}
					min_c_or_r := 10000
					max_c_or_r := -10000
					for _, c_or_r_coord := range side.n_or_c_coords {
						min_c_or_r = min(min_c_or_r, c_or_r_coord)
						max_c_or_r = max(max_c_or_r, c_or_r_coord)
					}
					if *_nCorR <= max_c_or_r && *_nCorR >= min_c_or_r {
						adding = false
						break
					}
					if *_nCorR-min_c_or_r == 1 {
						adding = false
						// extend it forward
						side.n_or_c_coords = append(side.n_or_c_coords, *_nCorR)
						break
					} else if min_c_or_r-*_nCorR == 1 {
						adding = false
						// extend it back
						side.n_or_c_coords = append(side.n_or_c_coords, *_nCorR)
						break
					}
				}
				if adding {
					(*_map)[*_nRorC] = append((*_map)[*_nRorC], &SideOfPlantRegion{[]int{*_nCorR}, i})
				}
			} else {
				(*_map)[*_nRorC] = []*SideOfPlantRegion{{[]int{*_nCorR}, i}}
			}
		}
	}

	/*

		It looks like this is working except for the fact some sides that should be combined are not combined
		so maybe loop through again to do some combining... idk

	*/

	var total int
	for _, sides := range side_rows {
		var subtract int
		for i, _side := range sides {
			for _, __side := range sides[i+1:] {
				if _side.side_dir != __side.side_dir {
					continue
				}
			check_loop:
				for _, coord := range _side.n_or_c_coords {
					for _, coord2 := range __side.n_or_c_coords {
						if math.Abs(float64(coord)-float64(coord2)) == 1 {
							subtract++
							break check_loop
						}
					}
				}
			}
		}
		total += len(sides) - subtract
	}
	for _, sides := range side_cols {
		var subtract int
		for i, _side := range sides {
			for _, __side := range sides[i+1:] {
				if _side.side_dir != __side.side_dir {
					continue
				}
			check_loop2:
				for _, coord := range _side.n_or_c_coords {
					for _, coord2 := range __side.n_or_c_coords {
						if math.Abs(float64(coord)-float64(coord2)) == 1 {
							subtract++
							break check_loop2
						}
					}
				}
			}
		}
		total += len(sides) - subtract
	}
	return total
}

func Part2_12(filepath string) string {
	grid := readStringGrid(filepath)
	regions, _ := FindRegionsAndNumBordersForEachPlant(grid)

	var total_cost int

	for _, region := range regions {
		area := len(region)
		num_sides := CalculateNumberOfSidesFromRegion(region)
		total_cost += area * num_sides
	}

	return fmt.Sprintf("%d", total_cost)
}
