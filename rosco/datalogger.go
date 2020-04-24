package rosco

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// MemsDataLogger logs the mems data to a CSV file
type MemsDataLogger struct {
	Filename string
	Logfile  *os.File
	IsOpen   bool
}

// NewMemsDataLogger logs the mems data to a CSV file
func NewMemsDataLogger() *MemsDataLogger {
	logger := &MemsDataLogger{}
	logger.setFilename()

	// check if this is a new file
	exist := logger.fileExists()

	// open the file
	logger.openFile()

	// if a new file then add a header to the file
	if !exist {
		logger.writeCSVHeader()
	}

	return logger
}

// WriteMemsDataToFile writes the mems data structure to a file
func (logger *MemsDataLogger) WriteMemsDataToFile(memsdata MemsData) {
	logger.writeCSVData(memsdata)
}

func (logger *MemsDataLogger) setFilename() {
	currentTime := time.Now()

	filename := fmt.Sprintf("logs/%s.csv", currentTime.Format("2006-01-02 15:04:05"))
	filename = strings.ReplaceAll(filename, ":", "")
	filename = strings.ReplaceAll(filename, " ", "-")

	logger.Filename = filename
}

// fileExists reports whether the named file or directory exists.
func (logger *MemsDataLogger) fileExists() bool {
	exists := false

	if _, err := os.Stat(logger.Filename); err != nil {
		if os.IsNotExist(err) {
			exists = false
		} else {
			exists = true
		}
	}

	LogI.Printf("%s exists %t", logger.Filename, exists)

	return exists
}

func (logger *MemsDataLogger) openFile() {
	var err error

	LogI.Printf("opening log file '%s'", logger.Filename)
	logger.Logfile, err = os.OpenFile(logger.Filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		LogE.Printf("%s", err)
		logger.IsOpen = false
	} else {
		logger.IsOpen = true
	}
}

// writeToFile will print any string of text to a file safely by
// checking for errors and syncing at the end.
func (logger *MemsDataLogger) writeToFile(data string) error {
	var err error

	if _, err = logger.Logfile.WriteString(data); err != nil {
		LogE.Printf("%s", err)
		return err
	}

	LogI.Printf("%s", data)

	return logger.Logfile.Sync()
}

func (logger *MemsDataLogger) writeCSVHeader() {
	header := "#time," +
		"80x01-02_engine-rpm,80x03_coolant_temp,80x04_ambient_temp,80x05_intake_air_temp,80x06_fuel_temp,80x07_map_kpa,80x08_battery_voltage,80x09_throttle_pot,80x0A_idle_switch,80x0B_uk1," +
		"80x0C_park_neutral_switch,80x0D-0E_fault_codes,80x0F_idle_set_point,80x10_idle_hot,80x11_uk2,80x12_iac_position,80x13-14_idle_error,80x15_ignition_advance_offset,80x16_ignition_advance,80x17-18_coil_time," +
		"80x19_crankshaft_position_sensor,80x1A_uk4,80x1B_uk5," +
		"7dx01_ignition_switch,7dx02_throttle_angle,7dx03_uk6,7dx04_air_fuel_ratio,7dx05_dtc2,7dx06_lambda_voltage,7dx07_lambda_sensor_frequency,7dx08_lambda_sensor_dutycycle,7dx09_lambda_sensor_status,7dx0A_closed_loop," +
		"7dx0B_long_term_fuel_trim,7dx0C_short_term_fuel_trim,7dx0D_carbon_canister_dutycycle,7dx0E_dtc3,7dx0F_idle_base_pos,7dx10_uk7,7dx11_dtc4,7dx12_ignition_advance2,7dx13_idle_speed_offset,7dx14_idle_error2," +
		"7dx14-15_uk10,7dx16_dtc5,7dx17_uk11,7dx18_uk12,7dx19_uk13,7dx1A_uk14,7dx1B_uk15,7dx1C_uk16,7dx1D_uk17,7dx1E_uk18,7dx1F_uk19\n"

	s := fmt.Sprintf("%s", header)
	logger.writeToFile(s)
}

func (logger *MemsDataLogger) writeCSVData(data MemsData) {
	s := fmt.Sprintf("%s,"+
		"%d,%d,%d,%d,%d,%f,%f,%f,%t,%t,"+
		"%t,%d,%d,%d,%d,%d,%d,%d,%f,%f,"+
		"%t,%d,%d,"+
		"%t,%d,%d,%d,%d,%d,%d,%d,%d,%t,"+
		"%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,"+
		"%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d\n",
		data.Time,
		data.EngineRPM,
		data.CoolantTemp,
		data.AmbientTemp,
		data.IntakeAirTemp,
		data.FuelTemp,
		data.ManifoldAbsolutePressure,
		data.BatteryVoltage,
		data.ThrottlePotSensor,
		data.IdleSwitch,
		data.AirconSwitch,
		data.ParkNeutralSwitch,
		data.DTC0,
		data.IdleSetPoint,
		data.IdleHot,
		data.Uk8011,
		data.IACPosition,
		data.IdleSpeedDeviation,
		data.IgnitionAdvanceOffset80,
		data.IgnitionAdvance,
		data.CoilTime,
		data.CrankshaftPositionSensor,
		data.Uk801a,
		data.Uk801b,
		data.IgnitionSwitch,
		data.ThrottleAngle,
		data.Uk7d03,
		data.AirFuelRatio,
		data.DTC2,
		data.LambdaVoltage,
		data.LambdaFrequency,
		data.LambdaDutycycle,
		data.LambdaStatus,
		data.ClosedLoop,
		data.LongTermFuelTrim,
		data.ShortTermFuelTrim,
		data.CarbonCanisterPurgeValve,
		data.DTC3,
		data.IdleBasePosition,
		data.Uk7d10,
		data.DTC4,
		data.IgnitionAdvanceOffset7d,
		data.IdleSpeedOffset,
		data.Uk7d14,
		data.Uk7d15,
		data.DTC5,
		data.Uk7d17,
		data.Uk7d18,
		data.Uk7d19,
		data.Uk7d1a,
		data.Uk7d1b,
		data.Uk7d1c,
		data.Uk7d1d,
		data.Uk7d1e,
		data.JackCount)

	logger.writeToFile(s)
}
