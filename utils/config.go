package utils

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/ini.v1"
)

// ReadmemsConfig readmems configuration
type ReadmemsConfig struct {
	// Config
	Port       string
	LogToFile  string
	LogFolder  string
	Loop       string
	Ports      []string
	Debug      string
	Frequency  string
	UseBrowser string
}

var config ReadmemsConfig
var homeFolder string
var appFolder string
var logFolder string

// NewConfig creates a new instance of readmems config
func NewConfig() *ReadmemsConfig {
	config.Port = "/dev/tty.serial"
	config.LogFolder = ""
	config.LogToFile = "true"
	config.Loop = "100000000"
	config.Debug = "false"
	config.UseBrowser = "true"
	config.Frequency = "950"

	return &config
}

func createFolder(path string) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			LogE.Panicf("unable to create folder %s (%s)", path, err)
		}
	}
}

func createDataFolders(home string) {
	homeFolder, _ = homedir.Dir()

	appFolder = fmt.Sprintf("%s/memsfcr", homeFolder)
	createFolder(appFolder)

	logFolder = fmt.Sprintf("%s/logs", appFolder)
	createFolder(logFolder)
}

// reads a whole file into memory and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// WriteConfig write the config file
func WriteConfig(c *ReadmemsConfig) {
	// create the folders if they don't exist
	createDataFolders(homeFolder)

	filename := fmt.Sprintf("%s/memsfcr.cfg", appFolder)

	os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)

	cfg, err := ini.LooseLoad(filename)
	if err != nil {
		LogI.Printf("failed to read file: %v", err)
	}

	cfg.Section("").Key("port").SetValue(c.Port)
	cfg.Section("").Key("loop").SetValue(c.Loop)
	cfg.Section("").Key("logtofile").SetValue(c.LogToFile)
	cfg.Section("").Key("logfolder").SetValue(c.LogFolder)
	cfg.Section("").Key("debug").SetValue(c.Debug)
	cfg.Section("").Key("frequency").SetValue(c.Frequency)
	cfg.Section("").Key("usebrowser").SetValue(c.UseBrowser)

	err = cfg.SaveTo(filename)

	if err != nil {
		LogI.Printf("failed to write file: %v", err)
	}

	LogI.Printf("updated config: %s", filename)
}

// ReadConfig readsthe config file
func ReadConfig() *ReadmemsConfig {
	// create the folders if they don't exist
	createDataFolders(homeFolder)

	filename := fmt.Sprintf("%s/memsfcr.cfg", appFolder)

	c := NewConfig()
	c.LogFolder = logFolder

	cfg, err := ini.Load(filename)
	if err != nil {
		LogI.Printf("failed to read file: %v", err)
		// couldn't read the config so write a new file
		WriteConfig(c)
		// return the default config
		return c
	}

	c.Port = cfg.Section("").Key("port").String()
	c.Loop = cfg.Section("").Key("loop").String()
	c.LogToFile = cfg.Section("").Key("logtofile").String()
	c.LogFolder = cfg.Section("").Key("logfolder").String()
	c.Debug = cfg.Section("").Key("debug").String()
	c.Frequency = cfg.Section("").Key("frequency").String()
	c.UseBrowser = cfg.Section("").Key("usebrowser").String()

	LogI.Println("MemsFCR Config", c)
	return c
}
