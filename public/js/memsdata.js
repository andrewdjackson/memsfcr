var sock = null;
var minLambda = false;
var maxLambda = false;
var minIAC = false;
var dataframeLoop;

var waitingForResponse = false;
var waitingForResponseTimeout;
const WaitForResponseInterval = 500

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

// settings
const SettingLogFolder = "logfolder"
const SettingLogToFile = "logtofile"
const LogToFileEnabled = "true"
const LogToFileDisabled = "false"
const SettingPort = "port"
const SettingPortList = "ports"

// duration in milliseconds between calls to the ECU for
// dataframes. the ECU will struggle to respond with a 
// value less than 450ms
const ECUQueryInterval = 900

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

// this function gets called as soon as the page load has completed
window.onload = function() {
    // get the url of the current page to build the websocket url
    wsuri = window.location.href.split("/").slice(0, 3).join("/");
    wsuri = wsuri.replace("http:", "ws:");

    // open the websock and set up listeners for
    // open, close and message events
    sock = new WebSocket(wsuri + "/ws");

    sock.onopen = function() {
        console.log("connected to " + wsuri);
        readConfig();
    };

    sock.onclose = function(e) {
        console.log("connection closed (" + e.code + ")");
    };

    sock.onmessage = function(e) {
        console.log("message received: " + e.data);
        clearWaitForResponse()
        parseMessage(e.data);
    };

    // draw the gauges
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
    gaugeIgnition.draw();

    // create the profiling line charts
    rpmChart = createChart(ChartRPM, "Engine (RPM)", 850, 1200);
    lambdaChart = createChart(ChartLambda, "Lambda Voltage (mV)");
    loopChart = createChart(ChartLoopIndicator, "Loop Indicator (0 Closed, 1 Open)");
    afrChart = createChart(ChartAFR, "Air : Fuel Ratio");
    coolantChart = createChart(ChartCoolant, "Coolant Temp (°C)", 80, 105);

    // wire the connect button to the relevant function
    // we have to do this in javascript, so we can change the onclick
    // event function programmatically
    $("#connectECUbtn").click(this.connectECU);
};

// parseMessage receives the websocket message as a json object
// in general the ECU operates in a synchronous command / response model
// as such once a command is sent, buttons are disabled until a response
// has been received. The serial interface has a timeout of a couple seconds
// so buttons may be disabled for this period of time if no response is
// received.
function parseMessage(m) {
    var msg = JSON.parse(m);
    var data = JSON.parse(msg.data);

    // config received
    if (msg.action == WebActionConfig) {
        console.log(data);
        setPort(data.Port);
        setSerialPortSelection(data.Ports);
        setLogToFile(data.LogToFile, data.LogFolder);
    }

    // connection status message received
    if (msg.action == WebActionConnection) {
        waitingForResponse = false;
        connected = data.Connnected & data.Initialised;
        updateConnected(data.Initialised);
    }

    // response received from a command sent to the ECU
    if (msg.action == WebActionResponse) {
        waitingForResponse = false
        parseECUResponse(data)
    }

    // new data received from the ECU, update the
    // gauges, graphs and status indicators 
    if (msg.action == WebActionData) {
        waitingForResponse = false;
        console.log(data);

        updateGauges(data);
        updateLEDs(data);
        updateGraphs(data);
        updateDataFrameValues(data);
        updateAdjustmentValues(data);
    }
}

function parseECUResponse(response) {
    var cmd = response.slice(0,2)
    var value = response.slice(2,)
    console.log("parsing response cmd : " + cmd + ", val : " + value)

    switch (cmd) {
        case ResponseIdleSpeedIncrement: 
        case ResponseIdleSpeedDecrement: 
            updateAdjustmentValue(AdjustmentIdleSpeed, value);
            break;
        case ResponseIgnitionAdvanceOffsetIncrement:
        case ResponseIgnitionAdvanceOffsetDecrement:
            updateAdjustmentValue(AdjustmentIgnitionAdvance, value);
            break;
        case ResponseIdleDecayIncrement:
        case ResponseIdleDecayDecrement:
            updateAdjustmentValue(AdjustmentIdleHot, value);
            break;
        case ResponseLTFTIncrement:
        case ResponseLTFTDecrement:
            updateAdjustmentValue(AdjustmentFuelTrim, value);
            break;
        case ResponseSTFTIncrement:
        case ResponseSTFTDecrement:
            updateAdjustmentValue(AdjustmentFuelTrim, value);
            break;
    }
}

function updateGauges(Responsedata) {
    gaugeRPM.value = Responsedata.EngineRPM;
    gaugeMap.value = Responsedata.ManifoldAbsolutePressure;
    gaugeThrottlePos.value = Responsedata.ThrottlePotSensor;
    gaugeIACPos.value = Responsedata.IACPosition;
    gaugeBattery.value = Responsedata.BatteryVoltage;
    gaugeCoolant.value = Responsedata.CoolantTemp;
    gaugeAir.value = Responsedata.IntakeAirTemp;
    gaugeLambda.value = Responsedata.LambdaVoltage;
    gaugeFuelTrim.value = Responsedata.FuelTrimCorrection;
    gaugeIgnition.value = Responsedata.IgnitionAdvance;
}

function updateGraphs(Responsedata) {
    addData(rpmChart, Responsedata.Time, Responsedata.EngineRPM);
    addData(lambdaChart, Responsedata.Time, Responsedata.LambdaVoltage);
    addData(loopChart, Responsedata.Time, Responsedata.ClosedLoop);
    addData(afrChart, Responsedata.Time, Responsedata.AirFuelRatio);
    addData(coolantChart, Responsedata.Time, Responsedata.CoolantTemp);
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

function updateConnected(connected) {
    console.log("connected " + connected);

    // enable all buttons
    enableAllButtons()

    setConnectionStatusMessage(connected)

    if (connected) {
        setStatusLED(true, IndicatorECUConnected, LEDStatus);

        // change the button operation to pause the data loop
        setConnectButtonStyle(
            "<i class='fa fa-pause-circle'>&nbsp</i>Pause Data Loop",
            "btn-outline-info",
            pauseECUDataLoop
        );

        // start the dataframe command loop
        startDataframeLoop();
    } else {
        setStatusLED(true, IndicatorECUConnected, LEDFault);

        // enable connect button
        setConnectButtonStyle("<i class='fa fa-plug'>&nbsp</i>Connect", "btn-outline-success", connectECU);
        $("#connectECUbtn").prop("disabled", false);
    }
}

function disableAllButtons() {
    // disable all buttons
    $(":button").prop("disabled", true);
}

function enableAllButtons() {
    // enable all buttons
    $(":button").prop("disabled", false);
}

function setConnectionStatusMessage(connected) {
    id = IndicatorConnectionMessage

    $('#' + id).removeClass("alert-light");
    $('#' + id).removeClass("alert-danger");
    $('#' + id).removeClass("alert-success");

    $('#' + id).removeClass("invisible");
    $('#' + id).addClass("visible");

    if (connected == true) {
        document.getElementById(id).textContent = "connected to " + document.getElementById("port").value
        $('#' + id).addClass("alert-success");
    } else {
        document.getElementById(id).textContent = "unable to connect to " + document.getElementById("port").value
        $('#' + id).addClass("alert-danger");
    }
}

// save the configuration settings
function Save() {
    folder = document.getElementById(SettingLogFolder).value;
    configPort = document.getElementById(SettingPort).value;

    if (document.getElementById(SettingLogToFile).checked == true) {
        logToFile = LogToFileEnabled;
    } else {
        logToFile = LogToFileDisabled;
    }

    var data = { Port: configPort, logFolder: folder, logtofile: logToFile };
    var msg = formatSocketMessage(WebActionSave, JSON.stringify(data));

    sendSocketMessage(msg);
}

function updateLEDs(data) {
    if (data.DTC0 > 0 || data.DTC1 > 0) {
        setStatusLED(true, IndicatorECUFault, LEDFault);
        setStatusLED(data.CoolantTempSensorFault, IndicatorCoolantFault, LEDFault);
        setStatusLED(data.AirIntakeTempSensorFault, IndicatorAirFault, LEDFault);
        setStatusLED(data.ThrottlePotCircuitFault, IndicatorThrottleFault, LEDFault);
        setStatusLED(data.FuelPumpCircuitFault, IndicatorFuelFault, LEDFault);
        setFaultStatusOnMenu(data);
    }

    setStatusLED(data.ClosedLoop, IndicatorClosedLoop, LEDStatus);
    setStatusLED(data.IdleSwitch, IndicatorIdleSwitch, LEDStatus);
    setStatusLED(data.ParkNeutralSwitch, IndicatorParkSwitch, LEDStatus);

    // derived warnings
    if (data.IACPosition == 0 && data.IdleError >= 50 && data.IdleSwitch == false) {
        minIAC = true;
    }

    // only evaluate lambda faults if we're in closed loop where
    // the lambda voltage has an effect
    if (data.ClosedLoop) {
        // evalute if a low lambda voltage has occured
        // if this has happened before trigger a fault indicator
        // this must be evaluated before we set the minLamda warning to ensure
        // we have at least one occurence first
        if (minLambda && data.LambdaVoltage <= 10) {
            setStatusLED(true, IndicatorLambdaLowFault, LEDFault);
        }
        if (data.LambdaVoltage <= 10) {
            minLambda = true;
        }

        // evalute if a high lambda voltage has occured
        // if this has happened before trigger a fault indicator
        // this must be evaluated before we set the maxLamda warning to ensure
        // we have at least one occurence first
        if (maxLambda && data.LambdaVoltage >= 900) {
            setStatusLED(true, IndicatorLambdaHighFault, LEDFault);
        }
        if (data.LambdaVoltage >= 900) {
            maxLambda = true;
        }
    }

    setStatusLED(data.Uk7d03 == 1, IndicatorRPMSensor, LEDWarning);
    setStatusLED(minLambda, IndicatorLambdaLow, LEDWarning);
    setStatusLED(maxLambda, IndicatorLambdaHigh, LEDWarning);
    setStatusLED(minIAC, IndicatorIACLow, LEDWarning);
}

function setFaultStatusOnMenu(data) {
    var count = 0

    if (data.CoolantTempSensorFault == true) count++
    if (data.AirIntakeTempSensorFault == true) count++
    if (data.ThrottlePotCircuitFault == true) count++
    if (data.FuelPumpCircuitFault == true) count++

    if (count > 0) {
        $("#ecu-fault-status").html(count.toString());
    } else {
        $("#ecu-fault-status").html('');
    }
}

function setStatusLED(status, id, statustype = LEDStatus) {
    led = "green";

    if (statustype == LEDWarning) led = "yellow";

    if (statustype == LEDFault) led = "red";

    console.log(id + " : " + status);

    if (status == true) {
        c = "led-" + led;
    } else {
        c = "led-" + led + "-off";
    }

    id = "#" + id;
    $(id).removeClass("led-green");
    $(id).removeClass("led-red");
    $(id).removeClass("led-yellow");
    $(id).removeClass("led-green-off");
    $(id).removeClass("led-red-off");
    $(id).removeClass("led-yellow-off");
    $(id).removeClass("led-" + led);
    $(id).removeClass("led-" + led + "-off");
    $(id).addClass(c);
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

    // assign new ones
    $(id).addClass(style);
    $(id).html(name);

    $(id).off().click(f);
}

function updateAdjustmentValues(Responsedata) {
    updateAdjustmentValue(AdjustmentIdleSpeed, Responsedata.IdleSpeedOffset);
    updateAdjustmentValue(AdjustmentIdleHot, Responsedata.IdleHot);
    updateAdjustmentValue(AdjustmentIgnitionAdvance, Responsedata.IgnitionAdvance);
    updateAdjustmentValue(AdjustmentFuelTrim, Responsedata.LongTermFuelTrim);
}

function updateAdjustmentValue(id, value) {
    console.log("updating "+ id + " to new value " + value.toString())

    $("input#" + id + ".range-slider__range").val(value);
    $("span#" + id + ".range-slider__value").html(value.toString());
}

function setSerialPortSelection(ports) {
    $.each(ports, function(key, value) {
        console.log("serial port added " + key + " : " + value);
        $("#ports").append('<a class="dropdown-item" href="#" onclick="selectPort(this)">' + value + '</a>');
    });
}

function selectPort(item) {
    console.log('selected ' + item.text)
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
    document.getElementById(SettingPort).value = port;
}

// request the config
function readConfig() {
    var msg = formatSocketMessage(WebActionConfig, CommandReadConfig);
    sendSocketMessage(msg);
}

//-------------------------------------
// ECU Command Requests 
//-------------------------------------

// Connect to the ECU
function connectECU() {
    var port = document.getElementById(SettingPort).value;
    var msg = formatSocketMessage(WebActionConnect, port);
    sendSocketMessage(msg);

    // show connecting
    setConnectButtonStyle("<i class='fa fa-plug'>&nbsp</i>Connecting..", "btn-warning", connectECU);

    // disable all buttons
    disableAllButtons()
}

// startDataframeLoop configures a timer interval to make
// a call to retieve the ECU dataframe
function startDataframeLoop() {
    dataframeLoop = setInterval(getDataframe, ECUQueryInterval);
}

// stop the interval timer when paused
function stopDataframeLoop() {
    clearInterval(dataframeLoop);
}

function startWaitForResponse() {
    waitingForResponse = true
    waitingForResponseTimeout = setInterval(clearWaitForResponse, WaitForResponseInterval)
}

// fail back if we don't get a response, so that the UI doesn't get blocked
function clearWaitForResponse() {
    waitingForResponse = false
    clearInterval(waitingForResponseTimeout);
}

// make a request for a Dataframe from the ECU
function getDataframe() {
    // if we're not waiting for a response then send the dataframe request
    var msg = formatSocketMessage(WebActionCommand, CommandDataFrame);
    sendSocketMessage(msg);
}

function increase(id) {
    // if we're not waiting for a response then send the ecu command
    var msg = formatSocketMessage(WebActionIncrease, id);
    sendSocketMessage(msg);
}

function decrease(id) {
    // if we're not waiting for a response then send the ecu command
    var msg = formatSocketMessage(WebActionDecrease, id);
    sendSocketMessage(msg);
}

function resetECU() {
    // if we're not waiting for a response then send the ecu command
    var msg = formatSocketMessage(WebActionCommand, CommandResetECU);
    sendSocketMessage(msg);
}

function resetAdj() {
    // if we're not waiting for a response then send the ecu command
    var msg = formatSocketMessage(WebActionCommand, CommandResetAdjustments);
    sendSocketMessage(msg); 
}

function clearFaultCodes() {
    // if we're not waiting for a response then send the ecu command
    var msg = formatSocketMessage(WebActionCommand, CommandClearFaults);
    sendSocketMessage(msg);

}

// Pause the Data Loop
function pauseECUDataLoop() {
    var msg = formatSocketMessage(WebActionCommand, CommandPause);
    sendSocketMessage(msg);

    // change the button operation to restart the data loop
    setConnectButtonStyle(
        "<i class='fa fa-play-circle'>&nbsp</i>Restart Data Loop",
        "btn-outline-warning",
        restartECUDataLoop
    );

    // stop the dataframe loop
    stopDataframeLoop();
}

// Restart the Data Loop
function restartECUDataLoop() {
    var msg = formatSocketMessage(WebActionCommand, CommandStart);
    sendSocketMessage(msg);

    // change the button operation back to pause the data loop
    setConnectButtonStyle(
        "<i class='fa fa-pause-circle'>&nbsp</i>Pause Data Loop",
        "btn-outline-info",
        pauseECUDataLoop
    );

    // restart the dataframe loop
    startDataframeLoop();
}

// send the formatted message over the websocket
function sendSocketMessage(msg) {
    if (!waitingForResponse) {
        console.log("sending socket message: " + msg);
        
        startWaitForResponse()
        sock.send(msg);
    } else {
        console.log("can't send whilst waiting for a response")
    }
}

// format messages to be sent over the websocket
// in json format as:
// {
//    action: '<verb>'
//    data: '<data payload'    
// }
function formatSocketMessage(a, d) {
    var msg = { action: a, data: d };
    return JSON.stringify(msg);
}
