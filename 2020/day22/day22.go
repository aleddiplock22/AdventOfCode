package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"os"
	"strconv"
	"strings"
	"time"
)

var USE_CACHE_IN_PART2 bool

func main() {
	const EXAMPLE_FILEPATH = "example.txt"
	const INPUT_FILEPATH = "input.txt"
	const INFINITE_EXAMPLE = "infinite_example.txt"

	fmt.Println("---Day 22---")

	fmt.Println("[Example P1] Expected: 306, Answer:", Part1(EXAMPLE_FILEPATH))
	fmt.Println("[Part 1] Answer:", Part1(INPUT_FILEPATH))

	fmt.Println("------------------------")
	fmt.Println("[Example P2] Expected: 291, Answer:", Part2(EXAMPLE_FILEPATH))
	fmt.Println("[Part 2] Anwer:", Part2(INPUT_FILEPATH))
	fmt.Println("infinite; ", Part2(INFINITE_EXAMPLE))

	fmt.Println("------------------------")
	USE_CACHE_IN_PART2 = true
	t := time.Now()
	for i := 0; i < 10; i++ {
		Part2(INPUT_FILEPATH)
	}
	fmt.Println("CACHED TOOK: ", time.Since(t)) // barely saves time lol, but clearly a game or two were quicker using a cached val :)
	USE_CACHE_IN_PART2 = false
	t2 := time.Now()
	for i := 0; i < 10; i++ {
		Part2(INPUT_FILEPATH)
	}
	fmt.Println("UNCACHED TOOK: ", time.Since(t2))
}

type Deck struct {
	cards []int
}

func (deck *Deck) RemoveTopCard() int {
	top := deck.cards[0]
	if len(deck.cards) > 1 {
		deck.cards = deck.cards[1:]
	} else {
		deck.cards = []int{}
	}
	return top
}

func (deck Deck) IsEmpty() bool { return len(deck.cards) == 0 }

func (deck *Deck) AddCard(card int) {
	deck.cards = append(deck.cards, card)
}

func (deck *Deck) AddTwoToBack(card1 int, card2 int) {
	deck.AddCard(card1)
	deck.AddCard(card2)
}

func (deck Deck) ComputeScore() int {
	score := 0
	for i, card := range deck.cards {
		score += (len(deck.cards) - i) * card
	}
	return score
}

func (deck Deck) Len() int {
	return len(deck.cards)
}

func (deck Deck) CopyDeck() Deck {
	new_cards := []int{}
	new_cards = append(new_cards, deck.cards...)
	return Deck{new_cards}
}

func (deck Deck) CopyDeckUptoNth(N int) Deck {
	new_cards := []int{}
	new_cards = append(new_cards, deck.cards[:N]...)

	return Deck{new_cards}
}

func PlayOnce(deck1 *Deck, deck2 *Deck) {
	card1 := deck1.RemoveTopCard()
	card2 := deck2.RemoveTopCard()
	if card1 == card2 {
		panic("Hadn't prepared for equal valued cards!?!??")
	} else if card1 > card2 {
		deck1.AddTwoToBack(card1, card2)
	} else {
		// card2 > card1
		deck2.AddTwoToBack(card2, card1)
	}
}

func Play(deck1 *Deck, deck2 *Deck, recursive bool) int {
	set := SeenSet{}
	seenPtr := &set
	cache := make(map[uint32]int)
	for !deck1.IsEmpty() && !deck2.IsEmpty() {
		if recursive {
			// Part 2
			if seenBefore(*deck1, *deck2, seenPtr) {
				fmt.Println("EARLY PLAYER 1 WIN")
				return deck1.ComputeScore()
			}
			copy_deck1 := deck1.CopyDeck()
			copy_deck2 := deck2.CopyDeck()
			seenPtr.addToDeck(1, copy_deck1)
			seenPtr.addToDeck(2, copy_deck2)
			PlayOnceRecursively(deck1, deck2, &cache)
		} else {
			// Part 1
			PlayOnce(deck1, deck2)
		}
	}
	if deck1.IsEmpty() {
		return deck2.ComputeScore()
	} else {
		if !deck2.IsEmpty() {
			panic("No one won??")
		}
		return deck1.ComputeScore()
	}
}

type SeenSet struct {
	deck1s []Deck
	deck2s []Deck
}

func Hash(deck1 Deck, deck2 Deck) uint32 {
	// // hopefully this is unique?
	// val := 1
	// for i, v := range deck1.cards {
	// 	val *= (i * v)
	// }
	// for j, v := range deck2.cards {
	// 	val += (j + v) * 3
	// }

	//above was not unique lol

	hashIntArrays := func(arrs ...[]int) uint32 {
		h := fnv.New32a()
		for _, arr := range arrs {
			for _, i := range arr {
				// Convert int to byte slice
				bs := []byte(fmt.Sprintf("%d", i))
				h.Write(bs)
			}
		}
		return h.Sum32()
	}

	return hashIntArrays(deck1.cards, deck2.cards)
}

func (seen *SeenSet) addToDeck(deck_num int, deck_to_add Deck) {
	deck1s := seen.deck1s
	deck2s := seen.deck2s
	if deck_num == 1 {
		deck1s = append(deck1s, deck_to_add)
	} else if deck_num == 2 {
		deck2s = append(deck2s, deck_to_add)
	}
	seen.deck1s = deck1s
	seen.deck2s = deck2s
}

func PlayRecursively(deck1 *Deck, deck2 *Deck, cache *map[uint32]int) int {
	set := SeenSet{}
	seenPtr := &set

	initial_hash := Hash(*deck1, *deck2)
	if USE_CACHE_IN_PART2 {
		result, exists := (*cache)[initial_hash]
		if exists {
			return result
		}
	}
	for !deck1.IsEmpty() && !deck2.IsEmpty() {
		if seenBefore(*deck1, *deck2, seenPtr) {
			return 1
		}
		copy_deck1 := deck1.CopyDeck()
		copy_deck2 := deck2.CopyDeck()
		seenPtr.addToDeck(1, copy_deck1)
		seenPtr.addToDeck(2, copy_deck2)
		PlayOnceRecursively(deck1, deck2, cache)
	}
	if deck1.IsEmpty() {
		(*cache)[initial_hash] = 2
		return 2
	} else {
		if !deck2.IsEmpty() {
			panic("No one won in recurisve play??")
		}
		(*cache)[initial_hash] = 1
		return 1
	}
}

func seenBefore(deck1 Deck, deck2 Deck, seen *SeenSet) bool {
	var seen_all_1s []int // idxs at which have had this deck
	for i, past_deck1 := range seen.deck1s {
		if deck1.Len() != past_deck1.Len() {
			continue
		}
		seen_all := true
		for i, card := range deck1.cards {
			if past_deck1.cards[i] != card {
				seen_all = false
				break
			}
		}
		if seen_all {
			seen_all_1s = append(seen_all_1s, i)
		}
	}

	if len(seen_all_1s) == 0 {
		return false
	}

	for _, j := range seen_all_1s {
		past_deck2 := seen.deck2s[j]
		if deck2.Len() != past_deck2.Len() {
			continue
		}
		seen_all := true
		for i, card := range deck2.cards {
			if past_deck2.cards[i] != card {
				seen_all = false
				break
			}
		}
		if seen_all {
			return true
		}
	}
	return false
}

func PlayOnceRecursively(deck1 *Deck, deck2 *Deck, cache *map[uint32]int) {
	var WINNER int

	// initial normal setup
	card1 := deck1.RemoveTopCard()
	card2 := deck2.RemoveTopCard()
	if card1 == card2 {
		panic("Hadn't prepared for equal valued cards!?!??")
	}

	// play
	if deck1.Len() >= card1 && deck2.Len() >= card2 { // recursive round!
		new_deck1 := deck1.CopyDeckUptoNth(card1)
		new_deck2 := deck2.CopyDeckUptoNth(card2)
		WINNER = PlayRecursively(&new_deck1, &new_deck2, cache)
	} else { // standard non-recursive round
		if card1 > card2 {
			WINNER = 1
		} else {
			// card2 > card1
			WINNER = 2
		}
	}

	if WINNER == 1 {
		deck1.AddTwoToBack(card1, card2)
	} else if WINNER == 2 {
		deck2.AddTwoToBack(card2, card1)
	}
}

func parseInput(filepath string) (Deck, Deck) {
	file, _ := os.Open(filepath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	deck1 := Deck{[]int{}}
	deck2 := Deck{[]int{}}
	player_mentions := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Player") {
			player_mentions++
		} else if line != "" {
			card_val, err := strconv.Atoi(line)
			if err != nil {
				panic("Couldn't parse card value!")
			}
			if player_mentions == 1 {
				deck1.AddCard(card_val)
			} else if player_mentions == 2 {
				deck2.AddCard(card_val)
			}
		}
	}
	return deck1, deck2
}

func Part1(filepath string) int {
	deck1, deck2 := parseInput(filepath)
	return Play(&deck1, &deck2, false)
}

func Part2(filepath string) int {
	deck1, deck2 := parseInput(filepath)
	return Play(&deck1, &deck2, true)
}
