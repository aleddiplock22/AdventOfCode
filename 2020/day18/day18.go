package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
P1: Strictly Left to right 'math', no operator precedence
P2: + has precedence over *
*/

func main() {
	const INPUT_FILEPATH = "input.txt"

	const EXAMPLE_1 = "1 + 2 * 3 + 4 * 5 + 6"                           // 71
	const EXAMPLE_2 = "5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))"       // 12240
	const EXAMPLE_3 = "((2 + 4 * 9) * (6 + 9 * 8 + 6) + 6) + 2 + 4 * 2" // 13632
	const EXAMPLE_4 = "5 + (8 * 3 + 9 + 3 * 4 * 3)"                     // 437

	if replaceBracketsWithStringOfComputedValue(EXAMPLE_1, false) != "71" ||
		replaceBracketsWithStringOfComputedValue(EXAMPLE_2, false) != "12240" ||
		replaceBracketsWithStringOfComputedValue(EXAMPLE_3, false) != "13632" ||
		replaceBracketsWithStringOfComputedValue(EXAMPLE_4, false) != "437" {
		fmt.Println("FAILED ON EXAMPLES FOR PART 1 - PANIC TIME!")
		panic("Failed Part 1 Examples :(")
	}

	const EXAMPLE_P2_1 = "2 * 3 + (4 * 5)"                                 // 46
	const EXAMPLE_P2_2 = "5 * 9 * (7 * 3 * 3 + 9 * 3 + (8 + 6 * 4))"       // 669060
	const EXAMPLE_P2_3 = "((2 + 4 * 9) * (6 + 9 * 8 + 6) + 6) + 2 + 4 * 2" // 23340

	// fmt.Println(replaceBracketsWithStringOfComputedValue(EXAMPLE_P2_1, true))

	if replaceBracketsWithStringOfComputedValue(EXAMPLE_P2_1, true) != "46" ||
		replaceBracketsWithStringOfComputedValue(EXAMPLE_P2_3, true) != "23340" ||
		replaceBracketsWithStringOfComputedValue(EXAMPLE_P2_2, true) != "669060" {
		fmt.Println("FAILED ON EXAMPLES FOR PART 2 - PANIC TIME")
		panic("Failed Part 2 Examples :(")
	}

	fmt.Println("---Day 18---")

	fmt.Println("[Part 1] Answer:", ComputeHomework(INPUT_FILEPATH, false))
	fmt.Println("[Part 2] Answer:", ComputeHomework(INPUT_FILEPATH, true))
}

func computeRawEquation(equation_without_parantheses string, isPart2 bool) int {
	const MULTIPLY = "*"
	const ADD = "+"
	eqn_segments := strings.Split(equation_without_parantheses, " ")
	value, err := strconv.Atoi(eqn_segments[0])
	if err != nil {
		panic("Trouble getting initial number in computeRawEquation.")
	}

	if isPart2 {
		var vals_to_multiply []int
		// + has precedence over *
		// 1 + 2 * 3 + 4 * 5 + 6
		// becomes (1 + 2) * (3 + 4) * (5 + 6)
		for _, segment := range eqn_segments[1:] {
			switch segment {
			case MULTIPLY:
				// hit a *, store what we had and continue!
				vals_to_multiply = append(vals_to_multiply, value)
				value = 0
			case ADD:
				continue
			case "":
				continue
			default:
				num, err := strconv.Atoi(segment)
				if err != nil {
					panic("Raw equation had non opeartor, space, num!")
				}
				value += num
			}
		}
		if value == 0 {
			value++
		}
		vals_to_multiply = append(vals_to_multiply, value)
		answer := 1
		for _, v := range vals_to_multiply {
			answer *= v
		}
		return answer
	} else {
		var current_operator string
		// Part 1: Strictly Left To Right Arithmetic
		for _, segment := range eqn_segments[1:] {
			switch segment {
			case MULTIPLY:
				current_operator = MULTIPLY
			case ADD:
				current_operator = ADD
			case "":
				continue
			default:
				num, err := strconv.Atoi(segment)
				if err != nil {
					panic("Raw equation had non opeartor, space, num!")
				}
				switch current_operator {
				case MULTIPLY:
					value *= num
				case ADD:
					value += num
				default:
					panic("UNKNOWN CURRENT OPERATOR TYPE!")
				}
			}
		}
		return value
	}

}

func replaceBracketsWithStringOfComputedValue(equation string, is_Part2 bool) string {
	if !strings.Contains(equation, "(") && !strings.Contains(equation, ")") {

		return strconv.Itoa(computeRawEquation(equation, is_Part2))
	}

	var currently_building_equation []rune
	var sub_eqn []rune

	num_brackets_left := 0
	building_eqn := true
	for _, component := range equation {
		// non bracket case:
		if component != '(' && component != ')' {
			if building_eqn {
				currently_building_equation = append(currently_building_equation, component)
			} else {
				sub_eqn = append(sub_eqn, component)
			}
			continue
		}
		// bracket case:
		if component == '(' {
			if num_brackets_left == 0 {
				building_eqn = false
			}
			num_brackets_left++
		} else if component == ')' {
			if num_brackets_left == 1 {
				rebuilt_sub_eqn := string(sub_eqn[1:])
				currently_building_equation = []rune(string(currently_building_equation) + replaceBracketsWithStringOfComputedValue(rebuilt_sub_eqn, is_Part2))
				component = ' '
				sub_eqn = []rune{}
				building_eqn = true
			}
			num_brackets_left--
		}
		if building_eqn {
			currently_building_equation = append(currently_building_equation, component)
		} else {
			sub_eqn = append(sub_eqn, component)
		}
	}

	var remaining string
	if currently_building_equation[len(currently_building_equation)-1] == ')' {
		remaining = string(currently_building_equation[:len(currently_building_equation)-1])
	} else {
		remaining = string(currently_building_equation)
	}
	return replaceBracketsWithStringOfComputedValue(remaining, is_Part2)
}

func ComputeHomework(filepath string, is_Part2 bool) int {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	total := 0
	for scanner.Scan() {
		line := scanner.Text()
		answer, err := strconv.Atoi(replaceBracketsWithStringOfComputedValue(line, is_Part2))
		if err != nil {
			panic("Error converting answer to integer in Part 1!")
		}
		total += answer
	}
	return total
}
