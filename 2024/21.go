package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

func day21(part2 bool) Solution {
	example_filepath := GetExamplePath(21)
	input_filepath := GetInputPath(21)
	if !part2 {
		example_p1 := SolveDay21(example_filepath, false)
		return Solution{
			"21",
			example_p1,
			SolveDay21(input_filepath, false),
		}
	} else {
		return Solution{
			"21",
			"example part 2",
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

func SolveDay21(filepath string, part2 bool) string {
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

		fmt.Println(shortest_path, value)
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
