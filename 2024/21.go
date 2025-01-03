package main

import (
	"fmt"
	"strconv"
	"strings"
)

func day21(part2 bool) Solution {
	example_filepath := GetExamplePath(21)
	// input_filepath := GetInputPath(21)
	if !part2 {
		example_p1 := SolveKeypadConundrum(example_filepath, 2)
		return Solution{
			"21",
			example_p1,
			"", //Part1_21(input_filepath),
		}
	} else {
		return Solution{
			"21",
			"", // SolveKeypadConundrum(example_filepath, 25),
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
const A_PRESS_KEYPAD = "A"
const A_PRESS_DIRPAD = "B"

// It's cool this all worked out. But too slow for the next part!
func SolveKeypadConundrum(filepath string, depth int) string {
	// Initial Setup: (this is a bit long, but idk how not to manually code given I want to use this format)
	var KeyPad_A KeyNode = KeyNode{A_PRESS_KEYPAD, nil, nil, nil, nil}
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
	DirPad_A = KeyNode{A_PRESS_DIRPAD, nil, nil, &DirPad_RIGHT, &DirPad_UP}
	DirPad_UP = KeyNode{UP, nil, &DirPad_A, &DirPad_DOWN, nil}
	DirPad_LEFT = KeyNode{LEFT, nil, &DirPad_DOWN, nil, nil}
	DirPad_DOWN = KeyNode{DOWN, &DirPad_UP, &DirPad_RIGHT, nil, &DirPad_LEFT}
	DirPad_RIGHT = KeyNode{RIGHT, &DirPad_A, nil, nil, &DirPad_DOWN}

	keypads := [11]*KeyNode{
		&KeyPad_A, &KeyPad_0, &KeyPad_1, &KeyPad_2, &KeyPad_3, &KeyPad_4, &KeyPad_5, &KeyPad_6, &KeyPad_7, &KeyPad_8, &KeyPad_9,
	}
	dirpads := [5]*KeyNode{
		&DirPad_A, &DirPad_UP, &DirPad_LEFT, &DirPad_DOWN, &DirPad_RIGHT,
	}

	pair_seqs := make(map[[2]string][]string)
	pair_lengths := make(map[[2]string]int)

	for _, key := range keypads {
		for _, key2 := range keypads {
			shortest_paths := ShortestPathsOnPad(key, key2)
			_key := [2]string{key.val, key2.val}
			pair_lengths[_key] = len(shortest_paths[0])
			pair_seqs[_key] = shortest_paths
		}
	}

	for _, key := range dirpads {
		for _, key2 := range dirpads {
			shortest_paths := ShortestPathsOnPad(key, key2)
			_key := [2]string{key.val, key2.val}
			pair_lengths[_key] = len(shortest_paths[0])
			pair_seqs[_key] = shortest_paths
		}
	}

	for _, key := range dirpads {
		for _, key2 := range dirpads {
			pair_lengths[[2]string{key.val, key2.val}] = len(ShortestPathsOnPad(key, key2)[0])
		}
	}

	// fmt.Println(pair_lengths)
	// fmt.Println("")

	// fmt.Println(pair_seqs)

	ALL_KEYS_MAP := map[string]*KeyNode{
		A_PRESS_KEYPAD: &KeyPad_A,
		"0":            &KeyPad_0,
		"1":            &KeyPad_1,
		"2":            &KeyPad_2,
		"3":            &KeyPad_3,
		"4":            &KeyPad_4,
		"5":            &KeyPad_5,
		"6":            &KeyPad_6,
		"7":            &KeyPad_7,
		"8":            &KeyPad_8,
		"9":            &KeyPad_9,
		A_PRESS_DIRPAD: &DirPad_A,
		UP:             &DirPad_UP,
		LEFT:           &DirPad_LEFT,
		DOWN:           &DirPad_DOWN,
		RIGHT:          &DirPad_RIGHT,
	}

	raw_input := readInput(filepath)

	cache := make(map[string]int)

	var answer int
	for _, instruction := range strings.Split(raw_input, "\r\n") {
		fmt.Println("DOING INSTRUCTION", instruction)
		value_str := instruction[:len(instruction)-1]
		value, err := strconv.Atoi(value_str)
		if err != nil {
			panic("trouble parsing value of instruction")
		}

		min_len := computeBestLengthOfSequence(instruction, depth, ALL_KEYS_MAP, pair_lengths, pair_seqs, cache)
		fmt.Println("MIN_LEN:", min_len)
		answer += value * min_len
	}

	return fmt.Sprintf("%v", answer)
}

func computeBestLengthOfSequence(sequence string, depth int, key_map map[string]*KeyNode, pair_lengths map[[2]string]int, pair_seqs map[[2]string][]string, cache map[string]int) int {
	if cached, exists := cache[sequence]; exists {
		return cached
	}

	A_seq := "A" + sequence
	if depth == 1 {
		total := 0
		for _, char := range strings.Split(A_seq, "") {
			for _, char2 := range strings.Split(sequence, "") {
				total += pair_lengths[[2]string{char, char2}]
			}
		}
		return total
	}

	// depth > 1 , we recurse
	length := 0
	for _, char := range strings.Split(A_seq, "") {
		for _, char2 := range strings.Split(sequence, "") {
			// min(compute_length(subseq, depth - 1) for subseq in pair_seqs[(x, y)])
			optimal := 9999999999999999
			if _, in_there := pair_seqs[[2]string{char, char2}]; !in_there {
				panic("What")
			}
			for _, subseq := range pair_seqs[[2]string{char, char2}] {
				optimal = min(optimal, computeBestLengthOfSequence(subseq, depth-1, key_map, pair_lengths, pair_seqs, cache))
			}
			length += optimal
		}
	}
	return length
}

type KeyPath struct {
	node       *KeyNode
	path       string
	target_idx int
}

func ShortestPathsOnPad(starting_node *KeyNode, target_node *KeyNode) []string {
	var shortest_paths []string
	queue := []KeyPath{{starting_node, "", 0}}
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
			next := KeyPath{edge.adjacent_node, new_path_str, curr.target_idx}
			if next.node.val == target_node.val {
				next.target_idx++
				next.path += A_PRESS_DIRPAD
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
			} else {
				queue = append(queue, next)
			}
		}
	}

	return shortest_paths
}
