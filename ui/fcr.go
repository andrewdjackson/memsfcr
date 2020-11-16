package ui

import (
	"strings"

	"github.com/andrewdjackson/memsfcr/rosco"
	"github.com/andrewdjackson/memsfcr/utils"
	"go.bug.st/serial.v1"
)

// MemsFCR structure
type MemsFCR struct {
	// Version of the build
	Version string
	Build   string

	// Config FCR configuration
	Config *utils.ReadmemsConfig

	// Paused dataframe read enabled / disabled
	Paused bool

	// Logging to file enabled / disabled
	Logging bool

	// ECU represents the serial connection to the ECU
	ECU *rosco.MemsConnection

	// channel for communication to the ECU
	FCRSendToECU chan rosco.MemsCommandResponse

	// channel for communication from the ECU
	ECUSendToFCR chan rosco.MemsCommandResponse
}

// NewMemsFCR creates an instance of a MEMs Fault Code Reader
func NewMemsFCR(version string, build string) *MemsFCR {
	memsfcr := &MemsFCR{}
	memsfcr.Version = version
	memsfcr.Build = build

	// set up the channels
	memsfcr.FCRSendToECU = make(chan rosco.MemsCommandResponse)
	memsfcr.ECUSendToFCR = make(chan rosco.MemsCommandResponse)

	memsfcr.Paused = false
	memsfcr.Logging = true

	// read and apply the configuration
	memsfcr.readConfig()

	return memsfcr
}

// read the configuration file and apply the values
func (memsfcr *MemsFCR) readConfig() {
	memsfcr.Config = utils.ReadConfig()

	// override config build versions with current versions
	// if the values are not in development values
	if strings.Compare(memsfcr.Version, "0.0.0") != 0 {
		memsfcr.Config.Version = memsfcr.Version
		memsfcr.Config.Build = memsfcr.Build
	}

	if memsfcr.Config.Loop == "inf" {
		// infitite loop, so set loop count to a very big number
		memsfcr.Config.Loop = "100000000"
	}

	if memsfcr.Config.LogToFile == "true" {
		memsfcr.Logging = true
	}

	// get the list of ports available
	memsfcr.Config.Ports = append(memsfcr.Config.Ports, memsfcr.Config.Port)
	memsfcr.Config.Ports = append(memsfcr.Config.Ports, memsfcr.getSerialPorts()...)
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

// ConnectFCR connects the FCR to the ECU
// on successful connection the FCR runs the initialisation sequence
func (memsfcr *MemsFCR) ConnectFCR() bool {
	// only connect and initialise if the ECU hasn't already been
	// initialised. This seems to do odd things to the ECU if you
	// run the sequence once initialised
	if !memsfcr.ECU.Initialised {
		memsfcr.ECU.ConnectAndInitialiseECU(memsfcr.Config.Port)

		// start the ECU loop
		if memsfcr.ECU.Initialised {
			go memsfcr.ECU.ListenTxECUChannelLoop()
		} else {
			utils.LogW.Printf("ECU not initialised")
		}
	}

	return memsfcr.ECU.Initialised
}

// TxECU send the command to the ECU from the FCR
func (memsfcr *MemsFCR) TxECU(cmd []byte) {
	var c rosco.MemsCommandResponse
	c.Command = cmd
	memsfcr.FCRSendToECU <- c
	utils.LogI.Printf("%s FCR sent command '%x' to ECU", utils.ECUCommandTrace, cmd)
	/*
		select {
		case memsfcr.FCRSendToECU <- c:
			utils.LogI.Printf("%s FCR sent command '%x' to ECU", utils.ECUCommandTrace, cmd)
		default:
			utils.LogW.Printf("%s FCR unable to send command '%x' to ECU on FCRSendToECU, blocked?", utils.ECUCommandTrace, cmd)
		}*/
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
		tx := <-memsfcr.FCRSendToECU
		utils.LogI.Printf("%s (Tx.1) FCR received command from FCR FCRSendToECU to send to ECU", utils.ECUCommandTrace)

		// block waiting for the command to be sent to the ECU
		memsfcr.ECU.TxECU <- tx
		utils.LogI.Printf("%s (Tx.2) FCR sent command to ECU TxECU channel", utils.ECUCommandTrace)

		// block waiting for the response to be received from the ECU
		rx := <-memsfcr.ECU.RxECU
		utils.LogI.Printf("%s (Rx.1) FCR received response from ECU RxECU channel", utils.ECUResponseTrace)

		// block waiting for the response to be collected for processing
		memsfcr.ECUSendToFCR <- rx
		utils.LogI.Printf("%s (Rx.2) ECU response sent to FCR ECUSendToFCR for processing", utils.ECUResponseTrace)
	}
}
