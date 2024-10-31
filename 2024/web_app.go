package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/wcharczuk/go-chart"
)

func WebApp() {
	home_page := HomePage()
	solutions_page := SolutionsPage(GetAllSolutions(DAY))
	stats_page := StatsPage()

	http.Handle("/", templ.Handler(home_page))
	http.Handle("/Solutions", templ.Handler(solutions_page))
	http.HandleFunc("/run-solution", runSolutionHandler)
	http.HandleFunc("/sse-progress", sseHandler)
	http.HandleFunc("/which-day", whichDayHandler)
	http.HandleFunc("/bar-chart", barChartHandler)

	http.Handle("/Stats", templ.Handler(stats_page))

	fmt.Println("Running on http://localhost:3000")
	http.ListenAndServe(":3000", nil)
}

var whichDayNumber int = 1

func whichDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		inputNumber := r.FormValue("whichDayText")
		if day_num, err := strconv.Atoi(inputNumber); err != nil {
			panic("Unrecognised day number - not a number!")
		} else {
			whichDayNumber = day_num
		}
		fmt.Printf("whichDayNumber set to: Day %v.\n", whichDayNumber)
	}
}

func barChartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// So first when hx-post is triggered, we create the image element
		html := `<img src="/bar-chart" width="500" height="300" style="border: 1px solid black;">`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	} else if r.Method == http.MethodGet {
		// then the above triggers a GET when trying to find src="/bar-chart" I think
		// and this renders to it
		graph := buildBarGraph()
		w.Header().Set("Content-Type", "image/png")
		graph.Render(chart.PNG, w)
	}
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
	go runSolution(whichDayNumber)

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

	// back up for getting the bar to 100 if the polling got blocked and discarded the update...
	ProgressBar(100).Render(r.Context(), w)
}

func wait(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
}

func runSolution(day int) {

	DayFunc := dayMap[strconv.Itoa(day)]
	wait(1)
	updateProgress(10)
	wait(2)
	part1 := DayFunc(false)
	updateProgress(50)
	wait(2)
	part2 := DayFunc(true)
	updateProgress(100)
	wait(1)

	// do something with the results...
	fmt.Println(part1.day, part2.day, "have run .,,,")
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
