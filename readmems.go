package main

import (
	"andrewj.com/readmems/rosco"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
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
	var memsCommandResponse []byte
	var commandInterval time.Duration
	var logging = false
	var paused = false

	const DataInterval = 950
	const HeartbeatInterval = 100

	commandInterval = DataInterval * time.Millisecond // time in ms between requests

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
			rosco.LogI.Println("Lost connection to ECU, exiting")
			// break
		}

		if mems.Exit == true {
			// exit if the serial port is disconnected
			rosco.LogI.Println("Exit requested, exiting")
			break
		}

		count, _ := strconv.Atoi(config.Loop)
		ecuID := hex.EncodeToString(mems.ECUID)

		// enter a command / response loop
		for loop := 0; loop < count; {
			// check if we have received a command from the web interface
			commandID := recieveMessageFromWebView()

			switch commandID {
			case commandPauseDataLoop:
				{
					paused = true
					rosco.LogI.Printf("Paused Data Loop, sending heartbeats to keep connection alive")
				}
			case commandStartDataLoop:
				{
					paused = false
					rosco.LogI.Printf("Resuming Data Loop")
				}
			case commandResetECU:
				{
					rosco.LogI.Printf("memsCommandResponseLoop sending Reset ECU")
					memsCommandResponse = rosco.MemsSendCommand(mems, rosco.MEMS_ResetECU)
					rosco.LogI.Printf("memsCommandResponseLoop recieved from Reset ECU %x", memsCommandResponse)
				}
			case commandClearFaults:
				{
					rosco.LogI.Printf("memsCommandResponseLoop sending Clear Faults")
					rosco.MemsSendCommand(mems, rosco.MEMS_ClearFaults)
				}
			case commandResetAdjustments:
				{
					rosco.LogI.Printf("memsCommandResponseLoop sending Reset Adjustments")
					rosco.MemsSendCommand(mems, rosco.MEMS_ResetAdj)
				}
			case commandIncreaseIdleSpeed:
				{
					rosco.LogI.Printf("memsCommandResponseLoop sending Increase Idle Speed")
					rosco.MemsSendCommand(mems, rosco.MEMS_IdleSpeed_Increment)
				}
			case commandIncreaseIdleHot:
				{
					rosco.LogI.Printf("memsCommandResponseLoop sending Increase Idle Decay (Hot)")
					rosco.MemsSendCommand(mems, rosco.MEMS_IdleDecay_Increment)
				}
			case commandIncreaseFuelTrim:
				{
					rosco.LogI.Printf("memsCommandResponseLoop sending Increase Fuel Trim (LTFT)")
					rosco.MemsSendCommand(mems, rosco.MEMS_LTFT_Increment)
				}
			case commandIncreaseIgnitionAdvance:
				{
					rosco.LogI.Printf("memsCommandResponseLoop sending Increase Ignition Advance")
					rosco.MemsSendCommand(mems, rosco.MEMS_IgnitionAdvanceOffset_Increment)
				}
			case commandDecreaseIdleSpeed:
				{
					rosco.LogI.Printf("memsCommandResponseLoop sending Decrease Idle Speed")
					rosco.MemsSendCommand(mems, rosco.MEMS_IdleSpeed_Decrement)
				}
			case commandDecreaseIdleHot:
				{
					rosco.LogI.Printf("memsCommandResponseLoop sending Decrease Idle Decay (Hot)")
					rosco.MemsSendCommand(mems, rosco.MEMS_IdleDecay_Decrement)
				}
			case commandDecreaseFuelTrim:
				{
					rosco.LogI.Printf("memsCommandResponseLoop sending Decrease Fuel Trim (LTFT)")
					rosco.MemsSendCommand(mems, rosco.MEMS_LTFT_Decrement)
				}
			case commandDecreaseIgnitionAdvance:
				{
					rosco.LogI.Printf("memsCommandResponseLoop sending Decrease Ignition Advance")
					rosco.MemsSendCommand(mems, rosco.MEMS_IgnitionAdvanceOffset_Decrement)
				}
			default:
			}

			if paused {
				// send a heartbeat when paused
				commandInterval = HeartbeatInterval * time.Millisecond
				//rosco.MemsSendCommand(mems, rosco.MEMS_Heartbeat)
				rosco.LogI.Printf("Sending Heatbeat")
			} else {
				rosco.LogI.Printf("Reading from ECU")
				// read data from the ECU
				memsdata = rosco.MemsRead(mems)
				// send it to the web interface
				sendDataToWebView(memsdata)

				if logging {
					// write to a log file if logging is enabled
					go rosco.WriteMemsDataToFile(ecuID, memsdata)
				}

				// increment count of data calls
				// don't increment if we're paused
				loop = loop + 1
			}

			// sleep between calls to give the ECU time to catch up
			// the ECU will get slower as load increases so this ensures
			// a regular time series for the data set
			time.Sleep(commandInterval)

		}

		// read loop complete, exit
		break
	}
}

func recieveMessageFromWebView() int {
	select {
	case m := <-webToMemsChannel:
		{
			rosco.LogI.Printf("Recieved command from channel")
			commandID := evaluateCommand(m)
			return commandID
		}
	case <-time.After(2 * time.Second):
		{
			rosco.LogE.Printf("Command channel receive timeout")
		}
	}

	return commandUnknown
}

func sendDataToWebView(memsdata rosco.MemsData) {
	var m wsMsg

	m.Action = "data"

	data, _ := json.Marshal(memsdata)
	m.Data = string(data)
	memsToWebChannel <- m
}

func sendConfigToWebView(config *rosco.ReadmemsConfig) {
	var m wsMsg

	m.Action = "config"

	data, _ := json.Marshal(config)
	m.Data = string(data)
	memsToWebChannel <- m
}

func main() {
	var showHelp bool

	// create a map of commands
	createCommandMap()

	// use if the readmems config is supplied
	config = rosco.ReadConfig()

	// parse the command line parameters and override config file settings
	flag.StringVar(&config.Port, "port", config.Port, "Name/path of the serial port")
	flag.StringVar(&config.Command, "command", config.Command, "Command to send")
	flag.StringVar(&config.Loop, "loop", config.Loop, "Read loop count, 'inf' for infinite")
	flag.BoolVar(&showHelp, "help", false, "A brief help message")
	flag.Parse()

	if showHelp {
		rosco.LogI.Println(helpMessage())
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
