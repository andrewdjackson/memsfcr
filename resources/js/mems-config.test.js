import {it, expect} from 'vitest';
import {MemsConfig} from "./mems-config.js"

var config = new MemsConfig("http://127.0.0.1:8081");

it('loads config', async () => {
    let response = await config.load();
    expect(parseInt(response.Frequency)).toBeGreaterThan(100);
});

it('loads available ports', async () => {
    let response = await config.ports();
    expect(response.ports.length).toBeGreaterThan(1);
});
