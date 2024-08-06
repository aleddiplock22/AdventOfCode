package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	const EXAMPLE_FILEPATH = "example.txt"
	const INPUT_FILEPATH = "input.txt"
	const EXAMPLE2_FILEPATH = "example2.txt"

	fmt.Println("---Day 14---")

	example_p1 := Part1(EXAMPLE_FILEPATH)
	part1_ans := Part1(INPUT_FILEPATH)

	fmt.Printf("[Example P1] Expected: 165, Answer: %v\n", example_p1)
	fmt.Printf("[Part 1] Answer: %v\n", part1_ans)

	example_p2 := Part2(EXAMPLE2_FILEPATH)
	part2_ans := Part2(INPUT_FILEPATH)

	fmt.Printf("[Example P2] Expected: 208, Answer; %v\n", example_p2)
	fmt.Printf("[Part 2] Answer: %v\n", part2_ans)
}

func isMask(line string) bool {
	const MASK_STRING = "mask"
	return strings.Split(line, " ")[0] == MASK_STRING
}

type memory_alloc struct {
	value   int
	address int
}

func parseLine(line string) memory_alloc {
	mem_address, err := strconv.Atoi(strings.Split(line[4:], "]")[0])
	if err != nil {
		panic("Error parsing mem_address of line for memory alloc.")
	}
	val, err := strconv.Atoi(strings.Split(line, " = ")[1])
	if err != nil {
		panic("Error parsing value of line for memory alloc.")
	}
	return memory_alloc{
		val,
		mem_address,
	}
}

func intToBinaryString(n int) string {
	var binaryString string
	for i := 35; i >= 0; i-- {
		bit := (n >> i) & 1
		// n >> i multiplies n by 2 i times, called 'bitwise right shift' by i positions
		// & 1: This part performs a 'bitwise AND' operation between the result of the right shift operation and 1.
		// Since 1 has a binary representation of 00000001, this operation effectively isolates the least significant bit,
		// effectively extracting the value of the bit being examined.
		binaryString += fmt.Sprintf("%d", bit)
	}
	return binaryString
}

func readUnsigned36BitStringToInt(bit_string string) int {
	result := 0
	for _, n := range bit_string {
		result = (result << 1) | (int(n - '0'))
	}
	// << is bitwise left shift operation on current value of result
	// effectively shifting its binary rep one to the left (same as multiplying by 2)
	// int(n - '0') is some bs ascii value maths to get the int number value
	// | is the bitwise OR operator, which combines the results of the left side shift
	// and the ascii number into integer representation
	// dont fully understand it mind...
	return result
}

func getBitMask(line string) func(int) int {
	mask_string := strings.Split(line, " = ")[1]
	mask_runes := []rune(mask_string)
	// the dumb shit that .reverse! macro gives lol. but it'll do I suppose
	for i, j := 0, len(mask_runes)-1; i < j; i, j = i+1, j-1 {
		mask_runes[i], mask_runes[j] = mask_runes[j], mask_runes[i]
	}

	return func(value int) int {
		value_as_binary_string := intToBinaryString(value)
		value_as_binary_rune_slice := []rune(value_as_binary_string)
		for i := 0; i < len(value_as_binary_string); i++ {
			// traverse it backwards basically
			position := len(value_as_binary_string) - i - 1
			if unicode.IsDigit(mask_runes[i]) {
				value_as_binary_rune_slice[position] = mask_runes[i]
			}
		}
		bitflipped_value_as_binary_string := string(value_as_binary_rune_slice)
		bit_flipped_value := readUnsigned36BitStringToInt(bitflipped_value_as_binary_string)
		return int(bit_flipped_value)
	}
}

func Part1(filepath string) int {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	memory := make(map[int]int)
	var current_bitmask func(int) int
	var current_mem_alloc memory_alloc
	for scanner.Scan() {
		line := scanner.Text()
		if isMask(line) {
			current_bitmask = getBitMask(line)
		} else {
			current_mem_alloc = parseLine(line)
			bit_flipped_val := current_bitmask(current_mem_alloc.value)
			memory[current_mem_alloc.address] = bit_flipped_val
		}
	}

	total := 0
	for _, val := range memory {
		total += val
	}
	return total
}

func getBitMaskPart2(line string) func(int) []int {
	mask_string := strings.Split(line, " = ")[1]
	mask_runes := []rune(mask_string)
	// the dumb shit that .reverse! macro gives lol. but it'll do I suppose
	for i, j := 0, len(mask_runes)-1; i < j; i, j = i+1, j-1 {
		mask_runes[i], mask_runes[j] = mask_runes[j], mask_runes[i]
	}

	type queued_value_as_binary_rune_slice struct {
		rune_slice     []rune
		starting_point int
	}

	return func(value int) []int {
		value_as_binary_string := intToBinaryString(value)
		value_as_binary_rune_slice := []rune(value_as_binary_string)
		var answers []int
		queue := []queued_value_as_binary_rune_slice{
			{
				value_as_binary_rune_slice,
				0,
			},
		}
		for len(queue) > 0 {
			current_item := queue[0]
			if len(queue) > 1 {
				queue = queue[1:]
			} else {
				queue = []queued_value_as_binary_rune_slice{}
			}
			for i := current_item.starting_point; i < len(current_item.rune_slice); i++ {
				// traverse it current slice backwards and mask forwards ? bloody mess this
				position := len(current_item.rune_slice) - i - 1
				if mask_runes[i] == '1' {
					current_item.rune_slice[position] = mask_runes[i]
				} else if mask_runes[i] == 'X' {
					new_rune_slice := make([]rune, len(current_item.rune_slice))
					copy(new_rune_slice, current_item.rune_slice)
					if new_rune_slice[position] == '0' {
						new_rune_slice[position] = '1'
					} else {
						new_rune_slice[position] = '0'
					}
					queue = append(queue, queued_value_as_binary_rune_slice{
						new_rune_slice,
						i + 1,
					})
				}
			}
			bitflipped_value_as_binary_string := string(current_item.rune_slice)
			bit_flipped_value := readUnsigned36BitStringToInt(bitflipped_value_as_binary_string)
			answers = append(answers, int(bit_flipped_value))
		}

		return answers
	}
}

func Part2(filepath string) int {
	// difference here is we alter memory address not value
	// 0 now does nothing, 1 still changes it to 1
	// and X isn't skipped but now means we add a version with both 0 and with 1!
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	memory := make(map[int]int)
	var current_bitmask func(int) []int
	var current_mem_alloc memory_alloc
	for scanner.Scan() {
		line := scanner.Text()
		if isMask(line) {
			current_bitmask = getBitMaskPart2(line)
		} else {
			current_mem_alloc = parseLine(line)
			bit_flipped_values := current_bitmask(current_mem_alloc.address)
			for _, bit_flipped_mem_address := range bit_flipped_values {
				memory[bit_flipped_mem_address] = current_mem_alloc.value
			}
		}
	}

	total := 0
	for _, val := range memory {
		total += val
	}
	return total
}
