package tests

import (
	"math/rand"
	"testing"

	"github.com/andrewdjackson/memsfcr/rosco"
)

func createTestData(items int) []rosco.MemsData {
	memsdata := []rosco.MemsData{}

	for i := 0; i < items; i++ {
		item := rosco.MemsData{}
		memsdata = append(memsdata, item)
	}

	return memsdata
}

func TestNewDiagnostic(t *testing.T) {
	d := rosco.NewMemsDiagnostics()
	if d == nil {
		t.Errorf("failed")
	}
}

func TestAddData(t *testing.T) {
	d := rosco.NewMemsDiagnostics()
	item := rosco.MemsData{}
	d.Add(item)
	if len(d.Dataset) != 1 {
		t.Errorf("failed item not added")
	}
}

func getRandomValues(min int, max int, qty int) []int {
	var v []int

	for i := 0; i < qty; i++ {
		v = append(v, rand.Intn(max-min)+min)
	}

	return v
}

func TestMovingAverageUnstable(t *testing.T) {
	temp := getRandomValues(20, 88, 20)
	rpm := getRandomValues(800, 1250, 20)

	expected := 54.0

	d := rosco.NewMemsDiagnostics()

	for i := 0; i < len(temp); i++ {
		item := rosco.MemsData{
			CoolantTemp: temp[i],
			EngineRPM:   rpm[i],
		}
		d.Add(item)
	}

	d.Analyse()
	s := d.GetMetricStatistics("CoolantTemp")

	if s.Mean != expected {
		t.Errorf("failed expected %d, got %f", temp, s.Mean)
	}
	t.Logf("%v", s)
}

func TestMovingAverageTrendStable(t *testing.T) {
	temp := getRandomValues(86, 88, 20)
	rpm := getRandomValues(1100, 1250, 20)

	expected := 87.0

	d := rosco.NewMemsDiagnostics()

	for i := 0; i < len(temp); i++ {
		item := rosco.MemsData{
			CoolantTemp: temp[i],
			EngineRPM:   rpm[i],
		}
		d.Add(item)
	}

	d.Analyse()
	s := d.GetMetricStatistics("CoolantTemp")

	if s.Mean != expected {
		t.Errorf("failed expected %d", temp)
	}
}
func TestGetSample(t *testing.T) {
	datasetSize := 20
	items := 10

	d := rosco.NewMemsDiagnostics()
	d.Dataset = createTestData(datasetSize)
	s := d.GetDataSetSample(items)

	if len(s) != items {
		t.Errorf("failed %d items, %d expected", len(s), items)
	}
}

func TestGetSampleWithSmallDataset(t *testing.T) {
	datasetSize := 9
	items := 10

	d := rosco.NewMemsDiagnostics()
	d.Dataset = createTestData(datasetSize)
	s := d.GetDataSetSample(items)

	if len(s) != datasetSize {
		t.Errorf("failed %d items, %d expected", len(s), items)
	}
}

func TestAnalyse(t *testing.T) {
	datasetSize := 20

	temp := getRandomValues(86, 88, datasetSize)
	rpm := getRandomValues(1100, 1250, datasetSize)

	d := rosco.NewMemsDiagnostics()

	for i := 0; i < len(temp); i++ {
		item := rosco.MemsData{
			CoolantTemp: temp[i],
			EngineRPM:   rpm[i],
		}
		d.Add(item)
	}

	d.Analyse()

	if len(d.Dataset) != datasetSize {
		t.Errorf("failed expected %d", datasetSize)
	}
}
