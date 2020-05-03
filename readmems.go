package main

import (
	"encoding/json"
	"fmt"

	"github.com/andrewdjackson/readmems/rosco"
	"github.com/andrewdjackson/readmems/ui"
	"github.com/andrewdjackson/readmems/utils"
	"github.com/zserge/webview"
)

/*
const version = "v0.1.0"

var header = fmt.Sprintf("\nMemsFCR %s\n", version)
var config *rosco.ReadmemsConfig
var paused = false
var logging = false
var mems *rosco.MemsConnection

func helpMessage() string {
	return fmt.Sprintf(`%s
	ROSCO MEMS 1.6 Fault Code Reader

    Usage:
	memsfcr [flags]

    Flags:
	-port		Name/path of the serial port
	-command	Command to execute on the ECU {read}
	-loop		Command execution loop count, use 'inf' for infinite
	-output		Use 'stdout' to send response to console, 'file' to log in CSV format to a file {stdout | file}
	-wait		Retry the connection until a connection is established {true | false}
	-help		This help message
	`, header)
}

func connectToECU(config *rosco.ReadmemsConfig) {
	const connectionRetryInterval = 2000 * time.Millisecond
	maxRetries := 0

	if config.Connection == "wait" {
		maxRetries = 1
	}

	for count := 0; count < maxRetries; count++ {
		// attempt to connect and initialise the ECU
		mems.ConnectAndInitialiseECU(config.Port)

		if mems.SerialPort == nil {
			// serial port is not available
			utils.LogE.Printf("serial port not connected, retrying (%d of %d)", count, maxRetries)
			time.Sleep(connectionRetryInterval)
		}
	}

	// update the web view
	sendConnectionStatusToWebView()
}

func memsCommandResponseLoop(config *rosco.ReadmemsConfig) {
	const DataInterval = 500 * time.Millisecond
	const HeartbeatInterval = 2000 * time.Millisecond
	var logger *rosco.MemsDataLogger

	// attempt to connect to the ECU
	connectToECU(config)

	// if the connection has been established and the ECU completed initialisation
	// then start the loop to send commands and recieve responses from the ECU
	if mems.Initialised == true {

		// enable logging if configured
		if config.Output == "file" {
			logging = true
			logger = rosco.NewMemsDataLogger()
		}

		for {
			// get how many dataframe calls to make
			maxDataFrameCalls, _ := strconv.Atoi(config.Loop)

			// start a loop for listening for events from the web interface
			go recieveMessageFromWebViewLoop()
			// start a loop to listen for data responses from the ECU
			go mems.ListenSendToECUChannelLoop()

			// enter a command / response loop
			for loop := 0; loop < maxDataFrameCalls; {
				if paused {
					// send a periodic heartbeat to keep the connection alive when paused
					utils.LogI.Printf("%s memsCommandResponseLoop sending heatbeat", utils.ECUCommandTrace)
					go sendCommandToMemsChannel(rosco.MEMS_Heartbeat)

					// send heatbeats at a slower interval to data frame requests
					time.Sleep(HeartbeatInterval)

				} else {
					// read data from the ECU
					utils.LogI.Printf("%s memsCommandResponseLoop sending dataframe request to ECU", utils.ECUCommandTrace)
					mems.ReadMemsData()

					// wait for response, this is built into the sendCommand function
					// but as we're reading the MemData we need to call this here
					data := receiveResponseFromMemsChannel()

					//utils.LogI.Printf("waiting for response from ECU")
					//data := <-mems.ReceivedFromECU
					//utils.LogI.Printf("received dataframe from ECU")

					// send it to the web interface
					sendDataToWebView(data.MemsDataFrame)

					if logging {
						// write to a log file if logging is enabled
						go logger.WriteMemsDataToFile(data.MemsDataFrame)
					}

					// increment count of data calls
					// don't increment if we're paused
					loop = loop + 1

					// sleep between calls to give the ECU time to catch up
					// the ECU will get slower as load increases so this ensures
					// a regular time series for the data set
					time.Sleep(DataInterval)
				}
			}

			// read loop complete, exit
			break
		}
	}
}

// loop listening for messages from the web interface
// send commands to the ecu if required via a go routine as to not block
func recieveMessageFromWebViewLoop() {
	for {
		utils.LogI.Printf("%s waiting for message from webview webToMemsChannel channel..", utils.ReceiveFromWebTrace)
		m := <-webToMemsChannel
		utils.LogI.Printf("%s recieved message from webToMemsChannel channel", utils.ReceiveFromWebTrace)

		c := evaluateCommand(m)

		switch c {
		case commandPauseDataLoop:
			{
				paused = true
				utils.LogI.Printf("Paused Data Loop, sending heartbeats to keep connection alive")
			}
		case commandStartDataLoop:
			{
				paused = false
				utils.LogI.Printf("Resuming Data Loop")
			}
		case commandResetECU:
			{
				utils.LogI.Printf("memsCommandResponseLoop sending Reset ECU")
				go sendCommandToMemsChannel(rosco.MEMS_ResetECU)
			}
		case commandClearFaults:
			{
				utils.LogI.Printf("memsCommandResponseLoop sending Clear Faults")
				go sendCommandToMemsChannel(rosco.MEMS_ClearFaults)
			}
		case commandResetAdjustments:
			{
				utils.LogI.Printf("memsCommandResponseLoop sending Reset Adjustments")
				go sendCommandToMemsChannel(rosco.MEMS_ResetAdj)
			}
		case commandIncreaseIdleSpeed:
			{
				utils.LogI.Printf("memsCommandResponseLoop sending Increase Idle Speed")
				go sendCommandToMemsChannel(rosco.MEMS_IdleSpeed_Increment)
			}
		case commandIncreaseIdleHot:
			{
				utils.LogI.Printf("memsCommandResponseLoop sending Increase Idle Decay (Hot)")
				go sendCommandToMemsChannel(rosco.MEMS_IdleDecay_Increment)
			}
		case commandIncreaseFuelTrim:
			{
				utils.LogI.Printf("memsCommandResponseLoop sending Increase Fuel Trim (LTFT)")
				go sendCommandToMemsChannel(rosco.MEMS_LTFT_Increment)
			}
		case commandIncreaseIgnitionAdvance:
			{
				utils.LogI.Printf("memsCommandResponseLoop sending Increase Ignition Advance")
				go sendCommandToMemsChannel(rosco.MEMS_IgnitionAdvanceOffset_Increment)
			}
		case commandDecreaseIdleSpeed:
			{
				utils.LogI.Printf("memsCommandResponseLoop sending Decrease Idle Speed")
				go sendCommandToMemsChannel(rosco.MEMS_IdleSpeed_Decrement)
			}
		case commandDecreaseIdleHot:
			{
				utils.LogI.Printf("memsCommandResponseLoop sending Decrease Idle Decay (Hot)")
				go sendCommandToMemsChannel(rosco.MEMS_IdleDecay_Decrement)
			}
		case commandDecreaseFuelTrim:
			{
				utils.LogI.Printf("memsCommandResponseLoop sending Decrease Fuel Trim (LTFT)")
				go sendCommandToMemsChannel(rosco.MEMS_LTFT_Decrement)
			}
		case commandDecreaseIgnitionAdvance:
			{
				utils.LogI.Printf("memsCommandResponseLoop sending Decrease Ignition Advance")
				go sendCommandToMemsChannel(rosco.MEMS_IgnitionAdvanceOffset_Decrement)
			}
		default:
		}
	}
}

// send the command to be executed by the ECU via a channel
func sendCommandToMemsChannel(command []byte) {
	var m rosco.MemsCommandResponse
	m.Command = command

	// send through channel
	utils.LogI.Printf("%s waiting to send mems command to the SendToECU channel", utils.ECUCommandTrace)
	select {
	case mems.SendToECU <- m:
		utils.LogI.Printf("%s mems command sent to the SendToECU channel", utils.ECUCommandTrace)
	default:
		utils.LogE.Printf("%s unable to send mems command to the SendToECU channel", utils.ECUCommandTrace)
	}

	// wait for response
	receiveResponseFromMemsChannel()
}

func receiveResponseFromMemsChannel() rosco.MemsCommandResponse {
	// wait for response
	utils.LogI.Printf("%s waiting for response from ECU", utils.ECUResponseTrace)

	data := <-mems.ReceivedFromECU

	utils.LogI.Printf("%s received dataframe from ECU", utils.ECUResponseTrace)

	return data
}

// send a connection status message back to the web interface via a channel
func sendConnectionStatusToWebView() {
	var c rosco.MemsConnectionStatus
	var m wsMsg

	c.Connected = mems.Connected
	c.Initialised = mems.Initialised

	m.Action = WebActionConnection

	data, _ := json.Marshal(c)
	m.Data = string(data)

	utils.LogI.Printf("%s waiting to send connection status to webview with memsToWebChannel channel", utils.SendToWebTrace)
	memsToWebChannel <- m
	utils.LogI.Printf("%s sent connection status to webview with memsToWebChannel channel", utils.SendToWebTrace)
}

// send a message back to the web interface via a channel
func sendDataToWebView(memsdata rosco.MemsData) {
	var m wsMsg

	m.Action = WebActionData

	data, _ := json.Marshal(memsdata)
	m.Data = string(data)

	utils.LogI.Printf("%s waiting to send data to webview with memsToWebChannel channel", utils.SendToWebTrace)
	memsToWebChannel <- m
	utils.LogI.Printf("%s sent data to webview with memsToWebChannel channel", utils.SendToWebTrace)
}

// send configuration to the web interace via a channel
func sendConfigToWebView(config *rosco.ReadmemsConfig) {
	var m wsMsg

	m.Action = WebActionConfig

	data, _ := json.Marshal(config)
	m.Data = string(data)
	utils.LogI.Printf("%s waiting to send config to webview with memsToWebChannel channel", utils.SendToWebTrace)
	memsToWebChannel <- m
	utils.LogI.Printf("%s sent config to webview with memsToWebChannel channel", utils.SendToWebTrace)
}

func getSerialPorts() []string {
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

func oldmain() {
	var showHelp bool

	// create a map of commands
	createCommandMap()

	// use if the readmems config is supplied
	config = rosco.ReadConfig()

	// parse the command line parameters and override config file settings
	flag.StringVar(&config.Port, "port", config.Port, "Name/path of the serial port")
	flag.StringVar(&config.Command, "command", config.Command, "Command to send")
	flag.StringVar(&config.Loop, "loop", config.Loop, "Read loop count, 'inf' for infinite")
	flag.BoolVar(&showHelp, "help", false, "A brief help message")
	flag.Parse()

	if showHelp {
		utils.LogI.Println(helpMessage())
		return
	}

	if config.Loop == "inf" {
		// infitite loop, so set loop count to a very big number
		config.Loop = "100000000"
	}

	// get the list of ports available
	config.Ports = append(config.Ports, config.Port)
	config.Ports = append(config.Ports, getSerialPorts()...)

	// create a mems instance
	mems = rosco.NewMemsConnection()

	// create and run the web interfacce
	wi := ui.NewWebInterface()
	utils.LogI.Printf("running web server %d", wi.HTTPPort)

	go wi.RunHTTPServer()

	displayWebView(wi)

	// run the http server
	//go RunHTTPServer()
	//go sendConfigToWebView(config)

	//ShowWebView(config)
}
*/

///////////////////////////////////////////

// MemsReader structure
type MemsReader struct {
	wi  *ui.WebInterface
	fcr *ui.MemsFCR
}

// NewMemsReader creates an instance of a MEMs Reader
func NewMemsReader() *MemsReader {
	r := &MemsReader{}

	// create the Mems Fault Code Reader
	r.fcr = ui.NewMemsFCR()

	// create a mems instance and assign it to the fault code reader instance
	r.fcr.ECU = rosco.NewMemsConnection()

	// create and run the web interfacce
	r.wi = ui.NewWebInterface()
	utils.LogI.Printf("running web server %d", r.wi.HTTPPort)

	return r
}

func (r *MemsReader) webLoop() {
	// busy clearing channels
	for {
		m := <-r.wi.FromWebChannel
		utils.LogI.Printf("%s received message FromWebChannel in main webLoop (%v)", utils.ReceiveFromWebTrace, m)

		switch m.Action {
		case ui.WebActionConfig:
			// configuration settings requested
			if m.Data == "read" {
				r.sendConfigToWebView()
			}
		case ui.WebActionConnect:
			// connect the ECU
			utils.LogI.Printf("connecting ecu")
			r.fcr.ConnectFCR()
		default:
		}
	}
}

func (r *MemsReader) sendConfigToWebView() {
	// pass configuration to the web interface
	m := ui.WebMsg{}
	m.Action = ui.WebActionConfig
	data, _ := json.Marshal(r.fcr.Config)
	m.Data = string(data)
	r.wi.ToWebChannel <- m
}

func (r *MemsReader) fcrLoop() {
	// busy clearing channels
	for {
		m := <-r.fcr.FromECUChannel
		utils.LogI.Printf("%s received message FromECUChannel (%v)", utils.ReceiveFromWebTrace, m)
	}
}

// displayWebView creates a webview
// this must be run in the main thread
func displayWebView(wi *ui.WebInterface) {
	w := webview.New(true)
	defer w.Destroy()

	w.SetTitle("MEMS Fault Code Reader")
	w.SetSize(1120, 920, webview.HintNone)

	w.Bind("quit", func() {
		w.Terminate()
	})

	url := fmt.Sprintf("http://127.0.0.1:%d/public/html/index.html", wi.HTTPPort)

	w.Navigate(url)
	w.Run()
}

func main() {
	memsReader := NewMemsReader()

	go memsReader.wi.RunHTTPServer()
	go memsReader.webLoop()
	go memsReader.fcrLoop()

	// run the listener for messages sent to the web interface from
	// the backend application
	go memsReader.wi.ListenToWebChannelLoop()

	// display the web interface
	displayWebView(memsReader.wi)
}
