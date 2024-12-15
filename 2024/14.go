package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func day14(part2 bool) Solution {
	example_filepath := GetExamplePath(14)
	input_filepath := GetInputPath(14)
	if !part2 {
		example_p1 := Part1_14(example_filepath, 11, 7)
		AssertExample("12", example_p1, 1)
		return Solution{
			"14",
			example_p1,
			Part1_14(input_filepath, 101, 103),
		}
	} else {
		// for visual: Part2_14(input_filepath, 101, 103, true)
		return Solution{
			"14",
			"N/A",
			Part2_14(input_filepath, 101, 103, false),
		}
	}
}

func ParseDay14(filepath string) [][4]int {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`-?\d+`)
	var output [][4]int
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindAllString(line, -1)
		var line_nums [4]int
		for i, num_str := range matches {
			num, err := strconv.Atoi(num_str)
			if err != nil {
				panic("error parsing num in line")
			}
			line_nums[i] = num
		}
		output = append(output, line_nums)
	}
	return output
}

func Part1_14(filepath string, x_max int, y_max int) string {
	robots := ParseDay14(filepath)
	// 7 tall 11 wide example
	X_HALF := x_max / 2
	Y_HALF := y_max / 2
	var ul, ur, bl, br int
	for _, robot := range robots {
		x, y, vx, vy := robot[0], robot[1], robot[2], robot[3]
		x = (x + vx*100) % x_max
		if x < 0 {
			x += x_max
		}
		y = (y + vy*100) % y_max
		if y < 0 {
			y += y_max
		}
		if x < X_HALF && y < Y_HALF {
			ul++
		} else if x < X_HALF && y > Y_HALF {
			bl++
		} else if x > X_HALF && y < Y_HALF {
			ur++
		} else if x > X_HALF && y > Y_HALF {
			br++
		}
	}
	return fmt.Sprintf("%d", ul*bl*ur*br)
}

func Part2_14(filepath string, x_max int, y_max int, visual bool) string {
	robots := ParseDay14(filepath)
	var answer int

	file, err := os.OpenFile("output.txt", os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic("trouble opening file..?")
	}
	defer file.Close()

	for i := 1; i < 10000; i++ {
		grid := [103][101]int{}
		for _, robot := range robots {
			x, y, vx, vy := robot[0], robot[1], robot[2], robot[3]
			x = (x + i*vx) % x_max
			if x < 0 {
				x += x_max
			}
			y = (y + i*vy) % y_max
			if y < 0 {
				y += y_max
			}
			grid[y][x] = 1
		}
		s := ""
		s += fmt.Sprint("\n", i, ":\n")
		for _, line := range grid {
			s += "\n"
			s += fmt.Sprintf("%v", line)
		}
		s += "\n"
		if visual {
			file.Write([]byte(s))
		} else {
			if strings.Contains(s, "1 1 1 1 1 1 1 1 1 1 1") {
				answer = i
				break
			}
		}
	}
	return fmt.Sprintf("%v", answer)
}
