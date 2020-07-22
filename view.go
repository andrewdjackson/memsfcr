package main

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/andrewdjackson/memsfcr/rosco"
	"github.com/andrewdjackson/memsfcr/scenarios"
	"github.com/andrewdjackson/memsfcr/ui"
	"github.com/andrewdjackson/memsfcr/utils"
	"github.com/pkg/browser"
	"github.com/zserge/webview"
)

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
			// save the configuration
			cfg := utils.ReadmemsConfig{}
			json.Unmarshal([]byte(m.Data), &cfg)

			utils.LogI.Printf("applying config (%+v)", cfg)

			r.fcr.Config.Port = cfg.Port
			r.fcr.Config.LogFolder = cfg.LogFolder
			r.fcr.Config.LogToFile = cfg.LogToFile
			r.fcr.Config.Frequency = cfg.Frequency

			if r.fcr.Config.LogToFile == "true" {
				r.fcr.Logging = false
			}

			utils.WriteConfig(r.fcr.Config)

		case ui.ConnectECU:
			// connect the ECU
			utils.LogI.Printf("connecting ecu")
			r.fcr.ConnectFCR()
			utils.LogI.Printf("sending connection status")
			r.sendConnectionStatusToWebView()

		case ui.Replay:
			utils.LogI.Printf("%s replay requested %+v", utils.EmulatorTrace, m.Data)
			r.fcr.Config.Port = scenarios.GetScenarioPath(m.Data)

		case ui.Dataframe:
			// request a dataframe from the ECU
			go r.fcr.TxECU(rosco.MEMSDataFrame)

		case ui.PauseDataLoop:
			// do nothing as this is handled in the web interface
			utils.LogI.Printf("Paused Data Loop")

		case ui.StartDataLoop:
			// do nothing as this is handled in the web interface
			utils.LogI.Printf("Resuming Data Loop")

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
