package rosco

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/tarm/serial"
	"math"
	"time"
)

// Mems communtication structure for MEMS
type Mems struct {
	// SerialPort the serial connection
	SerialPort  *serial.Port
	portReader  *bufio.Reader
	ECUID       []byte
	Command     []byte
	Response    []byte
	Connected   bool
	Initialised bool
	Exit        bool
}

// New creates a new mems structure
func New() *Mems {
	m := &Mems{}
	m.Connected = false
	m.Initialised = false
	m.Exit = false
	return m
}

// ConnectAndInitialiseECU connect and initialise the ECU
func ConnectAndInitialiseECU(mems *Mems, config *ReadmemsConfig) {
	if !mems.Connected {
		MemsConnect(mems, config.Port)
		if mems.Connected {
			MemsInitialise(mems)
		}
	}
}

// MemsConnect connect to MEMS via serial port
func MemsConnect(mems *Mems, port string) {
	// connect to the ecu
	c := &serial.Config{Name: port, Baud: 9600}

	LogI.Println("Opening ", port)

	s, err := serial.OpenPort(c)
	if err != nil {
		LogI.Printf("%s", err)
	} else {
		LogI.Println("Listening on ", port)

		mems.SerialPort = s
		mems.SerialPort.Flush()

		mems.Connected = true
	}
}

// checks the first byte of the response against the sent command
func isCommandEcho(mems *Mems) bool {
	return mems.Command[0] == mems.Response[0]
}

// MemsInitialise initialises the connection to the ECU
// The initialisation sequence is as follows:
//
// 1. Send command CA (MEMS_InitCommandA)
// 2. Recieve response CA
// 3. Send command 75 (MEMS_InitCommandB)
// 4. Recieve response 75
// 5. Send request ECU ID command D0 (MEMS_InitECUID)
// 6. Recieve response D0 XX XX XX XX
//
func MemsInitialise(mems *Mems) {
	if mems.SerialPort != nil {
		mems.SerialPort.Flush()

		MemsWriteSerial(mems, MEMS_InitCommandA)
		MemsReadSerial(mems)

		MemsWriteSerial(mems, MEMS_InitCommandB)
		MemsReadSerial(mems)

		MemsWriteSerial(mems, MEMS_Heartbeat)
		MemsReadSerial(mems)

		MemsWriteSerial(mems, MEMS_InitECUID)
		mems.ECUID = MemsReadSerial(mems)
	}

	mems.Initialised = true
}

// MemsReadSerial read from MEMS
// read 1 byte at a time until we have all the expected bytes
func MemsReadSerial(mems *Mems) []byte {
	var n int
	var e error

	size := GetResponseSize(mems.Command)

	// serial read buffer
	b := make([]byte, 100)

	//  data frame buffer
	data := make([]byte, 0)

	if mems.SerialPort != nil {

		// read all the expected bytes before returning the data
		for count := 0; count < size; {
			// wait for a response from MEMS
			n, e = mems.SerialPort.Read(b)

			if e != nil {
				LogI.Printf("error %s", e)
			} else {
				// append the read bytes to the data frame
				data = append(data, b[:n]...)
			}

			// increment by the number of bytes read
			count = count + n
			if count > size {
				LogI.Printf("data frame size mismatch (received %d, expected %d)", count, size)
			}
		}
	}

	LogI.Printf("ECU [%d] < %x", n, data)
	mems.Response = data

	if !isCommandEcho(mems) {
		LogI.Printf("Expecting command echo (%x)\n", mems.Command)
	}

	return data
}

// MemsWriteSerial write to MEMS
func MemsWriteSerial(mems *Mems, data []byte) {
	if mems.SerialPort != nil {
		// save the sent command
		mems.Command = data

		// write the response to the code reader
		n, e := mems.SerialPort.Write(data)

		if e != nil {
			LogI.Printf("FCR Send Error %s", e)
		}

		if n > 0 {
			LogI.Printf("FCR > %x", data)
		}
	}
}

// MemsSendCommand sends a command and returns the response
func MemsSendCommand(mems *Mems, cmd []byte) []byte {
	MemsWriteSerial(mems, cmd)
	return MemsReadSerial(mems)
}

// MemsRead reads the raw dataframes and returns structured data
func MemsRead(mems *Mems) MemsData {
	// read the raw dataframes
	d80, d7d := memsReadRaw(mems)

	// populate the DataFrame structure for command 0x80
	r := bytes.NewReader(d80)
	var df80 DataFrame80

	if err := binary.Read(r, binary.BigEndian, &df80); err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	// populate the DataFrame structure for command 0x7d
	r = bytes.NewReader(d7d)
	var df7d DataFrame7d

	if err := binary.Read(r, binary.BigEndian, &df7d); err != nil {
		fmt.Println("binary.Read failed:", err)
	}

	t := time.Now()

	// build the Mems Data frame using the raw data and applying the relevant
	// adjustments and calculations
	info := MemsData{
		Time:                     t.Format("15:04:05"),
		EngineRPM:                df80.EngineRpm,
		CoolantTemp:              df80.CoolantTemp - 55,
		AmbientTemp:              df80.AmbientTemp - 55,
		IntakeAirTemp:            df80.IntakeAirTemp - 55,
		FuelTemp:                 df80.FuelTemp - 55,
		ManifoldAbsolutePressure: float32(df80.ManifoldAbsolutePressure),
		BatteryVoltage:           float32(df80.BatteryVoltage / 10),
		ThrottlePotSensor:        float32(df80.ThrottlePotSensor) * 0.02,
		IdleSwitch:               df80.IdleSwitch == 1,
		AirconSwitch:             df80.AirconSwitch == 1,
		ParkNeutralSwitch:        df80.ParkNeutralSwitch == 1,
		DTC0:                     df80.Dtc0,
		DTC1:                     df80.Dtc1,
		IdleSetPoint:             df80.IdleSetPoint,
		IdleHot:                  df80.IdleHot,
		IACPosition:              (uint8(math.Round(float64(df80.IacPosition) / 1.8))),
		IdleSpeedDeviation:       df80.IdleSpeedDeviation,
		IgnitionAdvanceOffset80:  df80.IgnitionAdvanceOffset80,
		IgnitionAdvance:          (float32(df80.IgnitionAdvance) / 2) - 24,
		CoilTime:                 float32(df80.CoilTime) * 0.002,
		CrankshaftPositionSensor: df80.CrankshaftPositionSensor != 0,
		CoolantTempSensorFault:   df80.Dtc0&0x01 != 0,
		IntakeAirTempSensorFault: df80.Dtc0&0x02 != 0,
		FuelPumpCircuitFault:     df80.Dtc1&0x02 != 0,
		ThrottlePotCircuitFault:  df80.Dtc1&0x80 != 0,
		IgnitionSwitch:           df7d.IgnitionSwitch != 0,
		ThrottleAngle:            df7d.ThrottleAngle * 6 / 10,
		AirFuelRatio:             df7d.AirFuelRatio / 10,
		DTC2:                     df7d.Dtc2,
		LambdaVoltage:            df7d.LambdaVoltage * 5,
		LambdaFrequency:          df7d.LambdaFrequency,
		LambdaDutycycle:          df7d.LambdaDutyCycle,
		LambdaStatus:             df7d.LambdaStatus,
		ClosedLoop:               df7d.LoopIndicator != 0,
		LongTermFuelTrim:         df7d.LongTermFuelTrim,
		ShortTermFuelTrim:        df7d.ShortTermFuelTrim,
		FuelTrimCorrection:       df7d.ShortTermFuelTrim - 100,
		CarbonCanisterPurgeValve: df7d.CarbonCanisterPurgeValve,
		DTC3:                     df7d.Dtc3,
		IdleBasePosition:         df7d.IdleBasePos,
		DTC4:                     df7d.Dtc4,
		IgnitionAdvanceOffset7d:  df7d.IgnitionAdvanceOffset7d - 48,
		IdleSpeedOffset:          ((df7d.IdleSpeedOffset - 128) * 25),
		DTC5:                     df7d.Dtc5,
		JackCount:                df7d.JackCount,
		Dataframe80:              hex.EncodeToString(d80),
		Dataframe7d:              hex.EncodeToString(d7d),
	}

	return info
}

// MemsReadRaw reads dataframe 80 and then dataframe 7d as raw byte arrays
func memsReadRaw(mems *Mems) ([]byte, []byte) {
	MemsWriteSerial(mems, MEMS_ReqData80)
	dataframe80 := MemsReadSerial(mems)

	MemsWriteSerial(mems, MEMS_ReqData7D)
	dataframe7d := MemsReadSerial(mems)

	return dataframe80, dataframe7d
}
