package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func day10(part2 bool) Solution {
	example_filepath := GetExamplePath(10)
	input_filepath := GetInputPath(10)
	if !part2 {
		example_p1 := Part1_10(example_filepath)
		AssertExample("36", example_p1, 1)
		return Solution{
			"10",
			example_p1,
			Part1_10(input_filepath),
		}
	} else {
		example_p2 := Part2_10(example_filepath)
		AssertExample("81", example_p2, 2)
		return Solution{
			"10",
			example_p2,
			Part2_10(input_filepath),
		}
	}
}

func calculateTrailheadScore(grid *[][]int, sr, sc int, part2 bool) (score int) {
	R := len(*grid)
	C := len((*grid)[0])

	// point = [pos_r, pos_x, len_path]
	seen := make(map[[3]int]bool)
	queue := [][3]int{{sr, sc, 0}}
	for len(queue) > 0 {
		// pop
		tile := queue[0]
		queue = queue[1:]

		r, c, l := tile[0], tile[1], tile[2]
		height := (*grid)[r][c]
		// reached an end
		if height == 9 {
			score++
			continue
		}

		// find next step
		nl := l + 1
		for _, drdc := range [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
			dr, dc := drdc[0], drdc[1]
			nr := r + dr
			nc := c + dc
			if nr >= 0 && nr < R && nc >= 0 && nc < C {
				new_height := (*grid)[nr][nc]
				if new_height-height != 1 {
					continue
				}
				new_point := [3]int{nr, nc, nl}
				if !part2 {
					if already_travelled := seen[new_point]; already_travelled {
						continue
					}
					seen[new_point] = true
				}
				queue = append(queue, new_point)
			}
		}

	}
	return score
}

func calculateTrailheadScoreRecursively(grid *[][]int, sr, sc int) (score int) {
	// this is nice I guess. I can't think how to cleanly do p1 with recursion though
	// would have to keep some set of the '9's we reached, i.e. their coords, and count them
	// but tricky to combine a p1 and p2 function in that case. I'll leave it.
	R := len(*grid)
	C := len((*grid)[0])

	height := (*grid)[sr][sc]
	// reached an end
	if height == 9 {
		return 1
	}

	// find next step
	for _, drdc := range [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
		dr, dc := drdc[0], drdc[1]
		nr := sr + dr
		nc := sc + dc
		if nr >= 0 && nr < R && nc >= 0 && nc < C {
			new_height := (*grid)[nr][nc]
			if new_height-height != 1 {
				continue
			}
			score += calculateTrailheadScoreRecursively(grid, nr, nc)
		}
	}
	return score
}

func SolvePathCountingTask(filepath string, part2 bool, use_recursive_solution bool) string {
	grid := readIntGrid(filepath) // grid[r][c]

	var wg sync.WaitGroup
	scoreChan := make(chan int)
	for r, row := range grid {
		for c, value := range row {
			if value == 0 {
				// start path here!
				wg.Add(1)
				go func(sr, sc int) {
					defer wg.Done()
					if use_recursive_solution {
						if part2 {
							scoreChan <- calculateTrailheadScoreRecursively(&grid, sr, sc)
						}
					} else {
						scoreChan <- calculateTrailheadScore(&grid, sr, sc, part2)
					}
				}(r, c)
			}
		}
	}

	go func() {
		wg.Wait()
		close(scoreChan)
	}()

	var total int
	for score := range scoreChan {
		total += score
	}

	return fmt.Sprintf("%d", total)
}

func Part1_10(filepath string) string {
	return SolvePathCountingTask(filepath, false, false)
}

func Part2_10(filepath string) string {
	return SolvePathCountingTask(filepath, true, true)
}

// Visualisation
/*
	Using bubbletea package

	going to need:
		- Init, a function that returns an initial command for the application to run.
		- Update, a function that handles incoming events and updates the model accordingly.
		- View, a function that renders the UI based on the data in the model.
*/
type modelDay10Grid struct {
	grid                [][]int
	stylised_grid       [][]termenv.Style
	paths               [][10][2]int
	num_paths_displayed int
	stylised_paths      [][10][2]termenv.Style
}

func getPaths(sr, sc, R, C int, grid *[][]int, path *[10][2]int, paths *[][10][2]int) {
	height := (*grid)[sr][sc]
	// reached an end
	if height == 9 {
		*paths = append(*paths, *path)
		return
	}

	// find next step
	for _, drdc := range [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}} {
		dr, dc := drdc[0], drdc[1]
		nr := sr + dr
		nc := sc + dc
		if nr >= 0 && nr < R && nc >= 0 && nc < C {
			new_height := (*grid)[nr][nc]
			if new_height-height != 1 {
				continue
			}
			path[new_height] = [2]int{nr, nc}
			getPaths(nr, nc, R, C, grid, path, paths)
		}
	}
}

func initialModelDay10() modelDay10Grid {
	// default starting model should have all paths ready populated, 0 displayed
	grid := readIntGrid(GetInputPath(10))
	R := len(grid)
	C := len(grid)
	var paths [][10][2]int
	for r, row := range grid {
		for c, value := range row {
			path := [10][2]int{}
			path[0] = [2]int{r, c}
			if value == 0 {
				getPaths(r, c, R, C, &grid, &path, &paths)
			}
		}
	}

	p := termenv.ANSI256
	stylised_grid := make([][]termenv.Style, len(grid))
	for i, row := range grid {
		stylised_grid[i] = make([]termenv.Style, len(grid[0]))
		for j, val := range row {
			// background color
			bg := p.Color(fmt.Sprintf("%d", 24+3*(val+1)))
			out := termenv.String(fmt.Sprintf(" %d ", val))

			// apply colors
			out = out.Foreground(p.Color("7"))
			out = out.Background(bg)

			stylised_grid[i][j] = out
		}
	}

	return modelDay10Grid{
		grid,
		stylised_grid,
		paths,
		0,
		[][10][2]termenv.Style{},
	}
}

func (m modelDay10Grid) Init() tea.Cmd {
	// just returning nil, i.e. we don't need any I/O
	return nil
}

func (m modelDay10Grid) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// detect key press
	case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	m.num_paths_displayed++
	return m, nil
}

func (m modelDay10Grid) View() string {
	// heavy inspo from the colour grid in example from termenv
	// here: https://github.com/muesli/termenv/blob/master/examples/color-chart/main.go

	// The header
	s := termenv.String("\n\nDay 10 Visualisation! 🎅🟩⁉️ \n\n").Bold().String()

	p := termenv.ANSI256

	next_path := m.num_paths_displayed - 1
	if next_path > 0 && next_path < len(m.paths)-1 {
		color_str := fmt.Sprintf("%d", rand.Intn(231-161)+161)
		for j, position := range m.paths[m.num_paths_displayed] {
			if j == 0 || j == 9 {
				continue
			}
			r, c := position[0], position[1]
			replacement_str := termenv.String(fmt.Sprintf(" %d ", m.grid[r][c]))
			bg := p.Color(color_str)
			replacement_str = replacement_str.Foreground(p.Color("7"))
			replacement_str = replacement_str.Background(bg)

			m.stylised_grid[r][c] = replacement_str
		}
	}

	// Iterate over our choices
	for _, row := range m.stylised_grid {
		// render each row of grid
		s += "\n"
		for _, formatted_val := range row {
			s += formatted_val.String()
		}
	}

	// The footer
	s += termenv.String("\nHold any key or move mouse to display more paths!").Bold().String()
	if m.num_paths_displayed < len(m.paths) {
		s += termenv.String(fmt.Sprintf("\nNum paths displayed: %d", m.num_paths_displayed)).Bold().String()
	} else {
		s += termenv.String(fmt.Sprintf("You found all %d paths! Nice!\n", m.num_paths_displayed)).Blink().String()
	}
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func DoDay10Visualisation() {
	// enable ANSI parsing since on windows
	restoreConsole, err := termenv.EnableVirtualTerminalProcessing(termenv.DefaultOutput())
	if err != nil {
		panic(err)
	}
	defer restoreConsole()

	p := tea.NewProgram(initialModelDay10())
	if _, err := p.Run(); err != nil {
		fmt.Printf("uh oh error, rip: %v\n", err)
		os.Exit(1)
	}
}
