package ui

// WebAction constants
const (
	// WebActionConfig data packet is config
	WebActionConfig = "config"
	// WebActionConnection data packet is connection status
	WebActionConnection = "connection"
	// WebActionConnect action is to connect the ecu
	WebActionConnect = "connect"
	// WebActionECUCommand data packet is an ecu command
	WebActionECUCommand = "command"
	// WebActionECUCommandIncrease data packet is an increase adjustment command
	WebActionECUCommandIncrease = "command"
	// WebActionECUCommandDecrease data packet is an decrease adjustment command
	WebActionECUCommandDecrease = "command"
	// WebActionData data packet is ecu data
	WebActionData = "data"
)

const (
	// Unknown command
	Unknown                 = 0
	ConnectECU              = 1
	PauseDataLoop           = 2
	StartDataLoop           = 3
	ResetECU                = 4
	ResetAdjustments        = 5
	ClearFaults             = 6
	IncreaseIdleSpeed       = 7
	IncreaseIdleHot         = 8
	IncreaseFuelTrim        = 9
	IncreaseIgnitionAdvance = 10
	DecreaseIdleSpeed       = 11
	DecreaseIdleHot         = 12
	DecreaseFuelTrim        = 13
	DecreaseIgnitionAdvance = 14
	ConfigRead              = 15
	ConfigWrite             = 16
	Dataframe               = 17
)

// WebAction converts the JSON message from the
// web into a code
type WebAction struct {
	Msg   WebMsg
	Value int
}

// EvaluateWebMsg converts the JSON message from the
// web into a code
func EvaluateWebMsg(m WebMsg) WebAction {
	switch m.Action {
	// connect action is the same as command / connec
	case "connect":
		return WebAction{m, ConnectECU}
	// process ECU commands
	case "command":
		switch m.Data {
		// connect to ECU
		case "connect":
			return WebAction{m, ConnectECU}
		case "resetecu":
			return WebAction{m, ResetECU}
		case "resetadj":
			return WebAction{m, ResetAdjustments}
		case "clearfaults":
			return WebAction{m, ClearFaults}
		case "pause":
			return WebAction{m, PauseDataLoop}
		case "start":
			return WebAction{m, StartDataLoop}
		}
	case "increase":
		switch m.Data {
		case "idlespeed":
			return WebAction{m, IncreaseIdleSpeed}
		case "idlehot":
			return WebAction{m, IncreaseIdleHot}
		case "fueltrim":
			return WebAction{m, IncreaseFuelTrim}
		case "ignition":
			return WebAction{m, IncreaseIgnitionAdvance}
		}
	case "decrease":
		switch m.Data {
		case "idlespeed":
			return WebAction{m, DecreaseIdleSpeed}
		case "idlehot":
			return WebAction{m, DecreaseIdleHot}
		case "fueltrim":
			return WebAction{m, DecreaseFuelTrim}
		case "ignition":
			return WebAction{m, DecreaseIgnitionAdvance}
		}
	case "config":
		switch m.Data {
		case "read":
			return WebAction{m, ConfigRead}
		case "write":
			return WebAction{m, ConfigWrite}
		}
	default:
	}

	return WebAction{m, Unknown}
}
