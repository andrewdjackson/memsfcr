import {Endpoints, SendRequest} from "./mems-server.js"

export const Actuator = Object.freeze({
    FuelPump:   Symbol("fuelpump"),
    PTC:  Symbol("ptc"),
    AirCon: Symbol("aircon"),
    PurgeValve: Symbol("purgevalve"),
    BoostValve: Symbol("boostvalve"),
    Fan: Symbol("fan"),
    Fan1: Symbol("fan/1"),
    Fan2: Symbol("fan/2"),
    Injectors: Symbol("injectors"),
    Coil: Symbol("coil"),
});

export const Adjuster = Object.freeze({
    STFT: Symbol("stft"),
    LTFT: Symbol("ltft"),
    IdleDecay: Symbol("idledecay"),
    IdleSpeed: Symbol("idlespeed"),
    IgnitionAdvance: Symbol("ignitionadvance"),
    IAC: Symbol("iac"),
})

export class MemsReader {
    constructor(uri) {
        this._refreshInterval = 500;
        this._baseUri = uri;
        this._resetStatus();
    }

    _getEndpoint(endpoint) {
        return this._baseUri + endpoint;
    }

    disconnect() {
        let body = "";
        let endpoint = this._getEndpoint(Endpoints.disconnect);

        console.info(`disconnecting from the ecu`)

        return SendRequest('POST', endpoint, body)
            .then(data => this._disconnected(data))
            .catch(err => this._restError(err))
    };

    _disconnected(data) {
        console.info(`disconnected from the ecu`);
        this._resetStatus();
        return data;
    }

    pause() {
    };

    //
    // connect to the ecu
    //

    connect(port) {
        let body = {Port: port}
        let endpoint = this._getEndpoint(Endpoints.connect);

        console.info('connecting to ecu')

        return SendRequest('POST', endpoint, body)
            .then(data => this._connectedToECU(data))
            .then(response => { return response; })
            .catch(err => this._restError(err))
    }


    _connectedToECU(data) {
        console.info("connected to ecu (" + JSON.stringify(data) + ")");

        this._connected = data.Connected;
        this._ecuId = data.ECUID;
        this._ecuSerial = data.ECUSerial;
        this._iacPosition = data.IACPosition;

        return data;
    }

    //
    // activate / deactivate actuator
    //

    actuate(actuator, activate) {
        let body = {Activate: activate}
        let endpoint = this._getEndpoint(Endpoints.actuate + actuator.description);

        console.info(`actuator ${actuator.description} activating ${activate} -> ${endpoint}`)

        return SendRequest('POST', endpoint, body)
            .then(data => this._actuated(data))
            .then(response => { return response; })
            .catch(err => this._restError(err))
    }

    _actuated(data) {
        console.info("actuator response " + JSON.stringify(data));
        return data;
    }


    //
    // increment / decrement adjustment
    //

    adjust(adjuster, steps) {
        let body = {Steps: steps}
        let endpoint = this._getEndpoint(Endpoints.adjust + adjuster.description);

        console.info(`adjusting ${adjuster.description} by ${steps} steps -> ${endpoint}`)

        return SendRequest('POST', endpoint, body)
            .then(data => this._adjusted(data))
            .then(response => { return response; })
            .catch(err => this._restError(err))
    }

    _adjusted(data) {
        console.info("adjusted response " + JSON.stringify(data))
        return data;
    }

    //
    // get the connection status
    //

    status() {
        let endpoint = this._getEndpoint(Endpoints.status);

        console.info('getting ecu connection status')

        return SendRequest('GET', endpoint)
            .then(data => this._updateStatus(data))
            .then(response => { return response; })
            .catch(err => this._restError(err))
    }

    _updateStatus(state) {
        console.info("status updated " + JSON.stringify(state));

        this._ecuId = state.ECUID;
        this._ecuSerial = state.ECUSerial;
        this._connected = state.Connected;
        this._iacPosition = state.IACPosition;

        return state;
    }

    _resetStatus() {
        let state = {
            Connected: false,
            ECUID: "",
            ECUSerial: "",
            IACPosition: 0
        }

        this._updateStatus(state);
    }

    //
    // get the dataframes
    //

    dataframes() {
        let endpoint = this._getEndpoint(Endpoints.dataframe);

        console.info('getting dataframes from ecu')

        return SendRequest('GET', endpoint)
            .then(data => this._receivedDataframes(data))
            .then(response => { return response; })
            .catch(err => this._restError(err))
    }

    _receivedDataframes(data) {
        console.info("dataframes received " + JSON.stringify(data))
        return data;
    }

    //
    // throw an error if the rest call failed
    //

    _restError(err) {
        throw new Error(`request failed (${err})`);
    }

    //
    // asynchronously "sleep" for a period of time
    //

    sleep(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
}
