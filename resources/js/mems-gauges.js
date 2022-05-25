
const Spark = {
    RPM : "rpmspark",
    MAP : "mapspark",
    Throttle : "throttlespark",
    IAC : "iacspark",
    Battery : "batteryspark",
    Coolant : "coolantspark",
    Air : "airspark",
    Lambda : "lambdaspark",
    Fuel : "fuelspark",
    LTFuel : "ltfuelspark",
    AirFuel : "airfuelspark",
    Ignition : "ignitionspark",
};

const Chart = {
    RPM : "rpmchart",
    Lambda : "lambdachart",
    LoopIndicator : "loopchart",
    Coolant : "coolantchart",
    AFR : "afrchart",
    IdleBase : "idlebasechart",
    IdleError : "idleerrorchart",
    MAP : "mapchart",
    CoilTime : "coiltimechart",
    CAS : "caschart",
    Battery : "batterychart",
};

export class MemsGauges {
    constructor() {
        this.initialiseGauges();
        this.initialiseSparklines();
    }

    initialiseGauges() {
        gaugeRPM.draw();
        gaugeMap.draw();
        gaugeThrottlePos.draw();
        gaugeIACPos.draw();
        gaugeBattery.value = 11;
        gaugeBattery.draw();
        gaugeCoolant.draw();
        gaugeAir.draw();
        gaugeLambda.draw();
        gaugeFuelTrim.draw();
        gaugeLTFuelTrim.draw();
        gaugeAirFuel.draw();
        gaugeIgnition.draw();

        // draw adaptive value gauges
        gaugeAdaptiveIdleSpeed.draw();
        gaugeAdaptiveIACPos.draw();
        gaugeAdaptiveIdleDecay.draw();
        gaugeAdaptiveSTFT.draw();
        gaugeAdaptiveLTFT.draw();
        gaugeAdaptiveIgnition.draw();
    }

    initialiseSparklines() {
        this.rpmSpark = createSpark(Spark.RPM);
        this.mapSpark = createSpark(Spark.MAP);
        this.throttleSpark = createSpark(Spark.Throttle);
        this.iacSpark = createSpark(Spark.IAC);
        this.batterySpark = createSpark(Spark.Battery);
        this.coolantSpark = createSpark(Spark.Coolant);
        this.airSpark = createSpark(Spark.Air);
        this.lambdaSpark = createSpark(Spark.Lambda);
        this.fuelSpark = createSpark(Spark.Fuel);
        this.ltfuelSpark = createSpark(Spark.LTFuel);
        this.airfuelSpark = createSpark(Spark.AirFuel);
        this.ignitionSpark = createSpark(Spark.Ignition);
    }

    initialiseGraphs() {
        this.rpmChart = createChart(Chart.RPM, "Engine (RPM)")
        this.lambdaChart = createChart(Chart.Lambda, "Lambda (mV)");
        this.loopChart = createChart(Chart.LoopIndicator, "O2 Loop (0 = Active)");
        this.afrChart = createChart(Chart.AFR, "Air : Fuel Ratio");
        this.coolantChart = createChart(Chart.Coolant, "Coolant (°C)");
        this.idleBaseChart = createChart(Chart.IdleBase, "Idle Base (Steps)");
        this.idleErrorChart = createChart(Chart.IdleError, "Idle Error (Steps)");
        this.mapChart = createChart(Chart.MAP, "MAP (kPa)");
        this.coilTimeChart = createChart(Chart.CoilTime, "Coil Time (ms)");
        this.casChart = createChart(Chart.CAS, "Crankshaft Position");
        this.batteryChart = createChart(Chart.Battery, "Battery (V)");
    }
}
