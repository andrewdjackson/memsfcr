package rosco

import (
	"encoding/hex"
)

var Heartbeat = []byte{0xf4}
var Dataframe80 = []byte{0x80}
var Dataframe7d = []byte{0x7d}
var InitCommandA = []byte{0xca}
var InitCommandB = []byte{0x75}
var InitECUID = []byte{0xd0}

// MemsData struct
type (
	MemsData struct {
		Id string `json:"id"`

		// dataframe 0x80

		EngineRPM                uint16 `json:"80x01-02_engine-rpm"`
		CoolantTemp              uint8
		AmbientTemp              uint8
		IntakeAirTemp            uint8
		FuelTemp                 uint8
		MapKpa                   float32
		BatteryVoltage           float32
		ThrottlePotVoltage       float32
		IdleSwitch               uint8
		uk1                      uint8
		ParkNeutralSwitch        uint8
		FaultCodes               uint8
		IdleSetPoint             uint8
		IdleHot                  uint8
		uk2                      uint8
		IACPosition              uint8
		IdleError                uint16
		IgnitionAdvanceOffset    uint8
		IgnitionAdvance          uint8
		CoilTime                 uint16
		CrankshaftPositionSensor uint8
		uk4                      uint8
		uk5                      uint8

		CoolantTempSensorFault   bool
		IntakeAirTempSensorFault bool
		FuelPumpCircuitFault     bool
		ThrottlePotCircuitFault  bool

		// dataframe 0x7d

		IgnitionSwitch          uint8
		ThottleAngle            uint8
		uk6                     uint8
		AirFuelRatio            uint8
		DTC2                    uint8
		LambdaVoltage           uint8
		LambdaSensorFrequency   uint8
		LambdaSensorDutycycle   uint8
		LambdaSensorStatus      uint8
		ClosedLoop              uint8
		LongTermFuelTrim        uint8
		ShortTermFuelTrim       uint8
		CarbonCanisterDutycycle uint8
		DTC3                    uint8
		IdleBasePosition        uint8
		uk7                     uint8
		DTC4                    uint8
		IgnitionAdvance2        uint8
		IdleSpeedOffset         uint8
		IdleError2              uint8
		uk10                    uint8
		DTC5                    uint8
		uk11                    uint8
		uk12                    uint8
		uk13                    uint8
		uk14                    uint8
		uk15                    uint8
		uk16                    uint8
		uk1A                    uint8
		uk1B                    uint8
		uk1C                    uint8

		DataFrame80 string `json:"dataframe80"`
		DataFrame7d string `json:"dataframe7d"`
	}
)

var rosco = make(map[string][]byte)

func init() {
	rosco["0a"] = []byte{0x0A}
	rosco["80"] = []byte{0x80, 0x1C, 0x03, 0x5B, 0x8B, 0xFF, 0x56, 0xFF, 0x22, 0x8B, 0x1D, 0x00, 0x10, 0x01, 0x00, 0x00, 0x00, 0x24, 0x90, 0x2E, 0x00, 0x03, 0x00, 0x48, 0x06, 0x61, 0x10, 0x00, 0x00}
	rosco["7d"] = []byte{0x7d, 0x20, 0x10, 0x0D, 0xFF, 0x92, 0x00, 0x69, 0xFF, 0xFF, 0x00, 0x00, 0x96, 0x64, 0x00, 0xFF, 0x34, 0xFF, 0xFF, 0x30, 0x80, 0x7F, 0xFE, 0xFF, 0x19, 0x00, 0x1E, 0x80, 0x26, 0x40, 0x34, 0xC0, 0x1A}
	rosco["d0"] = []byte{0xD0, 0x99, 0x00, 0x03, 0x03}
	rosco["ca"] = []byte{0xCA}
	rosco["75"] = []byte{0x75}
	rosco["f0"] = []byte{0xF0, 0x00}
	rosco["f4"] = []byte{0xF4, 0x00}
}

func GetResponseSize(command []byte) int {
	c := hex.EncodeToString(command)
	r := rosco[c]

	if r != nil {
		return len(r) + 1
	}

	// default data returned is 2 bytes (echo of command and status)
	return 3
}
