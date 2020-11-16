function updateGauges(Responsedata) {
    gaugeRPM.value = Responsedata.EngineRPM;
    gaugeMap.value = Responsedata.ManifoldAbsolutePressure;
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

function updateFaultLEDs(data) {
    var derived = derivedFaultCount;

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

    derivedFaultCount = derived

    setStatusLED(data.LambdaStatus == 0, IndicatorO2SystemFault, LEDFault);
    setStatusLED(data.Uk7d03 == 1, IndicatorRPMSensor, LEDWarning);
    setStatusLED(minLambda, IndicatorLambdaLow, LEDWarning);
    setStatusLED(maxLambda, IndicatorLambdaHigh, LEDWarning);
    setStatusLED(minIAC, IndicatorIACLow, LEDWarning);

    setFaultCount(data);
    setFaultStatusOnMenu();
}

function setFaultCount(data) {
    var count = 0

    if (data.CoolantTempSensorFault == true) count++;
    if (data.AirIntakeTempSensorFault == true) count++;
    if (data.ThrottlePotCircuitFault == true) count++;
    if (data.FuelPumpCircuitFault == true) count++;
    if (data.LambdaStatus == 0) count++;

    faultCount = count;
}

function setFaultStatusOnMenu() {
    count = faultCount + derivedFaultCount + diagnosticsFaultCount;

    if (count > 0) {
        setStatusLED(true, IndicatorECUFault, LEDFault);
        $("#ecu-fault-status").html(count.toString());
    } else {
        setStatusLED(false, IndicatorECUFault, LEDFault);
        $("#ecu-fault-status").html('');
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

function updateVersionLabel(version, build) {
    $("li#version").html("Version " + version + "<br/>Build " + build);
}