package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/andrewdjackson/memsulator/utils"
	"github.com/gocarina/gocsv"
	"github.com/mitchellh/go-homedir"
)

// RawData represents the raw data from the log file
type RawData struct {
	Dataframe7d string `csv:"0x7d_raw"`
	Dataframe80 string `csv:"0x80_raw"`
}

// Scenario represents the scenario data
type Scenario struct {
	file *os.File
	// Rawdata from log
	Rawdata []*RawData
	// Position in the log
	Position int
	// Count of items in the log
	Count int
}

// NewScenario creates a new scenario
func NewScenario() *Scenario {
	scenario := &Scenario{}
	// initialise the log
	scenario.Rawdata = []*RawData{}
	// start at the beginning
	scenario.Position = 0
	// no items in the log
	scenario.Count = 0

	return scenario
}

// Open the CSV scenario file
func (scenario *Scenario) openFile(filepath string) {
	var err error

	scenario.file, err = os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)

	if err != nil {
		LogE.Printf("unable to open %s", err)
	}
}

// Load the scenario
func (scenario *Scenario) Load(filepath string) {
	scenario.openFile(filepath)

	if err := gocsv.Unmarshal(scenario.file, &scenario.Rawdata); err != nil {
		utils.LogE.Printf("unable to parse file %s", err)
	} else {
		scenario.Count = len(scenario.Rawdata)
		utils.LogI.Printf("loaded scenario %s (%d dataframes)", filepath, scenario.Count)
	}
}

// Next provides the next item in the log
func (scenario *Scenario) Next() *RawData {
	item := scenario.Rawdata[scenario.Position]
	scenario.Position = scenario.Position + 1

	// if we pass the end, loop back to the start
	if scenario.Position >= scenario.Count {
		utils.LogW.Printf("reached end of scenario, restarting from beginning")
		scenario.Position = 0
	}

	return item
}

// GetScenarios reads the directory and returns
// a list of scenario entries sorted by filename.
func GetScenarios() ([]string, error) {
	homeFolder, _ := homedir.Dir()
	appFolder := fmt.Sprintf("%s/memsfcr", homeFolder)
	logFolder := fmt.Sprintf("%s/logs", appFolder)

	var files []string
	fileInfo, err := ioutil.ReadDir(logFolder)

	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if strings.HasSuffix(file.Name(), ".csv") {
			files = append(files, file.Name())
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(files)))

	return files, nil
}
