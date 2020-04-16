package main

import (
	"andrewj.com/readmems/rosco"
	"andrewj.com/readmems/service"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// fileExists reports whether the named file or directory exists.
func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func connect(mems *rosco.Mems, config *rosco.ReadmemsConfig) {
	if !mems.Connected {
		rosco.MemsConnect(mems, config.Port)
		if mems.Connected {
			rosco.MemsInitialise(mems)
		}
	}
}

// WriteToFile will print any string of text to a file safely by
// checking for errors and syncing at the end.
func WriteToFile(filename string, data string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.WriteString(data); err != nil {
		return err
	}

	return file.Sync()
}

func main() {
	// use if the readmems config is supplied
	var config = rosco.ReadConfig()
	var memsdata rosco.MemsData

	// if argument is supplied then use that as the port id
	if len(os.Args) > 1 {
		config.Port = os.Args[1]
	}

	// connect to ECU
	mems := rosco.New()

	if config.WebPort != "0" {
		// start http service
		go service.StartService(mems, config)
	} else {
		fmt.Println("Disabling web interface")
	}

	connect(mems, config)

	for {
		// wait for comms

		if mems.SerialPort == nil {
			// exit if the serial port is disconnected
			fmt.Println("Lost connection to ECU, exiting")
			// break
		}

		if mems.Exit == true {
			// exit if the serial port is disconnected
			fmt.Println("Exit requested, exiting")
			break
		}

		if config.Loop == "inf" {
			memsdata = rosco.MemsRead(mems)
		} else {
			fmt.Println("Looping %s times", config.Loop)

			count, _ := strconv.Atoi(config.Loop)

			for loop := 0; loop < count; loop++ {
				memsdata = rosco.MemsRead(mems)
			}
			break
		}

		fmt.Println("%+v\n", memsdata)

		if config.Output == "file" {
			md, _ := json.Marshal(memsdata)
			WriteToFile("output.cvs", string(md))
		}
	}
}
