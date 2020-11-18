package ui

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/andrewdjackson/memsfcr/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// WebMsg structure fro sending / receiving over the websocket
type WebMsg struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

// WebInterface the web interface
type WebInterface struct {
	// mulitplex router interface
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
	// ServerRunning inidicates where the server is active
	ServerRunning bool
	// Pointer to FCR
	fcr *MemsFCR
}

// NewWebInterface creates a new web interface
func NewWebInterface(fcr *MemsFCR) *WebInterface {
	wi := &WebInterface{}
	wi.ToWebSocketChannel = make(chan WebMsg)
	wi.FromWebSocketChannel = make(chan WebMsg)
	wi.HTTPPort = 0
	wi.httpDir = ""
	wi.ServerRunning = false
	wi.fcr = fcr

	wi.upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return wi
}

func (wi *WebInterface) newRouter() *mux.Router {
	var webroot string

	// determine the path to find the local html files
	// based on the current executable path
	exepath, _ := os.Executable()
	path, err := filepath.Abs(filepath.Dir(exepath))

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

	wi.httpDir = filepath.ToSlash(webroot)

	utils.LogI.Printf("path to the local html files (%s)", wi.httpDir)

	if err != nil {
		utils.LogE.Printf("unable to find the current path to the local html files (%s)", err)
	}

	// set a router and a hander to accept messages over the websocket
	r := mux.NewRouter()
	r.HandleFunc("/ws", wi.wsHandler)
	r.HandleFunc("/scenario", wi.getScenariosHandler).Methods("GET")
	r.HandleFunc("/scenario", wi.putScenarioPlaybackHandler).Methods("PUT")
	r.HandleFunc("/scenario/{scenarioId}", wi.scenarioDataHandler).Methods("GET")

	r.HandleFunc("/rosco", wi.getECUConnectionStatus).Methods("GET")
	r.HandleFunc("/rosco/dataframe", wi.getECUDataframeHandler).Methods("GET")
	r.HandleFunc("/rosco/{command}", wi.getECUResponseHandler).Methods("GET")

	r.HandleFunc("/config", wi.getConfigHandler).Methods("GET")
	r.HandleFunc("/config", wi.updateConfigHandler).Methods("POST")

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir(wi.httpDir))
	r.Handle("/", fileServer)
	r.PathPrefix("/").Handler(fileServer).Methods("GET")

	return r
}

// RunHTTPServer run the server
func (wi *WebInterface) RunHTTPServer() {
	// Declare a new router
	wi.router = wi.newRouter()

	// We can then pass our router (after declaring all our routes) to this method
	// (where previously, we were leaving the secodn argument as nil)
	listener, err := net.Listen("tcp", ":8081")

	if err != nil {
		utils.LogE.Printf("error starting web interface (%s)", err)
	}

	wi.HTTPPort = listener.Addr().(*net.TCPAddr).Port

	utils.LogI.Printf("started http server on port %d", wi.HTTPPort)
	wi.ServerRunning = true

	err = http.Serve(listener, wi.router)

	if err != nil {
		utils.LogE.Printf("error starting web interface (%s)", err)
	}
}

// send message to the web interface over the websocket
func (wi *WebInterface) sendMessageToWebInterface(m WebMsg) {
	if wi.ws != nil {
		err := wi.ws.WriteJSON(m)
		if err != nil {
			utils.LogE.Printf("error sending message over websocket (%s)", err)
		} else {
			utils.LogI.Printf("%s send message over websocket", utils.SendToWebTrace)
		}
	} else {
		utils.LogW.Printf("%s unable to send message over websocket, connected?", utils.SendToWebTrace)
	}
}

// ListenToWebSocketChannelLoop loop for listening for messages over the ToWebSocketChannel
// these are messages that are to be passed to the web interface over the websocket
// from the backend application
// to be run as a go routine as the channel is coded to be non blocking
func (wi *WebInterface) ListenToWebSocketChannelLoop() {
	for {
		m := <-wi.ToWebSocketChannel
		wi.sendMessageToWebInterface(m)
		utils.LogI.Printf("%s sent message '%s : %s' on ToWebSocketChannel", utils.SendToWebTrace, m.Action, m.Data)
	}
}
