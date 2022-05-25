import {Endpoints, SendRequest} from "./mems-server.js"

export const Actuator = {
    FuelPump:  "fuelpump",
    PTC:  "ptc",
    AirCon: "aircon",
    PurgeValve: "purgevalve",
    BoostValve: "boostvalve",
    Fan: "fan",
    Fan1: "fan/1",
    Fan2: "fan/2",
    Injectors: "injectors",
    Coil: "coil",
};

export const Adjuster = {
    STFT: "stft",
    LTFT: "ltft",
    IdleDecay: "idledecay",
    IdleSpeed: "idlespeed",
    IgnitionAdvance: "ignitionadvance",
    IAC: "iac",
};

export class Mems16Reader {
    constructor(uri) {
        this.connected = false;
        this.ecuId = "";
        this.ecuSerial = "";
        this.iacPosition = 0;
        this._baseUri = uri;
        this._resetStatus();
    }

    _getEndpoint(endpoint) {
        return this._baseUri + endpoint;
    }

    async disconnect() {
        let body = "";
        let endpoint = this._getEndpoint(Endpoints.disconnect);

        console.info(`disconnecting from the ecu`)

        return await SendRequest('POST', endpoint, body)
            .then(data => this._disconnected(data))
            .catch(err => this._restError(err))
    };

    _disconnected(data) {
        console.info(`disconnected from the ecu`);
        this._resetStatus();
        return data;
    }

    //
    // connect to the ecu
    //

    async connect(port) {
        let body = {Port: port}
        let endpoint = this._getEndpoint(Endpoints.connect);

        console.info('connecting to ecu')

        return await SendRequest('POST', endpoint, body)
            .then(data => this._connectedToECU(data))
            .then(response => { return response; })
            .catch(err => this._restError(err))
    }


    _connectedToECU(data) {
        console.info("connected to ecu (" + JSON.stringify(data) + ")");

        this.connected = data.Connected;
        this.ecuId = data.ECUID;
        this.ecuSerial = data.ECUSerial;
        this.iacPosition = data.IACPosition;

        return data;
    }

    //
    // send a keep-alive heartbeat to the ecu
    //

    async heartbeat() {
        let endpoint = this._getEndpoint(Endpoints.heartbeat);

        console.info(`sending heartbeat -> ${endpoint}`)

        return await SendRequest('POST', endpoint)
            .then(response => { return response; })
            .catch(err => this._restError(err))
    }

    //
    // activate / deactivate actuator
    //

   async  actuate(actuator, activate) {
        let body = {Activate: activate}
        let endpoint = this._getEndpoint(Endpoints.actuate + actuator);

        console.info(`actuator ${actuator} activating ${activate} -> ${endpoint}`)

        return await SendRequest('POST', endpoint, body)
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

    async adjust(adjuster, steps) {
        let body = {Steps: steps}
        let endpoint = this._getEndpoint(Endpoints.adjust + adjuster);

        console.info(`adjusting ${adjuster} by ${steps} steps -> ${endpoint}`)

        return await SendRequest('POST', endpoint, body)
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

    async status() {
        let endpoint = this._getEndpoint(Endpoints.status);

        console.info('getting ecu connection status')

        return await SendRequest('GET', endpoint)
            .then(data => this._updateStatus(data))
            .then(response => { return response; })
            .catch(err => this._restError(err))
    }

    _updateStatus(state) {
        console.info("status updated " + JSON.stringify(state));

        this.ecuId = state.ECUID;
        this.ecuSerial = state.ECUSerial;
        this.connected = state.Connected;
        this.iacPosition = state.IACPosition;

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

    async dataframes() {
        let endpoint = this._getEndpoint(Endpoints.dataframe);

        console.info('getting dataframes from ecu')

        return await SendRequest('GET', endpoint)
            .then(data => this._receivedDataframes(data))
            .then(response => { return response; })
            .catch(err => this._restError(err))
    }

    _receivedDataframes(data) {
        console.info("dataframes received ");// + JSON.stringify(data));
        return data;
    }

    //
    // throw an error if the rest call failed
    //

    _restError(err) {
        console.error(`${err.message} ${JSON.stringify(err.response)})`);
        return err;
    }

    //
    // asynchronously "sleep" for a period of time
    //

    sleep(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
}
