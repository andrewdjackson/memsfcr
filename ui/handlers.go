package ui

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/andrewdjackson/memsfcr/scenarios"
	"github.com/andrewdjackson/memsfcr/utils"
)

func (wi *WebInterface) scenarioDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	scenarioID := vars["scenarioId"]

	data := scenarios.GetScenario(scenarioID)

	utils.LogI.Printf("%+v", data)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

func (wi *WebInterface) scenarioHandler(w http.ResponseWriter, r *http.Request) {
	scenarios, _ := scenarios.GetScenarios()

	utils.LogI.Printf("%+v", scenarios)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(scenarios); err != nil {
		panic(err)
	}
}

func (wi *WebInterface) configHandler(w http.ResponseWriter, r *http.Request) {
	config := utils.ReadConfig()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(config); err != nil {
		panic(err)
	}
}

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
		case wi.FromWebChannel <- m:
			utils.LogI.Printf("%s sent message to FromWebChannel (%v)", utils.ReceiveFromWebTrace, m)
		default:
		}
	}
}
