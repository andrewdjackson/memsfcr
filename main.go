package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/andrewdjackson/memsfcr/fcr"
	"github.com/mitchellh/go-homedir"
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
	// User Home Folder
	HomeFolder string
	// Application Binary Folder
	AppFolder string
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
		}
	}

	currentTime := time.Now()
	Build = currentTime.Format("2006-01-02")

	// get the users home directory
	dir, _ := homedir.Dir()
	HomeFolder = filepath.FromSlash(dir)

	// get the application binary current directory
	dir, _ = os.Getwd()
	AppFolder = filepath.FromSlash(dir)

	// intialise the logging
	setupLogging()
}

func setupLogging() {
	// create a log file using the current date and time
	// this saves trying to roll logs
	currentTime := time.Now()
	dateTime := currentTime.Format("2006-01-02 15:04:05")
	dateTime = strings.ReplaceAll(dateTime, ":", "")
	dateTime = strings.ReplaceAll(dateTime, " ", "-")
	filename := fmt.Sprintf("logs/debug-%s.log", dateTime)
	filename = filepath.FromSlash(filename)

	// write logs to file and console
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Warn("error opening log file")
	}

	multilogwriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multilogwriter)

	log.SetFormatter(&log.TextFormatter{
		ForceColors: false,
		DisableColors: false,
		FullTimestamp: true,
		TimestampFormat: "15:04:05.000",
	})

	// disable function logging
	log.SetReportCaller(false)
}

func main() {
	log.Infof("MemsFCR Version %s, Build %s", Version, Build)
	log.Infof("MemsFCR Home Folder %s", HomeFolder)
	log.Infof("MemsFCR App Folder %s", AppFolder)

	// set up and initialise the fault code reader
	reader := fcr.NewMemsReader()
	// start the web server
	reader.StartWebServer()
	// open the browser view
	// reader.OpenBrowser()

	// don't exit
	for { select{ } }
}
