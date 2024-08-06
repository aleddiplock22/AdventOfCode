package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	const EXAMPLE_FILEPATH = "example.txt"
	const INPUT_FILEPATH = "input.txt"

	const EXAMPLE_FILEPATH_P2 = "example_part2.txt"
	const INPUT_FILEPATH_P2 = "input_part2.txt"

	fmt.Println("---Day 19---")

	fmt.Println("[Example P1] Expected: 2, Answer:", Part1(EXAMPLE_FILEPATH))
	fmt.Println("[Part 1] Answer:", Part1(INPUT_FILEPATH))

	// Part 2 is incorrect. There's some edge case I'm missing I guess?
	// Honestly pretty similar solution to what Jonathon Paulson did but
	// not similar enough to know why I'm over-matching :|

	// 22/04/2024
	// Not motivated to re-write different solution just for part 2 or to spend ages debugging so we move on :(

	fmt.Println("[Example P2] Expected: 12, Answer:", Part1(EXAMPLE_FILEPATH_P2))
	fmt.Println("[Part 2] Answer:", Part1(INPUT_FILEPATH_P2))
}

type Rule struct {
	number      int
	mapper      bool
	actual_rule rune
	map_options [][]int
}

func parseInput(filepath string) (rules map[int]Rule, texts []string) {
	file, _ := os.ReadFile(filepath)
	file_content := string(file)
	parts := strings.Split(file_content, "\r\n\r\n")
	top, bottom := parts[0], parts[1]

	// Rules:
	rules = make(map[int]Rule)
	for _, rule_string := range strings.Split(top, "\r\n") {
		rule_number_str := strings.Split(rule_string, ":")[0]
		rule_number, err := strconv.Atoi(rule_number_str)
		if err != nil {
			panic("Couldn't parse as rule num as int!")
		}
		rule_to_parse := strings.Split(rule_string, ": ")[1]
		mapper := true
		actual_rule := ' '
		var map_options [][]int
		if strings.Contains(rule_to_parse, "a") {
			mapper = false
			actual_rule = 'a'
			map_options = [][]int{}
		} else if strings.Contains(rule_to_parse, "b") {
			mapper = false
			actual_rule = 'b'
			map_options = [][]int{}
		} else if strings.Contains(rule_to_parse, "|") {
			for _, map_option_string := range strings.Split(rule_to_parse, " | ") {
				var tmp_map_option []int
				for _, num := range strings.Split(map_option_string, " ") {
					tmp_num, err := strconv.Atoi(num)
					if err != nil {
						panic("Trouble parsing num in | case")
					}
					tmp_map_option = append(tmp_map_option, tmp_num)
				}
				map_options = append(map_options, tmp_map_option)
			}
		} else {
			var tmp_map_option []int
			for _, num := range strings.Split(rule_to_parse, " ") {
				tmp_num, err := strconv.Atoi(num)
				if err != nil {
					panic("Trouble parsing num in | case")
				}
				tmp_map_option = append(tmp_map_option, tmp_num)
			}
			map_options = append(map_options, tmp_map_option)
		}
		rules[rule_number] = Rule{
			rule_number,
			mapper,
			actual_rule,
			map_options,
		}
	}

	// texts
	texts = strings.Split(bottom, "\r\n")

	return rules, texts
}

func doubleCheckIsActuallyValid(text string, is_valid bool, validated_up_to int) bool {
	// validated_up_to does a ++ after last char is validated so do indeed check len not len - 1
	return is_valid && len(text) == validated_up_to
}

func textIsValid(text string, rule Rule, all_rules map[int]Rule, validating_char_idx int) (is_valid bool, validated_up_to int) {
	text_chars_to_validate := []rune(text)

	starting_validating_char_idx := validating_char_idx
	if starting_validating_char_idx == len(text) {
		// fmt.Println(text, rule, validating_char_idx)
		return true, starting_validating_char_idx
	}

	for _, map_option := range rule.map_options {
		is_valid = true
		validating_char_idx = starting_validating_char_idx
		for _, rule_number := range map_option {
			current_rule := all_rules[rule_number]
			if !current_rule.mapper {
				// actual rule
				if text_chars_to_validate[validating_char_idx] == current_rule.actual_rule {
					validating_char_idx++
				} else {
					is_valid = false
					validating_char_idx = 0
					break
				}
			} else {
				// mapper
				valid, valid_up_to := textIsValid(text, current_rule, all_rules, validating_char_idx)
				if valid {
					validating_char_idx = valid_up_to
				} else {
					// was found to be invalid, so this map option won't work out
					is_valid = false
					validating_char_idx = 0
					break
				}
			}
		}
		if is_valid {
			// no need to reset to other map_option as we've made it through!
			return is_valid, validating_char_idx
		}
	}

	return is_valid, validating_char_idx
}

func Part1(filepath string) int {
	// How many of the texts follow Rule 0?
	rules, texts := parseInput(filepath)
	rule_0 := rules[0]
	if rule_0.number != 0 {
		panic("Invalid rule 0: number was not 0!")
	}

	total_valid_texts := 0
	for _, text := range texts {
		is_valid, validated_up_to := textIsValid(text, rule_0, rules, 0)
		if doubleCheckIsActuallyValid(text, is_valid, validated_up_to) {
			total_valid_texts += 1
			// if text == "aaaabbaaaabbaaa" {
			// 	fmt.Println("Still saying valid for aaaabbaaaabbaaa")
			// }
		} else {
			continue
		}
	}

	return total_valid_texts
}
