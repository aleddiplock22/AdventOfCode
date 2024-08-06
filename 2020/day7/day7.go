package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	const EXAMPLE_FILEPATH = "example.txt"
	const INPUT_FILEPATH = "input.txt"

	fmt.Println("---Day 7---")

	example_p1 := Part1(EXAMPLE_FILEPATH)
	part1_ans := Part1(INPUT_FILEPATH)

	fmt.Println("[Example P1] Expected: 4, Answer: ", example_p1)
	fmt.Println("[Part 1] Answer: ", part1_ans)

	example_p2 := Part2(EXAMPLE_FILEPATH)
	part2_ans := Part2(INPUT_FILEPATH)

	fmt.Println("[Example P2] Expected: 32, Answer: ", example_p2)
	fmt.Println("[Part 2] Answer: ", part2_ans)
}

type bag struct {
	colour   string
	contents []contained_bag
}

type contained_bag struct {
	colour string
	amount int
}

func parseLine(line string) bag {
	split1 := strings.Split(line, " bags ")
	colour, rhs := split1[0], split1[1]

	if strings.Contains(rhs, "contain no other bags.") {
		// empty!
		return bag{
			colour,
			[]contained_bag{},
		}
	}

	rhs = strings.Split(rhs, "contain ")[1]
	containment_strings := strings.Split(rhs, ", ")
	var contents []contained_bag
	for _, containment_string := range containment_strings {
		content_string := strings.Split(containment_string, " bag")[0]
		amount, err := strconv.Atoi(content_string[0:1])
		if err != nil {
			panic("Couldn't parse bag amount!")
		}
		bag_colour := content_string[2:]
		contents = append(contents, contained_bag{
			bag_colour,
			amount,
		})
	}

	return bag{
		colour,
		contents,
	}
}

func parseInput(filepath string) (bags []bag) {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		bags = append(bags, parseLine(line))
	}
	return bags
}

func canContainBagColour(bag_colour string, candidate bag, all_bags map[string]bag) bool {
	if len(candidate.contents) == 0 {
		return false
	}
	for _, inner_bag := range candidate.contents {
		if inner_bag.colour == bag_colour {
			return true
		}
	}
	for _, inner_bag := range candidate.contents {
		if canContainBagColour(bag_colour, all_bags[inner_bag.colour], all_bags) {
			return true
		}
	}
	return false
}

const SHINY_GOLD = "shiny gold"

func getBagsAndAllBags(filepath string) (bags []bag, all_bags map[string]bag) {
	bags = parseInput(filepath)
	all_bags = make(map[string]bag, len(bags))
	for _, bag_ := range bags {
		all_bags[bag_.colour] = bag_
	}
	return bags, all_bags
}

func Part1(filepath string) int {
	bags, all_bags := getBagsAndAllBags(filepath)
	total := 0
	for _, candidate_bag := range bags {
		if canContainBagColour(SHINY_GOLD, candidate_bag, all_bags) {
			total += 1
		}
	}

	return total
}

func getTotalContainedBags(candidate_bag bag, all_bags map[string]bag) int {
	total := 0
	for _, inner_bag := range candidate_bag.contents {
		total += inner_bag.amount
		total += inner_bag.amount * getTotalContainedBags(all_bags[inner_bag.colour], all_bags)
	}
	return total
}

func Part2(filepath string) int {
	_, all_bags := getBagsAndAllBags(filepath)
	// how many bags contained within our shiny gold bag?
	shiny_gold_bag := all_bags[SHINY_GOLD]
	return getTotalContainedBags(shiny_gold_bag, all_bags)
}
