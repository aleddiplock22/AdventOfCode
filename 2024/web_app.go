package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
)

func WebApp() {
	home_page := HomePage()
	solutions_page := SolutionsPage(GetAllSolutions(DAY))
	stats_page := StatsPage()

	http.Handle("/", templ.Handler(home_page))
	http.Handle("/Solutions", templ.Handler(solutions_page))
	http.HandleFunc("/run-solution", runSolutionHandler)
	http.HandleFunc("/sse-progress", sseHandler)

	http.Handle("/Stats", templ.Handler(stats_page))

	fmt.Println("Running on http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}

/*
Progress handling explained...

		When the "Run Solution" button is clicked, it sends a POST request to /run-solution.
		The runSolutionHandler starts the solution process in a goroutine and renders the initial progress bar.
		The <div> with hx-ext="sse" establishes a Server-Sent Events connection to /sse-progress.
		The sseHandler sets up the SSE connection and listens for progress updates.
		The runSolution function simulates a long-running process, updating the progress every second.
		Progress updates are sent through a channel and then sent as SSE messages.
		HTMX receives these SSE messages and updates the progress bar accordingly.

*/

func runSolutionHandler(w http.ResponseWriter, r *http.Request) {
	// Start the solution process in a goroutine
	go runSolution()

	// Render initial progress bar
	ProgressBar(0).Render(r.Context(), w)
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	/*
		What does this do?
			- Sets up the SSE [Sever Sent Event]
			- Creates a channel progressChan to recieve progress updates
			- Starts 'listenForProgress' in a goroutine
			- Loops, waiting for progress updates on progressChan
			- For each update, it renders a new progress bar and sends it as an SSE message
	*/

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Create a channel to receive progress updates
	progressChan := make(chan int)

	// Start listening for progress updates
	go listenForProgress(progressChan)

	for progress := range progressChan {
		// Send SSE message
		fmt.Fprintf(w, "data: ")
		ProgressBar(progress).Render(r.Context(), w)
		fmt.Fprintf(w, "\n\n")
		flusher.Flush()

		if progress >= 100 {
			break
		}
	}
}

func runSolution() {
	/*
		TODO: update so this isn't a solution but actually has the AOC code running

		- I wonder if I can make it so we get to 25% for example p1, 50% full p1, 75% example p2, 100% full p2 ?

	*/

	// Simulate a long-running process
	for i := 0; i <= 10; i++ {
		progress := i * 10
		updateProgress(progress)
		time.Sleep(1 * time.Second) // Simulate work being done
	}
}

var progressUpdate = make(chan int, 1) // A buffered channel that can hold a single int value

func updateProgress(progress int) {
	select {
	case progressUpdate <- progress:
	default:
		// Channel is full, discard update (this avoids blocking)
	}
}

func listenForProgress(progressChan chan<- int) {
	// continuously reads from progressUpdate and sends the values to progressChan, which is connected to the SSE handler
	for progress := range progressUpdate {
		progressChan <- progress
	}
}
