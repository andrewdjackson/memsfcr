package rosco

// MEMS Command List

// MEMSHeartbeat command code for a communication heartbeat
var MEMSHeartbeat = []byte{0xf4}

// MEMSDataFrame request complete dataframe using 0x7d and 0x80 coomands
var MEMSDataFrame = []byte{0x80, 0x7d}

// MEMSReqData80 command code for requesting data frame 0x80
var MEMSReqData80 = []byte{0x80}

// MEMSReqData7D command code for requesting data frame 0x7D
var MEMSReqData7D = []byte{0x7d}

// MEMSInitCommandA command code to start initialisation sequence
var MEMSInitCommandA = []byte{0xca}

// MEMSInitCommandB command code forms second command as part of the initialisation sequence
var MEMSInitCommandB = []byte{0x75}

// MEMSInitECUID command code for retrieving the ECU ID as the final step in initialisation
var MEMSInitECUID = []byte{0xd0}

// MEMSClearFaults command code to clear fault codes
var MEMSClearFaults = []byte{0xCC}

// MEMSGetIACPosition command code to retrieve the Idle Air Control position
var MEMSGetIACPosition = []byte{0xFB}

// MEMSResetAdj command code that instructs the ECU to clear all adjustments
var MEMSResetAdj = []byte{0x0F}

// MEMSResetECU command code that instructs the ECU to clear all computed and learnt settings
var MEMSResetECU = []byte{0xFA}

// MEMS Adjustment Settings
//
// | Setting                 | Decrement | Increment |
// | ----------------------- | --------- | --------- |
// | Fuel trim (Short Term?) |     7A    |     79    |
// | Fuel trim (Long Term?)  |     7C    |     7B    |
// | Idle decay              |     8A    |     89    |
// | Idle speed              |     92    |     91    |
// | Ignition advance offset |     94    |     93    |

// MEMSSTFTDecrement command
var MEMSSTFTDecrement = []byte{0x7a}

// MEMSSTFTIncrement command
var MEMSSTFTIncrement = []byte{0x79}

// MEMSLTFTDecrement command
var MEMSLTFTDecrement = []byte{0x7c}

// MEMSLTFTIncrement command
var MEMSLTFTIncrement = []byte{0x7b}

// MEMSIdleDecayDecrement commad
var MEMSIdleDecayDecrement = []byte{0x7c}

// MEMSIdleDecayIncrement command
var MEMSIdleDecayIncrement = []byte{0x7b}

// MEMSIdleSpeedDecrement command
var MEMSIdleSpeedDecrement = []byte{0x92}

// MEMSIdleSpeedIncrement command
var MEMSIdleSpeedIncrement = []byte{0x93}

// MEMSIgnitionAdvanceOffsetDecrement command
var MEMSIgnitionAdvanceOffsetDecrement = []byte{0x94}

// MEMSIgnitionAdvanceOffsetIncrement command
var MEMSIgnitionAdvanceOffsetIncrement = []byte{0x93}

// Actuators
var MEMSFuelPumpOn = []byte{0x11}
var MEMSFuelPumpOff = []byte{0x01}
var MEMSPTCRelayOn = []byte{0x12}
var MEMSPTCRelayOff = []byte{0x02}
var MEMSACRelayOn = []byte{0x13}
var MEMSACRelayOff = []byte{0x03}
var MEMSPurgeValveOn = []byte{0x18}
var MEMSPurgeValveOff = []byte{0x08}
var MEMSO2HeaterOn = []byte{0x19}
var MEMSO2HeaterOff = []byte{0x09}
var MEMSBoostValveOn = []byte{0x1B}
var MEMSBoostValveOff = []byte{0x0B}
var MEMSFan1On = []byte{0x1D}
var MEMSFan1Off = []byte{0x0D}
var MEMSFan2On = []byte{0x1E}
var MEMSFan2Off = []byte{0x0E}
var MEMSTestInjectors = []byte{0xF7}
var MEMSFireCoil = []byte{0xF8}
var MEMSOpenIAC = []byte{0xFD}
var MEMSCloseIAC = []byte{0xFE}
