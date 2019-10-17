package rosco

import (
	"encoding/hex"
	"fmt"
	"log"

	"andrewj.com/readmems/config"
	"github.com/tarm/serial"
)

// Mems communtication structure for MEMS
type Mems struct {
	// SerialPort the serial connection
	SerialPort *serial.Port
	ECUID      []byte
	command    []byte
	response   []byte
}

// New creates a new mems structure
func New() *Mems {
	return &Mems{}
}

// MemsConnect connect to MEMS via serial port
func MemsConnect(mems *Mems, readmemsConfig config.ReadmemsConfig) {
	// connect to the ecu
	c := &serial.Config{Name: readmemsConfig.Port, Baud: 9600}

	fmt.Println("Opening ", readmemsConfig.Port)

	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Listening on ", readmemsConfig.Port)

	mems.SerialPort = s
	mems.SerialPort.Flush()
}

// checks the first byte of the response against the sent command
func isCommandEcho(mems *Mems) bool {
	return mems.command[0] == mems.response[0]
}

// MemsInitialise initialise the connection
func MemsInitialise(mems *Mems) bool {
	if mems.SerialPort != nil {
		MemsWriteSerial(mems, InitCommandA)
		MemsReadSerial(mems)
		MemsWriteSerial(mems, InitCommandB)
		MemsReadSerial(mems)
		MemsWriteSerial(mems, Heartbeat)
		MemsReadSerial(mems)
		MemsWriteSerial(mems, InitECUID)
		mems.ECUID = MemsReadSerial(mems)
	}

	return true
}

// MemsReadSerial read from MEMS
func MemsReadSerial(mems *Mems) []byte {
	size := GetResponseSize(mems.command)
	data := make([]byte, size)

	if mems.SerialPort != nil {
		// wait for a response from MEMS
		n, e := mems.SerialPort.Read(data)

		if e != nil {
			log.Printf("error %s", e)
		}

		if n > 0 {
			log.Printf("read (%d): %x", n, data[:n])
		}
	}

	mems.response = data

	if !isCommandEcho(mems) {
		log.Fatal("Expecting command echo")
	}

	return data
}

// MemsWriteSerial write to MEMS
func MemsWriteSerial(mems *Mems, data []byte) {
	if mems.SerialPort != nil {
		// save the sent command
		mems.command = data

		// write the response to the code reader
		n, e := mems.SerialPort.Write(data)

		if e != nil {
			log.Printf("error %s", e)
		}

		if n > 0 {
			log.Printf("write: %x", data)
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
	dataframe80, dataframe7d := MemsReadRaw(mems)

	info := MemsData{
		EngineRPM:   0,
		DataFrame80: hex.EncodeToString(dataframe80),
		DataFrame7d: hex.EncodeToString(dataframe7d),
	}

	return info
}

// MemsReadRaw reads dataframe 80 and then dataframe 7d as raw byte arrays
func MemsReadRaw(mems *Mems) ([]byte, []byte) {
	MemsWriteSerial(mems, Dataframe80)
	dataframe80 := MemsReadSerial(mems)

	MemsWriteSerial(mems, Dataframe7d)
	dataframe7d := MemsReadSerial(mems)

	return dataframe80, dataframe7d
}
