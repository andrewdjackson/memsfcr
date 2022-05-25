import {Mems16Reader} from "./mems-mems16reader.js";
import {ECUReader} from "./mems-ecureader.js";
import {EventQueue} from "./mems-queue.js";
import {Endpoints} from "./mems-server.js";
import {MemsGauges} from "./mems-gauges.js";

var handler;

// this function gets called as soon as the page load has completed
window.onload = function () {
    //uri = window.location.href.split("/").slice(0, 3).join("/");
    let uri = window.location.origin;

    handler = new MemsHandler(uri);
    handler.initialise();
};

const RunState = {
    Disconnected:   "disconnected",
    Connected:  "connected",
    Paused: "paused",
};

const RunMode = {
    ECU: "ecu",
    Playback: "playback",
}

class MemsHandler {
    constructor(uri) {
        this.uri = uri;
        this.eventQueue = new EventQueue();
        this.memsReader = new Mems16Reader(this.uri);
        this.ecuReader = new ECUReader(this.memsReader, this.eventQueue);
    }

    initialise() {
        this.initialiseServerEvents();
        this.initialiseGauges();
        this.initialiseClickHandlers();
    }

    initialiseServerEvents() {
        // connect to the server to establish a heartbeat link
        // if the user closes the browser, the server will detect no response
        // and terminate the application after a few seconds
        let server_event = new EventSource(this.uri + Endpoints.serverHeartbeat);

        server_event.onOpen = function () {
            console.debug("server-event connected");
        }

        server_event.onClose = function () {
            console.debug("server-event close");
        }

        // server-event message handler
        server_event.onMessage = function (e) {
            console.debug("server-event message " + e.data);
        }

        // listen for heartbeat events
        server_event.addEventListener('heartbeat', this.heartbeatHandler, false);

        // listen for ecu connection state changes
        server_event.addEventListener('status', this.statusHandler, false);
    }

    heartbeatHandler(e) {
        console.debug('server heartbeat')
    }

    statusHandler(e) {
        console.debug('server status change')
    }

    initialiseGauges() {
        this.gauges = new MemsGauges();
    }

    initialiseClickHandlers() {
        document.getElementById('connectECUbtn').addEventListener("click", this._connectECU);
    }

    _connectECU() {
        console.info("connect ecu clicked");
    }
}
