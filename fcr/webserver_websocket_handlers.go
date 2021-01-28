package fcr

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Websocket Handler
// listens to data sent to the websocket
// and is used to exit the application if the browser is closed
// triggering a websocket disconnect
func (webserver *WebServer) websocketHandler(w http.ResponseWriter, r *http.Request) {
	var m WebMsg
	var err error

	// upgrade the http connection to a websocket
	webserver.ws, err = webserver.upgrader.Upgrade(w, r, nil)
	defer webserver.ws.Close()

	if err != nil {
		log.Errorf("error in websocket (%s)", err)
	}

	// read loop, if a message is received over the websocket
	// then post it into the FromWeb communication channel
	// this is configured not to block if the channel is unable to
	// receive.
	for {
		err = webserver.ws.ReadJSON(&m)
		if err != nil {
			log.Errorf("error in websocket (%s)", err)
		} else {
			log.Infof("recieved websocket message (%v)", m)
		}

		select {
		case webserver.FromWebSocketChannel <- m:
			log.Infof("sent message to FromWebChannel (%v)", m)
		default:
		}
	}
}

func (webserver *WebServer) browserHeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "your browser doesn't support server-sent events", 503)
		return
	} else {
		log.Info("connected browser heartbeat")
	}

	// send a heartbeat to prevent connection timeout
	for {
		_, err := fmt.Fprintf(w, "event: heartbeat\ndata: heartbeat\n\n")

		if err != nil {
			log.Fatal("unable to sent heartbeat to browser, connection closed")
			return
		}
		flusher.Flush()
		time.Sleep(time.Second * 2)
	}
}
