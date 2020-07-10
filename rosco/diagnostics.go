package rosco

import (
	"reflect"

	"github.com/andrewdjackson/memsfcr/utils"
)

const (
	minIdleColdRPM      = 900  // Minimum expected RPM when running at Idle when cold
	maxIdleColdRPM      = 1200 // Maximum expected RPM when running at Idle when cold
	minIdleWarmRPM      = 700  // Minimum expected RPM when running at Idle when warm
	maxIdleWarmRPM      = 900  // Maximum expected RPM when running at Idle when warm
	minIdleMap          = 30   // Minimum MAP reading when the engine is running
	maxIdleMap          = 60   // Maximum MAP reading when the engine is running
	minMAPEngineOff     = 95   // Minumum MAP reading when the engine to not running
	engineOperatingTemp = 80   // Engine is at operating temp when coolant temp > 80C
	bestAFR             = 14.7 // Ideal Air to Fuel ratio
	lambdaLow           = 10   // Lambda minimum operating voltage
	lambdaHigh          = 900  // Lambda maximum operating voltage
	maxIdleError        = 50   // Max Idle Error
	maxSamples          = 30   // ~30 seconds
	minIAC              = 30   // Minimum normal operation steps for the IAC / Stepper Motor
	maxIAC              = 160  // Maximum normal operation steps for the IAC / Stepper Motor
)

const (
	codeOptimal      = "optimal"
	codeVacuumPipe   = "vacuum"
	codeThermostat   = "thermostat"
	codeCoolant      = "coolant"
	codeStepperMotor = "stepper"
	codeLambda       = "lambda"
	codeStepperMax   = "steppermax"
	codeStepperMin   = "steppermin"
)

// MemsAnalysisReport is the output from running the analysis
type MemsAnalysisReport struct {
	AnalysisCode             string
	IsEngineRunning          bool
	IsEngineWarming          bool
	IsAtOperatingTemp        bool
	IsEngineIdle             bool
	IsEngineIdleFault        bool
	IsCruising               bool
	IsClosedLoop             bool
	ClosedLoopExpected       bool
	MapFault                 bool
	VacuumFault              bool
	IdleAirControlFault      bool
	LambdaFault              bool
	CoolantTempSensorFault   bool
	IntakeAirTempSensorFault bool
	FuelPumpCircuitFault     bool
	ThrottlePotCircuitFault  bool
	IACPosition              int
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
	Stats map[string]utils.Stats
	// Analysis report
	Analysis MemsAnalysisReport
}

// NewMemsDiagnostics generates diagnostic reports
func NewMemsDiagnostics() *MemsDiagnostics {
	diagnostics := &MemsDiagnostics{}
	diagnostics.Dataset = []MemsData{}
	diagnostics.Analysis = MemsAnalysisReport{}
	diagnostics.Stats = make(map[string]utils.Stats)

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

	// get samples and associated stats for named metrics
	diagnostics.Stats["CoolantTemp"] = diagnostics.GetMetricStatistics("CoolantTemp")
	diagnostics.Stats["EngineRPM"] = diagnostics.GetMetricStatistics("EngineRPM")
	diagnostics.Stats["ManifoldAbsolutePressure"] = diagnostics.GetMetricStatistics("ManifoldAbsolutePressure")
	diagnostics.Stats["LambdaVoltage"] = diagnostics.GetMetricStatistics("LambdaVoltage")
	diagnostics.Stats["AirFuelRatio"] = diagnostics.GetMetricStatistics("AirFuelRatio")
	diagnostics.Stats["IACPosition"] = diagnostics.GetMetricStatistics("IACPosition")

	// default analysis outcome
	diagnostics.Analysis.AnalysisCode = codeOptimal

	// apply ECU detected faults
	diagnostics.Analysis.CoolantTempSensorFault = diagnostics.CurrentData.CoolantTempSensorFault
	diagnostics.Analysis.FuelPumpCircuitFault = diagnostics.CurrentData.FuelPumpCircuitFault
	diagnostics.Analysis.ThrottlePotCircuitFault = diagnostics.CurrentData.ThrottlePotCircuitFault
	diagnostics.Analysis.IntakeAirTempSensorFault = diagnostics.CurrentData.IntakeAirTempSensorFault

	// check the status of the sensors and running parameters
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
func (diagnostics *MemsDiagnostics) GetMetricStatistics(metricName string) utils.Stats {
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
	return *utils.NewStats(metricName, metricSample)
}

// IsEngineWarm uses the current engine temperature and the standard deviation in the sample to determine the
// stability of the temperature. If the reading is at the designated thermostat temp (88C) and the std deviation
// is low then deem the engine to be running at operating temperature
func (diagnostics *MemsDiagnostics) checkIsEngineWarm() {
	diagnostics.Analysis.IsAtOperatingTemp = (diagnostics.Stats["CoolantTemp"].Value >= engineOperatingTemp && diagnostics.Stats["CoolantTemp"].Stddev < 5)
}

// if the last reading of engine RPM is 0 then the engine is not running
// we don't use the sample set as the engine may have recently been stopped
func (diagnostics *MemsDiagnostics) checkIsEngineRunning() {
	diagnostics.Analysis.IsEngineRunning = !(diagnostics.CurrentData.EngineRPM == 0)
}

// IsIdle determines the correct idle speed parameters based on whether the engine is warm or cold
// if the RPM is within the parameters for the sample period then the engine is deemed to be at Idle
func (diagnostics *MemsDiagnostics) checkIsEngineIdle() {
	// if the engine is running and the RPM is stable, then deem to be in an idle state
	// if the engine RPM is above idle but stable then we're cruising
	if diagnostics.Analysis.IsEngineRunning && diagnostics.Stats["EngineRPM"].Stddev <= 10 {
		diagnostics.Analysis.IsEngineIdle = true
		diagnostics.Analysis.IsCruising = diagnostics.Stats["EngineRPM"].Mean > maxIdleWarmRPM
	}

	if diagnostics.Analysis.IsAtOperatingTemp {
		// use warm idle settings
		diagnostics.Analysis.IsEngineIdleFault = !(diagnostics.Stats["EngineRPM"].Mean >= minIdleWarmRPM && diagnostics.Stats["EngineRPM"].Mean <= maxIdleWarmRPM)
		diagnostics.Analysis.IsEngineWarming = false
	} else {
		// use cold idle settings
		diagnostics.Analysis.IsEngineIdleFault = !(diagnostics.Stats["EngineRPM"].Mean >= minIdleColdRPM && diagnostics.Stats["EngineRPM"].Mean <= maxIdleColdRPM)
		diagnostics.Analysis.IsEngineWarming = true
	}
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
	// expecting ECU to switch to closed loop when at operating temperature and either idling or cruising
	diagnostics.Analysis.ClosedLoopExpected = diagnostics.Analysis.IsAtOperatingTemp && (diagnostics.Analysis.IsEngineIdle || diagnostics.Analysis.IsCruising)
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
// A high number of steps indicates that the ECU is attempting to close the stepper or reduce the airflow
// a low number would indicate the inability to increase airflow

func (diagnostics *MemsDiagnostics) checkIdleAirControl() {
	if diagnostics.Analysis.IsEngineRunning {
		// IAC fault if the idle offset exceeds the max error, yet the IAC position remains at 0
		if diagnostics.CurrentData.IdleSpeedDeviation >= maxIdleError && diagnostics.CurrentData.IACPosition == 0 {
			diagnostics.Analysis.IdleAirControlFault = true
			diagnostics.Analysis.AnalysisCode = codeStepperMotor
		}

		// IAC fault if the stepper motor is wide open or closed
		if diagnostics.Stats["IACPosition"].Mean <= minIAC {
			diagnostics.Analysis.IdleAirControlFault = true
			diagnostics.Analysis.AnalysisCode = codeStepperMin
		}

		if diagnostics.Stats["IACPosition"].Mean >= maxIAC {
			diagnostics.Analysis.IdleAirControlFault = true
			diagnostics.Analysis.AnalysisCode = codeStepperMax
		}
	} else {
		diagnostics.Analysis.IdleAirControlFault = false
	}
}

//  At 2000 rpm it should be switching rapidly between the minimum and maximum figures as the MEMS controls the engine conditions.????
func (diagnostics *MemsDiagnostics) checkLambdaStatus() {
	if diagnostics.Analysis.IsEngineRunning && diagnostics.Analysis.IsClosedLoop {
		if diagnostics.Stats["LambdaVoltage"].Min >= lambdaLow && diagnostics.Stats["LambdaVoltage"].Max <= lambdaHigh {
			// lambda oprtating within exepcted parameters
			diagnostics.Analysis.LambdaFault = false
		} else {
			diagnostics.Analysis.LambdaFault = true
		}
	}
}

/**
 * Repeatedly send command to open or close the idle air control valve until
 * it is in the desired position. The valve does not necessarily move one full
 * step per serial command, depending on the rate at which the commands are
 * issued.
 */
/*
 bool mems_move_iac(mems_info *info, uint8_t desired_pos)
 {
   bool status = false;
   uint16_t attempts = 0;
   uint8_t current_pos = 0;
   actuator_cmd cmd;

   // read the current IAC position, and only take action
   // if we're not already at the desired point
   if (mems_read_iac_position(info, &current_pos))
   {
	 if ((desired_pos < current_pos) ||
		 ((desired_pos > current_pos) && (current_pos < IAC_MAXIMUM)))
	 {
	   cmd = (desired_pos > current_pos) ? MEMS_OpenIAC : MEMS_CloseIAC;

	   do
	   {
		 status = mems_test_actuator(info, cmd, &current_pos);
		 attempts += 1;
	   } while (status && (current_pos != desired_pos) && (attempts < 300));
	 }
   }

   status = (desired_pos == current_pos);

   return status;
 }*/
