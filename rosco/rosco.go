package rosco

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	"github.com/andrewdjackson/readmems/utils"
	"github.com/tarm/serial"
)

// MemsCommandResponse communication pair
type MemsCommandResponse struct {
	Command       []byte
	Response      []byte
	MemsDataFrame MemsData
}

// MemsConnection communtication structure for MEMS
type MemsConnection struct {
	// SerialPort the serial connection
	SerialPort      *serial.Port
	portReader      *bufio.Reader
	ECUID           []byte
	command         []byte
	response        []byte
	SendToECU       chan MemsCommandResponse
	ReceivedFromECU chan MemsCommandResponse
	Connected       bool
	Initialised     bool
	Exit            bool
}

// package init function
func init() {
	// Response formats for commands that do not respond with the format [COMMAND][VALUE]
	// Generally these are either part of the initialisation sequence or are ECU data frames
	responseMap["0a"] = []byte{0x0A}
	responseMap["ca"] = []byte{0xCA}
	responseMap["75"] = []byte{0x75}

	// Format for DataFrames starts with [Command Echo][Data Size][Data Bytes (28 for 0x80 and 32 for 0x7D)]
	responseMap["80"] = []byte{0x80, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B}
	responseMap["7d"] = []byte{0x7d, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F}
	responseMap["d0"] = []byte{0xD0, 0x99, 0x00, 0x03, 0x03}

	// generic response, expect command and single byte response
	responseMap["00"] = []byte{0x00, 0x00}
}

// NewMemsConnection creates a new mems structure
func NewMemsConnection() *MemsConnection {
	m := &MemsConnection{}
	m.Connected = false
	m.Initialised = false
	m.SendToECU = make(chan MemsCommandResponse, 2)
	m.ReceivedFromECU = make(chan MemsCommandResponse)

	return m
}

// ConnectAndInitialiseECU connect and initialise the ECU
func (mems *MemsConnection) ConnectAndInitialiseECU(config *ReadmemsConfig) {
	if !mems.Connected {
		mems.connect(config.Port)
		if mems.Connected {
			mems.initialise()
		}
	}
}

// connect to MEMS via serial port
func (mems *MemsConnection) connect(port string) {
	// connect to the ecu
	c := &serial.Config{Name: port, Baud: 9600}

	utils.LogI.Println("Opening ", port)

	s, err := serial.OpenPort(c)
	if err != nil {
		utils.LogI.Printf("%s", err)
	} else {
		utils.LogI.Println("Listening on ", port)

		mems.SerialPort = s
		mems.SerialPort.Flush()

		mems.Connected = true
	}
}

// checks the first byte of the response against the sent command
func (mems *MemsConnection) isCommandEcho() bool {
	return mems.command[0] == mems.response[0]
}

// initialises the connection to the ECU
// The initialisation sequence is as follows:
//
// 1. Send command CA (MEMS_InitCommandA)
// 2. Recieve response CA
// 3. Send command 75 (MEMS_InitCommandB)
// 4. Recieve response 75
// 5. Send request ECU ID command D0 (MEMS_InitECUID)
// 6. Recieve response D0 XX XX XX XX
//
func (mems *MemsConnection) initialise() {
	if mems.SerialPort != nil {
		mems.SerialPort.Flush()

		mems.writeSerial(MEMS_InitCommandA)
		mems.readSerial()

		mems.writeSerial(MEMS_InitCommandB)
		mems.readSerial()

		mems.writeSerial(MEMS_Heartbeat)
		mems.readSerial()

		mems.writeSerial(MEMS_InitECUID)
		mems.ECUID = mems.readSerial()
	}

	mems.Initialised = true
}

// readSerial read from MEMS
// read 1 byte at a time until we have all the expected bytes
func (mems *MemsConnection) readSerial() []byte {
	var n int
	var e error

	size := mems.getResponseSize(mems.command)

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
				utils.LogI.Printf("error %s", e)
			} else {
				// append the read bytes to the data frame
				data = append(data, b[:n]...)
			}

			// increment by the number of bytes read
			count = count + n
			if count > size {
				utils.LogI.Printf("data frame size mismatch (received %d, expected %d)", count, size)
			}
		}
	}

	utils.LogI.Printf("ECU [%d] < %x", n, data)
	mems.response = data

	if !mems.isCommandEcho() {
		utils.LogI.Printf("Expecting command echo (%x)\n", mems.command)
	}

	return data
}

// writeSerial write to MEMS
func (mems *MemsConnection) writeSerial(data []byte) {
	if mems.SerialPort != nil {
		// save the sent command
		mems.command = data

		// write the response to the code reader
		n, e := mems.SerialPort.Write(data)

		if e != nil {
			utils.LogI.Printf("FCR Send Error %s", e)
		}

		if n > 0 {
			utils.LogI.Printf("FCR > %x", data)
		}
	}
}

// ListenSendToECUChannelLoop for commands to be sent to the ECU
func (mems *MemsConnection) ListenSendToECUChannelLoop() {
	for {
		// wait for messages to be sent to the ECU
		utils.LogI.Printf(">>> waiting for mems command from the channel")
		m := <-mems.SendToECU
		utils.LogI.Printf(">>> mems command retrieved from the channel")
		// send the command
		reponse := mems.sendCommand(m.Command)
		// send back on the channel
		go mems.sendRecievedDataToChannel(reponse)
	}
}

func (mems *MemsConnection) sendRecievedDataToChannel(data []byte) {
	var m MemsCommandResponse
	m.Response = data

	utils.LogI.Printf("sending mems response to the channel")

	mems.ReceivedFromECU <- m
}

// sends a command and returns the response
func (mems *MemsConnection) sendCommand(cmd []byte) []byte {
	mems.writeSerial(cmd)
	return mems.readSerial()
}

// ReadMemsData reads the raw dataframes and returns structured data
func (mems *MemsConnection) ReadMemsData() {
	// read the raw dataframes
	d80, d7d := mems.readRaw()

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
	memsdata := MemsData{
		Time:                     t.Format("15:04:05.000"),
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

	// run as a go routine so it doesn't block this function completing
	go mems.sendMemsDataToChannel(memsdata)
}

func (mems *MemsConnection) sendMemsDataToChannel(memsdata MemsData) {
	var m MemsCommandResponse
	m.MemsDataFrame = memsdata

	utils.LogI.Printf("sending mems dataframe to the channel")

	mems.ReceivedFromECU <- m
}

// readRaw reads dataframe 80 and then dataframe 7d as raw byte arrays
func (mems *MemsConnection) readRaw() ([]byte, []byte) {
	mems.writeSerial(MEMS_ReqData80)
	dataframe80 := mems.readSerial()

	mems.writeSerial(MEMS_ReqData7D)
	dataframe7d := mems.readSerial()

	return dataframe80, dataframe7d
}

// getResponseSize returns the expected number of bytes for a given command
// The 'response' variable contains the formats for each command response pattern
// by default the response size is 2 bytes unless the command has a special format.
func (mems *MemsConnection) getResponseSize(command []byte) int {
	size := 2

	c := hex.EncodeToString(command)
	r := responseMap[c]

	if r != nil {
		size = len(r)
	} else {
		r = responseMap["00"]
		copy(r[0:], command)
	}

	utils.LogI.Printf("expecting %x -> o <- %x (%d)", command, r, size)
	return size
}
