package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func day17(part2 bool) Solution {
	example_filepath := GetExamplePath(17)
	input_filepath := GetInputPath(17)
	if !part2 {
		example_p1 := Part1_17(example_filepath)
		AssertExample("4,6,3,5,6,3,5,2,1,0", example_p1, 1)
		return Solution{
			"17",
			example_p1,
			Part1_17(input_filepath),
		}
	} else {
		example_p2 := Part2_17("./inputs/17/example2.txt")
		// AssertExample("117440", example_p2, 2)
		return Solution{
			"17",
			example_p2,
			"", // Part2_17(input_filepath),
		}
	}
}

type Day17 struct {
	A       int
	B       int
	C       int
	Program []int
}

func parseDay17Input(filepath string) Day17 {
	input := readInput(filepath)
	parts := strings.Split(input, "\r\n")
	re := regexp.MustCompile(`\d+`)
	a_str := re.FindString(parts[0])
	b_str := re.FindString(parts[1])
	c_str := re.FindString(parts[2])
	a, err := strconv.Atoi(a_str)
	if err != nil {
		panic("trouble parsing a as int")
	}
	b, err := strconv.Atoi(b_str)
	if err != nil {
		panic("trouble parsing b as int")
	}
	c, err := strconv.Atoi(c_str)
	if err != nil {
		panic("trouble parsing c as int")
	}

	program_strs := re.FindAllString(parts[len(parts)-1], -1)
	var program []int
	for _, prog_str := range program_strs {
		_op, err := strconv.Atoi(prog_str)
		if err != nil {
			panic("Trouble parsing op code as int")
		}
		program = append(program, _op)
	}

	return Day17{a, b, c, program}
}

func Part1_17(filepath string) string {
	input := parseDay17Input(filepath)
	return RunProgramFromInput(input)
}

func RunProgramFromInput(input Day17) string {
	register_A := input.A
	register_B := input.B
	register_C := input.C

	comboOperandMap := func(literal_op int) int {
		switch literal_op {
		case 0:
			return 0
		case 1:
			return 1
		case 2:
			return 2
		case 3:
			return 3
		case 4:
			return register_A
		case 5:
			return register_B
		case 6:
			return register_C
		case 7:
			panic("requested combo op from literal op 7")
		}
		panic(fmt.Sprintf("didnt find combo op from input literal: %v", literal_op))
	}
	var output []string

	i := 0
program_loop:
	for {
		if i >= len(input.Program)-1 {
			break
		}
		opcode := input.Program[i]
		literal_operand := input.Program[i+1]
		combo_operand := comboOperandMap(literal_operand)
		i += 2

		var ans int
		switch opcode {
		case 0:
			//adv
			ans = int(float64(register_A) / math.Pow(2, float64(combo_operand)))
			register_A = ans
		case 1:
			//bxl
			ans = register_B ^ literal_operand
			register_B = ans
		case 2:
			//bst
			ans = combo_operand % 8
			register_B = ans
		case 3:
			//jnz
			if register_A == 0 {
				break
			}
			i -= 2
			i = literal_operand
		case 4:
			//bxc
			ans = register_B ^ register_C
			register_B = ans
		case 5:
			// out
			ans = combo_operand % 8
			output = append(output, strconv.Itoa(ans))
		case 6:
			//bdv
			ans = int(float64(register_A) / math.Pow(2, float64(combo_operand)))
			register_B = ans
		case 7:
			//cdv
			ans = int(float64(register_A) / math.Pow(2, float64(combo_operand)))
			register_C = ans
		default:
			fmt.Println("EARLY EXIT?")
			break program_loop
		}
	}

	return strings.Join(output, ",")
}

func Part2_17(filepath string) string {
	input := parseDay17Input(filepath)

	var prog_strs []string
	for _, val := range input.Program {
		prog_strs = append(prog_strs, strconv.Itoa(val))
	}
	expected_output := strings.Join(prog_strs, ",")

	outputCheckChan := make(chan [2]string)
	doneChan := make(chan int, 1)

	var result string
	go func() {
		for output_check := range outputCheckChan {
			if expected_output == output_check[0] {
				result = output_check[1]
				doneChan <- 1
				return
			}
		}
	}()

	var a int
	for {
		select {
		case <-doneChan:
			return result
		default:
			go func(a_replace int) {
				tmp_input := Day17{
					a_replace,
					input.B,
					input.C,
					input.Program,
				}
				outputCheckChan <- [2]string{RunProgramFromInput(tmp_input), strconv.Itoa(a_replace)}
			}(a)
			a++
		}
	}
}
