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

		// evalute the message sent from the web interface
		// and determine the action

		action := ui.EvaluateWebMsg(m)
		utils.LogI.Printf("evaluated action (%v) as %d", action.Msg, action.Value)

		switch action.Value {
		case ui.ConfigRead:
			r.sendConfigToWebView()

		case ui.ConnectECU:
			// connect the ECU
			utils.LogI.Printf("connecting ecu")
			if r.fcr.ConnectFCR() {
				r.sendConnectionStatusToWebView()
			}

		case ui.PauseDataLoop:
			{
				//paused = true
				utils.LogI.Printf("Paused Data Loop, sending heartbeats to keep connection alive")
			}
		case ui.StartDataLoop:
			{
				//paused = false
				utils.LogI.Printf("Resuming Data Loop")
			}
		case ui.ResetECU:
			go r.fcr.SendToECU(rosco.MEMS_ResetECU)
		case ui.ClearFaults:
			go r.fcr.SendToECU(rosco.MEMS_ClearFaults)
		case ui.ResetAdjustments:
			go r.fcr.SendToECU(rosco.MEMS_ResetAdj)
		case ui.IncreaseIdleSpeed:
			go r.fcr.SendToECU(rosco.MEMS_IdleSpeed_Increment)
		case ui.IncreaseIdleHot:
			go r.fcr.SendToECU(rosco.MEMS_IdleDecay_Increment)
		case ui.IncreaseFuelTrim:
			go r.fcr.SendToECU(rosco.MEMS_LTFT_Increment)
		case ui.IncreaseIgnitionAdvance:
			go r.fcr.SendToECU(rosco.MEMS_IgnitionAdvanceOffset_Increment)
		case ui.DecreaseIdleSpeed:
			go r.fcr.SendToECU(rosco.MEMS_IdleSpeed_Decrement)
		case ui.DecreaseIdleHot:
			go r.fcr.SendToECU(rosco.MEMS_IdleDecay_Decrement)
		case ui.DecreaseFuelTrim:
			go r.fcr.SendToECU(rosco.MEMS_LTFT_Decrement)
		case ui.DecreaseIgnitionAdvance:
			go r.fcr.SendToECU(rosco.MEMS_IgnitionAdvanceOffset_Decrement)

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

// send a connection status message back to the web interface via a channel
func (r *MemsReader) sendConnectionStatusToWebView() {
	var c rosco.MemsConnectionStatus
	var m ui.WebMsg

	c.Connected = r.fcr.ECU.Connected
	c.Initialised = r.fcr.ECU.Initialised

	m.Action = ui.WebActionConnection

	data, _ := json.Marshal(c)
	m.Data = string(data)

	r.wi.ToWebChannel <- m
	utils.LogI.Printf("%s sent connection status to webview with memsToWebChannel channel", utils.SendToWebTrace)
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
	//go memsReader.fcrLoop()
	go memsReader.fcr.TxRxECULoop()

	// run the listener for messages sent to the web interface from
	// the backend application
	go memsReader.wi.ListenToWebChannelLoop()

	// display the web interface
	displayWebView(memsReader.wi)
}
