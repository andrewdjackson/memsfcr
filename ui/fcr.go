package ui

import (
	"github.com/andrewdjackson/readmems/rosco"
	"github.com/andrewdjackson/readmems/utils"
	"go.bug.st/serial.v1"
)

// MemsFCR structure
type MemsFCR struct {
	// FCR configuration
	config *rosco.ReadmemsConfig

	// dataframe read enabled / disabled
	paused bool

	// logging to file enabled / disabled
	logging bool

	// ECU represents the serial connection to the ECU
	ECU *rosco.MemsConnection

	// channel for communication to the ECU
	ToECUChannel chan rosco.MemsCommandResponse

	// channel for communication from the ECU
	FromECUChannel chan rosco.MemsCommandResponse
}

// NewMemsFCR creates an instance of a MEMs Fault Code Reader
func NewMemsFCR() *MemsFCR {
	memsfcr := &MemsFCR{}

	// set up the channels
	memsfcr.ToECUChannel = make(chan rosco.MemsCommandResponse)
	memsfcr.FromECUChannel = make(chan rosco.MemsCommandResponse)

	memsfcr.paused = false
	memsfcr.logging = false

	// read and apply the configuration
	memsfcr.readConfig()

	return memsfcr
}

// read the configuration file and apply the values
func (memsfcr *MemsFCR) readConfig() {
	memsfcr.config = rosco.ReadConfig()

	if memsfcr.config.Loop == "inf" {
		// infitite loop, so set loop count to a very big number
		memsfcr.config.Loop = "100000000"
	}

	// get the list of ports available
	memsfcr.config.Ports = append(memsfcr.config.Ports, memsfcr.config.Port)
	memsfcr.config.Ports = append(memsfcr.config.Ports, memsfcr.getSerialPorts()...)
}

// enumerate the available serial ports
// this won't enumerate virtual ports
func (memsfcr *MemsFCR) getSerialPorts() []string {
	ports, err := serial.GetPortsList()

	if err != nil {
		utils.LogE.Printf("error enumerating serial ports")
	}
	if len(ports) == 0 {
		utils.LogW.Printf("unable to find any serial ports")
	}
	for _, port := range ports {
		utils.LogI.Printf("found serial port %v", port)
	}

	return ports
}

// TxRxECULoop wraps the ECU send and recieve protocol
//
// The MEMs ECU uses a simple command / response protocol
// commands are single byte with a data response frame
//
// This loop should be run as a go routine the runs send and waiting for
// the ECU to respond. The channel blocking feature is used to manage flow
func (memsfcr *MemsFCR) TxRxECULoop() {
	for {
		// block waiting for an FCR command to send to the ECU
		tx := <-memsfcr.ToECUChannel

		// block waiting for the command to be sent to the ECU
		memsfcr.ECU.SendToECU <- tx

		// block waiting for the response to be received from the ECU
		rx := <-memsfcr.ECU.ReceivedFromECU

		// block waiting for the response to be collected for processing
		memsfcr.FromECUChannel <- rx
	}
}
