package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	component := HomePage()

	results := []Result{day01(), day02()}
	results_page := ResultsPage(results)

	http.Handle("/", templ.Handler(component))
	http.Handle("/Results", templ.Handler(results_page))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}

type Result struct {
	day   string
	part1 string
	part2 string
}
