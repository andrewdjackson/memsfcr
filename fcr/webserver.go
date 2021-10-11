package fcr

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type RelativePaths struct {
	Webroot string
	ExePath string
}

// WebServer the web interface
type WebServer struct {
	// multiplex router interface
	router *mux.Router
	// websocket interface
	httpDir  string
	paths    RelativePaths
	ws       *websocket.Conn
	upgrader websocket.Upgrader
	// HTTPPort used by the HTTP Server instance
	HTTPPort int
	// ServerRunning indicates where the server is active
	ServerRunning bool
	// Pointer to Mems Fault Code Reader
	reader *MemsReader
	// waiting for a response from the ECU
	waitingForECUResponse bool
}

const (
	indexTemplate = "index.template.html"
	indexData     = "index.template.json"
)

// NewWebInterface creates a new web interface
func NewWebServer(reader *MemsReader) *WebServer {
	webserver := &WebServer{}
	webserver.HTTPPort = 0
	webserver.httpDir = ""
	webserver.ServerRunning = false
	webserver.reader = reader
	webserver.paths = RelativePaths{}

	webserver.upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return webserver
}

func (webserver *WebServer) getRelativePaths() RelativePaths {
	paths := RelativePaths{}

	// determine the path to find the local html files
	// based on the current executable path
	dir, _ := os.Getwd()
	exepath := filepath.FromSlash(dir)
	paths.ExePath, _ = filepath.Abs(exepath)

	// use default browser on Windows until I can get the Webview to work
	if runtime.GOOS == "darwin" {
		// get the executable path on MacOS
		exepath, _ = os.Executable()
		paths.ExePath, _ = filepath.Abs(filepath.Dir(exepath))

		// MacOS use .app Resources
		if strings.Contains(paths.ExePath, "MacOS") {
			// packaged app
			paths.Webroot = strings.Replace(paths.ExePath, "MacOS", "Resources", -1)
		} else {
			// running a local or dev version
			paths.Webroot = fmt.Sprintf("%s/Resources", paths.ExePath)
		}
	} else if runtime.GOOS == "linux" {
		// linux path
		// get the executable path
		paths.Webroot = fmt.Sprintf("%s/resources", paths.ExePath)
	} else {
		// windows use the exe subdirectory
		paths.Webroot = fmt.Sprintf("%s\\resources", paths.ExePath)
	}

	paths.Webroot = filepath.ToSlash(paths.Webroot)

	log.Infof("path to the local html files (%s)", paths.Webroot)

	return paths
}

func (webserver *WebServer) newRouter() *mux.Router {
	webserver.paths = webserver.getRelativePaths()

	webserver.httpDir = webserver.paths.Webroot

	log.Infof("path to the local html files (%s) on (%s)", webserver.httpDir, runtime.GOOS)

	// set a router and a handler to accept messages over the websocket

	r := mux.NewRouter()
	//r.HandleFunc("/ws", webserver.websocketHandler)
	r.HandleFunc("/heartbeat", webserver.browserHeartbeatHandler)

	r.HandleFunc("/config", webserver.getConfigHandler).Methods(http.MethodGet)
	r.HandleFunc("/config/ports", webserver.getSerialPortsHandler).Methods(http.MethodGet)
	r.HandleFunc("/config", webserver.updateConfigHandler).Methods(http.MethodPut)

	r.HandleFunc("/scenario", webserver.getListofScenarios).Methods(http.MethodGet)
	r.HandleFunc("/scenario/play/{scenarioId}", webserver.getScenarioDetails).Methods(http.MethodGet)
	r.HandleFunc("/scenario/details/{scenarioId}", webserver.getPlaybackDetails).Methods(http.MethodGet)
	//r.HandleFunc("/scenario/seek", webserver.getPlaybackDetails).Methods(http.MethodPatch)

	r.HandleFunc("/rosco", webserver.getECUConnectionStatus).Methods(http.MethodGet)
	r.HandleFunc("/rosco/connect", webserver.postECUConnect).Methods(http.MethodPost)
	r.HandleFunc("/rosco/disconnect", webserver.postECUDisconnect).Methods(http.MethodPost)
	r.HandleFunc("/rosco/dataframe", webserver.getECUDataframes).Methods(http.MethodGet)
	r.HandleFunc("/rosco/heartbeat", webserver.postECUHeartbeat).Methods(http.MethodPost)
	r.HandleFunc("/rosco/iac", webserver.getECUIAC).Methods("GET")
	r.HandleFunc("/rosco/diagnostics", webserver.getDiagnostics).Methods(http.MethodGet)

	r.HandleFunc("/rosco/reset", webserver.postECUReset).Methods(http.MethodPost)
	r.HandleFunc("/rosco/reset/ecu", webserver.postECUReset).Methods(http.MethodPost)
	r.HandleFunc("/rosco/reset/faults", webserver.postECUClearFaults).Methods(http.MethodPost)
	r.HandleFunc("/rosco/reset/adjustments", webserver.postECUClearAdjustments).Methods(http.MethodPost)

	r.HandleFunc("/rosco/adjust/stft", webserver.postECUAdjustSTFT).Methods(http.MethodPost)
	r.HandleFunc("/rosco/adjust/ltft", webserver.postECUAdjustLTFT).Methods(http.MethodPost)
	r.HandleFunc("/rosco/adjust/idledecay", webserver.postECUAdjustIdleDecay).Methods(http.MethodPost)
	r.HandleFunc("/rosco/adjust/idlespeed", webserver.postECUAdjustIdleSpeed).Methods(http.MethodPost)
	r.HandleFunc("/rosco/adjust/ignitionadvance", webserver.postECUAdjustIgnitionAdvance).Methods(http.MethodPost)
	r.HandleFunc("/rosco/adjust/iac", webserver.postECUAdjustIAC).Methods(http.MethodPost)

	r.HandleFunc("/rosco/test/fuelpump", webserver.postECUTestFuelPump).Methods(http.MethodPost)
	r.HandleFunc("/rosco/test/ptc", webserver.postECUTestPTC).Methods(http.MethodPost)
	r.HandleFunc("/rosco/test/aircon", webserver.postECUTestAircon).Methods(http.MethodPost)
	r.HandleFunc("/rosco/test/purgevalve", webserver.postECUTestPurgeValve).Methods(http.MethodPost)
	r.HandleFunc("/rosco/test/boostvalve", webserver.postECUTestBoostValve).Methods(http.MethodPost)
	r.HandleFunc("/rosco/test/fan", webserver.postECUTestFan1).Methods(http.MethodPost)
	r.HandleFunc("/rosco/test/fan/1", webserver.postECUTestFan1).Methods(http.MethodPost)
	r.HandleFunc("/rosco/test/fan/2", webserver.postECUTestFan2).Methods(http.MethodPost)
	r.HandleFunc("/rosco/test/injectors", webserver.postECUTestInjectors).Methods(http.MethodPost)
	r.HandleFunc("/rosco/test/coil", webserver.postECUTestCoil).Methods(http.MethodPost)

	r.HandleFunc("/", webserver.renderIndex)

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.FileServer(http.Dir(webserver.httpDir))
	//r.Handle("/", fileServer)
	r.PathPrefix("/").Handler(fileServer).Methods("GET")

	return r
}

func (webserver *WebServer) renderIndex(w http.ResponseWriter, r *http.Request) {
	log.Infof("rendering html template")

	templateFile := fmt.Sprintf("%s/%s", webserver.paths.Webroot, indexTemplate)
	templateFile = filepath.ToSlash(templateFile)

	dataFile := fmt.Sprintf("%s/%s", webserver.paths.Webroot, indexData)
	dataFile = filepath.ToSlash(dataFile)

	page, err := template.ParseFiles(templateFile)

	if err != nil {
		log.Errorf("template error: ", err)
	}

	data := map[string]interface{}{}
	jsondata, err := ioutil.ReadFile(dataFile)

	if err != nil {
		log.Errorf("template error: ", err)
	}

	if err := json.Unmarshal(jsondata, &data); err != nil {
		log.Errorf("template error: ", err)
	}

	/*
		// blocks until version has been found and causes an exception if the network is
		// offline.
		// currently commented out until I get a chance to fix this
		//
		data["Version"] = webserver.reader.Config.Version
		latestVersion := webserver.newVersionAvailable()
		if latestVersion != webserver.reader.Config.Version {
			data["Version"] = "New Version! Click to Download"
			data["NewVersion"] = true
			log.Infof("%s", data["Version"])
		}
	*/

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Etag", webserver.reader.Config.Build)

	err = page.Execute(w, data)

	if err != nil {
		log.Errorf("\nRender Error: %v\n", err)
		return
	}
}

// RunHTTPServer run the server
func (webserver *WebServer) RunHTTPServer() {
	// Declare a new router
	webserver.router = webserver.newRouter()

	// We can then pass our router (after declaring all our routes) to this method
	// (where previously, we were leaving the second argument as nil)
	listener, err := net.Listen("tcp", ":8081")

	if err != nil {
		log.Errorf("error starting web interface (%s)", err)
	}

	webserver.HTTPPort = listener.Addr().(*net.TCPAddr).Port

	log.Infof("started http server on port %d", webserver.HTTPPort)
	webserver.ServerRunning = true

	err = http.Serve(listener, webserver.router)

	if err != nil {
		log.Errorf("error starting web interface (%s)", err)
	}
}

func (webserver *WebServer) newVersionAvailable() string {
	versionUrl := "https://raw.githubusercontent.com/andrewdjackson/memsfcr/master/version"
	latestVersion := webserver.reader.Config.Version
	response, err := http.Get(versionUrl)
	defer response.Body.Close()

	if err == nil {
		var lines []string
		scanner := bufio.NewScanner(response.Body)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		log.Infof("repo version %s", lines[0])
		latestVersion = lines[0]
	} else {
		log.Warnf("version check failed %s", err)
	}

	return latestVersion
}
