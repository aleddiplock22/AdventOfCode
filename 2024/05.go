package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func day05(part2 bool) Solution {
	if !part2 {
		example_p1 := Part1_05(GetExamplePath(5))

		return Solution{
			"05",
			example_p1,
			"input part 1",
		}
	} else {
		return Solution{
			"05",
			"example part 2",
			"input part 2",
		}
	}
}

type Day5 struct {
	ordering     map[int][]int // map[Y]={a,b,c} -> Y must be before any of a,b,c
	rev_ordering map[int][]int // map[X]={d,e,f} -> X must be after any of d,e,f
	sequences    [][]int
}

func readInput_05(filepath string) Day5 {
	file, _ := os.ReadFile(filepath)
	file_content := string(file)
	parts := strings.Split(file_content, "\r\n\r\n")

	mapping := map[int][]int{}
	rev_mapping := map[int][]int{}
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
		if afters, ok := rev_mapping[L]; ok {
			afters = append(afters, R)
			rev_mapping[L] = afters
		} else {
			rev_mapping[L] = []int{R}
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

	return Day5{mapping, rev_mapping, sequences}
}

func getMedian(seq []int) int {
	if len(seq)%2 == 0 {
		panic("assumed uneven seq length")
	}
	return seq[(len(seq)/2)+1]
}

func Part1_05(filepath string) string {
	input := readInput_05(filepath)

	var total int
	for _, sequence := range input.sequences {
		if sort.SliceIsSorted(sequence, func(i, j int) bool {
			if befores, exists := input.ordering[sequence[i]]; exists {
				for _, val := range befores {
					if sequence[j] == val {
						return true
					}
				}
			} else if afters, ok := input.rev_ordering[sequence[j]]; ok {
				for _, val := range afters {
					if sequence[j] == val {
						return false
					}
				}
			}
			return true
		}) {
			// slice was sorted in our sense of it
			total += getMedian(sequence)
		} else {
			fmt.Printf("Seq: %v, was not sorted\n", sequence)
		}
	}

	return fmt.Sprintf("%v", total)
}
