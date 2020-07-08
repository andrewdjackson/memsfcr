package utils

import (
	"math"

	"gonum.org/v1/gonum/stat"
)

// Stats structure
type Stats struct {
	Name       string
	Value      float64
	Max        float64
	Min        float64
	Mean       float64
	Stddev     float64
	Mode       float64
	ModeCount  float64
	TrendSlope float64
	Trend      float64
}

// NewStats generates stats from a sample of float64 values
func NewStats(name string, data []float64) *Stats {
	// the sample stats
	s := &Stats{
		Name:  name,
		Value: data[len(data)-1],
	}

	// get the sample stats
	s.Min, s.Max = findMinAndMax(data)
	s.Mean, s.Stddev = stat.MeanStdDev(data, nil)
	s.Mode, s.ModeCount = stat.Mode(data, nil)

	// round to 2 decimal places
	s.Mean = math.Round(s.Mean * 100 / 100)
	s.Stddev = math.Round(s.Stddev * 100 / 100)
	s.Mode = math.Round(s.Mode * 100 / 100)
	s.TrendSlope, s.Trend = linearRegression(data)

	LogI.Printf("%s Stats %+v", DiagnosticTrace, *s)

	return s
}

func findMinAndMax(data []float64) (min float64, max float64) {
	min = data[0]
	max = data[0]
	for _, value := range data {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}

func linearRegression(data []float64) (float64, float64) {
	origin := false

	xs := make([]float64, len(data))
	for i := range xs {
		xs[i] = float64(i)
	}

	ys := data

	return stat.LinearRegression(xs, ys, nil, origin)
}
