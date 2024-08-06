package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*

	TODO: Finish Part 2
	Need to spot 'sea monsters' in re-constructed image.

	NOTE: CURRENT APPROACH WONT WORK BECAUSE THE GRID WILL NOT BE A 3x3 IN REAL INPUT ffs

*/

func main() {
	const EXAMPLE_FILEPATH = "example.txt"
	const INPUT_FILEPATH = "input.txt"

	fmt.Println("---Day 20---")
	// Match up the tiles! What a nightmare.

	fmt.Println("[Example P1] Expected: 20899048083289, Answer:", Part1(EXAMPLE_FILEPATH))
	fmt.Println("[Part 1] Answer:", Part1(INPUT_FILEPATH))

	fmt.Println(Part2(INPUT_FILEPATH))
}

type Tile struct {
	name string
	tile [][]int
}

func (tile *Tile) Flip() {
	// INPLACE OPERATION
	// Flip horizontally all the rows!
	for _, row := range tile.tile {
		// reverse macro
		for i, j := 0, len(row)-1; i < j; i, j = i+1, j-1 {
			row[i], row[j] = row[j], row[i]
		}
	}
}

func (tile *Tile) Rotate(times int) {
	// Inplace clockwise rotation
	if times != 1 && times != 2 && times != 3 {
		panic("Invalid amount of rotations!")
	}
	// Rows become cols essentially !
	new_tile := make([][]int, len(tile.tile))
	for i := range new_tile {
		new_tile[i] = make([]int, len(tile.tile[0]))
	}
	for y, row := range tile.tile {
		for x, val := range row {
			// switched y and x
			new_tile[x][len(tile.tile[0])-1-y] = val
		}
	}
	tile.tile = new_tile
}

func parseInput(filepath string) (tiles []Tile) {
	file, _ := os.ReadFile(filepath)
	file_content := string(file)
	tiles_blocks := strings.Split(file_content, "\r\n\r\n")
	for _, block := range tiles_blocks {
		lines := strings.Split(block, "\r\n")
		name := lines[0][5 : len(lines[0])-1]
		tile := make([][]int, len(lines[1]))
		for i := range tile {
			tile[i] = make([]int, len(lines[1:]))
		}
		for y, line := range lines[1:] {
			for x, char := range strings.Split(line, "") {
				if char == "#" {
					tile[y][x] = 1
				}
			}
		}
		tiles = append(tiles, Tile{name, tile})
	}
	return tiles
}

type Grid struct {
	/*
		Grid in the following format:
			00 01 02
			10 11 12
			20 21 22
	*/
	POS_00 Tile
	POS_01 Tile
	POS_02 Tile
	POS_10 Tile
	POS_11 Tile
	POS_12 Tile
	POS_20 Tile
	POS_21 Tile
	POS_22 Tile
}

func (tile Tile) getEdges() (edges [][]int) {
	top_edge := tile.tile[0]
	bottom_edge := tile.tile[len(tile.tile)-1]
	var right_edge []int
	var left_edge []int
	for _, row := range tile.tile {
		left_edge = append(left_edge, row[0])
		right_edge = append(right_edge, row[len(row)-1])
	}
	edges = append(edges, top_edge, right_edge, bottom_edge, left_edge)
	return edges
}

func reversedEdge(edge []int) []int {
	new_edge := make([]int, len(edge))
	for i, j := 0, len(edge)-1; i < j; i, j = i+1, j-1 {
		new_edge[i], new_edge[j] = edge[j], edge[i]
	}
	return new_edge
}

func areEqual(v1 []int, v2 []int) bool {
	for i := range v1 {
		if v1[i] != v2[i] {
			return false
		}
	}
	return true
}

func appendIfNotIn[T string | int](v *[]T, val T) {
	do_append := true
	for _, x := range *v {
		if val == x {
			do_append = false
		}
	}
	if do_append {
		*v = append(*v, val)
	}
}

func Part1(filepath string) int {
	// Started thinking along the lines of finding every permutation and then checking each
	// this was getting unwieldly and after a look at the reddit it sounds like that's not the way to go
	// so lets try something with matching edges?

	tiles := parseInput(filepath)
	matched_edges_graph := getMatchedEdgesGraph(tiles)

	// CORNERS ALL HAVE 2 EDGES
	answer := 1
	for node_name, connections := range matched_edges_graph {
		if len(connections) == 2 {
			node_val, _ := strconv.Atoi(node_name)
			answer *= node_val
		}
	}
	return answer
}

func getMatchedEdgesGraph(tiles []Tile) map[string][]string {
	edge_dict := make(map[string][][]int, len(tiles))
	for _, tile := range tiles {
		edge_dict[tile.name] = tile.getEdges()
	}

	var MATCHED_EDGES []string
	for name1, edge_set_1 := range edge_dict {
		for _, edge1 := range edge_set_1 {
			for name2, edge_set_2 := range edge_dict {
				if name1 != name2 {
					for _, edge2 := range edge_set_2 {
						if areEqual(edge1, edge2) || areEqual(edge1, reversedEdge(edge2)) {
							MATCHED_EDGES = append(MATCHED_EDGES, name1+","+name2)
						}
					}
				}
			}
		}
	}
	// why am I using strings ffs
	matched_edges_graph := make(map[string][]string)
	for _, matched_edge := range MATCHED_EDGES {
		parts := strings.Split(matched_edge, ",")
		lhs, rhs := parts[0], parts[1]
		if lhs != rhs {
			current_lhs := matched_edges_graph[lhs]
			appendIfNotIn(&current_lhs, rhs)
			matched_edges_graph[lhs] = current_lhs
			current_rhs := matched_edges_graph[rhs]
			appendIfNotIn(&current_rhs, lhs)
			matched_edges_graph[rhs] = current_rhs
		}
	}

	return matched_edges_graph
}

func contains[T string | int](val T, slice []T) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}

type DetailedEdge struct {
	node1                  string
	node2                  string
	node_1_side            string // T R B L
	node_2_side            string
	node1_side_is_reversed bool
	node2_side_is_reversed bool
}

func getEdgeInfo(tiles []Tile) []DetailedEdge {
	edge_dict := make(map[string][][]int, len(tiles))
	for _, tile := range tiles {
		edge_dict[tile.name] = tile.getEdges()
	}

	which_edge := make(map[int]string)
	which_edge[0] = "T"
	which_edge[1] = "R"
	which_edge[2] = "B" // dont ask me why Top and Bottom seem flipped from what I'd expect I cba
	which_edge[3] = "L"

	var detailed_edge_info []DetailedEdge
	var _done_edges []int // sum of them to check
	for name1, edge_set_1 := range edge_dict {
		for i1, edge1 := range edge_set_1 {
			for name2, edge_set_2 := range edge_dict {
				if name1 != name2 {
					for i2, edge2 := range edge_set_2 {
						num_edge1, _ := strconv.Atoi(name1)
						num_edge2, _ := strconv.Atoi(name2)
						which1 := which_edge[i1]
						which2 := which_edge[i2]
						if areEqual(edge1, edge2) {
							if !contains(num_edge1+num_edge2, _done_edges) {
								detailed_edge_info = append(detailed_edge_info, DetailedEdge{name1, name2, which1, which2, false, false})
								_done_edges = append(_done_edges, num_edge1+num_edge2)
							}
						} else if areEqual(reversedEdge(edge1), edge2) {
							if !contains(num_edge1+num_edge2, _done_edges) {
								detailed_edge_info = append(detailed_edge_info, DetailedEdge{name1, name2, which1, which2, true, false})
								_done_edges = append(_done_edges, num_edge1+num_edge2)
							}
						}
					}
				}
			}
		}
	}
	return detailed_edge_info
}

const NONE = "NONE"

type DetailedNode struct {
	T            string // edge+side or NONE
	R            string
	B            string
	L            string
	FLIPPED_EDGE string // T R B L NONE
}

func getDetailedMap(edge_info []DetailedEdge) map[string]DetailedNode {
	detailed_map := make(map[string]DetailedNode)
	/*
		equiv as dict/json:
			{
				"edge1": {
					"T":  "some_edge_above",
					"B": "..."
					...
				},
				...
			}
	*/

	for _, info := range edge_info {
		node, exists := detailed_map[info.node1]
		if !exists {
			node = DetailedNode{NONE, NONE, NONE, NONE, NONE}
		}
		switch info.node_1_side {
		case "T":
			node.T = info.node2 + info.node_2_side
		case "R":
			node.R = info.node2 + info.node_2_side
		case "B":
			node.B = info.node2 + info.node_2_side
		case "L":
			node.L = info.node2 + info.node_2_side
		default:
			panic("?1")
		}
		if info.node1_side_is_reversed {
			node.FLIPPED_EDGE = info.node_1_side
		}

		detailed_map[info.node1] = node

		node2, exists2 := detailed_map[info.node2]
		if !exists2 {
			node2 = DetailedNode{NONE, NONE, NONE, NONE, NONE}
		}
		switch info.node_1_side {
		case "T":
			node2.T = info.node1 + info.node_1_side
		case "R":
			node2.R = info.node1 + info.node_1_side
		case "B":
			node2.B = info.node1 + info.node_1_side
		case "L":
			node2.L = info.node1 + info.node_1_side
		default:
			panic("?2")
		}
		if info.node2_side_is_reversed {
			node2.FLIPPED_EDGE = info.node_2_side
		}

		detailed_map[info.node2] = node2
	}
	return detailed_map
}

func formGridFromDetailedEdgeInfo(edge_info []DetailedEdge, center string, tiles_map map[string]Tile) Grid {
	grid := Grid{}
	detailed_map := getDetailedMap(edge_info)

	for k, v := range detailed_map {
		fmt.Println(k, v)
	}

	detailed_center := detailed_map[center]
	center_tile := tiles_map[center]
	switch detailed_center.FLIPPED_EDGE {
	case NONE: // do nothing
	case "T":
		center_tile.Flip()
	case "R":
		center_tile.Rotate(2)
		center_tile.Flip()
	case "L":

	}

	/*


		INCOMPLETE BC ITS TOO COMPLICATED TO CODE
		cant decide how to determine what to flip and how etc. ugh

	*/

	return grid
}

func getCenter(tiles []Tile) string {
	matched_edges_graph := getMatchedEdgesGraph(tiles)
	for k, v := range matched_edges_graph {
		if len(v) == 4 {
			return k
		}
	}
	panic("COULDNT FIND CENTER")
}

func getTilesMap(tiles []Tile) map[string]Tile {
	tiles_map := make(map[string]Tile, len(tiles))
	for _, tile := range tiles {
		tiles_map[tile.name] = tile
	}
	return tiles_map
}

func Part2(filepath string) int {
	/*
			1. find correct orientation to rebuild Grid / picture
		 	2. trim the borders off all the tiles (and squish together to form that one Grid/pic)
			3. rotate/flip entire Grid until can find a sea monster [ i think ]
		 	4. locate sea monsters of the form:
								#
				#    ##    ##    ###
				#  #  #  #  #  #
			(empty space can be ANYTHING here)
			5. Count number of #s that aren't part of a sea monster
	*/

	tiles := parseInput(filepath)
	center := getCenter(tiles)
	edge_info := getEdgeInfo(tiles)
	tiles_map := getTilesMap(tiles)
	formGridFromDetailedEdgeInfo(edge_info, center, tiles_map)
	return 0
}
