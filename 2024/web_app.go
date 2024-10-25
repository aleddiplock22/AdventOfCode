package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func WebApp() {
	home_page := HomePage()
	results_page := ResultsPage(GetAllResults(DAY))
	stats_page := StatsPage()

	http.Handle("/", templ.Handler(home_page))
	http.Handle("/Results", templ.Handler(results_page))
	http.Handle("/Stats", templ.Handler(stats_page))

	fmt.Println("Running on http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
