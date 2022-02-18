package fcr

import (
	"fmt"
	"github.com/andrewdjackson/rosco"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
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
			// error occurred because the heartbeat failed to send
			// we'll assume the browser session has been terminated, clean up and close the server
			log.Warnf("unable to sent heartbeat to browser")
			webserver.DisconnectAndSaveScenario()
			webserver.TerminateApplication()
			return
		}
		flusher.Flush()
		time.Sleep(time.Second * 2)
	}
}

func (webserver *WebServer) DisconnectAndSaveScenario() {
	ecu := webserver.reader.ECU

	// preform any clean up before the application terminates

	if !ecu.Status.Emulated {
		if ecu.Status.Connected {
			// disconnect the ECU
			ecu.Disconnect()
		}
	}

	webserver.SaveScenario()
}

func (webserver *WebServer) SaveScenario() {
	ecu := webserver.reader.ECU

	// save the log file as a scenario file
	if ecu.Datalogger.Filename != "" {
		if ecu.Datalogger.Filename[len(ecu.Datalogger.Filename)-3:] == "csv" {
			// use the same filepath as the logfile but replace the .csv with .fcr
			f := strings.Replace(ecu.Datalogger.Filepath, ".csv", ".fcr", 1)
			s := rosco.NewScenarioFile(f)
			err := s.ConvertLogToScenario(ecu.Datalogger.Filename)
			if err == nil {
				err = s.Write()
				log.Infof("saved scenario as %s", f)
			}
		}
	}
}

func (webserver *WebServer) TerminateApplication() {
	log.Info("shutting down application")
	os.Exit(0)
}
