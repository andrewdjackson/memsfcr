package rosco

import (
	"math"
	"reflect"

	"github.com/andrewdjackson/memsfcr/utils"
	"gonum.org/v1/gonum/stat"
)

const (
	minIdleColdRPM        = 1200
	maxIdleColdRPM        = 900
	minIdleWarmRPM        = 700
	maxIdleWarmRPM        = 900
	minIdleMap            = 30
	maxIdleMap            = 60
	bestAFR               = 14.7
	minMAPEngineOff       = 95
	engineWarmTemperature = 88
	lambdaLow             = 10
	lambdaHigh            = 900
	maxIdleError          = 50
	maxSamples            = 30 // ~30 seconds
)

// MemsSampleStats holds the statistics from a sample of a given metric
type MemsSampleStats struct {
	Name      string
	Value     float64
	Max       float64
	Min       float64
	Mean      float64
	Stddev    float64
	Mode      float64
	ModeCount float64
	Skew      float64
}

// NewMemsSampleStats generates stats from a sample of float64 values
func NewMemsSampleStats(metricName string, metricSample []float64) *MemsSampleStats {
	// the sample stats
	s := &MemsSampleStats{
		Name:  metricName,
		Value: metricSample[len(metricSample)-1],
	}

	// get the sample stats
	s.Min, s.Max = findMinAndMax(metricSample)
	s.Mean, s.Stddev = stat.MeanStdDev(metricSample, nil)
	s.Mode, s.ModeCount = stat.Mode(metricSample, nil)
	s.Skew = stat.Skew(metricSample, nil)

	// round to 2 decimal places
	s.Mean = math.Round(s.Mean * 100 / 100)
	s.Stddev = math.Round(s.Stddev * 100 / 100)
	s.Mode = math.Round(s.Mode * 100 / 100)
	s.Skew = math.Round(s.Skew * 100 / 100)

	utils.LogI.Printf("%s Engine stats %+v", utils.DiagnosticTrace, *s)

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

// MemsAnalysisReport is the output from running the analysis
type MemsAnalysisReport struct {
	IsEngineRunning     bool
	IsEngineWarming     bool
	IsEngineWarm        bool
	IsEngineIdle        bool
	IsCruising          bool
	IsClosedLoop        bool
	ClosedLoopExpected  bool
	MapFault            bool
	VacuumFault         bool
	IdleAirControlFault bool
	LambdaFault         bool
}

// MemsDiagnostics structure
type MemsDiagnostics struct {
	// CurrentData is the lastest reading
	CurrentData MemsData
	// Full DataSet of mems data
	Dataset []MemsData
	// Sample contains the last n readings
	Sample []MemsData
	// Stats of the sample
	Stats map[string]MemsSampleStats
	// Analysis report
	Analysis MemsAnalysisReport
}

// NewMemsDiagnostics generates diagnostic reports
func NewMemsDiagnostics() *MemsDiagnostics {
	diagnostics := &MemsDiagnostics{}
	diagnostics.Dataset = []MemsData{}
	diagnostics.Analysis = MemsAnalysisReport{}
	diagnostics.Stats = make(map[string]MemsSampleStats)

	utils.LogI.Printf("%s Starting diagnostics", utils.DiagnosticTrace)

	return diagnostics
}

// Add data to the data set for diagnosis
func (diagnostics *MemsDiagnostics) Add(data MemsData) {
	diagnostics.CurrentData = data
	diagnostics.Dataset = append(diagnostics.Dataset, data)
}

// Analyse runs a diagnostic review of the dataset
func (diagnostics *MemsDiagnostics) Analyse() {
	// work with a sample of the last n seconds of data
	diagnostics.Sample = diagnostics.GetDataSetSample(maxSamples)

	// get running stats
	diagnostics.Stats["CoolantTemp"] = diagnostics.GetMetricStatistics("CoolantTemp")
	diagnostics.Stats["EngineRPM"] = diagnostics.GetMetricStatistics("EngineRPM")
	diagnostics.Stats["ManifoldAbsolutePressure"] = diagnostics.GetMetricStatistics("ManifoldAbsolutePressure")
	diagnostics.Stats["LambdaVoltage"] = diagnostics.GetMetricStatistics("LambdaVoltage")
	diagnostics.Stats["AirFuelRatio"] = diagnostics.GetMetricStatistics("AirFuelRatio")

	diagnostics.checkIsEngineRunning()
	diagnostics.checkIsEngineWarm()
	diagnostics.checkIsEngineIdle()
	diagnostics.checkMapSensor()
	diagnostics.checkForExpectedClosedLoop()
	diagnostics.checkIdleAirControl()
	diagnostics.checkLambdaStatus()

	utils.LogI.Printf("%s Analysed engine data %+v", utils.DiagnosticTrace, diagnostics.Analysis)
}

// GetDataSetSample retrieves a slice of the dataset for the last n points
func (diagnostics *MemsDiagnostics) GetDataSetSample(points int) []MemsData {
	maxItems := len(diagnostics.Dataset)

	if points > maxItems {
		points = maxItems
	}

	return diagnostics.Dataset[maxItems-points:]
}

// GetMetricStatistics takes the sample and calculates the simple average
// this is useful to detect the trend for a metric
func (diagnostics *MemsDiagnostics) GetMetricStatistics(metricName string) MemsSampleStats {
	// get the fields available in the sample
	sampleValues := reflect.ValueOf(diagnostics.Sample)
	// an array to hold the sample
	metricSample := []float64{}

	// iterate the fields and create an array of values for the specific metric only
	for i := 0; i < sampleValues.Len(); i++ {
		sampleValue := sampleValues.Index(i)
		if sampleValue.Kind() == reflect.Struct {
			v := reflect.Indirect(sampleValue).FieldByName(metricName)
			// don't try to create metrics for strings, bools or uints (they're bit patterns)
			switch v.Interface().(type) {
			case int:
				metricSample = append(metricSample, float64(v.Interface().(int)))
			case float32:
				metricSample = append(metricSample, float64(v.Interface().(float32)))
			}
		}
	}

	// calculate the stats for this sample
	return *NewMemsSampleStats(metricName, metricSample)
}

// IsEngineWarm uses the current engine temperature and the standard deviation in the sample to determine the
// stability of the temperature. If the reading is at the designated thermostat temp (88C) and the std deviation
// is low then deem the engine to be running at operating temperature
func (diagnostics *MemsDiagnostics) checkIsEngineWarm() {
	diagnostics.Analysis.IsEngineWarm = (diagnostics.Stats["CoolantTemp"].Value >= engineWarmTemperature && diagnostics.Stats["CoolantTemp"].Stddev < 5)
}

// if the last reading of engine RPM is 0 then the engine is not running
// we don't use the sample set as the engine may have recently been stopped
func (diagnostics *MemsDiagnostics) checkIsEngineRunning() {
	diagnostics.Analysis.IsEngineRunning = !(diagnostics.CurrentData.EngineRPM == 0)
}

// IsIdle determines the correct idle speed parameters based on whether the engine is warm or cold
// if the RPM is within the parameters for the sample period then the engine is deemed to be at Idle
func (diagnostics *MemsDiagnostics) checkIsEngineIdle() {
	if diagnostics.Analysis.IsEngineWarm {
		// use warm idle settings
		diagnostics.Analysis.IsEngineIdle = diagnostics.Stats["EngineRPM"].Mean >= minIdleWarmRPM && diagnostics.Stats["EngineRPM"].Mean <= maxIdleWarmRPM
		diagnostics.Analysis.IsCruising = diagnostics.Stats["EngineRPM"].Mean > maxIdleWarmRPM
		diagnostics.Analysis.IsEngineWarming = false
	}

	// use cold idle settings
	diagnostics.Analysis.IsEngineIdle = diagnostics.Stats["EngineRPM"].Mean >= minIdleColdRPM && diagnostics.Stats["EngineRPM"].Mean <= maxIdleColdRPM
	diagnostics.Analysis.IsEngineWarming = true
}

// Manifold Pressure (KPa): This displays the pressure measured by the external MEMS air pressure sensor.
// Normal reading with the engine not running is approximately 100 KPa
// 30-40 KPa when the engine is idling.
// Very high values may indicate problems with the sensor or a blocked or disconnected vacuum pipe.
// Moderately raised values may indicate mechanical problems with the engine
func (diagnostics *MemsDiagnostics) checkMapSensor() {
	if diagnostics.Analysis.IsEngineRunning {
		// only check if engine is running at idle
		if diagnostics.Analysis.IsEngineIdle {
			// fault if the map readings are outside of expected when idling
			diagnostics.Analysis.MapFault = !(diagnostics.Stats["ManifoldAbsolutePressure"].Mean >= minIdleMap && diagnostics.Stats["ManifoldAbsolutePressure"].Mean <= maxIdleMap)
		}
	} else {
		// fault if the map is reading low when the engine is off
		diagnostics.Analysis.MapFault = diagnostics.Stats["ManifoldAbsolutePressure"].Mean < minMAPEngineOff
	}
}

// determines whether we're expecting the ECU to use closed loop.
// ECU will generally only use the lambda sensor’s output during two specific conditions
// (a) during idle, ie. when the engine is under no load apart from keeping itself running, and
// (b) during part-load conditions (which we usually term ‘cruising speed’) where the engine is keeping the car at a constant speed.
// Fast idle is typically 2500 - 3000 RPM
// Slow idle is typically  450 - 1500 RPM
func (diagnostics *MemsDiagnostics) checkForExpectedClosedLoop() {
	diagnostics.Analysis.IsClosedLoop = diagnostics.CurrentData.ClosedLoop
	diagnostics.Analysis.ClosedLoopExpected = diagnostics.Analysis.IsEngineWarm && (diagnostics.Analysis.IsEngineIdle || diagnostics.Analysis.IsCruising)
}

// if a hose is split the vacuum sensor in the ECU doesn't see true manifold pressure,
// but something of a slightly higher absolute pressure (a little closer to atmospheric).
// The ECU thinks then that the engine is more highly loaded, for the same RPM, than it really is and gives more fuel
func (diagnostics *MemsDiagnostics) checkForVacuumFault() {
	// wonder if this will be true if the AFR is rich and the MAP reading is high
	diagnostics.Analysis.VacuumFault = (diagnostics.Stats["ManifoldAbsolutePressure"].Mean >= maxIdleMap && diagnostics.Stats["AirFuelRatio"].Mean > bestAFR)
}

// Also known as stepping motor--idle air control valve (IACV)
// bolts on the side of the injection body housing to control engine idle speed
// and air flow from cold start up
func (diagnostics *MemsDiagnostics) checkIdleAirControl() {
	if diagnostics.Analysis.IsEngineRunning {
		// IAC fault if the idle offset exceeds the max error, yet the IAC position remains at 0
		diagnostics.Analysis.IdleAirControlFault = (diagnostics.CurrentData.IdleSpeedDeviation >= maxIdleError && diagnostics.CurrentData.IACPosition == 0)
	} else {
		diagnostics.Analysis.IdleAirControlFault = false
	}
}

//  At 2000 rpm it should be switching rapidly between the minimum and maximum figures as the MEMS controls the engine conditions.????
func (diagnostics *MemsDiagnostics) checkLambdaStatus() {
	if diagnostics.Analysis.IsEngineRunning && diagnostics.Analysis.IsClosedLoop {
		if diagnostics.Stats["LambdaVoltage"].Min >= lambdaLow && diagnostics.Stats["LambdaVoltage"].Max <= lambdaHigh {
			diagnostics.Analysis.LambdaFault = false
		} else {
			diagnostics.Analysis.LambdaFault = true
		}
	}
}
