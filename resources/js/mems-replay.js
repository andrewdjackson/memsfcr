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

    request.onload = function () {
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

// calls the resp api to get the list of available scenarios
function updateScenarios() {
    uri = window.location.href.split("/").slice(0, 3).join("/");

    // Create a request variable and assign a new XMLHttpRequest object to it.
    var request = new XMLHttpRequest()

    // Open a new connection, using the GET request on the URL endpoint
    request.open('GET', uri + '/scenario', true)

    request.onload = function () {
        // Begin accessing JSON data here
        var data = JSON.parse(this.response)

        var replay = $('#replayScenarios');
        $.each(data, function (val, text) {
            var i = $('<button class="dropdown-item replay" type="button"></button>').val(text).html(text)
            replay.append(i);
        });
    }

    $("#replayScenarios").off().click(replaySelectedScenario);
    // Send request
    request.send()
}
