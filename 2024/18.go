package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func day18(part2 bool) Solution {
	example_filepath := GetExamplePath(18)
	input_filepath := GetInputPath(18)
	if !part2 {
		example_p1 := Part1_18(example_filepath, 6, 12)
		return Solution{
			"18",
			example_p1,
			Part1_18(input_filepath, 70, 1024),
		}
	} else {
		example_p2 := Part2_18(example_filepath, 6, 12)
		AssertExample("6,1", example_p2, 2)
		return Solution{
			"18",
			example_p2,
			Part2_18(input_filepath, 70, 1024),
		}
	}
}

func GetGridDay18(filepath string, n int, num_bytes int) [][]string {
	// Initilaise an N + 1 x N +1 grid
	grid := make([][]string, n+1)
	for i := range grid {
		for range n + 1 {
			grid[i] = append(grid[i], ".")
		}
	}

	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	i := 0
	for scanner.Scan() {
		if i == num_bytes {
			break
		}
		line := scanner.Text()
		parts := strings.Split(line, ",")
		x, y := parts[0], parts[1]
		y_coord, err := strconv.Atoi(y)
		if err != nil {
			panic("trouble reading y coord")
		}
		x_coord, err := strconv.Atoi(x)
		if err != nil {
			panic("trouble reading x coord")
		}
		grid[y_coord][x_coord] = "#"
		i++
	}
	return grid
}

type Position18 struct {
	y     int
	x     int
	steps int
}

func Part1_18(filepath string, n int, num_bytes int) string {
	grid := GetGridDay18(filepath, n, num_bytes)

	Y := len(grid)
	X := len(grid[0])
	sy, sx := 0, 0
	start := Position18{sy, sx, 0} // y x steps
	queue := []Position18{start}
	seen := make(map[[2]int]bool)
	var least_steps int
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		y, x, steps := pos.y, pos.x, pos.steps
		for _, dydx := range FourSideDirs {
			dy, dx := dydx[0], dydx[1]
			ny := y + dy
			nx := x + dx
			if ny < 0 || ny >= Y || nx < 0 || nx >= X || grid[ny][nx] == "#" {
				continue
			}
			if _, exists := seen[[2]int{ny, nx}]; exists {
				continue
			}
			if ny == n && nx == n {
				least_steps = steps + 1
				break
			}
			seen[[2]int{ny, nx}] = true
			queue = append(queue, Position18{
				ny, nx, steps + 1,
			})
		}
	}

	return fmt.Sprintf("%d", least_steps)
}

func FallingBytesMazeSolver(filepath string, n int, num_bytes int) int {
	grid := GetGridDay18(filepath, n, num_bytes)
	Y := len(grid)
	X := len(grid[0])
	start := Position18{0, 0, 0} // y x steps
	queue := []Position18{start}
	seen := make(map[[2]int]bool)
	least_steps := -1
	for len(queue) > 0 {
		pos := queue[0]
		queue = queue[1:]
		y, x, steps := pos.y, pos.x, pos.steps
		for _, dydx := range FourSideDirs {
			dy, dx := dydx[0], dydx[1]
			ny := y + dy
			nx := x + dx
			if ny < 0 || ny >= Y || nx < 0 || nx >= X || grid[ny][nx] == "#" {
				continue
			}
			if _, exists := seen[[2]int{ny, nx}]; exists {
				continue
			}
			if ny == n && nx == n {
				least_steps = steps + 1
				break
			}
			seen[[2]int{ny, nx}] = true
			queue = append(queue, Position18{
				ny, nx, steps + 1,
			})
		}
	}
	return least_steps
}

func Part2_18(filepath string, n int, num_bytes int) string {
	max_bytes := n * n
	min_bytes := num_bytes
	var curr_bytes int
	for {
		curr_bytes = (max_bytes-min_bytes)/2 + min_bytes
		ans := FallingBytesMazeSolver(filepath, n, curr_bytes)
		if ans == -1 {
			if max_bytes == min_bytes {
				max_bytes--
				min_bytes--
			} else {
				max_bytes = curr_bytes
			}
		} else if max_bytes-min_bytes == 1 {
			min_bytes++
		} else if max_bytes == min_bytes {
			break
		} else {
			min_bytes = curr_bytes
		}
	}

	raw_input := readInput(filepath)
	ans := strings.Split(raw_input, "\r\n")[max_bytes]
	return fmt.Sprintf("%v", ans)
}
