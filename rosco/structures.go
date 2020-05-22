package rosco

var responseMap = make(map[string][]byte)

// MemsData is the mems information computed from dataframes 0x80 and 0x7d
type (
	MemsData struct {
		Time                     string
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
		IdleHot                  uint8
		Uk8011                   uint8
		IACPosition              uint8
		IdleSpeedDeviation       uint16
		IgnitionAdvanceOffset80  uint8
		IgnitionAdvance          float32
		CoilTime                 float32
		CrankshaftPositionSensor bool
		Uk801a                   uint8
		Uk801b                   uint8
		IgnitionSwitch           bool
		ThrottleAngle            uint8
		Uk7d03                   uint8
		AirFuelRatio             float32
		DTC2                     uint8
		LambdaVoltage            uint8
		LambdaFrequency          uint8
		LambdaDutycycle          uint8
		LambdaStatus             uint8
		ClosedLoop               bool
		LongTermFuelTrim         uint8
		ShortTermFuelTrim        uint8
		FuelTrimCorrection       uint8
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

		CoolantTempSensorFault   bool
		IntakeAirTempSensorFault bool
		FuelPumpCircuitFault     bool
		ThrottlePotCircuitFault  bool

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
		BytesinFrame             uint8 // 7dx00
		IgnitionSwitch           uint8 // 7dx01
		ThrottleAngle            uint8 // 7dx03
		Uk7d03                   uint8 // 7dx03
		AirFuelRatio             uint8 // 7dx04
		Dtc2                     uint8 // 7dx05
		LambdaVoltage            uint8 // 7dx06
		LambdaFrequency          uint8 // 7dx07
		LambdaDutyCycle          uint8 // 7dx08
		LambdaStatus             uint8 // 7dx09
		LoopIndicator            uint8 // 7dx0A
		LongTermFuelTrim         uint8 // 7dx0B
		ShortTermFuelTrim        uint8 // 7dx0C
		CarbonCanisterPurgeValve uint8 // 7dx0D
		Dtc3                     uint8 // 7dx0E
		IdleBasePos              uint8 // 7dx0F
		Uk7d10                   uint8 // 7dx10
		Dtc4                     uint8 // 7dx11
		IgnitionAdvanceOffset7d  uint8 // 7dx12
		IdleSpeedOffset          uint8 // 7dx13
		Uk7d14                   uint8 // 7dx14
		Uk7d15                   uint8 // 7dx15
		Dtc5                     uint8 // 7dx16
		Uk7d17                   uint8 // 7dx17
		Uk7d18                   uint8 // 7dx18
		Uk7d19                   uint8 // 7dx19
		Uk7d1a                   uint8 // 7dx1a
		Uk7d1b                   uint8 // 7dx1b
		Uk7d1c                   uint8 // 7dx1c
		Uk7d1d                   uint8 // 7dx1d
		Uk7d1e                   uint8 // 7dx1e
		JackCount                uint8 // 7dx1f
	}
)

type (
	// DataFrame80 data sequence returned by the ECU in reply to the command 0x80.
	// This structure represents the raw data from the ECU
	//
	DataFrame80 struct {
		Command                  uint8
		BytesinFrame             uint8  // 80x00
		EngineRpm                uint16 // 80x01 - 80x02
		CoolantTemp              uint8  // 80x03
		AmbientTemp              uint8  // 80x04
		IntakeAirTemp            uint8  // 80x05
		FuelTemp                 uint8  // 80x06
		ManifoldAbsolutePressure uint8  // 80x07
		BatteryVoltage           uint8  // 80x08
		ThrottlePotSensor        uint8  // 80x09
		IdleSwitch               uint8  // 80x0A
		AirconSwitch             uint8  // 80x0B
		ParkNeutralSwitch        uint8  // 80x0C
		Dtc0                     uint8  // 80x0D  Bit 0: Coolant temp sensor fault (Code 1) Bit 1: Inlet air temp sensor fault (Code 2)
		Dtc1                     uint8  // 80x0E  Bit 1: Fuel pump circuit fault (Code 10)  Bit 7: Throttle pot circuit fault (Code 16)
		IdleSetPoint             uint8  // 80x0F
		IdleHot                  uint8  // 80x10
		Uk8011                   uint8  // 80x11
		IacPosition              uint8  // 80x12
		IdleSpeedDeviation       uint16 // 80x13 - 80x14
		IgnitionAdvanceOffset80  uint8  // 80x15
		IgnitionAdvance          uint8  // 80x16
		CoilTime                 uint16 // 80x17 - 80x18
		CrankshaftPositionSensor uint8  // 80x19
		Uk801a                   uint8  // 80x1A
		Uk801b                   uint8  // 80x1B
	}
)

const (
	// AirSensorFaultCode 0x80 DTC0 Fault
	AirSensorFaultCode = byte(0b00000001)
	// CoolantSensorFaultCode 0x80 DTC0 Fault
	CoolantSensorFaultCode = byte(0b00000010)
	// FuelPumpFaultCode 0x80 DTC1 Fault
	FuelPumpFaultCode = byte(0b00000001)
	// ThrottlePotFaultCode 0x80 DTC1 Fault
	ThrottlePotFaultCode = byte(0b01000000)
)
