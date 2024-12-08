package main

import "fmt"

func day08(part2 bool) Solution {
	example_filepath := GetExamplePath(8)
	input_filepath := GetInputPath(8)
	if !part2 {
		example_p1 := Part1_08(example_filepath)
		AssertExample("14", example_p1, 1)
		return Solution{
			"08",
			example_p1,
			Part1_08(input_filepath),
		}
	} else {
		example_p2 := Part2_08(example_filepath)
		AssertExample("34", example_p2, 2)
		return Solution{
			"08",
			example_p2,
			Part2_08(input_filepath),
		}
	}
}

func FindAntinodes(locations [][2]int) [][2]int {
	/*
		An antinode occurs at any point that is perfectly
		in line with two antennas of the same frequency
		- but only when one of the antennas is twice as far away as the other.
	*/

	// 1. Create pairs from all the locations
	// 2. For each pair find the 2 antinodes either side of them

	var pairs [][2][2]int
	for i, loc := range locations {
		if i == len(locations)-1 {
			break
		}
		for _, loc2 := range locations[i+1:] {
			pairs = append(pairs, [2][2]int{loc, loc2})
		}
	}

	var antinodes [][2]int
	for _, pair := range pairs {
		a, b := pair[0], pair[1]
		dist_y := a[0] - b[0]
		dist_x := a[1] - b[1]
		antinodes = append(antinodes, [2]int{a[0] + dist_y, a[1] + dist_x})
		antinodes = append(antinodes, [2]int{b[0] - dist_y, b[1] - dist_x})
	}

	return antinodes
}

func FindAntinodes_p2(locations [][2]int, Y int, X int) [][2]int {
	/*
		An antinode occurs at any point that is perfectly
		in line with two antennas of the same frequency
		 ---- [ as many as fit on the grid now! ]  ----
	*/

	var pairs [][2][2]int
	for i, loc := range locations {
		if i == len(locations)-1 {
			break
		}
		for _, loc2 := range locations[i+1:] {
			pairs = append(pairs, [2][2]int{loc, loc2})
		}
	}

	var antinodes [][2]int
	for _, pair := range pairs {
		a, b := pair[0], pair[1]
		dist_y := a[0] - b[0]
		dist_x := a[1] - b[1]

		i := 1
		for {
			new_y := a[0] + i*dist_y
			new_x := a[1] + i*dist_x
			if new_y >= 0 && new_y < Y && new_x >= 0 && new_x < X {
				antinodes = append(antinodes, [2]int{new_y, new_x})
				i++
			} else {
				break
			}
		}
		j := 1
		for {
			new_y := b[0] - j*dist_y
			new_x := b[1] - j*dist_x
			if new_y >= 0 && new_y < Y && new_x >= 0 && new_x < X {
				antinodes = append(antinodes, [2]int{new_y, new_x})
				j++
			} else {
				break
			}
		}
	}

	return antinodes
}

func Part1_08(filepath string) string {
	grid := readStringGrid(filepath)
	Y := len(grid)
	X := len(grid[0])

	antennas := make(map[string][][2]int)
	for c, col := range grid {
		for r, char := range col {
			if char == "." {
				continue
			}
			antennas[char] = append(antennas[char], [2]int{c, r})
		}
	}

	antinodes_map := make(map[[2]int]bool)
	for _, locs := range antennas {
		antinodes := FindAntinodes(locs)
		for _, antinode := range antinodes {
			y, x := antinode[0], antinode[1]
			if y >= 0 && y < Y && x >= 0 && x < X {
				antinodes_map[antinode] = true
			}
		}
	}

	return fmt.Sprintf("%v", len(antinodes_map))
}

func Part2_08(filepath string) string {
	grid := readStringGrid(filepath)
	Y := len(grid)
	X := len(grid[0])

	antennas := make(map[string][][2]int)
	for c, col := range grid {
		for r, char := range col {
			if char == "." {
				continue
			}
			antennas[char] = append(antennas[char], [2]int{c, r})
		}
	}

	antinodes_map := make(map[[2]int]bool)
	for _, antennna_locs := range antennas {
		for _, loc := range antennna_locs {
			antinodes_map[loc] = true
		}
	}
	for _, locs := range antennas {
		antinodes := FindAntinodes_p2(locs, Y, X)
		for _, antinode := range antinodes {
			y, x := antinode[0], antinode[1]
			if !(y >= 0 && y < Y && x >= 0 && x < X) {
				panic("out of bounds antinode in p2??")
			} else {
				antinodes_map[antinode] = true
			}
		}
	}

	return fmt.Sprintf("%v", len(antinodes_map))
}
