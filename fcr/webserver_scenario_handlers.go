package fcr

import (
	"encoding/json"
	"github.com/andrewdjackson/rosco"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// REST API : GET Scenario
// returns the details of the specified scenario
func (webserver *WebServer) getScenarioDetails(w http.ResponseWriter, r *http.Request) {
	log.Info("rest get scenario details")

	vars := mux.Vars(r)
	scenarioID := vars["scenarioId"]

	defer r.Body.Close()

	data := rosco.GetScenario(scenarioID)

	log.Infof("%+v", data)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if data.Count == 0 {
		// return 404 not found
		w.WriteHeader(http.StatusNotFound)
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// REST API : GET Scenarios
// returns a list of available scenarios
func (webserver *WebServer) getListofScenarios(w http.ResponseWriter, r *http.Request) {
	log.Info("rest-get list of scenarios")
	scenarios, _ := rosco.GetScenarios()

	log.Infof("%+v", scenarios)
	webserver.sendResponse(w, r, scenarios)
}

func (webserver *WebServer) sendResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}
}
