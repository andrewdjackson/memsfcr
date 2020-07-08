package tests

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/andrewdjackson/memsfcr/utils"
)

func TestStats(t *testing.T) {
	min := 100
	max := 200

	data := getRandomDataset(min, max, 10)

	s := utils.NewStats("Test", data)

	if s.Max >= float64(max) {
		t.Errorf("failed")
	}
}

func TestIncreasingStats(t *testing.T) {
	var s *utils.Stats

	data := getIncreasingDataset(0, 1000, 60)

	s = utils.NewStats("Test", data)

	if s.Trend != 1 {
		t.Errorf("failed %+v", s)
	}
}

func TestDecreasingStats(t *testing.T) {
	var s *utils.Stats

	data := getDecreasingDataset(1000, 50, 60)

	s = utils.NewStats("Test", data)

	if s.Trend != -1 {
		t.Errorf("failed %+v", s)
	}
}

func TestStaticStats(t *testing.T) {
	var s *utils.Stats

	data := getStaticDataset(100, 10)

	s = utils.NewStats("Test", data)

	if s.Trend != -1 {
		t.Errorf("failed %+v", s)
	}
}

// helper functions

func getRandomDataset(min int, max int, qty int) []float64 {
	var v []float64

	for i := 0; i < qty; i++ {
		v = append(v, float64(rand.Intn(max-min)+min))
	}

	fmt.Printf("%+v\n", v)
	return v
}

func getStaticDataset(value int, qty int) []float64 {
	v := []float64{}

	for i := 0; i < qty; i++ {
		x := float64(value)
		v = append(v, x)
	}

	fmt.Printf("%+v\n", v)
	return v
}

func getIncreasingDataset(min int, max int, qty int) []float64 {
	v := []float64{}
	step := (max - min) / qty

	for i := 0; i < qty; i++ {
		x := float64(min + (i * step))
		v = append(v, x)
	}
	fmt.Printf("%+v\n", v)
	return v
}

func getDecreasingDataset(max int, min int, qty int) []float64 {
	v := []float64{}
	step := (max - min) / qty

	for i := 0; i < qty; i++ {
		x := float64(max - (i * step))
		v = append(v, x)
	}

	fmt.Printf("%+v\n", v)
	return v
}
