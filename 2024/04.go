package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func day04(part2 bool) Solution {
	input := GetInputPath(4)
	if !part2 {
		example_p1 := Part1_04(GetExamplePath(4))
		AssertExample("18", example_p1, 4)
		return Solution{
			"04",
			example_p1,
			Part1_04(input),
		}
	} else {
		example_p2 := Part2_04(GetExamplePath(4))
		AssertExample("9", example_p2, 4)
		return Solution{
			"04",
			example_p2,
			Part2_04(input),
		}
	}
}

// might move this function to utils.go
func readStringGrid(filepath string) [][]string {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	grid := [][]string{}
	for scanner.Scan() {
		line := scanner.Text()
		chars := strings.Split(line, "")
		grid = append(grid, chars)
	}
	return grid
}

func checkXMASLocs(locs [4][2]int, grid [][]string) bool {
	return grid[locs[0][0]][locs[0][1]] == "X" && grid[locs[1][0]][locs[1][1]] == "M" && grid[locs[2][0]][locs[2][1]] == "A" && grid[locs[3][0]][locs[3][1]] == "S"
}

func Part1_04(filepath string) string {
	grid := readStringGrid(filepath) //grid[y=col][x=row]

	Y := len(grid)
	X := len(grid[0])
	var total int
	for x, row := range grid {
		for y := range row {
			var locations [][4][2]int
			// horiz ->
			if x+3 < X {
				locations = append(locations, [4][2]int{
					{y, x}, {y, x + 1}, {y, x + 2}, {y, x + 3},
				})
			}
			// horiz <-
			if x-3 >= 0 {
				locations = append(locations, [4][2]int{
					{y, x}, {y, x - 1}, {y, x - 2}, {y, x - 3},
				})
			}
			// vert DOWN
			if y+3 < Y {
				locations = append(locations, [4][2]int{
					{y, x}, {y + 1, x}, {y + 2, x}, {y + 3, x},
				})
			}
			// vert UP
			if y-3 >= 0 {
				locations = append(locations, [4][2]int{
					{y, x}, {y - 1, x}, {y - 2, x}, {y - 3, x},
				})
			}
			// diag UPRIGHT
			if y-3 >= 0 && x+3 < X {
				locations = append(locations, [4][2]int{
					{y, x}, {y - 1, x + 1}, {y - 2, x + 2}, {y - 3, x + 3},
				})
			}
			// diag UPLEFT
			if y-3 >= 0 && x-3 >= 0 {
				locations = append(locations, [4][2]int{
					{y, x}, {y - 1, x - 1}, {y - 2, x - 2}, {y - 3, x - 3},
				})
			}
			// diag DOWNLEFT
			if y+3 < Y && x-3 >= 0 {
				locations = append(locations, [4][2]int{
					{y, x}, {y + 1, x - 1}, {y + 2, x - 2}, {y + 3, x - 3},
				})
			}
			// diag DOWNRIGHT
			if y+3 < Y && x+3 < X {
				locations = append(locations, [4][2]int{
					{y, x}, {y + 1, x + 1}, {y + 2, x + 2}, {y + 3, x + 3},
				})
			}
			for _, locs := range locations {
				if checkXMASLocs(locs, grid) {
					total++
				}
			}
		}
	}

	return fmt.Sprintf("%v", total)
}

func checkMasX(TL string, TR string, BL string, BR string) bool {
	MAS := "MAS"
	SAM := "SAM"
	diag1 := fmt.Sprintf("%vA%v", TL, BR)
	diag2 := fmt.Sprintf("%vA%v", TR, BL)
	return (diag1 == MAS || diag1 == SAM) && (diag2 == MAS || diag2 == SAM)
}

func Part2_04(filepath string) string {
	grid := readStringGrid(filepath) //grid[y=col][x=row]

	Y := len(grid)
	X := len(grid[0])
	var total int

	/*
		We're now looking for a diagonal X-MAX
		each M A S can be backwards

		M		S
			A
		M		S

		S		M
			A
		S		M

		both valid etc.
	*/
	for x, row := range grid {
		for y := range row {
			// [][TOPLEFT TOPRIGHT BOTTOMLEFT BOTTOMRIGHT]
			if grid[y][x] == "A" {
				// possible of central X-MAS
				if y-1 >= 0 && y+1 < Y && x-1 >= 0 && x+1 < X {
					if checkMasX(grid[y-1][x-1], grid[y-1][x+1], grid[y+1][x-1], grid[y+1][x+1]) {
						total++
					}
				}
			}
		}

	}

	return fmt.Sprintf("%v", total)
}
