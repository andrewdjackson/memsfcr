package fcr

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

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
