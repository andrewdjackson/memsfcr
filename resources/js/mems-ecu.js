//-------------------------------------
// ECU Command Requests 
//-------------------------------------

// Connect to the ECU
function wsConnectECU() {
    var port = document.getElementById(SettingPort).value;
    var msg = formatSocketMessage(WebActionConnect, port);
    sendSocketMessage(msg);

    // show connecting
    setConnectButtonStyle("<i class='fa fa-plug'>&nbsp</i>Connecting..", "btn-warning", connectECU);

    // disable all buttons
    disableAllButtons()
}

function connectECU() {
    var port = document.getElementById(SettingPort).value
    // show connecting
    setConnectButtonStyle("<i class='fa fa-plug'>&nbsp</i>Connecting..", "btn-warning", connectECU);

    uri = window.location.href.split("/").slice(0, 3).join("/");

    // Create a request variable and assign a new XMLHttpRequest object to it.
    var request = new XMLHttpRequest()

    // Open a new connection, using the GET request on the URL endpoint
    request.open('POST', uri + '/rosco/connect', true)
    request.setRequestHeader("Content-Type", "application/json;charset=UTF-8")

    request.onload = function () {
        // Begin accessing JSON data here
        var data = JSON.parse(this.response)
        console.info("connect request response " + JSON.stringify(data))
        updateConnected(data.Initialised);
        //setConnectionStatusMessage(data.Connected)
    }

    // Send request
    request.send(JSON.stringify({ "port": port }))
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