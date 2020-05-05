var sock = null;
var minLambda = false;
var maxLambda = false;
var minIAC = false;
var dataframeLoop;

const WebActionSave = "save";
const WebActionConfig = "config";
const WebActionConnection = "connection";
const WebActionData = "data";

window.onload = function () {
	wsuri = window.location.href.split("/").slice(0, 3).join("/");
	wsuri = wsuri.replace("http:", "ws:");

	sock = new WebSocket(wsuri);

	sock.onopen = function () {
		console.log("connected to " + wsuri);
		readConfig();
	};

	sock.onclose = function (e) {
		console.log("connection closed (" + e.code + ")");
	};

	sock.onmessage = function (e) {
		console.log("message received: " + e.data);
		parseMessage(e.data);
	};

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

	rpmChart = createChart("rpmchart", "Engine RPM", 850, 1200);
	lambdaChart = createChart("lambdachart", "Lambda Voltage (mV)");
	loopChart = createChart("loopchart", "Loop Indicator");
	coolantChart = createChart("coolantchart", "Coolant Temp (°C)", 80, 105);

	$("#connectECUbtn").click(this.connectECU);
};

function parseMessage(m) {
	var msg = JSON.parse(m);
	var data = JSON.parse(msg.data);

	if (msg.action == WebActionConfig) {
		console.log(data);
		setPort(data.Port);
		setSerialPortSelection(data.Ports);
		setLogToFile(data.Output, data.LogFolder);
	}

	if (msg.action == WebActionConnection) {
		connected = data.Connnected & data.Initialised;
		updateConnected(data.Initialised);
	}

	if (msg.action == WebActionData) {
		console.log(data);

		//memsdata = computeMemsData(data)

		updateGauges(data);
		updateLEDs(data);
		updateGraphs(data);
		updateDataFrameValues(data);
		updateAdjustmentValues(data);
	}
}

function computeMemsData(memsdata) {
	var d = Object.create(memsdata);

	// Dataframe 0x7d Compute Values
	d.ThrottleAngle = ((memsdata.ThrottleAngle * 6) / 10).toFixed(2);
	d.AirFuelRatio = (memsdata.AirFuelRatio / 10).toFixed(1);
	d.LambdaVoltage = (memsdata.LambdaVoltage * 5).toFixed(1);
	d.IgnitionAdvanceOffset7d = memsdata.IgnitionAdvanceOffset7d - 48;
	d.IdleSpeedOffset = ((memsdata.IdleSpeedOffset - 128) * 25).toFixed(2);

	// Dataframe 0x80 Compute Values
	d.CoolantTemp = (memsdata.CoolantTemp - 55).toFixed(1);
	d.AmbientTemp = (memsdata.AmbientTemp - 55).toFixed(1);
	d.IntakeAirTemp = (memsdata.IntakeAirTemp - 55).toFixed(1);
	d.FuelTemp = (memsdata.FuelTemp - 55).toFixed(1);
	d.BatteryVoltage = (memsdata.BatteryVoltage / 10).toFixed(1);
	d.ThrottlePotSensor = (memsdata.ThrottlePotSensor * 0.02).toFixed(2);
	d.IACPosition = (memsdata.IACPosition / 1.8).toFixed(2);
	d.IgnitionAdvance = (memsdata.IgnitionAdvance / 2 - 24).toFixed(2);
	d.CoilTime = (memsdata.CoilTime * 0.002).toFixed(2);

	// Additional Compute Values
	d.FuelTrimCorrection = memsdata.ShortTermFuelTrim - 100;

	return d;
}

function updateGauges(memsdata) {
	gaugeRPM.value = memsdata.EngineRPM;
	gaugeMap.value = memsdata.ManifoldAbsolutePressure;
	gaugeThrottlePos.value = memsdata.ThrottlePotSensor;
	gaugeIACPos.value = memsdata.IACPosition;
	gaugeBattery.value = memsdata.BatteryVoltage;
	gaugeCoolant.value = memsdata.CoolantTemp;
	gaugeAir.value = memsdata.IntakeAirTemp;
	gaugeLambda.value = memsdata.LambdaVoltage;
	gaugeFuelTrim.value = memsdata.FuelTrimCorrection;
	gaugeIgnition.value = memsdata.IgnitionAdvance;
}

function updateGraphs(memsdata) {
	addData(rpmChart, memsdata.Time, memsdata.EngineRPM);
	addData(lambdaChart, memsdata.Time, memsdata.LambdaVoltage);
	addData(loopChart, memsdata.Time, memsdata.ClosedLoop);
	addData(coolantChart, memsdata.Time, memsdata.CoolantTemp);
}

function updateDataFrameValues(memsdata) {
	Object.entries(memsdata).forEach((entry) => {
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

	if (connected) {
		setStatusLED(true, "ecudata", "status");
		// change the button operation to pause the data loop
		setConnectButtonStyle(
			"<i class='fa fa-pause-circle'>&nbsp</i>Pause Data Loop",
			"btn-outline-info",
			pauseECUDataLoop
		);
		// enable all buttons
		$(":button").prop("disabled", false);
		// start the dataframe command loop
		startDataframeLoop();
	} else {
		setStatusLED(true, "ecudata", "fault");
		// enable connect button
		setConnectButtonStyle("<i class='fa fa-plug'>&nbsp</i>Connect", "btn-outline-success", connectECU);
		$("#connectECUbtn").prop("disabled", false);
	}
}

function Save() {
	folder = document.getElementById("logfolder").value;
	configPort = document.getElementById("port").value;

	if (document.getElementById("logtofile").checked == true) {
		logToFile = "file";
	} else {
		logToFile = "stdout";
	}

	var data = { Port: configPort, logFolder: folder, Output: logToFile };
	var msg = formatSocketMessage("save", JSON.stringify(data));

	sendSocketMessage(msg);
}

function startDataframeLoop() {
	dataframeLoop = setInterval(getDataframe, 1000);
}

function stopDataframeLoop() {
	clearInterval(dataframeLoop);
}

function getDataframe() {
	paused = false;
	if (!paused) {
		var msg = formatSocketMessage("command", "dataframe");
		sendSocketMessage(msg);
	}
}

function updateLEDs(data) {
	if (data.DTC0 != 0 && data.DTC1 != 0 && data.DTC2 != 0) {
		setStatusLED(true, "ecufault", "fault");
		setStatusLED(data.CoolantTempSensorFault, "coolantfault", "fault");
		setStatusLED(data.AirIntakeTempSensorFault, "airfault", "fault");
		setStatusLED(data.ThrottlePotCircuitFault, "throttleault", "fault");
		setStatusLED(data.FuelPumpCircuitFault, "fuelfault", "fault");
	}

	setStatusLED(data.ClosedLoop, "closedloop", "status");
	setStatusLED(data.IdleSwitch, "idleswitch", "status");
	setStatusLED(data.ParkNeutralSwitch, "parkswitch", "status");

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
			setStatusLED(true, "lambdalowfault", "fault");
		}
		if (data.LambdaVoltage <= 10) {
			minLambda = true;
		}

		// evalute if a high lambda voltage has occured
		// if this has happened before trigger a fault indicator
		// this must be evaluated before we set the maxLamda warning to ensure
		// we have at least one occurence first
		if (maxLambda && data.LambdaVoltage >= 900) {
			setStatusLED(true, "lambdahighfault", "fault");
		}
		if (data.LambdaVoltage >= 900) {
			maxLambda = true;
		}
	}

	setStatusLED(data.Uk7d03 == 1, "rpmsensor", "warning");
	setStatusLED(minLambda, "lambdalow", "warning");
	setStatusLED(maxLambda, "lambdahigh", "warning");
	setStatusLED(minIAC, "iaclow", "warning");
}

function setStatusLED(status, id, statustype = "status") {
	led = "green";

	if (statustype == "warning") led = "yellow";

	if (statustype == "fault") led = "red";

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

function increase(id) {
	var msg = formatSocketMessage("increase", id);
	sendSocketMessage(msg);
}

function decrease(id) {
	var msg = formatSocketMessage("decrease", id);
	sendSocketMessage(msg);
}

function updateAdjustmentValues(memsdata) {
	updateAdjustmentValue("idlespeed", memsdata.IdleSpeedOffset);
	updateAdjustmentValue("idlehot", memsdata.IdleHot);
	updateAdjustmentValue("ignitionadvance", memsdata.IgnitionAdvance);
	updateAdjustmentValue("fueltrim", memsdata.LongTermFuelTrim);
}

function updateAdjustmentValue(id, value) {
	$("td#" + id + ".adjustment").html(value.toString());
}

function setSerialPortSelection(ports) {
	$.each(ports, function (key, value) {
		console.log("serial port added " + key + " : " + value);
        //$("#serialports").append($("<option></option>").attr("value", value).text(value));
        $("#ports").append('<a class="dropdown-item" href="#" onclick="selectPort(this)">' + value + '</a>');
	});
}

function selectPort(item) {
    console.log('selected ' + item.text)
    setPort(item.text)
}

function setLogToFile(logsetting, logfolder) {
	if (logsetting != "stdout") {
		$("#logtofile").attr("checked", true);
	} else {
		$("#logtofile").attr("checked", false);
	}

	document.getElementById("logfolder").value = logfolder;
}

function setPort(port) {
	document.getElementById("port").value = port;
}

// Connect to the ECU
function connectECU() {
	var port = document.getElementById("port").value;
	var msg = formatSocketMessage("connect", port);
	sendSocketMessage(msg);

	// show connecting
	setConnectButtonStyle("<i class='fa fa-plug'>&nbsp</i>Connecting..", "btn-warning", connectECU);
	// disable all buttons
	$(":button").prop("disabled", true);
}

function readConfig() {
	var msg = formatSocketMessage("config", "read");
	sendSocketMessage(msg);
}

function resetECU() {
	var msg = formatSocketMessage("command", "resetecu");
	sendSocketMessage(msg);
}

function resetAdj() {
	var msg = formatSocketMessage("command", "resetadj");
	sendSocketMessage(msg);
}

function clearFaultCodes() {
	var msg = formatSocketMessage("command", "clearfaults");
	sendSocketMessage(msg);
}

// Pause the Data Loop
function pauseECUDataLoop() {
	var msg = formatSocketMessage("command", "pause");
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
	var msg = formatSocketMessage("command", "start");
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

function sendSocketMessage(msg) {
	console.log("sending socket message: " + msg);
	sock.send(msg);
}

function formatSocketMessage(a, d) {
	var msg = { action: a, data: d };
	return JSON.stringify(msg);
}
