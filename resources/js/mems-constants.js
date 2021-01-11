// set to true to use the REST interface instead of the WebSocket
var useREST = false;

var sock = null;
var minLambda = false;
var maxLambda = false;
var minIAC = false;
var dataframeLoop;
var debug = false;
var replay = "";
var faultCount = 0
var derivedFaultCount = 0
var diagnosticsFaultCount = 0
var diagnosticReport = { "AnalysisCode": "optimal", "IsEngineRunning": false, "IsEngineWarming": true, "IsAtOperatingTemp": false, "IsEngineIdle": false, "IsEngineIdleFault": false, "IsCruising": false, "IsClosedLoop": false, "ClosedLoopExpected": false, "MapFault": false, "VacuumFault": false, "IdleAirControlFault": false, "LambdaFault": false, "CoolantTempSensorFault": false, "IntakeAirTempSensorFault": false, "FuelPumpCircuitFault": false, "ThrottlePotCircuitFault": false, "IACPosition": 0 }

// replay data
var replayCount = 0
var replayPosition = 0

// duration in milliseconds between calls to the ECU for
// dataframes. the ECU will struggle to respond with a 
// value less than 450ms
var ECUQueryInterval = 450

// wait time for the ECU to respond before sending another command
var waitingForResponse = false;
var waitingForResponseTimeout;

// Constants
const WaitForResponseInterval = ECUQueryInterval * 2

const AirSensorFaultCode = 0b00000001
const CoolantSensorFaultCode = 0b00000010
const FuelPumpFaultCode = 0b00000001
const ThrottlePotFaultCode = 0b01000000

const ResponseSTFTDecrement = "7a"
const ResponseSTFTIncrement = "79"
const ResponseLTFTDecrement = "7c"
const ResponseLTFTIncrement = "7b"
const ResponseIdleDecayDecrement = "7c"
const ResponseIdleDecayIncrement = "7b"
const ResponseIdleSpeedDecrement = "92"
const ResponseIdleSpeedIncrement = "93"
const ResponseIgnitionAdvanceOffsetDecrement = "94"
const ResponseIgnitionAdvanceOffsetIncrement = "93"

// web actions over the websocket protocol
const WebActionSave = "save";
const WebActionConfig = "config";
const WebActionConnection = "connection";
const WebActionConnect = "connect"
const WebActionData = "data";
const WebActionCommand = "command"
const WebActionResponse = "response"
const WebActionIncrease = "increase"
const WebActionDecrease = "decrease"
const WebActionDiagnostics = "diagnostics"
const WebActionReplay = "replay"

// WebActionCommand commands
const CommandStart = "start"
const CommandPause = "pause"
const CommandReadConfig = "read"
const CommandResetECU = "resetecu"
const CommandResetAdjustments = "resetadj"
const CommandClearFaults = "clearfaults"
const CommandDataFrame = "dataframe"

// adjustments
const AdjustmentIdleSpeed = "idlespeed"
const AdjustmentIdleHot = "idlehot"
const AdjustmentIgnitionAdvance = "ignitionadvance"
const AdjustmentFuelTrim = "fueltrim"

// actuators
const ActuatorFuelPumpOn = "11"
const ActuatorFuelPumpOff = "01"
const ActuatorPTCRelayOn = "12"
const ActuatorPTCRelayOff = "02"
const ActuatorACRelayOn = "13"
const ActuatorACRelayOff = "03"
const ActuatorPurgeValveOn = "18"
const ActuatorPurgeValveOff = "08"
const ActuatorO2HeaterOn = "19"
const ActuatorO2HeaterOff = "09"
const ActuatorBoostValveOn = "1B"
const ActuatorBoostValveOff = "0B"
const ActuatorFan1On = "1D"
const ActuatorFan1Off = "0D"
const ActuatorFan2On = "1E"
const ActuatorFan2Off = "0E"
const ActuatorFan3On = "6F"
const ActuatorFan3Off = "67"
const ActuatorWasteGateOn = "1B"
const ActuatorWasteGateOff = "0B"
const ActuatorTestInjectors = "F7"
const ActuatorTestInjectorsMPi = "EF"
const ActuatorFireCoil = "F8"
const ActuatorOpenIAC = "FD"
const ActuatorCloseIAC = "FE"
const ActuatorAllActuatorsOff = "F4"
const ActuatorFuelTrimPlus = "79"
const ActuatorFuelTrimMinus = "7A"
const ActuatorIdleDecayPlus = "89"
const ActuatorIdleDecayMinus = "8A"
const ActuatorIdleSpeedPlus = "91"
const ActuatorIdleSpeedMinus = "92"
const ActuatorIgnitionAdvancePlus = "93"
const ActuatorIgnitionAdvanceMinus = "94"

// settings
const SettingLogFolder = "logfolder"
const SettingLogToFile = "logtofile"
const LogToFileEnabled = "true"
const LogToFileDisabled = "false"
const SettingPort = "port"
const SettingPortList = "ports"
const SettingECUQueryFrequency = "ecuqueryfrequency"

// Indicators and Labels
const IndicatorConnectionMessage = "connectionMessage"
const IndicatorECUConnected = "ecudata"
const IndicatorECUFault = "ecufault"
const IndicatorCoolantFault = "coolantfault"
const IndicatorAirFault = "airfault"
const IndicatorThrottleFault = "throttlefault"
const IndicatorFuelFault = "fuelfault"
const IndicatorClosedLoop = "closedloop"
const IndicatorIdleSwitch = "idleswitch"
const IndicatorParkSwitch = "parkswitch"
const IndicatorLambdaLowFault = "lambdalowfault"
const IndicatorLambdaHighFault = "lambdahighfault"
const IndicatorLambdaLow = "lambdalow"
const IndicatorLambdaHigh = "lambdahigh"
const IndicatorRPMSensor = "rpmsensor"
const IndicatorIACLow = "iaclow"
const IndicatorO2SystemFault = "systemfault"

// LED statuses 
const LEDFault = "fault"
const LEDStatus = "status"
const LEDWarning = "warning"

// chart labels - must match id's used in the html
const ChartRPM = "rpmchart"
const ChartLambda = "lambdachart"
const ChartLoopIndicator = "loopchart"
const ChartCoolant = "coolantchart"
const ChartAFR = "afrchart"
const ReplayProgress = "replayprogress"
const ReplayProgressRemaining = "replayprogressremaining"

// spark labels - must match id's used in the html
const SparkRPM = "rpmspark"
const SparkMAP = "mapspark"
const SparkThrottle = "throttlespark"
const SparkIAC = "iacspark"
const SparkBattery = "batteryspark"
const SparkCoolant = "coolantspark"
const SparkAir = "airspark"
const SparkLambda = "lambdaspark"
const SparkFuel = "fuelspark"
const SparkLTFuel = "ltfuelspark"
const SparkAirFuel = "airfuelspark"
const SparkIgnition = "ignitionspark"

const IdleSpeedAdjustment = "idlespeed"
const IdleHotAdjustment = "idlehot"
const FuelTrimAdjustment = "fueltrim"
const IgnitionAdvanceAdjustment = "ignitionadvance"