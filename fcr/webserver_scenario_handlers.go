package fcr

import (
	"encoding/hex"
	"encoding/json"
	"github.com/andrewdjackson/rosco"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type ScenarioDetail struct {
	Timestamp   time.Time
	Position    int
	Dataframe7d string
	Dataframe80 string
}

type ScenarioDetails struct {
	First   ScenarioDetail
	Current ScenarioDetail
	Last    ScenarioDetail
}

type ScenarioSeekPosition struct {
	CurrentPosition int
	NewPosition     int
}

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

func (webserver *WebServer) getPlaybackDetails(w http.ResponseWriter, r *http.Request) {
	log.Info("rest-get scenario playback details")

	vars := mux.Vars(r)
	if len(vars) > 0 {
		scenarioID := vars["scenarioId"]
		log.Info("rest-get scenario playback id %v" + scenarioID)
	}

	details := ScenarioDetails{}

	d, err := webserver.reader.ECU.Responder.GetFirst()
	details.First = ScenarioDetail{}
	if err == nil {
		details.First.Timestamp = d.Timestamp
		details.First.Position = d.Position
		details.First.Dataframe80 = hex.EncodeToString(d.Dataframe80)
		details.First.Dataframe7d = hex.EncodeToString(d.Dataframe7d)
	}

	d, err = webserver.reader.ECU.Responder.GetCurrent()
	details.Current = ScenarioDetail{}
	if err == nil {
		details.Current.Timestamp = d.Timestamp
		details.Current.Position = d.Position
		details.Current.Dataframe80 = hex.EncodeToString(d.Dataframe80)
		details.Current.Dataframe7d = hex.EncodeToString(d.Dataframe7d)
	}

	d, err = webserver.reader.ECU.Responder.GetLast()
	details.Last = ScenarioDetail{}
	if err == nil {
		details.Last.Timestamp = d.Timestamp
		details.Last.Position = d.Position
		details.Last.Dataframe80 = hex.EncodeToString(d.Dataframe80)
		details.Last.Dataframe7d = hex.EncodeToString(d.Dataframe7d)
	}

	log.Infof("%+v", details)
	webserver.sendResponse(w, r, details)
}

func (webserver *WebServer) postPlaybackSeek(w http.ResponseWriter, r *http.Request) {
	// get the body of our request
	// unmarshal this into a new Config struct
	reqBody, _ := ioutil.ReadAll(r.Body)

	// get the current configuration
	position := ScenarioSeekPosition{}
	_ = json.Unmarshal(reqBody, &position)

	log.Infof("rest-post scenario playback seek (%+v)", position)

	if webserver.reader.ECU.Status.Connected && webserver.reader.ECU.Status.Emulated {
		last := webserver.reader.ECU.Responder.Playbook.Count

		if position.NewPosition < last {
			webserver.reader.ECU.Responder.MoveToPosition(position.NewPosition)
			detail, _ := webserver.reader.ECU.Responder.GetCurrent()
			log.Infof("rest-post scenario position moved from %v to %v", position.CurrentPosition, detail.Position)

			webserver.sendResponse(w, r, detail)
			w.WriteHeader(http.StatusOK)
		} else {
			log.Infof("rest-post scenario position too far (%v > %v)", position.NewPosition, last)
			// position not found
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		// service unavailable if we're not replaying
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func (webserver *WebServer) sendResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}
}
