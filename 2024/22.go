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
		example_p2 := Part2_22("./inputs/22/example2.txt")
		return Solution{
			"22",
			example_p2,
			Part2_22(input_filepath),
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

type MonkeyPortfolio map[[4]int]int

func PushTo4Sequence(x int, sequence *[4]int) {
	sequence[0], sequence[1], sequence[2], sequence[3] = sequence[1], sequence[2], sequence[3], x
}

func BuildMonkeyPortfolio(val int) MonkeyPortfolio {
	portfolio := make(MonkeyPortfolio)

	var prev_val int = val % 10
	secret := val
	var prev_seq [4]int
	for i := 0; i < 2000; i++ {
		mixer_1 := secret * 64
		secret = MonkeyMixAndPrune(mixer_1, secret)
		mixer_2 := secret / 32
		secret = MonkeyMixAndPrune(mixer_2, secret)
		mixer_3 := secret * 2048
		secret = MonkeyMixAndPrune(mixer_3, secret)

		value := secret % 10 // digits place
		PushTo4Sequence(value-prev_val, &prev_seq)
		prev_val = value

		if i >= 4 {
			if _, seen := portfolio[prev_seq]; seen {
				continue
			}
			portfolio[prev_seq] = value
		}
	}

	return portfolio
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

func Part2_22(filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		panic(fmt.Sprintf("Trouble opening file, err: %v.", err))
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var portfolios []MonkeyPortfolio
	for scanner.Scan() {
		line := scanner.Text()
		number, err := strconv.Atoi(line)
		if err != nil {
			panic("Trouble parsing number in input.")
		}
		portfolios = append(portfolios, BuildMonkeyPortfolio(number))
	}

	// brute force?
	var maximum_banans int
	seen := make(map[[4]int]bool)
	for i, portfolio := range portfolios {
		for sequence, value := range portfolio {
			if _, done := seen[sequence]; done {
				continue
			}
			seen[sequence] = true
			_sum := value
			for j, cmp_portfolio := range portfolios {
				if i == j {
					continue
				}
				_sum += cmp_portfolio[sequence]
			}
			maximum_banans = max(maximum_banans, _sum)
		}
	}

	return fmt.Sprintf("%d", maximum_banans)
}
