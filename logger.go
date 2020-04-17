package main

import (
	"andrewj.com/readmems/rosco"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// fileExists reports whether the named file or directory exists.
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// writeToFile will print any string of text to a file safely by
// checking for errors and syncing at the end.
func writeToFile(filename string, data string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)

	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(data); err != nil {
		return err
	}

	log.Printf(">> %s\n", data)

	return file.Sync()
}

func writeCSVHeader(filename string) {
	header := "#time," +
		"80x01-02_engine-rpm,80x03_coolant_temp,80x04_ambient_temp,80x05_intake_air_temp,80x06_fuel_temp,80x07_map_kpa,80x08_battery_voltage,80x09_throttle_pot,80x0A_idle_switch,80x0B_uk1," +
		"80x0C_park_neutral_switch,80x0D-0E_fault_codes,80x0F_idle_set_point,80x10_idle_hot,80x11_uk2,80x12_iac_position,80x13-14_idle_error,80x15_ignition_advance_offset,80x16_ignition_advance,80x17-18_coil_time," +
		"80x19_crankshaft_position_sensor,80x1A_uk4,80x1B_uk5," +
		"7dx01_ignition_switch,7dx02_throttle_angle,7dx03_uk6,7dx04_air_fuel_ratio,7dx05_dtc2,7dx06_lambda_voltage,7dx07_lambda_sensor_frequency,7dx08_lambda_sensor_dutycycle,7dx09_lambda_sensor_status,7dx0A_closed_loop," +
		"7dx0B_long_term_fuel_trim,7dx0C_short_term_fuel_trim,7dx0D_carbon_canister_dutycycle,7dx0E_dtc3,7dx0F_idle_base_pos,7dx10_uk7,7dx11_dtc4,7dx12_ignition_advance2,7dx13_idle_speed_offset,7dx14_idle_error2," +
		"7dx14-15_uk10,7dx16_dtc5,7dx17_uk11,7dx18_uk12,7dx19_uk13,7dx1A_uk14,7dx1B_uk15,7dx1C_uk16,7dx1D_uk17,7dx1E_uk18,7dx1F_uk19\n"

	s := fmt.Sprintf("%s", header)
	writeToFile(filename, s)
}

func writeCSVData(filename string, data rosco.MemsData) {
	s := fmt.Sprintf("%s,"+
		"%d,%d,%d,%d,%d,%f,%f,%f,%t,%d,"+
		"%t,%d,%d,%d,%d,%d,%d,%d,%f,%f,"+
		"%d,%d,%d,"+
		"%t,%d,%d,%d,%d,%d,%d,%d,%d,%t,"+
		"%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,"+
		"%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d\n",
		data.Time,
		data.EngineRPM,
		data.CoolantTemp,
		data.AmbientTemp,
		data.IntakeAirTemp,
		data.FuelTemp,
		data.MapKpa,
		data.BatteryVoltage,
		data.ThrottlePotVoltage,
		data.IdleSwitch,
		data.Uk1,
		data.ParkNeutralSwitch,
		data.FaultCodes,
		data.IdleSetPoint,
		data.IdleHot,
		data.Uk2,
		data.IACPosition,
		data.IdleError,
		data.IgnitionAdvanceOffset,
		data.IgnitionAdvance,
		data.CoilTime,
		data.CrankshaftPositionSensor,
		data.Uk4,
		data.Uk5,
		data.IgnitionSwitch,
		data.ThrottleAngle,
		data.Uk6,
		data.AirFuelRatio,
		data.DTC2,
		data.LambdaVoltage,
		data.LambdaSensorFrequency,
		data.LambdaSensorDutycycle,
		data.LambdaSensorStatus,
		data.ClosedLoop,
		data.LongTermFuelTrim,
		data.ShortTermFuelTrim,
		data.CarbonCanisterDutycycle,
		data.DTC3,
		data.IdleBasePosition,
		data.Uk7,
		data.DTC4,
		data.IgnitionAdvance2,
		data.IdleSpeedOffset,
		data.IdleError2,
		data.Uk10,
		data.DTC5,
		data.Uk11,
		data.Uk12,
		data.Uk13,
		data.Uk14,
		data.Uk15,
		data.Uk16,
		data.Uk1A,
		data.Uk1B,
		data.Uk1C)

	writeToFile(filename, s)
}

// WriteMemsDataToFile writes the mems data structure to a file
func WriteMemsDataToFile(memsdata rosco.MemsData) {
	currentTime := time.Now()
	filename := fmt.Sprintf("logs/%s.csv", currentTime.Format("2006-01-02 15:04:05"))
	filename = strings.ReplaceAll(filename, ":", "")
	filename = strings.ReplaceAll(filename, " ", "-")

	if !fileExists(filename) {
		writeCSVHeader(filename)
	}

	writeCSVData(filename, memsdata)

	//md, _ := json.Marshal(memsdata)
	//WriteToFile("output.cvs", string(md))
}
