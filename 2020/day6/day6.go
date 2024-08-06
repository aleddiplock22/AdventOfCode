package main

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	const EXAMPLE_FILEPATH = "example.txt"
	const INPUT_FILEPATH = "input.txt"
	fmt.Println("---Day 5---")
	example_ans := Part1(EXAMPLE_FILEPATH)
	input_ans := Part1(INPUT_FILEPATH)
	fmt.Printf("[Example P1] Expected: 11 Answer: %v\n", example_ans)
	fmt.Printf("[Part 1] Answer: %v\n", input_ans)

	example_ans = Part2(EXAMPLE_FILEPATH)
	fmt.Printf("[Example P2] Expected: 6 Answer: %v\n", example_ans)
	input_ans = Part2(INPUT_FILEPATH)
	fmt.Printf("[Part 2] Answer: %v", input_ans)
}

func readInput(filepath string) string {
	file, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("ERROR READING FILE!")
	}
	file_content := string(file)
	return file_content
}

func parseInputIntoBlocks(filepath string) []string {
	file_content := readInput(filepath)
	blocks := strings.Split(file_content, "\r\n\r\n")

	return blocks
}

func contains(x rune, collection []rune) bool {
	for _, char := range collection {
		if x == char {
			return true
		}
	}
	return false
}

func sumAnswersInBlock(block string) int {
	var collection []rune
	var count int = 0
	for _, char := range block {
		if unicode.IsLetter(char) {
			if contains(char, collection) {
				continue
			} else {
				collection = append(collection, char)
				count++
			}
		} else {
			continue
		}
	}

	return count
}

func sumEveryoneAnsweredInBlock(block string) int {
	individuals := strings.Split(block, "\r\n")

	if len(individuals) == 1 {
		// fmt.Printf("(Single) %v num=%v\n", individuals[0], len(individuals[0]))
		return len(individuals[0])
	}

	questions_answered := []rune(individuals[0])
	var eliminated []rune
	all_answered := len(questions_answered)
	for _, individual := range individuals[1:] {
		for _, char := range questions_answered {
			if !contains(char, []rune(individual)) {
				if contains(char, eliminated) {
					continue
				} else {
					all_answered = max(all_answered-1, 0)
					eliminated = append(eliminated, char)
				}
				// fmt.Printf("CURRENT_COUNT=%v, ELIMINATED=%v\n", all_answered, string(eliminated))
			}
		}
	}
	// fmt.Printf("(Multi) %v num=%v\n", string(questions_answered), all_answered)
	return all_answered
}

func Part1(filename string) (answer int) {
	blocks := parseInputIntoBlocks(filename)
	for _, block := range blocks {
		answer = answer + sumAnswersInBlock(block)
	}
	return answer
}

func Part2(filename string) (answer int) {
	blocks := parseInputIntoBlocks(filename)
	for _, block := range blocks {
		answer = answer + sumEveryoneAnsweredInBlock(block)
	}
	return answer
}
