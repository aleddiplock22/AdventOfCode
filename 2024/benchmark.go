package main

import (
	"fmt"
	"testing"
)

const WithAscii = false

func RunBenchmarking(day int, DayFunc SolutionFuncType) {
	p1_bench, p2_bench := RunBenchmark(DayFunc)
	fmt.Println(BOLD, RED_FORE, "--------------------------------------------------", RESET)
	fmt.Println(GREEN_FORE, "  🎅", "ho ho ho, it's day", day, "benchmarks time...", "❄️ 🎄", RESET)
	fmt.Println(RED_FORE, " --------------------------------------------------", RESET)
	fmt.Println(RED_FORE, "part 1:", GREEN_FORE, formatTimeForBenchmark(p1_bench))
	fmt.Println(GREEN_FORE, "part 2:", RED_FORE, formatTimeForBenchmark(p2_bench))
	if WithAscii {
		fmt.Println(GREEN_FORE, WasItSlowSantaAsciiArt, RESET, RESET_STYLE)
	} else {
		fmt.Print(RESET, RESET_STYLE)
	}
}

func formatTimeForBenchmark(bench testing.BenchmarkResult) string {
	us := bench.T.Microseconds()
	if us < 1000 {
		return fmt.Sprintf("%dμs", us)
	}
	if us < 1000000 {
		return fmt.Sprintf("%dms", bench.T.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", bench.T.Seconds())
}

func RunBenchmark(DayFunc SolutionFuncType) (testing.BenchmarkResult, testing.BenchmarkResult) {
	p1_bench := testing.Benchmark(func(b *testing.B) {
		DayFunc(false)
	})
	p2_bench := testing.Benchmark(func(b *testing.B) {
		DayFunc(true)
	})
	return p1_bench, p2_bench
}

const WasItSlowSantaAsciiArt = `
	_____________,--,
      | | | | | | |/ .-.\   Was it Slow?
      |_|_|_|_|_|_/ /   '.           ho ho
       |_|__|__|_; |      \
       |___|__|_/| |     .''}
       |_|__|__/ | |   .'.''\
       |__|__|/  ; ;  / /    \.-"-.
       ||__|_;   \ \  ||    /'___. \
       |_|___/\  /;.',\\   {_'___.;{}
       |__|_/ ';'__|'-.;|  |C' e e'\
       |___'L  \__|__|__|  | ''-o-' }
       ||___|\__)___|__||__|\   ^  /'\
       |__|__|__|__|__|_{___}'.__.'\_.'}
       ||___|__|__|__|__;\_)-''\   {_.-;
       |__|__|__|__|__|/' ('\__/     '-'
       |_|___|__|__/'      |
-------|__|___|__/'         \-------------------
-.__.-.|___|___;'            |.__.-.__.-.__.-.__
  |     |     ||             |  |     |     |
-' '---' '---' \             /-' '---' '---' '--
     |     |    '.        .' |     |     |     |
'---' '---' '---' '-===-''--' '---' '---' '---'
  |     |     |     |     |     |     |     |
-' '---' '---' '---' '---' '---' '---' '---' '--
     |     |     |     |     |     |     |     |
'---' '---' '---' '---' '---' '---' '---' '---'
`
