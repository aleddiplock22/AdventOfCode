package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const MAGIC_NUM = 20201227

func main() {
	const EXAMPLE_FILEPATH = "example.txt"
	const INPUT_FILEPATH = "input.txt"

	fmt.Println("---Day 25---")

	fmt.Println("[Example] Expected: 14897079, Answer:", Solve(EXAMPLE_FILEPATH))
	fmt.Println("[Day 25 Part 1] Answer:", Solve(INPUT_FILEPATH))
}

func parseInput(filepath string) (door_public_key int, card_public_key int) {
	file, _ := os.ReadFile(filepath)
	file_content := string(file)
	input_strings := strings.Split(file_content, "\r\n")
	door_string, card_string := input_strings[0], input_strings[1]
	door_public_key, err := strconv.Atoi(door_string)
	if err != nil {
		panic("!")
	}
	card_public_key, err = strconv.Atoi(card_string)
	if err != nil {
		panic("!!")
	}
	return door_public_key, card_public_key
}

func TransformSubject(subject int, quick_start int, loop_start int, loop_size int) int {
	value := quick_start
	for i := loop_start; i < loop_size; i++ {
		value *= subject
		value = value % MAGIC_NUM
	}
	return value
}

func FindLoopSize(public_key int) int {
	const SUBJECT = 7
	i := 0
	last_val := 1
	for {
		i++
		last_val = TransformSubject(SUBJECT, last_val, i-1, i)
		if last_val == public_key {
			return i
		}
	}
}

func Solve(filepath string) int {
	door_public_key, card_public_key := parseInput(filepath)
	door_loop := FindLoopSize(door_public_key)
	card_loop := FindLoopSize(card_public_key)

	encryption_from_door := TransformSubject(card_public_key, 1, 0, door_loop)
	encryption_from_card := TransformSubject(door_public_key, 1, 0, card_loop)
	if encryption_from_card != encryption_from_door {
		panic("Got differing final encryption keys!")
	}
	return encryption_from_card
}
