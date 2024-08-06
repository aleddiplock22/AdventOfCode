package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	const EXAMPLE_FILEPATH = "example.txt"
	const INPUT_FILEPATH = "input.txt"

	example_p1 := Part1(EXAMPLE_FILEPATH)
	part1_ans := Part1(INPUT_FILEPATH)

	fmt.Println("---Day 11---")
	fmt.Printf("[Example P1] Expected: 37, Answer: %v\n", example_p1)
	fmt.Printf("[Part 1] Answer: %v\n", part1_ans)

	example_p2 := Part2(EXAMPLE_FILEPATH)
	part2_ans := Part2(INPUT_FILEPATH)

	fmt.Printf("[Example P2] Expected: 26, Answer: %v\n", example_p2)
	fmt.Printf("[Part 2] Answer: %v\n", part2_ans)
}

type coords struct {
	y int
	x int
}

var MAX_Y int
var MAX_X int

const OCCUPADO = "#"
const EMPTY = "L"
const FLOOR = "."

func parseInput(filepath string) map[coords]string {
	floor_plan := make(map[coords]string)
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	i := 0
	j_max := 0
	i_max := 0
	for scanner.Scan() {
		i_max = max(i, i_max)
		line := scanner.Text()
		for j, char := range line {
			coord := coords{
				y: i,
				x: j,
			}
			floor_plan[coord] = string(char)
			j_max = max(j, j_max)
		}
		i++
	}
	MAX_Y = j_max // I know this is so dumb lol
	MAX_X = i_max
	return floor_plan
}

func countOccupiedSeats(floor_plan map[coords]string) int {
	count := 0
	for _, val := range floor_plan {
		if val == OCCUPADO {
			count++
		}
	}
	return count
}

func getNumOccupiedNeighbours(y int, x int, floor_plan map[coords]string) int {
	num_occupied_neighbours := 0
	for _, dy := range []int{-1, 0, 1} {
		for _, dx := range []int{-1, 0, 1} {
			if dy == 0 && dx == 0 {
				continue
			}
			current_state, exists := floor_plan[coords{
				y + dy,
				x + dx,
			}]
			if exists && current_state == OCCUPADO {
				num_occupied_neighbours++
			}
		}
	}
	return num_occupied_neighbours
}

func prettyPrintGrid(floor_plan map[coords]string) {
	fmt.Println()
	for y := 0; y <= MAX_Y; y++ {
		var row []string
		for x := 0; x <= MAX_X; x++ {
			row = append(row, floor_plan[coords{
				y,
				x,
			}])
		}
		fmt.Println(row)
	}
	fmt.Println()
}

func copyMap(map1 map[coords]string) map[coords]string {
	map2 := make(map[coords]string)
	for id, value := range map1 {
		map2[id] = value
	}
	return map2
}

func runSimulation(floor_plan map[coords]string, part2 bool) map[coords]string {
	MAX_NEIGHBOURS := 4
	if part2 {
		MAX_NEIGHBOURS = 5
	}
	// for {} is equivalent to for true, i.e. while true
	for {
		// have to make a copy since we need to make changes instantaneously based on grid at start of loop.
		new_floor_plan := copyMap(floor_plan)
		// prettyPrintGrid(floor_plan)
		made_a_change := false
		for y := 0; y <= MAX_Y; y++ {
			for x := 0; x <= MAX_X; x++ {
				current_state := floor_plan[coords{
					y,
					x,
				}]
				if current_state == FLOOR {
					continue
				}
				var occupied_neighbours int
				if part2 {
					occupied_neighbours = getNumOccupiedSeenNeighbours(y, x, floor_plan)
				} else {
					occupied_neighbours = getNumOccupiedNeighbours(y, x, floor_plan)
				}
				switch current_state {
				case EMPTY:
					if occupied_neighbours == 0 {
						new_floor_plan[coords{
							y,
							x,
						}] = OCCUPADO
						made_a_change = true
					}
				case OCCUPADO:
					if occupied_neighbours >= MAX_NEIGHBOURS {
						new_floor_plan[coords{
							y,
							x,
						}] = EMPTY
						made_a_change = true
					}
				default:
					panic("Unrecognised state, neither occupied nor empty!")
				}
			}
		}
		floor_plan = new_floor_plan
		if !made_a_change {
			return new_floor_plan
		}
	}
}

func Part1(filepath string) (answer int) {
	floor_plan := parseInput(filepath)
	floor_plan = runSimulation(floor_plan, false)
	return countOccupiedSeats(floor_plan)
}

func getNumOccupiedSeenNeighbours(y int, x int, floor_plan map[coords]string) int {
	num_occupied_neighbours := 0
	for _, dy := range []int{-1, 0, 1} {
		for _, dx := range []int{-1, 0, 1} {
			if dy == 0 && dx == 0 {
				continue
			}
		searching_inner_loop: // named loops!
			for multi := 1; multi < 1000; multi++ {
				current_state, exists := floor_plan[coords{
					y + multi*dy,
					x + multi*dx,
				}]
				if !exists {
					break
				}
				switch current_state {
				case FLOOR:
					continue
				case EMPTY:
					break searching_inner_loop
				case OCCUPADO:
					num_occupied_neighbours++
					break searching_inner_loop
				}
			}
		}
	}

	return num_occupied_neighbours
}

func Part2(filepath string) (answer int) {
	floor_plan := parseInput(filepath)
	floor_plan = runSimulation(floor_plan, true)
	return countOccupiedSeats(floor_plan)
}
