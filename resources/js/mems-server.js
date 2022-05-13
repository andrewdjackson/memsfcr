
export const SendRequest = async function (method, endpoint, body) {
    let init = {
        method: method,
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(body)
    }

    const response = await fetch(endpoint, init);
    if (!response.ok) throw new Error(`${response.status}`);

    return await response.json();
}

export const Endpoints = {
    // ROSCO ecu endpoints
    status: "/rosco",
    connect: "/rosco/connect",
    disconnect:  "/rosco/disconnect",
    heartbeat: "/rosco/heartbeat",
    dataframe:  "/rosco/dataframe",
    iacPosition: "/rosco/iac",
    diagnostics: "/rosco/diagnostics",
    reset: "/rosco/reset",
    resetECU: "/rosco/reset/ecu",
    resetFaults: "/rosco/reset/faults",
    resetAdjustments: "/rosco/reset/adjustments",
    adjust: "/rosco/adjust/",
    actuate: "/rosco/test/",

    // MemsFCR configuration endpoints
    config:  "/config",
    list_ports: "/config/ports",

    // MemsFCR scenario endpoints
    list: "/scenario",
    details: "/scenario/details",
    progress: "/scenario/progress",
    convert: "/scenario/convert",
    seek: "/scenario/seek",
}
