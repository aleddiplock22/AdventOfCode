package main

import (
	"fmt"
	"slices"
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
		example_p2 := Part2_15(example_filepath)
		AssertExample("9021", example_p2, 2)
		return Solution{
			"15",
			example_p2,
			Part2_15(input_filepath),
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
				tmp_y, tmp_x = doInstructionTransform(tmp_y, tmp_x, instruction)
				switch grid[tmp_y][tmp_x] {
				case "#":
					// wall, so scrap everything lol
					valid = false
					break transform_loop
				case ".":
					// ok everything can happen!
					transforms = append(transforms, [3]int{tmp_y, tmp_x, -3})
					break transform_loop
				case "O":
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

func MakeWideGridFromGrid(grid [][]string) (new_grid [][]string) {
	for _, row := range grid {
		new_row := []string{}
		for _, char := range row {
			switch char {
			case "#":
				new_row = append(new_row, "#", "#")
			case "O":
				new_row = append(new_row, "[", "]")
			case ".":
				new_row = append(new_row, ".", ".")
			case "@":
				new_row = append(new_row, "@", ".")
			default:
				panic("who is in my warehouse...")
			}
		}
		new_grid = append(new_grid, new_row)
	}

	return new_grid
}

/*
returns slice of [3]int{y_pos, x_pos, char_as_int} that should move in direction
where direction follows: (0 UP 1 RIGHT 2 DOWN 3 LEFT)
where char_as_int follows:

	-1 = .
	-2 = @
	-3 = [
	-4 = ]
*/
func findAllConnectedToMoveVertically(y, x int, direction int, grid *[][]string) (to_move [][2]int) {
	/*
		........
		..[][]..
		...[]...
		....@...

		->

		..[][]..
		...[]...
		....@...
		........

		y,x will be the ] or [ directly above or below an @
	*/
	if !(direction == 0 || direction == 2) {
		panic("invalid non up or down direction passed to findAllConnectedToMoveVertically")
	}

	to_move = append(to_move, [2]int{y, x})

	to_check_in_dir := [][2]int{{y, x}}
	var partner_y, partner_x int
	if (*grid)[y][x] == "]" {
		partner_y, partner_x = y, x-1
	} else {
		if (*grid)[y][x] != "[" {
			panic("expected [ bracket here, but found")
		}
		partner_y, partner_x = y, x+1
	}
	to_move = append(to_move, [2]int{partner_y, partner_x})
	to_check_in_dir = append(to_check_in_dir, [2]int{partner_y, partner_x})
	seen := map[[2]int]bool{{y, x}: true, {partner_y, partner_x}: true}

	var cy, cx int
	for len(to_check_in_dir) > 0 {
		_pos := to_check_in_dir[0]
		to_check_in_dir = to_check_in_dir[1:]
		cy, cx = _pos[0], _pos[1]

		char_at_pos := (*grid)[cy][cx]
		if char_at_pos == "." {
			// harmless
			continue
		}
		if char_at_pos == "#" {
			// wall - nothing should move
			return [][2]int{}
		}

		if char_at_pos == "[" {
			partner_y, partner_x = cy, cx+1
		} else if char_at_pos == "]" {
			partner_y, partner_x = cy, cx-1
		} else {
			panic("WHOM in warehouse")
		}

		_partner_pos := [2]int{partner_y, partner_x}

		if !slices.Contains(to_move, _pos) {
			to_move = append(to_move, _pos)
		}
		if _, exists := seen[_partner_pos]; !exists {
			to_check_in_dir = append(to_check_in_dir, _partner_pos)
			seen[_partner_pos] = true
		}
		_p_y, _p_x := doInstructionTransform(cy, cx, direction)
		_next := [2]int{_p_y, _p_x}
		if _, exists := seen[_next]; !exists {
			to_check_in_dir = append(to_check_in_dir, _next)
			seen[_next] = true
		}
	}
	return to_move
}

func Part2_15(filepath string) string {
	grid, instructions := ParseGridAndInstructions(filepath)
	grid = MakeWideGridFromGrid(grid)

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

		// boxes are where it differs from p1, as boxes are fat now
		if next_pos == "[" || next_pos == "]" {

			// so it is still simple for left right situation, lets do that first
			if instruction == 1 || instruction == 3 {
				/*
					-1 = .
					-2 = @
					-3 = [
					-4 = ]
				*/
				transforms := [][3]int{
					{cy, cx, -1},
					{ny, nx, -2},
				}
				tmp_y, tmp_x := ny, nx
				var LBOX, RBOX int
				if next_pos == "[" {
					LBOX, RBOX = -3, -4
				} else if next_pos == "]" {
					LBOX, RBOX = -4, -3
				}
				// ok so we have the initial transform done
				valid := true
			transform_loop:
				for {
					tmp_y, tmp_x = doInstructionTransform(tmp_y, tmp_x, instruction)
					switch grid[tmp_y][tmp_x] {
					case "#":
						// wall, so scrap everything lol
						valid = false
						break transform_loop
					case ".":
						// ok everything can happen!
						break transform_loop
					case "[":
						// another box
						transforms = append(transforms, [3]int{tmp_y, tmp_x, LBOX})
						tmp_y, tmp_x = doInstructionTransform(tmp_y, tmp_x, instruction)
						transforms = append(transforms, [3]int{tmp_y, tmp_x, RBOX})
						if grid[tmp_y][tmp_x] == "." {
							break transform_loop
						}
					case "]":
						// another box
						transforms = append(transforms, [3]int{tmp_y, tmp_x, LBOX})
						tmp_y, tmp_x = doInstructionTransform(tmp_y, tmp_x, instruction)
						transforms = append(transforms, [3]int{tmp_y, tmp_x, RBOX})
						if grid[tmp_y][tmp_x] == "." {
							break transform_loop
						}
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
							grid[ty][tx] = "["
						case -4:
							grid[ty][tx] = "]"
						}
					}
					cy, cx = doInstructionTransform(cy, cx, instruction)
				}
			} else if instruction == 0 || instruction == 2 {
				to_move := findAllConnectedToMoveVertically(ny, nx, instruction, &grid)
				if len(to_move) == 0 {
					// just dont do anything, wall situation
					continue
				}
				/*
					-1 = .
					-2 = @
					-3 = [
					-4 = ]
				*/
				transforms := make(map[[4]int]string)
				changing := map[[2]int]bool{{cy, cx}: true, {ny, nx}: true}
				for _, move := range to_move {
					_y, _x := move[0], move[1]
					_ny, _nx := doInstructionTransform(_y, _x, instruction)
					transforms[[4]int{_y, _x, _ny, _nx}] = grid[_y][_x]
				}
				for pos, transform := range transforms {
					og_y, og_x, pos_y, pos_x := pos[0], pos[1], pos[2], pos[3]
					if _, exist := changing[[2]int{og_y, og_x}]; !exist {
						// it's not getting changed, so we should make it a .
						changing[[2]int{og_y, og_x}] = true
						grid[og_y][og_x] = "."
					}
					// if _, exist := changing[[2]int{pos_y, pos_x}]; !exist {
					// 	// it's not getting changed, so we should make it a .
					// 	changing[[2]int{pos_y, pos_x}] = true
					// 	grid[pos_y][pos_x] = transform
					// }
					changing[[2]int{pos_y, pos_x}] = true
					grid[pos_y][pos_x] = transform
				}

				grid[cy][cx] = "."
				grid[ny][nx] = "@"
				cy, cx = ny, nx

			} else {
				panic("? instruction not in (0,1,2,3)")
			}
		}
	}

	/*
		For these larger boxes, distances are measured from the edge of the map
		to the closest edge of the box in question.
	*/
	box_pairs := make(map[[4]int]bool)
	var total int
	for y, row := range grid {
		for x, char := range row {
			if char == "[" {
				box_pairs[[4]int{y, x, y, x + 1}] = true
			} else if char == "]" {
				box_pairs[[4]int{y, x - 1, y, x}] = true
			}
		}
	}
	for pair := range box_pairs {
		ly, lx, ry, rx := pair[0], pair[1], pair[2], pair[3]
		y := min(ly, ry)
		x := min(lx, rx)
		total += 100*y + x
	}

	return fmt.Sprintf("%d", total)
}
