package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func day09(part2 bool) Solution {
	example_filepath := GetExamplePath(9)
	input_filepath := GetInputPath(9)
	if !part2 {
		example_p1 := Part1_09(example_filepath)
		AssertExample("1928", example_p1, 1)
		return Solution{
			"09",
			example_p1,
			Part1_09(input_filepath),
		}
	} else {
		example_p2 := Part2_09(example_filepath)
		return Solution{
			"09",
			example_p2,
			Part2_09(input_filepath),
		}
	}
}

func readDiscMapFromFilepath(filepath string) []string {
	return strings.Split(readInput(filepath), "")
}

func generateSpareDiscMap(disc_map []string) (spare_disc_map []string) {
	file_idx := 0
	for i, memory := range disc_map {
		if i%2 == 0 {
			// file
			file_size, err := strconv.Atoi(memory)
			if err != nil {
				panic("trouble reading file_size as int")
			}
			for range file_size {
				spare_disc_map = append(spare_disc_map, fmt.Sprintf("%d", file_idx))
			}
			file_idx++
		} else {
			free_mem, err := strconv.Atoi(memory)
			if err != nil {
				panic("trouble reading free_mem as int")
			}
			for range free_mem {
				spare_disc_map = append(spare_disc_map, ".")
			}
		}
	}
	return spare_disc_map
}

func CalculateChecksum(spare_disc_map []string, part2 bool) (total int) {
	for i, file_size := range spare_disc_map {
		if file_size == "#" || file_size == "." {
			if part2 {
				continue
			}
			break
		}
		file_size_int, err := strconv.Atoi(file_size)
		if err != nil {
			panic("invalid file size in final totalling")
		}
		total += i * file_size_int
	}
	return total
}

func Part1_09(filepath string) string {
	disc_map := readDiscMapFromFilepath(filepath)
	spare_disc_map := generateSpareDiscMap(disc_map)

	// now have our spare_disc_map
	// let's reorganise...
	to_re_alloc_idx := len(spare_disc_map) - 1
	var free_mem_idx int
	for {
		free_mem_idx = free_mem_idx + slices.Index(spare_disc_map[free_mem_idx:], ".")
		if free_mem_idx >= to_re_alloc_idx {
			break
		}
		file_mem := spare_disc_map[to_re_alloc_idx]
		if file_mem == "#" {
			panic("urhhh")
		}
		if file_mem == "." {
			to_re_alloc_idx--
			continue
		}
		spare_disc_map[free_mem_idx] = file_mem
		spare_disc_map[to_re_alloc_idx] = "#"
		to_re_alloc_idx--
	}

	total := CalculateChecksum(spare_disc_map, false)

	return fmt.Sprintf("%d", total)
}

func findBackwardsIndex(spare_disc_map *[]string, val string) int {
	for j := len(*spare_disc_map) - 1; j >= 0; j-- {
		if (*spare_disc_map)[j] == val {
			return j
		}
	}
	return -1
}

func CalculateChunkToRealloc(file_size_to_realloc_next int, spare_disc_map *[]string) int {
	fs_to_realloc := file_size_to_realloc_next
	idx_ := findBackwardsIndex(spare_disc_map, fmt.Sprintf("%d", fs_to_realloc))
	if idx_ == -1 {
		panic("Couldn't find the file size to realloc !")
	}
	chunk_to_realloc := 1
	_chunk_idx := idx_ - 1
	for {
		if (*spare_disc_map)[_chunk_idx] == fmt.Sprintf("%d", fs_to_realloc) {
			chunk_to_realloc++
			_chunk_idx--
		} else {
			break
		}
	}
	return chunk_to_realloc
}

func CalculateFreeMemoryChunkSize(spare_disc_map *[]string, min_idx int) int {
	// see how big a chunk is available
	free_chunk_size := 1
	_idx := min_idx + 1
	for {
		if _idx == len(*spare_disc_map) {
			break
		}
		if (*spare_disc_map)[_idx] == "." {
			free_chunk_size++
			_idx++
		} else {
			break
		}
	}
	return free_chunk_size
}

func Part2_09(filepath string) string {
	disc_map := readDiscMapFromFilepath(filepath)
	spare_disc_map := generateSpareDiscMap(disc_map)

	file_size_to_realloc_next_str := spare_disc_map[len(spare_disc_map)-1]
	file_size_to_realloc_next, err := strconv.Atoi(file_size_to_realloc_next_str)
	if err != nil {
		panic("initial file size to realloc next was not an int")
	}

	done_file_sizes := make(map[int]bool)

	var free_mem_idx int
	change := false
	_tmp_file_size_to_re_alloc := file_size_to_realloc_next
	var chunk_to_realloc int
	for {
		if _tmp_file_size_to_re_alloc <= 0 {
			if !change {
				break
			}
			change = false
			_tmp_file_size_to_re_alloc = file_size_to_realloc_next
		}
		if already_done := done_file_sizes[_tmp_file_size_to_re_alloc]; already_done {
			_tmp_file_size_to_re_alloc--
			continue
		}
		free_mem_idx = free_mem_idx + slices.Index(spare_disc_map[free_mem_idx:], ".")

		// size of memory chunk we want to move
		chunk_to_realloc = CalculateChunkToRealloc(_tmp_file_size_to_re_alloc, &spare_disc_map)

		// find space for chunk to fit
		_free_mem_idx := free_mem_idx
		var free_chunk_size int = 999999999
		moving := false

		_max_idx := findBackwardsIndex(&spare_disc_map, fmt.Sprintf("%d", _tmp_file_size_to_re_alloc))
		for {
			if _free_mem_idx >= _max_idx {
				// Failed for filesize
				_tmp_file_size_to_re_alloc--
				break
			}
			__idx_update := slices.Index(spare_disc_map[_free_mem_idx:], ".")
			if __idx_update == -1 {
				// __idx_update failure
				_tmp_file_size_to_re_alloc--
				break
			}
			_free_mem_idx = _free_mem_idx + __idx_update
			if _free_mem_idx >= _max_idx {
				// oob failure
				_tmp_file_size_to_re_alloc--
				break
			}
			free_chunk_size = CalculateFreeMemoryChunkSize(&spare_disc_map, _free_mem_idx)
			if chunk_to_realloc <= free_chunk_size {
				// it fits!
				if file_size_to_realloc_next == _tmp_file_size_to_re_alloc {
					file_size_to_realloc_next--
				}
				done_file_sizes[_tmp_file_size_to_re_alloc] = true
				moving = true
				break
			}
			_free_mem_idx++
		}
		if !moving {
			done_file_sizes[_tmp_file_size_to_re_alloc+1] = true // wah wah didnt move, never move?
		}

		if moving {
			change = true
			// actually move it
			mem_to_move_idx := findBackwardsIndex(&spare_disc_map, fmt.Sprintf("%d", _tmp_file_size_to_re_alloc))
			for range chunk_to_realloc {
				file_mem := spare_disc_map[mem_to_move_idx]
				if file_mem == "." {
					panic(". in MY memory?")
				}
				spare_disc_map[_free_mem_idx] = file_mem
				_free_mem_idx++
				spare_disc_map[mem_to_move_idx] = "."
				mem_to_move_idx--
			}
			_tmp_file_size_to_re_alloc = file_size_to_realloc_next
		}
	}

	total := CalculateChecksum(spare_disc_map, true)

	return fmt.Sprintf("%d", total)
}
