package main

import (
	"github.com/vicanso/go-charts/v2"
)

func buildBarGraph() []byte {
	// based on package's example https://github.com/vicanso/go-charts/blob/main/examples/bar_chart/main.go

	run_times := [][]float64{
		{
			2.0,
			4.9,
		},
		{
			2.6,
			5.9,
		},
	}
	p, err := charts.BarRender(
		run_times,
		charts.XAxisDataOptionFunc([]string{
			"Day 1",
			"Day 2",
		}),
		charts.LegendLabelsOptionFunc([]string{
			"Part 1",
			"Part 2",
		}, charts.PositionRight),
	)
	if err != nil {
		panic(err)
	}

	buf, err := p.Bytes()
	if err != nil {
		panic(err)
	}
	return buf
}
