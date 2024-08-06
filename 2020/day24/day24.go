package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	const EXAMPLE_FILEPATH = "example.txt"
	const INPUT_FILEPATH = "input.txt"

	fmt.Println("--- Day 24 ---")

	fmt.Println("[Example P1] Expected: 10, Answer:", Part1(EXAMPLE_FILEPATH))
	fmt.Println("[Part 1] Answer:", Part1(INPUT_FILEPATH))

	/*
		It's possible I got the coordinate system wrong first go around which didn't mess up part1 but did part2!
		Used JP https://github.com/jonathanpaulson/AdventOfCode/blob/master/2020/24.py to reimplement my direction vectors (and remove NORTH/SOUTH)

		>	I should've spotted there there should be 6 neighbours not 8... I was going directly up through the point of the pointed hexagon, which isn't a neighbour!

		Anyway despite this I still wasnt getting the right answer and mine was too slow.
		> 	Fixed slow by borrowing JP's idea of just tracking what was 'ON' (black). I think this approach to cellular automata of just tracking the ones that CAN change is a great optimisation
		>	Fixed why it was wrong by realising I'd simply typed in the coordinates for my HexCoord objects in the wrong order in a few places including the serialise/deserialise!!!!! Very silly
	*/

	fmt.Println("[Example P2] Expected: 2208, Answer:", Part2(EXAMPLE_FILEPATH))
	fmt.Println("[Part 2] Answer:", Part2(INPUT_FILEPATH))
}

func parseInput(filepath string) (instructions [][]HexCoord) {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		current_instructions := []HexCoord{}
		done_next := false
		line_chars := strings.Split(line, "")
		var next_char string

		// var INSTRUCTION_STRING string
		for i, char := range line_chars {
			if done_next {
				done_next = false
				continue
			}
			if i == len(line_chars)-1 {
				next_char = "..."
			} else {
				next_char = line_chars[i+1]
			}
			switch char {
			case "n":
				if next_char == "e" {
					// INSTRUCTION_STRING += "NorthEast,"
					done_next = true
					current_instructions = append(current_instructions, NORTH_EAST)
				} else if next_char == "w" {
					// INSTRUCTION_STRING += "NorthWest,"
					done_next = true
					current_instructions = append(current_instructions, NORTH_WEST)
				}
			case "s":
				if next_char == "e" {
					// INSTRUCTION_STRING += "SouthEast,"
					done_next = true
					current_instructions = append(current_instructions, SOUTH_EAST)
				} else if next_char == "w" {
					// INSTRUCTION_STRING += "SouthWest,"
					done_next = true
					current_instructions = append(current_instructions, SOUTH_WEST)
				}
			case "e":
				// INSTRUCTION_STRING += "East,"
				current_instructions = append(current_instructions, EAST)
			case "w":
				// INSTRUCTION_STRING += "West,"
				current_instructions = append(current_instructions, WEST)
			default:
				panic("Unknown instruction recieved!")
			}
		}
		// fmt.Println(INSTRUCTION_STRING)
		instructions = append(instructions, current_instructions)
	}
	return instructions
}

func Serialise(hex HexCoord) string {
	intSlice := []int{hex.r, hex.s, hex.q}
	// Convert int slice to string slice
	stringSlice := make([]string, len(intSlice))
	for i, num := range intSlice {
		stringSlice[i] = fmt.Sprintf("%d", num)
	}
	// Marshal the string slice to JSON
	jsonData, err := json.Marshal(stringSlice)
	if err != nil {
		panic(fmt.Sprint("error:", err))
	}
	return string(jsonData)
}

func Deserialise(str string) HexCoord {
	byte_slice_of_string := []byte(str)
	var value []string
	json.Unmarshal(byte_slice_of_string, &value)
	r, _ := strconv.Atoi(value[0])
	s, _ := strconv.Atoi(value[1])
	q, _ := strconv.Atoi(value[2])
	return HexCoord{
		r,
		s,
		q,
	}
}

func defineGridFromInstructions(instructions [][]HexCoord) (grid map[string]bool) {
	/*
		Have a dictionary (map) representing grid,
		{
			hex_coordinate_1[HASHED] : {
				is_flipped:bool (true=[black], false=[white])
			}
			...
		}
	*/
	grid = make(map[string]bool)
	for _, instruction_set := range instructions {
		tile := HexCoord{0, 0, 0} // start in center
		for _, step := range instruction_set {
			tile.Add(step)
		}
		serialised_tile := Serialise(tile)
		is_black, ok := grid[serialised_tile]
		if ok {
			// flip it from current
			grid[serialised_tile] = !is_black
		} else {
			// flip it to black
			grid[serialised_tile] = true
		}
	}
	return grid
}

func countBlackTiles(grid map[string]bool) int {
	count := 0
	for _, is_black := range grid {
		if is_black {
			count += 1
		}
	}
	return count
}

func NumBlacksAfterSimulation(grid map[string]bool) int {
	BLACK_COORDS := make(map[string]bool)
	for serialised_coord, is_black := range grid {
		if is_black {
			BLACK_COORDS[serialised_coord] = is_black
		}
	}

	for i := 0; i < 100; i++ {
		NEW_BLACKS := make(map[string]bool)
		to_check := []string{}
		for serialised_hex := range BLACK_COORDS {
			to_check = append(to_check, serialised_hex)
			hex := Deserialise(serialised_hex)
			for _, nbr := range hex.GetNeighbourCoords() {
				to_check = append(to_check, Serialise(nbr))
			}
		}

		for _, serialised_hex := range to_check {
			num_blacks := 0
			hex := Deserialise(serialised_hex)
			for _, nbr := range hex.GetNeighbourCoords() {
				_, exists := BLACK_COORDS[Serialise(nbr)]
				if exists {
					num_blacks++
				}
			}
			_, is_black := BLACK_COORDS[serialised_hex]
			if is_black && !(num_blacks == 0 || num_blacks > 2) {
				NEW_BLACKS[serialised_hex] = true
			}
			if !is_black && num_blacks == 2 {
				NEW_BLACKS[serialised_hex] = true
			}
		}
		BLACK_COORDS = NEW_BLACKS
	}
	return len(BLACK_COORDS)
}

func Part1(filepath string) int {
	instructions := parseInput(filepath)
	grid := defineGridFromInstructions(instructions)
	return countBlackTiles(grid)
}

func Part2(filepath string) int {
	instructions := parseInput(filepath)
	grid := defineGridFromInstructions(instructions)
	return NumBlacksAfterSimulation(grid)
}
