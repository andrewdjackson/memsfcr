var displayWidth = 200
var displayHeight = 200
var animationSpeed = 400
var gaugeFontFamily = '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji"'

var gaugeRPM = new RadialGauge({
    renderTo: 'gauge-rpm',
    title: 'Engine Speed',
    width: displayWidth,
    height: displayHeight,
    units: 'RPM',
    minValue: 0,
    maxValue: 7000,
    majorTicks: [
        '0',
        '1000',
        '2000',
        '3000',
        '4000',
        '5000',
        '6000',
        '7000'
    ],
    minorTicks: 10,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{
        from: 5000,
        to: 6000,
        color: 'rgba(78, 78, 76, 0.5)'
    }, {
        from: 6000,
        to: 7000,
        color: 'rgba(225, 7, 23, 0.75)'
    }],
    valueInt: 1,
    valueDec: 0,
    colorPlate: "#222",
    colorMajorTicks: "#f5f5f5",
    colorMinorTicks: "#ddd",
    colorTitle: "#ddd",
    colorUnits: "#ccc",
    colorNumbers: "#eee",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});

var gaugeMap = new RadialGauge({
    renderTo: 'gauge-map',
    title: 'Manifold Absolute Pressure',
    width: displayWidth,
    height: displayHeight,
    units: 'kPa',
    minValue: 0,
    maxValue: 100,
    majorTicks: [
        '0',
        '10',
        '20',
        '30',
        '40',
        '50',
        '60',
        '70',
        '80',
        '90',
        '100'
    ],
    minorTicks: 5,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{
        from: 0,
        to: 10,
        color: 'rgba(225, 7, 23, 0.75)'
    }, {
        from: 10,
        to: 60,
        color: 'rgba(10, 225, 6, 0.5)'
    }, {
        from: 80,
        to: 100,
        color: 'rgba(225, 7, 23, 0.75)'
    }],
    valueInt: 1,
    valueDec: 0,
    colorPlate: "#222",
    colorMajorTicks: "#f5f5f5",
    colorMinorTicks: "#ddd",
    colorTitle: "#ddd",
    colorUnits: "#ccc",
    colorNumbers: "#eee",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});

var gaugeThrottlePos = new RadialGauge({
    renderTo: 'gauge-throttlepos',
    title: 'Throttle Position',
    width: displayWidth,
    height: displayHeight,
    units: '%',
    minValue: 0,
    maxValue: 100,
    majorTicks: [
        '0',
        '10',
        '20',
        '30',
        '40',
        '50',
        '60',
        '70',
        '80',
        '90',
        '100'
    ],
    minorTicks: 5,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{
        from: 80,
        to: 100,
        color: 'rgba(225, 7, 23, 0.75)'
    }],
    valueInt: 1,
    valueDec: 0,
    colorPlate: "#222",
    colorMajorTicks: "#f5f5f5",
    colorMinorTicks: "#ddd",
    colorTitle: "#ddd",
    colorUnits: "#ccc",
    colorNumbers: "#eee",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});

var gaugeIACPos = new RadialGauge({
    renderTo: 'gauge-iacpos',
    title: 'Stepper Motor Position (IAC)',
    width: displayWidth,
    height: displayHeight,
    units: 'steps',
    minValue: 0,
    maxValue: 200,
    majorTicks: [
        '0',
        '20',
        '40',
        '60',
        '80',
        '100',
        '120',
        '140',
        '160',
        '180',
        '200'
    ],
    minorTicks: 5,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{
        from: 0,
        to: 30,
        color: 'rgba(225, 7, 23, 0.75)'
    }, {
        from: 170,
        to: 200,
        color: 'rgba(225, 7, 23, 0.75)'
    }],
    valueInt: 1,
    valueDec: 0,
    colorPlate: "#222",
    colorMajorTicks: "#f5f5f5",
    colorMinorTicks: "#ddd",
    colorTitle: "#ddd",
    colorUnits: "#ccc",
    colorNumbers: "#eee",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});

var gaugeBattery = new RadialGauge({
    renderTo: 'gauge-battery',
    title: 'Battery Voltage',
    width: displayWidth,
    height: displayHeight,
    units: 'V',
    minValue: 11,
    maxValue: 16,
    majorTicks: [
        '11.0',
        '12.0',
        '13.0',
        '14.0',
        '15.0',
        '16.0'
    ],
    minorTicks: 0.5,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{
        from: 12.5,
        to: 14.5,
        color: 'rgba(10, 225, 6, 0.5)'
    }, {
        from: 14.5,
        to: 16,
        color: 'rgba(225, 7, 23, 0.75)'
    }],
    valueInt: 0,
    valueDec: 1.0,
    colorPlate: "#222",
    colorMajorTicks: "#f5f5f5",
    colorMinorTicks: "#ddd",
    colorTitle: "#ddd",
    colorUnits: "#ccc",
    colorNumbers: "#eee",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});

var gaugeCoolant = new RadialGauge({
    renderTo: 'gauge-coolant',
    title: 'Coolant Temperature',
    width: displayWidth,
    height: displayHeight,
    units: '°C',
    minValue: 0,
    maxValue: 120,
    majorTicks: [
        '0',
        '10',
        '20',
        '30',
        '40',
        '50',
        '60',
        '70',
        '80',
        '90',
        '100',
        '110',
        '120'
    ],
    minorTicks: 5,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{
        from: 80,
        to: 95,
        color: 'rgba(10, 225, 6, 0.5)'
    }, {
        from: 95,
        to: 120,
        color: 'rgba(225, 7, 23, 0.75)'
    }],
    valueInt: 1,
    valueDec: 1.0,
    colorPlate: "#222",
    colorMajorTicks: "#f5f5f5",
    colorMinorTicks: "#ddd",
    colorTitle: "#ddd",
    colorUnits: "#ccc",
    colorNumbers: "#eee",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});

var gaugeAir = new RadialGauge({
    renderTo: 'gauge-air',
    title: 'Air Intake Temperature',
    width: displayWidth,
    height: displayHeight,
    units: '°C',
    minValue: -20,
    maxValue: 80,
    majorTicks: [
        '-20',
        '-10',
        '0',
        '10',
        '20',
        '30',
        '40',
        '50',
        '60',
        '70',
        '80'
    ],
    minorTicks: 5,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{
        from: 10,
        to: 50,
        color: 'rgba(10, 225, 6, 0.5)'
    }, {
        from: 50,
        to: 80,
        color: 'rgba(225, 7, 23, 0.75)'
    }],
    valueInt: 1,
    valueDec: 1.0,
    colorPlate: "#222",
    colorMajorTicks: "#f5f5f5",
    colorMinorTicks: "#ddd",
    colorTitle: "#ddd",
    colorUnits: "#ccc",
    colorNumbers: "#eee",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});

var gaugeLambda = new RadialGauge({
    renderTo: 'gauge-lambda',
    title: 'Lambda Voltage',
    width: displayWidth,
    height: displayHeight,
    units: 'mV',
    minValue: 0,
    maxValue: 900,
    majorTicks: [
        '0',
        '100',
        '200',
        '300',
        '400',
        '500',
        '600',
        '700',
        '800',
        '900'
    ],
    minorTicks: 5,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{
        from: 0,
        to: 50,
        color: 'rgba(225, 7, 23, 0.75)'
    }, {
        from: 50,
        to: 450,
        color: 'rgba(102,153,255,0.75)'
    }, {
        from: 450,
        to: 890,
        color: 'rgba(102,0,255,0.75)'
    }, {
        from: 850,
        to: 900,
        color: 'rgba(225, 7, 23, 0.75)'
    }],
    valueInt: 1,
    valueDec: 0,
    colorPlate: "#222",
    colorMajorTicks: "#f5f5f5",
    colorMinorTicks: "#ddd",
    colorTitle: "#ddd",
    colorUnits: "#ccc",
    colorNumbers: "#eee",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});

var gaugeFuelTrim = new RadialGauge({
    renderTo: 'gauge-fueltrim',
    title: 'STFT',
    width: displayWidth,
    height: displayHeight,
    units: '',
    minValue: -140,
    maxValue: 140,
    majorTicks: [
        '-140',
        '-120',
        '-100',
        '-80',
        '-60',
        '-40',
        '-20',
        '0',
        '20',
        '40',
        '60',
        '80',
        '100',
        '120',
        '140'
    ],
    minorTicks: 5,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{
        from: -128,
        to: 0,
        color: 'rgba(102,153,255,0.75)'
    }, {
        from: 0,
        to: 128,
        color: 'rgba(102,0,255,0.75)'
    }],
    valueInt: 1,
    valueDec: 0,
    colorPlate: "#222",
    colorMajorTicks: "#f5f5f5",
    colorMinorTicks: "#ddd",
    colorTitle: "#ddd",
    colorUnits: "#ccc",
    colorNumbers: "#eee",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});

var gaugeLTFuelTrim = new RadialGauge({
    renderTo: 'gauge-ltfueltrim',
    title: 'LTFT',
    width: displayWidth,
    height: displayHeight,
    units: '',
    minValue: -140,
    maxValue: 140,
    majorTicks: [
        '-140',
        '-120',
        '-100',
        '-80',
        '-60',
        '-40',
        '-20',
        '0',
        '20',
        '40',
        '60',
        '80',
        '100',
        '120',
        '140'
    ],
    minorTicks: 5,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{
        from: -128,
        to: 0,
        color: 'rgba(102,153,255,0.75)'
    }, {
        from: 0,
        to: 128,
        color: 'rgba(102,0,255,0.75)'
    }],
    valueInt: 1,
    valueDec: 0,
    colorPlate: "#222",
    colorMajorTicks: "#f5f5f5",
    colorMinorTicks: "#ddd",
    colorTitle: "#ddd",
    colorUnits: "#ccc",
    colorNumbers: "#eee",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});

var gaugeAirFuel = new RadialGauge({
    renderTo: 'gauge-airfuel',
    title: 'Air : Fuel Ratio',
    width: displayWidth,
    height: displayHeight,
    units: '',
    minValue: 10,
    maxValue: 20,
    majorTicks: [
        '10.0',
        '12.0',
        '14.0',
        '16.0',
        '18.0',
        '20.0'
    ],
    minorTicks: 5,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{
        from: 10,
        to: 11,
        color: 'rgba(225, 7, 23, 0.75)'
    }, {
        from: 11,
        to: 12,
        color: 'rgba(225, 185, 6, 0.75)'
    }, {
        from: 13.5,
        to: 16.5,
        color: 'rgba(10, 225, 6, 0.5)'
    }, {
        from: 18,
        to: 19,
        color: 'rgba(225, 185, 6, 0.75)'
    }, {
        from: 19,
        to: 20,
        color: 'rgba(225, 7, 23, 0.75)'
    }],
    valueInt: 1,
    valueDec: 1.0,
    colorPlate: "#222",
    colorMajorTicks: "#f5f5f5",
    colorMinorTicks: "#ddd",
    colorTitle: "#ddd",
    colorUnits: "#ccc",
    colorNumbers: "#eee",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});

var gaugeIgnition = new RadialGauge({
    renderTo: 'gauge-ignition',
    title: 'Ignition Advance',
    width: displayWidth,
    height: displayHeight,
    units: '°',
    minValue: 0,
    maxValue: 40,
    majorTicks: [
        '0',
        '5',
        '10',
        '15',
        '20',
        '25',
        '30',
        '35',
        '40'
    ],
    minorTicks: 5,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{
        from: 0,
        to: 20,
        color: 'rgba(10, 225, 6, 0.5)'
    }, {
        from: 20,
        to: 40,
        color: 'rgba(225, 7, 23, 0.75)'
    }],
    valueInt: 1,
    valueDec: 0,
    colorPlate: "#222",
    colorMajorTicks: "#f5f5f5",
    colorMinorTicks: "#ddd",
    colorTitle: "#ddd",
    colorUnits: "#ccc",
    colorNumbers: "#eee",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});

var gaugeAdaptiveIgnition = new RadialGauge({
    renderTo: 'ignitionadvance',
    title: 'Ignition Advance',
    width: displayWidth,
    height: displayHeight,
    units: '',
    minValue: -128,
    maxValue: 128,
    majorTicks: [
        '-130','-110','-90','-70','-50','-30','-10',
        '0',
        '10','30','50','70','90','110','130',
    ],
    minorTicks: 10,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{}],
    valueInt: 1,
    valueDec: 0,
    colorPlate: "#f5f5f5",
    colorMajorTicks: "#111",
    colorMinorTicks: "#999",
    colorTitle: "#111",
    colorUnits: "#111",
    colorNumbers: "#999",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});

var gaugeAdaptiveIdleSpeed = new RadialGauge({
    renderTo: 'idlespeed',
    title: 'Idle Speed Offset',
    width: displayWidth,
    height: displayHeight,
    units: '',
    minValue: -800,
    maxValue: 800,
    majorTicks: [
        '-800','-600','-400','-200',
        '0','200','400','600','800'
    ],
    minorTicks: 10,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{}],
    valueInt: 1,
    valueDec: 0,
    colorPlate: "#f5f5f5",
    colorMajorTicks: "#111",
    colorMinorTicks: "#999",
    colorTitle: "#111",
    colorUnits: "#111",
    colorNumbers: "#999",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});

var gaugeAdaptiveIdleDecay = new RadialGauge({
    renderTo: 'idledecay',
    title: 'Idle Decay',
    width: displayWidth,
    height: displayHeight,
    units: '',
    minValue: -128,
    maxValue: 128,
    majorTicks: [
        '-130','-110','-90','-70','-50','-30','-10',
        '0',
        '10','30','50','70','90','110','130',
    ],
    minorTicks: 10,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{}],
    valueInt: 1,
    valueDec: 0,
    colorPlate: "#f5f5f5",
    colorMajorTicks: "#111",
    colorMinorTicks: "#999",
    colorTitle: "#111",
    colorUnits: "#111",
    colorNumbers: "#999",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});

var gaugeAdaptiveLTFT = new RadialGauge({
    renderTo: 'ltft',
    title: 'LTFT',
    width: displayWidth,
    height: displayHeight,
    units: '',
    minValue: -128,
    maxValue: 128,
    majorTicks: [
        '-130','-110','-90','-70','-50','-30','-10',
        '0',
        '10','30','50','70','90','110','130',
    ],
    minorTicks: 10,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{}],
    valueInt: 1,
    valueDec: 0,
    colorPlate: "#f5f5f5",
    colorMajorTicks: "#111",
    colorMinorTicks: "#999",
    colorTitle: "#111",
    colorUnits: "#111",
    colorNumbers: "#999",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});

var gaugeAdaptiveSTFT = new RadialGauge({
    renderTo: 'stft',
    title: 'STFT',
    width: displayWidth,
    height: displayHeight,
    units: '',
    minValue: -128,
    maxValue: 128,
    majorTicks: [
        '-130','-110','-90','-70','-50','-30','-10',
        '0',
        '10','30','50','70','90','110','130',
    ],
    minorTicks: 10,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{}],
    valueInt: 1,
    valueDec: 0,
    colorPlate: "#f5f5f5",
    colorMajorTicks: "#111",
    colorMinorTicks: "#999",
    colorTitle: "#111",
    colorUnits: "#111",
    colorNumbers: "#999",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});

var gaugeAdaptiveIACPos = new RadialGauge({
    renderTo: 'iac',
    title: 'IAC Steps',
    width: displayWidth,
    height: displayHeight,
    units: '',
    minValue: 0,
    maxValue: 260,
    majorTicks: [
        '0','20','40','60','80','100','120','140','160','180','200','220','240','260'
    ],
    minorTicks: 10,
    ticksAngle: 250,
    startAngle: 55,
    strokeTicks: true,
    highlights: [{}],
    valueInt: 1,
    valueDec: 0,
    colorPlate: "#f5f5f5",
    colorMajorTicks: "#111",
    colorMinorTicks: "#999",
    colorTitle: "#111",
    colorUnits: "#111",
    colorNumbers: "#999",
    valueBox: true,
    valueBoxWidth: 30,
    colorValueText: "#000",
    colorValueBoxRect: "#ccc",
    colorValueBoxRectEnd: "#ccc",
    colorValueBoxBackground: "#ccc",
    colorValueBoxShadow: false,
    colorValueTextShadow: false,
    colorNeedleShadowUp: false,
    colorNeedleShadowDown: "333",
    colorNeedle: "rgba(240, 128, 128, 1)",
    colorNeedleEnd: "rgba(255, 160, 122, .9)",
    colorNeedleCircleOuter: "rgba(200, 200, 200, 1)",
    colorNeedleCircleOuterEnd: "rgba(200, 200, 200, 1)",
    borderShadowWidth: 0,
    borders: false,
    borderInnerWidth: 0,
    borderMiddleWidth: 0,
    borderOuterWidth: 5,
    colorBorderOuter: "#fafafa",
    colorBorderOuterEnd: "#cdcdcd",
    needleType: "arrow",
    needleWidth: 2,
    needleCircleSize: 7,
    needleCircleOuter: true,
    needleCircleInner: false,
    animationDuration: animationSpeed,
    animationRule: "bounce",
    fontNumbers: gaugeFontFamily,
    fontTitle: gaugeFontFamily,
    fontUnits: gaugeFontFamily,
    fontValue: gaugeFontFamily,
    fontValueStyle: 'normal',
    fontNumbersSize: 20,
    fontNumbersStyle: 'normal',
    fontNumbersWeight: 'normal',
    fontTitleSize: 14,
    fontUnitsSize: 22,
    fontValueSize: 30,
    animatedValue: true
});