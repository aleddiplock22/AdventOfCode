package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func WebApp() {
	home_page := HomePage()
	solutions_page := SolutionsPage(GetAllSolutions(DAY), currentProgress)
	stats_page := StatsPage()

	http.Handle("/", templ.Handler(home_page))
	http.Handle("/Solutions", templ.Handler(solutions_page))
	http.HandleFunc("/increment", incrementHandler)

	http.Handle("/Stats", templ.Handler(stats_page))

	fmt.Println("Running on http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}

var currentProgress = 0

func incrementHandler(w http.ResponseWriter, r *http.Request) {
	if currentProgress < 90 {
		currentProgress += 10
	} else if currentProgress < 100 {
		currentProgress++
	}

	if currentProgress > 100 {
		currentProgress = 100
	}

	w.Header().Set("Content-Type", "text/html")
	ProgressBar(currentProgress).Render(r.Context(), w)

	// If not complete, trigger another request
	// if currentProgress < 100 {
	// 	w.Header().Set("HX-Trigger", "nextIncrement")
	// } else {
	// 	currentProgress = 0
	// }
}
