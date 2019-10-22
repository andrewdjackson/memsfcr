package main

import (
	"andrewj.com/readmems/rosco"
	"andrewj.com/readmems/service"
	"fmt"
	"os"
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

func main() {
	// use if the readmems config is supplied
	var config = rosco.ReadConfig()

	// if argument is supplied then use that as the port id
	if len(os.Args) > 1 {
		config.Port = os.Args[1]
	}

	// connect to ECU
	mems := rosco.New()

	// start http service
	go service.StartService(mems, config)

	defer connect(mems, config)

	for {
		// wait for comms

		if mems.SerialPort == nil {
			// exit if the serial port is disconnected
			fmt.Println("Lost connection to ECU, exiting")
			break
		}

		if mems.Exit == true {
			// exit if the serial port is disconnected
			fmt.Println("Exit requested, exiting")
			break
		}
	}
}
