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

	fmt.Println("---Day 17---")

	fmt.Println("[Example P1] Expected: 45, Answer:", Part1(EXAMPLE_FILEPATH))
	fmt.Println("Part 1 Answer:", Part1(INPUT_FILEPATH))

	fmt.Println("[Example P2] Expected: 112, Answer:", Part2(EXAMPLE_FILEPATH))
	fmt.Println("Part 2 Answer:", Part2(INPUT_FILEPATH))
}

type State struct {
	x_coord int
	y_coord int
	x_vel   int
	y_vel   int
}

func (s *State) Step() {
	s.x_coord += s.x_vel
	s.y_coord += s.y_vel
	if s.x_vel > 0 {
		s.x_vel--
	} else if s.x_vel < 0 {
		s.x_vel++
	}
	s.y_vel--
}

type TargetArea struct {
	xl int // lower bound
	xu int
	yl int
	yu int
}

func Intersection(s *State, target *TargetArea) bool {
	return (s.x_coord >= target.xl &&
		s.x_coord <= target.xu &&
		s.y_coord >= target.yu &&
		s.y_coord <= target.yl)
}

func Beyond(s *State, target *TargetArea) bool {
	return s.x_coord > target.xu || s.y_coord < target.yu
}

func parseInput(filepath string) TargetArea {
	file, _ := os.ReadFile(filepath)
	contents := string(file)
	info := strings.Split(contents, "x=")[1]
	prts := strings.Split(info, "..")
	xl_str := prts[0]
	xl, _ := strconv.Atoi(xl_str)
	xu_str := strings.Split(prts[1], ",")[0]
	xu, _ := strconv.Atoi(xu_str)
	yinfo := strings.Split(prts[1], ", y=")[1]
	yprts := strings.Split(yinfo, "..")
	yu_str := yprts[0]
	yu, _ := strconv.Atoi(yu_str)
	yl_str := strings.Split(prts[2], "\r\n")[0]
	yl, _ := strconv.Atoi(yl_str)
	return TargetArea{
		xl, xu, yl, yu,
	}
}

func SimulateForMaxHeight(x_start_vel int, y_start_vel int, target *TargetArea) (int, error) {
	state := &State{
		x_coord: 0,
		y_coord: 0,
		x_vel:   x_start_vel,
		y_vel:   y_start_vel,
	}
	var max_y int = -9999
	err := fmt.Errorf("didn't find an intersection before going beyond")
	for {
		max_y = max(max_y, state.y_coord)
		if Intersection(state, target) {
			err = nil
		} else if Beyond(state, target) {
			break
		}
		state.Step()
	}
	return max_y, err
}

func Part1(filepath string) int {
	target := parseInput(filepath)
	best := -9999

	for vx := 1; vx <= target.xu; vx++ {
		for vy := -200; vy <= 200; vy++ {
			y_max, err := SimulateForMaxHeight(vx, vy, &target)
			if err != nil {
				continue
			}
			best = max(y_max, best)
		}
	}

	return best
}

func Part2(filepath string) int {
	target := parseInput(filepath)

	count := 0
	for vx := 1; vx <= target.xu; vx++ {
		for vy := -200; vy <= 200; vy++ {
			_, err := SimulateForMaxHeight(vx, vy, &target)
			if err != nil {
				continue
			}
			count++
		}
	}

	return count
}
