package ui

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/andrewdjackson/readmems/utils"
	"github.com/gorilla/mux"
	"golang.org/x/net/websocket"
)

type wsMsg struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

// WebAction constants
const (
	WebActionConfig             = "config"
	WebActionConnection         = "connection"
	WebActionConnect            = "connect"
	WebActionECUCommand         = "command"
	WebActionECUCommandIncrease = "command"
	WebActionECUCommandDecrease = "command"
	WebActionData               = "data"
)

// UI command map
const (
	commandUnknown                 = 0
	commandConnectECU              = 1
	commandPauseDataLoop           = 2
	commandStartDataLoop           = 3
	commandResetECU                = 4
	commandResetAdjustments        = 5
	commandClearFaults             = 6
	commandIncreaseIdleSpeed       = 7
	commandIncreaseIdleHot         = 8
	commandIncreaseFuelTrim        = 9
	commandIncreaseIgnitionAdvance = 10
	commandDecreaseIdleSpeed       = 11
	commandDecreaseIdleHot         = 12
	commandDecreaseFuelTrim        = 13
	commandDecreaseIgnitionAdvance = 14
)

type commandEnum struct {
	Cmd, Val string
}

// WebInterface the web interface
type WebInterface struct {
	// mulitplex router interface
	router *mux.Router

	// websocket interface
	ws      *websocket.Conn
	httpDir string

	// map of valid commands
	commandMap map[commandEnum]int

	// HTTPPort used by the HTTP Server instance
	HTTPPort int

	// channel for communication to the web interface
	ToWebChannel chan wsMsg

	// channel for communication from the web interface
	FromWebChannel chan wsMsg
}

// NewWebInterface creates a new web interface
func NewWebInterface() *WebInterface {
	wi := &WebInterface{}
	wi.ToWebChannel = make(chan wsMsg)
	wi.FromWebChannel = make(chan wsMsg)
	wi.intitialiseCommandMap()
	wi.HTTPPort = 0
	wi.httpDir = ""

	return wi
}

func (wi *WebInterface) intitialiseCommandMap() {
	wi.commandMap = make(map[commandEnum]int)
	wi.commandMap[commandEnum{"command", "connect"}] = commandConnectECU
	wi.commandMap[commandEnum{"command", "resetecu"}] = commandResetECU
	wi.commandMap[commandEnum{"command", "resetadj"}] = commandResetAdjustments
	wi.commandMap[commandEnum{"command", "clearfaults"}] = commandClearFaults
	wi.commandMap[commandEnum{"command", "pause"}] = commandPauseDataLoop
	wi.commandMap[commandEnum{"command", "start"}] = commandStartDataLoop
	wi.commandMap[commandEnum{"increase", "idlespeed"}] = commandIncreaseIdleSpeed
	wi.commandMap[commandEnum{"increase", "idlehot"}] = commandIncreaseIdleHot
	wi.commandMap[commandEnum{"increase", "fueltrim"}] = commandIncreaseFuelTrim
	wi.commandMap[commandEnum{"increase", "ignition"}] = commandIncreaseIgnitionAdvance
	wi.commandMap[commandEnum{"decrease", "idlespeed"}] = commandDecreaseIdleSpeed
	wi.commandMap[commandEnum{"decrease", "idlehot"}] = commandDecreaseIdleHot
	wi.commandMap[commandEnum{"decrease", "fueltrim"}] = commandDecreaseFuelTrim
	wi.commandMap[commandEnum{"decrease", "ignition"}] = commandDecreaseIgnitionAdvance
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
	ws := websocket.Handler(wi.recieveMessageFromWebInterface)
	r.Handle("/", ws)

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

// receive message handler to receive data sent over the websocket
// these are messages that are to be passed from the web interface
// to the backend application using the FromWebChannel
func (wi *WebInterface) recieveMessageFromWebInterface(ws *websocket.Conn) {
	var err error
	var msg string
	var m wsMsg

	if wi.ws == nil {
		wi.ws = ws
	}

	if err = websocket.Message.Receive(wi.ws, &msg); err != nil {
		utils.LogE.Printf("%s websocket connection broken", utils.SendToWebTrace)
	} else {
		// parse the received message
		json.Unmarshal([]byte(msg), &m)
		utils.LogI.Printf("parse: %s %s\r\n", m.Action, m.Data)

		// send over the FromWebChannel
		// use select to ensure the weboscket receiver is not blocked
		select {
		case wi.FromWebChannel <- m:
			utils.LogI.Printf("%s waiting to send message %s %s on FromWebChannel", utils.ReceiveFromWebTrace, m.Action, m.Data)
		default:
			utils.LogW.Printf("%s unable to send messgae on FromWebChannel, blocked?", utils.ReceiveFromWebTrace)
		}
	}
}

// send message to the web interface over the websocket
func (wi *WebInterface) sendMessageToWebInterface(m wsMsg) {
	msg, _ := json.Marshal(m)
	websocket.Message.Send(wi.ws, string(msg))
}

// loop for listening for messages over the ToWebChannel
// these are messages that are to be passed to the web interface over the websocket
// from the backend application
func (wi *WebInterface) listenToWebChannelLoop() {
	for {
		m := <-wi.ToWebChannel
		wi.sendMessageToWebInterface(m)
		utils.LogI.Printf("%s sent message %s %s on FromWebChannel", utils.SendToWebTrace, m.Action, m.Data)
	}
}

/*
// loop for listening for messages over the FromWebChannel
// these are messages received over the websocket that are from the web interface
// to be passed to the backend application
func (wi *WebInterface) listenFromWebChannelLoop() {
	for {
		utils.LogI.Printf("%s waiting to send message %s %s on FromWebChannel", utils.SendToWebTrace, m.Action, m.Data)
		wi.FromWebChannel <- m
		utils.LogI.Printf("%s sent message on FromWebChannel", utils.SendToWebTrace)
	}
}
*/
/////////////////////////////////////////////////////
//
//
/*
type commandEnum struct {
	Cmd, Val string
}

var commandMap = make(map[commandEnum]int)
var httpPort = 0

// loop waiting for data received from the ECU to be sent on to the web interface
func listenForDataSentFromECULoop(ws *websocket.Conn) {
	// wait for web interface to finish loading
	utils.LogI.Printf("%s waiting for data from memsToWebChannel..", utils.ReceiveFromWebTrace)

	for {
		select {
		case data := <-memsToWebChannel:
			utils.LogI.Printf("%s received %s %s from memsToWebChannel", utils.ReceiveFromWebTrace, data.Action, data.Data)
			SendMessage(ws, data)
		default:
		}

		// and breath..
		time.Sleep(50 * time.Millisecond)
	}
}

// RunHTTPServer run the server
func RunHTTPServer() {
	// Declare a new router
	r := newRouter()

	// We can then pass our router (after declaring all our routes) to this method
	// (where previously, we were leaving the secodn argument as nil)
	//http.ListenAndServe(":0", r)
	listener, err := net.Listen("tcp", ":0")

	if err != nil {
		panic(err)
	}

	httpPort = listener.Addr().(*net.TCPAddr).Port
	http.Serve(listener, r)
}

const commandUnknown = 0
const commandConnectECU = 1
const commandPauseDataLoop = 2
const commandStartDataLoop = 3
const commandResetECU = 4
const commandResetAdjustments = 5
const commandClearFaults = 6
const commandIncreaseIdleSpeed = 7
const commandIncreaseIdleHot = 8
const commandIncreaseFuelTrim = 9
const commandIncreaseIgnitionAdvance = 10
const commandDecreaseIdleSpeed = 11
const commandDecreaseIdleHot = 12
const commandDecreaseFuelTrim = 13
const commandDecreaseIgnitionAdvance = 14

func evaluateCommand(m wsMsg) int {
	c := commandMap[commandEnum{m.Action, m.Data}]
	utils.LogI.Printf("Evaluating %s, %s = %d", m.Action, m.Data, c)
	return c
}

func createCommandMap() {
	commandMap[commandEnum{"command", "connect"}] = commandConnectECU
	commandMap[commandEnum{"command", "resetecu"}] = commandResetECU
	commandMap[commandEnum{"command", "resetadj"}] = commandResetAdjustments
	commandMap[commandEnum{"command", "clearfaults"}] = commandClearFaults
	commandMap[commandEnum{"command", "pause"}] = commandPauseDataLoop
	commandMap[commandEnum{"command", "start"}] = commandStartDataLoop
	commandMap[commandEnum{"increase", "idlespeed"}] = commandIncreaseIdleSpeed
	commandMap[commandEnum{"increase", "idlehot"}] = commandIncreaseIdleHot
	commandMap[commandEnum{"increase", "fueltrim"}] = commandIncreaseFuelTrim
	commandMap[commandEnum{"increase", "ignition"}] = commandIncreaseIgnitionAdvance
	commandMap[commandEnum{"decrease", "idlespeed"}] = commandDecreaseIdleSpeed
	commandMap[commandEnum{"decrease", "idlehot"}] = commandDecreaseIdleHot
	commandMap[commandEnum{"decrease", "fueltrim"}] = commandDecreaseFuelTrim
	commandMap[commandEnum{"decrease", "ignition"}] = commandDecreaseIgnitionAdvance
}

// ShowWebView show the browser
func ShowWebView(config *rosco.ReadmemsConfig) {
	w := webview.New(true)
	defer w.Destroy()

	w.SetTitle("MEMSFCR")
	w.SetSize(1120, 920, webview.HintNone)

	w.Bind("quit", func() {
		w.Terminate()
	})

	url := fmt.Sprintf("http://127.0.0.1:%d/public/html/index.html", httpPort)

	w.Navigate(url)
	w.Run()
}
*/
