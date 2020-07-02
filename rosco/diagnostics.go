package rosco

import (
	"fmt"
	"reflect"

	"github.com/andrewdjackson/memsfcr/utils"
)

const (
	minIdleColdRPM        = 1200
	maxIdleColdRPM        = 900
	minIdleWarmRPM        = 700
	maxIdleWarmRPM        = 900
	engineWarmTemperature = 88
	lambdaLow             = 10
	lambdaHigh            = 900
	maxIdleError          = 50
)

// MemsLambdaStatus shows lambda diagnostics
type MemsLambdaStatus struct {
	minLambdaEvent bool
	maxLambdaEvent bool
	minLambdaFault bool
	maxLambdaFault bool
}

// MemsAnalysisReport is the output from running the analysis
type MemsAnalysisReport struct {
	// IsWarm is the engine warm
	IsWarm bool
	// IsIdle is the engine idle
	IsIdle bool
	// Lambda Status
	lambdaStatus MemsLambdaStatus
	// Sample Stats
	coolantTempAverage float32
}

// NewMemsAnalysisReport create a new report
func NewMemsAnalysisReport() *MemsAnalysisReport {
	analysis := &MemsAnalysisReport{}
	analysis.IsIdle = false
	analysis.IsWarm = false
	analysis.lambdaStatus = MemsLambdaStatus{}
	analysis.lambdaStatus.minLambdaEvent = false
	analysis.lambdaStatus.maxLambdaEvent = false
	analysis.lambdaStatus.minLambdaFault = false
	analysis.lambdaStatus.maxLambdaFault = false

	return analysis
}

// MemsDiagnostics logs the mems data to a CSV file
type MemsDiagnostics struct {
	// CurrentData is the lastest reading
	CurrentData MemsData
	// DataSet of mems data
	Dataset []MemsData
	// Analysis report
	Analysis MemsAnalysisReport
}

// NewMemsDiagnostics logs the mems data to a CSV file
func NewMemsDiagnostics() *MemsDiagnostics {
	diagnostics := &MemsDiagnostics{}
	diagnostics.Dataset = []MemsData{}
	diagnostics.Analysis = *NewMemsAnalysisReport()

	utils.LogI.Printf("starting diagnostics")

	return diagnostics
}

// Add data to the data set for diagnosis
func (diagnostics *MemsDiagnostics) Add(data MemsData) {
	diagnostics.CurrentData = data
	diagnostics.Dataset = append(diagnostics.Dataset, data)
}

// Analyse runs a diagnostic review of the dataset
func (diagnostics *MemsDiagnostics) Analyse() {
	// work with a sample of the last 60 seconds of data
	sample := diagnostics.getDataSetSample(60)
	diagnostics.Analysis.coolantTempAverage = diagnostics.getMovingAverage("CoolantTemp")
	diagnostics.Analysis.IsIdle = diagnostics.isIdle()
}

func (diagnostics *MemsDiagnostics) getDataSetSample(points int) []MemsData {
	maxItems := len(diagnostics.Dataset)

	if points > maxItems {
		points = maxItems
	}

	return diagnostics.Dataset[maxItems-points:]
}

func (diagnostics *MemsDiagnostics) getMovingAverage(sample []MemsData, metricName string) float32 {
	items := reflect.ValueOf(sample)

	for i := 0; i < items.Len(); i++ {
		item := items.Index(i)
		if item.Kind() == reflect.Struct {
			v := reflect.Indirect(item).FieldByName(metricName)
			for j := 0; j < v.NumField(); j++ {
				fmt.Println(v.Type().Field(j).Name, v.Field(j).Interface())
			}
		}
	}

	return 0
}

func (diagnostics *MemsDiagnostics) isIdle() bool {
	return diagnostics.CurrentData.CoolantTemp < engineWarmTemperature && diagnostics.CurrentData.EngineRPM < maxIdleWarmRPM
}

func (diagnostics *MemsDiagnostics) minIAC() bool {
	return (diagnostics.CurrentData.IdleSwitch == false && diagnostics.CurrentData.IdleSpeedDeviation >= maxIdleError && diagnostics.CurrentData.IACPosition == 0)
}

func (diagnostics *MemsDiagnostics) minLambda() bool {
	return int(diagnostics.CurrentData.LambdaVoltage) <= lambdaLow
}

func (diagnostics *MemsDiagnostics) maxLambda() bool {
	return int(diagnostics.CurrentData.LambdaVoltage) >= lambdaHigh
}
