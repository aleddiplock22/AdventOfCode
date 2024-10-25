package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func WebApp() {
	component := HomePage()

	results_page := ResultsPage(GetAllResults(DAY))

	http.Handle("/", templ.Handler(component))
	http.Handle("/Results", templ.Handler(results_page))

	fmt.Println("Running on http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}
