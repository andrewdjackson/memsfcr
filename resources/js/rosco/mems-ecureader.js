import {MemsAPIError} from "./mems-server.js";
import {EventTopic} from "./mems-queue.js";

export class ECUCommand {
    constructor(id, topic, command) {
        this.id = id;
        this.topic = topic;
        this.command = command;
    }
}

export class ECUReader {
    constructor(reader, eventQueue) {
        this.reader = reader;
        this._listening = true;
        this._paused = false;
        this._counter = 0;
        this._eventQueue = eventQueue;
        this._commandQueue = [];
        this._refreshInterval = 500;
    };

    //
    // pause getting dataframes from the ecu
    // heatbeats are sent whilst the loop is paused
    //

    get paused() {
        return this._paused;
    }

    set paused(pause) {
        this._paused = pause;
    };

    //
    // gets / sets the rate at which commands can be sent to the ecu
    //

    get interval() {
        return this._refreshInterval;
    }

    set interval(interval) {
        this._refreshInterval = interval;
    }

    //
    // start / stop the adding of commands to the queue
    //

    start() {
        this._listening = true;
    }

    stop() {
        this._listening = false;
    }

    //
    // runs the loop that sends a request for dataframes when not paused
    // when paused a keep-alive heartbeat is sent to the ecu to keep the connection open
    //

    async sendAndReceiveLoop() {
        while(this.reader.connected) {
            let ecuCommand;
            let id = this._counter++;

            if (!this.paused) {
                console.log(`queuing dataframe request`)
                // queue a request for the dataframes
                ecuCommand = new ECUCommand(id, EventTopic.Dataframe, this.reader.dataframes() )
            } else {
                console.log(`queuing heartbeat`)
                // queue a heartbeat command
                ecuCommand = new ECUCommand(id, EventTopic.Heartbeat, this.reader.heartbeat() )
            }

            // add command to the queue
            this.send(ecuCommand);

            // service the queue
            await this._sendToECU();
        }
    }

    //
    // add a command to the queue for sending to the ECU
    // returns true / false if the command was successfully added to the queue
    //

    send(ecuCommand) {
        if (this._listening) {
            // add command to the queue
            this._commandQueue.push(ecuCommand);
        }

        return this._listening;
    }

    //
    // executes the ecu command, rate if send is controlled here by waiting for
    // the timer to complete and the response from the ecu
    //

    async _sendToECU() {
        let ecuCommand = this._commandQueue.shift();
        console.info(`${Date.now().toString()} : sendToECU Executing ${ecuCommand.id}.${ecuCommand.topic}`);

        return Promise.allSettled([
            // execute the command and send the response to the receiver function
            // report any errors
            ecuCommand.command
                .then(response => {
                    console.debug(`${Date.now().toString()} : sendToECU Response ${ecuCommand.id}.${ecuCommand.topic}`);
                    this._receivedFromECU(ecuCommand.topic, response);
                })
                .catch(err => {
                    console.error(`${Date.now().toString()} : sendToECU Response ${ecuCommand.id}.${ecuCommand.topic} -> ${err}`);
                }),

            // wait for timer to expire, this essentially controls the
            // rate at which the commands are sent to the ecu
            this._sleep(this._refreshInterval).then(result => {
                console.debug(`${Date.now().toString()} : sendToECU Timer expired ${ecuCommand.id}.${ecuCommand.topic}`);
            }),
        ]).then(response => {
            console.debug(`${Date.now().toString()} : sendToECU Completed ${ecuCommand.id}.${ecuCommand.topic}`);
        });
    }

    //
    // publish the response received from the ecu on the event queue
    // if the response is not an error
    //

    _receivedFromECU(topic, response) {
        if (response instanceof MemsAPIError) {
            throw new MemsAPIError(response.message, response.response);
        } else {
            this._eventQueue.publish(topic, response);
        }
    }

    //
    // asynchronously "sleep" for a period of time
    //

    _sleep(ms) {
        console.debug(`${Date.now().toString()} : _sleep Timer set for ${ms}ms`);
        return new Promise(resolve => setTimeout(resolve, ms));
    }
}
