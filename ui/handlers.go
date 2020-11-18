package ui

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/andrewdjackson/memsfcr/rosco"
	"github.com/andrewdjackson/memsfcr/scenarios"
	"github.com/andrewdjackson/memsfcr/utils"
)

// REST API : GET Scenario
// returns the details of the specified scenario
func (wi *WebInterface) scenarioDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scenarioID := vars["scenarioId"]

	data := scenarios.GetScenario(scenarioID)

	utils.LogI.Printf("%+v", data)

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

// REST API : PUT Scenario Playback State
// changes the state of the scenario playback
func (wi *WebInterface) putScenarioPlaybackHandler(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	reqBody, _ := ioutil.ReadAll(r.Body)

	// get the current configuration
	var status scenarios.ScenarioDescription
	_ = json.Unmarshal(reqBody, &status)

	utils.LogI.Printf("%+v", status)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(status); err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// REST API : GET Scenarios
// returns a list of available scenarios
func (wi *WebInterface) getScenariosHandler(w http.ResponseWriter, r *http.Request) {
	scenarios, _ := scenarios.GetScenarios()

	utils.LogI.Printf("%+v", scenarios)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(scenarios); err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// REST API : GET Config
// returns the contents of the Config file as a JSON response
func (wi *WebInterface) getConfigHandler(w http.ResponseWriter, r *http.Request) {
	config := utils.ReadConfig()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(config); err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// REST API : POST Config
// updates the config
func (wi *WebInterface) updateConfigHandler(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Config struct
	reqBody, _ := ioutil.ReadAll(r.Body)

	// get the current configuration
	config := utils.ReadConfig()
	_ = json.Unmarshal(reqBody, &config)

	utils.LogI.Printf("%s REST updateConfig (%v)", utils.ReceiveFromWebTrace, config)
	// save the configuration
	utils.WriteConfig(config)

	// return a 200 status code
	w.WriteHeader(http.StatusOK)
}

// REST API : GET ECU Response
// send the specified command and returns the response data
func (wi *WebInterface) getECUResponseHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	command, err := hex.DecodeString(vars["command"])

	if err == nil {
		response, _ := wi.fcr.ECU.SendCommand(command)
		var ecu rosco.MemsCommandResponse
		ecu.Command = command
		ecu.Response = response

		utils.LogI.Printf("%s REST getECUResponseHandler (%v)", utils.ReceiveFromWebTrace, ecu)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if err := json.NewEncoder(w).Encode(r); err != nil {
			// return a error code
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		utils.LogE.Printf("%s error in REST getECUResponseHandler (%v)(%v)", utils.ReceiveFromWebTrace, command, err)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
	}
}

// REST API : GET ECU Response
// send the specified command and returns the response data
func (wi *WebInterface) getECUDataframeHandler(w http.ResponseWriter, r *http.Request) {
	memsdata := wi.fcr.ECU.GetDataframes()

	utils.LogI.Printf("%s REST getECUResponseHandler (%v)", utils.ReceiveFromWebTrace, memsdata)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(memsdata); err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// REST API : GET ECU Response
// send the specified command and returns the response data
func (wi *WebInterface) getECUConnectionStatus(w http.ResponseWriter, r *http.Request) {
	var status rosco.MemsConnectionStatus

	status.Connected = wi.fcr.ECU.Connected
	status.Initialised = wi.fcr.ECU.Initialised
	status.ECUID = fmt.Sprintf("%X", wi.fcr.ECU.ECUID)
	status.IACPosition = wi.fcr.ECU.Diagnostics.Analysis.IACPosition

	utils.LogI.Printf("%s REST getECUConnectionStatus (%v)", utils.ReceiveFromWebTrace, status)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(status); err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// WEBSOCKET Handler
// initite loop listening for data sent to the websocket
// and sendind data recieved from the channel over the websocket
func (wi *WebInterface) wsHandler(w http.ResponseWriter, r *http.Request) {
	var m WebMsg
	var err error

	// upgrade the http connection to a websocket
	wi.ws, err = wi.upgrader.Upgrade(w, r, nil)
	defer wi.ws.Close()

	if err != nil {
		utils.LogE.Printf("error in websocket (%s)", err)
	}

	// read loop, if a message is recieved over the websocket
	// then post it into the FromWeb communication channel
	// this is configured not to block if the channel is unable to
	// receive.
	for {
		err = wi.ws.ReadJSON(&m)
		if err != nil {
			utils.LogE.Fatalf("error in websocket (%s)", err)
		} else {
			utils.LogI.Printf("%s recieved websocket message (%v)", utils.ReceiveFromWebTrace, m)
		}

		select {
		case wi.FromWebSocketChannel <- m:
			utils.LogI.Printf("%s sent message to FromWebChannel (%v)", utils.ReceiveFromWebTrace, m)
		default:
		}
	}
}
