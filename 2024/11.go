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
		// example_p1 := Part2_11(example_filepath)
		return Solution{
			"11",
			"too slow... dw about it", //example_p1,
			Part2_11(input_filepath),
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
		fmt.Printf("%v / %v\n", i, blinks)
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

func Part2_11(filepath string) string {
	// raw_input := readInput(filepath)
	// nums_str := strings.Split(raw_input, " ")
	// num_rocks := HowManyRocks(nums_str, 75)  // oh no no nono nono

	// var nums []int
	// for _, num_str := range nums_str {
	// 	num, _ := strconv.Atoi(num_str)
	// 	nums = append(nums, num)
	// }

	num_rocks := 1
	return fmt.Sprintf("%d", num_rocks)
}
