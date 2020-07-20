package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"runtime"

	"github.com/andrewdjackson/memsfcr/rosco"
	"github.com/andrewdjackson/memsfcr/ui"
	"github.com/andrewdjackson/memsfcr/utils"
	"github.com/pkg/browser"
	"github.com/zserge/webview"
)

var (
	// Version of the application
	Version string
	// Build date
	Build string
)

// MemsReader structure
type MemsReader struct {
	wi         *ui.WebInterface
	fcr        *ui.MemsFCR
	dataLogger *rosco.MemsDataLogger
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

// webLoop services the channels processing messages send from the web interface
// run as a goroutine
func (r *MemsReader) webMainLoop() {
	for {
		m := <-r.wi.FromWebChannel
		utils.LogI.Printf("%s received message FromWebChannel in main webLoop (%v)", utils.ReceiveFromWebTrace, m)

		// evalute the message sent from the web interface
		// and determine the action

		action := ui.EvaluateWebMsg(m)
		utils.LogI.Printf("evaluated action (%+v) as %d", action.Msg, action.Value)

		switch action.Value {

		case ui.ConfigRead:
			r.sendConfigToWebView()

		case ui.Save:
			cfg := rosco.ReadmemsConfig{}
			json.Unmarshal([]byte(m.Data), &cfg)

			utils.LogI.Printf("applying config (%+v)", cfg)

			r.fcr.Config.Port = cfg.Port
			r.fcr.Config.LogFolder = cfg.LogFolder
			r.fcr.Config.LogToFile = cfg.LogToFile
			r.fcr.Config.Frequency = cfg.Frequency

			if r.fcr.Config.LogToFile == "true" {
				r.fcr.Logging = false
			}

			rosco.WriteConfig(r.fcr.Config)

		case ui.ConnectECU:
			// connect the ECU
			utils.LogI.Printf("connecting ecu")
			r.fcr.ConnectFCR()
			utils.LogI.Printf("sending connection status")
			r.sendConnectionStatusToWebView()

		case ui.Dataframe:
			go r.fcr.TxECU(rosco.MEMSDataFrame)

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
			go r.fcr.TxECU(rosco.MEMSResetECU)

		case ui.ClearFaults:
			go r.fcr.TxECU(rosco.MEMSClearFaults)

		case ui.ResetAdjustments:
			go r.fcr.TxECU(rosco.MEMSResetAdj)

		case ui.IncreaseIdleSpeed:
			go r.fcr.TxECU(rosco.MEMSIdleSpeedIncrement)

		case ui.IncreaseIdleHot:
			go r.fcr.TxECU(rosco.MEMSIdleDecayIncrement)

		case ui.IncreaseFuelTrim:
			go r.fcr.TxECU(rosco.MEMSLTFTIncrement)

		case ui.IncreaseIgnitionAdvance:
			go r.fcr.TxECU(rosco.MEMSIgnitionAdvanceOffsetIncrement)

		case ui.DecreaseIdleSpeed:
			go r.fcr.TxECU(rosco.MEMSIdleSpeedDecrement)

		case ui.DecreaseIdleHot:
			go r.fcr.TxECU(rosco.MEMSIdleDecayDecrement)

		case ui.DecreaseFuelTrim:
			go r.fcr.TxECU(rosco.MEMSLTFTDecrement)

		case ui.DecreaseIgnitionAdvance:
			go r.fcr.TxECU(rosco.MEMSIgnitionAdvanceOffsetDecrement)

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
	c.ECUID = fmt.Sprintf("%X", r.fcr.ECU.ECUID)
	c.IACPosition = r.fcr.ECU.Diagnostics.Analysis.IACPosition

	m.Action = ui.WebActionConnection

	data, _ := json.Marshal(c)
	m.Data = string(data)

	r.wi.ToWebChannel <- m
	utils.LogI.Printf("%s sent connection status to web with ToWebChannel channel", utils.SendToWebTrace)
}

func (r *MemsReader) fcrMainLoop() {
	var data []byte

	loggerOpen := false

	// busy clearing channels
	for {
		m := <-r.fcr.ECUSendToFCR
		utils.LogI.Printf("%s (Rx.3) received message ECUSendToFCR (%v)", utils.ReceiveFromWebTrace, m)

		// send to the web
		df := ui.WebMsg{}

		if bytes.Compare(m.Command, rosco.MEMSDataFrame) == 0 {
			// dataframe command
			df.Action = ui.WebActionData
			data, _ = json.Marshal(m.MemsDataFrame)
			if r.fcr.Logging {
				if r.fcr.ECU.Connected && !loggerOpen {
					prefix := fmt.Sprintf("%X-", r.fcr.ECU.ECUID)
					
					// create the data logger
					utils.LogI.Printf("opening log file with prefix %s", prefix)
					r.dataLogger = rosco.NewMemsDataLogger(r.fcr.Config.LogFolder, prefix)
					loggerOpen = true
				}

				// write data to log file
				r.dataLogger.WriteMemsDataToFile(m.MemsDataFrame)
			}
		} else {
			// send the response from the ECU to the web interface
			df.Action = ui.WebActionECUResponse
			ecuResponse := hex.EncodeToString(m.Response)
			data, _ = json.Marshal(ecuResponse)
		}

		df.Data = string(data)

		select {
		case r.wi.ToWebChannel <- df:
		default:
		}

		// send the diagnostics to the web interface
		r.fcrSendDiagnosticsToWebView()
	}
}

func (r *MemsReader) fcrSendDiagnosticsToWebView() {
	m := ui.WebMsg{}
	m.Action = ui.WebActionDiagnostics
	data, _ := json.Marshal(r.fcr.ECU.Diagnostics.Analysis)
	m.Data = string(data)

	utils.LogI.Printf("%s sending diagnostics to web (%v)", utils.SendToWebTrace, m)
	r.wi.ToWebChannel <- m
}

func openBrowser(url string) {
	var err error

	utils.LogI.Printf("opening browser (%s)", runtime.GOOS)
	err = browser.OpenURL(url)

	if err != nil {
		utils.LogE.Printf("error opening browser (%s)", err)
	}

	for {
	}
}

// displayWebView creates a webview
// this must be run in the main thread
func displayWebView(wi *ui.WebInterface, localView bool) {
	url := fmt.Sprintf("http://127.0.0.1:%d/index.html", wi.HTTPPort)

	if localView {
		w := webview.New(true)
		defer w.Destroy()

		w.SetTitle("MEMS Fault Code Reader")
		w.SetSize(1280, 1024, webview.HintNone)

		w.Bind("quit", func() {
			w.Terminate()
		})

		w.Navigate(url)
		w.Run()
	} else {
		openBrowser(url)
	}
}

func main() {
	utils.LogI.Printf("\nMemsFCR\nVersion %s (Build %s)\n\n", Version, Build)

	var debug bool
	flag.BoolVar(&debug, "debug", false, "enable debug")
	flag.Parse()

	memsReader := NewMemsReader()

	go memsReader.wi.RunHTTPServer()
	go memsReader.webMainLoop()
	go memsReader.fcrMainLoop()
	go memsReader.fcr.TxRxECULoop()

	// run the listener for messages sent to the web interface from
	// the backend application
	go memsReader.wi.ListenToWebChannelLoop()

	// display the web interface, wait for the HTTP Server to start
	for {
		if memsReader.wi.ServerRunning {
			break
		}
	}

	utils.LogI.Printf("starting webview.. (%v)", memsReader.wi.HTTPPort)

	// show the app in a local go webview window rather than in the web browser
	// unless debug is enabled
	showLocal := !debug

	// use default browser on Windows until I can get the Webview to work
	if runtime.GOOS == "windows" {
		showLocal = false
	}

	// use the browser if the user has configured this option
	if memsReader.fcr.Config.UseBrowser == "true" {
		showLocal = false
	}

	// if debug enabled use the full browser
	if memsReader.fcr.Config.Debug == "true" {
		showLocal = false
	}

	displayWebView(memsReader.wi, showLocal)
}
