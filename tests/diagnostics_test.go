package tests

import (
	"github.com/andrewdjackson/memsfcr/rosco"
	"testing"
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

func TestNewMemsAnalytics(t *testing.T) {
	r := rosco.NewMemsAnalysisReport()
	if r == nil {
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

func TestMovingAverage(t *testing.T) {
	temp := [3]int8{20, 30, 40}
	rpm := [3]uint16{1100, 1180, 1200}

	expected := 30.0

	d := rosco.NewMemsDiagnostics()

	for i := 0; i < len(temp); i++ {
		item := rosco.MemsData{
			CoolantTemp: temp[i],
			EngineRPM:   rpm[i],
		}
		d.Add(item)
	}

	a := d.GetMovingAverage("CoolantTemp")

	if a != expected {
		t.Errorf("failed expected %d, got %f", temp, a)
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
