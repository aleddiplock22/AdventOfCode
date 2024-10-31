package main

import (
	"github.com/wcharczuk/go-chart"
)

func buildBarGraph() chart.BarChart {
	graph := chart.BarChart{
		Width:  500,
		Height: 300,
		Bars: []chart.Value{
			{Value: 5, Label: "A"},
			{Value: 10, Label: "B"},
			{Value: 15, Label: "C"},
			{Value: 20, Label: "D"},
		},
	}
	return graph
}
