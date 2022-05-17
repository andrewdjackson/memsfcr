package main

import (
	"flag"
	"fmt"
	"github.com/andrewdjackson/rosco"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/andrewdjackson/memsfcr/fcr"
	log "github.com/sirupsen/logrus"
)

//
// This file is the main function within the main package. It sets up the logging
// and extracts the application information (e.g version, location etc.)
// The fault code reader functionality is managed in the memsreader.go within
// the main package
//

var (
	// Version of the application
	Version string
	// Build date
	Build string
)

func init() {
	// if the version is not written into the binary
	// then read the version from the version file and set the build date to Now
	if strings.Compare(Version, "") == 0 {
		version, err := ioutil.ReadFile("version")

		if err != nil {
			Version = "0.0.0"
		} else {
			Version = string(version)
			Version = strings.TrimSuffix(Version, "\n")
		}
	}

	currentTime := time.Now()
	Build = currentTime.Format("2006-01-02")
}

func setupLogging(debug bool) {
	if debug {
		// create a log file using the current date and time
		// this saves trying to roll logs
		currentTime := time.Now()
		dateTime := currentTime.Format("2006-01-02 15:04:05")
		dateTime = strings.ReplaceAll(dateTime, ":", "")
		dateTime = strings.ReplaceAll(dateTime, " ", "-")
		filename := fmt.Sprintf("%s/debug-%s.log", rosco.GetDebugFolder(), dateTime)
		filename = filepath.FromSlash(filename)

		// write logs to file and console
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Warn("error opening log file")
		}

		log.SetFormatter(&log.TextFormatter{
			DisableColors:   true,
			FullTimestamp:   true,
			TimestampFormat: "15:04:05.000",
		})

		multilogwriter := io.MultiWriter(os.Stdout, f)
		log.SetOutput(multilogwriter)
		log.Infof("debug logging to %s", filename)
	} else {
		log.SetOutput(os.Stdout)

		log.SetFormatter(&log.TextFormatter{
			ForceColors:     false,
			DisableColors:   false,
			FullTimestamp:   true,
			TimestampFormat: "15:04:05.000",
		})
	}

	// disable function logging
	log.SetReportCaller(false)
}

func main() {
	var debug bool
	var headless bool

	flag.BoolVar(&debug, "debug", true, "output to a debug file")
	flag.BoolVar(&headless, "headless", false, "headless server mode")
	flag.Parse()

	// initialise the logging
	fcr.CreateFolders()
	setupLogging(debug)

	log.Infof("MemsFCR Version %s, Build %s", Version, Build)
	log.Infof("MemsFCR Home Folder %s", rosco.GetHomeFolder())
	log.Infof("MemsFCR App Folder %s", rosco.GetAppFolder())
	log.Infof("MemsFCR Log Folder %s", rosco.GetLogFolder())
	log.Infof("MemsFCR Debug Folder %s", rosco.GetDebugFolder())

	// create a channel to notify app to exit
	exit := make(chan int)

	// set up and initialise the fault code reader
	reader := fcr.NewMemsReader(Version, Build, headless)
	// start the web server
	reader.StartWebServer()

	if !headless {
		// open the browser view
		reader.OpenBrowser()
	} else {
		log.Infof("MemsFCR started in headless mode")
	}

	// wait for exit on the channel
	for {
		<-exit
	}
}
