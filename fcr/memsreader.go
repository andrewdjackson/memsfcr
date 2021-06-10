package fcr

import (
	"fmt"
	"runtime"

	"github.com/andrewdjackson/rosco"
	"github.com/pkg/browser"
	log "github.com/sirupsen/logrus"
)

// MemsReader structure
type MemsReader struct {
	// Config FCR configuration
	Config *Config
	// ECU represents the serial connection to the ECU
	ECU *rosco.MemsConnection
	// Webserver
	WebServer *WebServer
}

func NewMemsReader(version string, build string) *MemsReader {
	reader := &MemsReader{}

	// read the config
	reader.Config = ReadConfig()
	reader.Config.Version = version
	reader.Config.Build = build

	// set up the connection to the ECU
	// this is also used to 'emulate' the ECU if
	// a pre-recorded scenario is played back
	reader.ECU = rosco.NewMemsConnection(reader.Config.LogFolder)

	// set up the webserver for websocket
	// and REST endpoints
	reader.WebServer = NewWebServer(reader)

	return reader
}

func (reader *MemsReader) StartWebServer() {
	// run the web server as a concurrent process
	go reader.WebServer.RunHTTPServer()

	// display the web interface, wait for the HTTP Server to start
	for {
		if reader.WebServer.ServerRunning {
			break
		}
	}
}

// OpenBrowser opens the browser
func (reader *MemsReader) OpenBrowser() {
	url := fmt.Sprintf("http://127.0.0.1:%d/index.html", reader.WebServer.HTTPPort)

	var err error

	log.Infof("opening browser (%s)", runtime.GOOS)
	err = browser.OpenURL(url)

	if err != nil {
		log.Errorf("error opening browser (%s)", err)
	}
}
