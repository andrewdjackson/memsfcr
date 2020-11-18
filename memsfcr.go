package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/andrewdjackson/memsfcr/rosco"
	"github.com/andrewdjackson/memsfcr/ui"
	"github.com/andrewdjackson/memsfcr/utils"
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

	// create a dynamic version label and build number when in development
	// these values get automatically set by the make file on production build
	if strings.Compare(Version, "") == 0 {
		Version = "0.0.0"
		currentTime := time.Now()
		Build = currentTime.Format("2006-01-02")
	}
	// create the Mems Fault Code Reader
	r.fcr = ui.NewMemsFCR(Version, Build)

	// create a mems instance and assign it to the fault code reader instance
	r.fcr.ECU = rosco.NewMemsConnection()

	// create and run the web interfacce
	r.wi = ui.NewWebInterface(r.fcr)
	utils.LogI.Printf("running web server %d", r.wi.HTTPPort)

	return r
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

		if bytes.Equal(m.Command, rosco.MEMSDataFrame) {
			// dataframe command
			df.Action = ui.WebActionData
			data, _ = json.Marshal(m.MemsDataFrame)
			if r.fcr.Logging {
				if r.fcr.ECU.Connected && !loggerOpen {
					prefix := fmt.Sprintf("%X-", r.fcr.ECU.ECUID)

					if r.fcr.ECU.Emulated {
						loggerOpen = false
					} else {
						// create the data logger
						utils.LogI.Printf("opening log file with prefix %s", prefix)
						r.dataLogger = rosco.NewMemsDataLogger(r.fcr.Config.LogFolder, prefix)
						loggerOpen = true
					}
				}

				// write data to log file
				if loggerOpen {
					r.dataLogger.WriteMemsDataToFile(m.MemsDataFrame)
				}
			}
		} else {
			// send the response from the ECU to the web interface
			df.Action = ui.WebActionECUResponse
			ecuResponse := hex.EncodeToString(m.Response)
			data, _ = json.Marshal(ecuResponse)
		}

		df.Data = string(data)

		select {
		case r.wi.ToWebSocketChannel <- df:
		default:
		}

		// send the diagnostics to the web interface
		r.fcrSendDiagnosticsToWebView()
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
	go memsReader.wi.ListenToWebSocketChannelLoop()

	// display the web interface, wait for the HTTP Server to start
	for {
		if memsReader.wi.ServerRunning {
			break
		}
	}

	utils.LogI.Printf("starting webview.. (%v)", memsReader.wi.HTTPPort)

	headless := false

	if memsReader.fcr.Config.Headless == "true" {
		headless = true
	}

	displayWebView(memsReader.wi, headless)
}
