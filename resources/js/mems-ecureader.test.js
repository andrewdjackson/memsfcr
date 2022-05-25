import {test, expect, beforeAll, afterAll} from 'vitest';
import {ECUCommand, ECUReader} from "./mems-ecureader.js";
import {Mems16Reader, Actuator, Adjuster} from "./mems-mems16reader.js"
import {EventQueue, EventTopic} from "./mems-queue.js";

var reader = new Mems16Reader("http://127.0.0.1:8081");
var eventQueue = new EventQueue();
var ecuReader = new ECUReader(reader, eventQueue);

beforeAll(async () => {
    let status = await reader.connect("/Users/andrew.jacksonglobalsign.com/ttyecu");
    expect(status.Connected).toEqual(true);

    // assign the event listeners
    eventQueue.subscribe(EventTopic.Dataframe, receivedDataframe);
    eventQueue.subscribe(EventTopic.Heartbeat, receivedHeartbeat);
    eventQueue.subscribe(EventTopic.Actuator, receivedActuator);
    eventQueue.subscribe(EventTopic.Adjustment, receivedAdjustment);
})

afterAll( async () => {
    await reader.disconnect();
})

/*
test('ecu dataframe loop', async () => {
    // start the dataframe loop, wait for 2 seconds before terminating
    // the dataframe loop should run i
    ecuReader.paused = false;
    // set the refresh interval to 1s. This controls the rate at which
    // commands can be sent to the ECU
    ecuReader.interval = 1000;

    ecuReader.sendAndReceiveLoop();
    await duration(4000);

    // force the dataframeLoop to exit
    ecuReader.reader.connected = false;
    ecuReader.reader.connected = true;
});

test('paused ecu dataframe loop', async () => {
    ecuReader.paused = true;
    ecuReader.reader.connected = true;

    ecuReader.sendAndReceiveLoop();
    await duration(2000);

    // force the dataframeLoop to exit
    ecuReader.reader.connected = false;
    ecuReader.reader.connected = true;
});
*/
test('mixed command loop', async () => {
    ecuReader.paused = false;
    ecuReader.reader.connected = true;

    ecuReader.sendAndReceiveLoop();

    let ecuCommand = new ECUCommand(ecuReader._counter++, EventTopic.Adjustment, reader.adjust(Adjuster.IdleDecay,1));
    ecuReader.send(ecuCommand);

    ecuCommand = new ECUCommand(ecuReader._counter++, EventTopic.Actuator, reader.actuate(Actuator.FuelPump, true));
    ecuReader.send(ecuCommand);

    await duration(4000);

    // force the dataframeLoop to exit
    ecuReader.reader.connected = false;
    ecuReader.reader.connected = true;
});

function receivedDataframe(data) {
    console.debug(`consumed ecu-dataframe`);
}

function receivedHeartbeat(data) {
    console.debug(`consumed ecu-heartbeat`);
}

function receivedActuator(data) {
    console.debug(`consumed ecu-actuator`);
}

function receivedAdjustment(data) {
    console.debug(`consumed ecu-adjustment`);
}

function apiError(e) {
    console.error("test api error " + e);
}

function duration(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}
