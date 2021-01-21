package fcr

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// REST API : GET Config
// returns the contents of the Config file as a JSON response
func (webserver *WebServer) getConfigHandler(w http.ResponseWriter, r *http.Request) {
	config := ReadConfig()
	log.Infof("REST GET Config (%v)", config)

	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(config); err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// REST API : POST Config
// updates the config
func (webserver *WebServer) updateConfigHandler(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Config struct
	reqBody, _ := ioutil.ReadAll(r.Body)

	// get the current configuration
	config := ReadConfig()
	_ = json.Unmarshal(reqBody, &config)

	log.Infof("REST POST Update Config (%v)", config)
	// save the configuration
	WriteConfig(config)

	// return a 200 status code
	w.WriteHeader(http.StatusOK)
}
