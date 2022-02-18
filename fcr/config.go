package fcr

import (
	"fmt"
	"github.com/andrewdjackson/rosco"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

// Config readmems configuration
type Config struct {
	Port       string
	Ports      []string
	Debug      string
	Frequency  string
	Version    string
	Build      string
	ServerPort string
}

var config Config

// NewConfig creates a new instance of readmems config
func NewConfig() *Config {
	config.Port = "/dev/tty.serial"
	config.Debug = "false"
	config.Frequency = "500"
	config.Version = "0.0.0"
	config.ServerPort = "0"

	currentTime := time.Now()
	config.Build = currentTime.Format("2006-01-02")

	return &config
}

// WriteConfig write the config file
func WriteConfig(c *Config) {
	filename := fmt.Sprintf("%s/memsfcr.cfg", rosco.GetHomeFolder())

	// create the file if it doesn't exist
	_, _ = os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)

	cfg, err := ini.LooseLoad(filename)
	if err != nil {
		log.Infof("failed to read file: %v", err)
	}

	cfg.Section("").Key("version").SetValue(c.Version)
	cfg.Section("").Key("build").SetValue(c.Build)
	cfg.Section("").Key("port").SetValue(c.Port)
	cfg.Section("").Key("debug").SetValue(c.Debug)
	cfg.Section("").Key("frequency").SetValue(c.Frequency)
	cfg.Section("").Key("serverport").SetValue(c.ServerPort)

	err = cfg.SaveTo(filename)

	if err != nil {
		log.Infof("failed to write file: %v", err)
	}

	log.Infof("updated config: %s", filename)
}

// ReadConfig reads the config file
func ReadConfig() *Config {
	filename := fmt.Sprintf("%s/memsfcr.cfg", rosco.GetHomeFolder())
	log.Infof("loading config from %s", filename)

	c := NewConfig()

	cfg, err := ini.Load(filename)
	if err != nil {
		log.Infof("failed to read file: %v", err)
		// couldn't read the config so write a new file
		WriteConfig(c)
		// return the default config
		return c
	}

	c.Port = cfg.Section("").Key("port").String()
	c.Debug = cfg.Section("").Key("debug").String()
	c.Frequency = cfg.Section("").Key("frequency").String()
	c.ServerPort = cfg.Section("").Key("serverport").String()

	log.Infof("MemsFCR Config %+v", c)
	return c
}

func CreateFolders() {
	err := createFolder(rosco.GetHomeFolder())
	if err == nil {
		_ = createFolder(rosco.GetDebugFolder())
		_ = createFolder(rosco.GetLogFolder())
		_ = createFolder(rosco.GetAppFolder())
	}
}

func createFolder(path string) error {
	info, err := os.Stat(path)

	if err != nil {
		log.Warnf("unable to find folder %s (%s)", path, err)
	} else {
		if info.IsDir() {
			log.Infof("found folder %s", path)
		}
	}

	if os.IsNotExist(err) {
		log.Errorf("folder %s does not exist, creating folder", path)

		err := os.MkdirAll(path, 0755)
		if err != nil {
			log.Errorf("unable to create folder %s (%s)", path, err)
		}
	}

	return err
}
