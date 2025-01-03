package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

func day21(part2 bool) Solution {
	example_filepath := GetExamplePath(21)
	// input_filepath := GetInputPath(21)
	if !part2 {
		example_p1 := Part1_21(example_filepath)
		return Solution{
			"21",
			example_p1,
			"", //Part1_21(input_filepath),
		}
	} else {
		return Solution{
			"21",
			Part2_21(example_filepath),
			"input part 2",
		}
	}
}

/*
Keypad:                              Directional:
+---+---+---+                            +---+---+
| 7 | 8 | 9 |                            | ^ | A |
+---+---+---+                        +---+---+---+
| 4 | 5 | 6 |                        | < | v | > |
+---+---+---+                        +---+---+---+
| 1 | 2 | 3 |
+---+---+---+
    | 0 | A |
    +---+---+

Human controls directional keypad to send directions to Robot1 (direcitonal, starts at A)
 -> Robot1 controls directional keypad to send directions to Robot2 (directional, starts at A)
	-> Robot2 controls directional keypad to send directions to Robot3 (Keypad, starts at A)
		-> Robot3 types code on keypad using given directions


Idea:
	- BFS classic shortest path problem, treat it like a Graph problem
		- start by finding shortest path for keypad, then work back?
			>potential issue could be that we need to find all the shortest paths
			each step due to cascading effect
			(i.e. equal length at one pad != equal length lower pads)
*/

type KeyNode struct {
	val   string
	up    *KeyNode
	right *KeyNode
	down  *KeyNode
	left  *KeyNode
}

type KeyEdge struct {
	adjacent_node *KeyNode
	direction     string
}

func (node KeyNode) Adjacents() []KeyEdge {
	var adjacents []KeyEdge
	if node.up != nil {
		adjacents = append(adjacents, KeyEdge{node.up, UP})
	}
	if node.right != nil {
		adjacents = append(adjacents, KeyEdge{node.right, RIGHT})
	}
	if node.down != nil {
		adjacents = append(adjacents, KeyEdge{node.down, DOWN})
	}
	if node.left != nil {
		adjacents = append(adjacents, KeyEdge{node.left, LEFT})
	}
	// add adjacency to itself ( no move )
	adjacents = append(adjacents, KeyEdge{&node, ""})
	return adjacents
}

const UP = "^"
const LEFT = "<"
const RIGHT = ">"
const DOWN = "v"
const A_PRESS = "A"

// It's cool this all worked out. But too slow for the next part!
func Part1_21(filepath string) string {
	// Initial Setup: (this is a bit long, but idk how not to manually code given I want to use this format)
	var KeyPad_A KeyNode = KeyNode{A_PRESS, nil, nil, nil, nil}
	var KeyPad_0 KeyNode = KeyNode{"0", nil, nil, nil, nil}
	var KeyPad_1 KeyNode = KeyNode{"1", nil, nil, nil, nil}
	var KeyPad_2 KeyNode = KeyNode{"2", nil, nil, nil, nil}
	var KeyPad_3 KeyNode = KeyNode{"3", nil, nil, nil, nil}
	var KeyPad_4 KeyNode = KeyNode{"4", nil, nil, nil, nil}
	var KeyPad_5 KeyNode = KeyNode{"5", nil, nil, nil, nil}
	var KeyPad_6 KeyNode = KeyNode{"6", nil, nil, nil, nil}
	var KeyPad_7 KeyNode = KeyNode{"7", nil, nil, nil, nil}
	var KeyPad_8 KeyNode = KeyNode{"8", nil, nil, nil, nil}
	var KeyPad_9 KeyNode = KeyNode{"9", nil, nil, nil, nil}
	KeyPad_A.left, KeyPad_A.up = &KeyPad_0, &KeyPad_3
	KeyPad_0.right, KeyPad_0.up = &KeyPad_A, &KeyPad_2
	KeyPad_1.right, KeyPad_1.up = &KeyPad_2, &KeyPad_4
	KeyPad_2.up, KeyPad_2.right, KeyPad_2.down, KeyPad_2.left = &KeyPad_5, &KeyPad_3, &KeyPad_0, &KeyPad_1
	KeyPad_3.up, KeyPad_3.down, KeyPad_3.left = &KeyPad_6, &KeyPad_A, &KeyPad_2
	KeyPad_4.up, KeyPad_4.right, KeyPad_4.down = &KeyPad_7, &KeyPad_5, &KeyPad_1
	KeyPad_5.up, KeyPad_5.right, KeyPad_5.down, KeyPad_5.left = &KeyPad_8, &KeyPad_6, &KeyPad_2, &KeyPad_4
	KeyPad_6.up, KeyPad_6.down, KeyPad_6.left = &KeyPad_9, &KeyPad_3, &KeyPad_5
	KeyPad_7.right, KeyPad_7.down = &KeyPad_8, &KeyPad_4
	KeyPad_8.right, KeyPad_8.down, KeyPad_8.left = &KeyPad_9, &KeyPad_5, &KeyPad_7
	KeyPad_9.down, KeyPad_9.left = &KeyPad_6, &KeyPad_8

	var DirPad_A, DirPad_UP, DirPad_LEFT, DirPad_RIGHT, DirPad_DOWN KeyNode
	DirPad_A = KeyNode{A_PRESS, nil, nil, &DirPad_RIGHT, &DirPad_UP}
	DirPad_UP = KeyNode{UP, nil, &DirPad_A, &DirPad_DOWN, nil}
	DirPad_LEFT = KeyNode{LEFT, nil, &DirPad_DOWN, nil, nil}
	DirPad_DOWN = KeyNode{DOWN, &DirPad_UP, &DirPad_RIGHT, nil, &DirPad_LEFT}
	DirPad_RIGHT = KeyNode{RIGHT, &DirPad_A, nil, nil, &DirPad_DOWN}

	raw_input := readInput(filepath)

	var answer int
	for _, instruction := range strings.Split(raw_input, "\r\n") {
		value_str := instruction[:len(instruction)-1]
		value, err := strconv.Atoi(value_str)
		if err != nil {
			panic("trouble parsing value of instruction")
		}
		// FIRST, HUMAN LAYER:
		first_shortest_paths := ShortestPathsOnPad(&KeyPad_A, instruction)

		// SECOND, ROBOT1 -> ROBOT2
		var second_shortest_paths []string
		for _, first_path := range first_shortest_paths {
			second_shortest_paths = append(second_shortest_paths, ShortestPathsOnPad(&DirPad_A, first_path)...)
		}

		/*
			infeasible to find ALL SHORTEST PATHS on this final search.
			so instead I think going to be necessary to just find a single shortest path on this last go, which I guess is obvious
		*/

		finalChan := make(chan string)
		var wg sync.WaitGroup
		// THIRD, ROBOT2 -> ROBOT3
		for _, second_path := range second_shortest_paths {
			wg.Add(1)
			go func() {
				defer wg.Done()
				finalChan <- ShortestPathOnPad(&DirPad_A, second_path)
			}()
		}

		go func() {
			wg.Wait()
			close(finalChan)
		}()

		var shortest_path int = 999999999999
		for shortest_path_string := range finalChan {
			shortest_path = min(shortest_path, len(shortest_path_string))
		}

		answer += value * shortest_path
	}

	return fmt.Sprintf("%v", answer)
}

func ShortestPathsOnPad(starting_node *KeyNode, instruction string) []string {
	targets := strings.Split(instruction, "")

	type Path struct {
		node       *KeyNode
		path       string
		target_idx int
	}
	final_target_idx := len(targets) - 1

	var shortest_paths []string
	queue := []Path{{starting_node, "", 0}}
	seen_paths := make(map[string]bool)
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, edge := range curr.node.Adjacents() {
			new_path_str := curr.path + edge.direction
			if _, seen := seen_paths[new_path_str]; seen {
				continue
			}
			seen_paths[new_path_str] = true
			next := Path{edge.adjacent_node, new_path_str, curr.target_idx}
			if next.node.val == targets[curr.target_idx] {
				next.target_idx++
				next.path += A_PRESS
				if curr.target_idx == final_target_idx {
					// done
					if len(shortest_paths) == 0 {
						shortest_paths = append(shortest_paths, next.path)
						continue
					}
					if len(shortest_paths[0]) == len(next.path) {
						shortest_paths = append(shortest_paths, next.path)
						continue
					} else {
						return shortest_paths
					}
				}
				// jump the queue
				new_queue := []Path{next}
				new_queue = append(new_queue, queue...)
				queue = new_queue
			} else {
				queue = append(queue, next)
			}
		}
	}

	return shortest_paths
}

func ShortestPathOnPad(starting_node *KeyNode, instruction string) string {
	targets := strings.Split(instruction, "")

	type Path struct {
		node *KeyNode
		path string
	}
	queue := []Path{{starting_node, ""}}
	seen := make(map[string]bool)
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		// up right left down options

		for _, edge := range curr.node.Adjacents() {
			if _, have_seen := seen[edge.adjacent_node.val]; have_seen {
				continue
			}
			seen[edge.adjacent_node.val] = true
			next := Path{edge.adjacent_node, curr.path + edge.direction}
			if next.node.val == targets[0] {
				next.path += A_PRESS
				if len(targets) == 1 {
					// done
					return next.path
				}
				// clear queue & seen
				queue = []Path{next}
				seen = map[string]bool{}
				// move to next target
				targets = targets[1:]
				break
			} else {
				queue = append(queue, next)
			}
		}
	}

	panic("Couldn't find shortest path!")
}

// PART 2:

const LEFT_DIR_21 = 0
const DOWN_DIR_21 = 1
const RIGHT_DIR_21 = 2
const UP_DIR_21 = 3
const A_DIR_21 = 4

var GLOB_CACHE = make(map[[2]int][][]int)

// Finds the fewest number of button presses for final robot to perform instruction
func SolveKeypadConundrum(instruction string, n_layers int) int {
	/*
		This time going to be smarter about path finding, since the distances between every single
		point on the 'graph' are KNOWN, and there's finite combinations between them.
	*/

	// index 10 is A, digits are as expected
	// each key has an array mapping distance to the other keys (including itself)
	// distance is represented as a [dx, dy] int array
	var KeyPad [][][2]int // size 11, 11, 2 but kept as slice so we can use both pads as slice args later

	// (x,y) cartesian, origin bottom left where there's no Key on keypad
	keypad_locs := [11][2]int{
		{1, 0},                 // 0
		{0, 1}, {1, 1}, {2, 1}, // 1 2 3
		{0, 2}, {1, 2}, {2, 2}, // 4 5 6
		{0, 3}, {1, 3}, {2, 3}, // 7 8 9
		{2, 0}, // index 10 = Key A
	}

	KeyPad = make([][][2]int, 11)
	for key, loc := range keypad_locs {
		KeyPad[key] = make([][2]int, 11)
		for cmp_key, cmp_loc := range keypad_locs {
			dx := cmp_loc[0] - loc[0]
			dy := cmp_loc[1] - loc[1]
			KeyPad[key][cmp_key] = [2]int{dx, dy}
		}
	}

	var DirectionalPad [][][2]int
	/*
		LEFT = 0
		DOWN = 1
		RIGHT = 2
		UP = 3
		A = 4
	*/
	directional_locs := [5][2]int{
		{0, 0}, {1, 0}, {2, 0},
		{1, 1}, {2, 1},
	}
	DirectionalPad = make([][][2]int, 5)
	for key, loc := range directional_locs {
		DirectionalPad[key] = make([][2]int, 5)
		for cmp_key, cmp_loc := range directional_locs {
			dx := cmp_loc[0] - loc[0]
			dy := cmp_loc[1] - loc[1]
			DirectionalPad[key][cmp_key] = [2]int{dx, dy}
		}
	}

	// ok now we're set up, still fairly verbose but I think better than before with Node/Edge

	GenerateInstructionsFromDist := func(dist [2]int) [][]int {
		if res, exists := GLOB_CACHE[dist]; exists {
			return res
		}
		// distance (1,3) => '>^^^' (in any combination!) + A
		dx, dy := dist[0], dist[1]
		if dx == 0 && dy == 0 {
			return [][]int{{A_DIR_21}}
		}
		var new_instrs []int
		if dx < 0 {
			for range -dx {
				new_instrs = append(new_instrs, LEFT_DIR_21)
			}
		} else if dx > 0 {
			for range dx {
				new_instrs = append(new_instrs, RIGHT_DIR_21)
			}
		}
		if dy < 0 {
			for range -dy {
				new_instrs = append(new_instrs, DOWN_DIR_21)
			}
		} else if dy > 0 {
			for range dy {
				new_instrs = append(new_instrs, UP_DIR_21)
			}
		}
		if len(new_instrs) == 0 {
			panic("Invalid new instructions generated!")
		}
		perms := getPermutations(new_instrs)
		for i := range perms {
			perms[i] = append(perms[i], A_DIR_21)
		}
		GLOB_CACHE[dist] = perms
		return perms
	}

	Solver := func(start int, pad [][][2]int, instruction []int) [][]int {
		current := start

		current_instructions := make([][]int, 1)

		for _, next := range instruction {
			dist_xy := pad[current][next]
			all_sub_instrs := GenerateInstructionsFromDist(dist_xy)
			new_instructions := [][]int{}
			for _, sub_instr := range all_sub_instrs {
				for _, instructions := range current_instructions {
					tmp := append([]int{}, instructions...)
					tmp = append(tmp, sub_instr...)
					new_instructions = append(new_instructions, tmp)
				}
			}
			current_instructions = new_instructions
			current = next
		}

		return current_instructions
	}

	// initial
	var starter_instruction []int
	for _, char := range strings.Split(instruction, "") {
		if char == "A" {
			starter_instruction = append(starter_instruction, 10)
			continue
		}
		instruction_int, err := strconv.Atoi(char)
		if err != nil {
			panic("trouble parsing starter instruction int")
		}
		starter_instruction = append(starter_instruction, instruction_int)
	}
	initial_instructions := Solver(10, KeyPad, starter_instruction)

	current_instructions := initial_instructions
	// n_layers robots
	for range 2 {
		fmt.Println("LAYER len ", len(current_instructions))
		var new_curr [][]int
		for _, instr := range current_instructions {
			new_curr = append(new_curr, Solver(A_DIR_21, DirectionalPad, instr)...)
		}
		current_instructions = minAndUniquePlease(new_curr)
	}

	min_len := 999999
	for _, instr := range current_instructions {
		min_len = min(min_len, len(instr))
	}

	return min_len
}

func minAndUniquePlease(v [][]int) [][]int {
	min_len := 99999999999
	var best [][]int
	for _, option := range v {
		if len(option) < min_len {
			min_len = len(option)
			best = [][]int{} // clear it
			best = append(best, option)
		} else if len(option) == min_len {
			best = append(best, option)
		}
	}
	return uniqueifyPerms(best)
}

func Part2_21(filepath string) string {
	const INNER_ROBOT_LAYERS = 25

	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var total int
	for scanner.Scan() {
		instruction := scanner.Text()
		value_str := instruction[:len(instruction)-1]
		value, err := strconv.Atoi(value_str)
		if err != nil {
			panic("trouble parsing value of instruction")
		}
		min_instruction_length := SolveKeypadConundrum(instruction, INNER_ROBOT_LAYERS)
		total += value * min_instruction_length
	}
	return fmt.Sprintf("%d", total)
}

var GLOB_CACHE_2 = make(map[string][][]int)

func getPermutations(x []int) [][]int {
	x_str := fmt.Sprintf("%v", x)
	if res, exists := GLOB_CACHE_2[x_str]; exists {
		return res
	}
	// base case, for one int y, all perms are [[y]]
	if len(x) == 1 {
		return [][]int{x}
	}

	current := x[0]    // current
	remaining := x[1:] // remaining

	perms := getPermutations(remaining) // get perms for remaining string

	allPerms := make([][]int, 0) // array to hold all perms of the string based on perms of substring

	// for every perm in the perms of substring
	for _, perm := range perms {
		// add current char at every possible position
		for i := 0; i <= len(perm); i++ {
			newPerm := insertAt(i, current, perm)
			allPerms = append(allPerms, newPerm)
		}
	}

	// unique-ify, bc repeating digits grr
	res := uniqueifyPerms(allPerms)
	GLOB_CACHE_2[x_str] = res
	return res
}

func uniqueifyPerms(allPerms [][]int) (uniquePerms [][]int) {
	m := make(map[string]bool)
	for _, perm := range allPerms {
		perm_str := fmt.Sprintf("%v", perm)
		if _, exists := m[perm_str]; exists {
			continue
		}
		m[perm_str] = true
		uniquePerms = append(uniquePerms, perm)
	}
	return uniquePerms
}

func insertAt(idx int, x int, s []int) []int {
	new_s := []int{}
	new_s = append(new_s, s[0:idx]...)
	new_s = append(new_s, x)
	new_s = append(new_s, s[idx:]...)
	return new_s
}
