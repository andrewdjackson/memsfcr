import {it, expect} from 'vitest';
import {MemsScenario} from "./mems-scenario.js";
import {MemsReader, Actuator, Adjuster} from "./mems-reader.js";

const testScenario = 'vacuum-fault.fcr';
const testScenarioCSV = 'fullrun.csv';
const serverUrl = "http://127.0.0.1:8081";

var scenario = new MemsScenario(serverUrl);
var reader = new MemsReader(serverUrl);

it('lists scenario', async() => {
    let response = await scenario.list();
    expect(response.length).toBeGreaterThan(1);
});

it('gets scenario details', async() => {
    let response = await scenario.details(testScenario);
    expect(response).toBeDefined();
});

it('connect and seek scenario', async() => {
    let response = await reader.connect(testScenario);
    expect(response.Connected).toEqual(true);

    response = await scenario.seek(200);
    expect(response.Position).toBe(200);

    await reader.disconnect();
});

it('convert scenario', async() => {
    let response = await scenario.convert(testScenarioCSV);
    expect(response).toBeDefined();
})
