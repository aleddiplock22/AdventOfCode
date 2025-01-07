package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func day23(part2 bool) Solution {
	example_filepath := GetExamplePath(23)
	input_filepath := GetInputPath(23)
	if !part2 {
		example_p1 := Part1_23(example_filepath)
		AssertExample("7", example_p1, 1)
		return Solution{
			"23",
			example_p1,
			Part1_23(input_filepath),
		}
	} else {
		example_p2 := Part2_23(example_filepath)
		AssertExample("co,de,ka,ta", example_p2, 2)
		return Solution{
			"23",
			example_p2,
			Part2_23(input_filepath),
		}
	}
}

type Computer struct {
	name        string
	connections []*Computer
}

// If not already in connections
func (c *Computer) AddConnection(cmp *Computer) {
	if c.name == cmp.name {
		return
	}
	for _, existing_connection := range c.connections {
		if existing_connection.name == cmp.name {
			// already in there! skip.
			return
		}
	}
	c.connections = append(c.connections, cmp)
}

func (c *Computer) IsConnectedTo(cmp *Computer) bool {
	for _, connection := range c.connections {
		if connection.name == cmp.name {
			return true
		}
	}
	return false
}

type Network map[string]*Computer

func (network Network) PrintNetwork() {
	for name, computer := range network {
		var connection_names []string
		for _, connection := range computer.connections {
			connection_names = append(connection_names, connection.name)
		}
		fmt.Printf("Computer: %v, Connections: %v\n", name, connection_names)
	}
}

func GenerateNetworkFromInputFilepath(filepath string, verbose bool) Network {
	file, err := os.Open(filepath)
	if err != nil {
		panic("Trouble opening file!")
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	network := make(Network)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "-")
		lhs, rhs := parts[0], parts[1]
		var cp_ptr_1, cp_ptr_2 *Computer
		if cp_ptr, in_network := network[lhs]; in_network {
			cp_ptr_1 = cp_ptr
		} else {
			cp_ptr_1 = &Computer{
				lhs,
				[]*Computer{},
			}
			network[lhs] = cp_ptr_1
		}
		if cp_ptr, in_network := network[rhs]; in_network {
			cp_ptr_2 = cp_ptr
		} else {
			cp_ptr_2 = &Computer{
				rhs,
				[]*Computer{},
			}
			network[rhs] = cp_ptr_2
		}

		(*cp_ptr_1).AddConnection(cp_ptr_2)
		(*cp_ptr_2).AddConnection(cp_ptr_1)
	}

	if verbose {
		fmt.Println("Generated the following network:")
		network.PrintNetwork()
	}

	return network
}

func Part1_23(filepath string) string {
	network := GenerateNetworkFromInputFilepath(filepath, false)
	three_parties := make(map[[3]string]bool)
	for computer_name, computer := range network {
		for i := 0; i < len(computer.connections); i++ {
			for j := 0; j < len(computer.connections); j++ {
				if i == j {
					continue
				}
				cmp2, cmp3 := computer.connections[i], computer.connections[j]
				if cmp2.IsConnectedTo(cmp3) {
					three := []string{computer_name, cmp2.name, cmp3.name}
					slices.Sort(three)
					three_parties[[3]string{three[0], three[1], three[2]}] = true
				}
			}
		}
	}
	var total_threes_with_a_t int
	for three_party := range three_parties {
		for _, name := range three_party {
			if name[:1] == "t" {
				total_threes_with_a_t++
				break
			}
		}
	}

	return fmt.Sprintf("%d", total_threes_with_a_t)
}

func (computer *Computer) FindLargestLanPartyForComputer() []string {
	_ = `
	START WITH MAX NUMBER (all of given computer's connections)
	FIND ALL PAIRS, if ALL CONNECTED, return this,

	then all n-1 combinations.... ETC!
	`
	combos := func(connections []*Computer) [][]*Computer {
		combinations := [][]*Computer{{}}
		for _, cmp := range connections {
			for _, combo := range combinations {
				newCombo := append([]*Computer{}, combo...)
				newCombo = append(newCombo, cmp)
				combinations = append(combinations, newCombo)
			}
		}
		return combinations
	}

	all_combos := combos(computer.connections)
	winning_combo := []string{computer.name}
outer_loop:
	for n := len(computer.connections); n > 1; n-- {
		for _, combo := range all_combos {
			if len(combo) != n {
				continue
			}
			// Get pairs
			valid := true
		inner_loop:
			for i := range combo {
				for j := range combo {
					if i == j {
						continue
					}
					if !combo[i].IsConnectedTo(combo[j]) {
						valid = false
						break inner_loop
					}
				}
			}
			if valid {
				for _, cmp := range combo {
					winning_combo = append(winning_combo, cmp.name)
				}
				break outer_loop
			}
		}
	}
	slices.Sort(winning_combo)
	return winning_combo
}

func Part2_23(filepath string) string {
	network := GenerateNetworkFromInputFilepath(filepath, false)
	// need to find largest LAN party, sort alphabetically, separate by commas
	var biggest_lan_party []string
	for _, computer := range network {
		biggest_for_computer := computer.FindLargestLanPartyForComputer()
		if len(biggest_for_computer) > len(biggest_lan_party) {
			biggest_lan_party = biggest_for_computer
		}
	}
	return strings.Join(biggest_lan_party, ",")
}
