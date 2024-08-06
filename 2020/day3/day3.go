package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("---Day 3---")
	example_ans := NumTreesHit("example.txt", 3, 1)
	fmt.Printf("[Example P1] Expected: 7 Got: %d\n", example_ans)
	input_ans := NumTreesHit("input.txt", 3, 1)
	fmt.Printf("Part 1: %d\n", input_ans)

	x_slopes := [5]int{1, 3, 5, 7, 1}
	y_slopes := [5]int{1, 1, 1, 1, 2}
	var y int
	trees_hit_example_p2 := 1
	trees_hit_p2 := 1
	for idx, x := range x_slopes {
		y = y_slopes[idx]
		trees_hit_example_p2 = trees_hit_example_p2 * NumTreesHit("example.txt", x, y)
		trees_hit_p2 = trees_hit_p2 * NumTreesHit("input.txt", x, y)
	}
	fmt.Printf("\n[Example P2] Expected: 336 Got: %d\n", trees_hit_example_p2)
	fmt.Printf("Part 2: %d\n", trees_hit_p2)
}

func NumTreesHit(filename string, slope_x int, slope_y int) (trees_hit int) {
	TREE_CHAR := "#"
	trees_hit = 0
	file, _ := os.Open(filename)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	first_scan := true
	var width int
	x_pos := 0
	y_pos := 0
	for scanner.Scan() {
		if y_pos%slope_y != 0 {
			y_pos++
			continue
		} else {
			y_pos++
		}
		line := scanner.Text()
		// setup
		if first_scan {
			width = len(line)
		}
		for x_pos > (width - 1) {
			x_pos = x_pos - width
		}
		// check for tree
		if line[x_pos:x_pos+1] == TREE_CHAR {
			trees_hit++
		}
		x_pos = x_pos + slope_x
	}
	return trees_hit
}
