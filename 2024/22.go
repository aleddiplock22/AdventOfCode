package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func day22(part2 bool) Solution {
	example_filepath := GetExamplePath(22)
	input_filepath := GetInputPath(22)
	if !part2 {
		example_p1 := Part1_22(example_filepath)
		AssertExample("37327623", example_p1, 1)
		return Solution{
			"22",
			example_p1,
			Part1_22(input_filepath),
		}
	} else {
		return Solution{
			"22",
			"example part 2",
			"input part 2",
		}
	}
}

func MonkeyMixAndPrune(val, secret int) int {
	// Mix (bitwise XOR)
	ans := val ^ secret
	// Prune ( modulo 16777216 )
	return ans % 16777216
}

func CalculateNthSecretNumber(val, n int) int {
	secret := val
	for range n {
		mixer_1 := secret * 64
		secret = MonkeyMixAndPrune(mixer_1, secret)
		mixer_2 := secret / 32
		secret = MonkeyMixAndPrune(mixer_2, secret)
		mixer_3 := secret * 2048
		secret = MonkeyMixAndPrune(mixer_3, secret)
	}
	return secret
}

func Part1_22(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		panic(fmt.Sprintf("Trouble opening file, err: %v.", err))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var answer int
	for scanner.Scan() {
		line := scanner.Text()
		number, err := strconv.Atoi(line)
		if err != nil {
			panic("Trouble parsing number in input.")
		}
		answer += CalculateNthSecretNumber(number, 2000)
	}

	return fmt.Sprintf("%d", answer)
}
