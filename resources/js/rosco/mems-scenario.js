import {Endpoints, SendRequest} from "./mems-server.js"

export class MemsScenario {
    constructor(uri) {
        this._baseUri = uri;
        this._name = "";
        this._count = 0;
        this._position = 0;
    }

    //
    // list available scenarios
    //

    list() {
        let endpoint = this._getEndpoint(Endpoints.list);

        console.info('getting available scenarios')

        return SendRequest('GET', endpoint)
            .then(data => this._availableScenarios(data))
            .then(response => { return response; })
            .catch(err => this._restError(err))
    }

    _availableScenarios(data) {
        console.info("available scenarios " + JSON.stringify(data));
        return data;
    }

    //
    // get details for the specified scenario
    //

    details(scenarioId) {
        let endpoint = this._getEndpoint(Endpoints.details);
        endpoint += `/${scenarioId}`;

        console.info(`playing scenario ${scenarioId} --> ${endpoint}`)

        return SendRequest('GET', endpoint)
            .then(data => this._scenarioDetails(data))
            .then(response => { return response; })
            .catch(err => this._restError(err))
    }

    _scenarioDetails(data) {
        console.info("scenario details " + JSON.stringify(data));
        return data;
    }

    //
    // get the information on the specified scenario
    //

    status(scenarioId) {
        let endpoint = this._getEndpoint(Endpoints.progress);
        endpoint += `/${scenarioId}`;

        console.info(`getting progress status on scenario ${scenarioId} --> ${endpoint}`)

        return SendRequest('GET', endpoint)
            .then(data => this._scenarioStatus(data))
            .then(response => { return response; })
            .catch(err => this._restError(err))
    }

    _scenarioStatus(data) {
        console.info("scenario status " + JSON.stringify(data));
        return data;
    }

    //
    // move the scenario to the specified position
    //

    seek(position) {
        let endpoint = this._getEndpoint(Endpoints.seek);

        // current position is returned but not used in the post
        let body = {
            CurrentPosition: 1,
            NewPosition: position,
        }

        console.info(`moving scenario to location ${position} --> ${endpoint}`)

        return SendRequest('POST', endpoint, body)
            .then(data => this._movedScenarioPosition(data))
            .then(response => { return response; })
            .catch(err => this._restError(err))
    }

    _movedScenarioPosition(data) {
        console.info("scenario moved location " + JSON.stringify(data));
        return data;
    }

    //
    // convert memsfcr .CVS file to a .FCR file
    //

    convert(scenarioId) {
        let endpoint = this._getEndpoint(Endpoints.convert);
        let body = { Source: scenarioId };

        console.info(`converting scenario ${scenarioId} to FCR file --> ${endpoint}`)
        return SendRequest('PUT', endpoint, body)
            .then(data => this._converted(data))
            .then(response => { return response; })
            .catch(err => this._restError(err))
    }

    _converted(data) {
        console.info("scenario converted " + JSON.stringify(data));
        return data;
    }

    //
    // set the scenario for replaying
    //

    replay(scenarioId) {
        this._name = scenarioId;
        this._count = this.details(scenarioId);
        this._position = 0;
    }

    //
    // transform the endpoint into a fqdn
    //

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
