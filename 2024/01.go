package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func day01(part2 bool) Solution {
	example_path := GetExamplePath(1)
	input_path := GetInputPath(1)

	example_input := readInput_01(example_path)
	example_p1 := SolveP1(example_input)
	if example_p1 != 11 {
		panic("P1 Example wrong!")
	}

	input := readInput_01(input_path)
	p1 := SolveP1(input)

	if !part2 {
		return Solution{
			"01",
			fmt.Sprintf("%v", example_p1),
			fmt.Sprintf("%v", p1),
		}
	} else {
		return Solution{
			"01",
			SolveP2(example_input),
			SolveP2(input),
		}
	}
}

func readInput_01(filepath string) [2][]int {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var numbers [2][]int
	for scanner.Scan() {
		line := scanner.Text()
		nums := strings.Split(line, "   ") // three spaces consistent gap
		num1, err1 := strconv.Atoi(nums[0])
		num2, err2 := strconv.Atoi(nums[1])
		if err1 != nil || err2 != nil {
			panic("Trouble converting to ints when reading")
		}
		numbers[0] = append(numbers[0], num1)
		numbers[1] = append(numbers[1], num2)
	}
	return numbers
}

func SolveP1(nums [2][]int) int {
	nums1, nums2 := nums[0], nums[1]
	slices.Sort(nums1)
	slices.Sort(nums2)

	var total_diff int
	for i, n1 := range nums1 {
		n2 := nums2[i]
		total_diff += int(math.Abs(float64(n1 - n2)))
	}
	return total_diff
}

func SolveP2(nums [2][]int) string {
	nums1, nums2 := nums[0], nums[1]
	var total int
	for _, num := range nums1 {
		var count int
		for _, num2 := range nums2 {
			if num == num2 {
				count++
			}
		}
		total += num * count
	}
	return strconv.Itoa(total)
}
