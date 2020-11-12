package utils

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"gopkg.in/ini.v1"
)

// ReadmemsConfig readmems configuration
type ReadmemsConfig struct {
	// Config
	Port      string
	LogToFile string
	LogFolder string
	Loop      string
	Ports     []string
	Debug     string
	Frequency string
	Headless  string
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
	config.Headless = "false"
	config.Frequency = "500"

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

// WriteConfig write the config file
func WriteConfig(c *ReadmemsConfig) {
	// create the folders if they don't exist
	createDataFolders(homeFolder)

	filename := fmt.Sprintf("%s/memsfcr.cfg", appFolder)

	// create the file if it doesn't exist
	_, _ = os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)

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
	cfg.Section("").Key("headless").SetValue(c.Headless)

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
	c.Headless = cfg.Section("").Key("headless").String()

	LogI.Println("MemsFCR Config", c)
	return c
}
