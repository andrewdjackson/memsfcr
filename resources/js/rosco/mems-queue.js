export const EventTopic = {
    Dataframe:   "dataframe",
    Heartbeat:  "heartbeat",
    Adjustment: "adjustment",
    Actuator: "actuator",
    Reset: "reset",
    IAC: "iac",
};

export class EventQueue {
    constructor() {
        this.handlers = [];
    }

    subscribe(event, handler) {
        this.handlers[event] = this.handlers[event] || [];
        this.handlers[event].push(handler);
    }

    publish(event, eventData) {
        const eventHandlers = this.handlers[event];

        if (eventHandlers) {
            for (var i = 0, l = eventHandlers.length; i < l; ++i) {
                eventHandlers[i].call({}, eventData);
            }
        }
    }
}
