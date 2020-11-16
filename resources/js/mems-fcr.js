// this function gets called as soon as the page load has completed
window.onload = function () {
    // get the url of the current page to build the websocket url
    wsuri = window.location.href.split("/").slice(0, 3).join("/");
    wsuri = wsuri.replace("http:", "ws:");

    // open the websock and set up listeners for
    // open, close and message events
    sock = new WebSocket(wsuri + "/ws");

    sock.onopen = function () {
        console.log("connected to " + wsuri);
        readConfig();
    };

    sock.onclose = function (e) {
        console.log("connection closed (" + e.code + ")");
    };

    sock.onmessage = function (e) {
        console.log("message received: " + e.data);
        clearWaitForResponse()
        parseMessage(e.data);
    };

    sock.onerror = function (error) {
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
        updateVersionLabel(data.Version, data.Build)

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
        updateFaultLEDs(data);
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
        console.log(data);
        parseDiagnosticsResponse(data)
    }
}

function parseDiagnosticsResponse(data) {
    diagnosticReport = data
    diagnosticsFaultCount = 0

    console.log("parseDiagnosticsResponse : " + diagnosticReport);

    console.log("lambda fault " + diagnosticReport.LambdaFault)
    if (diagnosticReport.LambdaFault == true) {
        diagnosticsFaultCount++;
    }

    console.log("engine idle fault " + diagnosticReport.IsEngineIdleFault)
    if (diagnosticReport.IsEngineIdleFault == true) {
        diagnosticsFaultCount++;
    }

    if (diagnosticReport.ClosedLoopExpected == true) {
        diagnosticsFaultCount++;
    }

    if (diagnosticReport.MapFault == true) {
        diagnosticsFaultCount++;
    }

    if (diagnosticReport.VacuumFault == true) {
        diagnosticsFaultCount++;
    }
}

function parseECUResponse(response) {
    var cmd = response.slice(0, 2)
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