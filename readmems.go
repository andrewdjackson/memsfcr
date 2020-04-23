package main

import (
	"andrewj.com/readmems/rosco"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"
)

const version = "v0.1.0"

var header = fmt.Sprintf("\nMemsFCR %s\n", version)
var config *rosco.ReadmemsConfig
var mems *rosco.Mems

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

func memsCommandResponseLoop(config *rosco.ReadmemsConfig) {
	var memsdata rosco.MemsData
	var logging = false
	var paused = false

	if config.Output == "file" {
		logging = true
	}

	// connect and initialise the ECU
	mems = rosco.New()
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
		ecuID := hex.EncodeToString(mems.ECUID)

		// enter a command / response loop
		for loop := 0; loop < count; loop++ {
			// check if we have received a pause / start command
			/*select {
			case m := <-commandChannel:
				if m.Data == "pause" {
					paused = true
					log.Printf("Paused Data Loop, sending heartbeats to keep connection alive")
				}
				if m.Data == "start" {
					paused = false
					log.Printf("Resuming Data Loop")
				}
			default:
			}*/

			if paused {
				// send a heartbeat when paused
				rosco.MemsSendCommand(mems, rosco.MEMS_Heartbeat)
				log.Printf("Sending Heatbeat")
			} else {
				// read data from the ECU
				memsdata = rosco.MemsRead(mems)
				// send it to the web interface
				sendDataToWebView(memsdata)

				if logging {
					// write to a log file if logging is enabled
					WriteMemsDataToFile(ecuID, memsdata)
				}

				// increment count of data calls
				// don't increment if we're paused
				loop = loop + 1
			}

			// sleep between calls to give the ECU time to catch up
			// the ECU will get slower as load increases so this ensures
			// a regular time series for the data set
			time.Sleep(950 * time.Millisecond)

		}

		// read loop complete, exit
		break
	}

}

func sendCommandToMems(m wsMsg) []byte {
	var command []byte

	if m.Action == "decrease" && m.Data == "idlespeed" {
		command = rosco.MEMS_IdleSpeed_Increment
	}
	if m.Action == "increase" && m.Data == "idlespeed" {
		command = rosco.MEMS_IdleSpeed_Decrement
	}

	if m.Action == "decrease" && m.Data == "idlehot" {
		command = rosco.MEMS_IdleDecay_Decrement
	}
	if m.Action == "increase" && m.Data == "idlehot" {
		command = rosco.MEMS_IdleDecay_Increment
	}

	if m.Action == "decrease" && m.Data == "ignitionadvance" {
		command = rosco.MEMS_IgnitionAdvanceOffset_Decrement
	}
	if m.Action == "increase" && m.Data == "ignitionadvance" {
		command = rosco.MEMS_IgnitionAdvanceOffset_Increment
	}

	if m.Action == "decrease" && m.Data == "fueltrim" {
		command = rosco.MEMS_LTFT_Decrement
	}
	if m.Action == "increase" && m.Data == "fueltrim" {
		command = rosco.MEMS_LTFT_Increment
	}

	if m.Action == "command" {
		if m.Data == "clear" {
			command = rosco.MEMS_ClearFaults
		}
		if m.Data == "resetecu" {
			command = rosco.MEMS_ResetECU
		}
		if m.Data == "resetadj" {
			command = rosco.MEMS_ResetAdj
		}
		if m.Data == "iacposition" {
			command = rosco.MEMS_GetIACPosition
		}
	}

	return rosco.MemsSendCommand(mems, command)
}

func sendDataToWebView(memsdata rosco.MemsData) {
	var m wsMsg

	m.Action = "data"

	data, _ := json.Marshal(memsdata)
	m.Data = string(data)
	memsChannel <- m
}

func sendConfigToWebView(config *rosco.ReadmemsConfig) {
	var m wsMsg

	m.Action = "config"

	data, _ := json.Marshal(config)
	m.Data = string(data)
	memsChannel <- m
}

func main() {
	var showHelp bool

	// use if the readmems config is supplied
	config = rosco.ReadConfig()

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

	if config.Loop == "inf" {
		// infitite loop, so set loop count to a very big number
		config.Loop = "10000000"
	}

	go RunHTTPServer()
	go sendConfigToWebView(config)

	ShowWebView(config)
}
