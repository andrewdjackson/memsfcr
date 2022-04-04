package fcr

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
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

type ECUActivateResponse struct {
	Actuator string `json:"actuator"`
	Activate bool   `json:"activate"`
}

type ActionResponse struct {
	Success bool `json:"success"`
}

type AdjustmentResponse struct {
	Adjustment string `json:"adjustment"`
	Value      int    `json:"value"`
}

const ActuatorFuelPump = "fuelpump"
const ActuatorPTC = "ptc"
const ActuatorAircon = "aircon"
const ActuatorPurgeValve = "purgevalve"
const ActuatorBoostValve = "boostvalve"
const ActuatorFan1 = "fan1"
const ActuatorFan2 = "fan2"
const ActuatorInjectors = "injectors"
const ActuatorCoil = "coil"

//
// Connection Status
// returns the status of the ecu connection along with the ecu id and the iac initial position
//
func (webserver *WebServer) getECUConnectionStatus(w http.ResponseWriter, r *http.Request) {
	log.Infof("rest-get read ecu status")

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("rest error closing response body")
		}
	}(r.Body)

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
	var connected bool
	var err error

	log.Infof("rest-post connect ecu")

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("rest error closing response body")
		}
	}(r.Body)

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

		if connected, err = webserver.reader.ECU.ConnectAndInitialiseECU(port.Port); err == nil {
			log.Infof("rest-post connected (%t) to the ecu", connected)
			// return a 200 status code
			w.WriteHeader(http.StatusOK)
		} else {
			log.Warnf("rest-post unable to connect to the ecu")
			// return service unavailable if unable to connect
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}

	if err = json.NewEncoder(w).Encode(webserver.reader.ECU.Status); err != nil {
		log.Warnf("rest-post ecu connect response failed")
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//
// Disconnect the ECU
//
func (webserver *WebServer) postECUDisconnect(w http.ResponseWriter, r *http.Request) {
	var err error

	log.Infof("rest-post disconnect ecu")

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("rest error closing response body")
		}
	}(r.Body)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if !webserver.reader.ECU.Status.Connected {
		// return status if already disconnected
		w.WriteHeader(http.StatusAlreadyReported)
	} else {
		// disconnect the ECU
		if err = webserver.reader.ECU.Disconnect(); err == nil {
			log.Infof("rest-post disconnected from the ecu")
			// return a 200 status code
			w.WriteHeader(http.StatusOK)
		} else {
			log.Warnf("rest-post unable to disconnect the ecu")
			// return service unavailable if unable to connect
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}

	if err := json.NewEncoder(w).Encode(webserver.reader.ECU.Status); err != nil {
		log.Warnf("rest-post disconnect ecu response failed")
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

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("rest error closing response body")
		}
	}(r.Body)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if webserver.isECUConnected(w) {
		if !webserver.waitingForECUResponse {
			// set the flag to prevent calls mid-protocol
			webserver.waitingForECUResponse = true
			// get the ECU data
			if memsdata, err := webserver.reader.ECU.GetDataframes(); err == nil {
				log.Infof("rest-get ecu dataframes (%+v)", memsdata)

				if err := json.NewEncoder(w).Encode(memsdata); err != nil {
					log.Warnf("rest-get read ecu dataframes response failed")
					// return a error code
					w.WriteHeader(http.StatusInternalServerError)
				}
			} else {
				log.Warnf("rest-get read ecu dataframes serial comms fault")
				// return a error code
				w.WriteHeader(http.StatusNotFound)
			}
			// clear the flag
			webserver.waitingForECUResponse = false
		} else {
			log.Warnf("rest-get already waiting for ECU")
			// return a error code
			w.WriteHeader(http.StatusTooManyRequests)
		}
	}
}

//
// Diagnostics
// returns the diagnostics
//
func (webserver *WebServer) getDiagnostics(w http.ResponseWriter, r *http.Request) {
	log.Infof("rest-get read diagnostics")

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("rest error closing response body")
		}
	}(r.Body)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// webserver.reader.ECU.Diagnostics.Analyse()
	diagnostics := webserver.reader.ECU.Diagnostics

	if err := json.NewEncoder(w).Encode(diagnostics); err != nil {
		log.Warnf("rest-post response failed (%+v)", err)
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
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

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warnf("rest error closing response body")
		}
	}(r.Body)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if webserver.isECUConnected(w) {
		if value, err := webserver.reader.ECU.GetIACPosition(); err == nil {
			response := AdjustmentResponse{Adjustment: "iac", Value: value}

			log.Infof("rest-get ecu iac position (%v)", value)

			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Warnf("rest-get iac position response failed")
				// return a error code
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			log.Warnf("rest-get iac position response failed")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

//
//  send heartbeat the ecu
//
func (webserver *WebServer) postECUHeartbeat(w http.ResponseWriter, r *http.Request) {
	value := false

	log.Infof("rest-post send heartbeat")
	if err := webserver.reader.ECU.SendHeartbeat(); err == nil {
		value = true
	}

	webserver.updateECUState(w, r, value)
}

//
//  reset the ecu
//
func (webserver *WebServer) postECUReset(w http.ResponseWriter, r *http.Request) {
	value := false

	log.Infof("rest-post reset ecu")
	if err := webserver.reader.ECU.ResetECU(); err == nil {
		value = true
	}
	webserver.updateECUState(w, r, value)
}

//
// clear the fault codes
//
func (webserver *WebServer) postECUClearFaults(w http.ResponseWriter, r *http.Request) {
	value := false
	log.Infof("rest-post clear ecu faults")
	if err := webserver.reader.ECU.ClearFaults(); err == nil {
		value = true
	}
	webserver.updateECUState(w, r, value)
}

//
// clear the adjustable values
//
func (webserver *WebServer) postECUClearAdjustments(w http.ResponseWriter, r *http.Request) {
	value := false
	log.Infof("rest-post clear ecu adjustable values")
	if err := webserver.reader.ECU.ResetAdjustments(); err == nil {
		value = true
	}
	webserver.updateECUState(w, r, value)
}

//
// update ecu state (clear faults, reset and heartbeat)
//
func (webserver *WebServer) updateECUState(w http.ResponseWriter, r *http.Request, value bool) {
	if webserver.isECUConnected(w) {

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Warnf("rest error closing response body")
			}
		}(r.Body)

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
// update the short term fuel trim
//
func (webserver *WebServer) postECUAdjustSTFT(w http.ResponseWriter, r *http.Request) {
	var data ECUAdjustment

	log.Infof("rest-post update ecu stft")

	// get the body of our POST request
	// unmarshal this into a new Config struct
	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &data)
	if err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	value, _ := webserver.reader.ECU.AdjustShortTermFuelTrim(data.Steps)
	adjustment := AdjustmentResponse{Adjustment: "stft", Value: value}
	webserver.updateAdjustableValue(w, r, adjustment)
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

	err := json.Unmarshal(body, &data)
	if err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	value, _ := webserver.reader.ECU.AdjustLongTermFuelTrim(data.Steps)
	adjustment := AdjustmentResponse{Adjustment: "ltft", Value: value}
	webserver.updateAdjustableValue(w, r, adjustment)
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

	err := json.Unmarshal(body, &data)
	if err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}

	value, _ := webserver.reader.ECU.AdjustIdleDecay(data.Steps)
	adjustment := AdjustmentResponse{Adjustment: "idledecay", Value: value}

	webserver.updateAdjustableValue(w, r, adjustment)
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

	err := json.Unmarshal(body, &data)
	if err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}

	value, _ := webserver.reader.ECU.AdjustIdleSpeed(data.Steps)
	adjustment := AdjustmentResponse{Adjustment: "idlespeed", Value: value}

	webserver.updateAdjustableValue(w, r, adjustment)
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

	err := json.Unmarshal(body, &data)
	if err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}

	value, _ := webserver.reader.ECU.AdjustIgnitionAdvanceOffset(data.Steps)
	adjustment := AdjustmentResponse{Adjustment: "ignitionadvance", Value: value}

	webserver.updateAdjustableValue(w, r, adjustment)
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

	err := json.Unmarshal(body, &data)
	if err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}

	value, _ := webserver.reader.ECU.AdjustIACPosition(data.Steps)
	adjustment := AdjustmentResponse{Adjustment: "iac", Value: value}

	webserver.updateAdjustableValue(w, r, adjustment)
}

//
// update the adjustable value
//
func (webserver *WebServer) updateAdjustableValue(w http.ResponseWriter, r *http.Request, adjustment AdjustmentResponse) {
	if webserver.isECUConnected(w) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Warnf("rest error closing response body")
			}
		}(r.Body)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		log.Infof("rest-post adjustable value response")
		response := AdjustmentResponse{Adjustment: adjustment.Adjustment, Value: adjustment.Value}

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

	err := json.Unmarshal(body, &data)
	if err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := webserver.reader.ECU.TestFuelPump(data.Activate); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	actuatorResponse := ECUActivateResponse{Actuator: ActuatorFuelPump, Activate: data.Activate}

	webserver.updateTestActuator(w, r, actuatorResponse)
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

	err := json.Unmarshal(body, &data)
	if err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := webserver.reader.ECU.TestPTCRelay(data.Activate); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	actuatorResponse := ECUActivateResponse{Actuator: ActuatorPTC, Activate: data.Activate}

	webserver.updateTestActuator(w, r, actuatorResponse)
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

	err := json.Unmarshal(body, &data)
	if err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := webserver.reader.ECU.TestACRelay(data.Activate); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	actuatorResponse := ECUActivateResponse{Actuator: ActuatorAircon, Activate: data.Activate}

	webserver.updateTestActuator(w, r, actuatorResponse)
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

	err := json.Unmarshal(body, &data)
	if err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := webserver.reader.ECU.TestPurgeValve(data.Activate); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	actuatorResponse := ECUActivateResponse{Actuator: ActuatorPurgeValve, Activate: data.Activate}

	webserver.updateTestActuator(w, r, actuatorResponse)
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

	err := json.Unmarshal(body, &data)
	if err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := webserver.reader.ECU.TestBoostValve(data.Activate); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	actuatorResponse := ECUActivateResponse{Actuator: ActuatorBoostValve, Activate: data.Activate}

	webserver.updateTestActuator(w, r, actuatorResponse)
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

	err := json.Unmarshal(body, &data)
	if err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := webserver.reader.ECU.TestFan1(data.Activate); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	actuatorResponse := ECUActivateResponse{Actuator: ActuatorFan1, Activate: data.Activate}

	webserver.updateTestActuator(w, r, actuatorResponse)
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

	err := json.Unmarshal(body, &data)
	if err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := webserver.reader.ECU.TestFan2(data.Activate); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	actuatorResponse := ECUActivateResponse{Actuator: ActuatorFan2, Activate: data.Activate}

	webserver.updateTestActuator(w, r, actuatorResponse)
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

	err := json.Unmarshal(body, &data)
	if err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := webserver.reader.ECU.TestInjectors(data.Activate); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	actuatorResponse := ECUActivateResponse{Actuator: ActuatorInjectors, Activate: data.Activate}

	webserver.updateTestActuator(w, r, actuatorResponse)
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

	err := json.Unmarshal(body, &data)
	if err != nil {
		// return a error code
		w.WriteHeader(http.StatusInternalServerError)
	}

	if err := webserver.reader.ECU.TestCoil(data.Activate); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	actuatorResponse := ECUActivateResponse{Actuator: ActuatorCoil, Activate: data.Activate}

	webserver.updateTestActuator(w, r, actuatorResponse)
}

//
// update the actuator status
//
func (webserver *WebServer) updateTestActuator(w http.ResponseWriter, r *http.Request, actuatorResponse ECUActivateResponse) {
	if webserver.isECUConnected(w) {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				log.Warnf("rest error closing response body")
			}
		}(r.Body)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		if err := json.NewEncoder(w).Encode(actuatorResponse); err != nil {
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
