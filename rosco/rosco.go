package rosco

import (
	"encoding/hex"
)

// MemsData is the mems information computed from dataframes 0x80 and 0x7d
type (
	MemsData struct {
		Time string

		// dataframe 0x80

		EngineRPM                uint16
		CoolantTemp              uint8
		AmbientTemp              uint8
		IntakeAirTemp            uint8
		FuelTemp                 uint8
		ManifoldAbsolutePressure float32
		BatteryVoltage           float32
		ThrottlePotSensor        float32
		ThrottlePosition         float32
		IdleSwitch               bool
		AirconSwitch             bool
		ParkNeutralSwitch        bool
		DTC0                     uint8
		DTC1                     uint8
		IdleSetPoint             uint8
		IdleDecay                uint8
		Uk8011                   uint8
		IACPosition              uint8
		IdleSpeedDeviation       uint16
		IgnitionAdvanceOffset80  uint8
		IgnitionAdvance          float32
		CoilTime                 float32
		CrankshaftPositionSensor bool
		Uk801a                   uint8
		Uk801b                   uint8

		CoolantTempSensorFault   bool
		IntakeAirTempSensorFault bool
		FuelPumpCircuitFault     bool
		ThrottlePotCircuitFault  bool

		// dataframe 0x7d

		IgnitionSwitch           bool
		ThrottleAngle            uint8
		Uk7d03                   uint8
		AirFuelRatio             uint8
		DTC2                     uint8
		LambdaVoltage            uint8
		LambdaFrequency          uint8
		LambdaDutycycle          uint8
		LambdaStatus             uint8
		ClosedLoop               bool
		LongTermFuelTrim         uint8
		ShortTermFuelTrim        uint8
		CarbonCanisterPurgeValve uint8
		DTC3                     uint8
		IdleBasePosition         uint8
		Uk7d10                   uint8
		DTC4                     uint8
		IgnitionAdvanceOffset7d  uint8
		IdleSpeedOffset          uint8
		Uk7d14                   uint8
		Uk7d15                   uint8
		DTC5                     uint8
		Uk7d17                   uint8
		Uk7d18                   uint8
		Uk7d19                   uint8
		Uk7d1a                   uint8
		Uk7d1b                   uint8
		Uk7d1c                   uint8
		Uk7d1d                   uint8
		Uk7d1e                   uint8
		JackCount                uint8

		Dataframe80 string `json:"dataframe80"`
		Dataframe7d string `json:"dataframe7d"`
	}
)

type (
	// DataFrame7d data sequence returned by the ECU in reply to the command 0x7D.
	// This structure represents the raw data from the ECU
	//
	DataFrame7d struct {
		Command                  uint8
		BytesinFrame             uint8
		IgnitionSwitch           uint8
		ThrottleAngle            uint8
		Uk7d03                   uint8
		AirFuelRatio             uint8
		Dtc2                     uint8
		LambdaVoltage            uint8
		LambdaFrequency          uint8
		LambdaDutyCycle          uint8
		LambdaStatus             uint8
		LoopIndicator            uint8
		LongTermFuelTrim         uint8
		ShortTermFuelTrim        uint8
		CarbonCanisterPurgeValve uint8
		Dtc3                     uint8
		IdleBasePos              uint8
		Uk7d10                   uint8
		Dtc4                     uint8
		IgnitionAdvanceOffset7d  uint8
		IdleSpeedOffset          uint8
		Uk7d14                   uint8
		Uk7d15                   uint8
		Dtc5                     uint8
		Uk7d17                   uint8
		Uk7d18                   uint8
		Uk7d19                   uint8
		Uk7d1a                   uint8
		Uk7d1b                   uint8
		Uk7d1c                   uint8
		Uk7d1d                   uint8
		Uk7d1e                   uint8
		JackCount                uint8
	}
)

type (
	// DataFrame80 data sequence returned by the ECU in reply to the command 0x80.
	// This structure represents the raw data from the ECU
	//
	DataFrame80 struct {
		Command                  uint8
		BytesinFrame             uint8
		EngineRpm                uint16
		CoolantTemp              uint8
		AmbientTemp              uint8
		IntakeAirTemp            uint8
		FuelTemp                 uint8
		ManifoldAbsolutePressure uint8
		BatteryVoltage           uint8
		ThrottlePotSensor        uint8
		IdleSwitch               uint8
		AirconSwitch             uint8
		ParkNeutralSwitch        uint8
		Dtc0                     uint8
		Dtc1                     uint8
		IdleSetPoint             uint8
		IdleDecay                uint8
		Uk8011                   uint8
		IacPosition              uint8
		IdleSpeedDeviation       uint16
		IgnitionAdvanceOffset80  uint8
		IgnitionAdvance          uint8
		CoilTime                 uint16
		CrankshaftPositionSensor uint8
		Uk801a                   uint8
		Uk801b                   uint8
	}
)

var response = make(map[string][]byte)

func init() {
	// Response formats for commands that do not respond with the format [COMMAND][VALUE]
	// Generally these are either part of the initialisation sequence or are ECU data frames
	response["0a"] = []byte{0x0A}
	response["80"] = []byte{0x80, 0x1C, 0x03, 0x5B, 0x8B, 0xFF, 0x56, 0xFF, 0x22, 0x8B, 0x1D, 0x00, 0x10, 0x01, 0x00, 0x00, 0x00, 0x24, 0x90, 0x2E, 0x00, 0x03, 0x00, 0x48, 0x06, 0x61, 0x10, 0x00, 0x00}
	response["7d"] = []byte{0x7d, 0x20, 0x10, 0x0D, 0xFF, 0x92, 0x00, 0x69, 0xFF, 0xFF, 0x00, 0x00, 0x96, 0x64, 0x00, 0xFF, 0x34, 0xFF, 0xFF, 0x30, 0x80, 0x7F, 0xFE, 0xFF, 0x19, 0x00, 0x1E, 0x80, 0x26, 0x40, 0x34, 0xC0, 0x1A}
	response["d0"] = []byte{0xD0, 0x99, 0x00, 0x03, 0x03}
	response["ca"] = []byte{0xCA}
	response["75"] = []byte{0x75}
	response["f0"] = []byte{0xF0, 0x00}
	response["f4"] = []byte{0xF4, 0x00}
}

// GetResponseSize returns the expected number of bytes for a given command
// The 'response' variable contains the formats for each command response pattern
// by default the response size is 2 bytes unless the command has a special format.
func GetResponseSize(command []byte) int {
	c := hex.EncodeToString(command)
	r := response[c]

	if r != nil {
		return len(r) + 1
	}

	// default data returned is 2 bytes (echo of command and status)
	return 3
}
