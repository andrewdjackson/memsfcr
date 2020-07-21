package ui

import (
	"encoding/json"
	"net/http"

	"github.com/andrewdjackson/memsfcr/utils"
)

func (wi *WebInterface) scenarioHandler(w http.ResponseWriter, r *http.Request) {
	scenarios, _ := utils.GetScenarios()

	utils.LogI.Printf("%+v", scenarios)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(scenarios); err != nil {
		panic(err)
	}
}

func (wi *WebInterface) configHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(""); err != nil {
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
