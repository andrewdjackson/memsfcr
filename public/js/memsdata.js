var sock = null;
var wsuri = "ws://127.0.0.1:1234";
var minLambda = false
var maxLambda = false
var minIAC = false

window.onload = function() {
    sock = new WebSocket(wsuri);

    sock.onopen = function() {
        console.log("connected to " + wsuri);
    }

    sock.onclose = function(e) {
        console.log("connection closed (" + e.code + ")");
    }

    sock.onmessage = function(e) {
        console.log("message received: " + e.data);
        parseMessage(e.data)
    }

    gaugeRPM.draw()
    gaugeMap.draw()
    gaugeThrottlePos.draw()
    gaugeIACPos.draw()
    gaugeBattery.value = 11
    gaugeBattery.draw()
    gaugeCoolant.draw()
    gaugeAir.draw()
    gaugeLambda.draw()
    gaugeFuelTrim.draw()
    gaugeIgnition.draw()

    rpmChart = createChart("rpmchart", "Engine RPM", 850, 1200)
    lambdaChart = createChart("lambdachart", "Lambda Voltage (mV)")
    loopChart = createChart("loopchart", "Loop Indicator")
    coolantChart = createChart("coolantchart", "Coolant Temp (°C)", 80, 105)
};

function parseMessage(m) {
    var msg = JSON.parse(m);
    var data = JSON.parse(msg.data)

    if (msg.action == "config") {
        console.log(data)
        setPort(data.Port)
    }

    if (msg.action == "data") {
        console.log(data)
        gaugeRPM.value = parseInt(data.EngineRPM)
        gaugeMap.value = parseInt(data.ManifoldAbsolutePressure)

        var throttleposition = parseInt(data.ThrottlePotSensor / 2)
        gaugeThrottlePos.value = throttleposition

        gaugeIACPos.value = parseInt(data.IACPosition)

        gaugeBattery.value = parseInt(data.BatteryVoltage)
        gaugeCoolant.value = parseInt(data.CoolantTemp)
        gaugeAir.value = parseInt(data.IntakeAirTemp)
        gaugeLambda.value = parseInt(data.LambdaVoltage)

        var fueltrimcorrection = parseInt(data.ShortTermFuelTrim - 100)
        gaugeFuelTrim.value = fueltrimcorrection
        gaugeIgnition.value = parseFloat(data.IgnitionAdvance)

        setLEDs(data)

        addData(rpmChart, data.Time, data.EngineRPM)
        addData(lambdaChart, data.Time, data.LambdaVoltage)
        addData(loopChart, data.Time, data.ClosedLoop)
        addData(coolantChart, data.Time, data.CoolantTemp)
    }
}

function setLEDs(data) {
    if ((data.DTC0 != 0) && (data.DTC1 != 0) && (data.DTC2 != 0)) {
        setStatusLED(true, "ecufault", "fault")
        setStatusLED(data.CoolantTempSensorFault, "coolantfault", "fault")
        setStatusLED(data.AirIntakeTempSensorFault, "airfault", "fault")
        setStatusLED(data.ThrottlePotCircuitFault, "throttleault", "fault")
        setStatusLED(data.FuelPumpCircuitFault, "fuelfault", "fault")
    }

    setStatusLED(data.ClosedLoop, "closedloop", "status")
    setStatusLED((data.IdleSwitch), "idleswitch", "status")
    setStatusLED(data.ParkNeutralSwitch, "parkswitch", "status")

    // derived warnings
    if ((data.IACPosition == 0) && (data.IdleError >= 50) && (data.IdleSwitch == false) && (data.fd7d03 != 0)) {
        minIAC = true
    }

    // only evaluate lambda faults if we're in closed loop where
    // the lambda voltage has an effect
    if (data.ClosedLoop) {
        // evalute if a low lambda voltage has occured
        // if this has happened before trigger a fault indicator
        // this must be evaluated before we set the minLamda warning to ensure
        // we have at least one occurence first
        if ((minLambda) && (data.LambdaVoltage <= 10)) {
            setStatusLED(true, "lambdalowfault", "fault")
        }
        if (data.LambdaVoltage <= 10) {
            minLambda = true
        }

        // evalute if a high lambda voltage has occured
        // if this has happened before trigger a fault indicator
        // this must be evaluated before we set the maxLamda warning to ensure
        // we have at least one occurence first
        if ((maxLambda) && (data.LambdaVoltage >= 900)) {
            setStatusLED(true, "lambdahighfault", "fault")
        }
        if (data.LambdaVoltage >= 900) {
            maxLambda = true
        }
    }

    setStatusLED((data.Uk7d03 == 1), "rpmsensor", "warning")
    setStatusLED(minLambda, "lambdalow", "warning")
    setStatusLED(maxLambda, "lambdahigh", "warning")
    setStatusLED(minIAC, "iaclow", "warning")
}

function setStatusLED(status, id, statustype = "status") {
    led = "green"

    if (statustype == "warning")
        led = "yellow"

    if (statustype == "fault")
        led = "red"

    console.log(id + ' : ' + status)

    if (status == true) {
        c = "led-" + led
    } else {
        c = "led-" + led + "-off"
    }

    id = "#" + id
    $(id).removeClass('led-' + led);
    $(id).removeClass('led-' + led + '-off');
    $(id).addClass(c);
}

function connectECU() {
    var port = document.getElementById('port').value;
    var msg = formatSocketMessage('connect', port)
    sendSocketMessage(msg)
}

function sendSocketMessage(msg) {
    sock.send(msg)
}

function formatSocketMessage(action, data) {
    var msg = '{"action":"' + action + '", "data":"' + data + '"}'
    return msg
}

function setPort(port) {
    document.getElementById('port').value = port
}