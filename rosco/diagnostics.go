package rosco

import (
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
	coolantTempAverage float64
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
	// Sample contains the last minute of readings
	Sample []MemsData
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
	diagnostics.Sample = diagnostics.GetDataSetSample(60)
}

// Analyse runs a diagnostic review of the dataset
func (diagnostics *MemsDiagnostics) Analyse() {
	// work with a sample of the last 60 seconds of data
	diagnostics.Sample = diagnostics.GetDataSetSample(60)
	diagnostics.Analysis.IsWarm = diagnostics.isEngineWarm()
	diagnostics.Analysis.IsIdle = diagnostics.isIdle()
}

// GetDataSetSample asda
func (diagnostics *MemsDiagnostics) GetDataSetSample(points int) []MemsData {
	maxItems := len(diagnostics.Dataset)

	if points > maxItems {
		points = maxItems
	}

	return diagnostics.Dataset[maxItems-points:]
}

// GetMovingAverage asda
func (diagnostics *MemsDiagnostics) GetMovingAverage(metricName string) float64 {
	items := reflect.ValueOf(diagnostics.Sample)
	count := 0.0
	total := 0.0

	for i := 0; i < items.Len(); i++ {
		item := items.Index(i)
		if item.Kind() == reflect.Struct {
			v := reflect.Indirect(item).FieldByName(metricName)
			total = total + float64(v.Interface().(int8))
			count = count + 1
		}
	}

	return total / count
}

func (diagnostics *MemsDiagnostics) isEngineWarm() bool {
	return diagnostics.GetMovingAverage("CoolantTemp") >= engineWarmTemperature
}
func (diagnostics *MemsDiagnostics) isIdle() bool {
	rpm := diagnostics.GetMovingAverage("EngineRPM")

	if diagnostics.isEngineWarm() {
		return rpm >= minIdleWarmRPM && rpm <= maxIdleWarmRPM
	}

	return rpm >= minIdleColdRPM && rpm <= maxIdleColdRPM
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
