var sock = null;
var minLambda = false;
var maxLambda = false;
var minIAC = false;
var dataframeLoop;
var debug = false;
var replay = "";

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

    sock.onerror = function(error) {
        alert(`[error] ${error.message}`);
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
    gaugeLTFuelTrim.draw();
    gaugeAirFuel.draw();
    gaugeIgnition.draw();

    // create gauge sparklines
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

    // create the profiling line charts
    rpmChart = createChart(ChartRPM, "Engine (RPM)", 850, 1200);
    lambdaChart = createChart(ChartLambda, "Lambda Voltage (mV)");
    loopChart = createChart(ChartLoopIndicator, "Loop Indicator (0 Closed, 1 Open)");
    afrChart = createChart(ChartAFR, "Air : Fuel Ratio");
    coolantChart = createChart(ChartCoolant, "Coolant Temp (°C)", 80, 105);

    // load the available scenarios
    updateScenarios();

    // wire the connect button to the relevant function
    // we have to do this in javascript, so we can change the onclick
    // event function programmatically
    $("#connectECUbtn").click(this.connectECU);
    $("#replayECUbtn").click(this.replayScenario);

    showProgressValues(false)
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
        setECUQueryFrequency(data.Frequency);

        if (data.Debug == "true") {
            debug = data.Debug
        } else {
            hideDebugValues()
        }
    }

    // connection status message received
    if (msg.action == WebActionConnection) {
        connected = data.Connnected & data.Initialised;
        updateConnected(data.Initialised);
    }

    // response received from a command sent to the ECU
    if (msg.action == WebActionResponse) {
        parseECUResponse(data)
    }

    // new data received from the ECU, update the
    // gauges, graphs and status indicators 
    if (msg.action == WebActionData) {
        console.log("Dataframe --> " + msg.data);

        updateGauges(data);
        updateLEDs(data);
        updateGraphs(data);
        updateDataFrameValues(data);
        updateAdjustmentValues(data);

        if (replay != "") {
            // increment the replay progress
            replayPosition = replayPosition + 1
                // loop back to the start
            if (replayPosition > replayCount)
                replayPosition = 1
                // update progress display
            updateReplayProgress();
        }
    }

    if (msg.action == WebActionDiagnostics) {
        //waitingForResponse = false;
        console.log(data);
    }
}

function parseECUResponse(response) {
    var cmd = response.slice(0, 2)
    var value = response.slice(2, )
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
    // no throttle = 0.6V - full throttle = ~5V
    //gaugeThrottlePos.value = (Responsedata.ThrottlePotSensor - 0.6) * 22.72;
    gaugeThrottlePos.value = (Responsedata.ThrottlePotSensor) * 20;
    gaugeIACPos.value = Responsedata.IACPosition;
    gaugeBattery.value = Responsedata.BatteryVoltage;
    gaugeCoolant.value = Responsedata.CoolantTemp;
    gaugeAir.value = Responsedata.IntakeAirTemp;
    gaugeLambda.value = Responsedata.LambdaVoltage;
    gaugeFuelTrim.value = Responsedata.FuelTrimCorrection;
    gaugeLTFuelTrim.value = Responsedata.LongTermFuelTrim;
    gaugeAirFuel.value = Responsedata.AirFuelRatio;
    gaugeIgnition.value = Responsedata.IgnitionAdvance;
}

function updateGraphs(Responsedata) {
    addData(rpmSpark, Responsedata.Time, Responsedata.EngineRPM);
    addData(mapSpark, Responsedata.Time, Responsedata.ManifoldAbsolutePressure);
    addData(throttleSpark, Responsedata.Time, (Responsedata.ThrottlePotSensor - 0.6) * 22.72);
    addData(iacSpark, Responsedata.Time, Responsedata.IACPosition);
    addData(batterySpark, Responsedata.Time, Responsedata.BatteryVoltage);
    addData(coolantSpark, Responsedata.Time, Responsedata.CoolantTemp);
    addData(airSpark, Responsedata.Time, Responsedata.IntakeAirTemp);
    addData(lambdaSpark, Responsedata.Time, Responsedata.LambdaVoltage);
    addData(fuelSpark, Responsedata.Time, Responsedata.FuelTrimCorrection);
    addData(ltfuelSpark, Responsedata.Time, Responsedata.LongTermFuelTrim);
    addData(airfuelSpark, Responsedata.Time, Responsedata.AirFuelRatio);
    addData(ignitionSpark, Responsedata.Time, Responsedata.IgnitionAdvance);

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
        // disable replay once connected
        $('#replayECUbtn').prop("disabled", true);
        setStatusLED(true, IndicatorECUConnected, LEDStatus);

        // change the button operation to pause the data loop
        setConnectButtonStyle(
            "<i class='fa fa-pause-circle'>&nbsp</i>Pause",
            "btn-outline-info",
            pauseECUDataLoop
        );

        // update the IAC start postion

        // update the ECUID

        // start the dataframe command loop
        startDataframeLoop();
    } else {
        setStatusLED(true, IndicatorECUConnected, LEDFault);

        // enable connect button
        setConnectButtonStyle("<i class='fa fa-plug'>&nbsp</i>Connect", "btn-outline-success", connectECU);
        $("#connectECUbtn").prop("disabled", false);
        $('#connectECUbtn').removeClass("flashing-button");
    }
}

// calls the resp api to get the list of available scenarios
function updateScenarios() {
    uri = window.location.href.split("/").slice(0, 3).join("/");

    // Create a request variable and assign a new XMLHttpRequest object to it.
    var request = new XMLHttpRequest()

    // Open a new connection, using the GET request on the URL endpoint
    request.open('GET', uri + '/scenario', true)

    request.onload = function() {
        // Begin accessing JSON data here
        var data = JSON.parse(this.response)

        var replay = $('#replayScenarios');
        $.each(data, function(val, text) {
            var i = $('<button class="dropdown-item replay" type="button"></button>').val(text).html(text)
            replay.append(i);
        });
    }

    $("#replayScenarios").off().click(replaySelectedScenario);
    // Send request
    request.send()
}

// set up the selected scenario for replay
// update the replay and connect button visuals
function replaySelectedScenario(e) {
    e = e || window.event;
    var targ = e.target || e.srcElement || e;
    if (targ.nodeType == 3) targ = targ.parentNode;

    // extract the filename from the selected item
    // replay is global and if set to a value indicates we're replaying a scenario
    replay = targ.value;

    // request replay 
    var msg = formatSocketMessage(WebActionReplay, replay);
    sendSocketMessage(msg);

    $('#replayECUbtn').removeClass("btn-outline-info");
    $('#replayECUbtn').removeClass("btn-success");
    $('#connectECUbtn').removeClass("flashing-button");
    $('#replayECUbtn').addClass("btn-success");
    $('#replayECUbtn').prop("disabled", true);

    // show the connect button as "Play" and flash the button
    setConnectButtonStyle("<i class='fa fa-play-circle'>&nbsp</i>Play", "btn-outline-success", connectECU);
    $('#connectECUbtn').addClass("flashing-button");

    replayCount = 0
    replayPosition = 0
        // show the replay progress bar
    showProgressValues(true)
        // get the replay scenario details
    getReplayScenarioDescription(replay)
}

// call the rest api to get the description of the scenario selected
function getReplayScenarioDescription(scenario) {
    uri = window.location.href.split("/").slice(0, 3).join("/");

    // Create a request variable and assign a new XMLHttpRequest object to it.
    var request = new XMLHttpRequest()

    // Open a new connection, using the GET request on the URL endpoint
    request.open('GET', uri + '/scenario/' + scenario, true)

    request.onload = function() {
        // Begin accessing JSON data here
        var data = JSON.parse(this.response)
        console.info("replay scenario description " + JSON.stringify(data))
        replayCount = data.count
        replayPosition = data.position
        updateReplayProgress()
    }

    // Send request
    request.send()
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
    console.log("freq " + frequency)
    f = parseInt(frequency)
    if (f > 200) {
        ECUQueryInterval = f
        updateAdjustmentValue(SettingECUQueryFrequency, ECUQueryInterval)
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
    var msg = formatSocketMessage(WebActionSave, JSON.stringify(data));

    sendSocketMessage(msg);
}

function updateLEDs(data) {
    var derived = 0;

    if (data.DTC0 > 0 || data.DTC1 > 0) {
        setStatusLED(data.CoolantTempSensorFault, IndicatorCoolantFault, LEDFault);
        setStatusLED(data.AirIntakeTempSensorFault, IndicatorAirFault, LEDFault);
        setStatusLED(data.ThrottlePotCircuitFault, IndicatorThrottleFault, LEDFault);
        setStatusLED(data.FuelPumpCircuitFault, IndicatorFuelFault, LEDFault);
    }

    setStatusLED(data.ClosedLoop, IndicatorClosedLoop, LEDStatus);
    setStatusLED(data.IdleSwitch, IndicatorIdleSwitch, LEDStatus);
    setStatusLED(data.ParkNeutralSwitch, IndicatorParkSwitch, LEDStatus);

    // derived warnings
    if (data.IACPosition == 0 && data.IdleError >= 50 && data.IdleSwitch == false) {
        minIAC = true;
        derived++;
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
            derived++;
        }

        if (data.LambdaVoltage <= 10) {
            minLambda = true;
            derived++;
        }

        // evalute if a high lambda voltage has occured
        // if this has happened before trigger a fault indicator
        // this must be evaluated before we set the maxLamda warning to ensure
        // we have at least one occurence first
        if (maxLambda && data.LambdaVoltage >= 900) {
            setStatusLED(true, IndicatorLambdaHighFault, LEDFault);
            derived++;
        }

        if (data.LambdaVoltage >= 900) {
            maxLambda = true;
            derived++;
        }
    }

    setStatusLED(data.LambdaStatus == 0, IndicatorO2SystemFault, LEDFault);
    setStatusLED(data.Uk7d03 == 1, IndicatorRPMSensor, LEDWarning);
    setStatusLED(minLambda, IndicatorLambdaLow, LEDWarning);
    setStatusLED(maxLambda, IndicatorLambdaHigh, LEDWarning);
    setStatusLED(minIAC, IndicatorIACLow, LEDWarning);

    setFaultStatusOnMenu(data, derived);
}

function setFaultStatusOnMenu(data, derived = 0) {
    var count = 0

    if (data.CoolantTempSensorFault == true) count++;
    if (data.AirIntakeTempSensorFault == true) count++;
    if (data.ThrottlePotCircuitFault == true) count++;
    if (data.FuelPumpCircuitFault == true) count++;
    if (data.LambdaStatus == 0) count++;

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
    $(id).removeClass("flashing-button");
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
    console.log("updating " + id + " to new value " + value.toString())

    $("input#" + id + ".range-slider__range").val(value);
    $("span#" + id + ".range-slider__value").html(value.toString());
}

function hideDebugValues() {
    console.log("hiding debug elements")
    for (let el of document.querySelectorAll('.debug')) el.style.display = 'none';
}

function showProgressValues(show) {
    console.log("hiding/showing progress elements")
    if (show) {
        d = 'visible'
    } else {
        d = 'hidden'
    }

    for (let el of document.querySelectorAll('.progressdisplay')) el.style.visibility = d;
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
    waitingForResponseTimeout = setInterval(waitForResponseTimedOut, WaitForResponseInterval)
}

// called if the wait for response times out
function waitForResponseTimedOut() {
    console.error("timed out waiting for response")
    clearWaitForResponse()
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
    // reset ECU command
    var msg = formatSocketMessage(WebActionCommand, CommandResetECU);
    sendSocketMessage(msg);
}

function resetAdj() {
    // reset ECU adjustments
    var msg = formatSocketMessage(WebActionCommand, CommandResetAdjustments);
    sendSocketMessage(msg);
}

function clearFaultCodes() {
    // clear fault codes
    var msg = formatSocketMessage(WebActionCommand, CommandClearFaults);
    sendSocketMessage(msg);

}

// Pause the Data Loop
function pauseECUDataLoop() {
    var msg = formatSocketMessage(WebActionCommand, CommandPause);
    sendSocketMessage(msg);

    // change the button operation to restart the data loop
    setConnectButtonStyle(
        "<i class='fa fa-play-circle'>&nbsp</i>Restart",
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
        "<i class='fa fa-pause-circle'>&nbsp</i>Pause",
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

        sock.send(msg);
        startWaitForResponse()
    } else {
        console.warn("can't send whilst waiting for a response")
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