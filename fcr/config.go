package fcr

import (
	"fmt"
	"os"
	"time"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

// Config readmems configuration
type Config struct {
	// Config
	Port      string
	LogToFile string
	LogFolder string
	Loop      string
	Ports     []string
	Debug     string
	Frequency string
	Headless  string
	Version   string
	Build     string
}

var config Config
var homeFolder string
var appFolder string
var logFolder string

// NewConfig creates a new instance of readmems config
func NewConfig() *Config {
	config.Port = "/dev/tty.serial"
	config.LogFolder = ""
	config.LogToFile = "true"
	config.Loop = "100000000"
	config.Debug = "false"
	config.Headless = "false"
	config.Frequency = "500"
	config.Version = "0.0.0"

	currentTime := time.Now()
	config.Build = currentTime.Format("2006-01-02")

	return &config
}

func getHomeFolder() string {
	homeFolder, _ = homedir.Dir()
	return homeFolder
}

func createFolder(path string) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(path, 0755)
		if errDir != nil {
			log.Errorf("unable to create folder %s (%s)", path, err)
		}
	}
}

func createDataFolders(folder string) {
	// sandbox folder
	//homeFolder = "./Documents"

	appFolder = fmt.Sprintf("%s/memsfcr", folder)
	createFolder(appFolder)

	logFolder = fmt.Sprintf("%s/logs", folder)
	createFolder(logFolder)
}

// WriteConfig write the config file
func WriteConfig(c *Config) {
	homeFolder = getHomeFolder()
	// create the folders if they don't exist
	createDataFolders(homeFolder)

	filename := fmt.Sprintf("%s/memsfcr.cfg", appFolder)

	// create the file if it doesn't exist
	_, _ = os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)

	cfg, err := ini.LooseLoad(filename)
	if err != nil {
		log.Infof("failed to read file: %v", err)
	}

	cfg.Section("").Key("version").SetValue(c.Version)
	cfg.Section("").Key("build").SetValue(c.Build)
	cfg.Section("").Key("port").SetValue(c.Port)
	cfg.Section("").Key("loop").SetValue(c.Loop)
	cfg.Section("").Key("logtofile").SetValue(c.LogToFile)
	cfg.Section("").Key("logfolder").SetValue(c.LogFolder)
	cfg.Section("").Key("debug").SetValue(c.Debug)
	cfg.Section("").Key("frequency").SetValue(c.Frequency)
	cfg.Section("").Key("headless").SetValue(c.Headless)

	err = cfg.SaveTo(filename)

	if err != nil {
		log.Infof("failed to write file: %v", err)
	}

	log.Infof("updated config: %s", filename)
}

// ReadConfig reads the config file
func ReadConfig() *Config {
	homeFolder = getHomeFolder()
	// create the folders if they don't exist
	createDataFolders(homeFolder)

	filename := fmt.Sprintf("%s/memsfcr.cfg", appFolder)

	c := NewConfig()
	c.LogFolder = logFolder

	cfg, err := ini.Load(filename)
	if err != nil {
		log.Infof("failed to read file: %v", err)
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

	log.Infof("MemsFCR Config %+v", c)
	return c
}
