package rosco

import (
	"bufio"
	"fmt"
	"os"

	"github.com/andrewdjackson/readmems/utils"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/ini.v1"
)

// ReadmemsConfig readmems configuration
type ReadmemsConfig struct {
	// Config
	Port       string
	Command    string
	Output     string
	LogFolder  string
	Loop       string
	Connection string
	Ports      []string
}

var config ReadmemsConfig

// NewConfig creates a new instance of readmems config
func NewConfig() *ReadmemsConfig {
	config.Port = "ttycodereader"
	config.Command = "read"
	config.Loop = "inf"
	config.Output = "stdout"
	config.Connection = "wait"
	config.LogFolder = "./logs"

	return &config
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
	folder, _ := homedir.Dir()
	filename := fmt.Sprintf("%s/.memsfcr.cfg", folder)

	os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)

	cfg, err := ini.LooseLoad(filename)
	if err != nil {
		utils.LogI.Printf("failed to read file: %v", err)
	}

	cfg.Section("").Key("port").SetValue(c.Port)
	cfg.Section("").Key("command").SetValue(c.Command)
	cfg.Section("").Key("loop").SetValue(c.Loop)
	cfg.Section("").Key("output").SetValue(c.Output)
	cfg.Section("").Key("logfolder").SetValue(c.LogFolder)
	cfg.Section("").Key("connection").SetValue(c.Connection)

	err = cfg.SaveTo(filename)

	if err != nil {
		utils.LogI.Printf("failed to write file: %v", err)
	}

	utils.LogI.Printf("updated config: %s", filename)
}

// ReadConfig readsthe config file
func ReadConfig() *ReadmemsConfig {
	folder, _ := homedir.Dir()
	filename := fmt.Sprintf("%s/.memsfcr.cfg", folder)
	c := NewConfig()

	cfg, err := ini.Load(filename)
	if err != nil {
		utils.LogI.Printf("failed to read file: %v", err)
		return c
	}

	c.Port = cfg.Section("").Key("port").String()
	c.Command = cfg.Section("").Key("command").String()
	c.Loop = cfg.Section("").Key("loop").String()
	c.Output = cfg.Section("").Key("output").String()
	c.LogFolder = cfg.Section("").Key("logfolder").String()
	c.Connection = cfg.Section("").Key("connection").String()

	utils.LogI.Println("MemsFCR Config", c)
	return c
}
