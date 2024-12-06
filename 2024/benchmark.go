package main

import (
	"fmt"
	"testing"
)

/*
	TODO: improve benchmarking
*/

func RunBenchmarking(DayFunc SolutionFuncType) {
	p1_bench, p2_bench := RunBenchmark(DayFunc)
	fmt.Println("part 1:", p1_bench.T.Microseconds(), "micro seconds")
	fmt.Println("part 2:", p2_bench.T.Microseconds(), "micro seconds")
	// fmt.Println(BOLD, RED_FORE, "---------------------------------", RESET)
	// fmt.Println(GREEN_FORE, "  🎅🎄❄️  AOC 2024 - Day", res_p1.day, "❄️ 🎄🎅", RESET)
	// fmt.Println(RED_FORE, " ---------------------------------", RESET)
	// fmt.Printf("%v%vPart 1:%v\n\tExample: %v%v\n\t%vInput: %v%v\n%vPart 2:%v\n\tExample: %v%v\n\t%vInput: %v%v\n%v%v", GREEN_FORE, BOLD, RED_FORE, GREEN_FORE, res_p1.example_ans, RED_FORE, GREEN_FORE, res_p1.input_ans, GREEN_FORE, RED_FORE, GREEN_FORE, res_p2.example_ans, RED_FORE, GREEN_FORE, res_p2.input_ans, RESET, RESET_STYLE)

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
