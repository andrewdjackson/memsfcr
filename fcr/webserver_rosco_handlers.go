package fcr

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type ECUConnectionPort struct {
	Port string `json:"port"`
}

type ECUAdjustment struct {
	Steps int `json:"steps"`
}

type ECUActivate struct {
	Activate bool `json:"activate"`
}

type ActionResponse struct {
	Success bool `json:"success"`
}

type AdjustmentResponse struct {
	Value int `json:"value"`
}

//
// Connection Status
// returns the status of the ecu connection along with the ecu id and the iac initial position
//
func (webserver *WebServer) getECUConnectionStatus(w http.ResponseWriter, r *http.Request) {
	log.Infof("rest-get read ecu status")
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	status := webserver.reader.ECU.Status
	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Warnf("rest-post response failed")
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//
// Connect and Initialise the ECU
//
func (webserver *WebServer) postECUConnect(w http.ResponseWriter, r *http.Request) {
	log.Infof("rest-post connect ecu")
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if webserver.reader.ECU.Status.Connected {
		log.Warnf("rest-post already connected to the ecu")
		// return status if already connected
		w.WriteHeader(http.StatusAlreadyReported)
	} else {
		// get the body of our POST request
		// unmarshal this into a new Config struct
		reqBody, _ := ioutil.ReadAll(r.Body)

		// get the current configuration
		var port ECUConnectionPort
		_ = json.Unmarshal(reqBody, &port)

		log.Infof("rest-post connecting ecu (%v)", port)

		webserver.reader.ECU.ConnectAndInitialiseECU(port.Port)

		if webserver.reader.ECU.Status.Connected {
			// return a 200 status code
			w.WriteHeader(http.StatusOK)
		} else {
			log.Warnf("rest-post unable to connect to the ecu")
			// return service unavailable if unable to connect
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}

	if err := json.NewEncoder(w).Encode(webserver.reader.ECU.Status); err != nil {
		log.Warnf("rest-post response failed")
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//
// Disconnect the ECU
//
func (webserver *WebServer) postECUDisconnect(w http.ResponseWriter, r *http.Request) {
	log.Infof("rest-post disconnect ecu")
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if !webserver.reader.ECU.Status.Connected {
		// return status if already disconnected
		w.WriteHeader(http.StatusAlreadyReported)
	} else {
		webserver.reader.ECU.Disconnect()

		if !webserver.reader.ECU.Status.Connected {
			// return a 200 status code
			w.WriteHeader(http.StatusOK)
		} else {
			log.Warnf("rest-post unable to disconnect the ecu")
			// return service unavailable if unable to connect
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}

	if err := json.NewEncoder(w).Encode(webserver.reader.ECU.Status); err != nil {
		log.Warnf("rest-post response failed")
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//
// Read the Dataframes from the ECU
// the dataframes contain the engine running parameters and fault codes
//
func (webserver *WebServer) getECUDataframes(w http.ResponseWriter, r *http.Request) {
	log.Infof("rest-get read ecu dataframes")
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if webserver.isECUConnected(w) {
		memsdata := webserver.reader.ECU.GetDataframes()

		log.Infof("rest-get ecu dataframes (%v)", memsdata)

		if err := json.NewEncoder(w).Encode(memsdata); err != nil {
			log.Warnf("rest-get response failed")
			// return a error code
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

//
// Read the IAC Position
// the iac position is used by the ecu to adjust the air fuel ratio
// by incrementing and decrementing the stepper motor.
// the ecu has no feedback from the stepper motor, the iac position is a calculated position
//
func (webserver *WebServer) getECUIAC(w http.ResponseWriter, r *http.Request) {
	log.Infof("rest-get read ecu iac position")
	defer r.Body.Close()
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if webserver.isECUConnected(w) {
		value := webserver.reader.ECU.GetIACPosition()
		response := AdjustmentResponse{Value: value}

		log.Infof("rest-get ecu iac position (%v)", value)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Warnf("rest-get response failed")
			// return a error code
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

//
//  send heartbeat the ecu
//
func (webserver *WebServer) postECUHeartbeat(w http.ResponseWriter, r *http.Request) {
	log.Infof("rest-post send heartbeat")
	value := webserver.reader.ECU.ResetECU()
	webserver.updateECUState(w, r, value)
}

//
//  reset the ecu
//
func (webserver *WebServer) postECUReset(w http.ResponseWriter, r *http.Request) {
	log.Infof("rest-post reset ecu")
	value := webserver.reader.ECU.ResetECU()
	webserver.updateECUState(w, r, value)
}

//
// clear the fault codes
//
func (webserver *WebServer) postECUClearFaults(w http.ResponseWriter, r *http.Request) {
	log.Infof("rest-post clear ecu faults")
	value := webserver.reader.ECU.ClearFaults()
	webserver.updateECUState(w, r, value)
}

//
// clear the adjustable values
//
func (webserver *WebServer) postECUClearAdjustments(w http.ResponseWriter, r *http.Request) {
	log.Infof("rest-post clear ecu adjustable values")
	value := webserver.reader.ECU.ResetAdjustments()
	webserver.updateECUState(w, r, value)
}

//
// update ecu state (clear faults, reset and heartbeat)
//
func (webserver *WebServer) updateECUState(w http.ResponseWriter, r *http.Request, value bool) {
	if webserver.isECUConnected(w) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		value := webserver.reader.ECU.ResetECU()
		response := ActionResponse{Success: value}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Warnf("rest-call response failed")
			// return a error code
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

//
// update the short term fuel trim
//
func (webserver *WebServer) postECUAdjustSTFT(w http.ResponseWriter, r *http.Request) {
	var data ECUAdjustment

	log.Infof("rest-post update ecu stft")

	// get the body of our POST request
	// unmarshal this into a new Config struct
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)

	value := webserver.reader.ECU.AdjustShortTermFuelTrim(data.Steps)
	webserver.updateAdjustableValue(w, r, value)
}

//
// update the long term fuel trim
//
func (webserver *WebServer) postECUAdjustLTFT(w http.ResponseWriter, r *http.Request) {
	var data ECUAdjustment

	log.Infof("rest-post update ecu ltft")

	// get the body of our POST request
	// unmarshal this into a new Config struct
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)

	value := webserver.reader.ECU.AdjustLongTermFuelTrim(data.Steps)
	webserver.updateAdjustableValue(w, r, value)
}

//
// update the idle decay
//
func (webserver *WebServer) postECUAdjustIdleDecay(w http.ResponseWriter, r *http.Request) {
	var data ECUAdjustment

	log.Infof("rest-post update ecu idle decay")

	// get the body of our POST request
	// unmarshal this into a new Config struct
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)

	value := webserver.reader.ECU.AdjustIdleDecay(data.Steps)
	webserver.updateAdjustableValue(w, r, value)
}

//
// update the idle speed
//
func (webserver *WebServer) postECUAdjustIdleSpeed(w http.ResponseWriter, r *http.Request) {
	var data ECUAdjustment

	log.Infof("rest-post update ecu idle speed")

	// get the body of our POST request
	// unmarshal this into a new Config struct
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)

	value := webserver.reader.ECU.AdjustIdleSpeed(data.Steps)
	webserver.updateAdjustableValue(w, r, value)
}

//
// update the ignition advance
//
func (webserver *WebServer) postECUAdjustIgnitionAdvance(w http.ResponseWriter, r *http.Request) {
	var data ECUAdjustment

	log.Infof("rest-post update ecu idle speed")

	// get the body of our POST request
	// unmarshal this into a new Config struct
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)

	value := webserver.reader.ECU.AdjustIgnitionAdvanceOffset(data.Steps)
	webserver.updateAdjustableValue(w, r, value)
}

//
// update the iac position
//
func (webserver *WebServer) postECUAdjustIAC(w http.ResponseWriter, r *http.Request) {
	var data ECUAdjustment

	log.Infof("rest-post update ecu idle speed")

	// get the body of our POST request
	// unmarshal this into a new Config struct
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)

	value := webserver.reader.ECU.AdjustIACPosition(data.Steps)
	webserver.updateAdjustableValue(w, r, value)
}

//
// update the adjustable value
//
func (webserver *WebServer) updateAdjustableValue(w http.ResponseWriter, r *http.Request, value int) {
	if webserver.isECUConnected(w) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		response := AdjustmentResponse{Value: value}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Warnf("rest-call response failed")
			// return a error code
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

//
// test the fuel pump
//
func (webserver *WebServer) postECUTestFuelPump(w http.ResponseWriter, r *http.Request) {
	var data ECUActivate

	log.Infof("rest-post test fuel pump")

	// get the body of our POST request
	// unmarshal this into a new Config struct
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)

	value := webserver.reader.ECU.TestFuelPump(data.Activate)
	webserver.updateTestActuator(w, r, value)
}

//
// test the ptc relay
//
func (webserver *WebServer) postECUTestPTC(w http.ResponseWriter, r *http.Request) {
	var data ECUActivate

	log.Infof("rest-post test PTC")

	// get the body of our POST request
	// unmarshal this into a new Config struct
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)

	value := webserver.reader.ECU.TestPTCRelay(data.Activate)
	webserver.updateTestActuator(w, r, value)
}

//
// test the aircon
//
func (webserver *WebServer) postECUTestAircon(w http.ResponseWriter, r *http.Request) {
	var data ECUActivate

	log.Infof("rest-post test aircon")

	// get the body of our POST request
	// unmarshal this into a new Config struct
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)

	value := webserver.reader.ECU.TestPTCRelay(data.Activate)
	webserver.updateTestActuator(w, r, value)
}

//
// test the purge valve
//
func (webserver *WebServer) postECUTestPurgeValve(w http.ResponseWriter, r *http.Request) {
	var data ECUActivate

	log.Infof("rest-post test purge valve")

	// get the body of our POST request
	// unmarshal this into a new Config struct
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)

	value := webserver.reader.ECU.TestPurgeValve(data.Activate)
	webserver.updateTestActuator(w, r, value)
}

//
// test the boost valve
//
func (webserver *WebServer) postECUTestBoostValve(w http.ResponseWriter, r *http.Request) {
	var data ECUActivate

	log.Infof("rest-post test boost valve")

	// get the body of our POST request
	// unmarshal this into a new Config struct
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)

	value := webserver.reader.ECU.TestBoostValve(data.Activate)
	webserver.updateTestActuator(w, r, value)
}

//
// test fan 1
//
func (webserver *WebServer) postECUTestFan1(w http.ResponseWriter, r *http.Request) {
	var data ECUActivate

	log.Infof("rest-post test fan 1")

	// get the body of our POST request
	// unmarshal this into a new Config struct
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)

	value := webserver.reader.ECU.TestFan1(data.Activate)
	webserver.updateTestActuator(w, r, value)
}

//
// test fan 2
//
func (webserver *WebServer) postECUTestFan2(w http.ResponseWriter, r *http.Request) {
	var data ECUActivate

	log.Infof("rest-post test fan 2")

	// get the body of our POST request
	// unmarshal this into a new Config struct
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)

	value := webserver.reader.ECU.TestFan2(data.Activate)
	webserver.updateTestActuator(w, r, value)
}

//
// test injectors
//
func (webserver *WebServer) postECUTestInjectors(w http.ResponseWriter, r *http.Request) {
	var data ECUActivate

	log.Infof("rest-post test injectors")

	// get the body of our POST request
	// unmarshal this into a new Config struct
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)

	value := webserver.reader.ECU.TestInjectors(data.Activate)
	webserver.updateTestActuator(w, r, value)
}

//
// test coil
//
func (webserver *WebServer) postECUTestCoil(w http.ResponseWriter, r *http.Request) {
	var data ECUActivate

	log.Infof("rest-post test coil")

	// get the body of our POST request
	// unmarshal this into a new Config struct
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)

	value := webserver.reader.ECU.TestCoil(data.Activate)
	webserver.updateTestActuator(w, r, value)
}

//
// update the actuator status
//
func (webserver *WebServer) updateTestActuator(w http.ResponseWriter, r *http.Request, value bool) {
	if webserver.isECUConnected(w) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		response := ActionResponse{Success: value}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Warnf("rest-call response failed")
			// return a error code
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

//
// checks if the ECU is connected and sets the headers accordingly
// create the response for methods that require the ecu to be connected.
//
func (webserver *WebServer) isECUConnected(w http.ResponseWriter) bool {
	if webserver.reader.ECU.Status.Connected {
		// return a 200 status code
		w.WriteHeader(http.StatusOK)
	} else {
		log.Infof("rest-call ecu is not connected")
		// return service unavailable if unable to connect
		w.WriteHeader(http.StatusServiceUnavailable)
		// put the ecu status in the body
		status := webserver.reader.ECU.Status

		if err := json.NewEncoder(w).Encode(status); err != nil {
			log.Warnf("rest-call response failed")
			// return a error code
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	return webserver.reader.ECU.Status.Connected
}