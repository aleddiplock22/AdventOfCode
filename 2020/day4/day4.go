package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

/*
Expected fields:
	byr (Birth Year)
	iyr (Issue Year)
	eyr (Expiration Year)
	hgt (Height)
	hcl (Hair Color)
	ecl (Eye Color)
	pid (Passport ID)
	cid (Country ID) XXX
*/

func main() {
	fmt.Println("---Day 4---")
	example_ans := Part1("example.txt")
	fmt.Printf("[Example P1] Expected: 2 Got: %d\n", example_ans)
	input_ans := Part1("input.txt")
	fmt.Printf("[Part 1] Answer: %d\n", input_ans)

	example_ans = Part2("example.txt")
	fmt.Printf("[Example P2] Answer: %d\n", example_ans)
	input_ans = Part2("input.txt")
	fmt.Printf("[Part 2] Answer: %d", input_ans) // 199 is too high
}

func readInput(filepath string) string {
	file, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Println("ERROR READING FILE!")
	}
	file_content := string(file)
	return file_content
}

func parseInputIntoPassportBlocks(filepath string) (passports []string) {
	file_content := readInput(filepath)
	passports = strings.Split(file_content, "\r\n\r\n")

	return passports
}

func passportHasReqs(passport string, reqs []string) bool {
	for _, req := range reqs {
		if !strings.Contains(passport, req) {
			return false
		}
	}
	return true
}

func Part1(filepath string) (answer int) {
	REQS := []string{"byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid"}
	passports := parseInputIntoPassportBlocks(filepath)

	answer = 0
	for _, pass := range passports {
		if passportHasReqs(pass, REQS) {
			answer++
		}
	}
	return answer
}

func passportIsValid(passport string) bool {
	/*
		byr (Birth Year) - four digits; at least 1920 and at most 2002.
		iyr (Issue Year) - four digits; at least 2010 and at most 2020.
		eyr (Expiration Year) - four digits; at least 2020 and at most 2030.
		hgt (Height) - a number followed by either cm or in:
		If cm, the number must be at least 150 and at most 193.
		If in, the number must be at least 59 and at most 76.
		hcl (Hair Color) - a # followed by exactly six characters 0-9 or a-f.
		ecl (Eye Color) - exactly one of: amb blu brn gry grn hzl oth.
		pid (Passport ID) - a nine-digit number, including leading zeroes.
		cid (Country ID) - ignored, missing or not.
	*/
	VALID_ECL := []string{"amb", "blu", "brn", "gry", "grn", "hzl", "oth"}

	//byr
	byr_idx := strings.Index(passport, "byr")
	if byr_idx == -1 {
		return false
	}
	byr_value := passport[byr_idx+4 : byr_idx+8]
	byr_value_numeric, err := strconv.Atoi(byr_value)
	if err != nil {
		return false
	}
	if !(1920 <= byr_value_numeric && byr_value_numeric <= 2002) {
		return false
	}
	if (byr_idx + 8) < (len(passport)) {
		byr_trailing_char := passport[byr_idx+8 : byr_idx+9]
		if (byr_trailing_char != " ") && (byr_trailing_char != "\n") && (byr_trailing_char != "\r") {
			return false
		}
	}

	//iyr
	iyr_idx := strings.Index(passport, "iyr")
	if iyr_idx == -1 {
		return false
	}
	iyr_value := passport[iyr_idx+4 : iyr_idx+8]
	iyr_value_numeric, err := strconv.Atoi(iyr_value)
	if err != nil {
		return false
	}
	if !(2010 <= iyr_value_numeric && iyr_value_numeric <= 2020) {
		return false
	}
	if (iyr_idx + 8) < (len(passport)) {
		iyr_trailing_char := passport[iyr_idx+8 : iyr_idx+9]
		if (iyr_trailing_char != " ") && (iyr_trailing_char != "\n") && (iyr_trailing_char != "\r") {
			return false
		}
	}

	//eyr
	eyr_idx := strings.Index(passport, "eyr")
	if eyr_idx == -1 {
		return false
	}
	eyr_value := passport[eyr_idx+4 : eyr_idx+8]
	eyr_value_numeric, err := strconv.Atoi(eyr_value)
	if err != nil {
		return false
	}
	if !(2020 <= eyr_value_numeric && eyr_value_numeric <= 2030) {
		return false
	}
	if (eyr_idx + 8) < (len(passport)) {
		eyr_trailing_char := passport[eyr_idx+8 : eyr_idx+9]
		if (eyr_trailing_char != " ") && (eyr_trailing_char != "\n") && (eyr_trailing_char != "\r") {
			return false
		}
	}

	//hgt
	hgt_index := strings.Index(passport, "hgt:")
	if hgt_index == -1 {
		return false
	}
	cm_idx := strings.Index(passport[hgt_index+4:], "cm")
	var hgt_value string
	var lbound int
	var rbound int
	if cm_idx != -1 {
		cm_idx = cm_idx + hgt_index + 4
		hgt_value = passport[hgt_index+4 : cm_idx]
		lbound = 150
		rbound = 193
	} else {
		in_idx := strings.Index(passport[hgt_index+4:], "in")
		if in_idx == -1 {
			return false
		}
		in_idx = in_idx + hgt_index + 4
		hgt_value = passport[hgt_index+4 : in_idx]
		lbound = 59
		rbound = 76
	}
	hgt_numeric, err := strconv.Atoi(hgt_value)
	if err != nil {
		return false
	}
	if !(lbound <= hgt_numeric && hgt_numeric <= rbound) {
		return false
	}

	//hcl
	hcl_index := strings.Index(passport, "hcl:")
	if hcl_index == -1 {
		return false
	}
	hcl_first_char := passport[hcl_index+4 : hcl_index+5]
	if hcl_first_char != "#" {
		return false
	}
	hcl_chars := passport[hcl_index+5 : hcl_index+11]
	for _, char := range hcl_chars {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) {
			return false
		}
	}
	if (hcl_index + 11) < (len(passport)) {
		hcl_trailing_char := passport[hcl_index+11 : hcl_index+12]
		if (hcl_trailing_char != " ") && (hcl_trailing_char != "\n") && (hcl_trailing_char != "\r") {
			return false
		}
	}

	//ecl
	ecl_index := strings.Index(passport, "ecl:")
	if ecl_index == -1 {
		return false
	}
	ecl_value := passport[ecl_index+4 : ecl_index+7]
	valid_ecl := false
	for _, _ecl := range VALID_ECL {
		if _ecl == ecl_value {
			valid_ecl = true
			break
		}
	}
	if !valid_ecl {
		return false
	}
	if (ecl_index + 7) < (len(passport)) {
		ecl_trailing_char := passport[ecl_index+7 : ecl_index+8]
		if (ecl_trailing_char != " ") && (ecl_trailing_char != "\n") && (ecl_trailing_char != "\r") {
			return false
		}
	}

	//pid
	pid_index := strings.Index(passport, "pid:")
	if pid_index == -1 {
		return false
	}
	pid_chars := passport[pid_index+4 : pid_index+13]
	for _, char := range pid_chars {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	if (pid_index + 13) < (len(passport)) {
		pid_trailing_char := passport[pid_index+13 : pid_index+14]
		if (pid_trailing_char != " ") && (pid_trailing_char != "\n") && (pid_trailing_char != "\r") {
			return false
		}
	}

	return true
}

func Part2(filepath string) (answer int) {
	passports := parseInputIntoPassportBlocks(filepath)

	answer = 0
	var passport_is_valid bool
	for _, pass := range passports {
		passport_is_valid = passportIsValid(pass)
		if passport_is_valid {
			answer++
		}
	}
	return answer
}
