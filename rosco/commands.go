package rosco

// MEMS Command List

// MEMS_Heartbeat command code for a communication heartbeat
var MEMS_Heartbeat = []byte{0xf4}

// MEMS_ReqData80 command code for requesting data frame 0x80
var MEMS_ReqData80 = []byte{0x80}

// MEMS_ReqData7D command code for requesting data frame 0x7D
var MEMS_ReqData7D = []byte{0x7d}

// MEMS_InitCommandA command code to start initialisation sequence
var MEMS_InitCommandA = []byte{0xca}

// MEMS_InitCommandB command code forms second command as part of the initialisation sequence
var MEMS_InitCommandB = []byte{0x75}

// MEMS_InitECUID command code for retrieving the ECU ID as the final step in initialisation
var MEMS_InitECUID = []byte{0xd0}

// MEMS_ClearFaults command code to clear fault codes
var MEMS_ClearFaults = []byte{0xCC}

// MEMS_GetIACPosition command code to retrieve the Idle Air Control position
var MEMS_GetIACPosition = []byte{0xFB}

// MEMS_ResetAdj command code that instructs the ECU to clear all adjustments
var MEMS_ResetAdj = []byte{0x0F}

// MEMS_ResetECU command code that instructs the ECU to clear all computed and learnt settings
var MEMS_ResetECU = []byte{0xFA}

// MEMS Adjustment Settings
//
// | Setting                 | Decrement | Increment |
// | ----------------------- | --------- | --------- |
// | Fuel trim (Short Term?) |     7A    |     79    |
// | Fuel trim (Long Term?)  |     7C    |     7B    |
// | Idle decay              |     8A    |     89    |
// | Idle speed              |     92    |     91    |
// | Ignition advance offset |     94    |     93    |

var MEMS_STFT_Decrement = []byte{0x7a}
var MEMS_STFT_Increment = []byte{0x79}
var MEMS_LTFT_Decrement = []byte{0x7c}
var MEMS_LTFT_Increment = []byte{0x7b}
var MEMS_IdleDecay_Decrement = []byte{0x7c}
var MEMS_IdleDecay_Increment = []byte{0x7b}
var MEMS_IdleSpeed_Decrement = []byte{0x92}
var MEMS_IdleSpeed_Increment = []byte{0x93}
var MEMS_IgnitionAdvanceOffset_Decrement = []byte{0x94}
var MEMS_IgnitionAdvanceOffset_Increment = []byte{0x93}
