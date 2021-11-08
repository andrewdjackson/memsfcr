
var minLambda = false
var maxLambda = false
var minIAC = false
var debug = false

// replay data
var replay = ""
var replayCount = 0
var replayPosition = 0

// duration in milliseconds between calls to the ECU for
// dataframes. the ECU will struggle to respond with a
// value less than 450ms
var ECUQueryInterval = 1000
var ECUHeartbeatInterval = 2000
var dataframeLoop

var resetDataframe = {
    "Time": "00:00:00.000",
    "EngineRPM": 0,
    "CoolantTemp": 0,
    "AmbientTemp": 0,
    "IntakeAirTemp": 0,
    "FuelTemp": 0,
    "ManifoldAbsolutePressure": 0,
    "BatteryVoltage": 0,
    "ThrottlePotSensor": 0.0,
    "ThrottlePosition": 0,
    "IdleSwitch": false,
    "AirconSwitch": false,
    "ParkNeutralSwitch": false,
    "DTC0": 0,
    "DTC1": 0,
    "IdleSetPoint": 0,
    "IdleHot": 0,
    "Uk8011": 0,
    "IACPosition": 0,
    "IdleSpeedDeviation": 0,
    "IgnitionAdvanceOffset80": 0,
    "IgnitionAdvance": 0,
    "CoilTime": 0,
    "CrankshaftPositionSensor": 0,
    "Uk801a": 0,
    "Uk801b": 0,
    "IgnitionSwitch": true,
    "ThrottleAngle": 0,
    "Uk7d03": 0,
    "AirFuelRatio": 0,
    "DTC2": 0,
    "LambdaVoltage": 0,
    "LambdaFrequency": 0,
    "LambdaDutycycle": 0,
    "LambdaStatus": 0,
    "ClosedLoop": false,
    "LongTermFuelTrim": 0,
    "ShortTermFuelTrim": 0,
    "FuelTrimCorrection": 0,
    "CarbonCanisterPurgeValve": 0,
    "DTC3": 0,
    "IdleBasePosition": 0,
    "Uk7d10": 0,
    "DTC4": 0,
    "IgnitionAdvanceOffset7d": 0,
    "IdleSpeedOffset": 0,
    "Uk7d14": 0,
    "Uk7d15": 0,
    "DTC5": 0,
    "Uk7d17": 0,
    "Uk7d18": 0,
    "Uk7d19": 0,
    "Uk7d1a": 0,
    "Uk7d1b": 0,
    "Uk7d1c": 0,
    "Uk7d1d": 0,
    "Uk7d1e": 0,
    "JackCount": 0,
    "CoolantTempSensorFault": false,
    "IntakeAirTempSensorFault": false,
    "FuelPumpCircuitFault": false,
    "ThrottlePotCircuitFault": false,
    "Analytics": {
    "ReadingFault": false,
        "IsEngineRunning": false,
        "IsEngineWarming": false,
        "IsAtOperatingTemp": false,
        "IsEngineIdle": false,
        "IsEngineIdleFault": false,
        "IdleSpeedFault": false,
        "IdleErrorFault": false,
        "IdleHotFault": false,
        "IdleBaseFault": false,
        "IsCruising": false,
        "IsClosedLoop": false,
        "IsClosedLoopExpected": false,
        "ClosedLoopFault": false,
        "IsThrottleActive": false,
        "MapFault": false,
        "VacuumFault": false,
        "IdleAirControlFault": false,
        "IdleAirControlRangeFault": false,
        "IdleAirControlJackFault": false,
        "O2SystemFault": false,
        "LambdaRangeFault": false,
        "LambdaOscillationFault": false,
        "ThermostatFault": false,
        "CoolantTempSensorFault": false,
        "IntakeAirTempSensorFault": false,
        "FuelPumpCircuitFault": false,
        "ThrottlePotCircuitFault": false,
        "CrankshaftSensorFault": false,
        "CoilFault": false,
        "IACPosition": 0
    },
    "Dataframe80": "801c000085ff4fff638e23001001000000208b60039d003808c1000000",
    "Dataframe7d": "7d201012ff92006effff0100996400ff3affff30807c63ff19401ec0264034c008"
}

const AirSensorFaultCode = 0b00000001
const CoolantSensorFaultCode = 0b00000010
const FuelPumpFaultCode = 0b00000001
const ThrottlePotFaultCode = 0b01000000

// adjustments
const AdjustmentIdleSpeed = "idlespeed"
const AdjustmentIdleDecay = "idledecay"
const AdjustmentIgnitionAdvance = "ignitionadvance"
const AdjustmentSTFT = "stft"
const AdjustmentLTFT = "ltft"
const AdjustmentIAC = "iac"

// settings
const SettingLogFolder = "logfolder"
const SettingLogToFile = "logtofile"
const LogToFileEnabled = "true"
const LogToFileDisabled = "false"
const SettingPort = "port"
const SettingPortList = "ports"
const SettingECUQueryFrequency = "ecuqueryfrequency"
const LabelECUQueryFrequency = "ecuqueryfrequencylabel"

// Indicators and Labels
const IndicatorConnectionMessage = "connectionMessage"
const IndicatorECUConnected = "ecudata"
const IndicatorECUFault = "ecufault"
const IndicatorCoolantFault = "coolantfault"
const IndicatorAirFault = "airfault"
const IndicatorThrottleFault = "throttlefault"
const IndicatorFuelFault = "fuelfault"
const IndicatorClosedLoop = "closedloop"
//const IndicatorIdleSwitch = "idleswitch"
const IndicatorParkSwitch = "parkswitch"
const IndicatorLambdaRangeFault = "lambdarangefault"
const IndicatorRPMSensor = "rpmsensor"
const IndicatorIACLow = "iaclow"
const IndicatorO2SystemFault = "systemfault"

// Analytics LEDs
const DashboardEngineRunning = "dashboard-enginerunning"
const DashboardCrankshaftSensorFault = "dashboard-crankshaftsensor"
const DashboardMapFault = "dashboard-mapfault"

const AnalyticsReadingFault = "analytics-readingfault"
const AnalyticsIsEngineRunning = "analytics-enginerunning"
const AnalyticsIsEngineWarming = "analytics-enginewarming"
const AnalyticsIsAtOperatingTemp = "analytics-operatingtemp"
const AnalyticsIsEngineIdle = "analytics-engineidle"
const AnalyticsIsEngineIdleFault = "analytics-engineidlefault"
const AnalyticsIdleSpeedFault = "analytics-idlespeedfault"
const AnalyticsIdleErrorFault = "analytics-idleerrorfault"
const AnalyticsIdleHotFault = "analytics-idlehotfault"
const AnalyticsIdleBaseFault = "analytics-idlebasefault"
const AnalyticsIsCruising = "analytics-cruising"
const AnalyticsIsClosedLoop = "analytics-closedloop"
const AnalyticsIsClosedLoopExpected = "analytics-closedloopexpected"
const AnalyticsClosedLoopFault = "analytics-closedloopfault"
const AnalyticsIsThrottleActive = "analytics-throttleactive"
const AnalyticsMapFault = "analytics-mapfault"
const AnalyticsVacuumFault = "analytics-vacuumfault"
const AnalyticsIdleAirControlFault = "analytics-iacfault"
const AnalyticsIdleAirControlRangeFault = "analytics-iacrangefault"
const AnalyticsIdleAirControlJackFault = "analytics-iacjackfault"
const AnalyticsO2SystemFault = "analytics-o2systemfault"
const AnalyticsLambdaRangeFault = "analytics-lambdarangefault"
const AnalyticsLambdaOscillationFault = "analytics-lambdaoscfault"
const AnalyticsThermostatFault = "analytics-thermostatfault"
const AnalyticsCoolantTempSensorFault = "analytics-coolanttempfault"
const AnalyticsIntakeAirTempSensorFault = "analytics-airtempfault"
const AnalyticsFuelPumpCircuitFault = "analytics-fuelpumpfault"
const AnalyticsThrottlePotCircuitFault = "analytics-throttlepotfault"
const AnalyticsCrankshaftSensorFault = "analytics-crankshaftfault"
const AnalyticsCoilFault = "analytics-coilfault"

// LED statuses
const LEDFault = "fault"
const LEDStatus = "status"
const LEDWarning = "warning"
const LEDInfo = "info"

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

const maxDebugLogLength = 75 // lines of debug log in the interface
const debugLogLineTerminator = "<br>"
var   debugLogLineCount = 0 // number of lines in the log

if (typeof console != "undefined") {
    var oldLogInfo = console.info
    var oldLogWarn = console.warn
    var oldLogError = console.error

    console.info = function(message) {
        oldLogInfo.apply(console, arguments);
        display("<span class='debugInfo'>INFO</span>", message);
    }

    console.warn = function(message) {
        oldLogWarn.apply(console, arguments);
        display("<span class='debugWarning'>WARN</span>", message);
    }

    console.error = function(message) {
        oldLogError.apply(console, arguments);
        display("<span class='debugError'>ERRO</span>", message);
    }

    function display(level, message) {
        date = new Date()
        time = String(date.getHours()).padStart(2,'0') + ":" + String(date.getMinutes()).padStart(2,'0') + ":" + String(date.getSeconds()).padStart(2,'0') + "." + String(date.getMilliseconds()).padStart(3, '0')
        debugLogLineCount++

        debugLogContent = document.getElementById('debugLog').innerHTML
        if (debugLogLineCount > maxDebugLogLength) {
            // set the line count to the max
            debugLogLineCount = maxDebugLogLength
            // find end of the line
            start = debugLogContent.indexOf(debugLogLineTerminator, 0) + debugLogLineTerminator.length
            // truncate the content
            debugLogContent = debugLogContent.substring(start, debugLogContent.length)

            document.getElementById('debugLog').innerHTML = debugLogContent
        }
        document.getElementById('debugLog').innerHTML += (level + "[" + time + "] " + message + debugLogLineTerminator);
    }
}

// enable tooltips
$(function () {
    $('[data-toggle="tooltip"]').tooltip()
})

class MemsReader {
    constructor(uri) {
        this.uri = {
            config: uri + "/config",
            ports: uri + "/config/ports",
            connect: uri + "/rosco/connect",
            disconnect: uri + "/rosco/disconnect",
            heartbeat: uri + "/rosco/heartbeat",
            dataframe: uri + "/rosco/dataframe",
            adjust: uri + "/rosco/adjust/",
            actuator: uri + "/rosco/test/",
            scenario: uri + "/scenario",
            play_scenario: uri + "/scenario/play",
            seek_scenario: uri + "/scenario/seek",
        }
        this.ecuid = ""
        this.iacposition = 0
        this.memsdata = {}
        this.status = {
            port: "",
            emulated: false,
            connected: false,
            paused: false,
            heartbeatActive: false,
        }
        // connect button
        this.connectButton = document.getElementById('connectECUbtn')
        this.connectButton.addEventListener('click', connectECU, {once: true});

        // play / pause button (disabled until connected)
        this.playPauseButton = document.getElementById('playPauseECUbtn')
        this.playPauseButton.disabled=true

        // replay button
        this.replayButton = document.getElementById('replayECUbtn')
        this.replayButton.addEventListener('click', loadScenarios, {once: true});
        this.replayButton.disabled=false
    }

    port() {
        if (memsreader.status.emulated) {
            return {"port": selectedScenario};
        } else {
            return {"port": document.getElementById(SettingPort).value};
        }
    }
}

// this function gets called as soon as the page load has completed
window.onload = function () {
    uri = window.location.href.split("/").slice(0, 3).join("/");
    memsreader = new MemsReader(uri)

    // establish a keep-alive heartbeat and listen
    // for server state changes
    initialiseServerEvents()

    // draw the gauges
    initialiseGauges()

    // create gauge sparklines
    initialiseSparklines()

    // create the profiling line charts
    initialiseGraphs()

    // get the configuration parameters
    readConfig()

    // hide the playback progress bar
    showProgressValues(false)
};

function initialiseServerEvents() {
    // connect to the server to establish a heartbeat link
    // if the user closes the browser, the server will detect no response
    // and terminate the application after a few seconds
    server_event = new EventSource(uri + "/heartbeat")

    server_event.onopen = function () {
        console.debug("server-event connected");
    }

    server_event.onclose = function () {
        console.debug("server-event close");
    }

    // server-event message handler
    server_event.onmessage = function (e) {
        console.debug("server-event message " + e.data);
    }

    // listen for heartbeat events
    server_event.addEventListener('heartbeat', heartbeatHandler, false);

    // listen for ecu connection state changes
    server_event.addEventListener('status', statusHandler, false);
}

function heartbeatHandler(e) {
    console.debug('server heartbeat')
}

function statusHandler(e) {
    console.debug('server status change')
}

function initialiseGauges() {
    gaugeRPM.draw();
    gaugeMap.draw();
    gaugeThrottlePos.draw();
    gaugeIACPos.draw();
    gaugeBattery.value = 11;
    gaugeBattery.draw();
    gaugeCoolant.draw();
    gaugeAir.draw();
    gaugeLambda.draw();
    gaugeFuelTrim.draw();
    gaugeLTFuelTrim.draw();
    gaugeAirFuel.draw();
    gaugeIgnition.draw();

    // draw adaptive value gauges
    gaugeAdaptiveIdleSpeed.draw()
    gaugeAdaptiveIACPos.draw()
    gaugeAdaptiveIdleDecay.draw()
    gaugeAdaptiveSTFT.draw()
    gaugeAdaptiveLTFT.draw()
    gaugeAdaptiveIgnition.draw()
}

function initialiseSparklines() {
    rpmSpark = createSpark(SparkRPM)
    mapSpark = createSpark(SparkMAP)
    throttleSpark = createSpark(SparkThrottle)
    iacSpark = createSpark(SparkIAC)
    batterySpark = createSpark(SparkBattery)
    coolantSpark = createSpark(SparkCoolant)
    airSpark = createSpark(SparkAir)
    lambdaSpark = createSpark(SparkLambda)
    fuelSpark = createSpark(SparkFuel)
    ltfuelSpark = createSpark(SparkLTFuel)
    airfuelSpark = createSpark(SparkAirFuel)
    ignitionSpark = createSpark(SparkIgnition)
}

function initialiseGraphs() {
    rpmChart = createChart(ChartRPM, "Engine (RPM)");
    lambdaChart = createChart(ChartLambda, "Lambda (mV)");
    loopChart = createChart(ChartLoopIndicator, "O2 Loop (0 = Active)");
    afrChart = createChart(ChartAFR, "Air : Fuel Ratio");
    coolantChart = createChart(ChartCoolant, "Coolant (°C)");
}

async function resetInterface() {
    console.info('resetting user interface')
    // wait for the dataframe interval to complete
    await new Promise(r => setTimeout(r, 500));
    updateECUDataframe(resetDataframe)
}

function updateGauges(data) {
    gaugeRPM.value = data.EngineRPM;
    gaugeMap.value = data.ManifoldAbsolutePressure;
    gaugeThrottlePos.value = (data.ThrottlePotSensor) * 20;
    gaugeIACPos.value = data.IACPosition;
    gaugeBattery.value = data.BatteryVoltage;
    gaugeCoolant.value = data.CoolantTemp;
    gaugeAir.value = data.IntakeAirTemp;
    gaugeLambda.value = data.LambdaVoltage;
    gaugeFuelTrim.value = data.FuelTrimCorrection;
    gaugeLTFuelTrim.value = data.LongTermFuelTrim;
    gaugeAirFuel.value = data.AirFuelRatio;
    gaugeIgnition.value = data.IgnitionAdvance;
}

function updateGraphs(data) {
    addData(rpmSpark, data.Time, data.EngineRPM);
    addData(mapSpark, data.Time, data.ManifoldAbsolutePressure);
    addData(throttleSpark, data.Time, data.ThrottlePotSensor);
    addData(iacSpark, data.Time, data.IACPosition);
    addData(batterySpark, data.Time, data.BatteryVoltage);
    addData(coolantSpark, data.Time, data.CoolantTemp);
    addData(airSpark, data.Time, data.IntakeAirTemp);
    addData(lambdaSpark, data.Time, data.LambdaVoltage);
    addData(fuelSpark, data.Time, data.FuelTrimCorrection);
    addData(ltfuelSpark, data.Time, data.LongTermFuelTrim);
    addData(airfuelSpark, data.Time, data.AirFuelRatio);
    addData(ignitionSpark, data.Time, data.IgnitionAdvance);

    addData(rpmChart, data.Time, data.EngineRPM);
    addData(lambdaChart, data.Time, data.LambdaVoltage);
    addData(loopChart, data.Time, data.ClosedLoop);
    addData(afrChart, data.Time, data.AirFuelRatio);
    addData(coolantChart, data.Time, data.CoolantTemp);
}

function setConnectionStatusMessage(connected) {
    id = IndicatorConnectionMessage

    $('#' + id).removeClass("alert-light");
    $('#' + id).removeClass("alert-danger");
    $('#' + id).removeClass("alert-success");

    $('#' + id).removeClass("invisible");
    $('#' + id).addClass("visible");

    if (connected == true) {
        if (replay == "") {
            msg = document.getElementById("port").value
        } else {
            msg = replay
        }

        document.getElementById(id).textContent = "connected to " + msg
        $('#' + id).addClass("alert-success");
    } else {
        document.getElementById(id).textContent = "unable to connect to " + msg
        $('#' + id).addClass("alert-danger");
    }
}

function setECUQueryFrequency(frequency) {
    console.info("freq " + frequency)
    f = parseInt(frequency)
    if (f > 200) {
        ECUQueryInterval = f
        updateAdjustmentValue(SettingECUQueryFrequency, ECUQueryInterval)
    }
}

function updateLEDs(data) {
    var derived = 0;

    if (data.DTC0 > 0 || data.DTC1 > 0) {
        setStatusLED(data.CoolantTempSensorFault, IndicatorCoolantFault, LEDFault);
        setStatusLED(data.IntakeAirTempSensorFault, IndicatorAirFault, LEDFault);
        setStatusLED(data.ThrottlePotCircuitFault, IndicatorThrottleFault, LEDFault);
        setStatusLED(data.FuelPumpCircuitFault, IndicatorFuelFault, LEDFault);
    }

    setStatusLED(data.ClosedLoop, IndicatorClosedLoop, LEDStatus);
    //setStatusLED(data.IdleSwitch, IndicatorIdleSwitch, LEDStatus);
    setStatusLED(data.ParkNeutralSwitch, IndicatorParkSwitch, LEDStatus);

    setStatusLED(memsreader.memsdata.Analytics.O2SystemFault, IndicatorO2SystemFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.CrankshaftSensorFault, IndicatorRPMSensor, LEDWarning);
    setStatusLED(memsreader.memsdata.Analytics.LambdaRangeFault, IndicatorLambdaRangeFault, LEDWarning);
    setStatusLED(memsreader.memsdata.Analytics.IdleAirControlFault, IndicatorIACLow, LEDWarning);
    setStatusLED(memsreader.memsdata.Analytics.ReadingFault, AnalyticsReadingFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.IsEngineRunning, AnalyticsIsEngineRunning, LEDStatus);
    setStatusLED(memsreader.memsdata.Analytics.IsEngineWarming, AnalyticsIsEngineWarming, LEDStatus);
    setStatusLED(memsreader.memsdata.Analytics.IsAtOperatingTemp, AnalyticsIsAtOperatingTemp, LEDStatus);
    setStatusLED(memsreader.memsdata.Analytics.IsEngineIdle, AnalyticsIsEngineIdle , LEDStatus);
    setStatusLED(memsreader.memsdata.Analytics.IsEngineIdleFault, AnalyticsIsEngineIdleFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.IdleSpeedFault, AnalyticsIdleSpeedFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.IdleErrorFault, AnalyticsIdleErrorFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.IdleHotFault, AnalyticsIdleHotFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.IdleBaseFault, AnalyticsIdleBaseFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.IsCruising, AnalyticsIsCruising, LEDStatus);
    setStatusLED(memsreader.memsdata.Analytics.IsClosedLoop, AnalyticsIsClosedLoop, LEDStatus);
    setStatusLED(memsreader.memsdata.Analytics.IsClosedLoopExpected, AnalyticsIsClosedLoopExpected, LEDInfo);
    setStatusLED(memsreader.memsdata.Analytics.ClosedLoopFault, AnalyticsClosedLoopFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.IsThrottleActive, AnalyticsIsThrottleActive, LEDStatus);
    setStatusLED(memsreader.memsdata.Analytics.MapFault, AnalyticsMapFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.VacuumFault, AnalyticsVacuumFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.IdleAirControlFault, AnalyticsIdleAirControlFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.IdleAirControlRangeFault, AnalyticsIdleAirControlRangeFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.IdleAirControlJackFaul, AnalyticsIdleAirControlJackFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.O2SystemFault, AnalyticsO2SystemFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.LambdaRangeFault, AnalyticsLambdaRangeFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.LambdaOscillationFault, AnalyticsLambdaOscillationFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.ThermostatFault, AnalyticsThermostatFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.CoolantTempSensorFault, AnalyticsCoolantTempSensorFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.IntakeAirTempSensorFault, AnalyticsIntakeAirTempSensorFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.FuelPumpCircuitFault, AnalyticsFuelPumpCircuitFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.ThrottlePotCircuitFault, AnalyticsThrottlePotCircuitFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.CrankshaftSensorFault, AnalyticsCrankshaftSensorFault, LEDFault);
    setStatusLED(memsreader.memsdata.Analytics.CoilFault, AnalyticsCoilFault, LEDFault);

    setFaultStatusOnMenu(data, derived);
}

function setFaultStatusOnMenu(data, derived = 0) {
    var count = 0

    if (data.CoolantTempSensorFault == true) count++;
    if (data.AirIntakeTempSensorFault == true) count++;
    if (data.ThrottlePotCircuitFault == true) count++;
    if (data.FuelPumpCircuitFault == true) count++;
    if (memsreader.memsdata.Analytics.O2SystemFault == true) count++;

    count = count + derived;

    if (count > 0) {
        setStatusLED(true, IndicatorECUFault, LEDFault);
        $("#ecu-fault-status").html(count.toString());
    } else {
        setStatusLED(false, IndicatorECUFault, LEDFault);
        $("#ecu-fault-status").html('');
    }
}

function setStatusLED(status, id, statustype = LEDStatus) {
    led = "green";

    if (statustype == LEDWarning) led = "yellow";
    if (statustype == LEDInfo) led = "blue";
    if (statustype == LEDFault) led = "red";

    console.debug("setting status led " + id + " : " + status);

    if (status == true) {
        c = "led-" + led;
    } else {
        c = "led-" + led + "-off";
    }

    id = "#" + id;
    $(id).removeClass("led-green");
    $(id).removeClass("led-red");
    $(id).removeClass("led-yellow");
    $(id).removeClass("led-blue");
    $(id).removeClass("led-green-off");
    $(id).removeClass("led-red-off");
    $(id).removeClass("led-yellow-off");
    $(id).removeClass("led-blue-off");
    $(id).removeClass("led-" + led);
    $(id).removeClass("led-" + led + "-off");
    $(id).addClass(c);
}

function setTooltip(id, message) {
    id = "#" + id;
    $(id).tooltip({title: message});
}

function updateDashboardAnalytics() {
    if (memsreader.memsdata.Analytics.IsEngineRunning) {
        setStatusLED(true, DashboardEngineRunning, LEDStatus);
        setTooltip(AnalyticsCrankshaftSensorFault, "Engine is running")
    }
    if (memsreader.memsdata.Analytics.CrankshaftSensorFault) {
        setStatusLED(true, DashboardCrankshaftSensorFault, LEDFault);
        setTooltip(AnalyticsCrankshaftSensorFault, "Crankshaft Sensor Fault, unable to start engine")
    }
    if (memsreader.memsdata.Analytics.MapFault) {
        setStatusLED(true, DashboardMapFault, LEDFault);
        setTooltip(AnalyticsCrankshaftSensorFault, "MAP Sensor Fault detected, check the vacuum pipes")
    }
}

function setConnectButtonStyle(name, style, f) {
    id = "#connectECUbtn";

    // remove all styles and handlers
    $(id).removeClass("btn-success");
    $(id).removeClass("btn-info");
    $(id).removeClass("btn-warning");
    $(id).removeClass("btn-outline-success");
    $(id).removeClass("btn-outline-info");
    $(id).removeClass("btn-outline-warning");
    $(id).removeClass("flashing-button");
    // assign new ones
    $(id).addClass(style);
    $(id).html(name);

    $(id).off().click(f);
}

function hideDebugValues() {
    console.debug("hiding debug elements")
    for (let el of document.querySelectorAll('.debug')) el.style.display = 'none';
}

function showProgressValues(show) {
    console.debug("hiding/showing progress elements")
    if (show) {
        d = 'visible'
    } else {
        d = 'hidden'
    }

    for (let el of document.querySelectorAll('.progressdisplay')) el.style.visibility = d;
}

function getAvailableSerialPorts() {
    console.info("requesting available serial ports")
    // Create a request variable and assign a new XMLHttpRequest object to it.
    var request = new XMLHttpRequest()

    request.open('GET', memsreader.uri.ports, true)
    request.setRequestHeader("Content-Type", "application/json;charset=UTF-8")
    request.addEventListener('load', updateAvailableSerialPorts)
    request.addEventListener('error', restError)

    // Send request
    request.send()
}

function updateAvailableSerialPorts(event) {
    var data = JSON.parse(this.response)
    var availablePorts = []

    // add current port
    availablePorts.push(memsreader.status.port)

    data.ports.forEach(function (p) {
        if (availablePorts.indexOf(p) === -1) {
            availablePorts.push(p);
        }
    });

    console.info("available serial ports " + JSON.stringify(availablePorts))

    $("#ports").empty()
    $.each(availablePorts, function(key, value) {
        console.info("serial port added " + key + " : " + value);
        $("#ports").append('<a class="dropdown-item" href="#" onclick="selectPort(this)">' + value + '</a>');
    });
}

function setSerialPortSelection(ports) {

    $.each(ports, function(key, value) {
        console.info("serial port added " + key + " : " + value);
        $("#ports").append('<a class="dropdown-item" href="#" onclick="selectPort(this)">' + value + '</a>');
    });
}

function selectPort(item) {
    console.info('selected ' + item.text)
    setPort(item.text)
}

function setLogToFile(logsetting, logfolder) {
    if (logsetting != LogToFileDisabled) {
        $("#logtofile").attr("checked", true);
    } else {
        $("#logtofile").attr("checked", false);
    }

    document.getElementById(SettingLogFolder).value = logfolder;
}

function setPort(port) {
    memsreader.status.port = port
    document.getElementById(SettingPort).value = port;
}

// request the config
function readConfig() {
    console.info("requesting configuration")
    // Create a request variable and assign a new XMLHttpRequest object to it.
    var request = new XMLHttpRequest()

    // Open a new connection, using the GET request on the URL endpoint
    request.open('GET', memsreader.uri.config, true)
    request.setRequestHeader("Content-Type", "application/json;charset=UTF-8")

    request.onload = function () {
        // Begin accessing JSON data here
        var data = JSON.parse(this.response)
        console.info("config request response " + JSON.stringify(data))
        updateConfigSettings(data)
    }

    request.onerror = function (e) {
        console.info("config request failed " + e)
    }

    // Send request
    request.send()
}

function updateECUQueryIntervalLabel(value) {
    $("#ecuqueryfrequencyvalue").html(value);
}

function updateConfigSettings(data) {
    console.info("Version " + data.Version + " Build " + data.Build)
    $("#version").text("Version " + data.Version + " Build " + data.Build)

    setPort(data.Port);
    setSerialPortSelection(data.Ports);
    setLogToFile(data.LogToFile, data.LogFolder);
    setECUQueryFrequency(data.Frequency);
    updateECUQueryIntervalLabel(data.Frequency);

    if (data.Debug == "true") {
        debug = data.Debug
    } else {
        hideDebugValues()
    }
}

// save the configuration settings
function Save() {
    folder = document.getElementById(SettingLogFolder).value;
    configPort = document.getElementById(SettingPort).value;
    setECUQueryFrequency(document.getElementById(SettingECUQueryFrequency).value)

    if (document.getElementById(SettingLogToFile).checked == true) {
        logToFile = LogToFileEnabled;
    } else {
        logToFile = LogToFileDisabled;
    }

    var data = { Port: configPort, logFolder: folder, logtofile: logToFile, frequency: ECUQueryInterval.toString() };

    // Create a request variable and assign a new XMLHttpRequest object to it.
    var request = new XMLHttpRequest()
    var url = uri + "/config"

    // Open a new connection, using the GET request on the URL endpoint
    request.open('PUT', url, true)
    request.setRequestHeader("Content-Type", "application/json;charset=UTF-8")

    request.onload = function () {
        console.info("post request successful (" + url + ")")
    }

    request.onerror = function () {
        console.warn("post request failed (" + url + ")")
    }

    // Send request
    request.send(JSON.stringify(data))
}

//
// Scenarios
//

// get the list of available scenarios
function loadScenarios() {
    console.info('load scenarios')

    // Create a request variable and assign a new XMLHttpRequest object to it.
    var request = new XMLHttpRequest()

    // Open a new connection, using the GET request on the URL endpoint
    request.open('GET', memsreader.uri.scenario, true)
    request.setRequestHeader("Content-Type", "application/json;charset=UTF-8")
    request.addEventListener('load', loadedScenarios)
    request.addEventListener('error', restError)

    // Send request
    request.send()
}

function loadedScenarios(event) {
    var data = JSON.parse(event.target.response)
    console.info("loaded scenarios " + JSON.stringify(data))

    // add scenarios to dropdown
    var replay = $('#replayList');

    data.forEach(function (s) {
        var scenario = '<a href="#" onclick=selectScenario("' + s.name + '") id="' + s.name + '" class="scenario list-group-item list-group-item-action">' + s.name + '</a>'
        console.debug("added scenario " + scenario)
        replay.append(scenario);
    });

    // handle highlighting selected items in the list
    $('.list-group a').click(function(e) {
        e.preventDefault()

        $that = $(this);

        $that.parent().find('a').removeClass('active');
        $that.addClass('active');
    })
}

function selectScenario(scenario) {
    console.info("select scenario " + scenario)
    selectedScenario = scenario
}

function loadScenario() {
    console.info("load scenario " + selectedScenario)

    // Create a request variable and assign a new XMLHttpRequest object to it.
    var request = new XMLHttpRequest()

    // Open a new connection, using the GET request on the URL endpoint
    var url = memsreader.uri.play_scenario + "/" + selectedScenario

    request.open('GET', url, true)
    request.setRequestHeader("Content-Type", "application/json;charset=UTF-8")
    request.addEventListener('load', runScenario)
    request.addEventListener('error', restError)

    // Send request
    request.send()
}

function runScenario(event) {
    console.info('run scenario ('+ event.target.response + ')')

    var scenarioInfo = JSON.parse(event.target.response)
    console.info("run scenario " + JSON.stringify(scenarioInfo))

    replayCount = scenarioInfo.Count
    replayPosition = scenarioInfo.Position

    console.info('connecting scenario ' + selectedScenario)

    showProgressValues(true)
    updateReplayProgress()

    var request = new XMLHttpRequest()
    var data = {"port" : selectedScenario}
    var url = memsreader.uri.connect

    memsreader.status.emulated = true

    console.info(url)

    request.open('POST', url, true)
    request.setRequestHeader("Content-Type", "application/json;charset=UTF-8")
    request.addEventListener('load', connected)
    request.addEventListener('error', restError)

    // Send request
    request.send(JSON.stringify(data))
}

// update the progress of the scenario replay
function updateReplayProgress() {
    console.info("replay " + replayPosition + " of " + replayCount)

    var percentProgress = Math.round((replayPosition / replayCount) * 100)
    var percentRemaining = 100 - percentProgress

    $("#" + ReplayProgress).width(percentProgress + "%")

    if (percentProgress < 87) {
        $("#" + ReplayProgress).html("")
        $("#" + ReplayProgressRemaining).html(percentProgress + "%")
    } else {
        $("#" + ReplayProgress).html(percentProgress + "%")
        $("#" + ReplayProgressRemaining).html("")
    }

    $("#" + ReplayProgressRemaining).width(percentRemaining + "%")
}

//-------------------------------------
// Clear the event handlers
//-------------------------------------

function removeConnectEventListeners() {
    memsreader.connectButton.removeEventListener('click', connectECU)
    memsreader.connectButton.removeEventListener('click', disconnectECU)
}

function removePlayEventListeners() {
    memsreader.playPauseButton.removeEventListener('click', pauseDataframeLoop)
    memsreader.playPauseButton.removeEventListener('click', startDataframeLoop)
}

//-------------------------------------
// ECU Command Requests
//-------------------------------------

function restError() {
    memsreader.status.emulated = false
    console.warn("post request failed (" + self.uri.connect + ")")
}

// Connect to the ECU
function connectECU() {
    console.info('connecting to ecu')

    // Create a request variable and assign a new XMLHttpRequest object to it.
    var request = new XMLHttpRequest()
    var data = memsreader.port()

    // Open a new connection, using the GET request on the URL endpoint
    request.open('POST', memsreader.uri.connect, true)
    request.setRequestHeader("Content-Type", "application/json;charset=UTF-8")
    request.addEventListener('load', connected)
    request.addEventListener('error', restError)

    // Send request
    request.send(JSON.stringify(data))
}

function connected(event) {
    var responseData = JSON.parse(event.target.response)
    console.info("connected to ecu (" + JSON.stringify(responseData) + ")")

    memsreader.status.connected = responseData.Initialised

    updateConnectMessage()
    setStatusLED(memsreader.status.connected, IndicatorECUConnected, LEDStatus)

    if (memsreader.status.connected) {
        // change the connect button action to disconnect
        if (memsreader.status.emulated) {
            updateButton("#connectECUbtn", "<i class='fa fa-power-off'></i>Stop", "btn-danger")
        } else {
            updateButton("#connectECUbtn", "<i class='fa fa-power-off'></i>Disconnect", "btn-danger")
        }

        removeConnectEventListeners()
        memsreader.connectButton.addEventListener('click', disconnectECU, {once: true})

        // set play button to pause the dataframe loop
        updateButton("#playPauseECUbtn", "<i class='fa fa-pause-circle'></i>Pause", "btn-warning")

        removePlayEventListeners()
        memsreader.playPauseButton.addEventListener('click', pauseDataframeLoop, {once: true})
        memsreader.playPauseButton.disabled = false

        // disable replay
        memsreader.replayButton.disabled = true

        setStatusLED(true, IndicatorECUConnected, LEDStatus);

        // start the dataframe command loop
        startDataframeLoop();
    } else {
        setStatusLED(true, IndicatorECUConnected, LEDFault);
        removeConnectEventListeners()
        memsreader.connectButton.addEventListener('click', connectECU, {once: true})
    }
}

function disconnectECU() {
    console.info('disconnecting from ecu')

    var request = new XMLHttpRequest()
    request.open('POST', memsreader.uri.disconnect, true)
    request.setRequestHeader("Content-Type", "application/json;charset=UTF-8")
    request.addEventListener('load', disconnected)
    request.addEventListener('error', restError)

    // Send request
    request.send()
}

function disconnected(event) {
    stopDataframeLoop()

    var responseData = JSON.parse(event.target.response)
    console.info("disconnected from the ecu (" + JSON.stringify(responseData) + ")")
    memsreader.status.connected = false
    setStatusLED(memsreader.status.connected, IndicatorECUConnected, LEDStatus)
    showProgressValues(false)
    //memsreader.status.emulated = false

    // disable the play / pause button
    updateButton("#playPauseECUbtn", "<i class='fa fa-play-circle'></i>Play", "btn-outline-success")
    memsreader.playPauseButton.disabled = true

    // enable replay
    memsreader.replayButton.disabled = false

    // update the connect button
    updateButton("#connectECUbtn", "<i class='fa fa-plug'></i>Connect", "btn-outline-success")
    clearConnectMessage()

    memsreader.status.emulated = false
    removeConnectEventListeners()
    memsreader.connectButton.addEventListener('click', connectECU, {once: true})

    resetInterface()
}

// startDataframeLoop configures a timer interval to make
// a call to retrieve the ECU dataframe
function startDataframeLoop() {
    console.info('start dataframe loop')
    memsreader.status.paused = false
    updateButton("#playPauseECUbtn", "<i class='fa fa-pause-circle'></i>Pause", "btn-warning")

    removePlayEventListeners()
    memsreader.playPauseButton.addEventListener('click', pauseDataframeLoop, {once: true})

    // reset interval
    clearInterval(dataframeLoop)
    dataframeLoop = setInterval(getDataframe, ECUQueryInterval)
}

// stop the interval timer when paused
function stopDataframeLoop() {
    console.info('stop dataframe loop')
    memsreader.status.paused = true
    memsreader.status.emulated = false
    clearInterval(dataframeLoop)
}

// Pause the Data Loop
function pauseDataframeLoop() {
    console.debug('pause dataframe loop')
    memsreader.status.paused = true
    updateButton("#playPauseECUbtn", "<i class='fa fa-play-circle'></i>Resume", "btn-success flashing-button")

    removePlayEventListeners()
    memsreader.playPauseButton.addEventListener('click', startDataframeLoop, {once: true})

    // set dataframe loop to send heartbeats
    clearInterval(dataframeLoop)
    dataframeLoop = setInterval(sendHeartbeat, ECUHeartbeatInterval)
}

// make a request for a Dataframe from the ECU
function getDataframe() {
    console.info('ecu dataframe')
    var request = new XMLHttpRequest()

    request.open('GET', memsreader.uri.dataframe, true)
    request.setRequestHeader("Content-Type", "application/json;charset=UTF-8")
    request.addEventListener('load', dataframeReceived)
    request.addEventListener('error', restError)

    // Send request
    request.send()
}

function dataframeReceived(event) {
    var data = JSON.parse(event.target.response)
    console.debug("dataframe request response " + JSON.stringify(data))

    // if the engine rpm or lambda status are unfeasibly wrong then
    // the data is corrupt
    if (data.LambdaStatus > 1 || data.EngineRPM > 7000) {
        console.error("exception dataframe is invalid!")
        console.error("memsdata " + JSON.stringify(data))
    } else {
        updateECUDataframe(data)
    }
}

function updateECUDataframe(data) {
    memsreader.memsdata = data
    console.info("memsdata " + JSON.stringify(memsreader.memsdata))

    updateGauges(data);
    updateLEDs(data);
    updateGraphs(data);
    updateDataFrameValues(data);
    updateAdjustmentValues(data);
    updateDashboardAnalytics()

    if (memsreader.status.emulated) {
        // increment the replay progress
        replayPosition = replayPosition + 1
        // loop back to the start
        if (replayPosition > replayCount)
            replayPosition = 1
        // update progress display
        updateReplayProgress();
    }
}

function updateDataFrameValues(Responsedata) {
    Object.entries(Responsedata).forEach((entry) => {
        let key = entry[0];
        let value = entry[1];
        updateDataFrameValue(key, value);
    });
}

function updateDataFrameValue(metric, data) {
    if (typeof data == "boolean") {
        data = data.toString();
    }

    $("td#" + metric + ".raw").html(data);
}

function updateAdjustmentValues(Responsedata) {
    updateAdjustmentValue(AdjustmentIdleSpeed, Responsedata.IdleSpeedOffset);
    updateAdjustmentValue(AdjustmentIdleDecay, Responsedata.IdleHot);
    updateAdjustmentValue(AdjustmentIgnitionAdvance, Responsedata.IgnitionAdvance);
    updateAdjustmentValue(AdjustmentSTFT, Responsedata.ShortTermFuelTrim);
    updateAdjustmentValue(AdjustmentLTFT, Responsedata.LongTermFuelTrim);
    updateAdjustmentValue(AdjustmentIAC, Responsedata.IACPosition);
}

function sendHeartbeat() {
    if (memsreader.status.heartbeatActive) {
        console.info('ecu heartbeat')
        var request = new XMLHttpRequest()

        // send heartbeat
        request.open('POST', memsreader.uri.heartbeat, true)
        request.setRequestHeader("Content-Type", "application/json;charset=UTF-8")
        request.addEventListener('load', heartbeatReceived)
        request.addEventListener('error', restError)

        // Send request
        request.send()
    }
}

function heartbeatReceived(event) {
    console.debug("heartbeat received")
}

function increase(id) {
    console.info('increase adjustable value ' + id)

    var data = { "Steps": 1}
    var url = memsreader.uri.adjust + id
    sendAdjustment(url, data)
}

function decrease(id) {
    console.info('decrease adjustable value ' + id)

    var data = { "Steps": -1}
    var url = memsreader.uri.adjust + id
    sendAdjustment(url, data)
}

function sendAdjustment(url, data) {
    if (memsreader.status.connected) {
        console.info('sending adjustment ' + JSON.stringify(data) + ' to ' + url)

        // Open a new connection, using the GET request on the URL endpoint
        var request = new XMLHttpRequest()
        request.open('POST', url, true)
        request.setRequestHeader("Content-Type", "application/json;charset=UTF-8")
        request.addEventListener('load', adjustmentComplete)
        request.addEventListener('error', restError)

        // Send request
        request.send(JSON.stringify(data))
    } else {
        console.warn('not connected, unable to send adjustment')
    }
}

function adjustmentComplete(event) {
    var data = JSON.parse(event.target.response)

    console.info("adjusting " + data.adjustment + " to " + data.value)
    updateAdjustmentValue(data.adjustment, data.value)
}

function updateAdjustmentValue(id, value) {
    console.info("updating " + id + " to new value " + value.toString())

    switch (id) {
        case AdjustmentSTFT: gaugeAdaptiveSTFT.value = value;
            break;
        case AdjustmentLTFT: gaugeAdaptiveLTFT.value = value;
            break;
        case AdjustmentIAC: gaugeAdaptiveIACPos.value = value;
            break;
        case AdjustmentIgnitionAdvance: gaugeAdaptiveIgnition.value = value;
            break;
        case AdjustmentIdleDecay: gaugeAdaptiveIdleDecay.value = value;
            break;
        case AdjustmentIdleSpeed: gaugeAdaptiveIdleSpeed.value = value;
            break;
    }

    // update slider
    $("input#" + id + ".range-slider__range").val(value);
    $("span#" + id + ".range-slider__value").html(value.toString());
}

function resetECU() {
    // reset ECU command
    console.info('reset ecu')
}

function resetAdj() {
    // reset ECU adjustments
    console.info('reset adjustable values')
}

function clearFaultCodes() {
    // clear fault codes
    console.info('reset fault codes')
}

function activateActuator(event) {
    if (memsreader.status.connected) {
        console.info('actuator ' + event.id + ' activate ' + event.checked)

        // Open a new connection, using the GET request on the URL endpoint
        var request = new XMLHttpRequest()
        var data = {'Activate': event.checked}
        var url = memsreader.uri.actuator + event.id

        request.open('POST', url, true)
        request.setRequestHeader("Content-Type", "application/json;charset=UTF-8")
        request.addEventListener('load', actuatorComplete)
        request.addEventListener('error', restError)

        // Send request
        request.send(JSON.stringify(data))
    } else {
        if (event.checked) {
            console.warn('not connected, unable to activate actuator')
            $('#' + event.id).bootstrapToggle('off')
        }
    }
}

async function actuatorComplete(event) {
    var response = JSON.parse(event.target.response)
    console.info("actuator response " + JSON.stringify(response))

    // if active, sleep for 2 seconds and then deactivate
    if (response.activate) {
        await sleep(2000)
        console.info('deactivating ' + response.actuator)
        $("#" + response.actuator).bootstrapToggle('off')
    }
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

function updateButton(id, name, style) {
    $(id).removeClass("btn-success");
    $(id).removeClass("btn-info");
    $(id).removeClass("btn-warning");
    $(id).removeClass("btn-danger");
    $(id).removeClass("btn-outline-success");
    $(id).removeClass("btn-outline-danger");
    $(id).removeClass("btn-outline-info");
    $(id).removeClass("btn-outline-warning");
    $(id).removeClass("flashing-button");
    // assign new ones
    $(id).addClass(style);
    $(id).html(name);
}

function updateConnectMessage() {
    var id = IndicatorConnectionMessage

    $('#' + id).removeClass("alert-light");
    $('#' + id).removeClass("alert-danger");
    $('#' + id).removeClass("alert-success");
    $('#' + id).removeClass("invisible");
    $('#' + id).addClass("visible");

    var port = memsreader.port().port
    console.info('connect message ' + memsreader.status + ' : ' + port)

    if (memsreader.status.connected) {
        if (memsreader.status.emulated) {
            document.getElementById(id).textContent = "replaying " + port
        } else {
            document.getElementById(id).textContent = "connected to " + port
        }

        $('#' + id).addClass("alert-success");
    } else {
        document.getElementById(id).textContent = "unable to connect to ECU, check connection and settings"
        $('#' + id).addClass("alert-danger");
    }
}

function clearConnectMessage() {
    var id = IndicatorConnectionMessage
    $('#' + id).removeClass("visible");
    $('#' + id).addClass("invisible");
}
