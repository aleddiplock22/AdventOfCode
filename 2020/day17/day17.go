package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

/*

Had to change like all the functions for part 2 to allow the extra dimension
so was simpler just to make an extra file/folder for part2 than try and adapt everything!

*/

func main() {
	const EXAMPLE_FILEPATH = "example.txt"
	const INPUT_FILEPATH = "input.txt"

	fmt.Println("---Day 17---")

	fmt.Println("[Example P1] Expected: 112, Answer:", Part1(EXAMPLE_FILEPATH))
	fmt.Println("[Part 1] Answer:", Part1(INPUT_FILEPATH))
}

const ACTIVE = 1
const INACTIVE = 0
const STARTER_POINT = 10
const ONE_SIDE_WIDTH = 15
const MAX_GRID_SPAN = STARTER_POINT + ONE_SIDE_WIDTH

func makeGrid() [][][]int {
	allocated_grid := make([][][]int, MAX_GRID_SPAN) // annoying pre-allocation
	for z := 0; z < MAX_GRID_SPAN; z++ {
		allocated_grid[z] = make([][]int, MAX_GRID_SPAN) // Make one inner slice per iteration
		for y := 0; y < MAX_GRID_SPAN; y++ {
			allocated_grid[z][y] = make([]int, MAX_GRID_SPAN)
		}
	}
	return allocated_grid
}

func parseInput(filepath string) (grid [][][]int) {
	// Returns a 3-D cartesian grid that can be accessed like grid[z][y][x]
	// where (STARTER_POINT,STARTER_POINT,STARTER_POINT) will be the TOP LEFT of input
	// (SP,SP-1,SP) is the leftmost entry second line
	// (SP,SP-1,SP+1) is the second entry in second line

	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	grid = makeGrid()

	const Z_PLANE = STARTER_POINT
	y := STARTER_POINT
	for scanner.Scan() {
		line := scanner.Text()
		states := strings.Split(line, "")
		for x, state := range states {
			if state == "." {
				grid[Z_PLANE][y][x+STARTER_POINT] = INACTIVE
			} else if state == "#" {
				grid[Z_PLANE][y][x+STARTER_POINT] = ACTIVE
			} else {
				panic("Unknown token in initial parsing! Not '.' or '#'.")
			}
		}
		y--
	}

	return grid
}

func getNeighbourStates(z int, y int, x int, grid [][][]int) (neighbouring_states []int) {
	adjacency_diffs := []int{-1, 0, 1}
	LOWER_BOUND := 0
	UPPER_BOUND := MAX_GRID_SPAN - 1

	for _, dz := range adjacency_diffs {
		for _, dy := range adjacency_diffs {
			for _, dx := range adjacency_diffs {
				if dz == 0 && dy == 0 && dx == 0 {
					// considering self
					continue
				}
				nz := z + dz
				ny := y + dy
				nx := x + dx
				if LOWER_BOUND <= nz && nz <= UPPER_BOUND &&
					LOWER_BOUND <= ny && ny <= UPPER_BOUND &&
					LOWER_BOUND <= nx && nx <= UPPER_BOUND {
					// within bounds
					neighbouring_states = append(neighbouring_states, grid[nz][ny][nx])
				}
			}
		}
	}
	return neighbouring_states
}

func numActive(neighbours []int) int {
	counter := 0
	for _, state := range neighbours {
		counter += state
	}
	return counter
}

func runSimulation(N int, grid [][][]int) [][][]int {
	for n := 1; n <= N; n++ {
		new_grid := makeGrid()
		for z, z_axis := range grid {
			for y, y_axis := range z_axis {
				for x := range y_axis {
					point := grid[z][y][x]
					point_neighbours := getNeighbourStates(z, y, x, grid)
					num_active := numActive(point_neighbours)

					var new_point int
					switch point {
					case ACTIVE:
						if num_active == 2 || num_active == 3 {
							// remain active
							new_point = ACTIVE
						} else {
							// become inactive
							new_point = INACTIVE
						}
					case INACTIVE:
						if num_active == 3 {
							// becomes active
							new_point = ACTIVE
						} else {
							// remains inactive
							new_point = INACTIVE
						}
					default:
						panic("Unknown non active or inactive state??")
					}

					new_grid[z][y][x] = new_point
				}
			}
		}
		copy(grid, new_grid)
	}

	return grid
}

func countTotalActive(grid [][][]int) int {
	total := 0
	for z, z_axis := range grid {
		for y, y_axis := range z_axis {
			for x := range y_axis {
				total += grid[z][y][x]
			}
		}
	}
	return total
}

func Part1(filepath string) int {
	const N_CYCLES = 6
	grid := parseInput(filepath)
	grid = runSimulation(N_CYCLES, grid)
	answer := countTotalActive(grid)

	return answer
}
