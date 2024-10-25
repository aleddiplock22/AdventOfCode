package main

import (
	"flag"
	"fmt"
)

// UPDATE DAY HERE FOR WEB APP ALL_RESULTS
const DAY = 2

func main() {
	day := flag.String("day", "NA", "Provide day number to run if quick mode")
	flag.Parse()

	if *day != "NA" {
		// quick dev on the day
		if dayFunc, exists := dayMap[*day]; exists {
			SingleDayDevelopment(dayFunc)
		} else {
			panic(fmt.Sprintf("Didn't recognise day: \"%v\". Panic time!", *day))
		}
	} else {
		// web app time :)
		WebApp()
	}
}
