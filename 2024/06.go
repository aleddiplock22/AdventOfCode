package main

import (
	"errors"
	"fmt"
	"sync"
)

const HASHTAG = "#"
const DOT = "."
const UP_ARROW = "^"

const UP6 = 0
const RIGHT6 = 1
const DOWN6 = 2
const LEFT6 = 3

func day06(part2 bool) Solution {
	example_filepath := GetExamplePath(6)
	input_filepath := GetInputPath(6)
	if !part2 {
		return Solution{
			"06",
			Part1_06(example_filepath),
			Part1_06(input_filepath),
		}
	} else {
		return Solution{
			"06",
			Part2_06(example_filepath),
			Part2_06(input_filepath),
		}
	}
}

func manouver(next *[3]int, grid *[][]string, count int) error {
	if count == 4 {
		return errors.New("tried all directions I guess?")
	}
	after := getNextMove6(*next)
	c, r := after[0], after[1]
	if c < 0 || c >= len(*grid) || r < 0 || r >= len((*grid)[0]) {
		return errors.New("out of bounds")
	}
	switch next[2] {
	case UP6:
		if (*grid)[next[0]-1][next[1]] == HASHTAG {
			next[2] = (next[2] + 1) % 4
			manouver(next, grid, count+1)
		}
	case RIGHT6:
		if (*grid)[next[0]][next[1]+1] == HASHTAG {
			next[2] = (next[2] + 1) % 4
			manouver(next, grid, count+1)
		}
	case DOWN6:
		if (*grid)[next[0]+1][next[1]] == HASHTAG {
			next[2] = (next[2] + 1) % 4
			manouver(next, grid, count+1)
		}
	case LEFT6:
		if (*grid)[next[0]][next[1]-1] == HASHTAG {
			next[2] = (next[2] + 1) % 4
			manouver(next, grid, count+1)
		}
	default:
		panic(fmt.Sprintf("Unknown direction in manouver: %v!\n", next[2]))
	}
	return nil
}

func getNextMove6(pos [3]int) (next [3]int) {
	switch pos[2] {
	case UP6:
		next = [3]int{pos[0] - 1, pos[1], pos[2]}
	case RIGHT6:
		next = [3]int{pos[0], pos[1] + 1, pos[2]}
	case DOWN6:
		next = [3]int{pos[0] + 1, pos[1], pos[2]}
	case LEFT6:
		next = [3]int{pos[0], pos[1] - 1, pos[2]}
	}
	return next
}

func Part1_06(filepath string) string {
	grid := readStringGrid(filepath)
	sc, sr := -1, -1
outer:
	for c, col := range grid {
		for r, val := range col {
			if val == UP_ARROW {
				sc, sr = c, r
				break outer
			}
		}
	}
	if sc == -1 || sr == -1 {
		panic("Couldn't find starting position.")
	}

	// [3]int = COL ROW DIR
	current_pos := [3]int{sc, sr, UP6}
	seen := map[[3]int]bool{current_pos: true}

	if grid[sc-1][sr] == HASHTAG {
		current_pos[2] = RIGHT6
	}

	basic_seen := map[[2]int]bool{{sc, sr}: true} // need this to not double count positions where we're just facing a diff direction
	var next [3]int
	for {
		// move
		next = getNextMove6(current_pos)
		c, r := next[0], next[1]
		if c < 0 || c >= len(grid) || r < 0 || r >= len(grid[0]) {
			break
		}
		if grid[c][r] == HASHTAG {
			panic("Unexpected obstacle")
		}

		if _, exists := seen[next]; exists {
			panic("This doesn't happen in p1?")
		}
		seen[next] = true
		basic_seen[[2]int{c, r}] = true

		// check next move orientation
		err := manouver(&next, &grid, 0)
		if err != nil {
			break
		}

		current_pos = next
	}
	return fmt.Sprintf("%v", len(basic_seen))

}

func getStartingPos(grid *[][]string) (starting_pos [2]int, err error) {
	sc, sr := -1, -1
outer:
	for c, col := range *grid {
		for r, val := range col {
			if val == UP_ARROW {
				sc, sr = c, r
				break outer
			}
		}
	}
	if sc == -1 || sr == -1 {
		return [2]int{-1, -1}, errors.New("couldn't find starting position")
	}
	return [2]int{sc, sr}, nil
}

func getPossibleBlockageLocations(grid *[][]string, starting_pos *[2]int) (positions map[[2]int]bool) {
	sc, sr := starting_pos[0], starting_pos[1]
	current_pos := [3]int{sc, sr, UP6}

	if (*grid)[sc-1][sr] == HASHTAG {
		current_pos[2] = RIGHT6
	}

	var next [3]int
	positions = make(map[[2]int]bool)
	for {
		// move
		next = getNextMove6(current_pos)
		c, r := next[0], next[1]
		if c < 0 || c >= len((*grid)) || r < 0 || r >= len((*grid)[0]) {
			break
		}
		if (*grid)[c][r] == HASHTAG {
			panic("Unexpected obstacle")
		}
		positions[[2]int{c, r}] = true

		// check next move orientation
		err := manouver(&next, grid, 0)
		if err != nil {
			break
		}

		current_pos = next
	}
	return positions
}

func Part2_06(filepath string) string {
	grid := readStringGrid(filepath)
	starting_pos, err := getStartingPos(&grid)
	if err != nil {
		panic("Couldn't find starting position!")
	}

	var total int
	var wg sync.WaitGroup
	results := make(chan bool)

	possible_blocks := getPossibleBlockageLocations(&grid, &starting_pos)

	for c, col := range grid {
		for r := range col {
			obstacle_pos := [2]int{c, r}
			if exists := possible_blocks[obstacle_pos]; !exists {
				continue
			}
			if grid[obstacle_pos[0]][obstacle_pos[1]] == HASHTAG || obstacle_pos == starting_pos {
				continue
			}
			wg.Add(1)
			go func(obstacle_pos [2]int) {
				defer wg.Done()
				result := doesP1SimulationGetLoop(&grid, &starting_pos, &obstacle_pos)
				results <- result
			}(obstacle_pos)
		}
	}

	// close channel when all goroutines are done:
	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		if result {
			total++
		}
	}

	return fmt.Sprintf("%v", total)
}

func doesP1SimulationGetLoop(base_grid *[][]string, starting_pos *[2]int, obstacle_pos *[2]int) bool {
	// Avoiding copies bc wtf is going on ....
	sc, sr := starting_pos[0], starting_pos[1]
	current_pos := [3]int{sc, sr, UP6}

	grid := make([][]string, len(*base_grid))
	for i := range *base_grid {
		grid[i] = make([]string, len((*base_grid)[i]))
		copy(grid[i], (*base_grid)[i])
	}

	// Add obstacle in
	grid[obstacle_pos[0]][obstacle_pos[1]] = HASHTAG

	seen := map[[3]int]bool{current_pos: true}

	err := manouver(&current_pos, &grid, 0)
	if err != nil {
		panic("dont think we should be stuck on the first move...?")
	}

	basic_seen := map[[2]int]bool{{sc, sr}: true} // need this to not double count positions where we're just facing a diff direction
	var next [3]int
	for {
		// move
		next = getNextMove6(current_pos)
		c, r := next[0], next[1]
		if c < 0 || c >= len(grid) || r < 0 || r >= len(grid[0]) {
			return true
		}
		if grid[c][r] == HASHTAG {
			fmt.Printf("c,r=%v,%v | obstacle_pos=%v, \n", c, r, obstacle_pos)
			panic("Unexpected obstacle")
		}

		if _, exists := seen[next]; exists {
			return true
		}
		seen[next] = true
		basic_seen[[2]int{c, r}] = true

		// check next move orientation
		err := manouver(&next, &grid, 0)
		if err != nil {
			return false
		}
		current_pos = next
	}
}
