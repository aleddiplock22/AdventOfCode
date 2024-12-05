package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func day05(part2 bool) Solution {
	example_path := GetExamplePath(5)
	input := GetInputPath(5)
	if !part2 {
		example_p1 := Part1_05(example_path)
		AssertExample("143", example_p1, 1)
		return Solution{
			"05",
			example_p1,
			Part1_05(input),
		}
	} else {
		example_p2 := Part2_05(example_path)
		AssertExample("123", example_p2, 2)
		return Solution{
			"05",
			example_p2,
			Part2_05(input),
		}
	}
}

type Day5 struct {
	ordering  map[int][]int // map[Y]={a,b,c} -> Y must be before any of a,b,c
	sequences [][]int
}

func readInput_05(filepath string) Day5 {
	file, _ := os.ReadFile(filepath)
	file_content := string(file)
	parts := strings.Split(file_content, "\r\n\r\n")

	mapping := map[int][]int{}
	p1 := strings.Split(parts[0], "\r\n")
	for _, mapping_str := range p1 {
		lr := strings.Split(mapping_str, "|")
		l, r := lr[0], lr[1]
		L, _ := strconv.Atoi(l)
		R, _ := strconv.Atoi(r)
		if befores, exist := mapping[R]; exist {
			befores = append(befores, L)
			mapping[R] = befores
		} else {
			mapping[R] = []int{L}
		}
	}

	sequences_as_str := strings.Split(parts[1], "\r\n")
	sequences := [][]int{}
	for _, str_sequence := range sequences_as_str {
		split_str_sequence := strings.Split(str_sequence, ",")
		seq := []int{}
		for _, str := range split_str_sequence {
			num, _ := strconv.Atoi(str)
			seq = append(seq, num)
		}
		sequences = append(sequences, seq)
	}

	return Day5{mapping, sequences}
}

func getMedian(seq []int) int {
	if len(seq)%2 == 0 {
		panic("assumed uneven seq length")
	}
	return seq[(len(seq) / 2)]
}

func isSortedSafetyManual(seq []int, ordering map[int][]int) bool {
	// is this like O(n!) lol
	for i, val := range seq {
		if befores, exist := ordering[val]; exist {
			for _, before := range befores {
				if slices.Contains(seq[i:], before) {
					return false
				}
			}
		}
	}
	return true
}

func sortedSafetyManual(seq []int, ordering map[int][]int) []int {
	// doc: This sort is not guaranteed to be stable. cmp(a, b) should return a negative number when a < b, a positive number when a > b and zero when a == b.
	slices.SortStableFunc(seq, func(a int, b int) int {
		befores, exist := ordering[a]
		if !exist {
			// no comparison
			return 0
		}
		for _, before := range befores {
			if b == before {
				// b was in a's befores, thus a > b
				return 1
			}
		}
		return -1
	})
	return seq
}

func Part1_05(filepath string) string {
	input := readInput_05(filepath)

	var total int
	for _, sequence := range input.sequences {
		if isSortedSafetyManual(sequence, input.ordering) {
			total += getMedian(sequence)
		}
	}

	return fmt.Sprintf("%v", total)
}

func Part2_05(filepath string) string {
	input := readInput_05(filepath)

	var total int
	for _, sequence := range input.sequences {
		if !isSortedSafetyManual(sequence, input.ordering) {
			// now we have to sort it
			sorted_seq := sortedSafetyManual(sequence, input.ordering)
			total += getMedian(sorted_seq)
		}
	}
	return fmt.Sprintf("%v", total)
}
