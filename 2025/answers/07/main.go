package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

const TEST_FILEPATH string = "../../data/7/test.txt"
const INPUT_FILEPATH string = "../../data/7/input.txt"

func main() {
	fmt.Println("🎄 Advent Of Code 2025 | Day 7 🎅")

	Part1Test := Part1(TEST_FILEPATH)
	Part1Answer := Part1(INPUT_FILEPATH)
	fmt.Printf("[Part 1] Test: %v | Real: %v\n", Part1Test, Part1Answer)

	Part2Test := Part2(TEST_FILEPATH)
	Part2Answer := Part2(INPUT_FILEPATH)
	fmt.Printf("[Part 2] Test: %v | Real: %v\n", Part2Test, Part2Answer)
}

func parseInput(filepath string) (grid [][]string) {
	file, err := os.Open(filepath)
	if err != nil {
		panic("trouble reading in file")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		grid = append(grid, strings.Split(line, ""))
	}

	return grid
}

type Beam struct {
	x int
	y int
}

func PrintGridWithBeams(grid [][]string, beams []Beam) {
	fmt.Println("------------------------------------------")
	for y, row := range grid {
		fmt.Println()
		for x, grid_entry := range row {
			char := grid_entry
			for _, beam := range beams {
				if beam.x == x && beam.y == y {
					char = "|"
					break
				}
			}
			fmt.Printf(" %v ", char)
		}
	}
	fmt.Println()
	fmt.Println("------------------------------------------")
}

func Part1(filepath string) (answer int) {
	grid := parseInput(filepath)
	StartingBeam := Beam{
		slices.Index(grid[0], "S"),
		0,
	}

	beams := []Beam{StartingBeam}
	bottomY := len(grid) - 1
	for len(beams) > 0 {
		// PrintGridWithBeams(grid, beams)

		beam := beams[0]
		if len(beams) < 1 {
			beams = []Beam{}
		} else {
			beams = beams[1:]
		}

		if beam.y == bottomY {
			// reached the bottom
			continue
		}

		nx := beam.x
		ny := beam.y + 1
		justSplit := false

		new_beams := []Beam{}
		if grid[ny][nx] == "^" {
			// split beam
			justSplit = true
			new_beams = append(new_beams, Beam{nx - 1, ny})
			new_beams = append(new_beams, Beam{nx + 1, ny})
		} else {
			new_beams = append(new_beams, Beam{nx, ny})
		}

	new_beam_loop:
		for _, new_beam := range new_beams {
			if !(new_beam.x >= 0 && new_beam.x < len(grid[0]) && new_beam.y >= 0 && new_beam.y < len(grid)) {
				continue
			}
			for _, existing_beam := range beams {
				if new_beam.x == existing_beam.x && new_beam.y == existing_beam.y {
					continue new_beam_loop
				}
			}
			// within grid & not existing, so we add to queue:
			if justSplit {
				justSplit = false
				answer++
			}
			beams = append(beams, new_beam)
		}
	}

	return answer
}

func processTimeline(currentTimelineBeam Beam, grid *[][]string) (new_timelines []Beam) {
	bottomY := len(*grid) - 1

	for currentTimelineBeam.y != bottomY {
		nx := currentTimelineBeam.x
		ny := currentTimelineBeam.y + 1

		candidateBeams := []Beam{}
		switch_timeline := false
		if (*grid)[ny][nx] == "^" {
			// split beam
			candidateBeams = append(candidateBeams, Beam{nx - 1, ny})
			candidateBeams = append(candidateBeams, Beam{nx + 1, ny})
		} else {
			currentTimelineBeam.x, currentTimelineBeam.y = nx, ny
			continue
		}

		for _, candidateBeam := range candidateBeams {
			if !(candidateBeam.x >= 0 && candidateBeam.x < len((*grid)[0]) && candidateBeam.y >= 0 && candidateBeam.y < len(*grid)) {
				continue
			}

			if switch_timeline {
				new_timelines = append(new_timelines, candidateBeam)
			} else {
				currentTimelineBeam.x, currentTimelineBeam.y = candidateBeam.x, candidateBeam.y
				switch_timeline = true
			}
		}
	}

	return new_timelines
}

func processTimelines(timelines []Beam, grid *[][]string, processed *int, cache *map[Beam]int) {
	for _, timeline := range timelines {
		if cachedCount, exists := (*cache)[timeline]; exists {
			*processed += cachedCount
			continue
		}

		beforeCount := *processed
		*processed++

		newTimelines := processTimeline(timeline, grid)
		processTimelines(newTimelines, grid, processed, cache)

		(*cache)[timeline] = *processed - beforeCount
	}
}

func Part2(filepath string) int {
	grid := parseInput(filepath)
	StartingBeam := Beam{
		slices.Index(grid[0], "S"),
		0,
	}

	timelines := []Beam{StartingBeam}
	var processed int

	cache := make(map[Beam]int)
	processTimelines(timelines, &grid, &processed, &cache)
	return processed
}
