package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("--- Day 1 ---")

	example_ans := day1("example.txt", false)
	if example_ans != 514579 {
		fmt.Println("Failed example part1, got: ", example_ans)
	}

	ans := day1("input.txt", true)
	if ans == -1 {
		fmt.Println("FAILED TO SOLVE PART 1!")
	} else {
		fmt.Println("Part 1 Answer: ", ans)
	}

	example_ans = day2("example.txt", false)
	if example_ans != 241861950 {
		fmt.Println("Failed example part1, got: ", example_ans)
	}

	ans = day2("input.txt", true)
	if ans == -1 {
		fmt.Println("FAILED TO SOLVE PART 1!")
	} else {
		fmt.Println("Part 1 Answer: ", ans)
	}
}

func read_input(filepath string) string {
	file, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("ERROR READING FILE!")
	}
	file_content := string(file)
	return file_content
}

func day1(filepath string, verbose bool) (ans int) {
	file_content := read_input(filepath)
	inputs_str := strings.Split(file_content, "\r\n")
	var inputs_num []int
	for _, str := range inputs_str {
		num, _ := strconv.Atoi(str)
		inputs_num = append(inputs_num, num)
	}
	for idx, num := range inputs_num {
		for i := 0; i < len(inputs_num) && idx != i; i++ {
			if num+inputs_num[i] == 2020 {
				if verbose {
					fmt.Printf("Found pair: (%d, %d)\n", num, inputs_num[i])
				}
				return num * inputs_num[i]
			}
		}
	}
	return -1
}

func day2(filepath string, verbose bool) (ans int) {
	file_content := read_input(filepath)
	inputs_str := strings.Split(file_content, "\r\n")
	var inputs_num []int
	for _, str := range inputs_str {
		num, _ := strconv.Atoi(str)
		inputs_num = append(inputs_num, num)
	}
	for outer_idx := range inputs_num {
		for middle_idx := 1; middle_idx < len(inputs_num) && middle_idx != outer_idx; middle_idx++ {
			for inner_idx := 2; inner_idx < len(inputs_num) && inner_idx != middle_idx && inner_idx != outer_idx; inner_idx++ {
				inner := inputs_num[outer_idx]
				middle := inputs_num[middle_idx]
				outer := inputs_num[inner_idx]
				if inner+middle+outer == 2020 {
					if verbose {
						fmt.Printf("Found trio: (%d, %d, %d)\n", inner, middle, outer)
					}
					return inputs_num[outer_idx] * inputs_num[middle_idx] * inputs_num[inner_idx]
				}
			}
		}
	}
	return -1
}
