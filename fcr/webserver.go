package fcr

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// WebMsg structure fro sending / receiving over the websocket
type WebMsg struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

// WebServer the web interface
type WebServer struct {
	// multiplex router interface
	router *mux.Router
	// websocket interface
	httpDir  string
	ws       *websocket.Conn
	upgrader websocket.Upgrader
	// HTTPPort used by the HTTP Server instance
	HTTPPort int
	// channels for communication over the websocket
	ToWebSocketChannel   chan WebMsg
	FromWebSocketChannel chan WebMsg
	// ServerRunning indicates where the server is active
	ServerRunning bool
	// Pointer to Mems Fault Code Reader
	reader *MemsReader
}

// NewWebInterface creates a new web interface
func NewWebServer(reader *MemsReader) *WebServer {
	webserver := &WebServer{}
	webserver.ToWebSocketChannel = make(chan WebMsg)
	webserver.FromWebSocketChannel = make(chan WebMsg)
	webserver.HTTPPort = 0
	webserver.httpDir = ""
	webserver.ServerRunning = false
	webserver.reader = reader

	webserver.upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return webserver
}

func (webserver *WebServer) newRouter() *mux.Router {
	var webroot string

	// determine the path to find the local html files
	// based on the current executable path
	dir, _ := os.Getwd()
	exepath := filepath.FromSlash(dir)
	path, err := filepath.Abs(exepath)

	log.Infof("path to the local html files (%s) on (%s)", path, runtime.GOOS)

	// use default browser on Windows until I can get the Webview to work
	if runtime.GOOS == "darwin" {
		// MacOS use .app Resources
		if strings.Contains(path, "MacOS") {
			// packaged app
			webroot = strings.Replace(path, "MacOS", "Resources", -1)
		} else {
			// running a local or dev version
			webroot = fmt.Sprintf("%s/Resources", path)
		}
	} else if runtime.GOOS == "linux" {
		// linux path
		webroot = fmt.Sprintf("%s/resources", path)
	} else {
		// windows use the exe subdirectory
		webroot = fmt.Sprintf("%s\\resources", path)
	}

	webserver.httpDir = filepath.ToSlash(webroot)

	log.Infof("path to the local html files (%s)", webserver.httpDir)

	if err != nil {
		log.Errorf("unable to find the current path to the local html files (%s)", err)
	}

	// set a router and a handler to accept messages over the websocket

	r := mux.NewRouter()
	//r.HandleFunc("/ws", webserver.websocketHandler)
	r.HandleFunc("/heartbeat", webserver.browserHeartbeatHandler)

	r.HandleFunc("/config", webserver.getConfigHandler).Methods("GET")
	r.HandleFunc("/config/ports", webserver.getSerialPortsHandler).Methods("GET")
	r.HandleFunc("/config", webserver.updateConfigHandler).Methods("POST")

	r.HandleFunc("/scenario", webserver.getListofScenarios).Methods("GET")
	r.HandleFunc("/scenario/{scenarioId}", webserver.getScenarioDetails).Methods("GET")

	r.HandleFunc("/rosco", webserver.getECUConnectionStatus).Methods("GET")
	r.HandleFunc("/rosco/connect", webserver.postECUConnect).Methods("POST")
	r.HandleFunc("/rosco/disconnect", webserver.postECUDisconnect).Methods("POST")
	r.HandleFunc("/rosco/dataframe", webserver.getECUDataframes).Methods("GET")
	r.HandleFunc("/rosco/heartbeat", webserver.postECUHeartbeat).Methods("POST")
	r.HandleFunc("/rosco/iac", webserver.getECUIAC).Methods("GET")

	r.HandleFunc("/rosco/reset", webserver.postECUReset).Methods("POST")
	r.HandleFunc("/rosco/reset/ecu", webserver.postECUReset).Methods("POST")
	r.HandleFunc("/rosco/reset/faults", webserver.postECUClearFaults).Methods("POST")
	r.HandleFunc("/rosco/reset/adjustments", webserver.postECUClearAdjustments).Methods("POST")

	r.HandleFunc("/rosco/adjust/stft", webserver.postECUAdjustSTFT).Methods("POST")
	r.HandleFunc("/rosco/adjust/ltft", webserver.postECUAdjustLTFT).Methods("POST")
	r.HandleFunc("/rosco/adjust/idledecay", webserver.postECUAdjustIdleDecay).Methods("POST")
	r.HandleFunc("/rosco/adjust/idlespeed", webserver.postECUAdjustIdleSpeed).Methods("POST")
	r.HandleFunc("/rosco/adjust/ignitionadvance", webserver.postECUAdjustIgnitionAdvance).Methods("POST")
	r.HandleFunc("/rosco/adjust/iac", webserver.postECUAdjustIAC).Methods("POST")

	r.HandleFunc("/rosco/test/fuelpump", webserver.postECUTestFuelPump).Methods("POST")
	r.HandleFunc("/rosco/test/ptc", webserver.postECUTestPTC).Methods("POST")
	r.HandleFunc("/rosco/test/aircon", webserver.postECUTestAircon).Methods("POST")
	r.HandleFunc("/rosco/test/purgevalve", webserver.postECUTestPurgeValve).Methods("POST")
	r.HandleFunc("/rosco/test/boostvalve", webserver.postECUTestBoostValve).Methods("POST")
	r.HandleFunc("/rosco/test/fan", webserver.postECUTestFan1).Methods("POST")
	r.HandleFunc("/rosco/test/fan/1", webserver.postECUTestFan1).Methods("POST")
	r.HandleFunc("/rosco/test/fan/2", webserver.postECUTestFan2).Methods("POST")
	r.HandleFunc("/rosco/test/injectors", webserver.postECUTestInjectors).Methods("POST")
	r.HandleFunc("/rosco/test/coil", webserver.postECUTestCoil).Methods("POST")

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir(webserver.httpDir))
	r.Handle("/", fileServer)
	r.PathPrefix("/").Handler(fileServer).Methods("GET")

	return r
}

// RunHTTPServer run the server
func (webserver *WebServer) RunHTTPServer() {
	// Declare a new router
	webserver.router = webserver.newRouter()

	// We can then pass our router (after declaring all our routes) to this method
	// (where previously, we were leaving the secodn argument as nil)
	listener, err := net.Listen("tcp", ":8081")

	if err != nil {
		log.Errorf("error starting web interface (%s)", err)
	}

	webserver.HTTPPort = listener.Addr().(*net.TCPAddr).Port

	log.Infof("started http server on port %d", webserver.HTTPPort)
	webserver.ServerRunning = true

	err = http.Serve(listener, webserver.router)

	if err != nil {
		log.Errorf("error starting web interface (%s)", err)
	}
}

// send message to the web interface over the websocket
func (webserver *WebServer) sendMessageToWebInterface(m WebMsg) {
	if webserver.ws != nil {
		err := webserver.ws.WriteJSON(m)
		if err != nil {
			log.Errorf("error sending message over websocket (%s)", err)
		} else {
			log.Infof("send message over websocket")
		}
	} else {
		log.Warnf("unable to send message over websocket, connected?")
	}
}

// ListenToWebSocketChannelLoop loop for listening for messages over the ToWebSocketChannel
// these are messages that are to be passed to the web interface over the websocket
// from the backend application
// to be run as a go routine as the channel is coded to be non blocking
func (webserver *WebServer) ListenToWebSocketChannelLoop() {
	for {
		m := <-webserver.ToWebSocketChannel
		webserver.sendMessageToWebInterface(m)
		log.Infof("sent message '%s : %s' on ToWebSocketChannel", m.Action, m.Data)
	}
}
