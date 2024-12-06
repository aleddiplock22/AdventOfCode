package main

import (
	"fmt"
	"testing"
)

func RunBenchmarking(day int, DayFunc SolutionFuncType) {
	p1_bench, p2_bench := RunBenchmark(DayFunc)
	fmt.Println(BOLD, RED_FORE, "---------------------------------", RESET)
	fmt.Println(GREEN_FORE, "  🎅", "ho ho ho, it's day", day, "benchmarks time...", "❄️ 🎄", RESET)
	fmt.Println(RED_FORE, " ---------------------------------", RESET)
	fmt.Println(RED_FORE, "part 1:", GREEN_FORE, formatTimeForBenchmark(p1_bench))
	fmt.Println(GREEN_FORE, "part 2:", RED_FORE, formatTimeForBenchmark(p2_bench), RESET, RESET_STYLE)
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
