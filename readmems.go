package main

import (
	"andrewj.com/readmems/rosco"
	"flag"
	"fmt"
	"strconv"
	"time"
)

const version = "v0.1.0"

var header = fmt.Sprintf("\nMemsFCR %s\n", version)

func helpMessage() string {
	return fmt.Sprintf(`%s
	ROSCO MEMS 1.6 Fault Code Reader

    Usage:
	memsfcr [flags]
	  
    Flags:
	-port		Name/path of the serial port 
	-command	Command to execute on the ECU {read}
	-loop		Command execution loop count, use 'inf' for infinite
	-output		Use 'stdout' to send response to console, 'file' to log in CSV format to a file {stdout | file}
	-wait		Retry the connection until a connection is established {true | false}
	-help		This help message
	`, header)
}

func main() {
	var showHelp bool

	// use if the readmems config is supplied
	var config = rosco.ReadConfig()

	// parse the command line parameters and override config file settings
	flag.StringVar(&config.Port, "port", config.Port, "Name/path of the serial port")
	flag.StringVar(&config.Command, "command", config.Command, "Command to send")
	flag.StringVar(&config.Loop, "loop", config.Loop, "Read loop count, 'inf' for infinite")
	flag.BoolVar(&showHelp, "help", false, "A brief help message")
	flag.Parse()

	if showHelp {
		fmt.Println(helpMessage())
		return
	}

	var memsdata rosco.MemsData
	var logging = false

	if config.Output == "file" {
		logging = true
	}

	if config.Loop == "inf" {
		// infitite loop, so set loop count to a very big number
		config.Loop = "10000000"
	}

	// connect and initialise the ECU
	mems := rosco.New()
	rosco.ConnectAndInitialiseECU(mems, config)

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

		count, _ := strconv.Atoi(config.Loop)

		for loop := 0; loop < count; loop++ {
			memsdata = rosco.MemsRead(mems)

			if logging {
				WriteMemsDataToFile(memsdata)
			}

			// sleep between calls to give the ECU time to catch up
			// the ECU will get slower as load increases so this ensures
			// a regular time series for the data set
			time.Sleep(950 * time.Millisecond)
		}
		break

	}
}
