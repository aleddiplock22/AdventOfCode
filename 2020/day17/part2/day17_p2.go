package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	const EXAMPLE_FILEPATH = "../example.txt"
	const INPUT_FILEPATH = "../input.txt"

	fmt.Println("---Day 17---")

	fmt.Println("[Example P2] Expected: 848, Answer:", Part2(EXAMPLE_FILEPATH))
	fmt.Println("[Part 2] Answer:", Part2(INPUT_FILEPATH))
}

const ACTIVE = 1
const INACTIVE = 0
const STARTER_POINT = 10
const ONE_SIDE_WIDTH = 12 // note with different input this is worth increasing/decreasing to balance performance with capturing possible further out changes!
const MAX_GRID_SPAN = STARTER_POINT + ONE_SIDE_WIDTH

func makeGrid() [][][][]int {
	allocated_grid := make([][][][]int, MAX_GRID_SPAN) // annoying pre-allocation
	for w := 0; w < MAX_GRID_SPAN; w++ {
		allocated_grid[w] = make([][][]int, MAX_GRID_SPAN)
		for z := 0; z < MAX_GRID_SPAN; z++ {
			allocated_grid[w][z] = make([][]int, MAX_GRID_SPAN) // Make one inner slice per iteration
			for y := 0; y < MAX_GRID_SPAN; y++ {
				allocated_grid[w][z][y] = make([]int, MAX_GRID_SPAN)
			}
		}
	}

	return allocated_grid
}

func parseInput(filepath string) (grid [][][][]int) {
	// Returns a 4-D cartesian grid that can be accessed like grid[w][z][y][x]

	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	grid = makeGrid()

	const W_PLANE = STARTER_POINT
	const Z_PLANE = STARTER_POINT
	y := STARTER_POINT
	for scanner.Scan() {
		line := scanner.Text()
		states := strings.Split(line, "")
		for x, state := range states {
			if state == "." {
				grid[W_PLANE][Z_PLANE][y][x+STARTER_POINT] = INACTIVE
			} else if state == "#" {
				grid[W_PLANE][Z_PLANE][y][x+STARTER_POINT] = ACTIVE
			} else {
				panic("Unknown token in initial parsing! Not '.' or '#'.")
			}
		}
		y--
	}

	return grid
}

func getNeighbourStates(w int, z int, y int, x int, grid [][][][]int) (neighbouring_states []int) {
	adjacency_diffs := []int{-1, 0, 1}
	LOWER_BOUND := 0
	UPPER_BOUND := MAX_GRID_SPAN - 1

	for _, dw := range adjacency_diffs {
		for _, dz := range adjacency_diffs {
			for _, dy := range adjacency_diffs {
				for _, dx := range adjacency_diffs {
					if dz == 0 && dy == 0 && dx == 0 && dw == 0 {
						// considering self
						continue
					}
					nw := w + dw
					nz := z + dz
					ny := y + dy
					nx := x + dx
					if LOWER_BOUND <= nw && nw <= UPPER_BOUND &&
						LOWER_BOUND <= nz && nz <= UPPER_BOUND &&
						LOWER_BOUND <= ny && ny <= UPPER_BOUND &&
						LOWER_BOUND <= nx && nx <= UPPER_BOUND {
						// within bounds
						neighbouring_states = append(neighbouring_states, grid[nw][nz][ny][nx])
					}
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

func runSimulation(N int, grid [][][][]int) [][][][]int {
	for n := 1; n <= N; n++ {
		new_grid := makeGrid()
		for w, w_axis := range grid {
			for z, z_axis := range w_axis {
				for y, y_axis := range z_axis {
					for x := range y_axis {
						point := grid[w][z][y][x]
						point_neighbours := getNeighbourStates(w, z, y, x, grid)
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

						new_grid[w][z][y][x] = new_point
					}
				}
			}
		}

		copy(grid, new_grid)
	}

	return grid
}

func countTotalActive(grid [][][][]int) int {
	total := 0
	for w, w_axis := range grid {
		for z, z_axis := range w_axis {
			for y, y_axis := range z_axis {
				for x := range y_axis {
					total += grid[w][z][y][x]
				}
			}
		}
	}
	return total
}

func Part2(filepath string) int {
	const N_CYCLES = 6
	grid := parseInput(filepath)
	grid = runSimulation(N_CYCLES, grid)
	answer := countTotalActive(grid)

	return answer
}
