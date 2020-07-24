package scenarios

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// GetScenarioPath returns the path to the scenario files
func GetScenarioPath(file string) string {
	var filename string

	homeFolder, _ := homedir.Dir()
	appFolder := fmt.Sprintf("%s/memsfcr", homeFolder)

	if file == "" {
		filename = fmt.Sprintf("%s/logs", appFolder)
	} else {
		filename = fmt.Sprintf("%s/logs/%s", appFolder, file)
	}

	return filepath.FromSlash(filename)
}

// GetScenarios reads the directory and returns
// a list of scenario entries sorted by filename.
func GetScenarios() ([]string, error) {
	logFolder := GetScenarioPath("")

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

// GetScenario returns the data for the given scenario
func GetScenario(id string) ScenarioDescription {
	file := GetScenarioPath(id)
	r := NewResponder()
	r.LoadScenario(file)

	scenario := ScenarioDescription{}
	scenario.Count = r.playbook.count
	scenario.Position = r.playbook.position
	scenario.Name = id

	return scenario
}
