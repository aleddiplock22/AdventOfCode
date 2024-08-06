package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

/*

	TODO: Finish Part 2

	there's something I should be able to do ( that is done in the python notebook )
	to make it to where my ing::possible_allergens map has an ing with just one ??

	once that's done can like eliminate one at a time or something. Just like with the Tickets in day16



*/

func main() {
	const EXAMPLE_FILEPATH = "example.txt"
	const INPUT_FILEPATH = "input.txt"
	fmt.Println("---Day 21---")

	fmt.Println("[Example P1] Expected: 5, Answer: ", Part1(EXAMPLE_FILEPATH))
	fmt.Println("[Part 1] Answer:", Part1(INPUT_FILEPATH))

	fmt.Println("[Example P2] Expected: mxmxvkd,sqjhc,fvjkl | Answer: ", Part2(EXAMPLE_FILEPATH))
	fmt.Println("[Part 2] Answer: ", Part2(INPUT_FILEPATH))
}

func parseInput(filepath string) (ingredient_entries [][]string, allergen_entries [][]string) {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " (contains ")
		ingredient_side, allergen_side := parts[0], parts[1][:len(parts[1])-1]
		ingredient_entries = append(ingredient_entries, strings.Split(ingredient_side, " "))
		allergen_entries = append(allergen_entries, strings.Split(allergen_side, ", "))
	}

	return ingredient_entries, allergen_entries
}

func contains[T string | int](slice []T, item T) bool {
	for _, val := range slice {
		if val == item {
			return true
		}
	}
	return false
}

func appendIfNotIn(slice *[]string, item string) {
	if !contains(*slice, item) {
		*slice = append(*slice, item)
	}
}

type Set struct {
	contents []string
}

func (set Set) Has(item string) bool {
	return contains(set.contents, item)
}

func (set *Set) Add(item string) {
	if !set.Has(item) {
		set.contents = append(set.contents, item)
	}
}

func (set Set) IsEmpty() bool {
	return len(set.contents) == 0
}

func (set *Set) Remove(item string) {
	if set.Has(item) {
		var new_vals []string
		for _, val := range set.contents {
			if val != item {
				new_vals = append(new_vals, val)
			}
		}
		set.contents = new_vals
	}
}

func getIngsToPossAlgnsMap(ingredient_entries [][]string, allergen_entries [][]string) map[string][]string {
	ing_to_poss_algns := make(map[string][]string)
	for i, ing_entries := range ingredient_entries {
		algn_entries := allergen_entries[i]
		for _, ing := range ing_entries {
			poss_algns, exists := ing_to_poss_algns[ing]
			if !exists {
				poss_algns = []string{}
			}
			for _, algn := range algn_entries {
				appendIfNotIn(&poss_algns, algn)
			}
			ing_to_poss_algns[ing] = poss_algns
		}
	}
	return ing_to_poss_algns
}

func findDefiniteNonAllergenIngrdients(ingredient_entries [][]string, allergen_entries [][]string) []string {
	// ALLERGEN FOUND IN *EXACTLY* ONE INGREDIENT
	// ALLERGENS ARE NOT ALWAYS MARKED

	ing_to_poss_algns := getIngsToPossAlgnsMap(ingredient_entries, allergen_entries)

	var definitely_dont_contain_allergens_list []string

	for ing, poss_algns := range ing_to_poss_algns {
		// allergens to eliminate the possibility of ('set'):
		algns_to_eliminate := Set{poss_algns}
		for i, ing_entries := range ingredient_entries {
			algn_entries := allergen_entries[i]
			if !contains(ing_entries, ing) {
				for _, algn := range algn_entries {
					if algns_to_eliminate.Has(algn) && contains(algn_entries, algn) {
						algns_to_eliminate.Remove(algn)
					}
				}
			}
		}
		if algns_to_eliminate.IsEmpty() {
			definitely_dont_contain_allergens_list = append(definitely_dont_contain_allergens_list, ing)
		}
	}

	// fmt.Println(definitely_dont_contain_allergens_list)
	return definitely_dont_contain_allergens_list
}

func countIngsApperances(non_allergen_list []string, ingredient_entries [][]string) int {
	count := 0
	for _, non_allergen := range non_allergen_list {
		for _, food := range ingredient_entries {
			for _, ing := range food {
				if non_allergen == ing {
					count += 1
				}
			}
		}
	}
	return count
}
func Part1(filepath string) int {
	ingredient_entries, allergen_entries := parseInput(filepath)
	definitely_dont_contain_allergens_list := findDefiniteNonAllergenIngrdients(ingredient_entries, allergen_entries)
	return countIngsApperances(definitely_dont_contain_allergens_list, ingredient_entries)
}

func updateAndImprovedIngsToPossAlgnsMap(ing_to_poss_algns map[string][]string, definitely_dont_contain_allergens_list []string, recipes_with_allgn_map map[string][]int, ingredient_entries [][]string) map[string]Set {

	/*
		for ingr, possible in possible_allers.items():
			impossible = set()

			for aller in possible:
				if any(ingr not in recipes[i] for i in recipes_with[aller]):
					impossible.add(aller)

			possible -= impossible
	*/

	improved_map := make(map[string]Set, len(ing_to_poss_algns))
	for ing, poss_algns := range ing_to_poss_algns {

		impossible_algns := Set{}
		for _, algn := range poss_algns {
			impossible := false
			for _, recipe_id := range recipes_with_allgn_map[algn] {
				if !contains(ingredient_entries[recipe_id], ing) {
					impossible = true
					break
				}
			}
			if impossible {
				impossible_algns.Add(algn)
			}
		}

		updated_poss_algns := Set{poss_algns}
		for _, algn := range updated_poss_algns.contents {
			if impossible_algns.Has(algn) {
				updated_poss_algns.Remove(algn)
			}
		}
		if !contains(definitely_dont_contain_allergens_list, ing) {
			improved_map[ing] = updated_poss_algns
		}
	}
	return improved_map
}

func getRecipesWithAllergenMap(allergen_entries [][]string) map[string][]int {
	recipes_with := make(map[string][]int)
	for i, allgn_entries := range allergen_entries {
		for _, allgn := range allgn_entries {
			current, ok := recipes_with[allgn]
			if !ok {
				current = []int{}
			}
			current = append(current, i)
			recipes_with[allgn] = current
		}
	}
	return recipes_with
}

type Pair struct {
	Key   string
	Value string
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value } // this works with strings for alphabetical comparison!
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func getCanonicalDangerousIngredientList(ingredient_entries [][]string, allergen_entries [][]string) string {
	// can we use this definitely dont contain allergen list to determine which ingredients definitely line up with which allergens...?

	/* in the example

		mxmxvkd kfcds sqjhc nhms (contains dairy, fish)
		trh fvjkl sbzzf mxmxvkd (contains dairy)
		sqjhc fvjkl (contains soy)
		sqjhc mxmxvkd sbzzf (contains fish)

	non-allergen (defo) - kfcds, nhms, sbzzf, or trh

	thus deduce:

		mxmxvkd contains dairy.
		sqjhc contains fish.
		fvjkl contains soy
	*/

	ings_to_poss_algns := getIngsToPossAlgnsMap(ingredient_entries, allergen_entries)
	definitely_dont_contain_allergens_list := findDefiniteNonAllergenIngrdients(ingredient_entries, allergen_entries)
	recipes_with_algn_map := getRecipesWithAllergenMap(allergen_entries)

	ings_to_poss_algns_updated := updateAndImprovedIngsToPossAlgnsMap(ings_to_poss_algns, definitely_dont_contain_allergens_list, recipes_with_algn_map, ingredient_entries)

	removeFromAllExcept := func(me_map map[string]Set, val string, exception string) {
		for k, set := range me_map {
			if k != exception {
				set.Remove(val)
				me_map[k] = set
			}
		}
	}

	update_map := func(me_map map[string]Set) {
		change := true
		done := []string{}
		for change {
			change = false
			for k, set := range me_map {
				if contains(done, k) {
					continue
				}
				if len(set.contents) == 1 {
					value_in_set := set.contents[0]
					removeFromAllExcept(me_map, value_in_set, k)
					change = true
					done = append(done, k)
					break
				}
			}
		}
	}

	update_map(ings_to_poss_algns_updated)

	// fmt.Println("UPDATED:")
	// for k, v := range ings_to_poss_algns_updated {
	// 	fmt.Println(k, v)
	// }

	sortByAllgn := func(ingsToAlgns map[string]Set) PairList {
		pl := make(PairList, len(ingsToAlgns))
		i := 0
		for k, v := range ingsToAlgns {
			pl[i] = Pair{k, v.contents[0]}
			i++
		}
		sort.Sort(pl)
		return pl
	}

	answer := ""
	for _, x := range sortByAllgn(ings_to_poss_algns_updated) {
		answer = answer + "," + x.Key
	}
	answer = answer[1:]

	return answer
}

func Part2(filepath string) string {
	ingredient_entries, allergen_entries := parseInput(filepath)
	return getCanonicalDangerousIngredientList(ingredient_entries, allergen_entries)
}
