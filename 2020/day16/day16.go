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

	fmt.Println("---Day 16---")

	example_p1 := Part1(EXAMPLE_FILEPATH)
	part1_ans := Part1(INPUT_FILEPATH)
	fmt.Println("[Example P1] Expected: 71, Answer: ", example_p1)
	fmt.Println("[Part 1] Answer: ", part1_ans)

	fmt.Println("[Part 2] Answer: ", Part2(INPUT_FILEPATH))
}

type InclusiveRange struct {
	lower int
	upper int
}

type TicketField struct {
	name             string
	inclusive_ranges []InclusiveRange
}

func inclusiveRangeFromString(range_string string) InclusiveRange {
	parts := strings.Split(range_string, "-")
	lower, err := strconv.Atoi(parts[0])
	if err != nil {
		panic("Trouble parsing lower integer in range.")
	}
	upper, err := strconv.Atoi(parts[1])
	if err != nil {
		panic("Trouble parsing upper integer in range.")
	}
	return InclusiveRange{
		lower,
		upper,
	}
}

func parseInput(filepath string) (fields []TicketField, my_ticket []int, nearby_tickets [][]int) {
	file, _ := os.ReadFile(filepath)
	parts := strings.Split(string(file), "\r\n\r\n")

	// Fields
	for _, line := range strings.Split(parts[0], "\r\n") {
		inner_parts := strings.Split(line, ": ")
		name, ranges_string := inner_parts[0], inner_parts[1]
		inner_parts = strings.Split(ranges_string, " or ")
		lh_range, rh_range := inner_parts[0], inner_parts[1]
		ranges := []InclusiveRange{
			inclusiveRangeFromString(lh_range), inclusiveRangeFromString(rh_range),
		}
		field := TicketField{
			name,
			ranges,
		}
		fields = append(fields, field)
	}

	// My Ticket
	for _, char := range strings.Split(strings.Split(parts[1], "\r\n")[1], ",") {
		num, err := strconv.Atoi(char)
		if err != nil {
			panic("Trouble parsing numer in my ticket!")
		}
		my_ticket = append(my_ticket, num)
	}

	// nearby tickets
	for _, line := range strings.Split(parts[2], "\r\n")[1:] {
		ticket := []int{}
		for _, char := range strings.Split(line, ",") {
			num, err := strconv.Atoi(char)
			if err != nil {
				panic("Trouble parsing numer in my ticket!")
			}
			ticket = append(ticket, num)
		}
		nearby_tickets = append(nearby_tickets, ticket)
	}

	return fields, my_ticket, nearby_tickets
}

func getInvalidVals(ticket []int, fields []TicketField) (invalid_vals []int) {
	// Returns False if ticket contains a value that wouldn't be
	// valid in any of the given fields
	for _, val := range ticket {
		valid := false

	field_loop:
		for _, field := range fields {
			for _, inclusive_range := range field.inclusive_ranges {
				if inclusive_range.lower <= val && val <= inclusive_range.upper {
					valid = true
					break field_loop
				}
			}
		}
		if !valid {
			invalid_vals = append(invalid_vals, val)
		}
	}
	return invalid_vals
}

func Part1(filepath string) int {
	fields, _, nearby_tickets := parseInput(filepath)
	total := 0
	for _, ticket := range nearby_tickets {
		invalid_vals_for_ticket := getInvalidVals(ticket, fields)
		if len(invalid_vals_for_ticket) > 0 {
			for _, val := range invalid_vals_for_ticket {
				total += val
			}
		}
	}
	return total
}

func getValidatedNearbyTickets(nearby_tickets [][]int, fields []TicketField) (validated_nearby_tickets [][]int) {
	for _, ticket := range nearby_tickets {
		invalid_vals_for_ticket := getInvalidVals(ticket, fields)
		if len(invalid_vals_for_ticket) == 0 {
			validated_nearby_tickets = append(validated_nearby_tickets, ticket)
		}
	}
	return validated_nearby_tickets
}

func checkIfAllNearbyTicketsValidForFieldAtPositition(field TicketField, position int, nearby_tickets [][]int) bool {
	for _, nearby_ticket := range nearby_tickets {
		val := nearby_ticket[position]
		valid := false
		for _, inclusive_range := range field.inclusive_ranges {
			if inclusive_range.lower <= val && val <= inclusive_range.upper {
				valid = true
				// fmt.Printf("Val=%v contained in field=%v\n", val, field)
				break
			}
		}
		if !valid {
			// found a ticket that's value at this position was invalid for this field
			// fmt.Printf("REJECTED %v from %v (pos: %v)\n", val, field, position)
			return false
		}
	}
	// found no such invalid tickets for value at this position for this field
	return true
}

func contains(x int, stash []int) bool {
	for _, v := range stash {
		if v == x {
			return true
		}
	}
	return false
}

func Part2(filepath string) int {
	fields, my_ticket, potentially_invalid_nearby_tickets := parseInput(filepath)
	nearby_tickets := getValidatedNearbyTickets(potentially_invalid_nearby_tickets, fields)

	// Figure out which field is which based on nearby tickets
	// Then assign my ticket accordingly to its fields, and multiply those that start with 'departure'

	// fields_map maps name of a field to its possible integer positions in ticket
	fields_map := make(map[string][]int, len(fields))
	num_potential_positions := len(my_ticket)
	for _, field := range fields {
		for i := 0; i < num_potential_positions; i++ {
			if checkIfAllNearbyTicketsValidForFieldAtPositition(field, i, nearby_tickets) {
				fields_map[field.name] = append(fields_map[field.name], i)
			}
		}
	}

	// Turns out if you just assign it to the first one it can go in, then you get the wrong answer
	// so instead we found all the possible positions each field could be
	// then assigned them more manually as below where we could be certain each time!
	left_to_do := len(fields)
	final_field_mapping := make(map[string]int, len(fields))
	for_certainty := 1
	var assigned_positions []int
	for left_to_do > 0 {
	field_loop:
		for field_name, possible_positions := range fields_map {
			if len(possible_positions) == for_certainty {
				for _, p := range possible_positions {
					if !contains(p, assigned_positions) {
						final_field_mapping[field_name] = p
						assigned_positions = append(assigned_positions, p)
						for_certainty++
						left_to_do--
						break field_loop
					}
				}
			}
		}
	}

	answer := 1
	for field_name, position := range final_field_mapping {
		if strings.Contains(field_name, "departure") {
			answer *= my_ticket[position]
		}
	}

	return answer
}
