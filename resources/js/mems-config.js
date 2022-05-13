import {Endpoints, SendRequest} from "./mems-server.js";

export class MemsConfig {
    constructor(uri) {
        this._baseUri = uri;
        this._port = "";
        this._debug = false;
        this._version = "0.0.0";
        this._build = "2022-01-01";
        this._serverPort = "8081";
        this._frequency = 0;
    }

    load() {
        let endpoint = this._getEndpoint(Endpoints.config);

        console.info('loading configuration')

        return SendRequest('GET', endpoint)
            .then(data => this._updateConfig(data))
            .then(response => { return response; })
            .catch(err => this._restError(err))
    }

    save() {
        let body = {
            Port: this._port,
            Debug: this._debug,
            Frequency: this._frequency,
            ServerPort: this._serverPort,
            Version: this._version,
            Build: this._build,
        }

        let endpoint = this._getEndpoint(Endpoints.config);

        console.info(`saving configuration ${JSON.stringify(body)}`)

        return SendRequest('PUT',endpoint, body)
            .then(response => { return response; })
            .catch(err => this._restError(err))
    }

    ports() {
        let endpoint = this._getEndpoint(Endpoints.list_ports);

        console.info('loading available ports')

        return SendRequest('GET',endpoint)
            .then(data => this._updateAvailablePorts(data))
            .then(data => { return data; })
            .catch(err => this._restError(err))
    }

    _updateAvailablePorts(data) {
        console.info("available ports " + JSON.stringify(data));
        return data;
    }

    _updateConfig(data) {
        console.info("config loaded " + JSON.stringify(data));

        this._port = data.Port;
        this._debug = data.Debug;
        this._version = data.Version;
        this._build = data.Build;
        this._serverPort = data.ServerPort;

        return data;
    }

    get port() {
        return this._port;
    }

    set port(serialPort) {
        this._port = serialPort;
    }

    get frequency() {
        return this._frequency;
    }

    set frequency(ms) {
        this._frequency = ms;
    }

    get serverPort() {
        return this._serverPort
    }

    _getEndpoint(endpoint) {
        return this._baseUri + endpoint;
    }

    //
    // throw an error if the rest call failed
    //

    _restError(err) {
        throw new Error(`request failed (${err})`);
    }

}
