package ui

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/andrewdjackson/readmems/utils"
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

	// channel for communication to the web interface
	ToWebChannel chan WebMsg

	// channel for communication from the web interface
	FromWebChannel chan WebMsg
}

// NewWebInterface creates a new web interface
func NewWebInterface() *WebInterface {
	wi := &WebInterface{}
	wi.ToWebChannel = make(chan WebMsg)
	wi.FromWebChannel = make(chan WebMsg)
	wi.HTTPPort = 0
	wi.httpDir = ""

	wi.upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return wi
}

func (wi *WebInterface) newRouter() *mux.Router {
	// determine the path to find the local html files
	// based on the current executable path
	exepath, err := os.Executable()
	path, err := filepath.Abs(filepath.Dir(exepath))

	if err != nil {
		utils.LogE.Printf("unable to find the current path to the local html files (%s)", err)
	}

	wi.httpDir = fmt.Sprintf("%s/public", path)

	// set a router and a hander to accept messages over the websocket
	r := mux.NewRouter()
	r.HandleFunc("/", wi.wsHandler)

	// Declare the static file directory and point it to the
	// directory we just made
	staticFileDirectory := http.Dir(wi.httpDir)

	// Declare the handler, that routes requests to their respective filename.
	staticFileHandler := http.StripPrefix("/public/", http.FileServer(staticFileDirectory))

	// The "PathPrefix" method acts as a matcher, and matches all routes starting
	// with "/pulic/", instead of the absolute route itself
	r.PathPrefix("/public").Handler(staticFileHandler).Methods("GET")

	return r
}

// RunHTTPServer run the server
func (wi *WebInterface) RunHTTPServer() {
	// Declare a new router
	wi.router = wi.newRouter()

	// We can then pass our router (after declaring all our routes) to this method
	// (where previously, we were leaving the secodn argument as nil)
	listener, err := net.Listen("tcp", ":0")

	if err != nil {
		utils.LogE.Printf("error starting web interface (%s)", err)
	}

	wi.HTTPPort = listener.Addr().(*net.TCPAddr).Port
	http.Serve(listener, wi.router)
}

func (wi *WebInterface) wsHandler(w http.ResponseWriter, r *http.Request) {
	var m WebMsg
	var err error

	// upgrade the http connection to a websocket
	wi.ws, err = wi.upgrader.Upgrade(w, r, nil)
	defer wi.ws.Close()

	if err != nil {
		utils.LogE.Printf("error in websocket (%s)", err)
	}

	// read loop, if a message is recieved over the websocket
	// then post it into the FromWeb communication channel
	// this is configured not to block if the channel is unable to
	// receive.
	for {
		wi.ws.ReadJSON(&m)
		utils.LogI.Printf("%s recieved websocket message (%v)", utils.ReceiveFromWebTrace, m)

		select {
		case wi.FromWebChannel <- m:
			utils.LogI.Printf("%s sent message to FromWebChannel (%v)", utils.ReceiveFromWebTrace, m)
		default:
		}
	}
}

// send message to the web interface over the websocket
func (wi *WebInterface) sendMessageToWebInterface(m WebMsg) {
	if wi.ws != nil {
		wi.ws.WriteJSON(m)
		utils.LogI.Printf("%s send message over websocket", utils.SendToWebTrace)
	} else {
		utils.LogW.Printf("%s unable to send message over websocket, connected?", utils.SendToWebTrace)
	}
}

// ListenToWebChannelLoop loop for listening for messages over the ToWebChannel
// these are messages that are to be passed to the web interface over the websocket
// from the backend application
// to be run as a go routine as the channel is coded to be non blocking
func (wi *WebInterface) ListenToWebChannelLoop() {
	for {
		m := <-wi.ToWebChannel
		wi.sendMessageToWebInterface(m)
		utils.LogI.Printf("%s sent message '%s : %s' on ToWebChannel", utils.SendToWebTrace, m.Action, m.Data)
	}
}
