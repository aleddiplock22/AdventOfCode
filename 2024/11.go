package main

import (
	"fmt"
	"strconv"
	"strings"
)

func day11(part2 bool) Solution {
	example_filepath := GetExamplePath(11)
	input_filepath := GetInputPath(11)
	if !part2 {
		example_p1 := Part1_11(example_filepath)
		AssertExample("55312", example_p1, 1)
		return Solution{
			"11",
			example_p1,
			Part1_11(input_filepath),
		}
	} else {
		ex_for_p1_using_p2 := Part2_11(example_filepath, 25)
		AssertExample("55312", ex_for_p1_using_p2, 1)
		return Solution{
			"11",
			"N/A",
			Part2_11(input_filepath, 75),
		}
	}
}

type Num11 struct {
	value int
	next  *Num11
}

func (n *Num11) String() string {
	return strconv.Itoa(n.value)
}

func (n *Num11) IsEvenLength() bool {
	// fmt.Printf("decided %v with length %v, isEven=%v\n", n.value, len(n.String()), len(n.String())%2 == 0)
	return len(n.String())%2 == 0
}

func (n *Num11) SplitVals() (int, int) {
	_str := n.String()
	_str1, _str2 := _str[:len(_str)/2], _str[len(_str)/2:]
	num1, err := strconv.Atoi(_str1)
	if err != nil {
		panic(fmt.Sprintf("trouble reading lhs of num: %d, after splitting.", n.value))
	}
	num2, err := strconv.Atoi(_str2)
	if err != nil {
		panic(fmt.Sprintf("trouble reading rhs of num: %d, after splitting.", n.value))

	}
	// fmt.Printf("%v became %v %v\n", n.value, num1, num2)
	return num1, num2
}

func HowManyRocks(nums_str []string, blinks int) int {
	/*
		This gets too slow after a while (~44 blinks for input)
	*/

	// SETUP INITIAL CHAIN
	last_int, err := strconv.Atoi(nums_str[len(nums_str)-1])
	if err != nil {
		panic("couldnt get int from last num")
	}
	last_num := Num11{
		last_int,
		nil,
	}
	var next *Num11 = &last_num
	var first_num *Num11
	for i := len(nums_str) - 2; i >= 0; i-- {
		val, err := strconv.Atoi(nums_str[i])
		if err != nil {
			panic("trouble reading in num as int")
		}
		tmp_num := Num11{
			val,
			next,
		}
		next = &tmp_num
		if i == 0 {
			first_num = &tmp_num
			// just keep 'last' (first in list) val so we go in order from here on
		}
	}

	// RUN SIMULATION
	var num_rocks int
	for i := range blinks + 1 {
		curr := first_num
		for {
			if i == blinks {
				num_rocks++
			}
			if curr.value == 0 {
				curr.value = 1
			} else if curr.IsEvenLength() {
				num1, num2 := curr.SplitVals()
				_existing_next := curr.next
				curr.value = num1
				curr.next = &Num11{
					num2,
					_existing_next,
				}
				curr = curr.next // extra skip to avoid processing newly split num
			} else {
				curr.value *= 2024
			}
			// reached end of the chain
			if curr.next == nil {
				break
			}
			curr = curr.next
		}
	}
	return num_rocks
}

func Part1_11(filepath string) string {
	raw_input := readInput(filepath)
	nums_str := strings.Split(raw_input, " ")
	num_rocks := HowManyRocks(nums_str, 25)

	return fmt.Sprintf("%d", num_rocks)
}

func IsEvenLengthAsString(x int) bool {
	return len(strconv.Itoa(x))%2 == 0
}

func SplitIntegerInHalf(x int) (int, int) {
	_str := strconv.Itoa(x)
	_str1, _str2 := _str[:len(_str)/2], _str[len(_str)/2:]
	num1, _ := strconv.Atoi(_str1)
	num2, _ := strconv.Atoi(_str2)
	return num1, num2
}

func Part2_11(filepath string, N int) string {
	raw_input := readInput(filepath)
	nums_str := strings.Split(raw_input, " ")

	/*
		what if we had a map with
		{
			number: appearences
		}
			so we just processed number and added to it each gen!
	*/

	num_to_count := make(map[int]int)
	for _, num_str := range nums_str {
		num, _ := strconv.Atoi(num_str)
		num_to_count[num] = 1
	}

	for range N {
		// create tmp copy to avoid editing within loop!
		tmp_nums_to_process := [][2]int{}
		for num, count := range num_to_count {
			tmp_nums_to_process = append(tmp_nums_to_process, [2]int{num, count})
		}
		// fmt.Println("TO PROCESS: ", tmp_nums_to_process)
		for _, numcount := range tmp_nums_to_process {
			// fmt.Println("Processing: ", numcount)
			// fmt.Println("nums_to_count at start: ", num_to_count)
			num, count := numcount[0], numcount[1]
			num_to_count[num] -= count // remove it the number of times we process it
			if num_to_count[num] == 0 {
				delete(num_to_count, num)
			}
			if num == 0 {
				num_to_count[1] += count
			} else if IsEvenLengthAsString(num) {
				num1, num2 := SplitIntegerInHalf(num)
				num_to_count[num1] += count
				num_to_count[num2] += count
			} else {
				num_to_count[num*2024] += count
			}
		}
	}

	var num_rocks int
	for _, count := range num_to_count {
		num_rocks += count
	}

	return fmt.Sprintf("%d", num_rocks)
}
