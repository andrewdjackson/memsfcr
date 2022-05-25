import {it, expect} from 'vitest';
import {Mems16Reader, Actuator, Adjuster} from "./mems-mems16reader.js"

var reader = new Mems16Reader("http://127.0.0.1:8081");

it('connect to mems reader', async () => {
    let status = await reader.connect("/Users/andrew.jacksonglobalsign.com/ttyecu");
    expect(status.Connected).toEqual(true);

    let activated = await reader.actuate(Actuator.FuelPump, true);
    expect(activated.activate).toBeTruthy();

    let adjusted = await reader.adjust(Adjuster.IdleDecay, -1);
    expect(adjusted.value).toBe(34);

    let dataframes = await reader.dataframes();
    expect(dataframes.Dataframe7d).toHaveLength(66);

    status = await reader.disconnect();
    expect(status.Connected).toEqual(false);
});
