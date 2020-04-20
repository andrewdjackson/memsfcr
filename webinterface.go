package main

import (
	"andrewj.com/readmems/rosco"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zserge/webview"
	"golang.org/x/net/websocket"
	"net/http"
)

type wsMsg struct {
	Action string `json:"action"`
	Data   string `json:"data"`
}

var memsChannel = make(chan wsMsg)

// Echo the messages
func Echo(ws *websocket.Conn) {
	var err error

	go listenForMems(ws)

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}

		parseMessage(reply)
		/*
			msg := reply
			fmt.Println("Sending to client: " + msg)

			if err = websocket.Message.Send(ws, msg); err != nil {
				fmt.Println("Can't send")
				break
			}
		*/
	}
}

func parseMessage(msg string) {
	var m wsMsg
	json.Unmarshal([]byte(msg), &m)
	fmt.Printf("parse: %s %s\r\n", m.Action, m.Data)

	if m.Action == "connect" {
		go memsCommandResponseLoop(config)
	}
}

// SendMessage to the
func SendMessage(ws *websocket.Conn, m wsMsg) {
	msg, _ := json.Marshal(m)
	websocket.Message.Send(ws, string(msg))
}

func listenForMems(ws *websocket.Conn) {
	for {
		data := <-memsChannel // receive from mems channel
		fmt.Printf("listen: %s %s\r\n", data.Action, data.Data)

		SendMessage(ws, data)
	}
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	ws := websocket.Handler(Echo)

	r.Handle("/", ws)

	//dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	//r.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(dir+"/public"))))

	// Declare the static file directory and point it to the
	// directory we just made
	staticFileDirectory := http.Dir("./public")

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
