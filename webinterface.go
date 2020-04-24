package main

import (
	"andrewj.com/readmems/rosco"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zserge/webview"
	"golang.org/x/net/websocket"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type wsMsg struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

// channel for communicating from the Mems interface to the web interface
var memsToWebChannel = make(chan wsMsg)

// channel for communicating from the web interface to the Mems interface
var webToMemsChannel = make(chan wsMsg)

type commandEnum struct {
	Cmd, Val string
}

var commandMap = make(map[commandEnum]int)

// recieveMessage the messages
func recieveMessage(ws *websocket.Conn) {
	var err error

	go listenForMems(ws)

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			rosco.LogI.Println("Websocket connection broken")
			break
		}

		parseMessage(ws, reply)
	}
}

func parseMessage(ws *websocket.Conn, msg string) {
	var m wsMsg
	json.Unmarshal([]byte(msg), &m)
	rosco.LogI.Printf("parse: %s %s\r\n", m.Action, m.Data)

	// connect to ECU and start the data loop
	if m.Action == "connect" {
		go memsCommandResponseLoop(config)
	}

	if m.Action == "increase" || m.Action == "decrease" || m.Action == "command" {
		// send to the CommandResponse loop
		rosco.LogI.Printf("parsing command %s %s, sending to channel", m.Action, m.Data)
		select {
		case webToMemsChannel <- m:
		case <-time.After(1000 * time.Millisecond):
			{
				rosco.LogE.Printf("Command channel blocked")
			}
		}
	}
}

// SendMessage to the web interface
func SendMessage(ws *websocket.Conn, m wsMsg) {
	msg, _ := json.Marshal(m)
	websocket.Message.Send(ws, string(msg))
}

func listenForMems(ws *websocket.Conn) {
	time.Sleep(1000 * time.Millisecond)

	for {
		data := <-memsToWebChannel // receive from mems interface
		rosco.LogI.Printf("listen: %s %s\r\n", data.Action, data.Data)

		SendMessage(ws, data)
	}
}

func newRouter() *mux.Router {
	exepath, err := os.Executable()
	path, err := filepath.Abs(filepath.Dir(exepath))
	if err != nil {
		rosco.LogI.Println(err)
	}

	httpdir := fmt.Sprintf("%s/public", path)

	r := mux.NewRouter()
	ws := websocket.Handler(recieveMessage)

	r.Handle("/", ws)

	//dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	//r.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(dir+"/public"))))

	// Declare the static file directory and point it to the
	// directory we just made
	staticFileDirectory := http.Dir(httpdir)

	// Declare the handler, that routes requests to their respective filename.
	// The fileserver is wrapped in the `stripPrefix` method, because we want to
	// remove the "/assets/" prefix when looking for files.
	// For example, if we type "/assets/index.html" in our browser, the file server
	// will look for only "index.html" inside the directory declared above.
	// If we did not strip the prefix, the file server would look for
	// "./assets/assets/index.html", and yield an error
	staticFileHandler := http.StripPrefix("/public/", http.FileServer(staticFileDirectory))

	// The "PathPrefix" method acts as a matcher, and matches all routes starting
	// with "/assets/", instead of the absolute route itself
	r.PathPrefix("/public").Handler(staticFileHandler).Methods("GET")

	return r
}

// RunHTTPServer run the server
func RunHTTPServer() {
	// Declare a new router
	r := newRouter()

	// We can then pass our router (after declaring all our routes) to this method
	// (where previously, we were leaving the secodn argument as nil)
	http.ListenAndServe(":1234", r)
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
	rosco.LogI.Printf("Evaluating %s, %s = %d", m.Action, m.Data, c)
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

	w.Navigate("http://127.0.0.1:1234/public/html/index.html")
	w.Run()
}
