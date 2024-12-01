package main

func day01(part2 bool) Solution {
	// fmt.Println("AOC 24 Day 1")
	// fmt.Println(GetInputPath(1))

	if !part2 {
		wait(1)
		return Solution{
			"01",
			"ExampleP1",
			"InputAnsP1",
		}
	} else {
		wait(2)
		return Solution{
			"01",
			"ExampleP2",
			"InputAnsP2",
		}
	}
}
