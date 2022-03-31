package fcr

import (
	"fmt"
	"github.com/andrewdjackson/rosco"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"reflect"
	"time"
)

func (webserver *WebServer) browserHeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	var flusher http.Flusher
	var supported bool

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if flusher, supported = w.(http.Flusher); !supported {
		http.Error(w, "your browser doesn't support server-sent events", 503)
		return
	} else {
		log.Info("connected browser heartbeat")
	}

	// send a heartbeat to prevent connection timeout
	for {
		if _, err := fmt.Fprintf(w, "event: heartbeat\ndata: heartbeat\n\n"); err != nil {
			// error occurred because the heartbeat failed to send
			// we'll assume the browser session has been terminated, clean up and close the server
			log.Warnf("unable to sent heartbeat to browser, terminating application")
			webserver.Disconnect()
			webserver.TerminateApplication()
		}

		flusher.Flush()
		// wait time between heartbeats
		time.Sleep(time.Second * 5)
	}
}

func (webserver *WebServer) Disconnect() {
	// disconnect the ECU
	if reflect.TypeOf(webserver.reader.ECU) == reflect.TypeOf(&rosco.MEMSReader{}) {
		log.Infof("diconnecting from the ecu")
		if err := webserver.reader.ECU.Disconnect(); err != nil {
			log.Warnf("error disconnecting from the ecu (%s)", err)
		}
	}
}

func (webserver *WebServer) TerminateApplication() {
	log.Info("shutting down application")
	os.Exit(0)
}
