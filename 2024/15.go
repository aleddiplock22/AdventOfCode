package main

import (
	"fmt"
	"strings"
)

func day15(part2 bool) Solution {
	example_filepath := GetExamplePath(15)
	input_filepath := GetInputPath(15)
	if !part2 {
		example_p1 := Part1_15(example_filepath)
		AssertExample("10092", example_p1, 1)
		return Solution{
			"15",
			example_p1,
			Part1_15(input_filepath),
		}
	} else {
		return Solution{
			"15",
			"example part 2",
			"input part 2",
		}
	}
}

func ParseGridAndInstructions(filepath string) (grid [][]string, instructions []int) {
	raw_input := readInput(filepath)
	parts := strings.Split(raw_input, "\r\n\r\n")
	grid_lines := strings.Split(parts[0], "\r\n")
	for _, line := range grid_lines {
		chars := strings.Split(line, "")
		grid = append(grid, chars)
	}
	for _, instruction_line := range strings.Split(parts[1], "\r\n") {
	instruction_loop:
		for _, instruction := range strings.Split(instruction_line, "") {
			var dir int
			switch instruction {
			case "^":
				dir = 0
			case ">":
				dir = 1
			case "v":
				dir = 2
			case "<":
				dir = 3
			default:
				break instruction_loop // reached end, trailing empty char for some reason idk
			}
			instructions = append(instructions, dir)
		}
	}

	return grid, instructions
}

func doInstructionTransform(y, x int, instruction int) (ny, nx int) {
	switch instruction {
	case 0:
		ny, nx = y-1, x
	case 1:
		ny, nx = y, x+1
	case 2:
		ny, nx = y+1, x
	case 3:
		ny, nx = y, x-1
	}
	return ny, nx
}

func Part1_15(filepath string) string {
	grid, instructions := ParseGridAndInstructions(filepath)

	var sy, sx int // starting position
outer_starting_loop:
	for y, row := range grid {
		for x, tile := range row {
			if tile == "@" {
				sy, sx = y, x
				break outer_starting_loop
			}
		}
	}

	cy, cx := sy, sx // current position

	for _, instruction := range instructions {
		// 0 UP 1 RIGHT 2 DOWN 3 LEFT
		// if robots can move, MOVE THEM
		ny, nx := doInstructionTransform(cy, cx, instruction)
		next_pos := grid[ny][nx]

		if next_pos == "#" {
			// dont move there, it's a wall!
			continue
		}

		if next_pos == "." {
			// free space :)
			grid[cy][cx] = "."
			grid[ny][nx] = "@"
			cy, cx = ny, nx
			continue
		}

		if next_pos == "O" {
			// robot situation....
			/*
				-1 = .
				-2 = @
				-3 = O
			*/
			transforms := [][3]int{
				{cy, cx, -1},
				{ny, nx, -2},
			}
			tmp_y, tmp_x := ny, nx
			valid := true
		transform_loop:
			for {
				// fmt.Println("in transform loop:", tmp_y, tmp_x)
				tmp_y, tmp_x = doInstructionTransform(tmp_y, tmp_x, instruction)
				// fmt.Print("\n... did instruction ...\n", tmp_y, tmp_x)
				switch grid[tmp_y][tmp_x] {
				case "#":
					// fmt.Println("WALL")
					// wall, so scrap everything lol
					valid = false
					break transform_loop
				case ".":
					// fmt.Println("DOT")
					// ok everything can happen!
					transforms = append(transforms, [3]int{tmp_y, tmp_x, -3})
					break transform_loop
				case "O":
					// fmt.Println("ROBOT")
					// another robot... we go agane
					transforms = append(transforms, [3]int{tmp_y, tmp_x, -3})
				}
			}
			if valid {
				for _, transform := range transforms {
					ty, tx, tt := transform[0], transform[1], transform[2]
					switch tt {
					case -1:
						grid[ty][tx] = "."
					case -2:
						grid[ty][tx] = "@"
					case -3:
						grid[ty][tx] = "O"
					}
				}
				cy, cx = doInstructionTransform(cy, cx, instruction)
			}
		}
	}

	var total int
	for y, row := range grid {
		for x, char := range row {
			if char == "O" {
				// robot.. or apparently box I just misunderstood the lore
				total += 100*y + x
			}
		}
	}

	return fmt.Sprintf("%d", total)
}
