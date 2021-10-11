package fcr

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial.v1"
	"io/ioutil"
	"net/http"
)

type AvailablePorts struct {
	Ports []string `json:"ports"`
}

// REST API : GET Config
// returns the contents of the Config file as a JSON response
func (webserver *WebServer) getConfigHandler(w http.ResponseWriter, r *http.Request) {
	config := webserver.reader.Config
	log.Infof("rest-get config (%v)", config)

	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(config); err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// REST API : PUT Config
// updates the config
func (webserver *WebServer) updateConfigHandler(w http.ResponseWriter, r *http.Request) {
	// get the body of our request
	// unmarshal this into a new Config struct
	reqBody, _ := ioutil.ReadAll(r.Body)

	// get the current configuration
	config := ReadConfig()
	_ = json.Unmarshal(reqBody, &config)

	log.Infof("rest-put update config (%v)", config)
	// save the configuration
	WriteConfig(config)

	// return a 200 status code
	w.WriteHeader(http.StatusOK)
}

// rest-api get list of available serial ports
func (webserver *WebServer) getSerialPortsHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("rest-get available serial ports")

	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	ports := AvailablePorts{Ports: webserver.getSerialPorts()}

	if err := json.NewEncoder(w).Encode(ports); err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// enumerate the available serial ports
// this won't enumerate virtual ports
func (webserver *WebServer) getSerialPorts() []string {
	ports, err := serial.GetPortsList()

	if err != nil {
		log.Error("error enumerating serial ports")
	}
	if len(ports) == 0 {
		log.Warn("unable to find any serial ports")
	}
	for _, port := range ports {
		log.Infof("found serial port %v", port)
	}

	return ports
}
