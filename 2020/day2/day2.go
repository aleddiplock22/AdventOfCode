package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("---Day 2---")

	example_ans := Part1("example.txt")
	input_ans := Part1("input.txt")

	if example_ans == 2 {
		fmt.Println("Got the example for part 1!")
	} else {
		fmt.Println("Failed the example!")
	}

	fmt.Printf("Answer for Part 1: %d\n", input_ans)

	example_ans = Part2("example.txt")
	if example_ans == 1 {
		fmt.Println("Got the example for part 2!")
	} else {
		fmt.Println("Failed the example for part 2!")
	}

	input_ans = Part2("input.txt")
	fmt.Printf("Answer for Part 2: %d\n", input_ans)

}

func Part1(filename string) (answer int) {
	file, _ := os.Open(filename)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	answer = 0
	for scanner.Scan() {
		line := scanner.Text()
		split_string := strings.Split(line, ": ")
		password := split_string[len(split_string)-1]
		lhs := split_string[0]
		lhs_arr := strings.Split(lhs, "-")
		minimum := lhs_arr[0]
		maximum := strings.Split(lhs_arr[1], " ")[0]
		char := strings.Split(lhs_arr[1], " ")[1]

		count_of_rule := strings.Count(password, char)

		min_num, _ := strconv.Atoi(minimum)
		max_num, _ := strconv.Atoi(maximum)

		if (min_num <= count_of_rule) && (count_of_rule <= max_num) {
			answer++
		}
	}
	return answer
}

func Part2(filename string) (answer int) {
	file, _ := os.Open(filename)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	answer = 0
	for scanner.Scan() {
		line := scanner.Text()
		split_string := strings.Split(line, ": ")
		password := split_string[len(split_string)-1]
		lhs := split_string[0]
		lhs_arr := strings.Split(lhs, "-")
		minimum := lhs_arr[0]
		maximum := strings.Split(lhs_arr[1], " ")[0]
		char := strings.Split(lhs_arr[1], " ")[1]

		pos_1_idx, _ := strconv.Atoi(minimum)
		pos_2_idx, _ := strconv.Atoi(maximum)

		pos_1 := password[pos_1_idx-1 : pos_1_idx]
		pos_2 := password[pos_2_idx-1 : pos_2_idx]

		if (pos_1 == char && pos_2 != char) || (pos_1 != char && pos_2 == char) {
			answer++
		}
	}
	return answer
}
