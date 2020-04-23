package mems

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"log"
	"time"
)

const memsBaudRate = 9600

// ECUCommandResponse structure to pass data over the channel
type ECUCommandResponse struct {
	command  []byte
	response []byte
}

var ecuChannel = make(chan []byte)

// MemsDataChannel is the channel to pass back ECU responses
var MemsDataChannel = make(chan MemsData)

// ECUConnection contains all the configuration necessary
// to open a serial port
type ECUConnection struct {
	config     *serial.Config
	port       *serial.Port
	portReader *bufio.Reader
	stateChan  chan error
	ecuid      []byte
}

// ConnectAndInitialiseECU connect and initialise the ECU
func ConnectAndInitialiseECU(config *ReadmemsConfig) {
	// create a serial port connection
	c, err := NewECUConnection(config.Port)

	if err != nil {
		log.Printf("FatalError: %v", err)
		return
	}

	// establish a connection to the ECU and start the read loop
	go c.Start()
	go c.readLoop()

	c.ecuid = c.initialiseECU()

	go c.memsCommandResponseLoop()
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

func (c *ECUConnection) initialiseECU() []byte {
	c.Write(MEMS_InitCommandA)
	// use the channel blocking to wait for response
	data := <-ecuChannel

	c.Write(MEMS_InitCommandB)
	// use the channel blocking to wait for response
	data = <-ecuChannel

	c.Write(MEMS_Heartbeat)
	// use the channel blocking to wait for response
	data = <-ecuChannel

	c.Write(MEMS_InitECUID)
	// use the channel blocking to wait for response
	data = <-ecuChannel

	// return the ecu id
	return data
}

// NewECUConnection returns a pointer to a ECUConnection instance
func NewECUConnection(portPath string) (*ECUConnection, error) {
	config := serial.Config{Name: portPath, Baud: memsBaudRate, ReadTimeout: time.Nanosecond}
	port, err := serial.OpenPort(&config)
	if err != nil {
		return nil, err
	}
	portReader := bufio.NewReader(port)
	stateChan := make(chan error)

	return &ECUConnection{
		config:     &config,
		port:       port,
		portReader: portReader,
		stateChan:  stateChan}, nil
}

// Start initializes a read loop that attempts to reconnect
// when the connection is broken
func (c *ECUConnection) Start() {
	for {
		select {
		case err := <-c.stateChan:
			if err != nil {
				fmt.Printf("Error connecting to %s", c.config.Name)
				go c.initialize()
			} else {
				fmt.Printf(" | Connection to %s reestablished!", c.config.Name)
			}
		}
	}
}

func (c *ECUConnection) initialize() {
	c.port.Close()
	for {
		time.Sleep(time.Second)
		port, err := serial.OpenPort(c.config)
		if err != nil {
			continue
		}
		c.port = port
		c.portReader = bufio.NewReader(port)
		c.stateChan <- nil
		return
	}
}

// MemsRead reads the raw dataframes and returns structured data
func (c *ECUConnection) memsCommandResponseLoop() {
	for {
		c.readmems()
		time.Sleep(800 * time.Millisecond)
	}
}

func (c *ECUConnection) readmems() {
	// read the raw dataframes
	d80, d7d := c.readRaw()

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
	/*
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
	*/

	info := MemsData{
		Time:                     t.Format("15:04:05"),
		EngineRPM:                df80.EngineRpm,
		CoolantTemp:              df80.CoolantTemp,
		AmbientTemp:              df80.AmbientTemp,
		IntakeAirTemp:            df80.IntakeAirTemp,
		FuelTemp:                 df80.FuelTemp,
		ManifoldAbsolutePressure: float32(df80.ManifoldAbsolutePressure),
		BatteryVoltage:           float32(df80.BatteryVoltage),
		ThrottlePotSensor:        float32(df80.ThrottlePotSensor),
		IdleSwitch:               df80.IdleSwitch == 1,
		AirconSwitch:             df80.AirconSwitch == 1,
		ParkNeutralSwitch:        df80.ParkNeutralSwitch == 1,
		DTC0:                     df80.Dtc0,
		DTC1:                     df80.Dtc1,
		IdleSetPoint:             df80.IdleSetPoint,
		IdleHot:                  df80.IdleHot,
		IACPosition:              df80.IacPosition,
		IdleSpeedDeviation:       df80.IdleSpeedDeviation,
		IgnitionAdvanceOffset80:  df80.IgnitionAdvanceOffset80,
		IgnitionAdvance:          float32(df80.IgnitionAdvance),
		CoilTime:                 float32(df80.CoilTime),
		CrankshaftPositionSensor: df80.CrankshaftPositionSensor != 0,
		CoolantTempSensorFault:   df80.Dtc0&0x01 != 0,
		IntakeAirTempSensorFault: df80.Dtc0&0x02 != 0,
		FuelPumpCircuitFault:     df80.Dtc1&0x02 != 0,
		ThrottlePotCircuitFault:  df80.Dtc1&0x80 != 0,
		IgnitionSwitch:           df7d.IgnitionSwitch != 0,
		ThrottleAngle:            df7d.ThrottleAngle,
		AirFuelRatio:             df7d.AirFuelRatio,
		DTC2:                     df7d.Dtc2,
		LambdaVoltage:            df7d.LambdaVoltage,
		LambdaFrequency:          df7d.LambdaFrequency,
		LambdaDutycycle:          df7d.LambdaDutyCycle,
		LambdaStatus:             df7d.LambdaStatus,
		ClosedLoop:               df7d.LoopIndicator != 0,
		LongTermFuelTrim:         df7d.LongTermFuelTrim,
		ShortTermFuelTrim:        df7d.ShortTermFuelTrim,
		CarbonCanisterPurgeValve: df7d.CarbonCanisterPurgeValve,
		DTC3:                     df7d.Dtc3,
		IdleBasePosition:         df7d.IdleBasePos,
		DTC4:                     df7d.Dtc4,
		IgnitionAdvanceOffset7d:  df7d.IgnitionAdvanceOffset7d,
		IdleSpeedOffset:          df7d.IdleSpeedOffset,
		DTC5:                     df7d.Dtc5,
		JackCount:                df7d.JackCount,
		Dataframe80:              hex.EncodeToString(d80),
		Dataframe7d:              hex.EncodeToString(d7d),
	}

	MemsDataChannel <- info
}

// reads dataframe 80 and then dataframe 7d as raw byte arrays
func (c *ECUConnection) readRaw() ([]byte, []byte) {
	c.Write(MEMS_ReqData80)
	dataframe80 := <-ecuChannel

	c.Write(MEMS_ReqData7D)
	dataframe7d := <-ecuChannel

	return dataframe80, dataframe7d
}

// Read loop from serial port
func (c *ECUConnection) readLoop() {
	for {
		response, err := c.portReader.ReadBytes('\n')
		// report the error
		if err != nil && err != io.EOF {
			c.stateChan <- err
			return
		}
		if len(response) > 0 {
			ecuChannel <- response
			log.Printf("ECU: %x\r\n", response)
		}
	}
}

func (c *ECUConnection) read() {
	for {
		response, err := c.portReader.ReadBytes('\n')
		// report the error
		if err != nil && err != io.EOF {
			c.stateChan <- err
			return
		}
		if len(response) > 0 {
			log.Printf("ECU: %x\r\n", response)
		}
	}
}

func (c *ECUConnection) Write(message []byte) {
	_, err := c.port.Write(message)
	if err != nil {
		fmt.Printf("Error writing to serial port: %v ", err)
	} else {
		log.Printf("FCR: %x\r\n", message)
	}
}
