package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	const EXAMPLE_FILEPATH = "example.txt"
	const INPUT_FILEPATH = "input.txt"

	fmt.Println("--- AOC 2022 | Day 12 ---")

	fmt.Printf("[Example P1] Expected: %v, Answer: %v\n", 31, Part1(EXAMPLE_FILEPATH))
	fmt.Printf("[Part 1] Answer: %v\n", Part1(INPUT_FILEPATH))

	fmt.Printf("[Example P2] Expected: %v, Answer: %v\n", 29, Part2(EXAMPLE_FILEPATH))
	fmt.Printf("[Part 2] Answer: %v\n", Part2(INPUT_FILEPATH))
}

type Path struct {
	current_loc Coord
	n_steps     int
	history     []Coord
}

type Queue struct {
	paths []Path
}

func (Q *Queue) Sort() {
	sort.Slice(Q.paths, func(i, j int) bool {
		return Q.paths[i].n_steps < Q.paths[j].n_steps
	})
}

func (Q *Queue) Pop() Path {
	Q.Sort()
	best_path := Q.paths[0]
	if len(Q.paths) == 1 {
		Q.paths = []Path{}
	} else {
		Q.paths = Q.paths[1:]
	}
	return best_path
}

type Coord struct {
	y int
	x int
}

type Grid struct {
	grid         [][]rune
	starting_pos Coord
	ending_pos   Coord
}

func readInput(filepath string) Grid {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var grid Grid
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		chars := []rune(line)
		grid.grid = append(grid.grid, chars)
		for x, char := range chars {
			if char == 'S' {
				grid.starting_pos = Coord{y, x}
				grid.grid[y][x] = 'a'
			}
			if char == 'E' {
				grid.ending_pos = Coord{y, x}
				grid.grid[y][x] = 'z'
			}
		}
		y++
	}

	return grid
}

type Breadcrumb struct {
	position  Coord
	num_steps int
}

func ContainsBetter(crumbs []Breadcrumb, candidate Breadcrumb) bool {
	for _, crumb := range crumbs {
		if crumb.position.y == candidate.position.y && crumb.position.x == candidate.position.x && crumb.num_steps <= candidate.num_steps {
			return true
		}
	}
	return false
}

func FindShortestPathLength(grid Grid) int {
	queue := Queue{
		[]Path{
			{grid.starting_pos, 0, []Coord{grid.starting_pos}},
		},
	}

	seen := []Breadcrumb{
		{grid.starting_pos, 0},
	}
	for len(queue.paths) > 0 {
		current_path := queue.Pop()

		current_val := grid.grid[current_path.current_loc.y][current_path.current_loc.x]

		if current_path.current_loc == grid.ending_pos {
			return current_path.n_steps
		}

		for _, next_pos := range []Coord{
			{current_path.current_loc.y - 1, current_path.current_loc.x},
			{current_path.current_loc.y + 1, current_path.current_loc.x},
			{current_path.current_loc.y, current_path.current_loc.x - 1},
			{current_path.current_loc.y, current_path.current_loc.x + 1},
		} {
			if next_pos.x < 0 || next_pos.y < 0 || next_pos.x >= len(grid.grid[0]) || next_pos.y >= len(grid.grid) {
				continue
			} else {
				// in bounds next pos
				next_val := grid.grid[next_pos.y][next_pos.x]
				if next_val-current_val == 1 || current_val >= next_val {
					if ContainsBetter(seen, Breadcrumb{next_pos, current_path.n_steps + 1}) {
						continue
					} else {
						seen = append(seen, Breadcrumb{next_pos, current_path.n_steps + 1})
					}
					// valid step!
					history := append(current_path.history, next_pos)
					queue.paths = append(queue.paths, Path{
						next_pos,
						current_path.n_steps + 1,
						history,
					})
				}
			}
		}
	}

	return 0
}

func Part1(filepath string) int {
	grid := readInput(filepath)
	return FindShortestPathLength(grid)
}

func Part2(filepath string) int {
	initial_grid := readInput(filepath)
	all_grids := []Grid{initial_grid}
	for y, row := range initial_grid.grid {
		for x, val := range row {
			if val == 'a' {
				all_grids = append(all_grids, Grid{
					grid:         initial_grid.grid,
					starting_pos: Coord{y, x},
					ending_pos:   initial_grid.ending_pos,
				})
			}
		}
	}
	best_distance := 99999
	for _, grid := range all_grids {
		dist := FindShortestPathLength(grid)
		if dist < best_distance && dist != 0 {
			best_distance = dist
		}
	}
	return best_distance
}
